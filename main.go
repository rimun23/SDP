package main

import "fmt"

func main() {
	rates := &SimulatedExchangeRateService{}
	alice := Student{Name: "Alice", GPA: 3.6, Year: 4, IsInternational: false, HasScholarship: true}
	credits := 5
	fmt.Println("1. Single Input Processed by All Tuition Strategies")
	strats := []TuitionStrategy{
		&StandardTuition{},
		&ScholarshipTuition{},
		&InternationalTuition{rates: rates},
	}
	for _, s := range strats {
		fmt.Printf("%-28T -> %.2f\n", s, s.Calculate(alice, credits))
	}
	fmt.Println("\n2. Manual Selection and Runtime Strategy Switching")
	ctx := NewTuitionContext(&StandardTuition{})
	fmt.Printf("Initial (Standard): %.2f\n", ctx.Execute(alice, credits))

	ctx.SetStrategy(&ScholarshipTuition{})
	fmt.Printf("After SetStrategy (Scholarship): %.2f\n", ctx.Execute(alice, credits))

	fmt.Println("\n3. Automatic Strategy Selection via Selector")
	selector := &TuitionSelector{rates: rates}
	selectedStrategy := selector.Select(alice)
	ctx.SetStrategy(selectedStrategy)
	fmt.Printf("Selector assigned %T -> Total: %.2f\n", selectedStrategy, ctx.Execute(alice, credits))

	fmt.Println("\n4. Strategy Family 2: Course Enrollment Priority")
	pSel := &PrioritySelector{}
	pCtx := NewPriorityContext(pSel.Select("REGULAR"))
	fmt.Printf("Regular Registration Phase: %d points\n", pCtx.Execute(alice))

	pCtx.SetStrategy(pSel.Select("GRADUATION"))
	fmt.Printf("Graduation Deadline Phase:  %d points\n", pCtx.Execute(alice))
}
