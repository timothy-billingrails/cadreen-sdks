from __future__ import annotations

from typing import Any

from ..types import (
    CreatePolicyRequest,
    CreatePolicyResponse,
    ConfirmPolicyResponse,
    EvaluatePolicyRequest,
    EvaluatePolicyResponse,
    ListPoliciesResponse,
    PolicyBundle,
)
from .policies import PoliciesResource


class GuardrailsResource:
    def __init__(self, policies: PoliciesResource) -> None:
        self._policies = policies

    async def check(self, request: EvaluatePolicyRequest) -> EvaluatePolicyResponse:
        return await self._policies.evaluate(request)

    async def add(self, request: CreatePolicyRequest) -> CreatePolicyResponse:
        return await self._policies.create(request)

    async def require_approval(self, description: str) -> CreatePolicyResponse:
        return await self._policies.require_approval(description)

    async def approve(self, id: str) -> ConfirmPolicyResponse:
        return await self._policies.confirm(id)

    async def list(self) -> ListPoliciesResponse:
        return await self._policies.list()

    async def get(self, id: str) -> PolicyBundle:
        return await self._policies.get(id)
