package price

import (
	"git.jasonc.me/main/money/db"
)

func UpdateInvestment(investment *db.Investment) error {
	if (investment.InvestmentType == db.InvestmentType_Crypto.String()) {
		return UpdateCryptoInvestmentFromCoinMarketCap(investment)
	} else {
		return UpdateStockInvestmentFromGoogleFinance(investment)
	}
}
