package cadreen

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type ListTracesOptions struct {
	Domain   string
	Decision string
	From     string
	To       string
	Limit    int
	Offset   int
}

func (c *Client) GetTrace(ctx context.Context, traceID string, opts ...RequestOption) (*IntelligenceTraceEntry, error) {
	path := fmt.Sprintf("/api/v1/cadreen/intelligence/%s", traceID)
	var result IntelligenceTraceEntry
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get trace: %w", err)
	}
	return &result, nil
}

func (c *Client) ListTraces(ctx context.Context, options *ListTracesOptions, opts ...RequestOption) (*ListIntelligenceResponse, error) {
	u, err := url.Parse("/api/v1/cadreen/intelligence")
	if err != nil {
		return nil, fmt.Errorf("list traces: parse url: %w", err)
	}

	q := u.Query()
	if options != nil {
		if options.Domain != "" {
			q.Set("domain", options.Domain)
		}
		if options.Decision != "" {
			q.Set("decision", options.Decision)
		}
		if options.From != "" {
			q.Set("from", options.From)
		}
		if options.To != "" {
			q.Set("to", options.To)
		}
		if options.Limit > 0 {
			q.Set("limit", strconv.Itoa(options.Limit))
		}
		if options.Offset > 0 {
			q.Set("offset", strconv.Itoa(options.Offset))
		}
	}
	u.RawQuery = q.Encode()

	var result ListIntelligenceResponse
	if err := c.do(ctx, "GET", u.String(), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list traces: %w", err)
	}
	return &result, nil
}

func (c *Client) IntelligenceStats(ctx context.Context, opts ...RequestOption) (*IntelligenceStats, error) {
	var result IntelligenceStats
	if err := c.do(ctx, "GET", "/api/v1/cadreen/intelligence/stats", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("intelligence stats: %w", err)
	}
	return &result, nil
}
