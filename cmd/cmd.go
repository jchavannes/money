package cmd

import (
	"github.com/spf13/cobra"
)

var (
	moneyCmd = &cobra.Command{
		Use:   "iiproject",
		Short: "Run iiProject application",
	}

	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Main browser application",
		RunE: func(c *cobra.Command, args []string) error {
			return CmdWeb()
		},
	}
)

func Execute() error {
	moneyCmd.AddCommand(webCmd)
	return moneyCmd.Execute()
}
