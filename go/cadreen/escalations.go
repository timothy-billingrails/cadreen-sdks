package cadreen

import (
	"context"
	"fmt"
	"net/url"
)

type ResolveEscalationRequest struct {
	Decision string `json:"decision"`
}

func (c *Client) ListEscalations(ctx context.Context, opts ...RequestOption) (*ListEscalationsResponse, error) {
	var result ListEscalationsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/escalations", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list escalations: %w", err)
	}
	return &result, nil
}

func (c *Client) GetEscalation(ctx context.Context, id string, opts ...RequestOption) (*Escalation, error) {
	var result Escalation
	if err := c.do(ctx, "GET", "/api/v1/cadreen/escalations/"+url.PathEscape(id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get escalation: %w", err)
	}
	return &result, nil
}

func (c *Client) ResolveEscalation(ctx context.Context, id string, decision string, opts ...RequestOption) (*Escalation, error) {
	req := ResolveEscalationRequest{Decision: decision}
	var result Escalation
	if err := c.do(ctx, "POST", "/api/v1/cadreen/escalations/"+url.PathEscape(id)+"/resolve", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("resolve escalation: %w", err)
	}
	return &result, nil
}
