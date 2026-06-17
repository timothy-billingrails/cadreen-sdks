import type { CadreenConfig, RequestOptions } from "./types";
import type { TelemetryProvider } from "./telemetry";
import { NoOpProvider, wrapWithTelemetry } from "./telemetry";

type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

interface ApiErrorResponse {
  error: {
    type: string;
    code: string;
    message: string;
    details?: Array<{ field: string; message: string }>;
  };
  intelligence?: unknown;
}

export class CadreenError extends Error {
  public readonly status: number;
  public readonly code: string;
  public readonly errorType: string;
  public readonly details?: Array<{ field: string; message: string }>;
  public readonly intelligence?: unknown;

  constructor(
    status: number,
    code: string,
    errorType: string,
    message: string,
    details?: Array<{ field: string; message: string }>,
    intelligence?: unknown
  ) {
    super(message);
    this.name = "CadreenError";
    this.status = status;
    this.code = code;
    this.errorType = errorType;
    this.details = details;
    this.intelligence = intelligence;
  }
}

export class CadreenBlockedError extends CadreenError {
  public readonly reason_code?: string;
  public readonly policy_id?: string;
  public readonly traceId: string;
  public readonly intelligence: unknown;

  constructor(params: {
    reason_code?: string;
    policy_id?: string;
    intelligence: unknown;
    traceId: string;
  }) {
    super(
      403,
      params.reason_code || "blocked_by_policy",
      "blocked",
      `Action blocked by governance policy${params.policy_id ? `: ${params.policy_id}` : ""}`
    );
    this.name = "CadreenBlockedError";
    this.reason_code = params.reason_code;
    this.policy_id = params.policy_id;
    this.intelligence = params.intelligence;
    this.traceId = params.traceId;
  }
}

export class CadreenClarifyError extends CadreenError {
  public readonly questions: Array<{ id: string; question: string; type: string; required: boolean }>;
  public readonly conversationId: string;
  public readonly traceId: string;
  public readonly intelligence: unknown;

  constructor(params: {
    questions: Array<{ id: string; question: string; type: string; required: boolean }>;
    conversationId: string;
    intelligence: unknown;
    traceId: string;
  }) {
    super(422, "needs_input", "clarify", "System needs clarification before proceeding");
    this.name = "CadreenClarifyError";
    this.questions = params.questions;
    this.conversationId = params.conversationId;
    this.intelligence = params.intelligence;
    this.traceId = params.traceId;
  }
}

const DEFAULT_BASE_URL = "https://accomplishanything.today";
const DEFAULT_MAX_RETRIES = 2;
const DEFAULT_TIMEOUT = 30000;
const RETRYABLE_STATUS_CODES = new Set([408, 429, 502, 503, 504]);
const IDEMPOTENT_METHODS = new Set<HttpMethod>(["GET", "PUT"]);

function generateIdempotencyKey(): string {
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID();
  }
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

