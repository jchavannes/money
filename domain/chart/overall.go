package chart

import (
	"github.com/jchavannes/money/data/db"
	"github.com/jchavannes/money/object/portfolio"
	"github.com/jchavannes/jgo/jerr"
	"time"
	"github.com/jchavannes/money/object/price"
	"github.com/jchavannes/jgo/jtime"
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

	for _, investmentTransaction := range investmentTransactions {
		history, err := price.GetHistory(&investmentTransaction.Investment)
		if err != nil {
			return nil, jerr.Get("Error getting history", err)
		}
		for _, investmentPrice := range history {
			investmentPriceTimestamp := time.Unix(jtime.RoundTimeToDay(investmentPrice.Timestamp), 0)
			inSellRange := investmentTransaction.Type == uint(db.InvestmentTransactionType_Sell) && investmentPriceTimestamp.Unix() < investmentTransaction.Date.Unix()
			inBuyRange := investmentTransaction.Type == uint(db.InvestmentTransactionType_Buy) && investmentPriceTimestamp.Unix() > investmentTransaction.Date.Unix()
			if ! inSellRange && ! inBuyRange {
				continue
			}

		}
	}

	chartItems := []*ChartItem{}

	for _, portfolioItem := range userPortfolio.Items {
		individualChartItems, err := GetIndividualChartItems(&portfolioItem.Investment, investmentTransactions)
		if err != nil {
			return nil, jerr.Get("Error getting individual chart items", err)
		}
		chartItems = append(chartItems, individualChartItems...)
	}
	chartData := ChartData{
		Title: "Overall",
		ChartItems: chartItems,
	}
	return chartData.GetChartDataOutput(), nil
}
