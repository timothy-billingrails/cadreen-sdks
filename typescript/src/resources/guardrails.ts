import type {
  CreatePolicyRequest,
  CreatePolicyResponse,
  ConfirmPolicyResponse,
  EvaluatePolicyRequest,
  EvaluatePolicyResponse,
  ListPoliciesResponse,
  PolicyBundle,
} from "../types";
import { PoliciesResource } from "./policies";

export class GuardrailsResource {
  constructor(private policies: PoliciesResource) {}

  async check(request: EvaluatePolicyRequest): Promise<EvaluatePolicyResponse> {
    return this.policies.evaluate(request);
  }

  async add(request: CreatePolicyRequest): Promise<CreatePolicyResponse> {
    return this.policies.create(request);
  }

  async requireApproval(description: string): Promise<CreatePolicyResponse> {
    return this.policies.requireApproval(description);
  }

  async approve(id: string): Promise<ConfirmPolicyResponse> {
    return this.policies.confirm(id);
  }

  async list(): Promise<ListPoliciesResponse> {
    return this.policies.list();
  }

  async get(id: string): Promise<PolicyBundle> {
    return this.policies.get(id);
  }
}
