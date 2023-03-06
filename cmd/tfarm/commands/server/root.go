package server

import (
	"os"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Interface with the tfarmd server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(certsCmd)

	return rootCmd
}
