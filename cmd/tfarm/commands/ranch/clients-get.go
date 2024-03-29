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

func ClientsGetCmd(tokenDir, endpoint string, oidcConfig *auth.OIDCClientConfig) *cobra.Command {
	var outCredentials bool

	clientsGetCmd := &cobra.Command{
		Use:           "get [id]",
		Short:         "Get a ranch client",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			return ClientsGet(tokenDir, endpoint, oidcConfig, args[0], outCredentials)
		},
	}

	clientsGetCmd.Flags().BoolVar(&outCredentials, "credentials", false, "output in credentials.json format")

	return clientsGetCmd
}

func ClientsGet(tokenDir, endpoint string, oidcConfig *auth.OIDCClientConfig, id string, outCredentials bool) error {
	if id == "" {
		return fmt.Errorf("client id is required")
	}

	ctx := context.Background()
	oc, err := auth.NewOIDCClient(ctx, oidcConfig)
	if err != nil {
		return fmt.Errorf("error creating OIDC client: %s", err)
	}

	token := utils.TryGetToken(ctx, oc, tokenDir)
	if !token.Valid() {
		return fmt.Errorf("not logged in")
	}

	apiClient := api.NewClient(oc.HTTPClient(ctx, token), endpoint)

	var b []byte
	if outCredentials {
		b, err = apiClient.GetClientCredentialsJson(&api.ClientRequestParams{
			ID: id,
		})
		if err != nil {
			return fmt.Errorf("error getting client credentials: %s", err)
		}
	} else {
		client, err := apiClient.GetClient(&api.ClientRequestParams{
			ID: id,
		})
		if err != nil {
			return fmt.Errorf("error getting client: %s", err)
		}

		b, err = json.Marshal(client)
		if err != nil {
			return fmt.Errorf("error marshaling client: %s", err)
		}
	}

	fmt.Print(string(b))

	return nil
}
