package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/cbodonnell/tfarm/pkg/auth"
	"github.com/cbodonnell/tfarm/pkg/version"
)

type APIClient struct {
	endpoint   string
	httpClient *http.Client
	configDir  string
}

type APIInfo struct {
	ClientVersion  string `json:"client_version"`
	ServerVersion  string `json:"server_version"`
	ServerEndpoint string `json:"server_endpoint"`
	ConfigDir      string `json:"config_dir"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type APIRequest struct {
	// nothing
}

type CreateRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	LocalIP   string `json:"local_ip"`
	LocalPort int    `json:"local_port"`
	ProxyID   string // client-side identifier
}

type DeleteRequest struct {
	Name string `json:"name"`
}

func NewClient(endpoint string, configDir string) (*APIClient, error) {
	tlsFiles := &TLSFiles{
		CertFile: path.Join(configDir, "tls", "client.crt"),
		KeyFile:  path.Join(configDir, "tls", "client.key"),
		CAFile:   path.Join(configDir, "tls", "ca.crt"),
	}

	// Load the server certificate and key
	cert, err := tls.LoadX509KeyPair(tlsFiles.CertFile, tlsFiles.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %s", err)
	}

	// Load the CA certificate
	caCert, err := ioutil.ReadFile(tlsFiles.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration with the server certificate/key and CA certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// Create an HTTP client with the TLS configuration
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &APIClient{
		endpoint:   endpoint,
		httpClient: httpClient,
		configDir:  configDir,
	}, nil
}

func (c *APIClient) Info() (*APIInfo, error) {
	return &APIInfo{
		ClientVersion:  version.Version,
		ServerVersion:  "TODO",
		ServerEndpoint: c.endpoint,
		ConfigDir:      c.configDir,
	}, nil
}

//  TODO: Refactor this to use a generic Do method
func (c *APIClient) Status(req *APIRequest) (*APIResponse, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}

func (c *APIClient) Verify(req *APIRequest) (*APIResponse, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/verify")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}

func (c *APIClient) Reload(req *APIRequest) (*APIResponse, error) {
	resp, err := c.httpClient.Post(c.endpoint+"/api/reload", "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}

func (c *APIClient) Restart(req *APIRequest) (*APIResponse, error) {
	resp, err := c.httpClient.Post(c.endpoint+"/api/restart", "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}

func (c *APIClient) Configure(credentials *auth.ConfigureCredentials) (*APIResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(credentials); err != nil {
		return nil, err
	}
	body := bytes.NewReader(buf.Bytes())

	req, err := http.NewRequest("PUT", c.endpoint+"/api/configure", body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}

func (c *APIClient) Create(req *CreateRequest) (*APIResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return nil, err
	}
	body := bytes.NewReader(buf.Bytes())

	resp, err := c.httpClient.Post(c.endpoint+"/api/tunnel", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}

func (c *APIClient) Delete(opts *DeleteRequest) (*APIResponse, error) {
	req, err := http.NewRequest("DELETE", c.endpoint+fmt.Sprintf("/api/tunnel/%s", opts.Name), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	return &response, nil
}
