package factorymethod

import (
	"fmt"
	"strings"
)
type Order struct {
	Customer string
	Drink    string
	Price    float64
	Ticket   string
}

type Creator interface {
	CreateDrink() Drink
	ServeOrder(customer string) Order
}

type drinkFactory interface {
	CreateDrink() Drink
}

type BaseCreator struct {
	factory drinkFactory
}

func (c *BaseCreator) ServeOrder(customer string) Order {
	drink := c.factory.CreateDrink()
	return Order{
		Customer: customer,
		Drink:    drink.Name(),
		Price:    drink.Price(),
		Ticket:   formatTicket(customer, drink),
	}
}

func formatTicket(customer string, drink Drink) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Ticket for %s: %s (%d ml), %.2f\n",
		customer, drink.Name(), drink.VolumeML(), drink.Price())
	fmt.Fprintf(&b, "Ready in about %d seconds\n", drink.BrewSeconds())
	return b.String()
}
