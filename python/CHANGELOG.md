# Changelog

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

- IntelligenceMeta shape aligned with server — `capability` (not `capability_assessment`), `humility` (not `epistemic_humility`)

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
