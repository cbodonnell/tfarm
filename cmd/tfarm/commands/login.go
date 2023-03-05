package commands

import (
	"errors"
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/term"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:           "login",
	Short:         "Configure tfarmd with your tunnel.farm credentials",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Login()
	},
}

var username string
var password string

func init() {
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "tunnel.farm username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "tunnel.farm password")
}

func Login() error {
	if username == "" {
		// cli prompt for username
		username = term.StringPrompt("Username:")
	}

	if username == "" {
		return errors.New("username is required")
	}

	if password == "" {
		// cli prompt for password
		password = term.PasswordPrompt("Password:")
	}

	if password == "" {
		return errors.New("password is required")
	}

	req := &api.LoginRequest{
		Username: username,
		Password: password,
	}
	status, err := client.Login(req)
	if err != nil {
		return fmt.Errorf("error logging in: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
