import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.agents import AgentsResource
from cadreen.types import (
    Agent,
    AgentConfig,
    AgentCapabilities,
    AgentMessage,
    AgentKnowledge,
    AgentGovernancePolicy,
    AgentAuditEntry,
    AgentNegotiation,
    CreateAgentRequest,
    UpdateAgentRequest,
    SendMessageRequest,
    CreateExecutionRequest,
    CreateAgentKnowledgeRequest,
    CreateAgentGovernanceRequest,
    UpdateAgentGovernanceRequest,
    StartNegotiationRequest,
    RespondNegotiationRequest,
    ListAgentsResponse,
    ListAgentMessagesResponse,
    ListAgentExecutionsResponse,
    ListAgentKnowledgeResponse,
    ListAgentGovernanceResponse,
    ListAgentAuditResponse,
    ListAgentNegotiationsResponse,
    SearchAgentKnowledgeResponse,
    ExecutionStatus,
)


AGENT_FIXTURES = {
    "POST /api/v1/cadreen/agents": {
        "id": "agt_abc",
        "name": "Support Bot",
        "status": "active",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T00:00:00Z",
        "description": "Customer support agent",
        "model": "gpt-4o",
        "domain": "support",
        "tags": ["support", "billing"],
        "version": 1,
    },
    "GET /api/v1/cadreen/agents": {
        "agents": [
            {
                "id": "agt_1",
                "name": "Support Bot",
                "status": "active",
                "created_at": "2026-07-01T00:00:00Z",
                "updated_at": "2026-07-01T00:00:00Z",
                "description": "Handles support tickets",
                "model": "gpt-4o",
            },
            {
                "id": "agt_2",
                "name": "Sales Bot",
                "status": "draft",
                "created_at": "2026-07-02T00:00:00Z",
                "updated_at": "2026-07-02T00:00:00Z",
            },
        ],
        "count": 2,
    },
    "GET /api/v1/cadreen/agents/agt_abc": {
        "id": "agt_abc",
        "name": "Support Bot",
        "status": "active",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T00:00:00Z",
        "description": "Customer support agent",
        "model": "gpt-4o",
        "domain": "support",
        "version": 1,
    },
    "PATCH /api/v1/cadreen/agents/agt_abc": {
        "id": "agt_abc",
        "name": "Updated Bot",
        "status": "active",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T01:00:00Z",
        "description": "Updated description",
        "model": "gpt-4o",
        "version": 2,
    },
    "DELETE /api/v1/cadreen/agents/agt_abc": None,
    "GET /api/v1/cadreen/agents/agt_abc/config": {
        "model": "gpt-4o",
        "domain": "support",
        "config": {"temperature": 0.7},
        "version": 1,
    },
    "POST /api/v1/cadreen/agents/agt_abc/deploy": {
        "id": "agt_abc",
        "name": "Support Bot",
        "status": "deployed",
        "created_at": "2026-07-09T00:00:00Z",
        "updated_at": "2026-07-09T02:00:00Z",
        "version": 1,
    },
    "GET /api/v1/cadreen/agents/agt_abc/capabilities": {
        "tools": ["search_docs", "create_ticket"],
        "integrations": ["zendesk", "slack"],
        "knowledge_count": 42,
        "governance_count": 3,
    },
    "POST /api/v1/cadreen/agents/agt_abc/send": {
        "id": "msg_001",
        "role": "user",
        "content": "Hello, agent!",
        "created_at": "2026-07-09T03:00:00Z",
    },
    "GET /api/v1/cadreen/agents/agt_abc/messages": {
        "messages": [
            {
                "id": "msg_001",
                "role": "user",
                "content": "Hello",
                "created_at": "2026-07-09T03:00:00Z",
            },
            {
                "id": "msg_002",
                "role": "assistant",
                "content": "Hi! How can I help?",
                "created_at": "2026-07-09T03:00:01Z",
            },
        ],
        "count": 2,
    },
    "GET /api/v1/cadreen/agents/agt_abc/executions": {
        "executions": [
            {
                "id": "exec_001",
                "status": "completed",
                "progress": 1.0,
                "result": {"output": "done"},
            }
        ],
        "count": 1,
    },
    "POST /api/v1/cadreen/agents/agt_abc/executions": {
        "id": "exec_002",
        "status": "running",
    },
    "GET /api/v1/cadreen/agents/agt_abc/knowledge": {
        "knowledge": [
            {
                "id": "know_001",
                "type": "reference",
                "content": {"text": "Company policy is..."},
                "created_at": "2026-07-09T00:00:00Z",
                "domain": "support",
            }
        ],
        "count": 1,
    },
    "POST /api/v1/cadreen/agents/agt_abc/knowledge": {
        "id": "know_002",
        "type": "reference",
        "content": {"text": "New knowledge"},
        "created_at": "2026-07-09T04:00:00Z",
    },
    "POST /api/v1/cadreen/agents/agt_abc/knowledge/search": {
        "results": [
            {
                "id": "know_001",
                "type": "reference",
                "content": {"text": "Matching result"},
                "created_at": "2026-07-09T00:00:00Z",
            }
        ],
        "count": 1,
    },
    "DELETE /api/v1/cadreen/agents/agt_abc/knowledge/know_001": None,
    "GET /api/v1/cadreen/agents/agt_abc/governance": {
        "policies": [
            {
                "id": "gov_001",
                "name": "No PII Exposure",
                "rules": [{"action": "redact_pii", "severity": "high"}],
                "created_at": "2026-07-09T00:00:00Z",
                "domain": "support",
                "priority": 1,
            }
        ],
        "count": 1,
    },
    "POST /api/v1/cadreen/agents/agt_abc/governance": {
        "id": "gov_002",
        "name": "Spending Limit",
        "rules": [{"action": "limit_spend", "max": 1000}],
        "created_at": "2026-07-09T05:00:00Z",
    },
    "PATCH /api/v1/cadreen/agents/agt_abc/governance/gov_001": {
        "id": "gov_001",
        "name": "Updated Policy",
        "rules": [{"action": "redact_pii", "severity": "critical"}],
        "created_at": "2026-07-09T00:00:00Z",
        "priority": 1,
    },
    "DELETE /api/v1/cadreen/agents/agt_abc/governance/gov_001": None,
    "GET /api/v1/cadreen/agents/agt_abc/audit": {
        "entries": [
            {
                "id": "aud_001",
                "action": "agent.created",
                "timestamp": "2026-07-09T00:00:00Z",
                "actor": "user_123",
                "detail": "Agent created via API",
            }
        ],
        "count": 1,
    },
    "POST /api/v1/cadreen/agents/agt_abc/negotiate": {
        "id": "neg_001",
        "status": "pending",
        "initiator_agent_id": "agt_abc",
        "target_agent_id": "agt_xyz",
        "proposal": {"task": "share_knowledge", "scope": "support"},
        "created_at": "2026-07-09T06:00:00Z",
        "updated_at": "2026-07-09T06:00:00Z",
    },
    "GET /api/v1/cadreen/agents/agt_abc/negotiations": {
        "negotiations": [
            {
                "id": "neg_001",
                "status": "pending",
                "initiator_agent_id": "agt_abc",
                "target_agent_id": "agt_xyz",
                "proposal": {"task": "share_knowledge"},
                "created_at": "2026-07-09T06:00:00Z",
                "updated_at": "2026-07-09T06:00:00Z",
            }
        ],
        "count": 1,
    },
    "GET /api/v1/cadreen/agents/agt_abc/negotiations/neg_001": {
        "id": "neg_001",
        "status": "accepted",
        "initiator_agent_id": "agt_abc",
        "target_agent_id": "agt_xyz",
        "proposal": {"task": "share_knowledge"},
        "created_at": "2026-07-09T06:00:00Z",
        "updated_at": "2026-07-09T07:00:00Z",
        "response": {"accepted": True},
    },
    "POST /api/v1/cadreen/agents/agt_abc/negotiations/neg_001/respond": {
        "id": "neg_001",
        "status": "accepted",
        "initiator_agent_id": "agt_abc",
        "target_agent_id": "agt_xyz",
        "proposal": {"task": "share_knowledge"},
        "created_at": "2026-07-09T06:00:00Z",
        "updated_at": "2026-07-09T07:00:00Z",
        "response": {"accepted": True},
    },
}


