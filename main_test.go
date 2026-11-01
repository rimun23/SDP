package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSalaryCalculatorVisitor(t *testing.T) {
	team := &Team{
		Name: "Test Team",
		Members: []SportsElement{
			&Forward{Name: "F", Salary: 100},
			&Defender{Name: "D", Salary: 150},
			&Goalkeeper{Name: "G", Salary: 80},
			&Coach{Name: "C", Salary: 200},
		},
	}
	v := &SalaryCalculatorVisitor{}
	team.Accept(v)

	const expected = 530.0
	if v.TotalSalary != expected {
		t.Fatalf("Expected total salary %.2f, got %.2f", expected, v.TotalSalary)
	}
}

func TestRatingAccumulatorVisitor(t *testing.T) {
	team := &Team{
		Name: "Test Team",
		Members: []SportsElement{
			&Forward{Name: "F", Rating: 9.0},
			&Defender{Name: "D", Rating: 8.0},
			&Goalkeeper{Name: "G", Rating: 7.0},
			&Coach{Name: "C"},
		},
	}
	v := &RatingAccumulatorVisitor{}
	team.Accept(v)

	const expected = 8.0
	if v.AverageRating() != expected {
		t.Fatalf("Expected avg rating %.2f, got %.2f", expected, v.AverageRating())
	}
}

func TestStarFilterVisitor(t *testing.T) {
	team := &Team{
		Name: "Test Team",
		Members: []SportsElement{
			&Forward{Name: "F1", Rating: 9.5},
			&Forward{Name: "F2", Rating: 6.5},
			&Defender{Name: "D1", Rating: 8.8},
		},
	}
	v := &StarFilterVisitor{Threshold: 8.5}
	team.Accept(v)

	if len(v.Found) != 2 {
		t.Fatalf("Expected 2 star players, got %d", len(v.Found))
	}
}

func TestJsonExportVisitor(t *testing.T) {
	club := &Club{
		Name: "Club1",
		Teams: []*Team{
			{
				Name: "TeamA",
				Members: []SportsElement{
					&Forward{Name: "Player1", Salary: 100, Rating: 9.0},
				},
			},
		},
	}
	v := &JsonExportVisitor{}
	club.Accept(v)
	res := v.Finalize()

	if !strings.Contains(res, `"club":"Club1"`) || !strings.Contains(res, `"Forward"`) {
		t.Fatalf("Invalid JSON produced: %s", res)
	}
}

type SelfUnsubscribingListener struct {
	broker *EventBroker
	called bool
}

func (s *SelfUnsubscribingListener) OnEvent(e Event) {
	s.called = true
	s.broker.Unsubscribe(e.Type, s)
}

func TestUnsubscribeInsideHandler(t *testing.T) {
	broker := NewEventBroker()
	listener := &SelfUnsubscribingListener{broker: broker}
	broker.Subscribe(EventGoal, listener)

	broker.Notify(Event{Type: EventGoal})
	if !listener.called {
		t.Fatal("Listener should have been called")
	}

	listener.called = false
	broker.Notify(Event{Type: EventGoal})
	if listener.called {
		t.Fatal("Listener should have unsubscribed itself")
	}
}

type CountingObserver struct {
	count int64
}

func (c *CountingObserver) OnEvent(e Event) {
	atomic.AddInt64(&c.count, 1)
}

func TestStressConcurrentDispatch(t *testing.T) {
	broker := NewEventBroker()
	counter := &CountingObserver{}
	broker.Subscribe(EventGoal, counter)

	threads := 4
	iterations := 1000
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				broker.Notify(Event{
					Type:    EventGoal,
					ClubID:  fmt.Sprintf("Club-%d", id),
					Payload: "Goal",
				})
			}
		}(i)
	}

	wg.Wait()

	expected := int64(threads * iterations)
	if atomic.LoadInt64(&counter.count) != expected {
		t.Fatalf("Expected exact %d events, but got %d", expected, counter.count)
	}
}
