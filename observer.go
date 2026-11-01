package main

import (
	"fmt"
	"sync"
)

type EventType string

const (
	EventGoal       EventType = "GOAL"
	EventInjury     EventType = "INJURY"
	EventMatchCard  EventType = "CARD"
	EventAggregated EventType = "AGGREGATED_STATS"
)

type Event struct {
	Type    EventType
	ClubID  string
	TeamID  string
	Payload string
	Value   float64
}

type Observer interface {
	OnEvent(event Event)
}

type Subject interface {
	Subscribe(eventType EventType, obs Observer)
	Unsubscribe(eventType EventType, obs Observer)
	Notify(event Event)
}

type EventBroker struct {
	mu        sync.RWMutex
	listeners map[EventType][]Observer
}

func NewEventBroker() *EventBroker {
	return &EventBroker{listeners: make(map[EventType][]Observer)}
}

func (b *EventBroker) Subscribe(eventType EventType, obs Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[eventType] = append(b.listeners[eventType], obs)
}

func (b *EventBroker) Unsubscribe(eventType EventType, obs Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	list := b.listeners[eventType]
	for i, l := range list {
		if l == obs {
			b.listeners[eventType] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func (b *EventBroker) Notify(event Event) {
	b.mu.RLock()
	rawListeners := b.listeners[event.Type]
	targets := make([]Observer, len(rawListeners))
	copy(targets, rawListeners)
	b.mu.RUnlock()

	for _, obs := range targets {
		obs.OnEvent(event)
	}
}

type Scoreboard struct{}

func (s *Scoreboard) OnEvent(e Event) {
	if e.Type == EventGoal {
		fmt.Printf("Scoreboard: GOAL! Club %s: %s\n", e.ClubID, e.Payload)
	}
}

type OddsCalculator struct{}

func (o *OddsCalculator) OnEvent(e Event) {
	fmt.Printf("OddsCalculator: Updating coefficients due to %s on %s\n", e.Type, e.Payload)
}

type FantasyPointsTracker struct{}

func (f *FantasyPointsTracker) OnEvent(e Event) {
	fmt.Printf("FantasyTracker: Updating points for team %s (event: %s)\n", e.TeamID, e.Type)
}

type PushNotifier struct{}

func (p *PushNotifier) OnEvent(e Event) {
	fmt.Printf("PushNotification: Alert sent to fans: %s -> %s\n", e.Type, e.Payload)
}

type MatchAggregator struct {
	broker    *EventBroker
	outBroker *EventBroker
	mu        sync.Mutex
	goalCount int
}

func NewMatchAggregator(out *EventBroker) *MatchAggregator {
	return &MatchAggregator{outBroker: out}
}

func (m *MatchAggregator) OnEvent(e Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e.Type == EventGoal {
		m.goalCount++
		m.outBroker.Notify(Event{
			Type:    EventAggregated,
			ClubID:  e.ClubID,
			Payload: fmt.Sprintf("Total match goals so far: %d", m.goalCount),
			Value:   float64(m.goalCount),
		})
	}
}
