package socks5

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config without auth",
			config: &Config{
				ProxyAddr: "127.0.0.1:1080",
				Timeout:   30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid config with auth",
			config: &Config{
				ProxyAddr: "127.0.0.1:1080",
				Username:  "user",
				Password:  "pass",
				Timeout:   30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "default timeout applied",
			config: &Config{
				ProxyAddr: "127.0.0.1:1080",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "empty proxy address",
			config: &Config{
				ProxyAddr: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
			if !tt.wantErr && tt.config.Timeout == 0 {
				if client.config.Timeout != 30*time.Second {
					t.Errorf("Default timeout not applied: got %v, want 30s", client.config.Timeout)
				}
			}
		})
	}
}

func TestNewClientFromURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		timeout      time.Duration
		wantAddr     string
		wantUsername string
		wantPassword string
		wantErr      bool
	}{
		{
			name:     "simple URL without auth",
			url:      "socks5://127.0.0.1:1080",
			timeout:  30 * time.Second,
			wantAddr: "127.0.0.1:1080",
			wantErr:  false,
		},
		{
			name:         "URL with auth",
			url:          "socks5://user:pass@127.0.0.1:1080",
			timeout:      30 * time.Second,
			wantAddr:     "127.0.0.1:1080",
			wantUsername: "user",
			wantPassword: "pass",
			wantErr:      false,
		},
		{
			name:     "URL with hostname",
			url:      "socks5://proxy.example.com:9050",
			timeout:  30 * time.Second,
			wantAddr: "proxy.example.com:9050",
			wantErr:  false,
		},
		{
			name:    "invalid scheme",
			url:     "http://127.0.0.1:1080",
			timeout: 30 * time.Second,
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "not a valid url",
			timeout: 30 * time.Second,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientFromURL(tt.url, tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if client.config.ProxyAddr != tt.wantAddr {
					t.Errorf("ProxyAddr = %v, want %v", client.config.ProxyAddr, tt.wantAddr)
				}
				if tt.wantUsername != "" && client.config.Username != tt.wantUsername {
					t.Errorf("Username = %v, want %v", client.config.Username, tt.wantUsername)
				}
				if tt.wantPassword != "" && client.config.Password != tt.wantPassword {
					t.Errorf("Password = %v, want %v", client.config.Password, tt.wantPassword)
				}
			}
		})
	}
}

func TestClientHelpers(t *testing.T) {
	// Test GetProxyAddr
	client, _ := NewClient(&Config{
		ProxyAddr: "127.0.0.1:1080",
		Timeout:   30 * time.Second,
	})

	if addr := client.GetProxyAddr(); addr != "127.0.0.1:1080" {
		t.Errorf("GetProxyAddr() = %v, want 127.0.0.1:1080", addr)
	}

	// Test HasAuth - no auth
	if client.HasAuth() {
		t.Error("HasAuth() should return false when no auth configured")
	}

	// Test HasAuth - with auth
	clientWithAuth, _ := NewClient(&Config{
		ProxyAddr: "127.0.0.1:1080",
		Username:  "user",
		Password:  "pass",
		Timeout:   30 * time.Second,
	})

	if !clientWithAuth.HasAuth() {
		t.Error("HasAuth() should return true when auth configured")
	}

	// Test GetConfig (password should be redacted)
	cfg := clientWithAuth.GetConfig()
	if cfg.Password != "***" {
		t.Errorf("GetConfig() should redact password, got %v", cfg.Password)
	}
	if cfg.Username != "user" {
		t.Errorf("GetConfig() Username = %v, want user", cfg.Username)
	}

	// Test GetDialer
	if dialer := client.GetDialer(); dialer == nil {
		t.Error("GetDialer() returned nil")
	}
}

func TestDialContextValidation(t *testing.T) {
	client, _ := NewClient(&Config{
		ProxyAddr: "127.0.0.1:1080",
		Timeout:   30 * time.Second,
	})

	tests := []struct {
		name    string
		network string
		wantErr bool
	}{
		{
			name:    "valid tcp",
			network: "tcp",
			wantErr: false,
		},
		{
			name:    "valid tcp4",
			network: "tcp4",
			wantErr: false,
		},
		{
			name:    "valid tcp6",
			network: "tcp6",
			wantErr: false,
		},
		{
			name:    "invalid udp",
			network: "udp",
			wantErr: true,
		},
		{
			name:    "invalid unix",
			network: "unix",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a very short timeout and invalid address to fail fast
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer cancel()

			_, err := client.DialContext(ctx, tt.network, "invalid:12345")

			// For valid networks, we expect connection failure (not validation error)
			// For invalid networks, we expect immediate validation error
			if tt.wantErr {
				if err == nil {
					t.Error("DialContext() should return error for invalid network")
				} else if err.Error()[:29] != "unsupported network type" {
					t.Errorf("Expected unsupported network error, got: %v", err)
				}
			}
		})
	}
}

func TestDialContextCancellation(t *testing.T) {
	client, _ := NewClient(&Config{
		ProxyAddr: "127.0.0.1:1080",
		Timeout:   30 * time.Second,
	})

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.DialContext(ctx, "tcp", "example.com:80")
	if err == nil {
		t.Error("DialContext() should return error for cancelled context")
	}
}

// TestDialWithMockSOCKS5 tests dial functionality with a mock SOCKS5 server
func TestDialWithMockSOCKS5(t *testing.T) {
	// Start a simple TCP echo server to act as the target
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()

	// Start echo server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				buf := make([]byte, 1024)
				n, _ := c.Read(buf)
				if n > 0 {
					c.Write(buf[:n])
				}
			}(conn)
		}
	}()

	// Note: This test requires a real SOCKS5 proxy running on localhost:1080
	// In production tests, you would use a mock SOCKS5 server
	// For now, we just validate that the client can be created
	client, err := NewClient(&Config{
		ProxyAddr: "127.0.0.1:1080",
		Timeout:   5 * time.Second,
	})

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// We can't actually dial through the proxy without a real SOCKS5 server
	// but we can verify the client is properly configured
	if client.GetProxyAddr() != "127.0.0.1:1080" {
		t.Errorf("Unexpected proxy address: %v", client.GetProxyAddr())
	}

	t.Logf("Server started on %s (for manual testing)", serverAddr)
}

func TestDialTimeout(t *testing.T) {
	client, _ := NewClient(&Config{
		ProxyAddr: "127.0.0.1:1080",
		Timeout:   30 * time.Second,
	})

	// Try to dial with a very short timeout to an unreachable address
	_, err := client.DialTimeout("tcp", "192.0.2.1:12345", 1*time.Millisecond)
	if err == nil {
		t.Error("DialTimeout() should fail with short timeout")
	}
}

// Benchmark tests
func BenchmarkNewClient(b *testing.B) {
	cfg := &Config{
		ProxyAddr: "127.0.0.1:1080",
		Username:  "user",
		Password:  "pass",
		Timeout:   30 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(cfg)
	}
}

func BenchmarkNewClientFromURL(b *testing.B) {
	url := "socks5://user:pass@127.0.0.1:1080"
	timeout := 30 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewClientFromURL(url, timeout)
	}
}
