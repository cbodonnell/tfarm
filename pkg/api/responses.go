package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func RespondWithSuccess(w http.ResponseWriter, message string) {
	response := &Response{
		Success: true,
		Message: message,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write HTTP response: %s", err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, errMsg string) {
	response := &Response{
		Success: false,
		Error:   errMsg,
	}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write HTTP response: %s", err)
	}
}
