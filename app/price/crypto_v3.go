package price

import (
	"encoding/json"
	"fmt"
	"github.com/jchavannes/money/app/db"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func CmcHistoryV3(investment *db.Investment) error {
	lastItemTimestamp := getLastInvestmentTimestamp(investment)
	if isTimestampUpdatedRecently(lastItemTimestamp) {
		return nil
	}

	cmcBody, err := fetchCmcHistoryV3(investment)
	if err != nil {
		return fmt.Errorf("error getting cmc history v3; %w", err)
	}

	prices, err := ParseCoinMarketCapV3(cmcBody)
	if err != nil {
		return fmt.Errorf("error getting prices from cmc history v3; %w", err)
	}

	if err := saveInvestmentPrices(investment, prices, lastItemTimestamp); err != nil {
		return fmt.Errorf("error saving cmc history v3; %w", err)
	}
	return nil
}

func fetchCmcHistoryV3(investment *db.Investment) ([]byte, error) {
	url := GetCoinMarketCapUrlV3(investment.Symbol)

	fmt.Printf("Fetching data from: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting crypto data; %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body; %w", err)
	}

	return body, nil

}

type CryptoJsonV3 struct {
	Data struct {
		Points map[string]map[string][]float64 `json:"points"`
	}
}

func ParseCoinMarketCapV3(body []byte) ([]*TimePrice, error) {
	var cryptoJson CryptoJsonV3
	if err := json.Unmarshal(body, &cryptoJson); err != nil {
		return nil, fmt.Errorf("error parsing json coinmarketcap v3; %w", err)
	}

	var prices []*TimePrice
	for unixString, data := range cryptoJson.Data.Points {
		unixInt, err := strconv.ParseInt(unixString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing unix time; %w", err)
		}

		if _, ok := data["v"]; !ok || len(data["v"]) < 1 {
			return nil, fmt.Errorf("unable to find usd price")
		}

		prices = append(prices, &TimePrice{
			Time:  time.Unix(unixInt, 0),
			Price: data["v"][0],
		})
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Time.Before(prices[j].Time)
	})

	return prices, nil
}

func saveInvestmentPrices(investment *db.Investment, prices []*TimePrice, lastItemTimestamp int64) error {
	var totalRowsAdded int
	for _, p := range prices {
		if p.Price > 100000 || p.Price < 0.0001 || p.Time.Unix() < lastItemTimestamp {
			continue
		}

		investmentPrice := &db.InvestmentPrice{
			Investment:   *investment,
			InvestmentId: investment.Id,
			Price:        float32(p.Price),
			Timestamp:    p.Time.Unix(),
		}

		if err := investmentPrice.AddOrUpdate(); err != nil {
			return fmt.Errorf("error updating investment price: %w", err)
		}

		totalRowsAdded++
	}

	if totalRowsAdded == 0 {
		return fmt.Errorf("no rows added (%d prices)", len(prices))
	}

	fmt.Printf("Rows added/updated: %d\n", totalRowsAdded)
	return nil
}
