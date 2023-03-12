package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/term"
	"github.com/spf13/cobra"
)

func ConfigureCmd() *cobra.Command {
	var clientID string
	var clientSecret string
	var credentialsStdin bool

	configureCmd := &cobra.Command{
		Use:           "configure",
		Short:         "Configure tfarm server with your ranch credentials",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Configure(clientID, clientSecret, credentialsStdin)
		},
	}

	configureCmd.Flags().StringVar(&clientID, "client-id", "", "ranch client id")
	configureCmd.Flags().StringVar(&clientSecret, "client-secret", "", "ranch client secret")
	configureCmd.Flags().BoolVar(&credentialsStdin, "credentials-stdin", false, "read credentials from stdin")

	return configureCmd
}

func Configure(clientID, clientSecret string, credentialsStdin bool) error {
	credentials := &auth.ConfigureCredentials{}

	if credentialsStdin {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading input: %s", err)
		}

		if err := json.Unmarshal(input, credentials); err != nil {
			return fmt.Errorf("error unmarshaling input: %s", err)
		}
	} else {
		if clientID == "" {
			clientID = term.StringPrompt("Client ID:")
		}

		if clientSecret == "" {
			clientSecret = term.PasswordPrompt("Client Secret:")
		}

		credentials.ClientID = clientID
		credentials.ClientSecret = clientSecret
	}

	if credentials.ClientID == "" {
		return errors.New("client id is required")
	}

	if credentials.ClientSecret == "" {
		return errors.New("client secret is required")
	}

	status, err := client.Configure(credentials)
	if err != nil {
		return fmt.Errorf("error creating: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
