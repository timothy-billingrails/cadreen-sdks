# Changelog

## 0.7.2

**Features:**

- Add `DevicesResource`: full device lifecycle — list, create, get, delete, get_status, update_state, get_map, get_map_stats, update_map, list_tasks, create_task, complete_task, assign_tasks, detect_collisions, get_avoidance, diagnose, ask, get_model_stats, get_capabilities, get_sync_status, get_sync_pending, get_sync_conflicts, get_blackboard

**Fixes:**

- Fixed `DiagnoseRequest` type shadowing — renamed device-specific `DiagnoseRequest` (with `readings`) to `DeviceDiagnoseRequest` to avoid conflict with healing `DiagnoseRequest`
- Fixed missing `AsyncIterator` import in `resources/intent.py`
- Fixed Python connection pooling — `HttpClient` now reuses a persistent `httpx.AsyncClient` instead of creating a new one per request; added `__aenter__`/`__aexit__` for async context manager support

## 0.7.1

**Fixes:**

- Sandbox mode: `post_stream` and `stream` now raise `CadreenError(404)` instead of making network requests in sandbox mode.
- Release workflow: fixed `git stash pop` ordering in Homebrew formula update step.

**Known server-side events (for consumers polling/streaming executions):**

- `execution_steps_complete` — step loop finished; synthesis, quality review, and delivery governance still follow. NOT terminal.
- `mission_complete` — mission completed successfully. Terminal.
- `mission_completed_with_gaps` — steps succeeded but delivery governance blocked the deliverable. Terminal.
- `mission_blocked` — execution failed. Terminal.
- `mission_partial` — execution partially completed. Terminal.

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

- Add `AgentsResource`: full agent lifecycle — create, list, get, update, delete, deploy, get_config, get_capabilities
- Add agent messaging — send_message, list_messages
- Add agent executions — create_execution, list_executions
- Add agent knowledge — create_knowledge, list_knowledge, search_knowledge, delete_knowledge
- Add agent governance — create_governance, list_governance, update_governance, delete_governance
- Add agent audit trail — list_audit
- Add agent negotiations — start_negotiation, list_negotiations, get_negotiation, respond_to_negotiation
- Add `FederationResource`: cross-workspace federation — create, list, get, approve, suspend, revoke
- Add federation permissions — get_permissions, update_permissions
- Add federation agent linking — link_agent, list_agents, unlink_agent
- Add `ExternalAgentsResource`: A2A external agent connections — connect, list, get, approve, suspend, revoke, delete
- Add external agent interactions — list_interactions
- Add external agent settings — get_settings, update_settings, list_all
- Add `ResponsesResource`: OpenAI-compatible responses API — create, retrieve, stream
- Add 37 new types: `Agent`, `AgentKnowledge`, `AgentGovernancePolicy`, `AgentNegotiation`, `FederationLink`, `FederationAgent`, `ExternalAgentConnection`, `ExternalAgentInteraction`, `ExternalAgentSettings`, `ExternalAgentSkill`, `ExternalAgentCapabilities`, etc.

## 0.6.3

**Fixes:**

- Fix `invoke_stream()` crash — was using nonexistent `_client._session`, now uses `httpx-sse` `aconnect_sse`
- Fix `chat.completions()` and `completions_stream()` silently dropping `user_id` from request body

## 0.6.2

**Features:**

- Add `intelligence` field to `ChatCompletionResponse` — full intelligence metadata
- Add `conversation_id` field to `ChatCompletionResponse` — conversation continuity
- Add `reasoning` field to `ChatDelta` — streaming reasoning from thinking models
- Add `reasoning_tokens`, `cache_write_tokens`, `prompt_tokens_details` to `ChatUsage`
- Add `BlueprintsResource`: list, get, create, update, delete, run, get_runs
- Add `SchedulesResource`: list, get, create, update, delete, pause, resume, get_runs
- Add 16 new types: `Blueprint`, `BlueprintRun`, `Schedule`, `ScheduleRun`, etc.
- Document `reasoning_delta` streaming event in README

**Fixes:**

- IntelligenceMeta shape aligned with server

## 0.6.1

- Add `user_id` optional field to `IntentRequest` and `ChatCompletionRequest` — pass end-user identity for per-user context and memory filtering
- Add `WorkspaceUsersResource`: list, invite, update_role, remove — manage workspace team members
- Add `WorkspaceUser`, `WorkspaceRole`, `InviteUserRequest`, `UpdateRoleRequest`, `ListWorkspaceUsersResponse` types
- Add `HttpClient.patch()` method for PATCH requests

## 0.6.0

**Breaking changes:**

- `ResolveEscalationRequest.resolution` renamed to `decision` — aligns SDK with API contract
- `DiagnoseRequest.error` renamed to `error_message` — aligns SDK with API contract

**Fixes:**

- Retry default increased from 2 to 3

## 0.5.5

- Add `SetupSessionsResource`: create, list, get, add_resources, apply for stateful setup sessions
- Add `SetupSession`, `SetupSessionCreateRequest`, `SetupSessionAddRequest`, `SetupSessionApplyRequest`, `SetupSessionApplyResult` types

## 0.5.4

- Add `ProposalsResource`: list, get, accept, dismiss, stats for task proposals
