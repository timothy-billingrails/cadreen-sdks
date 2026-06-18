import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig, IntentRequest, IntentMessage, IntentContext
from cadreen.resources.intent import IntentResource, _map_intent_response, _default_intelligence
from cadreen.types import (
    DirectResult,
    ClarifyResult,
    ExecutionResult,
    BlockedResult,
    ConnectRequiredResult,
    ClarificationQuestion,
)


@pytest.fixture
def intent_sandbox_client(fixtures):
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
    return HttpClient(config)


class TestMapIntentResponse:
    def test_map_direct(self):
        raw = {
            "type": "direct",
            "message": {"role": "assistant", "content": "Hello world"},
            "trace_id": "tr_1",
        }
        result = _map_intent_response(raw)
        assert isinstance(result, DirectResult)
        assert result.message.content == "Hello world"
        assert result.trace_id == "tr_1"

    def test_map_direct_fallback(self):
        """Unrecognized type falls back to DirectResult"""
        raw = {
            "type": "unknown",
            "message": {"role": "assistant", "content": "default response"},
            "trace_id": "tr_fallback",
        }
        result = _map_intent_response(raw)
        assert isinstance(result, DirectResult)
        assert result.message.content == "default response"

    def test_map_direct_without_type_field(self):
        """Missing 'type' defaults to direct"""
        raw = {
            "message": {"role": "assistant", "content": "implicit direct"},
            "trace_id": "tr_no_type",
        }
        result = _map_intent_response(raw)
        assert isinstance(result, DirectResult)

    def test_map_clarify(self):
        raw = {
            "type": "clarify",
            "trace_id": "tr_clarify",
            "clarification": {
                "conversation_id": "conv_123",
                "questions": [
                    {"id": "q1", "question": "Budget?", "type": "text", "required": True, "reason": "Missing"},
                    {"id": "q2", "question": "Region?", "type": "choice", "required": False},
                ],
            },
        }
        result = _map_intent_response(raw)
        assert isinstance(result, ClarifyResult)
        assert result.conversation_id == "conv_123"
        assert len(result.questions) == 2
        assert result.questions[0].id == "q1"
        assert result.questions[0].question == "Budget?"
        assert result.questions[0].required is True
        assert result.questions[0].reason == "Missing"

    def test_map_clarify_string_questions(self):
        """Questions can be strings instead of dicts"""
        raw = {
            "type": "clarify",
            "trace_id": "tr_str",
            "clarification": {
                "conversation_id": "conv_s",
                "questions": ["What is your name?", "Where are you located?"],
            },
        }
        result = _map_intent_response(raw)
        assert isinstance(result, ClarifyResult)
        assert len(result.questions) == 2
        assert result.questions[0].question == "What is your name?"
        assert result.questions[0].id == ""

    def test_map_mission(self):
        raw = {
            "type": "mission",
            "trace_id": "tr_mis",
            "mission": {
                "id": "mis_001",
                "status": "running",
                "stream_url": "/stream/mis_001",
                "poll_url": "/poll/mis_001",
            },
        }
        result = _map_intent_response(raw)
        assert isinstance(result, ExecutionResult)
        assert result.execution["id"] == "mis_001"
        assert result.execution["status"] == "running"
        assert result.execution["stream_url"] == "/stream/mis_001"

    def test_map_blocked(self):
        raw = {
            "type": "blocked",
            "trace_id": "tr_blocked",
            "meta": {
                "governance": {
                    "decision": "human_approval_required",
                    "reason": "wallet_lockdown_policy",
                }
            },
        }
        result = _map_intent_response(raw)
        assert isinstance(result, BlockedResult)
        assert result.reason_code == "human_approval_required"
        assert result.policy_id == "wallet_lockdown_policy"

    def test_map_connect_required(self):
        raw = {
            "type": "connect_required",
            "trace_id": "tr_connect",
            "mission": {"stream_url": "/connect/github"},
            "meta": {"governance": {"reason": "GitHub integration needed"}},
        }
        result = _map_intent_response(raw)
        assert isinstance(result, ConnectRequiredResult)
        assert result.endpoint == "/connect/github"
        assert result.reason == "GitHub integration needed"

    def test_map_with_intelligence(self):
        raw = {
            "type": "direct",
            "message": {"role": "assistant", "content": "Smart answer"},
            "trace_id": "tr_intel",
            "intelligence": {
                "version": "2026-06-03",
                "summary": "Processed successfully",
                "capability": {"total_available": 10, "healthy_count": 8, "active_integrations": ["stripe"]},
                "reasoning": {"capability_matches": 3},
                "memory": {"healthy": True, "knowledge_queried": 5},
                "governance": {"active": True, "decision": "auto", "confidence": 0.95},
                "humility": {"gaps_detected": 2, "blocking": 0},
                "process": {"started_at": "2026-06-17T00:00:00Z", "duration_ms": 150},
            },
        }
        result = _map_intent_response(raw)
        assert isinstance(result, DirectResult)
        assert result.intelligence.version == "2026-06-03"
        assert result.intelligence.summary == "Processed successfully"
        assert result.intelligence.capability.total_available == 10
        assert result.intelligence.capability.healthy_count == 8
        assert result.intelligence.governance.confidence == 0.95

    def test_default_intelligence(self):
        intel = _default_intelligence()
        assert intel.capability.total_available == 0
        assert intel.memory.healthy is True
        assert intel.governance.active is False
        assert intel.process.duration_ms == 0


