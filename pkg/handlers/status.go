package handlers

import (
	"log"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/frpc"
)

func HandleStatus(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := f.Status()
		if err != nil {
			log.Printf("failed to get frpc status: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to get frpc status")
			return
		}
		api.RespondWithSuccess(w, string(output))
	}
}
