package proxy

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/RevEngine3r/whatsapp-proxy-go/internal/protocol"
)

// Metrics holds proxy server metrics
type Metrics struct {
	// Connection counters
	connectionsTotal    atomic.Uint64
	connectionsActive   atomic.Int64
	connectionsFailed   atomic.Uint64

	// Protocol-specific counters
	httpConnections     atomic.Uint64
	httpsConnections    atomic.Uint64
	jabberConnections   atomic.Uint64
	unknownConnections  atomic.Uint64

	// Data transfer counters
	bytesSent           atomic.Uint64
	bytesReceived       atomic.Uint64

	// Error counters
	errorsTotal         atomic.Uint64

	// Server start time
	startTime           time.Time
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

// IncrementConnections increments active connections
func (m *Metrics) IncrementConnections() {
	m.connectionsTotal.Add(1)
	m.connectionsActive.Add(1)
}

// DecrementConnections decrements active connections
func (m *Metrics) DecrementConnections() {
	m.connectionsActive.Add(-1)
}

// IncrementConnectionsFailed increments failed connection counter
func (m *Metrics) IncrementConnectionsFailed() {
	m.connectionsFailed.Add(1)
}

// IncrementProtocol increments the counter for a specific protocol
func (m *Metrics) IncrementProtocol(proto protocol.Protocol) {
	switch proto {
	case protocol.ProtocolHTTP:
		m.httpConnections.Add(1)
	case protocol.ProtocolHTTPS:
		m.httpsConnections.Add(1)
	case protocol.ProtocolJabber:
		m.jabberConnections.Add(1)
	default:
		m.unknownConnections.Add(1)
	}
}

// AddBytesSent adds to bytes sent counter
func (m *Metrics) AddBytesSent(bytes uint64) {
	m.bytesSent.Add(bytes)
}

// AddBytesReceived adds to bytes received counter
func (m *Metrics) AddBytesReceived(bytes uint64) {
	m.bytesReceived.Add(bytes)
}

// IncrementErrors increments error counter
func (m *Metrics) IncrementErrors() {
	m.errorsTotal.Add(1)
}

// GetConnectionsTotal returns total connections count
func (m *Metrics) GetConnectionsTotal() uint64 {
	return m.connectionsTotal.Load()
}

// GetConnectionsActive returns active connections count
func (m *Metrics) GetConnectionsActive() int64 {
	return m.connectionsActive.Load()
}

// GetUptime returns server uptime
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

// ServeHTTP implements http.Handler for metrics endpoint
func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	// OpenMetrics format
	fmt.Fprintf(w, "# HELP whatsapp_proxy_connections_total Total number of connections\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_connections_total counter\n")
	fmt.Fprintf(w, "whatsapp_proxy_connections_total %d\n", m.connectionsTotal.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_connections_active Number of active connections\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_connections_active gauge\n")
	fmt.Fprintf(w, "whatsapp_proxy_connections_active %d\n", m.connectionsActive.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_connections_failed Total number of failed connections\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_connections_failed counter\n")
	fmt.Fprintf(w, "whatsapp_proxy_connections_failed %d\n", m.connectionsFailed.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_protocol_connections Connections by protocol\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_protocol_connections counter\n")
	fmt.Fprintf(w, "whatsapp_proxy_protocol_connections{protocol=\"http\"} %d\n", m.httpConnections.Load())
	fmt.Fprintf(w, "whatsapp_proxy_protocol_connections{protocol=\"https\"} %d\n", m.httpsConnections.Load())
	fmt.Fprintf(w, "whatsapp_proxy_protocol_connections{protocol=\"jabber\"} %d\n", m.jabberConnections.Load())
	fmt.Fprintf(w, "whatsapp_proxy_protocol_connections{protocol=\"unknown\"} %d\n", m.unknownConnections.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_bytes_sent_total Total bytes sent\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_bytes_sent_total counter\n")
	fmt.Fprintf(w, "whatsapp_proxy_bytes_sent_total %d\n", m.bytesSent.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_bytes_received_total Total bytes received\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_bytes_received_total counter\n")
	fmt.Fprintf(w, "whatsapp_proxy_bytes_received_total %d\n", m.bytesReceived.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_errors_total Total number of errors\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_errors_total counter\n")
	fmt.Fprintf(w, "whatsapp_proxy_errors_total %d\n", m.errorsTotal.Load())
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP whatsapp_proxy_uptime_seconds Server uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE whatsapp_proxy_uptime_seconds gauge\n")
	fmt.Fprintf(w, "whatsapp_proxy_uptime_seconds %.0f\n", m.GetUptime().Seconds())
	fmt.Fprintf(w, "\n")

	// EOF marker for OpenMetrics
	fmt.Fprintf(w, "# EOF\n")
}
