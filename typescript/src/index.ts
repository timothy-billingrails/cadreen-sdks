import { HttpClient } from "./client";
import { IntentResource } from "./resources/intent";
import { MemoryResource } from "./resources/memory";
import { PoliciesResource } from "./resources/policies";
import { ConnectionsResource } from "./resources/connections";
import { TracesResource } from "./resources/traces";
import { ExecutionsResource } from "./resources/executions";
import { GuardrailsResource } from "./resources/guardrails";
import { SkillsResource } from "./resources/skills";
import { FailuresResource } from "./resources/failures";
import { WebhooksResource } from "./resources/webhooks";
import { ChatResource } from "./resources/chat";
import { DocumentsResource } from "./resources/documents";
import { EscalationsResource } from "./resources/escalations";
import { HealingResource } from "./resources/healing";
import { LearningResource } from "./resources/learning";
import { CredentialsResource } from "./resources/credentials";
import { BlueprintsResource } from "./resources/blueprints";
import { SchedulesResource } from "./resources/schedules";
import { ProposalsResource } from "./resources/proposals";
import { SetupSessionsResource } from "./resources/setup-sessions";
import { WorkspaceUsersResource } from "./resources/workspace-users";
import { AgentsResource } from "./resources/agents";
import { FederationResource } from "./resources/federation";
import { ResponsesResource } from "./resources/responses";
import { ExternalAgentsResource } from "./resources/external-agents";
import { DevicesResource } from "./resources/devices";
import type {
  CadreenConfig,
  IntentRequest,
  IntentResult,
  IntentContext,
  RememberRequest,
  CreateMemoryResponse,
  SearchMemoryRequest,
  SearchMemoryResponse,
  ConnectResult,
  SetupRequest,
  SetupResult,
  ListCapabilitiesResponse,
  Assessment,
} from "./types";

export { CadreenError, CadreenBlockedError, CadreenClarifyError } from "./client";
export { requiresHuman, handoffReason, explainTrace, redactTrace, redactMessages, redactResult } from "./intelligence_helpers";
export { intentStatus } from "./types";
export type { RedactOptions } from "./intelligence_helpers";
export type { TelemetryProvider, TelemetrySpan, TelemetryMeter } from "./telemetry";
export { OpenTelemetryAdapter, NoOpProvider, NoOpSpan, NoOpMeter } from "./telemetry";

export class Cadreen {
  public readonly intent: IntentResource;
  public readonly memory: MemoryResource;
  public readonly policies: PoliciesResource;
  public readonly connections: ConnectionsResource;
  public readonly traces: TracesResource;
  public readonly executions: ExecutionsResource;
  public readonly guardrails: GuardrailsResource;
  public readonly skills: SkillsResource;
  public readonly failures: FailuresResource;
  public readonly webhooks: WebhooksResource;
  public readonly chat: ChatResource;
  public readonly documents: DocumentsResource;
  public readonly escalations: EscalationsResource;
  public readonly healing: HealingResource;
  public readonly learning: LearningResource;
  public readonly credentials: CredentialsResource;
  public readonly blueprints: BlueprintsResource;
  public readonly schedules: SchedulesResource;
  public readonly proposals: ProposalsResource;
  public readonly setupSessions: SetupSessionsResource;
  public readonly workspaceUsers: WorkspaceUsersResource;
  public readonly agents: AgentsResource;
  public readonly federation: FederationResource;
  public readonly responses: ResponsesResource;
  public readonly externalAgents: ExternalAgentsResource;
  public readonly devices: DevicesResource;

  private readonly client: HttpClient;

