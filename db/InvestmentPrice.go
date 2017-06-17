package db

import (
	"time"
	"fmt"
)

type InvestmentPrice struct {
	Id           uint `gorm:"primary_key"`
	Investment   Investment
	InvestmentId uint
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
		Id: i.Id,
		InvestmentId: i.InvestmentId,
		Timestamp: i.Timestamp,
	}
	result := find(findPrice, findPrice)
	if result.Error == nil {
		i.Id = findPrice.Id
	} else if ! isRecordNotFoundError(result.Error) {
		return fmt.Errorf("Error looking for existing record: %s", result.Error)
	}
	result = save(i)
	if result.Error != nil {
		return fmt.Errorf("Error saving investment price: %s", result.Error)
	}
	return nil
}

func GetLastInvestmentPrice(investment *Investment) (*InvestmentPrice, error) {
	db, err := getDb()
	if err != nil {
		return nil, fmt.Errorf("Error getting db: %s", err)
	}
	var lastItem InvestmentPrice
	result := db.Last(&lastItem, InvestmentPrice{InvestmentId: investment.Id})
	if result.Error != nil {
		return nil, fmt.Errorf("Error getting last item price: %s", result.Error)
	}
	return &lastItem, nil
}
