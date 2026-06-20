import type {
  ListDocumentsResponse,
  Document,
} from "../types";
import { HttpClient } from "../client";

export class DocumentsResource {
  constructor(private client: HttpClient) {}

  async list(): Promise<ListDocumentsResponse> {
    return this.client.get<ListDocumentsResponse>("/api/v1/cadreen/documents");
  }

  async get(id: string): Promise<Document> {
    return this.client.get<Document>(`/api/v1/cadreen/documents/${encodeURIComponent(id)}`);
  }

  async download(id: string): Promise<Response> {
    const url = `${(this.client as any).baseUrl}/api/v1/cadreen/documents/${encodeURIComponent(id)}/download`;
    const response = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${(this.client as any).apiKey}`,
      },
    });
    return response;
  }
}
