package server

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
)

var dashboardRoute = web.Route{
	Pattern: UrlDashboard,
	Handler: func(r *web.Response) {
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(getUrlWithBaseUrl(UrlIndex, r))
			return
		}
		r.Render()
	},
}
