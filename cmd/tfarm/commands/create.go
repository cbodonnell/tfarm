package commands

import (
	"fmt"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/spf13/cobra"
)

func CreateCmd() *cobra.Command {
	var tunnelType string
	var localIP string
	var localPort int

	createCmd := &cobra.Command{
		Use:           "create [NAME]",
		Short:         "Create a new tunnel",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("name is required")
			}
			return Create(args[0], tunnelType, localIP, localPort)
		},
	}

	createCmd.Flags().StringVarP(&tunnelType, "type", "t", "http", "tunnel type (http, tcp, udp)")
	createCmd.Flags().StringVarP(&localIP, "local-ip", "l", "127.0.0.1", "local ip address")
	createCmd.Flags().IntVarP(&localPort, "local-port", "p", 0, "local port (required)")

	return createCmd
}

func Create(name string, tunnelType string, localIP string, localPort int) error {
	if localPort == 0 {
		return fmt.Errorf("local port is required")
	}

	client, err := getClient()
	if err != nil {
		return fmt.Errorf("error creating client: %s", err)
	}

	req := &api.CreateRequest{
		Name:      name,
		Type:      tunnelType,
		LocalIP:   localIP,
		LocalPort: localPort,
	}
	status, err := client.Create(req)
	if err != nil {
		return fmt.Errorf("error creating: %s", err)
	}

	if status.Success {
		fmt.Println(status.Message)
	} else {
		fmt.Println(status.Error)
	}

	return nil
}
