# Changelog

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
