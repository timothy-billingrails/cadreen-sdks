# Changelog

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
