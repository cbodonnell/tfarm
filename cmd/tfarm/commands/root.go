package commands

import (
	"log"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

var version = "dev"
var rootCmd = &cobra.Command{
	Use:     "tfarm",
	Short:   "tfarm - a CLI to interact with the tfarmd daemon",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

// client is the API client used by all commands
var client *api.APIClient

func InitAndExecute() {
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

	tlsFiles := &api.TLSFiles{
		CertFile: path.Join(configDir, "tls", "client.crt"),
		KeyFile:  path.Join(configDir, "tls", "client.key"),
		CAFile:   path.Join(configDir, "tls", "ca.crt"),
	}

	newClient, err := api.NewClient(endpoint, tlsFiles)
	if err != nil {
		log.Fatalf("error creating API client: %s", err)
	}

	client = newClient

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
