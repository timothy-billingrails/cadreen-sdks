from __future__ import annotations

from typing import Any, Optional

from ..client import HttpClient
from ..resources.intent import _parse_intelligence
from ..types import (
    IntelligenceTraceEntry,
    ListIntelligenceResponse,
    IntelligenceStats,
    Pagination,
    ReplayResult,
    HandoffPacket,
    PromoteResult,
)


class TracesResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def get(self, id: str) -> IntelligenceTraceEntry:
        raw = await self._client.get(f"/api/v1/cadreen/intelligence/{id}")
        entry = IntelligenceTraceEntry(
            id=raw["id"],
            domain=raw.get("domain", ""),
            request_path=raw.get("request_path", ""),
            request_method=raw.get("request_method", ""),
            meta=_parse_intelligence(raw.get("meta", {})),
            created_at=raw.get("created_at"),
        )
        return entry

    async def list(
        self,
        domain: Optional[str] = None,
        decision: Optional[str] = None,
        from_: Optional[str] = None,
        to: Optional[str] = None,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> ListIntelligenceResponse:
        params: dict[str, Any] = {}
        if domain is not None:
            params["domain"] = domain
        if decision is not None:
            params["decision"] = decision
        if from_ is not None:
            params["from"] = from_
        if to is not None:
            params["to"] = to
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset

        raw = await self._client.get("/api/v1/cadreen/intelligence", params)
        traces = [
            IntelligenceTraceEntry(
                id=t["id"],
                domain=t.get("domain", ""),
                request_path=t.get("request_path", ""),
                request_method=t.get("request_method", ""),
                meta=_parse_intelligence(t.get("meta", {})),
                created_at=t.get("created_at"),
            )
            for t in raw.get("traces", [])
        ]
        pagination = None
        if raw.get("pagination"):
            pagination = Pagination(
                limit=raw["pagination"]["limit"],
                offset=raw["pagination"]["offset"],
                has_more=raw["pagination"]["has_more"],
            )
        return ListIntelligenceResponse(
            traces=traces,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def stats(self) -> IntelligenceStats:
        raw = await self._client.get("/api/v1/cadreen/intelligence/stats")
        return IntelligenceStats(
            traces_24h=raw.get("traces_24h", 0),
            traces_7d=raw.get("traces_7d", 0),
            traces_30d=raw.get("traces_30d", 0),
            avg_confidence_by_domain=raw.get("avg_confidence_by_domain", {}),
            gap_detection_rate=raw.get("gap_detection_rate", 0.0),
            governance_decisions=raw.get("governance_decisions", {}),
        )

    async def replay(self, id: str, mode: str = "current") -> ReplayResult:
        raw = await self._client.post(f"/api/v1/cadreen/intelligence/{id}/replay", {"mode": mode})
        return ReplayResult(
            trace_id=raw["trace_id"],
            mode=raw.get("mode", ""),
            domain=raw.get("domain", ""),
            original_gate=raw.get("original_gate", ""),
            original_confidence=raw.get("original_confidence", 0.0),
            current_gate=raw.get("current_gate", ""),
            current_confidence=raw.get("current_confidence", 0.0),
            gate_changed=raw.get("gate_changed", False),
            change_summary=raw.get("change_summary", ""),
            current_capability=raw.get("current_capability", {}),
            current_memory=raw.get("current_memory", {}),
            current_gaps=raw.get("current_gaps", {}),
            replay_note=raw.get("replay_note", ""),
        )

    async def handoff(self, id: str) -> HandoffPacket:
        raw = await self._client.get(f"/api/v1/cadreen/intelligence/{id}/handoff")
        return HandoffPacket(
            trace_id=raw["trace_id"],
            domain=raw.get("domain", ""),
            created_at=raw.get("created_at", ""),
            governance=raw.get("governance", {}),
            what_the_system_knew=raw.get("what_the_system_knew", {}),
            what_the_system_didnt_know=raw.get("what_the_system_didnt_know", {}),
            what_happened=raw.get("what_happened", {}),
            suggested_actions=raw.get("suggested_actions", []),
            next_action=raw.get("next_action", {}),
            trace_url=raw.get("trace_url", ""),
        )

    async def promote(self, id: str, kind: str, name: str | None = None) -> PromoteResult:
        body: dict[str, Any] = {"kind": kind}
        if name:
            body["name"] = name
        raw = await self._client.post(f"/api/v1/cadreen/intelligence/{id}/promote", body)
        return PromoteResult(
            id=raw["id"],
            kind=raw.get("kind", ""),
            status=raw.get("status", ""),
            tool_name=raw.get("tool_name"),
            tool_sequence=raw.get("tool_sequence"),
            source_trace_id=raw.get("source_trace_id", ""),
        )

