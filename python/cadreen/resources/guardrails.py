from __future__ import annotations

from typing import Any

from ..types import (
    CreatePolicyResponse,
    ConfirmPolicyResponse,
    EvaluatePolicyResponse,
    ListPoliciesResponse,
    PolicyBundle,
)
from .policies import PoliciesResource


class GuardrailsResource:
    def __init__(self, policies: PoliciesResource) -> None:
        self._policies = policies

    async def check(
        self,
        action: str,
        *,
        domain: str | None = None,
        context: dict[str, object] | None = None,
    ) -> EvaluatePolicyResponse:
        return await self._policies.evaluate(action, domain=domain, context=context)

    async def add(
        self,
        name: str,
        *,
        rules: list[dict[str, object]] | None = None,
        domain: str | None = None,
        auto_draft: bool | None = None,
    ) -> CreatePolicyResponse:
        return await self._policies.create(name, rules=rules, domain=domain, auto_draft=auto_draft)

    async def require_approval(self, description: str) -> CreatePolicyResponse:
        return await self._policies.require_approval(description)

    async def approve(self, id: str) -> ConfirmPolicyResponse:
        return await self._policies.confirm(id)

    async def list(self) -> ListPoliciesResponse:
        return await self._policies.list()

    async def get(self, id: str) -> PolicyBundle:
        return await self._policies.get(id)
