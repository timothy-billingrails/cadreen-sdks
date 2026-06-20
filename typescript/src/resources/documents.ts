import type {
  ListDocumentsResponse,
  Document,
  UploadDocumentResponse,
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

  async upload(file: File | Blob | ReadableStream, filename?: string): Promise<UploadDocumentResponse> {
    const formData = new FormData();
    if (file instanceof ReadableStream) {
      const blob = await new Response(file).blob();
      formData.append("document", blob, filename);
    } else {
      formData.append("document", file, filename || (file instanceof File ? file.name : undefined));
    }
    return this.client.postMultipart<UploadDocumentResponse>("/api/v1/cadreen/documents/upload", formData);
  }
}
