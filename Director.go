package main

type TravelPackageDirector struct{}

func newTravelPackageDirector() *TravelPackageDirector {
	return &TravelPackageDirector{}
}
func (p *TravelPackageDirector) makeBudgetTravelPackage(t TravelPackageBuilder) {
	t.setApartment("2 star hotel").setFlight("economic class").setInsurance("econom")
}
func (p *TravelPackageDirector) makeRichTravelPackage(t TravelPackageBuilder) {
	t.setApartment("5 star hotel").setFlight("business class").setInsurance("max")
}
