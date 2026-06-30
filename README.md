# Cadreen SDKs

Official SDKs for the [Cadreen API](https://accomplishanything.today/infra/docs) — Intelligence as a Service.


## Packages

| Language | Package | Version | Install |
|----------|---------|---------|---------|
| TypeScript | `@cadreen/sdk` | 0.6.1 | `npm install @cadreen/sdk` |
| Python | `cadreen-sdk` | 0.6.1 | `pip install cadreen-sdk` |
| Go | `cadreen` | 0.6.1 | `go get github.com/timothy-billingrails/cadreen-sdks/go/cadreen@latest` |

## Quick Start

### Intent — the primary door

```ts
// TypeScript
import { Cadreen } from "@cadreen/sdk";

const cadreen = new Cadreen({ apiKey: process.env.CADREEN_API_KEY });

const result = await cadreen.intent.invoke({
  messages: [{ role: "user", content: "Handle refund for invoice inv_123" }],
});

switch (result.type) {
  case "direct":
    console.log(result.message.content);
    break;
  case "clarify":
    for (const q of result.questions) console.log(q);
    break;
  case "execution":
    for await (const event of cadreen.executions.stream(result.execution.id)) {
      console.log(event.type, event.data);
    }
    break;
}
```

```python
# Python
import asyncio
from cadreen import Cadreen

cadreen = Cadreen(api_key="...")

async def main():
    result = await cadreen.intent("Handle refund for invoice inv_123")
    print(result.explain())

asyncio.run(main())
```

```go
// Go
c := cadreen.NewClient(cadreen.CadreenConfig{
    APIKey: os.Getenv("CADREEN_API_KEY"),
})

result, err := c.IntentInvoke(context.Background(), cadreen.IntentRequest{
    Messages: []cadreen.IntentMessage{
        {Role: "user", Content: "Handle refund for invoice inv_123"},
    },
})
```

### Chat Completions — OpenAI-compatible with governance

```ts
// TypeScript
const response = await cadreen.chat.completions({
  messages: [{ role: "user", content: "Summarize my recent orders" }],
  tools: [{
    type: "function",
    function: {
      name: "get_orders",
      description: "Fetch recent orders",
      parameters: { type: "object", properties: { limit: { type: "number" } } },
    },
  }],
});
console.log(response.choices[0].message);
```

```python
# Python
from cadreen.resources.chat import ChatCompletionRequest, ChatMessage, ChatToolDefinition, ChatFunctionDefinition

response = await cadreen.chat.completions(ChatCompletionRequest(
    messages=[ChatMessage(role="user", content="Summarize my recent orders")],
    tools=[ChatToolDefinition(function=ChatFunctionDefinition(
        name="get_orders",
        description="Fetch recent orders",
        parameters={"type": "object", "properties": {"limit": {"type": "number"}}},
    ))],
))
print(response.choices[0].message)
```

```go
// Go
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
    Messages: []cadreen.ChatMessage{
        {Role: "user", Content: "Summarize my recent orders"},
    },
    Tools: []cadreen.ChatToolDefinition{{
        Type: "function",
        Function: cadreen.ChatFunctionDefinition{
            Name:       "get_orders",
            Description: "Fetch recent orders",
            Parameters: map[string]any{"type": "object", "properties": map[string]any{"limit": map[string]any{"type": "number"}}},
        },
    }},
})
```

### Tool Discovery

```ts
// TypeScript
const tools = await cadreen.chat.listTools();
```

```python
# Python
tools = await cadreen.chat.list_tools()
```

```go
// Go
tools, err := c.ListTools(ctx)
```

## Marketplace

All three SDKs support the unified integration marketplace:

```ts
// TypeScript
const catalog = await cadreen.connections.catalog();
const install = await cadreen.connections.install("slack");
```

```python
# Python
catalog = await cadreen.connections.catalog()
install = await cadreen.connections.install("slack")
```

```go
// Go
catalog, _ := c.Catalog(ctx)
install, _ := c.Install(ctx, "slack")
```

## Documentation

- [API Docs](https://accomplishanything.today/infra/docs)
- [OpenAPI Spec](https://accomplishanything.today/api/v1/cadreen/docs/openapi.json)
- [Cadreen API Reference](https://accomplishanything.today/cadreen-api.md)

## License

MIT
