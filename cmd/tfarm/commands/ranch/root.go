package ranch

import (
	"log"
	"os"
	"path"

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

	rootCmd.AddCommand(ClientsCmd(tokenDir, endpoint))
	rootCmd.AddCommand(LoginCmd(tokenDir))
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
