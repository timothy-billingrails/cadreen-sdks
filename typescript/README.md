# @cadreen/sdk

TypeScript SDK for Cadreen — Intelligence as a Service.

## Install

```bash
npm install @cadreen/sdk
```

## Quick Start

```ts
import { Cadreen } from "@cadreen/sdk";

const cadreen = new Cadreen({ apiKey: process.env.CADREEN_API_KEY });

// Intent — the primary door
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
  case "blocked":
    console.log(result.policy.reason);
    break;
  case "connect_required":
    console.log(result.connection.endpoint);
    break;
}

// Memory — remember things
await cadreen.memory.remember({
  type: "reference",
  content: { text: "GDPR Article 17...", title: "GDPR Art. 17" },
  authority: 10,
});

// Policies — require approval
await cadreen.policies.requireApproval("Refunds over $500 require human approval");

// Traces — inspect what happened
const trace = await cadreen.traces.get(result.traceId);
console.log(trace.explain());

// Connections — connect tools
await cadreen.connections.registerOpenAPI({
  name: "internal-erp",
  specUrl: "https://erp.example.com/openapi.json",
});
```

## Configuration

```ts
const cadreen = new Cadreen({
  apiKey: "cadreen_...",                    // required
  baseUrl: "https://accomplishanything.today", // optional, default shown
  maxRetries: 2,                             // optional, default 2
  timeout: 30000,                            // optional, default 30s
  profile: "lean",                           // optional: "lean" | "audit" | "full" (default "full")
});
```

### Response Profiles

Control how much intelligence metadata you get back:

| Profile | What you get | Use when |
|---------|-------------|----------|
| `"full"` (default) | Full intelligence envelope with capability, reasoning, memory, governance, humility, process | You want full transparency |
| `"audit"` | Only governance decision + confidence + blocking gaps | You need to react to gates, not inspect internals |
| `"lean"` | No envelope. `trace_id` in body + `X-Cadreen-Trace-ID` header | Hot-looping, minimal payload |

```ts
// Lean: skip the envelope entirely
const lean = new Cadreen({ apiKey: "...", profile: "lean" });

// Audit: just the action-bearing fields
const audit = new Cadreen({ apiKey: "...", profile: "audit" });
```

## Error Handling

```ts
import { CadreenError } from "@cadreen/sdk";

try {
  const result = await cadreen.intent.invoke({ messages: [...] });
} catch (err) {
  if (err instanceof CadreenError) {
    console.log(err.status);       // HTTP status
    console.log(err.code);         // machine-readable code
    console.log(err.intelligence); // trace context when available
  }
}
```

## Resources

| Resource | Methods |
|----------|---------|
| `cadreen.intent` | `invoke(request)` |
| `cadreen.memory` | `remember(request)`, `search(request)`, `get(id)` |
| `cadreen.policies` | `create(request)`, `evaluate(request)`, `confirm(id)`, `list()`, `get(id)`, `requireApproval(desc)` |
| `cadreen.connections` | `registerOpenAPI(request)`, `registerMCP(request)`, `installComposio(request)`, `searchComposio(query)`, `composioStatus(...)`, `list()`, `delete(id)` |
| `cadreen.traces` | `get(id)`, `list(options?)`, `stats()` |
| `cadreen.executions` | `stream(id)`, `getStatus(id)` |

## Shorthand

`cadreen.invoke(request)` is an alias for `cadreen.intent.invoke(request)`.

## License

MIT
