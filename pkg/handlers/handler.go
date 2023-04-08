package handlers

import (
	"log"
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/frpc"
	"github.com/gorilla/mux"
)

func NewMuxHandler(f *frpc.Frpc) http.Handler {
	r := mux.NewRouter()

	// pre-configure routes
	preConfigure := r.NewRoute().Subrouter()
	preConfigure.HandleFunc("/api/info", HandleInfo()).Methods("GET")
	preConfigure.HandleFunc("/api/configure", HandleConfigure(f)).Methods("PUT")

	// post-configure routes
	postConfigure := r.NewRoute().Subrouter()
	postConfigure.HandleFunc("/api/status", HandleStatus(f)).Methods("GET")
	postConfigure.HandleFunc("/api/verify", HandleVerify(f)).Methods("GET")
	postConfigure.HandleFunc("/api/reload", HandleReload(f)).Methods("POST")
	postConfigure.HandleFunc("/api/restart", HandleRestart(f)).Methods("POST")
	postConfigure.HandleFunc("/api/tunnel", HandleCreate(f)).Methods("POST")
	postConfigure.HandleFunc("/api/tunnel/{name}", HandleDelete(f)).Methods("DELETE")
	postConfigure.Use(isConfiguredMiddleware, isCmdMiddlware(f))

	return r
}

func isConfiguredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsConfigured() {
			log.Printf("not configured")
			api.RespondWithError(w, http.StatusUnauthorized, "tfarmd not configured. run `tfarmd configure`")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isCmdMiddlware(f *frpc.Frpc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !f.IsCmd() {
				log.Printf("frpc not running")
				api.RespondWithError(w, http.StatusUnauthorized, "frpc not running. check tfarm server logs for more information")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
