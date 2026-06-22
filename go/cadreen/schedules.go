package cadreen

import (
	"context"
	"fmt"
)

type Schedule struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	BlueprintID      string                 `json:"blueprint_id"`
	BlueprintVersion int                    `json:"blueprint_version"`
	Status           string                 `json:"status"`
	Trigger          map[string]interface{} `json:"trigger"`
	Timezone         string                 `json:"timezone"`
	Params           map[string]interface{} `json:"params,omitempty"`
	NextRunAt        string                 `json:"next_run_at,omitempty"`
	LastRunAt        string                 `json:"last_run_at,omitempty"`
	PauseReason      string                 `json:"pause_reason,omitempty"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

type CreateScheduleRequest struct {
	BlueprintID string                 `json:"blueprint_id"`
	Name        string                 `json:"name"`
	Trigger     map[string]interface{} `json:"trigger"`
	Timezone    string                 `json:"timezone,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
}

type ListSchedulesResponse struct {
	Schedules []Schedule `json:"schedules"`
	Count     int        `json:"count"`
}

func (c *Client) ListSchedules(ctx context.Context, opts ...RequestOption) (*ListSchedulesResponse, error) {
	var result ListSchedulesResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/schedules", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}
	return &result, nil
}

func (c *Client) GetSchedule(ctx context.Context, id string, opts ...RequestOption) (*Schedule, error) {
	var result Schedule
	if err := c.do(ctx, "GET", fmt.Sprintf("/api/v1/cadreen/schedules/%s", id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	return &result, nil
}

func (c *Client) CreateSchedule(ctx context.Context, req CreateScheduleRequest, opts ...RequestOption) (*Schedule, error) {
	var result Schedule
	if err := c.do(ctx, "POST", "/api/v1/cadreen/schedules", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}
	return &result, nil
}

func (c *Client) PauseSchedule(ctx context.Context, id string, reason string, opts ...RequestOption) error {
	payload := map[string]interface{}{}
	if reason != "" {
		payload["reason"] = reason
	}
	if err := c.do(ctx, "POST", fmt.Sprintf("/api/v1/cadreen/schedules/%s/pause", id), payload, nil, opts...); err != nil {
		return fmt.Errorf("pause schedule: %w", err)
	}
	return nil
}

func (c *Client) ResumeSchedule(ctx context.Context, id string, opts ...RequestOption) error {
	if err := c.do(ctx, "POST", fmt.Sprintf("/api/v1/cadreen/schedules/%s/resume", id), nil, nil, opts...); err != nil {
		return fmt.Errorf("resume schedule: %w", err)
	}
	return nil
}
