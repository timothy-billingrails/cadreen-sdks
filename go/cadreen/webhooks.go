package cadreen

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
)

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
	Secret string   `json:"secret,omitempty"`
}

func (c *Client) CreateWebhook(ctx context.Context, req CreateWebhookRequest, opts ...RequestOption) (*Webhook, error) {
	var result Webhook
	if err := c.do(ctx, "POST", "/api/v1/cadreen/webhooks", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}
	return &result, nil
}

func (c *Client) ListWebhooks(ctx context.Context, opts ...RequestOption) (*ListWebhooksResponse, error) {
	var result ListWebhooksResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/webhooks", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, id string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/webhooks/"+url.PathEscape(id), nil, nil, opts...); err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}
	return nil
}

// VerifyWebhookSignature verifies the HMAC-SHA256 signature of a webhook payload.
//
// Parameters:
//   - rawBody: the raw request body bytes (do NOT unmarshal JSON first)
//   - signature: the value of the X-Cadreen-Signature header
//   - secret: the secret you set when creating the webhook subscription
//
// Returns true if the signature is valid.
func VerifyWebhookSignature(rawBody []byte, signature string, secret string) bool {
	if len(rawBody) == 0 || signature == "" || secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expected))
}
