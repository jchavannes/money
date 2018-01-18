package price

import (
	"strings"
	"github.com/jchavannes/money/data/db"
)

func GetGoogleFinanceUrl(s db.Investment) string {
	var url = "https://www.google.com/finance/getprices?&i=86400&p=10Y&f=d,c,v,k,o,h,l&df=cpct"
	return url + "&q=" + strings.ToUpper(s.Symbol) + "&x=" + strings.ToUpper(s.InvestmentType)
}

func GetCoinMarketCapUrl(s db.Investment) string {
	var url = "https://graphs2.coinmarketcap.com/currencies/"
	return url + strings.ToLower(s.Symbol) + "/"
}
