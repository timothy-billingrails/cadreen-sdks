import { HttpClient } from "../client";

export interface Schedule {
  id: string;
  name: string;
  blueprint_id: string;
  blueprint_version: number;
  status: string;
  trigger: { type: string; frequency?: string; time?: string; expression?: string; weekdays?: string[] };
  timezone: string;
  params?: Record<string, unknown>;
  next_run_at?: string;
  last_run_at?: string;
  pause_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateScheduleRequest {
  blueprint_id: string;
  name: string;
  trigger: { type: string; frequency?: string; time?: string; expression?: string };
  timezone?: string;
  params?: Record<string, unknown>;
}

export interface ListSchedulesResponse {
  schedules: Schedule[];
  count: number;
}

export class SchedulesResource {
  constructor(private client: HttpClient) {}

  async list(): Promise<ListSchedulesResponse> {
    return this.client.get<ListSchedulesResponse>("/api/v1/cadreen/schedules");
  }

  async get(id: string): Promise<Schedule> {
    return this.client.get<Schedule>(`/api/v1/cadreen/schedules/${encodeURIComponent(id)}`);
  }

  async create(request: CreateScheduleRequest): Promise<Schedule> {
    return this.client.post<Schedule>("/api/v1/cadreen/schedules", request);
  }

  async pause(id: string, reason?: string): Promise<{ id: string; status: string }> {
    return this.client.post<{ id: string; status: string }>(`/api/v1/cadreen/schedules/${encodeURIComponent(id)}/pause`, reason ? { reason } : {});
  }

  async resume(id: string): Promise<{ id: string; status: string; next_run_at?: string }> {
    return this.client.post<{ id: string; status: string; next_run_at?: string }>(`/api/v1/cadreen/schedules/${encodeURIComponent(id)}/resume`);
  }
}
