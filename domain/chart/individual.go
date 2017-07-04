package chart

import (
	"github.com/jchavannes/money/data/db"
	"github.com/jchavannes/money/object/price"
	"github.com/jchavannes/jgo/jerr"
	"time"
	"github.com/jchavannes/jgo/jtime"
	"sort"
)

func GetIndividualChartItems(investment *db.Investment, investmentTransactions []*db.InvestmentTransaction) ([]*ChartItem, error) {
	costChartItem := &ChartItem{
		Name: investment.Symbol + " - Cost",
		ChartDataPoints: []*ChartDataPoint{},
	}
	valueChartItem := &ChartItem{
		Name: investment.Symbol + " - Value",
		ChartDataPoints: []*ChartDataPoint{},
	}
	history, err := price.GetHistory(investment)
	if err != nil {
		return nil, jerr.Get("Error getting history", err)
	}
	for _, investmentPrice := range history {
		investmentPriceTimestamp := time.Unix(jtime.RoundTimeToDay(investmentPrice.Timestamp), 0)
		var costDataPointAmount float32
		var valueDataPointAmount float32
		for _, investmentTransaction := range investmentTransactions {
			if investmentTransaction.InvestmentId == investment.Id && investmentTransaction.Date.Unix() < investmentPrice.Timestamp {
				value := investmentPrice.Price * investmentTransaction.Quantity
				cost := investmentTransaction.Price * investmentTransaction.Quantity
				if investmentTransaction.Type == uint(db.InvestmentTransactionType_Buy) {
					costDataPointAmount += cost
					valueDataPointAmount += value
				} else {
					costDataPointAmount -= cost
					valueDataPointAmount -= value
				}
			}
		}
		if valueDataPointAmount <= 0.0 {
			continue
		}
		costDataPoint, err := costChartItem.GetDataPoint(investmentPriceTimestamp)
		if err != nil {
			costDataPoint = &ChartDataPoint{
				Timestamp: investmentPriceTimestamp,
				Amount: costDataPointAmount,
			}
			costChartItem.ChartDataPoints = append(costChartItem.ChartDataPoints, costDataPoint)
		}
		valueDataPoint, err := valueChartItem.GetDataPoint(investmentPriceTimestamp)
		if err != nil {
			valueDataPoint = &ChartDataPoint{
				Timestamp: investmentPriceTimestamp,
				Amount: valueDataPointAmount,
			}
			valueChartItem.ChartDataPoints = append(valueChartItem.ChartDataPoints, valueDataPoint)
		}
	}
	fillInChartItem(costChartItem)
	fillInChartItem(valueChartItem)
	return []*ChartItem{
		costChartItem,
		valueChartItem,
	}, nil
}

func fillInChartItem(chartItem *ChartItem) {
	if chartItem == nil {
		return
	}
	var lastTimeStamp int64
	var lastAmount float32
	// Two scenarios of gaps (e.g. weekend):
	// 1: Middle of data
	// 2: End of data
	additionalDataPoints := []*ChartDataPoint{}
	for _, dataPoint := range chartItem.ChartDataPoints {
		nextTimeStamp := lastTimeStamp + jtime.SecondsInDay
		maxLoopsLeft := 3
		for lastTimeStamp > 0 && nextTimeStamp < dataPoint.Timestamp.Unix() && maxLoopsLeft > 0 {
			additionalDataPoints = append(additionalDataPoints, &ChartDataPoint{
				Timestamp: time.Unix(nextTimeStamp, 0),
				Amount: lastAmount,
			})
			maxLoopsLeft--
			nextTimeStamp += jtime.SecondsInDay
		}
		lastTimeStamp = dataPoint.Timestamp.Unix()
		lastAmount = dataPoint.Amount
	}
	today := time.Now().Unix()
	nextTimeStamp := lastTimeStamp + jtime.SecondsInDay
	maxLoopsLeft := 3
	for lastTimeStamp > 0 && nextTimeStamp <= today && maxLoopsLeft > 0 {
		additionalDataPoints = append(additionalDataPoints, &ChartDataPoint{
			Timestamp: time.Unix(nextTimeStamp, 0),
			Amount: lastAmount,
		})
		maxLoopsLeft--
		nextTimeStamp += jtime.SecondsInDay
	}
	chartItem.ChartDataPoints = append(chartItem.ChartDataPoints, additionalDataPoints...)
	sort.Sort(ChartDataPointSorter(chartItem.ChartDataPoints))
}
