import type { HttpClient } from "../client";
import type {
  ResponseRequest,
  ResponsesCompletion,
  ResponseStreamEvent,
} from "../types";

export class ResponsesResource {
  constructor(private client: HttpClient) {}

  /**
   * Create a non-streaming response.
   *
   * @example
   * ```ts
   * const response = await cadreen.responses.create({
   *   model: "gpt-4.1",
   *   input: "What is the capital of France?",
   * });
   * console.log(response.output_text);
   * ```
   */
  async create(request: ResponseRequest): Promise<ResponsesCompletion> {
    return this.client.post<ResponsesCompletion>("/api/v1/cadreen/responses", {
      ...request,
      stream: false,
    });
  }

  /**
   * Retrieve a response by ID.
   *
   * @example
   * ```ts
   * const response = await cadreen.responses.retrieve("resp_abc123");
   * console.log(response.status);
   * ```
   */
  async retrieve(responseId: string): Promise<ResponsesCompletion> {
    return this.client.get<ResponsesCompletion>(`/api/v1/cadreen/responses/${responseId}`);
  }

  /**
   * Create a streaming response.
   * Returns an async iterable of SSE events.
   *
   * @example
   * ```ts
   * const stream = await cadreen.responses.stream({
   *   model: "gpt-4.1",
   *   input: "Write a haiku.",
   * });
   * for await (const event of stream) {
   *   if (event.type === "response.output_text.delta") {
   *     process.stdout.write(event.delta || "");
   *   }
   * }
   * ```
   */
  async stream(request: ResponseRequest): Promise<AsyncIterable<ResponseStreamEvent>> {
    const response = await this.client.postStream("/api/v1/cadreen/responses", {
      ...request,
      stream: true,
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Responses stream failed: ${response.status} ${body}`);
    }

    return parseResponsesSSEStream(response);
  }
}

async function* parseResponsesSSEStream(
  response: globalThis.Response
): AsyncIterable<ResponseStreamEvent> {
  const reader = response.body?.getReader();
  if (!reader) {
    throw new Error("No response body");
  }

  const decoder = new TextDecoder();
  let buffer = "";
  let currentEvent = "message";

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() || "";

      for (const line of lines) {
        const trimmed = line.trim();
        if (trimmed.startsWith("event: ")) {
          currentEvent = trimmed.slice(7).trim();
          continue;
        }
        if (!trimmed.startsWith("data: ")) continue;

        const data = trimmed.slice(6);
        if (data === "[DONE]") {
          return;
        }

        try {
          const parsed = JSON.parse(data) as Record<string, unknown>;
          yield {
            type: currentEvent,
            sequence: parsed.sequence as number | undefined,
            response: parsed.response as ResponsesCompletion | undefined,
            item: parsed.item as ResponseStreamEvent["item"],
            output_index: parsed.output_index as number | undefined,
            content_index: parsed.content_index as number | undefined,
            delta: parsed.delta as string | undefined,
          };
          currentEvent = "message";
        } catch {
          yield { type: currentEvent, delta: data };
          currentEvent = "message";
        }
      }
    }
  } finally {
    reader.releaseLock();
  }
}
