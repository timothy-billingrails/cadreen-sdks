import type {
  CreatePolicyRequest,
  CreatePolicyResponse,
  ConfirmPolicyResponse,
  EvaluatePolicyRequest,
  EvaluatePolicyResponse,
  ListPoliciesResponse,
  PolicyBundle,
} from "../types";
import { HttpClient } from "../client";

export class PoliciesResource {
  constructor(private client: HttpClient) {}

  async create(request: CreatePolicyRequest): Promise<CreatePolicyResponse> {
    return this.client.post<CreatePolicyResponse>("/api/v1/cadreen/policies", request);
  }

  async evaluate(request: EvaluatePolicyRequest): Promise<EvaluatePolicyResponse> {
    return this.client.post<EvaluatePolicyResponse>("/api/v1/cadreen/policies/evaluate", request);
  }

  async confirm(id: string): Promise<ConfirmPolicyResponse> {
    return this.client.post<ConfirmPolicyResponse>(`/api/v1/cadreen/policies/${encodeURIComponent(id)}/confirm`);
  }

  async list(): Promise<ListPoliciesResponse> {
    return this.client.get<ListPoliciesResponse>("/api/v1/cadreen/policies");
  }

  async get(id: string): Promise<PolicyBundle> {
    return this.client.get<PolicyBundle>(`/api/v1/cadreen/policies/${encodeURIComponent(id)}`);
  }

  async requireApproval(description: string): Promise<CreatePolicyResponse> {
    return this.create({
      name: description,
      auto_draft: true,
    });
  }
}
