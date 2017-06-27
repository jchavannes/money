package price

import (
	"github.com/jchavannes/money/data/db"
	"github.com/jchavannes/jgo/jerr"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type CryptoJson struct {
	PriceUSD [][2]float32 `json:"price_usd"`
}

func UpdateCryptoInvestmentFromCoinMarketCap(investment *db.Investment) error {
	var lastItemTimestamp int64

	lastItem, err := db.GetLastInvestmentPrice(investment)
	if err != nil {
		lastItemTimestamp = 0
	} else {
		lastItemTimestamp = lastItem.Timestamp
	}

	url := GetCoinMarketCapUrl(*investment)

	fmt.Printf("Fetching data from: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return jerr.Get("Error getting crypto data", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var cryptoJson CryptoJson
	json.Unmarshal(body, &cryptoJson)

	totalRowsAdded := 0

	fmt.Println("Finished parsing data")

	for i := range cryptoJson.PriceUSD {
		timestamp := int64(cryptoJson.PriceUSD[i][0]) / 1000
		price := cryptoJson.PriceUSD[i][1]

		if price > 10000 || price < 0.0001 || timestamp < lastItemTimestamp {
			continue
		}

		investmentPrice := &db.InvestmentPrice{
			Timestamp: timestamp,
			Price: price,
			InvestmentId: investment.Id,
			Investment: *investment,
		}

		investmentPrice.Print()

		err = investmentPrice.AddOrUpdate()
		if err != nil {
			return jerr.Get(fmt.Sprintf("Error updating investment price: %#v", investmentPrice), err)
		}

		totalRowsAdded++
	}

	if totalRowsAdded == 0 {
		return jerr.New("No rows added")
	}

	fmt.Printf("Rows added/updated: %d\n", totalRowsAdded)
	return nil
}
