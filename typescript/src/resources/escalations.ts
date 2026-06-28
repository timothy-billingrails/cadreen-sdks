import type {
  ListEscalationsResponse,
  Escalation,
  ResolveEscalationRequest,
} from "../types";
import { HttpClient } from "../client";

export class EscalationsResource {
  constructor(private client: HttpClient) {}

  async list(): Promise<ListEscalationsResponse> {
    return this.client.get<ListEscalationsResponse>("/api/v1/cadreen/escalations");
  }

  async get(id: string): Promise<Escalation> {
    return this.client.get<Escalation>(`/api/v1/cadreen/escalations/${encodeURIComponent(id)}`);
  }

  async resolve(id: string, decision: string): Promise<Escalation> {
    return this.client.post<Escalation>(
      `/api/v1/cadreen/escalations/${encodeURIComponent(id)}/resolve`,
      { decision } as ResolveEscalationRequest
    );
  }
}
