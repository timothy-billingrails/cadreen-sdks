from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Literal, Optional, TypedDict, Union


class ExecutionData(TypedDict, total=False):
    """Typed dict for execution data in IntentResult."""
    id: str
    status: str
    stream_url: Optional[str]
    poll_url: Optional[str]


HealthStatus = Literal["healthy", "degraded", "unhealthy", "unknown", "latent"]
ConnectorType = Literal["mcp", "openapi", "composio", "native_rest", "utp", "builtin", "zenrows", "bug0", "instavm", "direct"]
CapabilitySource = Literal["detected", "built_in", "mcp", "openapi", "composio", "native_rest", "utp", "builtin", "direct"]
TransportType = Literal["http", "sse", "stdio"]
EscalationStatus = Literal["pending", "resolved", "rejected"]
CredentialType = Literal["api_key", "bearer", "basic", "oauth2", "session"]
AtomScope = Literal["tenant", "global", "personal"]
AtomCategory = Literal["reference", "preference", "episode", "precedent", "note", "project"]
AtomType = Literal[
    "reference", "preference", "episode", "precedent", "note", "project",
    "policy", "procedure", "decision", "fact", "metric",
    "constraint", "event", "observation", "error", "mission", "module",
    "answer", "image", "video", "audio", "visualization", "code", "test",
    "dataset", "document", "research", "prompt", "billing_event", "opinion",
    "instruction", "definition", "question", "tool_invocation", "tool_failure_pattern",
]
ErrorCategory = Literal[
    "auth", "network", "resource", "logic", "config", "dependency",
    "capability", "rate_limit", "timeout", "validation", "parsing",
    "external", "unknown",
]
RecoveryStrategyType = Literal[
    "retry", "sub_execution", "human_handoff", "skip", "reconfigure",
    "regenerate", "coerce", "repair", "re_decide", "try_alternative", "none",
]
StackItemSource = Literal["user_data", "cadreen", "connector", "gap", "detected", "built_in"]
StackItemStatus = Literal[
    "existing", "ready", "healthy", "degraded", "unhealthy", "unknown",
    "pending_auth", "connected", "registered",
]
GovernanceDecisionType = Literal["auto", "auto_complete", "handoff", "escalate", "clarify_requester", "abstain"]
RecoveryStatus = Literal["diagnosing", "recovering", "sub_execution", "escalating", "recovered", "failed", "skipped"]
IntentMode = Literal["auto", "chat", "execution"]
GapSeverity = Literal["blocking", "high", "medium", "low", "optional"]
AssessmentQuality = Literal["insufficient_data", "partial", "complete"]


@dataclass
class Pagination:
    limit: int
    offset: int
    has_more: bool


@dataclass
class Pathway:
    id: str
    capability: str
    connector: ConnectorType
    transport: TransportType
    health: HealthStatus
    tool_id: str


@dataclass
class ConnectionGroup:
    capability: str
    pathways: Optional[list[Pathway]] = None
    status: HealthStatus = "unknown"


@dataclass
class ListConnectionsResponse:
    connections: list[ConnectionGroup]
    total_capabilities: int
    total_pathways: int
    pagination: Optional[Pagination] = None


@dataclass
class AtomContent:
    text: Optional[str] = None
    source: Optional[str] = None
    subject: Optional[str] = None
    constraint: Optional[str] = None
    query: Optional[str] = None
    tools_used: Optional[list[str]] = None
    outcome: Optional[str] = None
    situation: Optional[str] = None
    action: Optional[str] = None
    result: Optional[str] = None
    name: Optional[str] = None
    constraints: Optional[list[str]] = None
    deadline: Optional[str] = None
    is_private: Optional[bool] = None


@dataclass
class Atom:
    id: str
    type: str
    domain: str
    authority: int
    version: int
    scope: Optional[AtomScope] = None
    content: Optional[AtomContent] = None
    tags: Optional[list[str]] = None
    created_at: Optional[str] = None


