package cmd

import (
	"github.com/jchavannes/money/web"
)

func CmdWeb() error {
	return web.RunWeb()
}
