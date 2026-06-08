from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    RegisterOpenAPIRequest,
    RegisterOpenAPIResponse,
    RegisterMCPRequest,
    RegisterMCPResponse,
    ListConnectionsResponse,
    InstallComposioRequest,
    ConnectionGroup,
    Pathway,
    Pagination,
    ConnectResult,
    ConnectPrebuiltDetail,
    ConnectSchemaRequiredDetail,
    ConnectManualDetail,
    ConnectPathway,
    ConnectUnknownDetail,
)


class ConnectionsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def register_openapi(
        self,
        name: str,
        *,
        spec_url: str | None = None,
        spec_content: str | None = None,
        credential_id: str | None = None,
    ) -> RegisterOpenAPIResponse:
        body: dict[str, Any] = {"name": name}
        if spec_url is not None:
            body["spec_url"] = spec_url
        if spec_content is not None:
            body["spec_content"] = spec_content
        if credential_id is not None:
            body["credential_id"] = credential_id

        raw = await self._client.post("/api/v1/cadreen/connections/openapi", body)
        return RegisterOpenAPIResponse(
            id=raw["id"],
            name=raw["name"],
            type=raw.get("type", ""),
            status=raw.get("status", ""),
            tools_generated=raw.get("tools_generated"),
            tools_registered=raw.get("tools_registered"),
            functions=raw.get("functions"),
            spec_url=raw.get("spec_url"),
        )

    async def register_mcp(
        self,
        name: str,
        url: str,
        *,
        transport: str | None = None,
        headers: dict[str, str] | None = None,
    ) -> RegisterMCPResponse:
        body: dict[str, Any] = {"name": name, "url": url}
        if transport is not None:
            body["transport"] = transport
        if headers is not None:
            body["headers"] = headers

        raw = await self._client.post("/api/v1/cadreen/connections/mcp", body)
        return RegisterMCPResponse(
            id=raw["id"],
            name=raw["name"],
            type=raw.get("type", ""),
            status=raw.get("status", ""),
            transport=raw.get("transport"),
            url=raw.get("url"),
        )

    async def install_composio(self, toolkit: str, *, user_id: str | None = None) -> dict[str, Any]:
        body: dict[str, Any] = {"toolkit": toolkit}
        if user_id is not None:
            body["user_id"] = user_id
        return await self._client.post("/api/v1/cadreen/connections/composio/install", body)

    async def search_composio(self, query: str) -> dict[str, Any]:
        return await self._client.post("/api/v1/cadreen/connections/composio/search", {"query": query})

    async def composio_status(self, toolkit: str | None = None, user_id: str | None = None) -> dict[str, Any]:
        params: dict[str, Any] = {}
        if toolkit is not None:
            params["toolkit"] = toolkit
        if user_id is not None:
            params["user_id"] = user_id
        return await self._client.get("/api/v1/cadreen/connections/composio/status", params)

    async def list(self) -> ListConnectionsResponse:
        raw = await self._client.get("/api/v1/cadreen/connections")
        connections: list[ConnectionGroup] = []
        for cg in raw.get("connections", []):
            pathways = None
            if cg.get("pathways"):
                pathways = [
                    Pathway(
                        id=p["id"],
                        capability=p["capability"],
                        connector=p["connector"],
                        transport=p["transport"],
                        health=p["health"],
                        tool_id=p["tool_id"],
                    )
                    for p in cg["pathways"]
                ]
            connections.append(ConnectionGroup(
                capability=cg["capability"],
                pathways=pathways,
                status=cg.get("status", "unknown"),
            ))
        pagination = None
        if raw.get("pagination"):
            pagination = Pagination(
                limit=raw["pagination"]["limit"],
                offset=raw["pagination"]["offset"],
                has_more=raw["pagination"]["has_more"],
            )
        return ListConnectionsResponse(
            connections=connections,
            total_capabilities=raw.get("total_capabilities", 0),
            total_pathways=raw.get("total_pathways", 0),
            pagination=pagination,
        )

    async def delete(self, id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/connections/{id}")

    async def connect(self, capability: str) -> ConnectResult:
        raw = await self._client.post("/api/v1/cadreen/connections", {"capability": capability})
        result_type = raw.get("type", "unknown")
        detail_raw = raw.get("detail", {})

        detail: ConnectPrebuiltDetail | ConnectSchemaRequiredDetail | ConnectManualDetail | ConnectUnknownDetail
        if result_type == "prebuilt":
            detail = ConnectPrebuiltDetail(
                tool_id=detail_raw.get("tool_id", ""),
                tool_name=detail_raw.get("tool_name", ""),
                service_id=detail_raw.get("service_id", ""),
                service_name=detail_raw.get("service_name", ""),
                auth_type=detail_raw.get("auth_type", ""),
                account_id=detail_raw.get("account_id"),
                source=detail_raw.get("source", ""),
            )
        elif result_type == "schema_required":
            detail = ConnectSchemaRequiredDetail(
                tool_id=detail_raw.get("tool_id", ""),
                tool_name=detail_raw.get("tool_name", ""),
                auth_url=detail_raw.get("auth_url", ""),
                connector=detail_raw.get("connector", ""),
            )
        elif result_type == "manual":
            pathways = [
                ConnectPathway(
                    id=p["id"],
                    connector=p["connector"],
                    tool_id=p["tool_id"],
                    health=p["health"],
                    priority=p["priority"],
                )
                for p in detail_raw.get("pathways", [])
            ]
            detail = ConnectManualDetail(pathways=pathways)
        else:
            detail = ConnectUnknownDetail(
                searched=detail_raw.get("searched", ""),
                hints=detail_raw.get("hints"),
            )

        return ConnectResult(type=result_type, capability=raw.get("capability", capability), detail=detail)
