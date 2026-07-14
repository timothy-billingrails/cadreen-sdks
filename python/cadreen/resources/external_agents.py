from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    ExternalAgentConnection,
    ExternalAgentInteraction,
    ExternalAgentSkill,
    ExternalAgentCapabilities,
    ExternalAgentSettings,
    ListExternalConnectionsResponse,
    ListExternalInteractionsResponse,
)


def _parse_skill(raw: dict[str, Any]) -> ExternalAgentSkill:
    return ExternalAgentSkill(
        id=raw.get("id", ""),
        name=raw.get("name", ""),
        description=raw.get("description"),
        tags=raw.get("tags"),
    )


def _parse_capabilities(raw: dict[str, Any]) -> ExternalAgentCapabilities:
    return ExternalAgentCapabilities(
        streaming=raw.get("streaming", False),
        push_notifications=raw.get("pushNotifications", False),
        state_transition_history=raw.get("stateTransitionHistory", False),
    )


def _parse_connection(raw: dict[str, Any]) -> ExternalAgentConnection:
    skills = [_parse_skill(s) for s in raw.get("skills", [])]
    caps = _parse_capabilities(raw.get("capabilities", {}))
    return ExternalAgentConnection(
        id=raw.get("id", ""),
        agent_id=raw.get("agentId", ""),
        agent_card_url=raw.get("agentCardUrl", ""),
        agent_name=raw.get("agentName", ""),
        agent_description=raw.get("agentDescription", ""),
        agent_system=raw.get("agentSystem", ""),
        agent_version=raw.get("agentVersion", ""),
        agent_card_json=raw.get("agentCardJson", {}),
        skills=skills,
        capabilities=caps,
        status=raw.get("status", "pending_approval"),
        health=raw.get("health", "unknown"),
        created_at=raw.get("createdAt", ""),
        updated_at=raw.get("updatedAt", ""),
        last_used_at=raw.get("lastUsedAt"),
        last_health_check_at=raw.get("lastHealthCheckAt"),
        error_message=raw.get("errorMessage"),
        approved_by=raw.get("approvedBy"),
        approved_at=raw.get("approvedAt"),
    )


def _parse_interaction(raw: dict[str, Any]) -> ExternalAgentInteraction:
    return ExternalAgentInteraction(
        id=raw.get("id", ""),
        connection_id=raw.get("connectionId", ""),
        agent_id=raw.get("agentId", ""),
        direction=raw.get("direction", ""),
        operation=raw.get("operation", ""),
        status=raw.get("status", ""),
        created_at=raw.get("createdAt", ""),
        task_id=raw.get("taskId"),
        message=raw.get("message"),
        duration_ms=raw.get("durationMs"),
        error_message=raw.get("errorMessage"),
        governance_result=raw.get("governanceResult"),
    )


class ExternalAgentsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def connect(self, agent_id: str, agent_card_url: str) -> ExternalAgentConnection:
        """Connect to an external A2A agent by providing its Agent Card URL."""
        raw = await self._client.post(
            f"/api/v1/cadreen/agents/{agent_id}/external",
            {"agentCardUrl": agent_card_url},
        )
        return _parse_connection(raw)

    async def list(
        self,
        agent_id: str,
        *,
        status: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListExternalConnectionsResponse:
        """List external agent connections for an agent."""
        params: dict[str, Any] = {}
        if status is not None:
            params["status"] = status
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/external",
            params=params or None,
        )
        connections = [_parse_connection(c) for c in raw.get("connections", [])]
        return ListExternalConnectionsResponse(
            connections=connections,
            total=raw.get("total", 0),
            limit=raw.get("limit", 20),
            offset=raw.get("offset", 0),
        )

    async def get(self, agent_id: str, connection_id: str) -> ExternalAgentConnection:
        """Get a specific external agent connection."""
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/external/{connection_id}"
        )
        return _parse_connection(raw)

    async def approve(self, agent_id: str, connection_id: str) -> dict[str, str]:
        """Approve a pending external agent connection."""
        return await self._client.post(
            f"/api/v1/cadreen/agents/{agent_id}/external/{connection_id}/approve"
        )

    async def suspend(self, agent_id: str, connection_id: str) -> dict[str, str]:
        """Suspend an active external agent connection."""
        return await self._client.post(
            f"/api/v1/cadreen/agents/{agent_id}/external/{connection_id}/suspend"
        )

    async def revoke(self, agent_id: str, connection_id: str) -> dict[str, str]:
        """Revoke an external agent connection (permanent)."""
        return await self._client.post(
            f"/api/v1/cadreen/agents/{agent_id}/external/{connection_id}/revoke"
        )

    async def delete(self, agent_id: str, connection_id: str) -> None:
        """Delete an external agent connection."""
        await self._client.delete(
            f"/api/v1/cadreen/agents/{agent_id}/external/{connection_id}"
        )

    async def list_interactions(
        self,
        agent_id: str,
        connection_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListExternalInteractionsResponse:
        """List interactions for an external agent connection."""
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/external/{connection_id}/interactions",
            params=params or None,
        )
        interactions = [_parse_interaction(i) for i in raw.get("interactions", [])]
        return ListExternalInteractionsResponse(
            interactions=interactions,
            total=raw.get("total", 0),
            limit=raw.get("limit", 20),
            offset=raw.get("offset", 0),
        )

    async def get_settings(self) -> ExternalAgentSettings:
        """Get external agent settings for the workspace."""
        raw = await self._client.get("/api/v1/cadreen/external-agents/settings")
        return ExternalAgentSettings(enabled=raw.get("enabled", False))

    async def update_settings(self, enabled: bool) -> ExternalAgentSettings:
        """Enable or disable external agents for the workspace."""
        raw = await self._client.put(
            "/api/v1/cadreen/external-agents/settings", {"enabled": enabled}
        )
        return ExternalAgentSettings(enabled=raw.get("enabled", False))

    async def list_all(self) -> ListExternalConnectionsResponse:
        """List all external agent connections in the workspace."""
        raw = await self._client.get("/api/v1/cadreen/external-agents/connections")
        connections = [_parse_connection(c) for c in raw.get("connections", [])]
        return ListExternalConnectionsResponse(
            connections=connections,
            total=raw.get("total", len(connections)),
            limit=raw.get("limit", len(connections)),
            offset=raw.get("offset", 0),
        )
