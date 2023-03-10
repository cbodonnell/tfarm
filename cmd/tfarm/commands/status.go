package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

func StatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:           "status",
		Short:         "Get the status of all tunnels",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Status()
		},
	}

	return statusCmd
}

func Status() error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("error creating client: %s", err)
	}

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
