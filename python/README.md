# cadreen-sdk

Python SDK for [Cadreen](https://accomplishanything.today/infra/docs) — Intelligence as a Service.

Cadreen is a cognitive operating system. Send messages describing what you want done, and Cadreen reasons, connects tools, recalls knowledge, governs actions, and escalates to humans when needed. The SDK handles authentication, retries, idempotency, streaming, and error classification.

## Install

```bash
pip install cadreen-sdk
```

## Quick Start

```python
import asyncio
from cadreen import Cadreen

cadreen = Cadreen(api_key="sk_cadreen_...")

async def main():
    # Intent — the main entry point
    result = await cadreen.ask(
        "Handle a refund request for invoice inv_123",
    )
    print(result.explain())

    # Discriminated union with match/case
    match result.type:
        case "direct":
            print(result.message.content)
        case "clarify":
            for q in result.questions:
                print(q)
        case "execution":
            async for event in cadreen.executions.stream(result.execution["id"]):
                print(event.type, event.data)
        case "blocked":
            print(result.reason_code)
        case "connect_required":
            print(result.endpoint)

asyncio.run(main())
```

## Chat Completions

OpenAI-compatible chat completions with built-in governance. Every tool call goes through governance before execution — auto-approved calls execute silently, blocked calls become conversation.

```python
from cadreen.resources.chat import (
    ChatCompletionRequest, ChatMessage,
    ChatToolDefinition, ChatFunctionDefinition,
)

# Basic completion
response = await cadreen.chat.completions(ChatCompletionRequest(
    messages=[ChatMessage(role="user", content="Hello!")],
))
print(response.choices[0].message.content)

# With tool calling
response = await cadreen.chat.completions(ChatCompletionRequest(
    messages=[ChatMessage(role="user", content="Refund order 456")],
    tools=[ChatToolDefinition(function=ChatFunctionDefinition(
        name="process_refund",
        description="Process a refund for an order",
        parameters={
            "type": "object",
            "properties": {"order_id": {"type": "string"}},
            "required": ["order_id"],
        },
    ))],
))

# If governance needs approval, the response contains a text message
# asking the user to confirm — no tool_calls field
msg = response.choices[0].message
if msg.tool_calls:
    for tc in msg.tool_calls:
        print(f"{tc.function.name}({tc.function.arguments})")

# Resume a conversation
follow_up = await cadreen.chat.completions(ChatCompletionRequest(
    messages=[ChatMessage(role="user", content="What about order 789?")],
    conversation_id=response.conversation_id,
))
```

### Streaming

```python
async for chunk in await cadreen.chat.completions_stream(ChatCompletionRequest(
    messages=[ChatMessage(role="user", content="Hello!")],
)):
    if chunk.choices and chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

### Tool Discovery

```python
tools = await cadreen.chat.list_tools()
for tool in tools.data:
    print(f"{tool.function.name}: {tool.function.description}")
```

### Tool Chaining

When the model proposes tool calls, send results back for follow-up:

```python
response = await cadreen.chat.completions(ChatCompletionRequest(
    messages=[
        ChatMessage(role="user", content="What's the weather in NYC?"),
        ChatMessage(role="assistant", tool_calls=[ChatToolCall(
            id="tc_1",
            function=ChatFunctionCall(name="get_weather", arguments='{"city":"NYC"}'),
        )]),
        ChatMessage(role="tool", tool_call_id="tc_1", content='{"temp": 72, "condition": "sunny"}'),
    ],
))
# Model may propose more tools or return a final text response
```

## Configuration

```python
cadreen = Cadreen(
    api_key="sk_cadreen_...",
    base_url="https://accomplishanything.today",  # default
    max_retries=2,      # default 2
    timeout=30,          # default 30s
    profile="lean",     # optional: "lean" | "audit" | "full" (default "full")
)
```

### Response Profiles

| Profile | What you get | Use when |
|---------|-------------|----------|
| `"full"` (default) | Full response metadata | You want full transparency |
| `"audit"` | Only governance decision + confidence + blocking gaps | You need to react to gates |
| `"lean"` | No envelope. Just `trace_id` | Hot-looping, minimal payload |

## Marketplace

Browse and install integrations without knowing which provider powers them:

```python
# Browse available integrations
catalog = await cadreen.connections.catalog()
for category in catalog.categories:
    print(f"{category.name}: {len(category.integrations)} integrations")

# One-click install (returns OAuth URL)
install = await cadreen.connections.install("slack")
if install.status == "pending_auth":
    print(f"Authenticate at: {install.auth_url}")

# Check what's installed
print(catalog.installed)  # ["stripe", "github"]
```

## Memory

```python
# Store knowledge
await cadreen.memory.remember(
    type="reference",
    content={"text": "GDPR Article 17: Right to erasure", "title": "GDPR Art. 17"},
    authority=10,
)

# Search
results = await cadreen.memory.search("data deletion rules")

# Get by ID
item = await cadreen.memory.get("mem_abc123")
```

## Policies

```python
# Create a policy
await cadreen.policies.create(
    name="refund_threshold",
    rules=[{"condition": "refund_amount > 500", "action": "require_approval"}],
)

# Evaluate an action
evaluation = await cadreen.policies.evaluate(
    action="Process $750 refund for order 456"
)
```

## Connections

```python
# Register from OpenAPI spec
await cadreen.connections.register_openapi(
    name="internal-erp",
    spec_url="https://erp.example.com/openapi.json",
)

# Register MCP server
await cadreen.connections.register_mcp(
    name="my-mcp-server",
    url="https://mcp.example.com/sse",
    transport="sse",
)

# List installed
connections = await cadreen.connections.list()
```

## Traces

```python
trace = await cadreen.traces.get(result.trace_id)
print(trace.explain())

recent = await cadreen.traces.list(limit=10)
stats = await cadreen.traces.stats()
```

## Documents

```python
# List documents
docs = await cadreen.documents.list()
for doc in docs.documents:
    print(f"{doc.name} ({doc.content_type})")

# Upload a document from file path
result = await cadreen.documents.upload("/path/to/report.pdf")
print(f"Uploaded: {result.name} (ID: {result.id})")

# Upload from bytes
result = await cadreen.documents.upload_bytes(b"file contents", "data.txt")
print(f"Uploaded: {result.name} (ID: {result.id})")

# Get document details
doc = await cadreen.documents.get(result.id)
print(f"Status: {doc.status}, Size: {doc.size} bytes")
```

## Requirements

- Python 3.10+
- httpx >= 0.25
- httpx-sse >= 0.4

## Changelog

### v0.5.2
- Added `dry_run` mode to `setup()` — preview what would be created without persisting
- Added `notice` and `dry_run` fields to `SetupResult`
- Added `cadreen.list_blueprints()`, `cadreen.get_blueprint()`, `cadreen.create_blueprint()`, `cadreen.delete_blueprint()`, `cadreen.run_blueprint()`, `cadreen.list_blueprint_runs()`
- Added `cadreen.list_schedules()`, `cadreen.get_schedule()`, `cadreen.create_schedule()`, `cadreen.pause_schedule()`, `cadreen.resume_schedule()`

### v0.5.1
- Added `cadreen.documents.upload(file_path)` and `cadreen.documents.upload_bytes(content, filename)` — upload documents
- Added `cadreen.documents` — document management (list, get)
- Added `cadreen.escalations` — escalation management (list, get, resolve)
- Added `cadreen.healing` — self-healing (stats, precedents, diagnose)
- Added `cadreen.webhooks` — webhook CRUD (create, list, delete)
- Added `cadreen.learning` — learning insights (patterns, episodes, suggestions)
- Added `cadreen.credentials` — credential management (list, create, delete)
- Added `cadreen.list_capabilities()` — list available capabilities
- Added `cadreen.assess(task, domain=)` — assess task readiness

### v0.5.0 (BREAKING)
- **BREAKING:** API endpoints moved to Cadreen surface (`/api/v1/cadreen/`):
  - `/api/v1/chat/completions` → `/api/v1/cadreen/chat/completions`
  - `/api/v1/tools` → `/api/v1/cadreen/tools`
- All external API calls now route through the Cadreen surface
- Removed "response metadata" terminology (was "intelligence envelope")

### v0.4.0
- Added `cadreen.chat.completions()` — OpenAI-compatible chat completions with governance
- Added `cadreen.chat.completions_stream()` — streaming chat completions via SSE
- Added `cadreen.chat.list_tools()` — discover available tools as OpenAI function definitions
- Added tool calling support: `tools` param, `tool_calls` in responses, tool chaining
- Added `conversation_id` for persistent conversations across requests
- Added `ChatCompletionRequest`, `ChatCompletionResponse`, `ChatToolDefinition` and related types

### v0.3.0
- Added `catalog()` — browse the unified integration marketplace
- Added `install(integration_id)` — one-click install with OAuth flow
- Added `CatalogResponse`, `InstallResponse` types

### v0.2.1
- Initial public release

## License

MIT
