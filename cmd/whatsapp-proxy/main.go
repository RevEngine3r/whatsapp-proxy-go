package main

import (
	"fmt"
	"os"

	"github.com/RevEngine3r/whatsapp-proxy-go/internal/config"
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
  - Single port operation for all protocols
  - Upstream SOCKS5 proxy with authentication
  - Auto-generated SSL certificates
  - Metrics endpoint for monitoring`,
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
	fmt.Printf("WhatsApp Proxy Go v%s\n", Version)
	fmt.Println("===============================================")
	fmt.Printf("Server: %s\n", cfg.Server.GetAddress())
	fmt.Printf("SOCKS5 Proxy: ")
	if cfg.SOCKS5.Enabled {
		fmt.Printf("Enabled (%s)\n", cfg.SOCKS5.GetAddress())
	} else {
		fmt.Println("Disabled")
	}
	fmt.Printf("SSL: Auto-generate=%v\n", cfg.SSL.AutoGenerate)
	fmt.Printf("Log Level: %s\n", cfg.Logging.Level)
	if cfg.Metrics.Enabled {
		fmt.Printf("Metrics: http://%s/metrics\n", cfg.Metrics.GetAddress())
	} else {
		fmt.Println("Metrics: Disabled")
	}
	fmt.Println("===============================================")

	fmt.Println("\n✅ Configuration loaded and validated successfully!")
	fmt.Println("\n🚧 Step 2 Complete - Configuration Management")
	fmt.Println("Next: SOCKS5 client implementation")

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
