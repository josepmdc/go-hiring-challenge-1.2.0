package api

import (
	"encoding/json"
	"net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	writeJSON(w, map[string]any{"error": message})
}

func CreatedResponse(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, data)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
