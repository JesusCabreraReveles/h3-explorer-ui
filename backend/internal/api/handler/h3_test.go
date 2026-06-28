package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// fakeService is a hand-written test double implementing domain.H3Service. It
// lets handler tests assert HTTP behavior without depending on h3-go.
type fakeService struct {
	info       domain.CellInfo
	boundary   domain.Boundary
	resolution []domain.ResolutionInfo
	geometry   domain.CellGeometry
	cells      []domain.CellGeometry
	multi      domain.MultiPolygon
	err        error
}

func (f *fakeService) FromCoordinates(_ context.Context, _, _ float64, _ int) (domain.CellInfo, error) {
	return f.info, f.err
}
func (f *fakeService) Inspect(_ context.Context, _ string) (domain.CellInfo, error) {
	return f.info, f.err
}
func (f *fakeService) Boundary(_ context.Context, _ string) (domain.Boundary, error) {
	return f.boundary, f.err
}
func (f *fakeService) Resolutions(_ context.Context) ([]domain.ResolutionInfo, error) {
	return f.resolution, f.err
}
func (f *fakeService) GridDisk(_ context.Context, _ string, _ int) ([]domain.CellGeometry, error) {
	return f.cells, f.err
}
func (f *fakeService) GridRing(_ context.Context, _ string, _ int) ([]domain.CellGeometry, error) {
	return f.cells, f.err
}
func (f *fakeService) GridPath(_ context.Context, _, _ string) ([]domain.CellGeometry, error) {
	return f.cells, f.err
}
func (f *fakeService) Parent(_ context.Context, _ string, _ int) (domain.CellGeometry, error) {
	return f.geometry, f.err
}
func (f *fakeService) Children(_ context.Context, _ string, _ int) ([]domain.CellGeometry, error) {
	return f.cells, f.err
}
func (f *fakeService) Neighbors(_ context.Context, _ string) ([]domain.CellGeometry, error) {
	return f.cells, f.err
}
func (f *fakeService) PolygonToCells(_ context.Context, _ domain.Boundary, _ []domain.Boundary, _ int) ([]domain.CellGeometry, error) {
	return f.cells, f.err
}
func (f *fakeService) CellsToMultiPolygon(_ context.Context, _ []string) (domain.MultiPolygon, error) {
	return f.multi, f.err
}

func TestFromCoordinatesHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		svc        *fakeService
		wantStatus int
		wantCode   string // error code expected in the envelope, if any
	}{
		{
			name:       "valid request returns 200",
			body:       `{"lat":37.77,"lng":-122.41,"resolution":9}`,
			svc:        &fakeService{info: domain.CellInfo{Index: "8928308280fffff", Resolution: 9}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing coordinates returns 400",
			body:       `{"resolution":9}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name:       "missing resolution returns 400",
			body:       `{"lat":1,"lng":2}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name:       "unknown field returns 400",
			body:       `{"lat":1,"lng":2,"resolution":9,"foo":true}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name:       "malformed JSON returns 400",
			body:       `{not json`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "bad_request",
		},
		{
			name:       "service resolution error maps to 400",
			body:       `{"lat":1,"lng":2,"resolution":99}`,
			svc:        &fakeService{err: domain.ErrInvalidResolution},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_resolution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewH3Handler(tt.svc)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/h3/from-coordinates", strings.NewReader(tt.body))

			h.FromCoordinates(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: want %d, got %d (body=%s)", tt.wantStatus, rec.Code, rec.Body.String())
			}
			if tt.wantCode != "" {
				assertErrorCode(t, rec.Body.Bytes(), tt.wantCode)
			}
		})
	}
}

func TestInspectHandler(t *testing.T) {
	t.Parallel()

	t.Run("valid index returns cell info", func(t *testing.T) {
		t.Parallel()
		svc := &fakeService{info: domain.CellInfo{Index: "8928308280fffff", Resolution: 9, BaseCell: 20}}
		h := NewH3Handler(svc)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/inspect", strings.NewReader(`{"index":"8928308280fffff"}`))

		h.Inspect(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: want 200, got %d", rec.Code)
		}
		var got domain.CellInfo
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got.BaseCell != 20 {
			t.Errorf("baseCell: want 20, got %d", got.BaseCell)
		}
	})

	t.Run("empty index returns 400", func(t *testing.T) {
		t.Parallel()
		h := NewH3Handler(&fakeService{})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/h3/inspect", strings.NewReader(`{"index":""}`))

		h.Inspect(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: want 400, got %d", rec.Code)
		}
	})
}

func TestResolutionsHandler(t *testing.T) {
	t.Parallel()
	svc := &fakeService{resolution: []domain.ResolutionInfo{{Resolution: 0, AvgAreaKm2: 1, AvgEdgeLengthKm: 1, TotalCells: 122}}}
	h := NewH3Handler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/h3/resolutions", nil)

	h.Resolutions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var got resolutionsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Resolutions) != 1 || got.Resolutions[0].TotalCells != 122 {
		t.Errorf("unexpected payload: %+v", got)
	}
}

// assertErrorCode decodes the error envelope and checks the machine-readable code.
func assertErrorCode(t *testing.T, body []byte, want string) {
	t.Helper()
	var env errorEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("decode error envelope: %v (body=%s)", err, body)
	}
	if env.Error.Code != want {
		t.Errorf("error code: want %q, got %q", want, env.Error.Code)
	}
}
