from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    CreatePolicyResponse,
    ConfirmPolicyResponse,
    EvaluatePolicyResponse,
    ListPoliciesResponse,
    PolicyBundle,
    Policy,
    GovernanceDecision,
)


class PoliciesResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def create(
        self,
        name: str,
        *,
        rules: list[dict[str, object]] | None = None,
        domain: str | None = None,
        auto_draft: bool | None = None,
    ) -> CreatePolicyResponse:
        body: dict[str, Any] = {"name": name}
        if rules is not None:
            body["rules"] = rules
        if domain is not None:
            body["domain"] = domain
        if auto_draft is not None:
            body["auto_draft"] = auto_draft

        raw = await self._client.post("/api/v1/cadreen/policies", body)
        return CreatePolicyResponse(
            id=raw["id"],
            name=raw["name"],
            version=raw.get("version", 0),
            status=raw.get("status", ""),
            confirmation_required=raw.get("confirmation_required"),
            approve_url=raw.get("approve_url"),
        )

    async def evaluate(
        self,
        action: str,
        *,
        domain: str | None = None,
        context: dict[str, object] | None = None,
    ) -> EvaluatePolicyResponse:
        body: dict[str, Any] = {"action": action}
        if domain is not None:
            body["domain"] = domain
        if context is not None:
            body["context"] = context

        raw = await self._client.post("/api/v1/cadreen/policies/evaluate", body)
        gov = raw.get("result", {})
        return EvaluatePolicyResponse(
            action=raw.get("action", ""),
            domain=raw.get("domain", ""),
            result=GovernanceDecision(
                type=gov.get("type", "abstain"),
                confidence=gov.get("confidence", 0.0),
                reason=gov.get("reason", ""),
            ),
        )

    async def confirm(self, id: str) -> ConfirmPolicyResponse:
        raw = await self._client.post(f"/api/v1/cadreen/policies/{id}/confirm")
        return ConfirmPolicyResponse(
            id=raw["id"],
            version=raw.get("version", 0),
            status=raw.get("status", ""),
            previous_version=raw.get("previous_version"),
            already_active=raw.get("already_active"),
            confirmed_at=raw.get("confirmed_at"),
        )

    async def list(self) -> ListPoliciesResponse:
        raw = await self._client.get("/api/v1/cadreen/policies")
        policies = [
            Policy(
                id=p["id"],
                name=p["name"],
                domain=p.get("domain", ""),
                priority=p.get("priority", 0),
                requires_human=p.get("requires_human", False),
                approver_role=p.get("approver_role"),
                sla_hours=p.get("sla_hours"),
                rationale=p.get("rationale"),
            )
            for p in raw.get("policies", [])
        ]
        return ListPoliciesResponse(
            policies=policies,
            version=raw.get("version"),
        )

    async def get(self, id: str) -> PolicyBundle:
        raw = await self._client.get(f"/api/v1/cadreen/policies/{id}")
        policies = [
            Policy(
                id=p["id"],
                name=p["name"],
                domain=p.get("domain", ""),
                priority=p.get("priority", 0),
                requires_human=p.get("requires_human", False),
                approver_role=p.get("approver_role"),
                sla_hours=p.get("sla_hours"),
                rationale=p.get("rationale"),
            )
            for p in raw.get("policies", [])
        ]
        return PolicyBundle(
            id=raw["id"],
            version=raw.get("version", 0),
            name=raw.get("name", ""),
            policies=policies,
            created_at=raw.get("created_at"),
        )

    async def require_approval(self, description: str) -> CreatePolicyResponse:
        return await self.create(description, auto_draft=True)
