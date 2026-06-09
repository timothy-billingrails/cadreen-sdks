package cadreen

import (
	"context"
	"fmt"
)

func (c *Client) Setup(ctx context.Context, req SetupRequest, opts ...RequestOption) (*SetupResult, error) {
	var result SetupResult
	if err := c.do(ctx, "POST", "/api/v1/cadreen/setup", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}
	return &result, nil
}
