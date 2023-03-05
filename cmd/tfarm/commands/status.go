package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:           "status",
	Short:         "Get the status of all tunnels",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Status()
	},
}

func Status() error {
	req := &api.APIRequest{}
	status, err := client.Status(req)
	if err != nil {
		return fmt.Errorf("error getting status: %s", err)
	}

	if status.Success {
		fmt.Print(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
