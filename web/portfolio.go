package web

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/domain/auth"
	"github.com/jchavannes/money/object/portfolio"
	"net/http"
	"sort"
)

var portfolioGetRoute = web.Route{
	Pattern: URL_PORTFOLIO_GET,
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
		userPortfolio, err := portfolio.Get(user.Id)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		sort.Sort(portfolio.PortfolioItemSorter(userPortfolio.Items))
		r.WriteJson(userPortfolio, false)
	},
}
