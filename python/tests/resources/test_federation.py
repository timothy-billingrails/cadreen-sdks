import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.federation import FederationResource
from cadreen.types import (
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


FEDERATION_FIXTURES = {
    "POST /api/v1/cadreen/federation": {
        "id": "fed_abc",
        "name": "Partner Corp Link",
        "status": "pending",
        "targetWorkspaceId": "tenant_xyz",
        "createdAt": "2026-07-09T00:00:00Z",
        "updatedAt": "2026-07-09T00:00:00Z",
        "description": "Federation with Partner Corp",
        "permissions": ["read:knowledge", "write:messages"],
    },
    "GET /api/v1/cadreen/federation": {
        "links": [
            {
                "id": "fed_1",
                "name": "Partner A",
                "status": "active",
                "targetWorkspaceId": "tenant_a",
                "createdAt": "2026-07-01T00:00:00Z",
                "updatedAt": "2026-07-01T00:00:00Z",
            },
            {
                "id": "fed_2",
                "name": "Partner B",
                "status": "pending",
                "targetWorkspaceId": "tenant_b",
                "createdAt": "2026-07-02T00:00:00Z",
                "updatedAt": "2026-07-02T00:00:00Z",
            },
        ],
        "count": 2,
    },
    "GET /api/v1/cadreen/federation/fed_abc": {
        "id": "fed_abc",
        "name": "Partner Corp Link",
        "status": "active",
        "target_tenant_id": "tenant_xyz",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T01:00:00Z",
        "description": "Federation with Partner Corp",
        "permissions": ["read:knowledge"],
    },
    "POST /api/v1/cadreen/federation/fed_abc/approve": {
        "id": "fed_abc",
        "name": "Partner Corp Link",
        "status": "active",
        "target_tenant_id": "tenant_xyz",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T02:00:00Z",
    },
    "POST /api/v1/cadreen/federation/fed_abc/suspend": {
        "id": "fed_abc",
        "name": "Partner Corp Link",
        "status": "suspended",
        "target_tenant_id": "tenant_xyz",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T03:00:00Z",
    },
    "POST /api/v1/cadreen/federation/fed_abc/revoke": {
        "id": "fed_abc",
        "name": "Partner Corp Link",
        "status": "revoked",
        "target_tenant_id": "tenant_xyz",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T04:00:00Z",
    },
    "GET /api/v1/cadreen/federation/fed_abc/permissions": {
        "federation_id": "fed_abc",
        "permissions": ["read:knowledge", "write:messages"],
        "updated_at": "2026-07-09T00:00:00Z",
    },
    "PUT /api/v1/cadreen/federation/fed_abc/permissions": {
        "federation_id": "fed_abc",
        "permissions": ["read:knowledge", "write:messages", "execute:agents"],
        "updated_at": "2026-07-09T05:00:00Z",
    },
    "POST /api/v1/cadreen/federation/fed_abc/agents": {
        "id": "link_001",
        "localAgentId": "agt_xyz",
        "federationLinkId": "fed_abc",
        "status": "linked",
        "createdAt": "2026-07-09T06:00:00Z",
        "capabilities": ["negotiate", "share_knowledge"],
    },
    "GET /api/v1/cadreen/federation/fed_abc/agents": {
        "agents": [
            {
                "id": "link_001",
                "localAgentId": "agt_xyz",
                "federationLinkId": "fed_abc",
                "status": "linked",
                "createdAt": "2026-07-09T06:00:00Z",
                "capabilities": ["negotiate"],
            }
        ],
        "count": 1,
    },
    "DELETE /api/v1/cadreen/federation/fed_abc/agents/link_001": None,
}


@pytest.fixture
def federation_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=FEDERATION_FIXTURES)
    return HttpClient(config)


class TestFederationResource:
    @pytest.mark.asyncio
    async def test_create(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.create(CreateFederationRequest(
            target_workspace_id="ws_xyz",
            description="Federation with Partner Corp",
            permissions=["read:knowledge", "write:messages"],
        ))
        assert isinstance(result, FederationLink)
        assert result.id == "fed_abc"
        assert result.status == "pending"

    @pytest.mark.asyncio
    async def test_list(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.list()
        assert isinstance(result, ListFederationResponse)
        assert result.count == 2
        assert len(result.links) == 2
        assert result.links[0].id == "fed_1"
        assert result.links[1].id == "fed_2"

    @pytest.mark.asyncio
    async def test_get(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.get("fed_abc")
        assert isinstance(result, FederationLink)
        assert result.id == "fed_abc"
        assert result.status == "active"

    @pytest.mark.asyncio
    async def test_approve(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.approve("fed_abc")
        assert isinstance(result, FederationLink)
        assert result.status == "active"

    @pytest.mark.asyncio
    async def test_suspend(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.suspend("fed_abc", SuspendFederationRequest(reason="Policy violation"))
        assert isinstance(result, FederationLink)
        assert result.status == "suspended"

    @pytest.mark.asyncio
    async def test_suspend_no_reason(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.suspend("fed_abc")
        assert isinstance(result, FederationLink)
        assert result.status == "suspended"

    @pytest.mark.asyncio
    async def test_revoke(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.revoke("fed_abc", RevokeFederationRequest(reason="Security breach"))
        assert isinstance(result, FederationLink)
        assert result.status == "revoked"

    @pytest.mark.asyncio
    async def test_revoke_no_reason(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.revoke("fed_abc")
        assert isinstance(result, FederationLink)
        assert result.status == "revoked"

    @pytest.mark.asyncio
    async def test_get_permissions(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.get_permissions("fed_abc")
        assert isinstance(result, FederationPermissions)
        assert result.federation_id == "fed_abc"
        assert result.permissions == ["read:knowledge", "write:messages"]

    @pytest.mark.asyncio
    async def test_update_permissions(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.update_permissions("fed_abc", UpdateFederationPermissionsRequest(
            permissions=["read:knowledge", "write:messages", "execute:agents"],
        ))
        assert isinstance(result, FederationPermissions)
        assert len(result.permissions) == 3
        assert "execute:agents" in result.permissions

    @pytest.mark.asyncio
    async def test_link_agent(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.link_agent("fed_abc", LinkFederationAgentRequest(
            local_agent_id="agt_local",
            remote_agent_id="agt_remote",
        ))
        assert isinstance(result, FederationAgent)
        assert result.id == "link_001"
        assert result.status == "linked"

    @pytest.mark.asyncio
    async def test_list_agents(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.list_agents("fed_abc")
        assert isinstance(result, ListFederationAgentsResponse)
        assert result.count == 1
        assert result.agents[0].agent_id == "agt_xyz"

    @pytest.mark.asyncio
    async def test_unlink_agent(self, federation_client):
        resource = FederationResource(federation_client)
        result = await resource.unlink_agent("fed_abc", "link_001")
        assert result is None

    @pytest.mark.asyncio
    async def test_list_empty(self):
        fixtures = {"GET /api/v1/cadreen/federation": {"links": [], "count": 0}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = FederationResource(client)
        result = await resource.list()
        assert isinstance(result, ListFederationResponse)
        assert result.count == 0
        assert len(result.links) == 0
