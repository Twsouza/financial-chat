package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	STOCK_API_ENDPOINT = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

type Stock struct {
	Symbol string
	Date   string
	Time   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

const (
	STOCK_BOT_NAME = "stock-bot"
)

func (s *Stock) Quote() string {
	return fmt.Sprintf("%s quote is $%.2f per share", s.Symbol, s.Close)
}

func HandleStockCommand(symbol string) (string, string, error) {
	res, err := httpGet(fmt.Sprintf(STOCK_API_ENDPOINT, symbol))
	if err != nil {
		fmt.Printf("error getting stock quote: %v\n", err)
		return STOCK_BOT_NAME, "", fmt.Errorf("error getting stock %s quote, please check the stock quote and try again", symbol)
	}
	defer res.Body.Close()

	stock, err := parseStockResponse(res.Body)
	if err != nil {
		fmt.Printf("error getting stock quote: %v\n", err)
		return STOCK_BOT_NAME, "", fmt.Errorf("error getting stock %s quote, please check the stock quote and try again", symbol)
	}

	return STOCK_BOT_NAME, stock.Quote(), nil
}

func parseStockResponse(body io.Reader) (*Stock, error) {
	records, err := csv.NewReader(body).ReadAll()
	if err != nil {
		return nil, err
	}

	open, err := strconv.ParseFloat(records[1][3], 64)
	if err != nil {
		return nil, err
	}

	high, err := strconv.ParseFloat(records[1][4], 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(records[1][5], 64)
	if err != nil {
		return nil, err
	}

	close, err := strconv.ParseFloat(records[1][6], 64)
	if err != nil {
		return nil, err
	}

	volume, err := strconv.ParseInt(records[1][7], 10, 64)
	if err != nil {
		return nil, err
	}

	stock := &Stock{
		Symbol: records[1][0],
		Date:   records[1][1],
		Time:   records[1][2],
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: volume,
	}

	return stock, nil
}

func httpGet(url string) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}
