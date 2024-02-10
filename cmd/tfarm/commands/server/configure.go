package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/term"
	"github.com/spf13/cobra"
)

func ConfigureCmd() *cobra.Command {
	var clientID string
	var clientSecret string
	var clientCACert string
	var clientTLSCert string
	var clientTLSKey string
	var credentialsStdin bool

	configureCmd := &cobra.Command{
		Use:           "configure",
		Short:         "Configure tfarm server with your ranch credentials",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Configure(clientID, clientSecret, clientCACert, clientTLSCert, clientTLSKey, credentialsStdin)
		},
	}

	configureCmd.Flags().StringVar(&clientID, "client-id", "", "ranch client id")
	configureCmd.Flags().StringVar(&clientSecret, "client-secret", "", "ranch client secret")
	configureCmd.Flags().StringVar(&clientCACert, "client-ca-cert", "", "base64 encoded ranch client ca cert")
	configureCmd.Flags().StringVar(&clientTLSCert, "client-tls-cert", "", "base64 encoded ranch client tls cert")
	configureCmd.Flags().StringVar(&clientTLSKey, "client-tls-key", "", "base64 encoded ranch client tls key")
	configureCmd.Flags().BoolVar(&credentialsStdin, "credentials-stdin", false, "read credentials from stdin")

	return configureCmd
}

func Configure(clientID, clientSecret, clientCACert, clientTLSCert, clientTLSKey string, credentialsStdin bool) error {
	workDir := os.Getenv("TFARMD_WORK_DIR")
	if workDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %s", err)
		}
		workDir = pwd
	}

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

		if clientCACert == "" {
			clientCACert = term.StringPrompt("Client CA Cert (base64 encoded):")
		}

		if clientTLSCert == "" {
			clientTLSCert = term.StringPrompt("Client TLS Cert (base64 encoded):")
		}

		if clientTLSKey == "" {
			clientTLSKey = term.PasswordPrompt("Client TLS Key (base64 encoded):")
		}

		credentials.ClientID = clientID
		credentials.ClientSecret = clientSecret
		credentials.ClientCACert = clientCACert
		credentials.ClientTLSCert = clientTLSCert
		credentials.ClientTLSKey = clientTLSKey
	}

	if credentials.ClientID == "" {
		return errors.New("client id is required")
	}

	if credentials.ClientSecret == "" {
		return errors.New("client secret is required")
	}

	if credentials.ClientCACert == "" {
		return errors.New("client ca cert is required")
	}

	if credentials.ClientTLSCert == "" {
		return errors.New("client tls cert is required")
	}

	if credentials.ClientTLSKey == "" {
		return errors.New("client tls key is required")
	}

	credPath := path.Join(workDir, "credentials.json")
	credBytes, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("error marshaling credentials: %s", err)
	}

	if err := os.WriteFile(credPath, credBytes, 0600); err != nil {
		return fmt.Errorf("error writing credentials: %s", err)
	}

	fmt.Println("tfarm server configured")

	return nil
}
