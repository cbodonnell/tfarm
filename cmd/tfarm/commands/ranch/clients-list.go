package ranch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cbodonnell/oauth2utils/pkg/utils"
	"github.com/cbodonnell/tfarm/pkg/ranch/api"
	"github.com/cbodonnell/tfarm/pkg/ranch/auth"
	"github.com/spf13/cobra"
)

func ClientsListCmd(tokenDir string) *cobra.Command {
	clientsListCmd := &cobra.Command{
		Use:           "list",
		Short:         "List ranch clients",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ClientsList(tokenDir)
		},
	}

	return clientsListCmd
}

func ClientsList(tokenDir string) error {
	ctx := context.Background()
	oc, err := auth.NewOIDCClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating OIDC client: %s", err)
	}

	token := utils.TryGetToken(ctx, oc, tokenDir)
	if !token.Valid() {
		return fmt.Errorf("not logged in")
	}

	// TODO: set this as a higher scope
	endpoint := os.Getenv("RANCH_API_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.tunnel.farm"
	}

	apiClient := api.NewClient(oc.HTTPClient(ctx, token), endpoint)
	clients, err := apiClient.ListClients(&api.APIRequestParams{})
	if err != nil {
		return fmt.Errorf("error listing clients: %s", err)
	}

	// marshal and print client
	b, err := json.Marshal(clients)
	if err != nil {
		return fmt.Errorf("error marshaling clients: %s", err)
	}

	fmt.Print(string(b))

	return nil
}
