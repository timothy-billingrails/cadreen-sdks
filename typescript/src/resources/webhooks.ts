import { timingSafeEqual } from "crypto";
import { createHmac } from "crypto";

export class WebhooksResource {
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
