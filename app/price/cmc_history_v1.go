package price

import (
	"encoding/json"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/app/db"
	"io"
	"net/http"
	"sort"
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

	prices, err := GetPricesFromCmcV1History(cmcJson)
	if err != nil {
		return jerr.Get("Error getting prices from cmc history new", err)
	}

	if err := saveInvestmentPrices(investment, prices, lastItemTimestamp); err != nil {
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

func GetPricesFromCmcV1History(cmcJson *CmcHistoryJsonV1) ([]*TimePrice, error) {
	var timePrices []*TimePrice
	for dateString, prices := range cmcJson.Data {
		date, err := time.Parse("2006-01-02T15:04:05.000Z", dateString)
		if err != nil {
			return nil, jerr.Get("Error parsing time", err)
		}

		usd, ok := prices["USD"]
		if !ok || len(usd) < 1 {
			return nil, jerr.New("Unable to find usd price")
		}
		price := usd[0]

		timePrices = append(timePrices, &TimePrice{
			Time:  date,
			Price: float64(price),
		})
	}

	sort.Slice(timePrices, func(i, j int) bool {
		return timePrices[i].Time.Before(timePrices[j].Time)
	})

	return timePrices, nil

}
