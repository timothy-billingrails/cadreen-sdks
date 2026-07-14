# Changelog

## 0.1.0

- Initial release
- 22 tools covering the full Cadreen API surface:
  - Intent (cadreen_intent)
  - Agents (list, create, get)
  - Knowledge (search, add)
  - Governance (list, create)
  - Federation (list, create)
  - Responses (create — OpenAI-compatible)
  - External Agents/A2A (list, connect, approve, suspend, revoke, delete, settings, interactions)
- Resources: agent-card:///{agent_id}
- Prompts: create_agent, teach_agent
- stdio transport for Claude Desktop, Cursor, and any MCP-compatible client
- Environment variables: CADREEN_API_KEY (required), CADREEN_BASE_URL (optional, default: https://accomplishanything.today)
- Hosted MCP SSE endpoint available at `/api/v1/cadreen/mcp/sse` — connect without installing the npm package
