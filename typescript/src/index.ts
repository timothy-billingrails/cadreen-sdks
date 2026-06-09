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
} from "./types";

export { CadreenError } from "./client";
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
  InstallComposioRequest,
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
  ConnectPathway,
  ConnectUnknownDetail,
  SetupRequest,
  SetupResult,
  SetupConnectionResult,
  SetupCredentialResult,
  SetupMemoryResult,
  SetupPolicyResult,
} from "./types";

export { GuardrailsResource } from "./resources/guardrails";
