package price

import (
	"bufio"
	"fmt"
	"git.jasonc.me/main/money/db"
	"net/http"
	"regexp"
	"strconv"
	"github.com/jchavannes/jgo/jerr"
)

func UpdateStockInvestmentFromGoogleFinance(investment *db.Investment) error {
	var lastItemTimestamp int64

	lastItem, err := db.GetLastInvestmentPrice(investment)
	if err != nil {
		lastItemTimestamp = 0
	} else {
		lastItemTimestamp = lastItem.Timestamp
	}

	url := investment.GetGoogleFinanceUrl()

	fmt.Printf("Fetching data from: %s\n", url)
	resp, err := http.Get(url)

	if err != nil {
		return jerr.Get("Error getting stock data", err)
	}

	scanner := bufio.NewScanner(resp.Body)

	startsWithA := regexp.MustCompile(`^a`)
	startsWithNumber := regexp.MustCompile(`^\d`)

	var currentMarker int64

	totalRowsAdded := 0

	for scanner.Scan() {
		line := scanner.Text()

		var investmentPrice *db.InvestmentPrice
		var getDayPrice *regexp.Regexp

		switch {
		case startsWithA.MatchString(line):
			getDayPrice = regexp.MustCompile("^a(\\d+),([\\d.]+)")
		case startsWithNumber.MatchString(line):
			getDayPrice = regexp.MustCompile("^(\\d+),([\\d.]+)")
		default:
			continue
		}

		matches := getDayPrice.FindStringSubmatch(line)

		if matches == nil {
			return jerr.New("Error parsing data")
		}

		timestamp, _ := strconv.ParseInt(matches[1], 10, 64)
		price64, _ := strconv.ParseFloat(matches[2], 32)
		price := float32(price64)

		if timestamp < 1e4 {
			timestamp = currentMarker + (timestamp * 86400)
		} else {
			currentMarker = timestamp
		}

		if timestamp < lastItemTimestamp {
			continue
		}

		investmentPrice = &db.InvestmentPrice{
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

	err = scanner.Err()

	if err != nil {
		return jerr.Get("Error with scanner", err)
	}

	if totalRowsAdded == 0 {
		return jerr.New("No rows added")
	}

	fmt.Printf("Rows added/updated: %d\n", totalRowsAdded)
	return nil
}
