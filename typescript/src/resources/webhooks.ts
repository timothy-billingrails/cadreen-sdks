import { timingSafeEqual } from "crypto";
import { createHmac } from "crypto";
import type {
  Webhook,
  ListWebhooksResponse,
} from "../types";
import { HttpClient } from "../client";

export class WebhooksResource {
  constructor(private client?: HttpClient) {}

  async create(request: { url: string; events?: string[]; secret?: string }): Promise<Webhook> {
    if (!this.client) throw new Error("WebhooksResource not initialized with HttpClient");
    return this.client.post<Webhook>("/api/v1/cadreen/webhooks", request);
  }

  async list(): Promise<ListWebhooksResponse> {
    if (!this.client) throw new Error("WebhooksResource not initialized with HttpClient");
    return this.client.get<ListWebhooksResponse>("/api/v1/cadreen/webhooks");
  }

  async delete(id: string): Promise<void> {
    if (!this.client) throw new Error("WebhooksResource not initialized with HttpClient");
    return this.client.delete<void>(`/api/v1/cadreen/webhooks/${encodeURIComponent(id)}`);
  }

  /**
   * Verify the HMAC-SHA256 signature of a webhook payload.
   *
   * @param rawBody - The raw request body string (do NOT JSON.parse first)
   * @param signature - The value of the X-Cadreen-Signature header
   * @param secret - The secret you set when creating the webhook subscription
   * @returns true if the signature is valid, false otherwise
   */
  verifySignature(rawBody: string, signature: string, secret: string): boolean {
    if (!rawBody || !signature || !secret) {
      return false;
    }

    try {
      const expected = createHmac("sha256", secret).update(rawBody, "utf-8").digest("hex");
      return timingSafeEqual(Buffer.from(signature), Buffer.from(expected));
    } catch {
      return false;
    }
  }
}
