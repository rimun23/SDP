package abstractfactory

import "fmt"

const classicBrand = "Urban Roast"

type classicCup struct{}

func (classicCup) Brand() string      { return classicBrand }
func (classicCup) RimDiameterMM() int { return 80 }
func (classicCup) Describe() string   { return "white paper cup with a red stripe" }

type classicLid struct{}

func (classicLid) Brand() string    { return classicBrand }
func (classicLid) DiameterMM() int  { return 80 }
func (classicLid) Describe() string { return "black plastic sip lid" }

type classicReceipt struct{}

func (classicReceipt) Brand() string { return classicBrand }

func (classicReceipt) Render(customer, item string, price float64) string {
	return fmt.Sprintf("--- %s ---\n"+
		"Customer: %s\n"+
		"Item:     %s\n"+
		"Total:    %.2f\n"+
		"Thank you, see you again!\n",
		classicBrand, customer, item, price)
}

type classicFactory struct{}

func NewClassicFactory() ServingSetFactory { return classicFactory{} }

func (classicFactory) CreateCup() Cup         { return classicCup{} }
func (classicFactory) CreateLid() Lid         { return classicLid{} }
func (classicFactory) CreateReceipt() Receipt { return classicReceipt{} }
