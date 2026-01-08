package ssl

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
)

// saveCertificateToCache saves a certificate to the cache directory
func saveCertificateToCache(cacheDir, key string, cert *tls.Certificate) error {
	if cert == nil || len(cert.Certificate) == 0 {
		return fmt.Errorf("invalid certificate")
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	certPath := filepath.Join(cacheDir, key+".crt")
	keyPath := filepath.Join(cacheDir, key+".key")

	// Parse certificate to get PEM encoding
	certPEM, keyPEM, err := encodeCertificateToPEM(cert)
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	// Write certificate file
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	// Write key file (with restricted permissions)
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

// loadCertificateFromCache loads a certificate from the cache directory
func loadCertificateFromCache(cacheDir, key string) (*tls.Certificate, error) {
	certPath := filepath.Join(cacheDir, key+".crt")
	keyPath := filepath.Join(cacheDir, key+".key")

	// Check if files exist
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("certificate file not found")
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file not found")
	}

	// Load certificate
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	return &cert, nil
}

// encodeCertificateToPEM encodes a tls.Certificate to PEM format
func encodeCertificateToPEM(cert *tls.Certificate) ([]byte, []byte, error) {
	if cert == nil {
		return nil, nil, fmt.Errorf("certificate is nil")
	}

	// For generated certificates, we need to get the PEM encoding
	// This is a simplified approach - in production you might want to store
	// the original PEM data during generation

	// Note: This won't work perfectly for certificates without PEM data
	// But for our use case with fresh generation, it's acceptable

	if len(cert.Certificate) == 0 {
		return nil, nil, fmt.Errorf("no certificate data")
	}

	// Since we generate certificates with PEM encoding, we can use a workaround
	// In a real implementation, you'd store PEM data alongside the certificate
	return nil, nil, fmt.Errorf("PEM encoding not available (use fresh generation)")
}
