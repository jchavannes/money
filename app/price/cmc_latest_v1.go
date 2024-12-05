package price

import (
	"encoding/json"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/app/db"
	"io/ioutil"
	"net/http"
	"time"
)

type CmcLatestJsonV1 struct {
	Data map[string]struct {
		Id     int
		Name   string
		Symbol string
		Slug   string
		Quote  map[string]struct {
			Price       float32
			LastUpdated string `json:"last_updated"`
		}
	}
}

func CmcLatestV1(investments []db.Investment) error {
	url := GetCmcLatestUrlV1(investments)

	fmt.Printf("Fetching data from: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return jerr.Get("Error getting crypto data", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var cryptoJson CmcLatestJsonV1
	json.Unmarshal(body, &cryptoJson)

	totalRowsAdded := 0

	fmt.Println("Finished parsing data")

	for _, item := range cryptoJson.Data {
		quote, ok := item.Quote["USD"]
		if !ok {
			return jerr.New("unable to find usd price")
		}
		date, err := time.Parse("2006-01-02T15:04:05.000Z", quote.LastUpdated)
		if err != nil {
			return jerr.Get("error parsing last updated time", err)
		}
		timestamp := date.Unix()
		price := quote.Price

		if quote.Price > 1000000 || quote.Price < 0.0001 {
			continue
		}
		var foundInvestment db.Investment
		for _, investment := range investments {
			if item.Id == GetIdFromSymbol(investment.Symbol) {
				foundInvestment = investment
			}
		}
		if foundInvestment.Id == 0 {
			return jerr.New("unable to find investment for price")
		}

		investmentPrice := &db.InvestmentPrice{
			Timestamp:    timestamp,
			Price:        price,
			InvestmentId: foundInvestment.Id,
			Investment:   foundInvestment,
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
