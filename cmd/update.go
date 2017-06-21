package cmd

import (
	"git.jasonc.me/main/money/db/investment"
	"git.jasonc.me/main/money/db/price"
	"github.com/jchavannes/jgo/jerr"
)

func CmdUpdate(userId uint) error {
	investmentTransactions, err := investment.GetTransactionsForUser(userId)
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

func intInSlice(findItem uint, slice []uint) bool {
	for _, item := range slice {
		if item == findItem {
			return true
		}
	}
	return false
}
