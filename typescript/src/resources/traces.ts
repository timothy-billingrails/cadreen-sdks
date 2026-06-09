import type {
  IntelligenceTraceEntry,
  ListIntelligenceResponse,
  IntelligenceStats,
  TraceExplain,
  ReplayResult,
  HandoffPacket,
  PromoteResult,
} from "../types";
import { HttpClient } from "../client";

export class TracesResource {
  constructor(private client: HttpClient) {}

  async get(id: string): Promise<IntelligenceTraceEntry & { explain: () => TraceExplain }> {
    const trace = await this.client.get<IntelligenceTraceEntry>(
      `/api/v1/cadreen/intelligence/${encodeURIComponent(id)}`
    );
    return {
      ...trace,
      explain: (): TraceExplain => {
        const meta = trace.meta;
        const steps: string[] = [];
        if (meta.capability.total_available > 0) {
          steps.push(`${meta.capability.healthy_count}/${meta.capability.total_available} capabilities healthy`);
        }
        if (meta.governance.active) {
          steps.push(`Governance: ${meta.governance.decision || "active"} (confidence: ${meta.governance.confidence || 0})`);
        }
        if (meta.humility.gaps_detected && meta.humility.gaps_detected > 0) {
          steps.push(`${meta.humility.gaps_detected} gaps detected (${meta.humility.blocking || 0} blocking)`);
        }
        if (meta.memory.knowledge_queried) {
          steps.push(`${meta.memory.knowledge_queried} knowledge items queried`);
        }
        const recommendations: string[] = [];
        if (meta.humility.blocking && meta.humility.blocking > 0) {
          recommendations.push("Resolve blocking capability gaps before proceeding");
        }
        if (!meta.memory.healthy) {
          recommendations.push("Knowledge store is degraded — expect reduced context quality");
        }
        return {
          summary: meta.summary || `Trace ${trace.id}: ${steps.join("; ")}`,
          steps,
          recommendations: recommendations.length > 0 ? recommendations : undefined,
        };
      },
    };
  }

  async list(options?: {
    domain?: string;
    decision?: string;
    from?: string;
    to?: string;
    limit?: number;
    offset?: number;
  }): Promise<ListIntelligenceResponse> {
    const params: Record<string, string | number | undefined> = {};
    if (options?.domain) params.domain = options.domain;
    if (options?.decision) params.decision = options.decision;
    if (options?.from) params.from = options.from;
    if (options?.to) params.to = options.to;
    if (options?.limit) params.limit = options.limit;
    if (options?.offset) params.offset = options.offset;
    return this.client.get<ListIntelligenceResponse>("/api/v1/cadreen/intelligence", params);
  }

  async stats(): Promise<IntelligenceStats> {
    return this.client.get<IntelligenceStats>("/api/v1/cadreen/intelligence/stats");
  }

  async replay(id: string, mode?: "current" | "historical"): Promise<ReplayResult> {
    return this.client.post<ReplayResult>(`/api/v1/cadreen/intelligence/${encodeURIComponent(id)}/replay`, {
      mode: mode || "current",
    });
  }

  async handoff(id: string): Promise<HandoffPacket> {
    return this.client.get<HandoffPacket>(`/api/v1/cadreen/intelligence/${encodeURIComponent(id)}/handoff`);
  }

  async promote(id: string, kind: "healing" | "memory" | "procedure", name?: string): Promise<PromoteResult> {
    return this.client.post<PromoteResult>(`/api/v1/cadreen/intelligence/${encodeURIComponent(id)}/promote`, {
      kind,
      name,
    });
  }
}
