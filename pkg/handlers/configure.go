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
		if configureCredentials.ClientID == "" || configureCredentials.ClientSecret == "" {
			log.Printf("client_id and client_secret are required")
			api.RespondWithError(w, http.StatusBadRequest, "client_id and client_secret are required")
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
		if err := os.WriteFile(path.Join(f.WorkDir, "credentials.json"), configureCredentialsJSON, 0644); err != nil {
			log.Printf("failed to write configure request to file: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to write configure request to file")
			return
		}

		api.RespondWithSuccess(w, "tfarmd configured")
	}
}
