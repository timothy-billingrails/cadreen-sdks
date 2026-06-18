import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.policies import PoliciesResource
from cadreen.resources.guardrails import GuardrailsResource
from cadreen.types import (
    EvaluatePolicyResponse,
    CreatePolicyResponse,
    ConfirmPolicyResponse,
    ListPoliciesResponse,
    PolicyBundle,
    GovernanceDecision,
)


GUARDRAILS_FIXTURES = {
    "POST /api/v1/cadreen/policies/evaluate": {
        "action": "delete_user",
        "domain": "security",
        "result": {
            "type": "handoff",
            "confidence": 0.6,
            "reason": "User deletion requires human approval",
        },
    },
    "POST /api/v1/cadreen/policies": {
        "id": "pol_guard",
        "name": "Data Deletion Guard",
        "version": 1,
        "status": "active",
        "confirmation_required": True,
        "approve_url": "https://app.example.com/policies/pol_guard/approve",
    },
    "POST /api/v1/cadreen/policies/pol_confirm/confirm": {
        "id": "pol_confirm",
        "version": 2,
        "status": "active",
        "previous_version": 1,
        "already_active": False,
        "confirmed_at": "2026-06-17T12:00:00Z",
    },
    "GET /api/v1/cadreen/policies": {
        "policies": [
            {
                "id": "pol_g1",
                "name": "Data Retention",
                "domain": "compliance",
                "priority": 1,
                "requires_human": True,
                "approver_role": "compliance_officer",
                "rationale": "GDPR Article 5",
            }
        ],
        "version": 1,
    },
    "GET /api/v1/cadreen/policies/bundle_g": {
        "id": "bundle_g",
        "version": 1,
        "name": "Security Bundle",
        "policies": [
            {
                "id": "pol_g1",
                "name": "Data Retention",
                "domain": "compliance",
                "priority": 1,
                "requires_human": True,
            }
        ],
        "created_at": "2026-06-01T00:00:00Z",
    },
}


@pytest.fixture
def guardrails_policies_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=GUARDRAILS_FIXTURES)
    return HttpClient(config)


class TestGuardrailsResource:
    @pytest.mark.asyncio
    async def test_check_delegates_to_evaluate(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.check("delete_user", domain="security")
        assert isinstance(result, EvaluatePolicyResponse)
        assert result.action == "delete_user"
        assert result.domain == "security"
        assert isinstance(result.result, GovernanceDecision)
        assert result.result.type == "handoff"
        assert result.result.reason == "User deletion requires human approval"

    @pytest.mark.asyncio
    async def test_check_with_context(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.check(
            "delete_user",
            domain="security",
            context={"user_id": "user_abc", "role": "admin"},
        )
        assert result.result.type == "handoff"

    @pytest.mark.asyncio
    async def test_add_delegates_to_create(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.add("Data Deletion Guard", domain="security")
        assert isinstance(result, CreatePolicyResponse)
        assert result.id == "pol_guard"
        assert result.name == "Data Deletion Guard"
        assert result.status == "active"
        assert result.confirmation_required is True

    @pytest.mark.asyncio
    async def test_add_with_rules(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.add(
            "MFA Required",
            rules=[{"action": "*", "require": "mfa"}],
            domain="auth",
            auto_draft=False,
        )
        assert result.id == "pol_guard"

    @pytest.mark.asyncio
    async def test_require_approval(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.require_approval("Require VP sign-off for data deletion")
        assert isinstance(result, CreatePolicyResponse)
        assert result.id == "pol_guard"

    @pytest.mark.asyncio
    async def test_approve_delegates_to_confirm(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.approve("pol_confirm")
        assert isinstance(result, ConfirmPolicyResponse)
        assert result.id == "pol_confirm"
        assert result.version == 2
        assert result.status == "active"
        assert result.already_active is False

    @pytest.mark.asyncio
    async def test_list(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.list()
        assert isinstance(result, ListPoliciesResponse)
        assert len(result.policies) == 1
        assert result.policies[0].name == "Data Retention"
        assert result.policies[0].requires_human is True
        assert result.policies[0].approver_role == "compliance_officer"

    @pytest.mark.asyncio
    async def test_get(self, guardrails_policies_client):
        policies = PoliciesResource(guardrails_policies_client)
        guardrails = GuardrailsResource(policies)
        result = await guardrails.get("bundle_g")
        assert isinstance(result, PolicyBundle)
        assert result.id == "bundle_g"
        assert result.name == "Security Bundle"
        assert len(result.policies) == 1
        assert result.policies[0].name == "Data Retention"
