package subsystemC

import (
	"errors"
	"fmt"
	"time"
)

type ConfirmationStatus string

const (
	StatusConfPending ConfirmationStatus = "PENDING"
	StatusConfirmed   ConfirmationStatus = "CONFIRMED"
	StatusVoided      ConfirmationStatus = "VOIDED"
)

type EnrollmentConfirmation struct {
	RegistrationID string             `json:"registration_id"`
	BilledAmount   float64            `json:"billed_amount"`
	InvoiceNumber  string             `json:"invoice_number"`
	Status         ConfirmationStatus `json:"status"`
	ConfirmedAt    time.Time          `json:"confirmed_at"`
}

var ErrInvalidConfirmTransition = errors.New("subsystem_c: invalid confirmation status transition")

func NewEnrollmentConfirmation(regID string, amount float64, invoice string) *EnrollmentConfirmation {
	return &EnrollmentConfirmation{
		RegistrationID: regID,
		BilledAmount:   amount,
		InvoiceNumber:  invoice,
		Status:         StatusConfPending,
		ConfirmedAt:    time.Now(),
	}
}

func (c *EnrollmentConfirmation) TransitionTo(newStatus ConfirmationStatus) error {
	switch c.Status {
	case StatusConfPending:
		if newStatus == StatusConfirmed || newStatus == StatusVoided {
			c.Status = newStatus
			return nil
		}
	case StatusConfirmed:
		if newStatus == StatusVoided {
			c.Status = newStatus
			return nil
		}
	case StatusVoided:
		return fmt.Errorf("%w: cannot modify VOIDED confirmation", ErrInvalidConfirmTransition)
	}
	return fmt.Errorf("%w: invalid transition %s -> %s", ErrInvalidConfirmTransition, c.Status, newStatus)
}
