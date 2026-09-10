package main

type TravelPackageBuilder interface {
	setApartment(apartment string) TravelPackageBuilder
	setFlight(flight string) TravelPackageBuilder
	setInsurance(insurance string) TravelPackageBuilder
}
