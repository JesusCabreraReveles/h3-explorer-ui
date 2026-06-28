package h3

import (
	"context"
	"errors"
	"testing"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

func TestGridDisk(t *testing.T) {
	t.Parallel()
	svc := New()

	tests := []struct {
		name      string
		index     string
		k         int
		wantCount int // 1 + 3k(k+1) for a non-pentagon cell
		wantErr   error
	}{
		{name: "k=0 is just the origin", index: sfRes9Index, k: 0, wantCount: 1},
		{name: "k=1 disk has 7 cells", index: sfRes9Index, k: 1, wantCount: 7},
		{name: "k=2 disk has 19 cells", index: sfRes9Index, k: 2, wantCount: 19},
		{name: "negative k rejected", index: sfRes9Index, k: -1, wantErr: domain.ErrInvalidK},
		{name: "k over max rejected", index: sfRes9Index, k: domain.MaxGridK + 1, wantErr: domain.ErrInvalidK},
		{name: "invalid index rejected", index: "nope", k: 1, wantErr: domain.ErrInvalidCell},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cells, err := svc.GridDisk(context.Background(), tt.index, tt.k)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(cells) != tt.wantCount {
				t.Errorf("count: want %d, got %d", tt.wantCount, len(cells))
			}
			for _, c := range cells {
				if c.Index == "" || len(c.Boundary) < 5 {
					t.Errorf("cell missing geometry: %+v", c)
				}
			}
		})
	}
}

func TestGridRing(t *testing.T) {
	t.Parallel()
	svc := New()

	// For a non-pentagon cell, ring(k) has exactly 6k cells (ring(0) == 1).
	for k, want := range map[int]int{0: 1, 1: 6, 2: 12, 3: 18} {
		cells, err := svc.GridRing(context.Background(), sfRes9Index, k)
		if err != nil {
			t.Fatalf("k=%d: unexpected error: %v", k, err)
		}
		if len(cells) != want {
			t.Errorf("ring(%d): want %d cells, got %d", k, want, len(cells))
		}
	}

	if _, err := svc.GridRing(context.Background(), sfRes9Index, -1); !errors.Is(err, domain.ErrInvalidK) {
		t.Errorf("want ErrInvalidK, got %v", err)
	}
}

func TestNeighbors(t *testing.T) {
	t.Parallel()
	svc := New()

	cells, err := svc.Neighbors(context.Background(), sfRes9Index)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 6 {
		t.Errorf("want 6 neighbors, got %d", len(cells))
	}
}

func TestParent(t *testing.T) {
	t.Parallel()
	svc := New()

	t.Run("resolution 8 parent matches golden value", func(t *testing.T) {
		t.Parallel()
		parent, err := svc.Parent(context.Background(), sfRes9Index, 8)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if parent.Index != sfRes9Parent {
			t.Errorf("parent: want %s, got %s", sfRes9Parent, parent.Index)
		}
	})

	t.Run("finer parent resolution rejected", func(t *testing.T) {
		t.Parallel()
		_, err := svc.Parent(context.Background(), sfRes9Index, 10)
		if !errors.Is(err, domain.ErrInvalidResolution) {
			t.Errorf("want ErrInvalidResolution, got %v", err)
		}
	})
}

func TestChildren(t *testing.T) {
	t.Parallel()
	svc := New()

	t.Run("immediate children count is 7", func(t *testing.T) {
		t.Parallel()
		cells, err := svc.Children(context.Background(), sfRes9Index, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cells) != 7 {
			t.Errorf("want 7 children, got %d", len(cells))
		}
	})

	t.Run("coarser child resolution rejected", func(t *testing.T) {
		t.Parallel()
		_, err := svc.Children(context.Background(), sfRes9Index, 8)
		if !errors.Is(err, domain.ErrInvalidResolution) {
			t.Errorf("want ErrInvalidResolution, got %v", err)
		}
	})

	t.Run("excessive resolution gap rejected", func(t *testing.T) {
		t.Parallel()
		// res 0 -> res 15 would be ~7^15 cells.
		res0, _ := svc.FromCoordinates(context.Background(), sfLat, sfLng, 0)
		_, err := svc.Children(context.Background(), res0.Index, 15)
		if !errors.Is(err, domain.ErrResultTooLarge) {
			t.Errorf("want ErrResultTooLarge, got %v", err)
		}
	})
}

