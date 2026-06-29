# Changelog

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
