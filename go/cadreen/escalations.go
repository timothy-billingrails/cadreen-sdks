package cadreen

import (
	"context"
	"fmt"
)

type ResolveEscalationRequest struct {
	Resolution string `json:"resolution"`
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
	if err := c.do(ctx, "GET", "/api/v1/cadreen/escalations/"+id, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get escalation: %w", err)
	}
	return &result, nil
}

func (c *Client) ResolveEscalation(ctx context.Context, id string, resolution string, opts ...RequestOption) (*Escalation, error) {
	req := ResolveEscalationRequest{Resolution: resolution}
	var result Escalation
	if err := c.do(ctx, "POST", "/api/v1/cadreen/escalations/"+id+"/resolve", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("resolve escalation: %w", err)
	}
	return &result, nil
}
