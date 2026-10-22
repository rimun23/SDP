package main

import "math"

type TuitionStrategy interface {
	Calculate(s Student, credits int) float64
}

type StandardTuition struct{}

func (st *StandardTuition) Calculate(s Student, credits int) float64 {
	return float64(credits) * 100.0
}

type ScholarshipTuition struct{}

func (sc *ScholarshipTuition) Calculate(s Student, credits int) float64 {
	return math.Max(50.0, (float64(credits)*100.0)-(s.GPA*60.0))
}

type InternationalTuition struct{ rates ExchangeRateService }

func (it *InternationalTuition) Calculate(s Student, credits int) float64 {
	return (float64(credits) * 120.0 * it.rates.GetRate()) + 1000.0
}

type TuitionContext struct{ strategy TuitionStrategy }

func NewTuitionContext(s TuitionStrategy) *TuitionContext { return &TuitionContext{strategy: s} }
func (c *TuitionContext) SetStrategy(s TuitionStrategy)   { c.strategy = s }
func (c *TuitionContext) Execute(s Student, credits int) float64 {
	return c.strategy.Calculate(s, credits)
}

type TuitionSelector struct{ rates ExchangeRateService }

func (ts *TuitionSelector) Select(s Student) TuitionStrategy {
	if s.IsInternational {
		return &InternationalTuition{rates: ts.rates}
	}
	if s.HasScholarship && s.GPA >= 3.5 {
		return &ScholarshipTuition{}
	}
	return &StandardTuition{}
}

type LegacyTuitionCalculator struct{ rates ExchangeRateService }

func (l *LegacyTuitionCalculator) Calculate(s Student, credits int, mode string) float64 {
	if mode == "INT" {
		return (float64(credits) * 120.0 * l.rates.GetRate()) + 1000.0
	} else if mode == "SCHOLAR" {
		return math.Max(50.0, (float64(credits)*100.0)-(s.GPA*60.0))
	}
	return float64(credits) * 100.0
}
