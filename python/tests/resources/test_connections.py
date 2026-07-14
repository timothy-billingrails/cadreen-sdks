import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.connections import ConnectionsResource
from cadreen.types import (
    RegisterOpenAPIResponse,
    RegisterMCPResponse,
    ListConnectionsResponse,
    ConnectResult,
    ConnectPrebuiltDetail,
    ConnectSchemaRequiredDetail,
    ConnectManualDetail,
    ConnectUnknownDetail,
    CatalogResponse,
    CatalogCategory,
    CatalogIntegration,
    InstallResponse,
)


CONNECTIONS_FIXTURES = {
    "POST /api/v1/cadreen/connections/openapi": {
        "id": "openapi_1",
        "name": "Stripe API",
        "type": "openapi",
        "status": "active",
        "tools_generated": 42,
        "tools_registered": 42,
        "functions": ["create_payment", "list_customers"],
        "spec_url": "https://api.stripe.com/openapi.json",
    },
    "POST /api/v1/cadreen/connections/mcp": {
        "id": "mcp_1",
        "name": "GitHub MCP",
        "type": "mcp",
        "status": "connected",
        "transport": "stdio",
        "url": "https://github-mcp.example.com",
    },
    "GET /api/v1/cadreen/connections": {
        "connections": [
            {
                "capability": "payments",
                "status": "healthy",
            }
        ],
        "total_capabilities": 1,
        "pagination": {"limit": 50, "offset": 0, "has_more": False},
    },
    "DELETE /api/v1/cadreen/connections/conn_1": None,
    "POST /api/v1/cadreen/connections": {
        "type": "prebuilt",
        "capability": "payments",
        "detail": {
            "tool_id": "stripe_v1",
            "tool_name": "Stripe",
            "service_id": "svc_stripe",
            "service_name": "Stripe Payments",
            "auth_type": "oauth2",
            "source": "catalog",
            "account_id": "acct_123",
        },
    },
    "GET /api/v1/cadreen/connections/catalog": {
        "categories": [
            {
                "name": "Payments",
                "description": "Payment processing integrations",
                "integrations": [
                    {
                        "id": "stripe",
                        "name": "Stripe",
                        "description": "Online payment processing",
                        "category": "Payments",
                        "provider": "Stripe Inc.",
                        "status": "available",
                        "auth_type": "oauth2",
                        "install_time": "~2 min",
                        "capabilities": ["payments", "invoicing"],
                        "tags": ["payments", "fintech"],
                        "popularity": 95,
                        "featured": True,
                    }
                ],
            }
        ],
        "installed": ["github", "slack"],
        "total_available": 12,
    },
    "POST /api/v1/cadreen/connections/install": {
        "status": "pending_auth",
        "provider": "stripe",
        "auth_url": "https://connect.stripe.com/oauth/authorize?client_id=xxx",
        "estimated_time": "~2 min",
    },
}


@pytest.fixture
def connections_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=CONNECTIONS_FIXTURES)
    return HttpClient(config)


