package cadreen

import (
	"context"
	"fmt"
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
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/webhooks/"+id, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}
	return nil
}
