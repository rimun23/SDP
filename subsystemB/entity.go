package subsystemB

import (
	"errors"
	"fmt"
	"time"
)

type AllocationStatus string

const (
	StatusAllocPending AllocationStatus = "PENDING"
	StatusAllocated    AllocationStatus = "ALLOCATED"
	StatusReleased     AllocationStatus = "RELEASED"
)

type SeatAllocation struct {
	RegistrationID string           `json:"registration_id"`
	CourseID       string           `json:"course_id"`
	SeatNumber     int              `json:"seat_number"`
	Status         AllocationStatus `json:"status"`
	AllocatedAt    time.Time        `json:"allocated_at"`
}

var ErrInvalidSeatTransition = errors.New("subsystem_b: invalid seat status transition")

func NewSeatAllocation(regID, courseID string, seatNum int) *SeatAllocation {
	return &SeatAllocation{
		RegistrationID: regID,
		CourseID:       courseID,
		SeatNumber:     seatNum,
		Status:         StatusAllocPending,
		AllocatedAt:    time.Now(),
	}
}

func (s *SeatAllocation) TransitionTo(newStatus AllocationStatus) error {
	switch s.Status {
	case StatusAllocPending:
		if newStatus == StatusAllocated || newStatus == StatusReleased {
			s.Status = newStatus
			return nil
		}
	case StatusAllocated:
		if newStatus == StatusReleased {
			s.Status = newStatus
			return nil
		}
	case StatusReleased:
		return fmt.Errorf("%w: cannot modify RELEASED seat", ErrInvalidSeatTransition)
	}
	return fmt.Errorf("%w: invalid transition %s -> %s", ErrInvalidSeatTransition, s.Status, newStatus)
}
