package price

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/app/db"
)

func UpdateInvestment(investment *db.Investment) error {
	if investment.InvestmentType == db.InvestmentType_Crypto.String() {
		return CmcHistoryV1(investment)
	} else {
		return UpdateStockInvestmentFromGoogleFinance(investment)
	}
}

func UpdateForUser(userId uint) error {
	investmentTransactions, err := db.GetInvestmentTransactionsForUser(userId)
	if err != nil {
		return jerr.Get("Error getting transactions for user", err)
	}
	completedInvestmentIds := []uint{}
	var cryptoInvestments []db.Investment
	for _, investmentTransaction := range investmentTransactions {
		if intInSlice(investmentTransaction.InvestmentId, completedInvestmentIds) {
			continue
		}
		if investmentTransaction.Investment.InvestmentType == string(db.InvestmentType_Crypto) {
			cryptoInvestments = append(cryptoInvestments, investmentTransaction.Investment)
		}
		err = UpdateInvestment(&investmentTransaction.Investment)
		if err != nil {
			return jerr.Get("Error updating stock investments", err)
		}
		completedInvestmentIds = append(completedInvestmentIds, investmentTransaction.InvestmentId)
	}
	if len(cryptoInvestments) > 0 {
		err := CmcLatestV1(cryptoInvestments)
		if err != nil {
			return jerr.Get("error updating latest crypto investments", err)
		}
	}
	return nil
}

func UpdateInvestmentById(investmentId uint) error {
	investmentToUpdate := db.Investment{
		Id: investmentId,
	}
	err := investmentToUpdate.Load()
	if err != nil {
		return jerr.Get("Error loading investment", err)
	}

	err = UpdateInvestment(&investmentToUpdate)
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
