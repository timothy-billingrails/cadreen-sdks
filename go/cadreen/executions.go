package cadreen

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ExecutionEvent represents a single SSE event from an execution stream.
type ExecutionEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// ExecutionStatus represents the current status of an execution.
type ExecutionStatus struct {
	ID       string                 `json:"id"`
	Status   string                 `json:"status"`
	Progress *float64               `json:"progress,omitempty"`
	Result   map[string]interface{} `json:"result,omitempty"`
	Error    string                 `json:"error,omitempty"`
}

// ExecutionStream streams SSE events from an execution.
func (c *Client) ExecutionStream(ctx context.Context, executionID string) (<-chan ExecutionEvent, error) {
	if c.config.Sandbox {
		return nil, fmt.Errorf("streaming is not available in sandbox mode")
	}

	url := fmt.Sprintf("%s/api/v1/cadreen/executions/%s/stream", c.config.BaseURL, executionID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create execution stream request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execution stream request: %w", err)
	}

	c.setLastTrace(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, parseError(resp.StatusCode, respBody)
	}

	ch := make(chan ExecutionEvent, 16)
	go readExecutionSSEStream(ctx, resp, ch)
	return ch, nil
}

// GetExecutionStatus returns the current status of an execution.
func (c *Client) GetExecutionStatus(ctx context.Context, executionID string) (*ExecutionStatus, error) {
	var resp ExecutionStatus
	if err := c.do(ctx, "GET", fmt.Sprintf("/api/v1/cadreen/executions/%s", executionID), nil, &resp); err != nil {
		return nil, fmt.Errorf("get execution status: %w", err)
	}
	return &resp, nil
}

func readExecutionSSEStream(ctx context.Context, resp *http.Response, ch chan<- ExecutionEvent) {
	defer close(ch)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var currentEvent string
	var dataBuffer strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if dataBuffer.Len() > 0 {
				sendExecutionSSEToChan(ctx, ch, currentEvent, dataBuffer.String())
			}
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			if dataBuffer.Len() > 0 {
				sendExecutionSSEToChan(ctx, ch, currentEvent, dataBuffer.String())
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

func sendExecutionSSEToChan(ctx context.Context, ch chan<- ExecutionEvent, eventType string, data string) {
	if eventType == "" {
		eventType = "message"
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(data), &parsed); err != nil {
		parsed = map[string]interface{}{"raw": data}
	}

	select {
	case ch <- ExecutionEvent{Type: eventType, Data: parsed}:
	case <-ctx.Done():
		return
	}
}
