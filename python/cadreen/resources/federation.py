from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    FederationLink,
    FederationAgent,
    FederationPermissions,
    CreateFederationRequest,
    SuspendFederationRequest,
    RevokeFederationRequest,
    UpdateFederationPermissionsRequest,
    LinkFederationAgentRequest,
    ListFederationResponse,
    ListFederationAgentsResponse,
)


class FederationResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def create(self, request: CreateFederationRequest) -> FederationLink:
        body: dict[str, Any] = {"targetWorkspaceId": request.target_workspace_id}
        if request.description is not None:
            body["description"] = request.description
        if request.permissions is not None:
            body["permissions"] = request.permissions
        raw = await self._client.post("/api/v1/cadreen/federation", body)
        return _parse_federation_link(raw)

    async def list(self) -> ListFederationResponse:
        raw = await self._client.get("/api/v1/cadreen/federation")
        links = [_parse_federation_link(l) for l in raw.get("links", [])]
        return ListFederationResponse(links=links, count=raw.get("count", 0))

    async def get(self, federation_id: str) -> FederationLink:
        raw = await self._client.get(f"/api/v1/cadreen/federation/{federation_id}")
        return _parse_federation_link(raw)

    async def approve(self, federation_id: str) -> FederationLink:
        raw = await self._client.post(f"/api/v1/cadreen/federation/{federation_id}/approve")
        return _parse_federation_link(raw)

    async def suspend(
        self, federation_id: str, request: SuspendFederationRequest | None = None
    ) -> FederationLink:
        body: dict[str, Any] = {}
        if request is not None and request.reason is not None:
            body["reason"] = request.reason
        raw = await self._client.post(
            f"/api/v1/cadreen/federation/{federation_id}/suspend", body or None
        )
        return _parse_federation_link(raw)

    async def revoke(
        self, federation_id: str, request: RevokeFederationRequest | None = None
    ) -> FederationLink:
        body: dict[str, Any] = {}
        if request is not None and request.reason is not None:
            body["reason"] = request.reason
        raw = await self._client.post(
            f"/api/v1/cadreen/federation/{federation_id}/revoke", body or None
        )
        return _parse_federation_link(raw)

    async def get_permissions(self, federation_id: str) -> FederationPermissions:
        raw = await self._client.get(f"/api/v1/cadreen/federation/{federation_id}/permissions")
        return _parse_federation_permissions(raw, federation_id)

    async def update_permissions(
        self, federation_id: str, request: UpdateFederationPermissionsRequest
    ) -> FederationPermissions:
        body: dict[str, Any] = {"permissions": request.permissions}
        if request.allowed_agents is not None:
            body["allowed_agents"] = request.allowed_agents
        if request.allowed_tools is not None:
            body["allowed_tools"] = request.allowed_tools
        if request.allowed_knowledge is not None:
            body["allowed_knowledge"] = request.allowed_knowledge
        if request.max_requests_per_day is not None:
            body["max_requests_per_day"] = request.max_requests_per_day
        raw = await self._client.put(
            f"/api/v1/cadreen/federation/{federation_id}/permissions", body
        )
        return _parse_federation_permissions(raw, federation_id)

    async def link_agent(
        self, federation_id: str, request: LinkFederationAgentRequest
    ) -> FederationAgent:
        body: dict[str, Any] = {"local_agent_id": request.local_agent_id, "remote_agent_id": request.remote_agent_id}
        if request.capabilities is not None:
            body["capabilities"] = request.capabilities
        raw = await self._client.post(
            f"/api/v1/cadreen/federation/{federation_id}/agents", body
        )
        return _parse_federation_agent(raw)

    async def list_agents(self, federation_id: str) -> ListFederationAgentsResponse:
        raw = await self._client.get(f"/api/v1/cadreen/federation/{federation_id}/agents")
        agents = [_parse_federation_agent(a) for a in raw.get("agents", [])]
        return ListFederationAgentsResponse(agents=agents, count=raw.get("count", 0))

    async def unlink_agent(self, federation_id: str, agent_link_id: str) -> None:
        await self._client.delete(
            f"/api/v1/cadreen/federation/{federation_id}/agents/{agent_link_id}"
        )


def _parse_federation_link(raw: dict[str, Any]) -> FederationLink:
    return FederationLink(
        id=raw["id"],
        name=raw.get("name", ""),
        status=raw.get("status", ""),
        target_workspace_id=raw.get("targetWorkspaceId", ""),
        created_at=raw.get("createdAt", ""),
        updated_at=raw.get("updatedAt", ""),
        description=raw.get("description"),
        permissions=raw.get("permissions"),
        created_by_user_id=raw.get("createdByUserId"),
        approved_by_user_id=raw.get("approvedByUserId"),
        suspended_by_user_id=raw.get("suspendedByUserId"),
        suspension_reason=raw.get("suspensionReason"),
        last_activity_at=raw.get("lastActivityAt"),
        revoked_at=raw.get("revokedAt"),
        revoke_reason=raw.get("revokeReason"),
        source_workspace_name=raw.get("sourceWorkspaceName"),
        source_workspace_slug=raw.get("sourceWorkspaceSlug"),
        target_workspace_name=raw.get("targetWorkspaceName"),
    )


def _parse_federation_agent(raw: dict[str, Any]) -> FederationAgent:
    return FederationAgent(
        id=raw["id"],
        federation_link_id=raw.get("federationLinkId", ""),
        local_agent_id=raw.get("localAgentId", ""),
        remote_agent_id=raw.get("remoteAgentId", ""),
        status=raw.get("status", ""),
        created_at=raw.get("createdAt", ""),
        updated_at=raw.get("updatedAt", ""),
        local_agent_name=raw.get("localAgentName"),
        remote_agent_name=raw.get("remoteAgentName"),
        capabilities=raw.get("capabilities"),
    )


def _parse_federation_permissions(raw: dict[str, Any], federation_id: str) -> FederationPermissions:
    return FederationPermissions(
        federation_link_id=raw.get("federation_link_id", raw.get("federationLinkId", federation_id)),
        permissions=raw.get("permissions", []),
        updated_at=raw.get("updated_at", raw.get("updatedAt", "")),
        allowed_agents=raw.get("allowed_agents", raw.get("allowedAgents")),
        allowed_tools=raw.get("allowed_tools", raw.get("allowedTools")),
        allowed_knowledge=raw.get("allowed_knowledge", raw.get("allowedKnowledge")),
        max_requests_per_day=raw.get("max_requests_per_day", raw.get("maxRequestsPerDay")),
    )
