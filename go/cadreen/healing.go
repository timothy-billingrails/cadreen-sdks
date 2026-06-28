package cadreen

import (
	"context"
	"fmt"
)

type DiagnoseRequest struct {
	ErrorMessage string  `json:"error_message"`
	ToolName     string  `json:"tool_name,omitempty"`
	TraceID      string  `json:"trace_id,omitempty"`
}

type ListHealingPrecedentsResponse struct {
	Precedents []HealingPrecedent `json:"precedents"`
	Count      int                `json:"count"`
	Pagination Pagination         `json:"pagination,omitempty"`
}

func (c *Client) HealingStats(ctx context.Context, opts ...RequestOption) (*HealingStatsResponse, error) {
	var result HealingStatsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/healing/stats", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("healing stats: %w", err)
	}
	return &result, nil
}

func (c *Client) ListHealingPrecedents(ctx context.Context, opts ...RequestOption) (*ListHealingPrecedentsResponse, error) {
	var result ListHealingPrecedentsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/healing/precedents", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list healing precedents: %w", err)
	}
	return &result, nil
}

func (c *Client) Diagnose(ctx context.Context, req DiagnoseRequest, opts ...RequestOption) (*HealingDiagnosis, error) {
	var result HealingDiagnosis
	if err := c.do(ctx, "POST", "/api/v1/cadreen/healing/diagnose", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("diagnose: %w", err)
	}
	return &result, nil
}
