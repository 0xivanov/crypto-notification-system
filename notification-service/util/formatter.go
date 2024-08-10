package util

import (
	"fmt"

	"github.com/0xivanov/crypto-notification-system/common/model"
)

func FormatMessage(tickerData model.TickerData) string {
	return fmt.Sprintf("Ticker update for %s:\nLast Price: %f\nBid: %f\nAsk: %f\nChange: %f%%",
		tickerData.Symbol, tickerData.Last, tickerData.Bid, tickerData.Ask, tickerData.ChangePct)
}
