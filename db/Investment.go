package db

import (
	"fmt"
	"strings"
)

type InvestmentType string

func (i InvestmentType) String() string {
	return string(i)
}

const (
	InvestmentType_NYSEMKT InvestmentType = "nysemkt"
	InvestmentType_NYSE InvestmentType = "nyse"
	InvestmentType_NASDAQ InvestmentType = "nasdaq"
	InvestmentType_Index InvestmentType = "indexsp"
	InvestmentType_Crypto InvestmentType = "crypto"
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
			return fmt.Errorf("Error saving investment: %s", result.Error)
		}
	}
	return nil
}

func (s *Investment) GetUrl() string {
	var url = "https://www.google.com/finance/getprices?&i=86400&p=10Y&f=d,c,v,k,o,h,l&df=cpct"
	return url + "&q=" + s.Symbol + "&x=" + strings.ToUpper(s.InvestmentType)
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
