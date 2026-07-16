# Changelog

## 0.7.1

**Fixes:**

- Sandbox mode: `PostStream` and `Stream` now return `ErrSandboxStreaming` instead of making network requests in sandbox mode.
- Release workflow: fixed `git stash pop` ordering in Homebrew formula update step.

**Known server-side events (for consumers polling/streaming executions):**

- `execution_steps_complete` — step loop finished; synthesis, quality review, and delivery governance still follow. NOT terminal.
- `mission_complete` — mission completed successfully. Terminal.
- `mission_completed_with_gaps` — steps succeeded but delivery governance blocked the deliverable. Terminal.
- `mission_blocked` — execution failed. Terminal.
- `mission_partial` — execution partially completed. Terminal.

## 0.7.0

**Breaking Changes:**

- Removed `Pathways` and `TotalPathways` from connection responses. `ConnectionGroup` now returns only `Capability` and `Status`.
- Removed `Pathway` type. Internal routing details (connector, transport, tool_id) are no longer exposed.
- Changed `ConnectManualDetail` from `{Pathways: [...]}` to `{Capability, Available, Health}`.
- Removed `WorkspaceID` from response types: `SetupResult`, `SetupSession`, `WebhookSubscription`, `WebhookPayload`. (Still accepted on request types.)
- Removed `AuthScheme` from `ExternalAgentConnection` responses.
- Removed `AtomsConsulted`, `EpisodesMatched`, `PrecedentsApplied` from memory trace in intelligence metadata.
- `SourcesConsulted` renamed to `KnowledgeQueried` in MemoryTrace.
- All entity responses (Agent, Knowledge, Governance, Federation, Negotiation, ExternalAgentConnection) no longer include `WorkspaceID`.
- New MCP SSE endpoint: `GET /api/v1/cadreen/mcp/sse` + `POST /api/v1/cadreen/mcp/message`. Connect Cadreen as an MCP server without installing the npm package.

**Features:**

- Add Agents resource — full agent lifecycle management:
  - `CreateAgent`, `GetAgent`, `ListAgents`, `UpdateAgent`, `DeleteAgent`
  - `GetAgentConfig`, `DeployAgent`, `GetAgentCapabilities`
  - `SendAgentMessage`, `ListAgentMessages`
  - `ListAgentExecutions`, `CreateAgentExecution`
  - `ListAgentKnowledge`, `CreateAgentKnowledge`, `SearchAgentKnowledge`, `DeleteAgentKnowledge`
  - `ListAgentGovernance`, `CreateAgentGovernance`, `UpdateAgentGovernance`, `DeleteAgentGovernance`
  - `ListAgentAudit`
  - `StartNegotiation`, `ListNegotiations`, `GetNegotiation`, `RespondToNegotiation`
- Add Federation resource — cross-organization agent linking:
  - `CreateFederation`, `GetFederation`, `ListFederations`
  - `ApproveFederation`, `SuspendFederation`, `RevokeFederation`
  - `GetFederationPermissions`, `UpdateFederationPermissions`
  - `LinkFederationAgent`, `ListFederationAgents`, `UnlinkFederationAgent`
- Add External Agents resource — A2A external agent connections:
  - `ConnectExternalAgent`, `ListExternalConnections`, `GetExternalConnection`
  - `ApproveExternalConnection`, `SuspendExternalConnection`, `RevokeExternalConnection`, `DeleteExternalConnection`
  - `ListExternalInteractions`
  - `GetExternalAgentSettings`, `UpdateExternalAgentSettings`, `ListAllExternalConnections`
- Add Responses resource — OpenAI-compatible responses API:
  - `CreateResponse`, `GetResponse`, `CreateResponseStreaming`
- Add types: `Agent`, `AgentKnowledge`, `AgentGovernancePolicy`, `AgentAuditEntry`, `AgentNegotiation`, `AgentMessage`, `FederationLink`, `FederationAgent`, `ExternalAgentConnection`, `ExternalAgentInteraction`, `ExternalAgentSettings`, `ExternalAgentSkill`, `ExternalAgentCapabilities`

## 0.6.3

**Fixes:**

- Fix retry body reuse bug — `bytes.Reader` was consumed on first attempt, retries sent empty body. Now creates fresh `bytes.NewReader` for each retry.
- Add `VerifyWebhookSignature()` function for HMAC-SHA256 webhook signature verification

## 0.6.2

**Features:**

- Add `Intelligence` field to `ChatCompletionResponse` — full intelligence metadata
- Add `ConversationID` field to `ChatCompletionResponse` — conversation continuity
- Add `Reasoning` field to `ChatDelta` — streaming reasoning from thinking models
- Add `ReasoningTokens`, `CacheWriteTokens`, `PromptTokensDetails` to `ChatUsage`
- Add `Reasoning` field to `ChatStreamEvent` — non-empty on reasoning_delta events
- Document `reasoning_delta` streaming event in README

**Fixes:**

- IntelligenceMeta shape aligned with server

## 0.6.1

- Add `UserID` field to `IntentRequest` and `ChatCompletionRequest` — pass end-user identity for per-user context and memory filtering
- Add `ListWorkspaceUsers()`, `InviteUser()`, `UpdateUserRole()`, `RemoveUser()` — manage workspace team members
- Add `WorkspaceUser`, `WorkspaceRole`, `InviteUserRequest`, `UpdateRoleRequest`, `ListWorkspaceUsersResponse` types

## 0.6.0

**Breaking changes:**

- `ResolveEscalationRequest.Resolution` renamed to `Decision` — aligns SDK with API contract
- `DiagnoseRequest.Error` renamed to `ErrorMessage` — aligns SDK with API contract

**Fixes:**

- Retry default was already at 3 — no change needed

## 0.5.5

- Add `CreateSetupSession()`, `ListSetupSessions()`, `GetSetupSession()`, `AddToSetupSession()`, `ApplySetupSession()` — stateful setup sessions
- Add `SetupSession`, `SetupSessionCreateRequest`, `SetupSessionAddRequest`, `SetupSessionApplyRequest`, `SetupSessionApplyResult` types

## 0.5.4

- Add ProposalsResource: ListProposals, GetProposal, AcceptProposal, DismissProposal, ProposalStats
