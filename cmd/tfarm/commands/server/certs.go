package server

import (
	"github.com/spf13/cobra"
)

func CertsCmd() *cobra.Command {
	certsCmd := &cobra.Command{
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

	certsCmd.AddCommand(CertsRegenerateCmd())
	certsCmd.AddCommand(CertsClientCmd())

	return certsCmd
}
