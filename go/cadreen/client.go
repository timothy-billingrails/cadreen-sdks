package cadreen

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	DefaultBaseURL     = "https://accomplishanything.today"
	DefaultMaxRetries  = 3
	DefaultTimeout     = 30 * time.Second
	DefaultProfile     = "full"
)

var retryableStatusCodes = map[int]bool{
	408: true,
	429: true,
	502: true,
	503: true,
	504: true,
}

var idempotentMethods = map[string]bool{
	"GET": true,
	"PUT": true,
}

func generateIdempotencyKey() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "idem-" + fmt.Sprintf("%d", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type RequestOption func(*requestOptions)

type requestOptions struct {
	idempotencyKey string
}

func WithIdempotencyKey(key string) RequestOption {
	return func(o *requestOptions) {
		o.idempotencyKey = key
	}
}

type CadreenConfig struct {
	APIKey     string
	BaseURL    string
	MaxRetries int
	Timeout    time.Duration
	Profile    string
	Sandbox    bool
	Fixtures   map[string]any
}

type Client struct {
	httpClient   *http.Client
	config       CadreenConfig
	lastResponse *http.Response
	lastTraceID  string
	mu           sync.Mutex
}

func NewClient(config CadreenConfig) *Client {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	maxRetries := config.MaxRetries
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	profile := config.Profile
	if profile == "" {
		profile = DefaultProfile
	}

	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		config: CadreenConfig{
			APIKey:     config.APIKey,
			BaseURL:    baseURL,
			MaxRetries: maxRetries,
			Timeout:    timeout,
			Profile:    profile,
			Sandbox:    config.Sandbox,
			Fixtures:   config.Fixtures,
		},
	}
}

func (c *Client) LastTraceID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lastTraceID
}

func (c *Client) setLastTrace(resp *http.Response) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastResponse = resp
	c.lastTraceID = resp.Header.Get("X-Cadreen-Trace-ID")
}

func (c *Client) do(ctx context.Context, method, path string, body, result any, opts ...RequestOption) error {
	options := &requestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if c.config.Sandbox {
		fixtureKey := method + " " + path
		if v, ok := c.config.Fixtures[fixtureKey]; ok {
			return populateResult(v, result)
		}
		if v, ok := c.config.Fixtures[path]; ok {
			return populateResult(v, result)
		}
		return &APIError{
			StatusCode: 404,
			Type:       "not_found",
			Code:       "not_found",
			Message:    fmt.Sprintf("No fixture for %s. Provide fixtures via CadreenConfig.Fixtures keyed by 'METHOD /path' or '/path'.", fixtureKey),
		}
	}

	reqURL := c.config.BaseURL + path

	var bodyBytes []byte
	var bodyReader io.Reader
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Accept", fmt.Sprintf(`application/json; profile="%s"`, c.config.Profile))

	if method == "POST" || method == "PUT" || method == "PATCH" {
		idempotencyKey := options.idempotencyKey
		if idempotencyKey == "" {
			idempotencyKey = generateIdempotencyKey()
		}
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	isIdempotent := idempotentMethods[method] || req.Header.Get("Idempotency-Key") != ""
	maxAttempts := 1
	if isIdempotent {
		maxAttempts = c.config.MaxRetries + 1
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			backoff := 1 << uint(attempt-1)
			if backoff > 8 {
				backoff = 8
			}
			delay := time.Duration(backoff) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if isIdempotent && attempt < maxAttempts-1 {
				lastErr = &APIError{Code: "network_error", Type: "network", Message: err.Error()}
				freshReq, newReqErr := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewReader(bodyBytes))
				if newReqErr != nil {
					return lastErr
				}
				freshReq.Header = req.Header.Clone()
				req = freshReq
				continue
			}
			return &APIError{Code: "network_error", Type: "network", Message: err.Error()}
		}

		c.setLastTrace(resp)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if resp.StatusCode == http.StatusNoContent {
				resp.Body.Close()
				return nil
			}
			defer resp.Body.Close()
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read response body: %w", err)
			}
			if result != nil {
				if err := json.Unmarshal(respBytes, result); err != nil {
					return fmt.Errorf("unmarshal response: %w", err)
				}
			}
			return nil
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		apiErr := parseError(resp.StatusCode, respBody)

		if retryableStatusCodes[resp.StatusCode] && isIdempotent && attempt < maxAttempts-1 {
			lastErr = apiErr
			freshReq, newReqErr := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewReader(bodyBytes))
			if newReqErr != nil {
				return lastErr
			}
			freshReq.Header = req.Header.Clone()
			req = freshReq
			continue
		}

		return apiErr
	}

	if lastErr != nil {
		return lastErr
	}
	return &APIError{Code: "network_error", Type: "network", Message: "request failed after retries"}
}

func parseError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		Type:       "error",
		Code:       "unknown",
		Message:    http.StatusText(statusCode),
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return apiErr
	}

	apiErr.Type = errResp.Error.Type
	apiErr.Code = errResp.Error.Code
	apiErr.Message = errResp.Error.Message
	apiErr.Hint = errResp.Error.Hint
	apiErr.NextAction = errResp.Error.NextAction
	apiErr.Details = errResp.Error.Details
	apiErr.RetryAfter = errResp.Error.RetryAfter

	return apiErr
}

func populateResult(src any, dst any) error {
	if dst == nil {
		return nil
	}
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal fixture: %w", err)
	}
	return json.Unmarshal(data, dst)
}
