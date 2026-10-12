package subsystemA

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type IntakeRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewIntakeRepository(storageDir string) (*IntakeRepository, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}
	path := filepath.Join(storageDir, "intake.json")
	return &IntakeRepository{filePath: path}, nil
}
func (r *IntakeRepository) loadAll() (map[string]*IntakeRecord, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return make(map[string]*IntakeRecord), nil
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return make(map[string]*IntakeRecord), nil
	}
	var records map[string]*IntakeRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}
func (r *IntakeRepository) Save(record *IntakeRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	records, err := r.loadAll()
	if err != nil {
		return err
	}
	records[record.RegistrationID] = record
	bytes, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, bytes, 0644)
}
func (r *IntakeRepository) FindByID(regID string) (*IntakeRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	records, err := r.loadAll()
	if err != nil {
		return nil, err
	}
	rec, exists := records[regID]
	if !exists {
		return nil, fmt.Errorf("subsystem_a: record %s not found", regID)
	}
	return rec, nil
}
func (r *IntakeRepository) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return os.Remove(r.filePath)
}
