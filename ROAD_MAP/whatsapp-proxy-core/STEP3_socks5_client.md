# Step 3: SOCKS5 Client Implementation

## Objective
Implement robust SOCKS5 client for upstream proxy connections with authentication, connection pooling, and error handling.

## Tasks

### 1. SOCKS5 Protocol Implementation
- SOCKS5 handshake (RFC 1928)
- Authentication methods (no auth, username/password)
- Connection establishment
- Command handling (CONNECT)

### 2. Connection Management
- Dial through SOCKS5 proxy
- Connection wrapper
- Timeout handling
- Keep-alive settings

### 3. Connection Pooling
- Pool implementation
- Connection reuse
- Idle connection cleanup
- Pool size limits

### 4. Error Handling
- SOCKS5 error codes
- Retry logic with backoff
- Connection failure fallback
- Detailed error messages

### 5. Testing and Validation
- Connection test utility
- Proxy reachability check
- Authentication validation

## Implementation Details

### SOCKS5 Client Interface
```go
type Client struct {
    config     *Config
    dialer     *net.Dialer
    pool       *connpool.Pool
}

type Config struct {
    ProxyAddr  string
    Username   string
    Password   string
    Timeout    time.Duration
    PoolSize   int
}

func NewClient(config *Config) (*Client, error)
func (c *Client) Dial(network, addr string) (net.Conn, error)
func (c *Client) DialContext(ctx context.Context, network, addr string) (net.Conn, error)
func (c *Client) Close() error
func (c *Client) Test() error
```

### SOCKS5 Handshake Flow
1. Client greeting (methods)
2. Server method selection
3. Authentication (if required)
4. Connection request
5. Server reply
6. Data transfer

### Authentication Methods
- 0x00: No authentication
- 0x02: Username/password (RFC 1929)

### Connection Pool Strategy
- Max idle connections: 10
- Max open connections: 100
- Idle timeout: 90 seconds
- Connection lifetime: 10 minutes

### Error Handling Strategy
- Network errors: retry 3 times with exponential backoff
- Auth errors: fail immediately
- SOCKS5 errors: log and return detailed message

## Testing
- SOCKS5 handshake without auth
- SOCKS5 handshake with username/password
- Connection through proxy succeeds
- Pool reuses connections
- Error handling works correctly
- Timeout handling

## Deliverables
- internal/socks5 package
- Client implementation
- Connection pool
- Comprehensive unit tests
- Integration tests with mock SOCKS5 server
