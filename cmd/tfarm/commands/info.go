package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/term"
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

	info := client.Info()

	switch outputFormat {
	// TODO: make this yaml so it can be more dynamic
	case "text":
		fmt.Println("Client:")
		fmt.Println("  Version:", info.Client.Version)
		fmt.Println("  Config:", info.Client.Config)
		fmt.Println("Server:")
		fmt.Println("  Version:", info.Server.Version)
		fmt.Println("  Endpoint:", info.Server.Endpoint)
		if info.Server.Error != "" {
			fmt.Println("  Error:", info.Server.Error)
		}
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
