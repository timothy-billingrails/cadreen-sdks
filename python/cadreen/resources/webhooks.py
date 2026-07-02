from __future__ import annotations

import hashlib
import hmac
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

    @staticmethod
    def verify_signature(raw_body: str | bytes, signature: str, secret: str) -> bool:
        """Verify the HMAC-SHA256 signature of a webhook payload.

        Args:
            raw_body: The raw request body (do NOT parse JSON first).
            signature: The value of the X-Cadreen-Signature header.
            secret: The secret you set when creating the webhook subscription.

        Returns:
            True if the signature is valid, False otherwise.
        """
        if not raw_body or not signature or not secret:
            return False

        try:
            if isinstance(raw_body, str):
                raw_body = raw_body.encode("utf-8")
            expected = hmac.new(secret.encode("utf-8"), raw_body, hashlib.sha256).hexdigest()
            return hmac.compare_digest(signature, expected)
        except Exception:
            return False
