package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:           "verify",
	Short:         "Verify the frpc configuration",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Verify()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func Verify() error {
	req := &api.APIRequest{}
	status, err := client.Verify(req)
	if err != nil {
		return fmt.Errorf("error verifying: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
