# @cadreen/mcp

MCP server for [Cadreen](https://accomplishanything.today/infra/docs) — intelligence infrastructure for agents.

This server lets any MCP-compatible LLM use Cadreen's capabilities: intent routing, agent management, knowledge, governance, federation, external agents, and the Responses API.

**22 tools** covering the full Cadreen API surface.

## Install

```bash
npm install -g @cadreen/mcp
```

## Configure

Set your Cadreen API key:

```bash
export CADREEN_API_KEY="sk_cadreen_..."
```

Optional: set a custom base URL (default: `https://accomplishanything.today`):

```bash
export CADREEN_BASE_URL="https://your-instance.cadreen.com"
```

## Use with Claude Desktop

Add to your Claude Desktop config (`~/.claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "cadreen": {
      "command": "cadreen-mcp",
      "env": {
        "CADREEN_API_KEY": "sk_cadreen_..."
      }
    }
  }
}
```

## Use with any MCP client

```bash
# Run the server
cadreen-mcp

# Or with npx
npx @cadreen/mcp
```

The server communicates over stdio using the MCP protocol.

## Hosted SSE

Prefer not to install anything? Cadreen exposes a hosted MCP SSE endpoint — connect any MCP-compatible client directly to Cadreen without installing the npm package.

- **SSE endpoint:** `GET /api/v1/cadreen/mcp/sse`
- **Message endpoint:** `POST /api/v1/cadreen/mcp/message`

Point your MCP client at `https://accomplishanything.today/api/v1/cadreen/mcp/sse` with your `CADREEN_API_KEY` and Cadreen acts as an MCP server over SSE.

## Available Tools

| Tool | Description |
|------|-------------|
| **Intent** | |
| `cadreen_intent` | Send a request. Describe what you want done. |
| **Agents** | |
| `cadreen_agents_list` | List all agents in your workspace. |
| `cadreen_agents_create` | Create a new agent. |
| `cadreen_agents_get` | Get agent details, status, and health. |
| **Knowledge** | |
| `cadreen_knowledge_search` | Search an agent's knowledge base. |
| `cadreen_knowledge_add` | Teach an agent something new. |
| **Governance** | |
| `cadreen_governance_list` | List governance policies. |
| `cadreen_governance_create` | Create a governance policy. |
| **Federation** | |
| `cadreen_federation_list` | List federation links. |
| `cadreen_federation_create` | Create a federation link. |
| **Responses** | |
| `cadreen_responses_create` | Create an OpenAI-compatible response. |
| **External Agents (A2A)** | |
| `cadreen_external_agents_list` | List external A2A connections for an agent. |
| `cadreen_external_agents_list_all` | List all external connections across workspace. |
| `cadreen_external_agents_connect` | Connect to an external agent via Agent Card URL. |
| `cadreen_external_agents_get` | Get connection details, status, and capabilities. |
| `cadreen_external_agents_approve` | Approve a pending connection. |
| `cadreen_external_agents_suspend` | Suspend an active connection. |
| `cadreen_external_agents_revoke` | Revoke a connection permanently. |
| `cadreen_external_agents_delete` | Delete a connection. |
| `cadreen_external_agents_list_interactions` | List tasks sent/received for a connection. |
| `cadreen_external_agents_get_settings` | Get workspace external agent settings. |
| `cadreen_external_agents_update_settings` | Enable/disable external agents. |

## Resources

| Resource | Description |
|----------|-------------|
| `agent-card:///{agent_id}` | Get agent details as a resource. |

## Prompts

| Prompt | Description |
|--------|-------------|
| `create_agent` | Create a new agent with knowledge and governance. |
| `teach_agent` | Teach an agent something new. |

## Examples

### Ask Cadreen to do something

```
Use cadreen_intent with message: "Refund order 12345 and notify the customer"
```

### Teach an agent something

```
Use cadreen_knowledge_add with:
  agent_id: "agent_123"
  subject: "refund policy"
  content: "Refunds require manager approval for amounts over $100"
```

### Create a governance policy

```
Use cadreen_governance_create with:
  agent_id: "agent_123"
  name: "refund approval"
  description: "Require approval for refunds over $50"
```

## License

MIT
