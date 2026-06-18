import { describe, it, expect } from "vitest";
import { MemoryResource } from "../../resources/memory";
import { HttpClient } from "../../client";
import type { CreateMemoryResponse, SearchMemoryResponse } from "../../types";

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("MemoryResource remember()", () => {
  it("returns fixture data", async () => {
    const fixture: CreateMemoryResponse = {
      id: "mem-1",
      type: "note",
      domain: "general",
      content: { text: "Remember this" },
      authority: 50,
      version: 1,
    };
    const client = buildClient({ "POST /api/v1/cadreen/memory": fixture });
    const resource = new MemoryResource(client);
    const result = await resource.remember({
      type: "note",
      content: { text: "Remember this" },
    });
    expect(result.id).toBe("mem-1");
    expect(result.type).toBe("note");
  });
});

describe("MemoryResource search()", () => {
  it("returns fixture data for search", async () => {
    const fixture: SearchMemoryResponse = {
      results: [
        {
          id: "atom-1",
          type: "note",
          domain: "general",
          authority: 50,
          version: 1,
          content: { text: "Found note" },
        },
      ],
      count: 1,
    };
    const client = buildClient({
      "GET /api/v1/cadreen/memory/search?query=test": fixture,
    });
    const resource = new MemoryResource(client);
    const result = await resource.search({ query: "test" });
    expect(result.count).toBe(1);
    expect(result.results[0].id).toBe("atom-1");
  });

  it("builds query string with domain, tag, limit", async () => {
    const fixture: SearchMemoryResponse = {
      results: [],
      count: 0,
    };
    const client = buildClient({
      "GET /api/v1/cadreen/memory/search?query=test&domain=support&tag=urgent&limit=10": fixture,
    });
    const resource = new MemoryResource(client);
    const result = await resource.search({
      query: "test",
      domain: "support",
      tag: "urgent",
      limit: 10,
    });
    expect(result.count).toBe(0);
  });
});
