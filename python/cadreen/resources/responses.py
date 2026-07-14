from __future__ import annotations

import json
from typing import Any, AsyncIterator

from ..client import HttpClient
from ..types import (
    Response,
    ResponseOutputContent,
    ResponseOutputItem,
    ResponseRequest,
    ResponseStreamEvent,
    ResponseUsage,
)


def _parse_usage(raw: dict[str, Any] | None) -> ResponseUsage | None:
    if not raw:
        return None
    return ResponseUsage(
        input_tokens=raw.get("input_tokens", 0),
        output_tokens=raw.get("output_tokens", 0),
        total_tokens=raw.get("total_tokens", 0),
    )


def _parse_output_content(raw: dict[str, Any]) -> ResponseOutputContent:
    return ResponseOutputContent(
        type=raw.get("type", ""),
        text=raw.get("text"),
        annotations=raw.get("annotations"),
    )


def _parse_output_item(raw: dict[str, Any]) -> ResponseOutputItem:
    content_raw = raw.get("content")
    content = [_parse_output_content(c) for c in content_raw] if content_raw else None
    return ResponseOutputItem(
        id=raw.get("id", ""),
        type=raw.get("type", ""),
        status=raw.get("status"),
        role=raw.get("role"),
        content=content,
        name=raw.get("name"),
        call_id=raw.get("call_id"),
        arguments=raw.get("arguments"),
    )


def _parse_response(raw: dict[str, Any]) -> Response:
    output_raw = raw.get("output", [])
    return Response(
        id=raw.get("id", ""),
        object=raw.get("object", "response"),
        created_at=raw.get("created_at", 0),
        model=raw.get("model", ""),
        output=[_parse_output_item(o) for o in output_raw],
        output_text=raw.get("output_text"),
        usage=raw.get("usage"),
        status=raw.get("status", "completed"),
        previous_response_id=raw.get("previous_response_id"),
        metadata=raw.get("metadata"),
    )


def _request_to_body(request: ResponseRequest) -> dict[str, Any]:
    body: dict[str, Any] = {
        "model": request.model,
        "input": request.input,
        "stream": request.stream,
    }
    if request.instructions is not None:
        body["instructions"] = request.instructions
    if request.tools is not None:
        body["tools"] = request.tools
    if request.previous_response_id is not None:
        body["previous_response_id"] = request.previous_response_id
    if request.store is not None:
        body["store"] = request.store
    if request.max_output_tokens is not None:
        body["max_output_tokens"] = request.max_output_tokens
    if request.temperature is not None:
        body["temperature"] = request.temperature
    if request.metadata is not None:
        body["metadata"] = request.metadata
    return body


class ResponsesResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def create(self, request: ResponseRequest) -> Response:
        body = _request_to_body(request)
        body["stream"] = False
        raw = await self._client.post("/api/v1/cadreen/responses", body)
        return _parse_response(raw)

    async def retrieve(self, response_id: str) -> Response:
        raw = await self._client.get(f"/api/v1/cadreen/responses/{response_id}")
        return _parse_response(raw)

    async def stream(self, request: ResponseRequest) -> AsyncIterator[ResponseStreamEvent]:
        body = _request_to_body(request)
        body["stream"] = True

        url = f"{self._client._base_url}/api/v1/cadreen/responses"
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self._client._api_key}",
            "Accept": "text/event-stream",
        }

        import httpx
        from httpx_sse import aconnect_sse

        async with httpx.AsyncClient(timeout=None) as client:
            async with aconnect_sse(client, "POST", url, headers=headers, json=body) as event_source:
                async for event in event_source.aiter_sse():
                    if event.data == "[DONE]":
                        return
                    try:
                        data = json.loads(event.data)
                        yield ResponseStreamEvent(
                            type=data.get("type", event.event or ""),
                            sequence=data.get("sequence"),
                            response=data.get("response"),
                            item=data.get("item"),
                            output_index=data.get("output_index"),
                            content_index=data.get("content_index"),
                            delta=data.get("delta"),
                        )
                    except Exception:
                        continue
