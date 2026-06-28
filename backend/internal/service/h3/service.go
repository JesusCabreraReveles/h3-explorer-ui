// Package h3 implements domain.H3Service by adapting Uber's h3-go library to
// the application's domain types. It is the only place in the codebase allowed
// to import h3-go, so all H3 business logic is centralized and the rest of the
// app stays decoupled from the binding.
package h3

import (
	"context"
	"fmt"

	lib "github.com/uber/h3-go/v4"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// Service is the concrete H3 implementation. It is stateless and safe for
// concurrent use, so a single instance is shared across all requests.
type Service struct{}

// New constructs the H3 service. It returns the interface to encourage callers
// to depend on the abstraction rather than the concrete type.
func New() domain.H3Service {
	return &Service{}
}

// FromCoordinates indexes a coordinate at the requested resolution.
func (s *Service) FromCoordinates(ctx context.Context, lat, lng float64, resolution int) (domain.CellInfo, error) {
	if err := validateResolution(resolution); err != nil {
		return domain.CellInfo{}, err
	}
	if err := validateCoordinates(lat, lng); err != nil {
		return domain.CellInfo{}, err
	}

	cell, err := lib.LatLngToCell(lib.LatLng{Lat: lat, Lng: lng}, resolution)
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("indexing coordinate: %w", err)
	}
	return s.describe(ctx, cell)
}

// Inspect returns the full description for an existing cell index.
func (s *Service) Inspect(ctx context.Context, index string) (domain.CellInfo, error) {
	cell, err := parseCell(index)
	if err != nil {
		return domain.CellInfo{}, err
	}
	return s.describe(ctx, cell)
}

// Boundary returns just the polygon ring for a cell index.
func (s *Service) Boundary(_ context.Context, index string) (domain.Boundary, error) {
	cell, err := parseCell(index)
	if err != nil {
		return nil, err
	}
	ring, err := cell.Boundary()
	if err != nil {
		return nil, fmt.Errorf("computing boundary: %w", err)
	}
	return toBoundary(ring), nil
}

// Resolutions returns aggregate metadata for every H3 resolution.
func (s *Service) Resolutions(_ context.Context) ([]domain.ResolutionInfo, error) {
	out := make([]domain.ResolutionInfo, 0, domain.MaxResolution+1)
	for res := domain.MinResolution; res <= domain.MaxResolution; res++ {
		area, err := lib.HexagonAreaAvgKm2(res)
		if err != nil {
			return nil, fmt.Errorf("avg area for res %d: %w", res, err)
		}
		edge, err := lib.HexagonEdgeLengthAvgKm(res)
		if err != nil {
			return nil, fmt.Errorf("avg edge for res %d: %w", res, err)
		}
		out = append(out, domain.ResolutionInfo{
			Resolution:      res,
			AvgAreaKm2:      area,
			AvgEdgeLengthKm: edge,
			TotalCells:      lib.NumCells(res),
		})
	}
	return out, nil
}

// describe assembles the complete CellInfo for a valid cell. It is the single
// source of truth for what the inspector exposes.
func (s *Service) describe(_ context.Context, cell lib.Cell) (domain.CellInfo, error) {
	res := cell.Resolution()

	center, err := cell.LatLng()
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("cell center: %w", err)
	}
	ring, err := cell.Boundary()
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("cell boundary: %w", err)
	}
	areaKm2, err := lib.CellAreaKm2(cell)
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("cell area km2: %w", err)
	}
	areaM2, err := lib.CellAreaM2(cell)
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("cell area m2: %w", err)
	}
	edgeKm, err := lib.HexagonEdgeLengthAvgKm(res)
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("edge length km: %w", err)
	}
	edgeM, err := lib.HexagonEdgeLengthAvgM(res)
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("edge length m: %w", err)
	}
	faces, err := cell.IcosahedronFaces()
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("icosahedron faces: %w", err)
	}

	info := domain.CellInfo{
		Index:            cell.String(),
		Resolution:       res,
		Center:           domain.LatLng{Lat: center.Lat, Lng: center.Lng},
		Boundary:         toBoundary(ring),
		AreaKm2:          areaKm2,
		AreaM2:           areaM2,
		EdgeLengthKm:     edgeKm,
		EdgeLengthM:      edgeM,
		BaseCell:         cell.BaseCellNumber(),
		IcosahedronFaces: faces,
		IsPentagon:       cell.IsPentagon(),
		IsClassIII:       cell.IsResClassIII(),
		Children:         []string{},
		Neighbors:        []string{},
	}

	// Parent (only exists below the coarsest resolution).
	if res > domain.MinResolution {
		parent, perr := cell.ImmediateParent()
		if perr != nil {
			return domain.CellInfo{}, fmt.Errorf("immediate parent: %w", perr)
		}
		info.Parent = parent.String()
	}

	// Immediate children (only exist above the finest resolution).
	if res < domain.MaxResolution {
		children, cerr := cell.ImmediateChildren()
		if cerr != nil {
			return domain.CellInfo{}, fmt.Errorf("immediate children: %w", cerr)
		}
		info.Children = cellsToStrings(children)
		info.NumChildren = len(children)
	}

	// Neighbors: the gridDisk of radius 1 minus the origin cell itself.
	disk, err := cell.GridDisk(1)
	if err != nil {
		return domain.CellInfo{}, fmt.Errorf("grid disk: %w", err)
	}
	neighbors := make([]string, 0, len(disk))
	for _, n := range disk {
		if n == cell {
			continue
		}
		neighbors = append(neighbors, n.String())
	}
	info.Neighbors = neighbors

	return info, nil
}

// validateResolution guards the [0,15] range before touching the binding.
func validateResolution(res int) error {
	if res < domain.MinResolution || res > domain.MaxResolution {
		return fmt.Errorf("%w: got %d", domain.ErrInvalidResolution, res)
	}
	return nil
}

// validateCoordinates rejects coordinates outside valid geographic ranges.
func validateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return fmt.Errorf("%w: lat=%.6f lng=%.6f", domain.ErrInvalidCoordinates, lat, lng)
	}
	return nil
}

// parseCell converts a hex index string into a validated h3 Cell.
func parseCell(index string) (lib.Cell, error) {
	cell := lib.CellFromString(index)
	if !cell.IsValid() {
		return 0, fmt.Errorf("%w: %q", domain.ErrInvalidCell, index)
	}
	return cell, nil
}

// toBoundary converts an h3 CellBoundary into the domain Boundary type.
func toBoundary(ring lib.CellBoundary) domain.Boundary {
	out := make(domain.Boundary, len(ring))
	for i, v := range ring {
		out[i] = domain.LatLng{Lat: v.Lat, Lng: v.Lng}
	}
	return out
}

// cellsToStrings maps a slice of cells to their hex string representations.
func cellsToStrings(cells []lib.Cell) []string {
	out := make([]string, len(cells))
	for i, c := range cells {
		out[i] = c.String()
	}
	return out
}
