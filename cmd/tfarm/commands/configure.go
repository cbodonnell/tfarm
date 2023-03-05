package commands

import (
	"errors"
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/term"
	"github.com/spf13/cobra"
)

func ConfigureCmd() *cobra.Command {
	var clientID string
	var clientSecret string

	configureCmd := &cobra.Command{
		Use:           "configure",
		Short:         "Configure tfarm server with your ranch credentials",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Configure(clientID, clientSecret)
		},
	}

	configureCmd.Flags().StringVar(&clientID, "client-id", "", "ranch client id")
	configureCmd.Flags().StringVar(&clientSecret, "client-secret", "", "ranch client secret")

	return configureCmd
}

func Configure(clientID, clientSecret string) error {
	if clientID == "" {
		clientID = term.StringPrompt("Client ID:")
	}

	if clientID == "" {
		return errors.New("client id is required")
	}

	if clientSecret == "" {
		clientSecret = term.PasswordPrompt("Client Secret:")
	}

	if clientSecret == "" {
		return errors.New("client secret is required")
	}

	req := &api.ConfigureRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	status, err := client.Configure(req)
	if err != nil {
		return fmt.Errorf("error configuring: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
