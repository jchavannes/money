package db

import (
	"strings"
	"github.com/jchavannes/jgo/jerr"
)

type InvestmentType string

func (i InvestmentType) String() string {
	return string(i)
}

const (
	InvestmentType_NYSEMKT InvestmentType = "nysemkt"
	InvestmentType_NYSE    InvestmentType = "nyse"
	InvestmentType_NASDAQ  InvestmentType = "nasdaq"
	InvestmentType_Index   InvestmentType = "indexsp"
	InvestmentType_Crypto  InvestmentType = "crypto"
)

type Investment struct {
	Id             uint `gorm:"primary_key"`
	Symbol         string `gorm:"unique_index:investment_type_symbol"`
	InvestmentType string `gorm:"unique_index:investment_type_symbol"`
}

func (i *Investment) Load() error {
	result := find(i, i)
	if result.Error != nil {
		result = save(i)
		if result.Error != nil {
			return jerr.Get("Error saving investment", result.Error)
		}
	}
	return nil
}

func GetInvestment(investmentType string, symbol string) (*Investment, error) {
	investment := &Investment{
		Symbol:         strings.ToLower(symbol),
		InvestmentType: strings.ToLower(investmentType),
	}
	err := investment.Load()
	if err != nil {
		return nil, jerr.Get("Error loading investment", err)
	}
	return investment, nil
}

func GetInvestmentsForType(investmentType string) ([]*Investment, error) {
	var investments []*Investment
	result := find(&investments, &Investment{
		InvestmentType: investmentType,
	})
	if result.Error != nil {
		return []*Investment{}, result.Error
	}
	return investments, nil
}
