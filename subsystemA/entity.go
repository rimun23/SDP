package subsystemA

import (
	"errors"
	"fmt"
	"time"
)

type IntakeStatus string

const (
	StatusNew       IntakeStatus = "NEW"
	StatusReserved  IntakeStatus = "RESERVED"
	StatusCancelled IntakeStatus = "CANCELLED"
)

type IntakeRecord struct {
	RegistrationID string       `json:"registration_id"`
	StudentID      string       `json:"student_id"`
	CourseID       string       `json:"course_id"`
	Status         IntakeStatus `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
}

var (
	ErrInvalidTransition = errors.New("subsystem_a: illegal status transition")
)

func NewIntakeRecord(regID, studentID, courseID string) *IntakeRecord {
	return &IntakeRecord{
		RegistrationID: regID,
		StudentID:      studentID,
		CourseID:       courseID,
		Status:         StatusNew,
		CreatedAt:      time.Now(),
	}
}
func (r *IntakeRecord) TransitionTo(newStatus IntakeStatus) error {
	switch r.Status {
	case StatusNew:
		if newStatus == StatusReserved || newStatus == StatusCancelled {
			r.Status = newStatus
			return nil
		}
	case StatusReserved:
		if newStatus == StatusCancelled {
			r.Status = newStatus
			return nil
		}
	case StatusCancelled:
		return fmt.Errorf("%w: cannot transition from CANCELLED", ErrInvalidTransition)
	}
	return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidTransition, r.Status, newStatus)
}
