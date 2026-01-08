package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Server.Port != 8443 {
		t.Errorf("expected default port 8443, got %d", cfg.Server.Port)
	}

	if cfg.Server.BindAddr != "0.0.0.0" {
		t.Errorf("expected default bind address 0.0.0.0, got %s", cfg.Server.BindAddr)
	}

	if cfg.SOCKS5.Enabled {
		t.Error("expected SOCKS5 to be disabled by default")
	}

	if !cfg.SSL.AutoGenerate {
		t.Error("expected SSL auto-generate to be enabled by default")
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("expected default log level info, got %s", cfg.Logging.Level)
	}

	if !cfg.Metrics.Enabled {
		t.Error("expected metrics to be enabled by default")
	}
}

func TestServerConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ServerConfig{
				Port:           8443,
				BindAddr:       "0.0.0.0",
				IdleTimeout:    300 * time.Second,
				MaxConnections: 1000,
			},
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			config: ServerConfig{
				Port:           0,
				BindAddr:       "0.0.0.0",
				IdleTimeout:    300 * time.Second,
				MaxConnections: 1000,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: ServerConfig{
				Port:           70000,
				BindAddr:       "0.0.0.0",
				IdleTimeout:    300 * time.Second,
				MaxConnections: 1000,
			},
			wantErr: true,
		},
		{
			name: "invalid bind address",
			config: ServerConfig{
				Port:           8443,
				BindAddr:       "invalid",
				IdleTimeout:    300 * time.Second,
				MaxConnections: 1000,
			},
			wantErr: true,
		},
		{
			name: "negative max connections",
			config: ServerConfig{
				Port:           8443,
				BindAddr:       "0.0.0.0",
				IdleTimeout:    300 * time.Second,
				MaxConnections: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSOCKS5ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  SOCKS5Config
		wantErr bool
	}{
		{
			name: "valid config with IP",
			config: SOCKS5Config{
				Enabled: true,
				Host:    "127.0.0.1",
				Port:    1080,
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid config with hostname",
			config: SOCKS5Config{
				Enabled: true,
				Host:    "localhost",
				Port:    1080,
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: SOCKS5Config{
				Enabled: true,
				Host:    "127.0.0.1",
				Port:    70000,
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "empty host",
			config: SOCKS5Config{
				Enabled: true,
				Host:    "",
				Port:    1080,
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseSOCKS5URL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantHost string
		wantPort int
		wantUser string
		wantPass string
		wantErr  bool
	}{
		{
			name:     "simple URL",
			url:      "socks5://127.0.0.1:1080",
			wantHost: "127.0.0.1",
			wantPort: 1080,
			wantErr:  false,
		},
		{
			name:     "URL with auth",
			url:      "socks5://user:pass@127.0.0.1:1080",
			wantHost: "127.0.0.1",
			wantPort: 1080,
			wantUser: "user",
			wantPass: "pass",
			wantErr:  false,
		},
		{
			name:     "URL with hostname",
			url:      "socks5://proxy.example.com:9050",
			wantHost: "proxy.example.com",
			wantPort: 9050,
			wantErr:  false,
		},
		{
			name:    "invalid scheme",
			url:     "http://127.0.0.1:1080",
			wantErr: true,
		},
		{
			name:    "missing port",
			url:     "socks5://127.0.0.1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := parseSOCKS5URL(tt.url, cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseSOCKS5URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg.SOCKS5.Host != tt.wantHost {
					t.Errorf("host = %v, want %v", cfg.SOCKS5.Host, tt.wantHost)
				}
				if cfg.SOCKS5.Port != tt.wantPort {
					t.Errorf("port = %v, want %v", cfg.SOCKS5.Port, tt.wantPort)
				}
				if tt.wantUser != "" && cfg.SOCKS5.Username != tt.wantUser {
					t.Errorf("username = %v, want %v", cfg.SOCKS5.Username, tt.wantUser)
				}
				if tt.wantPass != "" && cfg.SOCKS5.Password != tt.wantPass {
					t.Errorf("password = %v, want %v", cfg.SOCKS5.Password, tt.wantPass)
				}
			}
		})
	}
}

func TestSSLConfigValidation(t *testing.T) {
	// Create temp directory for testing
	tmpDir := t.TempDir()

	// Create dummy cert and key files
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")
	_ = os.WriteFile(certFile, []byte("dummy cert"), 0644)
	_ = os.WriteFile(keyFile, []byte("dummy key"), 0644)

	tests := []struct {
		name    string
		config  SSLConfig
		wantErr bool
	}{
		{
			name: "auto-generate enabled",
			config: SSLConfig{
				AutoGenerate: true,
				ValidityDays: 365,
				CacheDir:     tmpDir,
			},
			wantErr: false,
		},
		{
			name: "custom cert - valid files",
			config: SSLConfig{
				AutoGenerate: false,
				CertFile:     certFile,
				KeyFile:      keyFile,
				ValidityDays: 365,
			},
			wantErr: false,
		},
		{
			name: "custom cert - missing cert file",
			config: SSLConfig{
				AutoGenerate: false,
				CertFile:     "",
				KeyFile:      keyFile,
				ValidityDays: 365,
			},
			wantErr: true,
		},
		{
			name: "invalid IP address",
			config: SSLConfig{
				AutoGenerate: true,
				IPAddresses:  []string{"invalid-ip"},
				ValidityDays: 365,
				CacheDir:     tmpDir,
			},
			wantErr: true,
		},
		{
			name: "invalid validity days",
			config: SSLConfig{
				AutoGenerate: true,
				ValidityDays: 5000,
				CacheDir:     tmpDir,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoggingConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  LoggingConfig
		wantErr bool
	}{
		{
			name: "valid - stdout",
			config: LoggingConfig{
				Level:  "info",
				Format: "text",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "valid - json format",
			config: LoggingConfig{
				Level:  "debug",
				Format: "json",
				Output: "stderr",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: LoggingConfig{
				Level:  "invalid",
				Format: "text",
				Output: "stdout",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: LoggingConfig{
				Level:  "info",
				Format: "xml",
				Output: "stdout",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigHelpers(t *testing.T) {
	// Test SOCKS5Config.GetAddress
	socks5 := SOCKS5Config{Host: "127.0.0.1", Port: 1080}
	if addr := socks5.GetAddress(); addr != "127.0.0.1:1080" {
		t.Errorf("GetAddress() = %s, want 127.0.0.1:1080", addr)
	}

	// Test SOCKS5Config.HasAuth
	socks5NoAuth := SOCKS5Config{}
	if socks5NoAuth.HasAuth() {
		t.Error("HasAuth() should be false for empty credentials")
	}

	socks5WithAuth := SOCKS5Config{Username: "user", Password: "pass"}
	if !socks5WithAuth.HasAuth() {
		t.Error("HasAuth() should be true with credentials")
	}

	// Test SSLConfig.GetIPAddresses
	ssl := SSLConfig{IPAddresses: []string{"127.0.0.1", "192.168.1.1", "invalid"}}
	ips := ssl.GetIPAddresses()
	if len(ips) != 2 {
		t.Errorf("GetIPAddresses() returned %d IPs, want 2", len(ips))
	}
}
