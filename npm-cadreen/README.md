# @cadreen/cli

Cadreen CLI — intelligence infrastructure for developers.

> **⚠️ Compatibility Advisory:** Version 0.4.0 has known contract mismatches with the server. Version 0.4.0 is unsupported. **Upgrade to 0.4.2.**

Cadreen remembers things, follows rules, connects to services, and heals itself when things go wrong. All governed. All observable.

## Install

```bash
npm install -g @cadreen/cli
```

## Quick Start

```bash
cadreen init        # Set up your account
cadreen ask "what can you do?"  # Ask a question
cadreen chat        # Interactive chat
cadreen doctor      # Check readiness
```

## Commands

| Command | What it does |
|---------|-------------|
| `cadreen init` | Set up your account |
| `cadreen login` | Authenticate |
| `cadreen ask "..."` | One-shot question |
| `cadreen chat` | Interactive chat |
| `cadreen status` | System health |
| `cadreen doctor` | Readiness check |
| `cadreen memory list` | What Cadreen knows |
| `cadreen memory add "..."` | Teach something |
| `cadreen policies list` | Active rules |
| `cadreen policies evaluate "..."` | Test an action |
| `cadreen tools` | Available tools |
| `cadreen traces` | What happened |
| `cadreen documents list` | List documents |
| `cadreen documents get [id]` | Get document details |
| `cadreen documents upload [filepath]` | Upload a document |
| `cadreen escalations list` | List escalations |
| `cadreen escalations get [id]` | Get escalation details |
| `cadreen escalations resolve [id] "..."` | Resolve an escalation |
| `cadreen healing stats` | Healing statistics |
| `cadreen healing precedents` | Healing precedents |
| `cadreen healing diagnose "..."` | Diagnose a failure |
| `cadreen webhooks list` | List webhooks |
| `cadreen webhooks create [url]` | Create webhook |
| `cadreen webhooks delete [id]` | Delete webhook |
| `cadreen learning patterns` | Detected patterns |
| `cadreen learning episodes` | Learning episodes |
| `cadreen learning suggestions` | Improvement suggestions |
| `cadreen credentials list` | List credentials |
| `cadreen credentials delete [id]` | Delete credential |
| `cadreen proposals list` | Task proposals waiting for you |
| `cadreen proposals get [id]` | Get proposal details |
| `cadreen proposals accept [id]` | Accept a proposal |
| `cadreen proposals dismiss [id]` | Dismiss a proposal |
| `cadreen proposals stats` | Proposal statistics |
| `cadreen setup run` | One-shot workspace setup |
| `cadreen setup session create` | Create a setup session |
| `cadreen setup session list` | List setup sessions |
| `cadreen setup session get [id]` | Get session details |
| `cadreen setup session add [id]` | Add resources to a session |
| `cadreen setup session apply [id]` | Apply a session atomically |
| `cadreen workspace users list` | List workspace users |
| `cadreen workspace users invite [email]` | Invite a user to the workspace |
| `cadreen workspace users role [id] [role]` | Update a user's role |
| `cadreen workspace users remove [id]` | Remove a user from the workspace |
| `cadreen agents create [name]` | Create an agent |
| `cadreen agents list` | List agents |
| `cadreen agents get [id]` | Get agent details |
| `cadreen agents deploy [id]` | Deploy an agent |
| `cadreen agents ask [id] "..."` | Ask an agent a question |
| `cadreen agents messages [id]` | List agent messages |
| `cadreen agents executions [id]` | List agent executions |
| `cadreen agents knowledge [id]` | List agent knowledge |
| `cadreen agents governance [id]` | List agent governance rules |
| `cadreen agents audit [id]` | List agent audit entries |
| `cadreen federation create` | Create a federation link |
| `cadreen federation list` | List federation links |
| `cadreen federation get [id]` | Get federation link details |
| `cadreen federation approve [id]` | Approve a federation link |
| `cadreen federation suspend [id]` | Suspend a federation link |
| `cadreen external-agents connect [url]` | Connect to an external agent |
| `cadreen external-agents list` | List external agent connections |
| `cadreen external-agents approve [id]` | Approve an external agent |
| `cadreen external-agents suspend [id]` | Suspend an external agent |
| `cadreen external-agents revoke [id]` | Revoke an external agent |
| `cadreen external-agents settings` | Manage external agent settings |
| `cadreen responses create "..."` | Create a response |
| `cadreen responses get [id]` | Get a response |
| `cadreen devices list` | List devices |
| `cadreen devices get [id]` | Get device details |
| `cadreen devices diagnose` | Diagnose a device |
| `cadreen config` | Local settings |
| `cadreen update` | Update CLI |

