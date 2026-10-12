package subsystemC

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type ConfirmationRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewConfirmationRepository(storageDir string) (*ConfirmationRepository, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}
	path := filepath.Join(storageDir, "confirmation.json")
	return &ConfirmationRepository{filePath: path}, nil
}

func (r *ConfirmationRepository) loadAll() (map[string]*EnrollmentConfirmation, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return make(map[string]*EnrollmentConfirmation), nil
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return make(map[string]*EnrollmentConfirmation), nil
	}
	var records map[string]*EnrollmentConfirmation
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *ConfirmationRepository) Save(conf *EnrollmentConfirmation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	records, err := r.loadAll()
	if err != nil {
		return err
	}
	records[conf.RegistrationID] = conf
	bytes, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, bytes, 0644)
}

func (r *ConfirmationRepository) FindByID(regID string) (*EnrollmentConfirmation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	records, err := r.loadAll()
	if err != nil {
		return nil, err
	}
	conf, exists := records[regID]
	if !exists {
		return nil, fmt.Errorf("subsystem_c: record %s not found", regID)
	}
	return conf, nil
}

func (r *ConfirmationRepository) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return os.Remove(r.filePath)
}
