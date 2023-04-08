package info

type Info struct {
	Ready    bool     `json:"ready"`
	Error    string   `json:"error,omitempty"`
	Version  string   `json:"version"`
	TokenDir string   `json:"token_dir"`
	Endpoint string   `json:"endpoint"`
	OIDC     OIDCInfo `json:"oidc"`
}

type OIDCInfo struct {
	Issuer   string `json:"issuer"`
	ClientID string `json:"client_id"`
}
