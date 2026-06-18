# Cadreen SDKs

Official SDKs for the [Cadreen API](https://accomplishanything.today/infra/docs) — Intelligence as a Service.

## Packages

| Language | Package | Version | Install |
|----------|---------|---------|---------|
| TypeScript | `@cadreen/sdk` | 0.4.0 | `npm install @cadreen/sdk` |
| Python | `cadreen-sdk` | 0.4.0 | `pip install cadreen-sdk` |
| Go | `cadreen` | 0.4.0 | `go get github.com/timothy-billingrails/cadreen-sdks/go/cadreen@latest` |

## Quick Start

### TypeScript

```ts
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

### Python

```python
import asyncio
from cadreen import Cadreen

cadreen = Cadreen(api_key="...")

async def main():
    result = await cadreen.intent("Handle refund for invoice inv_123")
    print(result.explain())

asyncio.run(main())
```

### Go

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
		fmt.Println("Execution started:", result.Execution.ID)
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
