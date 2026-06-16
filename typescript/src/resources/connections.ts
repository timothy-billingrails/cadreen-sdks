import type {
  RegisterOpenAPIRequest,
  RegisterOpenAPIResponse,
  RegisterMCPRequest,
  RegisterMCPResponse,
  ListConnectionsResponse,
  InstallComposioRequest,
  InstallComposioResponse,
  SearchComposioResponse,
  ComposioStatusResponse,
  ConnectResult,
  CatalogResponse,
  InstallResponse,
} from "../types";
import { HttpClient } from "../client";

export class ConnectionsResource {
  constructor(private client: HttpClient) {}

  async connect(capability: string): Promise<ConnectResult> {
    return this.client.post<ConnectResult>("/api/v1/cadreen/connections", { capability });
  }

  /** Browse the unified marketplace catalog. */
  async catalog(): Promise<CatalogResponse> {
    return this.client.get<CatalogResponse>("/api/v1/cadreen/connections/catalog");
  }

  /** One-click install an integration from the catalog. Returns OAuth URL for auth. */
  async install(integrationId: string): Promise<InstallResponse> {
    return this.client.post<InstallResponse>("/api/v1/cadreen/connections/install", {
      integration_id: integrationId,
    });
  }

  async registerOpenAPI(request: RegisterOpenAPIRequest): Promise<RegisterOpenAPIResponse> {
    return this.client.post<RegisterOpenAPIResponse>("/api/v1/cadreen/connections/openapi", request);
  }

  async registerMCP(request: RegisterMCPRequest): Promise<RegisterMCPResponse> {
    return this.client.post<RegisterMCPResponse>("/api/v1/cadreen/connections/mcp", request);
  }

  async installComposio(request: InstallComposioRequest): Promise<InstallComposioResponse> {
    return this.client.post<InstallComposioResponse>("/api/v1/cadreen/connections/composio/install", {
      toolkit: request.toolkit,
      user_id: request.user_id,
    });
  }

  async searchComposio(query: string): Promise<SearchComposioResponse> {
    return this.client.post<SearchComposioResponse>("/api/v1/cadreen/connections/composio/search", { query });
  }

  async composioStatus(toolkit?: string, userId?: string): Promise<ComposioStatusResponse> {
    const params: Record<string, string | undefined> = {};
    if (toolkit) params.toolkit = toolkit;
    if (userId) params.user_id = userId;
    return this.client.get<ComposioStatusResponse>("/api/v1/cadreen/connections/composio/status", params);
  }

  async list(): Promise<ListConnectionsResponse> {
    return this.client.get<ListConnectionsResponse>("/api/v1/cadreen/connections");
  }

  async delete(id: string): Promise<void> {
    return this.client.delete<void>(`/api/v1/cadreen/connections/${encodeURIComponent(id)}`);
  }
}
