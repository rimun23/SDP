package main

import (
	"SDP/src/abstractfactory"
	"SDP/src/factorymethod"
	"fmt"
)

var baristas = map[string]func() factorymethod.Creator{
	"espresso":   factorymethod.NewEspressoBarista,
	"latte":      factorymethod.NewLatteBarista,
	"cappuccino": factorymethod.NewCappuccinoBarista,
}
var servingSets = map[string]func() abstractfactory.ServingSetFactory{
	"classic": abstractfactory.NewClassicFactory,
	"eco":     abstractfactory.NewEcoFactory,
	"premium": abstractfactory.NewPremiumFactory,
}

func main() {
	drinkName, brandName, customer := "latte", "classic", "Alex"
	newBarista, ok := baristas[drinkName]
	if !ok {
		fmt.Println("unknown drink:", drinkName, "(try espresso, latte or cappuccino)")
		return
	}
	newFactory, ok := servingSets[brandName]
	if !ok {
		fmt.Println("unknown brand:", brandName, "(try classic, eco or premium)")
		return
	}
	order := newBarista().ServeOrder(customer)
	fmt.Println("Factory Method:")
	fmt.Print(order.Ticket)
	counter := abstractfactory.NewCounter(newFactory())
	packed, err := counter.Serve(order.Customer, order.Drink, order.Price)
	if err != nil {
		fmt.Println("cannot pack the order:", err)
		return
	}
	fmt.Println()
	fmt.Println("Abstract Factory:")
	fmt.Print(packed)
}
