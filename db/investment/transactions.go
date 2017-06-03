package investment

import (
	"git.jasonc.me/main/money/db"
	"time"
	"fmt"
)

func GetTransactionsForUser(userId uint) ([]*db.InvestmentTransaction, error) {
	return db.GetInvestmentTransactionsForUser(userId)
}

func AddTransaction(userId uint, investmentType string, symbol string, transactionType db.InvestmentTransactionType, date time.Time, price float32, quantity float32) error {
	investment := db.Investment{
		Symbol: symbol,
		InvestmentType: investmentType,
	}
	err := investment.Load()
	if err != nil {
		return fmt.Errorf("Error loading investment: %s", err)
	}
	investmentTransaction := db.InvestmentTransaction{
		UserId: userId,
		Type: transactionType.Uint(),
		InvestmentId: investment.Id,
		Investment: investment,
		Date: date,
		Price: price,
		Quantity: quantity,
	}
	err = investmentTransaction.Save()
	if err != nil {
		return fmt.Errorf("Error saving investment transaction: %s", err)
	}
	fmt.Printf("Saved transaction: %#v\n", investmentTransaction)
	return nil
}
