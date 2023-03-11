package handlers

import (
	"net/http"

	"github.com/cbodonnell/tfarm/pkg/frpc"
	"github.com/gorilla/mux"
)

func NewMuxHandler(f *frpc.Frpc) http.Handler {
	r := mux.NewRouter()

	// TODO: pre-configure routes
	// preConfigure := r.NewRoute().Subrouter()

	// post-configure routes
	postConfigure := r.NewRoute().Subrouter()
	postConfigure.HandleFunc("/api/status", HandleStatus(f)).Methods("GET")
	postConfigure.HandleFunc("/api/verify", HandleVerify(f)).Methods("GET")
	postConfigure.HandleFunc("/api/reload", HandleReload(f)).Methods("POST")
	postConfigure.HandleFunc("/api/restart", HandleRestart(f)).Methods("POST")
	postConfigure.HandleFunc("/api/tunnel", HandleCreate(f)).Methods("POST")
	postConfigure.HandleFunc("/api/tunnel/{name}", HandleDelete(f)).Methods("DELETE")

	return r
}
