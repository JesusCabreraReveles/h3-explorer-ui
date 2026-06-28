package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// stubService is the minimal H3Service needed to exercise the router wiring.
type stubService struct{}

func (stubService) FromCoordinates(context.Context, float64, float64, int) (domain.CellInfo, error) {
	return domain.CellInfo{Index: "8928308280fffff"}, nil
}
func (stubService) Inspect(context.Context, string) (domain.CellInfo, error) {
	return domain.CellInfo{}, nil
}
func (stubService) Boundary(context.Context, string) (domain.Boundary, error) {
	return domain.Boundary{}, nil
}
func (stubService) Resolutions(context.Context) ([]domain.ResolutionInfo, error) {
	return nil, nil
}
func (stubService) GridDisk(context.Context, string, int) ([]domain.CellGeometry, error) {
	return nil, nil
}
func (stubService) GridRing(context.Context, string, int) ([]domain.CellGeometry, error) {
	return nil, nil
}
func (stubService) GridPath(context.Context, string, string) ([]domain.CellGeometry, error) {
	return nil, nil
}
func (stubService) Parent(context.Context, string, int) (domain.CellGeometry, error) {
	return domain.CellGeometry{}, nil
}
func (stubService) Children(context.Context, string, int) ([]domain.CellGeometry, error) {
	return nil, nil
}
func (stubService) Neighbors(context.Context, string) ([]domain.CellGeometry, error) {
	return nil, nil
}
func (stubService) PolygonToCells(context.Context, domain.Boundary, []domain.Boundary, int) ([]domain.CellGeometry, error) {
	return nil, nil
}
func (stubService) CellsToMultiPolygon(context.Context, []string) (domain.MultiPolygon, error) {
	return domain.MultiPolygon{}, nil
}

func newTestRouter() http.Handler {
	return NewRouter(Dependencies{
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
		H3Service:          stubService{},
		Version:            "test",
		CORSAllowedOrigins: []string{"*"},
		OpenAPISpec:        []byte("openapi: 3.1.0"),
	})
}

func TestRouterRoutes(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(newTestRouter())
	t.Cleanup(srv.Close)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"health", http.MethodGet, "/health", "", http.StatusOK},
		{"openapi", http.MethodGet, "/openapi.yaml", "", http.StatusOK},
		{"resolutions", http.MethodGet, "/api/h3/resolutions", "", http.StatusOK},
		{"from-coordinates", http.MethodPost, "/api/h3/from-coordinates", `{"lat":1,"lng":2,"resolution":9}`, http.StatusOK},
		{"grid-disk", http.MethodPost, "/api/h3/grid-disk", `{"index":"8928308280fffff","k":1}`, http.StatusOK},
		{"grid-path", http.MethodPost, "/api/h3/grid-path", `{"origin":"a","destination":"b"}`, http.StatusOK},
		{"parent", http.MethodPost, "/api/h3/parent", `{"index":"8928308280fffff","resolution":8}`, http.StatusOK},
		{"neighbors", http.MethodPost, "/api/h3/neighbors", `{"index":"8928308280fffff"}`, http.StatusOK},
		{"cells-to-multi-polygon", http.MethodPost, "/api/h3/cells-to-multi-polygon", `{"indexes":["8928308280fffff"]}`, http.StatusOK},
		{"unknown route 404", http.MethodGet, "/api/h3/nope", "", http.StatusNotFound},
		{"wrong method 405", http.MethodGet, "/api/h3/from-coordinates", "", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(tt.method, srv.URL+tt.path, bodyReader(tt.body))
			if err != nil {
				t.Fatalf("build request: %v", err)
			}
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, resp.Body)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func bodyReader(s string) io.Reader {
	if s == "" {
		return nil
	}
	return strings.NewReader(s)
}
