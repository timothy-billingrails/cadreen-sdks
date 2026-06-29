# Changelog

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
