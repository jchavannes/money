package server

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
	"github.com/jchavannes/money/app/db"
	"github.com/jchavannes/money/app/price"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	FormInputInvestmentType      = "type"
	FormInputInvestmentSymbol    = "symbol"
	FormInputTransactionId       = "transactionId"
	FormInputTransactionDate     = "date"
	FormInputTransactionPrice    = "price"
	FormInputTransactionQuantity = "quantity"
	FormInputTransactionType     = "transaction-type"
	FormInputInvestmentId        = "investmentId"
)

var investmentUpdateRoute = web.Route{
	Pattern:     UrlInvestmentUpdate,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		investmentId := r.Request.GetFormValueInt(FormInputInvestmentId)
		err := price.UpdateInvestmentById(uint(investmentId))
		if err != nil {
			r.Error(jerr.Get("Error updating investment", err), http.StatusInternalServerError)
		}
	},
}

var investmentUpdateAllRoute = web.Route{
	Pattern:     UrlInvestmentUpdateAll,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
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
	Pattern:     UrlInvestmentTransactionsGet,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
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
	Pattern:     UrlInvestmentSymbolsGet,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		investmentType := r.Request.GetFormValue(FormInputInvestmentType)
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
	Pattern:     UrlInvestmentTransactionAdd,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		investmentType := r.Request.GetFormValue(FormInputInvestmentType)
		investmentSymbol := r.Request.GetFormValue(FormInputInvestmentSymbol)

		transactionDateString := r.Request.GetFormValue(FormInputTransactionDate)
		layout := "01/02/2006"
		transactionDate, err := time.Parse(layout, transactionDateString)

		transactionPriceString := r.Request.GetFormValue(FormInputTransactionPrice)
		transactionPrice, err := strconv.ParseFloat(transactionPriceString, 32)
		if err != nil {
			r.Error(jerr.Get("Error converting price string", err), http.StatusUnprocessableEntity)
			return
		}

		transactionQuantityString := r.Request.GetFormValue(FormInputTransactionQuantity)
		transactionQuantity, err := strconv.ParseFloat(transactionQuantityString, 32)
		if err != nil {
			r.Error(jerr.Get("Error converting quantity string", err), http.StatusUnprocessableEntity)
			return
		}

		transactionTypeString := r.Request.GetFormValue(FormInputTransactionType)
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
	Pattern:     UrlInvestmentTransactionDelete,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}

		transactionIdString := r.Request.GetFormValue(FormInputTransactionId)
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
