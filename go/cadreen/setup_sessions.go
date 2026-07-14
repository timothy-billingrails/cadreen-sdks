package cadreen

import (
	"context"
	"fmt"
)

type SetupSessionCreateRequest struct {
	WorkspaceID string   `json:"workspace_id,omitempty"`
	Purpose     string   `json:"purpose,omitempty"`
	Constraints []string `json:"constraints,omitempty"`
}

type SetupSessionAddRequest struct {
	Connections []SetupConnection `json:"connections,omitempty"`
	Credentials []SetupCredential `json:"credentials,omitempty"`
	Memory      []SetupMemory     `json:"memory,omitempty"`
	Policies    []SetupPolicy     `json:"policies,omitempty"`
}

type SetupSessionApplyRequest struct {
	Confirm bool `json:"confirm"`
}

type SetupSession struct {
	ID           string            `json:"id"`
	Status       string            `json:"status"`
	Purpose      string            `json:"purpose,omitempty"`
	Constraints  []string          `json:"constraints,omitempty"`
	Connections  []SetupConnection `json:"connections"`
	Credentials  []SetupCredential `json:"credentials"`
	Memory       []SetupMemory     `json:"memory"`
	Policies     []SetupPolicy     `json:"policies"`
	Proposals    []SetupProposal   `json:"proposals,omitempty"`
	AppliedCount int               `json:"applied_count"`
	FailedCount  int               `json:"failed_count"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	AppliedAt    string            `json:"applied_at,omitempty"`
}

type SetupSessionApplyResult struct {
	SessionID string       `json:"session_id"`
	Status    string       `json:"status"`
	Applied   int          `json:"applied"`
	Failed    int          `json:"failed"`
	Result    *SetupResult `json:"result,omitempty"`
}

type ListSetupSessionsResponse struct {
	Sessions []SetupSession `json:"sessions"`
}

func (c *Client) CreateSetupSession(ctx context.Context, req SetupSessionCreateRequest, opts ...RequestOption) (*SetupSession, error) {
	var result SetupSession
	if err := c.do(ctx, "POST", "/api/v1/cadreen/setup/sessions", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create setup session: %w", err)
	}
	return &result, nil
}

func (c *Client) ListSetupSessions(ctx context.Context, opts ...RequestOption) (*ListSetupSessionsResponse, error) {
	var result ListSetupSessionsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/setup/sessions", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list setup sessions: %w", err)
	}
	return &result, nil
}

func (c *Client) GetSetupSession(ctx context.Context, id string, opts ...RequestOption) (*SetupSession, error) {
	var result SetupSession
	if err := c.do(ctx, "GET", "/api/v1/cadreen/setup/sessions/"+id, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get setup session: %w", err)
	}
	return &result, nil
}

func (c *Client) AddToSetupSession(ctx context.Context, id string, req SetupSessionAddRequest, opts ...RequestOption) (*SetupSession, error) {
	var result SetupSession
	if err := c.do(ctx, "POST", "/api/v1/cadreen/setup/sessions/"+id, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("add to setup session: %w", err)
	}
	return &result, nil
}

func (c *Client) ApplySetupSession(ctx context.Context, id string, req SetupSessionApplyRequest, opts ...RequestOption) (*SetupSessionApplyResult, error) {
	var result SetupSessionApplyResult
	if err := c.do(ctx, "POST", "/api/v1/cadreen/setup/sessions/"+id+"/apply", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("apply setup session: %w", err)
	}
	return &result, nil
}
