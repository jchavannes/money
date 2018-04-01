package price

import (
	"github.com/jchavannes/money/app/db"
)

func GetHistory(investment *db.Investment) ([]*db.InvestmentPrice, error) {
	return db.GetAllInvestmentPricesForInvestment(investment)
}
