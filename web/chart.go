package web

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/domain/auth"
	"net/http"
	"github.com/jchavannes/money/domain/chart"
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
		overallChartData, err := chart.GetOverallChartData(user.Id)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		r.WriteJson(overallChartData, false)
	},
}
