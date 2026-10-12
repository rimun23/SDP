package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ChainStepState string

const (
	StepNotStarted  ChainStepState = "NOT_STARTED"
	StepACommitted  ChainStepState = "PENDING_B"
	StepBCommitted  ChainStepState = "PENDING_C"
	StepCompleted   ChainStepState = "COMPLETED"
	StepCompensated ChainStepState = "COMPENSATED"
	StepFailed      ChainStepState = "FAILED"
)

type ChainLog struct {
	RegistrationID string         `json:"registration_id"`
	StudentID      string         `json:"student_id"`
	CourseID       string         `json:"course_id"`
	CurrentState   ChainStepState `json:"current_state"`
	FailureReason  string         `json:"failure_reason"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type JournalStore struct {
	filePath string
	mu       sync.RWMutex
}

func NewJournalStore(dir string) (*JournalStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &JournalStore{filePath: filepath.Join(dir, "journal.json")}, nil
}

func (j *JournalStore) loadAll() (map[string]*ChainLog, error) {
	if _, err := os.Stat(j.filePath); os.IsNotExist(err) {
		return make(map[string]*ChainLog), nil
	}
	data, err := os.ReadFile(j.filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return make(map[string]*ChainLog), nil
	}
	var logs map[string]*ChainLog
	if err := json.Unmarshal(data, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (j *JournalStore) RecordStep(regID, studentID, courseID string, state ChainStepState, reason string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	logs, err := j.loadAll()
	if err != nil {
		return err
	}
	logs[regID] = &ChainLog{
		RegistrationID: regID,
		StudentID:      studentID,
		CourseID:       courseID,
		CurrentState:   state,
		FailureReason:  reason,
		UpdatedAt:      time.Now(),
	}
	bytes, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(j.filePath, bytes, 0644)
}

func (j *JournalStore) GetAll() (map[string]*ChainLog, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.loadAll()
}

func (j *JournalStore) Clear() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return os.Remove(j.filePath)
}
