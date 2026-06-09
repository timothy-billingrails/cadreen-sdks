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

type ReplayRequest struct {
	Mode string `json:"mode"`
}

type ReplayResult struct {
	TraceID           string                 `json:"trace_id"`
	Mode              string                 `json:"mode"`
	Domain            string                 `json:"domain"`
	OriginalGate      string                 `json:"original_gate"`
	OriginalConfidence float64               `json:"original_confidence"`
	CurrentGate       string                 `json:"current_gate"`
	CurrentConfidence float64                `json:"current_confidence"`
	GateChanged       bool                   `json:"gate_changed"`
	ChangeSummary     string                 `json:"change_summary"`
	CurrentCapability map[string]interface{} `json:"current_capability"`
	CurrentMemory     map[string]interface{} `json:"current_memory"`
	CurrentGaps       map[string]interface{} `json:"current_gaps"`
	ReplayNote        string                 `json:"replay_note"`
}

func (c *Client) Replay(ctx context.Context, traceID string, mode string, opts ...RequestOption) (*ReplayResult, error) {
	path := fmt.Sprintf("/api/v1/cadreen/intelligence/%s/replay", traceID)
	req := ReplayRequest{Mode: mode}
	if mode == "" {
		req.Mode = "current"
	}
	var result ReplayResult
	if err := c.do(ctx, "POST", path, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("replay: %w", err)
	}
	return &result, nil
}

type HandoffPacket struct {
	TraceID               string                   `json:"trace_id"`
	Domain                string                   `json:"domain"`
	CreatedAt             string                   `json:"created_at"`
	Governance            map[string]interface{}   `json:"governance"`
	WhatTheSystemKnew     map[string]interface{}   `json:"what_the_system_knew"`
	WhatTheSystemDidntKnow map[string]interface{}  `json:"what_the_system_didnt_know"`
	WhatHappened          map[string]interface{}   `json:"what_happened"`
	SuggestedActions      []map[string]interface{} `json:"suggested_actions"`
	NextAction            map[string]interface{}   `json:"next_action"`
	TraceURL              string                   `json:"trace_url"`
}

func (c *Client) Handoff(ctx context.Context, traceID string, opts ...RequestOption) (*HandoffPacket, error) {
	path := fmt.Sprintf("/api/v1/cadreen/intelligence/%s/handoff", traceID)
	var result HandoffPacket
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("handoff: %w", err)
	}
	return &result, nil
}

type PromoteRequest struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
}

type PromoteResult struct {
	ID            string   `json:"id"`
	Kind          string   `json:"kind"`
	Status        string   `json:"status"`
	ToolName      string   `json:"tool_name,omitempty"`
	ToolSequence  []string `json:"tool_sequence,omitempty"`
	SourceTraceID string   `json:"source_trace_id"`
}

func (c *Client) Promote(ctx context.Context, traceID string, kind string, name string, opts ...RequestOption) (*PromoteResult, error) {
	path := fmt.Sprintf("/api/v1/cadreen/intelligence/%s/promote", traceID)
	req := PromoteRequest{Kind: kind, Name: name}
	var result PromoteResult
	if err := c.do(ctx, "POST", path, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("promote: %w", err)
	}
	return &result, nil
}
