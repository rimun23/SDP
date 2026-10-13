package subsystemC

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Confirmation struct {
	RegistrationID string  `json:"registration_id"`
	Amount         float64 `json:"amount"`
	Invoice        string  `json:"invoice"`
	Status         string  `json:"status"`
}

func (c *Confirmation) TransitionTo(newStatus string) error {
	if c.Status == "VOIDED" {
		return errors.New("confirmation is voided")
	}
	c.Status = newStatus
	return nil
}

type BillingSystem struct{}

func (b *BillingSystem) Charge(studentID string, credits int) (float64, string, error) {
	if studentID == "BROKE_STUDENT" {
		return 0, "", errors.New("tuition billing debt limit exceeded")
	}
	return float64(credits * 500), fmt.Sprintf("INV-%s", studentID), nil
}

type ConfirmationRepository struct{ path string }

func NewConfirmationRepository(dir string) *ConfirmationRepository {
	return &ConfirmationRepository{path: dir + "confirmation.json"}
}

func (r *ConfirmationRepository) load() map[string]*Confirmation {
	data, _ := os.ReadFile(r.path)
	res := make(map[string]*Confirmation)
	_ = json.Unmarshal(data, &res)
	return res
}

func (r *ConfirmationRepository) Save(conf *Confirmation) {
	all := r.load()
	all[conf.RegistrationID] = conf
	b, _ := json.MarshalIndent(all, "", "  ")
	_ = os.WriteFile(r.path, b, 0644)
}

func (r *ConfirmationRepository) Find(id string) *Confirmation {
	return r.load()[id]
}

func (r *ConfirmationRepository) Clear() { _ = os.Remove(r.path) }