@dataclass
class CreateMemoryResponse:
    id: str
    type: str
    domain: str
    authority: int
    version: int
    scope: Optional[AtomScope] = None
    content: Optional[AtomContent] = None
    indexed: Optional[bool] = None
    tags: Optional[list[str]] = None
    created_at: Optional[str] = None


@dataclass
class SearchMemoryResponse:
    results: list[Atom]
    count: int


@dataclass
class MemoryTypesResponse:
    type_values: list[str]
    kind_values: list[str]
    description: str


@dataclass
class Policy:
    id: str
    name: str
    domain: str
    priority: int
    requires_human: bool
    approver_role: Optional[str] = None
    sla_hours: Optional[int] = None
    rationale: Optional[str] = None


@dataclass
class PolicyBundle:
    id: str
    version: int
    name: str
    policies: list[Policy]
    created_at: Optional[str] = None


@dataclass
class GovernanceDecision:
    type: GovernanceDecisionType
    confidence: float
    reason: str


@dataclass
class EvaluatePolicyResponse:
    action: str
    domain: str
    result: GovernanceDecision


@dataclass
class CreatePolicyResponse:
    id: str
    name: str
    version: int
    status: str
    confirmation_required: Optional[bool] = None
    approve_url: Optional[str] = None


@dataclass
class ConfirmPolicyResponse:
    id: str
    version: int
    status: str
    previous_version: Optional[int] = None
    already_active: Optional[bool] = None
    confirmed_at: Optional[str] = None


@dataclass
class Escalation:
    id: str
    status: EscalationStatus
    intent: Optional[str] = None
    category: Optional[str] = None
    execution_id: Optional[str] = None
    tool_name: Optional[str] = None
    error_message: Optional[str] = None
    severity: Optional[str] = None
    human_prompt: Optional[str] = None
    suggestions: Optional[list[str]] = None
    created_at: Optional[str] = None
    resolved_at: Optional[str] = None
    resolved_by: Optional[str] = None
    resolution: Optional[str] = None


@dataclass
class ListEscalationsResponse:
    escalations: list[Escalation]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class CredentialMetadata:
    id: str
    provider: str
    credential_name: str
    is_active: bool
    has_credential_data: bool
    type: Optional[CredentialType] = None


@dataclass
class ListCredentialsResponse:
    credentials: list[CredentialMetadata]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class CapabilityMatch:
    name: str
    human_name: Optional[str] = None
    description: Optional[str] = None
    score: Optional[float] = None
    matched_on: Optional[list[str]] = None
    health: Optional[HealthStatus] = None
    source: Optional[CapabilitySource] = None
    status: Optional[HealthStatus] = None
    functions: Optional[list[str]] = None
    category: Optional[str] = None


@dataclass
class Gap:
    capability: str
    severity: GapSeverity
    blocking: bool
    reason: Optional[str] = None
    description: Optional[str] = None
    source: Optional[str] = None


@dataclass
class ListCapabilitiesResponse:
    available: list[CapabilityMatch]
    count: int
    gaps: Optional[list[Gap]] = None
    pagination: Optional[Pagination] = None


@dataclass
class ListPoliciesResponse:
    policies: list[Policy]
    version: Optional[int] = None
    pagination: Optional[Pagination] = None


@dataclass
class Assessment:
    task: str
    can_do: float
    assessment_quality: AssessmentQuality
    ready_capabilities: int
    total_capabilities: int
    gap_count: int
    ready_for_deployment: bool
    capabilities: Optional[list[CapabilityMatch]] = None
    gaps: Optional[list[Gap]] = None
    gap_filling_tasks: Optional[list[Any]] = None
    blocking_gaps: Optional[int] = None
    policies_recommended: Optional[list[PolicyRecommendation]] = None
    needs_clarification: Optional[list[str]] = None
    stack: Optional[StackBreakdown] = None
    governance_decision: Optional[GovernanceDecision] = None
    outcomes: Optional[list[Outcome]] = None
    intelligence: Optional[IntelligenceMeta] = None


