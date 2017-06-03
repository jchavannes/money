package web

import (
	"github.com/jchavannes/jgo/web"
)

var indexRoute = web.Route{
	Pattern: URL_INDEX,
	Handler: func(r *web.Response) {
		r.Render()
	},
}
