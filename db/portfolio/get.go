package portfolio

import (
	"git.jasonc.me/main/money/db/investment"
	"fmt"
)

func Get(userId uint) (*Portfolio, error) {
	investmentTransactions, err := investment.GetTransactionsForUser(userId)
	if err != nil {
		return nil, fmt.Errorf("Error getting transactions for user: %s", err)
	}
	portfolioItems := []PortfolioItem{}
	InvestmentTransactionsLoop:
	for _, transaction := range investmentTransactions {
		for i := range portfolioItems {
			if portfolioItems[i].Investment == transaction.Investment {
				portfolioItems[i].Price = transaction.Price
				portfolioItems[i].Quantity += transaction.Quantity
				portfolioItems[i].Cost += transaction.Quantity * transaction.Price
				continue InvestmentTransactionsLoop
			}
		}
		portfolioItem := PortfolioItem{
			Investment: transaction.Investment,
			Quantity: transaction.Quantity,
			Price: transaction.Price,
			Cost: transaction.Quantity * transaction.Price,
			//Value: transaction.Quantity * transaction.Price,
			//NetGainLoss: -624.80,
			//NetGainLossPercent: -46.02,
			//DistributionPercent: 0.67,
			//NetGainLossWeighted: -1.20,
		}
		portfolioItems = append(portfolioItems, portfolioItem)
	}
	return &Portfolio{
		Items: portfolioItems,
	}, nil
}
