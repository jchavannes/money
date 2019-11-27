package price

import (
	"github.com/jchavannes/money/app/db"
)

func GetHistory(investment *db.Investment) ([]*db.InvestmentPrice, error) {
	return db.GetAllInvestmentPricesForInvestment(investment)
}

func GetRecentPrice(investment *db.Investment) (*db.InvestmentPrice, error) {
	return db.GetLastInvestmentPrice(investment)
}
