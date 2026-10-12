package subsystemB

import (
	"errors"
	"fmt"
)

var (
	ErrNoSeatsAvailable = errors.New("subsystem_b: course capacity reached, no seats available")
	ErrWaitlistFull     = errors.New("subsystem_b: waitlist queue is completely full")
)

type SeatAllocator struct {
	capacityPerCourse map[string]int
	occupiedSeats     map[string]int
}

func NewSeatAllocator() *SeatAllocator {
	return &SeatAllocator{
		capacityPerCourse: map[string]int{
			"CS201":       2,
			"FULL_COURSE": 0,
		},
		occupiedSeats: make(map[string]int),
	}
}
func (sa *SeatAllocator) Allocate(courseID string) (int, error) {
	cap, exists := sa.capacityPerCourse[courseID]
	if !exists {
		cap = 30
	}
	occupied := sa.occupiedSeats[courseID]
	if occupied >= cap {
		return 0, fmt.Errorf("%w for course %s", ErrNoSeatsAvailable, courseID)
	}
	sa.occupiedSeats[courseID]++
	return sa.occupiedSeats[courseID], nil
}
func (sa *SeatAllocator) Deallocate(courseID string) {
	if sa.occupiedSeats[courseID] > 0 {
		sa.occupiedSeats[courseID]--
	}
}
func (sa *SeatAllocator) Reset() {
	sa.occupiedSeats = make(map[string]int)
}

type WaitlistManager struct {
	maxWaitlist int
	currentWait map[string]int
}

func NewWaitlistManager() *WaitlistManager {
	return &WaitlistManager{
		maxWaitlist: 5,
		currentWait: make(map[string]int),
	}
}
func (w *WaitlistManager) CanJoin(courseID string) bool {
	return w.currentWait[courseID] < w.maxWaitlist
}
