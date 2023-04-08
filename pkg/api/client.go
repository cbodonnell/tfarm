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
	"github.com/cbodonnell/tfarm/pkg/certs"
	"github.com/cbodonnell/tfarm/pkg/version"
)

type APIClient struct {
	endpoint   string
	httpClient *http.Client
	configDir  string
}

// TODO: move this to the info package and differentiate between tfarm and ranch info
type Info struct {
	Client ClientInfo `json:"client"`
	Server ServerInfo `json:"server"`
}

type ClientInfo struct {
	Version string `json:"version"`
	Config  string `json:"config"`
}

type ServerInfo struct {
	Version  string `json:"version"`
	Endpoint string `json:"endpoint"`
	Error    string `json:"error,omitempty"`
}

type ServerInfoResponse struct {
	Version string `json:"version"`
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
	client, err := certs.LoadClientFromFile(path.Join(configDir, "client.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to load client: %s", err)
	}

	cert, err := tls.X509KeyPair(client.Cert, client.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(client.CA)

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

func (c *APIClient) Info() *Info {
	info := &Info{
		Client: ClientInfo{
			Version: version.Version,
			Config:  c.configDir,
		},
		Server: ServerInfo{
			Endpoint: c.endpoint,
		},
	}

	if serverInfo, err := c.getServerInfo(); err != nil {
		info.Server.Error = fmt.Sprintf("error getting info: %s", err)
	} else {
		info.Server.Version = serverInfo.Version
	}

	return info
}

func (c *APIClient) getServerInfo() (*ServerInfoResponse, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/info")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response with status code %d: %s", resp.StatusCode, err)
	}

	serverInfo := &ServerInfoResponse{}
	if err := json.Unmarshal([]byte(response.Message), serverInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response message: %s", err)
	}

	return serverInfo, nil
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
