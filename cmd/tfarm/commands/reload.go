package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:           "reload",
	Short:         "Reload the frpc configuration",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Reload()
	},
}

func Reload() error {
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
