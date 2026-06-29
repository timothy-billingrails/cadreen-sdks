from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    WorkspaceUser,
    InviteUserRequest,
    UpdateRoleRequest,
    ListWorkspaceUsersResponse,
)


class WorkspaceUsersResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(self) -> ListWorkspaceUsersResponse:
        raw = await self._client.get("/api/v1/cadreen/workspace/users")
        users = [
            WorkspaceUser(
                id=u["id"],
                workspace_id=u["workspace_id"],
                user_id=u["user_id"],
                role=u["role"],
                invited_at=u["invited_at"],
                created_at=u["created_at"],
                updated_at=u["updated_at"],
                invited_by=u.get("invited_by"),
            )
            for u in raw.get("users", [])
        ]
        return ListWorkspaceUsersResponse(
            users=users,
            count=raw.get("count", 0),
        )

    async def invite(self, *, email: str, role: str | None = None) -> WorkspaceUser:
        body: dict[str, Any] = {"email": email}
        if role is not None:
            body["role"] = role
        raw = await self._client.post("/api/v1/cadreen/workspace/users", body)
        return WorkspaceUser(
            id=raw["id"],
            workspace_id=raw["workspace_id"],
            user_id=raw["user_id"],
            role=raw["role"],
            invited_at=raw["invited_at"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            invited_by=raw.get("invited_by"),
        )

    async def update_role(self, id: str, *, role: str) -> WorkspaceUser:
        raw = await self._client.patch(
            f"/api/v1/cadreen/workspace/users/{id}",
            {"role": role},
        )
        return WorkspaceUser(
            id=raw["id"],
            workspace_id=raw["workspace_id"],
            user_id=raw["user_id"],
            role=raw["role"],
            invited_at=raw["invited_at"],
            created_at=raw["created_at"],
            updated_at=raw["updated_at"],
            invited_by=raw.get("invited_by"),
        )

    async def remove(self, id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/workspace/users/{id}")
