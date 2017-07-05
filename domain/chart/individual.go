package chart

import (
	"github.com/jchavannes/money/data/db"
	"github.com/jchavannes/money/object/price"
	"github.com/jchavannes/jgo/jerr"
	"time"
	"github.com/jchavannes/jgo/jtime"
	"sort"
	"github.com/jchavannes/money/object/portfolio"
	"strings"
)

func GetIndividualChartData(userId uint, symbol string, market string) (*ChartDataOutput, error) {
	investmentTransactions, err := db.GetInvestmentTransactionsForUser(userId)
	if err != nil {
		return nil, jerr.Get("Error getting investment transactions", err)
	}
	userPortfolio, err := portfolio.Get(userId)
	if err != nil {
		return nil, jerr.Get("Error getting user portfolio", err)
	}

	var investment *db.Investment
	for _, portfolioItem := range userPortfolio.Items {
		if portfolioItem.Investment.Symbol == symbol && portfolioItem.Investment.InvestmentType == market {
			investment = &portfolioItem.Investment
		}
	}
	if investment == nil {
		return nil, jerr.New("Unable to find investment")
	}

	individualChartItems, err := GetIndividualChartItems(investment, investmentTransactions)
	if err != nil {
		return nil, jerr.Get("Error getting individual chart items", err)
	}
	if len(individualChartItems) != 2 {
		return nil, jerr.New("Unexpected number of chart items")
	}

	chartData := ChartData{
		Title: strings.ToUpper(investment.InvestmentType) + " - " + strings.ToUpper(investment.Symbol),
		ChartItems: []*ChartItem{
			individualChartItems[0],
			individualChartItems[1],
		},
	}
	return chartData.GetChartDataOutput(), nil
}

func GetIndividualChartItems(investment *db.Investment, investmentTransactions []*db.InvestmentTransaction) ([]*ChartItem, error) {
	costChartItem := &ChartItem{
		Name: strings.ToUpper(investment.Symbol) + " - Cost",
		ChartDataPoints: []*ChartDataPoint{},
	}
	valueChartItem := &ChartItem{
		Name: strings.ToUpper(investment.Symbol) + " - Value",
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
