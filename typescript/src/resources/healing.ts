import type {
  HealingStatsResponse,
  ListHealingPrecedentsResponse,
  HealingDiagnosis,
  DiagnoseRequest,
} from "../types";
import { HttpClient } from "../client";

export class HealingResource {
  constructor(private client: HttpClient) {}

  async stats(): Promise<HealingStatsResponse> {
    return this.client.get<HealingStatsResponse>("/api/v1/cadreen/healing/stats");
  }

  async precedents(): Promise<ListHealingPrecedentsResponse> {
    return this.client.get<ListHealingPrecedentsResponse>("/api/v1/cadreen/healing/precedents");
  }

  async diagnose(request: DiagnoseRequest): Promise<HealingDiagnosis> {
    return this.client.post<HealingDiagnosis>("/api/v1/cadreen/healing/diagnose", request);
  }
}
