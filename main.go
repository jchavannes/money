package main

import (
	"errors"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/object/price"
	"github.com/jchavannes/money/web"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

const FLAG_INSECURE = "insecure"

var (
	moneyCmd = &cobra.Command{
		Use:   "money",
		Short: "Run money application",
	}

	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Main browser application",
		RunE: func(c *cobra.Command, args []string) error {
			sessionCookieInsecure, _ := c.Flags().GetBool(FLAG_INSECURE)
			err := web.RunWeb(sessionCookieInsecure)
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
				return jerr.Get("Error parsing userId: "+args[0], err)
			}
			err = price.UpdateForUser(uint(userId))
			if err != nil {
				return jerr.Get("Error updating user", err)
			}
			return nil
		},
	}
)

func init() {
	webCmd.Flags().Bool(FLAG_INSECURE, false, "Allow session cookie over unencrypted HTTP")
}

func main() {
	moneyCmd.AddCommand(webCmd)
	moneyCmd.AddCommand(updateCmd)
	err := moneyCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
