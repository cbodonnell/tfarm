package auth

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/cbodonnell/tfarm/pkg/crypto"
	"github.com/cbodonnell/tfarm/pkg/frpc"
)

type ClientCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ClientSecretWorker struct {
	clientEndpoint       string
	workDir              string                // directory to store tokens.json
	loginCredentialsChan chan LoginCredentials // channel for receiving login credentials
	loginResultChan      chan LoginResult      // channel for sending login result
	isAuthenticated      bool                  // flag for whether or not the worker is logged in
	needsLogin           bool
}

func NewClientSecretWorker(workDir string) *ClientSecretWorker {
	return &ClientSecretWorker{
		workDir:              workDir,
		loginCredentialsChan: make(chan LoginCredentials),
		loginResultChan:      make(chan LoginResult),
		isAuthenticated:      false,
		needsLogin:           false,
	}
}

func (o *ClientSecretWorker) NeedsLogin() bool {
	return o.needsLogin
}

func (o *ClientSecretWorker) IsAuthenticated() bool {
	return o.isAuthenticated
}

func (o *ClientSecretWorker) LoginCredentialsChan() chan LoginCredentials {
	return o.loginCredentialsChan
}

func (o *ClientSecretWorker) LoginResultChan() chan LoginResult {
	return o.loginResultChan
}

func (o *ClientSecretWorker) WaitForLogin() {
	for {
		log.Println("checking for credentials...")
		credPath := path.Join(o.workDir, "credentials.json")
		credBytes, err := os.ReadFile(credPath)
		if err != nil {
			if os.IsNotExist(err) {
				o.needsLogin = true
				log.Println("credentials not found, waiting for login")
				loginCreds := <-o.loginCredentialsChan
				o.needsLogin = false

				creds := &ClientCredentials{
					ClientID:     loginCreds.Username,
					ClientSecret: loginCreds.Password,
				}

				log.Println("received login, saving credentials...")
				credBytes, err := json.Marshal(creds)
				if err != nil {
					log.Printf("error marshaling credentials: %s", err)
					o.loginResultChan <- LoginResult{Error: err}
					continue
				}

				if err := os.WriteFile(credPath, credBytes, 0644); err != nil {
					log.Printf("error writing credentials.json: %s", err)
					o.loginResultChan <- LoginResult{Error: err}
					continue
				}

				log.Println("login successful")
				o.loginResultChan <- LoginResult{Success: true}
				continue
			}

			log.Printf("error reading credentials.json: %s", err)
			continue
		}

		log.Println("unmarshaling credentials...")
		creds := &ClientCredentials{}
		if err := json.Unmarshal(credBytes, creds); err != nil {
			log.Printf("error unmarshaling credentials: %s", err)
			if err := os.Remove(credPath); err != nil {
				log.Printf("error deleting credentials.json: %s", err)
			}
			continue
		}

		decodedSecret, err := base64.StdEncoding.DecodeString(creds.ClientSecret)
		if err != nil {
			log.Printf("error decoding client secret: %s", err)
			continue
		}

		log.Println("writing credentials to frpc config...")
		cfg, err := frpc.ParseFrpcCommonConfig(path.Join(o.workDir, "frpc.ini"))
		if err != nil {
			log.Printf("error parsing frpc.ini: %s", err)
			continue
		}

		if cfg.Metas == nil {
			cfg.Metas = make(map[string]string)
		}

		cfg.Metas["client_id"] = creds.ClientID
		cfg.Metas["client_signature"] = crypto.HMAC(decodedSecret, []byte(creds.ClientID))

		if err := frpc.SaveFrpcCommonConfig(cfg, path.Join(o.workDir, "frpc.ini")); err != nil {
			log.Printf("error writing frpc.ini: %s", err)
			continue
		}

		log.Println("authentication successful")
		o.isAuthenticated = true
		return // we're authenticated
	}
}
