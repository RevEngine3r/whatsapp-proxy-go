package socks5

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// Client wraps a SOCKS5 proxy connection with additional features
type Client struct {
	config *Config
	dialer proxy.Dialer
}

// Config holds SOCKS5 client configuration
type Config struct {
	// ProxyAddr is the SOCKS5 proxy address (host:port)
	ProxyAddr string

	// Username for SOCKS5 authentication (optional)
	Username string

	// Password for SOCKS5 authentication (optional)
	Password string

	// Timeout for connection operations
	Timeout time.Duration

	// ForwardDialer is optional dialer to use for upstream connection
	// If nil, uses direct connection
	ForwardDialer proxy.Dialer
}

// NewClient creates a new SOCKS5 client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if cfg.ProxyAddr == "" {
		return nil, fmt.Errorf("proxy address cannot be empty")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	// Create authentication if provided
	var auth *proxy.Auth
	if cfg.Username != "" || cfg.Password != "" {
		auth = &proxy.Auth{
			User:     cfg.Username,
			Password: cfg.Password,
		}
	}

	// Create forward dialer (direct connection if not specified)
	forward := cfg.ForwardDialer
	if forward == nil {
		forward = &net.Dialer{
			Timeout:   cfg.Timeout,
			KeepAlive: 30 * time.Second,
		}
	}

	// Create SOCKS5 dialer using golang.org/x/net/proxy
	// This supports both SOCKS5 and SOCKS5h (hostname resolution on proxy)
	dialer, err := proxy.SOCKS5("tcp", cfg.ProxyAddr, auth, forward)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	return &Client{
		config: cfg,
		dialer: dialer,
	}, nil
}

// NewClientFromURL creates a SOCKS5 client from a URL string
// Format: socks5://[username:password@]host:port
func NewClientFromURL(proxyURL string, timeout time.Duration) (*Client, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	if u.Scheme != "socks5" {
		return nil, fmt.Errorf("invalid scheme: expected socks5, got %s", u.Scheme)
	}

	cfg := &Config{
		ProxyAddr: u.Host,
		Timeout:   timeout,
	}

	// Extract authentication if present
	if u.User != nil {
		cfg.Username = u.User.Username()
		cfg.Password, _ = u.User.Password()
	}

	return NewClient(cfg)
}

// Dial connects to the address on the named network through the SOCKS5 proxy
// Network can be "tcp", "tcp4", or "tcp6"
// Address format is "host:port" where host can be hostname or IP
// When using hostname, DNS resolution happens on the SOCKS5 proxy (SOCKS5h)
func (c *Client) Dial(network, address string) (net.Conn, error) {
	return c.DialContext(context.Background(), network, address)
}

// DialContext connects through SOCKS5 proxy with context support
// This allows for cancellation and timeout control
func (c *Client) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	// Validate network type
	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, fmt.Errorf("unsupported network type: %s (must be tcp, tcp4, or tcp6)", network)
	}

	// Create a channel for the dial result
	type dialResult struct {
		conn net.Conn
		err  error
	}
	resultCh := make(chan dialResult, 1)

	// Dial in goroutine to support context cancellation
	go func() {
		conn, err := c.dialer.Dial(network, address)
		resultCh <- dialResult{conn: conn, err: err}
	}()

	// Wait for dial or context cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("dial cancelled: %w", ctx.Err())
	case result := <-resultCh:
		if result.err != nil {
			return nil, fmt.Errorf("failed to dial through SOCKS5: %w", result.err)
		}
		return result.conn, nil
	}
}

// DialTimeout connects through SOCKS5 proxy with a timeout
func (c *Client) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.DialContext(ctx, network, address)
}

// Test tests the SOCKS5 proxy connection by attempting to connect
// to a well-known test endpoint
func (c *Client) Test() error {
	// Try to connect to a well-known endpoint through the proxy
	// We use google.com:80 as it's reliable and supports TCP
	testEndpoint := "google.com:80"

	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	conn, err := c.DialContext(ctx, "tcp", testEndpoint)
	if err != nil {
		return fmt.Errorf("proxy test failed: %w", err)
	}
	defer conn.Close()

	// Successfully connected
	return nil
}

// GetProxyAddr returns the configured proxy address
func (c *Client) GetProxyAddr() string {
	return c.config.ProxyAddr
}

// HasAuth returns true if authentication is configured
func (c *Client) HasAuth() bool {
	return c.config.Username != "" || c.config.Password != ""
}

// GetDialer returns the underlying proxy.Dialer
// This can be used to create custom http.Transport with the SOCKS5 proxy
func (c *Client) GetDialer() proxy.Dialer {
	return c.dialer
}

// GetConfig returns a copy of the client configuration
func (c *Client) GetConfig() Config {
	// Return a copy to prevent external modification
	return Config{
		ProxyAddr:     c.config.ProxyAddr,
		Username:      c.config.Username,
		Password:      "***", // Redact password
		Timeout:       c.config.Timeout,
		ForwardDialer: c.config.ForwardDialer,
	}
}
