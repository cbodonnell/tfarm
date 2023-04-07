package info

type Info struct {
	Ready    bool     `json:"ready"`
	Version  string   `json:"version"`
	TokenDir string   `json:"token_dir"`
	Endpoint string   `json:"endpoint"`
	OIDC     OIDCInfo `json:"oidc"`
}

type OIDCInfo struct {
	Issuer   string `json:"issuer"`
	ClientID string `json:"client_id"`
}
