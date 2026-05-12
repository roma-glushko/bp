package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/roma-glushko/bp/internal/domain"
)

const maxBodySize = 1 << 20 // 1MB

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeValidationErrors(w http.ResponseWriter, errs []domain.ValidationError) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
}

func readJSON(r *http.Request, dst any) error {
	body := io.LimitReader(r.Body, maxBodySize)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}
