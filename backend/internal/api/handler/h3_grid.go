package handler

import (
	"context"
	"net/http"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// --- Request DTOs ---------------------------------------------------------

// gridKRequest is shared by grid-disk and grid-ring.
type gridKRequest struct {
	Index string `json:"index"`
	K     *int   `json:"k"`
}

type gridPathRequest struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

// hierarchyRequest is shared by parent and children: an index plus a target
// resolution.
type hierarchyRequest struct {
	Index      string `json:"index"`
	Resolution *int   `json:"resolution"`
}

type polygonToCellsRequest struct {
	Polygon    domain.Boundary   `json:"polygon"`
	Holes      []domain.Boundary `json:"holes"`
	Resolution *int              `json:"resolution"`
}

type cellsRequest struct {
	Indexes []string `json:"indexes"`
}

// --- Response DTOs --------------------------------------------------------

// cellsResponse is the common envelope for endpoints returning many cells. The
// count is included so clients can display it without re-counting.
type cellsResponse struct {
	Count int                   `json:"count"`
	Cells []domain.CellGeometry `json:"cells"`
}

func newCellsResponse(cells []domain.CellGeometry) cellsResponse {
	return cellsResponse{Count: len(cells), Cells: cells}
}

// --- Handlers -------------------------------------------------------------

// GridDisk returns all cells within grid distance k of the origin.
func (h *H3Handler) GridDisk(w http.ResponseWriter, r *http.Request) {
	h.handleGridK(w, r, h.svc.GridDisk)
}

// GridRing returns the hollow ring of cells at exactly grid distance k.
func (h *H3Handler) GridRing(w http.ResponseWriter, r *http.Request) {
	h.handleGridK(w, r, h.svc.GridRing)
}

// handleGridK factors out the shared decode/validate/dispatch flow for the two
// k-radius endpoints, which differ only in the service call.
func (h *H3Handler) handleGridK(
	w http.ResponseWriter,
	r *http.Request,
	op func(ctx context.Context, index string, k int) ([]domain.CellGeometry, error),
) {
	var req gridKRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Index == "" {
		writeError(w, r, badRequest("field 'index' is required"))
		return
	}
	if req.K == nil {
		writeError(w, r, badRequest("field 'k' is required"))
		return
	}

	cells, err := op(r.Context(), req.Index, *req.K)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, newCellsResponse(cells))
}

// GridPath returns the line of cells connecting origin and destination.
func (h *H3Handler) GridPath(w http.ResponseWriter, r *http.Request) {
	var req gridPathRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Origin == "" || req.Destination == "" {
		writeError(w, r, badRequest("fields 'origin' and 'destination' are required"))
		return
	}

	cells, err := h.svc.GridPath(r.Context(), req.Origin, req.Destination)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, newCellsResponse(cells))
}

// Parent returns the ancestor of a cell at the requested resolution.
func (h *H3Handler) Parent(w http.ResponseWriter, r *http.Request) {
	var req hierarchyRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Index == "" {
		writeError(w, r, badRequest("field 'index' is required"))
		return
	}
	if req.Resolution == nil {
		writeError(w, r, badRequest("field 'resolution' is required"))
		return
	}

	cell, err := h.svc.Parent(r.Context(), req.Index, *req.Resolution)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, cell)
}

// Children returns the descendants of a cell at the requested resolution.
func (h *H3Handler) Children(w http.ResponseWriter, r *http.Request) {
	var req hierarchyRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Index == "" {
		writeError(w, r, badRequest("field 'index' is required"))
		return
	}
	if req.Resolution == nil {
		writeError(w, r, badRequest("field 'resolution' is required"))
		return
	}

	cells, err := h.svc.Children(r.Context(), req.Index, *req.Resolution)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, newCellsResponse(cells))
}

// Neighbors returns the immediately adjacent cells.
func (h *H3Handler) Neighbors(w http.ResponseWriter, r *http.Request) {
	var req cellRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if req.Index == "" {
		writeError(w, r, badRequest("field 'index' is required"))
		return
	}

	cells, err := h.svc.Neighbors(r.Context(), req.Index)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, newCellsResponse(cells))
}

// PolygonToCells returns the cells covering a polygon at a resolution.
func (h *H3Handler) PolygonToCells(w http.ResponseWriter, r *http.Request) {
	var req polygonToCellsRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if len(req.Polygon) == 0 {
		writeError(w, r, badRequest("field 'polygon' is required"))
		return
	}
	if req.Resolution == nil {
		writeError(w, r, badRequest("field 'resolution' is required"))
		return
	}

	cells, err := h.svc.PolygonToCells(r.Context(), req.Polygon, req.Holes, *req.Resolution)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, newCellsResponse(cells))
}

// CellsToMultiPolygon merges a set of cells into their outline.
func (h *H3Handler) CellsToMultiPolygon(w http.ResponseWriter, r *http.Request) {
	var req cellsRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, r, err)
		return
	}
	if len(req.Indexes) == 0 {
		writeError(w, r, badRequest("field 'indexes' must contain at least one cell"))
		return
	}

	multi, err := h.svc.CellsToMultiPolygon(r.Context(), req.Indexes)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, r, http.StatusOK, multi)
}
