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
	for _, transaction := range investmentTransactions {
		portfolioItem := PortfolioItem{
			Investment: transaction.Investment,
			Quantity: transaction.Quantity,
			Price: transaction.Price,
			Value: transaction.Quantity * transaction.Price,
			Cost: transaction.Quantity * transaction.Price,
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
