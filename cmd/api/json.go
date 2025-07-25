package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int) error {
	type envelope struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	return writeJSON(w, status, envelope{
		Status:  "error",
		Message: http.StatusText(status),
	})
}
