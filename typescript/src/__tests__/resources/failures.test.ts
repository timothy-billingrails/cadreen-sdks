import { describe, it, expect } from "vitest";
import { FailuresResource } from "../../resources/failures";
import { TracesResource } from "../../resources/traces";
import { HttpClient, CadreenBlockedError } from "../../client";
import type { IntelligenceMeta, IntelligenceTraceEntry } from "../../types";

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

describe("FailuresResource explain()", () => {
  it("returns trace with explain helper", async () => {
    const fixture: IntelligenceTraceEntry = {
      id: "fail-1",
      domain: "general",
      request_path: "/api/v1/cadreen/intent",
      request_method: "POST",
      meta: {
        ...baseMeta(),
        summary: "Failure trace",
      },
    };
    const client = new HttpClient({
      apiKey: "sk_test",
      sandbox: true,
      fixtures: { "GET /api/v1/cadreen/intelligence/fail-1": fixture },
    });
    const traces = new TracesResource(client);
    const failures = new FailuresResource(traces);
    const result = await failures.explain("fail-1");
    expect(result.id).toBe("fail-1");
    expect(typeof result.explain).toBe("function");
    const exp = result.explain();
    expect(exp.summary).toBe("Failure trace");
  });
});

describe("FailuresResource recent()", () => {
  it("returns trace list via delegate", async () => {
    const fixture = {
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
    const client = new HttpClient({
      apiKey: "sk_test",
      sandbox: true,
      fixtures: { "GET /api/v1/cadreen/intelligence": fixture },
    });
    const traces = new TracesResource(client);
    const failures = new FailuresResource(traces);
    const result = await failures.recent();
    expect(result.count).toBe(1);
    expect(result.traces[0].id).toBe("t1");
  });

  it("passes filter options via delegate", async () => {
    const fixture = { traces: [], count: 0 };
    const client = new HttpClient({
      apiKey: "sk_test",
      sandbox: true,
      fixtures: {
        "GET /api/v1/cadreen/intelligence?domain=support&limit=10": fixture,
      },
    });
    const traces = new TracesResource(client);
    const failures = new FailuresResource(traces);
    const result = await failures.recent({ domain: "support", limit: 10 });
    expect(result.count).toBe(0);
  });
});

describe("FailuresResource why()", () => {
  const traces = new TracesResource(
    new HttpClient({ apiKey: "sk_test", sandbox: true, fixtures: {} })
  );
  const failures = new FailuresResource(traces);

  it("extracts info from CadreenBlockedError", () => {
    const err = new CadreenBlockedError({
      reason_code: "high_risk",
      policy_id: "pol_1",
      intelligence: {},
      traceId: "trace-xyz",
    });
    const result = failures.why(err);
    expect(result.summary).toContain("high_risk");
    expect(result.traceId).toBe("trace-xyz");
    expect(result.recommendation).toContain("Look up trace trace-xyz");
  });

  it("extracts info from generic error object", () => {
    const err = new Error("Network timeout");
    (err as Record<string, unknown>).code = "timeout";
    (err as Record<string, unknown>).trace_id = "trace-timeout";
    const result = failures.why(err);
    expect(result.summary).toBe("timeout: Network timeout");
    expect(result.traceId).toBe("trace-timeout");
    expect(result.recommendation).toContain("Retry with a longer timeout");
  });

  it("handles rate_limited error code", () => {
    const err = new Error("Too many requests");
    (err as Record<string, unknown>).code = "rate_limited";
    const result = failures.why(err);
    expect(result.recommendation).toContain("Back off and retry");
  });

  it("handles plain string error", () => {
    const result = failures.why("Something went wrong");
    expect(result.summary).toBe("Something went wrong");
    expect(result.recommendation).toBe("Check recent traces for context.");
  });

  it("handles null/undefined", () => {
    const result = failures.why(null);
    expect(result.summary).toBe("null");
    expect(result.recommendation).toBe("Check recent traces for context.");
  });

  it("handles missing code in error object", () => {
    const err = new Error("Generic failure");
    const result = failures.why(err);
    expect(result.summary).toBe("unknown: Generic failure");
  });
});
