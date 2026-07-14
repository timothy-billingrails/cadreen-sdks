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

// CreateResponse sends a non-streaming Responses API request.
// Endpoint: POST /api/v1/cadreen/responses
func (c *Client) CreateResponse(ctx context.Context, req ResponseRequest, opts ...RequestOption) (*Response, error) {
	req.Stream = false
	var resp Response
	if err := c.do(ctx, "POST", "/api/v1/cadreen/responses", req, &resp, opts...); err != nil {
		return nil, fmt.Errorf("create response: %w", err)
	}
	return &resp, nil
}

// GetResponse retrieves a previously created response by ID.
// Endpoint: GET /api/v1/cadreen/responses/{responseID}
func (c *Client) GetResponse(ctx context.Context, responseID string, opts ...RequestOption) (*Response, error) {
	var resp Response
	path := "/api/v1/cadreen/responses/" + responseID
	if err := c.do(ctx, "GET", path, nil, &resp, opts...); err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}
	return &resp, nil
}

// ResponseStreamIterator iterates over SSE events from a streaming Responses API call.
type ResponseStreamIterator struct {
	ch      <-chan ResponseStreamEvent
	current ResponseStreamEvent
	err     error
	closed  bool
}

// CreateResponseStreaming sends a streaming Responses API request and returns an iterator.
// Endpoint: POST /api/v1/cadreen/responses (with stream: true)
func (c *Client) CreateResponseStreaming(ctx context.Context, req ResponseRequest, opts ...RequestOption) (*ResponseStreamIterator, error) {
	if c.config.Sandbox {
		return nil, fmt.Errorf("streaming is not available in sandbox mode")
	}
	req.Stream = true

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal response request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/api/v1/cadreen/responses", bytes.NewReader(bodyBytes))
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

	ch := make(chan ResponseStreamEvent, 16)
	go readResponseSSEStream(ctx, resp, ch)

	return &ResponseStreamIterator{ch: ch}, nil
}

// Next advances the iterator to the next event. Returns false when the stream is exhausted.
func (it *ResponseStreamIterator) Next() bool {
	if it.closed {
		return false
	}
	event, ok := <-it.ch
	if !ok {
		it.closed = true
		return false
	}
	it.current = event
	return true
}

// Current returns the most recent event from the stream.
func (it *ResponseStreamIterator) Current() ResponseStreamEvent {
	return it.current
}

// Err returns the first error encountered during streaming, if any.
func (it *ResponseStreamIterator) Err() error {
	return it.err
}

func readResponseSSEStream(ctx context.Context, resp *http.Response, ch chan<- ResponseStreamEvent) {
	defer close(ch)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var currentEvent string
	var dataBuffer strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if dataBuffer.Len() > 0 {
				sendResponseSSEToChan(ctx, ch, currentEvent, dataBuffer.String())
			}
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			if dataBuffer.Len() > 0 {
				sendResponseSSEToChan(ctx, ch, currentEvent, dataBuffer.String())
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

func sendResponseSSEToChan(ctx context.Context, ch chan<- ResponseStreamEvent, eventType string, data string) {
	var event ResponseStreamEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		event = ResponseStreamEvent{Type: eventType}
	}
	if event.Type == "" {
		event.Type = eventType
	}

	select {
	case ch <- event:
	case <-ctx.Done():
		return
	}
}
