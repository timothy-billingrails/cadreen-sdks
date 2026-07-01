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

// ── Chat Completions Types (OpenAI-compatible) ──

// ChatMessage represents a message in a chat conversation.
type ChatMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	Reasoning  string           `json:"reasoning,omitempty"` // Model reasoning (thinking models: DeepSeek, MiMo, Anthropic)
	Name       string           `json:"name,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"` // for "tool" role
	ToolCalls  []ChatToolCall   `json:"tool_calls,omitempty"`   // for "assistant" role
}

// ChatToolCall represents a tool call proposed by the model.
type ChatToolCall struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"` // "function"
	Function ChatFunctionCall   `json:"function"`
}

// ChatFunctionCall represents a function call with name and arguments.
type ChatFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// ChatToolDefinition describes a tool the model can call.
type ChatToolDefinition struct {
	Type     string               `json:"type"` // "function"
	Function ChatFunctionDefinition `json:"function"`
}

// ChatFunctionDefinition describes a callable function.
type ChatFunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters"` // JSON Schema
}

// ChatCompletionRequest is the request for POST /api/v1/cadreen/chat/completions.
type ChatCompletionRequest struct {
	Model          string               `json:"model,omitempty"`
	Messages       []ChatMessage        `json:"messages"`
	Stream         bool                 `json:"stream,omitempty"`
	Tools          []ChatToolDefinition `json:"tools,omitempty"`
	Context        map[string]any       `json:"context,omitempty"`
	ConversationID string               `json:"conversation_id,omitempty"`
	UserID         string               `json:"user_id,omitempty"`
	MaxTokens      int                  `json:"max_tokens,omitempty"`
}

// ChatCompletionResponse is the response from POST /api/v1/cadreen/chat/completions.
type ChatCompletionResponse struct {
	ID             string             `json:"id"`
	Object         string             `json:"object"`
	Created        int64              `json:"created"`
	Model          string             `json:"model"`
	Choices        []ChatChoice       `json:"choices"`
	Usage          *ChatUsage         `json:"usage,omitempty"`
	Intelligence   *IntelligenceMeta  `json:"intelligence,omitempty"`
	ConversationID string             `json:"conversation_id,omitempty"`
}

// ChatChoice represents a single choice in the response.
type ChatChoice struct {
	Index        int              `json:"index"`
	Message      ChatMessage      `json:"message"`
	FinishReason string           `json:"finish_reason"`
}

// ChatUsage represents token usage.
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details,omitempty"`
	ReasoningTokens  int `json:"reasoning_tokens,omitempty"`
	CacheWriteTokens int `json:"cache_write_tokens,omitempty"`
}

// ChatCompletionChunk is a streaming chunk from POST /api/v1/cadreen/chat/completions.
type ChatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChatChunkChoice `json:"choices"`
	Usage   *ChatUsage    `json:"usage,omitempty"`
}

// ChatChunkChoice represents a single choice in a streaming chunk.
type ChatChunkChoice struct {
	Index        int          `json:"index"`
	Delta        ChatDelta    `json:"delta"`
	FinishReason *string      `json:"finish_reason"`
}

// ChatDelta represents the delta content in a streaming chunk.
type ChatDelta struct {
	Role      string         `json:"role,omitempty"`
	Content   string         `json:"content,omitempty"`
	Reasoning string         `json:"reasoning,omitempty"`
	ToolCalls []ChatToolCall `json:"tool_calls,omitempty"`
}

// ChatStreamEvent represents an event from the chat completion stream.
type ChatStreamEvent struct {
	Chunk    *ChatCompletionChunk
	Reasoning string // Non-empty when this is a reasoning_delta event
	RawJSON  []byte // Raw JSON before typed parsing; includes Cadreen-specific fields (pending_actions, conversation_id, intelligence)
	Error    error
}

// ── Tool Discovery Types ──

// ToolEntry represents a tool from GET /api/v1/cadreen/tools.
type ToolEntry struct {
	Type     string               `json:"type"`
	Function ChatFunctionDefinition `json:"function"`
}

// ListToolsResponse is the response from GET /api/v1/cadreen/tools.
type ListToolsResponse struct {
	Object string     `json:"object"`
	Data   []ToolEntry `json:"data"`
}

// ── Chat Completions Methods ──

// ChatCompletions sends a chat completion request and returns the response.
// This is the non-streaming version.
func (c *Client) ChatCompletions(ctx context.Context, req ChatCompletionRequest, opts ...RequestOption) (*ChatCompletionResponse, error) {
	req.Stream = false
	var resp ChatCompletionResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/chat/completions", req, &resp, opts...); err != nil {
		return nil, fmt.Errorf("chat completions: %w", err)
	}
	return &resp, nil
}

// ChatCompletionsStream sends a streaming chat completion request and returns
// a channel of chunks. The channel is closed when the stream ends.
func (c *Client) ChatCompletionsStream(ctx context.Context, req ChatCompletionRequest) (<-chan ChatStreamEvent, error) {
	if c.config.Sandbox {
		return nil, fmt.Errorf("streaming is not available in sandbox mode")
	}
	req.Stream = true

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal chat request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/api/v1/cadreen/chat/completions", bytes.NewReader(bodyBytes))
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

	ch := make(chan ChatStreamEvent, 16)
	go readChatSSEStream(ctx, resp, ch)
	return ch, nil
}

func readChatSSEStream(ctx context.Context, resp *http.Response, ch chan<- ChatStreamEvent) {
	defer close(ch)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	currentEvent := "message"

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				select {
				case ch <- ChatStreamEvent{Error: fmt.Errorf("read SSE: %w", err)}:
				case <-ctx.Done():
				}
			}
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return
		}

		if currentEvent == "reasoning_delta" {
			var rd struct {
				Reasoning string `json:"reasoning"`
			}
			if err := json.Unmarshal([]byte(data), &rd); err == nil {
				select {
				case ch <- ChatStreamEvent{Reasoning: rd.Reasoning, RawJSON: []byte(data)}:
				case <-ctx.Done():
					return
				}
			}
			currentEvent = "message"
			continue
		}

		var chunk ChatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			select {
			case ch <- ChatStreamEvent{RawJSON: []byte(data), Error: fmt.Errorf("parse SSE chunk: %w", err)}:
			case <-ctx.Done():
			}
			return
		}

		currentEvent = "message"
		select {
		case ch <- ChatStreamEvent{Chunk: &chunk, RawJSON: []byte(data)}:
		case <-ctx.Done():
			return
		}
	}
}

// ── Tool Discovery Methods ──

// ListTools returns the available tools as OpenAI-compatible function definitions.
func (c *Client) ListTools(ctx context.Context, opts ...RequestOption) (*ListToolsResponse, error) {
	var resp ListToolsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/tools", nil, &resp, opts...); err != nil {
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &resp, nil
}
