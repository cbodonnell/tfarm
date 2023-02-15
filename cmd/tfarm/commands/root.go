package commands

import (
	"os"

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
	client = api.NewClient(endpoint)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
