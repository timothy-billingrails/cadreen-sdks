import pytest

from cadreen import Cadreen, CadreenError
from cadreen.types import (
    IntentRequest,
    IntentMessage,
    IntentContext,
    CadreenConfig,
    SetupRequest,
    SetupConnection,
    SetupCredential,
    SetupMemory,
    SetupPolicy,
    RequestOptions,
)
from cadreen.resources.intent import IntentResource
from cadreen.resources.memory import MemoryResource
from cadreen.resources.policies import PoliciesResource
from cadreen.resources.connections import ConnectionsResource
from cadreen.resources.traces import TracesResource
from cadreen.resources.executions import ExecutionsResource
from cadreen.resources.guardrails import GuardrailsResource


INTEGRATION_FIXTURES = {
    "POST /api/v1/cadreen/intent": {
        "type": "direct",
        "message": {"role": "assistant", "content": "Hello! How can I help you?"},
        "trace_id": "trace_int_001",
    },
    "POST /api/v1/cadreen/memory": {
        "id": "mem_int_1",
        "type": "fact",
        "domain": "general",
        "authority": 5,
        "version": 1,
        "content": {"text": "Stored for integration test"},
        "indexed": True,
        "tags": ["test"],
    },
    "GET /api/v1/cadreen/memory/search?query=test": {
        "results": [
            {"id": "atom_found", "type": "fact", "domain": "general", "content": {"text": "Found item"}},
        ],
        "count": 1,
    },
    "GET /api/v1/cadreen/connections": {
        "connections": [],
        "total_capabilities": 0,
        "total_pathways": 0,
    },
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
        },
    },
    "POST /api/v1/cadreen/setup": {
        "connections": [],
        "credentials": [],
        "memory": [],
        "policies": [],
        "applied": 3,
        "failed": 0,
        "workspace_id": "ws_test",
    },
}


@pytest.fixture
def cadreen():
    return Cadreen(api_key="test_key", sandbox=True, fixtures=INTEGRATION_FIXTURES)


@pytest.fixture
def cadreen_small():
    """Cadreen with minimal fixture set"""
    return Cadreen(api_key="test_key", sandbox=True, fixtures={
        "POST /api/v1/cadreen/intent": {
            "type": "direct",
            "message": {"role": "assistant", "content": "Minimal response"},
            "trace_id": "tr_small",
        }
    })


class TestCadreenIntegration:
    @pytest.mark.asyncio
    async def test_ask_returns_direct_result(self, cadreen):
        result = await cadreen.ask("Hello")
        from cadreen.types import DirectResult
        assert isinstance(result, DirectResult)
        assert result.message.content == "Hello! How can I help you?"
        assert result.trace_id == "trace_int_001"

    @pytest.mark.asyncio
    async def test_ask_with_conversation_id(self, cadreen_small):
        result = await cadreen_small.ask("Hello", conversation_id="conv_abc")
        assert result.message.content == "Minimal response"

    @pytest.mark.asyncio
    async def test_act_uses_execution_mode(self, cadreen_small):
        result = await cadreen_small.act("Process orders", conversation_id="conv_exec")
        assert result.message.content == "Minimal response"

    @pytest.mark.asyncio
    async def test_remember_delegates_to_memory(self, cadreen):
        result = await cadreen.remember("fact", {"text": "Important data"}, domain="general", tags=["test"])
        assert result.id == "mem_int_1"
        assert result.type == "fact"
        assert result.content.text == "Stored for integration test"
        assert result.indexed is True

    @pytest.mark.asyncio
    async def test_context_delegates_to_memory_search(self, cadreen):
        result = await cadreen.context("test")
        assert result.count == 1
        assert len(result.results) == 1
        assert result.results[0].id == "atom_found"

    @pytest.mark.asyncio
    async def test_connect_delegates_to_connections(self, cadreen):
        result = await cadreen.connect("payments")
        from cadreen.types import ConnectResult, ConnectPrebuiltDetail
        assert isinstance(result, ConnectResult)
        assert result.type == "prebuilt"
        assert isinstance(result.detail, ConnectPrebuiltDetail)
        assert result.detail.tool_id == "stripe_v1"

    @pytest.mark.asyncio
    async def test_setup(self, cadreen):
        req = SetupRequest(
            connections=[SetupConnection(capability="stripe")],
            purpose="Test workspace setup",
            workspace_id="ws_test",
        )
        result = await cadreen.setup(req)
        from cadreen.types import SetupResult
        assert isinstance(result, SetupResult)
        assert result.applied == 3
        assert result.failed == 0
        assert result.workspace_id == "ws_test"

    @pytest.mark.asyncio
    async def test_invoke_delegates_to_intent(self, cadreen):
        msg = IntentMessage(role="user", content="Direct invoke")
        req = IntentRequest(messages=[msg])
        result = await cadreen.invoke(req)
        from cadreen.types import DirectResult
        assert isinstance(result, DirectResult)
        assert result.message.content == "Hello! How can I help you?"

    def test_public_resource_members(self, cadreen):
        assert isinstance(cadreen.intent, IntentResource)
        assert isinstance(cadreen.memory, MemoryResource)
        assert isinstance(cadreen.policies, PoliciesResource)
        assert isinstance(cadreen.connections, ConnectionsResource)
        assert isinstance(cadreen.traces, TracesResource)
        assert isinstance(cadreen.executions, ExecutionsResource)
        assert isinstance(cadreen.guardrails, GuardrailsResource)

    def test_cadreen_config_custom_base_url(self):
        c = Cadreen(api_key="key", base_url="http://localhost:8080", sandbox=True)
        assert c._client._base_url == "http://localhost:8080"

    def test_cadreen_config_custom_retries_and_timeout(self):
        c = Cadreen(api_key="key", sandbox=True, max_retries=5, timeout=60)
        assert c._client._max_retries == 5
        assert c._client._timeout == 60

    @pytest.mark.asyncio
    async def test_setup_with_all_resource_types(self, cadreen):
        req = SetupRequest(
            connections=[SetupConnection(capability="stripe"), SetupConnection(capability="github")],
            credentials=[SetupCredential(provider="stripe", key_data={"api_key": "sk_test"}, name="stripe_key")],
            memory=[SetupMemory(content={"text": "Test"}, type="fact", domain="general", tags=["setup"])],
            policies=[SetupPolicy(name="Approval Policy", rule="requires_human: true", description="Test")],
            workspace_id="ws_full",
            purpose="Full workspace setup",
            examples=["Example usage"],
            constraints=["Budget < $1000"],
        )
        result = await cadreen.setup(req)
        assert result.applied == 3
        assert result.failed == 0

    @pytest.mark.asyncio
    async def test_ask_with_context(self, cadreen_small):
        ctx = IntentContext(domain="finance", constraints={"budget": 5000})
        result = await cadreen_small.ask("What's my budget?", context=ctx)
        assert result.message.content == "Minimal response"

    @pytest.mark.asyncio
    async def test_act_with_context_and_stream(self, cadreen_small):
        ctx = IntentContext(existing_connectors=["stripe", "github"])
        result = await cadreen_small.act("Deploy service", context=ctx, stream=False)
        assert result.message.content == "Minimal response"
