package cmd

import (
	"git.jasonc.me/main/money/db/investment"
)

func CmdUpdate(userId uint) error {
	return investment.UpdateForUser(userId)
}
