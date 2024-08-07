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

const (
	SymbolBitcoin        = "bitcoin"
	SymbolLitecoin       = "litecoin"
	SymbolRipple         = "ripple"
	SymbolEthereum       = "ethereum"
	SymbolBitcoinCash    = "bitcoin-cash"
	SymbolBitcoinSv      = "bitcoin-sv"
	SymbolBitcoinCashAbc = "bitcoin-cash-abc"
	SymbolEcash          = "ecash"

	IdBitcoin        = 1
	IdLitecoin       = 2
	IdRipple         = 52
	IdEthereum       = 1027
	IdBitcoinCash    = 1831
	IdBitcoinSv      = 3602
	IdBitcoinCashAbc = 7686
	IdEcash          = 10791
)

var symbolTickers = map[string]int{
	SymbolBitcoin:        IdBitcoin,
	SymbolLitecoin:       IdLitecoin,
	SymbolRipple:         IdRipple,
	SymbolEthereum:       IdEthereum,
	SymbolBitcoinCash:    IdBitcoinCash,
	SymbolBitcoinSv:      IdBitcoinSv,
	SymbolBitcoinCashAbc: IdBitcoinCashAbc,
	SymbolEcash:          IdEcash,
}

func GetIdFromSymbol(symbol string) int {
	ticker, _ := symbolTickers[symbol]
	return ticker
}

func GetCmcHistoryUrlV1(s db.Investment) string {
	var baseUrl = "https://web-api.coinmarketcap.com/v1/cryptocurrency/quotes/historical?" +
		"format=chart_crypto_details&convert=USD&interval=1d&time_start=2013-04-28"
	return fmt.Sprintf("%s&id=%d&time_end=%s",
		baseUrl, GetIdFromSymbol(s.Symbol), time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
}

func GetCmcLatestUrlV1(investments []db.Investment) string {
	var ids []string
	for _, investment := range investments {
		ids = append(ids, fmt.Sprintf("%d", GetIdFromSymbol(investment.Symbol)))
	}
	return fmt.Sprintf("https://web-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?id=%s",
		strings.Join(ids, ","))
}

func GetCoinMarketCapUrlV3(symbol string) string {
	return fmt.Sprintf(
		"https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?range=1Y&id=%d",
		GetIdFromSymbol(symbol),
	)
}

func GetCmcPushPriceUrl() string {
	return "wss://push.coinmarketcap.com/ws?device=web&client_source=coin_detail_page"
}

func GetCmcPushPriceSubscribeMessage(investments []db.Investment) string {
	var ids = make([]string, len(investments))
	for i := range investments {
		ids[i] = fmt.Sprintf("%d", GetIdFromSymbol(investments[i].Symbol))
	}
	return fmt.Sprintf(`{"method":"RSUBSCRIPTION","params":["main-site@crypto_price_5s@{}@normal","%s"]}`,
		strings.Join(ids, ","))
}
