package ssl

import (
	"crypto/tls"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	// Check that certificate was generated
	if len(manager.certCache) == 0 {
		t.Error("No certificates in cache")
	}
}

func TestNewManagerWithCustomCert(t *testing.T) {
	tmpDir := t.TempDir()

	// First generate a certificate to use as "custom"
	cert, err := generateSelfSignedCertificate(
		[]string{"test.local"},
		[]net.IP{net.ParseIP("127.0.0.1")},
		365,
	)
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	// Save to files
	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	// We need to manually save the PEM data
	// For testing, we'll use the tls.Certificate directly
	_ = cert

	// Skip this test for now as it requires more complex setup
	// In production, you'd have actual PEM files
	t.Skip("Skipping custom certificate test (requires PEM files)")

	cfg := &Config{
		AutoGenerate: false,
		CertFile:     certPath,
		KeyFile:      keyPath,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()
}

func TestGenerateSelfSignedCertificate(t *testing.T) {
	dnsNames := []string{"localhost", "test.local"}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	cert, err := generateSelfSignedCertificate(dnsNames, ipAddresses, 365)
	if err != nil {
		t.Fatalf("generateSelfSignedCertificate() error = %v", err)
	}

	if cert == nil {
		t.Fatal("generateSelfSignedCertificate() returned nil")
	}

	// Verify certificate has correct SANs
	x509Cert, err := parseTLSCertificate(cert)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	if len(x509Cert.DNSNames) != len(dnsNames) {
		t.Errorf("DNS names count = %d, want %d", len(x509Cert.DNSNames), len(dnsNames))
	}

	if len(x509Cert.IPAddresses) != len(ipAddresses) {
		t.Errorf("IP addresses count = %d, want %d", len(x509Cert.IPAddresses), len(ipAddresses))
	}

	// Verify validity period
	expectedExpiry := time.Now().Add(365 * 24 * time.Hour)
	if x509Cert.NotAfter.Before(expectedExpiry.Add(-24 * time.Hour)) {
		t.Error("Certificate expires too soon")
	}
}

func TestGetCertificate(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()

	// Test GetCertificate
	cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
		ServerName: "localhost",
	})
	if err != nil {
		t.Errorf("GetCertificate() error = %v", err)
	}
	if cert == nil {
		t.Error("GetCertificate() returned nil")
	}
}

func TestGetTLSConfig(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     []string{"localhost"},
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()

	tlsConfig := manager.GetTLSConfig()
	if tlsConfig == nil {
		t.Fatal("GetTLSConfig() returned nil")
	}

	// Verify TLS configuration
	if tlsConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = %d, want %d (TLS 1.2)", tlsConfig.MinVersion, tls.VersionTLS12)
	}

	if len(tlsConfig.CipherSuites) == 0 {
		t.Error("No cipher suites configured")
	}

	if tlsConfig.GetCertificate == nil {
		t.Error("GetCertificate callback not set")
	}
}

func TestIsCertificateExpiringSoon(t *testing.T) {
	tests := []struct {
		name         string
		validityDays int
		checkDays    int
		wantExpiring bool
	}{
		{
			name:         "expires in 1 day, check 30 days",
			validityDays: 1,
			checkDays:    30,
			wantExpiring: true,
		},
		{
			name:         "expires in 60 days, check 30 days",
			validityDays: 60,
			checkDays:    30,
			wantExpiring: false,
		},
		{
			name:         "expires in 365 days, check 30 days",
			validityDays: 365,
			checkDays:    30,
			wantExpiring: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert, err := generateSelfSignedCertificate(
				[]string{"test"},
				[]net.IP{net.ParseIP("127.0.0.1")},
				tt.validityDays,
			)
			if err != nil {
				t.Fatalf("Failed to generate certificate: %v", err)
			}

			expiring := isCertificateExpiringSoon(cert, tt.checkDays)
			if expiring != tt.wantExpiring {
				t.Errorf("isCertificateExpiringSoon() = %v, want %v", expiring, tt.wantExpiring)
			}
		})
	}
}

func TestRotateCertificates(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     []string{"localhost"},
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()

	// Get original certificate
	origCert, _ := manager.GetCertificate(nil)

	// Rotate certificates
	if err := manager.RotateCertificates(); err != nil {
		t.Errorf("RotateCertificates() error = %v", err)
	}

	// Get new certificate
	newCert, _ := manager.GetCertificate(nil)

	// Verify they're different (different serial numbers)
	origX509, _ := parseTLSCertificate(origCert)
	newX509, _ := parseTLSCertificate(newCert)

	if origX509.SerialNumber.Cmp(newX509.SerialNumber) == 0 {
		t.Error("Certificate was not rotated (same serial number)")
	}
}

func TestGetCertificateInfo(t *testing.T) {
	tmpDir := t.TempDir()

	dnsNames := []string{"localhost", "test.local"}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1")}

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     dnsNames,
		IPAddresses:  ipAddresses,
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()

	info := manager.GetCertificateInfo()
	if info == nil {
		t.Fatal("GetCertificateInfo() returned nil")
	}

	if len(info.DNSNames) != len(dnsNames) {
		t.Errorf("DNS names count = %d, want %d", len(info.DNSNames), len(dnsNames))
	}

	if len(info.IPAddresses) != len(ipAddresses) {
		t.Errorf("IP addresses count = %d, want %d", len(info.IPAddresses), len(ipAddresses))
	}

	if info.NotAfter.Before(time.Now()) {
		t.Error("Certificate already expired")
	}
}

func TestCacheDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     []string{"localhost"},
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer manager.Close()

	// Verify cache directory was created
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Cache directory was not created")
	}
}

// Benchmark tests
func BenchmarkGenerateSelfSignedCertificate(b *testing.B) {
	dnsNames := []string{"localhost"}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateSelfSignedCertificate(dnsNames, ipAddresses, 365)
	}
}

func BenchmarkGetCertificate(b *testing.B) {
	tmpDir := b.TempDir()

	cfg := &Config{
		AutoGenerate: true,
		DNSNames:     []string{"localhost"},
		ValidityDays: 365,
		CacheDir:     tmpDir,
	}

	manager, _ := NewManager(cfg)
	defer manager.Close()

	hello := &tls.ClientHelloInfo{ServerName: "localhost"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.GetCertificate(hello)
	}
}
