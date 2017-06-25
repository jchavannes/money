package investment

import (
	"git.jasonc.me/main/money/db"
	"github.com/jchavannes/jgo/jerr"
	"strings"
)

func Get(investmentType string, symbol string) (*db.Investment, error) {
	investment := &db.Investment{
		Symbol: strings.ToLower(symbol),
		InvestmentType: strings.ToLower(investmentType),
	}
	err := investment.Load()
	if err != nil {
		return nil, jerr.Get("Error loading investment", err)
	}
	return investment, nil
}
