# @cadreen/cli

Cadreen CLI — intelligence infrastructure for developers.

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
curl -fsSL https://raw.githubusercontent.com/timothy-billingrails/cadreen-sdks/main/install.sh | sh
```

**Go:**
```bash
go install github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen@latest
```

## Links

- [Docs](https://accomplishanything.today/infra/docs)
- [GitHub](https://github.com/timothy-billingrails/cadreen-sdks)

## Changelog

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
