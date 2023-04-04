package ranch

import (
	"github.com/cbodonnell/tfarm/pkg/ranch/auth"
	"github.com/spf13/cobra"
)

func ClientsCmd(tokenDir, endpoint string, oidcConfig *auth.OIDCClientConfig) *cobra.Command {
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

	clientsCmd.AddCommand(ClientsCreateCmd(tokenDir, endpoint, oidcConfig))
	clientsCmd.AddCommand(ClientsDeleteCmd(tokenDir, endpoint, oidcConfig))
	clientsCmd.AddCommand(ClientsGetCmd(tokenDir, endpoint, oidcConfig))
	clientsCmd.AddCommand(ClientsListCmd(tokenDir, endpoint, oidcConfig))

	return clientsCmd
}