@dataclass
class Outcome:
    title: str
    description: str
    confidence: float
    ready: bool
    blocked_by: Optional[list[str]] = None


@dataclass
class PolicyRecommendation:
    policy: str
    reason: str
    action: str
    blocking: bool


@dataclass
class StackItem:
    name: str
    type: Optional[str] = None
    source: Optional[StackItemSource] = None
    status: Optional[StackItemStatus] = None
    description: Optional[str] = None
    contains: Optional[list[str]] = None
    functions: Optional[list[str]] = None


@dataclass
class StackBreakdown:
    user_data: Optional[list[StackItem]] = None
    cadreen: Optional[list[StackItem]] = None
    connectors: Optional[list[StackItem]] = None
    gaps: Optional[list[StackItem]] = None


@dataclass
class CapabilityTrace:
    total_available: int
    healthy_count: int
    active_integrations: Optional[list[str]] = None


@dataclass
class ReasoningTrace:
    capability_matches: Optional[int] = None


@dataclass
class MemoryTrace:
    healthy: bool
    knowledge_queried: Optional[int] = None


@dataclass
class GovernanceTrace:
    active: bool
    decision: Optional[str] = None
    confidence: Optional[float] = None
    reason_code: Optional[str] = None
    policy_id: Optional[str] = None
    next_actions: Optional[list] = None


@dataclass
class HumilityTrace:
    gaps_detected: Optional[int] = None
    blocking: Optional[int] = None


@dataclass
class ProcessTrace:
    started_at: str
    duration_ms: int
    components: Optional[dict[str, bool]] = None


@dataclass
class IntelligenceStage:
    name: str
    status: str
    duration_ms: Optional[int] = None
    detail: Optional[str] = None
    inputs: Optional[dict[str, object]] = None
    outputs: Optional[dict[str, object]] = None


@dataclass
class FieldStability:
    stable: list[str]
    evolving: list[str]
    internal: list[str]


@dataclass
class NextAction:
    type: str
    label: str
    reason: str
    endpoint: Optional[str] = None


@dataclass
class IntelligenceMeta:
    capability: CapabilityTrace
    reasoning: ReasoningTrace
    memory: MemoryTrace
    governance: GovernanceTrace
    humility: HumilityTrace
    process: ProcessTrace
    field_stability: FieldStability
    version: Optional[str] = None
    summary: Optional[str] = None
    next_action: Optional[NextAction] = None
    stages: Optional[list[IntelligenceStage]] = None

    def requires_human(self) -> bool:
        if self.governance.decision == "blocked" and (self.humility.blocking or 0) > 0:
            return True
        if self.stages:
            for stage in self.stages:
                if stage.name in ("escalation", "human_handoff") and stage.status != "skipped":
                    return True
        return False

    def handoff_reason(self) -> str | None:
        if self.stages:
            for stage in self.stages:
                if stage.name in ("escalation", "human_handoff") and stage.status != "skipped":
                    return stage.detail
        if self.next_action and self.next_action.reason:
            return self.next_action.reason
        return None

    def explain(self) -> TraceExplain:
        steps: list[str] = []
        if self.capability.total_available > 0:
            steps.append(f"{self.capability.healthy_count}/{self.capability.total_available} capabilities healthy")
        if self.governance.active:
            steps.append(f"Governance: {self.governance.decision or 'active'} (confidence: {self.governance.confidence or 0})")
        if self.humility.gaps_detected and self.humility.gaps_detected > 0:
            steps.append(f"{self.humility.gaps_detected} gaps detected ({self.humility.blocking or 0} blocking)")
        if self.memory.knowledge_queried:
            steps.append(f"{self.memory.knowledge_queried} knowledge items queried")
        if self.stages:
            for stage in self.stages:
                line = f"{stage.name}: {stage.status}"
                if stage.detail:
                    line += f" — {stage.detail}"
                steps.append(line)
        recommendations: list[str] = []
        if self.humility.blocking and self.humility.blocking > 0:
            recommendations.append("Resolve blocking capability gaps before proceeding")
        if not self.memory.healthy:
            recommendations.append("Knowledge store is degraded — expect reduced context quality")
        if self.requires_human():
            recommendations.append("Human intervention required")
        return TraceExplain(
            summary=self.summary or f"Intelligence: {'; '.join(steps)}",
            steps=steps,
            recommendations=recommendations if recommendations else None,
        )


