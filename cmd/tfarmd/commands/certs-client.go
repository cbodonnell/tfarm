package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var certsClientCmd = &cobra.Command{
	Use:           "client [name]",
	Short:         "Generate a client certificate",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}
		return CertsClient(args[0])
	},
}

func init() {
	certsCmd.AddCommand(certsClientCmd)
}

func CertsClient(name string) error {
	workDir := os.Getenv("TFARMD_WORK_DIR")
	if workDir == "" {
		workDir = "/var/lib/tfarmd"
	}

	fmt.Println("Generating client certificate...")

	// read the ca.key and ca.crt files from the workDir/tls directory
	caKeyFile, err := os.Open(path.Join(workDir, "tls", "ca.key"))
	if err != nil {
		return fmt.Errorf("error opening CA key file: %s", err)
	}
	caKeyPEM, err := io.ReadAll(caKeyFile)
	if err != nil {
		return fmt.Errorf("error reading CA key file: %s", err)
	}
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing CA key: %s", err)
	}

	caFile, err := os.Open(path.Join(workDir, "tls", "ca.crt"))
	if err != nil {
		return fmt.Errorf("error opening CA file: %s", err)
	}
	caCertPEM, err := io.ReadAll(caFile)
	if err != nil {
		return fmt.Errorf("error reading CA file: %s", err)
	}
	caCertBlock, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing CA certificate: %s", err)
	}

	// Generate a new client certificate/key pair
	clientTemplate := x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		Issuer:                caCert.Subject,
		BasicConstraintsValid: true,
	}
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating client key: %s", err)
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &clientTemplate, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("error generating client certificate: %s", err)
	}

	// make sure the tls/clients directory exists
	err = os.MkdirAll(path.Join(workDir, "tls", "clients", name), 0755)
	if err != nil {
		return fmt.Errorf("error creating clients directory: %s", err)
	}

	// Write the client certificate and key to files
	certFile, err := os.Create(path.Join(workDir, "tls", "clients", name, "client.crt"))
	if err != nil {
		return fmt.Errorf("error creating client certificate file: %s", err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: clientCert})
	certFile.Close()

	keyFile, err := os.Create(path.Join(workDir, "tls", "clients", name, "client.key"))
	if err != nil {
		return fmt.Errorf("error creating client key file: %s", err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})
	keyFile.Close()

	absPath, err := filepath.Abs(path.Join(workDir, "tls", "clients", name))
	if err != nil {
		return fmt.Errorf("error getting absolute path of client certificate: %s", err)
	}

	fmt.Println("Client certificate and key written to:")
	fmt.Println("  ", path.Join(absPath, "client.crt"))
	fmt.Println("  ", path.Join(absPath, "client.key"))

	return nil
}
