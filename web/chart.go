package web

import (
	"github.com/jchavannes/jgo/web"
	"git.jasonc.me/main/money/db/auth"
	"net/http"
	"git.jasonc.me/main/money/db/price"
)

var chartGetRoute = web.Route{
	Pattern: URL_CHART_GET,
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
		overallChartData, err := price.GetOverallChartData(user.Id)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		r.WriteJson(overallChartData, false)
	},
}

var chartRoute = web.Route{
	Pattern: URL_CHART,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetResponseCode(http.StatusUnauthorized)
			return
		}
		r.Render()
	},
}
