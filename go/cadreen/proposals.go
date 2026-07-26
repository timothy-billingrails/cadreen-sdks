package cadreen

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type ListProposalsOptions struct {
	Status string
	Limit  int
}

type ListProposalsResponse struct {
	Proposals []TaskProposal `json:"proposals"`
	Count     int            `json:"count"`
}

type AcceptProposalResponse struct {
	Status       string                 `json:"status"`
	ExecutionID  string                 `json:"execution_id"`
	Action       string                 `json:"action"`
	Intent       string                 `json:"intent"`
	NextStep     string                 `json:"next_step"`
	AutoApproved *bool                  `json:"auto_approved,omitempty"`
	Result       map[string]interface{} `json:"result,omitempty"`
}

type DismissProposalRequest struct {
	Reason string `json:"reason,omitempty"`
}

type DismissProposalResponse struct {
	Status string `json:"status"`
}

type ProposalStatsResponse struct {
	Proposed  int `json:"proposed"`
	Accepted  int `json:"accepted"`
	Dismissed int `json:"dismissed"`
	Expired   int `json:"expired"`
}

func (c *Client) ListProposals(ctx context.Context, opts ListProposalsOptions, reqOpts ...RequestOption) (*ListProposalsResponse, error) {
	path := "/api/v1/cadreen/proposals"
	params := url.Values{}
	if opts.Status != "" {
		params.Set("status", opts.Status)
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result ListProposalsResponse
	if err := c.do(ctx, "GET", path, nil, &result, reqOpts...); err != nil {
		return nil, fmt.Errorf("list proposals: %w", err)
	}
	return &result, nil
}

func (c *Client) GetProposal(ctx context.Context, id string, opts ...RequestOption) (*TaskProposal, error) {
	var result TaskProposal
	if err := c.do(ctx, "GET", "/api/v1/cadreen/proposals/"+url.PathEscape(id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get proposal: %w", err)
	}
	return &result, nil
}

func (c *Client) AcceptProposal(ctx context.Context, id string, opts ...RequestOption) (*AcceptProposalResponse, error) {
	var result AcceptProposalResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/proposals/"+url.PathEscape(id)+"/accept", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("accept proposal: %w", err)
	}
	return &result, nil
}

func (c *Client) DismissProposal(ctx context.Context, id string, req DismissProposalRequest, opts ...RequestOption) (*DismissProposalResponse, error) {
	var result DismissProposalResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/proposals/"+url.PathEscape(id)+"/dismiss", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("dismiss proposal: %w", err)
	}
	return &result, nil
}

func (c *Client) ProposalStats(ctx context.Context, opts ...RequestOption) (*ProposalStatsResponse, error) {
	var result ProposalStatsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/proposals/stats", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("proposal stats: %w", err)
	}
	return &result, nil
}
