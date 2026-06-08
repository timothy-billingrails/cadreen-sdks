from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import ExecutionEvent, ExecutionStatus


class ExecutionsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def stream(self, execution_id: str):
        path = f"/api/v1/cadreen/executions/{execution_id}/stream"
        async for event in self._client.stream(path):
            yield ExecutionEvent(
                type=event["type"],
                data=event["data"],
            )

    async def get_status(self, execution_id: str) -> ExecutionStatus:
        raw = await self._client.get(f"/api/v1/cadreen/executions/{execution_id}")
        return ExecutionStatus(
            id=raw["id"],
            status=raw.get("status", ""),
            progress=raw.get("progress"),
            result=raw.get("result"),
            error=raw.get("error"),
        )
