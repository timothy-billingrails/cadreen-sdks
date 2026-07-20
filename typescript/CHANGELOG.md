# Changelog

## 0.7.2

**Features:**

- Add `DevicesResource`: full device lifecycle — list, create, get, delete, getStatus, updateState, getMap, getMapStats, updateMap, listTasks, createTask, completeTask, assignTasks, detectCollisions, getAvoidance, diagnose, ask, getModelStats, getCapabilities, getSyncStatus, getSyncPending, getSyncConflicts, getBlackboard

**Fixes:**

- Fixed `DiagnoseRequest` type shadowing — renamed device-specific `DiagnoseRequest` (with `readings: SensorReading[]`) to `DeviceDiagnoseRequest` to avoid conflict with healing `DiagnoseRequest`

## 0.7.1

**Fixes:**

- Sandbox mode: `postStream` and `stream` now throw `CadreenError(404)` instead of making network requests in sandbox mode.
- Release workflow: fixed `git stash pop` ordering in Homebrew formula update step.
- Homebrew formula SHA256 update now runs before commit in release workflow.

**Known server-side events (for consumers polling/streaming executions):**

- `execution_steps_complete` — step loop finished; synthesis, quality review, and delivery governance still follow. NOT terminal.
- `mission_complete` — mission completed successfully. Terminal.
- `mission_completed_with_gaps` — steps succeeded but delivery governance blocked the deliverable. Terminal.
- `mission_blocked` — execution failed. Terminal.
- `mission_partial` — execution partially completed. Terminal.
- `execution_complete` — legacy event; prefer `execution_steps_complete`.

## 0.7.0

**Breaking Changes:**

- Removed `pathways` and `total_pathways` from connection responses. `ConnectionGroup` now returns only `capability` and `status`.
- Removed `Pathway` type. Internal routing details (connector, transport, tool_id) are no longer exposed.
- Changed `ConnectManualDetail` from `{pathways: [...]}` to `{capability, available, health}`.
- Removed `workspace_id` from response types: `SetupResult`, `SetupSession`, `WebhookSubscription`, `WebhookPayload`. (Still accepted on request types.)
- Removed `authScheme` from `ExternalAgentConnection` responses.
- Removed `atoms_consulted`, `episodes_matched`, `precedents_applied` from memory trace in intelligence metadata.
- `sources_consulted` renamed to `knowledge_queried` in MemoryTrace.
- All entity responses (Agent, Knowledge, Governance, Federation, Negotiation, ExternalAgentConnection) no longer include `workspace_id` or `workspaceId`.
- New MCP SSE endpoint: `GET /api/v1/cadreen/mcp/sse` + `POST /api/v1/cadreen/mcp/message`. Connect Cadreen as an MCP server without installing the npm package.

**Features:**

- Add `AgentsResource`: full agent lifecycle — create, list, get, update, delete, deploy, getConfig, getCapabilities
- Add agent messaging: `sendMessage`, `listMessages`
- Add agent executions: `createExecution`, `listExecutions`
- Add agent knowledge: `createKnowledge`, `searchKnowledge`, `listKnowledge`, `deleteKnowledge`
- Add agent governance: `createGovernance`, `updateGovernance`, `listGovernance`, `deleteGovernance`
- Add agent audit: `listAudit`
- Add agent negotiations: `startNegotiation`, `respondToNegotiation`, `getNegotiation`, `listNegotiations`
- Add `FederationResource`: cross-workspace federation — create, list, get, approve, suspend, revoke
- Add federation permissions: `getPermissions`, `updatePermissions`
- Add federation agents: `linkAgent`, `listAgents`, `unlinkAgent`
- Add `ExternalAgentsResource`: A2A external agent connections — connect, list, get, approve, suspend, revoke, delete
- Add external agent interactions: `listInteractions`
- Add external agent settings: `getSettings`, `updateSettings`, `listAll`
- Add `ResponsesResource`: OpenAI-compatible responses API — create, retrieve, stream
- Add types: `Agent`, `AgentKnowledge`, `AgentGovernancePolicy`, `AgentAuditEntry`, `AgentNegotiation`, `AgentMessage`, `AgentExecution`, `FederationLink`, `FederationAgent`, `FederationPermissions`, `ExternalAgentConnection`, `ExternalAgentInteraction`, `ExternalAgentSettings`, `ExternalAgentSkill`, `ExternalAgentCapabilities` and all related request/response types

## 0.6.3

**Fixes:**

- Add missing `reasoning_tokens`, `cache_write_tokens`, `prompt_tokens_details` fields to `ChatUsage` — were claimed in v0.6.2 changelog but not present in code

## 0.6.2

**Features:**

- Add `intelligence` field to `ChatCompletionResponse` — full intelligence metadata (capability, reasoning, memory, governance, humility, process traces)
- Add `conversation_id` field to `ChatCompletionResponse` — conversation continuity across requests
- Add `reasoning` field to `ChatDelta` — streaming reasoning from thinking models
- Add `reasoning_tokens`, `cache_write_tokens`, `prompt_tokens_details` to `ChatUsage`
- Document `reasoning_delta` streaming event in README

**Fixes:**

- IntelligenceMeta shape aligned with server

## 0.6.1

- Add `user_id` optional field to `IntentRequest` and `ChatCompletionRequest` — pass end-user identity for per-user context and memory filtering
- Add `WorkspaceUsersResource`: list, invite, updateRole, remove — manage workspace team members
- Add `WorkspaceUser`, `WorkspaceRole`, `InviteUserRequest`, `UpdateRoleRequest`, `ListWorkspaceUsersResponse` types
- Add `HttpClient.patch()` method for PATCH requests

## 0.6.0

**Breaking changes:**

- `ResolveEscalationRequest.resolution` renamed to `decision` — aligns SDK with API contract
- `DiagnoseRequest.error` renamed to `error_message` — aligns SDK with API contract

**Fixes:**

- Retry default increased from 2 to 3

## 0.5.5

- Add `SetupSessionsResource`: create, list, get, addResources, apply for stateful setup sessions
- Add `SetupSession`, `SetupSessionCreateRequest`, `SetupSessionAddRequest`, `SetupSessionApplyRequest`, `SetupSessionApplyResult` types

## 0.5.4

- Add `ProposalsResource`: list, get, accept, dismiss, stats for task proposals
