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
| `"full"` (default) | Full intelligence envelope | You want full transparency |
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
results = await cadreen.memory.search(SearchMemoryRequest(query="data deletion rules"))

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

## Requirements

- Python 3.10+
- httpx >= 0.25
- httpx-sse >= 0.4

## Changelog

### v0.3.0
- Added `catalog()` — browse the unified integration marketplace
- Added `install(integration_id)` — one-click install with OAuth flow
- Added `CatalogResponse`, `InstallResponse` types

### v0.2.1
- Initial public release

## License

MIT
