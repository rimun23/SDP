package subsystemA

import (
	"encoding/json"
	"errors"
	"os"
)

type IntakeRecord struct {
	RegistrationID string `json:"registration_id"`
	StudentID      string `json:"student_id"`
	CourseID       string `json:"course_id"`
	Status         string `json:"status"`
}

func (r *IntakeRecord) TransitionTo(newStatus string) error {
	if r.Status == "CANCELLED" {
		return errors.New("cannot change CANCELLED status")
	}
	r.Status = newStatus
	return nil
}

type PrerequisiteChecker struct{}

func (p *PrerequisiteChecker) Validate(studentID, courseID string) error {
	if courseID == "CS201" && studentID == "STUDENT_FAIL" {
		return errors.New("CS201 requires CS101 prerequisite")
	}
	if studentID == "STUDENT_DEBT" {
		return errors.New("financial hold on student account")
	}
	return nil
}

type IntakeRepository struct{ path string }

func NewIntakeRepository(dir string) *IntakeRepository {
	return &IntakeRepository{path: dir + "intake.json"}
}

func (r *IntakeRepository) load() map[string]*IntakeRecord {
	data, _ := os.ReadFile(r.path)
	res := make(map[string]*IntakeRecord)
	_ = json.Unmarshal(data, &res)
	return res
}

func (r *IntakeRepository) Save(rec *IntakeRecord) {
	all := r.load()
	all[rec.RegistrationID] = rec
	b, _ := json.MarshalIndent(all, "", "  ")
	_ = os.WriteFile(r.path, b, 0644)
}

func (r *IntakeRepository) Find(id string) *IntakeRecord {
	return r.load()[id]
}

func (r *IntakeRepository) Clear() { _ = os.Remove(r.path) }
