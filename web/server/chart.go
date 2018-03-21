package server

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
	"github.com/jchavannes/money/app/chart"
	"net/http"
)

const (
	FormInputSymbol = "symbol"
	FormInputMarket = "market"
)

var individualChartsRoute = web.Route{
	Pattern:    UrlIndividual,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var chartGetRoute = web.Route{
	Pattern:     UrlChartGet,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
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
	Pattern:     UrlIndividualChartGet,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		symbol := r.Request.GetFormValue(FormInputSymbol)
		market := r.Request.GetFormValue(FormInputMarket)

		chartData, err := chart.GetIndividualChartData(user.Id, symbol, market)
		if err != nil {
			r.Error(err, http.StatusInternalServerError)
			return
		}
		r.WriteJson(chartData, false)
	},
}
