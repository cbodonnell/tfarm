package handlers

import (
	"log"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/frpc"
)

func HandleReload(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := f.Output("reload"); err != nil {
			log.Printf("failed to reload: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to reload")
			return
		}
		api.RespondWithSuccess(w, "frpc configuration reloaded")
	}
}
