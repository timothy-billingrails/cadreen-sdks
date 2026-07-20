from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
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
    CreateAgentKnowledgeRequest,
    CreateAgentGovernanceRequest,
    UpdateAgentGovernanceRequest,
    StartNegotiationRequest,
    RespondNegotiationRequest,
    SendMessageRequest,
    CreateExecutionRequest,
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


class AgentsResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def create(self, request: CreateAgentRequest) -> Agent:
        body: dict[str, Any] = {"name": request.name}
        if request.description is not None:
            body["description"] = request.description
        if request.model is not None:
            body["model"] = request.model
        if request.domain is not None:
            body["domain"] = request.domain
        if request.config is not None:
            body["config"] = request.config
        if request.tags is not None:
            body["tags"] = request.tags
        raw = await self._client.post("/api/v1/cadreen/agents", body)
        return _parse_agent(raw)

    async def list(
        self,
        *,
        search: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListAgentsResponse:
        params: dict[str, Any] = {}
        if search is not None:
            params["search"] = search
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get("/api/v1/cadreen/agents", params=params or None)
        agents = [_parse_agent(a) for a in raw.get("agents", [])]
        return ListAgentsResponse(agents=agents, count=raw.get("count", 0))

    async def get(self, agent_id: str) -> Agent:
        raw = await self._client.get(f"/api/v1/cadreen/agents/{agent_id}")
        return _parse_agent(raw)

    async def update(self, agent_id: str, request: UpdateAgentRequest) -> Agent:
        body: dict[str, Any] = {}
        if request.name is not None:
            body["name"] = request.name
        if request.description is not None:
            body["description"] = request.description
        if request.model is not None:
            body["model"] = request.model
        if request.domain is not None:
            body["domain"] = request.domain
        if request.config is not None:
            body["config"] = request.config
        if request.tags is not None:
            body["tags"] = request.tags
        if request.status is not None:
            body["status"] = request.status
        raw = await self._client.patch(f"/api/v1/cadreen/agents/{agent_id}", body)
        return _parse_agent(raw)

    async def delete(self, agent_id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/agents/{agent_id}")

    async def get_config(self, agent_id: str) -> AgentConfig:
        raw = await self._client.get(f"/api/v1/cadreen/agents/{agent_id}/config")
        return AgentConfig(
            model=raw.get("model", ""),
            domain=raw.get("domain", ""),
            config=raw.get("config", {}),
            version=raw.get("version", 0),
        )

    async def deploy(self, agent_id: str, config_snapshot: dict[str, Any], change_summary: str | None = None) -> Agent:
        body: dict[str, Any] = {"configSnapshot": config_snapshot}
        if change_summary is not None:
            body["changeSummary"] = change_summary
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/deploy", body)
        return _parse_agent(raw)

    async def get_capabilities(self, agent_id: str) -> AgentCapabilities:
        raw = await self._client.get(f"/api/v1/cadreen/agents/{agent_id}/capabilities")
        return AgentCapabilities(
            tools=raw.get("tools", []),
            integrations=raw.get("integrations", []),
            knowledge_count=raw.get("knowledge_count", 0),
            governance_count=raw.get("governance_count", 0),
        )

    async def send_message(self, agent_id: str, request: SendMessageRequest) -> AgentMessage:
        body: dict[str, Any] = {"fromAgentId": request.from_agent_id, "content": request.content}
        if request.context is not None:
            body["context"] = request.context
        if request.execution_id is not None:
            body["executionId"] = request.execution_id
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/send", body)
        return _parse_agent_message(raw)

    async def list_messages(
        self,
        agent_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListAgentMessagesResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/messages",
            params=params or None,
        )
        messages = [_parse_agent_message(m) for m in raw.get("messages", [])]
        return ListAgentMessagesResponse(messages=messages, count=raw.get("count", 0))

    async def list_executions(
        self,
        agent_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListAgentExecutionsResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/executions",
            params=params or None,
        )
        executions = [
            ExecutionStatus(
                id=e["id"],
                status=e["status"],
                progress=e.get("progress"),
                result=e.get("result"),
                error=e.get("error"),
            )
            for e in raw.get("executions", [])
        ]
        return ListAgentExecutionsResponse(executions=executions, count=raw.get("count", 0))

    async def create_execution(self, agent_id: str, request: CreateExecutionRequest) -> ExecutionStatus:
        body: dict[str, Any] = {"intent": request.intent}
        if request.context is not None:
            body["context"] = request.context
        if request.max_budget_usd is not None:
            body["maxBudgetUsd"] = request.max_budget_usd
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/executions", body)
        return ExecutionStatus(
            id=raw["id"],
            status=raw["status"],
            progress=raw.get("progress"),
            result=raw.get("result"),
            error=raw.get("error"),
        )

    async def list_knowledge(
        self,
        agent_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListAgentKnowledgeResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/knowledge",
            params=params or None,
        )
        items = [_parse_agent_knowledge(k) for k in raw.get("knowledge", [])]
        return ListAgentKnowledgeResponse(knowledge=items, count=raw.get("count", 0))

    async def create_knowledge(self, agent_id: str, request: CreateAgentKnowledgeRequest) -> AgentKnowledge:
        body: dict[str, Any] = {"factType": request.fact_type, "subject": request.subject}
        if request.predicate is not None:
            body["predicate"] = request.predicate
        if request.object is not None:
            body["object"] = request.object
        if request.source is not None:
            body["source"] = request.source
        if request.confidence is not None:
            body["confidence"] = request.confidence
        if request.tags is not None:
            body["tags"] = request.tags
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/knowledge", body)
        return _parse_agent_knowledge(raw)

    async def search_knowledge(self, agent_id: str, request: dict[str, Any]) -> SearchAgentKnowledgeResponse:
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/knowledge/search", request)
        results = [_parse_agent_knowledge(k) for k in raw.get("results", [])]
        return SearchAgentKnowledgeResponse(results=results, count=raw.get("count", 0))

    async def delete_knowledge(self, agent_id: str, knowledge_id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/agents/{agent_id}/knowledge/{knowledge_id}")

    async def list_governance(self, agent_id: str) -> ListAgentGovernanceResponse:
        raw = await self._client.get(f"/api/v1/cadreen/agents/{agent_id}/governance")
        policies = [_parse_agent_governance_policy(p) for p in raw.get("policies", [])]
        return ListAgentGovernanceResponse(policies=policies, count=raw.get("count", 0))

    async def create_governance(self, agent_id: str, request: CreateAgentGovernanceRequest) -> AgentGovernancePolicy:
        body: dict[str, Any] = {"name": request.name, "rules": request.rules}
        if request.domain is not None:
            body["domain"] = request.domain
        if request.priority is not None:
            body["priority"] = request.priority
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/governance", body)
        return _parse_agent_governance_policy(raw)

    async def update_governance(
        self, agent_id: str, policy_id: str, request: UpdateAgentGovernanceRequest
    ) -> AgentGovernancePolicy:
        body: dict[str, Any] = {}
        if request.name is not None:
            body["name"] = request.name
        if request.rules is not None:
            body["rules"] = request.rules
        if request.domain is not None:
            body["domain"] = request.domain
        if request.priority is not None:
            body["priority"] = request.priority
        raw = await self._client.patch(
            f"/api/v1/cadreen/agents/{agent_id}/governance/{policy_id}", body
        )
        return _parse_agent_governance_policy(raw)

    async def delete_governance(self, agent_id: str, policy_id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/agents/{agent_id}/governance/{policy_id}")

    async def list_audit(
        self,
        agent_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListAgentAuditResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/audit",
            params=params or None,
        )
        entries = [_parse_agent_audit_entry(e) for e in raw.get("entries", [])]
        return ListAgentAuditResponse(entries=entries, count=raw.get("count", 0))

    async def start_negotiation(self, agent_id: str, request: StartNegotiationRequest) -> AgentNegotiation:
        body: dict[str, Any] = {"target_agent_id": request.target_agent_id, "proposal": request.proposal}
        if request.context is not None:
            body["context"] = request.context
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/negotiate", body)
        return _parse_agent_negotiation(raw)

    async def list_negotiations(
        self,
        agent_id: str,
        *,
        limit: int | None = None,
        offset: int | None = None,
    ) -> ListAgentNegotiationsResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/negotiations",
            params=params or None,
        )
        negotiations = [_parse_agent_negotiation(n) for n in raw.get("negotiations", [])]
        return ListAgentNegotiationsResponse(negotiations=negotiations, count=raw.get("count", 0))

    async def get_negotiation(self, agent_id: str, negotiation_id: str) -> AgentNegotiation:
        raw = await self._client.get(
            f"/api/v1/cadreen/agents/{agent_id}/negotiations/{negotiation_id}"
        )
        return _parse_agent_negotiation(raw)

    async def respond_to_negotiation(
        self, agent_id: str, negotiation_id: str, request: RespondNegotiationRequest
    ) -> AgentNegotiation:
        body: dict[str, Any] = {"response": request.response}
        if request.counter_proposal is not None:
            body["counter_proposal"] = request.counter_proposal
        raw = await self._client.post(
            f"/api/v1/cadreen/agents/{agent_id}/negotiations/{negotiation_id}/respond", body
        )
        return _parse_agent_negotiation(raw)


def _parse_agent(raw: dict[str, Any]) -> Agent:
    return Agent(
        id=raw["id"],
        name=raw["name"],
        status=raw.get("status", ""),
        created_at=raw.get("created_at", ""),
        updated_at=raw.get("updated_at", ""),
        description=raw.get("description"),
        model=raw.get("model"),
        domain=raw.get("domain"),
        tags=raw.get("tags"),
        version=raw.get("version"),
    )


def _parse_agent_message(raw: dict[str, Any]) -> AgentMessage:
    return AgentMessage(
        id=raw["id"],
        role=raw["role"],
        content=raw["content"],
        created_at=raw.get("created_at", ""),
        metadata=raw.get("metadata"),
    )


def _parse_agent_knowledge(raw: dict[str, Any]) -> AgentKnowledge:
    return AgentKnowledge(
        id=raw["id"],
        fact_type=raw.get("factType", raw.get("fact_type", "")),
        subject=raw.get("subject", ""),
        predicate=raw.get("predicate"),
        object=raw.get("object"),
        source=raw.get("source"),
        confidence=raw.get("confidence"),
        created_at=raw.get("created_at", ""),
        agent_id=raw.get("agent_id"),
        visibility=raw.get("visibility"),
        updated_at=raw.get("updated_at"),
        domain=raw.get("domain"),
        tags=raw.get("tags"),
        authority=raw.get("authority"),
    )


def _parse_agent_governance_policy(raw: dict[str, Any]) -> AgentGovernancePolicy:
    return AgentGovernancePolicy(
        id=raw["id"],
        name=raw["name"],
        rules=raw.get("rules", []),
        created_at=raw.get("created_at", ""),
        domain=raw.get("domain"),
        priority=raw.get("priority"),
    )


def _parse_agent_audit_entry(raw: dict[str, Any]) -> AgentAuditEntry:
    return AgentAuditEntry(
        id=raw["id"],
        action=raw["action"],
        timestamp=raw.get("timestamp", ""),
        actor=raw.get("actor", ""),
        detail=raw.get("detail"),
        policy_id=raw.get("policy_id"),
    )


def _parse_agent_negotiation(raw: dict[str, Any]) -> AgentNegotiation:
    return AgentNegotiation(
        id=raw["id"],
        status=raw["status"],
        initiator_agent_id=raw.get("initiator_agent_id", ""),
        target_agent_id=raw.get("target_agent_id", ""),
        proposal=raw.get("proposal", {}),
        created_at=raw.get("created_at", ""),
        updated_at=raw.get("updated_at", ""),
        response=raw.get("response"),
        counter_proposal=raw.get("counter_proposal"),
    )
