package abstractfactory

import (
	"fmt"
	"strings"
)

type Counter struct {
	factory ServingSetFactory
}

func NewCounter(factory ServingSetFactory) *Counter {
	return &Counter{factory: factory}
}

func (c *Counter) Serve(customer, item string, price float64) (string, error) {
	cup := c.factory.CreateCup()
	lid := c.factory.CreateLid()
	receipt := c.factory.CreateReceipt()

	if lid.DiameterMM() != cup.RimDiameterMM() {
		return "", fmt.Errorf("lid (%d mm) does not fit the cup (%d mm)",
			lid.DiameterMM(), cup.RimDiameterMM())
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Cup:  %s\n", cup.Describe())
	fmt.Fprintf(&b, "Lid:  %s (%d mm, fits the cup)\n", lid.Describe(), lid.DiameterMM())
	b.WriteString(receipt.Render(customer, item, price))
	return b.String(), nil
}
