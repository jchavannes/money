package portfolio

import (
	"github.com/jchavannes/money/app/db"
	"time"
)

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
	LastUpdate          time.Time
}

type PortfolioItemSorter []*PortfolioItem

func (pis PortfolioItemSorter) Len() int {
	return len(pis)
}
func (pis PortfolioItemSorter) Swap(i, j int) {
	pis[i], pis[j] = pis[j], pis[i]
}
func (pis PortfolioItemSorter) Less(i, j int) bool {
	return pis[i].Value > pis[j].Value
}