@dataclass
class IntelligenceTraceEntry:
    id: str
    domain: str
    request_path: str
    request_method: str
    meta: IntelligenceMeta
    created_at: Optional[str] = None

    @property
    def intelligence(self) -> IntelligenceMeta:
        return self.meta


@dataclass
class ListIntelligenceResponse:
    traces: list[IntelligenceTraceEntry]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class IntelligenceStats:
    traces_24h: int
    traces_7d: int
    traces_30d: int
    avg_confidence_by_domain: dict[str, float]
    gap_detection_rate: float
    governance_decisions: dict[str, int]


@dataclass
class ReplayResult:
    trace_id: str
    mode: str
    domain: str
    original_gate: str
    original_confidence: float
    current_gate: str
    current_confidence: float
    gate_changed: bool
    change_summary: str
    current_capability: dict[str, Any]
    current_memory: dict[str, Any]
    current_gaps: dict[str, Any]
    replay_note: str


@dataclass
class HandoffPacket:
    trace_id: str
    domain: str
    created_at: str
    governance: dict[str, Any]
    what_the_system_knew: dict[str, Any]
    what_the_system_didnt_know: dict[str, Any]
    what_happened: dict[str, Any]
    suggested_actions: list[dict[str, Any]]
    next_action: dict[str, Any]
    trace_url: str


@dataclass
class PromoteResult:
    id: str
    kind: str
    status: str
    source_trace_id: str
    tool_name: Optional[str] = None
    tool_sequence: Optional[list[str]] = None


@dataclass
class IntentMessage:
    role: Literal["user", "assistant", "system"]
    content: str


@dataclass
class IntentContext:
    existing_connectors: Optional[list[str]] = None
    constraints: Optional[dict[str, object]] = None
    domain: Optional[str] = None


@dataclass
class IntentRequest:
    messages: list[IntentMessage]
    conversation_id: Optional[str] = None
    context: Optional[IntentContext] = None
    mode: Optional[IntentMode] = None
    stream: Optional[bool] = None
    dry_run: Optional[bool] = None
    user_id: Optional[str] = None


@dataclass
class ResponseMessage:
    role: str
    content: str
    metadata: Optional[dict[str, object]] = None


@dataclass
class ClarificationQuestion:
    id: str
    question: str
    type: str
    required: bool
    reason: Optional[str] = None


@dataclass
class TraceExplain:
    summary: str
    steps: list[str]
    recommendations: Optional[list[str]] = None


@dataclass
class ExecutionEvent:
    type: str
    data: dict[str, object]


@dataclass
class ExecutionStatus:
    id: str
    status: str
    progress: Optional[float] = None
    result: Optional[dict[str, object]] = None
    error: Optional[str] = None


@dataclass
class RegisterOpenAPIRequest:
    name: str
    spec_url: Optional[str] = None
    spec_content: Optional[str] = None
    credential_id: Optional[str] = None


@dataclass
class RegisterOpenAPIResponse:
    id: str
    name: str
    type: str
    status: str
    tools_generated: Optional[int] = None
    tools_registered: Optional[int] = None
    functions: Optional[list[str]] = None
    spec_url: Optional[str] = None


@dataclass
class RegisterMCPRequest:
    name: str
    url: str
    transport: Optional[TransportType] = None
    headers: Optional[dict[str, str]] = None


@dataclass
class RegisterMCPResponse:
    id: str
    name: str
    type: str
    status: str
    transport: Optional[str] = None
    url: Optional[str] = None


