package frpc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/crypto"
	"github.com/cbodonnell/tfarm/pkg/logging"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/rodaine/table"
)

type Frpc struct {
	binPath      string
	WorkDir      string
	stdout       io.Writer
	stderr       io.Writer
	cmd          *exec.Cmd
	StartErrChan chan error
	ErrChan      chan error
	ExitChan     chan struct{}
	restarting   bool
}

type ErrCredentialsNotFound struct {
	Err error
}

func (e *ErrCredentialsNotFound) Error() string {
	return fmt.Sprintf("credentials not found: %s", e.Err)
}

func New(binPath, workDir string) *Frpc {
	return &Frpc{
		binPath:      binPath,
		WorkDir:      workDir,
		stdout:       os.Stdout,
		stderr:       os.Stderr,
		cmd:          nil,
		StartErrChan: make(chan error),
		ErrChan:      make(chan error),
		ExitChan:     make(chan struct{}),
		restarting:   false,
	}
}

func (f *Frpc) IsCmd() bool {
	return f.cmd != nil
}

func (f *Frpc) StartLoop() {
	go func() {
		restartDelay := 5 * time.Second
		for {
			creds, err := auth.WaitForCredentials(f.WorkDir)
			if err != nil {
				f.StartErrChan <- fmt.Errorf("error waiting for credentials: %s", err)
			}

			if err := f.SignConfig(creds); err != nil {
				f.StartErrChan <- fmt.Errorf("error signing frpc config: %s", err)
			}

			f.StartAndWait()

			select {
			case err = <-f.ErrChan:
				log.Printf("frpc exited: %s", err)
				log.Printf("restarting frpc in %s", restartDelay.String())
				time.Sleep(restartDelay)
			}
		}
	}()
}

func (f *Frpc) SignConfig(creds *auth.ConfigureCredentials) error {
	decodedSecret, err := base64.StdEncoding.DecodeString(creds.ClientSecret)
	if err != nil {
		return fmt.Errorf("error decoding client secret: %s", err)
	}

	cfg, err := ParseFrpcCommonConfig(path.Join(f.WorkDir, "frpc.ini"))
	if err != nil {
		return fmt.Errorf("error reading frpc.ini: %s", err)
	}

	if cfg.Metas == nil {
		cfg.Metas = make(map[string]string)
	}

	cfg.Metas["client_id"] = creds.ClientID
	cfg.Metas["client_signature"] = crypto.HMAC(decodedSecret, []byte(creds.ClientID))

	if err := SaveFrpcCommonConfig(cfg, path.Join(f.WorkDir, "frpc.ini")); err != nil {
		return fmt.Errorf("error writing frpc.ini: %s", err)
	}

	return nil
}

func (f *Frpc) Start() error {
	log.Println("starting frpc")

	if f.cmd != nil {
		return errors.New("frpc already running")
	}

	f.cmd = exec.Command(f.binPath, "-c", "frpc.ini")
	f.cmd.Dir = f.WorkDir
	stdout, _ := f.cmd.StdoutPipe()
	stderr, _ := f.cmd.StderrPipe()
	go logging.LogReaderWithPrefix(stdout, "frpc stdout: ")
	go logging.LogReaderWithPrefix(stderr, "frpc stderr: ")

	if err := f.cmd.Start(); err != nil {
		f.cmd = nil
		return fmt.Errorf("failed to start frpc: %s", err)
	}

	log.Println("frpc started")

	return nil
}

func (f *Frpc) Wait() error {
	if f.cmd == nil {
		return errors.New("frpc not running")
	}

	if err := f.cmd.Wait(); err != nil {
		if f.restarting {
			return nil
		}
		f.cmd = nil
		return fmt.Errorf("frpc exited unexpectedly: %s", err)
	}

	f.cmd = nil
	return errors.New("frpc exited unexpectedly with no error")
}

func (f *Frpc) StartAndWait() {
	go func() {
		if err := f.Start(); err != nil {
			f.ErrChan <- fmt.Errorf("failed to start frpc: %s", err)
			return
		}

		if err := f.Wait(); err != nil {
			f.ErrChan <- fmt.Errorf("frpc exited unexpectedly: %s", err)
			return
		}

		f.ExitChan <- struct{}{}
	}()
}

