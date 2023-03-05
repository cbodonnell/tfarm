package handlers

import (
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/api"
	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/frpc"
	"github.com/gorilla/mux"
)

func NewMuxHandler(o auth.AuthWorker, f *frpc.Frpc) http.Handler {
	r := mux.NewRouter()

	// pre-configure routes
	preConfigure := r.NewRoute().Subrouter()
	preConfigure.HandleFunc("/api/configure", HandleConfigure(o)).Methods("PUT")

	// post-configure routes
	postConfigure := r.NewRoute().Subrouter()
	postConfigure.HandleFunc("/api/status", HandleStatus(f)).Methods("GET")
	postConfigure.HandleFunc("/api/verify", HandleVerify(f)).Methods("GET")
	postConfigure.HandleFunc("/api/reload", HandleReload(f)).Methods("POST")
	postConfigure.HandleFunc("/api/restart", HandleRestart(f)).Methods("POST")
	postConfigure.HandleFunc("/api/tunnel", HandleCreate(f)).Methods("POST")
	postConfigure.HandleFunc("/api/tunnel/{name}", HandleDelete(f)).Methods("DELETE")
	postConfigure.Use(IsAuthenticated(o))

	return r
}

// middleware that checks if the oauth worker is authenticated
func IsAuthenticated(o auth.AuthWorker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !o.IsAuthenticated() {
				api.RespondWithError(w, http.StatusUnauthorized, "not logged in")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
