package cadreen

import (
	"context"
	"fmt"
)

type AssessRequest struct {
	Task   string `json:"task"`
	Domain string `json:"domain,omitempty"`
}

func (c *Client) Assess(ctx context.Context, task string, opts ...RequestOption) (*Assessment, error) {
	req := AssessRequest{Task: task}
	var result Assessment
	if err := c.do(ctx, "POST", "/api/v1/cadreen/assess", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("assess: %w", err)
	}
	return &result, nil
}

func (c *Client) AssessWithDomain(ctx context.Context, task, domain string, opts ...RequestOption) (*Assessment, error) {
	req := AssessRequest{Task: task, Domain: domain}
	var result Assessment
	if err := c.do(ctx, "POST", "/api/v1/cadreen/assess", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("assess: %w", err)
	}
	return &result, nil
}
