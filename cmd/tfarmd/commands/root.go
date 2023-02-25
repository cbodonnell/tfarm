package commands

import (
	"os"

	"github.com/cbodonnell/tfarm/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "tfarmd",
	Short:   "tfarmd - a daemon to manage frpc tunnels",
	Version: version.Version,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func InitAndExecute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
