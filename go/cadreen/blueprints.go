package cadreen

import (
	"context"
	"fmt"
)

type Blueprint struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	Intent      string `json:"intent,omitempty"`
	SourceType  string `json:"source_type,omitempty"`
	SourceID    string `json:"source_id,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type BlueprintRun struct {
	ID               string                 `json:"id"`
	BlueprintID      string                 `json:"blueprint_id"`
	BlueprintVersion int                    `json:"blueprint_version"`
	Status           string                 `json:"status"`
	Params           map[string]interface{} `json:"params,omitempty"`
	ResultSummary    string                 `json:"result_summary,omitempty"`
	TraceID          string                 `json:"trace_id,omitempty"`
	CreatedAt        string                 `json:"created_at"`
}

type CreateBlueprintRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description,omitempty"`
	Source          *BlueprintSource       `json:"source,omitempty"`
	ParameterSchema map[string]interface{} `json:"parameter_schema,omitempty"`
	DefaultParams   map[string]interface{} `json:"default_params,omitempty"`
}

type BlueprintSource struct {
	Type        string `json:"type"`
	TraceID     string `json:"trace_id,omitempty"`
	ExecutionID string `json:"execution_id,omitempty"`
}

type ListBlueprintsResponse struct {
	Blueprints []Blueprint `json:"blueprints"`
	Count      int         `json:"count"`
}

type ListBlueprintRunsResponse struct {
	Runs  []BlueprintRun `json:"runs"`
	Count int            `json:"count"`
}

func (c *Client) ListBlueprints(ctx context.Context, status string, limit int, opts ...RequestOption) (*ListBlueprintsResponse, error) {
	var result ListBlueprintsResponse
	path := fmt.Sprintf("/api/v1/cadreen/blueprints?status=%s&limit=%d", status, limit)
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list blueprints: %w", err)
	}
	return &result, nil
}

func (c *Client) GetBlueprint(ctx context.Context, id string, opts ...RequestOption) (*Blueprint, error) {
	var result Blueprint
	if err := c.do(ctx, "GET", fmt.Sprintf("/api/v1/cadreen/blueprints/%s", id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get blueprint: %w", err)
	}
	return &result, nil
}

func (c *Client) CreateBlueprint(ctx context.Context, req CreateBlueprintRequest, opts ...RequestOption) (*Blueprint, error) {
	var result Blueprint
	if err := c.do(ctx, "POST", "/api/v1/cadreen/blueprints", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create blueprint: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteBlueprint(ctx context.Context, id string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/cadreen/blueprints/%s", id), nil, nil, opts...); err != nil {
		return fmt.Errorf("delete blueprint: %w", err)
	}
	return nil
}

func (c *Client) RunBlueprint(ctx context.Context, id string, params map[string]interface{}, opts ...RequestOption) (*BlueprintRun, error) {
	var result BlueprintRun
	payload := map[string]interface{}{}
	if len(params) > 0 {
		payload["params"] = params
	}
	if err := c.do(ctx, "POST", fmt.Sprintf("/api/v1/cadreen/blueprints/%s/runs", id), payload, &result, opts...); err != nil {
		return nil, fmt.Errorf("run blueprint: %w", err)
	}
	return &result, nil
}

func (c *Client) ListBlueprintRuns(ctx context.Context, id string, limit int, opts ...RequestOption) (*ListBlueprintRunsResponse, error) {
	var result ListBlueprintRunsResponse
	if err := c.do(ctx, "GET", fmt.Sprintf("/api/v1/cadreen/blueprints/%s/runs?limit=%d", id, limit), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list blueprint runs: %w", err)
	}
	return &result, nil
}
