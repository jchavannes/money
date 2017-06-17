package cmd

import (
	"github.com/spf13/cobra"
	"errors"
	"strconv"
	"fmt"
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
			err := CmdWeb()
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
			err = CmdUpdate(uint(userId))
			if err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}
)

func Execute() error {
	moneyCmd.AddCommand(webCmd)
	moneyCmd.AddCommand(updateCmd)
	return moneyCmd.Execute()
}
