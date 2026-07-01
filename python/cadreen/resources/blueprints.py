from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    Blueprint,
    BlueprintRun,
    BlueprintSource,
    CreateBlueprintRequest,
    UpdateBlueprintRequest,
    ListBlueprintsResponse,
    ListBlueprintRunsResponse,
)


class BlueprintsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(
        self,
        *,
        status: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListBlueprintsResponse:
        params: dict[str, Any] = {}
        if status is not None:
            params["status"] = status
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get("/api/v1/cadreen/blueprints", params=params or None)
        blueprints = [
            Blueprint(
                id=b["id"],
                name=b["name"],
                status=b["status"],
                version=b["version"],
                created_at=b["created_at"],
                updated_at=b["updated_at"],
                description=b.get("description"),
                intent=b.get("intent"),
                source_type=b.get("source_type"),
                source_id=b.get("source_id"),
            )
            for b in raw.get("blueprints", [])
        ]
        return ListBlueprintsResponse(
            blueprints=blueprints,
            count=raw.get("count", 0),
        )

    async def get(self, blueprint_id: str) -> Blueprint:
        raw = await self._client.get(f"/api/v1/cadreen/blueprints/{blueprint_id}")
        return Blueprint(
            id=raw["id"],
            name=raw["name"],
            status=raw["status"],
            version=raw["version"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            description=raw.get("description"),
            intent=raw.get("intent"),
            source_type=raw.get("source_type"),
            source_id=raw.get("source_id"),
        )

    async def create(self, request: CreateBlueprintRequest) -> Blueprint:
        body: dict[str, Any] = {"name": request.name}
        if request.description is not None:
            body["description"] = request.description
        if request.source is not None:
            body["source"] = {"type": request.source.type}
            if request.source.trace_id is not None:
                body["source"]["trace_id"] = request.source.trace_id
            if request.source.execution_id is not None:
                body["source"]["execution_id"] = request.source.execution_id
        if request.parameter_schema is not None:
            body["parameter_schema"] = request.parameter_schema
        if request.default_params is not None:
            body["default_params"] = request.default_params
        raw = await self._client.post("/api/v1/cadreen/blueprints", body)
        return Blueprint(
            id=raw["id"],
            name=raw["name"],
            status=raw["status"],
            version=raw["version"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            description=raw.get("description"),
            intent=raw.get("intent"),
            source_type=raw.get("source_type"),
            source_id=raw.get("source_id"),
        )

    async def update(self, blueprint_id: str, request: UpdateBlueprintRequest) -> Blueprint:
        body: dict[str, Any] = {}
        if request.name is not None:
            body["name"] = request.name
        if request.description is not None:
            body["description"] = request.description
        if request.parameter_schema is not None:
            body["parameter_schema"] = request.parameter_schema
        if request.default_params is not None:
            body["default_params"] = request.default_params
        raw = await self._client.patch(f"/api/v1/cadreen/blueprints/{blueprint_id}", body)
        return Blueprint(
            id=raw["id"],
            name=raw["name"],
            status=raw["status"],
            version=raw["version"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            description=raw.get("description"),
            intent=raw.get("intent"),
            source_type=raw.get("source_type"),
            source_id=raw.get("source_id"),
        )

    async def delete(self, blueprint_id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/blueprints/{blueprint_id}")

    async def run(self, blueprint_id: str, params: dict[str, Any] | None = None) -> BlueprintRun:
        body: dict[str, Any] = {}
        if params is not None:
            body["params"] = params
        raw = await self._client.post(f"/api/v1/cadreen/blueprints/{blueprint_id}/runs", body)
        return BlueprintRun(
            id=raw["id"],
            blueprint_id=raw["blueprint_id"],
            blueprint_version=raw["blueprint_version"],
            status=raw["status"],
            created_at=raw["created_at"],
            params=raw.get("params"),
            result_summary=raw.get("result_summary"),
            trace_id=raw.get("trace_id"),
        )

    async def get_runs(
        self,
        blueprint_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListBlueprintRunsResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/blueprints/{blueprint_id}/runs",
            params=params or None,
        )
        runs = [
            BlueprintRun(
                id=r["id"],
                blueprint_id=r["blueprint_id"],
                blueprint_version=r["blueprint_version"],
                status=r["status"],
                created_at=r["created_at"],
                params=r.get("params"),
                result_summary=r.get("result_summary"),
                trace_id=r.get("trace_id"),
            )
            for r in raw.get("runs", [])
        ]
        return ListBlueprintRunsResponse(
            runs=runs,
            count=raw.get("count", 0),
        )
