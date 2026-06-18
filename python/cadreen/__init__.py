from .client import CadreenError
from .types import (
    CadreenConfig,
    HealthStatus,
    ConnectorType,
    CapabilitySource,
    TransportType,
    EscalationStatus,
    CredentialType,
    AtomScope,
    AtomType,
    AtomCategory,
    MemoryTypesResponse,
    ErrorCategory,
    RecoveryStrategyType,
    StackItemSource,
    StackItemStatus,
    GovernanceDecisionType,
    RecoveryStatus,
    IntentMode,
    Pagination,
    Pathway,
    ConnectionGroup,
    ListConnectionsResponse,
    AtomContent,
    Atom,
    CreateMemoryResponse,
    SearchMemoryResponse,
    Policy,
    PolicyBundle,
    GovernanceDecision,
    EvaluatePolicyResponse,
    CreatePolicyResponse,
    ConfirmPolicyResponse,
    Escalation,
    ListEscalationsResponse,
    CredentialMetadata,
    ListCredentialsResponse,
    CapabilityMatch,
    ListCapabilitiesResponse,
    ListPoliciesResponse,
    Gap,
    Outcome,
    Assessment,
    PolicyRecommendation,
    StackItem,
    StackBreakdown,
    CapabilityTrace,
    ReasoningTrace,
    MemoryTrace,
    GovernanceTrace,
    HumilityTrace,
    ProcessTrace,
    IntelligenceMeta,
    IntelligenceTraceEntry,
    ListIntelligenceResponse,
    IntelligenceStats,
    ReplayResult,
    HandoffPacket,
    PromoteResult,
    IntentMessage,
    IntentContext,
    IntentRequest,
    ResponseMessage,
    TraceExplain,
    ExecutionEvent,
    ExecutionStatus,
    IntentResult,
    IntentResultType,
    DirectResult,
    ClarifyResult,
    ExecutionResult,
    BlockedResult,
    ConnectRequiredResult,
    RegisterOpenAPIRequest,
    RegisterOpenAPIResponse,
    RegisterMCPRequest,
    RegisterMCPResponse,
    SearchMemoryRequest,
    RememberRequest,
    CreatePolicyRequest,
    EvaluatePolicyRequest,
    RequestOptions,
    ConnectResult,
    ConnectResultType,
    ConnectPrebuiltDetail,
    ConnectSchemaRequiredDetail,
    ConnectManualDetail,
    ConnectPathway,
    ConnectUnknownDetail,
    intent_status,
    SetupRequest,
    SetupResult,
    SetupConnection,
    SetupCredential,
    SetupMemory,
    SetupPolicy,
    SetupConnectionResult,
    SetupCredentialResult,
    SetupMemoryResult,
    SetupPolicyResult,
    SetupProposal,
    CatalogResponse,
    CatalogCategory,
    CatalogIntegration,
    InstallResponse,
)

from .client import HttpClient
from .resources.intent import IntentResource
from .resources.memory import MemoryResource
from .resources.policies import PoliciesResource
from .resources.connections import ConnectionsResource
from .resources.traces import TracesResource
from .resources.executions import ExecutionsResource
from .resources.guardrails import GuardrailsResource
from .resources.chat import ChatResource
from .resources.chat import (
    ChatMessage,
    ChatToolCall,
    ChatFunctionCall,
    ChatToolDefinition,
    ChatFunctionDefinition,
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatChoice,
    ChatUsage,
    ChatCompletionChunk,
    ChatChunkChoice,
    ChatDelta,
    ToolEntry,
    ListToolsResponse,
)
from .redaction import redact_string, redact_trace, redact_messages, RedactOptions
from .telemetry import TelemetryProvider, TelemetrySpan, TelemetryMeter, OpenTelemetryAdapter, NoOpProvider


