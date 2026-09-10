package main

import (
	"fmt"
	"strings"
)

type TravelPackageText struct {
	sb        strings.Builder
	apartment string
	flight    string
}

func newTravelPackageText() *TravelPackageText {
	b := &TravelPackageText{}
	b.sb.WriteString("This travel package included \n")
	return b
}
func (b *TravelPackageText) setApartment(apartment string) TravelPackageBuilder {
	b.apartment = apartment
	b.sb.WriteString(fmt.Sprintf("A place where you will live is %s \n", apartment))
	return b
}
func (b *TravelPackageText) setFlight(flight string) TravelPackageBuilder {
	b.flight = flight
	b.sb.WriteString(fmt.Sprintf("Ticket for your flight will be: %s \n", flight))
	return b
}
func (b *TravelPackageText) setInsurance(insurance string) TravelPackageBuilder {
	b.sb.WriteString(fmt.Sprintf("Insurance for this trip: %s \n", insurance))
	return b
}
func (b *TravelPackageText) getResult() string {
	return b.sb.String()
}
