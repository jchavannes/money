package cmd

import (
	"git.jasonc.me/main/money/db/price"
)

func CmdUpdate(userId uint) error {
	return price.UpdateForUser(userId)
}
