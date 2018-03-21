package server

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
	"net/http"
)

const (
	UrlIndex        = "/"
	UrlDashboard    = "/dashboard"
	UrlSignup       = "/signup"
	UrlSignupSubmit = "/signup-submit"
	UrlLogin        = "/login"
	UrlLoginSubmit  = "/login-submit"
	UrlLogout       = "/logout"
)

const (
	UrlChartGet           = "/chart-get"
	UrlIndividual         = "/individual"
	UrlIndividualChartGet = "/individual-chart-get"
	UrlPortfolioGet       = "/portfolio-get"
)

const (
	UrlInvestmentUpdate            = "/investment-update"
	UrlInvestmentUpdateAll         = "/investment-update-all"
	UrlInvestmentTransactionsGet   = "/investment-transactions-get"
	UrlInvestmentTransactionAdd    = "/investment-transaction-add"
	UrlInvestmentTransactionDelete = "/investment-transaction-delete"
	UrlInvestmentSymbolsGet        = "/investment-symbols-get"
)

func isLoggedIn(r *web.Response) bool {
	if ! auth.IsLoggedIn(r.Session.CookieId) {
		r.SetResponseCode(http.StatusUnauthorized)
		return false
	}
	return true
}

func preHandler(r *web.Response) {
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

func notFoundHandler(r *web.Response) {
	r.SetResponseCode(http.StatusNotFound)
	r.RenderTemplate("404")
}

func getBaseUrl(r *web.Response) string {
	baseUrl := r.Request.GetHeader("AppPath")
	if baseUrl == "" {
		baseUrl = "/"
	}
	return baseUrl
}

func getUrlWithBaseUrl(url string, r *web.Response) string {
	baseUrl := getBaseUrl(r)
	baseUrl = baseUrl[:len(baseUrl)-1]
	return baseUrl + url
}

func RunWeb(sessionCookieInsecure bool) error {
	server := web.Server{
		CookiePrefix:    "money_tracker",
		IsLoggedIn:      isLoggedIn,
		NotFoundHandler: notFoundHandler,
		Port:            8247,
		UseSessions:     true,
		TemplatesDir:    "web/templates",
		StaticFilesDir:  "web/public",
		PreHandler:      preHandler,
		InsecureCookie:  sessionCookieInsecure,
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
