package ranch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cbodonnell/oauth2utils/pkg/utils"
	"github.com/cbodonnell/tfarm/pkg/ranch/api"
	"github.com/cbodonnell/tfarm/pkg/ranch/auth"
	"github.com/spf13/cobra"
)

func ClientsDeleteCmd(tokenDir string) *cobra.Command {
	clientsDeleteCmd := &cobra.Command{
		Use:           "delete [id]",
		Short:         "Delete a ranch client",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			return ClientsDelete(tokenDir, args[0])
		},
	}

	return clientsDeleteCmd
}

func ClientsDelete(tokenDir, id string) error {
	_, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid client id: %s", err)
	}

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
	client, err := apiClient.DeleteClient(&api.ClientRequestParams{
		ID: id,
	})
	if err != nil {
		return fmt.Errorf("error deleting client: %s", err)
	}

	b, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("error marshaling client: %s", err)
	}

	fmt.Print(string(b))

	return nil
}