func (f *Frpc) Stop() error {
	log.Println("stopping frpc")

	if f.cmd == nil {
		return errors.New("frpc not running")
	}

	if err := f.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to send interrupt signal to frpc: %s", err)
	}

	select {
	case <-time.After(5 * time.Second):
		log.Println("frpc did not exit gracefully after 5 seconds, killing...")
		f.cmd.Process.Kill()
		<-f.ExitChan
		log.Println("frpc killed")
	case <-f.ExitChan:
		log.Println("frpc exited gracefully")
	}

	f.cmd = nil

	return nil
}

func (f *Frpc) Restart() error {
	log.Println("restarting frpc")

	f.restarting = true
	if err := f.Stop(); err != nil {
		return fmt.Errorf("failed to stop frpc: %s", err)
	}
	f.restarting = false

	creds, err := auth.WaitForCredentials(f.WorkDir)
	if err != nil {
		return fmt.Errorf("error waiting for credentials: %s", err)
	}

	if err := f.SignConfig(creds); err != nil {
		return fmt.Errorf("error signing config: %s", err)
	}

	f.StartAndWait()

	return nil
}

func (f *Frpc) Output(cmd string) ([]byte, error) {
	frpcCmd := exec.Command(f.binPath, cmd, "-c", "frpc.ini")
	frpcCmd.Dir = f.WorkDir

	output, err := frpcCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute frpc verify: %s", err)
	}

	return output, nil
}

func (f *Frpc) Status() ([]byte, error) {
	clientCfg, err := config.UnmarshalClientConfFromIni(path.Join(f.WorkDir, "frpc.ini"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse frpc config: %s", err)
	}

	if clientCfg.AdminPort == 0 {
		return nil, fmt.Errorf("admin_port shoud be set if you want to get proxy status")
	}

	endpoint := fmt.Sprintf("http://%s:%d/api/status", clientCfg.AdminAddr, clientCfg.AdminPort)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %s", err)
	}

	req.SetBasicAuth(clientCfg.AdminUser, clientCfg.AdminPwd)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http request failed with status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response error: %s", err)
	}
	res := &client.StatusResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("parse http response error: %s", err)
	}

	buf := new(bytes.Buffer)
	tbl := table.New("Name", "Type", "Status", "LocalAddr", "Plugin", "RemoteAddr", "Error").WithWriter(buf)

	if len(res.TCP) > 0 {
		for _, ps := range res.TCP {
			tbl.AddRow(ps.Name, "TCP", ps.Status, ps.LocalAddr, ps.Plugin, ps.RemoteAddr, ps.Err)
		}
	}
	if len(res.UDP) > 0 {
		for _, ps := range res.UDP {
			tbl.AddRow(ps.Name, "UDP", ps.Status, ps.LocalAddr, ps.Plugin, ps.RemoteAddr, ps.Err)
		}
	}
	if len(res.HTTP) > 0 {
		for _, ps := range res.HTTP {
			if ps.RemoteAddr != "" {
				ps.RemoteAddr = fmt.Sprintf("http://%s", ps.RemoteAddr)
			}
			tbl.AddRow(ps.Name, "HTTP", ps.Status, ps.LocalAddr, ps.Plugin, ps.RemoteAddr, ps.Err)
		}
	}
	if len(res.HTTPS) > 0 {
		for _, ps := range res.HTTPS {
			if ps.RemoteAddr != "" {
				ps.RemoteAddr = fmt.Sprintf("https://%s", ps.RemoteAddr)
			}
			tbl.AddRow(ps.Name, "HTTPS", ps.Status, ps.LocalAddr, ps.Plugin, ps.RemoteAddr, ps.Err)
		}
	}
	if len(res.STCP) > 0 {
		for _, ps := range res.STCP {
			tbl.AddRow(ps.Name, "STCP", ps.Status, ps.LocalAddr, ps.Plugin, ps.RemoteAddr, ps.Err)
		}
	}
	if len(res.XTCP) > 0 {
		for _, ps := range res.XTCP {
			tbl.AddRow(ps.Name, "XTCP", ps.Status, ps.LocalAddr, ps.Plugin, ps.RemoteAddr, ps.Err)
		}
	}

	tbl.Print()

	return buf.Bytes(), nil
}
