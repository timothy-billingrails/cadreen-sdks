from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    AcceptProposalResponse,
    DismissProposalResponse,
    ListProposalsResponse,
    ProposalEvidence,
    ProposalStatsResponse,
    TaskProposal,
)


class ProposalsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(self, *, status: str | None = None, limit: int | None = None) -> ListProposalsResponse:
        params: list[str] = []
        if status is not None:
            params.append(f"status={status}")
        if limit is not None:
            params.append(f"limit={limit}")
        query = "?" + "&".join(params) if params else ""
        raw = await self._client.get(f"/api/v1/cadreen/proposals{query}")
        proposals = [
            TaskProposal(
                id=p["id"],
                title=p["title"],
                description=p["description"],
                intent=p["intent"],
                proposal_type=p["proposal_type"],
                trigger_type=p["trigger_type"],
                trigger_source=p["trigger_source"],
                confidence=p["confidence"],
                priority=p["priority"],
                status=p["status"],
                created_at=p["created_at"],
                domain=p.get("domain"),
                mission_intent=p.get("mission_intent"),
                trigger_details=p.get("trigger_details"),
                evidence=[
                    ProposalEvidence(
                        description=e["description"],
                        source=e.get("source"),
                        count=e.get("count"),
                        confidence=e.get("confidence"),
                    )
                    for e in p["evidence"]
                ] if p.get("evidence") else None,
                expires_at=p.get("expires_at"),
                accepted_at=p.get("accepted_at"),
                dismissed_at=p.get("dismissed_at"),
                dismissal_reason=p.get("dismissal_reason"),
                execution_id=p.get("execution_id"),
                dedup_key=p.get("dedup_key"),
                requires_review=p.get("requires_review"),
            )
            for p in raw.get("proposals", [])
        ]
        return ListProposalsResponse(
            proposals=proposals,
            count=raw.get("count", 0),
        )

    async def get(self, id: str) -> TaskProposal:
        raw = await self._client.get(f"/api/v1/cadreen/proposals/{id}")
        evidence = None
        if raw.get("evidence"):
            evidence = [
                ProposalEvidence(
                    description=e["description"],
                    source=e.get("source"),
                    count=e.get("count"),
                    confidence=e.get("confidence"),
                )
                for e in raw["evidence"]
            ]
        return TaskProposal(
            id=raw["id"],
            title=raw["title"],
            description=raw["description"],
            intent=raw["intent"],
            proposal_type=raw["proposal_type"],
            trigger_type=raw["trigger_type"],
            trigger_source=raw["trigger_source"],
            confidence=raw["confidence"],
            priority=raw["priority"],
            status=raw["status"],
            created_at=raw["created_at"],
            domain=raw.get("domain"),
            mission_intent=raw.get("mission_intent"),
            trigger_details=raw.get("trigger_details"),
            evidence=evidence,
            expires_at=raw.get("expires_at"),
            accepted_at=raw.get("accepted_at"),
            dismissed_at=raw.get("dismissed_at"),
            dismissal_reason=raw.get("dismissal_reason"),
            execution_id=raw.get("execution_id"),
            dedup_key=raw.get("dedup_key"),
            requires_review=raw.get("requires_review"),
        )

    async def accept(self, id: str) -> AcceptProposalResponse:
        raw = await self._client.post(f"/api/v1/cadreen/proposals/{id}/accept", {})
        return AcceptProposalResponse(
            status=raw["status"],
            execution_id=raw["execution_id"],
            action=raw["action"],
            intent=raw["intent"],
            next_step=raw["next_step"],
            auto_approved=raw.get("auto_approved"),
            result=raw.get("result"),
        )

    async def dismiss(self, id: str, *, reason: str | None = None) -> DismissProposalResponse:
        body: dict[str, Any] = {}
        if reason is not None:
            body["reason"] = reason
        raw = await self._client.post(f"/api/v1/cadreen/proposals/{id}/dismiss", body)
        return DismissProposalResponse(
            status=raw["status"],
        )

    async def stats(self) -> ProposalStatsResponse:
        raw = await self._client.get("/api/v1/cadreen/proposals/stats")
        return ProposalStatsResponse(
            proposed=raw.get("proposed", 0),
            accepted=raw.get("accepted", 0),
            dismissed=raw.get("dismissed", 0),
            expired=raw.get("expired", 0),
        )
