package abstractfactory

type Cup interface {
	Brand() string
	Describe() string
	RimDiameterMM() int
}
type Lid interface {
	Brand() string
	Describe() string
	DiameterMM() int
}
type Receipt interface {
	Brand() string
	Render(customer, item string, price float64) string
}
type ServingSetFactory interface {
	CreateCup() Cup
	CreateLid() Lid
	CreateReceipt() Receipt
}
