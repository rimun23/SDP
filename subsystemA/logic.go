package subsystemA

import (
	"errors"
	"fmt"
)

var (
	ErrPrerequisiteNotMet = errors.New("subsystemA: student does not meet course prerequisites")
	ErrFinancialHoldFound = errors.New("subsystemA: student has an unresolved financial hold")
)

type PrerequisiteChecker struct {
	studentHistory map[string][]string
}

func NewPrerequisiteChecker() *PrerequisiteChecker {
	return &PrerequisiteChecker{
		studentHistory: map[string][]string{
			"STUDENT_OK":   {"CS101", "MATH101"},
			"STUDENT_FAIL": {},
		},
	}
}
func (p *PrerequisiteChecker) ValidateEligibility(studentID, courseID string) error {
	if courseID == "CS201" {
		courses := p.studentHistory[studentID]
		found := false
		for _, c := range courses {
			if c == "CS101" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%w: CS201 requires CS101", ErrPrerequisiteNotMet)
		}
	}
	return nil
}

type FinancialHoldChecker struct {
	studentsWithHold map[string]bool
}

func NewFinancialHoldChecker() *FinancialHoldChecker {
	return &FinancialHoldChecker{
		studentsWithHold: map[string]bool{
			"STUDENT_DEBT": true,
		},
	}
}
func (f *FinancialHoldChecker) CheckHold(studentID string) error {
	if f.studentsWithHold[studentID] {
		return fmt.Errorf("%w for student %s", ErrFinancialHoldFound, studentID)
	}
	return nil
}
