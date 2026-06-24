import type {
  TaskProposal,
  ListProposalsResponse,
  AcceptProposalResponse,
  DismissProposalResponse,
  ProposalStatsResponse,
  ListProposalsOptions,
} from "../types";
import { HttpClient } from "../client";

export class ProposalsResource {
  constructor(private client: HttpClient) {}

  async list(options?: ListProposalsOptions): Promise<ListProposalsResponse> {
    const params = new URLSearchParams();
    if (options?.status) params.set("status", options.status);
    if (options?.limit) params.set("limit", String(options.limit));
    const qs = params.toString();
    const path = qs ? `/api/v1/cadreen/proposals?${qs}` : "/api/v1/cadreen/proposals";
    return this.client.get<ListProposalsResponse>(path);
  }

  async get(id: string): Promise<TaskProposal> {
    return this.client.get<TaskProposal>(
      `/api/v1/cadreen/proposals/${encodeURIComponent(id)}`
    );
  }

  async accept(id: string): Promise<AcceptProposalResponse> {
    return this.client.post<AcceptProposalResponse>(
      `/api/v1/cadreen/proposals/${encodeURIComponent(id)}/accept`
    );
  }

  async dismiss(id: string, reason?: string): Promise<DismissProposalResponse> {
    return this.client.post<DismissProposalResponse>(
      `/api/v1/cadreen/proposals/${encodeURIComponent(id)}/dismiss`,
      reason ? { reason } : undefined
    );
  }

  async stats(): Promise<ProposalStatsResponse> {
    return this.client.get<ProposalStatsResponse>("/api/v1/cadreen/proposals/stats");
  }
}
