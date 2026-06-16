export type HealthStatus = "healthy" | "degraded" | "unhealthy" | "unknown" | "latent";
export type ConnectorType = "mcp" | "openapi" | "composio" | "native_rest" | "utp" | "builtin" | "zenrows" | "bug0" | "instavm" | "direct";
export type CapabilitySource = "detected" | "built_in" | "mcp" | "openapi" | "composio" | "native_rest" | "utp" | "builtin" | "direct";
export type TransportType = "http" | "sse" | "stdio";
export type EscalationStatus = "pending" | "resolved" | "rejected";
export type CredentialType = "api_key" | "bearer" | "basic" | "oauth2" | "session";
export type AtomScope = "tenant" | "global" | "personal";
export type AtomCategory = "reference" | "preference" | "episode" | "precedent" | "note" | "project";
export type AtomType = AtomCategory | "policy" | "procedure" | "decision" | "fact" | "metric" | "constraint" | "event" | "observation" | "error" | "mission" | "module" | "answer" | "image" | "video" | "audio" | "visualization" | "code" | "test" | "dataset" | "document" | "research" | "prompt" | "billing_event" | "opinion" | "instruction" | "definition" | "question" | "tool_invocation" | "tool_failure_pattern";
export type ErrorCategory = "auth" | "network" | "resource" | "logic" | "config" | "dependency" | "capability" | "rate_limit" | "timeout" | "validation" | "parsing" | "external" | "unknown";
export type RecoveryStrategyType = "retry" | "sub_execution" | "human_handoff" | "skip" | "reconfigure" | "regenerate" | "coerce" | "repair" | "re_decide" | "try_alternative" | "none";
export type StackItemSource = "user_data" | "cadreen" | "connector" | "gap" | "detected" | "built_in";
export type StackItemStatus = "existing" | "ready" | "healthy" | "degraded" | "unhealthy" | "unknown" | "pending_auth" | "connected" | "registered";
export type GovernanceDecisionType = "auto" | "auto_complete" | "handoff" | "escalate" | "clarify_requester" | "abstain";
export type RecoveryStatus = "diagnosing" | "recovering" | "sub_execution" | "escalating" | "recovered" | "failed" | "skipped";
export type IntentMode = "auto" | "chat" | "execution";

export interface Pagination {
  limit: number;
  offset: number;
  has_more: boolean;
}

export interface Pathway {
  id: string;
  capability: string;
  connector: ConnectorType;
  transport: TransportType;
  health: HealthStatus;
  tool_id: string;
}

export interface ConnectionGroup {
  capability: string;
  pathways?: Pathway[];
  status: HealthStatus;
}

export interface ListConnectionsResponse {
  connections: ConnectionGroup[];
  total_capabilities: number;
  total_pathways: number;
  pagination?: Pagination;
}

export interface AtomContent {
  text?: string;
  source?: string;
  subject?: string;
  constraint?: string;
  query?: string;
  tools_used?: string[];
  outcome?: string;
  situation?: string;
  action?: string;
  result?: string;
  name?: string;
  constraints?: string[];
  deadline?: string;
  is_private?: boolean;
}

export interface Atom {
  id: string;
  type: string;
  domain: string;
  scope?: AtomScope;
  content?: AtomContent;
  authority: number;
  version: number;
  tags?: string[];
  created_at?: string;
}

export interface CreateMemoryResponse {
  id: string;
  type: string;
  domain: string;
  scope?: AtomScope;
  content?: AtomContent;
  authority: number;
  version: number;
  indexed?: boolean;
  tags?: string[];
  created_at?: string;
}

export interface SearchMemoryResponse {
  results: Atom[];
  count: number;
}

export interface MemoryTypesResponse {
  type_values: AtomCategory[];
  kind_values: AtomType[];
  description: string;
}

export interface Policy {
  id: string;
  name: string;
  domain: string;
  priority: number;
  requires_human: boolean;
  approver_role?: string;
  sla_hours?: number;
  rationale?: string;
}

export interface PolicyBundle {
  id: string;
  version: number;
  name: string;
  policies: Policy[];
  created_at?: string;
}

export interface GovernanceDecision {
  type: GovernanceDecisionType;
  confidence: number;
  reason: string;
}

export interface EvaluatePolicyResponse {
  action: string;
  domain: string;
  result: GovernanceDecision;
}

export interface CreatePolicyResponse {
  id: string;
  name: string;
  version: number;
  status: string;
  confirmation_required?: boolean;
  approve_url?: string;
}

export interface ConfirmPolicyResponse {
  id: string;
  version: number;
  previous_version?: number;
  status: string;
  already_active?: boolean;
  confirmed_at?: string;
}

