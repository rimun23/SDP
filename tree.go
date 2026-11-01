package main

type SportsElement interface {
	Accept(v SportsVisitor)
	GetName() string
}

type SportsVisitor interface {
	VisitClub(c *Club)
	VisitTeam(t *Team)
	VisitForward(f *Forward)
	VisitDefender(d *Defender)
	VisitGoalkeeper(g *Goalkeeper)
	VisitCoach(c *Coach)
}

type Club struct {
	ID    string
	Name  string
	Teams []*Team
}

func (c *Club) GetName() string { return c.Name }
func (c *Club) Accept(v SportsVisitor) {
	v.VisitClub(c)
	for _, team := range c.Teams {
		team.Accept(v)
	}
}

type Team struct {
	ID      string
	Name    string
	Members []SportsElement
}

func (t *Team) GetName() string { return t.Name }
func (t *Team) Accept(v SportsVisitor) {
	v.VisitTeam(t)
	for _, member := range t.Members {
		member.Accept(v)
	}
}

type Forward struct {
	Name   string
	Salary float64
	Rating float64
	Goals  int
}

func (f *Forward) GetName() string        { return f.Name }
func (f *Forward) Accept(v SportsVisitor) { v.VisitForward(f) }

type Defender struct {
	Name    string
	Salary  float64
	Rating  float64
	Tackles int
}

func (d *Defender) GetName() string        { return d.Name }
func (d *Defender) Accept(v SportsVisitor) { v.VisitDefender(d) }

type Goalkeeper struct {
	Name       string
	Salary     float64
	Rating     float64
	CleanSheet int
}

func (g *Goalkeeper) GetName() string        { return g.Name }
func (g *Goalkeeper) Accept(v SportsVisitor) { v.VisitGoalkeeper(g) }

type Coach struct {
	Name       string
	Salary     float64
	TacticsExp int
}

func (c *Coach) GetName() string        { return c.Name }
func (c *Coach) Accept(v SportsVisitor) { v.VisitCoach(c) }
