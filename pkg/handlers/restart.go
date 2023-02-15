package handlers

import (
	"log"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/frpc"
)

func HandleRestart(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f.Restart(); err != nil {
			log.Printf("failed to restart: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to restart")
			return
		}
		api.RespondWithSuccess(w, "frpc restarted")
	}
}
