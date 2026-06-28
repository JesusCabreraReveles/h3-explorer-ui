package handler

import (
	"net/http"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// H3Handler exposes the H3 endpoints. It depends only on the domain port, so it
// can be unit-tested with a fake service and is fully decoupled from h3-go.
type H3Handler struct {
	svc domain.H3Service
}

// NewH3Handler wires the handler to an H3 service implementation (DI).
func NewH3Handler(svc domain.H3Service) *H3Handler {
	return &H3Handler{svc: svc}
}

// --- Request DTOs ---------------------------------------------------------

// fromCoordinatesRequest is the body for POST /api/h3/from-coordinates.
// Pointers distinguish "omitted" from a legitimate 0 value (lat/lng 0 is valid).
type fromCoordinatesRequest struct {
	Lat        *float64 `json:"lat"`
	Lng        *float64 `json:"lng"`
	Resolution *int     `json:"resolution"`
}

// cellRequest is the body for endpoints that operate on an existing index.
type cellRequest struct {
	Index string `json:"index"`
}

// --- Handlers -------------------------------------------------------------

// FromCoordinates indexes a coordinate at a resolution and returns the cell.
func (h *H3Handler) FromCoordinates(w http.ResponseWriter, r *http.Request) {
	var req fromCoordinatesRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Lat == nil || req.Lng == nil {
		writeError(w, r, badRequest("fields 'lat' and 'lng' are required"))
		return
	}
	if req.Resolution == nil {
		writeError(w, r, badRequest("field 'resolution' is required"))
		return
	}

	info, err := h.svc.FromCoordinates(r.Context(), *req.Lat, *req.Lng, *req.Resolution)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, info)
}

// Inspect returns the full description of a cell index.
func (h *H3Handler) Inspect(w http.ResponseWriter, r *http.Request) {
	var req cellRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Index == "" {
		writeError(w, r, badRequest("field 'index' is required"))
		return
	}

	info, err := h.svc.Inspect(r.Context(), req.Index)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, info)
}

// boundaryResponse wraps the boundary ring so the payload is a JSON object
// rather than a bare array (more forward-compatible).
type boundaryResponse struct {
	Index    string          `json:"index"`
	Boundary domain.Boundary `json:"boundary"`
}

// ToBoundary returns just the polygon ring for a cell index.
func (h *H3Handler) ToBoundary(w http.ResponseWriter, r *http.Request) {
	var req cellRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Index == "" {
		writeError(w, r, badRequest("field 'index' is required"))
		return
	}

	ring, err := h.svc.Boundary(r.Context(), req.Index)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, boundaryResponse{Index: req.Index, Boundary: ring})
}

// resolutionsResponse wraps the resolution table.
type resolutionsResponse struct {
	Resolutions []domain.ResolutionInfo `json:"resolutions"`
}

// Resolutions returns aggregate metadata for every H3 resolution.
func (h *H3Handler) Resolutions(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.Resolutions(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, resolutionsResponse{Resolutions: res})
}
