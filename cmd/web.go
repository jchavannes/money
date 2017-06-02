package cmd

import (
	"github.com/jchavannes/jgo/web"
)

const port = 8247

var (
	indexRoute = web.Route{
		Pattern: "/",
		Handler: func(r *web.Response) {
			r.Helper["CsrfToken"] = r.Session.GetCsrfToken()
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

func CmdWeb() error {
	server := web.Server{
		Port: port,
		UseSessions: true,
		TemplatesDir: "templates",
		StaticFilesDir: "public",
		Routes: []web.Route{
			indexRoute,
			postRoute,
		},
	}
	return server.Run()
}
