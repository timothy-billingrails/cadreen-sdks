package cadreen

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type RememberRequest struct {
	Type      string                 `json:"type"`
	Content   map[string]any `json:"content"`
	Domain    string                 `json:"domain,omitempty"`
	Scope     string                 `json:"scope,omitempty"`
	Authority int                    `json:"authority,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
}

type SearchMemoryOptions struct {
	Domain string
	Tag    string
	Limit  int
}

func (c *Client) Teach(ctx context.Context, req RememberRequest, opts ...RequestOption) (*CreateMemoryResponse, error) {
	var result CreateMemoryResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/memory", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("teach: %w", err)
	}
	return &result, nil
}

func (c *Client) Search(ctx context.Context, query string, options *SearchMemoryOptions, opts ...RequestOption) (*SearchMemoryResponse, error) {
	u, err := url.Parse("/api/v1/cadreen/memory/search")
	if err != nil {
		return nil, fmt.Errorf("search: parse url: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	if options != nil {
		if options.Domain != "" {
			q.Set("domain", options.Domain)
		}
		if options.Tag != "" {
			q.Set("tag", options.Tag)
		}
		if options.Limit > 0 {
			q.Set("limit", strconv.Itoa(options.Limit))
		}
	}
	u.RawQuery = q.Encode()

	var result SearchMemoryResponse
	if err := c.do(ctx, "GET", u.String(), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	return &result, nil
}

func (c *Client) GetAtom(ctx context.Context, id string, opts ...RequestOption) (*Atom, error) {
	var result Atom
	if err := c.do(ctx, "GET", "/api/v1/cadreen/memory/"+url.PathEscape(id), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get atom: %w", err)
	}
	return &result, nil
}

func (c *Client) Profile(ctx context.Context, userID string, opts ...RequestOption) (*MemoryProfileResponse, error) {
	var result MemoryProfileResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/memory/profile/"+url.PathEscape(userID), nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("profile: %w", err)
	}
	return &result, nil
}

func (c *Client) MemoryTypes(ctx context.Context, opts ...RequestOption) (*MemoryTypesResponse, error) {
	var result MemoryTypesResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/memory/types", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("memory types: %w", err)
	}
	return &result, nil
}
