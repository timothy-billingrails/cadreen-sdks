import type { HttpClient } from "../client";

// ── Chat Completions Types (OpenAI-compatible) ──

export interface ChatMessage {
  role: "system" | "user" | "assistant" | "tool";
  content?: string;
  name?: string;
  tool_call_id?: string; // for "tool" role
  tool_calls?: ChatToolCall[]; // for "assistant" role
}

export interface ChatToolCall {
  id: string;
  type: "function";
  function: ChatFunctionCall;
}

export interface ChatFunctionCall {
  name: string;
  arguments: string; // JSON string
}

export interface ChatToolDefinition {
  type: "function";
  function: ChatFunctionDefinition;
}

export interface ChatFunctionDefinition {
  name: string;
  description?: string;
  parameters?: unknown; // JSON Schema
}

export interface ChatCompletionRequest {
  model?: string;
  messages: ChatMessage[];
  stream?: boolean;
  tools?: ChatToolDefinition[];
  context?: Record<string, unknown>;
  conversation_id?: string;
  user_id?: string;
}

export interface ChatCompletionResponse {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: ChatChoice[];
  usage?: ChatUsage;
}

export interface ChatChoice {
  index: number;
  message: ChatMessage;
  finish_reason: string;
}

export interface ChatUsage {
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
}

export interface ChatCompletionChunk {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: ChatChunkChoice[];
  usage?: ChatUsage;
}

export interface ChatChunkChoice {
  index: number;
  delta: ChatDelta;
  finish_reason: string | null;
}

export interface ChatDelta {
  role?: string;
  content?: string;
  tool_calls?: ChatToolCall[];
}

// ── Tool Discovery Types ──

export interface ToolEntry {
  type: string;
  function: ChatFunctionDefinition;
}

export interface ListToolsResponse {
  object: string;
  data: ToolEntry[];
}

// ── Chat Resource ──

export class ChatResource {
  constructor(private client: HttpClient) {}

  /**
   * Send a chat completion request (non-streaming).
   *
   * @example
   * ```ts
   * const response = await cadreen.chat.completions({
   *   messages: [{ role: "user", content: "Hello!" }],
   * });
   * console.log(response.choices[0].message.content);
   * ```
   */
  async completions(request: ChatCompletionRequest): Promise<ChatCompletionResponse> {
    return this.client.post<ChatCompletionResponse>("/api/v1/cadreen/chat/completions", {
      ...request,
      stream: false,
    });
  }

  /**
   * Send a streaming chat completion request.
   * Returns an async iterable of chunks.
   *
   * @example
   * ```ts
   * const stream = await cadreen.chat.completionsStream({
   *   messages: [{ role: "user", content: "Hello!" }],
   * });
   * for await (const event of stream) {
   *   if (event.type === "chunk") {
   *     process.stdout.write(event.chunk.choices[0]?.delta?.content || "");
   *   }
   * }
   * ```
   */
  async completionsStream(
    request: ChatCompletionRequest
  ): Promise<AsyncIterable<ChatStreamEvent>> {
    const response = await this.client.postStream("/api/v1/cadreen/chat/completions", {
      ...request,
      stream: true,
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Chat stream failed: ${response.status} ${body}`);
    }

    return parseChatSSEStream(response);
  }

  /**
   * List available tools as OpenAI-compatible function definitions.
   */
  async listTools(): Promise<ListToolsResponse> {
    return this.client.get<ListToolsResponse>("/api/v1/cadreen/tools");
  }
}

export type ChatStreamEvent =
  | { type: "chunk"; chunk: ChatCompletionChunk }
  | { type: "done" }
  | { type: "error"; error: Error };

async function* parseChatSSEStream(
  response: Response
): AsyncIterable<ChatStreamEvent> {
  const reader = response.body?.getReader();
  if (!reader) {
    yield { type: "error", error: new Error("No response body") };
    return;
  }

  const decoder = new TextDecoder();
  let buffer = "";

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() || "";

      for (const line of lines) {
        const trimmed = line.trim();
        if (!trimmed.startsWith("data: ")) continue;

        const data = trimmed.slice(6);
        if (data === "[DONE]") {
          yield { type: "done" };
          return;
        }

        try {
          const chunk = JSON.parse(data) as ChatCompletionChunk;
          yield { type: "chunk", chunk };
        } catch (e) {
          yield { type: "error", error: new Error(`Parse error: ${e}`) };
        }
      }
    }
  } finally {
    reader.releaseLock();
  }
}
