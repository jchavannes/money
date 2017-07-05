package web

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/domain/auth"
	"net/http"
	"github.com/jchavannes/money/domain/chart"
)

const (
	FORM_INPUT_SYMBOL = "symbol"
	FORM_INPUT_MARKET = "market"
)

var individualChartsRoute = web.Route{
	Pattern: URL_INDIVIDUAL,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

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

var individualChartGetRoute = web.Route{
	Pattern: URL_INDIVIDUAL_CHART_GET,
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
		symbol := r.Request.GetFormValue(FORM_INPUT_SYMBOL)
		market := r.Request.GetFormValue(FORM_INPUT_MARKET)

		chartData, err := chart.GetIndividualChartData(user.Id, symbol, market)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		r.WriteJson(chartData, false)
	},
}
