package price

import (
	"encoding/json"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/app/db"
	"io"
	"net/http"
	"time"
)

type CmcHistoryJsonV1 struct {
	Data map[string]map[string][]float32
}

func CmcHistoryV1(investment *db.Investment) error {
	lastItemTimestamp := getLastInvestmentTimestamp(investment)
	if isTimestampUpdatedRecently(lastItemTimestamp) {
		return nil
	}

	cmcJson, err := fetchCmcHistoryJsonV1(investment)
	if err != nil {
		return jerr.Get("Error getting crypto json new", err)
	}

	if err := saveCmcHistoryJsonV1(investment, cmcJson, lastItemTimestamp); err != nil {
		return jerr.Get("Error saving crypto json new", err)
	}
	return nil
}

func fetchCmcHistoryJsonV1(investment *db.Investment) (*CmcHistoryJsonV1, error) {
	url := GetCmcHistoryUrlV1(*investment)

	fmt.Printf("Fetching data from: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, jerr.Get("Error getting crypto data", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, jerr.Get("Error reading body", err)
	}

	var cmcJson CmcHistoryJsonV1
	if err := json.Unmarshal(body, &cmcJson); err != nil {
		return nil, jerr.Get("Error parsing json coinmarketcap new", err)
	}

	fmt.Println("Finished parsing data")

	return &cmcJson, nil
}

func saveCmcHistoryJsonV1(investment *db.Investment, cmcJson *CmcHistoryJsonV1, lastItemTimestamp int64) error {
	var totalRowsAdded int
	for dateString, prices := range cmcJson.Data {
		investmentPrice, err := getInvestmentPriceFromCmcJsonV1(investment, dateString, prices)
		if err != nil {
			return jerr.Get("Error getting investment price", err)
		}

		if investmentPrice.Price > 100000 || investmentPrice.Price < 0.0001 ||
			investmentPrice.Timestamp < lastItemTimestamp {
			continue
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

func getInvestmentPriceFromCmcJsonV1(investment *db.Investment, dateString string, prices map[string][]float32) (*db.InvestmentPrice, error) {
	date, err := time.Parse("2006-01-02T15:04:05.000Z", dateString)
	if err != nil {
		return nil, jerr.Get("error parsing time", err)
	}
	timestamp := date.Unix()

	usd, ok := prices["USD"]
	if !ok || len(usd) < 1 {
		return nil, jerr.New("unable to find usd price")
	}
	price := usd[0]

	var investmentPrice = &db.InvestmentPrice{
		Timestamp:    timestamp,
		Price:        price,
		InvestmentId: investment.Id,
		Investment:   *investment,
	}

	return investmentPrice, nil
}
