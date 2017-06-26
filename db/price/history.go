package price

import (
	"git.jasonc.me/main/money/db"
)

func GetHistory(investment *db.Investment) ([]*db.InvestmentPrice, error) {
	return db.GetAllInvestmentPricesForInvestment(investment)
}
