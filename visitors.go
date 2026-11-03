package main

import "fmt"

type SportsVisitor interface {
	VisitClub(c *Club)
	VisitTeam(t *Team)
	VisitForward(f *Forward)
	VisitDefender(d *Defender)
	VisitGoalkeeper(g *Goalkeeper)
	VisitCoach(c *Coach)
}

func traverseTeam(v SportsVisitor, t *Team) {
	for _, f := range t.Forwards {
		f.Accept(v)
	}
	for _, d := range t.Defenders {
		d.Accept(v)
	}
	for _, g := range t.Goalkeepers {
		g.Accept(v)
	}
	for _, c := range t.Coaches {
		c.Accept(v)
	}
}

type SalaryCalculatorVisitor struct {
	TotalSalary float64
}

func (v *SalaryCalculatorVisitor) VisitClub(c *Club) {
	for _, team := range c.Teams {
		team.Accept(v)
	}
}
func (v *SalaryCalculatorVisitor) VisitTeam(t *Team)         { traverseTeam(v, t) }
func (v *SalaryCalculatorVisitor) VisitForward(f *Forward)   { v.TotalSalary += f.Salary }
func (v *SalaryCalculatorVisitor) VisitDefender(d *Defender) { v.TotalSalary += d.Salary }
func (v *SalaryCalculatorVisitor) VisitGoalkeeper(g *Goalkeeper) {
	v.TotalSalary += g.Salary
}
func (v *SalaryCalculatorVisitor) VisitCoach(c *Coach) { v.TotalSalary += c.Salary }

type RatingAccumulatorVisitor struct {
	totalRating float64
	count       int
}

func (v *RatingAccumulatorVisitor) VisitClub(c *Club) {
	for _, team := range c.Teams {
		team.Accept(v)
	}
}
func (v *RatingAccumulatorVisitor) VisitTeam(t *Team) { traverseTeam(v, t) }
func (v *RatingAccumulatorVisitor) VisitForward(f *Forward) {
	v.totalRating += f.Rating
	v.count++
}
func (v *RatingAccumulatorVisitor) VisitDefender(d *Defender) {
	v.totalRating += d.Rating
	v.count++
}
func (v *RatingAccumulatorVisitor) VisitGoalkeeper(g *Goalkeeper) {
	v.totalRating += g.Rating
	v.count++
}
func (v *RatingAccumulatorVisitor) VisitCoach(c *Coach) {}

func (v *RatingAccumulatorVisitor) AverageRating() float64 {
	if v.count == 0 {
		return 0
	}
	return v.totalRating / float64(v.count)
}

type StarFilterVisitor struct {
	Threshold float64
	Found     []string
}

func (v *StarFilterVisitor) VisitClub(c *Club) {
	for _, team := range c.Teams {
		team.Accept(v)
	}
}
func (v *StarFilterVisitor) VisitTeam(t *Team) { traverseTeam(v, t) }
func (v *StarFilterVisitor) VisitForward(f *Forward) {
	if f.Rating >= v.Threshold {
		v.Found = append(v.Found, f.Name)
	}
}
func (v *StarFilterVisitor) VisitDefender(d *Defender) {
	if d.Rating >= v.Threshold {
		v.Found = append(v.Found, d.Name)
	}
}
func (v *StarFilterVisitor) VisitGoalkeeper(g *Goalkeeper) {
	if g.Rating >= v.Threshold {
		v.Found = append(v.Found, g.Name)
	}
}
func (v *StarFilterVisitor) VisitCoach(c *Coach) {}

type JsonExportVisitor struct {
	result string
}

func (v *JsonExportVisitor) VisitClub(c *Club) {
	v.result += fmt.Sprintf(`{"club":"%s","teams":[`, c.Name)
	for i, team := range c.Teams {
		if i > 0 {
			v.result += ","
		}
		team.Accept(v)
	}
	v.result += "]}"
}

func (v *JsonExportVisitor) VisitTeam(t *Team) {
	v.result += fmt.Sprintf(`{"team":"%s","members":[`, t.Name)
	first := true
	appendItem := func(jsonStr string) {
		if !first {
			v.result += ","
		}
		v.result += jsonStr
		first = false
	}

	for _, f := range t.Forwards {
		appendItem(fmt.Sprintf(`{"type":"Forward","name":"%s"}`, f.Name))
	}
	for _, d := range t.Defenders {
		appendItem(fmt.Sprintf(`{"type":"Defender","name":"%s"}`, d.Name))
	}
	for _, g := range t.Goalkeepers {
		appendItem(fmt.Sprintf(`{"type":"Goalkeeper","name":"%s"}`, g.Name))
	}
	for _, c := range t.Coaches {
		appendItem(fmt.Sprintf(`{"type":"Coach","name":"%s"}`, c.Name))
	}
	v.result += "]}"
}

func (v *JsonExportVisitor) VisitForward(f *Forward)       {}
func (v *JsonExportVisitor) VisitDefender(d *Defender)     {}
func (v *JsonExportVisitor) VisitGoalkeeper(g *Goalkeeper) {}
func (v *JsonExportVisitor) VisitCoach(c *Coach)           {}

func (v *JsonExportVisitor) Output() string { return v.result }
