import { describe, it, expect } from "vitest";
import { IntentResource } from "../../resources/intent";
import { HttpClient, CadreenBlockedError, CadreenClarifyError } from "../../client";
import type { IntentResult, IntelligenceMeta } from "../../types";

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

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("IntentResource invoke() with direct result", () => {
  it("returns direct result", async () => {
    const fixture: IntentResult & { type: "direct" } = {
      type: "direct",
      message: { role: "assistant", content: "Hello, how can I help?" },
      intelligence: baseIntelligence(),
      traceId: "trace-direct-1",
    };
    const raw = {
      id: "resp-1",
      type: "direct",
      trace_id: "trace-direct-1",
      message: { role: "assistant", content: "Hello, how can I help?" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invoke({ messages: [{ role: "user", content: "Hi" }] });
    expect(result.type).toBe("direct");
    if (result.type === "direct") {
      expect(result.message.content).toBe("Hello, how can I help?");
      expect(result.traceId).toBe("trace-direct-1");
    }
  });
});

describe("IntentResource invoke() throws CadreenBlockedError on blocked", () => {
  it("throws CadreenBlockedError", async () => {
    const raw = {
      id: "resp-2",
      type: "blocked",
      trace_id: "trace-blocked-1",
      meta: { governance: { decision: "high_risk", reason: "pol_block" } },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    await expect(
      resource.invoke({ messages: [{ role: "user", content: "Delete everything" }] })
    ).rejects.toThrow(CadreenBlockedError);
    try {
      await resource.invoke({ messages: [{ role: "user", content: "Delete everything" }] });
    } catch (err) {
      expect(err).toBeInstanceOf(CadreenBlockedError);
      if (err instanceof CadreenBlockedError) {
        expect(err.reason_code).toBe("high_risk");
        expect(err.policy_id).toBe("pol_block");
        expect(err.traceId).toBe("trace-blocked-1");
      }
    }
  });
});

describe("IntentResource invoke() throws CadreenClarifyError on clarify", () => {
  it("throws CadreenClarifyError", async () => {
    const raw = {
      id: "resp-3",
      type: "clarify",
      trace_id: "trace-clarify-1",
      clarification: {
        questions: [{ id: "q1", question: "What date?", type: "date", required: true }],
        conversation_id: "conv-1",
      },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    await expect(
      resource.invoke({ messages: [{ role: "user", content: "Schedule something" }] })
    ).rejects.toThrow(CadreenClarifyError);
    try {
      await resource.invoke({ messages: [{ role: "user", content: "Schedule something" }] });
    } catch (err) {
      expect(err).toBeInstanceOf(CadreenClarifyError);
      if (err instanceof CadreenClarifyError) {
        expect(err.questions).toHaveLength(1);
        expect(err.questions[0].question).toBe("What date?");
        expect(err.conversationId).toBe("conv-1");
        expect(err.traceId).toBe("trace-clarify-1");
      }
    }
  });
});

describe("IntentResource invokeResult() returns all types", () => {
  it("returns direct without throwing", async () => {
    const raw = {
      id: "resp-4",
      type: "direct",
      trace_id: "trace-dir",
      message: { role: "assistant", content: "OK" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "Test" }] });
    expect(result.type).toBe("direct");
  });

  it("returns blocked without throwing", async () => {
    const raw = {
      id: "resp-5",
      type: "blocked",
      trace_id: "trace-block",
      meta: { governance: { decision: "high_risk", reason: "pol_1" } },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "Test" }] });
    expect(result.type).toBe("blocked");
    if (result.type === "blocked") {
      expect(result.reason_code).toBe("high_risk");
    }
  });

  it("returns clarify without throwing", async () => {
    const raw = {
      id: "resp-6",
      type: "clarify",
      trace_id: "trace-clar",
      clarification: {
        questions: [{ id: "q1", question: "When?", type: "date", required: true }],
        conversation_id: "conv-2",
      },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "Test" }] });
    expect(result.type).toBe("clarify");
  });

  it("returns connect_required", async () => {
    const raw = {
      id: "resp-7",
      type: "connect_required",
      trace_id: "trace-conn",
      mission: { id: "", status: "", stream_url: "https://api.example.com" },
      meta: { governance: { reason: "connection required" } },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "Test" }] });
    expect(result.type).toBe("connect_required");
  });

  it("returns execution for mission type", async () => {
    const raw = {
      id: "resp-8",
      type: "mission",
      trace_id: "trace-exec",
      mission: { id: "exec-1", status: "running", stream_url: "https://stream.example.com" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "Do a task" }] });
    expect(result.type).toBe("execution");
    if (result.type === "execution") {
      expect(result.execution.id).toBe("exec-1");
      expect(result.execution.status).toBe("running");
    }
  });

  it("fallback to direct for unknown type", async () => {
    const raw = {
      id: "resp-9",
      type: "unknown_type",
      message: { role: "assistant", content: "Default" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "?" }] });
    expect(result.type).toBe("direct");
  });

  it("handles string clarification questions", async () => {
    const raw = {
      id: "resp-10",
      type: "clarify",
      clarification: {
        questions: ["What is your name?"],
        conversation_id: "conv-str",
      },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "?" }] });
    expect(result.type).toBe("clarify");
    if (result.type === "clarify") {
      expect(result.questions[0].question).toBe("What is your name?");
      expect(result.questions[0].id).toBe("");
      expect(result.questions[0].type).toBe("open");
      expect(result.questions[0].required).toBe(false);
    }
  });
});

describe("IntentResource sandbox fixture matching", () => {
  it("matches by METHOD /path key", async () => {
    const raw = {
      id: "match-1",
      type: "direct",
      message: { role: "assistant", content: "Matched" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const resource = new IntentResource(client);
    const result = await resource.invokeResult({ messages: [{ role: "user", content: "Hi" }] });
    expect(result.type).toBe("direct");
  });
});
