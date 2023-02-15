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

	// pre-login routes
	preLogin := r.NewRoute().Subrouter()
	preLogin.HandleFunc("/api/login", HandleLogin(o)).Methods("PUT")

	// post-login routes
	postLogin := r.NewRoute().Subrouter()
	postLogin.HandleFunc("/api/status", HandleStatus(f)).Methods("GET")
	postLogin.HandleFunc("/api/verify", HandleVerify(f)).Methods("GET")
	postLogin.HandleFunc("/api/reload", HandleReload(f)).Methods("POST")
	postLogin.HandleFunc("/api/restart", HandleRestart(f)).Methods("POST")
	postLogin.HandleFunc("/api/tunnel", HandleCreate(f)).Methods("POST")
	postLogin.HandleFunc("/api/tunnel/{name}", HandleDelete(f)).Methods("DELETE")
	postLogin.Use(IsAuthenticated(o))

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
