# cadreen-sdk

Python SDK for Cadreen — intelligence-native automation infrastructure.

## Install

```bash
pip install cadreen-sdk
```

## Quick Start

```python
import os
import asyncio
from cadreen import Cadreen

cadreen = Cadreen(api_key=os.environ["CADREEN_API_KEY"])

async def main():
    # Intent — the main entry point
    result = await cadreen.intent(
        "Handle a refund request for invoice inv_123",
        domain="support",
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
            print(result.policy["reason"])
        case "connect_required":
            print(result.connection["endpoint"])

    # Memory — store knowledge
    await cadreen.memory.remember(
        type="reference",
        content={"text": "GDPR Article 17...", "title": "GDPR Art. 17"},
        authority=10,
    )

    # Policies — governance guardrails
    await cadreen.policies.require_approval("Refunds over $500 require human approval")

    # Traces — intelligence observability
    trace = await cadreen.traces.get(result.trace_id)
    explanation = trace.explain()
    print(explanation.summary)

    # Connections — register external services
    await cadreen.connections.register_openapi(
        name="internal-erp",
        spec_url="https://erp.example.com/openapi.json",
    )

asyncio.run(main())
```

## Configuration

```python
cadreen = Cadreen(
    api_key="...",
    base_url="https://accomplishanything.today",  # default, configurable
    max_retries=2,      # default 2
    timeout=30,          # default 30s
    profile="lean",     # optional: "lean" | "audit" | "full" (default "full")
)
```

### Response Profiles

Control how much intelligence metadata you get back:

| Profile | What you get | Use when |
|---------|-------------|----------|
| `"full"` (default) | Full intelligence envelope | You want full transparency |
| `"audit"` | Only governance decision + confidence + blocking gaps | You need to react to gates |
| `"lean"` | No envelope. Just `trace_id` | Hot-looping, minimal payload |

## Requirements

- Python 3.10+
- httpx >= 0.25
- httpx-sse >= 0.4
