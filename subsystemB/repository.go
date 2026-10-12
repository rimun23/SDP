package subsystemB

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type SeatRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewSeatRepository(storageDir string) (*SeatRepository, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}
	path := filepath.Join(storageDir, "seats.json")
	return &SeatRepository{filePath: path}, nil
}

func (r *SeatRepository) loadAll() (map[string]*SeatAllocation, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return make(map[string]*SeatAllocation), nil
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return make(map[string]*SeatAllocation), nil
	}
	var records map[string]*SeatAllocation
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *SeatRepository) Save(alloc *SeatAllocation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	records, err := r.loadAll()
	if err != nil {
		return err
	}
	records[alloc.RegistrationID] = alloc
	bytes, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, bytes, 0644)
}

func (r *SeatRepository) FindByID(regID string) (*SeatAllocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	records, err := r.loadAll()
	if err != nil {
		return nil, err
	}
	alloc, exists := records[regID]
	if !exists {
		return nil, fmt.Errorf("subsystem_b: allocation %s not found", regID)
	}
	return alloc, nil
}

func (r *SeatRepository) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return os.Remove(r.filePath)
}
