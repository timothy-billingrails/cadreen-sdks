import type { IntelligenceMeta, IntelligenceStage, IntentResult, IntentMessage } from "./types";

export function requiresHuman(intel: IntelligenceMeta): boolean {
  if (intel.governance.decision === "blocked" && (intel.humility.blocking ?? 0) > 0) {
    return true;
  }
  if (intel.stages) {
    for (const stage of intel.stages) {
      if ((stage.name === "escalation" || stage.name === "human_handoff") && stage.status !== "skipped") {
        return true;
      }
    }
  }
  return false;
}

export function handoffReason(intel: IntelligenceMeta): string | null {
  if (intel.stages) {
    for (const stage of intel.stages) {
      if ((stage.name === "escalation" || stage.name === "human_handoff") && stage.status !== "skipped") {
        return stage.detail ?? null;
      }
    }
  }
  if (intel.next_action?.reason) {
    return intel.next_action.reason;
  }
  return null;
}

export function explainTrace(intel: IntelligenceMeta): string {
  const parts: string[] = [];
  if (intel.summary) {
    parts.push(intel.summary);
  }
  if (intel.stages && intel.stages.length > 0) {
    const stageParts = intel.stages.map((s: IntelligenceStage) => {
      let line = `${s.name}: ${s.status}`;
      if (s.detail) {
        line += ` — ${s.detail}`;
      }
      return line;
    });
    parts.push(stageParts.join("; "));
  }
  if (requiresHuman(intel)) {
    parts.push("Requires human intervention");
  }
  return parts.length > 0 ? parts.join("\n") : "No intelligence trace available";
}

const EMAIL_RE = /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/g;
const PHONE_RE = /(?:\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}/g;
const UUID_RE = /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/gi;
const API_KEY_RE = /sk_[a-zA-Z]+_[a-zA-Z0-9]{8,}/g;
const IP_RE = /\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b/g;

function redactString(text: string, options?: RedactOptions): string {
  let result = text;
  result = result.replace(EMAIL_RE, "[email]");
  result = result.replace(PHONE_RE, "[phone]");
  result = result.replace(API_KEY_RE, "[api_key]");
  result = result.replace(IP_RE, "[ip]");
  if (!options?.preserveUUIDs) {
    result = result.replace(UUID_RE, "[id]");
  }
  return result;
}

function redactValue(value: unknown, options?: RedactOptions): unknown {
  if (typeof value === "string") {
    return redactString(value, options);
  }
  if (Array.isArray(value)) {
    return value.map((v) => redactValue(v, options));
  }
    if (value && typeof value === "object") {
      const result: Record<string, unknown> = {};
      for (const [key, val] of Object.entries(value)) {
        if (key === "__proto__" || key === "constructor" || key === "prototype") {
          continue;
        }
      const keysToRedact = options?.keysToRedact ?? ["content", "message", "text", "body", "email", "phone", "address", "name"];
      if (keysToRedact.includes(key.toLowerCase()) && typeof val === "string") {
        result[key] = redactString(val, options);
      } else {
        result[key] = redactValue(val, options);
      }
    }
    return result;
  }
  return value;
}

export interface RedactOptions {
  preserveUUIDs?: boolean;
  keysToRedact?: string[];
}

export function redactTrace(intel: IntelligenceMeta, options?: RedactOptions): IntelligenceMeta {
  return redactValue(intel, options) as IntelligenceMeta;
}

export function redactMessages(messages: IntentMessage[], options?: RedactOptions): IntentMessage[] {
  return messages.map((msg) => ({
    role: msg.role,
    content: redactString(msg.content, options),
  }));
}

export function redactResult(result: IntentResult, options?: RedactOptions): IntentResult {
  const redacted: Record<string, unknown> = { ...result };
  if ("intelligence" in redacted && redacted.intelligence) {
    redacted.intelligence = redactTrace(redacted.intelligence as IntelligenceMeta, options);
  }
  if ("message" in redacted && redacted.message && typeof redacted.message === "object" && "content" in (redacted.message as Record<string, unknown>)) {
    const msg = redacted.message as Record<string, unknown>;
    redacted.message = { ...msg, content: redactString(msg.content as string, options) };
  }
  return redacted as IntentResult;
}
