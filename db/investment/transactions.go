package investment

import "git.jasonc.me/main/money/db"

func GetInvestmentTransactionsForUser(userId uint) ([]*db.InvestmentTransaction, error) {
	return db.GetInvestmentTransactionsForUser(userId)
}