@dataclass
class SearchMemoryRequest:
    query: str
    domain: Optional[str] = None
    tag: Optional[str] = None
    limit: Optional[int] = None


@dataclass
class RememberRequest:
    type: AtomType
    content: dict[str, object]
    domain: Optional[str] = None
    scope: Optional[AtomScope] = None
    authority: Optional[int] = None
    tags: Optional[list[str]] = None


@dataclass
class CreatePolicyRequest:
    name: str
    rules: Optional[list[dict[str, object]]] = None
    domain: Optional[str] = None
    auto_draft: Optional[bool] = None


@dataclass
class EvaluatePolicyRequest:
    action: str
    domain: Optional[str] = None
    context: Optional[dict[str, object]] = None


@dataclass
class RequestOptions:
    idempotency_key: Optional[str] = None
    headers: Optional[dict[str, str]] = None


@dataclass
class CadreenConfig:
    api_key: str
    base_url: Optional[str] = None
    max_retries: Optional[int] = None
    timeout: Optional[int] = None
    telemetry: Any = None
    sandbox: bool = False
    fixtures: Optional[dict[str, Any]] = None
    profile: Literal["lean", "audit", "full"] = "full"


IntentResultType = Literal["direct", "clarify", "execution", "blocked", "connect_required"]


@dataclass
class DirectResult:
    type: Literal["direct"]
    message: ResponseMessage
    intelligence: IntelligenceMeta
    trace_id: str

    def explain(self) -> str:
        return self.message.content

    @property
    def ready(self) -> bool:
        return True

    @property
    def needs(self) -> list[str]:
        return []

    @property
    def next(self) -> str:
        na = getattr(self.intelligence, "next_action", None)
        if na and getattr(na, "type", "none") != "none":
            return na.label
        return "done"


@dataclass
class ClarifyResult:
    type: Literal["clarify"]
    questions: list[ClarificationQuestion]
    conversation_id: str
    intelligence: IntelligenceMeta
    trace_id: str

    def explain(self) -> str:
        return "Clarification needed: " + "; ".join(q.question for q in self.questions)

    @property
    def ready(self) -> bool:
        return False

    @property
    def needs(self) -> list[str]:
        return [q.question for q in self.questions]

    @property
    def next(self) -> str:
        n = len(self.questions)
        return f"answer {n} question{'s' if n != 1 else ''}"


@dataclass
class ExecutionResult:
    type: Literal["execution"]
    execution: ExecutionData
    intelligence: IntelligenceMeta
    trace_id: str

    def explain(self) -> str:
        return f"Execution started: {self.execution['id']}"

    @property
    def ready(self) -> bool:
        return True

    @property
    def needs(self) -> list[str]:
        return []

    @property
    def next(self) -> str:
        e = self.execution
        if isinstance(e, dict) and e.get("stream_url"):
            return "stream execution"
        return f"poll {e.get('poll_url', e.get('id', ''))}"


@dataclass
class BlockedResult:
    type: Literal["blocked"]
    intelligence: IntelligenceMeta
    trace_id: str
    reason_code: Optional[str] = None
    policy_id: Optional[str] = None
    status: str = ""

    def explain(self) -> str:
        return self.status or f"Blocked by policy: {self.reason_code or 'governance gate'}"

    @property
    def ready(self) -> bool:
        return False

    @property
    def needs(self) -> list[str]:
        return [f"blocked: {self.reason_code or 'governance gate'}"]

    @property
    def next(self) -> str:
        return self.policy_id or "resolve policy block"


@dataclass
class ConnectRequiredResult:
    type: Literal["connect_required"]
    intelligence: IntelligenceMeta
    trace_id: str
    endpoint: Optional[str] = None
    reason: Optional[str] = None
    status: str = ""
    next_action: Optional[dict] = None

    def explain(self) -> str:
        return self.status or f"Connection required: {self.endpoint}"

    @property
    def ready(self) -> bool:
        return False

    @property
    def needs(self) -> list[str]:
        return [f"connect {self.endpoint}"] if self.endpoint else []

    @property
    def next(self) -> str:
        return self.endpoint or ""


