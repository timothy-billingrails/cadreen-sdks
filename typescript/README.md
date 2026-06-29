# @cadreen/sdk

TypeScript SDK for [Cadreen](https://accomplishanything.today/infra/docs) — Intelligence as a Service.

Cadreen is a cognitive operating system. Send messages describing what you want done, and Cadreen reasons, connects tools, recalls knowledge, governs actions, and escalates to humans when needed. The SDK handles authentication, retries, idempotency, streaming, and error classification.

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
    console.log(result.reason_code);
    break;
  case "connect_required":
    console.log(result.endpoint);
    break;
}
```

## Chat Completions

OpenAI-compatible chat completions with built-in governance. Every tool call goes through governance before execution — auto-approved calls execute silently, blocked calls become conversation.

```ts
// Basic completion
const response = await cadreen.chat.completions({
  messages: [{ role: "user", content: "Hello!" }],
});
console.log(response.choices[0].message.content);

// With tool calling
const response = await cadreen.chat.completions({
  messages: [{ role: "user", content: "Refund order 456" }],
  tools: [{
    type: "function",
    function: {
      name: "process_refund",
      description: "Process a refund for an order",
      parameters: {
        type: "object",
        properties: { order_id: { type: "string" } },
        required: ["order_id"],
      },
    },
  }],
});

// If governance needs approval, the response contains a text message
// asking the user to confirm — no tool_calls field
if (response.choices[0].message.tool_calls) {
  // Tool calls were auto-approved, handle them
  for (const tc of response.choices[0].message.tool_calls) {
    console.log(`${tc.function.name}(${tc.function.arguments})`);
  }
}

// Resume a conversation
const followUp = await cadreen.chat.completions({
  messages: [{ role: "user", content: "What about order 789?" }],
  conversation_id: response.conversation_id,
});
```

### Streaming

```ts
const stream = await cadreen.chat.completionsStream({
  messages: [{ role: "user", content: "Hello!" }],
});

for await (const event of stream) {
  if (event.type === "chunk") {
    process.stdout.write(event.chunk.choices[0]?.delta?.content || "");
  }
}
```

### Tool Discovery

```ts
const tools = await cadreen.chat.listTools();
for (const tool of tools.data) {
  console.log(`${tool.function.name}: ${tool.function.description}`);
}
```

### Tool Chaining

When the model proposes tool calls, send results back for follow-up:

```ts
const response = await cadreen.chat.completions({
  messages: [
    { role: "user", content: "What's the weather in NYC?" },
    { role: "assistant", content: null, tool_calls: [{ id: "tc_1", type: "function", function: { name: "get_weather", arguments: '{"city":"NYC"}' } }] },
    { role: "tool", tool_call_id: "tc_1", content: '{"temp": 72, "condition": "sunny"}' },
  ],
});
// Model may propose more tools or return a final text response
```

## Configuration

```ts
const cadreen = new Cadreen({
  apiKey: "sk_cadreen_...",                       // required
  baseUrl: "https://accomplishanything.today",    // optional, default shown
  maxRetries: 2,                                  // optional, default 2
  timeout: 30000,                                 // optional, default 30s
  profile: "lean",                                // optional: "lean" | "audit" | "full" (default "full")
});
```

### Response Profiles

Control how much intelligence metadata you get back:

| Profile | What you get | Use when |
|---------|-------------|----------|
| `"full"` (default) | Full response metadata with capability, reasoning, memory, governance, humility, process | You want full transparency |
| `"audit"` | Only governance decision + confidence + blocking gaps | You need to react to gates, not inspect internals |
| `"lean"` | No envelope. `trace_id` in body + `X-Cadreen-Trace-ID` header | Hot-looping, minimal payload |

```ts
const lean = new Cadreen({ apiKey: "...", profile: "lean" });
const audit = new Cadreen({ apiKey: "...", profile: "audit" });
```

## Marketplace

Browse and install integrations without knowing which provider powers them:

```ts
// Browse available integrations
const catalog = await cadreen.connections.catalog();
for (const category of catalog.categories) {
  console.log(`${category.name}: ${category.integrations.length} integrations`);
}

// One-click install (returns OAuth URL)
const install = await cadreen.connections.install("slack");
if (install.status === "pending_auth") {
  console.log(`Authenticate at: ${install.auth_url}`);
}

