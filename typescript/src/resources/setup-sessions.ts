import type {
  SetupSession,
  SetupSessionCreateRequest,
  SetupSessionAddRequest,
  SetupSessionApplyRequest,
  SetupSessionApplyResult,
} from "../types";
import { HttpClient } from "../client";

export class SetupSessionsResource {
  constructor(private client: HttpClient) {}

  async create(request: SetupSessionCreateRequest): Promise<SetupSession> {
    return this.client.post<SetupSession>(
      "/api/v1/cadreen/setup/sessions",
      request
    );
  }

  async list(): Promise<{ sessions: SetupSession[] }> {
    return this.client.get<{ sessions: SetupSession[] }>(
      "/api/v1/cadreen/setup/sessions"
    );
  }

  async get(id: string): Promise<SetupSession> {
    return this.client.get<SetupSession>(
      `/api/v1/cadreen/setup/sessions/${encodeURIComponent(id)}`
    );
  }

  async addResources(
    id: string,
    request: SetupSessionAddRequest
  ): Promise<SetupSession> {
    return this.client.post<SetupSession>(
      `/api/v1/cadreen/setup/sessions/${encodeURIComponent(id)}`,
      request
    );
  }

  async apply(
    id: string,
    request: SetupSessionApplyRequest
  ): Promise<SetupSessionApplyResult> {
    return this.client.post<SetupSessionApplyResult>(
      `/api/v1/cadreen/setup/sessions/${encodeURIComponent(id)}/apply`,
      request
    );
  }
}
