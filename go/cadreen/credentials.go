package cadreen

import (
	"context"
	"fmt"
)

type CreateCredentialRequest struct {
	Provider string                 `json:"provider"`
	Name     string                 `json:"name,omitempty"`
	KeyData  map[string]any `json:"key_data"`
}

func (c *Client) ListCredentials(ctx context.Context, opts ...RequestOption) (*ListCredentialsResponse, error) {
	var result ListCredentialsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/credentials", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list credentials: %w", err)
	}
	return &result, nil
}

func (c *Client) CreateCredential(ctx context.Context, req CreateCredentialRequest, opts ...RequestOption) (*CredentialMetadata, error) {
	var result CredentialMetadata
	if err := c.do(ctx, "POST", "/api/v1/cadreen/credentials", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create credential: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteCredential(ctx context.Context, id string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/credentials/"+id, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete credential: %w", err)
	}
	return nil
}
