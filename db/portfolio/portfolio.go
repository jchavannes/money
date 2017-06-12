package portfolio

type Portfolio struct {
	Items              []PortfolioItem
	TotalValue         float32
	TotalCost          float32
	NetGainLoss        float32
	NetGainLossPercent float32
}
