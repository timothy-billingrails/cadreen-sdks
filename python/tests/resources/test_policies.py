import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.policies import PoliciesResource
from cadreen.types import (
    CreatePolicyResponse,
    ConfirmPolicyResponse,
    EvaluatePolicyResponse,
    ListPoliciesResponse,
    PolicyBundle,
    Policy,
    GovernanceDecision,
)


POLICIES_FIXTURES = {
    "POST /api/v1/cadreen/policies": {
        "id": "pol_abc",
        "name": "Budget Approval",
        "version": 1,
        "status": "active",
        "confirmation_required": True,
        "approve_url": "https://app.example.com/policies/pol_abc/approve",
    },
    "POST /api/v1/cadreen/policies/evaluate": {
        "action": "transfer_funds",
        "domain": "finance",
        "result": {
            "type": "auto",
            "confidence": 0.95,
            "reason": "Amount within authorized limit",
        },
    },
    "POST /api/v1/cadreen/policies/pol_123/confirm": {
        "id": "pol_123",
        "version": 2,
        "status": "active",
        "previous_version": 1,
        "already_active": False,
        "confirmed_at": "2026-06-17T10:00:00Z",
    },
    "GET /api/v1/cadreen/policies": {
        "policies": [
            {
                "id": "pol_1",
                "name": "Data Privacy",
                "domain": "compliance",
                "priority": 1,
                "requires_human": True,
                "approver_role": "admin",
                "sla_hours": 24,
                "rationale": "GDPR compliance required",
            },
            {
                "id": "pol_2",
                "name": "Spending Limit",
                "domain": "finance",
                "priority": 2,
                "requires_human": False,
                "rationale": "Prevent overspend",
            },
        ],
        "version": 3,
    },
    "GET /api/v1/cadreen/policies/bundle_1": {
        "id": "bundle_1",
        "version": 1,
        "name": "Finance Bundle",
        "policies": [
            {
                "id": "pol_a",
                "name": "Spending Limit",
                "domain": "finance",
                "priority": 1,
                "requires_human": False,
            }
        ],
        "created_at": "2026-06-01T00:00:00Z",
    },
}


@pytest.fixture
def policies_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=POLICIES_FIXTURES)
    return HttpClient(config)


class TestPoliciesResource:
    @pytest.mark.asyncio
    async def test_create(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.create("Budget Approval", domain="finance", auto_draft=False)
        assert isinstance(result, CreatePolicyResponse)
        assert result.id == "pol_abc"
        assert result.name == "Budget Approval"
        assert result.version == 1
        assert result.status == "active"
        assert result.confirmation_required is True
        assert result.approve_url == "https://app.example.com/policies/pol_abc/approve"

    @pytest.mark.asyncio
    async def test_create_with_rules(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.create(
            "Security Policy",
            rules=[{"action": "elevate", "require": "mfa"}],
            domain="security",
        )
        assert result.id == "pol_abc"

    @pytest.mark.asyncio
    async def test_evaluate(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.evaluate("transfer_funds", domain="finance")
        assert isinstance(result, EvaluatePolicyResponse)
        assert result.action == "transfer_funds"
        assert result.domain == "finance"
        assert isinstance(result.result, GovernanceDecision)
        assert result.result.type == "auto"
        assert result.result.confidence == 0.95
        assert result.result.reason == "Amount within authorized limit"

    @pytest.mark.asyncio
    async def test_evaluate_with_context(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.evaluate(
            "transfer_funds",
            domain="finance",
            context={"amount": 500, "recipient": "vendor_xyz"},
        )
        assert result.result.type == "auto"

    @pytest.mark.asyncio
    async def test_confirm(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.confirm("pol_123")
        assert isinstance(result, ConfirmPolicyResponse)
        assert result.id == "pol_123"
        assert result.version == 2
        assert result.status == "active"
        assert result.previous_version == 1
        assert result.already_active is False
        assert result.confirmed_at == "2026-06-17T10:00:00Z"

    @pytest.mark.asyncio
    async def test_list(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.list()
        assert isinstance(result, ListPoliciesResponse)
        assert len(result.policies) == 2
        assert result.policies[0].name == "Data Privacy"
        assert result.policies[0].requires_human is True
        assert result.policies[0].approver_role == "admin"
        assert result.policies[0].sla_hours == 24

    @pytest.mark.asyncio
    async def test_get(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.get("bundle_1")
        assert isinstance(result, PolicyBundle)
        assert result.id == "bundle_1"
        assert result.name == "Finance Bundle"
        assert len(result.policies) == 1
        assert result.policies[0].name == "Spending Limit"

    @pytest.mark.asyncio
    async def test_require_approval(self, policies_client):
        resource = PoliciesResource(policies_client)
        result = await resource.require_approval("Require manager sign-off for refunds")
        assert isinstance(result, CreatePolicyResponse)
        assert result.id == "pol_abc"
        assert result.name == "Budget Approval"

    @pytest.mark.asyncio
    async def test_list_empty(self):
        """List when policies array is empty"""
        fixtures = {"GET /api/v1/cadreen/policies": {"policies": []}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = PoliciesResource(client)
        result = await resource.list()
        assert isinstance(result, ListPoliciesResponse)
        assert len(result.policies) == 0

    @pytest.mark.asyncio
    async def test_evaluate_abstain_result(self):
        """Evaluate with abstain governance decision"""
        fixtures = {
            "POST /api/v1/cadreen/policies/evaluate": {
                "action": "unknown_action",
                "domain": "mystery",
                "result": {"type": "abstain", "confidence": 0.0, "reason": "No matching policy"},
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = PoliciesResource(client)
        result = await resource.evaluate("unknown_action")
        assert result.result.type == "abstain"
        assert result.result.confidence == 0.0
