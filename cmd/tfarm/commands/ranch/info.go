package ranch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/info"
	"github.com/cbodonnell/tfarm/pkg/ranch/api"
	"github.com/spf13/cobra"
)

func InfoCmd(tokenDir, endpoint string) *cobra.Command {
	infoCmd := &cobra.Command{
		Use:           "info",
		Short:         "Print ranch info",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Info(tokenDir, endpoint)
		},
	}

	return infoCmd
}

func Info(tokenDir, endpoint string) error {
	info := getInfo(tokenDir, endpoint)
	b, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("error marshaling clients: %s", err)
	}

	fmt.Print(string(b))

	return nil
}

func getInfo(tokenDir, endpoint string) *info.Info {
	info := info.Info{
		TokenDir: tokenDir,
		Endpoint: endpoint,
	}

	apiClient := api.NewClient(http.DefaultClient, endpoint)
	res, err := apiClient.GetInfo(&api.APIRequestParams{})
	if err != nil {
		return &info
	}

	info.Ready = res.Ready
	info.Version = res.Version
	info.OIDC.Issuer = res.OIDC.Issuer
	info.OIDC.ClientID = res.OIDC.ClientID

	return &info
}
