package price

import (
	"git.jasonc.me/main/money/db"
	"git.jasonc.me/main/money/db/portfolio"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func GetHistory(investment *db.Investment) ([]*db.InvestmentPrice, error) {
	return db.GetAllInvestmentPricesForInvestment(investment)
}

type ChartDataOutput struct {
	Title string
	Data  map[string][2]float64
}

type ChartData struct {
	Title string
	Data map[string]ChartDataItem
}

func (c *ChartData) GetChartDataOutput() ChartDataOutput {

}

type ChartDataItem struct {
	Name string
	Prices map[int64]float32
}

func GetOverallChartData(userId uint) (*ChartDataOutput, error) {
	userPortfolio, err := portfolio.Get(userId)
	if err != nil {
		return nil, err
	}
	costData := ChartDataItem{
		Name: "Cost",
		Prices: map[int64]float32,
	}
	valueData := ChartDataItem{
		Name: "Cost",
		Prices: map[int64]float32,
	}
	for _, portfolioItem := range userPortfolio.Items {
		history, err := GetHistory(&portfolioItem.Investment)
		if err != nil {
			return nil, jerr.Get("Error getting history", err)
		}
		for _, investmentPrice := range history {
			if _, ok := costData.Prices[investmentPrice.Timestamp]; ! ok {
				costData.Prices[investmentPrice.Timestamp] = 0.0
			}
			costData.Prices[investmentPrice.Timestamp] += investmentPrice.Price * 
		}
		// history
	}
	chartData := &ChartDataOutput{
		Title: "Overall",
		Data: data,
	}
	return chartData, nil
}
