package ranch

import (
	"github.com/spf13/cobra"
)

func ClientsCmd(tokenDir, endpoint string) *cobra.Command {
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

	clientsCmd.AddCommand(ClientsCreateCmd(tokenDir, endpoint))
	clientsCmd.AddCommand(ClientsDeleteCmd(tokenDir, endpoint))
	clientsCmd.AddCommand(ClientsGetCmd(tokenDir, endpoint))
	clientsCmd.AddCommand(ClientsListCmd(tokenDir, endpoint))

	return clientsCmd
}
