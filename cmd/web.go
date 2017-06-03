package cmd

import (
	"git.jasonc.me/main/money/web"
)

func CmdWeb() error {
	return web.RunWeb()
}
