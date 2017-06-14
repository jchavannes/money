package web

import (
	"github.com/jchavannes/jgo/web"
	"git.jasonc.me/main/money/db/auth"
	"git.jasonc.me/main/money/db/portfolio"
	"net/http"
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
		r.WriteJson(userPortfolio, false)
	},
}
