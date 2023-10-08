package core_test

import (
	"bot/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStockQuote(t *testing.T) {
	s := core.Stock{
		Symbol: "AAPL.US",
		Date:   "2023-10-05",
		Time:   "20:25:47",
		Open:   100.00,
		High:   110.00,
		Low:    90.00,
		Close:  93.42,
		Volume: 25899606,
	}

	expected := "AAPL.US quote is $93.42 per share"
	require.Equal(t, expected, s.Quote())
}
