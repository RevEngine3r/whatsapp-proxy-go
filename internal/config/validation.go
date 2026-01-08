package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// Validate validates the entire configuration
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if c.SOCKS5.Enabled {
		if err := c.SOCKS5.Validate(); err != nil {
			return fmt.Errorf("socks5 config: %w", err)
		}
	}

	if err := c.SSL.Validate(); err != nil {
		return fmt.Errorf("ssl config: %w", err)
	}

	if err := c.Logging.Validate(); err != nil {
		return fmt.Errorf("logging config: %w", err)
	}

	if c.Metrics.Enabled {
		if err := c.Metrics.Validate(); err != nil {
			return fmt.Errorf("metrics config: %w", err)
		}
	}

	// Check for port conflicts
	if c.Metrics.Enabled && c.Server.Port == c.Metrics.Port {
		return fmt.Errorf("server port and metrics port cannot be the same")
	}

	return nil
}

// Validate validates server configuration
func (c *ServerConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	if c.BindAddr == "" {
		return fmt.Errorf("bind address cannot be empty")
	}

	// Validate bind address is a valid IP
	if ip := net.ParseIP(c.BindAddr); ip == nil {
		return fmt.Errorf("invalid bind address: %s", c.BindAddr)
	}

	if c.IdleTimeout < 0 {
		return fmt.Errorf("idle timeout cannot be negative")
	}

	if c.MaxConnections < 1 {
		return fmt.Errorf("max connections must be at least 1, got %d", c.MaxConnections)
	}

	return nil
}

// Validate validates SOCKS5 configuration
func (c *SOCKS5Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("socks5 host cannot be empty")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("socks5 port must be between 1 and 65535, got %d", c.Port)
	}

	if c.Timeout < 0 {
		return fmt.Errorf("socks5 timeout cannot be negative")
	}

	// Validate that host is resolvable or valid IP
	if ip := net.ParseIP(c.Host); ip == nil {
		// Not an IP, try to resolve hostname
		if _, err := net.LookupHost(c.Host); err != nil {
			return fmt.Errorf("cannot resolve socks5 host %s: %w", c.Host, err)
		}
	}

	return nil
}

// Validate validates SSL configuration
func (c *SSLConfig) Validate() error {
	if !c.AutoGenerate {
		// Custom certificates - must provide both cert and key
		if c.CertFile == "" {
			return fmt.Errorf("cert_file must be specified when auto_generate is false")
		}
		if c.KeyFile == "" {
			return fmt.Errorf("key_file must be specified when auto_generate is false")
		}

		// Check files exist
		if _, err := os.Stat(c.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("cert file does not exist: %s", c.CertFile)
		}
		if _, err := os.Stat(c.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("key file does not exist: %s", c.KeyFile)
		}
	} else {
		// Auto-generate - validate cache directory
		if c.CacheDir != "" {
			// Try to create cache directory if it doesn't exist
			if err := os.MkdirAll(c.CacheDir, 0755); err != nil {
				return fmt.Errorf("cannot create cache directory %s: %w", c.CacheDir, err)
			}
		}
	}

	// Validate IP addresses
	for _, ipStr := range c.IPAddresses {
		if ip := net.ParseIP(ipStr); ip == nil {
			return fmt.Errorf("invalid IP address: %s", ipStr)
		}
	}

	if c.ValidityDays < 1 || c.ValidityDays > 3650 {
		return fmt.Errorf("validity days must be between 1 and 3650, got %d", c.ValidityDays)
	}

	return nil
}

// Validate validates logging configuration
func (c *LoggingConfig) Validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[c.Level] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Level)
	}

	validFormats := map[string]bool{
		"text": true,
		"json": true,
	}

	if !validFormats[c.Format] {
		return fmt.Errorf("invalid log format: %s (must be text or json)", c.Format)
	}

	// Validate output
	if c.Output != "stdout" && c.Output != "stderr" {
		// Assume it's a file path - validate directory exists
		dir := filepath.Dir(c.Output)
		if dir != "." && dir != "" {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return fmt.Errorf("log output directory does not exist: %s", dir)
			}
		}
	}

	return nil
}

// Validate validates metrics configuration
func (c *MetricsConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("metrics port must be between 1 and 65535, got %d", c.Port)
	}

	if c.BindAddr == "" {
		return fmt.Errorf("metrics bind address cannot be empty")
	}

	// Validate bind address is a valid IP
	if ip := net.ParseIP(c.BindAddr); ip == nil {
		return fmt.Errorf("invalid metrics bind address: %s", c.BindAddr)
	}

	return nil
}
