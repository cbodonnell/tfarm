package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/frpc"
	"github.com/google/uuid"
)

const httpTunnelTemplate = `[{{ .Name }}]
type = {{ .Type }}
local_ip = {{ .LocalIP }}
local_port = {{ .LocalPort }}
subdomain = TBD
meta_proxy_id = {{ .ProxyID }}
`

const tcpTunnelTemplate = `[{{ .Name }}]
type = {{ .Type }}
local_ip = {{ .LocalIP }}
local_port = {{ .LocalPort }}
remote_port = {{ .RemotePort }}
meta_proxy_id = {{ .ProxyID }}
`

func HandleCreate(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var createRequest api.CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			log.Printf("failed to decode request body: %s", err)
			api.RespondWithError(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		tunnelConfigPath := filepath.Join(f.WorkDir, "conf.d", createRequest.Name+".ini")
		if _, err := os.Stat(tunnelConfigPath); err == nil {
			log.Printf("tunnel already exists: %s", createRequest.Name)
			api.RespondWithError(w, http.StatusConflict, fmt.Sprintf("tunnel already exists: %s", createRequest.Name))
			return
		}

		createRequest.ProxyID = uuid.New().String()

		var tunnelConfigTemplate *template.Template
		switch createRequest.Type {
		case "http", "https":
			tunnelConfigTemplate = template.Must(template.New("tunnel").Parse(httpTunnelTemplate))
		case "tcp":
			tunnelConfigTemplate = template.Must(template.New("tunnel").Parse(tcpTunnelTemplate))
		default:
			log.Printf("invalid tunnel type: %s", createRequest.Type)
			api.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid tunnel type: %s", createRequest.Type))
			return
		}

		tunnelConfig := &bytes.Buffer{}
		if err := tunnelConfigTemplate.Execute(tunnelConfig, createRequest); err != nil {
			log.Printf("failed to execute template: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to execute template")
			return
		}

		if err := os.WriteFile(tunnelConfigPath, tunnelConfig.Bytes(), 0600); err != nil {
			log.Printf("failed to write tunnel config: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to write tunnel config")
			return
		}

		if _, err := f.Output("verify"); err != nil {
			log.Printf("failed to verify: %s", err)
			if err := os.Remove(tunnelConfigPath); err != nil {
				log.Printf("failed to delete tunnel config file: %s", err)
			}
			api.RespondWithError(w, http.StatusInternalServerError, "failed to verify")
			return
		}

		if _, err := f.Output("reload"); err != nil {
			log.Printf("failed to reload: %s", err)
			if err := os.Remove(tunnelConfigPath); err != nil {
				log.Printf("failed to delete tunnel config file: %s", err)
			}
			api.RespondWithError(w, http.StatusInternalServerError, "failed to reload")
			return
		}

		api.RespondWithSuccess(w, "tunnel created")
	}
}
