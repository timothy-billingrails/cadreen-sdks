import { HttpClient } from "../client";

export interface Blueprint {
  id: string;
  name: string;
  description?: string;
  status: string;
  version: number;
  intent?: string;
  source_type?: string;
  source_id?: string;
  created_at: string;
  updated_at: string;
}

export interface BlueprintRun {
  id: string;
  blueprint_id: string;
  blueprint_version: number;
  status: string;
  params?: Record<string, unknown>;
  result_summary?: string;
  trace_id?: string;
  created_at: string;
}

export interface CreateBlueprintRequest {
  name: string;
  description?: string;
  source?: { type: string; trace_id?: string; execution_id?: string };
  parameter_schema?: Record<string, unknown>;
  default_params?: Record<string, unknown>;
}

export interface ListBlueprintsResponse {
  blueprints: Blueprint[];
  count: number;
}

export interface ListBlueprintRunsResponse {
  runs: BlueprintRun[];
  count: number;
}

export class BlueprintsResource {
  constructor(private client: HttpClient) {}

  async list(options?: { status?: string; limit?: number; offset?: number }): Promise<ListBlueprintsResponse> {
    const params = new URLSearchParams();
    if (options?.status) params.set("status", options.status);
    if (options?.limit) params.set("limit", String(options.limit));
    if (options?.offset) params.set("offset", String(options.offset));
    const qs = params.toString();
    return this.client.get<ListBlueprintsResponse>(`/api/v1/cadreen/blueprints${qs ? `?${qs}` : ""}`);
  }

  async get(id: string): Promise<Blueprint> {
    return this.client.get<Blueprint>(`/api/v1/cadreen/blueprints/${encodeURIComponent(id)}`);
  }

  async create(request: CreateBlueprintRequest): Promise<Blueprint> {
    return this.client.post<Blueprint>("/api/v1/cadreen/blueprints", request);
  }

  async delete(id: string): Promise<void> {
    return this.client.delete(`/api/v1/cadreen/blueprints/${encodeURIComponent(id)}`);
  }

  async run(id: string, params?: Record<string, unknown>): Promise<BlueprintRun> {
    return this.client.post<BlueprintRun>(`/api/v1/cadreen/blueprints/${encodeURIComponent(id)}/runs`, params ? { params } : {});
  }

  async listRuns(id: string, options?: { limit?: number; offset?: number }): Promise<ListBlueprintRunsResponse> {
    const params = new URLSearchParams();
    if (options?.limit) params.set("limit", String(options.limit));
    if (options?.offset) params.set("offset", String(options.offset));
    const qs = params.toString();
    return this.client.get<ListBlueprintRunsResponse>(`/api/v1/cadreen/blueprints/${encodeURIComponent(id)}/runs${qs ? `?${qs}` : ""}`);
  }
}
