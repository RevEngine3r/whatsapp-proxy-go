package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/RevEngine3r/whatsapp-proxy-go/internal/config"
	"github.com/RevEngine3r/whatsapp-proxy-go/internal/socks5"
)

// Server represents the proxy server
type Server struct {
	config       *config.Config
	listener     net.Listener
	socks5Client *socks5.Client
	metrics      *Metrics
	metricsServer *http.Server
	wg           sync.WaitGroup
	shutdown     chan struct{}
}

// New creates a new proxy server
func New(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	s := &Server{
		config:   cfg,
		metrics:  NewMetrics(),
		shutdown: make(chan struct{}),
	}

	// Create SOCKS5 client if enabled
	if cfg.SOCKS5.Enabled {
		socks5Cfg := &socks5.Config{
			ProxyAddr: cfg.SOCKS5.GetAddress(),
			Username:  cfg.SOCKS5.Username,
			Password:  cfg.SOCKS5.Password,
			Timeout:   cfg.SOCKS5.Timeout,
		}

		client, err := socks5.NewClient(socks5Cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create SOCKS5 client: %w", err)
		}

		s.socks5Client = client
		log.Printf("[INFO] SOCKS5 proxy enabled: %s", cfg.SOCKS5.GetAddress())

		// Test SOCKS5 connection
		if err := client.Test(); err != nil {
			log.Printf("[WARN] SOCKS5 proxy test failed: %v", err)
		}
	}

	return s, nil
}

// Start starts the proxy server
func (s *Server) Start() error {
	// Create TCP listener
	listener, err := net.Listen("tcp", s.config.Server.GetAddress())
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.listener = listener

	log.Printf("[INFO] Proxy server listening on %s", s.config.Server.GetAddress())

	// Start metrics server if enabled
	if s.config.Metrics.Enabled {
		if err := s.startMetricsServer(); err != nil {
			log.Printf("[WARN] Failed to start metrics server: %v", err)
		}
	}

	// Accept connections
	s.wg.Add(1)
	go s.acceptLoop()

	return nil
}

// acceptLoop accepts incoming connections
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			// Set accept deadline to allow checking shutdown signal
			s.listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))

			conn, err := s.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					// Timeout is expected, continue
					continue
				}
				select {
				case <-s.shutdown:
					return
				default:
					log.Printf("[ERROR] Accept error: %v", err)
					continue
				}
			}

			// Handle connection in goroutine
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				s.handleConnection(conn)
			}()
		}
	}
}

// startMetricsServer starts the metrics HTTP server
func (s *Server) startMetricsServer() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", s.metrics)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	s.metricsServer = &http.Server{
		Addr:    s.config.Metrics.GetAddress(),
		Handler: mux,
	}

	go func() {
		log.Printf("[INFO] Metrics server listening on http://%s/metrics", s.config.Metrics.GetAddress())
		if err := s.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[ERROR] Metrics server error: %v", err)
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("[INFO] Shutting down proxy server...")

	// Signal shutdown
	close(s.shutdown)

	// Close listener
	if s.listener != nil {
		s.listener.Close()
	}

	// Shutdown metrics server
	if s.metricsServer != nil {
		if err := s.metricsServer.Shutdown(ctx); err != nil {
			log.Printf("[WARN] Metrics server shutdown error: %v", err)
		}
	}

	// Wait for all connections to finish with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[INFO] All connections closed gracefully")
	case <-ctx.Done():
		log.Println("[WARN] Shutdown timeout, forcing close")
	}

	return nil
}

// GetMetrics returns the server metrics
func (s *Server) GetMetrics() *Metrics {
	return s.metrics
}