## Other Install Methods

**Homebrew:**
```bash
brew tap timothy-billingrails/cadreen-sdks
brew install cadreen
```

**curl:**
```bash
curl -fsSL https://raw.githubusercontent.com/timothy-billingrails/cadreen-sdks/master/install.sh | sh
```

**Go:**
```bash
go install github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen@latest
```

## Links

- [Docs](https://accomplishanything.today/infra/docs)
- [GitHub](https://github.com/timothy-billingrails/cadreen-sdks)

## Changelog

### v0.4.2
- Added `devices` command group — full device management (list, create, get, delete, status, state, map, tasks, collisions, avoidance, diagnose, ask, sync, blackboard)

### v0.4.0
- **BREAKING:** Removed `pathways` and `total_pathways` from connection responses. `ConnectionGroup` now returns only `capability` and `status`.
- **BREAKING:** Removed `Pathway` type. Internal routing details (connector, transport, tool_id) are no longer exposed.
- **BREAKING:** Changed `ConnectManualDetail` from `{pathways: [...]}` to `{capability, available, health}`.
- **BREAKING:** Removed `workspace_id` from response types: `SetupResult`, `SetupSession`, `WebhookSubscription`, `WebhookPayload`. (Still accepted on request types.)
- **BREAKING:** Removed `authScheme` from `ExternalAgentConnection` responses.
- **BREAKING:** Removed `atoms_consulted`, `episodes_matched`, `precedents_applied` from memory trace in intelligence metadata.
- **BREAKING:** `sources_consulted` renamed to `knowledge_queried` in MemoryTrace.
- **BREAKING:** All entity responses (Agent, Knowledge, Governance, Federation, Negotiation, ExternalAgentConnection) no longer include `workspace_id` or `workspaceId`.
- New MCP SSE endpoint: `GET /api/v1/cadreen/mcp/sse` + `POST /api/v1/cadreen/mcp/message`. Connect Cadreen as an MCP server without installing the npm package.
- Fix: CLI federation command `target_workspace_id` → `targetWorkspaceId` (was reading wrong JSON key)
- Added `cadreen agents` — 17 subcommands: create, list, get, update, delete, deploy, capabilities, ask, messages, executions, knowledge, knowledge-add, knowledge-search, governance, audit, negotiate, negotiations
- Added `cadreen federation` — 9 subcommands: create, list, get, approve, suspend, revoke, permissions, link-agent, agents
- Human language in all command descriptions

### v0.3.3
- Fix: CLI version string was `0.2.5`, now matches npm package version
- Fix: `traces` output — add newline between confidence and duration lines

### v0.3.2
- Added reasoning display to `chat` and `ask` commands — thinking model reasoning shown in gray during streaming

### v0.3.1
- Added `--user-id` flag to `ask` and `chat` commands — pass end-user identity for per-user context
- Added `workspace users` command group: `list`, `invite`, `role`, `remove` — manage workspace team members

### v0.3.0 (BREAKING)
- **BREAKING:** `cadreen setup-session` renamed to `cadreen setup session` — nested under `setup` as a subcommand
- **BREAKING:** `cadreen setup` flags (`--purpose`, `--dry-run`, etc.) moved to `cadreen setup run`
- `cadreen proposals`: added description and usage examples
- `cadreen learning`: added usage examples

### v0.2.5
- Added `setup-session` commands: create, list, get, add, apply — stateful setup sessions for incremental workspace configuration

### v0.2.4
- Added `proposals` commands: list, get, accept, dismiss, stats
- Flags: `--status`, `--limit`, `--reason`

### v0.2.3
- Updated npm description and keywords for discoverability (openai, ai, automation, orchestration, llm, agent, tool-calling, policy)

### v0.2.2
- Added `setup` command with `--dry-run` flag — preview without creating
- Added `setup --purpose`, `--memory`, `--policy` flags
- Added `blueprints` commands (list, create, run, show, runs, archive)
- Added `schedules` commands (list, create, pause, resume, show)
- Added `policies create` command

### v0.2.1
- Added `documents upload` command — upload files via multipart POST
- Added `documents` commands (list, get)
- Added `escalations` commands (list, get, resolve)
- Added `healing` commands (stats, precedents, diagnose)
- Added `webhooks` commands (list, create, delete)
- Added `learning` commands (patterns, episodes, suggestions)
- Added `credentials` commands (list, delete)

### v0.2.0
- Updated to SDK v0.5.0 (Cadreen surface paths)
- `cadreen tools` now calls `/api/v1/cadreen/tools`

### v0.1.0
- Initial release

## License

UNLICENSED
