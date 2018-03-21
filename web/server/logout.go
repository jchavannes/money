package server

import (
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/money/app/auth"
	"net/http"
)

var logoutRoute = web.Route{
	Pattern: UrlLogout,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			err := auth.Logout(r.Session.CookieId)
			if err != nil {
				r.Error(err, http.StatusInternalServerError)
				return
			}
		}
		r.SetRedirect(getUrlWithBaseUrl(UrlIndex, r))
	},
}
