package web

import (
	"github.com/jchavannes/jgo/web"
	"git.jasonc.me/main/money/db/auth"
	"net/http"
	"git.jasonc.me/main/money/db/investment"
)

var investmentTransactionsGetRoute = web.Route{
	Pattern: URL_INVESTMENT_TRANSACTIONS_GET,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(getUrlWithBaseUrl(URL_INDEX, r))
			return
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		investmentTransactions, err := investment.GetInvestmentTransactionsForUser(user.Id)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		r.WriteJson(investmentTransactions, false)
	},
}