IntentResult = Union[DirectResult, ClarifyResult, ExecutionResult, BlockedResult, ConnectRequiredResult]

ConnectResultType = Literal["prebuilt", "schema_required", "manual", "unknown"]


@dataclass
class ConnectPrebuiltDetail:
    tool_id: str
    tool_name: str
    service_id: str
    service_name: str
    auth_type: str
    source: str
    account_id: Optional[str] = None


@dataclass
class ConnectSchemaRequiredDetail:
    tool_id: str
    tool_name: str
    auth_url: str
    connector: str


@dataclass
class ConnectPathway:
    id: str
    connector: str
    tool_id: str
    health: str
    priority: int


@dataclass
class ConnectManualDetail:
    pathways: list[ConnectPathway]


@dataclass
class ConnectUnknownDetail:
    searched: str
    hints: Optional[list[str]] = None


@dataclass
class ConnectResult:
    type: ConnectResultType
    capability: str
    detail: Union[ConnectPrebuiltDetail, ConnectSchemaRequiredDetail, ConnectManualDetail, ConnectUnknownDetail]


def intent_status(result: IntentResult) -> dict:
    return {"ready": result.ready, "needs": result.needs, "next": result.next}


@dataclass
class SetupConnection:
    capability: str


@dataclass
class SetupCredential:
    provider: str
    key_data: dict
    name: Optional[str] = None


@dataclass
class SetupMemory:
    content: dict
    type: Optional[str] = None
    domain: Optional[str] = None
    tags: Optional[list[str]] = None
    authority: Optional[int] = None


@dataclass
class SetupPolicy:
    name: str
    rule: str
    description: Optional[str] = None
    severity: Optional[str] = None


@dataclass
class SetupProposal:
    type: str
    description: str
    detail: str


@dataclass
class SetupRequest:
    connections: Optional[list[SetupConnection]] = None
    credentials: Optional[list[SetupCredential]] = None
    memory: Optional[list[SetupMemory]] = None
    policies: Optional[list[SetupPolicy]] = None
    workspace_id: Optional[str] = None
    purpose: Optional[str] = None
    examples: Optional[list[str]] = None
    constraints: Optional[list[str]] = None
    confirm: Optional[bool] = None
    dry_run: Optional[bool] = None


@dataclass
class SetupConnectionResult:
    capability: str
    status: str
    detail: Optional[Any] = None
    error: Optional[str] = None


@dataclass
class SetupCredentialResult:
    provider: str
    name: str
    status: str
    id: Optional[str] = None
    error: Optional[str] = None


@dataclass
class SetupMemoryResult:
    id: str
    type: str
    classified: bool
    status: str
    kind: Optional[str] = None
    error: Optional[str] = None


@dataclass
class SetupPolicyResult:
    name: str
    status: str
    id: Optional[str] = None
    error: Optional[str] = None


@dataclass
class SetupResult:
    connections: list[SetupConnectionResult]
    credentials: list[SetupCredentialResult]
    memory: list[SetupMemoryResult]
    policies: list[SetupPolicyResult]
    applied: int
    failed: int
    workspace_id: Optional[str] = None
    proposals: Optional[list[SetupProposal]] = None
    notice: Optional[str] = None
    dry_run: Optional[bool] = None


@dataclass
class SetupSessionCreateRequest:
    purpose: Optional[str] = None
    constraints: Optional[list[str]] = None
    workspace_id: Optional[str] = None


@dataclass
class SetupSessionAddRequest:
    connections: Optional[list[dict[str, Any]]] = None
    credentials: Optional[list[dict[str, Any]]] = None
    memory: Optional[list[dict[str, Any]]] = None
    policies: Optional[list[dict[str, Any]]] = None


@dataclass
class SetupSessionApplyRequest:
    confirm: bool = True


