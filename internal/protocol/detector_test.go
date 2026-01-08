package protocol

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDetectHTTP(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected Protocol
	}{
		{
			name:     "GET request",
			data:     []byte("GET / HTTP/1.1\r\nHost: example.com\r\n"),
			expected: ProtocolHTTP,
		},
		{
			name:     "POST request",
			data:     []byte("POST /api HTTP/1.1\r\n"),
			expected: ProtocolHTTP,
		},
		{
			name:     "CONNECT request",
			data:     []byte("CONNECT example.com:443 HTTP/1.1\r\n"),
			expected: ProtocolHTTP,
		},
		{
			name:     "OPTIONS request",
			data:     []byte("OPTIONS * HTTP/1.1\r\n"),
			expected: ProtocolHTTP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader(tt.data))
			protocol, err := Detect(reader)
			if err != nil {
				t.Errorf("Detect() error = %v", err)
			}
			if protocol != tt.expected {
				t.Errorf("Detect() = %v, want %v", protocol, tt.expected)
			}
		})
	}
}

func TestDetectHTTPS(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected Protocol
	}{
		{
			name:     "TLS 1.0 handshake",
			data:     []byte{0x16, 0x03, 0x01, 0x00, 0x00},
			expected: ProtocolHTTPS,
		},
		{
			name:     "TLS 1.2 handshake",
			data:     []byte{0x16, 0x03, 0x03, 0x00, 0x00},
			expected: ProtocolHTTPS,
		},
		{
			name:     "TLS 1.3 handshake",
			data:     []byte{0x16, 0x03, 0x04, 0x00, 0x00},
			expected: ProtocolHTTPS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader(tt.data))
			protocol, err := Detect(reader)
			if err != nil {
				t.Errorf("Detect() error = %v", err)
			}
			if protocol != tt.expected {
				t.Errorf("Detect() = %v, want %v", protocol, tt.expected)
			}
		})
	}
}

func TestDetectJabber(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected Protocol
	}{
		{
			name:     "XML declaration",
			data:     []byte("<?xml version='1.0'?>"),
			expected: ProtocolJabber,
		},
		{
			name:     "Stream opening with namespace",
			data:     []byte("<stream:stream xmlns='jabber:client'>"),
			expected: ProtocolJabber,
		},
		{
			name:     "Stream opening simple",
			data:     []byte("<stream>"),
			expected: ProtocolJabber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader(tt.data))
			protocol, err := Detect(reader)
			if err != nil {
				t.Errorf("Detect() error = %v", err)
			}
			if protocol != tt.expected {
				t.Errorf("Detect() = %v, want %v", protocol, tt.expected)
			}
		})
	}
}

func TestDetectUnknown(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "random binary data",
			data: []byte{0xFF, 0xFE, 0xFD, 0xFC},
		},
		{
			name: "plain text",
			data: []byte("Hello World"),
		},
		{
			name: "empty data",
			data: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.data) == 0 {
				// Empty data should return error
				reader := bufio.NewReader(bytes.NewReader(tt.data))
				_, err := Detect(reader)
				if err == nil {
					t.Error("Detect() should return error for empty data")
				}
			} else {
				reader := bufio.NewReader(bytes.NewReader(tt.data))
				protocol, err := Detect(reader)
				if err != nil {
					t.Errorf("Detect() error = %v", err)
				}
				if protocol != ProtocolUnknown {
					t.Errorf("Detect() = %v, want Unknown", protocol)
				}
			}
		})
	}
}

func TestDetectFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected Protocol
	}{
		{
			name:     "HTTP GET",
			data:     []byte("GET / HTTP/1.1"),
			expected: ProtocolHTTP,
		},
		{
			name:     "HTTPS TLS",
			data:     []byte{0x16, 0x03, 0x01},
			expected: ProtocolHTTPS,
		},
		{
			name:     "Jabber XML",
			data:     []byte("<?xml version='1.0'?>"),
			expected: ProtocolJabber,
		},
		{
			name:     "Unknown",
			data:     []byte{0xFF, 0xFE},
			expected: ProtocolUnknown,
		},
		{
			name:     "Empty",
			data:     []byte{},
			expected: ProtocolUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protocol := DetectFromBytes(tt.data)
			if protocol != tt.expected {
				t.Errorf("DetectFromBytes() = %v, want %v", protocol, tt.expected)
			}
		})
	}
}

func TestProtocolString(t *testing.T) {
	tests := []struct {
		protocol Protocol
		expected string
	}{
		{ProtocolHTTP, "HTTP"},
		{ProtocolHTTPS, "HTTPS"},
		{ProtocolJabber, "Jabber"},
		{ProtocolUnknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.protocol.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkDetect(b *testing.B) {
	data := []byte("GET / HTTP/1.1\r\nHost: example.com\r\n")
	reader := bufio.NewReader(bytes.NewReader(data))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Reset(bytes.NewReader(data))
		Detect(reader)
	}
}

func BenchmarkDetectFromBytes(b *testing.B) {
	data := []byte("GET / HTTP/1.1\r\nHost: example.com\r\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectFromBytes(data)
	}
}
