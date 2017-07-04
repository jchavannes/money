package main

import (
	"github.com/jchavannes/money/web"
	"github.com/jchavannes/money/object/price"
	"github.com/spf13/cobra"
	"errors"
	"strconv"
	"log"
	"github.com/jchavannes/jgo/jerr"
)

var (
	moneyCmd = &cobra.Command{
		Use:   "money",
		Short: "Run money application",
	}

	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Main browser application",
		RunE: func(c *cobra.Command, args []string) error {
			err := web.RunWeb()
			if err != nil {
				return jerr.Get("Error running web", err)
			}
			return nil
		},
	}

	updateCmd = &cobra.Command{
		Use:   "update [userId]",
		Short: "Update stock and currency data",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Must specify a userId.")
			}
			userId, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return jerr.Get("Error parsing userId: " + args[0], err)
			}
			err = price.UpdateForUser(uint(userId))
			if err != nil {
				return jerr.Get("Error updating user", err)
			}
			return nil
		},
	}
)

func main() {
	moneyCmd.AddCommand(webCmd)
	moneyCmd.AddCommand(updateCmd)
	err := moneyCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
