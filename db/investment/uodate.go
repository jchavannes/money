package investment

import (
	"github.com/jchavannes/jgo/jerr"
	"git.jasonc.me/main/money/db/price"
	"git.jasonc.me/main/money/db"
)

func UpdateForUser(userId uint) error {
	investmentTransactions, err := GetTransactionsForUser(userId)
	if err != nil {
		return jerr.Get("Error getting transactions for user", err)
	}
	completedInvestmentIds := []uint{}
	for _, investmentTransaction := range investmentTransactions {
		if intInSlice(investmentTransaction.InvestmentId, completedInvestmentIds) {
			continue
		}
		err = price.UpdateStockInvestmentFromGoogleFinance(&investmentTransaction.Investment)
		if err != nil {
			return jerr.Get("Error updating stock investments", err)
		}
		completedInvestmentIds = append(completedInvestmentIds, investmentTransaction.InvestmentId)
	}
	return nil
}

func UpdateInvestment(investmentId uint) error {
	investment := db.Investment{
		Id: investmentId,
	}
	err := investment.Load()
	if err != nil {
		return jerr.Get("Error loading investment", err)
	}

	err = price.UpdateStockInvestmentFromGoogleFinance(&investment)
	if err != nil {
		return jerr.Get("Error updating stock investments", err)
	}
	return nil
}

func intInSlice(findItem uint, slice []uint) bool {
	for _, item := range slice {
		if item == findItem {
			return true
		}
	}
	return false
}
