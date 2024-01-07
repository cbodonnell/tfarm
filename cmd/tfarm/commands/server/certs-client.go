package server

import (
	"fmt"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/certs"
	"github.com/spf13/cobra"
)

func CertsClientCmd() *cobra.Command {
	certsClientCmd := &cobra.Command{
		Use:           "client [name]",
		Short:         "Generate a client certificate",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			return CertsClient(args[0])
		},
	}

	return certsClientCmd
}

func CertsClient(name string) error {
	workDir := os.Getenv("TFARMD_WORK_DIR")
	if workDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %s", err)
		}
		workDir = pwd
	}

	return certs.GenerateClientCerts(path.Join(workDir, "tls"), name)
}
