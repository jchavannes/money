package cmd

import (
	"github.com/jchavannes/jgo/web"
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
)

func CmdWeb() error {
	server := web.Server{
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
