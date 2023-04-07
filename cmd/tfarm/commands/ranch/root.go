package ranch

import (
	"log"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/ranch/auth"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ranch",
		Short: "Interface with the ranch api",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	tokenDir := getRanchTokenDir()
	endpoint := getRanchAPIEndpoint()
	// TODO: Make this discover the OIDC config from the ranch API
	oidcConfig := getOIDCConfig()

	rootCmd.AddCommand(InfoCmd(tokenDir, endpoint))
	rootCmd.AddCommand(ClientsCmd(tokenDir, endpoint, oidcConfig))
	rootCmd.AddCommand(LoginCmd(tokenDir, oidcConfig))
	rootCmd.AddCommand(LogoutCmd(tokenDir))

	return rootCmd
}

func getRanchTokenDir() string {
	configDir := os.Getenv("TFARM_CONFIG_DIR")
	if configDir == "" {
		// get the user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting user's home directory: %s", err)
		}

		configDir = path.Join(home, ".tfarm")
	}

	return path.Join(configDir, "ranch")
}

func getRanchAPIEndpoint() string {
	endpoint := os.Getenv("RANCH_API_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.tunnel.farm"
	}

	return endpoint
}

func getOIDCConfig() *auth.OIDCClientConfig {
	ranchOauthIssuer := os.Getenv("RANCH_OIDC_ISSUER")
	if ranchOauthIssuer == "" {
		ranchOauthIssuer = "https://auth.tunnel.farm/realms/tunnel.farm"
	}

	ranchOauthClientID := os.Getenv("RANCH_OIDC_CLIENT_ID")
	if ranchOauthClientID == "" {
		ranchOauthClientID = "tfarm-cli"
	}

	return &auth.OIDCClientConfig{
		Issuer:   ranchOauthIssuer,
		ClientID: ranchOauthClientID,
	}
}
