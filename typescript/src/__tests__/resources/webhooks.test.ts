import { describe, it, expect } from "vitest";
import { WebhooksResource } from "../../resources/webhooks";
import { createHmac } from "crypto";

describe("WebhooksResource verifySignature()", () => {
  const resource = new WebhooksResource();
  const secret = "whsec_test_secret_12345";

  it("returns true for valid signature", () => {
    const body = '{"event":"intent.completed","data":{"id":"123"}}';
    const signature = createHmac("sha256", secret).update(body, "utf-8").digest("hex");
    expect(resource.verifySignature(body, signature, secret)).toBe(true);
  });

  it("returns false for invalid signature", () => {
    const body = '{"event":"intent.completed","data":{"id":"123"}}';
    const wrongSignature = createHmac("sha256", "wrong_secret").update(body, "utf-8").digest("hex");
    expect(resource.verifySignature(body, wrongSignature, secret)).toBe(false);
  });

  it("returns false for tampered body", () => {
    const originalBody = '{"event":"intent.completed","data":{"id":"123"}}';
    const signature = createHmac("sha256", secret).update(originalBody, "utf-8").digest("hex");
    const tamperedBody = '{"event":"intent.completed","data":{"id":"456"}}';
    expect(resource.verifySignature(tamperedBody, signature, secret)).toBe(false);
  });

  it("returns false for empty body", () => {
    const body = '{"event":"test"}';
    const signature = createHmac("sha256", secret).update(body, "utf-8").digest("hex");
    expect(resource.verifySignature("", signature, secret)).toBe(false);
  });

  it("returns false for empty signature", () => {
    const body = '{"event":"test"}';
    expect(resource.verifySignature(body, "", secret)).toBe(false);
  });

  it("returns false for empty secret", () => {
    const body = '{"event":"test"}';
    const signature = createHmac("sha256", secret).update(body, "utf-8").digest("hex");
    expect(resource.verifySignature(body, signature, "")).toBe(false);
  });

  it("handles unicode body", () => {
    const body = '{"event":"intent.completed","message":"café résumé"}';
    const signature = createHmac("sha256", secret).update(body, "utf-8").digest("hex");
    expect(resource.verifySignature(body, signature, secret)).toBe(true);
  });
});
