package utils

import (
	"encoding/json"
	"net/http"
)

func ToPtr[T any](value T) *T {
	return &value
}

func RespondWithJSON(w http.ResponseWriter, body interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}
