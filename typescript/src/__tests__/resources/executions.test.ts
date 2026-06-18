import { describe, it, expect } from "vitest";
import { ExecutionsResource } from "../../resources/executions";
import { HttpClient } from "../../client";
import type { ExecutionStatus } from "../../types";

function buildClient(fixtures: Record<string, unknown>) {
  return new HttpClient({
    apiKey: "sk_test",
    sandbox: true,
    fixtures,
  });
}

describe("ExecutionsResource getStatus()", () => {
  it("returns fixture data", async () => {
    const fixture: ExecutionStatus = {
      id: "exec-1",
      status: "completed",
      progress: 100,
      result: { output: "Done" },
    };
    const client = buildClient({ "GET /api/v1/cadreen/executions/exec-1": fixture });
    const resource = new ExecutionsResource(client);
    const result = await resource.getStatus("exec-1");
    expect(result.id).toBe("exec-1");
    expect(result.status).toBe("completed");
    expect(result.progress).toBe(100);
  });

  it("returns running status", async () => {
    const fixture: ExecutionStatus = {
      id: "exec-2",
      status: "running",
      progress: 45,
    };
    const client = buildClient({ "GET /api/v1/cadreen/executions/exec-2": fixture });
    const resource = new ExecutionsResource(client);
    const result = await resource.getStatus("exec-2");
    expect(result.status).toBe("running");
    expect(result.progress).toBe(45);
  });

  it("returns failed status with error", async () => {
    const fixture: ExecutionStatus = {
      id: "exec-3",
      status: "failed",
      error: "Tool execution timed out",
    };
    const client = buildClient({ "GET /api/v1/cadreen/executions/exec-3": fixture });
    const resource = new ExecutionsResource(client);
    const result = await resource.getStatus("exec-3");
    expect(result.status).toBe("failed");
    expect(result.error).toBe("Tool execution timed out");
  });
});
