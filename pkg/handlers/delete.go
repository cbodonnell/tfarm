package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/frpc"
	"github.com/gorilla/mux"
)

func HandleDelete(f *frpc.Frpc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tunnelName := vars["name"]

		// delete file, reload, and restore if failed
		tunnelConfigPath := filepath.Join(f.WorkDir, "conf.d", tunnelName+".ini")
		if _, err := os.Stat(tunnelConfigPath); err != nil {
			log.Printf("tunnel does not exist: %s", tunnelName)
			api.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("tunnel does not exist: %s", tunnelName))
			return
		}

		tunnelConfig, err := ioutil.ReadFile(tunnelConfigPath)
		if err != nil {
			log.Printf("failed to read tunnel config file: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to read tunnel config file")
			return
		}

		if err := os.Remove(tunnelConfigPath); err != nil {
			log.Printf("failed to delete tunnel config file: %s", err)
			api.RespondWithError(w, http.StatusInternalServerError, "failed to delete tunnel config file")
			return
		}

		if _, err = f.Output("verify"); err != nil {
			log.Printf("failed to verify: %s", err)
			if err := os.Remove(tunnelConfigPath); err != nil {
				log.Printf("failed to delete tunnel config file: %s", err)
			}
			api.RespondWithError(w, http.StatusInternalServerError, "failed to verify")
			return
		}

		if _, err = f.Output("reload"); err != nil {
			log.Printf("failed to reload: %s", err)
			if err := ioutil.WriteFile(tunnelConfigPath, tunnelConfig, 0644); err != nil {
				log.Printf("failed to restore tunnel config file: %s", err)
			}
			api.RespondWithError(w, http.StatusInternalServerError, "failed to reload")
			return
		}

		api.RespondWithSuccess(w, "tunnel deleted")
	}
}
