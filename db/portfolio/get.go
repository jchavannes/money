package portfolio

import "git.jasonc.me/main/money/db"

func Get(userId uint) (*Portfolio, error) {
	return &Portfolio{
		Items: []PortfolioItem{{
			Investment: db.Investment{
				Symbol: "fslr",
				InvestmentType: string(db.InvestmentType_NASDAQ),
			},
			Quantity: 20,
			Price: 36.65,
			Value: 733,
			Cost: 1357.80,
			NetGainLoss: -624.80,
			NetGainLossPercent: -46.02,
			DistributionPercent: 0.67,
			NetGainLossWeighted: -1.20,
		}},
	}, nil
}
