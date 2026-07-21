package cadreen

import (
	"context"
	"fmt"
)

type CreateFederationRequest struct {
	TargetWorkspaceID string         `json:"targetWorkspaceId"`
	Description       string         `json:"description,omitempty"`
	Permissions       map[string]any `json:"permissions,omitempty"`
}

type SuspendFederationRequest struct {
	Reason string `json:"reason,omitempty"`
}

type RevokeFederationRequest struct {
	Reason string `json:"reason,omitempty"`
}

type UpdateFederationPermissionsRequest struct {
	Permissions map[string]any `json:"permissions"`
}

type LinkFederationAgentRequest struct {
	LocalAgentID  string `json:"localAgentId"`
	RemoteAgentID string `json:"remoteAgentId"`
}

func (c *Client) CreateFederation(ctx context.Context, req CreateFederationRequest, opts ...RequestOption) (*FederationLink, error) {
	var result FederationLink
	if err := c.do(ctx, "POST", "/api/v1/cadreen/federation", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create federation: %w", err)
	}
	return &result, nil
}

func (c *Client) ListFederations(ctx context.Context, opts ...RequestOption) (*ListFederationsResponse, error) {
	var result ListFederationsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/federation", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list federations: %w", err)
	}
	return &result, nil
}

func (c *Client) GetFederation(ctx context.Context, federationID string, opts ...RequestOption) (*FederationLink, error) {
	var result FederationLink
	if err := c.do(ctx, "GET", "/api/v1/cadreen/federation/"+federationID, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get federation: %w", err)
	}
	return &result, nil
}

func (c *Client) ApproveFederation(ctx context.Context, federationID string, opts ...RequestOption) (*FederationLink, error) {
	var result FederationLink
	if err := c.do(ctx, "POST", "/api/v1/cadreen/federation/"+federationID+"/approve", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("approve federation: %w", err)
	}
	return &result, nil
}

func (c *Client) SuspendFederation(ctx context.Context, federationID string, req SuspendFederationRequest, opts ...RequestOption) (*FederationLink, error) {
	var result FederationLink
	if err := c.do(ctx, "POST", "/api/v1/cadreen/federation/"+federationID+"/suspend", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("suspend federation: %w", err)
	}
	return &result, nil
}

func (c *Client) RevokeFederation(ctx context.Context, federationID string, req RevokeFederationRequest, opts ...RequestOption) (*FederationLink, error) {
	var result FederationLink
	if err := c.do(ctx, "POST", "/api/v1/cadreen/federation/"+federationID+"/revoke", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("revoke federation: %w", err)
	}
	return &result, nil
}

func (c *Client) GetFederationPermissions(ctx context.Context, federationID string, opts ...RequestOption) (*FederationPermissions, error) {
	var result FederationPermissions
	if err := c.do(ctx, "GET", "/api/v1/cadreen/federation/"+federationID+"/permissions", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get federation permissions: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateFederationPermissions(ctx context.Context, federationID string, req UpdateFederationPermissionsRequest, opts ...RequestOption) (*FederationPermissions, error) {
	var result FederationPermissions
	if err := c.do(ctx, "PUT", "/api/v1/cadreen/federation/"+federationID+"/permissions", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("update federation permissions: %w", err)
	}
	return &result, nil
}

func (c *Client) LinkFederationAgent(ctx context.Context, federationID string, req LinkFederationAgentRequest, opts ...RequestOption) (*FederationAgent, error) {
	var result FederationAgent
	if err := c.do(ctx, "POST", "/api/v1/cadreen/federation/"+federationID+"/agents", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("link federation agent: %w", err)
	}
	return &result, nil
}

func (c *Client) ListFederationAgents(ctx context.Context, federationID string, opts ...RequestOption) (*ListFederationAgentsResponse, error) {
	var result ListFederationAgentsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/federation/"+federationID+"/agents", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list federation agents: %w", err)
	}
	return &result, nil
}

func (c *Client) UnlinkFederationAgent(ctx context.Context, federationID, agentLinkID string, opts ...RequestOption) error {
	path := fmt.Sprintf("/api/v1/cadreen/federation/%s/agents/%s", federationID, agentLinkID)
	if err := c.do(ctx, "DELETE", path, nil, nil, opts...); err != nil {
		return fmt.Errorf("unlink federation agent: %w", err)
	}
	return nil
}
