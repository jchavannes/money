package cmd

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/jchavannes/money/app/price"
	"github.com/jchavannes/money/web/server"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

const FlagInsecure = "insecure"

var (
	moneyCmd = &cobra.Command{
		Use:   "money",
		Short: "Run money application",
	}

	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Main browser application",
		RunE: func(c *cobra.Command, args []string) error {
			sessionCookieInsecure, _ := c.Flags().GetBool(FlagInsecure)
			err := server.RunWeb(sessionCookieInsecure)
			if err != nil {
				return jerr.Get("error running web", err)
			}
			return nil
		},
	}

	updateCmd = &cobra.Command{
		Use:   "update [userId]",
		Short: "Update stock and currency data",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return jerr.New("must specify a userId.")
			}
			userId, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return jerr.Get("error parsing userId: "+args[0], err)
			}
			jlog.Logf("Updating for user %d\n", userId)
			err = price.UpdateForUser(uint(userId))
			if err != nil {
				return jerr.Get("error updating user", err)
			}
			return nil
		},
	}
)

func init() {
	webCmd.Flags().Bool(FlagInsecure, false, "Allow session cookie over unencrypted HTTP")
}

func Execute() {
	moneyCmd.AddCommand(webCmd)
	moneyCmd.AddCommand(updateCmd)
	err := moneyCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
