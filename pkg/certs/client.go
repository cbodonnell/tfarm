package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Client struct {
	CA   []byte `json:"ca"`
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

type ClientFile struct {
	CA   string `json:"ca"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func LoadClientFromFile(path string) (*Client, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening client file: %s", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading client file: %s", err)
	}
	clientFile := &ClientFile{}
	if err := json.Unmarshal(b, clientFile); err != nil {
		return nil, fmt.Errorf("error unmarshaling client file: %s", err)
	}

	c := &Client{}
	c.CA, err = base64.URLEncoding.DecodeString(clientFile.CA)
	if err != nil {
		return nil, fmt.Errorf("error decoding CA: %s", err)
	}
	c.Cert, err = base64.URLEncoding.DecodeString(clientFile.Cert)
	if err != nil {
		return nil, fmt.Errorf("error decoding cert: %s", err)
	}
	c.Key, err = base64.URLEncoding.DecodeString(clientFile.Key)
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %s", err)
	}

	return c, nil
}

func (c *Client) SaveToFile(path string) error {
	clientFile := &ClientFile{
		CA:   base64.URLEncoding.EncodeToString(c.CA),
		Cert: base64.URLEncoding.EncodeToString(c.Cert),
		Key:  base64.URLEncoding.EncodeToString(c.Key),
	}
	b, err := json.Marshal(clientFile)
	if err != nil {
		return fmt.Errorf("error marshaling client file: %s", err)
	}
	if err := os.WriteFile(path, b, 0600); err != nil {
		return fmt.Errorf("error writing client file: %s", err)
	}
	return nil
}

func GenerateClientCerts(dir string, name string) error {
	fmt.Println("Generating client certificate...")

	// read the ca.key and ca.crt files from the workDir/tls directory
	caKeyFile, err := os.Open(path.Join(dir, "ca.key"))
	if err != nil {
		return fmt.Errorf("error opening CA key file: %s", err)
	}
	defer caKeyFile.Close()
	caKeyPEM, err := io.ReadAll(caKeyFile)
	if err != nil {
		return fmt.Errorf("error reading CA key file: %s", err)
	}
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing CA key: %s", err)
	}

	caFile, err := os.Open(path.Join(dir, "ca.crt"))
	if err != nil {
		return fmt.Errorf("error opening CA file: %s", err)
	}
	defer caFile.Close()
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
	err = os.MkdirAll(path.Join(dir, "clients", name), 0755)
	if err != nil {
		return fmt.Errorf("error creating clients directory: %s", err)
	}

	certBuff := bytes.NewBuffer(nil)
	if err := pem.Encode(certBuff, &pem.Block{Type: "CERTIFICATE", Bytes: clientCert}); err != nil {
		return fmt.Errorf("error pem encoding client certificate: %s", err)
	}

	keyBuff := bytes.NewBuffer(nil)
	if err := pem.Encode(keyBuff, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)}); err != nil {
		return fmt.Errorf("error pem encoding client key: %s", err)
	}

	client := &Client{
		CA:   caCertPEM,
		Cert: certBuff.Bytes(),
		Key:  keyBuff.Bytes(),
	}
	if err := client.SaveToFile(path.Join(dir, "clients", name, "client.json")); err != nil {
		return fmt.Errorf("error saving client to file: %s", err)
	}

	absPath, err := filepath.Abs(path.Join(dir, "clients", name))
	if err != nil {
		return fmt.Errorf("error getting absolute path of client certificate: %s", err)
	}

	fmt.Println("Client certificate and key written to:")
	fmt.Println("  ", path.Join(absPath, "client.json"))

	return nil
}
