# Step 5: SSL Certificate Management

## Objective
Implement automatic SSL certificate generation and management for HTTPS termination with caching and rotation support.

## Tasks

### 1. Certificate Generation
- Self-signed certificate creation
- RSA key pair generation (2048-bit)
- Certificate attributes (Subject, SANs)
- Validity period (1 year default)

### 2. Certificate Caching
- File-based cache
- In-memory cache
- Cache invalidation
- Automatic regeneration

### 3. TLS Configuration
- TLS version restrictions (1.2+)
- Cipher suite selection
- Certificate loading
- Dynamic certificate serving

### 4. Custom Certificate Support
- Load from file
- Validation
- Private key protection

### 5. Certificate Rotation
- Expiry detection
- Automatic renewal
- Hot reload without restart

## Implementation Details

### SSL Manager Interface
```go
type Manager struct {
    config    *Config
    certCache map[string]*tls.Certificate
    mutex     sync.RWMutex
}

type Config struct {
    AutoGenerate bool
    CertFile     string
    KeyFile      string
    DNSNames     []string
    IPAddresses  []net.IP
    CacheDir     string
    ValidityDays int
}

func NewManager(cfg *Config) (*Manager, error)
func (m *Manager) GetCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error)
func (m *Manager) GenerateCertificate(dnsNames []string, ips []net.IP) (*tls.Certificate, error)
func (m *Manager) LoadCertificate(certFile, keyFile string) error
func (m *Manager) RotateCertificates() error
```

### Certificate Generation
```go
// RSA 2048-bit key
priv, _ := rsa.GenerateKey(rand.Reader, 2048)

// Certificate template
template := &x509.Certificate{
    SerialNumber: big.NewInt(rand.Int63()),
    Subject: pkix.Name{
        Organization: []string{"WhatsApp Proxy"},
        CommonName:   "WhatsApp Proxy Server",
    },
    NotBefore:             time.Now(),
    NotAfter:              time.Now().Add(365 * 24 * time.Hour),
    KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
    DNSNames:              dnsNames,
    IPAddresses:           ips,
}
```

### TLS Config
```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
    },
    GetCertificate: manager.GetCertificate,
}
```

### Cache Strategy
- Cache key: hash of DNS names + IPs
- File format: PEM (cert + key)
- Cache location: ~/.whatsapp-proxy/certs/
- Check expiry on load
- Regenerate if expired

### Certificate Rotation
- Background goroutine checks expiry daily
- Regenerate 30 days before expiry
- Hot reload: update cache, connections use new cert
- Log rotation events

## Testing
- Generate self-signed certificate
- Load custom certificate
- Cache works correctly
- TLS handshake succeeds
- Rotation triggers correctly
- Expired certificate regenerates

## Deliverables
- internal/ssl package
- Certificate manager
- TLS configuration
- Cache implementation
- Rotation logic
- Unit tests
