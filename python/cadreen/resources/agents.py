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
    AgentExecution,
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
        if request.system_prompt is not None:
            body["system_prompt"] = request.system_prompt
        if request.capabilities is not None:
            body["capabilities"] = request.capabilities
        if request.tags is not None:
            body["tags"] = request.tags
        if request.config is not None:
            body["config"] = request.config
        if request.metadata is not None:
            body["metadata"] = request.metadata
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
        if request.system_prompt is not None:
            body["system_prompt"] = request.system_prompt
        if request.capabilities is not None:
            body["capabilities"] = request.capabilities
        if request.tags is not None:
            body["tags"] = request.tags
        if request.config is not None:
            body["config"] = request.config
        if request.metadata is not None:
            body["metadata"] = request.metadata
        raw = await self._client.patch(f"/api/v1/cadreen/agents/{agent_id}", body)
        return _parse_agent(raw)

    async def delete(self, agent_id: str) -> None:
        await self._client.delete(f"/api/v1/cadreen/agents/{agent_id}")

    async def get_config(self, agent_id: str) -> AgentConfig:
        raw = await self._client.get(f"/api/v1/cadreen/agents/{agent_id}/config")
        return AgentConfig(
            agent_id=raw.get("agent_id", agent_id),
            model=raw.get("model"),
            system_prompt=raw.get("system_prompt"),
            temperature=raw.get("temperature"),
            max_tokens=raw.get("max_tokens"),
            tools=raw.get("tools"),
            connections=raw.get("connections"),
            memory_scope=raw.get("memory_scope"),
            governance_policy_id=raw.get("governance_policy_id"),
            metadata=raw.get("metadata"),
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
            agent_id=raw.get("agent_id", agent_id),
            tools=raw.get("tools", []),
            connections=raw.get("connections", []),
            knowledge_count=raw.get("knowledge_count", 0),
            governance_policies=raw.get("governance_policies", 0),
            can_execute=raw.get("can_execute", False),
            can_federate=raw.get("can_federate", False),
            can_negotiate=raw.get("can_negotiate", False),
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
            AgentExecution(
                id=e["id"],
                agent_id=e.get("agent_id", agent_id),
                status=e["status"],
                started_at=e.get("started_at", ""),
                input=e.get("input"),
                output=e.get("output"),
                error=e.get("error"),
                completed_at=e.get("completed_at"),
            )
            for e in raw.get("executions", [])
        ]
        return ListAgentExecutionsResponse(executions=executions, count=raw.get("count", 0))

    async def create_execution(self, agent_id: str, request: CreateExecutionRequest) -> AgentExecution:
        body: dict[str, Any] = {"intent": request.intent}
        if request.context is not None:
            body["context"] = request.context
        if request.max_budget_usd is not None:
            body["maxBudgetUsd"] = request.max_budget_usd
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/executions", body)
        return AgentExecution(
            id=raw["id"],
            agent_id=raw.get("agent_id", agent_id),
            status=raw["status"],
            started_at=raw.get("started_at", ""),
            input=raw.get("input"),
            output=raw.get("output"),
            error=raw.get("error"),
            completed_at=raw.get("completed_at"),
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
        if request.visibility is not None:
            body["visibility"] = request.visibility
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
        body: dict[str, Any] = {"name": request.name, "scope": request.scope, "rules": request.rules}
        if request.description is not None:
            body["description"] = request.description
        if request.enabled is not None:
            body["enabled"] = request.enabled
        raw = await self._client.post(f"/api/v1/cadreen/agents/{agent_id}/governance", body)
        return _parse_agent_governance_policy(raw)

    async def update_governance(
        self, agent_id: str, policy_id: str, request: UpdateAgentGovernanceRequest
    ) -> AgentGovernancePolicy:
        body: dict[str, Any] = {}
        if request.name is not None:
            body["name"] = request.name
        if request.description is not None:
            body["description"] = request.description
        if request.rules is not None:
            body["rules"] = request.rules
        if request.enabled is not None:
            body["enabled"] = request.enabled
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
        body: dict[str, Any] = {"to_agent_id": request.to_agent_id, "proposal": request.proposal}
        if request.max_rounds is not None:
            body["max_rounds"] = request.max_rounds
        if request.deadline is not None:
            body["deadline"] = request.deadline
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
        body: dict[str, Any] = {"action": request.action}
        if request.counter_proposal is not None:
            body["counter_proposal"] = request.counter_proposal
        if request.reason is not None:
            body["reason"] = request.reason
        raw = await self._client.post(
            f"/api/v1/cadreen/agents/{agent_id}/negotiations/{negotiation_id}/respond", body
        )
        return _parse_agent_negotiation(raw)


def _parse_agent(raw: dict[str, Any]) -> Agent:
    return Agent(
        id=raw["id"],
        name=raw["name"],
        status=raw.get("status", ""),
        created_at=raw.get("created_at", raw.get("createdAt", "")),
        updated_at=raw.get("updated_at", raw.get("updatedAt", "")),
        description=raw.get("description"),
        model=raw.get("model"),
        health=raw.get("health"),
        system_prompt=raw.get("system_prompt", raw.get("systemPrompt")),
        capabilities=raw.get("capabilities"),
        tags=raw.get("tags"),
        config=raw.get("config"),
        metadata=raw.get("metadata"),
        deployed_at=raw.get("deployed_at", raw.get("deployedAt")),
    )


def _parse_agent_message(raw: dict[str, Any]) -> AgentMessage:
    return AgentMessage(
        id=raw["id"],
        from_agent_id=raw.get("from_agent_id", raw.get("fromAgentId", "")),
        to_agent_id=raw.get("to_agent_id", raw.get("toAgentId", "")),
        content=raw["content"],
        status=raw.get("status", ""),
        message_type=raw.get("message_type", raw.get("messageType", "")),
        created_at=raw.get("created_at", raw.get("createdAt", "")),
        context=raw.get("context"),
        response=raw.get("response"),
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
        scope=raw.get("scope", "agent"),
        rules=raw.get("rules", []),
        created_at=raw.get("created_at", raw.get("createdAt", "")),
        updated_at=raw.get("updated_at", raw.get("updatedAt", "")),
        description=raw.get("description"),
        agent_id=raw.get("agent_id", raw.get("agentId")),
        enabled=raw.get("enabled", True),
    )


def _parse_agent_audit_entry(raw: dict[str, Any]) -> AgentAuditEntry:
    return AgentAuditEntry(
        id=raw["id"],
        agent_id=raw.get("agent_id", raw.get("agentId", "")),
        action=raw["action"],
        created_at=raw.get("created_at", raw.get("createdAt", "")),
        target_type=raw.get("target_type", raw.get("targetType")),
        target_id=raw.get("target_id", raw.get("targetId")),
        details=raw.get("details"),
        policy_action=raw.get("policy_action", raw.get("policyAction")),
    )


def _parse_agent_negotiation(raw: dict[str, Any]) -> AgentNegotiation:
    return AgentNegotiation(
        id=raw["id"],
        from_agent_id=raw.get("from_agent_id", raw.get("fromAgentId", "")),
        to_agent_id=raw.get("to_agent_id", raw.get("toAgentId", "")),
        proposal=raw.get("proposal", {}),
        status=raw.get("status", ""),
        current_round=raw.get("current_round", raw.get("currentRound", 0)),
        max_rounds=raw.get("max_rounds", raw.get("maxRounds", 0)),
        created_at=raw.get("created_at", raw.get("createdAt", "")),
        updated_at=raw.get("updated_at", raw.get("updatedAt", "")),
        deadline=raw.get("deadline"),
        resolution=raw.get("resolution"),
    )
