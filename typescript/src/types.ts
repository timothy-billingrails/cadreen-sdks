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

export interface ConnectionGroup {
  capability: string;
  status: HealthStatus;
}

export interface ListConnectionsResponse {
  connections: ConnectionGroup[];
  total_capabilities: number;
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
  components?: Record<string, unknown>;
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
  dry_run?: boolean;
  user_id?: string;
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
  capability: string;
  available: boolean;
  health?: string;
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
  dry_run?: boolean;
}

export interface SetupProposal {
  type: string;
  description: string;
  detail: string;
}

export interface SetupResult {
  connections: Array<SetupConnectionResult>;
  credentials: Array<SetupCredentialResult>;
  memory: Array<SetupMemoryResult>;
  policies: Array<SetupPolicyResult>;
  proposals?: Array<SetupProposal>;
  notice?: string;
  applied: number;
  failed: number;
  dry_run?: boolean;
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
  status: "created" | "would_create" | "failed";
  error?: string;
}

export interface SetupMemoryResult {
  id: string;
  type: string;
  kind?: string;
  classified: boolean;
  status: "created" | "would_create" | "failed";
  error?: string;
}

export interface SetupPolicyResult {
  name: string;
  id?: string;
  status: "created" | "would_create" | "failed";
  error?: string;
}

export interface SetupSessionCreateRequest {
  workspace_id?: string;
  purpose?: string;
  constraints?: string[];
}

export interface SetupSessionAddRequest {
  connections?: Array<{ capability: string }>;
  credentials?: Array<{ provider: string; name?: string; key_data: Record<string, unknown> }>;
  memory?: Array<{ type?: string; content: Record<string, unknown>; domain?: string; tags?: string[]; authority?: number }>;
  policies?: Array<{ name: string; description?: string; rule: string; severity?: string }>;
}

export interface SetupSessionApplyRequest {
  confirm: boolean;
}

export interface SetupSession {
  id: string;
  status: "draft" | "applying" | "applied" | "failed";
  purpose?: string;
  constraints?: string[];
  connections: Array<{ capability: string }>;
  credentials: Array<{ provider: string; name?: string; key_data?: Record<string, unknown> }>;
  memory: Array<{ type?: string; content: Record<string, unknown> }>;
  policies: Array<{ name: string; rule?: string }>;
  proposals?: Array<SetupProposal>;
  applied_count: number;
  failed_count: number;
  created_at: string;
  updated_at: string;
  applied_at?: string;
}