// Check what's installed
console.log(catalog.installed); // ["stripe", "github"]
```

## Documents

```ts
// List documents
const docs = await cadreen.documents.list();
for (const doc of docs.documents) {
  console.log(`${doc.name} (${doc.content_type})`);
}

// Upload a document (File, Blob, or ReadableStream)
const file = new File(["content"], "report.pdf", { type: "application/pdf" });
const result = await cadreen.documents.upload(file);
console.log(`Uploaded: ${result.name} (ID: ${result.id})`);

// Get document details
const doc = await cadreen.documents.get(result.id);
console.log(`Status: ${doc.status}, Size: ${doc.size} bytes`);
```

## Memory

Store and retrieve knowledge:

```ts
// Remember something
await cadreen.memory.remember({
  type: "reference",
  content: { text: "GDPR Article 17: Right to erasure", title: "GDPR Art. 17" },
  authority: 10,
});

// Search knowledge
const results = await cadreen.memory.search({ query: "data deletion rules" });

// Get by ID
const item = await cadreen.memory.get("mem_abc123");
```

## Policies

Set governance guardrails:

```ts
// Create a policy
await cadreen.policies.create({
  name: "refund_threshold",
  rules: [{ condition: "refund_amount > 500", action: "require_approval" }],
});

// Evaluate an action against policies
const evaluation = await cadreen.policies.evaluate({
  action: "Process $750 refund for order 456",
});
```

## Connections

Register external tools:

```ts
// Register from OpenAPI spec
await cadreen.connections.registerOpenAPI({
  name: "internal-erp",
  specUrl: "https://erp.example.com/openapi.json",
});

// Register MCP server
await cadreen.connections.registerMCP({
  name: "my-mcp-server",
  url: "https://mcp.example.com/sse",
  transport: "sse",
});

// List installed connections
const connections = await cadreen.connections.list();
```

## Traces

Inspect what happened:

```ts
const trace = await cadreen.traces.get(result.traceId);
console.log(trace.explain()); // human-readable summary

// List recent traces
const recent = await cadreen.traces.list({ limit: 10 });

// Aggregated stats
const stats = await cadreen.traces.stats();
```

## Proposals

Cadreen watches your usage and suggests improvements — actions to automate, schedules to set, rules to relax. You decide what runs.

```ts
// List proposals waiting for your decision
const { proposals } = await cadreen.proposals.list();
for (const p of proposals) {
  console.log(`[${p.status}] ${p.title} (${Math.round(p.confidence * 100)}% confidence)`);
}

// Get a specific proposal
const proposal = await cadreen.proposals.get("550e8400-...");

// Accept — executes via the intent engine
const result = await cadreen.proposals.accept("550e8400-...");
console.log(`Execution: ${result.execution_id}, Action: ${result.action}`);

// Dismiss — teaches Cadreen what you don't want
await cadreen.proposals.dismiss("550e8400-...", "We handle this manually");

// See counts by status
const stats = await cadreen.proposals.stats();
console.log(`Waiting: ${stats.proposed}, Accepted: ${stats.accepted}`);
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
| `cadreen.chat` | `completions(request)`, `completionsStream(request)`, `listTools()` |
| `cadreen.memory` | `remember(request)`, `search(request)`, `get(id)` |
| `cadreen.policies` | `create(request)`, `evaluate(request)`, `confirm(id)`, `list()`, `get(id)`, `requireApproval(desc)` |
| `cadreen.connections` | `catalog()`, `install(id)`, `registerOpenAPI(request)`, `registerMCP(request)`, `list()`, `delete(id)` |
| `cadreen.traces` | `get(id)`, `list(options?)`, `stats()` |
| `cadreen.executions` | `stream(id)`, `getStatus(id)` |
| `cadreen.documents` | `list()`, `get(id)`, `download(id)`, `upload(file)` |
| `cadreen.escalations` | `list()`, `get(id)`, `resolve(id, resolution)` |
| `cadreen.healing` | `stats()`, `precedents()`, `diagnose(request)` |
| `cadreen.webhooks` | `create(request)`, `list()`, `delete(id)`, `verifySignature()` |
| `cadreen.learning` | `patterns()`, `episodes()`, `suggestions()` |
| `cadreen.proposals` | `list(options?)`, `get(id)`, `accept(id)`, `dismiss(id, reason?)`, `stats()` |
| `cadreen.setupSessions` | `create(request)`, `list()`, `get(id)`, `addResources(id, request)`, `apply(id, request)` |
| `cadreen.credentials` | `list()`, `create(request)`, `delete(id)` |
| `cadreen.listCapabilities()` | List available capabilities |
| `cadreen.assess(task, domain?)` | Assess task readiness |

