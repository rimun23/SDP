package main

import "errors"

var ErrNotFullTravelPackage = errors.New("apartment and flight are requirment")

type TravelPackageObject struct {
	apartment string
	flight    string
	insurance string
}

func TravelPackageObjectBuilder() *TravelPackageObject {
	return &TravelPackageObject{}
}
func (t *TravelPackageObject) setApartment(apartment string) TravelPackageBuilder {
	t.apartment = apartment
	return t
}
func (t *TravelPackageObject) setFlight(flight string) TravelPackageBuilder {
	t.flight = flight
	return t
}
func (t *TravelPackageObject) setInsurance(insurance string) TravelPackageBuilder {
	t.insurance = insurance
	return t
}
func (t *TravelPackageObject) getResult() (*TravelPackage, error) {
	if t.apartment == "" || t.flight == "" {
		return nil, ErrNotFullTravelPackage
	}
	return newTravelPackage(t.apartment, t.flight, t.insurance), nil
}
