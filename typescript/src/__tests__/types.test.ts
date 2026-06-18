import { describe, it, expect } from "vitest";
import { intentStatus } from "../types";
import type { IntentResult } from "../types";

function directResult(overrides: Partial<IntentResult & { type: "direct" }> = {}): IntentResult & { type: "direct" } {
  return {
    type: "direct",
    message: { role: "assistant", content: "Hello!" },
    intelligence: {
      capability: { total_available: 10, healthy_count: 8 },
      reasoning: {},
      memory: { healthy: true },
      governance: { active: false },
      humility: {},
      process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 100 },
      field_stability: { stable: [], evolving: [], internal: [] },
    },
    traceId: "trace-direct",
    ...overrides,
  };
}

function clarifyResult(): IntentResult & { type: "clarify" } {
  return {
    type: "clarify",
    questions: [
      { id: "q1", question: "What date?", type: "date", required: true },
      { id: "q2", question: "Which team?", type: "select", required: false },
    ],
    conversationId: "conv-1",
    intelligence: {
      capability: { total_available: 5, healthy_count: 5 },
      reasoning: {},
      memory: { healthy: true },
      governance: { active: false },
      humility: {},
      process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 50 },
      field_stability: { stable: [], evolving: [], internal: [] },
    },
    traceId: "trace-clarify",
  };
}

function executionResult(): IntentResult & { type: "execution" } {
  return {
    type: "execution",
    execution: { id: "exec-1", status: "running", stream_url: "https://example.com/stream" },
    intelligence: {
      capability: { total_available: 3, healthy_count: 3 },
      reasoning: {},
      memory: { healthy: true },
      governance: { active: false },
      humility: {},
      process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 200 },
      field_stability: { stable: [], evolving: [], internal: [] },
    },
    traceId: "trace-exec",
  };
}

function blockedResult(): IntentResult & { type: "blocked" } {
  return {
    type: "blocked",
    reason_code: "high_risk",
    policy_id: "pol_block",
    intelligence: {
      capability: { total_available: 5, healthy_count: 5 },
      reasoning: {},
      memory: { healthy: true },
      governance: { active: true, decision: "blocked", reason_code: "high_risk" },
      humility: { blocking: 1 },
      process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 30 },
      field_stability: { stable: [], evolving: [], internal: [] },
    },
    traceId: "trace-blocked",
  };
}

function connectRequiredResult(): IntentResult & { type: "connect_required" } {
  return {
    type: "connect_required",
    endpoint: "https://api.example.com",
    reason: "connection required",
    intelligence: {
      capability: { total_available: 0, healthy_count: 0 },
      reasoning: {},
      memory: { healthy: true },
      governance: { active: false },
      humility: {},
      process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 10 },
      field_stability: { stable: [], evolving: [], internal: [] },
    },
    traceId: "trace-connect",
  };
}

describe("intentStatus", () => {
  it("returns ready for direct result", () => {
    const status = intentStatus(directResult());
    expect(status.ready).toBe(true);
    expect(status.needs).toEqual([]);
    expect(status.next).toBe("done");
  });

  it("returns next action for direct result with next_action", () => {
    const status = intentStatus(
      directResult({
        intelligence: {
          ...directResult().intelligence,
          next_action: { type: "add_memory", label: "Record this", reason: "useful" },
        },
      })
    );
    expect(status.ready).toBe(true);
    expect(status.next).toBe("Record this");
  });

  it("returns ready for execution result with stream_url", () => {
    const status = intentStatus(executionResult());
    expect(status.ready).toBe(true);
    expect(status.next).toBe("stream execution");
  });

  it("returns poll fallback for execution without stream_url", () => {
    const result = {
      type: "execution" as const,
      execution: { id: "exec-2", status: "running" },
      intelligence: executionResult().intelligence,
      traceId: "trace-exec-2",
    };
    const status = intentStatus(result);
    expect(status.ready).toBe(true);
    expect(status.next).toBe("poll exec-2");
  });

  it("returns not ready for clarify result", () => {
    const status = intentStatus(clarifyResult());
    expect(status.ready).toBe(false);
    expect(status.needs).toEqual(["What date?", "Which team?"]);
    expect(status.next).toBe("answer 2 questions");
  });

  it("returns singular question language for single question", () => {
    const result: IntentResult = {
      type: "clarify",
      questions: [{ id: "q1", question: "When?", type: "date", required: true }],
      conversationId: "c1",
      intelligence: clarifyResult().intelligence,
      traceId: "t1",
    };
    const status = intentStatus(result);
    expect(status.next).toBe("answer 1 question");
  });

  it("returns not ready for blocked result", () => {
    const status = intentStatus(blockedResult());
    expect(status.ready).toBe(false);
    expect(status.needs).toEqual(["blocked: high_risk"]);
    expect(status.next).toBe("policy: pol_block");
  });

  it("uses governance reason_code fallback for blocked without explicit reason_code", () => {
    const result: IntentResult = {
      type: "blocked",
      intelligence: blockedResult().intelligence,
      traceId: "t2",
    };
    const status = intentStatus(result);
    expect(status.needs).toEqual(["blocked: high_risk"]);
  });

  it("uses generic fallback for blocked without any reason_code", () => {
    const result: IntentResult = {
      type: "blocked",
      intelligence: {
        ...blockedResult().intelligence,
        governance: { active: true },
      },
      traceId: "t3",
    };
    const status = intentStatus(result);
    expect(status.needs).toEqual(["blocked: governance gate"]);
    expect(status.next).toBe("resolve policy block");
  });

  it("returns not ready for connect_required", () => {
    const status = intentStatus(connectRequiredResult());
    expect(status.ready).toBe(false);
    expect(status.needs).toEqual(["connect https://api.example.com"]);
    expect(status.next).toBe("https://api.example.com");
  });

  it("uses fallback for connect_required without endpoint", () => {
    const result: IntentResult = {
      type: "connect_required",
      intelligence: connectRequiredResult().intelligence,
      traceId: "t4",
    };
    const status = intentStatus(result);
    expect(status.needs).toEqual(["connect required tool"]);
    expect(status.next).toBe("");
  });
});
