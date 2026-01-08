package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Protocol represents the detected protocol type
type Protocol int

const (
	// ProtocolUnknown represents an unidentified protocol
	ProtocolUnknown Protocol = iota
	// ProtocolHTTP represents HTTP protocol
	ProtocolHTTP
	// ProtocolHTTPS represents HTTPS/TLS protocol
	ProtocolHTTPS
	// ProtocolJabber represents Jabber/XMPP protocol
	ProtocolJabber
)

// String returns the string representation of the protocol
func (p Protocol) String() string {
	switch p {
	case ProtocolHTTP:
		return "HTTP"
	case ProtocolHTTPS:
		return "HTTPS"
	case ProtocolJabber:
		return "Jabber"
	default:
		return "Unknown"
	}
}

// Detector detects the protocol from initial connection bytes
type Detector struct {
	buffer []byte
}

// NewDetector creates a new protocol detector
func NewDetector() *Detector {
	return &Detector{
		buffer: make([]byte, 0, 1024),
	}
}

// Detect attempts to detect the protocol from the reader
// It peeks at the first few bytes without consuming them
func Detect(reader *bufio.Reader) (Protocol, error) {
	// Peek at first few bytes (up to 32 bytes should be enough)
	peek, err := reader.Peek(32)
	if err != nil && err != io.EOF {
		return ProtocolUnknown, fmt.Errorf("failed to peek bytes: %w", err)
	}

	if len(peek) == 0 {
		return ProtocolUnknown, fmt.Errorf("no data to detect")
	}

	// Check for TLS/HTTPS (starts with 0x16 for TLS handshake)
	if peek[0] == 0x16 {
		// TLS handshake record
		if len(peek) > 1 && peek[1] == 0x03 {
			// SSL/TLS version (0x03 0x01 = TLS 1.0, 0x03 0x03 = TLS 1.2, etc.)
			return ProtocolHTTPS, nil
		}
	}

	// Check for HTTP methods
	httpMethods := []string{
		"GET ",
		"POST ",
		"PUT ",
		"DELETE ",
		"HEAD ",
		"OPTIONS ",
		"CONNECT ",
		"PATCH ",
		"TRACE ",
	}

	for _, method := range httpMethods {
		if bytes.HasPrefix(peek, []byte(method)) {
			return ProtocolHTTP, nil
		}
	}

	// Check for Jabber/XMPP (starts with <?xml or <stream)
	if bytes.HasPrefix(peek, []byte("<?xml")) ||
		bytes.HasPrefix(peek, []byte("<stream:stream")) ||
		bytes.HasPrefix(peek, []byte("<stream")) {
		return ProtocolJabber, nil
	}

	return ProtocolUnknown, nil
}

// DetectFromBytes detects protocol from a byte slice
func DetectFromBytes(data []byte) Protocol {
	if len(data) == 0 {
		return ProtocolUnknown
	}

	// Check for TLS/HTTPS
	if data[0] == 0x16 && len(data) > 1 && data[1] == 0x03 {
		return ProtocolHTTPS
	}

	// Check for HTTP methods
	httpMethods := []string{
		"GET ", "POST ", "PUT ", "DELETE ", "HEAD ",
		"OPTIONS ", "CONNECT ", "PATCH ", "TRACE ",
	}

	for _, method := range httpMethods {
		if bytes.HasPrefix(data, []byte(method)) {
			return ProtocolHTTP
		}
	}

	// Check for Jabber/XMPP
	if bytes.HasPrefix(data, []byte("<?xml")) ||
		bytes.HasPrefix(data, []byte("<stream:stream")) ||
		bytes.HasPrefix(data, []byte("<stream")) {
		return ProtocolJabber
	}

	return ProtocolUnknown
}
