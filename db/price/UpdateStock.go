package price

import (
	"bufio"
	"fmt"
	"git.jasonc.me/main/money/db"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"errors"
)

func UpdateStockInvestment(investment *db.Investment) error {
	lastItem, err := db.GetLastInvestmentPrice(investment)
	var lastItemTimestamp int64

	if err != nil {
		lastItemTimestamp = 0
		return err
	} else {
		lastItemTimestamp = lastItem.Timestamp
	}

	url := investment.GetUrl()

	fmt.Printf("Fetching data from: %s\n", url)
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Error getting stock data: %s", err)
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
			return errors.New("Error parsing data")
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
			return fmt.Errorf("Error updating investment price: %s", err)
		}

		totalRowsAdded++
	}

	err = scanner.Err()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rows added: %d\n", totalRowsAdded)
	return nil
}
