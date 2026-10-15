package main

import (
	"SDP/subsystemA"
	"SDP/subsystemB"
	"SDP/subsystemC"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type JournalStore struct{ path string }

func (j *JournalStore) Log(regID, state, reason string) {
	data, _ := os.ReadFile(j.path)
	logs := make(map[string]map[string]string)
	_ = json.Unmarshal(data, &logs)
	logs[regID] = map[string]string{"state": state, "reason": reason}
	b, _ := json.MarshalIndent(logs, "", "  ")
	_ = os.WriteFile(j.path, b, 0644)
}

func (j *JournalStore) GetAll() map[string]map[string]string {
	data, _ := os.ReadFile(j.path)
	logs := make(map[string]map[string]string)
	_ = json.Unmarshal(data, &logs)
	return logs
}

func (j *JournalStore) Clear() { _ = os.Remove(j.path) }

type StudentRegistrationFacade struct {
	checker *subsystemA.PrerequisiteChecker
	repoA   *subsystemA.IntakeRepository
	allocB  *subsystemB.SeatAllocator
	repoB   *subsystemB.SeatRepository
	billC   *subsystemC.BillingSystem
	repoC   *subsystemC.ConfirmationRepository
	journal *JournalStore
}

func (f *StudentRegistrationFacade) RegisterCourse(regID, studentID, courseID string, credits int) error {
	log.Printf("[FACADE] Starting chain for %s (Course: %s)", studentID, courseID)

	if err := f.checker.Validate(studentID, courseID); err != nil {
		f.journal.Log(regID, "FAILED", err.Error())
		return fmt.Errorf("step 0 check failed: %w", err)
	}

	recA := &subsystemA.IntakeRecord{RegistrationID: regID, StudentID: studentID, CourseID: courseID, Status: "RESERVED"}
	f.repoA.Save(recA)
	f.journal.Log(regID, "PENDING_B", "A committed")

	seatNum, err := f.allocB.Allocate(courseID)
	if err != nil {
		log.Printf("[COMPENSATION] Step B failed! Rolling back Subsystem A...")
		f.compensateA(regID, err.Error())
		return fmt.Errorf("step B failed, A rolled back: %w", err)
	}
	f.repoB.Save(&subsystemB.SeatAllocation{RegistrationID: regID, CourseID: courseID, SeatNumber: seatNum, Status: "ALLOCATED"})
	f.journal.Log(regID, "PENDING_C", "B committed")

	amt, inv, err := f.billC.Charge(studentID, credits)
	if err != nil {
		log.Printf("[COMPENSATION] Step C failed! Rolling back B then A...")
		f.compensateBThenA(regID, courseID, err.Error())
		return fmt.Errorf("step C failed, B and A rolled back: %w", err)
	}
	f.repoC.Save(&subsystemC.Confirmation{RegistrationID: regID, Amount: amt, Invoice: inv, Status: "CONFIRMED"})

	f.journal.Log(regID, "COMPLETED", "All 3 steps successful")
	log.Printf("[FACADE] SUCCESS: All 3 stages committed!")
	return nil
}

func (f *StudentRegistrationFacade) compensateA(regID, reason string) {
	if a := f.repoA.Find(regID); a != nil {
		_ = a.TransitionTo("CANCELLED")
		f.repoA.Save(a)
	}
	f.journal.Log(regID, "COMPENSATED", "B failed: "+reason)
}

func (f *StudentRegistrationFacade) compensateBThenA(regID, courseID, reason string) {
	if b := f.repoB.Find(regID); b != nil {
		_ = b.TransitionTo("RELEASED")
		f.repoB.Save(b)
		f.allocB.Release(courseID)
	}
	if a := f.repoA.Find(regID); a != nil {
		_ = a.TransitionTo("CANCELLED")
		f.repoA.Save(a)
	}
	f.journal.Log(regID, "COMPENSATED", "C failed: "+reason)
}

type AdminInspectionFacade struct {
	repoA   *subsystemA.IntakeRepository
	repoB   *subsystemB.SeatRepository
	repoC   *subsystemC.ConfirmationRepository
	journal *JournalStore
	student *StudentRegistrationFacade
}

func (a *AdminInspectionFacade) Inspect(regID string) {
	logEntry := a.journal.GetAll()[regID]
	statA, statB, statC := "NOT_FOUND", "NOT_FOUND", "NOT_FOUND"
	if rec := a.repoA.Find(regID); rec != nil {
		statA = rec.Status
	}
	if rec := a.repoB.Find(regID); rec != nil {
		statB = rec.Status
	}
	if rec := a.repoC.Find(regID); rec != nil {
		statC = rec.Status
	}
	fmt.Printf("\nPERSISTED STATE FOR %s\n", regID)
	fmt.Printf("Journal:      %v\n", logEntry["state"])
	fmt.Printf("Subsystem A:  %s\n", statA)
	fmt.Printf("Subsystem B:  %s\n", statB)
	fmt.Printf("Subsystem C:  %s\n", statC)
}

func (a *AdminInspectionFacade) Recover() {
	for id, l := range a.journal.GetAll() {
		switch l["state"] {
		case "PENDING_B":
			a.student.compensateA(id, "Crash reboot recovery")
			fmt.Printf("Recovered stalled chain %s (A rolled back)\n", id)
		case "PENDING_C":
			a.student.compensateBThenA(id, "CS201", "Crash reboot recovery")
			fmt.Printf("Recovered stalled chain %s (B & A rolled back)\n", id)
		default:
			fmt.Println("Nothing to recover")
		}

	}
}
