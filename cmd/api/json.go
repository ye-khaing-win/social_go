package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

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

func (app *application) errorResponse(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	return writeJSON(w, status, envelope{
		Status:  "error",
		Message: message,
	})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelop struct {
		Status string `json:"status"`
		Data   any    `json:"data"`
	}

	return writeJSON(w, status, envelop{
		Status: "success",
		Data:   data,
	})
}
