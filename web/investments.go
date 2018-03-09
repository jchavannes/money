package web

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/domain/auth"
	"net/http"
	"github.com/jchavannes/money/data/db"
	"time"
	"strconv"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/object/price"
	"sort"
)

const (
	FORM_INPUT_INVESTMENT_TYPE = "type"
	FORM_INPUT_INVESTMENT_SYMBOL = "symbol"
	FORM_INPUT_TRANSACTION_ID = "transactionId"
	FORM_INPUT_TRANSACTION_DATE = "date"
	FORM_INPUT_TRANSACTION_PRICE = "price"
	FORM_INPUT_TRANSACTION_QUANTITY = "quantity"
	FORM_INPUT_TRANSACTION_TYPE = "transaction-type"
	FORM_INPUT_INVESTMENT_ID = "investmentId"
)

var investmentUpdateRoute = web.Route{
	Pattern: URL_INVESTMENT_UPDATE,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		investmentId := r.Request.GetFormValueInt(FORM_INPUT_INVESTMENT_ID)
		err := price.UpdateInvestmentById(uint(investmentId))
		if err != nil {
			r.Error(jerr.Get("Error updating investment", err), http.StatusInternalServerError)
		}
	},
}

var investmentUpdateAllRoute = web.Route{
	Pattern: URL_INVESTMENT_UPDATE_ALL,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		err = price.UpdateForUser(user.Id)
		if err != nil {
			r.Error(jerr.Get("Error updating user investments", err), http.StatusInternalServerError)
		}
	},
}

var investmentTransactionsGetRoute = web.Route{
	Pattern: URL_INVESTMENT_TRANSACTIONS_GET,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		investmentTransactions, err := db.GetInvestmentTransactionsForUser(user.Id)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		sort.Sort(db.InvestmentTransactionSorter(investmentTransactions))
		r.WriteJson(investmentTransactions, false)
	},
}

var investmentSymbolsGetRoute = web.Route{
	Pattern: URL_INVESTMENT_SYMBOLS_GET,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		investmentType := r.Request.GetFormValue(FORM_INPUT_INVESTMENT_TYPE)
		investments, err := db.GetInvestmentsForType(investmentType)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		var tags []string
		for _, inv := range investments {
			tags = append(tags, inv.Symbol)
		}
		if investmentType == db.InvestmentType_Crypto.String() {
			tags = append(tags, "bitcoin", "ethereum", "litecoin")
		}
		r.WriteJson(tags, false)
	},
}

var investmentTransactionAddRoute = web.Route{
	Pattern: URL_INVESTMENT_TRANSACTION_ADD,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		investmentType := r.Request.GetFormValue(FORM_INPUT_INVESTMENT_TYPE)
		investmentSymbol := r.Request.GetFormValue(FORM_INPUT_INVESTMENT_SYMBOL)

		transactionDateString := r.Request.GetFormValue(FORM_INPUT_TRANSACTION_DATE)
		layout := "01/02/2006"
		transactionDate, err := time.Parse(layout, transactionDateString)

		transactionPriceString := r.Request.GetFormValue(FORM_INPUT_TRANSACTION_PRICE)
		transactionPrice, err := strconv.ParseFloat(transactionPriceString, 32)
		if err != nil {
			r.Error(jerr.Get("Error converting price string", err), http.StatusUnprocessableEntity)
			return
		}

		transactionQuantityString := r.Request.GetFormValue(FORM_INPUT_TRANSACTION_QUANTITY)
		transactionQuantity, err := strconv.ParseFloat(transactionQuantityString, 32)
		if err != nil {
			r.Error(jerr.Get("Error converting quantity string", err), http.StatusUnprocessableEntity)
			return
		}

		transactionTypeString := r.Request.GetFormValue(FORM_INPUT_TRANSACTION_TYPE)
		var transactionType db.InvestmentTransactionType
		if transactionTypeString == "buy" {
			transactionType = db.InvestmentTransactionType_Buy
		} else {
			transactionType = db.InvestmentTransactionType_Sell
		}

		transactionInvestment, err := db.GetInvestment(investmentType, investmentSymbol)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}

		err = price.UpdateInvestmentById(transactionInvestment.Id)
		if err != nil {
			r.Error(jerr.Get("Error updating investment", err), http.StatusUnprocessableEntity)
			r.Write("Unable to update investment.")
			return
		}

		err = db.AddInvestmentTransaction(
			user.Id,
			transactionInvestment,
			transactionType,
			transactionDate,
			float32(transactionPrice),
			float32(transactionQuantity),
		)

		if err != nil {
			r.Error(err, http.StatusUnprocessableEntity)
			return
		}
	},
}

var investmentTransactionDeleteRoute = web.Route{
	Pattern: URL_INVESTMENT_TRANSACTION_DELETE,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}

		transactionIdString := r.Request.GetFormValue(FORM_INPUT_TRANSACTION_ID)
		transactionId, err := strconv.Atoi(transactionIdString)
		if err != nil {
			r.Error(err, http.StatusUnprocessableEntity)
			return
		}

		err = db.DeleteInvestmentTransaction(
			user.Id,
			uint(transactionId),
		)
		if err != nil {
			r.Error(err, http.StatusUnprocessableEntity)
			return
		}
	},
}
