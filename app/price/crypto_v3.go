package price

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func GetCoinMarketCapUrlV3(symbol string) string {
	return fmt.Sprintf(
		"https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?range=1Y&id=%d",
		GetIdFromSymbol(symbol),
	)
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

	prices, err := GetPricesFromCoinMarketCapV3Json(cryptoJson)
	if err != nil {
		return nil, fmt.Errorf("error getting prices from coinmarketcap v3 json; %w", err)
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Time.Before(prices[j].Time)
	})

	return prices, nil
}

func GetPricesFromCoinMarketCapV3Json(cryptoJson CryptoJsonV3) ([]*TimePrice, error) {
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

	return prices, nil
}
