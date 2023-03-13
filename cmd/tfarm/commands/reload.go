package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

func ReloadCmd() *cobra.Command {
	reloadCmd := &cobra.Command{
		Use:           "reload",
		Short:         "Reload the frpc configuration",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Reload()
		},
	}

	return reloadCmd
}

func Reload() error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("error creating client: %s", err)
	}

	req := &api.APIRequest{}
	status, err := client.Reload(req)
	if err != nil {
		return fmt.Errorf("error reloading: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
