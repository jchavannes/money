package server

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
)

var indexRoute = web.Route{
	Pattern: UrlIndex,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(getUrlWithBaseUrl(UrlDashboard, r))
			return
		}
		r.Render()
	},
}
