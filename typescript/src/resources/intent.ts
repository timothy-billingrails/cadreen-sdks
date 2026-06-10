import type {
  IntentRequest,
  IntentResult,
  ResponseMessage,
  IntelligenceMeta,
  ClarificationQuestion,
} from "../types";
import { HttpClient } from "../client";

function mapIntentResponse(raw: {
  id: string;
  type: string;
  trace_id?: string;
  message?: ResponseMessage;
  mission?: { id: string; status: string; stream_url?: string; poll_url?: string };
  clarification?: { questions?: ClarificationQuestion[]; conversation_id?: string };
  meta?: { governance?: { decision?: string; reason?: string } };
  intelligence?: IntelligenceMeta;
}): IntentResult {
  const intelligence = raw.intelligence || {
    capability: { total_available: 0, healthy_count: 0 },
    reasoning: {},
    memory: { healthy: true },
    governance: { active: false },
    humility: {},
    process: { started_at: "", duration_ms: 0 },
    field_stability: { stable: [], evolving: [], internal: [] },
  };

  const traceId = raw.trace_id || raw.id;

  switch (raw.type) {
    case "direct":
      return {
        type: "direct",
        message: raw.message || { role: "assistant", content: "" },
        intelligence,
        traceId,
      };
    case "clarify":
      return {
        type: "clarify",
        questions: (raw.clarification?.questions || []).map((q) =>
          typeof q === "string" ? { id: "", question: q, type: "open", required: false } : q
        ),
        conversationId: raw.clarification?.conversation_id || "",
        intelligence,
        traceId,
      };
    case "mission":
      return {
        type: "execution",
        execution: {
          id: raw.mission?.id || "",
          status: raw.mission?.status || "",
          stream_url: raw.mission?.stream_url,
          poll_url: raw.mission?.poll_url,
        },
        intelligence,
        traceId,
      };
    case "blocked":
      const gov = raw.meta?.governance || {};
      return {
        type: "blocked",
        reason_code: gov.decision || raw.type,
        policy_id: gov.reason,
        intelligence,
        traceId,
      };
    case "connect_required":
      return {
        type: "connect_required",
        endpoint: raw.mission?.stream_url || "",
        reason: raw.meta?.governance?.reason || "connection required",
        intelligence,
        traceId,
      };
    default:
      return {
        type: "direct",
        message: raw.message || { role: "assistant", content: "" },
        intelligence,
        traceId,
      };
  }
}

export class IntentResource {
  constructor(private client: HttpClient) {}

  async invoke(request: IntentRequest): Promise<IntentResult> {
    const raw = await this.client.post<{
      id: string;
      type: string;
      trace_id?: string;
      message?: ResponseMessage;
      mission?: { id: string; status: string; stream_url?: string; poll_url?: string };
      clarification?: { questions?: ClarificationQuestion[]; conversation_id?: string };
      meta?: { governance?: { decision?: string; reason?: string } };
      intelligence?: IntelligenceMeta;
    }>("/api/v1/cadreen/intent", request);

    return mapIntentResponse(raw);
  }
}
