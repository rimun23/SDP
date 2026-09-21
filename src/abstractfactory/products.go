// Package abstractfactory is Part B of the assignment: the Abstract Factory
// pattern applied to a coffee kiosk.
//
// A "serving set" is a family of three products that must share one brand:
// Cup + Lid + Receipt. A lid only fits the cup of its own family, so mixing
// families must be impossible.
//
// Pattern roles:
//
//	AbstractProduct  -> Cup, Lid, Receipt                    (products.go)
//	AbstractFactory  -> ServingSetFactory                    (products.go)
//	ConcreteFactory  -> classicFactory, ecoFactory, premiumFactory
//	ConcreteProduct  -> classicCup, ecoLid, premiumReceipt, ...
//	Client           -> Counter                              (client.go)
//
// Concrete products and factories are unexported, so a client can only reach
// them through the New...Factory constructors and the interfaces below.
package abstractfactory

// Cup is the first abstract product of the family.
type Cup interface {
	Brand() string
	Describe() string
	RimDiameterMM() int
}

// Lid is the second abstract product. It must fit the Cup of the same family.
type Lid interface {
	Brand() string
	Describe() string
	DiameterMM() int
}

// Receipt is the third abstract product.
type Receipt interface {
	Brand() string
	Render(customer, item string, price float64) string
}

// ServingSetFactory is the AbstractFactory: exactly one create-method per
// product in the family.
type ServingSetFactory interface {
	CreateCup() Cup
	CreateLid() Lid
	CreateReceipt() Receipt
}
