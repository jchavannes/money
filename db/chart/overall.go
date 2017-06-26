package chart

import (
	"git.jasonc.me/main/money/db"
	"git.jasonc.me/main/money/db/portfolio"
	"github.com/jchavannes/jgo/jerr"
	"git.jasonc.me/main/money/db/investment"
	"time"
	"git.jasonc.me/main/money/db/price"
	"github.com/jchavannes/jgo/jtime"
)

func GetOverallChartData(userId uint) (*ChartDataOutput, error) {
	investmentTransactions, err := investment.GetTransactionsForUser(userId)
	if err != nil {
		return nil, jerr.Get("Error getting investment transactions", err)
	}
	userPortfolio, err := portfolio.Get(userId)
	if err != nil {
		return nil, jerr.Get("Error getting user portfolio", err)
	}
	costChartItem := ChartItem{
		Name: "Cost",
		ChartDataPoints: []*ChartDataPoint{},
	}
	valueChartItem := ChartItem{
		Name: "Value",
		ChartDataPoints: []*ChartDataPoint{},
	}
	for _, portfolioItem := range userPortfolio.Items {
		history, err := price.GetHistory(&portfolioItem.Investment)
		if err != nil {
			return nil, jerr.Get("Error getting history", err)
		}
		for _, investmentPrice := range history {
			investmentPriceTimestamp := time.Unix(jtime.RoundTimeToDay(investmentPrice.Timestamp), 0)
			costDataPoint, err := costChartItem.GetDataPoint(investmentPriceTimestamp)
			if err != nil {
				costDataPoint = &ChartDataPoint{
					Timestamp: investmentPriceTimestamp,
				}
				costChartItem.ChartDataPoints = append(costChartItem.ChartDataPoints, costDataPoint)
			}
			valueDataPoint, err := valueChartItem.GetDataPoint(investmentPriceTimestamp)
			if err != nil {
				valueDataPoint = &ChartDataPoint{
					Timestamp: investmentPriceTimestamp,
				}
				valueChartItem.ChartDataPoints = append(valueChartItem.ChartDataPoints, valueDataPoint)
			}
			for _, investmentTransactions := range investmentTransactions {
				if investmentTransactions.InvestmentId == portfolioItem.Investment.Id {
					value := investmentPrice.Price * investmentTransactions.Quantity
					cost := investmentTransactions.Price * investmentTransactions.Quantity
					if investmentTransactions.Type == uint(db.InvestmentTransactionType_Buy) {
						valueDataPoint.Amount += value
						costDataPoint.Amount += cost
					} else {
						valueDataPoint.Amount -= value
						costDataPoint.Amount -= cost
					}
				}
			}
		}
	}
	chartData := ChartData{
		Title: "Overall",
		ChartItems: []*ChartItem{
			&costChartItem,
			&valueChartItem,
		},
	}
	return chartData.GetChartDataOutput(), nil
}
