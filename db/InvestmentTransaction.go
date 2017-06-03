package db

import "time"

type InvestmentTransactionType uint

const (
	InvestmentTransactionType_Sell InvestmentTransactionType = 0
	InvestmentTransactionType_Buy InvestmentTransactionType = 1
)

type InvestmentTransaction struct {
	Id            uint `gorm:"primary_key"`
	InvestmentId  uint
	Price         float32
	Quantity      float32
	TransactionTs time.Time
	Type          InvestmentTransactionType
	UserId        uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func GetInvestmentTransactionsForUser(userId uint) ([]*InvestmentTransaction, error) {
	var investmentTransactions []*InvestmentTransaction
	result := find(&investmentTransactions, InvestmentTransaction{
		UserId: userId,
	})
	if result.Error != nil {
		return []*InvestmentTransaction{}, result.Error
	}
	return investmentTransactions, nil
}
