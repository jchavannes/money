package cmd

import (
	"github.com/jchavannes/jgo/web"
	"net/http"
)

const port = 8247

var (
	indexRoute = web.Route{
		Pattern: "/",
		Handler: func(r *web.Response) {
			r.Render()
		},
	}

	postRoute = web.Route{
		Pattern: "/post",
		CsrfProtect: true,
		Handler: func(r *web.Response) {
			r.Write("Posts")
		},
	}
)

var (
	preHandler = func(r *web.Response) {
		r.Helper["CsrfToken"] = r.Session.GetCsrfToken()
		r.Helper["BaseUrl"] = getBaseUrl(r)
	}

	getBaseUrl = func(r *web.Response) string {
		baseUrl := r.Request.GetHeader("AppPath")
		if baseUrl == "" {
			baseUrl = "/"
		}
		return baseUrl
	}

	notFoundHandler = func(r *web.Response) {
		r.SetResponseCode(http.StatusNotFound)
		r.RenderTemplate("404")
	}
)

func CmdWeb() error {
	server := web.Server{
		NotFoundHandler: notFoundHandler,
		Port: port,
		UseSessions: true,
		TemplatesDir: "templates",
		StaticFilesDir: "pub",
		PreHandler: preHandler,
		Routes: []web.Route{
			indexRoute,
			postRoute,
		},
	}
	return server.Run()
}