## Shorthand

`cadreen.invoke(request)` is an alias for `cadreen.intent.invoke(request)`.

## Changelog

### v0.6.1
- Added `user_id` optional field to `IntentRequest` and `ChatCompletionRequest` — pass end-user identity for per-user context and memory filtering
- Added `WorkspaceUsersResource`: list, invite, updateRole, remove — manage workspace team members
- Added `WorkspaceUser`, `WorkspaceRole`, `InviteUserRequest`, `UpdateRoleRequest` types
- Added `HttpClient.patch()` method for PATCH requests

### v0.6.0 (BREAKING)
- **BREAKING:** `ResolveEscalationRequest.resolution` renamed to `decision` — aligns SDK with API contract
- **BREAKING:** `DiagnoseRequest.error` renamed to `error_message` — aligns SDK with API contract
- Retry default increased from 2 to 3

### v0.5.5
- Added `cadreen.setupSessions` — stateful setup sessions (create, list, get, addResources, apply)
- Added `SetupSession`, `SetupSessionCreateRequest`, `SetupSessionAddRequest`, `SetupSessionApplyRequest`, `SetupSessionApplyResult` types

### v0.5.4
- Added `cadreen.proposals` — task proposals (list, get, accept, dismiss, stats)
- Added `TaskProposal`, `ProposalType`, `ProposalStatus`, `ProposalEvidence` types

### v0.5.3
- Added `dry_run` field to `IntentRequest` — preview intent classification, governance, and capability assessment without creating a mission or persisting conversation
- Fixed Go SDK `Version` constant (was stuck at 0.5.0)

### v0.5.2
- Added `dry_run` mode to `setup()` — preview what would be created without persisting
- Added `notice` and `dry_run` fields to `SetupResult`
- Added `BlueprintsResource` — list, get, create, delete, run, listRuns
- Added `SchedulesResource` — list, get, create, pause, resume
- Added `would_create` to status union types for dry_run responses

### v0.5.1
- Added `cadreen.documents.upload(file)` — upload documents via multipart POST
- Added `cadreen.documents` — document management (list, get, download)
- Added `cadreen.escalations` — escalation management (list, get, resolve)
- Added `cadreen.healing` — self-healing (stats, precedents, diagnose)
- Added `cadreen.webhooks` — webhook CRUD (create, list, delete) + signature verification
- Added `cadreen.learning` — learning insights (patterns, episodes, suggestions)
- Added `cadreen.credentials` — credential management (list, create, delete)
- Added `cadreen.listCapabilities()` — list available capabilities
- Added `cadreen.assess(task, domain?)` — assess task readiness

### v0.5.0 (BREAKING)
- **BREAKING:** API endpoints moved to Cadreen surface (`/api/v1/cadreen/`):
  - `/api/v1/chat/completions` → `/api/v1/cadreen/chat/completions`
  - `/api/v1/tools` → `/api/v1/cadreen/tools`
- All external API calls now route through the Cadreen surface
- Removed "response metadata" terminology (was "intelligence envelope")

### v0.4.0
- Added `cadreen.chat.completions()` — OpenAI-compatible chat completions with governance
- Added `cadreen.chat.completionsStream()` — streaming chat completions via SSE
- Added `cadreen.chat.listTools()` — discover available tools as OpenAI function definitions
- Added tool calling support: `tools` param, `tool_calls` in responses, tool chaining
- Added `conversation_id` for persistent conversations across requests
- Added `ChatCompletionRequest`, `ChatCompletionResponse`, `ChatToolDefinition` and related types

### v0.3.0
- Added `catalog()` — browse the unified integration marketplace
- Added `install(integrationId)` — one-click install with OAuth flow
- Added `CatalogResponse`, `InstallResponse` types

### v0.2.0
- Initial public release

## License

MIT
