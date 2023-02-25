package commands

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/frpc"
	"github.com/cbodonnell/tfarm/pkg/handlers"
	"github.com/cbodonnell/tfarm/pkg/version"
	"github.com/fatedier/frp/pkg/config"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tfarmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Start()
	},
}

var port int
var frpcAdminAddr string
var frpcAdminPort int
var frpcLogLevel string
var frpsServerAddr string
var frpsServerPort int
var frpsToken string

func init() {
	startCmd.Flags().IntVarP(&port, "port", "p", api.DefaultPort, "port to listen on")
	startCmd.Flags().StringVar(&frpcAdminAddr, "frpc-admin-addr", "127.0.0.1", "address of frpc admin interface")
	startCmd.Flags().IntVar(&frpcAdminPort, "frpc-admin-port", 7400, "frpc admin port")
	startCmd.Flags().StringVar(&frpcLogLevel, "frpc-log-level", "info", "frpc log level")
	startCmd.Flags().StringVar(&frpsServerAddr, "frps-server-addr", "127.0.0.1", "frps server address")
	startCmd.Flags().IntVar(&frpsServerPort, "frps-server-port", 7000, "frps server port")
	startCmd.Flags().StringVar(&frpsToken, "frps-token", "", "frps token")
	rootCmd.AddCommand(startCmd)
}

func Start() error {
	log.Printf("starting tfarmd version %s", version.Version)

	frpcBinPath := os.Getenv("TFARMD_FRPC_BIN_PATH")
	if frpcBinPath == "" {
		frpcBinPath = "/usr/bin/frpc"
	}

	if _, err := os.Stat(frpcBinPath); os.IsNotExist(err) {
		return fmt.Errorf("frpc binary not found at %s", frpcBinPath)
	}

	workDir := os.Getenv("TFARMD_WORK_DIR")
	if workDir == "" {
		workDir = "/var/lib/tfarmd"
	}

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		return fmt.Errorf("work directory not found at %s", workDir)
	}

	cfg := config.ClientCommonConf{
		ServerAddr: frpsServerAddr,
		ServerPort: frpsServerPort,
		AdminAddr:  frpcAdminAddr,
		AdminPort:  frpcAdminPort,
		LogLevel:   frpcLogLevel,
		Metas:      make(map[string]string),
	}
	cfg.Token = frpsToken

	if err := frpc.SaveFrpcCommonConfig(cfg, path.Join(workDir, "frpc.ini")); err != nil {
		return fmt.Errorf("error saving frpc config: %s", err)
	}

	o := auth.NewClientSecretWorker(workDir)
	f := frpc.New(frpcBinPath, workDir)

	h := handlers.NewMuxHandler(o, f)
	tls := &api.TLSFiles{
		CertFile: path.Join(workDir, "tls", "server.crt"),
		KeyFile:  path.Join(workDir, "tls", "server.key"),
		CAFile:   path.Join(workDir, "tls", "ca.crt"),
	}

	a, err := api.NewServer(h, port, tls)
	if err != nil {
		return fmt.Errorf("error starting api server: %s", err)
	}
	a.Start()

	o.WaitForLogin()

	// TODO: propagate errors from the frpc process
	f.StartAndWait()

	select {
	case err := <-f.ErrChan:
		return fmt.Errorf("frpc exited: %s", err)
	case err := <-a.ErrChan:
		return fmt.Errorf("api server exited: %s", err)
	}
}
