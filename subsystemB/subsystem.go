package subsystemB

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type SeatAllocation struct {
	RegistrationID string `json:"registration_id"`
	CourseID       string `json:"course_id"`
	SeatNumber     int    `json:"seat_number"`
	Status         string `json:"status"`
}

func (s *SeatAllocation) TransitionTo(newStatus string) error {
	if s.Status == "RELEASED" {
		return errors.New("seat already released")
	}
	s.Status = newStatus
	return nil
}

type SeatAllocator struct {
	occupied map[string]int
}

func NewSeatAllocator() *SeatAllocator {
	return &SeatAllocator{occupied: make(map[string]int)}
}

func (s *SeatAllocator) Allocate(courseID string) (int, error) {
	if courseID == "FULL_COURSE" || s.occupied[courseID] >= 2 {
		return 0, fmt.Errorf("no seats available for %s", courseID)
	}
	s.occupied[courseID]++
	return s.occupied[courseID], nil
}

func (s *SeatAllocator) Release(courseID string) {
	if s.occupied[courseID] > 0 {
		s.occupied[courseID]--
	}
}

func (s *SeatAllocator) Reset() { s.occupied = make(map[string]int) }

type SeatRepository struct{ path string }

func NewSeatRepository(dir string) *SeatRepository {
	return &SeatRepository{path: dir + "seats.json"}
}

func (r *SeatRepository) load() map[string]*SeatAllocation {
	data, _ := os.ReadFile(r.path)
	res := make(map[string]*SeatAllocation)
	_ = json.Unmarshal(data, &res)
	return res
}

func (r *SeatRepository) Save(alloc *SeatAllocation) {
	all := r.load()
	all[alloc.RegistrationID] = alloc
	b, _ := json.MarshalIndent(all, "", "  ")
	_ = os.WriteFile(r.path, b, 0644)
}

func (r *SeatRepository) Find(id string) *SeatAllocation {
	return r.load()[id]
}

func (r *SeatRepository) Clear() { _ = os.Remove(r.path) }
