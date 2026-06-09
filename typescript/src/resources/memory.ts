import type {
  RememberRequest,
  CreateMemoryResponse,
  SearchMemoryRequest,
  SearchMemoryResponse,
  MemoryTypesResponse,
  Atom,
} from "../types";
import { HttpClient } from "../client";

export class MemoryResource {
  constructor(private client: HttpClient) {}

  async remember(request: RememberRequest): Promise<CreateMemoryResponse> {
    return this.client.post<CreateMemoryResponse>("/api/v1/cadreen/memory", {
      type: request.type,
      content: request.content,
      domain: request.domain,
      scope: request.scope,
      authority: request.authority,
      tags: request.tags,
    });
  }

  async search(request: SearchMemoryRequest): Promise<SearchMemoryResponse> {
    const params: Record<string, string | number | undefined> = {
      query: request.query,
    };
    if (request.domain) params.domain = request.domain;
    if (request.tag) params.tag = request.tag;
    if (request.limit) params.limit = request.limit;
    return this.client.get<SearchMemoryResponse>("/api/v1/cadreen/memory/search", params);
  }

  async get(id: string): Promise<Atom> {
    return this.client.get<Atom>(`/api/v1/cadreen/memory/${encodeURIComponent(id)}`);
  }

  async types(): Promise<MemoryTypesResponse> {
    return this.client.get<MemoryTypesResponse>("/api/v1/cadreen/memory/types");
  }
}
