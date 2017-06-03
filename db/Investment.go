package db

type InvestmentType string

const (
	InvestmentType_NYSEMKT InvestmentType = "NYSEMKT"
	InvestmentType_NYSE InvestmentType = "NYSE"
	InvestmentType_NASDAQ InvestmentType = "NASDAQ"
	InvestmentType_Index InvestmentType = "INDEXSP"
	InvestmentType_Crypto InvestmentType = "CRYPTO"
)

type Investment struct {
	Id             uint `gorm:"primary_key"`
	Symbol         string `gorm:"unique_index:investment_type_symbol"`
	InvestmentType InvestmentType `gorm:"unique_index:investment_type_symbol"`
}

func (s *Investment) GetUrl() string {
	var url = "https://www.google.com/finance/getprices?&i=86400&p=10Y&f=d,c,v,k,o,h,l&df=cpct"
	return url + "&q=" + s.Symbol + "&x=" + string(s.InvestmentType)
}
