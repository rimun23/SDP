package main

import "fmt"

type SalaryCalculatorVisitor struct {
	TotalSalary float64
}

func (v *SalaryCalculatorVisitor) VisitClub(c *Club)         {}
func (v *SalaryCalculatorVisitor) VisitTeam(t *Team)         {}
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

func (v *RatingAccumulatorVisitor) VisitClub(c *Club)             {}
func (v *RatingAccumulatorVisitor) VisitTeam(t *Team)             {}
func (v *RatingAccumulatorVisitor) VisitForward(f *Forward)       { v.add(f.Rating) }
func (v *RatingAccumulatorVisitor) VisitDefender(d *Defender)     { v.add(d.Rating) }
func (v *RatingAccumulatorVisitor) VisitGoalkeeper(g *Goalkeeper) { v.add(g.Rating) }
func (v *RatingAccumulatorVisitor) VisitCoach(c *Coach)           {}

func (v *RatingAccumulatorVisitor) add(r float64) {
	v.totalRating += r
	v.count++
}

func (v *RatingAccumulatorVisitor) AverageRating() float64 {
	if v.count == 0 {
		return 0
	}
	return v.totalRating / float64(v.count)
}

type StarFilterVisitor struct {
	Threshold float64
	Found     []SportsElement
}

func (v *StarFilterVisitor) VisitClub(c *Club) {}
func (v *StarFilterVisitor) VisitTeam(t *Team) {}
func (v *StarFilterVisitor) VisitForward(f *Forward) {
	if f.Rating >= v.Threshold {
		v.Found = append(v.Found, f)
	}
}
func (v *StarFilterVisitor) VisitDefender(d *Defender) {
	if d.Rating >= v.Threshold {
		v.Found = append(v.Found, d)
	}
}
func (v *StarFilterVisitor) VisitGoalkeeper(g *Goalkeeper) {
	if g.Rating >= v.Threshold {
		v.Found = append(v.Found, g)
	}
}
func (v *StarFilterVisitor) VisitCoach(c *Coach) {}

type JsonExportVisitor struct {
	result string
}

func (v *JsonExportVisitor) VisitClub(c *Club) {
	v.result += fmt.Sprintf(`{"club":"%s","teams":[`, c.Name)
}
func (v *JsonExportVisitor) VisitTeam(t *Team) {
	if v.result[len(v.result)-1] != '[' {
		v.result += ","
	}
	v.result += fmt.Sprintf(`{"team":"%s","members":[`, t.Name)
}
func (v *JsonExportVisitor) VisitForward(f *Forward) {
	v.appendLeaf(fmt.Sprintf(`{"type":"Forward","name":"%s","salary":%.0f,"rating":%.1f}`, f.Name, f.Salary, f.Rating))
}
func (v *JsonExportVisitor) VisitDefender(d *Defender) {
	v.appendLeaf(fmt.Sprintf(`{"type":"Defender","name":"%s","salary":%.0f,"rating":%.1f}`, d.Name, d.Salary, d.Rating))
}
func (v *JsonExportVisitor) VisitGoalkeeper(g *Goalkeeper) {
	v.appendLeaf(fmt.Sprintf(`{"type":"Goalkeeper","name":"%s","salary":%.0f,"rating":%.1f}`, g.Name, g.Salary, g.Rating))
}
func (v *JsonExportVisitor) VisitCoach(c *Coach) {
	v.appendLeaf(fmt.Sprintf(`{"type":"Coach","name":"%s","salary":%.0f}`, c.Name, c.Salary))
}

func (v *JsonExportVisitor) appendLeaf(jsonStr string) {
	if v.result[len(v.result)-1] != '[' {
		v.result += ","
	}
	v.result += jsonStr
}

func (v *JsonExportVisitor) Finalize() string {
	return v.result + "]}]}"
}
