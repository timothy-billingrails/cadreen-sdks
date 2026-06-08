import type { ExecutionEvent, ExecutionStatus } from "../types";
import { HttpClient } from "../client";

export class ExecutionsResource {
  constructor(private client: HttpClient) {}

  async *stream(executionId: string): AsyncGenerator<ExecutionEvent> {
    const path = `/api/v1/cadreen/executions/${encodeURIComponent(executionId)}/stream`;
    for await (const event of this.client.stream(path)) {
      yield {
        type: event.type,
        data: event.data,
      };
    }
  }

  async getStatus(executionId: string): Promise<ExecutionStatus> {
    return this.client.get<ExecutionStatus>(
      `/api/v1/cadreen/executions/${encodeURIComponent(executionId)}`
    );
  }
}
