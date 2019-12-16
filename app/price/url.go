package price

import (
	"fmt"
	"github.com/jchavannes/money/app/db"
	"strings"
	"time"
)

func GetGoogleFinanceUrl(s db.Investment) string {
	var url = "https://www.google.com/finance/getprices?&i=86400&p=10Y&f=d,c,v,k,o,h,l&df=cpct"
	return url + "&q=" + strings.ToUpper(s.Symbol) + "&x=" + strings.ToUpper(s.InvestmentType)
}

func GetCoinMarketCapUrl(s db.Investment) string {
	var url = "https://graphs2.coinmarketcap.com/currencies/"
	return url + strings.ToLower(s.Symbol) + "/"
}

var symbolTickers = map[string]int{
	"bitcoin":      1,
	"litecoin":     2,
	"ripple":       52,
	"ethereum":     1027,
	"bitcoin-cash": 1831,
	"bitcoin-sv":   3602,
}

func GetIdFromSymbol(symbol string) int {
	ticker, _ := symbolTickers[symbol]
	return ticker
}

func GetCoinMarketCapUrlNew(s db.Investment) string {
	var baseUrl = "https://web-api.coinmarketcap.com/v1/cryptocurrency/quotes/historical?" +
		"format=chart_crypto_details&convert=USD&interval=1d&time_start=2013-04-28"
	return fmt.Sprintf("%s&id=%d&time_end=%s",
		baseUrl, GetIdFromSymbol(s.Symbol), time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
}

func GetCoinMarketCapUrlLatest(investments []db.Investment) string {
	var ids []string
	for _, investment := range investments {
		ids = append(ids, fmt.Sprintf("%d", GetIdFromSymbol(investment.Symbol)))
	}
	return fmt.Sprintf("https://web-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?id=%s", strings.Join(ids, ","))
}
