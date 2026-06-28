package handler

import (
	"net/http"
	"time"
)

// HealthHandler serves liveness/readiness checks.
type HealthHandler struct {
	version string
	started time.Time
}

// NewHealthHandler constructs a health handler stamped with the build version.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{version: version, started: time.Now()}
}

type healthResponse struct {
	Status        string `json:"status"`
	Version       string `json:"version"`
	UptimeSeconds int64  `json:"uptimeSeconds"`
}

// Health responds with service status, version, and uptime.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, healthResponse{
		Status:        "ok",
		Version:       h.version,
		UptimeSeconds: int64(time.Since(h.started).Seconds()),
	})
}
