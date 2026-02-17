package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	pkglog "sb-scanner/pkg/logger"
)

type AuthenticatedClient struct {
	*DefaultClient

	// GitHub App credentials
	appID          string
	installationID string
	keyBytes       []byte // PEM encoded RSA private key

	// issued installation access token
	token          string
	tokenExpiresAt time.Time
}

// Creates authenticated GitHub API client using GitHub App credentials.
func NewAuthenticatedClient(appID, installationID string, keyBytes []byte) (*AuthenticatedClient, error) {
	authenticatedCli := &AuthenticatedClient{
		DefaultClient: &DefaultClient{
			logger:  pkglog.GetLogger().With("pkg", "github"),
			httpcli: http.DefaultClient,
		},
		appID:          appID,
		installationID: installationID,
		keyBytes:       keyBytes,
	}
	if err := authenticatedCli.refreshToken(); err != nil {
		return nil, fmt.Errorf("failed to refresh github auth token: %w", err)
	}

	return authenticatedCli, nil
}

func (c *AuthenticatedClient) refreshToken() error {
	pk, err := jwt.ParseRSAPrivateKeyFromPEM(c.keyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse RSA private key: %w", err)
	}
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    c.appID,
		IssuedAt:  jwt.NewNumericDate(now.Add(-1 * time.Minute)),
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
	})
	signedJWT, err := token.SignedString(pk)
	if err != nil {
		return fmt.Errorf("failed to sign JWT: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", c.installationID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Accept", acceptHeaderValue)
	req.Header.Add(apiVersionHeaderKey, apiVersionHeaderValue)
	req.Header.Add("User-Agent", userAgentHeaderValue)
	req.Header.Add("Authorization", "Bearer "+signedJWT)

	resp, err := c.httpcli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	var result accessTokenResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	c.token = result.Token
	c.tokenExpiresAt = result.ExpiresAt
	c.logger.Debug("refreshed github installation access token")

	return nil
}

func (c *AuthenticatedClient) SearchCommits(ctx context.Context, keywords []string, opts ...SearchOption) (SearchResult, error) {
	// not considering concurrent use for now
	if time.Until(c.tokenExpiresAt) < 5*time.Minute {
		if err := c.refreshToken(); err != nil {
			return SearchResult{}, fmt.Errorf("failed to refresh github auth token: %w", err)
		}
	}
	opts = append(opts, withAuthToken(c.token))

	return c.DefaultClient.SearchCommits(ctx, keywords, opts...)
}
