import type {
  IntelligenceTraceEntry,
  ListIntelligenceResponse,
  TraceExplain,
} from "../types";
import { TracesResource } from "./traces";

/**
 * FailuresResource — Builder-shaped naming facade for understanding errors.
 *
 * When something goes wrong, developers ask `cadreen.failures.explain(id)`
 * or `cadreen.failures.recent()` instead of navigating the trace API directly.
 *
 * This is a thin wrapper; no new backend routes are introduced.
 */
export class FailuresResource {
  constructor(private traces: TracesResource) {}

  /** Fetch a trace and return it with an `.explain()` helper. */
  async explain(id: string): Promise<IntelligenceTraceEntry & { explain: () => TraceExplain }> {
    return this.traces.get(id);
  }

  /** List recent traces with optional filtering. */
  async recent(options?: {
    domain?: string;
    decision?: string;
    from?: string;
    to?: string;
    limit?: number;
    offset?: number;
  }): Promise<ListIntelligenceResponse> {
    return this.traces.list(options);
  }

  /** Summarize why an error occurred, using intelligence data when available. */
  why(error: unknown): { summary: string; traceId?: string; recommendation?: string } {
    if (error && typeof error === "object") {
      const e = error as Record<string, unknown>;
      const traceId =
        typeof e.traceId === "string"
          ? e.traceId
          : typeof e.trace_id === "string"
          ? e.trace_id
          : undefined;
      const code = typeof e.code === "string" ? e.code : "unknown";
      const message = typeof e.message === "string" ? e.message : String(error);

      return {
        summary: `${code}: ${message}`,
        traceId,
        recommendation:
          code === "timeout"
            ? "Retry with a longer timeout or check network connectivity."
            : code === "rate_limited"
            ? "Back off and retry with exponential delay."
            : code === "policy_blocked"
            ? "Review the policy that blocked this request or request an override."
            : traceId
            ? `Look up trace ${traceId} for full provenance.`
            : "Check recent traces for context.",
      };
    }

    return {
      summary: String(error),
      recommendation: "Check recent traces for context.",
    };
  }
}
