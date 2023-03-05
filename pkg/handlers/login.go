package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/auth"
)

func HandleConfigure(o auth.AuthWorker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Should be able to configure at an
		if !o.NeedsConfigure() {
			api.RespondWithError(w, http.StatusUnauthorized, "already configured")
			return
		}

		var configureRequest api.ConfigureRequest
		if err := json.NewDecoder(r.Body).Decode(&configureRequest); err != nil {
			log.Printf("failed to decode request body: %s", err)
			api.RespondWithError(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		credentials := auth.ConfigureCredentials{
			ClientID:     configureRequest.ClientID,
			ClientSecret: configureRequest.ClientSecret,
		}

		log.Println("sending configure credentials to channel")
		o.ConfigureCredentialsChan() <- credentials
		log.Println("sent configure credentials to channel")

		select {
		case result := <-o.ConfigureResultChan():
			if !result.Success {
				api.RespondWithError(w, http.StatusUnauthorized, result.Error.Error())
				return
			}
		case <-time.After(10 * time.Second):
			api.RespondWithError(w, http.StatusUnauthorized, "failed to configure: timeout")
			return
		}

		api.RespondWithSuccess(w, "logged in")
	}
}
