package main

// Интерфейс исключительно для реализации Double Dispatch в Visitor
type Element interface {
	Accept(v SportsVisitor)
}

// 2 Контейнера

type Club struct {
	ID    string
	Name  string
	Teams []*Team
}

func (c *Club) Accept(v SportsVisitor) { v.VisitClub(c) }

type Team struct {
	ID          string
	Name        string
	Forwards    []*Forward
	Defenders   []*Defender
	Goalkeepers []*Goalkeeper
	Coaches     []*Coach
}

func (t *Team) Accept(v SportsVisitor) { v.VisitTeam(t) }

// 4 Листовых типа (только данные)

type Forward struct {
	Name   string
	Salary float64
	Rating float64
	Goals  int
}

func (f *Forward) Accept(v SportsVisitor) { v.VisitForward(f) }

type Defender struct {
	Name    string
	Salary  float64
	Rating  float64
	Tackles int
}

func (d *Defender) Accept(v SportsVisitor) { v.VisitDefender(d) }

type Goalkeeper struct {
	Name        string
	Salary      float64
	Rating      float64
	CleanSheets int
}

func (g *Goalkeeper) Accept(v SportsVisitor) { v.VisitGoalkeeper(g) }

type Coach struct {
	Name       string
	Salary     float64
	TacticsExp int
}

func (c *Coach) Accept(v SportsVisitor) { v.VisitCoach(c) }
