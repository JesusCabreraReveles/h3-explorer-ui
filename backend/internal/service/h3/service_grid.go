package h3

import (
	"context"
	"fmt"

	lib "github.com/uber/h3-go/v4"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// GridDisk returns every cell within grid distance k of the origin.
func (s *Service) GridDisk(_ context.Context, index string, k int) ([]domain.CellGeometry, error) {
	cell, err := parseCell(index)
	if err != nil {
		return nil, err
	}
	if err := validateK(k); err != nil {
		return nil, err
	}
	cells, err := cell.GridDisk(k)
	if err != nil {
		return nil, fmt.Errorf("grid disk: %w", err)
	}
	return toGeometries(cells)
}

// GridRing returns the hollow ring of cells at exactly grid distance k.
func (s *Service) GridRing(_ context.Context, index string, k int) ([]domain.CellGeometry, error) {
	cell, err := parseCell(index)
	if err != nil {
		return nil, err
	}
	if err := validateK(k); err != nil {
		return nil, err
	}
	// Cell.GridRing is the safe variant: it falls back gracefully near pentagons
	// where the unsafe ring algorithm fails.
	cells, err := cell.GridRing(k)
	if err != nil {
		return nil, fmt.Errorf("grid ring: %w", err)
	}
	return toGeometries(cells)
}

// GridPath returns the line of cells connecting origin and destination.
func (s *Service) GridPath(_ context.Context, origin, destination string) ([]domain.CellGeometry, error) {
	from, err := parseCell(origin)
	if err != nil {
		return nil, fmt.Errorf("origin: %w", err)
	}
	to, err := parseCell(destination)
	if err != nil {
		return nil, fmt.Errorf("destination: %w", err)
	}
	cells, err := from.GridPath(to)
	if err != nil {
		// A path failure is a client-side concern (unreachable / different
		// resolutions), so surface it as a domain error.
		return nil, fmt.Errorf("%w: %s", domain.ErrGridPath, err.Error())
	}
	return toGeometries(cells)
}

// Parent returns the ancestor of a cell at the given coarser resolution.
func (s *Service) Parent(_ context.Context, index string, parentResolution int) (domain.CellGeometry, error) {
	cell, err := parseCell(index)
	if err != nil {
		return domain.CellGeometry{}, err
	}
	if err := validateResolution(parentResolution); err != nil {
		return domain.CellGeometry{}, err
	}
	if parentResolution > cell.Resolution() {
		return domain.CellGeometry{}, fmt.Errorf("%w: parent resolution %d must be <= cell resolution %d",
			domain.ErrInvalidResolution, parentResolution, cell.Resolution())
	}
	parent, err := cell.Parent(parentResolution)
	if err != nil {
		return domain.CellGeometry{}, fmt.Errorf("parent: %w", err)
	}
	return describeGeometry(parent)
}

// Children returns the descendants of a cell at the given finer resolution.
func (s *Service) Children(_ context.Context, index string, childResolution int) ([]domain.CellGeometry, error) {
	cell, err := parseCell(index)
	if err != nil {
		return nil, err
	}
	if err := validateResolution(childResolution); err != nil {
		return nil, err
	}
	if childResolution < cell.Resolution() {
		return nil, fmt.Errorf("%w: child resolution %d must be >= cell resolution %d",
			domain.ErrInvalidResolution, childResolution, cell.Resolution())
	}
	// Guard against exponential blow-up before allocating: each finer resolution
	// multiplies the child count by ~7.
	if estimateChildren(childResolution-cell.Resolution()) > domain.MaxResultCells {
		return nil, fmt.Errorf("%w: resolution gap too large", domain.ErrResultTooLarge)
	}
	children, err := cell.Children(childResolution)
	if err != nil {
		return nil, fmt.Errorf("children: %w", err)
	}
	return toGeometries(children)
}

// Neighbors returns the immediately adjacent cells (the ring at distance 1).
func (s *Service) Neighbors(_ context.Context, index string) ([]domain.CellGeometry, error) {
	cell, err := parseCell(index)
	if err != nil {
		return nil, err
	}
	ring, err := cell.GridRing(1)
	if err != nil {
		return nil, fmt.Errorf("neighbors: %w", err)
	}
	return toGeometries(ring)
}

// PolygonToCells returns the cells covering a polygon at the given resolution.
func (s *Service) PolygonToCells(_ context.Context, outer domain.Boundary, holes []domain.Boundary, resolution int) ([]domain.CellGeometry, error) {
	if err := validateResolution(resolution); err != nil {
		return nil, err
	}
	if len(outer) < 3 {
		return nil, fmt.Errorf("%w: got %d", domain.ErrInvalidPolygon, len(outer))
	}

	poly := lib.GeoPolygon{GeoLoop: toGeoLoop(outer)}
	for _, h := range holes {
		if len(h) < 3 {
			return nil, fmt.Errorf("%w: hole has %d vertices", domain.ErrInvalidPolygon, len(h))
		}
		poly.Holes = append(poly.Holes, toGeoLoop(h))
	}

	cells, err := lib.PolygonToCells(poly, resolution)
	if err != nil {
		return nil, fmt.Errorf("polygon to cells: %w", err)
	}
	if len(cells) > domain.MaxResultCells {
		return nil, fmt.Errorf("%w: %d cells", domain.ErrResultTooLarge, len(cells))
	}
	return toGeometries(cells)
}

// CellsToMultiPolygon merges a set of cells into their multi-polygon outline.
func (s *Service) CellsToMultiPolygon(_ context.Context, indexes []string) (domain.MultiPolygon, error) {
	if len(indexes) == 0 {
		return domain.MultiPolygon{}, fmt.Errorf("%w: at least one cell index is required", domain.ErrInvalidCell)
	}
	cells := make([]lib.Cell, 0, len(indexes))
	for _, idx := range indexes {
		cell, err := parseCell(idx)
		if err != nil {
			return domain.MultiPolygon{}, err
		}
		cells = append(cells, cell)
	}

	polygons, err := lib.CellsToMultiPolygon(cells)
	if err != nil {
		return domain.MultiPolygon{}, fmt.Errorf("cells to multi polygon: %w", err)
	}

	out := domain.MultiPolygon{Polygons: make([]domain.Polygon, 0, len(polygons))}
	for _, p := range polygons {
		poly := domain.Polygon{Outer: geoLoopToBoundary(p.GeoLoop), Holes: make([]domain.Boundary, 0, len(p.Holes))}
		for _, h := range p.Holes {
			poly.Holes = append(poly.Holes, geoLoopToBoundary(h))
		}
		out.Polygons = append(out.Polygons, poly)
	}
	return out, nil
}

// --- helpers --------------------------------------------------------------

// validateK guards the grid radius against negative and excessive values.
func validateK(k int) error {
	if k < 0 || k > domain.MaxGridK {
		return fmt.Errorf("%w: got %d (max %d)", domain.ErrInvalidK, k, domain.MaxGridK)
	}
	return nil
}

// estimateChildren returns an upper bound on the number of children produced by
// descending `gap` resolutions (7^gap). It saturates at MaxResultCells+1 to
// avoid integer overflow on large gaps.
func estimateChildren(gap int) int {
	count := 1
	for i := 0; i < gap; i++ {
		count *= 7
		if count > domain.MaxResultCells {
			return domain.MaxResultCells + 1
		}
	}
	return count
}

// describeGeometry builds the renderable geometry for a single cell.
func describeGeometry(cell lib.Cell) (domain.CellGeometry, error) {
	center, err := cell.LatLng()
	if err != nil {
		return domain.CellGeometry{}, fmt.Errorf("cell center: %w", err)
	}
	ring, err := cell.Boundary()
	if err != nil {
		return domain.CellGeometry{}, fmt.Errorf("cell boundary: %w", err)
	}
	return domain.CellGeometry{
		Index:    cell.String(),
		Center:   domain.LatLng{Lat: center.Lat, Lng: center.Lng},
		Boundary: toBoundary(ring),
	}, nil
}

// toGeometries maps a slice of cells to their renderable geometries.
func toGeometries(cells []lib.Cell) ([]domain.CellGeometry, error) {
	out := make([]domain.CellGeometry, 0, len(cells))
	for _, c := range cells {
		g, err := describeGeometry(c)
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

// toGeoLoop converts a domain boundary into an h3 GeoLoop.
func toGeoLoop(b domain.Boundary) lib.GeoLoop {
	loop := make(lib.GeoLoop, len(b))
	for i, v := range b {
		loop[i] = lib.LatLng{Lat: v.Lat, Lng: v.Lng}
	}
	return loop
}

// geoLoopToBoundary converts an h3 GeoLoop into a domain boundary.
func geoLoopToBoundary(loop lib.GeoLoop) domain.Boundary {
	out := make(domain.Boundary, len(loop))
	for i, v := range loop {
		out[i] = domain.LatLng{Lat: v.Lat, Lng: v.Lng}
	}
	return out
}
