package server

import (
	"fmt"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/certs"
	"github.com/spf13/cobra"
)

func CertsRegenerateCmd() *cobra.Command {
	certsRegenerateCmd := &cobra.Command{
		Use:           "regenerate",
		Short:         "Regenerate TLS certificates",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return CertsRegenerate()
		},
	}

	return certsRegenerateCmd
}

func CertsRegenerate() error {
	workDir := os.Getenv("TFARMD_WORK_DIR")
	if workDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %s", err)
		}
		workDir = pwd
	}

	return certs.GenerateCerts(path.Join(workDir, "tls"))
}