@pytest.fixture
def agents_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=AGENT_FIXTURES)
    return HttpClient(config)


class TestAgentsResource:
    @pytest.mark.asyncio
    async def test_create(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.create(CreateAgentRequest(
            name="Support Bot",
            description="Customer support agent",
            model="gpt-4o",
            domain="support",
            tags=["support", "billing"],
        ))
        assert isinstance(result, Agent)
        assert result.id == "agt_abc"
        assert result.name == "Support Bot"
        assert result.status == "active"
        assert result.model == "gpt-4o"
        assert result.domain == "support"
        assert result.tags == ["support", "billing"]

    @pytest.mark.asyncio
    async def test_list(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list()
        assert isinstance(result, ListAgentsResponse)
        assert result.count == 2
        assert len(result.agents) == 2
        assert result.agents[0].id == "agt_1"
        assert result.agents[1].id == "agt_2"

    @pytest.mark.asyncio
    async def test_get(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.get("agt_abc")
        assert isinstance(result, Agent)
        assert result.id == "agt_abc"
        assert result.name == "Support Bot"
        assert result.version == 1

    @pytest.mark.asyncio
    async def test_update(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.update("agt_abc", UpdateAgentRequest(name="Updated Bot", description="Updated description"))
        assert isinstance(result, Agent)
        assert result.name == "Updated Bot"
        assert result.version == 2

    @pytest.mark.asyncio
    async def test_delete(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.delete("agt_abc")
        assert result is None

    @pytest.mark.asyncio
    async def test_get_config(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.get_config("agt_abc")
        assert isinstance(result, AgentConfig)
        assert result.model == "gpt-4o"
        assert result.domain == "support"
        assert result.config == {"temperature": 0.7}
        assert result.version == 1

    @pytest.mark.asyncio
    async def test_deploy(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.deploy("agt_abc")
        assert isinstance(result, Agent)
        assert result.status == "deployed"

    @pytest.mark.asyncio
    async def test_get_capabilities(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.get_capabilities("agt_abc")
        assert isinstance(result, AgentCapabilities)
        assert result.tools == ["search_docs", "create_ticket"]
        assert result.integrations == ["zendesk", "slack"]
        assert result.knowledge_count == 42
        assert result.governance_count == 3

    @pytest.mark.asyncio
    async def test_send_message(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.send_message("agt_abc", SendMessageRequest(content="Hello, agent!"))
        assert isinstance(result, AgentMessage)
        assert result.id == "msg_001"
        assert result.content == "Hello, agent!"
        assert result.role == "user"

    @pytest.mark.asyncio
    async def test_list_messages(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list_messages("agt_abc")
        assert isinstance(result, ListAgentMessagesResponse)
        assert result.count == 2
        assert result.messages[0].role == "user"
        assert result.messages[1].role == "assistant"

    @pytest.mark.asyncio
    async def test_list_executions(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list_executions("agt_abc")
        assert isinstance(result, ListAgentExecutionsResponse)
        assert result.count == 1
        assert result.executions[0].id == "exec_001"
        assert result.executions[0].status == "completed"

    @pytest.mark.asyncio
    async def test_create_execution(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.create_execution("agt_abc", CreateExecutionRequest(task="Summarize tickets"))
        assert isinstance(result, ExecutionStatus)
        assert result.id == "exec_002"
        assert result.status == "running"

    @pytest.mark.asyncio
    async def test_list_knowledge(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list_knowledge("agt_abc")
        assert isinstance(result, ListAgentKnowledgeResponse)
        assert result.count == 1
        assert result.knowledge[0].type == "reference"

    @pytest.mark.asyncio
    async def test_create_knowledge(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.create_knowledge("agt_abc", CreateAgentKnowledgeRequest(
            type="reference", content={"text": "New knowledge"},
        ))
        assert isinstance(result, AgentKnowledge)
        assert result.id == "know_002"

    @pytest.mark.asyncio
    async def test_search_knowledge(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.search_knowledge("agt_abc", {"query": "policy"})
        assert isinstance(result, SearchAgentKnowledgeResponse)
        assert result.count == 1
        assert result.results[0].id == "know_001"

    @pytest.mark.asyncio
    async def test_delete_knowledge(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.delete_knowledge("agt_abc", "know_001")
        assert result is None

    @pytest.mark.asyncio
    async def test_list_governance(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list_governance("agt_abc")
        assert isinstance(result, ListAgentGovernanceResponse)
        assert result.count == 1
        assert result.policies[0].name == "No PII Exposure"

    @pytest.mark.asyncio
    async def test_create_governance(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.create_governance("agt_abc", CreateAgentGovernanceRequest(
            name="Spending Limit", rules=[{"action": "limit_spend", "max": 1000}],
        ))
        assert isinstance(result, AgentGovernancePolicy)
        assert result.id == "gov_002"

    @pytest.mark.asyncio
    async def test_update_governance(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.update_governance("agt_abc", "gov_001", UpdateAgentGovernanceRequest(
            name="Updated Policy", rules=[{"action": "redact_pii", "severity": "critical"}],
        ))
        assert isinstance(result, AgentGovernancePolicy)
        assert result.name == "Updated Policy"

    @pytest.mark.asyncio
    async def test_delete_governance(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.delete_governance("agt_abc", "gov_001")
        assert result is None

    @pytest.mark.asyncio
    async def test_list_audit(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list_audit("agt_abc")
        assert isinstance(result, ListAgentAuditResponse)
        assert result.count == 1
        assert result.entries[0].action == "agent.created"
        assert result.entries[0].actor == "user_123"

    @pytest.mark.asyncio
    async def test_start_negotiation(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.start_negotiation("agt_abc", StartNegotiationRequest(
            target_agent_id="agt_xyz", proposal={"task": "share_knowledge", "scope": "support"},
        ))
        assert isinstance(result, AgentNegotiation)
        assert result.id == "neg_001"
        assert result.status == "pending"
        assert result.target_agent_id == "agt_xyz"

    @pytest.mark.asyncio
    async def test_list_negotiations(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.list_negotiations("agt_abc")
        assert isinstance(result, ListAgentNegotiationsResponse)
        assert result.count == 1
        assert result.negotiations[0].id == "neg_001"

    @pytest.mark.asyncio
    async def test_get_negotiation(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.get_negotiation("agt_abc", "neg_001")
        assert isinstance(result, AgentNegotiation)
        assert result.status == "accepted"
        assert result.response == {"accepted": True}

    @pytest.mark.asyncio
    async def test_respond_to_negotiation(self, agents_client):
        resource = AgentsResource(agents_client)
        result = await resource.respond_to_negotiation("agt_abc", "neg_001", RespondNegotiationRequest(
            response="accepted",
        ))
        assert isinstance(result, AgentNegotiation)
        assert result.status == "accepted"

    @pytest.mark.asyncio
    async def test_list_empty(self):
        fixtures = {"GET /api/v1/cadreen/agents": {"agents": [], "count": 0}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = AgentsResource(client)
        result = await resource.list()
        assert isinstance(result, ListAgentsResponse)
        assert result.count == 0
        assert len(result.agents) == 0
