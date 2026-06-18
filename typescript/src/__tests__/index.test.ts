import { describe, it, expect } from "vitest";
import { Cadreen } from "../index";
import type { IntelligenceMeta, IntentResult } from "../types";

function baseIntelligence(): IntelligenceMeta {
  return {
    capability: { total_available: 10, healthy_count: 8 },
    reasoning: {},
    memory: { healthy: true },
    governance: { active: false },
    humility: {},
    process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 100 },
    field_stability: { stable: [], evolving: [], internal: [] },
  };
}

function fixtures() {
  return {
    // Intent fixtures (ask/act -> chat/execution mode)
    "POST /api/v1/cadreen/intent": {
      id: "resp-integration",
      type: "direct",
      trace_id: "trace-int",
      message: { role: "assistant", content: "Integration response" },
      intelligence: baseIntelligence(),
    },
    // Memory fixtures
    "POST /api/v1/cadreen/memory": {
      id: "mem-int",
      type: "note",
      domain: "integration",
      content: { text: "Integration memory" },
      authority: 50,
      version: 1,
    },
    // Connections fixtures
    "POST /api/v1/cadreen/connections": {
      type: "prebuilt",
      capability: "test_capability",
      detail: {
        tool_id: "tool-test",
        tool_name: "Test Tool",
        service_id: "test-svc",
        service_name: "Test Service",
        auth_type: "api_key",
        source: "cadreen",
      },
    },
    // Policies fixtures
    "GET /api/v1/cadreen/policies": {
      policies: [{ id: "p-int", name: "Integration Policy", domain: "general", priority: 1, requires_human: false }],
    },
    // Traces fixtures
    "GET /api/v1/cadreen/intelligence": {
      traces: [
        {
          id: "trace-int",
          domain: "integration",
          request_path: "/api/v1/cadreen/intent",
          request_method: "POST",
          meta: baseIntelligence(),
        },
      ],
      count: 1,
    },
    // Executions fixtures
    "GET /api/v1/cadreen/executions/exec-int": {
      id: "exec-int",
      status: "completed",
      progress: 100,
      result: { output: "Done" },
    },
    // Setup
    "POST /api/v1/cadreen/setup": {
      workspace_id: "ws-int",
      connections: [],
      credentials: [],
      memory: [],
      policies: [],
      applied: 0,
      failed: 0,
    },
  };
}

describe("Cadreen integration", () => {
  const cadreen = new Cadreen({
    apiKey: "sk_integration_test",
    sandbox: true,
    fixtures: fixtures(),
  });

  it("all public resource members are accessible", () => {
    expect(cadreen.intent).toBeDefined();
    expect(cadreen.memory).toBeDefined();
    expect(cadreen.policies).toBeDefined();
    expect(cadreen.connections).toBeDefined();
    expect(cadreen.traces).toBeDefined();
    expect(cadreen.executions).toBeDefined();
    expect(cadreen.guardrails).toBeDefined();
    expect(cadreen.skills).toBeDefined();
    expect(cadreen.failures).toBeDefined();
    expect(cadreen.webhooks).toBeDefined();
  });

  it("cadreen.ask() returns intent result", async () => {
    const result = await cadreen.ask("Hello");
    expect(result.type).toBe("direct");
    if (result.type === "direct") {
      expect(result.message.content).toBe("Integration response");
    }
  });

  it("cadreen.act() returns intent result", async () => {
    const result = await cadreen.act("Do something");
    expect(result.type).toBe("direct");
  });

  it("cadreen.remember() returns memory response", async () => {
    const result = await cadreen.remember({
      type: "note",
      content: { text: "Remember integration test" },
    });
    expect(result.id).toBe("mem-int");
  });

  it("cadreen.context() returns search response", async () => {
    // We need a search fixture
    const cadreenSearch = new Cadreen({
      apiKey: "sk_test",
      sandbox: true,
      fixtures: {
        ...fixtures(),
        "GET /api/v1/cadreen/memory/search?query=test": {
          results: [{ id: "a1", type: "note", domain: "int", authority: 50, version: 1, content: { text: "found" } }],
          count: 1,
        },
      },
    });
    const result = await cadreenSearch.context({ query: "test" });
    expect(result.count).toBe(1);
  });

  it("cadreen.connect() returns connect result", async () => {
    const result = await cadreen.connect("test_capability");
    expect(result.type).toBe("prebuilt");
  });

  it("cadreen.invoke() delegates to intent", async () => {
    const result = await cadreen.invoke({
      messages: [{ role: "user", content: "Hello" }],
    });
    expect(result.type).toBe("direct");
  });

  it("cadreen.setup() returns setup result", async () => {
    const result = await cadreen.setup({ workspace_id: "ws-int" });
    expect(result.workspace_id).toBe("ws-int");
    expect(result.applied).toBe(0);
  });

  it("policies.list() works through cadreen", async () => {
    const result = await cadreen.policies.list();
    expect(result.policies).toHaveLength(1);
  });

  it("traces.list() works through cadreen", async () => {
    const result = await cadreen.traces.list();
    expect(result.count).toBe(1);
  });

  it("executions.getStatus() works through cadreen", async () => {
    const result = await cadreen.executions.getStatus("exec-int");
    expect(result.id).toBe("exec-int");
    expect(result.status).toBe("completed");
  });

  it("webhooks.verifySignature() works through cadreen", () => {
    expect(typeof cadreen.webhooks.verifySignature).toBe("function");
    expect(cadreen.webhooks.verifySignature("", "", "")).toBe(false);
  });
});
