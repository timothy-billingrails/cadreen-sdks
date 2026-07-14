package cadreen

import (
	"context"
	"fmt"
	"net/url"
)

type ConnectExternalAgentRequest struct {
	AgentCardURL string `json:"agentCardUrl"`
}

type ListExternalConnectionsParams struct {
	Status string `json:"status,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type ListExternalInteractionsParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func (c *Client) ConnectExternalAgent(ctx context.Context, agentID string, req ConnectExternalAgentRequest) (*ExternalAgentConnection, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external", url.PathEscape(agentID))
	var result ExternalAgentConnection
	if err := c.do(ctx, "POST", path, req, &result); err != nil {
		return nil, fmt.Errorf("connect external agent: %w", err)
	}
	return &result, nil
}

func (c *Client) ListExternalConnections(ctx context.Context, agentID string, params ListExternalConnectionsParams) (*ListExternalConnectionsResponse, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external", url.PathEscape(agentID))
	var result ListExternalConnectionsResponse
	if err := c.do(ctx, "GET", path, params, &result); err != nil {
		return nil, fmt.Errorf("list external connections: %w", err)
	}
	return &result, nil
}

func (c *Client) GetExternalConnection(ctx context.Context, agentID, connectionID string) (*ExternalAgentConnection, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s", url.PathEscape(agentID), url.PathEscape(connectionID))
	var result ExternalAgentConnection
	if err := c.do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("get external connection: %w", err)
	}
	return &result, nil
}

func (c *Client) ApproveExternalConnection(ctx context.Context, agentID, connectionID string) (map[string]string, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/approve", url.PathEscape(agentID), url.PathEscape(connectionID))
	var result map[string]string
	if err := c.do(ctx, "POST", path, nil, &result); err != nil {
		return nil, fmt.Errorf("approve external connection: %w", err)
	}
	return result, nil
}

func (c *Client) SuspendExternalConnection(ctx context.Context, agentID, connectionID string) (map[string]string, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/suspend", url.PathEscape(agentID), url.PathEscape(connectionID))
	var result map[string]string
	if err := c.do(ctx, "POST", path, nil, &result); err != nil {
		return nil, fmt.Errorf("suspend external connection: %w", err)
	}
	return result, nil
}

func (c *Client) RevokeExternalConnection(ctx context.Context, agentID, connectionID string) (map[string]string, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/revoke", url.PathEscape(agentID), url.PathEscape(connectionID))
	var result map[string]string
	if err := c.do(ctx, "POST", path, nil, &result); err != nil {
		return nil, fmt.Errorf("revoke external connection: %w", err)
	}
	return result, nil
}

func (c *Client) DeleteExternalConnection(ctx context.Context, agentID, connectionID string) error {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s", url.PathEscape(agentID), url.PathEscape(connectionID))
	return c.do(ctx, "DELETE", path, nil, nil)
}

func (c *Client) ListExternalInteractions(ctx context.Context, agentID, connectionID string, params ListExternalInteractionsParams) (*ListExternalInteractionsResponse, error) {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/interactions", url.PathEscape(agentID), url.PathEscape(connectionID))
	var result ListExternalInteractionsResponse
	if err := c.do(ctx, "GET", path, params, &result); err != nil {
		return nil, fmt.Errorf("list external interactions: %w", err)
	}
	return &result, nil
}

func (c *Client) GetExternalAgentSettings(ctx context.Context) (*ExternalAgentSettings, error) {
	var result ExternalAgentSettings
	if err := c.do(ctx, "GET", "/api/v1/cadreen/external-agents/settings", nil, &result); err != nil {
		return nil, fmt.Errorf("get external agent settings: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateExternalAgentSettings(ctx context.Context, enabled bool) (*ExternalAgentSettings, error) {
	var result ExternalAgentSettings
	if err := c.do(ctx, "PUT", "/api/v1/cadreen/external-agents/settings", map[string]bool{"enabled": enabled}, &result); err != nil {
		return nil, fmt.Errorf("update external agent settings: %w", err)
	}
	return &result, nil
}

func (c *Client) ListAllExternalConnections(ctx context.Context) (*ListExternalConnectionsResponse, error) {
	var result ListExternalConnectionsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/external-agents/connections", nil, &result); err != nil {
		return nil, fmt.Errorf("list all external connections: %w", err)
	}
	return &result, nil
}
