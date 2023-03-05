package ranch

import (
	"github.com/spf13/cobra"
)

func ClientsCmd(tokenDir string) *cobra.Command {
	clientsCmd := &cobra.Command{
		Use:           "clients",
		Short:         "Manage ranch clients",
		SilenceUsage:  true,
		SilenceErrors: false,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}

	clientsCmd.AddCommand(ClientsCreateCmd(tokenDir))
	clientsCmd.AddCommand(ClientsDeleteCmd(tokenDir))
	clientsCmd.AddCommand(ClientsGetCmd(tokenDir))
	clientsCmd.AddCommand(ClientsListCmd(tokenDir))

	return clientsCmd
}
