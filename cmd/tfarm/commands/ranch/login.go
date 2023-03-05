package ranch

import (
	"context"
	"fmt"

	"github.com/cbodonnell/oauth2utils/pkg/persistence"
	"github.com/cbodonnell/oauth2utils/pkg/utils"
	"github.com/cbodonnell/tfarm/pkg/ranch/auth"
	"github.com/cbodonnell/tfarm/pkg/term"
	"github.com/spf13/cobra"
)

func LoginCmd(tokenDir string) *cobra.Command {
	var username string
	var password string

	loginCmd := &cobra.Command{
		Use:           "login",
		Short:         "Login to ranch",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Login(tokenDir, username, password)
		},
	}

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "ranch username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "ranch password")

	return loginCmd
}

func Login(tokenDir, username, password string) error {
	ctx := context.Background()
	oc, err := auth.NewOIDCClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating OIDC client: %s", err)
	}

	var message string
	token := utils.TryGetToken(ctx, oc, tokenDir)
	if !token.Valid() {
		if username == "" {
			username = term.StringPrompt("Username:")
		}
		if password == "" {
			password = term.PasswordPrompt("Password:")
		}
		newToken, err := oc.Password(ctx, username, password)
		if err != nil {
			return fmt.Errorf("error logging in: %s", err)
		}
		if err := persistence.SaveToken(newToken, tokenDir); err != nil {
			return fmt.Errorf("error saving token: %s", err)
		}
		token = newToken
		message = "logged in"
	} else {
		message = "already logged in"
	}

	fmt.Println(message)

	return nil
}
