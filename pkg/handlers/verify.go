package handlers

import (
	"log"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/frpc"
)

func HandleVerify(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := f.Output("verify"); err != nil {
			log.Printf("failed to verify: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to verify")
			return
		}
		api.RespondWithSuccess(w, "frpc configuration verified")
	}
}
