package investment

import (
	"git.jasonc.me/main/money/db"
	"time"
	"fmt"
	"errors"
)

func GetTransactionsForUser(userId uint) ([]*db.InvestmentTransaction, error) {
	transactions, err := db.GetInvestmentTransactionsForUser(userId)
	if err != nil {
		return []*db.InvestmentTransaction{}, fmt.Errorf("Error getting investment transactions for user: %s", err)
	}
	for _, transaction := range transactions {
		transaction.Investment.Id = transaction.InvestmentId
		transaction.Investment.Load()
	}
	return transactions, nil
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

func DeleteTransaction(userId uint, investmentId uint) error {
	investmentTransaction := db.InvestmentTransaction{
		Id: investmentId,
	}
	err := investmentTransaction.Load()
	if err != nil {
		return fmt.Errorf("Error loading transaction: %s", err)
	}
	if investmentTransaction.UserId != userId {
		return errors.New("UserId does not match")
	}
	err = investmentTransaction.Delete()
	if err != nil {
		return fmt.Errorf("Error deleting transaction: %s", err)
	}
	return nil
}
