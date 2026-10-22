package main

import "testing"

func TestStrategies(t *testing.T) {
	fake := &FakeExchangeRateService{Rate: 1.0}
	s := Student{Name: "Bob", GPA: 3.8, Year: 4, HasScholarship: true, IsInternational: true}

	if val := (&StandardTuition{}).Calculate(s, 3); val != 300.0 {
		t.Errorf("expected 300, got %f", val)
	}
	if val := (&ScholarshipTuition{}).Calculate(s, 3); val != 72.0 {
		t.Errorf("expected 72, got %f", val)
	}
	if val := (&InternationalTuition{rates: fake}).Calculate(s, 3); val != 1360.0 {
		t.Errorf("expected 1360, got %f", val)
	}

	if (&SeniorityPriority{}).GetScore(s) != 400 {
		t.Error("invalid seniority score")
	}
	if (&MeritPriority{}).GetScore(s) != 380 {
		t.Error("invalid merit score")
	}
	if (&GraduationPriority{}).GetScore(s) != 999 {
		t.Error("invalid graduation score")
	}
}
func TestSelectors(t *testing.T) {
	tSel := &TuitionSelector{rates: &FakeExchangeRateService{Rate: 1.0}}

	if _, ok := tSel.Select(Student{IsInternational: true}).(*InternationalTuition); !ok {
		t.Error("selector failed for international student")
	}

	scholar := Student{HasScholarship: true, GPA: 3.8, IsInternational: false}
	if _, ok := tSel.Select(scholar).(*ScholarshipTuition); !ok {
		t.Error("selector failed for scholarship student")
	}

	regular := Student{HasScholarship: false, GPA: 2.0, IsInternational: false}
	if _, ok := tSel.Select(regular).(*StandardTuition); !ok {
		t.Error("selector failed for standard student")
	}

	pSel := &PrioritySelector{}
	if _, ok := pSel.Select("GRADUATION").(*GraduationPriority); !ok {
		t.Error("priority selector failed for graduation phase")
	}
	if _, ok := pSel.Select("HONORS").(*MeritPriority); !ok {
		t.Error("priority selector failed for honors phase")
	}
	if _, ok := pSel.Select("REGULAR").(*SeniorityPriority); !ok {
		t.Error("priority selector failed for regular phase")
	}
}

func TestAllTuitionStrategiesPositiveFee(t *testing.T) {
	fake := &FakeExchangeRateService{Rate: 1.0}
	s := Student{Name: "Test", GPA: 4.0, Year: 2}

	tests := []struct {
		name     string
		strategy TuitionStrategy
	}{
		{"Standard", &StandardTuition{}},
		{"Scholarship", &ScholarshipTuition{}},
		{"International", &InternationalTuition{rates: fake}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if fee := tt.strategy.Calculate(s, 5); fee <= 0 {
				t.Errorf("%s failed: fee must be positive, got %f", tt.name, fee)
			}
		})
	}
}
