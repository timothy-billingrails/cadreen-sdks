package cadreen

import (
	"context"
	"fmt"
)

func (c *Client) ListLearningPatterns(ctx context.Context, opts ...RequestOption) (*ListLearningPatternsResponse, error) {
	var result ListLearningPatternsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/learning/patterns", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list learning patterns: %w", err)
	}
	return &result, nil
}

func (c *Client) ListLearningEpisodes(ctx context.Context, opts ...RequestOption) (*ListLearningEpisodesResponse, error) {
	var result ListLearningEpisodesResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/learning/episodes", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list learning episodes: %w", err)
	}
	return &result, nil
}

func (c *Client) ListLearningSuggestions(ctx context.Context, opts ...RequestOption) (*ListLearningSuggestionsResponse, error) {
	var result ListLearningSuggestionsResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/learning/suggestions", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list learning suggestions: %w", err)
	}
	return &result, nil
}