export interface Escalation {
  id: string;
  intent?: string;
  status: EscalationStatus;
  category?: string;
  execution_id?: string;
  tool_name?: string;
  error_message?: string;
  severity?: string;
  human_prompt?: string;
  suggestions?: string[];
  created_at?: string;
  resolved_at?: string;
  resolved_by?: string;
  resolution?: string;
}

export interface ListEscalationsResponse {
  escalations: Escalation[];
  count: number;
  pagination?: Pagination;
}

export interface CredentialMetadata {
  id: string;
  provider: string;
  credential_name: string;
  type?: CredentialType;
  is_active: boolean;
  has_credential_data: boolean;
}

export interface ListCredentialsResponse {
  credentials: CredentialMetadata[];
  count: number;
  pagination?: Pagination;
}

export interface CapabilityMatch {
  name: string;
  human_name?: string;
  description?: string;
  score?: number;
  matched_on?: string[];
  health?: HealthStatus;
  source?: CapabilitySource;
  status?: HealthStatus;
  functions?: string[];
  category?: string;
}

export interface ListCapabilitiesResponse {
  available: CapabilityMatch[];
  gaps?: Gap[];
  count: number;
  pagination?: Pagination;
}

export interface ListPoliciesResponse {
  policies: Policy[];
  version?: number;
  pagination?: Pagination;
}

export interface Gap {
  capability: string;
  reason?: string;
  description?: string;
  blocking: boolean;
  severity: "blocking" | "high" | "medium" | "low" | "optional";
  source?: string;
}

export interface Outcome {
  title: string;
  description: string;
  confidence: number;
  ready: boolean;
  blocked_by?: string[];
}

export interface Assessment {
  task: string;
  capabilities?: CapabilityMatch[];
  gaps?: Gap[];
  gap_filling_tasks?: unknown[];
  blocking_gaps?: number;
  policies_recommended?: PolicyRecommendation[];
  needs_clarification?: string[];
  can_do: number;
  assessment_quality: "insufficient_data" | "partial" | "complete";
  ready_capabilities: number;
  total_capabilities: number;
  gap_count: number;
  ready_for_deployment: boolean;
  stack?: StackBreakdown;
  governance_decision?: GovernanceDecision;
  outcomes?: Outcome[];
  intelligence?: IntelligenceMeta;
}

export interface PolicyRecommendation {
  policy: string;
  reason: string;
  action: string;
  blocking: boolean;
}

export interface StackItem {
  name: string;
  type?: string;
  source?: StackItemSource;
  status?: StackItemStatus;
  description?: string;
  contains?: string[];
  functions?: string[];
}

export interface StackBreakdown {
  user_data?: StackItem[];
  cadreen?: StackItem[];
  connectors?: StackItem[];
  gaps?: StackItem[];
}

export interface CapabilityTrace {
  total_available: number;
  healthy_count: number;
  active_integrations?: string[];
}

export interface ReasoningTrace {
  capability_matches?: number;
}

export interface MemoryTrace {
  healthy: boolean;
  knowledge_queried?: number;
}

export interface GovernanceTrace {
  active: boolean;
  decision?: string;
  confidence?: number;
  reason_code?: string;
  policy_id?: string;
  next_actions?: NextAction[];
}

export interface HumilityTrace {
  gaps_detected?: number;
  blocking?: number;
}

export interface ProcessTrace {
  started_at: string;
  duration_ms: number;
  components?: Record<string, boolean>;
}

export interface NextAction {
  type: "connect_tool" | "add_policy" | "add_memory" | "resolve_gap" | "check_auth" | "none";
  label: string;
  endpoint?: string;
  reason: string;
}

export interface IntelligenceStage {
  name: "capability_check" | "memory_lookup" | "governance_eval" | "gap_analysis" | "healing_attempt" | "execution" | "escalation" | "human_handoff";
  status: "passed" | "failed" | "skipped" | "degraded" | "blocked";
  duration_ms?: number;
  detail?: string;
  inputs?: Record<string, unknown>;
  outputs?: Record<string, unknown>;
}

export interface FieldStability {
  stable: string[];
  evolving: string[];
  internal: string[];
}

export interface IntelligenceMeta {
  version?: string;
  summary?: string;
  capability: CapabilityTrace;
  reasoning: ReasoningTrace;
  memory: MemoryTrace;
  governance: GovernanceTrace;
  humility: HumilityTrace;
  process: ProcessTrace;
  next_action?: NextAction;
  stages?: IntelligenceStage[];
  field_stability: FieldStability;
}

export interface IntelligenceTraceEntry {
  id: string;
  domain: string;
  request_path: string;
  request_method: string;
  meta: IntelligenceMeta;
  created_at?: string;
}

