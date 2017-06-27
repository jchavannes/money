package price

import (
	"github.com/jchavannes/money/db"
)

func GetHistory(investment *db.Investment) ([]*db.InvestmentPrice, error) {
	return db.GetAllInvestmentPricesForInvestment(investment)
}
