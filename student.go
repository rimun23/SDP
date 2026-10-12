package main

import (
	"SDP/subsystemA"
	"SDP/subsystemB"
	"SDP/subsystemC"
	"fmt"
	"log"
)

type StudentRegistrationFacade struct {
	prereqChecker *subsystemA.PrerequisiteChecker
	holdChecker   *subsystemA.FinancialHoldChecker
	intakeRepo    *subsystemA.IntakeRepository
	seatAllocator *subsystemB.SeatAllocator
	seatRepo      *subsystemB.SeatRepository
	billingSystem *subsystemC.BillingSystem
	transcriptSvc *subsystemC.TranscriptService
	confirmRepo   *subsystemC.ConfirmationRepository
	journal       *JournalStore
}

func NewStudentRegistrationFacade(
	storageDir string,
	pa *subsystemA.PrerequisiteChecker,
	ha *subsystemA.FinancialHoldChecker,
	ra *subsystemA.IntakeRepository,
	sb *subsystemB.SeatAllocator,
	rb *subsystemB.SeatRepository,
	bc *subsystemC.BillingSystem,
	tc *subsystemC.TranscriptService,
	rc *subsystemC.ConfirmationRepository) (*StudentRegistrationFacade, error) {
	j, err := NewJournalStore(storageDir)
	if err != nil {
		return nil, err
	}
	return &StudentRegistrationFacade{
		prereqChecker: pa,
		holdChecker:   ha,
		intakeRepo:    ra,
		seatAllocator: sb,
		seatRepo:      rb,
		billingSystem: bc,
		transcriptSvc: tc,
		confirmRepo:   rc,
		journal:       j,
	}, nil
}
func (f *StudentRegistrationFacade) RegisterCourse(regID, studentID, courseID string, credits int) error {
	log.Printf("[FACADE] Starting registration chain %s for student %s in course %s", regID, studentID, courseID)
	if err := f.prereqChecker.ValidateEligibility(studentID, courseID); err != nil {
		f.journal.RecordStep(regID, studentID, courseID, StepFailed, err.Error())
		return fmt.Errorf("step 0 eligibility failed: %w", err)
	}
	if err := f.holdChecker.CheckHold(studentID); err != nil {
		f.journal.RecordStep(regID, studentID, courseID, StepFailed, err.Error())
		return fmt.Errorf("step 0 financial hold: %w", err)
	}
	intake := subsystemA.NewIntakeRecord(regID, studentID, courseID)
	if err := intake.TransitionTo(subsystemA.StatusReserved); err != nil {
		return err
	}
	if err := f.intakeRepo.Save(intake); err != nil {
		return err
	}
	f.journal.RecordStep(regID, studentID, courseID, StepACommitted, "A committed, waiting B")
	log.Printf("[FACADE] Step A committed: slot reserved for %s", studentID)
	seatNum, err := f.seatAllocator.Allocate(courseID)
	if err != nil {
		log.Printf("[FACADE] ALERT: Step B failed (%v)! Triggering 1-step rollback of A...", err)
		f.compensateA(regID, studentID, courseID, err.Error())
		return fmt.Errorf("step B failed, chain compensated: %w", err)
	}
	seat := subsystemB.NewSeatAllocation(regID, courseID, seatNum)
	if err := seat.TransitionTo(subsystemB.StatusAllocated); err != nil {
		f.compensateA(regID, studentID, courseID, err.Error())
		return err
	}
	if err := f.seatRepo.Save(seat); err != nil {
		f.compensateA(regID, studentID, courseID, err.Error())
		return err
	}
	f.journal.RecordStep(regID, studentID, courseID, StepBCommitted, "B committed, waiting C")
	log.Printf("[FACADE] Step B committed: seat #%d allocated for %s", seatNum, studentID)
	amount, invoice, err := f.billingSystem.ProcessTuitionCharge(studentID, credits)
	if err != nil {
		log.Printf("[FACADE] ALERT: Step C failed (%v)! Triggering 2-step rollback: B -> A...", err)
		f.compensateBThenA(regID, studentID, courseID, err.Error())
		return fmt.Errorf("step C failed, chain fully compensated (B then A): %w", err)
	}
	if err := f.transcriptSvc.RecordTermEntry(studentID, courseID); err != nil {
		log.Printf("[FACADE] ALERT: Step C transcript failed (%v)! Triggering 2-step rollback: B -> A...", err)
		f.compensateBThenA(regID, studentID, courseID, err.Error())
		return fmt.Errorf("step C failed, chain compensated: %w", err)
	}
	confirm := subsystemC.NewEnrollmentConfirmation(regID, amount, invoice)
	if err := confirm.TransitionTo(subsystemC.StatusConfirmed); err != nil {
		f.compensateBThenA(regID, studentID, courseID, err.Error())
		return err
	}
	if err := f.confirmRepo.Save(confirm); err != nil {
		f.compensateBThenA(regID, studentID, courseID, err.Error())
		return err
	}
	f.journal.RecordStep(regID, studentID, courseID, StepCompleted, "All 3 steps successful")
	log.Printf("[FACADE] Step C committed! Chain successfully finished for %s.", studentID)
	return nil
}
func (f *StudentRegistrationFacade) compensateA(regID, studentID, courseID, reason string) {
	log.Printf("[COMPENSATION] Rolling back Subsystem A for %s...", regID)
	intake, err := f.intakeRepo.FindByID(regID)
	if err == nil {
		_ = intake.TransitionTo(subsystemA.StatusCancelled)
		_ = f.intakeRepo.Save(intake)
	}
	_ = f.journal.RecordStep(regID, studentID, courseID, StepCompensated, fmt.Sprintf("B failed: %s; A rolled back", reason))
}
func (f *StudentRegistrationFacade) compensateBThenA(regID, studentID, courseID, reason string) {
	log.Printf("[COMPENSATION] Initiating 2-step reverse rollback: B then A...")
	seat, err := f.seatRepo.FindByID(regID)
	if err == nil {
		_ = seat.TransitionTo(subsystemB.StatusReleased)
		_ = f.seatRepo.Save(seat)
		f.seatAllocator.Deallocate(courseID)
		log.Printf("[COMPENSATION] Step 1/2 complete: Subsystem B seat released.")
	}
	intake, err := f.intakeRepo.FindByID(regID)
	if err == nil {
		_ = intake.TransitionTo(subsystemA.StatusCancelled)
		_ = f.intakeRepo.Save(intake)
		log.Printf("[COMPENSATION] Step 2/2 complete: Subsystem A intake reservation cancelled.")
	}
	_ = f.journal.RecordStep(regID, studentID, courseID, StepCompensated, fmt.Sprintf("C failed: %s; B then A rolled back", reason))
}
