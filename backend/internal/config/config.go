// Package config loads and validates runtime configuration from the
// environment, following the 12-factor approach. All defaults live here so the
// rest of the application can assume a fully-populated, validated Config.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all tunables for the HTTP server.
type Config struct {
	// Host:Port the server binds to.
	Host string
	Port int

	// ReadTimeout / WriteTimeout / IdleTimeout protect the server from slow or
	// abusive clients.
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// ShutdownTimeout bounds graceful shutdown.
	ShutdownTimeout time.Duration

	// LogLevel is one of debug|info|warn|error.
	LogLevel string
	// LogFormat is json|text.
	LogFormat string

	// CORSAllowedOrigins is the list of origins permitted to call the API.
	CORSAllowedOrigins []string
}

// Addr returns the host:port string for net/http.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Load reads configuration from the environment, applying sane defaults, and
// validates the result. It returns an error rather than panicking so the
// composition root controls process exit.
func Load() (Config, error) {
	cfg := Config{
		Host:               getEnv("H3_HOST", "0.0.0.0"),
		Port:               getEnvInt("H3_PORT", 8080),
		ReadTimeout:        getEnvDuration("H3_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:       getEnvDuration("H3_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:        getEnvDuration("H3_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:    getEnvDuration("H3_SHUTDOWN_TIMEOUT", 15*time.Second),
		LogLevel:           getEnv("H3_LOG_LEVEL", "info"),
		LogFormat:          getEnv("H3_LOG_FORMAT", "json"),
		CORSAllowedOrigins: getEnvList("H3_CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// validate enforces invariants that the rest of the app relies on.
func (c Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port %d: must be in [1,65535]", c.Port)
	}
	if len(c.CORSAllowedOrigins) == 0 {
		return fmt.Errorf("at least one CORS origin must be configured")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvList(key string, fallback []string) []string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		parts := strings.Split(v, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return fallback
}
