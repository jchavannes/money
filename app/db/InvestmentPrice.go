package db

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type InvestmentPrice struct {
	Id           uint `gorm:"primary_key"`
	Investment   Investment
	InvestmentId uint  `gorm:"unique_index:stock_id_timestamp"`
	Timestamp    int64 `gorm:"unique_index:stock_id_timestamp"`
	Price        float32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (i *InvestmentPrice) Print() {
	message := "Investment: %s Type: %s - Date: %s Price: $%.2f\n"
	tm := time.Unix(i.Timestamp, 0).Format(time.RFC3339)
	fmt.Printf(message, i.Investment.Symbol, i.Investment.InvestmentType, tm, i.Price)
}

func (i *InvestmentPrice) AddOrUpdate() error {
	findPrice := &InvestmentPrice{
		Id:           i.Id,
		InvestmentId: i.InvestmentId,
		Timestamp:    i.Timestamp,
	}
	result := find(findPrice, findPrice)
	if result.Error == nil {
		i.Id = findPrice.Id
	} else if ! IsRecordNotFoundError(result.Error) {
		return jerr.Get("Error looking for existing record", result.Error)
	}
	result = save(i)
	if result.Error != nil {
		return jerr.Get("Error saving investment price", result.Error)
	}
	return nil
}

func GetLastInvestmentPrice(investment *Investment) (*InvestmentPrice, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("Error getting db", err)
	}
	var lastInvestmentPrice InvestmentPrice
	result := db.Order("timestamp DESC").First(&lastInvestmentPrice, InvestmentPrice{InvestmentId: investment.Id})
	if result.Error != nil {
		return nil, jerr.Get("Error getting last investment price", result.Error)
	}
	return &lastInvestmentPrice, nil
}

func GetAllInvestmentPricesForInvestment(investment *Investment) ([]*InvestmentPrice, error) {
	var investmentPrices []*InvestmentPrice
	result := find(&investmentPrices, &InvestmentPrice{
		InvestmentId: investment.Id,
	})
	if result.Error != nil {
		return nil, jerr.Get("Error getting investment prices", result.Error)
	}
	return investmentPrices, nil
}
