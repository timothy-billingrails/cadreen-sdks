import type {
  FederationLink,
  FederationAgent,
  FederationPermissions,
  CreateFederationRequest,
  SuspendFederationRequest,
  RevokeFederationRequest,
  UpdatePermissionsRequest,
  LinkAgentRequest,
  ListFederationResponse,
  ListFederationAgentsResponse,
} from "../types";
import { HttpClient } from "../client";

export class FederationResource {
  constructor(private client: HttpClient) {}

  async create(request: CreateFederationRequest): Promise<FederationLink> {
    return this.client.post<FederationLink>("/api/v1/cadreen/federation", request);
  }

  async list(): Promise<ListFederationResponse> {
    return this.client.get<ListFederationResponse>("/api/v1/cadreen/federation");
  }

  async get(federationId: string): Promise<FederationLink> {
    return this.client.get<FederationLink>(`/api/v1/cadreen/federation/${encodeURIComponent(federationId)}`);
  }

  async approve(federationId: string): Promise<FederationLink> {
    return this.client.post<FederationLink>(`/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/approve`);
  }

  async suspend(federationId: string, request?: SuspendFederationRequest): Promise<FederationLink> {
    return this.client.post<FederationLink>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/suspend`,
      request
    );
  }

  async revoke(federationId: string, request?: RevokeFederationRequest): Promise<FederationLink> {
    return this.client.post<FederationLink>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/revoke`,
      request
    );
  }

  async getPermissions(federationId: string): Promise<FederationPermissions> {
    return this.client.get<FederationPermissions>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/permissions`
    );
  }

  async updatePermissions(federationId: string, request: UpdatePermissionsRequest): Promise<FederationPermissions> {
    return this.client.put<FederationPermissions>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/permissions`,
      request
    );
  }

  async linkAgent(federationId: string, request: LinkAgentRequest): Promise<FederationAgent> {
    return this.client.post<FederationAgent>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/agents`,
      request
    );
  }

  async listAgents(federationId: string): Promise<ListFederationAgentsResponse> {
    return this.client.get<ListFederationAgentsResponse>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/agents`
    );
  }

  async unlinkAgent(federationId: string, agentLinkId: string): Promise<void> {
    return this.client.delete<void>(
      `/api/v1/cadreen/federation/${encodeURIComponent(federationId)}/agents/${encodeURIComponent(agentLinkId)}`
    );
  }
}
