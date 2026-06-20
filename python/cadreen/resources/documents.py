from __future__ import annotations

from ..client import HttpClient
from ..types import (
    ListDocumentsResponse,
    Document,
    Pagination,
)


class DocumentsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(self) -> ListDocumentsResponse:
        raw = await self._client.get("/api/v1/cadreen/documents")
        documents = [
            Document(
                id=d["id"],
                name=d["name"],
                content_type=d.get("content_type"),
                size=d.get("size"),
                status=d.get("status"),
                created_at=d.get("created_at"),
            )
            for d in raw.get("documents", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListDocumentsResponse(
            documents=documents,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def get(self, id: str) -> Document:
        raw = await self._client.get(f"/api/v1/cadreen/documents/{id}")
        return Document(
            id=raw["id"],
            name=raw["name"],
            content_type=raw.get("content_type"),
            size=raw.get("size"),
            status=raw.get("status"),
            created_at=raw.get("created_at"),
        )
