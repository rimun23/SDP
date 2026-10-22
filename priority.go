package main

type PriorityStrategy interface {
	GetScore(s Student) int
}

type SeniorityPriority struct{}

func (sp *SeniorityPriority) GetScore(s Student) int { return s.Year * 100 }

type MeritPriority struct{}

func (mp *MeritPriority) GetScore(s Student) int { return int(s.GPA * 100) }

type GraduationPriority struct{}

func (gp *GraduationPriority) GetScore(s Student) int {
	if s.Year == 4 {
		return 999
	}
	return 10
}

type PriorityContext struct{ strategy PriorityStrategy }

func NewPriorityContext(s PriorityStrategy) *PriorityContext { return &PriorityContext{strategy: s} }
func (c *PriorityContext) SetStrategy(s PriorityStrategy)    { c.strategy = s }
func (c *PriorityContext) Execute(s Student) int             { return c.strategy.GetScore(s) }

type PrioritySelector struct{}

func (ps *PrioritySelector) Select(phase string) PriorityStrategy {
	switch phase {
	case "GRADUATION":
		return &GraduationPriority{}
	case "HONORS":
		return &MeritPriority{}
	default:
		return &SeniorityPriority{}
	}
}
