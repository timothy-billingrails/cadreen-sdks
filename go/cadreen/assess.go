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
	var wrapper struct {
		Assessment Assessment `json:"assessment"`
	}
	if err := c.do(ctx, "POST", "/api/v1/cadreen/assess", req, &wrapper, opts...); err != nil {
		return nil, fmt.Errorf("assess: %w", err)
	}
	return &wrapper.Assessment, nil
}

func (c *Client) AssessWithDomain(ctx context.Context, task, domain string, opts ...RequestOption) (*Assessment, error) {
	req := AssessRequest{Task: task, Domain: domain}
	var wrapper2 struct {
		Assessment Assessment `json:"assessment"`
	}
	if err := c.do(ctx, "POST", "/api/v1/cadreen/assess", req, &wrapper2, opts...); err != nil {
		return nil, fmt.Errorf("assess: %w", err)
	}
	return &wrapper2.Assessment, nil
}
