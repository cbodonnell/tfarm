package commands

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
	"path"
	"time"

	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:   "certs [type]",
	Short: "Generate TLS certificates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Certs()
	},
}

func init() {
	rootCmd.AddCommand(certsCmd)
}

func Certs() error {
	workDir := os.Getenv("TFARMD_WORK_DIR")
	if workDir == "" {
		workDir = "/var/lib/tfarmd"
	}

	// Generate a new CA certificate
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "tfarmd CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating CA key: %s", err)
	}
	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("error generating CA certificate: %s", err)
	}

	// TODO: make this configurable through a config file (like /etc/tfarm/config)
	// Write the CA certificate to a file
	caFile, err := os.Create(path.Join(workDir, "tls", "ca.crt"))
	if err != nil {
		return fmt.Errorf("error creating CA file: %s", err)
	}
	pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCert})
	caFile.Close()

	// Generate a new server certificate/key pair
	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		Issuer:                caTemplate.Subject,
		BasicConstraintsValid: true,
	}
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating server key: %s", err)
	}
	serverCert, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &caTemplate, &serverKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("error generating server certificate: %s", err)
	}

	// Write the server certificate and key to files
	certFile, err := os.Create(path.Join(workDir, "tls", "server.crt"))
	if err != nil {
		return fmt.Errorf("error creating server certificate file: %s", err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: serverCert})
	certFile.Close()

	keyFile, err := os.Create(path.Join(workDir, "tls", "server.key"))
	if err != nil {
		return fmt.Errorf("error creating server key file: %s", err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})
	keyFile.Close()

	// Generate a new client certificate/key pair
	clientTemplate := x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName: "tfarmd client",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		Issuer:                caTemplate.Subject,
		BasicConstraintsValid: true,
	}
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating client key: %s", err)
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &clientTemplate, &caTemplate, &clientKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("error generating client certificate: %s", err)
	}

	// Write the client certificate and key to files
	certFile, err = os.Create(path.Join(workDir, "tls", "client.crt"))
	if err != nil {
		return fmt.Errorf("error creating client certificate file: %s", err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: clientCert})
	certFile.Close()

	keyFile, err = os.Create(path.Join(workDir, "tls", "client.key"))
	if err != nil {
		return fmt.Errorf("error creating client key file: %s", err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})
	keyFile.Close()

	return nil
}
