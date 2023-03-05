package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
)

// TODO: RANCH_OIDC_ISSUER and RANCH_OIDC_CLIENT_ID should be arguments to this function?
func NewOIDCClient(ctx context.Context) (*oauth.OIDCClient, error) {
	ranchOauthIssuer := os.Getenv("RANCH_OIDC_ISSUER")
	if ranchOauthIssuer == "" {
		ranchOauthIssuer = "https://auth.tunnel.farm/realms/tunnel.farm"
	}

	ranchOauthClientID := os.Getenv("RANCH_OIDC_CLIENT_ID")
	if ranchOauthClientID == "" {
		ranchOauthClientID = "tfarm-cli"
	}

	oc, err := oauth.NewOIDCClient(ctx, ranchOauthIssuer, ranchOauthClientID, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating OIDC client: %w", err)
	}

	return oc, nil
}
