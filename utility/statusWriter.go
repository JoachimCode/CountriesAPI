package utility

import (
	"encoding/json"
	"net/http"
)

func StatusWriter(w http.ResponseWriter, status int, message string) {
	statusResponse := StatusResponse{
		Message:    message,
		StatusCode: status,
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse)
}
