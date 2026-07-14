import { describe, it, expect } from "vitest";
import { ConnectionsResource } from "../../resources/connections";
import { HttpClient } from "../../client";
import type {
  CatalogResponse,
  InstallResponse,
  RegisterOpenAPIResponse,
  RegisterMCPResponse,
  ListConnectionsResponse,
  ConnectResult,
} from "../../types";

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("ConnectionsResource catalog()", () => {
  it("returns fixture data", async () => {
    const fixture: CatalogResponse = {
      categories: [
        {
          name: "Productivity",
          description: "Productivity tools",
          integrations: [
            {
              id: "int-1",
              name: "Slack",
              description: "Messaging",
              category: "Productivity",
              provider: "slack",
              status: "available",
              auth_type: "oauth2",
              install_time: "2 min",
            },
          ],
        },
      ],
      installed: [],
      total_available: 42,
    };
    const client = buildClient({ "GET /api/v1/cadreen/connections/catalog": fixture });
    const resource = new ConnectionsResource(client);
    const result = await resource.catalog();
    expect(result.categories).toHaveLength(1);
    expect(result.total_available).toBe(42);
  });
});

describe("ConnectionsResource install()", () => {
  it("returns fixture data", async () => {
    const fixture: InstallResponse = {
      status: "installing",
      auth_url: "https://auth.example.com",
      provider: "slack",
      estimated_time: "2 min",
    };
    const client = buildClient({ "POST /api/v1/cadreen/connections/install": fixture });
    const resource = new ConnectionsResource(client);
    const result = await resource.install("int-1");
    expect(result.status).toBe("installing");
    expect(result.provider).toBe("slack");
  });
});

describe("ConnectionsResource registerOpenAPI()", () => {
  it("returns fixture data", async () => {
    const fixture: RegisterOpenAPIResponse = {
      id: "openapi-1",
      name: "My API",
      type: "openapi",
      tools_generated: 5,
      tools_registered: 5,
      functions: ["listUsers", "getUser"],
      status: "registered",
    };
    const client = buildClient({ "POST /api/v1/cadreen/connections/openapi": fixture });
    const resource = new ConnectionsResource(client);
    const result = await resource.registerOpenAPI({
      name: "My API",
      spec_url: "https://api.example.com/openapi.json",
    });
    expect(result.id).toBe("openapi-1");
    expect(result.tools_registered).toBe(5);
  });
});

describe("ConnectionsResource registerMCP()", () => {
  it("returns fixture data", async () => {
    const fixture: RegisterMCPResponse = {
      id: "mcp-1",
      name: "My MCP",
      type: "mcp",
      transport: "stdio",
      status: "registered",
    };
    const client = buildClient({ "POST /api/v1/cadreen/connections/mcp": fixture });
    const resource = new ConnectionsResource(client);
    const result = await resource.registerMCP({
      name: "My MCP",
      url: "http://localhost:8080",
      transport: "stdio",
    });
    expect(result.id).toBe("mcp-1");
    expect(result.transport).toBe("stdio");
  });
});

describe("ConnectionsResource list()", () => {
  it("returns fixture data", async () => {
    const fixture: ListConnectionsResponse = {
      connections: [
        {
          capability: "send_email",
          status: "healthy",
        },
      ],
      total_capabilities: 1,
    };
    const client = buildClient({ "GET /api/v1/cadreen/connections": fixture });
    const resource = new ConnectionsResource(client);
    const result = await resource.list();
    expect(result.connections).toHaveLength(1);
    expect(result.total_capabilities).toBe(1);
  });
});

describe("ConnectionsResource connect()", () => {
  it("returns fixture data", async () => {
    const fixture: ConnectResult = {
      type: "prebuilt",
      capability: "send_email",
      detail: {
        tool_id: "tool-email",
        tool_name: "Send Email",
        service_id: "gmail",
        service_name: "Gmail",
        auth_type: "oauth2",
        source: "cadreen",
      },
    };
    const client = buildClient({ "POST /api/v1/cadreen/connections": fixture });
    const resource = new ConnectionsResource(client);
    const result = await resource.connect("send_email");
    expect(result.type).toBe("prebuilt");
    expect(result.capability).toBe("send_email");
  });
});
