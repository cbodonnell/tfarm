package api

import (
	"fmt"
	"log"
	"net/http"
)

type APIServer struct {
	handler http.Handler
	port    int
	ErrChan chan error
}

func NewServer(handler http.Handler, port int) *APIServer {
	return &APIServer{
		handler: handler,
		port:    port,
		ErrChan: make(chan error),
	}
}

func (a *APIServer) Start() {
	go func() {
		log.Printf("api server listening on port %d", a.port)
		a.ErrChan <- http.ListenAndServe(fmt.Sprintf(":%d", a.port), a.handler)
	}()
}
