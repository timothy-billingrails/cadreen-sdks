# Changelog

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
