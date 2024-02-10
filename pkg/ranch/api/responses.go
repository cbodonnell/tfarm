package api

import "time"

type InfoResponse struct {
	Ready   bool        `json:"ready"`
	Version string      `json:"version"`
	OIDC    OIDCReponse `json:"oidc"`
}

type OIDCReponse struct {
	Issuer   string `json:"issuer"`
	ClientID string `json:"client_id"`
}

type ClientResponse struct {
	ClientID      string     `json:"client_id"`
	ClientSecret  string     `json:"client_secret,omitempty"`
	ClientTLSCert string     `json:"client_tls_cert"`
	ClientTLSKey  string     `json:"client_tls_key,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
}
