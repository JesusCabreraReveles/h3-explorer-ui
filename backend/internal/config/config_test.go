package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// No env vars set (t.Setenv guarantees isolation by failing on parallel use).
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("default port: want 8080, got %d", cfg.Port)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("default log format: want json, got %s", cfg.LogFormat)
	}
	if len(cfg.CORSAllowedOrigins) == 0 {
		t.Error("expected a default CORS origin")
	}
	if cfg.Addr() != "0.0.0.0:8080" {
		t.Errorf("addr: want 0.0.0.0:8080, got %s", cfg.Addr())
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("H3_PORT", "9090")
	t.Setenv("H3_LOG_LEVEL", "debug")
	t.Setenv("H3_READ_TIMEOUT", "5s")
	t.Setenv("H3_CORS_ALLOWED_ORIGINS", "https://a.com, https://b.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("port: want 9090, got %d", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("log level: want debug, got %s", cfg.LogLevel)
	}
	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("read timeout: want 5s, got %s", cfg.ReadTimeout)
	}
	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Errorf("cors origins: want 2, got %d (%v)", len(cfg.CORSAllowedOrigins), cfg.CORSAllowedOrigins)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	t.Setenv("H3_PORT", "70000")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestLoadInvalidPortFallsBackOnGarbage(t *testing.T) {
	// Non-numeric values should fall back to the default rather than error.
	t.Setenv("H3_PORT", "not-a-number")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("port: want fallback 8080, got %d", cfg.Port)
	}
}
