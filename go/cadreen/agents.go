package cadreen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type CreateAgentRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Type        string         `json:"type,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

type UpdateAgentRequest struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

type ListAgentsParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

type SendAgentMessageRequest struct {
	FromAgentID string `json:"fromAgentId"`
	Content     string `json:"content"`
	Context     string `json:"context,omitempty"`
	ExecutionID string `json:"executionId,omitempty"`
}

type ListAgentMessagesParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type ListAgentExecutionsParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Status string `json:"status,omitempty"`
}

type CreateAgentExecutionRequest struct {
	Intent       string  `json:"intent"`
	Context      string  `json:"context,omitempty"`
	MaxBudgetUSD float64 `json:"maxBudgetUsd,omitempty"`
}

type ListAgentKnowledgeParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Type   string `json:"type,omitempty"`
}

type CreateAgentKnowledgeRequest struct {
	FactType    string   `json:"factType"`
	Subject     string   `json:"subject"`
	Predicate   string   `json:"predicate,omitempty"`
	Object      string   `json:"object,omitempty"`
	Source      string   `json:"source,omitempty"`
	Confidence  float64  `json:"confidence,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type SearchAgentKnowledgeRequest struct {
	Query      string `json:"query"`
	Type       string `json:"type,omitempty"`
	MaxResults int    `json:"max_results,omitempty"`
}

type CreateAgentGovernanceRequest struct {
	Name     string           `json:"name"`
	Rules    []map[string]any `json:"rules,omitempty"`
	Priority int              `json:"priority,omitempty"`
}

type UpdateAgentGovernanceRequest struct {
	Name     string           `json:"name,omitempty"`
	Rules    []map[string]any `json:"rules,omitempty"`
	Priority *int             `json:"priority,omitempty"`
	Active   *bool            `json:"active,omitempty"`
}

type ListAgentAuditParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Action string `json:"action,omitempty"`
}

type StartNegotiationRequest struct {
	Responder string         `json:"responder"`
	Topic     string         `json:"topic,omitempty"`
	Proposal  map[string]any `json:"proposal,omitempty"`
}

type ListNegotiationsParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Status string `json:"status,omitempty"`
}

type RespondToNegotiationRequest struct {
	Accepted bool           `json:"accepted"`
	Response map[string]any `json:"response,omitempty"`
	Reason   string         `json:"reason,omitempty"`
}

func (c *Client) CreateAgent(ctx context.Context, req CreateAgentRequest, opts ...RequestOption) (*Agent, error) {
	var result Agent
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create agent: %w", err)
	}
	return &result, nil
}

func (c *Client) ListAgents(ctx context.Context, params ListAgentsParams, opts ...RequestOption) (*ListAgentsResponse, error) {
	path := "/api/v1/cadreen/agents"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Type != "" {
		q.Set("type", params.Type)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var result ListAgentsResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	return &result, nil
}

func (c *Client) GetAgent(ctx context.Context, agentID string, opts ...RequestOption) (*Agent, error) {
	var result Agent
	if err := c.do(ctx, "GET", "/api/v1/cadreen/agents/"+agentID, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get agent: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateAgent(ctx context.Context, agentID string, req UpdateAgentRequest, opts ...RequestOption) (*Agent, error) {
	var result Agent
	if err := c.do(ctx, "PATCH", "/api/v1/cadreen/agents/"+agentID, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("update agent: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteAgent(ctx context.Context, agentID string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/agents/"+agentID, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete agent: %w", err)
	}
	return nil
}

func (c *Client) GetAgentConfig(ctx context.Context, agentID string, opts ...RequestOption) (*AgentConfig, error) {
	var result AgentConfig
	if err := c.do(ctx, "GET", "/api/v1/cadreen/agents/"+agentID+"/config", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get agent config: %w", err)
	}
	return &result, nil
}

type DeployAgentRequest struct {
	ConfigSnapshot json.RawMessage `json:"configSnapshot"`
	ChangeSummary  string          `json:"changeSummary,omitempty"`
}

func (c *Client) DeployAgent(ctx context.Context, agentID string, req DeployAgentRequest, opts ...RequestOption) (*Agent, error) {
	var result Agent
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/deploy", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("deploy agent: %w", err)
	}
	return &result, nil
}

func (c *Client) GetAgentCapabilities(ctx context.Context, agentID string, opts ...RequestOption) (*AgentCapabilitiesResponse, error) {
	var result AgentCapabilitiesResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/agents/"+agentID+"/capabilities", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get agent capabilities: %w", err)
	}
	return &result, nil
}

func (c *Client) SendAgentMessage(ctx context.Context, agentID string, req SendAgentMessageRequest, opts ...RequestOption) (*AgentMessage, error) {
	var result AgentMessage
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/send", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("send agent message: %w", err)
	}
	return &result, nil
}

func (c *Client) ListAgentMessages(ctx context.Context, agentID string, params ListAgentMessagesParams, opts ...RequestOption) (*ListAgentMessagesResponse, error) {
	path := "/api/v1/cadreen/agents/" + agentID + "/messages"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var result ListAgentMessagesResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list agent messages: %w", err)
	}
	return &result, nil
}

