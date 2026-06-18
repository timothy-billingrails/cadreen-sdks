import { describe, it, expect } from "vitest";
import { PoliciesResource } from "../../resources/policies";
import { GuardrailsResource } from "../../resources/guardrails";
import { HttpClient } from "../../client";
import type {
  CreatePolicyResponse,
  EvaluatePolicyResponse,
  ConfirmPolicyResponse,
  ListPoliciesResponse,
} from "../../types";

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("PoliciesResource create()", () => {
  it("returns fixture data", async () => {
    const fixture: CreatePolicyResponse = {
      id: "pol-1",
      name: "Test Policy",
      version: 1,
      status: "active",
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies": fixture });
    const resource = new PoliciesResource(client);
    const result = await resource.create({ name: "Test Policy" });
    expect(result.id).toBe("pol-1");
    expect(result.status).toBe("active");
  });
});

describe("PoliciesResource evaluate()", () => {
  it("returns fixture data", async () => {
    const fixture: EvaluatePolicyResponse = {
      action: "send_email",
      domain: "general",
      result: { type: "auto", confidence: 95, reason: "Safe action" },
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies/evaluate": fixture });
    const resource = new PoliciesResource(client);
    const result = await resource.evaluate({ action: "send_email" });
    expect(result.action).toBe("send_email");
    expect(result.result.type).toBe("auto");
  });
});

describe("PoliciesResource confirm()", () => {
  it("returns fixture data", async () => {
    const fixture: ConfirmPolicyResponse = {
      id: "pol-2",
      version: 2,
      previous_version: 1,
      status: "confirmed",
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies/pol-2/confirm": fixture });
    const resource = new PoliciesResource(client);
    const result = await resource.confirm("pol-2");
    expect(result.id).toBe("pol-2");
    expect(result.status).toBe("confirmed");
  });
});

describe("PoliciesResource list()", () => {
  it("returns fixture data", async () => {
    const fixture: ListPoliciesResponse = {
      policies: [
        { id: "p1", name: "Policy 1", domain: "general", priority: 1, requires_human: false },
      ],
    };
    const client = buildClient({ "GET /api/v1/cadreen/policies": fixture });
    const resource = new PoliciesResource(client);
    const result = await resource.list();
    expect(result.policies).toHaveLength(1);
    expect(result.policies[0].name).toBe("Policy 1");
  });
});

describe("PoliciesResource requireApproval()", () => {
  it("delegates to create with auto_draft", async () => {
    const fixture: CreatePolicyResponse = {
      id: "pol-3",
      name: "Approval Needed",
      version: 1,
      status: "pending_approval",
      confirmation_required: true,
      approve_url: "https://example.com/approve",
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies": fixture });
    const resource = new PoliciesResource(client);
    const result = await resource.requireApproval("Approval Needed");
    expect(result.id).toBe("pol-3");
    expect(result.confirmation_required).toBe(true);
  });
});

describe("GuardrailsResource delegates to policies", () => {
  it("check() delegates to evaluate()", async () => {
    const fixture: EvaluatePolicyResponse = {
      action: "test",
      domain: "general",
      result: { type: "handoff", confidence: 50, reason: "Needs review" },
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies/evaluate": fixture });
    const policies = new PoliciesResource(client);
    const guardrails = new GuardrailsResource(policies);
    const result = await guardrails.check({ action: "test" });
    expect(result.result.type).toBe("handoff");
  });

  it("add() delegates to create()", async () => {
    const fixture: CreatePolicyResponse = {
      id: "pol-4",
      name: "New Rule",
      version: 1,
      status: "active",
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies": fixture });
    const policies = new PoliciesResource(client);
    const guardrails = new GuardrailsResource(policies);
    const result = await guardrails.add({ name: "New Rule" });
    expect(result.id).toBe("pol-4");
  });

  it("approve() delegates to confirm()", async () => {
    const fixture: ConfirmPolicyResponse = {
      id: "pol-5",
      version: 1,
      status: "confirmed",
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies/pol-5/confirm": fixture });
    const policies = new PoliciesResource(client);
    const guardrails = new GuardrailsResource(policies);
    const result = await guardrails.approve("pol-5");
    expect(result.status).toBe("confirmed");
  });

  it("requireApproval() delegates correctly", async () => {
    const fixture: CreatePolicyResponse = {
      id: "pol-6",
      name: "Need Approval",
      version: 1,
      status: "pending",
      confirmation_required: true,
    };
    const client = buildClient({ "POST /api/v1/cadreen/policies": fixture });
    const policies = new PoliciesResource(client);
    const guardrails = new GuardrailsResource(policies);
    const result = await guardrails.requireApproval("Need Approval");
    expect(result.id).toBe("pol-6");
  });

  it("list() delegates to list()", async () => {
    const fixture: ListPoliciesResponse = {
      policies: [{ id: "p2", name: "P2", domain: "general", priority: 1, requires_human: true }],
    };
    const client = buildClient({ "GET /api/v1/cadreen/policies": fixture });
    const policies = new PoliciesResource(client);
    const guardrails = new GuardrailsResource(policies);
    const result = await guardrails.list();
    expect(result.policies[0].requires_human).toBe(true);
  });
});
