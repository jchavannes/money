package main

import (
	"fmt"
	"github.com/jchavannes/jgo/web"
)

const port = 8247

func main() {
	server := web.Server{
		Port: port,
		TemplateDirectory: "templates",
		StaticDirectory: "public",
		Routes: []web.Route{{
			Pattern: "/post",
			CsrfProtect: true,
			Handler: func(r *web.Request) {
				r.Write("Posts")
			},
		}},
	}
	fmt.Printf("Starting money web server on port %d\n", port)
	server.Run()
}
