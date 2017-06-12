package portfolio

import "git.jasonc.me/main/money/db"

type PortfolioItem struct {
	Investment          db.Investment
	Url                 string
	Quantity            float32
	Price               float32
	Value               float32
	Cost                float32
	NetGainLoss         float32
	NetGainLossPercent  float32
	DistributionPercent float32
	NetGainLossWeighted float32
}
