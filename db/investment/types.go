package investment

import "git.jasonc.me/main/money/db"

func GetInvestmentsForType(investmentType string) ([]*db.Investment, error) {
	return db.GetInvestmentsForType(investmentType)
}
