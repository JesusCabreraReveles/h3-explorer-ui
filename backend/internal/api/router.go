// Package api wires the HTTP routing tree together. It is the transport
// composition layer: it knows about chi, middleware, and handlers, but contains
// no business logic itself.
package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/api/handler"
	appmw "github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/api/middleware"
	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/domain"
)

// Dependencies is the set of collaborators the router needs. Passing a struct
// keeps the constructor signature stable as the app grows.
type Dependencies struct {
	Logger             *slog.Logger
	H3Service          domain.H3Service
	Version            string
	CORSAllowedOrigins []string
	// OpenAPISpec is the raw bytes of the OpenAPI document served at /openapi.yaml.
	OpenAPISpec []byte
}

// NewRouter builds the fully-configured HTTP handler for the application.
func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	// --- Global middleware (order matters) ---
	r.Use(chimw.RequestID)
	r.Use(appmw.RequestLogger(deps.Logger))
	r.Use(appmw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.CORSAllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	health := handler.NewHealthHandler(deps.Version)
	h3h := handler.NewH3Handler(deps.H3Service)

	// --- Operational endpoints ---
	r.Get("/health", health.Health)
	if len(deps.OpenAPISpec) > 0 {
		r.Get("/openapi.yaml", serveOpenAPI(deps.OpenAPISpec))
	}

	// --- H3 API ---
	r.Route("/api/h3", func(r chi.Router) {
		// Indexing & inspection.
		r.Get("/resolutions", h3h.Resolutions)
		r.Post("/from-coordinates", h3h.FromCoordinates)
		r.Post("/to-boundary", h3h.ToBoundary)
		r.Post("/inspect", h3h.Inspect)

		// Traversal & hierarchy.
		r.Post("/grid-disk", h3h.GridDisk)
		r.Post("/grid-ring", h3h.GridRing)
		r.Post("/grid-path", h3h.GridPath)
		r.Post("/parent", h3h.Parent)
		r.Post("/children", h3h.Children)
		r.Post("/neighbors", h3h.Neighbors)

		// Polygon operations.
		r.Post("/polygon-to-cells", h3h.PolygonToCells)
		r.Post("/cells-to-multi-polygon", h3h.CellsToMultiPolygon)
	})

	return r
}

// serveOpenAPI returns the embedded OpenAPI document.
func serveOpenAPI(spec []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(spec)
	}
}
