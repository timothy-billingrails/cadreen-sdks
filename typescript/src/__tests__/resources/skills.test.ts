import { describe, it, expect } from "vitest";
import { SkillsResource } from "../../resources/skills";
import { IntentResource } from "../../resources/intent";
import { MemoryResource } from "../../resources/memory";
import { ConnectionsResource } from "../../resources/connections";
import { HttpClient } from "../../client";
import type { IntelligenceMeta } from "../../types";

function baseIntelligence(): IntelligenceMeta {
  return {
    capability: { total_available: 10, healthy_count: 8 },
    reasoning: {},
    memory: { healthy: true },
    governance: { active: false },
    humility: {},
    process: { started_at: "2026-01-01T00:00:00Z", duration_ms: 100 },
    field_stability: { stable: [], evolving: [], internal: [] },
  };
}

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("SkillsResource ask()", () => {
  it("delegates to intent.invoke with mode=chat", async () => {
    const raw = {
      id: "resp-ask",
      type: "direct",
      trace_id: "trace-ask",
      message: { role: "assistant", content: "Answer from skills" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const intent = new IntentResource(client);
    const memory = new MemoryResource(client);
    const connections = new ConnectionsResource(client);
    const skills = new SkillsResource(intent, memory, connections);

    const result = await skills.ask("What's the weather?");
    expect(result.type).toBe("direct");
    if (result.type === "direct") {
      expect(result.message.content).toBe("Answer from skills");
    }
  });
});

describe("SkillsResource act()", () => {
  it("delegates to intent.invoke with mode=execution", async () => {
    const raw = {
      id: "resp-act",
      type: "mission",
      trace_id: "trace-act",
      mission: { id: "exec-10", status: "running" },
      intelligence: baseIntelligence(),
    };
    const client = buildClient({ "POST /api/v1/cadreen/intent": raw });
    const intent = new IntentResource(client);
    const memory = new MemoryResource(client);
    const connections = new ConnectionsResource(client);
    const skills = new SkillsResource(intent, memory, connections);

    const result = await skills.act("Send an email");
    expect(result.type).toBe("execution");
  });
});

describe("SkillsResource remember()", () => {
  it("delegates to memory.remember", async () => {
    const fixture = {
      id: "mem-skill",
      type: "note",
      domain: "general",
      content: { text: "Skill note" },
      authority: 50,
      version: 1,
    };
    const client = buildClient({ "POST /api/v1/cadreen/memory": fixture });
    const intent = new IntentResource(client);
    const memory = new MemoryResource(client);
    const connections = new ConnectionsResource(client);
    const skills = new SkillsResource(intent, memory, connections);

    const result = await skills.remember({ type: "note", content: { text: "Skill note" } });
    expect(result.id).toBe("mem-skill");
    expect(result.type).toBe("note");
  });
});

describe("SkillsResource recall()", () => {
  it("delegates to memory.search", async () => {
    const fixture = {
      results: [{ id: "atom-skill", type: "note", domain: "general", authority: 50, version: 1, content: { text: "Found" } }],
      count: 1,
    };
    const client = buildClient({ "GET /api/v1/cadreen/memory/search?query=find": fixture });
    const intent = new IntentResource(client);
    const memory = new MemoryResource(client);
    const connections = new ConnectionsResource(client);
    const skills = new SkillsResource(intent, memory, connections);

    const result = await skills.recall({ query: "find" });
    expect(result.count).toBe(1);
    expect(result.results[0].id).toBe("atom-skill");
  });
});

describe("SkillsResource connectOpenAPI()", () => {
  it("delegates to connections.registerOpenAPI", async () => {
    const fixture = {
      id: "oa-skill",
      name: "Skill API",
      type: "openapi",
      tools_registered: 3,
      status: "registered",
    };
    const client = buildClient({ "POST /api/v1/cadreen/connections/openapi": fixture });
    const intent = new IntentResource(client);
    const memory = new MemoryResource(client);
    const connections = new ConnectionsResource(client);
    const skills = new SkillsResource(intent, memory, connections);

    const result = await skills.connectOpenAPI({ name: "Skill API", spec_url: "https://example.com/spec" });
    expect(result.id).toBe("oa-skill");
    expect(result.status).toBe("registered");
  });
});

describe("SkillsResource connectMCP()", () => {
  it("delegates to connections.registerMCP", async () => {
    const fixture = {
      id: "mcp-skill",
      name: "Skill MCP",
      type: "mcp",
      transport: "stdio",
      status: "registered",
    };
    const client = buildClient({ "POST /api/v1/cadreen/connections/mcp": fixture });
    const intent = new IntentResource(client);
    const memory = new MemoryResource(client);
    const connections = new ConnectionsResource(client);
    const skills = new SkillsResource(intent, memory, connections);

    const result = await skills.connectMCP({ name: "Skill MCP", url: "http://localhost:8080" });
    expect(result.id).toBe("mcp-skill");
  });
});
