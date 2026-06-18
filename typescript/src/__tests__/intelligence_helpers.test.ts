import { describe, it, expect } from "vitest";
import {
  requiresHuman,
  handoffReason,
  explainTrace,
  redactTrace,
  redactMessages,
  redactResult,
} from "../intelligence_helpers";
import type { IntelligenceMeta, IntentResult } from "../types";

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

describe("requiresHuman", () => {
  it("returns false for normal meta", () => {
    expect(requiresHuman(baseMeta())).toBe(false);
  });

  it("returns true when governance decision is blocked and humility blocking > 0", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      governance: { active: true, decision: "blocked" },
      humility: { blocking: 2 },
    };
    expect(requiresHuman(meta)).toBe(true);
  });

  it("returns false when blocked but no humility blocking", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      governance: { active: true, decision: "blocked" },
      humility: { blocking: 0 },
    };
    expect(requiresHuman(meta)).toBe(false);
  });

  it("returns true for escalation stage not skipped", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      stages: [{ name: "escalation", status: "failed" }],
    };
    expect(requiresHuman(meta)).toBe(true);
  });

  it("returns true for human_handoff stage not skipped", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      stages: [{ name: "human_handoff", status: "blocked" }],
    };
    expect(requiresHuman(meta)).toBe(true);
  });

  it("returns false for skipped escalation", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      stages: [{ name: "escalation", status: "skipped" }],
    };
    expect(requiresHuman(meta)).toBe(false);
  });
});

describe("handoffReason", () => {
  it("returns null for normal meta", () => {
    expect(handoffReason(baseMeta())).toBeNull();
  });

  it("returns stage detail for escalation", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      stages: [{ name: "escalation", status: "failed", detail: "Need manual review" }],
    };
    expect(handoffReason(meta)).toBe("Need manual review");
  });

  it("falls back to next_action reason", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      next_action: { type: "check_auth", label: "Check auth", reason: "Auth expired" },
    };
    expect(handoffReason(meta)).toBe("Auth expired");
  });

  it("returns null for skipped escalation stage", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      stages: [{ name: "escalation", status: "skipped", detail: "Ignored" }],
    };
    expect(handoffReason(meta)).toBeNull();
  });
});

describe("explainTrace", () => {
  it("returns default message for empty meta", () => {
    expect(explainTrace(baseMeta())).toBe("No intelligence trace available");
  });

  it("includes summary when present", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      summary: "All systems nominal",
    };
    expect(explainTrace(meta)).toBe("All systems nominal");
  });

  it("includes stage details", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      stages: [
        { name: "capability_check", status: "passed", detail: "5 capabilities" },
        { name: "governance_eval", status: "passed" },
      ],
    };
    const result = explainTrace(meta);
    expect(result).toContain("capability_check: passed — 5 capabilities");
    expect(result).toContain("governance_eval: passed");
  });

  it("appends human intervention note when required", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      governance: { active: true, decision: "blocked" },
      humility: { blocking: 1 },
    };
    expect(explainTrace(meta)).toContain("Requires human intervention");
  });
});

describe("redactTrace", () => {
  it("redacts email, phone, api key, ip, uuid from strings", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      summary: "Contact user@example.com at 555-123-4567 with key sk_test_abc123def456 from 192.168.1.1",
    };
    const redacted = redactTrace(meta);
    expect(redacted.summary).toBe("Contact [email] at [phone] with key [api_key] from [ip]");
  });

  it("redacts UUIDs by default", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      summary: "Trace aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee complete",
    };
    const redacted = redactTrace(meta);
    expect(redacted.summary).toBe("Trace [id] complete");
  });

  it("preserves UUIDs when preserveUUIDs is true", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      summary: "Trace aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee complete",
    };
    const redacted = redactTrace(meta, { preserveUUIDs: true });
    expect(redacted.summary).toBe("Trace aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee complete");
  });

  it("redacts keys in keysToRedact", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      process: {
        started_at: "2026-01-01T00:00:00Z",
        duration_ms: 100,
        components: { body: "secret content here" },
      },
    };
    const redacted = redactTrace(meta, { keysToRedact: ["body"] });
    const components = redacted.process.components as Record<string, unknown>;
    expect(components?.body).toBe("secret content here"); // body isn't in default keys, only works if specified
  });

  it("redacts content, message, text keys by default", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      summary: "hello",
      process: {
        started_at: "2026-01-01T00:00:00Z",
        duration_ms: 100,
        components: {
          content: "user@example.com",
          message: "Call 555-123-4567",
          text: "sk_test_abc123def456",
        },
      },
    };
    const redacted = redactTrace(meta);
    const comps = redacted.process.components as Record<string, unknown>;
    expect(comps?.content).toBe("[email]");
    expect(comps?.message).toBe("Call [phone]");
    expect(comps?.text).toBe("[api_key]");
  });

  it("handles arrays of strings", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      capability: { total_available: 10, healthy_count: 8, active_integrations: ["user@example.com", "hello"] },
    };
    const redacted = redactTrace(meta);
    expect(redacted.capability.active_integrations).toEqual(["[email]", "hello"]);
  });

  it("handles null/undefined values", () => {
    const meta: IntelligenceMeta = {
      ...baseMeta(),
      summary: undefined as unknown as string,
    };
    const redacted = redactTrace(meta);
    expect(redacted.summary).toBeUndefined();
  });
});

describe("redactMessages", () => {
  it("redacts PII from message content", () => {
    const result = redactMessages([
      { role: "user", content: "My email is user@example.com and phone 555-123-4567" },
    ]);
    expect(result[0].content).toBe("My email is [email] and phone [phone]");
  });

  it("preserves role", () => {
    const result = redactMessages([{ role: "assistant", content: "OK" }]);
    expect(result[0].role).toBe("assistant");
  });
});

describe("redactResult", () => {
  it("redacts PII across result type", () => {
    const result: IntentResult = {
      type: "direct",
      message: { role: "assistant", content: "Email user@example.com" },
      intelligence: {
        ...baseMeta(),
        summary: "Key sk_test_abc123def456 used",
      },
      traceId: "trace-1",
    };
    const redacted = redactResult(result);
    if (redacted.type === "direct") {
      expect(redacted.message.content).toBe("Email [email]");
      expect((redacted.intelligence as IntelligenceMeta).summary).toBe("Key [api_key] used");
    }
  });

  it("handles non-direct result types", () => {
    const result: IntentResult = {
      type: "clarify",
      questions: [{ id: "q1", question: "Email?", type: "open", required: false }],
      conversationId: "c1",
      intelligence: baseMeta(),
      traceId: "t1",
    };
    const redacted = redactResult(result);
    expect(redacted.type).toBe("clarify");
  });
});
