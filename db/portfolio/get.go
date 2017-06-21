package portfolio

import (
	"git.jasonc.me/main/money/db/investment"
	"git.jasonc.me/main/money/db"
	"github.com/jchavannes/jgo/jerr"
	"fmt"
)

func Get(userId uint) (*Portfolio, error) {
	investmentTransactions, err := investment.GetTransactionsForUser(userId)
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
			//Value: transaction.Quantity * lastInvestmentPrice.Price,
			//NetGainLoss: -624.80,
			//NetGainLossPercent: -46.02,
			//DistributionPercent: 0.67,
			//NetGainLossWeighted: -1.20,
		}
		portfolioItems = append(portfolioItems, portfolioItem)
	}
	for _, portfolioItem := range portfolioItems {
		portfolioItem.Value = portfolioItem.Quantity * portfolioItem.Price
		portfolioItem.NetGainLoss = portfolioItem.Value - portfolioItem.Cost
	}
	return &Portfolio{
		Items: portfolioItems,
	}, nil
}
