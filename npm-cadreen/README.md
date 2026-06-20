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

## License

UNLICENSED
