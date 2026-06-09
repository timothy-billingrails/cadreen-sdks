from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    RememberRequest,
    CreateMemoryResponse,
    SearchMemoryRequest,
    SearchMemoryResponse,
    MemoryTypesResponse,
    Atom,
    AtomContent,
)


def _parse_atom(raw: dict[str, Any]) -> Atom:
    content = None
    if raw.get("content"):
        c = raw["content"]
        content = AtomContent(
            text=c.get("text"),
            source=c.get("source"),
            subject=c.get("subject"),
            constraint=c.get("constraint"),
            query=c.get("query"),
            tools_used=c.get("tools_used"),
            outcome=c.get("outcome"),
            situation=c.get("situation"),
            action=c.get("action"),
            result=c.get("result"),
            name=c.get("name"),
            constraints=c.get("constraints"),
            deadline=c.get("deadline"),
            is_private=c.get("is_private"),
        )
    return Atom(
        id=raw["id"],
        type=raw["type"],
        domain=raw["domain"],
        authority=raw.get("authority", 0),
        version=raw.get("version", 0),
        scope=raw.get("scope"),
        content=content,
        tags=raw.get("tags"),
        created_at=raw.get("created_at"),
    )


def _parse_create_response(raw: dict[str, Any]) -> CreateMemoryResponse:
    content = None
    if raw.get("content"):
        c = raw["content"]
        content = AtomContent(
            text=c.get("text"),
            source=c.get("source"),
            subject=c.get("subject"),
            constraint=c.get("constraint"),
            query=c.get("query"),
            tools_used=c.get("tools_used"),
            outcome=c.get("outcome"),
            situation=c.get("situation"),
            action=c.get("action"),
            result=c.get("result"),
            name=c.get("name"),
            constraints=c.get("constraints"),
            deadline=c.get("deadline"),
            is_private=c.get("is_private"),
        )
    return CreateMemoryResponse(
        id=raw["id"],
        type=raw["type"],
        domain=raw["domain"],
        authority=raw.get("authority", 0),
        version=raw.get("version", 0),
        scope=raw.get("scope"),
        content=content,
        indexed=raw.get("indexed"),
        tags=raw.get("tags"),
        created_at=raw.get("created_at"),
    )


class MemoryResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def remember(
        self,
        type: str,
        content: dict[str, Any],
        *,
        domain: str | None = None,
        scope: str | None = None,
        authority: int | None = None,
        tags: list[str] | None = None,
    ) -> CreateMemoryResponse:
        body: dict[str, Any] = {"type": type, "content": content}
        if domain is not None:
            body["domain"] = domain
        if scope is not None:
            body["scope"] = scope
        if authority is not None:
            body["authority"] = authority
        if tags is not None:
            body["tags"] = tags

        raw = await self._client.post("/api/v1/cadreen/memory", body)
        return _parse_create_response(raw)

    async def search(self, request: SearchMemoryRequest) -> SearchMemoryResponse:
        params: dict[str, Any] = {"query": request.query}
        if request.domain is not None:
            params["domain"] = request.domain
        if request.tag is not None:
            params["tag"] = request.tag
        if request.limit is not None:
            params["limit"] = request.limit

        raw = await self._client.get("/api/v1/cadreen/memory/search", params)
        return SearchMemoryResponse(
            results=[_parse_atom(r) for r in raw.get("results", [])],
            count=raw.get("count", 0),
        )

    async def get(self, id: str) -> Atom:
        raw = await self._client.get(f"/api/v1/cadreen/memory/{id}")
        return _parse_atom(raw)

    async def types(self) -> MemoryTypesResponse:
        raw = await self._client.get("/api/v1/cadreen/memory/types")
        return MemoryTypesResponse(
            type_values=raw.get("type_values", []),
            kind_values=raw.get("kind_values", []),
            description=raw.get("description", ""),
        )
