package investment

import (
	"git.jasonc.me/main/money/db"
	"time"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

func GetTransactionsForUser(userId uint) ([]*db.InvestmentTransaction, error) {
	transactions, err := db.GetInvestmentTransactionsForUser(userId)
	if err != nil {
		return []*db.InvestmentTransaction{}, jerr.Get("Error getting investment transactions for user", err)
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
		return jerr.Get("Error loading investment", err)
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
		return jerr.Get("Error saving investment transaction", err)
	}
	fmt.Printf("Saved transaction: %#v\n", investmentTransaction)
	return nil
}

func DeleteTransaction(userId uint, investmentTransactionId uint) error {
	investmentTransaction := db.InvestmentTransaction{
		Id: investmentTransactionId,
	}
	err := investmentTransaction.Load()
	if err != nil {
		return jerr.Get("Error loading transaction", err)
	}
	if investmentTransaction.UserId != userId {
		return jerr.New("UserId does not match")
	}
	err = investmentTransaction.Delete()
	if err != nil {
		return jerr.Get("Error deleting transaction", err)
	}
	return nil
}
