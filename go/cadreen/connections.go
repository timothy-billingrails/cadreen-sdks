package cadreen

import (
	"context"
	"fmt"
)

func (c *Client) RegisterOpenAPI(ctx context.Context, req RegisterOpenAPIRequest, opts ...RequestOption) (*RegisterOpenAPIResponse, error) {
	var result RegisterOpenAPIResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/connections/openapi", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("register openapi: %w", err)
	}
	return &result, nil
}

func (c *Client) ListConnections(ctx context.Context, opts ...RequestOption) (*ListConnectionsResponse, error) {
	var result ListConnectionsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/connections", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list connections: %w", err)
	}
	return &result, nil
}

func (c *Client) ListCapabilities(ctx context.Context, opts ...RequestOption) (*ListCapabilitiesResponse, error) {
	var result ListCapabilitiesResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/capabilities", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list capabilities: %w", err)
	}
	return &result, nil
}
