import type {
  IntentResult,
  IntentContext,
  RememberRequest,
  CreateMemoryResponse,
  SearchMemoryRequest,
  SearchMemoryResponse,
  RegisterOpenAPIRequest,
  RegisterOpenAPIResponse,
  RegisterMCPRequest,
  RegisterMCPResponse,
  InstallComposioRequest,
} from "../types";
import { IntentResource } from "./intent";
import { MemoryResource } from "./memory";
import { ConnectionsResource } from "./connections";

/**
 * SkillsResource — Builder-shaped naming facade for the system's capabilities.
 *
 * Instead of remembering that `intent` is the chat endpoint and `memory`
 * is the knowledge store, developers work with `cadreen.skills.ask()`,
 * `cadreen.skills.remember()`, and `cadreen.skills.connect()`.
 *
 * This is a thin wrapper; no new backend routes are introduced.
 */
export class SkillsResource {
  constructor(
    private intent: IntentResource,
    private memory: MemoryResource,
    private connections: ConnectionsResource
  ) {}

  /** Ask a question and get a direct answer or clarifying questions. */
  async ask(
    prompt: string,
    options?: { conversation_id?: string; context?: IntentContext; stream?: boolean }
  ): Promise<IntentResult> {
    return this.intent.invoke({
      messages: [{ role: "user", content: prompt }],
      mode: "chat",
      conversation_id: options?.conversation_id,
      context: options?.context,
      stream: options?.stream,
    });
  }

  /** Execute an action (mission-bound intent). */
  async act(
    prompt: string,
    options?: { conversation_id?: string; context?: IntentContext; stream?: boolean }
  ): Promise<IntentResult> {
    return this.intent.invoke({
      messages: [{ role: "user", content: prompt }],
      mode: "execution",
      conversation_id: options?.conversation_id,
      context: options?.context,
      stream: options?.stream,
    });
  }

  /** Store a memory atom. */
  async remember(request: RememberRequest): Promise<CreateMemoryResponse> {
    return this.memory.remember(request);
  }

  /** Search stored memories. */
  async recall(request: SearchMemoryRequest): Promise<SearchMemoryResponse> {
    return this.memory.search(request);
  }

  /** Register an OpenAPI connector. */
  async connectOpenAPI(request: RegisterOpenAPIRequest): Promise<RegisterOpenAPIResponse> {
    return this.connections.registerOpenAPI(request);
  }

  /** Register an MCP connector. */
  async connectMCP(request: RegisterMCPRequest): Promise<RegisterMCPResponse> {
    return this.connections.registerMCP(request);
  }

  /** Install a Composio integration. */
  async connectComposio(request: InstallComposioRequest): Promise<Record<string, unknown>> {
    return this.connections.installComposio(request);
  }
}
