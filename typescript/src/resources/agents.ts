import type {
  Agent,
  AgentConfig,
  AgentCapabilities,
  AgentKnowledge,
  AgentGovernancePolicy,
  AgentNegotiation,
  AgentMessage,
  CreateAgentRequest,
  UpdateAgentRequest,
  SendMessageRequest,
  CreateKnowledgeRequest,
  SearchKnowledgeRequest,
  CreateGovernanceRequest,
  UpdateGovernanceRequest,
  StartNegotiationRequest,
  RespondNegotiationRequest,
  ListAgentsResponse,
  ListAgentMessagesResponse,
  ListAgentExecutionsResponse,
  ListAgentKnowledgeResponse,
  ListAgentGovernanceResponse,
  ListAgentAuditResponse,
  ListAgentNegotiationsResponse,
  AgentExecution,
  CreateExecutionRequest,
} from "../types";
import { HttpClient } from "../client";

export class AgentsResource {
  constructor(private client: HttpClient) {}

  async create(request: CreateAgentRequest): Promise<Agent> {
    return this.client.post<Agent>("/api/v1/cadreen/agents", request);
  }

  async list(params?: { search?: string; limit?: number; offset?: number }): Promise<ListAgentsResponse> {
    return this.client.get<ListAgentsResponse>("/api/v1/cadreen/agents", params as Record<string, string | number | boolean | undefined>);
  }

  async get(agentId: string): Promise<Agent> {
    return this.client.get<Agent>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}`);
  }

  async update(agentId: string, request: UpdateAgentRequest): Promise<Agent> {
    return this.client.patch<Agent>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}`, request);
  }

  async delete(agentId: string): Promise<void> {
    return this.client.delete<void>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}`);
  }

  async getConfig(agentId: string): Promise<AgentConfig> {
    return this.client.get<AgentConfig>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/config`);
  }

  async deploy(agentId: string, request: { configSnapshot: unknown; changeSummary?: string }): Promise<Agent> {
    return this.client.post<Agent>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/deploy`, request);
  }

  async getCapabilities(agentId: string): Promise<AgentCapabilities> {
    return this.client.get<AgentCapabilities>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/capabilities`);
  }

  async sendMessage(agentId: string, request: SendMessageRequest): Promise<AgentMessage> {
    return this.client.post<AgentMessage>(`/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/send`, request);
  }

  async listMessages(agentId: string, params?: { limit?: number; offset?: number }): Promise<ListAgentMessagesResponse> {
    return this.client.get<ListAgentMessagesResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/messages`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  async listExecutions(agentId: string, params?: { limit?: number; offset?: number }): Promise<ListAgentExecutionsResponse> {
    return this.client.get<ListAgentExecutionsResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/executions`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  async createExecution(agentId: string, request: CreateExecutionRequest): Promise<AgentExecution> {
    return this.client.post<AgentExecution>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/executions`,
      request
    );
  }

  async listKnowledge(agentId: string, params?: { limit?: number; offset?: number }): Promise<ListAgentKnowledgeResponse> {
    return this.client.get<ListAgentKnowledgeResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/knowledge`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  async createKnowledge(agentId: string, request: CreateKnowledgeRequest): Promise<AgentKnowledge> {
    return this.client.post<AgentKnowledge>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/knowledge`,
      request
    );
  }

  async searchKnowledge(agentId: string, request: SearchKnowledgeRequest): Promise<ListAgentKnowledgeResponse> {
    return this.client.post<ListAgentKnowledgeResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/knowledge/search`,
      request
    );
  }

  async deleteKnowledge(agentId: string, knowledgeId: string): Promise<void> {
    return this.client.delete<void>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/knowledge/${encodeURIComponent(knowledgeId)}`
    );
  }

  async listGovernance(agentId: string): Promise<ListAgentGovernanceResponse> {
    return this.client.get<ListAgentGovernanceResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/governance`
    );
  }

  async createGovernance(agentId: string, request: CreateGovernanceRequest): Promise<AgentGovernancePolicy> {
    return this.client.post<AgentGovernancePolicy>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/governance`,
      request
    );
  }

  async updateGovernance(agentId: string, policyId: string, request: UpdateGovernanceRequest): Promise<AgentGovernancePolicy> {
    return this.client.patch<AgentGovernancePolicy>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/governance/${encodeURIComponent(policyId)}`,
      request
    );
  }

  async deleteGovernance(agentId: string, policyId: string): Promise<void> {
    return this.client.delete<void>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/governance/${encodeURIComponent(policyId)}`
    );
  }

  async listAudit(agentId: string, params?: { limit?: number; offset?: number }): Promise<ListAgentAuditResponse> {
    return this.client.get<ListAgentAuditResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/audit`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  async startNegotiation(agentId: string, request: StartNegotiationRequest): Promise<AgentNegotiation> {
    return this.client.post<AgentNegotiation>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/negotiate`,
      request
    );
  }

  async listNegotiations(agentId: string, params?: { limit?: number; offset?: number }): Promise<ListAgentNegotiationsResponse> {
    return this.client.get<ListAgentNegotiationsResponse>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/negotiations`,
      params as Record<string, string | number | boolean | undefined>
    );
  }

  async getNegotiation(agentId: string, negotiationId: string): Promise<AgentNegotiation> {
    return this.client.get<AgentNegotiation>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/negotiations/${encodeURIComponent(negotiationId)}`
    );
  }

  async respondToNegotiation(agentId: string, negotiationId: string, request: RespondNegotiationRequest): Promise<AgentNegotiation> {
    return this.client.post<AgentNegotiation>(
      `/api/v1/cadreen/agents/${encodeURIComponent(agentId)}/negotiations/${encodeURIComponent(negotiationId)}/respond`,
      request
    );
  }
}
