package ranch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cbodonnell/oauth2utils/pkg/utils"
	"github.com/cbodonnell/tfarm/pkg/ranch/api"
	"github.com/cbodonnell/tfarm/pkg/ranch/auth"
	"github.com/spf13/cobra"
)

func ClientsListCmd(tokenDir, endpoint string) *cobra.Command {
	clientsListCmd := &cobra.Command{
		Use:           "list",
		Short:         "List ranch clients",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ClientsList(tokenDir, endpoint)
		},
	}

	return clientsListCmd
}

func ClientsList(tokenDir, endpoint string) error {
	ctx := context.Background()
	oc, err := auth.NewOIDCClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating OIDC client: %s", err)
	}

	token := utils.TryGetToken(ctx, oc, tokenDir)
	if !token.Valid() {
		return fmt.Errorf("not logged in")
	}

	apiClient := api.NewClient(oc.HTTPClient(ctx, token), endpoint)
	clients, err := apiClient.ListClients(&api.APIRequestParams{})
	if err != nil {
		return fmt.Errorf("error listing clients: %s", err)
	}

	b, err := json.Marshal(clients)
	if err != nil {
		return fmt.Errorf("error marshaling clients: %s", err)
	}

	fmt.Print(string(b))

	return nil
}
