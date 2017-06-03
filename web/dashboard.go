package web

import (
	"github.com/jchavannes/jgo/web"
	"git.jasonc.me/main/money/db/auth"
)

var dashboardRoute = web.Route{
	Pattern: URL_DASHBOARD,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(getUrlWithBaseUrl(URL_INDEX, r))
			return
		}
		r.Render()
	},
}
