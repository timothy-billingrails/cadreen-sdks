package cadreen

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) IntentInvoke(ctx context.Context, req IntentRequest, opts ...RequestOption) (*IntentResult, error) {
	var raw unifiedIntentResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/intent", req, &raw, opts...); err != nil {
		return nil, fmt.Errorf("intent invoke: %w", err)
	}

	traceID := c.LastTraceID()
	return mapIntentResponse(raw, traceID), nil
}

type unifiedIntentResponse struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Status        string                `json:"status,omitempty"`
	NextAction    *NextAction           `json:"next_action,omitempty"`
	TraceID       string                `json:"trace_id,omitempty"`
	Message       *ResponseMessage      `json:"message,omitempty"`
	Mission       *ResponseExecution    `json:"mission,omitempty"`
	Execution     *ResponseExecution    `json:"execution,omitempty"`
	Clarification *unifiedClarification `json:"clarification,omitempty"`
	Meta          *unifiedMeta          `json:"meta,omitempty"`
	Intelligence  *IntelligenceMeta     `json:"intelligence,omitempty"`
}

type unifiedClarification struct {
	Questions      ClarificationQuestions `json:"questions,omitempty"`
	ConversationID string                 `json:"conversation_id,omitempty"`
}

type unifiedMeta struct {
	Governance *unifiedGovernance `json:"governance,omitempty"`
}

type unifiedGovernance struct {
	Decision   string `json:"decision,omitempty"`
	Reason     string `json:"reason,omitempty"`
	ReasonCode string `json:"reason_code,omitempty"`
	PolicyID   string `json:"policy_id,omitempty"`
}

var defaultIntelligence = &IntelligenceMeta{
	Capability: CapabilityTrace{TotalAvailable: 0, HealthyCount: 0},
	Reasoning:  ReasoningTrace{},
	Memory:     MemoryTrace{Healthy: true},
	Governance: GovernanceTrace{Active: false},
	Humility:   HumilityTrace{},
	Process:    ProcessTrace{StartedAt: "", DurationMs: 0},
	FieldStability: FieldStability{
		Stable:   []string{},
		Evolving: []string{},
		Internal: []string{},
	},
}

func deriveStatus(rawType string) string {
	switch rawType {
	case "direct":
		return "answered"
	case "clarify":
		return "needs_input"
	case "mission":
		return "ready_to_execute"
	case "execution":
		return "ready_to_execute"
	case "blocked":
		return "blocked_by_policy"
	case "connect_required":
		return "blocked_by_missing_connection"
	default:
		return "failed"
	}
}

func mapIntentResponse(raw unifiedIntentResponse, traceID string) *IntentResult {
	intelligence := raw.Intelligence
	if intelligence == nil {
		intelligence = defaultIntelligence
	}

	if traceID == "" {
		traceID = raw.TraceID
		if traceID == "" {
			traceID = raw.ID
		}
	}

	status := raw.Status
	if status == "" {
		status = deriveStatus(raw.Type)
	}

	nextAction := raw.NextAction
	if nextAction == nil {
		nextAction = intelligence.NextAction
	}

	result := &IntentResult{
		Status:       status,
		NextAction:   nextAction,
		TraceID:      traceID,
		Intelligence: intelligence,
	}

	switch raw.Type {
	case "direct":
		result.Type = IntentResultDirect
		msg := raw.Message
		if msg == nil {
			msg = &ResponseMessage{Role: "assistant", Content: ""}
		}
		result.Message = msg

	case "clarify":
		result.Type = IntentResultClarify
		if raw.Clarification != nil {
			result.Questions = raw.Clarification.Questions
			result.ConversationID = raw.Clarification.ConversationID
		} else {
			result.Questions = []ResponseClarificationQuestion{}
			result.ConversationID = ""
		}

	case "execution":
		result.Type = IntentResultExecution
		if raw.Execution != nil {
			result.Execution = raw.Execution
		} else if raw.Mission != nil {
			result.Execution = raw.Mission
		} else {
			result.Execution = &ResponseExecution{}
		}

	case "mission":
		result.Type = IntentResultExecution
		if raw.Mission != nil {
			result.Execution = raw.Mission
		} else {
			result.Execution = &ResponseExecution{}
		}

	case "blocked":
		result.Type = IntentResultBlocked
		if raw.Meta != nil && raw.Meta.Governance != nil {
			result.ReasonCode = raw.Meta.Governance.ReasonCode
			result.PolicyID = raw.Meta.Governance.PolicyID
		}

	case "connect_required":
		result.Type = IntentResultConnectRequired
		if result.NextAction == nil {
			result.NextAction = &NextAction{
				Endpoint: "/api/v1/cadreen/connections",
				Reason:   "connection required",
			}
		}

	default:
		result.Type = IntentResultDirect
		msg := raw.Message
		if msg == nil {
			msg = &ResponseMessage{Role: "assistant", Content: ""}
		}
		result.Message = msg
	}

	return result
}

type IntentEvent struct {
	Type string
	Data map[string]any
}

func (c *Client) IntentStream(ctx context.Context, req IntentRequest) (<-chan IntentEvent, error) {
	if c.config.Sandbox {
		return nil, fmt.Errorf("streaming is not available in sandbox mode")
	}
	req.Stream = true

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal intent request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/api/v1/cadreen/intent", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create stream request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("stream request: %w", err)
	}

	c.setLastTrace(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, parseError(resp.StatusCode, respBody)
	}

	ch := make(chan IntentEvent, 16)
	go readSSEStream(ctx, resp, ch)
	return ch, nil
}

func readSSEStream(ctx context.Context, resp *http.Response, ch chan<- IntentEvent) {
	defer close(ch)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var currentEvent string
	var dataBuffer strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if dataBuffer.Len() > 0 {
				sendSSEToChan(ctx, ch, currentEvent, dataBuffer.String())
			}
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			if dataBuffer.Len() > 0 {
				sendSSEToChan(ctx, ch, currentEvent, dataBuffer.String())
				dataBuffer.Reset()
				currentEvent = "message"
			}
			continue
		}

		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if dataBuffer.Len() > 0 {
				dataBuffer.WriteString("\n")
			}
			dataBuffer.WriteString(data)
			continue
		}
	}
}

func sendSSEToChan(ctx context.Context, ch chan<- IntentEvent, eventType string, data string) {
	if eventType == "" {
		eventType = "message"
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(data), &parsed); err != nil {
		parsed = map[string]any{"raw": data}
	}

	select {
	case ch <- IntentEvent{Type: eventType, Data: parsed}:
	case <-ctx.Done():
		return
	}
}
