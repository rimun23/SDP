package factorymethod

type espressoBarista struct{ BaseCreator }

func NewEspressoBarista() Creator {
	barista := &espressoBarista{}
	barista.factory = barista
	return barista
}
func (*espressoBarista) CreateDrink() Drink { return espresso{} }

type latteBarista struct{ BaseCreator }

func NewLatteBarista() Creator {
	barista := &latteBarista{}
	barista.factory = barista
	return barista
}
func (*latteBarista) CreateDrink() Drink { return latte{} }

type cappuccinoBarista struct{ BaseCreator }

func NewCappuccinoBarista() Creator {
	barista := &cappuccinoBarista{}
	barista.factory = barista
	return barista
}
func (*cappuccinoBarista) CreateDrink() Drink { return cappuccino{} }