func (c *Client) ListAgentExecutions(ctx context.Context, agentID string, params ListAgentExecutionsParams, opts ...RequestOption) (*ListAgentExecutionsResponse, error) {
	path := "/api/v1/cadreen/agents/" + agentID + "/executions"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var result ListAgentExecutionsResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list agent executions: %w", err)
	}
	return &result, nil
}

func (c *Client) CreateAgentExecution(ctx context.Context, agentID string, req CreateAgentExecutionRequest, opts ...RequestOption) (*AgentExecution, error) {
	var result AgentExecution
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/executions", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create agent execution: %w", err)
	}
	return &result, nil
}

func (c *Client) ListAgentKnowledge(ctx context.Context, agentID string, params ListAgentKnowledgeParams, opts ...RequestOption) (*ListAgentKnowledgeResponse, error) {
	path := "/api/v1/cadreen/agents/" + agentID + "/knowledge"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Type != "" {
		q.Set("type", params.Type)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var result ListAgentKnowledgeResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list agent knowledge: %w", err)
	}
	return &result, nil
}

func (c *Client) CreateAgentKnowledge(ctx context.Context, agentID string, req CreateAgentKnowledgeRequest, opts ...RequestOption) (*AgentKnowledge, error) {
	var result AgentKnowledge
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/knowledge", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create agent knowledge: %w", err)
	}
	return &result, nil
}

func (c *Client) SearchAgentKnowledge(ctx context.Context, agentID string, req SearchAgentKnowledgeRequest, opts ...RequestOption) (*SearchAgentKnowledgeResponse, error) {
	var result SearchAgentKnowledgeResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/knowledge/search", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("search agent knowledge: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteAgentKnowledge(ctx context.Context, agentID, knowledgeID string, opts ...RequestOption) error {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/knowledge/%s", agentID, knowledgeID)
	if err := c.do(ctx, "DELETE", path, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete agent knowledge: %w", err)
	}
	return nil
}

func (c *Client) ListAgentGovernance(ctx context.Context, agentID string, opts ...RequestOption) (*ListAgentGovernanceResponse, error) {
	var result ListAgentGovernanceResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/agents/"+agentID+"/governance", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list agent governance: %w", err)
	}
	return &result, nil
}

func (c *Client) CreateAgentGovernance(ctx context.Context, agentID string, req CreateAgentGovernanceRequest, opts ...RequestOption) (*AgentGovernancePolicy, error) {
	var result AgentGovernancePolicy
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/governance", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create agent governance: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateAgentGovernance(ctx context.Context, agentID, policyID string, req UpdateAgentGovernanceRequest, opts ...RequestOption) (*AgentGovernancePolicy, error) {
	var result AgentGovernancePolicy
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/governance/%s", agentID, policyID)
	if err := c.do(ctx, "PATCH", path, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("update agent governance: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteAgentGovernance(ctx context.Context, agentID, policyID string, opts ...RequestOption) error {
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/governance/%s", agentID, policyID)
	if err := c.do(ctx, "DELETE", path, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete agent governance: %w", err)
	}
	return nil
}

func (c *Client) ListAgentAudit(ctx context.Context, agentID string, params ListAgentAuditParams, opts ...RequestOption) (*ListAgentAuditResponse, error) {
	path := "/api/v1/cadreen/agents/" + agentID + "/audit"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Action != "" {
		q.Set("action", params.Action)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var result ListAgentAuditResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list agent audit: %w", err)
	}
	return &result, nil
}

func (c *Client) StartNegotiation(ctx context.Context, agentID string, req StartNegotiationRequest, opts ...RequestOption) (*AgentNegotiation, error) {
	var result AgentNegotiation
	if err := c.do(ctx, "POST", "/api/v1/cadreen/agents/"+agentID+"/negotiate", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("start negotiation: %w", err)
	}
	return &result, nil
}

func (c *Client) ListNegotiations(ctx context.Context, agentID string, params ListNegotiationsParams, opts ...RequestOption) (*ListNegotiationsResponse, error) {
	path := "/api/v1/cadreen/agents/" + agentID + "/negotiations"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var result ListNegotiationsResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list negotiations: %w", err)
	}
	return &result, nil
}

func (c *Client) GetNegotiation(ctx context.Context, agentID, negotiationID string, opts ...RequestOption) (*AgentNegotiation, error) {
	var result AgentNegotiation
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/negotiations/%s", agentID, negotiationID)
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get negotiation: %w", err)
	}
	return &result, nil
}

func (c *Client) RespondToNegotiation(ctx context.Context, agentID, negotiationID string, req RespondToNegotiationRequest, opts ...RequestOption) (*AgentNegotiation, error) {
	var result AgentNegotiation
	path := fmt.Sprintf("/api/v1/cadreen/agents/%s/negotiations/%s/respond", agentID, negotiationID)
	if err := c.do(ctx, "POST", path, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("respond to negotiation: %w", err)
	}
	return &result, nil
}
