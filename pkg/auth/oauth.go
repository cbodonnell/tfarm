package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"golang.org/x/oauth2"
)

type OAuthWorker struct {
	clientEndpoint       string
	workDir              string                // directory to store tokens.json
	loginCredentialsChan chan LoginCredentials // channel for receiving login credentials
	loginResultChan      chan LoginResult      // channel for sending login result
	isAuthenticated      bool                  // flag for whether or not the worker is logged in
	needsLogin           bool
}

func NewOAuthWorker(clientEndpoint, workDir string) *OAuthWorker {
	// create an oauth2 config without a client ID or secret
	return &OAuthWorker{
		clientEndpoint:       clientEndpoint,
		workDir:              workDir,
		loginCredentialsChan: make(chan LoginCredentials),
		loginResultChan:      make(chan LoginResult),
		isAuthenticated:      false,
		needsLogin:           false,
	}
}

func (o *OAuthWorker) NeedsLogin() bool {
	return o.needsLogin
}

func (o *OAuthWorker) IsAuthenticated() bool {
	return o.isAuthenticated
}

func (o *OAuthWorker) LoginCredentialsChan() chan LoginCredentials {
	return o.loginCredentialsChan
}

func (o *OAuthWorker) LoginResultChan() chan LoginResult {
	return o.loginResultChan
}

func (o *OAuthWorker) Start() {
	go func() {
		for {
			log.Println("checking for token...")
			tokenPath := path.Join(o.workDir, "tokens.json")
			tokenBytes, err := os.ReadFile(tokenPath)
			if err != nil {
				if os.IsNotExist(err) {
					o.needsLogin = true
					log.Println("token not found, waiting for login")
					credentials := <-o.loginCredentialsChan
					o.needsLogin = false
					log.Println("received login, getting token...")

					token, err := o.GetOAuthToken(credentials.Username, credentials.Password)
					if err != nil {
						log.Printf("error getting token: %s", err)
						// TODO: check if the error is an invalid credentials error
						// if so, send the error back to the channel
						// otherwise, just send a generic error
						o.loginResultChan <- LoginResult{Error: err}
						continue
					}

					log.Println("saving token...")
					tokenBytes, err := json.Marshal(token)
					if err != nil {
						log.Printf("error marshaling token: %s", err)
						o.loginResultChan <- LoginResult{Error: err}
						continue
					}

					if err := os.WriteFile(tokenPath, tokenBytes, 0644); err != nil {
						log.Printf("error writing tokens.json: %s", err)
						o.loginResultChan <- LoginResult{Error: err}
						continue
					}

					log.Println("login successful")
					o.loginResultChan <- LoginResult{Success: true}
					continue
				}

				log.Printf("error reading tokens.json: %s", err)
				continue
			}

			log.Println("unmarshaling token...")
			token := &oauth2.Token{}
			if err := json.Unmarshal(tokenBytes, token); err != nil {
				log.Printf("error unmarshaling token: %s", err)
				if err := os.Remove(tokenPath); err != nil {
					log.Printf("error deleting tokens.json: %s", err)
				}
				continue
			}

			log.Println("checking if token is valid...")
			if !token.Valid() {
				// TODO: refresh the token, but for now delete token and wait for a login
				log.Println("token is invalid, deleting token and waiting for login")
				if err := os.Remove(tokenPath); err != nil {
					log.Printf("error deleting tokens.json: %s", err)
				}
				continue
			}

			o.isAuthenticated = true
			log.Println("token is valid, sleeping until token expires...")
			time.Sleep(token.Expiry.Sub(time.Now().Add(30 * time.Second)))
		}
	}()
}

func (o *OAuthWorker) WaitForLogin() {
	for !o.isAuthenticated {
		time.Sleep(time.Second)
	}
}

func (o *OAuthWorker) GetOAuthToken(username, password string) (*oauth2.Token, error) {
	req, err := http.NewRequest("POST", o.clientEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		// read the response body to avoid leaking connections
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response error: %s", err)
		}
		return nil, &ErrInvalidCredentials{Err: fmt.Errorf("invalid credentials: %s", string(body))}
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %s", err)
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %s", err)
	}

	return token, nil
}
