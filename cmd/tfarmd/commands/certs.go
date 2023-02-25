package commands

import (
	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:           "certs",
	Short:         "Manage TLS certificates",
	SilenceUsage:  true,
	SilenceErrors: false,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(certsCmd)
}
