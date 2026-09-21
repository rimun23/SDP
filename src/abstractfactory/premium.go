package abstractfactory

import "fmt"

const premiumBrand = "Midnight Reserve"

type premiumCup struct{}

func (premiumCup) Brand() string      { return premiumBrand }
func (premiumCup) RimDiameterMM() int { return 73 }
func (premiumCup) Describe() string   { return "matte black double-wall cup with a gold logo" }

type premiumLid struct{}

func (premiumLid) Brand() string    { return premiumBrand }
func (premiumLid) DiameterMM() int  { return 73 }
func (premiumLid) Describe() string { return "matte black lid with a silicone seal" }

type premiumReceipt struct{}

func (premiumReceipt) Brand() string { return premiumBrand }

func (premiumReceipt) Render(customer, item string, price float64) string {
	return fmt.Sprintf("/*** %s ***/\n"+
		"| Customer: %s\n"+
		"| Item:     %s\n"+
		"| Total:    %.2f\n"+
		"/*************************/\n"+
		"Enjoy your reserve selection.\n",
		premiumBrand, customer, item, price)
}

type premiumFactory struct{}

func NewPremiumFactory() ServingSetFactory { return premiumFactory{} }

func (premiumFactory) CreateCup() Cup         { return premiumCup{} }
func (premiumFactory) CreateLid() Lid         { return premiumLid{} }
func (premiumFactory) CreateReceipt() Receipt { return premiumReceipt{} }
