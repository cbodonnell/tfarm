package certs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
