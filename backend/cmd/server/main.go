// Command server is the composition root for the h3-explorer-ui backend. It
// wires configuration, logging, the H3 service, and the HTTP router together
// and runs the server with graceful shutdown. All dependency construction
// happens here and nowhere else.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/api"
	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/config"
	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/openapi"
	h3svc "github.com/JesusCabreraReveles/h3-explorer-ui/backend/internal/service/h3"
	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/pkg/logging"
)

// version is overridable at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	// -healthcheck lets the binary probe its own /health endpoint, which is how
	// the distroless container (no shell, no curl) implements its Docker
	// HEALTHCHECK.
	healthcheck := flag.Bool("healthcheck", false, "probe /health and exit 0 if healthy")
	flag.Parse()

	if *healthcheck {
		os.Exit(runHealthcheck())
	}

	if err := run(); err != nil {
		// Logger may not exist yet, so fall back to the default.
		slog.Error("fatal", slog.Any("err", err))
		os.Exit(1)
	}
}

// runHealthcheck performs an HTTP GET against the local /health endpoint and
// returns a process exit code (0 healthy, 1 otherwise).
func runHealthcheck() int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck: config:", err)
		return 1
	}
	client := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("http://127.0.0.1:%d/health", cfg.Port)

	resp, err := client.Get(url) //nolint:noctx // short-lived CLI probe
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck: status", resp.StatusCode)
		return 1
	}
	return 0
}

// run holds the real logic so that deferred cleanup runs and we can return
// errors instead of calling os.Exit directly (which would skip defers).
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := logging.New(os.Stdout, cfg.LogLevel, logging.Format(cfg.LogFormat))
	slog.SetDefault(logger)

	// Signal-aware root context: cancelled on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := api.NewRouter(api.Dependencies{
		Logger:             logger,
		H3Service:          h3svc.New(),
		Version:            version,
		CORSAllowedOrigins: cfg.CORSAllowedOrigins,
		OpenAPISpec:        openapi.Spec,
	})

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		BaseContext:  func(net.Listener) context.Context { return ctx },
	}

	// Run the server in the background and surface listen errors.
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting",
			slog.String("addr", cfg.Addr()),
			slog.String("version", version),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Block until either the server fails or a shutdown signal arrives.
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received, draining connections")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return errors.Join(errors.New("graceful shutdown failed"), err)
	}

	logger.Info("server stopped cleanly")
	return nil
}
