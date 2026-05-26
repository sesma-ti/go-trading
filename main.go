package main

import (
	"github.com/sesma-ti/go-trading/brokers"
	"github.com/sesma-ti/go-trading/candle"
	"github.com/sesma-ti/go-trading/indicators"
	"github.com/sesma-ti/go-trading/orders"
	"github.com/sesma-ti/go-trading/providers"
)

func main() {
	var provider candle.CandleProvider = providers.NewYahooFinance()
	var broker orders.OrderBroker = brokers.NewIOLBroker()

	ma20 := indicators.NewMovingAverage(20)
	ma50 := indicators.NewMovingAverage(50)
	ma200 := indicators.NewMovingAverage(200)

	for c := range provider.Candles("SPY", 500) {
		ma20val := ma20.Add(c.Close)
		ma50val := ma50.Add(c.Close)
		ma200val := ma200.Add(c.Close)

		if c.Open < ma20val &&
			c.Close > c.Open &&
			c.Close > ma20val &&
			ma20val > ma50val &&
			ma50val > ma200val {
			broker.Send(orders.Order{
				Side:   "COMPRA",
				Ticker: "SPY",
				Amount: 1,
				Price:  c.Open,
			})
		}

		if c.Open > ma20val &&
			c.Close < ma20val &&
			ma20val < ma50val &&
			ma50val < ma200val {
			broker.Send(orders.Order{
				Side:   "VENTA",
				Ticker: "SPY",
				Amount: 1,
				Price:  c.Open,
			})
		}

		c.Printf("MA20 %f MA50 %f MA200 %f", ma20val, ma50val, ma200val)
	}
}