export interface ListIntelligenceResponse {
  traces: IntelligenceTraceEntry[];
  count: number;
  pagination?: Pagination;
}

export interface IntelligenceStats {
  traces_24h: number;
  traces_7d: number;
  traces_30d: number;
  avg_confidence_by_domain: Record<string, number>;
  gap_detection_rate: number;
  governance_decisions: Record<string, number>;
}

export interface ReplayResult {
  trace_id: string;
  mode: string;
  domain: string;
  original_gate: string;
  original_confidence: number;
  current_gate: string;
  current_confidence: number;
  gate_changed: boolean;
  change_summary: string;
  current_capability: Record<string, unknown>;
  current_memory: Record<string, unknown>;
  current_gaps: Record<string, unknown>;
  replay_note: string;
}

export interface HandoffPacket {
  trace_id: string;
  domain: string;
  created_at: string;
  governance: Record<string, unknown>;
  what_the_system_knew: Record<string, unknown>;
  what_the_system_didnt_know: Record<string, unknown>;
  what_happened: Record<string, unknown>;
  suggested_actions: Record<string, unknown>[];
  next_action: Record<string, unknown>;
  trace_url: string;
}

export interface PromoteResult {
  id: string;
  kind: string;
  status: string;
  tool_name?: string;
  tool_sequence?: string[];
  source_trace_id: string;
}

export interface IntentMessage {
  role: "user" | "assistant" | "system";
  content: string;
}

export interface IntentContext {
  existing_connectors?: string[];
  constraints?: Record<string, unknown>;
  domain?: string;
}

export interface IntentRequest {
  messages: IntentMessage[];
  conversation_id?: string;
  context?: IntentContext;
  mode?: IntentMode;
  stream?: boolean;
}

export interface ResponseMessage {
  role: string;
  content: string;
  metadata?: Record<string, unknown>;
}

export interface ResponseExecution {
  id: string;
  status: string;
  stream_url?: string;
  poll_url?: string;
}

export interface ClarificationQuestion {
  id: string;
  question: string;
  type: string;
  required: boolean;
  reason?: string;
}

export type IntentResult =
  | { type: "direct"; message: ResponseMessage; intelligence: IntelligenceMeta; traceId: string }
  | { type: "clarify"; questions: ClarificationQuestion[]; conversationId: string; intelligence: IntelligenceMeta; traceId: string }
  | { type: "execution"; execution: ResponseExecution; intelligence: IntelligenceMeta; traceId: string }
  | { type: "blocked"; reason_code?: string; policy_id?: string; intelligence: IntelligenceMeta; traceId: string }
  | { type: "connect_required"; endpoint?: string; reason?: string; intelligence: IntelligenceMeta; traceId: string };

export interface IntentStatus {
  ready: boolean;
  needs: string[];
  next: string;
}

export function intentStatus(result: IntentResult): IntentStatus {
  switch (result.type) {
    case "direct":
      return {
        ready: true,
        needs: [],
        next: result.intelligence?.next_action?.type && result.intelligence.next_action.type !== "none"
          ? result.intelligence.next_action.label
          : "done",
      };
    case "execution":
      return {
        ready: true,
        needs: [],
        next: result.execution.stream_url
          ? "stream execution"
          : `poll ${result.execution.poll_url ?? result.execution.id}`,
      };
    case "clarify":
      return {
        ready: false,
        needs: result.questions.map(q => q.question),
        next: `answer ${result.questions.length} question${result.questions.length === 1 ? "" : "s"}`,
      };
    case "blocked":
      return {
        ready: false,
        needs: [`blocked: ${result.reason_code || result.intelligence?.governance?.reason_code || "governance gate"}`],
        next: result.policy_id ? `policy: ${result.policy_id}` : "resolve policy block",
      };
    case "connect_required":
      return {
        ready: false,
        needs: [`connect ${result.endpoint || "required tool"}`],
        next: result.endpoint || "",
      };
  }
}

export interface RegisterOpenAPIRequest {
  name: string;
  spec_url?: string;
  spec_content?: string;
  credential_id?: string;
}

export interface RegisterOpenAPIResponse {
  id: string;
  name: string;
  type: string;
  tools_generated?: number;
  tools_registered?: number;
  functions?: string[];
  spec_url?: string;
  status: string;
}

export interface RegisterMCPRequest {
  name: string;
  url: string;
  transport?: TransportType;
  headers?: Record<string, string>;
}

export interface RegisterMCPResponse {
  id: string;
  name: string;
  type: string;
  transport?: string;
  url?: string;
  status: string;
}

export interface SearchMemoryRequest {
  query: string;
  domain?: string;
  tag?: string;
  limit?: number;
}

