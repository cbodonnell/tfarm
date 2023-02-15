package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:           "delete [NAME]",
	Short:         "Delete a tunnel",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("name is required")
		}
		return Delete(args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func Delete(args []string) error {
	req := &api.DeleteRequest{
		Name: args[0],
	}
	status, err := client.Delete(req)
	if err != nil {
		return fmt.Errorf("error deleting: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
