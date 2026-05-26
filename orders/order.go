package orders

type Side string

const (
	Buy  Side = "COMPRA"
	Sell Side = "VENTA"
)

type Order struct {
	Side   Side
	Ticker string
	Amount float64
	Price  float64
}
