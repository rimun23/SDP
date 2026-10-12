package subsystemC

import (
	"errors"
	"fmt"
)

var (
	ErrBillingLimitExceeded = errors.New("subsystem_c: student exceeded credit tuition billing threshold")
	ErrTranscriptLocked     = errors.New("subsystem_c: registrar locked transcript for modifications")
)

type BillingSystem struct {
	ratePerCredit float64
	maxCreditDebt map[string]float64
}

func NewBillingSystem() *BillingSystem {
	return &BillingSystem{
		ratePerCredit: 500.0,
		maxCreditDebt: map[string]float64{
			"BROKE_STUDENT": 100.0,
		},
	}
}
func (b *BillingSystem) ProcessTuitionCharge(studentID string, credits int) (float64, string, error) {
	amount := float64(credits) * b.ratePerCredit
	if limit, exists := b.maxCreditDebt[studentID]; exists && amount > limit {
		return 0, "", fmt.Errorf("%w: charge %.2f exceeds max allowed %.2f", ErrBillingLimitExceeded, amount, limit)
	}
	invoice := fmt.Sprintf("INV-%s-%d", studentID, credits)
	return amount, invoice, nil
}

type TranscriptService struct{}

func NewTranscriptService() *TranscriptService {
	return &TranscriptService{}
}
func (t *TranscriptService) RecordTermEntry(studentID, courseID string) error {
	if studentID == "LOCKED_STUDENT" {
		return ErrTranscriptLocked
	}
	return nil
}
