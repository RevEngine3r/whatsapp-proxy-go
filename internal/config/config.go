package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config holds all configuration for the proxy server
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	SOCKS5  SOCKS5Config  `mapstructure:"socks5"`
	SSL     SSLConfig     `mapstructure:"ssl"`
	Logging LoggingConfig `mapstructure:"logging"`
	Metrics MetricsConfig `mapstructure:"metrics"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port           int           `mapstructure:"port"`
	BindAddr       string        `mapstructure:"bind_addr"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	MaxConnections int           `mapstructure:"max_connections"`
}

// SOCKS5Config holds SOCKS5 upstream proxy settings
type SOCKS5Config struct {
	Enabled  bool          `mapstructure:"enabled"`
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// SSLConfig holds SSL/TLS certificate settings
type SSLConfig struct {
	AutoGenerate bool     `mapstructure:"auto_generate"`
	CertFile     string   `mapstructure:"cert_file"`
	KeyFile      string   `mapstructure:"key_file"`
	DNSNames     []string `mapstructure:"dns_names"`
	IPAddresses  []string `mapstructure:"ip_addresses"`
	ValidityDays int      `mapstructure:"validity_days"`
	CacheDir     string   `mapstructure:"cache_dir"`
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// MetricsConfig holds metrics endpoint settings
type MetricsConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Port     int    `mapstructure:"port"`
	BindAddr string `mapstructure:"bind_addr"`
}

// Default returns a Config with sensible defaults
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".whatsapp-proxy", "certs")

	return &Config{
		Server: ServerConfig{
			Port:           8443,
			BindAddr:       "0.0.0.0",
			IdleTimeout:    300 * time.Second,
			MaxConnections: 1000,
		},
		SOCKS5: SOCKS5Config{
			Enabled: false,
			Host:    "127.0.0.1",
			Port:    1080,
			Timeout: 30 * time.Second,
		},
		SSL: SSLConfig{
			AutoGenerate: true,
			DNSNames:     []string{"localhost"},
			IPAddresses:  []string{"127.0.0.1"},
			ValidityDays: 365,
			CacheDir:     cacheDir,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Metrics: MetricsConfig{
			Enabled:  true,
			Port:     8199,
			BindAddr: "127.0.0.1",
		},
	}
}

// Load loads configuration from file and CLI flags
func Load(cmd *cobra.Command) (*Config, error) {
	cfg := Default()

	// Set up viper
	v := viper.New()
	v.SetConfigType("yaml")

	// Bind environment variables
	v.SetEnvPrefix("WHATSAPP_PROXY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Load config file if specified
	configFile, _ := cmd.Flags().GetString("config")
	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Search for config in common locations
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.whatsapp-proxy")
		v.AddConfigPath("/etc/whatsapp-proxy")
		_ = v.ReadInConfig() // Ignore error if no config file found
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with CLI flags if set
	if err := overrideFromFlags(cmd, cfg); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// overrideFromFlags overrides config values with CLI flags
func overrideFromFlags(cmd *cobra.Command, cfg *Config) error {
	if cmd.Flags().Changed("port") {
		port, _ := cmd.Flags().GetInt("port")
		cfg.Server.Port = port
	}

	if cmd.Flags().Changed("bind") {
		bind, _ := cmd.Flags().GetString("bind")
		cfg.Server.BindAddr = bind
	}

	if cmd.Flags().Changed("socks5-proxy") {
		proxyURL, _ := cmd.Flags().GetString("socks5-proxy")
		if err := parseSOCKS5URL(proxyURL, cfg); err != nil {
			return fmt.Errorf("invalid SOCKS5 proxy URL: %w", err)
		}
		cfg.SOCKS5.Enabled = true
	}

	if cmd.Flags().Changed("log-level") {
		level, _ := cmd.Flags().GetString("log-level")
		cfg.Logging.Level = level
	}

	if cmd.Flags().Changed("metrics-port") {
		port, _ := cmd.Flags().GetInt("metrics-port")
		cfg.Metrics.Port = port
	}

	if cmd.Flags().Changed("disable-metrics") {
		disabled, _ := cmd.Flags().GetBool("disable-metrics")
		cfg.Metrics.Enabled = !disabled
	}

	return nil
}

// parseSOCKS5URL parses a SOCKS5 proxy URL
// Format: socks5://[username:password@]host:port
func parseSOCKS5URL(proxyURL string, cfg *Config) error {
	if !strings.HasPrefix(proxyURL, "socks5://") {
		return fmt.Errorf("proxy URL must start with socks5://")
	}

	// Remove scheme
	rest := strings.TrimPrefix(proxyURL, "socks5://")

	// Check for auth
	var auth, hostPort string
	if idx := strings.Index(rest, "@"); idx != -1 {
		auth = rest[:idx]
		hostPort = rest[idx+1:]

		// Parse username:password
		parts := strings.SplitN(auth, ":", 2)
		if len(parts) == 2 {
			cfg.SOCKS5.Username = parts[0]
			cfg.SOCKS5.Password = parts[1]
		}
	} else {
		hostPort = rest
	}

	// Parse host:port
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return fmt.Errorf("invalid host:port format: %w", err)
	}

	cfg.SOCKS5.Host = host

	// Parse port
	var portNum int
	if _, err := fmt.Sscanf(port, "%d", &portNum); err != nil {
		return fmt.Errorf("invalid port number: %w", err)
	}
	cfg.SOCKS5.Port = portNum

	return nil
}

// GetSOCKS5Address returns the full SOCKS5 proxy address
func (c *SOCKS5Config) GetAddress() string {
	return net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))
}

// HasAuth returns true if SOCKS5 authentication is configured
func (c *SOCKS5Config) HasAuth() bool {
	return c.Username != "" || c.Password != ""
}

// GetIPAddresses parses IP address strings into net.IP
func (c *SSLConfig) GetIPAddresses() []net.IP {
	var ips []net.IP
	for _, ipStr := range c.IPAddresses {
		if ip := net.ParseIP(ipStr); ip != nil {
			ips = append(ips, ip)
		}
	}
	return ips
}

// GetServerAddress returns the server listen address
func (c *ServerConfig) GetAddress() string {
	return net.JoinHostPort(c.BindAddr, fmt.Sprintf("%d", c.Port))
}

// GetMetricsAddress returns the metrics listen address
func (c *MetricsConfig) GetAddress() string {
	return net.JoinHostPort(c.BindAddr, fmt.Sprintf("%d", c.Port))
}
