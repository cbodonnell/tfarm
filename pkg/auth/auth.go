package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var isConfigured bool
var isConfiguredMu sync.RWMutex

type ConfigureCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	// TODO: add client certificate, key, and ca to this and save to ./tls/frps
}

func IsConfigured() bool {
	isConfiguredMu.RLock()
	defer isConfiguredMu.RUnlock()
	return isConfigured
}

func WaitForCredentials(workDir string) (*ConfigureCredentials, error) {
	isConfiguredMu.Lock()
	defer isConfiguredMu.Unlock()
	credsPath := path.Join(workDir, "credentials.json")
	if _, err := os.Stat(credsPath); err != nil {
		// check if the file is not found
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking for credentials.json: %s", err)
		}

		isConfigured = false
		log.Println("waiting for credentials.json to be created")
		for {
			if _, err := os.Stat(credsPath); err != nil {
				if !os.IsNotExist(err) {
					return nil, fmt.Errorf("error checking for credentials.json: %s", err)
				}
			} else {
				break
			}
			time.Sleep(1 * time.Second)
		}
		log.Println("credentials.json created")
	}
	isConfigured = true

	b, err := os.ReadFile(credsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading credentials.json: %s", err)
	}

	creds := &ConfigureCredentials{}
	if err := json.Unmarshal(b, creds); err != nil {
		return nil, fmt.Errorf("error unmarshaling credentials.json: %s", err)
	}

	return creds, nil
}
