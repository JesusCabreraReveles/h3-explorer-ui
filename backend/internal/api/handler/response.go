// Package handler contains the HTTP adapters that translate requests into
// service calls and service results into JSON responses. It depends on the
// domain port, never on the concrete service implementation.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/pkg/logging"
)

// maxBodyBytes caps request bodies to defend against abusive payloads.
const maxBodyBytes = 1 << 20 // 1 MiB

// errorEnvelope is the single, consistent error shape returned by every endpoint.
type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON serializes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logging.FromContext(r.Context()).Error("encoding response", slog.Any("err", err))
	}
}

// writeError maps a domain error to the appropriate HTTP status and a stable
// machine-readable code. Centralizing this keeps handlers free of HTTP/status
// boilerplate and guarantees consistent error semantics across the API.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, code := http.StatusInternalServerError, "internal_error"
	message := "internal server error"

	switch {
	case errors.Is(err, domain.ErrInvalidResolution):
		status, code, message = http.StatusBadRequest, "invalid_resolution", err.Error()
	case errors.Is(err, domain.ErrInvalidCoordinates):
		status, code, message = http.StatusBadRequest, "invalid_coordinates", err.Error()
	case errors.Is(err, domain.ErrInvalidCell):
		status, code, message = http.StatusBadRequest, "invalid_cell", err.Error()
	case errors.Is(err, domain.ErrInvalidK):
		status, code, message = http.StatusBadRequest, "invalid_k", err.Error()
	case errors.Is(err, domain.ErrInvalidPolygon):
		status, code, message = http.StatusBadRequest, "invalid_polygon", err.Error()
	case errors.Is(err, domain.ErrGridPath):
		status, code, message = http.StatusBadRequest, "grid_path_failed", err.Error()
	case errors.Is(err, domain.ErrResultTooLarge):
		// 422: the request is well-formed but the result is too large to return.
		status, code, message = http.StatusUnprocessableEntity, "result_too_large", err.Error()
	case errors.Is(err, errBadRequest):
		status, code, message = http.StatusBadRequest, "bad_request", err.Error()
	default:
		// Log unexpected errors with detail but never leak internals to clients.
		logging.FromContext(r.Context()).Error("unhandled error", slog.Any("err", err))
	}

	writeJSON(w, r, status, errorEnvelope{Error: errorBody{Code: code, Message: message}})
}

// errBadRequest marks malformed-request errors (bad JSON, missing fields).
var errBadRequest = errors.New("bad request")

// decodeJSON strictly decodes a JSON request body into dst, rejecting unknown
// fields and oversized or malformed payloads with a bad-request error.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return badRequest("malformed JSON body: %v", err)
	}
	// Ensure the body contains exactly one JSON value.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return badRequest("request body must contain a single JSON object")
	}
	return nil
}

// badRequest wraps a formatted message as a bad-request error that writeError
// maps to HTTP 400.
func badRequest(format string, args ...any) error {
	return fmt.Errorf("%w: %s", errBadRequest, fmt.Sprintf(format, args...))
}
