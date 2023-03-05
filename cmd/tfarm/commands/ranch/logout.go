package ranch

import (
	"fmt"

	"github.com/cbodonnell/oauth2utils/pkg/persistence"
	"github.com/spf13/cobra"
)

func LogoutCmd(tokenDir string) *cobra.Command {
	logoutCmd := &cobra.Command{
		Use:           "logout",
		Short:         "Logout of ranch",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Logout(tokenDir)
		},
	}

	return logoutCmd
}

func Logout(tokenDir string) error {
	if _, err := persistence.LoadToken(tokenDir); err != nil {
		fmt.Println("not logged in")
		return nil
	}

	if err := persistence.DeleteToken(tokenDir); err != nil {
		return fmt.Errorf("error deleting token: %s", err)
	}

	fmt.Println("logged out")

	return nil
}
