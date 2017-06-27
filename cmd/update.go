package cmd

import (
	"github.com/jchavannes/money/db/price"
)

func CmdUpdate(userId uint) error {
	return price.UpdateForUser(userId)
}
