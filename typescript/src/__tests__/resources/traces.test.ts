import { describe, it, expect } from "vitest";
import { TracesResource } from "../../resources/traces";
import { HttpClient } from "../../client";
import type {
  IntelligenceTraceEntry,
  ListIntelligenceResponse,
  IntelligenceStats,
  ReplayResult,
  HandoffPacket,
  PromoteResult,
  TraceExplain,
  IntelligenceMeta,
} from "../../types";

function baseMeta(): IntelligenceMeta {
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

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("TracesResource get()", () => {
  it("returns trace with explain helper", async () => {
    const fixture: IntelligenceTraceEntry = {
      id: "trace-1",
      domain: "general",
      request_path: "/api/v1/cadreen/intent",
      request_method: "POST",
      meta: {
        ...baseMeta(),
        summary: "Test trace",
        capability: { total_available: 5, healthy_count: 5 },
        governance: { active: true, decision: "auto", confidence: 90 },
        humility: { gaps_detected: 0 },
      },
      created_at: "2026-01-01T00:00:00Z",
    };
    const client = buildClient({ "GET /api/v1/cadreen/intelligence/trace-1": fixture });
    const resource = new TracesResource(client);
    const result = await resource.get("trace-1");
    expect(result.id).toBe("trace-1");
    expect(typeof result.explain).toBe("function");

    const explanation: TraceExplain = result.explain();
    expect(explanation.summary).toBe("Test trace");
    expect(explanation.steps).toContain("5/5 capabilities healthy");
    expect(explanation.steps).toContain("Governance: auto (confidence: 90)");
  });

  it("explain() provides recommendations for blocking gaps", async () => {
    const fixture: IntelligenceTraceEntry = {
      id: "trace-2",
      domain: "general",
      request_path: "/api/v1/cadreen/intent",
      request_method: "POST",
      meta: {
        ...baseMeta(),
        humility: { gaps_detected: 3, blocking: 2 },
        memory: { healthy: false },
      },
    };
    const client = buildClient({ "GET /api/v1/cadreen/intelligence/trace-2": fixture });
    const resource = new TracesResource(client);
    const result = await resource.get("trace-2");
    const explanation = result.explain();
    expect(explanation.recommendations).toBeDefined();
    expect(explanation.recommendations?.length).toBeGreaterThan(0);
    expect(explanation.recommendations).toContain("Resolve blocking capability gaps before proceeding");
    expect(explanation.recommendations).toContain("Knowledge store is degraded — expect reduced context quality");
  });

  it("explain() summary fallback when no summary", async () => {
    const fixture: IntelligenceTraceEntry = {
      id: "trace-3",
      domain: "general",
      request_path: "/api/v1/cadreen/intent",
      request_method: "GET",
      meta: {
        ...baseMeta(),
        capability: { total_available: 0, healthy_count: 0 },
        memory: { healthy: true, knowledge_queried: 5 },
      },
    };
    const client = buildClient({ "GET /api/v1/cadreen/intelligence/trace-3": fixture });
    const resource = new TracesResource(client);
    const result = await resource.get("trace-3");
    const explanation = result.explain();
    expect(explanation.summary).toContain("Trace trace-3:");
    expect(explanation.steps).toContain("5 knowledge items queried");
  });
});

describe("TracesResource list()", () => {
  it("returns fixture data", async () => {
    const fixture: ListIntelligenceResponse = {
      traces: [
        {
          id: "t1",
          domain: "support",
          request_path: "/api/v1/cadreen/intent",
          request_method: "POST",
          meta: baseMeta(),
        },
      ],
      count: 1,
    };
    const client = buildClient({ "GET /api/v1/cadreen/intelligence": fixture });
    const resource = new TracesResource(client);
    const result = await resource.list();
    expect(result.count).toBe(1);
    expect(result.traces[0].id).toBe("t1");
  });

  it("passes query filters", async () => {
    const fixture: ListIntelligenceResponse = { traces: [], count: 0 };
    const client = buildClient({
      "GET /api/v1/cadreen/intelligence?domain=support&limit=5&offset=10": fixture,
    });
    const resource = new TracesResource(client);
    const result = await resource.list({ domain: "support", limit: 5, offset: 10 });
    expect(result.count).toBe(0);
  });
});

describe("TracesResource stats()", () => {
  it("returns fixture data", async () => {
    const fixture: IntelligenceStats = {
      traces_24h: 100,
      traces_7d: 500,
      traces_30d: 2000,
      avg_confidence_by_domain: { general: 85, support: 92 },
      gap_detection_rate: 0.15,
      governance_decisions: { auto: 400, handoff: 100 },
    };
    const client = buildClient({ "GET /api/v1/cadreen/intelligence/stats": fixture });
    const resource = new TracesResource(client);
    const result = await resource.stats();
    expect(result.traces_24h).toBe(100);
    expect(result.avg_confidence_by_domain.general).toBe(85);
  });
});

describe("TracesResource replay()", () => {
  it("returns fixture data", async () => {
    const fixture: ReplayResult = {
      trace_id: "trace-1",
      mode: "current",
      domain: "general",
      original_gate: "blocked",
      original_confidence: 30,
      current_gate: "auto",
      current_confidence: 85,
      gate_changed: true,
      change_summary: "New capabilities enabled",
      current_capability: { total: 10 },
      current_memory: { healthy: true },
      current_gaps: {},
      replay_note: "System evolved since original trace",
    };
    const client = buildClient({
      "POST /api/v1/cadreen/intelligence/trace-1/replay": fixture,
    });
    const resource = new TracesResource(client);
    const result = await resource.replay("trace-1", "current");
    expect(result.gate_changed).toBe(true);
    expect(result.current_confidence).toBe(85);
  });
});

describe("TracesResource handoff()", () => {
  it("returns fixture data", async () => {
    const fixture: HandoffPacket = {
      trace_id: "trace-1",
      domain: "general",
      created_at: "2026-01-01T00:00:00Z",
      governance: { decision: "escalate" },
      what_the_system_knew: { capabilities: 5 },
      what_the_system_didnt_know: { gaps: 2 },
      what_happened: { status: "escalated" },
      suggested_actions: [{ type: "add_policy", label: "Add approval policy" }],
      next_action: { type: "human_handoff", label: "Review", reason: "Need input" },
      trace_url: "https://example.com/trace-1",
    };
    const client = buildClient({
      "GET /api/v1/cadreen/intelligence/trace-1/handoff": fixture,
    });
    const resource = new TracesResource(client);
    const result = await resource.handoff("trace-1");
    expect(result.trace_id).toBe("trace-1");
    expect(result.trace_url).toBe("https://example.com/trace-1");
  });
});

describe("TracesResource promote()", () => {
  it("returns fixture data", async () => {
    const fixture: PromoteResult = {
      id: "promo-1",
      kind: "healing",
      status: "promoted",
      tool_name: "heal_timeout",
      source_trace_id: "trace-1",
    };
    const client = buildClient({
      "POST /api/v1/cadreen/intelligence/trace-1/promote": fixture,
    });
    const resource = new TracesResource(client);
    const result = await resource.promote("trace-1", "healing");
    expect(result.id).toBe("promo-1");
    expect(result.kind).toBe("healing");
    expect(result.status).toBe("promoted");
  });
});