  constructor(config: CadreenConfig) {
    this.client = new HttpClient(config);
    this.intent = new IntentResource(this.client);
    this.memory = new MemoryResource(this.client);
    this.policies = new PoliciesResource(this.client);
    this.connections = new ConnectionsResource(this.client);
    this.traces = new TracesResource(this.client);
    this.executions = new ExecutionsResource(this.client);
    this.guardrails = new GuardrailsResource(this.policies);
    this.skills = new SkillsResource(this.intent, this.memory, this.connections);
    this.failures = new FailuresResource(this.traces);
    this.webhooks = new WebhooksResource(this.client);
    this.chat = new ChatResource(this.client);
    this.documents = new DocumentsResource(this.client);
    this.escalations = new EscalationsResource(this.client);
    this.healing = new HealingResource(this.client);
    this.learning = new LearningResource(this.client);
    this.credentials = new CredentialsResource(this.client);
    this.blueprints = new BlueprintsResource(this.client);
    this.schedules = new SchedulesResource(this.client);
    this.proposals = new ProposalsResource(this.client);
    this.setupSessions = new SetupSessionsResource(this.client);
    this.workspaceUsers = new WorkspaceUsersResource(this.client);
    this.agents = new AgentsResource(this.client);
    this.federation = new FederationResource(this.client);
    this.responses = new ResponsesResource(this.client);
    this.externalAgents = new ExternalAgentsResource(this.client);
    this.devices = new DevicesResource(this.client);
  }

  async invoke(request: IntentRequest): Promise<IntentResult> {
    return this.intent.invoke(request);
  }

  async ask(
    prompt: string,
    options?: { conversation_id?: string; context?: IntentContext; stream?: boolean }
  ): Promise<IntentResult> {
    return this.intent.invoke({
      messages: [{ role: "user", content: prompt }],
      mode: "chat",
      conversation_id: options?.conversation_id,
      context: options?.context,
      stream: options?.stream,
    });
  }

  async act(
    prompt: string,
    options?: { conversation_id?: string; context?: IntentContext; stream?: boolean }
  ): Promise<IntentResult> {
    return this.intent.invoke({
      messages: [{ role: "user", content: prompt }],
      mode: "execution",
      conversation_id: options?.conversation_id,
      context: options?.context,
      stream: options?.stream,
    });
  }

  async remember(request: RememberRequest): Promise<CreateMemoryResponse> {
    return this.memory.remember(request);
  }

  async context(request: SearchMemoryRequest): Promise<SearchMemoryResponse> {
    return this.memory.search(request);
  }

  async connect(capability: string): Promise<ConnectResult> {
    return this.connections.connect(capability);
  }

  async setup(request: SetupRequest): Promise<SetupResult> {
    return this.client.post<SetupResult>("/api/v1/cadreen/setup", request);
  }

  async listCapabilities(): Promise<ListCapabilitiesResponse> {
    return this.client.get<ListCapabilitiesResponse>("/api/v1/cadreen/capabilities");
  }

  async assess(task: string, domain?: string): Promise<Assessment> {
    const body: Record<string, string> = { task };
    if (domain) body.domain = domain;
    const wrapper = await this.client.post<{ assessment: Assessment }>("/api/v1/cadreen/assess", body);
    return wrapper.assessment;
  }
}

