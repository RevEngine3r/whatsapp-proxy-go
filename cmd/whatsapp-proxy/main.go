package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RevEngine3r/whatsapp-proxy-go/internal/config"
	"github.com/RevEngine3r/whatsapp-proxy-go/internal/proxy"
	"github.com/spf13/cobra"
)

// Version information (set via ldflags)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "whatsapp-proxy",
	Short: "WhatsApp Proxy Server",
	Long: `A lightweight, cross-platform WhatsApp proxy server written in Go.

Supports:
  - Single port operation for all protocols (HTTP, HTTPS, Jabber)
  - Upstream SOCKS5 proxy with authentication
  - Auto-generated SSL certificates
  - Metrics endpoint for monitoring
  - Protocol detection and routing`,
	Version: Version,
	RunE:    run,
}

func init() {
	// Server flags
	rootCmd.Flags().IntP("port", "p", 8443, "Server port")
	rootCmd.Flags().String("bind", "0.0.0.0", "Bind address")

	// Config file
	rootCmd.Flags().StringP("config", "c", "", "Config file path")

	// SOCKS5 proxy
	rootCmd.Flags().String("socks5-proxy", "", "Upstream SOCKS5 proxy (format: socks5://[user:pass@]host:port)")

	// Logging
	rootCmd.Flags().String("log-level", "info", "Log level (debug, info, warn, error)")

	// Metrics
	rootCmd.Flags().Int("metrics-port", 8199, "Metrics endpoint port")
	rootCmd.Flags().Bool("disable-metrics", false, "Disable metrics endpoint")

	// Version template
	rootCmd.SetVersionTemplate(fmt.Sprintf(
		"WhatsApp Proxy Go v%s\nBuild Time: %s\nGit Commit: %s\n",
		Version, BuildTime, GitCommit,
	))
}

func run(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(cmd)
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Display configuration summary
	printBanner()
	printConfig(cfg)

	// Create proxy server
	server, err := proxy.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Start server
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	log.Println("[INFO] Server started successfully")
	log.Println("[INFO] Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[INFO] Interrupt received, shutting down...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Println("[INFO] Server stopped gracefully")
	return nil
}

func printBanner() {
	fmt.Println()
	fmt.Println("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
	fmt.Println("┃  WhatsApp Proxy Go                    ┃")
	fmt.Printf("┃  Version: %-28s ┃\n", Version)
	fmt.Println("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")
	fmt.Println()
}

func printConfig(cfg *config.Config) {
	fmt.Println("🚀 Configuration:")
	fmt.Println("===============================================")
	fmt.Printf("🎯 Server:        %s\n", cfg.Server.GetAddress())

	fmt.Printf("🔌 SOCKS5 Proxy:  ")
	if cfg.SOCKS5.Enabled {
		fmt.Printf("Enabled (%s)\n", cfg.SOCKS5.GetAddress())
		if cfg.SOCKS5.HasAuth() {
			fmt.Println("             with authentication")
		}
	} else {
		fmt.Println("Disabled (direct connection)")
	}

	fmt.Printf("🔐 SSL:           Auto-generate=%v\n", cfg.SSL.AutoGenerate)
	fmt.Printf("📝 Log Level:     %s\n", cfg.Logging.Level)

	if cfg.Metrics.Enabled {
		fmt.Printf("📊 Metrics:       http://%s/metrics\n", cfg.Metrics.GetAddress())
		fmt.Printf("             http://%s/health\n", cfg.Metrics.GetAddress())
	} else {
		fmt.Println("📊 Metrics:       Disabled")
	}

	fmt.Println("===============================================")
	fmt.Println()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
