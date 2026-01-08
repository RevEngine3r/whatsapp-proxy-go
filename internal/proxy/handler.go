package proxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/RevEngine3r/whatsapp-proxy-go/internal/protocol"
)

// handleConnection handles an incoming connection
func (s *Server) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	s.metrics.IncrementConnections()
	defer s.metrics.DecrementConnections()

	// Set read deadline for protocol detection
	clientConn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Wrap connection in buffered reader for protocol detection
	reader := bufio.NewReader(clientConn)

	// Detect protocol
	proto, err := protocol.Detect(reader)
	if err != nil {
		s.logError("protocol detection failed", err)
		s.metrics.IncrementErrors()
		s.metrics.IncrementConnectionsFailed()
		return
	}

	s.metrics.IncrementProtocol(proto)
	s.logInfo(fmt.Sprintf("detected protocol: %s from %s", proto, clientConn.RemoteAddr()))

	// Remove read deadline for actual data transfer
	clientConn.SetReadDeadline(time.Time{})

	// Route to appropriate handler
	switch proto {
	case protocol.ProtocolHTTP:
		s.handleHTTP(clientConn, reader)
	case protocol.ProtocolHTTPS:
		s.handleHTTPS(clientConn, reader)
	case protocol.ProtocolJabber:
		s.handleJabber(clientConn, reader)
	default:
		s.handleUnknown(clientConn, reader)
	}
}

