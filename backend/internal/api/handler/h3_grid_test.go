package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

func TestGridDiskHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		svc        *fakeService
		wantStatus int
		wantCode   string
	}{
		{
			name:       "valid request returns cells",
			body:       `{"index":"8928308280fffff","k":1}`,
			svc:        &fakeService{cells: []domain.CellGeometry{{Index: "a"}, {Index: "b"}}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing k returns 400",
			body:       `{"index":"8928308280fffff"}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name:       "missing index returns 400",
			body:       `{"k":1}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name:       "service invalid-k maps to 400 invalid_k",
			body:       `{"index":"8928308280fffff","k":-1}`,
			svc:        &fakeService{err: domain.ErrInvalidK},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_k",
		},
		{
			name:       "k=0 is accepted (zero value, pointer disambiguates)",
			body:       `{"index":"8928308280fffff","k":0}`,
			svc:        &fakeService{cells: []domain.CellGeometry{{Index: "a"}}},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewH3Handler(tt.svc)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/h3/grid-disk", strings.NewReader(tt.body))

			h.GridDisk(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: want %d, got %d (body=%s)", tt.wantStatus, rec.Code, rec.Body.String())
			}
			if tt.wantCode != "" {
				assertErrorCode(t, rec.Body.Bytes(), tt.wantCode)
			}
		})
	}
}

func TestChildrenHandlerResultTooLarge(t *testing.T) {
	t.Parallel()
	h := NewH3Handler(&fakeService{err: domain.ErrResultTooLarge})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/h3/children", strings.NewReader(`{"index":"8928308280fffff","resolution":15}`))

	h.Children(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: want 422, got %d", rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "result_too_large")
}

func TestGridPathHandlerErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing destination returns 400", func(t *testing.T) {
		t.Parallel()
		h := NewH3Handler(&fakeService{})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/grid-path", strings.NewReader(`{"origin":"8928308280fffff"}`))
		h.GridPath(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want 400, got %d", rec.Code)
		}
	})

	t.Run("unreachable path maps to grid_path_failed", func(t *testing.T) {
		t.Parallel()
		h := NewH3Handler(&fakeService{err: domain.ErrGridPath})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/grid-path", strings.NewReader(`{"origin":"a","destination":"b"}`))
		h.GridPath(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want 400, got %d", rec.Code)
		}
		assertErrorCode(t, rec.Body.Bytes(), "grid_path_failed")
	})
}

func TestPolygonToCellsHandler(t *testing.T) {
	t.Parallel()

	t.Run("valid polygon returns cells", func(t *testing.T) {
		t.Parallel()
		h := NewH3Handler(&fakeService{cells: []domain.CellGeometry{{Index: "a"}}})
		rec := httptest.NewRecorder()
		body := `{"polygon":[{"lat":37.77,"lng":-122.43},{"lat":37.77,"lng":-122.41},{"lat":37.79,"lng":-122.41}],"resolution":9}`
		req := httptest.NewRequest(http.MethodPost, "/api/h3/polygon-to-cells", strings.NewReader(body))
		h.PolygonToCells(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("want 200, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("empty polygon returns 400", func(t *testing.T) {
		t.Parallel()
		h := NewH3Handler(&fakeService{})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/polygon-to-cells", strings.NewReader(`{"resolution":9}`))
		h.PolygonToCells(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want 400, got %d", rec.Code)
		}
	})
}

func TestCellsToMultiPolygonHandler(t *testing.T) {
	t.Parallel()

	t.Run("empty indexes returns 400", func(t *testing.T) {
		t.Parallel()
		h := NewH3Handler(&fakeService{})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/cells-to-multi-polygon", strings.NewReader(`{"indexes":[]}`))
		h.CellsToMultiPolygon(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want 400, got %d", rec.Code)
		}
	})

	t.Run("valid request returns multipolygon", func(t *testing.T) {
		t.Parallel()
		multi := domain.MultiPolygon{Polygons: []domain.Polygon{{Outer: domain.Boundary{{Lat: 1, Lng: 2}}}}}
		h := NewH3Handler(&fakeService{multi: multi})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/cells-to-multi-polygon", strings.NewReader(`{"indexes":["8928308280fffff"]}`))
		h.CellsToMultiPolygon(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("want 200, got %d", rec.Code)
		}
	})
}
