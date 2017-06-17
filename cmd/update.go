package cmd

import (
	"git.jasonc.me/main/money/db/investment"
	"fmt"
	"git.jasonc.me/main/money/db/price"
)

func CmdUpdate(userId uint) error {
	investmentTransactions, err := investment.GetTransactionsForUser(userId)
	if err != nil {
		return fmt.Errorf("Error getting transactions for user: %s", err)
	}
	err = price.UpdateStockInvestment(&investmentTransactions[0].Investment)
	if err != nil {
		return fmt.Errorf("Error updating stock investments: %s", err)
	}
	return nil
}
