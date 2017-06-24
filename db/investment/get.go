package investment

import (
	"git.jasonc.me/main/money/db"
	"github.com/jchavannes/jgo/jerr"
)

func Get(investmentType string, symbol string) (*db.Investment, error) {
	investment := &db.Investment{
		Symbol: symbol,
		InvestmentType: investmentType,
	}
	err := investment.Load()
	if err != nil {
		return nil, jerr.Get("Error loading investment", err)
	}
	return investment, nil
}
