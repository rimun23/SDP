package main

import "fmt"

type TravelPackage struct {
	apartment string
	flight    string
	insurance string
}

func newTravelPackage(apartment, flight, insurance string) *TravelPackage {
	return &TravelPackage{apartment: apartment, flight: flight, insurance: insurance}
}
func (t *TravelPackage) Apartment() string { return t.apartment }
func (t *TravelPackage) Flight() string    { return t.flight }
func (t *TravelPackage) Insurance() string { return t.insurance }
func (t *TravelPackage) String() string {
	return fmt.Sprintf("This travel package has %s, %s, %s", t.apartment, t.flight, t.insurance)
}
