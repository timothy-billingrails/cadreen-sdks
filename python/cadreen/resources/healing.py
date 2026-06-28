from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    HealingStatsResponse,
    ListHealingPrecedentsResponse,
    HealingDiagnosis,
    HealingPrecedent,
    StrategyCount,
    ToolHealingStats,
    TimeRange,
    Pagination,
)


class HealingResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def stats(self) -> HealingStatsResponse:
        raw = await self._client.get("/api/v1/cadreen/healing/stats")
        strategies = None
        if raw.get("common_strategies"):
            strategies = [StrategyCount(strategy=s["strategy"], count=s["count"]) for s in raw["common_strategies"]]
        top_tools = None
        if raw.get("top_tools"):
            top_tools = [
                ToolHealingStats(
                    tool_name=t["tool_name"],
                    total=t["total"],
                    successful=t["successful"],
                    failed=t["failed"],
                    success_rate=t["success_rate"],
                    top_strategy=t.get("top_strategy"),
                )
                for t in raw["top_tools"]
            ]
        time_range = None
        if raw.get("time_range"):
            tr = raw["time_range"]
            time_range = TimeRange(first_precedent=tr.get("first_precedent"), last_precedent=tr.get("last_precedent"))
        return HealingStatsResponse(
            total_precedents=raw.get("total_precedents"),
            successful_recoveries=raw.get("successful_recoveries"),
            failed_recoveries=raw.get("failed_recoveries"),
            success_rate=raw.get("success_rate"),
            avg_duration_ms=raw.get("avg_duration_ms"),
            common_strategies=strategies,
            top_tools=top_tools,
            by_category=raw.get("by_category"),
            time_range=time_range,
        )

    async def precedents(self) -> ListHealingPrecedentsResponse:
        raw = await self._client.get("/api/v1/cadreen/healing/precedents")
        precedents = [
            HealingPrecedent(
                id=p["id"],
                error_type=p["error_type"],
                success=p["success"],
                attempts=p.get("attempts", 0),
                confidence=p.get("confidence", 0.0),
                tool_name=p.get("tool_name"),
                error_category=p.get("error_category"),
                semantic_reason=p.get("semantic_reason"),
                root_cause=p.get("root_cause"),
                recovery_strategy=p.get("recovery_strategy"),
                what_worked=p.get("what_worked"),
                what_failed=p.get("what_failed"),
                duration_ms=p.get("duration_ms"),
                created_at=p.get("created_at"),
                domain=p.get("domain"),
                tags=p.get("tags"),
            )
            for p in raw.get("precedents", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListHealingPrecedentsResponse(
            precedents=precedents,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def diagnose(
        self,
        error_message: str,
        *,
        tool_name: str | None = None,
        trace_id: str | None = None,
    ) -> HealingDiagnosis:
        body: dict[str, Any] = {"error_message": error_message}
        if tool_name is not None:
            body["tool_name"] = tool_name
        if trace_id is not None:
            body["trace_id"] = trace_id
        raw = await self._client.post("/api/v1/cadreen/healing/diagnose", body)
        return HealingDiagnosis(
            error_category=raw.get("error_category"),
            semantic_reason=raw.get("semantic_reason"),
            root_cause=raw.get("root_cause"),
            can_retry=raw.get("can_retry"),
            needs_sub_execution=raw.get("needs_sub_execution"),
            needs_human=raw.get("needs_human"),
            should_skip=raw.get("should_skip"),
            needs_re_decide=raw.get("needs_re_decide"),
            needs_try_alternative=raw.get("needs_try_alternative"),
            retry_delay_ms=raw.get("retry_delay_ms"),
            confidence=raw.get("confidence"),
        )
