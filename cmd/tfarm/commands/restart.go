package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

func RestartCmd() *cobra.Command {
	restartCmd := &cobra.Command{
		Use:           "restart",
		Short:         "Restart frpc",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Restart()
		},
	}

	return restartCmd
}

func Restart() error {
	req := &api.APIRequest{}
	status, err := client.Restart(req)
	if err != nil {
		return fmt.Errorf("error restarting: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
