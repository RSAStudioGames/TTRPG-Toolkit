package api

import (
	"encoding/json"
	"net/http"
)

// Envelope is the standardized API response format.
type Envelope struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, body Envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteSuccess sends a success envelope.
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, Envelope{Status: "success", Data: data})
}

// WriteError sends an error envelope.
func WriteError(w http.ResponseWriter, status int, message string, errors []string) {
	writeJSON(w, status, Envelope{Status: "error", Message: message, Errors: errors})
}
