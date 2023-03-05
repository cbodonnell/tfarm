package commands

import (
	"log"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/cmd/tfarm/commands/server"
	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/version"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "tfarm",
		Short:   "tfarm - a CLI to interact with the tfarmd daemon",
		Version: version.Version,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.AddCommand(CreateCmd())
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(verifyCmd)

	// add the server subcommand
	rootCmd.AddCommand(server.RootCmd())

	return rootCmd
}

// client is the API client used by all commands
var client *api.APIClient
var configDir string

func InitAndExecute() {
	endpoint := os.Getenv("TFARM_API_ENDPOINT")
	if endpoint == "" {
		endpoint = api.DefaultEndpoint
	}

	// TODO: make this configurable through a config file (like ~/.tfarm/tls/)
	configDir = os.Getenv("TFARM_CONFIG_DIR")
	if configDir == "" {
		// get the user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting user's home directory: %s", err)
		}

		configDir = path.Join(home, ".tfarm")
	}

	newClient, err := api.NewClient(endpoint, configDir)
	if err != nil {
		log.Fatalf("error creating API client: %s", err)
	}

	client = newClient

	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
