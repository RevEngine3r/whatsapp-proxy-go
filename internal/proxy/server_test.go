package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/RevEngine3r/whatsapp-proxy-go/internal/config"
)

func TestNew(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Port = 0 // Use random port

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if server == nil {
		t.Fatal("New() returned nil server")
	}

	if server.metrics == nil {
		t.Error("Server metrics not initialized")
	}
}

func TestNewWithSOCKS5(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Port = 0
	cfg.SOCKS5.Enabled = true
	cfg.SOCKS5.Host = "127.0.0.1"
	cfg.SOCKS5.Port = 1080

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("New() with SOCKS5 error = %v", err)
	}

	if server.socks5Client == nil {
		t.Error("SOCKS5 client not initialized")
	}
}

func TestServerStartShutdown(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Port = 0 // Use random port
	cfg.Metrics.Enabled = false // Disable metrics for simpler test

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Start server
	if err := server.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestMetricsServer(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Port = 0
	cfg.Metrics.Enabled = true
	cfg.Metrics.Port = 0 // Random port

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := server.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Give metrics server time to start
	time.Sleep(100 * time.Millisecond)

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func TestGetMetrics(t *testing.T) {
	cfg := config.Default()
	server, _ := New(cfg)

	metrics := server.GetMetrics()
	if metrics == nil {
		t.Error("GetMetrics() returned nil")
	}

	// Test initial values
	if metrics.GetConnectionsTotal() != 0 {
		t.Error("Initial connections total should be 0")
	}

	if metrics.GetConnectionsActive() != 0 {
		t.Error("Initial active connections should be 0")
	}
}