func TestGridPath(t *testing.T) {
	t.Parallel()
	svc := New()

	t.Run("path from a cell to itself is a single cell", func(t *testing.T) {
		t.Parallel()
		cells, err := svc.GridPath(context.Background(), sfRes9Index, sfRes9Index)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cells) != 1 || cells[0].Index != sfRes9Index {
			t.Errorf("want single origin cell, got %+v", cells)
		}
	})

	t.Run("path between adjacent neighbors has length 2", func(t *testing.T) {
		t.Parallel()
		ring, err := svc.GridRing(context.Background(), sfRes9Index, 1)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		cells, err := svc.GridPath(context.Background(), sfRes9Index, ring[0].Index)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cells) != 2 {
			t.Errorf("want path length 2, got %d", len(cells))
		}
	})

	t.Run("path across resolutions fails", func(t *testing.T) {
		t.Parallel()
		_, err := svc.GridPath(context.Background(), sfRes9Index, sfRes9Parent)
		if !errors.Is(err, domain.ErrGridPath) {
			t.Errorf("want ErrGridPath, got %v", err)
		}
	})
}

func TestPolygonToCells(t *testing.T) {
	t.Parallel()
	svc := New()

	// A small box around downtown SF.
	box := domain.Boundary{
		{Lat: 37.770, Lng: -122.430},
		{Lat: 37.770, Lng: -122.410},
		{Lat: 37.790, Lng: -122.410},
		{Lat: 37.790, Lng: -122.430},
	}

	t.Run("covers the polygon with resolution-9 cells", func(t *testing.T) {
		t.Parallel()
		cells, err := svc.PolygonToCells(context.Background(), box, nil, 9)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cells) == 0 {
			t.Fatal("expected at least one cell")
		}
		// Every covering cell must report resolution 9.
		for _, c := range cells {
			info, ierr := svc.Inspect(context.Background(), c.Index)
			if ierr != nil {
				t.Fatalf("inspect %s: %v", c.Index, ierr)
			}
			if info.Resolution != 9 {
				t.Errorf("cell %s has resolution %d, want 9", c.Index, info.Resolution)
			}
		}
	})

	t.Run("degenerate polygon rejected", func(t *testing.T) {
		t.Parallel()
		_, err := svc.PolygonToCells(context.Background(), box[:2], nil, 9)
		if !errors.Is(err, domain.ErrInvalidPolygon) {
			t.Errorf("want ErrInvalidPolygon, got %v", err)
		}
	})

	t.Run("invalid resolution rejected", func(t *testing.T) {
		t.Parallel()
		_, err := svc.PolygonToCells(context.Background(), box, nil, 99)
		if !errors.Is(err, domain.ErrInvalidResolution) {
			t.Errorf("want ErrInvalidResolution, got %v", err)
		}
	})
}

func TestCellsToMultiPolygon(t *testing.T) {
	t.Parallel()
	svc := New()

	// The outline of a gridDisk should form a single, hole-free polygon.
	disk, err := svc.GridDisk(context.Background(), sfRes9Index, 2)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	indexes := make([]string, len(disk))
	for i, c := range disk {
		indexes[i] = c.Index
	}

	multi, err := svc.CellsToMultiPolygon(context.Background(), indexes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(multi.Polygons) != 1 {
		t.Fatalf("want 1 polygon, got %d", len(multi.Polygons))
	}
	if len(multi.Polygons[0].Outer) < 6 {
		t.Errorf("outer ring too small: %d vertices", len(multi.Polygons[0].Outer))
	}

	t.Run("empty input rejected", func(t *testing.T) {
		t.Parallel()
		_, err := svc.CellsToMultiPolygon(context.Background(), nil)
		if !errors.Is(err, domain.ErrInvalidCell) {
			t.Errorf("want ErrInvalidCell, got %v", err)
		}
	})
}
