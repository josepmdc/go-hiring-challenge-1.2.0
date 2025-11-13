package api

import (
	"encoding/json"
	"net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// we shouldn't really return internal error messages in the response
		// like this since they can expose critical information of the system.
		// We should probably have some middleware that cleans error messages,
		// but for the simplicity of the task we'll leave it like this
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]any{"error": message}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
