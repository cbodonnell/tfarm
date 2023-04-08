package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/version"
)

func HandleInfo() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		info := &api.ServerInfoResponse{
			Version: version.Version,
		}
		output, err := json.Marshal(info)
		if err != nil {
			log.Printf("failed to marshal info: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to marshal info")
			return
		}
		api.RespondWithSuccess(w, string(output))
	}
}
