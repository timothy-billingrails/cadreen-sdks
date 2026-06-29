# Changelog

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
