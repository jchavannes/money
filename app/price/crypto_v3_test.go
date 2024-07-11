package price_test

import (
	"github.com/jchavannes/money/app/db"
	"github.com/jchavannes/money/app/price"
	"log"
	"os"
	"testing"
	"time"
)

func TestCmcHistoryV3(t *testing.T) {
	t.Skip("Skipping test TestCmcHistoryV3 since it will make a network call to fetch data from coinmarketcap.com")
	history, err := price.FetchCmcHistoryV3(&db.Investment{
		Symbol: price.SymbolBitcoinCash,
	})
	if err != nil {
		t.Errorf("error fetching cmc history v3; %v", err)
		return
	}
	log.Printf("history: %s\n", history)
}

func TestGetCoinMarketCapUrlV3(t *testing.T) {
	url := price.GetCoinMarketCapUrlV3(price.SymbolBitcoinCash)
	expectedUrl := "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?range=1Y&id=1831"

	if url != expectedUrl {
		t.Errorf("Expected url to be %s, got %s", expectedUrl, url)
	}
}

func TestCryptoV3HandleResponse(t *testing.T) {
	cryptoV3SampleResponse, err := os.ReadFile("./testdata/crypto_v3_response.json")
	if err != nil {
		t.Errorf("error reading crypto v3 sample response; %v", err)
		return
	}

	prices, err := price.ParseCoinMarketCapV3(cryptoV3SampleResponse)
	if err != nil {
		t.Errorf("error parsing coinmarketcap v3 response; %v", err)
	} else if len(prices) == 0 {
		t.Error("error coinmarketcap v3 response empty")
	}

	log.Printf("prices count: %d\n", len(prices))

	for _, p := range prices {
		log.Printf("time: %s, price: %f\n", p.Time.Format(time.RFC3339), p.Price)
	}
}