@dataclass
class SetupSession:
    id: str
    workspace_id: str
    status: str
    created_at: str
    updated_at: str
    connections: list[dict[str, Any]] = field(default_factory=list)
    credentials: list[dict[str, Any]] = field(default_factory=list)
    memory: list[dict[str, Any]] = field(default_factory=list)
    policies: list[dict[str, Any]] = field(default_factory=list)
    applied_count: int = 0
    failed_count: int = 0
    purpose: Optional[str] = None
    constraints: Optional[list[str]] = None
    proposals: Optional[list[dict[str, Any]]] = None
    applied_at: Optional[str] = None


@dataclass
class SetupSessionApplyResult:
    session_id: str
    status: str
    applied: int
    failed: int
    result: Optional[dict[str, Any]] = None


# ---------------------------------------------------------------------------
# Abstraction aliases: Atom → MemoryItem
# ---------------------------------------------------------------------------

MemoryType = AtomCategory


# ---------------------------------------------------------------------------
# Marketplace catalog types
# ---------------------------------------------------------------------------

@dataclass
class CatalogIntegration:
    id: str
    name: str
    description: str
    category: str
    provider: str
    status: str
    auth_type: str
    install_time: str
    capabilities: Optional[list[str]] = None
    tags: Optional[list[str]] = None
    popularity: int = 0
    featured: bool = False


@dataclass
class CatalogCategory:
    name: str
    description: str
    integrations: list[CatalogIntegration]


@dataclass
class CatalogResponse:
    categories: list[CatalogCategory]
    installed: list[str]
    total_available: int


@dataclass
class InstallResponse:
    status: str
    provider: str
    auth_url: Optional[str] = None
    estimated_time: Optional[str] = None


# ---------------------------------------------------------------------------
# Document types
# ---------------------------------------------------------------------------


@dataclass
class Document:
    id: str
    name: str
    content_type: Optional[str] = None
    size: Optional[int] = None
    status: Optional[str] = None
    created_at: Optional[str] = None


@dataclass
class ListDocumentsResponse:
    documents: list[Document]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class UploadDocumentResponse:
    id: str
    name: str
    content_type: Optional[str] = None
    size: Optional[int] = None
    status: Optional[str] = None


# ---------------------------------------------------------------------------
# Webhook types
# ---------------------------------------------------------------------------


@dataclass
class Webhook:
    id: str
    url: str
    is_active: bool
    events: Optional[list[str]] = None
    secret: Optional[str] = None
    created_at: Optional[str] = None


@dataclass
class ListWebhooksResponse:
    webhooks: list[Webhook]
    count: int
    pagination: Optional[Pagination] = None


# ---------------------------------------------------------------------------
# Learning types
# ---------------------------------------------------------------------------


@dataclass
class LearningPattern:
    id: str
    pattern: str
    confidence: float
    occurrences: Optional[int] = None
    domain: Optional[str] = None
    tags: Optional[list[str]] = None
    created_at: Optional[str] = None


@dataclass
class LearningEpisode:
    id: str
    description: str
    outcome: Optional[str] = None
    trace_id: Optional[str] = None
    domain: Optional[str] = None
    created_at: Optional[str] = None


@dataclass
class LearningSuggestion:
    id: str
    type: str
    description: str
    impact: Optional[str] = None
    domain: Optional[str] = None


@dataclass
class ListLearningPatternsResponse:
    patterns: list[LearningPattern]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class ListLearningEpisodesResponse:
    episodes: list[LearningEpisode]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class ListLearningSuggestionsResponse:
    suggestions: list[LearningSuggestion]
    count: int
    pagination: Optional[Pagination] = None


# ---------------------------------------------------------------------------
# Healing types
# ---------------------------------------------------------------------------


@dataclass
class HealingDiagnosis:
    error_category: Optional[str] = None
    semantic_reason: Optional[str] = None
    root_cause: Optional[str] = None
    can_retry: Optional[bool] = None
    needs_sub_execution: Optional[bool] = None
    needs_human: Optional[bool] = None
    should_skip: Optional[bool] = None
    needs_re_decide: Optional[bool] = None
    needs_try_alternative: Optional[bool] = None
    retry_delay_ms: Optional[int] = None
    confidence: Optional[float] = None


