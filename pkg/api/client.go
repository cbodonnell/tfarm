package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIClient struct {
	endpoint   string
	httpClient *http.Client
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type APIRequest struct {
	// nothing
}

type LoginRequest struct {
	AdminPort int    `json:"admin_port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
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

func NewClient(endpoint string) *APIClient {
	return &APIClient{
		endpoint:   endpoint,
		httpClient: http.DefaultClient,
	}
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
		return nil, err
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
		return nil, err
	}
	return &response, nil
}

func (c *APIClient) Reload(req *APIRequest) (*APIResponse, error) {
	resp, err := c.httpClient.Post(c.endpoint+"/api/reload", "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
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
		return nil, err
	}
	return &response, nil
}

func (c *APIClient) Login(opts *LoginRequest) (*APIResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(opts); err != nil {
		return nil, err
	}
	body := bytes.NewReader(buf.Bytes())

	req, err := http.NewRequest("PUT", c.endpoint+"/api/login", body)
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	return &response, nil
}
