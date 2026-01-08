# Step 4: Proxy Server Core

## Objective
Implement main proxy server with protocol detection, single port operation, request routing, and upstream forwarding.

## Tasks

### 1. Protocol Detection
- Detect HTTP (plaintext)
- Detect HTTPS (TLS handshake)
- Detect Jabber/XMPP (port 5222 protocol)
- Fallback handling

### 2. Single Port Listener
- TCP listener on configured port
- Accept connections
- Initial bytes inspection
- Route to appropriate handler

### 3. HTTP/HTTPS Handler
- HTTP CONNECT method for tunneling
- Direct HTTP proxy
- TLS termination and re-encryption
- Header manipulation

### 4. Jabber/XMPP Handler
- XMPP stream detection
- Transparent proxying
- Connection persistence

### 5. Upstream Forwarding
- Connect via SOCKS5 client
- Bidirectional data copy
- Connection lifecycle management
- Graceful shutdown

### 6. Metrics Endpoint
- HTTP metrics server
- Connection statistics
- Error counters
- OpenMetrics format

## Implementation Details

### Server Structure
```go
type Server struct {
    config      *config.Config
    listener    net.Listener
    socks5      *socks5.Client
    sslManager  *ssl.Manager
    metrics     *Metrics
    logger      *log.Logger
}

func NewServer(cfg *config.Config) (*Server, error)
func (s *Server) Start() error
func (s *Server) Shutdown(ctx context.Context) error
func (s *Server) handleConnection(conn net.Conn)
func (s *Server) detectProtocol(conn net.Conn) (Protocol, error)
```

### Protocol Detection Logic
```
1. Read first few bytes (peek)
2. Check for TLS handshake (0x16)
3. Check for HTTP methods (GET, POST, CONNECT)
4. Check for XMPP stream start
5. Default to raw TCP proxy
```

### Connection Handling Flow
```
1. Accept connection
2. Detect protocol
3. Route to handler:
   - HTTP: handleHTTP()
   - HTTPS: handleHTTPS()
   - Jabber: handleJabber()
   - Unknown: handleRaw()
4. Connect upstream via SOCKS5
5. Bidirectional copy
6. Close on completion
```

### Metrics Exposed
- `whatsapp_proxy_connections_total` - Total connections
- `whatsapp_proxy_connections_active` - Active connections
- `whatsapp_proxy_bytes_sent_total` - Bytes sent
- `whatsapp_proxy_bytes_received_total` - Bytes received
- `whatsapp_proxy_errors_total` - Error count
- `whatsapp_proxy_protocol_connections{protocol}` - Per-protocol stats

### Graceful Shutdown
1. Stop accepting new connections
2. Wait for active connections (with timeout)
3. Force close remaining after timeout
4. Clean up resources

## Testing
- HTTP proxy request works
- HTTPS CONNECT tunnel works
- Jabber connection proxies
- Protocol detection accurate
- Metrics endpoint accessible
- Graceful shutdown completes
- Concurrent connections handled

## Deliverables
- internal/proxy package
- Server implementation
- Protocol handlers
- Metrics collector
- Integration tests
- Load testing results
