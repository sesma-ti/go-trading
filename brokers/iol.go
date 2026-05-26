package brokers

import (
	"fmt"

	"github.com/sesma-ti/go-trading/orders"
)

type IOLBroker struct {
	Capital float64
	Stocks  float64
}

func NewIOLBroker() *IOLBroker {
	return &IOLBroker{
		Capital: 10000,
		Stocks:  0,
	}
}

func (iol *IOLBroker) Send(order orders.Order) {
	// fmt.Printf("IOL - Sending Order | Ticker=%s Side=%s Amount=%f Price=%f\n", order.Ticker, order.Side, order.Amount, order.Price)

	if order.Side == orders.Buy {
		total := order.Amount * order.Price
		iol.Capital -= total
		iol.Stocks += order.Amount
	} else if order.Side == orders.Sell {
		total := order.Amount * order.Price
		iol.Capital += total
		iol.Stocks -= order.Amount
	}

	fmt.Printf("IOL - Account | Capital=%f Stocks=%f\n", iol.Capital, iol.Stocks)
}
