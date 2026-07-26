package cadreen

import (
	"context"
	"fmt"
	"net/url"
)

type CreatePolicyRequest struct {
	Name      string                   `json:"name"`
	Rules     []map[string]any `json:"rules,omitempty"`
	Domain    string                   `json:"domain,omitempty"`
	AutoDraft bool                     `json:"auto_draft,omitempty"`
}

type EvaluatePolicyRequest struct {
	Action  string                 `json:"action"`
	Domain  string                 `json:"domain,omitempty"`
	Context map[string]any `json:"context,omitempty"`
}

func (c *Client) Draft(ctx context.Context, req CreatePolicyRequest, opts ...RequestOption) (*CreatePolicyResponse, error) {
	var result CreatePolicyResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/policies", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("draft policy: %w", err)
	}
	return &result, nil
}

func (c *Client) Evaluate(ctx context.Context, action, domain string, opts ...RequestOption) (*EvaluatePolicyResponse, error) {
	req := EvaluatePolicyRequest{
		Action: action,
		Domain: domain,
	}
	var result EvaluatePolicyResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/policies/evaluate", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("evaluate policy: %w", err)
	}
	return &result, nil
}

func (c *Client) EvaluatePolicy(ctx context.Context, req EvaluatePolicyRequest, opts ...RequestOption) (*EvaluatePolicyResponse, error) {
	var result EvaluatePolicyResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/policies/evaluate", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("evaluate policy: %w", err)
	}
	return &result, nil
}

func (c *Client) Confirm(ctx context.Context, policyID string, opts ...RequestOption) (*ConfirmPolicyResponse, error) {
	var result ConfirmPolicyResponse
	path := fmt.Sprintf("/api/v1/cadreen/policies/%s/confirm", policyID)
	if err := c.do(ctx, "POST", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("confirm policy: %w", err)
	}
	return &result, nil
}

func (c *Client) ListPolicies(ctx context.Context, opts ...RequestOption) (*ListPoliciesResponse, error) {
	var result ListPoliciesResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/policies", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list policies: %w", err)
	}
	return &result, nil
}

func (c *Client) GetPolicy(ctx context.Context, id string, opts ...RequestOption) (*PolicyBundle, error) {
	var result PolicyBundle
	if err := c.do(ctx, "GET", "/api/v1/cadreen/policies/"+url.PathEscape(id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get policy: %w", err)
	}
	return &result, nil
}
