package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"time"

	"github.com/sesma-ti/go-trading/candle"
)

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol string `json:"symbol"`
			} `json:"meta"`
			Timestamps []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

type YahooFinance struct{}

func NewYahooFinance() *YahooFinance {
	return &YahooFinance{}
}

func (*YahooFinance) Candles(symbol string, lookbackDays int) iter.Seq[*candle.Candle] {
	if lookbackDays < 200 {
		lookbackDays = 200
	}

	end := time.Now()
	start := end.AddDate(0, 0, -lookbackDays)

	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&period1=%d&period2=%d",
		symbol,
		start.Unix(),
		end.Unix(),
	)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("error creando request: %w", err))
	}
	// Yahoo requiere un User-Agent razonable
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; trading-bot/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("error descargando datos de Yahoo Finance: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("error leyendo respuesta: %w", err))
	}

	var yr yahooResponse
	if err := json.Unmarshal(body, &yr); err != nil {
		panic(fmt.Errorf("error parseando JSON: %w", err))
	}

	if yr.Chart.Error != nil {
		panic(fmt.Errorf("Yahoo Finance error [%s]: %s", yr.Chart.Error.Code, yr.Chart.Error.Description))
	}

	if len(yr.Chart.Result) == 0 {
		panic(fmt.Errorf("no se encontraron datos para el símbolo '%s'", symbol))
	}

	result := yr.Chart.Result[0]
	quotes := result.Indicators.Quote[0]
	timestamps := result.Timestamps

	if len(timestamps) == 0 {
		panic(fmt.Errorf("el símbolo '%s' no devolvió barras OHLCV", symbol))
	}

	return func(yield func(*candle.Candle) bool) {
		for i, ts := range timestamps {
			if i >= len(quotes.Close) || quotes.Close[i] == 0 {
				continue
			}

			c := &candle.Candle{
				Time:   time.Unix(ts, 0).UTC(),
				Open:   quotes.Open[i],
				High:   quotes.High[i],
				Low:    quotes.Low[i],
				Close:  quotes.Close[i],
				Volume: quotes.Volume[i],
			}

			if !yield(c) {
				return
			}
		}
	}
}
