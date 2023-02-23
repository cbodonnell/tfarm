package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type APIServer struct {
	server  *http.Server
	port    int
	ErrChan chan error
}

type TLSFiles struct {
	CertFile string
	KeyFile  string
	CAFile   string
}

func NewServer(handler http.Handler, port int, tlsFiles *TLSFiles) (*APIServer, error) {
	cert, err := tls.LoadX509KeyPair(tlsFiles.CertFile, tlsFiles.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %s", err)
	}

	caCert, err := ioutil.ReadFile(tlsFiles.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
	return &APIServer{
		server:  server,
		port:    port,
		ErrChan: make(chan error),
	}, nil
}

func (a *APIServer) Start() {
	go func() {
		log.Printf("api server listening on %s", a.server.Addr)
		a.ErrChan <- a.server.ListenAndServeTLS("", "")
	}()
}
