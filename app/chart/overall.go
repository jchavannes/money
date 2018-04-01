package chart

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/app/db"
	"github.com/jchavannes/money/app/portfolio"
)

func GetOverallChartData(userId uint) (*ChartDataOutput, error) {
	investmentTransactions, err := db.GetInvestmentTransactionsForUser(userId)
	if err != nil {
		return nil, jerr.Get("Error getting investment transactions", err)
	}
	userPortfolio, err := portfolio.Get(userId)
	if err != nil {
		return nil, jerr.Get("Error getting user portfolio", err)
	}

	costChartItem := &ChartItem{Name: "Cost"}
	valueChartItem := &ChartItem{Name: "Value"}
	for _, portfolioItem := range userPortfolio.Items {
		individualChartItems, err := GetIndividualChartItems(&portfolioItem.Investment, investmentTransactions)
		if err != nil {
			return nil, jerr.Get("Error getting individual chart items", err)
		}
		if len(individualChartItems) != 2 {
			return nil, jerr.New("Unexpected number of chart items")
		}
		for _, individualCostDataPoint := range individualChartItems[0].ChartDataPoints {
			var overallChartDataPoint *ChartDataPoint
			overallChartDataPoint, err := costChartItem.GetDataPoint(individualCostDataPoint.Timestamp)
			if err != nil {
				overallChartDataPoint = &ChartDataPoint{
					Timestamp: individualCostDataPoint.Timestamp,
				}
				costChartItem.ChartDataPoints = append(costChartItem.ChartDataPoints, overallChartDataPoint)
			}
			overallChartDataPoint.Amount += individualCostDataPoint.Amount
		}
		for _, individualValueDataPoint := range individualChartItems[1].ChartDataPoints {
			var overallChartDataPoint *ChartDataPoint
			overallChartDataPoint, err := valueChartItem.GetDataPoint(individualValueDataPoint.Timestamp)
			if err != nil {
				overallChartDataPoint = &ChartDataPoint{
					Timestamp: individualValueDataPoint.Timestamp,
				}
				valueChartItem.ChartDataPoints = append(valueChartItem.ChartDataPoints, overallChartDataPoint)
			}
			overallChartDataPoint.Amount += individualValueDataPoint.Amount
		}
	}
	chartData := ChartData{
		Title: "Overall",
		ChartItems: []*ChartItem{
			valueChartItem,
			costChartItem,
		},
	}
	return chartData.GetChartDataOutput(), nil
}
