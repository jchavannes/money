package portfolio

import (
	"github.com/jchavannes/money/data/db"
	"github.com/jchavannes/jgo/jerr"
	"fmt"
)

func Get(userId uint) (*Portfolio, error) {
	investmentTransactions, err := db.GetTransactionsForUser(userId)
	if err != nil {
		return nil, jerr.Get("Error getting transactions for user", err)
	}
	portfolioItems := []*PortfolioItem{}
	InvestmentTransactionsLoop:
	for _, transaction := range investmentTransactions {
		for i := range portfolioItems {
			if portfolioItems[i].Investment == transaction.Investment {
				portfolioItems[i].Quantity += transaction.Quantity
				portfolioItems[i].Cost += transaction.Quantity * transaction.Price
				continue InvestmentTransactionsLoop
			}
		}
		lastInvestmentPrice, err := db.GetLastInvestmentPrice(&transaction.Investment)
		if err != nil {
			errMessage := fmt.Sprintf("Error getting last investment price for investment: %#v", transaction.Investment)
			return nil, jerr.Get(errMessage, err)
		}
		portfolioItem := &PortfolioItem{
			Investment: transaction.Investment,
			Quantity: transaction.Quantity,
			Price: lastInvestmentPrice.Price,
			Cost: transaction.Quantity * transaction.Price,
			LastUpdate: lastInvestmentPrice.UpdatedAt,
		}
		portfolioItems = append(portfolioItems, portfolioItem)
	}
	var totalValue float32
	var totalCost float32
	for _, portfolioItem := range portfolioItems {
		portfolioItem.Value = portfolioItem.Quantity * portfolioItem.Price
		portfolioItem.NetGainLoss = portfolioItem.Value - portfolioItem.Cost
		if portfolioItem.Cost > 0 {
			portfolioItem.NetGainLossPercent = (portfolioItem.Value - portfolioItem.Cost) / portfolioItem.Cost
		}
		totalValue += portfolioItem.Value
		totalCost += portfolioItem.Cost
	}
	portfolio := &Portfolio{
		Items: portfolioItems,
		TotalValue: totalValue,
		TotalCost: totalCost,
		NetGainLoss: totalValue - totalCost,
	}
	if totalCost > 0 {
		portfolio.NetGainLossPercent = (totalValue - totalCost) / totalCost
	}
	if totalValue > 0 {
		for _, portfolioItem := range portfolioItems {
			portfolioItem.DistributionPercent = portfolioItem.Value / totalValue
			portfolioItem.NetGainLossWeighted = portfolioItem.NetGainLoss / (totalValue - totalCost) * portfolio.NetGainLossPercent
		}
	}
	return portfolio, nil
}
