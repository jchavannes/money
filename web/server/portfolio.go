package server

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
	"github.com/jchavannes/money/app/portfolio"
	"net/http"
	"sort"
)

var portfolioGetRoute = web.Route{
	Pattern:     UrlPortfolioGet,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
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
