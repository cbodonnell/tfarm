package ranch

import (
	"fmt"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/info"
	"github.com/cbodonnell/tfarm/pkg/ranch/api"
	"github.com/cbodonnell/tfarm/pkg/term"
	"github.com/spf13/cobra"
)

func InfoCmd(tokenDir, endpoint string) *cobra.Command {
	var outputFormat string

	infoCmd := &cobra.Command{
		Use:           "info",
		Short:         "Print ranch info",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Info(tokenDir, endpoint, outputFormat)
		},
	}

	infoCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json)")

	return infoCmd
}

func Info(tokenDir, endpoint, outputFormat string) error {
	info := getInfo(tokenDir, endpoint)

	switch outputFormat {
	// TODO: make this yaml so it can be more dynamic
	case "text":
		fmt.Println("Ready:", info.Ready)
		if info.Error != "" {
			fmt.Println("Error:", info.Error)
		}
		fmt.Println("Version:", info.Version)
		fmt.Println("Config:", info.TokenDir)
		fmt.Println("Endpoint:", info.Endpoint)
		fmt.Println("OIDC:")
		fmt.Println("  Issuer:", info.OIDC.Issuer)
		fmt.Println("  Client ID:", info.OIDC.ClientID)
	case "json":
		b, err := term.PrettyJSON(info)
		if err != nil {
			return fmt.Errorf("error marshaling info to json: %s", err)
		}
		fmt.Println(string(b))
	default:
		return fmt.Errorf("invalid output format: %s", outputFormat)
	}

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
		info.Error = err.Error()
		return &info
	}

	info.Ready = res.Ready
	info.Version = res.Version
	info.OIDC.Issuer = res.OIDC.Issuer
	info.OIDC.ClientID = res.OIDC.ClientID

	return &info
}