@dataclass
class StrategyCount:
    strategy: str
    count: int


@dataclass
class ToolHealingStats:
    tool_name: str
    total: int
    successful: int
    failed: int
    success_rate: float
    top_strategy: Optional[str] = None


@dataclass
class TimeRange:
    first_precedent: Optional[str] = None
    last_precedent: Optional[str] = None


@dataclass
class HealingStatsResponse:
    total_precedents: Optional[int] = None
    successful_recoveries: Optional[int] = None
    failed_recoveries: Optional[int] = None
    success_rate: Optional[float] = None
    avg_duration_ms: Optional[int] = None
    common_strategies: Optional[list[StrategyCount]] = None
    top_tools: Optional[list[ToolHealingStats]] = None
    by_category: Optional[dict[str, Any]] = None
    time_range: Optional[TimeRange] = None


@dataclass
class HealingPrecedent:
    id: str
    error_type: str
    success: bool
    attempts: int
    confidence: float
    tool_name: Optional[str] = None
    error_category: Optional[str] = None
    semantic_reason: Optional[str] = None
    root_cause: Optional[str] = None
    recovery_strategy: Optional[str] = None
    what_worked: Optional[str] = None
    what_failed: Optional[str] = None
    duration_ms: Optional[int] = None
    created_at: Optional[str] = None
    domain: Optional[str] = None
    tags: Optional[list[str]] = None


@dataclass
class ListHealingPrecedentsResponse:
    precedents: list[HealingPrecedent]
    count: int
    pagination: Optional[Pagination] = None


@dataclass
class CreateCredentialRequest:
    provider: str
    key_data: dict
    name: Optional[str] = None


@dataclass
class ResolveEscalationRequest:
    decision: str


@dataclass
class DiagnoseRequest:
    error_message: str
    tool_name: Optional[str] = None
    trace_id: Optional[str] = None


# ---------------------------------------------------------------------------
# Proposal types
# ---------------------------------------------------------------------------


@dataclass
class ProposalEvidence:
    description: str
    source: Optional[str] = None
    count: Optional[int] = None
    confidence: Optional[float] = None


@dataclass
class TaskProposal:
    id: str
    title: str
    description: str
    intent: str
    proposal_type: str
    trigger_type: str
    trigger_source: str
    confidence: float
    priority: int
    status: str
    created_at: str
    domain: Optional[str] = None
    mission_intent: Optional[str] = None
    trigger_details: Optional[str] = None
    evidence: Optional[list[ProposalEvidence]] = None
    expires_at: Optional[str] = None
    accepted_at: Optional[str] = None
    dismissed_at: Optional[str] = None
    dismissal_reason: Optional[str] = None
    execution_id: Optional[str] = None
    dedup_key: Optional[str] = None
    requires_review: Optional[bool] = None


@dataclass
class ListProposalsResponse:
    proposals: list[TaskProposal]
    count: int


@dataclass
class AcceptProposalResponse:
    status: str
    execution_id: str
    action: str
    intent: str
    next_step: str
    auto_approved: Optional[bool] = None
    result: Optional[dict] = None


@dataclass
class DismissProposalResponse:
    status: str


@dataclass
class ProposalStatsResponse:
    proposed: int = 0
    accepted: int = 0
    dismissed: int = 0
    expired: int = 0


WorkspaceRole = Literal["admin", "operator", "member", "viewer"]


@dataclass
class WorkspaceUser:
    id: str
    workspace_id: str
    user_id: str
    role: str
    invited_at: str
    created_at: str
    updated_at: str
    invited_by: Optional[str] = None


@dataclass
class InviteUserRequest:
    email: str
    role: Optional[str] = None


@dataclass
class UpdateRoleRequest:
    role: str


@dataclass
class ListWorkspaceUsersResponse:
    users: list[WorkspaceUser]
    count: int
