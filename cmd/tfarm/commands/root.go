package commands

import (
	"log"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/cmd/tfarm/commands/ranch"
	"github.com/cbodonnell/tfarm/cmd/tfarm/commands/server"
	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/version"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "tfarm",
		Short:   "tfarm - a CLI for creating and managing tunnels",
		Version: version.Version,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.AddCommand(ConfigureCmd())
	rootCmd.AddCommand(CreateCmd())
	rootCmd.AddCommand(DeleteCmd())
	rootCmd.AddCommand(InfoCmd())
	rootCmd.AddCommand(ReloadCmd())
	rootCmd.AddCommand(RestartCmd())
	rootCmd.AddCommand(StatusCmd())
	rootCmd.AddCommand(VerifyCmd())

	// add the server subcommand
	rootCmd.AddCommand(server.RootCmd())

	// add the ranch subcommand
	rootCmd.AddCommand(ranch.RootCmd())

	return rootCmd
}

func InitAndExecute() {

	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func getClient() (*api.APIClient, error) {
	endpoint := os.Getenv("TFARM_API_ENDPOINT")
	if endpoint == "" {
		endpoint = api.DefaultEndpoint
	}

	// TODO: make this configurable through a config file (like ~/.tfarm/tls/)
	configDir := os.Getenv("TFARM_CONFIG_DIR")
	if configDir == "" {
		// get the user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting user's home directory: %s", err)
		}

		configDir = path.Join(home, ".tfarm")
	}

	client, err := api.NewClient(endpoint, configDir)
	if err != nil {
		log.Fatalf("error creating API client: %s", err)
	}

	return client, nil
}
