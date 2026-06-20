from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    Webhook,
    ListWebhooksResponse,
    Pagination,
)


class WebhooksResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def create(
        self,
        url: str,
        *,
        events: list[str] | None = None,
        secret: str | None = None,
    ) -> Webhook:
        body: dict[str, Any] = {"url": url}
        if events is not None:
            body["events"] = events
        if secret is not None:
            body["secret"] = secret
        raw = await self._client.post("/api/v1/cadreen/webhooks", body)
        return Webhook(
            id=raw["id"],
            url=raw["url"],
            is_active=raw.get("is_active", True),
            events=raw.get("events"),
            secret=raw.get("secret"),
            created_at=raw.get("created_at"),
        )

    async def list(self) -> ListWebhooksResponse:
        raw = await self._client.get("/api/v1/cadreen/webhooks")
        webhooks = [
            Webhook(
                id=w["id"],
                url=w["url"],
                is_active=w.get("is_active", True),
                events=w.get("events"),
                secret=w.get("secret"),
                created_at=w.get("created_at"),
            )
            for w in raw.get("webhooks", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListWebhooksResponse(
            webhooks=webhooks,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def delete(self, id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/webhooks/{id}")
