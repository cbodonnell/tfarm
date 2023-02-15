package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/auth"
)

func HandleLogin(o auth.AuthWorker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !o.NeedsLogin() {
			api.RespondWithError(w, http.StatusUnauthorized, "already logged in")
			return
		}

		var loginRequest api.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			log.Printf("failed to decode request body: %s", err)
			api.RespondWithError(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		credentials := auth.LoginCredentials{
			Username: loginRequest.Username,
			Password: loginRequest.Password,
		}

		log.Println("sending login credentials to channel")
		o.LoginCredentialsChan() <- credentials
		log.Println("sent login credentials to channel")

		select {
		case result := <-o.LoginResultChan():
			if !result.Success {
				api.RespondWithError(w, http.StatusUnauthorized, result.Error.Error())
				return
			}
		case <-time.After(10 * time.Second):
			api.RespondWithError(w, http.StatusUnauthorized, "failed to login: timeout")
			return
		}

		api.RespondWithSuccess(w, "logged in")
	}
}