function buildQueryString(params: Record<string, string | number | boolean | undefined>): string {
  const parts: string[] = [];
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") {
      parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`);
    }
  }
  return parts.length > 0 ? `?${parts.join("&")}` : "";
}

function parseSSELine(line: string): { event?: string; data?: string } {
  if (line.startsWith("event:")) {
    return { event: line.slice(6).trim() };
  }
  if (line.startsWith("data:")) {
    return { data: line.slice(5).trim() };
  }
  return {};
}

export class HttpClient {
  private readonly baseUrl: string;
  private readonly apiKey: string;
  private readonly maxRetries: number;
  private readonly timeout: number;
  private readonly telemetry: ReturnType<typeof wrapWithTelemetry>;
  private readonly sandbox: boolean;
  private readonly fixtures: Record<string, unknown>;
  private readonly profile: string;

  constructor(config: CadreenConfig) {
    this.baseUrl = (config.baseUrl || DEFAULT_BASE_URL).replace(/\/+$/, "");
    this.apiKey = config.apiKey;
    this.maxRetries = config.maxRetries ?? DEFAULT_MAX_RETRIES;
    this.timeout = config.timeout ?? DEFAULT_TIMEOUT;
    this.sandbox = config.sandbox ?? false;
    this.fixtures = config.fixtures ?? {};
    this.profile = config.profile ?? "full";
    const provider: TelemetryProvider = (config.telemetry as TelemetryProvider) ?? new NoOpProvider();
    this.telemetry = wrapWithTelemetry(provider);
  }

  private async request<T>(
    method: HttpMethod,
    path: string,
    body?: unknown,
    options?: RequestOptions
  ): Promise<T> {
    if (this.sandbox) {
      const fixtureKey = `${method} ${path}`;
      if (fixtureKey in this.fixtures) {
        return this.fixtures[fixtureKey] as T;
      }
      if (path in this.fixtures) {
        return this.fixtures[path] as T;
      }
      throw new CadreenError(404, "not_found", "not_found", `No fixture for ${fixtureKey}. Provide fixtures via CadreenConfig.fixtures keyed by "METHOD /path" or "/path".`);
    }

    const url = `${this.baseUrl}${path}`;
    const span = this.telemetry.onRequestStart(method, path);
    const startTime = Date.now();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      Authorization: `Bearer ${this.apiKey}`,
      Accept: `application/json; profile="${this.profile}"`,
      ...(options?.headers || {}),
    };

    if (method === "POST" || method === "PUT" || method === "PATCH") {
      headers["Idempotency-Key"] = options?.idempotencyKey || generateIdempotencyKey();
    }

    const isIdempotent = IDEMPOTENT_METHODS.has(method) || !!headers["Idempotency-Key"];
    const maxAttempts = isIdempotent ? this.maxRetries + 1 : 1;

    let lastError: CadreenError | null = null;

    for (let attempt = 0; attempt < maxAttempts; attempt++) {
      if (attempt > 0) {
        const delay = Math.min(1000 * Math.pow(2, attempt - 1), 10000);
        this.telemetry.onRetry(method, path, attempt);
        await new Promise((resolve) => setTimeout(resolve, delay));
      }

      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), this.timeout);

        const response = await fetch(url, {
          method,
          headers,
          body: body !== undefined ? JSON.stringify(body) : undefined,
          signal: controller.signal,
        });

        clearTimeout(timeoutId);

        if (!response.ok) {
          let errorBody: ApiErrorResponse | null = null;
          try {
            errorBody = await response.json();
          } catch {
            // not JSON
          }

          const err = new CadreenError(
            response.status,
            errorBody?.error?.code || "unknown",
            errorBody?.error?.type || "error",
            errorBody?.error?.message || response.statusText,
            errorBody?.error?.details,
            errorBody?.intelligence
          );

          if (RETRYABLE_STATUS_CODES.has(response.status) && isIdempotent && attempt < maxAttempts - 1) {
            lastError = err;
            continue;
          }

          this.telemetry.onError(span, err);
          throw err;
        }

        if (response.status === 204) {
          this.telemetry.onRequestEnd(span, method, path, 204, Date.now() - startTime);
          return undefined as T;
        }

        const result = (await response.json()) as T;
        this.telemetry.onRequestEnd(span, method, path, response.status, Date.now() - startTime);
        return result;
      } catch (error) {
        if (error instanceof CadreenError) {
          throw error;
        }
        if (error instanceof DOMException && error.name === "AbortError") {
          throw new CadreenError(408, "timeout", "timeout", "Request timed out");
        }
        if (isIdempotent && attempt < maxAttempts - 1) {
          lastError = new CadreenError(0, "network_error", "network", (error as Error).message);
          continue;
        }
        throw new CadreenError(0, "network_error", "network", (error as Error).message);
      }
    }

    throw lastError || new CadreenError(0, "network_error", "network", "Request failed after retries");
  }

  async get<T>(path: string, params?: Record<string, string | number | boolean | undefined>): Promise<T> {
    const qs = params ? buildQueryString(params) : "";
    return this.request<T>("GET", `${path}${qs}`);
  }

  async post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>("POST", path, body, options);
  }

  async put<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>("PUT", path, body, options);
  }

  async delete<T>(path: string): Promise<T> {
    return this.request<T>("DELETE", path);
  }

  async *stream(path: string): AsyncGenerator<{ type: string; data: Record<string, unknown> }> {
    const url = `${this.baseUrl}${path}`;
    const headers: Record<string, string> = {
      Authorization: `Bearer ${this.apiKey}`,
      Accept: "text/event-stream",
    };

    const response = await fetch(url, {
      method: "GET",
      headers,
    });

    if (!response.ok) {
      let errorBody: ApiErrorResponse | null = null;
      try {
        errorBody = await response.json();
      } catch {
        // not JSON
      }
      throw new CadreenError(
        response.status,
        errorBody?.error?.code || "unknown",
        errorBody?.error?.type || "error",
        errorBody?.error?.message || response.statusText,
        errorBody?.error?.details,
        errorBody?.intelligence
      );
    }

    const reader = response.body?.getReader();
    if (!reader) {
      throw new CadreenError(0, "stream_error", "network", "Response body is not readable");
    }

    const decoder = new TextDecoder();
    let buffer = "";
    let currentEvent = "message";

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() || "";

      for (const line of lines) {
        const parsed = parseSSELine(line);
        if (parsed.event) {
          currentEvent = parsed.event;
        }
        if (parsed.data) {
          try {
            const data = JSON.parse(parsed.data) as Record<string, unknown>;
            yield { type: currentEvent, data };
          } catch {
            yield { type: currentEvent, data: { raw: parsed.data } };
          }
          currentEvent = "message";
        }
      }
    }
  }
}
