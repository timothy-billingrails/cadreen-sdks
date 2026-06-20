package cadreen

import (
	"context"
	"fmt"
)

func (c *Client) ListDocuments(ctx context.Context, opts ...RequestOption) (*ListDocumentsResponse, error) {
	var result ListDocumentsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/documents", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return &result, nil
}

func (c *Client) GetDocument(ctx context.Context, id string, opts ...RequestOption) (*Document, error) {
	var result Document
	if err := c.do(ctx, "GET", "/api/v1/cadreen/documents/"+id, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &result, nil
}
