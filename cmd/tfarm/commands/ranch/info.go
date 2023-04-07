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
	apiClient := api.NewClient(http.DefaultClient, endpoint)
	res, err := apiClient.GetInfo(&api.APIRequestParams{})
	if err != nil {
		return fmt.Errorf("error listing clients: %s", err)
	}

	info := info.Info{
		Ready:    res.Ready,
		Version:  res.Version,
		TokenDir: tokenDir,
		Endpoint: endpoint,
		OIDC: info.OIDCInfo{
			Issuer:   res.OIDC.Issuer,
			ClientID: res.OIDC.ClientID,
		},
	}

	b, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("error marshaling clients: %s", err)
	}

	fmt.Print(string(b))

	return nil
}
