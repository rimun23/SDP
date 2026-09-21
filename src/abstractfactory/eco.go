package abstractfactory

import "fmt"

const ecoBrand = "Green Leaf"

type ecoCup struct{}

func (ecoCup) Brand() string      { return ecoBrand }
func (ecoCup) RimDiameterMM() int { return 90 }
func (ecoCup) Describe() string   { return "kraft cup with a compostable plant-based lining" }

type ecoLid struct{}

func (ecoLid) Brand() string    { return ecoBrand }
func (ecoLid) DiameterMM() int  { return 90 }
func (ecoLid) Describe() string { return "moulded sugarcane-pulp lid" }

type ecoReceipt struct{}

func (ecoReceipt) Brand() string { return ecoBrand }

func (ecoReceipt) Render(customer, item string, price float64) string {
	return fmt.Sprintf("<<< %s >>>\n"+
		"Customer: %s\n"+
		"Item:     %s\n"+
		"Total:    %.2f\n"+
		"Printed on recycled paper. Please reuse your cup!\n",
		ecoBrand, customer, item, price)
}

type ecoFactory struct{}

func NewEcoFactory() ServingSetFactory { return ecoFactory{} }

func (ecoFactory) CreateCup() Cup         { return ecoCup{} }
func (ecoFactory) CreateLid() Lid         { return ecoLid{} }
func (ecoFactory) CreateReceipt() Receipt { return ecoReceipt{} }
