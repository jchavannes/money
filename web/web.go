package web

import (
	"github.com/jchavannes/jgo/web"
	"git.jasonc.me/main/money/db/auth"
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
	URL_INVESTMENT_TRANSACTIONS_GET = "/investment-transactions-get"
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

func RunWeb() error {
	server := web.Server{
		NotFoundHandler: notFoundHandler,
		Port: 8247,
		UseSessions: true,
		TemplatesDir: "templates",
		StaticFilesDir: "pub",
		PreHandler: preHandler,
		Routes: []web.Route{
			indexRoute,
			dashboardRoute,
			signupRoute,
			signupSubmitRoute,
			loginRoute,
			loginSubmitRoute,
			logoutRoute,
			investmentTransactionsGetRoute,
		},
	}
	return server.Run()
}
