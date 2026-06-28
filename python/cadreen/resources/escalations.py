from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    ListEscalationsResponse,
    Escalation,
    Pagination,
)


class EscalationsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(self) -> ListEscalationsResponse:
        raw = await self._client.get("/api/v1/cadreen/escalations")
        escalations = [
            Escalation(
                id=e["id"],
                status=e["status"],
                intent=e.get("intent"),
                category=e.get("category"),
                execution_id=e.get("execution_id"),
                tool_name=e.get("tool_name"),
                error_message=e.get("error_message"),
                severity=e.get("severity"),
                human_prompt=e.get("human_prompt"),
                suggestions=e.get("suggestions"),
                created_at=e.get("created_at"),
                resolved_at=e.get("resolved_at"),
                resolved_by=e.get("resolved_by"),
                resolution=e.get("resolution"),
            )
            for e in raw.get("escalations", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListEscalationsResponse(
            escalations=escalations,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def get(self, id: str) -> Escalation:
        raw = await self._client.get(f"/api/v1/cadreen/escalations/{id}")
        return Escalation(
            id=raw["id"],
            status=raw["status"],
            intent=raw.get("intent"),
            category=raw.get("category"),
            execution_id=raw.get("execution_id"),
            tool_name=raw.get("tool_name"),
            error_message=raw.get("error_message"),
            severity=raw.get("severity"),
            human_prompt=raw.get("human_prompt"),
            suggestions=raw.get("suggestions"),
            created_at=raw.get("created_at"),
            resolved_at=raw.get("resolved_at"),
            resolved_by=raw.get("resolved_by"),
            resolution=raw.get("resolution"),
        )

    async def resolve(self, id: str, decision: str) -> Escalation:
        raw = await self._client.post(
            f"/api/v1/cadreen/escalations/{id}/resolve",
            {"decision": decision},
        )
        return Escalation(
            id=raw["id"],
            status=raw["status"],
            intent=raw.get("intent"),
            category=raw.get("category"),
            execution_id=raw.get("execution_id"),
            tool_name=raw.get("tool_name"),
            error_message=raw.get("error_message"),
            severity=raw.get("severity"),
            human_prompt=raw.get("human_prompt"),
            suggestions=raw.get("suggestions"),
            created_at=raw.get("created_at"),
            resolved_at=raw.get("resolved_at"),
            resolved_by=raw.get("resolved_by"),
            resolution=raw.get("resolution"),
        )