export type {
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
  IntentMessage,
  IntentContext,
  IntentRequest,
  ResponseMessage,
  ResponseExecution,
  IntentResult,
  IntentStatus,
  RegisterOpenAPIRequest,
  RegisterOpenAPIResponse,
  RegisterMCPRequest,
  RegisterMCPResponse,
  CatalogResponse,
  CatalogCategory,
  CatalogIntegration,
  InstallResponse,
  SearchMemoryRequest,
  RememberRequest,
  CreatePolicyRequest,
  EvaluatePolicyRequest,
  ExecutionEvent,
  ExecutionStatus,
  TraceExplain,
  RequestOptions,
  NextAction,
  IntelligenceStage,
  FieldStability,
  ClarificationQuestion,
  ConnectResult,
  ConnectResultType,
  ConnectPrebuiltDetail,
  ConnectSchemaRequiredDetail,
  ConnectManualDetail,
  ConnectUnknownDetail,
  SetupRequest,
  SetupResult,
  SetupConnectionResult,
  SetupCredentialResult,
  SetupMemoryResult,
  SetupPolicyResult,
  SetupProposal,
  SetupSessionCreateRequest,
  SetupSessionAddRequest,
  SetupSessionApplyRequest,
  SetupSession,
  SetupSessionApplyResult,
  ReplayResult,
  HandoffPacket,
  PromoteResult,
  Document,
  ListDocumentsResponse,
  UploadDocumentResponse,
  Webhook,
  ListWebhooksResponse,
  LearningPattern,
  LearningEpisode,
  LearningSuggestion,
  ListLearningPatternsResponse,
  ListLearningEpisodesResponse,
  ListLearningSuggestionsResponse,
  HealingDiagnosis,
  HealingStatsResponse,
  HealingPrecedent,
  ListHealingPrecedentsResponse,
  StrategyCount,
  ToolHealingStats,
  TimeRange,
  CreateCredentialRequest,
  ResolveEscalationRequest,
  DiagnoseRequest,
  DeviceDiagnoseRequest,
  TaskProposal,
  ProposalType,
  ProposalStatus,
  ProposalEvidence,
  ListProposalsOptions,
  ListProposalsResponse,
  AcceptProposalResponse,
  DismissProposalResponse,
  ProposalStatsResponse,
  WorkspaceRole,
  WorkspaceUser,
  InviteUserRequest,
  UpdateRoleRequest,
  ListWorkspaceUsersResponse,
  AgentStatus,
  AgentHealth,
  AgentMessageType,
  AgentMessageStatus,
  NegotiationStatus,
  FactType,
  GovernanceScope,
  PolicyAction,
  Agent,
  AgentConfig,
  AgentCapabilities,
  AgentKnowledge,
  AgentGovernancePolicy,
  AgentAuditEntry,
  AgentNegotiation,
  AgentMessage,
  AgentExecution,
  ListAgentsResponse,
  ListAgentMessagesResponse,
  ListAgentExecutionsResponse,
  ListAgentKnowledgeResponse,
  ListAgentGovernanceResponse,
  ListAgentAuditResponse,
  ListAgentNegotiationsResponse,
  CreateAgentRequest,
  UpdateAgentRequest,
  SendMessageRequest,
  CreateExecutionRequest,
  CreateKnowledgeRequest,
  SearchKnowledgeRequest,
  CreateGovernanceRequest,
  UpdateGovernanceRequest,
  StartNegotiationRequest,
  RespondNegotiationRequest,
  FederationStatus,
  FederationAgentStatus,
  FederationLink,
  FederationAgent,
  FederationPermissions,
  ListFederationResponse,
  ListFederationAgentsResponse,
  CreateFederationRequest,
  SuspendFederationRequest,
  RevokeFederationRequest,
  UpdatePermissionsRequest,
  LinkAgentRequest,
  ResponseRequest,
  ResponseInputItem,
  ResponseContentPart,
  ResponseTool,
  ResponsesCompletion,
  ResponseOutputItem,
  ResponseOutputContent,
  ResponseAnnotation,
  ResponseUsage,
  ResponseStreamEvent,
  ExternalAgentStatus,
  ExternalAgentHealth,
  ExternalAgentDirection,
  ExternalAgentOperation,
  ExternalAgentTaskStatus,
  ExternalAgentGovernanceResult,
  ExternalAgentSkill,
  ExternalAgentCapabilities,
  ExternalAgentConnection,
  ExternalAgentInteraction,
  ExternalAgentSettings,
  ListExternalConnectionsResponse,
  ListExternalInteractionsResponse,
} from "./types";

export { ProposalsResource } from "./resources/proposals";
export { SetupSessionsResource } from "./resources/setup-sessions";
export { WorkspaceUsersResource } from "./resources/workspace-users";
export { GuardrailsResource } from "./resources/guardrails";
export { WebhooksResource } from "./resources/webhooks";
export { ChatResource } from "./resources/chat";
export { DocumentsResource } from "./resources/documents";
export { EscalationsResource } from "./resources/escalations";
export { HealingResource } from "./resources/healing";
export { LearningResource } from "./resources/learning";
export { CredentialsResource } from "./resources/credentials";
export { AgentsResource } from "./resources/agents";
export { FederationResource } from "./resources/federation";
export { ResponsesResource } from "./resources/responses";
export { ExternalAgentsResource } from "./resources/external-agents";
export type {
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
  ChatStreamEvent,
  ToolEntry,
  ListToolsResponse,
} from "./resources/chat";
