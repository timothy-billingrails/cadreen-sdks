# cadreen

Go SDK for [Cadreen](https://accomplishanything.today/infra/docs) — Intelligence as a Service.

Cadreen is a cognitive operating system. Send messages describing what you want done, and Cadreen reasons, connects tools, recalls knowledge, governs actions, and escalates to humans when needed. The SDK handles authentication, retries, idempotency, streaming, sandbox mode, and error classification.

## Install

```bash
go get github.com/timothy-billingrails/cadreen-sdks/go/cadreen@latest
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"
)

func main() {
	c := cadreen.NewClient(cadreen.CadreenConfig{
		APIKey: os.Getenv("CADREEN_API_KEY"),
	})

	result, err := c.IntentInvoke(context.Background(), cadreen.IntentRequest{
		Messages: []cadreen.IntentMessage{
			{Role: "user", Content: "Handle refund for invoice inv_123"},
		},
	})
	if err != nil {
		panic(err)
	}

	switch result.Type {
	case cadreen.IntentResultDirect:
		fmt.Println(result.Message.Content)
	case cadreen.IntentResultExecution:
		stream, _ := c.ExecutionStream(context.Background(), result.Execution.ID)
		for event := range stream {
			fmt.Println(event.Type, event.Data)
		}
	case cadreen.IntentResultBlocked:
		fmt.Printf("Blocked by %s: %s\n", result.PolicyID, result.ReasonCode)
	case cadreen.IntentResultClarify:
		for _, q := range result.Questions {
			fmt.Println(q.Question)
		}
	case cadreen.IntentResultConnectRequired:
		fmt.Println("Connect:", result.Connection.Endpoint)
	}
}
```

## Chat Completions

OpenAI-compatible chat completions with built-in governance. Every tool call goes through governance before execution — auto-approved calls execute silently, blocked calls become conversation.

```go
// Basic completion
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "Hello!"},
	},
})
fmt.Println(resp.Choices[0].Message.Content)

// With tool calling
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "Refund order 456"},
	},
	Tools: []cadreen.ChatToolDefinition{{
		Type: "function",
		Function: cadreen.ChatFunctionDefinition{
			Name:        "process_refund",
			Description: "Process a refund for an order",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]any{"order_id": map[string]any{"type": "string"}},
				"required":   []string{"order_id"},
			},
		},
	}},
})

// Check for tool calls
if len(resp.Choices[0].Message.ToolCalls) > 0 {
	for _, tc := range resp.Choices[0].Message.ToolCalls {
		fmt.Printf("%s(%s)\n", tc.Function.Name, tc.Function.Arguments)
	}
}

// Resume a conversation
followUp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "What about order 789?"},
	},
	ConversationID: resp.ID, // use conversation_id from prior response
})
```

### Streaming

```go
ch, err := c.ChatCompletionsStream(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "Hello!"},
	},
})
for event := range ch {
	if event.Error != nil {
		log.Fatal(event.Error)
	}
	if len(event.Chunk.Choices) > 0 {
		fmt.Print(event.Chunk.Choices[0].Delta.Content)
	}
}
```

### Tool Discovery

```go
tools, err := c.ListTools(ctx)
for _, tool := range tools.Data {
	fmt.Printf("%s: %s\n", tool.Function.Name, tool.Function.Description)
}
```

### Tool Chaining

When the model proposes tool calls, send results back for follow-up:

```go
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "What's the weather in NYC?"},
		{Role: "assistant", ToolCalls: []cadreen.ChatToolCall{{
			ID:       "tc_1",
			Type:     "function",
			Function: cadreen.ChatFunctionCall{Name: "get_weather", Arguments: `{"city":"NYC"}`},
		}}},
		{Role: "tool", ToolCallID: "tc_1", Content: `{"temp": 72, "condition": "sunny"}`},
	},
})
// Model may propose more tools or return a final text response
```

## Configuration

```go
c := cadreen.NewClient(cadreen.CadreenConfig{
	APIKey:     "sk_cadreen_...",                    // required
	BaseURL:    "https://accomplishanything.today",  // optional, default shown
	MaxRetries: 3,                                   // optional, default 3
	Timeout:    30 * time.Second,                    // optional, default 30s
	Profile:    "lean",                              // optional: "lean" | "audit" | "full" (default "full")
})
```

### Response Profiles

| Profile | What you get | Use when |
|---------|-------------|----------|
| `"full"` (default) | Full response metadata | You want full transparency |
| `"audit"` | Only governance decision + confidence + blocking gaps | You need to react to gates |
| `"lean"` | No envelope. Just `trace_id` | Hot-looping, minimal payload |

### Sandbox Mode

Test without hitting the API:

```go
c := cadreen.NewClient(cadreen.CadreenConfig{
	Sandbox: true,
	Fixtures: map[string]any{
		"POST /api/v1/cadreen/intent": cadreen.IntentResult{
			Type:    cadreen.IntentResultDirect,
			Status:  "answered",
			TraceID: "sandbox-trace",
			Message: &cadreen.ResponseMessage{Role: "assistant", Content: "It's done."},
		},
	},
})
```

## Marketplace

```go
catalog, _ := c.Catalog(ctx)
for _, category := range catalog.Categories {
	fmt.Printf("%s: %d integrations\n", category.Name, len(category.Integrations))
}

install, _ := c.Install(ctx, "slack")
if install.Status == "pending_auth" {
	fmt.Printf("Authenticate at: %s\n", install.AuthURL)
}
```

## Resources

| Resource | Methods |
|----------|---------|
| `Client.IntentInvoke()` | Intent — primary entry point |
| `Client.ChatCompletions()` | Non-streaming chat completions |
| `Client.ChatCompletionsStream()` | Streaming chat completions (SSE) |
| `Client.ListTools()` | Tool discovery |
| `Client.ExecutionStream()` | Stream execution events |
| `Client.GetExecutionStatus()` | Get execution status |
| `Client.Setup()` | One-time setup |
| `Client.Catalog()` | Browse marketplace |
| `Client.Install()` | Install integration |

## Changelog

### v0.5.0 (BREAKING)
- **BREAKING:** API endpoints moved to Cadreen surface (`/api/v1/cadreen/`):
  - `/api/v1/chat/completions` → `/api/v1/cadreen/chat/completions`
  - `/api/v1/tools` → `/api/v1/cadreen/tools`
- All external API calls now route through the Cadreen surface
- Removed "response metadata" terminology (was "intelligence envelope")

### v0.4.1
- Added `ChatStreamEvent.RawJSON` — raw JSON bytes before typed parsing; enables extracting `pending_actions`, `conversation_id`, and `intelligence` fields without SDK type changes
- Fixed unused `bytes` import in `executions.go`

### v0.4.0
- Added `ChatCompletions()` — OpenAI-compatible chat completions with governance
- Added `ChatCompletionsStream()` — streaming chat completions via SSE
- Added `ListTools()` — discover available tools as OpenAI function definitions
- Added tool calling support: `Tools` param, `ToolCalls` in responses, tool chaining
- Added `ConversationID` for persistent conversations across requests
- Added `ChatMessage`, `ChatToolDefinition`, `ChatToolCall` and related types
- Added `ExecutionStream()` and `GetExecutionStatus()` for execution monitoring

### v0.3.0
- Added `Catalog()` — browse the unified integration marketplace
- Added `Install()` — one-click install with OAuth flow

### v0.2.0
- Initial public release

## License

MIT
