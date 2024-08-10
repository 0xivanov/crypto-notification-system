package util

import (
	"testing"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/stretchr/testify/assert"
)

func TestFormatMessage(t *testing.T) {
	tickerData := model.TickerData{
		Symbol:    "BTC-USD",
		Last:      30000.50,
		Bid:       29950.00,
		Ask:       30010.00,
		ChangePct: 2.5,
	}

	expectedMessage := "Ticker update for BTC-USD:\nLast Price: 30000.500000\nBid: 29950.000000\nAsk: 30010.000000\nChange: 2.500000%"

	actualMessage := FormatMessage(tickerData)

	assert.Equal(t, expectedMessage, actualMessage)
}

func TestFormatMessageNegativeChange(t *testing.T) {
	tickerData := model.TickerData{
		Symbol:    "ETH-USD",
		Last:      2000.00,
		Bid:       1995.00,
		Ask:       2005.00,
		ChangePct: -1.25,
	}

	expectedMessage := "Ticker update for ETH-USD:\nLast Price: 2000.000000\nBid: 1995.000000\nAsk: 2005.000000\nChange: -1.250000%"

	actualMessage := FormatMessage(tickerData)

	assert.Equal(t, expectedMessage, actualMessage)
}

func TestFormatMessageZeroChange(t *testing.T) {
	tickerData := model.TickerData{
		Symbol:    "XRP-USD",
		Last:      0.50,
		Bid:       0.49,
		Ask:       0.51,
		ChangePct: 0.0,
	}

	expectedMessage := "Ticker update for XRP-USD:\nLast Price: 0.500000\nBid: 0.490000\nAsk: 0.510000\nChange: 0.000000%"

	actualMessage := FormatMessage(tickerData)

	assert.Equal(t, expectedMessage, actualMessage)
}

func TestFormatMessageLargeNumbers(t *testing.T) {
	tickerData := model.TickerData{
		Symbol:    "LTC-USD",
		Last:      123456789.123456,
		Bid:       123456788.123456,
		Ask:       123456790.123456,
		ChangePct: 99.999999,
	}

	expectedMessage := "Ticker update for LTC-USD:\nLast Price: 123456789.123456\nBid: 123456788.123456\nAsk: 123456790.123456\nChange: 99.999999%"

	actualMessage := FormatMessage(tickerData)

	assert.Equal(t, expectedMessage, actualMessage)
}
