package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MedicalInspectionService struct {
	rootTree *Club
}

func (m *MedicalInspectionService) OnEvent(e Event) {
	if e.Type != EventInjury {
		return
	}
	fmt.Printf("\nInjury detected for %s in team %s. Inspecting team...\n", e.Payload, e.TeamID)
	var targetTeam *Team
	for _, t := range m.rootTree.Teams {
		if t.ID == e.TeamID {
			targetTeam = t
			break
		}
	}
	if targetTeam != nil {
		filter := &StarFilterVisitor{Threshold: 8.5}
		targetTeam.Accept(filter)

		salaryCalc := &SalaryCalculatorVisitor{}
		targetTeam.Accept(salaryCalc)

		fmt.Printf("Inspection Result: Target team %s has %d high-rated players. Payroll: $%.0f\n\n",
			targetTeam.Name, len(filter.Found), salaryCalc.TotalSalary)
	}
}
func main() { // Инициализация структуры данных в main.go
	firstTeam := &Team{
		ID:   "T1",
		Name: "First Squad",
		Forwards: []*Forward{
			{Name: "Haaland", Salary: 500000, Rating: 9.3, Goals: 30},
		},
		Defenders: []*Defender{
			{Name: "Dias", Salary: 250000, Rating: 8.7, Tackles: 45},
		},
		Goalkeepers: []*Goalkeeper{
			{Name: "Ederson", Salary: 200000, Rating: 8.2, CleanSheets: 12},
		},
		Coaches: []*Coach{
			{Name: "Guardiola", Salary: 600000, TacticsExp: 20},
		},
	}
	reserves := &Team{
		ID:   "T2",
		Name: "Reserves",
		Forwards: []*Forward{
			{Name: "Bobb", Salary: 50000, Rating: 7.4, Goals: 4},
		},
	}
	club := &Club{ID: "C1", Name: "Manchester City", Teams: []*Team{firstTeam, reserves}}

	broker := NewEventBroker()
	aggregatorBroker := NewEventBroker()

	scoreboard := &Scoreboard{}
	oddsCalc := &OddsCalculator{}
	fantasy := &FantasyPointsTracker{}
	push := &PushNotifier{}
	aggregator := NewMatchAggregator(aggregatorBroker)
	inspection := &MedicalInspectionService{rootTree: club}

	broker.Subscribe(EventGoal, scoreboard)
	broker.Subscribe(EventGoal, oddsCalc)
	broker.Subscribe(EventGoal, fantasy)
	broker.Subscribe(EventGoal, push)
	broker.Subscribe(EventGoal, aggregator)

	broker.Subscribe(EventInjury, push)
	broker.Subscribe(EventInjury, inspection)

	aggregatorBroker.Subscribe(EventAggregated, push)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		broker.Notify(Event{Type: EventGoal, ClubID: "C1", TeamID: "T1", Payload: "Haaland 12'"})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		broker.Notify(Event{Type: EventMatchCard, ClubID: "C1", TeamID: "T1", Payload: "Yellow card Dias"})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		broker.Notify(Event{Type: EventInjury, ClubID: "C1", TeamID: "T1", Payload: "Haaland"})
	}()

	<-ctx.Done()
	wg.Wait()
	fmt.Println("All producer events completed gracefully.")
}
