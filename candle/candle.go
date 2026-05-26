package candle

import (
	"fmt"
	"time"
)

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func NewEmptyCandle() *Candle {
	return &Candle{}
}

func (c Candle) Printf(extraFormat string, args ...any) {
	extra := ""

	if extraFormat != "" {
		extra = " " + fmt.Sprintf(extraFormat, args...)
	}

	fmt.Printf(
		"T: %s O: %f H: %f L: %f C: %f V: %d %s\n",
		c.Time.Format(time.DateTime),
		c.Open,
		c.High,
		c.Low,
		c.Close,
		c.Volume,
		extra,
	)
}
