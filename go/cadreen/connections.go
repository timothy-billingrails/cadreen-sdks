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

// Catalog returns the unified marketplace catalog of all available integrations.
func (c *Client) Catalog(ctx context.Context, opts ...RequestOption) (*CatalogResponse, error) {
	var result CatalogResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/connections/catalog", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}
	return &result, nil
}

// Install starts a one-click install flow for an integration from the catalog.
// Returns an OAuth URL for marketplace integrations.
func (c *Client) Install(ctx context.Context, integrationID string, opts ...RequestOption) (*InstallResponse, error) {
	var result InstallResponse
	body := map[string]string{"integration_id": integrationID}
	if err := c.do(ctx, "POST", "/api/v1/cadreen/connections/install", body, &result, opts...); err != nil {
		return nil, fmt.Errorf("install: %w", err)
	}
	return &result, nil
}

// RegisterMCP registers a Model Context Protocol server.
func (c *Client) RegisterMCP(ctx context.Context, req RegisterMCPRequest, opts ...RequestOption) (*RegisterMCPResponse, error) {
	var result RegisterMCPResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/connections/mcp", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("register mcp: %w", err)
	}
	return &result, nil
}

// DeleteConnection removes a connection by ID.
func (c *Client) DeleteConnection(ctx context.Context, id string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/connections/"+id, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete connection: %w", err)
	}
	return nil
}

// Connect resolves a capability name to the best available connection.
func (c *Client) Connect(ctx context.Context, capability string, opts ...RequestOption) (*ConnectResult, error) {
	var result ConnectResult
	body := map[string]string{"capability": capability}
	if err := c.do(ctx, "POST", "/api/v1/cadreen/connections", body, &result, opts...); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return &result, nil
}
