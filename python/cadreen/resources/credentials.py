from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    CredentialMetadata,
    ListCredentialsResponse,
    Pagination,
)


class CredentialsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(self) -> ListCredentialsResponse:
        raw = await self._client.get("/api/v1/cadreen/credentials")
        credentials = [
            CredentialMetadata(
                id=c["id"],
                provider=c["provider"],
                credential_name=c.get("credential_name", ""),
                is_active=c.get("is_active", False),
                has_credential_data=c.get("has_credential_data", False),
                type=c.get("type"),
            )
            for c in raw.get("credentials", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListCredentialsResponse(
            credentials=credentials,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def create(
        self,
        provider: str,
        key_data: dict,
        *,
        name: str | None = None,
    ) -> CredentialMetadata:
        body: dict[str, Any] = {"provider": provider, "key_data": key_data}
        if name is not None:
            body["name"] = name
        raw = await self._client.post("/api/v1/cadreen/credentials", body)
        return CredentialMetadata(
            id=raw["id"],
            provider=raw["provider"],
            credential_name=raw.get("credential_name", ""),
            is_active=raw.get("is_active", False),
            has_credential_data=raw.get("has_credential_data", False),
            type=raw.get("type"),
        )

    async def delete(self, id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/credentials/{id}")
