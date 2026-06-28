package h3

import (
	"context"
	"errors"
	"testing"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// Golden values computed directly from uber/h3-go v4.5.0 for the coordinate
// (37.775938, -122.41795) — a well-known reference point near downtown SF.
const (
	sfLat        = 37.775938
	sfLng        = -122.41795
	sfRes9Index  = "8928308280fffff"
	sfRes9Parent = "8828308281fffff"
)

func TestFromCoordinates(t *testing.T) {
	t.Parallel()
	svc := New()

	tests := []struct {
		name       string
		lat, lng   float64
		resolution int
		wantErr    error
		wantIndex  string
	}{
		{
			name:       "known SF cell at resolution 9",
			lat:        sfLat,
			lng:        sfLng,
			resolution: 9,
			wantIndex:  sfRes9Index,
		},
		{
			name:       "resolution 0 is valid",
			lat:        sfLat,
			lng:        sfLng,
			resolution: 0,
		},
		{
			name:       "resolution 15 is valid",
			lat:        sfLat,
			lng:        sfLng,
			resolution: 15,
		},
		{
			name:       "resolution above max is rejected",
			lat:        sfLat,
			lng:        sfLng,
			resolution: 16,
			wantErr:    domain.ErrInvalidResolution,
		},
		{
			name:       "negative resolution is rejected",
			lat:        sfLat,
			lng:        sfLng,
			resolution: -1,
			wantErr:    domain.ErrInvalidResolution,
		},
		{
			name:       "latitude out of range is rejected",
			lat:        91,
			lng:        0,
			resolution: 9,
			wantErr:    domain.ErrInvalidCoordinates,
		},
		{
			name:       "longitude out of range is rejected",
			lat:        0,
			lng:        181,
			resolution: 9,
			wantErr:    domain.ErrInvalidCoordinates,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := svc.FromCoordinates(context.Background(), tt.lat, tt.lng, tt.resolution)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Resolution != tt.resolution {
				t.Errorf("resolution: want %d, got %d", tt.resolution, got.Resolution)
			}
			if tt.wantIndex != "" && got.Index != tt.wantIndex {
				t.Errorf("index: want %s, got %s", tt.wantIndex, got.Index)
			}
			// Invariant valid at every resolution: a cell's center must re-index
			// back to the same cell.
			roundTrip, err := svc.FromCoordinates(context.Background(), got.Center.Lat, got.Center.Lng, tt.resolution)
			if err != nil {
				t.Fatalf("re-indexing center: %v", err)
			}
			if roundTrip.Index != got.Index {
				t.Errorf("center did not round-trip: %s -> %s", got.Index, roundTrip.Index)
			}
			// A cell boundary is a 5- or 6-vertex ring.
			if n := len(got.Boundary); n < 5 || n > 10 {
				t.Errorf("unexpected boundary vertex count: %d", n)
			}
		})
	}
}

func TestInspect(t *testing.T) {
	t.Parallel()
	svc := New()

	t.Run("known SF cell exposes full topology", func(t *testing.T) {
		t.Parallel()
		info, err := svc.Inspect(context.Background(), sfRes9Index)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Index != sfRes9Index {
			t.Errorf("index: want %s, got %s", sfRes9Index, info.Index)
		}
		if info.Resolution != 9 {
			t.Errorf("resolution: want 9, got %d", info.Resolution)
		}
		if info.Parent != sfRes9Parent {
			t.Errorf("parent: want %s, got %s", sfRes9Parent, info.Parent)
		}
		if info.BaseCell != 20 {
			t.Errorf("baseCell: want 20, got %d", info.BaseCell)
		}
		if info.IsPentagon {
			t.Error("expected non-pentagon cell")
		}
		if info.AreaKm2 <= 0 {
			t.Errorf("area must be positive, got %f", info.AreaKm2)
		}
		// A resolution-9 hexagon has exactly 7 immediate children and 6 neighbors.
		if info.NumChildren != 7 {
			t.Errorf("numChildren: want 7, got %d", info.NumChildren)
		}
		if len(info.Neighbors) != 6 {
			t.Errorf("neighbors: want 6, got %d", len(info.Neighbors))
		}
		if len(info.IcosahedronFaces) == 0 {
			t.Error("expected at least one icosahedron face")
		}
	})

	t.Run("invalid index is rejected", func(t *testing.T) {
		t.Parallel()
		_, err := svc.Inspect(context.Background(), "not-a-cell")
		if !errors.Is(err, domain.ErrInvalidCell) {
			t.Fatalf("want ErrInvalidCell, got %v", err)
		}
	})

	t.Run("resolution 0 cell has no parent", func(t *testing.T) {
		t.Parallel()
		res0, err := svc.FromCoordinates(context.Background(), sfLat, sfLng, 0)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		info, err := svc.Inspect(context.Background(), res0.Index)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Parent != "" {
			t.Errorf("res-0 cell should have no parent, got %q", info.Parent)
		}
	})
}

func TestBoundary(t *testing.T) {
	t.Parallel()
	svc := New()

	ring, err := svc.Boundary(context.Background(), sfRes9Index)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ring) != 6 {
		t.Errorf("want 6 boundary vertices, got %d", len(ring))
	}

	if _, err := svc.Boundary(context.Background(), "garbage"); !errors.Is(err, domain.ErrInvalidCell) {
		t.Errorf("want ErrInvalidCell, got %v", err)
	}
}

func TestResolutions(t *testing.T) {
	t.Parallel()
	svc := New()

	res, err := svc.Resolutions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != domain.MaxResolution+1 {
		t.Fatalf("want %d resolutions, got %d", domain.MaxResolution+1, len(res))
	}
	for i, r := range res {
		if r.Resolution != i {
			t.Errorf("entry %d has resolution %d", i, r.Resolution)
		}
		if r.AvgAreaKm2 <= 0 || r.AvgEdgeLengthKm <= 0 || r.TotalCells <= 0 {
			t.Errorf("resolution %d has non-positive metadata: %+v", i, r)
		}
		// Finer resolutions must have strictly smaller average cell area.
		if i > 0 && res[i].AvgAreaKm2 >= res[i-1].AvgAreaKm2 {
			t.Errorf("area not decreasing at resolution %d", i)
		}
	}
}
