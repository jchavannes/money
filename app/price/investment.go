package price

import (
	"github.com/jchavannes/money/app/db"
	"time"
)

func getLastInvestmentTimestamp(investment *db.Investment) int64 {
	var lastItemTimestamp int64

	lastItem, err := db.GetLastInvestmentPrice(investment)
	if err != nil {
		lastItemTimestamp = 0
	} else {
		lastItemTimestamp = lastItem.Timestamp
	}

	return lastItemTimestamp
}

func isTimestampUpdatedRecently(timestamp int64) bool {
	return timestamp > time.Now().AddDate(0, 0, -1).Unix()
}
