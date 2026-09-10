package main

import (
	"fmt"
	"log"
)

func main() {
	director := newTravelPackageDirector()
	objectBuilder := TravelPackageObjectBuilder()
	director.makeBudgetTravelPackage(objectBuilder)
	budgetTravelPackage, err := objectBuilder.getResult()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(budgetTravelPackage)
	fmt.Println("\n")
	textBuilder := newTravelPackageText()
	director.makeRichTravelPackage(textBuilder)
	richTravelPackage := textBuilder.getResult()
	fmt.Println(richTravelPackage)
}
