#!/usr/bin/env node

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const CADREEN_BASE_URL = process.env.CADREEN_BASE_URL || "https://accomplishanything.today";
const CADREEN_API_KEY = process.env.CADREEN_API_KEY || "";

if (!CADREEN_API_KEY) {
  console.error("CADREEN_API_KEY environment variable is required");
  process.exit(1);
}

async function cadreenRequest(method: string, path: string, body?: Record<string, unknown>) {
  const url = `${CADREEN_BASE_URL}${path}`;
  const headers: Record<string, string> = {
    "Authorization": `Bearer ${CADREEN_API_KEY}`,
    "Content-Type": "application/json",
  };

  const response = await fetch(url, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Cadreen API error (${response.status}): ${text}`);
  }

  return response.json();
}

const server = new McpServer({
  name: "cadreen",
  version: "0.1.2",
  description: "MCP server for Cadreen — intelligence infrastructure for AI agents",
});

// ─── Intent Tool ───

server.tool(
  "cadreen_intent",
  "Send a request to Cadreen. Describe what you want done, and Cadreen will classify, plan, and execute it. Returns a direct answer, clarification questions, or starts a mission.",
  {
    message: z.string().describe("What you want Cadreen to do"),
    mode: z.enum(["auto", "chat", "execution"]).optional().describe("Processing mode: auto (default), chat (conversation), or execution (run tools)"),
  },
  async ({ message, mode }) => {
    const result = await cadreenRequest("POST", "/api/v1/cadreen/intent", {
      messages: [{ role: "user", content: message }],
      mode: mode || "auto",
    });

    const typedResult = result as Record<string, unknown>;
    return {
      content: [{
        type: "text",
        text: JSON.stringify(typedResult, null, 2),
      }],
    };
  }
);

// ─── Agent Tools ───

server.tool(
  "cadreen_agents_list",
  "List all agents in your workspace. Agents are autonomous workers that handle tasks, follow rules, and learn from outcomes.",
  {},
  async () => {
    const result = await cadreenRequest("GET", "/api/v1/cadreen/agents");
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_agents_create",
  "Create a new agent. Give it a name and description of what it does.",
  {
    name: z.string().describe("Agent name"),
    description: z.string().optional().describe("What this agent does"),
  },
  async ({ name, description }) => {
    const result = await cadreenRequest("POST", "/api/v1/cadreen/agents", {
      name,
      description,
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_agents_get",
  "Get details about a specific agent, including its status, health, and configuration.",
  {
    agent_id: z.string().describe("Agent ID"),
  },
  async ({ agent_id }) => {
    const result = await cadreenRequest("GET", `/api/v1/cadreen/agents/${agent_id}`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── Knowledge Tools ───

server.tool(
  "cadreen_knowledge_search",
  "Search an agent's knowledge base. Find facts, procedures, and precedents the agent has learned.",
  {
    agent_id: z.string().describe("Agent ID"),
    query: z.string().describe("What to search for"),
  },
  async ({ agent_id, query }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/knowledge/search`, {
      query,
      limit: 10,
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_knowledge_add",
  "Teach an agent something new. Add knowledge that the agent can use in future tasks.",
  {
    agent_id: z.string().describe("Agent ID"),
    subject: z.string().describe("What you're teaching (e.g., 'refund policy')"),
    content: z.string().describe("The knowledge content"),
    type: z.enum(["reference", "procedure", "preference"]).optional().describe("Knowledge type"),
  },
  async ({ agent_id, subject, content, type }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/knowledge`, {
      subject,
      content,
      type: type || "reference",
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── Governance Tools ───

server.tool(
  "cadreen_governance_list",
  "List governance policies for an agent. Policies control what the agent can and cannot do.",
  {
    agent_id: z.string().describe("Agent ID"),
  },
  async ({ agent_id }) => {
    const result = await cadreenRequest("GET", `/api/v1/cadreen/agents/${agent_id}/governance`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_governance_create",
  "Create a governance policy for an agent. Control what actions require approval, what's blocked, and what's auto-approved.",
  {
    agent_id: z.string().describe("Agent ID"),
    name: z.string().describe("Policy name"),
    description: z.string().optional().describe("What this policy does"),
    rules: z.array(z.record(z.unknown())).optional().describe("Policy rules"),
  },
  async ({ agent_id, name, description, rules }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/governance`, {
      name,
      description,
      rules,
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── Federation Tools ───

server.tool(
  "cadreen_federation_list",
  "List federation links. Federation lets workspaces share agents and knowledge.",
  {},
  async () => {
    const result = await cadreenRequest("GET", "/api/v1/cadreen/federation");
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_federation_create",
  "Create a federation link to another workspace. Enables cross-workspace agent collaboration.",
  {
    target_workspace_id: z.string().describe("Target workspace ID"),
  },
  async ({ target_workspace_id }) => {
    const result = await cadreenRequest("POST", "/api/v1/cadreen/federation", {
      targetWorkspaceId: target_workspace_id,
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── Responses Tool (OpenAI-compatible) ───

server.tool(
  "cadreen_responses_create",
  "Create an OpenAI-compatible response. Cadreen handles tool calling, governance, and memory automatically.",
  {
    input: z.string().describe("The user's message or prompt"),
    model: z.string().optional().describe("Model to use (default: cadreen)"),
    stream: z.boolean().optional().describe("Whether to stream the response"),
  },
  async ({ input, model, stream }) => {
    const result = await cadreenRequest("POST", "/api/v1/cadreen/responses", {
      input,
      model: model || "cadreen",
      stream: stream || false,
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── External Agents (A2A) Tools ───

server.tool(
  "cadreen_external_agents_list",
  "List external agent connections. External agents are from other systems (LangChain, CrewAI, etc.) that can send tasks to your agents.",
  {
    agent_id: z.string().describe("Agent ID"),
  },
  async ({ agent_id }) => {
    const result = await cadreenRequest("GET", `/api/v1/cadreen/agents/${agent_id}/external`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_connect",
  "Connect to an external agent by providing its Agent Card URL. The connection must be approved before it becomes active.",
  {
    agent_id: z.string().describe("Agent ID"),
    agent_card_url: z.string().describe("URL to the agent's Agent Card (/.well-known/agent.json)"),
  },
  async ({ agent_id, agent_card_url }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/external`, {
      agentCardUrl: agent_card_url,
    });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_get",
  "Get details about a specific external agent connection, including its status, health, and capabilities.",
  {
    agent_id: z.string().describe("Agent ID"),
    connection_id: z.string().describe("Connection ID"),
  },
  async ({ agent_id, connection_id }) => {
    const result = await cadreenRequest("GET", `/api/v1/cadreen/agents/${agent_id}/external/${connection_id}`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_approve",
  "Approve a pending external agent connection. The connection becomes active and the external agent can send tasks.",
  {
    agent_id: z.string().describe("Agent ID"),
    connection_id: z.string().describe("Connection ID"),
  },
  async ({ agent_id, connection_id }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/external/${connection_id}/approve`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_suspend",
  "Suspend an active external agent connection. The external agent can no longer send tasks until re-approved.",
  {
    agent_id: z.string().describe("Agent ID"),
    connection_id: z.string().describe("Connection ID"),
  },
  async ({ agent_id, connection_id }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/external/${connection_id}/suspend`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_revoke",
  "Revoke an external agent connection permanently. This cannot be undone.",
  {
    agent_id: z.string().describe("Agent ID"),
    connection_id: z.string().describe("Connection ID"),
  },
  async ({ agent_id, connection_id }) => {
    const result = await cadreenRequest("POST", `/api/v1/cadreen/agents/${agent_id}/external/${connection_id}/revoke`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_delete",
  "Delete an external agent connection. Use revoke first if the connection is active.",
  {
    agent_id: z.string().describe("Agent ID"),
    connection_id: z.string().describe("Connection ID"),
  },
  async ({ agent_id, connection_id }) => {
    await cadreenRequest("DELETE", `/api/v1/cadreen/agents/${agent_id}/external/${connection_id}`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify({ status: "deleted" }, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_list_interactions",
  "List interactions (tasks sent/received) for an external agent connection.",
  {
    agent_id: z.string().describe("Agent ID"),
    connection_id: z.string().describe("Connection ID"),
  },
  async ({ agent_id, connection_id }) => {
    const result = await cadreenRequest("GET", `/api/v1/cadreen/agents/${agent_id}/external/${connection_id}/interactions`);
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_get_settings",
  "Get external agent settings for the workspace. Shows whether external agents are enabled.",
  {},
  async () => {
    const result = await cadreenRequest("GET", "/api/v1/cadreen/external-agents/settings");
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_update_settings",
  "Enable or disable external agents for the workspace. Must be enabled before connecting to external agents.",
  {
    enabled: z.boolean().describe("Whether to enable or disable external agents"),
  },
  async ({ enabled }) => {
    const result = await cadreenRequest("PUT", "/api/v1/cadreen/external-agents/settings", { enabled });
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

server.tool(
  "cadreen_external_agents_list_all",
  "List all external agent connections across all agents in the workspace.",
  {},
  async () => {
    const result = await cadreenRequest("GET", "/api/v1/cadreen/external-agents/connections");
    return {
      content: [{
        type: "text",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── Resources ───

server.resource(
  "agent-card",
  "agent-card:///{agent_id}",
  async (uri) => {
    const agentId = uri.pathname.split("/").pop();
    const result = await cadreenRequest("GET", `/api/v1/cadreen/agents/${agentId}`);
    return {
      contents: [{
        uri: uri.href,
        mimeType: "application/json",
        text: JSON.stringify(result, null, 2),
      }],
    };
  }
);

// ─── Prompts ───

server.prompt(
  "create_agent",
  "Create a new agent with knowledge and governance",
  {
    name: z.string().describe("Agent name"),
    purpose: z.string().describe("What the agent does"),
    rules: z.string().optional().describe("Governance rules for the agent"),
  },
  async ({ name, purpose, rules }) => {
    const messages = [
      {
        role: "user" as const,
        content: `Create an agent called "${name}" that ${purpose}.${rules ? ` Rules: ${rules}` : ""}`,
      },
    ];
    return { messages };
  }
);

server.prompt(
  "teach_agent",
  "Teach an agent something new",
  {
    agent_id: z.string().describe("Agent ID"),
    subject: z.string().describe("What to teach"),
    content: z.string().describe("The knowledge content"),
  },
  async ({ agent_id, subject, content }) => {
    const messages = [
      {
        role: "user" as const,
        content: `Add knowledge to agent ${agent_id}: ${subject} = ${content}`,
      },
    ];
    return { messages };
  }
);

// ─── Start Server ───

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("Cadreen MCP server running on stdio");
}

main().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});
