package ssl

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Manager manages SSL certificates for the proxy server
type Manager struct {
	config    *Config
	certCache map[string]*tls.Certificate
	mutex     sync.RWMutex
	rotationDone chan struct{}
}

// Config holds SSL manager configuration
type Config struct {
	// AutoGenerate enables automatic certificate generation
	AutoGenerate bool

	// CertFile is the path to custom certificate file
	CertFile string

	// KeyFile is the path to custom private key file
	KeyFile string

	// DNSNames for certificate Subject Alternative Names
	DNSNames []string

	// IPAddresses for certificate Subject Alternative Names
	IPAddresses []net.IP

	// ValidityDays is the certificate validity period in days
	ValidityDays int

	// CacheDir is the directory for caching generated certificates
	CacheDir string

	// EnableRotation enables automatic certificate rotation
	EnableRotation bool

	// RotationCheckInterval is how often to check for expiring certificates
	RotationCheckInterval time.Duration
}

// NewManager creates a new SSL certificate manager
func NewManager(cfg *Config) (*Manager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Set defaults
	if cfg.ValidityDays == 0 {
		cfg.ValidityDays = 365
	}
	if cfg.RotationCheckInterval == 0 {
		cfg.RotationCheckInterval = 24 * time.Hour
	}
	if cfg.CacheDir == "" {
		homeDir, _ := os.UserHomeDir()
		cfg.CacheDir = filepath.Join(homeDir, ".whatsapp-proxy", "certs")
	}

	m := &Manager{
		config:       cfg,
		certCache:    make(map[string]*tls.Certificate),
		rotationDone: make(chan struct{}),
	}

	// Create cache directory if auto-generating
	if cfg.AutoGenerate {
		if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}
	}

	// Load or generate initial certificate
	if err := m.initialize(); err != nil {
		return nil, err
	}

	// Start rotation goroutine if enabled
	if cfg.EnableRotation && cfg.AutoGenerate {
		go m.rotationLoop()
	}

	return m, nil
}

// initialize loads or generates the initial certificate
func (m *Manager) initialize() error {
	if m.config.AutoGenerate {
		// Try to load from cache first
		cacheKey := m.getCacheKey()
		cachedCert, err := loadCertificateFromCache(m.config.CacheDir, cacheKey)
		if err == nil && !isCertificateExpiringSoon(cachedCert, 30) {
			log.Printf("[INFO] Loaded certificate from cache")
			m.certCache["default"] = cachedCert
			return nil
		}

		// Generate new certificate
		log.Printf("[INFO] Generating new self-signed certificate")
		cert, err := generateSelfSignedCertificate(
			m.config.DNSNames,
			m.config.IPAddresses,
			m.config.ValidityDays,
		)
		if err != nil {
			return fmt.Errorf("failed to generate certificate: %w", err)
		}

		// Cache the certificate
		if err := saveCertificateToCache(m.config.CacheDir, cacheKey, cert); err != nil {
			log.Printf("[WARN] Failed to cache certificate: %v", err)
		}

		m.certCache["default"] = cert
		log.Printf("[INFO] Certificate generated successfully")
	} else {
		// Load custom certificate
		log.Printf("[INFO] Loading custom certificate from %s", m.config.CertFile)
		cert, err := loadCertificateFromFiles(m.config.CertFile, m.config.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load certificate: %w", err)
		}

		m.certCache["default"] = cert
		log.Printf("[INFO] Custom certificate loaded successfully")
	}

	return nil
}

// GetCertificate returns a certificate for the given ClientHello
// This implements the tls.Config.GetCertificate callback
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	m.mutex.RLock()
	cert, exists := m.certCache["default"]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no certificate available")
	}

	return cert, nil
}

// GetTLSConfig returns a TLS configuration using this manager
func (m *Manager) GetTLSConfig() *tls.Config {
	return &tls.Config{
		GetCertificate: m.GetCertificate,
		MinVersion:     tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
		PreferServerCipherSuites: true,
	}
}

// RotateCertificates manually rotates certificates
func (m *Manager) RotateCertificates() error {
	if !m.config.AutoGenerate {
		return fmt.Errorf("certificate rotation only available with auto-generation")
	}

	log.Println("[INFO] Rotating certificates...")

	// Generate new certificate
	cert, err := generateSelfSignedCertificate(
		m.config.DNSNames,
		m.config.IPAddresses,
		m.config.ValidityDays,
	)
	if err != nil {
		return fmt.Errorf("failed to generate certificate: %w", err)
	}

	// Update cache
	m.mutex.Lock()
	m.certCache["default"] = cert
	m.mutex.Unlock()

	// Save to disk cache
	cacheKey := m.getCacheKey()
	if err := saveCertificateToCache(m.config.CacheDir, cacheKey, cert); err != nil {
		log.Printf("[WARN] Failed to cache rotated certificate: %v", err)
	}

	log.Println("[INFO] Certificate rotation complete")
	return nil
}

// rotationLoop periodically checks for expiring certificates
func (m *Manager) rotationLoop() {
	ticker := time.NewTicker(m.config.RotationCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mutex.RLock()
			cert, exists := m.certCache["default"]
			m.mutex.RUnlock()

			if exists && isCertificateExpiringSoon(cert, 30) {
				log.Println("[INFO] Certificate expiring soon, rotating...")
				if err := m.RotateCertificates(); err != nil {
					log.Printf("[ERROR] Certificate rotation failed: %v", err)
				}
			}
		case <-m.rotationDone:
			return
		}
	}
}

// Close stops the certificate manager
func (m *Manager) Close() error {
	if m.config.EnableRotation {
		close(m.rotationDone)
	}
	return nil
}

// getCacheKey generates a cache key based on configuration
func (m *Manager) getCacheKey() string {
	// Simple key based on DNS names and IPs
	// In production, you might want a more sophisticated key
	return "default"
}

// GetCertificateInfo returns information about the current certificate
func (m *Manager) GetCertificateInfo() *CertificateInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cert, exists := m.certCache["default"]
	if !exists {
		return nil
	}

	// Parse the certificate
	x509Cert, err := parseTLSCertificate(cert)
	if err != nil {
		return nil
	}

	return &CertificateInfo{
		Subject:    x509Cert.Subject.String(),
		Issuer:     x509Cert.Issuer.String(),
		NotBefore:  x509Cert.NotBefore,
		NotAfter:   x509Cert.NotAfter,
		DNSNames:   x509Cert.DNSNames,
		IPAddresses: x509Cert.IPAddresses,
	}
}

// CertificateInfo holds certificate information
type CertificateInfo struct {
	Subject     string
	Issuer      string
	NotBefore   time.Time
	NotAfter    time.Time
	DNSNames    []string
	IPAddresses []net.IP
}
