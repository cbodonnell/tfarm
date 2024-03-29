package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/frpc"
)

func HandleConfigure(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		configureCredentials := &auth.ConfigureCredentials{}
		if err := json.NewDecoder(r.Body).Decode(&configureCredentials); err != nil {
			log.Printf("failed to decode request body: %s", err)
			api.RespondWithError(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		// chech that configure request has all required fields
		if configureCredentials.ClientID == "" ||
			configureCredentials.ClientSecret == "" ||
			configureCredentials.ClientCACert == "" ||
			configureCredentials.ClientTLSCert == "" ||
			configureCredentials.ClientTLSKey == "" {
			log.Printf("client_id, client_secret, client_ca_cert, client_tls_cert, and client_tls_key are required")
			api.RespondWithError(w, http.StatusBadRequest, "client_id, client_secret, client_ca_cert, client_tls_cert, and client_tls_key are required")
			return
		}

		// marshal configure request to json
		configureCredentialsJSON, err := json.Marshal(configureCredentials)
		if err != nil {
			log.Printf("failed to marshal configure request: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to marshal configure request")
			return
		}

		// write configure request to file
		if err := os.WriteFile(path.Join(f.WorkDir, "credentials.json"), configureCredentialsJSON, 0600); err != nil {
			log.Printf("failed to write configure request to file: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to write configure request to file")
			return
		}

		if err := f.SignConfig(configureCredentials); err != nil {
			log.Printf("failed to sign frpc config: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to sign frpc config")
			return
		}

		if err := f.Restart(); err != nil {
			log.Printf("failed to restart frpc: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to restart frpc")
			return
		}

		api.RespondWithSuccess(w, "tfarmd configured")
	}
}