export interface SetupSessionApplyResult {
  session_id: string;
  status: string;
  applied: number;
  failed: number;
  result?: SetupResult;
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

// ---------------------------------------------------------------------------
// Document types
// ---------------------------------------------------------------------------

export interface Document {
  id: string;
  name: string;
  content_type?: string;
  size?: number;
  status?: string;
  created_at?: string;
}

export interface ListDocumentsResponse {
  documents: Document[];
  count: number;
  pagination?: Pagination;
}

export interface UploadDocumentResponse {
  id: string;
  name: string;
  content_type?: string;
  size?: number;
  status?: string;
}

// ---------------------------------------------------------------------------
// Webhook types
// ---------------------------------------------------------------------------

export interface Webhook {
  id: string;
  url: string;
  events?: string[];
  secret?: string;
  is_active: boolean;
  created_at?: string;
}

export interface ListWebhooksResponse {
  webhooks: Webhook[];
  count: number;
  pagination?: Pagination;
}

// ---------------------------------------------------------------------------
// Learning types
// ---------------------------------------------------------------------------

export interface LearningPattern {
  id: string;
  pattern: string;
  confidence: number;
  occurrences?: number;
  domain?: string;
  tags?: string[];
  created_at?: string;
}

export interface LearningEpisode {
  id: string;
  description: string;
  outcome?: string;
  trace_id?: string;
  domain?: string;
  created_at?: string;
}

export interface LearningSuggestion {
  id: string;
  type: string;
  description: string;
  impact?: string;
  domain?: string;
}

export interface ListLearningPatternsResponse {
  patterns: LearningPattern[];
  count: number;
  pagination?: Pagination;
}

export interface ListLearningEpisodesResponse {
  episodes: LearningEpisode[];
  count: number;
  pagination?: Pagination;
}

export interface ListLearningSuggestionsResponse {
  suggestions: LearningSuggestion[];
  count: number;
  pagination?: Pagination;
}

// ---------------------------------------------------------------------------
// Healing types
// ---------------------------------------------------------------------------

export interface HealingDiagnosis {
  error_category?: string;
  semantic_reason?: string;
  root_cause?: string;
  can_retry?: boolean;
  needs_sub_execution?: boolean;
  needs_human?: boolean;
  should_skip?: boolean;
  needs_re_decide?: boolean;
  needs_try_alternative?: boolean;
  retry_delay_ms?: number;
  confidence?: number;
}

export interface StrategyCount {
  strategy: string;
  count: number;
}

export interface ToolHealingStats {
  tool_name: string;
  total: number;
  successful: number;
  failed: number;
  success_rate: number;
  top_strategy?: string;
}

export interface TimeRange {
  first_precedent?: string;
  last_precedent?: string;
}

export interface HealingStatsResponse {
  total_precedents?: number;
  successful_recoveries?: number;
  failed_recoveries?: number;
  success_rate?: number;
  avg_duration_ms?: number;
  common_strategies?: StrategyCount[];
  top_tools?: ToolHealingStats[];
  by_category?: Record<string, unknown>;
  time_range?: TimeRange;
}

export interface HealingPrecedent {
  id: string;
  tool_name?: string;
  error_type: string;
  error_category?: string;
  semantic_reason?: string;
  root_cause?: string;
  recovery_strategy?: string;
  success: boolean;
  what_worked?: string;
  what_failed?: string;
  attempts: number;
  duration_ms?: number;
  confidence: number;
  created_at?: string;
  domain?: string;
  tags?: string[];
}

export interface ListHealingPrecedentsResponse {
  precedents: HealingPrecedent[];
  count: number;
  pagination?: Pagination;
}

// ---------------------------------------------------------------------------
// Credential types (request)
// ---------------------------------------------------------------------------

export interface CreateCredentialRequest {
  provider: string;
  name?: string;
  key_data: Record<string, unknown>;
}

export interface ResolveEscalationRequest {
  decision: string;
}

export interface DiagnoseRequest {
  error_message: string;
  tool_name?: string;
  trace_id?: string;
}

// ---------------------------------------------------------------------------
// Proposal types
// ---------------------------------------------------------------------------

export type ProposalType = "automation" | "governance" | "learning" | "blueprint" | "task";
export type ProposalStatus = "proposed" | "accepted" | "dismissed" | "expired" | "superseded";

export interface ProposalEvidence {
  description: string;
  source?: string;
  count?: number;
  confidence?: number;
}

export interface TaskProposal {
  id: string;
  title: string;
  description: string;
  intent: string;
  domain?: string;
  proposal_type: ProposalType;
  mission_intent?: string;
  trigger_type: string;
  trigger_source: string;
  trigger_details?: string;
  evidence?: ProposalEvidence[];
  confidence: number;
  priority: number;
  status: ProposalStatus;
  created_at: string;
  expires_at?: string;
  accepted_at?: string;
  dismissed_at?: string;
  dismissal_reason?: string;
  execution_id?: string;
  dedup_key?: string;
  requires_review?: boolean;
}

export interface ListProposalsOptions {
  status?: ProposalStatus | "all";
  limit?: number;
}

export interface ListProposalsResponse {
  proposals: TaskProposal[];
  count: number;
}

export interface AcceptProposalResponse {
  status: string;
  execution_id: string;
  action: string;
  intent: string;
  next_step: string;
  auto_approved?: boolean;
  result?: Record<string, unknown>;
}

export interface DismissProposalResponse {
  status: string;
}

export interface ProposalStatsResponse {
  proposed: number;
  accepted: number;
  dismissed: number;
  expired: number;
}

export type WorkspaceRole = "admin" | "operator" | "member" | "viewer";

export interface WorkspaceUser {
  id: string;
  workspace_id: string;
  user_id: string;
  role: WorkspaceRole;
  invited_by?: string;
  invited_at: string;
  created_at: string;
  updated_at: string;
}

export interface InviteUserRequest {
  email: string;
  role?: WorkspaceRole;
}

export interface UpdateRoleRequest {
  role: WorkspaceRole;
}

export interface ListWorkspaceUsersResponse {
  users: WorkspaceUser[];
  count: number;
}

// ---------------------------------------------------------------------------
// Agent types
// ---------------------------------------------------------------------------

export type AgentStatus = "draft" | "active" | "deploying" | "deployed" | "paused" | "error" | "archived";
export type AgentHealth = "healthy" | "degraded" | "unhealthy" | "unknown";
export type AgentMessageType = "direct" | "broadcast" | "request" | "response" | "notification";
export type AgentMessageStatus = "pending" | "delivered" | "read" | "failed" | "expired";
export type NegotiationStatus = "pending" | "in_progress" | "accepted" | "rejected" | "counter_proposed" | "expired" | "cancelled";
export type FactType = "fact" | "preference" | "constraint" | "instruction" | "context" | "relationship";
export type GovernanceScope = "agent" | "workspace" | "federation";
export type PolicyAction = "allow" | "deny" | "require_approval" | "log_only";

export interface Agent {
  id: string;
  name: string;
  description?: string;
  status: AgentStatus;
  health: AgentHealth;
  model?: string;
  system_prompt?: string;
  capabilities?: string[];
  tags?: string[];
  config?: Record<string, unknown>;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
  deployed_at?: string;
}

export interface AgentConfig {
  agent_id: string;
  model?: string;
  system_prompt?: string;
  temperature?: number;
  max_tokens?: number;
  tools?: string[];
  connections?: string[];
  memory_scope?: string;
  governance_policy_id?: string;
  metadata?: Record<string, unknown>;
}

export interface AgentCapabilities {
  agent_id: string;
  tools: string[];
  connections: string[];
  knowledge_count: number;
  governance_policies: number;
  can_execute: boolean;
  can_federate: boolean;
  can_negotiate: boolean;
}

export interface AgentKnowledge {
  id: string;
  agent_id: string;
  fact_type: FactType;
  subject: string;
  predicate: string;
  object: string;
  source?: string;
  confidence: number;
  tags?: string[];
  visibility: "private" | "workspace" | "federation";
  created_at: string;
  updated_at: string;
}

export interface AgentGovernancePolicy {
  id: string;
  name: string;
  description?: string;
  scope: GovernanceScope;
  agent_id?: string;
  rules: Record<string, unknown>[];
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface AgentAuditEntry {
  id: string;
  agent_id: string;
  action: string;
  target_type?: string;
  target_id?: string;
  details?: Record<string, unknown>;
  policy_action?: PolicyAction;
  created_at: string;
}

export interface AgentNegotiation {
  id: string;
  from_agent_id: string;
  to_agent_id: string;
  proposal: Record<string, unknown>;
  status: NegotiationStatus;
  current_round: number;
  max_rounds: number;
  deadline?: string;
  resolution?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface AgentMessage {
  id: string;
  from_agent_id: string;
  to_agent_id: string;
  content: string;
  context?: Record<string, unknown>;
  status: AgentMessageStatus;
  response?: string;
  message_type: AgentMessageType;
  created_at: string;
}

export interface AgentExecution {
  id: string;
  agent_id: string;
  status: string;
  input?: Record<string, unknown>;
  output?: Record<string, unknown>;
  error?: string;
  started_at: string;
  completed_at?: string;
}

export interface ListAgentsResponse {
  agents: Agent[];
  count: number;
  pagination?: Pagination;
}

export interface ListAgentMessagesResponse {
  messages: AgentMessage[];
  count: number;
  pagination?: Pagination;
}

export interface ListAgentExecutionsResponse {
  executions: AgentExecution[];
  count: number;
  pagination?: Pagination;
}

export interface ListAgentKnowledgeResponse {
  knowledge: AgentKnowledge[];
  count: number;
  pagination?: Pagination;
}

export interface ListAgentGovernanceResponse {
  policies: AgentGovernancePolicy[];
  count: number;
}

export interface ListAgentAuditResponse {
  entries: AgentAuditEntry[];
  count: number;
  pagination?: Pagination;
}

export interface ListAgentNegotiationsResponse {
  negotiations: AgentNegotiation[];
  count: number;
  pagination?: Pagination;
}

export interface CreateAgentRequest {
  name: string;
  description?: string;
  model?: string;
  system_prompt?: string;
  capabilities?: string[];
  tags?: string[];
  config?: Record<string, unknown>;
  metadata?: Record<string, unknown>;
}

export interface UpdateAgentRequest {
  name?: string;
  description?: string;
  model?: string;
  system_prompt?: string;
  capabilities?: string[];
  tags?: string[];
  config?: Record<string, unknown>;
  metadata?: Record<string, unknown>;
}

export interface SendMessageRequest {
  content: string;
  message_type?: AgentMessageType;
  context?: Record<string, unknown>;
}

export interface CreateExecutionRequest {
  input: Record<string, unknown>;
  stream?: boolean;
}

export interface CreateKnowledgeRequest {
  fact_type: FactType;
  subject: string;
  predicate: string;
  object: string;
  source?: string;
  confidence?: number;
  tags?: string[];
  visibility?: "private" | "workspace" | "federation";
}

export interface SearchKnowledgeRequest {
  query: string;
  fact_type?: FactType;
  limit?: number;
}

export interface CreateGovernanceRequest {
  name: string;
  description?: string;
  scope: GovernanceScope;
  rules: Record<string, unknown>[];
  enabled?: boolean;
}

export interface UpdateGovernanceRequest {
  name?: string;
  description?: string;
  rules?: Record<string, unknown>[];
  enabled?: boolean;
}

export interface StartNegotiationRequest {
  to_agent_id: string;
  proposal: Record<string, unknown>;
  max_rounds?: number;
  deadline?: string;
}

export interface RespondNegotiationRequest {
  action: "accept" | "reject" | "counter";
  counter_proposal?: Record<string, unknown>;
  reason?: string;
}

// ---------------------------------------------------------------------------
// Federation types
// ---------------------------------------------------------------------------

export type FederationStatus = "pending" | "active" | "suspended" | "revoked" | "expired";
export type FederationAgentStatus = "active" | "paused" | "error" | "unlinked";

export interface FederationLink {
  id: string;
  targetWorkspaceId: string;
  status: FederationStatus;
  permissions: Record<string, unknown>;
  name?: string;
  description?: string;
  metadata?: Record<string, unknown>;
  createdByUserId?: string;
  approvedByUserId?: string;
  suspendedByUserId?: string;
  suspensionReason?: string;
  lastActivityAt?: string;
  revokedAt?: string;
  revokeReason?: string;
  sourceWorkspaceName?: string;
  sourceWorkspaceSlug?: string;
  targetWorkspaceName?: string;
  createdAt: string;
  updatedAt: string;
}

export interface FederationAgent {
  id: string;
  federationLinkId: string;
  localAgentId: string;
  remoteAgentId: string;
  localAgentName?: string;
  remoteAgentName?: string;
  status: FederationAgentStatus;
  capabilities?: string[];
  createdAt: string;
  updatedAt: string;
}

export interface FederationPermissions {
  federation_link_id: string;
  permissions: string[];
  allowed_agents?: string[];
  allowed_tools?: string[];
  allowed_knowledge?: string[];
  max_requests_per_day?: number;
  updated_at: string;
}

export interface ListFederationResponse {
  links: FederationLink[];
  count: number;
}

export interface ListFederationAgentsResponse {
  agents: FederationAgent[];
  count: number;
}

export interface CreateFederationRequest {
  target_workspace_id: string;
  name?: string;
  description?: string;
  permissions?: string[];
  metadata?: Record<string, unknown>;
}

export interface SuspendFederationRequest {
  reason?: string;
}

export interface RevokeFederationRequest {
  reason?: string;
}

export interface UpdatePermissionsRequest {
  permissions: string[];
  allowed_agents?: string[];
  allowed_tools?: string[];
  allowed_knowledge?: string[];
  max_requests_per_day?: number;
}

export interface LinkAgentRequest {
  local_agent_id: string;
  remote_agent_id: string;
  capabilities?: string[];
}

// ---------------------------------------------------------------------------
// Responses API types
// ---------------------------------------------------------------------------

export interface ResponseRequest {
  model: string;
  input: string | ResponseInputItem[];
  instructions?: string;
  tools?: ResponseTool[];
  previous_response_id?: string;
  store?: boolean;
  stream?: boolean;
  max_output_tokens?: number;
  temperature?: number;
  metadata?: Record<string, unknown>;
}

export interface ResponseInputItem {
  type: 'message' | 'function_call_output';
  role?: 'user' | 'assistant' | 'system';
  content?: string | ResponseContentPart[];
  call_id?: string;
  output?: string;
}

export interface ResponseContentPart {
  type: 'input_text' | 'output_text' | 'image_url';
  text?: string;
  image_url?: string;
}

export interface ResponseTool {
  type: 'function';
  name: string;
  description?: string;
  parameters?: Record<string, unknown>;
  strict?: boolean;
}

export interface ResponsesCompletion {
  id: string;
  object: 'response';
  created_at: number;
  model: string;
  output: ResponseOutputItem[];
  output_text?: string;
  usage?: ResponseUsage;
  status: 'completed' | 'failed' | 'in_progress' | 'cancelled';
  previous_response_id?: string;
  metadata?: Record<string, unknown>;
}

export interface ResponseOutputItem {
  id: string;
  type: 'message' | 'reasoning' | 'function_call';
  status?: 'completed' | 'in_progress' | 'incomplete';
  role?: 'assistant';
  content?: ResponseOutputContent[];
  name?: string;
  call_id?: string;
  arguments?: string;
}

export interface ResponseOutputContent {
  type: 'output_text' | 'refusal';
  text?: string;
  annotations?: ResponseAnnotation[];
}

export interface ResponseAnnotation {
  type: 'url_citation' | 'file_citation';
  url?: string;
}

export interface ResponseUsage {
  input_tokens: number;
  output_tokens: number;
  total_tokens: number;
}

export interface ResponseStreamEvent {
  type: string;
  sequence?: number;
  response?: ResponsesCompletion;
  item?: ResponseOutputItem;
  output_index?: number;
  content_index?: number;
  delta?: string;
}

// ─── External Agents (A2A) ───

export type ExternalAgentStatus = "pending_approval" | "active" | "suspended" | "revoked" | "error";
export type ExternalAgentHealth = "healthy" | "degraded" | "unhealthy" | "unknown";
export type ExternalAgentDirection = "outbound" | "inbound";
export type ExternalAgentOperation = "message/send" | "message/stream" | "tasks/get" | "tasks/list" | "tasks/cancel";
export type ExternalAgentTaskStatus = "submitted" | "working" | "input-required" | "completed" | "failed" | "canceled" | "rejected";
export type ExternalAgentGovernanceResult = "allowed" | "blocked" | "requires_approval" | "no_policies";

export interface ExternalAgentSkill {
  id: string;
  name: string;
  description?: string;
  tags?: string[];
}

export interface ExternalAgentCapabilities {
  streaming?: boolean;
  pushNotifications?: boolean;
  stateTransitionHistory?: boolean;
}

export interface ExternalAgentConnection {
  id: string;
  agentId: string;
  agentCardUrl: string;
  agentName: string;
  agentDescription: string;
  agentSystem: string;
  agentVersion: string;
  agentCardJson: Record<string, unknown>;
  skills: ExternalAgentSkill[];
  capabilities: ExternalAgentCapabilities;
  status: ExternalAgentStatus;
  health: ExternalAgentHealth;
  lastUsedAt?: string;
  lastHealthCheckAt?: string;
  errorMessage?: string;
  approvedBy?: string;
  approvedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ExternalAgentInteraction {
  id: string;
  connectionId: string;
  agentId: string;
  direction: ExternalAgentDirection;
  operation: ExternalAgentOperation;
  taskId?: string;
  message?: string;
  status: ExternalAgentTaskStatus;
  durationMs?: number;
  errorMessage?: string;
  governanceResult?: ExternalAgentGovernanceResult;
  createdAt: string;
}

export interface ExternalAgentSettings {
  enabled: boolean;
}

export interface ListExternalConnectionsResponse {
  connections: ExternalAgentConnection[];
  total: number;
  limit: number;
  offset: number;
}

export interface ListExternalInteractionsResponse {
  interactions: ExternalAgentInteraction[];
  total: number;
  limit: number;
  offset: number;
}
