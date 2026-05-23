package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	log.Printf("http error status=%d message=%q", status, message)
	writeJSON(w, status, errorResponse{Error: message})
}

func writeInternalError(w http.ResponseWriter, message string, err error) {
	log.Printf("internal error message=%q err=%v", message, err)
	writeJSON(w, http.StatusInternalServerError, errorResponse{Error: message})
}
