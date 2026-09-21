package factorymethod

type Drink interface {
	Name() string
	VolumeML() int
	Price() float64
	BrewSeconds() int
}
type espresso struct{}

func (espresso) Name() string     { return "Espresso" }
func (espresso) VolumeML() int    { return 30 }
func (espresso) Price() float64   { return 2.50 }
func (espresso) BrewSeconds() int { return 28 }

type latte struct{}

func (latte) Name() string     { return "Latte" }
func (latte) VolumeML() int    { return 250 }
func (latte) Price() float64   { return 3.80 }
func (latte) BrewSeconds() int { return 45 }

type cappuccino struct{}

func (cappuccino) Name() string     { return "Cappuccino" }
func (cappuccino) VolumeML() int    { return 180 }
func (cappuccino) Price() float64   { return 3.50 }
func (cappuccino) BrewSeconds() int { return 50 }