class TestConnectionsResource:
    @pytest.mark.asyncio
    async def test_catalog(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.catalog()
        assert isinstance(result, CatalogResponse)
        assert result.total_available == 12
        assert len(result.categories) == 1
        assert result.categories[0].name == "Payments"
        assert result.installed == ["github", "slack"]
        integration = result.categories[0].integrations[0]
        assert integration.id == "stripe"
        assert integration.popularity == 95
        assert integration.featured is True

    @pytest.mark.asyncio
    async def test_install(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.install("stripe")
        assert isinstance(result, InstallResponse)
        assert result.status == "pending_auth"
        assert result.provider == "stripe"
        assert "stripe.com" in result.auth_url
        assert result.estimated_time == "~2 min"

    @pytest.mark.asyncio
    async def test_register_openapi(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.register_openapi(
            "Stripe API",
            spec_url="https://api.stripe.com/openapi.json",
        )
        assert isinstance(result, RegisterOpenAPIResponse)
        assert result.id == "openapi_1"
        assert result.name == "Stripe API"
        assert result.type == "openapi"
        assert result.status == "active"
        assert result.tools_generated == 42
        assert result.tools_registered == 42
        assert result.functions == ["create_payment", "list_customers"]

    @pytest.mark.asyncio
    async def test_register_openapi_with_spec_content(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.register_openapi(
            "Custom API",
            spec_content='{"openapi": "3.0"}',
            credential_id="cred_123",
        )
        assert result.id == "openapi_1"

    @pytest.mark.asyncio
    async def test_register_mcp(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.register_mcp(
            "GitHub MCP",
            "https://github-mcp.example.com",
            transport="stdio",
        )
        assert isinstance(result, RegisterMCPResponse)
        assert result.id == "mcp_1"
        assert result.name == "GitHub MCP"
        assert result.type == "mcp"
        assert result.status == "connected"
        assert result.transport == "stdio"
        assert result.url == "https://github-mcp.example.com"

    @pytest.mark.asyncio
    async def test_list(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.list()
        assert isinstance(result, ListConnectionsResponse)
        assert len(result.connections) == 1
        assert result.total_capabilities == 1
        cg = result.connections[0]
        assert cg.capability == "payments"
        assert cg.status == "healthy"
        assert result.pagination.limit == 50
        assert result.pagination.has_more is False

    @pytest.mark.asyncio
    async def test_connect_prebuilt(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.connect("payments")
        assert isinstance(result, ConnectResult)
        assert result.type == "prebuilt"
        assert result.capability == "payments"
        detail = result.detail
        assert isinstance(detail, ConnectPrebuiltDetail)
        assert detail.tool_id == "stripe_v1"
        assert detail.tool_name == "Stripe"
        assert detail.service_name == "Stripe Payments"
        assert detail.auth_type == "oauth2"
        assert detail.source == "catalog"
        assert detail.account_id == "acct_123"

    @pytest.mark.asyncio
    async def test_connect_schema_required(self):
        fixtures = {
            "POST /api/v1/cadreen/connections": {
                "type": "schema_required",
                "capability": "custom_api",
                "detail": {
                    "tool_id": "custom_v1",
                    "tool_name": "Custom API",
                    "auth_url": "https://auth.example.com/connect",
                    "connector": "openapi",
                },
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = ConnectionsResource(client)
        result = await resource.connect("custom_api")
        assert result.type == "schema_required"
        assert isinstance(result.detail, ConnectSchemaRequiredDetail)
        assert result.detail.tool_name == "Custom API"
        assert result.detail.auth_url == "https://auth.example.com/connect"

    @pytest.mark.asyncio
    async def test_connect_manual(self):
        fixtures = {
            "POST /api/v1/cadreen/connections": {
                "type": "manual",
                "capability": "database",
                "detail": {
                    "capability": "database",
                    "available": True,
                    "health": "unknown",
                },
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = ConnectionsResource(client)
        result = await resource.connect("database")
        assert result.type == "manual"
        assert isinstance(result.detail, ConnectManualDetail)
        assert result.detail.capability == "database"
        assert result.detail.available is True
        assert result.detail.health == "unknown"

    @pytest.mark.asyncio
    async def test_connect_unknown(self):
        fixtures = {
            "POST /api/v1/cadreen/connections": {
                "type": "unknown",
                "capability": "teleportation",
                "detail": {
                    "searched": "catalog, community",
                    "hints": ["Try a different capability name", "Check the catalog for alternatives"],
                },
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = ConnectionsResource(client)
        result = await resource.connect("teleportation")
        assert result.type == "unknown"
        assert isinstance(result.detail, ConnectUnknownDetail)
        assert result.detail.searched == "catalog, community"
        assert len(result.detail.hints) == 2

    @pytest.mark.asyncio
    async def test_delete(self, connections_client):
        resource = ConnectionsResource(connections_client)
        result = await resource.delete("conn_1")
        assert result is None

    @pytest.mark.asyncio
    async def test_list_empty(self):
        fixtures = {
            "GET /api/v1/cadreen/connections": {
                "connections": [],
                "total_capabilities": 0,
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = ConnectionsResource(client)
        result = await resource.list()
        assert len(result.connections) == 0
        assert result.total_capabilities == 0
