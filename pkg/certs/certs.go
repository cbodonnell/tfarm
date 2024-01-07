package certs

import (
	"bytes"
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
	"path/filepath"
	"time"
)

// TODO: Split this into a CA, server, and client cert generator
func GenerateServerCerts(dir string) error {
	fmt.Println("Generating CA certificate...")

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

	// make sure the tls directory exists
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating clients directory: %s", err)
	}

	// Write the CA key to a file
	caKeyFile, err := os.OpenFile(path.Join(dir, "ca.key"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error creating CA key file: %s", err)
	}
	defer caKeyFile.Close()
	if err := pem.Encode(caKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caKey)}); err != nil {
		return fmt.Errorf("error writing CA key file: %s", err)
	}

	// Write the CA certificate to a file
	caFile, err := os.OpenFile(path.Join(dir, "ca.crt"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating CA file: %s", err)
	}
	defer caFile.Close()
	if err := pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCert}); err != nil {
		return fmt.Errorf("error writing CA file: %s", err)
	}

	fmt.Println("Generating server certificate...")

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
	certFile, err := os.OpenFile(path.Join(dir, "server.crt"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating server certificate file: %s", err)
	}
	defer certFile.Close()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: serverCert}); err != nil {
		return fmt.Errorf("error writing server certificate file: %s", err)
	}

	keyFile, err := os.OpenFile(path.Join(dir, "server.key"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error creating server key file: %s", err)
	}
	defer keyFile.Close()
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)}); err != nil {
		return fmt.Errorf("error writing server key file: %s", err)
	}

	fmt.Println("Generating admin client certificate...")

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

	caBuff := bytes.NewBuffer(nil)
	if err := pem.Encode(caBuff, &pem.Block{Type: "CERTIFICATE", Bytes: caCert}); err != nil {
		return fmt.Errorf("error encoding CA certificate: %s", err)
	}

	certBuff := bytes.NewBuffer(nil)
	if err := pem.Encode(certBuff, &pem.Block{Type: "CERTIFICATE", Bytes: clientCert}); err != nil {
		return fmt.Errorf("error encoding client certificate: %s", err)
	}

	keyBuff := bytes.NewBuffer(nil)
	if err := pem.Encode(keyBuff, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)}); err != nil {
		return fmt.Errorf("error encoding client key: %s", err)
	}

	client := &Client{
		CA:   caBuff.Bytes(),
		Cert: certBuff.Bytes(),
		Key:  keyBuff.Bytes(),
	}
	if err := client.SaveToFile(path.Join(dir, "client.json")); err != nil {
		return fmt.Errorf("error saving client to file: %s", err)
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("error getting absolute path of tls directory: %s", err)
	}

	fmt.Println("Certificates written to:")
	fmt.Println("  ", absPath)

	return nil
}
