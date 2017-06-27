package investment

import "github.com/jchavannes/money/db"

func GetInvestmentsForType(investmentType string) ([]*db.Investment, error) {
	return db.GetInvestmentsForType(investmentType)
}
