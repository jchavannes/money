package web

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/domain/auth"
	"net/http"
)

const (
	URL_INDEX = "/"
	URL_DASHBOARD = "/dashboard"
	URL_SIGNUP = "/signup"
	URL_SIGNUP_SUBMIT = "/signup-submit"
	URL_LOGIN = "/login"
	URL_LOGIN_SUBMIT = "/login-submit"
	URL_LOGOUT = "/logout"
	URL_CHART_GET = "/chart-get"
	URL_INDIVIDUAL = "/individual"
	URL_INDIVIDUAL_CHART_GET = "/individual-chart-get"
	URL_PORTFOLIO_GET = "/portfolio-get"
	URL_INVESTMENT_UPDATE = "/investment-update"
	URL_INVESTMENT_UPDATE_ALL = "/investment-update-all"
	URL_INVESTMENT_TRANSACTIONS_GET = "/investment-transactions-get"
	URL_INVESTMENT_TRANSACTION_ADD = "/investment-transaction-add"
	URL_INVESTMENT_TRANSACTION_DELETE = "/investment-transaction-delete"
	URL_INVESTMENT_SYMBOLS_GET = "/investment-symbols-get"
)

var (
	preHandler = func(r *web.Response) {
		r.Helper["BaseUrl"] = getBaseUrl(r)
		if auth.IsLoggedIn(r.Session.CookieId) {
			user, err := auth.GetSessionUser(r.Session.CookieId)
			if err != nil {
				r.Error(err, http.StatusUnprocessableEntity)
				return
			}
			r.Helper["Username"] = user.Username
		}
	}

	notFoundHandler = func(r *web.Response) {
		r.SetResponseCode(http.StatusNotFound)
		r.RenderTemplate("404")
	}

	getBaseUrl = func(r *web.Response) string {
		baseUrl := r.Request.GetHeader("AppPath")
		if baseUrl == "" {
			baseUrl = "/"
		}
		return baseUrl
	}

	getUrlWithBaseUrl = func(url string, r *web.Response) string {
		baseUrl := getBaseUrl(r)
		baseUrl = baseUrl[:len(baseUrl) - 1]
		return baseUrl + url
	}
)

func RunWeb(sessionCookieInsecure bool) error {
	server := web.Server{
		CookiePrefix: "money_tracker",
		NotFoundHandler: notFoundHandler,
		Port: 8247,
		UseSessions: true,
		TemplatesDir: "web/templates",
		StaticFilesDir: "web/pub",
		PreHandler: preHandler,
		InsecureCookie: sessionCookieInsecure,
		Routes: []web.Route{
			indexRoute,
			dashboardRoute,
			signupRoute,
			signupSubmitRoute,
			loginRoute,
			loginSubmitRoute,
			logoutRoute,
			individualChartsRoute,
			individualChartGetRoute,
			chartGetRoute,
			portfolioGetRoute,
			investmentUpdateRoute,
			investmentUpdateAllRoute,
			investmentTransactionsGetRoute,
			investmentSymbolsGetRoute,
			investmentTransactionAddRoute,
			investmentTransactionDeleteRoute,
		},
	}
	return server.Run()
}
