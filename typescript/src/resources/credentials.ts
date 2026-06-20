import type {
  ListCredentialsResponse,
  CredentialMetadata,
  CreateCredentialRequest,
} from "../types";
import { HttpClient } from "../client";

export class CredentialsResource {
  constructor(private client: HttpClient) {}

  async list(): Promise<ListCredentialsResponse> {
    return this.client.get<ListCredentialsResponse>("/api/v1/cadreen/credentials");
  }

  async create(request: CreateCredentialRequest): Promise<CredentialMetadata> {
    return this.client.post<CredentialMetadata>("/api/v1/cadreen/credentials", request);
  }

  async delete(id: string): Promise<void> {
    return this.client.delete<void>(`/api/v1/cadreen/credentials/${encodeURIComponent(id)}`);
  }
}
