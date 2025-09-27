package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func SendResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if body == nil {
		 return
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		 log.Error().Err(err).Msg("failed to encode response")
	}
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
