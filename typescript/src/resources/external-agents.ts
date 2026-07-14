import type {
  ExternalAgentConnection,
  ExternalAgentSettings,
  ListExternalConnectionsResponse,
  ListExternalInteractionsResponse,
} from "../types";
import { HttpClient } from "../client";

export class ExternalAgentsResource {
  constructor(private client: HttpClient) {}

  /**
   * Connect to an external A2A agent by providing its Agent Card URL.
   * The connection starts in pending_approval status and must be approved.
   */
  async connect(agentId: string, agentCardUrl: string): Promise<ExternalAgentConnection> {
    return this.client.post<ExternalAgentConnection>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external`,
      { agentCardUrl }
    );
  }

  /**
   * List external agent connections for an agent.
   */
  async list(
    agentId: string,
    params?: { status?: string; limit?: number; offset?: number }
  ): Promise<ListExternalConnectionsResponse> {
    return this.client.get<ListExternalConnectionsResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  /**
   * Get a specific external agent connection.
   */
  async get(agentId: string, connectionId: string): Promise<ExternalAgentConnection> {
    return this.client.get<ExternalAgentConnection>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external/${encodeURIComponent(connectionId)}`
    );
  }

  /**
   * Approve a pending external agent connection.
   */
  async approve(agentId: string, connectionId: string): Promise<{ status: string }> {
    return this.client.post<{ status: string }>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external/${encodeURIComponent(connectionId)}/approve`
    );
  }

  /**
   * Suspend an active external agent connection.
   */
  async suspend(agentId: string, connectionId: string): Promise<{ status: string }> {
    return this.client.post<{ status: string }>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external/${encodeURIComponent(connectionId)}/suspend`
    );
  }

  /**
   * Revoke an external agent connection (permanent).
   */
  async revoke(agentId: string, connectionId: string): Promise<{ status: string }> {
    return this.client.post<{ status: string }>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external/${encodeURIComponent(connectionId)}/revoke`
    );
  }

  /**
   * Delete an external agent connection.
   */
  async delete(agentId: string, connectionId: string): Promise<void> {
    return this.client.delete<void>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external/${encodeURIComponent(connectionId)}`
    );
  }

  /**
   * List interactions for an external agent connection.
   */
  async listInteractions(
    agentId: string,
    connectionId: string,
    params?: { limit?: number; offset?: number }
  ): Promise<ListExternalInteractionsResponse> {
    return this.client.get<ListExternalInteractionsResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/external/${encodeURIComponent(connectionId)}/interactions`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  /**
   * Get external agent settings for the workspace.
   */
  async getSettings(): Promise<ExternalAgentSettings> {
    return this.client.get<ExternalAgentSettings>("/api/v1/cadreen/external-agents/settings");
  }

  /**
   * Enable or disable external agents for the workspace.
   */
  async updateSettings(enabled: boolean): Promise<ExternalAgentSettings> {
    return this.client.put<ExternalAgentSettings>("/api/v1/cadreen/external-agents/settings", { enabled });
  }

  /**
   * List all external agent connections in the workspace.
   */
  async listAll(): Promise<ListExternalConnectionsResponse> {
    return this.client.get<ListExternalConnectionsResponse>("/api/v1/cadreen/external-agents/connections");
  }
}
