package db

import (
	"time"
	"github.com/jchavannes/jgo/jerr"
)

type InvestmentTransactionType uint

func (i InvestmentTransactionType) Uint() uint {
	return uint(i)
}

const (
	InvestmentTransactionType_Sell InvestmentTransactionType = 0
	InvestmentTransactionType_Buy InvestmentTransactionType = 1
)

type InvestmentTransaction struct {
	Id           uint `gorm:"primary_key"`
	Date         time.Time
	Investment   Investment
	InvestmentId uint
	Price        float32
	Quantity     float32
	Type         uint
	UserId       uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (i *InvestmentTransaction) Save() error {
	result := save(i)
	if result.Error != nil {
		return jerr.Get("Error saving investment transaction", result.Error)
	}
	return nil
}

func (i *InvestmentTransaction) Load() error {
	result := find(i, i)
	if result.Error != nil {
		return jerr.Get("Error finding investment transaction", result.Error)
	}
	return nil
}

func (i *InvestmentTransaction) Delete() error {
	result := remove(i)
	if result.Error != nil {
		return jerr.Get("Error removing investment transaction", result.Error)
	}
	return nil
}

func GetInvestmentTransactionsForUser(userId uint) ([]*InvestmentTransaction, error) {
	var investmentTransactions []*InvestmentTransaction
	result := find(&investmentTransactions, &InvestmentTransaction{
		UserId: userId,
	})
	if result.Error != nil {
		return []*InvestmentTransaction{}, result.Error
	}
	return investmentTransactions, nil
}
