package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"assignment2-design-patterns/src/abstractfactory"
	"assignment2-design-patterns/src/factorymethod"
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
	drinkName := flag.String("drink", "espresso", "drink to brew: "+optionList(baristas))
	brandName := flag.String("brand", "classic", "serving-set brand: "+optionList(servingSets))
	customer := flag.String("customer", "Alex", "customer name")
	flag.Parse()

	newBarista, ok := baristas[*drinkName]
	if !ok {
		fail("drink", *drinkName, baristas)
	}

	order := newBarista().ServeOrder(*customer)
	fmt.Println("Factory Method:")
	fmt.Print(order.Ticket)
	newFactory, ok := servingSets[*brandName]
	if !ok {
		fail("brand", *brandName, servingSets)
	}

	counter := abstractfactory.NewCounter(newFactory())
	packed, err := counter.Serve(order.Customer, order.Drink, order.Price)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot pack the order:", err)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println("Abstract Factory:")
	fmt.Print(packed)
}

func optionList[T any](options map[string]T) string {
	names := make([]string, 0, len(options))
	for name := range options {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func fail[T any](kind, got string, options map[string]T) {
	fmt.Fprintf(os.Stderr, "unknown %s %q (available: %s)\n", kind, got, optionList(options))
	os.Exit(2)
}
