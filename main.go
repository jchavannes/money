package main

import (
	"github.com/jchavannes/money/web"
	"github.com/jchavannes/money/object/price"
	"github.com/spf13/cobra"
	"errors"
	"strconv"
	"fmt"
	"log"
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
				fmt.Println(err)
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
				return fmt.Errorf("Error parsing userId: %s", err)
			}
			err = price.UpdateForUser(uint(userId))
			if err != nil {
				fmt.Println(err)
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
