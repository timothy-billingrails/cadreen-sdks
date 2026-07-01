from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    Schedule,
    ScheduleRun,
    CreateScheduleRequest,
    UpdateScheduleRequest,
    ListSchedulesResponse,
    ListScheduleRunsResponse,
    PauseScheduleResponse,
    ResumeScheduleResponse,
)


class SchedulesResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(self) -> ListSchedulesResponse:
        raw = await self._client.get("/api/v1/cadreen/schedules")
        schedules = [
            Schedule(
                id=s["id"],
                name=s["name"],
                blueprint_id=s["blueprint_id"],
                blueprint_version=s["blueprint_version"],
                status=s["status"],
                trigger=s["trigger"],
                timezone=s["timezone"],
                created_at=s["created_at"],
                updated_at=s["updated_at"],
                params=s.get("params"),
                next_run_at=s.get("next_run_at"),
                last_run_at=s.get("last_run_at"),
                pause_reason=s.get("pause_reason"),
            )
            for s in raw.get("schedules", [])
        ]
        return ListSchedulesResponse(
            schedules=schedules,
            count=raw.get("count", 0),
        )

    async def get(self, schedule_id: str) -> Schedule:
        raw = await self._client.get(f"/api/v1/cadreen/schedules/{schedule_id}")
        return Schedule(
            id=raw["id"],
            name=raw["name"],
            blueprint_id=raw["blueprint_id"],
            blueprint_version=raw["blueprint_version"],
            status=raw["status"],
            trigger=raw["trigger"],
            timezone=raw["timezone"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            params=raw.get("params"),
            next_run_at=raw.get("next_run_at"),
            last_run_at=raw.get("last_run_at"),
            pause_reason=raw.get("pause_reason"),
        )

    async def create(self, request: CreateScheduleRequest) -> Schedule:
        body: dict[str, Any] = {
            "blueprint_id": request.blueprint_id,
            "name": request.name,
            "trigger": request.trigger,
        }
        if request.timezone is not None:
            body["timezone"] = request.timezone
        if request.params is not None:
            body["params"] = request.params
        raw = await self._client.post("/api/v1/cadreen/schedules", body)
        return Schedule(
            id=raw["id"],
            name=raw["name"],
            blueprint_id=raw["blueprint_id"],
            blueprint_version=raw["blueprint_version"],
            status=raw["status"],
            trigger=raw["trigger"],
            timezone=raw["timezone"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            params=raw.get("params"),
            next_run_at=raw.get("next_run_at"),
            last_run_at=raw.get("last_run_at"),
            pause_reason=raw.get("pause_reason"),
        )

    async def update(self, schedule_id: str, request: UpdateScheduleRequest) -> Schedule:
        body: dict[str, Any] = {}
        if request.name is not None:
            body["name"] = request.name
        if request.trigger is not None:
            body["trigger"] = request.trigger
        if request.timezone is not None:
            body["timezone"] = request.timezone
        if request.params is not None:
            body["params"] = request.params
        raw = await self._client.patch(f"/api/v1/cadreen/schedules/{schedule_id}", body)
        return Schedule(
            id=raw["id"],
            name=raw["name"],
            blueprint_id=raw["blueprint_id"],
            blueprint_version=raw["blueprint_version"],
            status=raw["status"],
            trigger=raw["trigger"],
            timezone=raw["timezone"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            params=raw.get("params"),
            next_run_at=raw.get("next_run_at"),
            last_run_at=raw.get("last_run_at"),
            pause_reason=raw.get("pause_reason"),
        )

    async def delete(self, schedule_id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/schedules/{schedule_id}")

    async def pause(self, schedule_id: str, reason: str | None = None) -> PauseScheduleResponse:
        body: dict[str, Any] = {}
        if reason is not None:
            body["reason"] = reason
        raw = await self._client.post(f"/api/v1/cadreen/schedules/{schedule_id}/pause", body)
        return PauseScheduleResponse(
            id=raw["id"],
            status=raw["status"],
        )

    async def resume(self, schedule_id: str) -> ResumeScheduleResponse:
        raw = await self._client.post(f"/api/v1/cadreen/schedules/{schedule_id}/resume")
        return ResumeScheduleResponse(
            id=raw["id"],
            status=raw["status"],
            next_run_at=raw.get("next_run_at"),
        )

    async def get_runs(
        self,
        schedule_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListScheduleRunsResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/schedules/{schedule_id}/runs",
            params=params or None,
        )
        runs = [
            ScheduleRun(
                id=r["id"],
                schedule_id=r["schedule_id"],
                status=r["status"],
                started_at=r.get("started_at"),
                finished_at=r.get("finished_at"),
                result_summary=r.get("result_summary"),
                error=r.get("error"),
            )
            for r in raw.get("runs", [])
        ]
        return ListScheduleRunsResponse(
            runs=runs,
            count=raw.get("count", 0),
        )
