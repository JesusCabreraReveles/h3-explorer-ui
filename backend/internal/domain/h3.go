// Package domain holds the core, framework-agnostic types and ports of the
// application. It must not import transport (HTTP) or infrastructure (h3-go)
// packages — those depend on it, never the other way around. This is the
// stable center of the Clean Architecture onion.
package domain

import (
	"context"
	"errors"
)

// Resolution bounds as defined by the H3 specification.
const (
	MinResolution = 0
	MaxResolution = 15
)

// Safety limits. Grid and hierarchy operations can grow exponentially, so the
// service guards inputs to keep responses bounded and the process protected
// from accidental (or malicious) blow-ups.
const (
	// MaxGridK bounds the radius k accepted by grid-disk and grid-ring.
	// gridDisk(50) already yields ~7,651 cells.
	MaxGridK = 50
	// MaxResultCells caps the number of cells any single operation may return.
	MaxResultCells = 100_000
)

// Sentinel errors. Transport layers map these to HTTP status codes via
// errors.Is, keeping HTTP concerns out of the service layer.
var (
	// ErrInvalidResolution is returned when a resolution is outside [0,15].
	ErrInvalidResolution = errors.New("resolution must be between 0 and 15")
	// ErrInvalidCoordinates is returned for out-of-range latitude/longitude.
	ErrInvalidCoordinates = errors.New("coordinates out of range")
	// ErrInvalidCell is returned when an H3 index string is malformed or invalid.
	ErrInvalidCell = errors.New("invalid h3 cell index")
	// ErrInvalidK is returned when a grid radius k is negative or exceeds MaxGridK.
	ErrInvalidK = errors.New("k must be between 0 and the configured maximum")
	// ErrResultTooLarge is returned when an operation would exceed MaxResultCells.
	ErrResultTooLarge = errors.New("operation would return too many cells; narrow the request")
	// ErrGridPath is returned when no grid path exists between two cells (e.g.
	// different resolutions or cells too far apart).
	ErrGridPath = errors.New("no grid path between the given cells")
	// ErrInvalidPolygon is returned when a polygon has fewer than 3 vertices.
	ErrInvalidPolygon = errors.New("polygon must have at least 3 vertices")
)

// LatLng is a geographic coordinate in decimal degrees.
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Boundary is the ordered ring of vertices describing a cell's polygon.
// For most cells this is 6 vertices; for pentagons it is 5 (and may include
// extra vertices where a cell crosses an icosahedron edge).
type Boundary []LatLng

// CellInfo is the complete, transport-agnostic description of a single H3 cell.
// It is the payload behind the inspector and from-coordinates features.
type CellInfo struct {
	Index            string   `json:"index"`
	Resolution       int      `json:"resolution"`
	Center           LatLng   `json:"center"`
	Boundary         Boundary `json:"boundary"`
	AreaKm2          float64  `json:"areaKm2"`
	AreaM2           float64  `json:"areaM2"`
	EdgeLengthKm     float64  `json:"edgeLengthKm"`
	EdgeLengthM      float64  `json:"edgeLengthM"`
	BaseCell         int      `json:"baseCell"`
	IcosahedronFaces []int    `json:"icosahedronFaces"`
	IsPentagon       bool     `json:"isPentagon"`
	IsClassIII       bool     `json:"isClassIII"`
	Parent           string   `json:"parent,omitempty"`
	Children         []string `json:"children"`
	NumChildren      int      `json:"numChildren"`
	Neighbors        []string `json:"neighbors"`
}

// CellGeometry is the minimal renderable description of a cell: enough to draw
// and label it on a map without a follow-up request. Grid and hierarchy
// endpoints return collections of these so the frontend never has to recompute
// H3 geometry itself.
type CellGeometry struct {
	Index    string   `json:"index"`
	Center   LatLng   `json:"center"`
	Boundary Boundary `json:"boundary"`
}

// Polygon is a single polygon with an outer ring and zero or more holes.
type Polygon struct {
	Outer Boundary   `json:"outer"`
	Holes []Boundary `json:"holes"`
}

// MultiPolygon is the merged outline of a set of cells (cells-to-multi-polygon).
type MultiPolygon struct {
	Polygons []Polygon `json:"polygons"`
}

// ResolutionInfo describes the aggregate characteristics of an H3 resolution.
// It powers the resolution explorer.
type ResolutionInfo struct {
	Resolution      int     `json:"resolution"`
	AvgAreaKm2      float64 `json:"avgAreaKm2"`
	AvgEdgeLengthKm float64 `json:"avgEdgeLengthKm"`
	TotalCells      int     `json:"totalCells"`
}

// H3Service is the port the transport layer depends on. The concrete
// implementation lives in internal/service/h3. Defining the interface here
// (consumer side) is what lets handlers be tested with a fake and keeps the
// dependency arrows pointing inward.
type H3Service interface {
	// FromCoordinates indexes a lat/lng at the given resolution and returns the
	// full cell description.
	FromCoordinates(ctx context.Context, lat, lng float64, resolution int) (CellInfo, error)
	// Inspect returns the full description of an existing cell index.
	Inspect(ctx context.Context, index string) (CellInfo, error)
	// Boundary returns just the polygon ring for a cell index.
	Boundary(ctx context.Context, index string) (Boundary, error)
	// Resolutions returns aggregate metadata for every H3 resolution 0..15.
	Resolutions(ctx context.Context) ([]ResolutionInfo, error)

	// GridDisk returns every cell within grid distance k of the origin
	// (including the origin).
	GridDisk(ctx context.Context, index string, k int) ([]CellGeometry, error)
	// GridRing returns the hollow ring of cells at exactly grid distance k.
	GridRing(ctx context.Context, index string, k int) ([]CellGeometry, error)
	// GridPath returns the line of cells connecting origin and destination.
	GridPath(ctx context.Context, origin, destination string) ([]CellGeometry, error)
	// Parent returns the ancestor of a cell at the given (coarser) resolution.
	Parent(ctx context.Context, index string, parentResolution int) (CellGeometry, error)
	// Children returns the descendants of a cell at the given (finer) resolution.
	Children(ctx context.Context, index string, childResolution int) ([]CellGeometry, error)
	// Neighbors returns the immediately adjacent cells (the ring at distance 1).
	Neighbors(ctx context.Context, index string) ([]CellGeometry, error)
	// PolygonToCells returns the cells covering a polygon (with optional holes)
	// at the given resolution.
	PolygonToCells(ctx context.Context, outer Boundary, holes []Boundary, resolution int) ([]CellGeometry, error)
	// CellsToMultiPolygon merges a set of cells into their multi-polygon outline.
	CellsToMultiPolygon(ctx context.Context, indexes []string) (MultiPolygon, error)
}
