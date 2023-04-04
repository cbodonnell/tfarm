package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	DefaultEndpoint = "http://localhost:9090"
)

type APIClient struct {
	endpoint   string
	httpClient *http.Client
	token      *oauth2.Token
}

type APIRequestParams struct {
	QueryParams map[string]string
}

type ClientRequestParams struct {
	APIRequestParams
	ID string
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewClient(httpClient *http.Client, endpoint string) *APIClient {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	return &APIClient{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

func (c *APIClient) SetToken(token *oauth2.Token) {
	c.token = token
}

func (c *APIClient) ListClients(params *APIRequestParams) ([]*ClientResponse, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/clients")
	if err != nil {
		return nil, fmt.Errorf("error listing clients: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var response []*ClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return response, nil
}

func (c *APIClient) GetClient(params *ClientRequestParams) (*ClientResponse, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/clients/" + params.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var response ClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &response, nil
}

func (c *APIClient) GetClientCredentialsJson(params *ClientRequestParams) ([]byte, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/clients/" + params.ID + "/credentials.json")
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return b, nil
}

func (c *APIClient) CreateClient(params *APIRequestParams) (*ClientResponse, error) {
	resp, err := c.httpClient.Post(c.endpoint+"/api/clients", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	response := &ClientResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

func (c *APIClient) CreateClientJSON(params *APIRequestParams) ([]byte, error) {
	u, err := url.Parse(c.endpoint + "/api/clients")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %w", err)
	}

	q := u.Query()
	for k, v := range params.QueryParams {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Post(u.String(), "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return b, nil
}

func (c *APIClient) DeleteClient(params *ClientRequestParams) (*ClientResponse, error) {
	req, err := http.NewRequest("DELETE", c.endpoint+"/api/clients/"+params.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error deleting client: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var response ClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &response, nil
}
