package auth

import (
	"context"
	"fmt"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
)

// OIDCClientConfig is a struct that contains the configuration for an OIDCClient.
type OIDCClientConfig struct {
	Issuer   string
	ClientID string
}

func NewOIDCClient(ctx context.Context, cfg *OIDCClientConfig) (*oauth.OIDCClient, error) {
	oc, err := oauth.NewOIDCClient(ctx, cfg.Issuer, cfg.ClientID, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating OIDC client: %w", err)
	}

	return oc, nil
}
