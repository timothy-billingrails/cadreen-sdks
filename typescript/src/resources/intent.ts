import type {
  IntentRequest,
  IntentResult,
  ResponseMessage,
  IntelligenceMeta,
  ClarificationQuestion,
} from "../types";
import { HttpClient, CadreenBlockedError, CadreenClarifyError } from "../client";

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

  private async rawInvoke(request: IntentRequest): Promise<IntentResult> {
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

  /**
   * Invoke intent and throw on blocked/clarify outcomes.
   *
   * Throws:
   * - CadreenBlockedError when governance blocks the action
   * - CadreenClarifyError when the system needs clarification
   *
   * Use try/catch to handle governed outcomes:
   * ```ts
   * try {
   *   const result = await cadreen.intent.invoke({ messages: [...] });
   *   // result is always type "direct" or "execution"
   * } catch (err) {
   *   if (err instanceof CadreenBlockedError) {
   *     // show user the block reason and escalation path
   *   } else if (err instanceof CadreenClarifyError) {
   *     // present err.questions to the user
   *   }
   * }
   * ```
   */
  async invoke(request: IntentRequest): Promise<IntentResult> {
    const result = await this.rawInvoke(request);

    if (result.type === "blocked") {
      throw new CadreenBlockedError({
        reason_code: result.reason_code,
        policy_id: result.policy_id,
        intelligence: result.intelligence,
        traceId: result.traceId,
      });
    }

    if (result.type === "clarify") {
      throw new CadreenClarifyError({
        questions: result.questions.map((q) => ({
          id: q.id || "",
          question: q.question || "",
          type: q.type || "open",
          required: q.required ?? false,
        })),
        conversationId: result.conversationId,
        intelligence: result.intelligence,
        traceId: result.traceId,
      });
    }

    return result;
  }

  /**
   * Invoke intent and return the full result without throwing.
   * Use this when you want to handle all outcome types (direct, clarify,
   * execution, blocked, connect_required) as data rather than errors.
   */
  async invokeResult(request: IntentRequest): Promise<IntentResult> {
    return this.rawInvoke(request);
  }
}
