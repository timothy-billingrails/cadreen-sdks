# Changelog

All notable changes to the Go SDK will be documented in this file.

## v0.7.2
- Added `devices.go` — full device lifecycle (CreateDevice, ListDevices, GetDevice, DeleteDevice, GetDeviceStatus, UpdateDeviceState, GetDeviceMap, etc.)
- `GetAgentConfig` now returns `*AgentConfig` instead of `map[string]any`
- `CreateAgentExecution` now returns `*AgentExecution` instead of `map[string]any`
- `GetFederationPermissions` and `UpdateFederationPermissions` now return `*FederationPermissions` instead of `map[string]any`
- Added typed structs: `AgentConfig`, `AgentExecution`, `FederationPermissions`

## v0.7.1
- Fix: sandbox mode now returns `ErrSandboxStreaming` for `PostStream` and `Stream` instead of making network requests
- Fix: release workflow Homebrew formula update ordering
- Documented new server-side execution events: `execution_steps_complete`, `mission_completed_with_gaps`

## v0.7.0
- **BREAKING:** Removed `Pathways` and `TotalPathways` from connection responses. `ConnectionGroup` now returns only `Capability` and `Status`.
- **BREAKING:** Removed `Pathway` type. Internal routing details (connector, transport, tool_id) are no longer exposed.
- **BREAKING:** Changed `ConnectManualDetail` from `{Pathways: [...]}` to `{Capability, Available, Health}`.
- **BREAKING:** Removed `WorkspaceID` from response types: `SetupResult`, `SetupSession`, `WebhookSubscription`, `WebhookPayload`. (Still accepted on request types.)
- **BREAKING:** Removed `AuthScheme` from `ExternalAgentConnection` responses.
- **BREAKING:** Removed `AtomsConsulted`, `EpisodesMatched`, `PrecedentsApplied` from memory trace in intelligence metadata.
- **BREAKING:** `SourcesConsulted` renamed to `KnowledgeQueried` in MemoryTrace.
- **BREAKING:** All entity responses (Agent, Knowledge, Governance, Federation, Negotiation, ExternalAgentConnection) no longer include `WorkspaceID`.
- New MCP SSE endpoint: `GET /api/v1/cadreen/mcp/sse` + `POST /api/v1/cadreen/mcp/message`. Connect Cadreen as an MCP server without installing the npm package.
- Added `agents.go` — full agent lifecycle (CreateAgent, ListAgents, GetAgent, UpdateAgent, DeleteAgent, DeployAgent, GetAgentConfig, GetAgentCapabilities)
- Added agent messaging (SendAgentMessage, ListAgentMessages)
- Added agent executions (CreateAgentExecution, ListAgentExecutions)
- Added agent knowledge (CreateAgentKnowledge, ListAgentKnowledge, SearchAgentKnowledge, DeleteAgentKnowledge)
- Added agent governance (CreateAgentGovernance, ListAgentGovernance, UpdateAgentGovernance, DeleteAgentGovernance)
- Added agent audit (ListAgentAudit)
- Added negotiations (StartNegotiation, ListNegotiations, GetNegotiation, RespondToNegotiation)
- Added `federation.go` — cross-workspace federation (CreateFederation, ListFederations, GetFederation, ApproveFederation, SuspendFederation, RevokeFederation)
- Added federation permissions (GetFederationPermissions, UpdateFederationPermissions)
- Added federation agent linking (LinkFederationAgent, ListFederationAgents, UnlinkFederationAgent)
- Added `external_agents.go` — A2A external agent connections (ConnectExternalAgent, ListExternalConnections, GetExternalConnection, ApproveExternalConnection, SuspendExternalConnection, RevokeExternalConnection, DeleteExternalConnection)
- Added external agent interactions (ListExternalInteractions)
- Added external agent settings (GetExternalAgentSettings, UpdateExternalAgentSettings, ListAllExternalConnections)
- Added `responses.go` — OpenAI-compatible responses API (CreateResponse, GetResponse, CreateResponseStreaming)
- Added 31 new types: Agent, AgentKnowledge, AgentGovernancePolicy, AgentAuditEntry, AgentNegotiation, FederationLink, FederationAgent, ExternalAgentConnection, ExternalAgentInteraction, ExternalAgentSettings, ExternalAgentSkill, ExternalAgentCapabilities, etc.

## v0.6.3
- Fix: retry body reuse bug — retries were sending empty body because `bytes.Reader` was consumed on first attempt
- Add: `VerifyWebhookSignature()` function for HMAC-SHA256 webhook signature verification

