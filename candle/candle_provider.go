package candle

import "iter"

type CandleProvider interface {
	Candles(symbol string, lookbackDays int) iter.Seq[*Candle]
}
