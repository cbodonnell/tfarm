package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func InfoCmd() *cobra.Command {
	var outputFormat string

	infoCmd := &cobra.Command{
		Use:           "info",
		Short:         "Print information about the daemon",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Info(outputFormat)
		},
	}

	infoCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json)")

	return infoCmd
}

func Info(outputFormat string) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("error creating client: %s", err)
	}

	info, err := client.Info()
	if err != nil {
		return fmt.Errorf("error getting info: %s", err)
	}

	switch outputFormat {
	case "text":
		fmt.Println("Client Version:")
		fmt.Println("  ", info.ClientVersion)
		fmt.Println("Server Version:")
		fmt.Println("  ", info.ServerVersion)
		fmt.Println("Server Endpoint:")
		fmt.Println("  ", info.ServerEndpoint)
		fmt.Println("Configuration:")
		fmt.Println("  ", info.ConfigDir)
	case "json":
		b, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("error marshaling info to json: %s", err)
		}
		fmt.Println(string(b))
	default:
		return fmt.Errorf("invalid output format: %s", outputFormat)
	}

	return nil
}