## v0.6.2
- Added `Intelligence` field to `ChatCompletionResponse` — full intelligence metadata
- Added `ConversationID` field to `ChatCompletionResponse` — conversation continuity
- Added `Reasoning` field to `ChatDelta` — streaming reasoning from thinking models
- Added `ReasoningTokens`, `CacheWriteTokens`, `PromptTokensDetails` to `ChatUsage`
- Added `Reasoning` field to `ChatStreamEvent` — non-empty on reasoning_delta events
- Document `reasoning_delta` streaming event in README
- IntelligenceMeta shape aligned with server

## v0.6.1
- Added `UserID` field to `IntentRequest` and `ChatCompletionRequest` — pass end-user identity for per-user context and memory filtering
- Added `ListWorkspaceUsers()`, `InviteUser()`, `UpdateUserRole()`, `RemoveUser()` — manage workspace team members
- Added `WorkspaceUser`, `WorkspaceRole`, `InviteUserRequest`, `UpdateRoleRequest` types

## v0.6.0 (BREAKING)
- **BREAKING:** `ResolveEscalationRequest.Resolution` renamed to `Decision` — aligns SDK with API contract
- **BREAKING:** `DiagnoseRequest.Error` renamed to `ErrorMessage` — aligns SDK with API contract
- Retry default was already at 3 — no change needed

## v0.5.5
- Added `CreateSetupSession()`, `ListSetupSessions()`, `GetSetupSession()`, `AddToSetupSession()`, `ApplySetupSession()` — stateful setup sessions
- Added `SetupSession`, `SetupSessionCreateRequest`, `SetupSessionAddRequest`, `SetupSessionApplyRequest`, `SetupSessionApplyResult` types

## v0.5.4
- Added `ListProposals()`, `GetProposal()`, `AcceptProposal()`, `DismissProposal()`, `ProposalStats()` — task proposals
- Added `TaskProposal`, `ProposalEvidence`, `ListProposalsResponse`, `AcceptProposalResponse`, `DismissProposalResponse`, `ProposalStatsResponse` types

## v0.5.3
- Added `DryRun` field to `IntentRequest` — preview intent classification, governance, and capability assessment without creating a mission or persisting conversation
- Fixed `Version` constant (was stuck at 0.5.0)

## v0.5.2
- Added `DryRun` field to `SetupRequest` — preview what would be created without persisting
- Added `DryRun` and `Notice` fields to `SetupResult`
- Added `ListBlueprints()`, `GetBlueprint()`, `CreateBlueprint()`, `DeleteBlueprint()`, `RunBlueprint()`, `ListBlueprintRuns()`
- Added `ListSchedules()`, `GetSchedule()`, `CreateSchedule()`, `PauseSchedule()`, `ResumeSchedule()`

## v0.5.1
- Added `UploadDocument()` — upload a document from file path via multipart POST
- Added `ListDocuments()`, `GetDocument()` — document management
- Added `ListEscalations()`, `GetEscalation()`, `ResolveEscalation()` — escalation management
- Added `HealingStats()`, `ListHealingPrecedents()`, `Diagnose()` — self-healing
- Added `CreateWebhook()`, `ListWebhooks()`, `DeleteWebhook()` — webhook subscriptions
- Added `ListLearningPatterns()`, `ListLearningEpisodes()`, `ListLearningSuggestions()` — learning insights
- Added `ListCredentials()`, `CreateCredential()`, `DeleteCredential()` — credential management
- Added `ListPolicies()`, `GetPolicy()` — policy listing (was missing, TS/Python had it)

## v0.5.0 (BREAKING)
- **BREAKING:** API endpoints moved to Cadreen surface (`/api/v1/cadreen/`):
  - `/api/v1/chat/completions` → `/api/v1/cadreen/chat/completions`
  - `/api/v1/tools` → `/api/v1/cadreen/tools`
- All external API calls now route through the Cadreen surface
- Renamed response profile levels for clarity

## v0.4.1
- Added `ChatStreamEvent.RawJSON` — raw JSON bytes before typed parsing; enables extracting `pending_actions`, `conversation_id`, and `intelligence` fields without SDK type changes
- Fixed unused `bytes` import in `executions.go`

## v0.4.0
- Added `ChatCompletions()` — OpenAI-compatible chat completions with governance
- Added `ChatCompletionsStream()` — streaming chat completions via SSE
- Added `ListTools()` — discover available tools as OpenAI function definitions
- Added tool calling support: `Tools` param, `ToolCalls` in responses, tool chaining
- Added `ConversationID` for persistent conversations across requests
- Added `ChatMessage`, `ChatToolDefinition`, `ChatToolCall` and related types
- Added `ExecutionStream()` and `GetExecutionStatus()` for execution monitoring

## v0.3.0
- Added `Catalog()` — browse the unified integration marketplace
- Added `Install()` — one-click install with OAuth flow

## v0.2.0
- Initial public release
