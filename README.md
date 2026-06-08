# Cadreen SDKs

Official SDKs for the [Cadreen API](https://accomplishanything.today/infra/docs) — Intelligence as a Service.

## Packages

| Language | Package | Install |
|----------|---------|---------|
| TypeScript | `@cadreen/sdk` | `npm install @cadreen/sdk` |
| Python | `cadreen-sdk` | `pip install cadreen-sdk` |
| Scaffold | `create-cadreen-app` | `npx create-cadreen-app my-project` |

## Quick Start

### One command

```bash
npx create-cadreen-app my-project
cd my-project
# Set CADREEN_API_KEY in .env
npm run dev
```

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

## Documentation

- [API Docs](https://accomplishanything.today/infra/docs)
- [OpenAPI Spec](https://accomplishanything.today/api/v1/cadreen/docs/openapi.json)

## License

MIT
