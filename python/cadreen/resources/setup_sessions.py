from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    SetupSession,
    SetupSessionApplyResult,
)


class SetupSessionsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def create(
        self,
        *,
        purpose: str | None = None,
        constraints: list[str] | None = None,
        workspace_id: str | None = None,
    ) -> SetupSession:
        body: dict[str, Any] = {}
        if purpose is not None:
            body["purpose"] = purpose
        if constraints is not None:
            body["constraints"] = constraints
        if workspace_id is not None:
            body["workspace_id"] = workspace_id
        raw = await self._client.post("/api/v1/cadreen/setup/sessions", body)
        return self._parse_session(raw)

    async def list(self) -> list[SetupSession]:
        raw = await self._client.get("/api/v1/cadreen/setup/sessions")
        return [self._parse_session(s) for s in raw.get("sessions", [])]

    async def get(self, id: str) -> SetupSession:
        raw = await self._client.get(f"/api/v1/cadreen/setup/sessions/{id}")
        return self._parse_session(raw)

    async def add_resources(
        self,
        id: str,
        *,
        connections: list[dict[str, Any]] | None = None,
        credentials: list[dict[str, Any]] | None = None,
        memory: list[dict[str, Any]] | None = None,
        policies: list[dict[str, Any]] | None = None,
    ) -> SetupSession:
        body: dict[str, Any] = {}
        if connections is not None:
            body["connections"] = connections
        if credentials is not None:
            body["credentials"] = credentials
        if memory is not None:
            body["memory"] = memory
        if policies is not None:
            body["policies"] = policies
        raw = await self._client.post(f"/api/v1/cadreen/setup/sessions/{id}", body)
        return self._parse_session(raw)

    async def apply(self, id: str, *, confirm: bool = True) -> SetupSessionApplyResult:
        raw = await self._client.post(
            f"/api/v1/cadreen/setup/sessions/{id}/apply",
            {"confirm": confirm},
        )
        return SetupSessionApplyResult(
            session_id=raw["session_id"],
            status=raw["status"],
            applied=raw["applied"],
            failed=raw["failed"],
            result=raw.get("result"),
        )

    def _parse_session(self, raw: dict[str, Any]) -> SetupSession:
        return SetupSession(
            id=raw["id"],
            status=raw["status"],
            purpose=raw.get("purpose"),
            constraints=raw.get("constraints"),
            connections=raw.get("connections", []),
            credentials=raw.get("credentials", []),
            memory=raw.get("memory", []),
            policies=raw.get("policies", []),
            proposals=raw.get("proposals"),
            applied_count=raw.get("applied_count", 0),
            failed_count=raw.get("failed_count", 0),
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            applied_at=raw.get("applied_at"),
        )
