package main

import (
	"SDP/subsystemA"
	"SDP/subsystemB"
	"SDP/subsystemC"
	"fmt"
)

type SystemStatusSummary struct {
	RegistrationID string
	JournalState   ChainStepState
	IntakeStatus   string
	SeatStatus     string
	ConfirmStatus  string
}
type AdminInspectionFacade struct {
	intakeRepo  *subsystemA.IntakeRepository
	seatRepo    *subsystemB.SeatRepository
	confirmRepo *subsystemC.ConfirmationRepository
	journal     *JournalStore
	studentFac  *StudentRegistrationFacade
}

func NewAdminInspectionFacade(
	ra *subsystemA.IntakeRepository,
	rb *subsystemB.SeatRepository,
	rc *subsystemC.ConfirmationRepository,
	j *JournalStore,
	sf *StudentRegistrationFacade,
) *AdminInspectionFacade {
	return &AdminInspectionFacade{
		intakeRepo:  ra,
		seatRepo:    rb,
		confirmRepo: rc,
		journal:     j,
		studentFac:  sf,
	}
}
func (a *AdminInspectionFacade) InspectRegistration(regID string) (*SystemStatusSummary, error) {
	summary := &SystemStatusSummary{RegistrationID: regID}
	if logs, err := a.journal.GetAll(); err == nil {
		if l, ok := logs[regID]; ok {
			summary.JournalState = l.CurrentState
		}
	}
	if rec, err := a.intakeRepo.FindByID(regID); err == nil {
		summary.IntakeStatus = string(rec.Status)
	} else {
		summary.IntakeStatus = "NOT_FOUND"
	}
	if seat, err := a.seatRepo.FindByID(regID); err == nil {
		summary.SeatStatus = string(seat.Status)
	} else {
		summary.SeatStatus = "NOT_FOUND"
	}
	if conf, err := a.confirmRepo.FindByID(regID); err == nil {
		summary.ConfirmStatus = string(conf.Status)
	} else {
		summary.ConfirmStatus = "NOT_FOUND"
	}
	return summary, nil
}
func (a *AdminInspectionFacade) RecoverStalledChains() ([]string, error) {
	logs, err := a.journal.GetAll()
	if err != nil {
		return nil, err
	}
	var recovered []string
	for regID, logEntry := range logs {
		switch logEntry.CurrentState {
		case StepACommitted:
			a.studentFac.compensateA(regID, logEntry.StudentID, logEntry.CourseID, "Recovered from reboot crash at Step A")
			recovered = append(recovered, fmt.Sprintf("%s (repaired PENDING_B -> COMPENSATED)", regID))
		case StepBCommitted:
			a.studentFac.compensateBThenA(regID, logEntry.StudentID, logEntry.CourseID, "Recovered from reboot crash at Step B")
			recovered = append(recovered, fmt.Sprintf("%s (repaired PENDING_C -> COMPENSATED)", regID))
		}
	}
	return recovered, nil
}