class TestIntentResourceInvoke:
    @pytest.mark.asyncio
    async def test_invoke_direct(self, direct_result_fixture):
        fixture = {"POST /api/v1/cadreen/intent": direct_result_fixture["POST /api/v1/cadreen/intent"]}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        msg = IntentMessage(role="user", content="Hello")
        req = IntentRequest(messages=[msg])
        result = await resource.invoke(req)
        assert isinstance(result, DirectResult)
        assert result.message.content == "Here is your answer."
        assert result.trace_id == "trace_abc123"
        assert result.ready is True
        assert result.needs == []

    @pytest.mark.asyncio
    async def test_invoke_blocked(self, blocked_result_fixture):
        fixture = {"POST /api/v1/cadreen/intent": blocked_result_fixture["POST /api/v1/cadreen/intent"]}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        msg = IntentMessage(role="user", content="Do risky thing")
        req = IntentRequest(messages=[msg])
        result = await resource.invoke(req)
        assert isinstance(result, BlockedResult)
        assert result.ready is False
        assert "human_approval_required" in result.needs[0]

    @pytest.mark.asyncio
    async def test_invoke_clarify(self, clarify_result_fixture):
        fixture = {"POST /api/v1/cadreen/intent": clarify_result_fixture["POST /api/v1/cadreen/intent"]}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        msg = IntentMessage(role="user", content="Build something")
        req = IntentRequest(messages=[msg])
        result = await resource.invoke(req)
        assert isinstance(result, ClarifyResult)
        assert len(result.questions) == 2
        assert result.ready is False
        assert "What is the budget?" in result.needs

    @pytest.mark.asyncio
    async def test_invoke_mission(self, mission_result_fixture):
        fixture = {"POST /api/v1/cadreen/intent": mission_result_fixture["POST /api/v1/cadreen/intent"]}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        msg = IntentMessage(role="user", content="Execute task")
        req = IntentRequest(messages=[msg])
        result = await resource.invoke(req)
        assert isinstance(result, ExecutionResult)
        assert result.execution["id"] == "mis_abc"
        assert result.execution["status"] == "running"
        assert result.ready is True

    @pytest.mark.asyncio
    async def test_invoke_connect_required(self, connect_required_result_fixture):
        fixture = {"POST /api/v1/cadreen/intent": connect_required_result_fixture["POST /api/v1/cadreen/intent"]}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        msg = IntentMessage(role="user", content="Use Stripe")
        req = IntentRequest(messages=[msg])
        result = await resource.invoke(req)
        assert isinstance(result, ConnectRequiredResult)
        assert result.endpoint == "/connect/stripe"
        assert result.ready is False

    @pytest.mark.asyncio
    async def test_invoke_with_mode_and_context(self):
        """Invoke with mode and context sends them in the request body"""
        fixture = {
            "POST /api/v1/cadreen/intent": {
                "type": "direct",
                "message": {"role": "assistant", "content": "Acknowledged"},
                "trace_id": "tr_mode",
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        msg = IntentMessage(role="user", content="Process order")
        ctx = IntentContext(domain="sales", existing_connectors=["stripe"])
        req = IntentRequest(messages=[msg], mode="execution", context=ctx, stream=False)
        result = await resource.invoke(req)
        assert isinstance(result, DirectResult)
        assert result.message.content == "Acknowledged"

    @pytest.mark.asyncio
    async def test_call_convenience(self):
        fixture = {
            "POST /api/v1/cadreen/intent": {
                "type": "direct",
                "message": {"role": "assistant", "content": "Hello back!"},
                "trace_id": "tr_call",
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        result = await resource("Hi there", domain="general", mode="chat")
        assert isinstance(result, DirectResult)
        assert result.message.content == "Hello back!"

    @pytest.mark.asyncio
    async def test_call_with_conversation_id(self):
        fixture = {
            "POST /api/v1/cadreen/intent": {
                "type": "direct",
                "message": {"role": "assistant", "content": "Continuing"},
                "trace_id": "tr_conv",
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixture)
        client = HttpClient(config)
        resource = IntentResource(client)
        result = await resource("Continue", conversation_id="conv_xyz")
        assert isinstance(result, DirectResult)
        assert result.message.content == "Continuing"