export interface RememberRequest {
  type: AtomType;
  content: Record<string, unknown>;
  domain?: string;
  scope?: AtomScope;
  authority?: number;
  tags?: string[];
}

export interface CreatePolicyRequest {
  name: string;
  rules?: Record<string, unknown>[];
  domain?: string;
  auto_draft?: boolean;
}

export interface EvaluatePolicyRequest {
  action: string;
  domain?: string;
  context?: Record<string, unknown>;
}

export interface ExecutionEvent {
  type: string;
  data: Record<string, unknown>;
}

export interface ExecutionStatus {
  id: string;
  status: string;
  progress?: number;
  result?: Record<string, unknown>;
  error?: string;
}

export interface TraceExplain {
  summary: string;
  steps: string[];
  recommendations?: string[];
}

export interface CadreenConfig {
  apiKey: string;
  baseUrl?: string;
  maxRetries?: number;
  timeout?: number;
  telemetry?: unknown;
  sandbox?: boolean;
  fixtures?: Record<string, unknown>;
  /** Response profile: "lean" (no envelope), "audit" (action-bearing only), "full" (default) */
  profile?: "lean" | "audit" | "full";
}

export interface RequestOptions {
  idempotencyKey?: string;
  headers?: Record<string, string>;
}

export type ConnectResultType = "prebuilt" | "schema_required" | "manual" | "unknown";

export interface ConnectResult {
  type: ConnectResultType;
  capability: string;
  detail: ConnectPrebuiltDetail | ConnectSchemaRequiredDetail | ConnectManualDetail | ConnectUnknownDetail;
}

export interface ConnectPrebuiltDetail {
  tool_id: string;
  tool_name: string;
  service_id: string;
  service_name: string;
  auth_type: string;
  account_id?: string;
  source: string;
}

export interface ConnectSchemaRequiredDetail {
  tool_id: string;
  tool_name: string;
  auth_url: string;
  connector: string;
}

export interface ConnectManualDetail {
  pathways: ConnectPathway[];
}

export interface ConnectPathway {
  id: string;
  connector: string;
  tool_id: string;
  health: string;
  priority: number;
}

export interface ConnectUnknownDetail {
  searched: string;
  hints?: string[];
}

export interface SetupRequest {
  workspace_id?: string;
  purpose?: string;
  examples?: string[];
  constraints?: string[];
  connections?: Array<{ capability: string }>;
  credentials?: Array<{ provider: string; name?: string; key_data: Record<string, unknown> }>;
  memory?: Array<{ type?: string; content: Record<string, unknown>; domain?: string; tags?: string[]; authority?: number }>;
  policies?: Array<{ name: string; description?: string; rule: string; severity?: string }>;
  confirm?: boolean;
}

export interface SetupProposal {
  type: string;
  description: string;
  detail: string;
}

export interface SetupResult {
  workspace_id?: string;
  connections: Array<SetupConnectionResult>;
  credentials: Array<SetupCredentialResult>;
  memory: Array<SetupMemoryResult>;
  policies: Array<SetupPolicyResult>;
  proposals?: Array<SetupProposal>;
  applied: number;
  failed: number;
}

export interface SetupConnectionResult {
  capability: string;
  status: "connected" | "schema_required" | "manual" | "unknown" | "failed";
  detail?: ConnectPrebuiltDetail | ConnectSchemaRequiredDetail | ConnectManualDetail | ConnectUnknownDetail;
  error?: string;
}

export interface SetupCredentialResult {
  provider: string;
  name: string;
  id?: string;
  status: "created" | "failed";
  error?: string;
}

export interface SetupMemoryResult {
  id: string;
  type: string;
  kind?: string;
  classified: boolean;
  status: "created" | "failed";
  error?: string;
}

export interface SetupPolicyResult {
  name: string;
  id?: string;
  status: "created" | "failed";
  error?: string;
}

// ---------------------------------------------------------------------------
// Abstraction aliases: Atom → MemoryItem
// ---------------------------------------------------------------------------

/** @deprecated Use AtomCategory instead */
export type MemoryType = AtomCategory;

// ---------------------------------------------------------------------------
// Marketplace catalog types
// ---------------------------------------------------------------------------

export interface CatalogIntegration {
  id: string;
  name: string;
  description: string;
  category: string;
  capabilities?: string[];
  tags?: string[];
  provider: string;
  status: string;
  auth_type: string;
  install_time: string;
  popularity?: number;
  featured?: boolean;
}

export interface CatalogCategory {
  name: string;
  description: string;
  integrations: CatalogIntegration[];
}

export interface CatalogResponse {
  categories: CatalogCategory[];
  installed: string[];
  total_available: number;
}

export interface InstallResponse {
  status: string;
  auth_url?: string;
  provider: string;
  estimated_time?: string;
}