// handleHTTP handles HTTP protocol connections
func (s *Server) handleHTTP(clientConn net.Conn, reader *bufio.Reader) {
	// Parse HTTP request
	req, err := http.ReadRequest(reader)
	if err != nil {
		s.logError("failed to read HTTP request", err)
		s.metrics.IncrementErrors()
		return
	}

	s.logInfo(fmt.Sprintf("HTTP %s %s", req.Method, req.RequestURI))

	// Handle CONNECT method (for HTTPS tunneling)
	if req.Method == http.MethodConnect {
		s.handleHTTPConnect(clientConn, req)
		return
	}

	// For other HTTP methods, establish connection to target
	target := req.Host
	if target == "" {
		target = req.URL.Host
	}
	if !strings.Contains(target, ":") {
		target += ":80"
	}

	// Connect to upstream
	upstreamConn, err := s.dialUpstream("tcp", target)
	if err != nil {
		s.logError(fmt.Sprintf("failed to connect to %s", target), err)
		s.metrics.IncrementErrors()
		http.Error(clientConn.(*net.TCPConn), "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer upstreamConn.Close()

	// Write request to upstream
	if err := req.Write(upstreamConn); err != nil {
		s.logError("failed to write request to upstream", err)
		s.metrics.IncrementErrors()
		return
	}

	// Bidirectional copy
	s.bidirectionalCopy(clientConn, upstreamConn)
}

// handleHTTPConnect handles HTTP CONNECT method (HTTPS tunneling)
func (s *Server) handleHTTPConnect(clientConn net.Conn, req *http.Request) {
	target := req.Host
	if !strings.Contains(target, ":") {
		target += ":443"
	}

	s.logInfo(fmt.Sprintf("CONNECT tunnel to %s", target))

	// Connect to upstream
	upstreamConn, err := s.dialUpstream("tcp", target)
	if err != nil {
		s.logError(fmt.Sprintf("failed to connect to %s", target), err)
		s.metrics.IncrementErrors()
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer upstreamConn.Close()

	// Send success response
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Bidirectional copy
	s.bidirectionalCopy(clientConn, upstreamConn)
}

// handleHTTPS handles HTTPS/TLS protocol connections
func (s *Server) handleHTTPS(clientConn net.Conn, reader *bufio.Reader) {
	// For HTTPS, we need SNI (Server Name Indication) from TLS ClientHello
	// For now, we'll just proxy the TLS handshake through
	// In a full implementation, you'd extract SNI from the ClientHello

	// Read the first peek to get destination (simplified)
	// In production, parse TLS ClientHello for SNI
	s.logInfo("HTTPS connection - proxying TLS handshake")

	// For WhatsApp, common endpoints are:
	// - e*.whatsapp.net:443
	// - web.whatsapp.com:443
	// Since we can't determine the exact host without SNI parsing,
	// we'll need to handle this differently in production

	// For now, proxy to a default WhatsApp endpoint
	target := "web.whatsapp.com:443"

	upstreamConn, err := s.dialUpstream("tcp", target)
	if err != nil {
		s.logError(fmt.Sprintf("failed to connect to %s", target), err)
		s.metrics.IncrementErrors()
		return
	}
	defer upstreamConn.Close()

	// Copy any buffered data first
	if reader.Buffered() > 0 {
		buffered := make([]byte, reader.Buffered())
		reader.Read(buffered)
		upstreamConn.Write(buffered)
	}

	// Bidirectional copy
	s.bidirectionalCopy(clientConn, upstreamConn)
}

// handleJabber handles Jabber/XMPP protocol connections
func (s *Server) handleJabber(clientConn net.Conn, reader *bufio.Reader) {
	s.logInfo("Jabber/XMPP connection")

	// WhatsApp uses Jabber protocol on port 5222
	// Connect to WhatsApp's Jabber server
	target := "e1.whatsapp.net:5222"

	upstreamConn, err := s.dialUpstream("tcp", target)
	if err != nil {
		s.logError(fmt.Sprintf("failed to connect to %s", target), err)
		s.metrics.IncrementErrors()
		return
	}
	defer upstreamConn.Close()

	// Copy any buffered data first
	if reader.Buffered() > 0 {
		buffered := make([]byte, reader.Buffered())
		reader.Read(buffered)
		upstreamConn.Write(buffered)
	}

	// Bidirectional copy
	s.bidirectionalCopy(clientConn, upstreamConn)
}

// handleUnknown handles unknown protocol connections
func (s *Server) handleUnknown(clientConn net.Conn, reader *bufio.Reader) {
	s.logInfo("unknown protocol - attempting transparent proxy")

	// For unknown protocols, we can't determine the destination
	// This would need additional configuration or rules
	s.metrics.IncrementErrors()
	s.logError("cannot proxy unknown protocol without destination", nil)
}

// bidirectionalCopy copies data bidirectionally between two connections
func (s *Server) bidirectionalCopy(conn1, conn2 net.Conn) {
	done := make(chan struct{}, 2)

	// Copy from conn1 to conn2
	go func() {
		defer func() { done <- struct{}{} }()
		n, err := io.Copy(conn2, conn1)
		if err != nil && err != io.EOF {
			s.logError("copy error (client->upstream)", err)
		}
		s.metrics.AddBytesReceived(uint64(n))
	}()

	// Copy from conn2 to conn1
	go func() {
		defer func() { done <- struct{}{} }()
		n, err := io.Copy(conn1, conn2)
		if err != nil && err != io.EOF {
			s.logError("copy error (upstream->client)", err)
		}
		s.metrics.AddBytesSent(uint64(n))
	}()

	// Wait for both directions to complete
	<-done
	<-done
}

// dialUpstream dials the upstream server, optionally through SOCKS5
func (s *Server) dialUpstream(network, address string) (net.Conn, error) {
	if s.socks5Client != nil {
		// Dial through SOCKS5 proxy
		return s.socks5Client.DialTimeout(network, address, 30*time.Second)
	}

	// Direct connection
	return net.DialTimeout(network, address, 30*time.Second)
}

// Logging helpers
func (s *Server) logInfo(msg string) {
	if s.config.Logging.Level == "debug" || s.config.Logging.Level == "info" {
		log.Printf("[INFO] %s", msg)
	}
}

func (s *Server) logError(msg string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", msg, err)
	} else {
		log.Printf("[ERROR] %s", msg)
	}
}
