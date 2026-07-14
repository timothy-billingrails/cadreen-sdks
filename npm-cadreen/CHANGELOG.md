# Changelog

## 0.4.0

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

**Fixes:**

- Fixed CLI federation command: `target_workspace_id` → `targetWorkspaceId` (was reading wrong JSON key).

**Features:**

- Add `agents` command group: create, list, get, update, delete, deploy, capabilities, ask, messages, executions, knowledge, knowledge-add, knowledge-search, governance, audit, negotiate, negotiations
- Add `external-agents` command group: connect, list, approve, suspend, revoke, settings
- Add `responses` command group: create, get, stream
- Add `federation` command group: create, list, get, approve, suspend, revoke, permissions, agents

## 0.3.3

**Fixes:**

- Fix CLI version string — was `0.2.5`, now matches npm package version
- Fix `traces` output formatting — add newline between confidence and duration lines

## 0.3.2

- Add reasoning display to `chat` and `ask` commands — thinking model reasoning shown in gray during streaming

## 0.3.1

- Add `--user-id` flag to `ask` and `chat` commands — pass end-user identity for per-user context
- Add `workspace users` command group: `list`, `invite`, `role`, `remove` — manage workspace team members

## 0.3.0

**Breaking changes:**

- `cadreen setup-session <cmd>` renamed to `cadreen setup session <cmd>` — nested under `setup` as a subcommand
- `cadreen setup` flags (`--purpose`, `--dry-run`, etc.) moved to `cadreen setup run`

**Improvements:**

- `cadreen proposals`: added description and 5 usage examples
- `cadreen learning`: added 3 usage examples (patterns, episodes, suggestions)

## 0.2.5

- Added `setup-session` commands: create, list, get, add, apply — stateful setup sessions for incremental workspace configuration
