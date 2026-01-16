package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
)

const (
	CAValidityDays   = 3650 // 10 years
	CertValidityDays = 365  // 1 year
	KeySize          = 2048
)

// GenerateCA creates a new Certificate Authority
func GenerateCA() error {
	caDir := config.GetCADir()
	if err := os.MkdirAll(caDir, 0700); err != nil {
		return fmt.Errorf("failed to create CA directory: %w", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create CA certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"InstantTLS Local CA"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    "InstantTLS Local Development CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, CAValidityDays),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		MaxPathLen:            1,
	}

	// Self-sign the CA certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	// Save CA certificate
	certPath := filepath.Join(caDir, "ca.crt")
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write CA certificate: %w", err)
	}

	// Save CA private key
	keyPath := filepath.Join(caDir, "ca.key")
	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create CA key file: %w", err)
	}
	defer keyFile.Close()

	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}); err != nil {
		return fmt.Errorf("failed to write CA key: %w", err)
	}

	return nil
}

// CAExists checks if the CA has been generated
func CAExists() bool {
	caDir := config.GetCADir()
	certPath := filepath.Join(caDir, "ca.crt")
	keyPath := filepath.Join(caDir, "ca.key")

	_, certErr := os.Stat(certPath)
	_, keyErr := os.Stat(keyPath)

	return certErr == nil && keyErr == nil
}

// LoadCA loads the CA certificate and key
func LoadCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	caDir := config.GetCADir()

	// Load CA certificate
	certPEM, err := os.ReadFile(filepath.Join(caDir, "ca.crt"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, nil, fmt.Errorf("failed to decode CA certificate PEM")
	}

	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Load CA private key
	keyPEM, err := os.ReadFile(filepath.Join(caDir, "ca.key"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CA key: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("failed to decode CA key PEM")
	}

	caKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA key: %w", err)
	}

	return caCert, caKey, nil
}

// GenerateCert creates a certificate for the given domain
func GenerateCert(domain string) (string, error) {
	if !CAExists() {
		return "", fmt.Errorf("CA not found. Run 'instanttls init' first")
	}

	caCert, caKey, err := LoadCA()
	if err != nil {
		return "", err
	}

	// Sanitize domain for directory name
	sanitizedDomain := sanitizeDomain(domain)
	certDir := filepath.Join(config.GetCertsDir(), sanitizedDomain)
	if err := os.MkdirAll(certDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create cert directory: %w", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return "", fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"InstantTLS"},
			CommonName:   domain,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, CertValidityDays),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Handle wildcard and regular domains
	if strings.HasPrefix(domain, "*.") {
		baseDomain := strings.TrimPrefix(domain, "*.")
		template.DNSNames = []string{domain, baseDomain}
	} else {
		template.DNSNames = []string{domain}
		// Check if it's an IP address
		if ip := net.ParseIP(domain); ip != nil {
			template.IPAddresses = []net.IP{ip}
			template.DNSNames = nil
		}
	}

	// Sign the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &privateKey.PublicKey, caKey)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save certificate
	certPath := filepath.Join(certDir, "cert.pem")
	certFile, err := os.Create(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return "", fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key
	keyPath := filepath.Join(certDir, "key.pem")
	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}); err != nil {
		return "", fmt.Errorf("failed to write key: %w", err)
	}

	return certDir, nil
}

// ListCerts returns all generated certificates
func ListCerts() ([]CertInfo, error) {
	certsDir := config.GetCertsDir()
	var certs []CertInfo

	entries, err := os.ReadDir(certsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return certs, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		certPath := filepath.Join(certsDir, entry.Name(), "cert.pem")
		certPEM, err := os.ReadFile(certPath)
		if err != nil {
			continue
		}

		certBlock, _ := pem.Decode(certPEM)
		if certBlock == nil {
			continue
		}

		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			continue
		}

		certs = append(certs, CertInfo{
			Domain:    entry.Name(),
			NotBefore: cert.NotBefore,
			NotAfter:  cert.NotAfter,
			Path:      filepath.Join(certsDir, entry.Name()),
		})
	}

	return certs, nil
}

// CountWildcardCerts returns the number of wildcard certificates
func CountWildcardCerts() int {
	certs, err := ListCerts()
	if err != nil {
		return 0
	}

	count := 0
	for _, cert := range certs {
		if strings.HasPrefix(cert.Domain, "_") { // sanitized wildcard
			count++
		}
	}
	return count
}

// RenewExpiring renews certificates expiring within the given days
func RenewExpiring(daysThreshold int) ([]string, error) {
	certs, err := ListCerts()
	if err != nil {
		return nil, err
	}

	var renewed []string
	threshold := time.Now().AddDate(0, 0, daysThreshold)

	for _, cert := range certs {
		if cert.NotAfter.Before(threshold) {
			// Re-generate the certificate
			domain := unsanitizeDomain(cert.Domain)
			if _, err := GenerateCert(domain); err != nil {
				return renewed, fmt.Errorf("failed to renew %s: %w", domain, err)
			}
			renewed = append(renewed, domain)
		}
	}

	return renewed, nil
}

type CertInfo struct {
	Domain    string
	NotBefore time.Time
	NotAfter  time.Time
	Path      string
}

func sanitizeDomain(domain string) string {
	// Replace * with _ for filesystem
	sanitized := strings.ReplaceAll(domain, "*", "_")
	// Replace other problematic characters
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return reg.ReplaceAllString(sanitized, "_")
}

func unsanitizeDomain(sanitized string) string {
	// Replace leading _ back to *
	if strings.HasPrefix(sanitized, "_") {
		return "*" + strings.TrimPrefix(sanitized, "_")
	}
	return sanitized
}
