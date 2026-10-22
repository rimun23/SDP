package main

type Student struct {
	Name            string
	GPA             float64
	Year            int
	IsInternational bool
	HasScholarship  bool
}

type ExchangeRateService interface {
	GetRate() float64
}

type SimulatedExchangeRateService struct{}

func (s *SimulatedExchangeRateService) GetRate() float64 { return 500.0 }

type FakeExchangeRateService struct{ Rate float64 }

func (f *FakeExchangeRateService) GetRate() float64 { return f.Rate }