class Cadreen:
    def __init__(self, api_key: str, base_url: str | None = None, max_retries: int | None = None, timeout: int | None = None, *, sandbox: bool = False, fixtures: dict | None = None) -> None:
        config = CadreenConfig(api_key=api_key, base_url=base_url, max_retries=max_retries, timeout=timeout, sandbox=sandbox, fixtures=fixtures)
        self._client = HttpClient(config)
        self.intent = IntentResource(self._client)
        self.memory = MemoryResource(self._client)
        self.policies = PoliciesResource(self._client)
        self.connections = ConnectionsResource(self._client)
        self.traces = TracesResource(self._client)
        self.executions = ExecutionsResource(self._client)
        self.guardrails = GuardrailsResource(self.policies)
        self.chat = ChatResource(self._client)

    async def invoke(self, request: IntentRequest) -> IntentResult:
        return await self.intent.invoke(request)

    async def ask(
        self,
        prompt: str,
        *,
        conversation_id: str | None = None,
        context: IntentContext | None = None,
        stream: bool | None = None,
    ) -> IntentResult:
        return await self.intent.invoke(
            IntentRequest(
                messages=[IntentMessage(role="user", content=prompt)],
                mode="chat",
                conversation_id=conversation_id,
                context=context,
                stream=stream,
            )
        )

    async def act(
        self,
        prompt: str,
        *,
        conversation_id: str | None = None,
        context: IntentContext | None = None,
        stream: bool | None = None,
    ) -> IntentResult:
        return await self.intent.invoke(
            IntentRequest(
                messages=[IntentMessage(role="user", content=prompt)],
                mode="execution",
                conversation_id=conversation_id,
                context=context,
                stream=stream,
            )
        )

    async def remember(
        self,
        type: str,
        content: dict,
        *,
        domain: str | None = None,
        scope: str | None = None,
        authority: int | None = None,
        tags: list[str] | None = None,
    ) -> CreateMemoryResponse:
        return await self.memory.remember(type, content, domain=domain, scope=scope, authority=authority, tags=tags)

    async def context(self, query: str, *, domain: str | None = None, tag: str | None = None, limit: int | None = None) -> SearchMemoryResponse:
        return await self.memory.search(query, domain=domain, tag=tag, limit=limit)

    async def connect(self, capability: str) -> ConnectResult:
        return await self.connections.connect(capability)

    async def setup(self, request: SetupRequest) -> SetupResult:
        payload: dict = {}
        if request.workspace_id:
            payload["workspace_id"] = request.workspace_id
        if request.purpose:
            payload["purpose"] = request.purpose
        if request.examples:
            payload["examples"] = request.examples
        if request.constraints:
            payload["constraints"] = request.constraints
        if request.connections:
            payload["connections"] = [{"capability": c.capability} for c in request.connections]
        if request.credentials:
            payload["credentials"] = [{"provider": c.provider, "name": c.name, "key_data": c.key_data} for c in request.credentials]
        if request.memory:
            mem_items = []
            for m in request.memory:
                item: dict = {"content": m.content}
                if m.type:
                    item["type"] = m.type
                if m.domain:
                    item["domain"] = m.domain
                if m.tags:
                    item["tags"] = m.tags
                if m.authority:
                    item["authority"] = m.authority
                mem_items.append(item)
            payload["memory"] = mem_items
        if request.policies:
            payload["policies"] = [{"name": p.name, "rule": p.rule, "description": p.description, "severity": p.severity} for p in request.policies]
        resp = await self._client.post("/api/v1/cadreen/setup", payload)
        conns = [SetupConnectionResult(capability=c["capability"], status=c["status"], detail=c.get("detail"), error=c.get("error")) for c in resp.get("connections", [])]
        creds = [SetupCredentialResult(provider=c["provider"], name=c.get("name", ""), status=c["status"], id=c.get("id"), error=c.get("error")) for c in resp.get("credentials", [])]
        mems = [SetupMemoryResult(id=m["id"], type=m["type"], classified=m["classified"], status=m["status"], kind=m.get("kind"), error=m.get("error")) for m in resp.get("memory", [])]
        pols = [SetupPolicyResult(name=p["name"], status=p["status"], id=p.get("id"), error=p.get("error")) for p in resp.get("policies", [])]
        props = [SetupProposal(type=p["type"], description=p["description"], detail=p["detail"]) for p in resp.get("proposals", [])]
        return SetupResult(
            connections=conns, credentials=creds, memory=mems, policies=pols,
            applied=resp["applied"], failed=resp["failed"],
            workspace_id=resp.get("workspace_id"),
            proposals=props or None,
        )
