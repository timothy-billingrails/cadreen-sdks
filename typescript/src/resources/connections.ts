import type {
  RegisterOpenAPIRequest,
  RegisterOpenAPIResponse,
  RegisterMCPRequest,
  RegisterMCPResponse,
  ListConnectionsResponse,
  InstallComposioRequest,
  ConnectResult,
} from "../types";
import { HttpClient } from "../client";

export class ConnectionsResource {
  constructor(private client: HttpClient) {}

  async connect(capability: string): Promise<ConnectResult> {
    return this.client.post<ConnectResult>("/api/v1/cadreen/connections", { capability });
  }

  async registerOpenAPI(request: RegisterOpenAPIRequest): Promise<RegisterOpenAPIResponse> {
    return this.client.post<RegisterOpenAPIResponse>("/api/v1/cadreen/connections/openapi", request);
  }

  async registerMCP(request: RegisterMCPRequest): Promise<RegisterMCPResponse> {
    return this.client.post<RegisterMCPResponse>("/api/v1/cadreen/connections/mcp", request);
  }

  async installComposio(request: InstallComposioRequest): Promise<Record<string, unknown>> {
    return this.client.post<Record<string, unknown>>("/api/v1/cadreen/connections/composio/install", {
      toolkit: request.toolkit,
      user_id: request.user_id,
    });
  }

  async searchComposio(query: string): Promise<Record<string, unknown>> {
    return this.client.post<Record<string, unknown>>("/api/v1/cadreen/connections/composio/search", { query });
  }

  async composioStatus(toolkit?: string, userId?: string): Promise<Record<string, unknown>> {
    const params: Record<string, string | undefined> = {};
    if (toolkit) params.toolkit = toolkit;
    if (userId) params.user_id = userId;
    return this.client.get<Record<string, unknown>>("/api/v1/cadreen/connections/composio/status", params);
  }

  async list(): Promise<ListConnectionsResponse> {
    return this.client.get<ListConnectionsResponse>("/api/v1/cadreen/connections");
  }

  async delete(id: string): Promise<void> {
    return this.client.delete<void>(`/api/v1/cadreen/connections/${encodeURIComponent(id)}`);
  }
}
