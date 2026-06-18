package cadreen

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name           string
		config         CadreenConfig
		wantBaseURL    string
		wantMaxRetries int
		wantTimeout    time.Duration
		wantProfile    string
		wantSandbox    bool
	}{
		{
			name:           "default values",
			config:         CadreenConfig{APIKey: "test-key"},
			wantBaseURL:    DefaultBaseURL,
			wantMaxRetries: DefaultMaxRetries,
			wantTimeout:    DefaultTimeout,
			wantProfile:    DefaultProfile,
			wantSandbox:    false,
		},
		{
			name: "custom base URL",
			config: CadreenConfig{
				APIKey:  "test-key",
				BaseURL: "https://custom.example.com",
			},
			wantBaseURL:    "https://custom.example.com",
			wantMaxRetries: DefaultMaxRetries,
			wantTimeout:    DefaultTimeout,
			wantProfile:    DefaultProfile,
			wantSandbox:    false,
		},
		{
			name: "sandbox mode activated",
			config: CadreenConfig{
				APIKey:  "sandbox_key",
				Sandbox: true,
			},
			wantBaseURL:    DefaultBaseURL,
			wantMaxRetries: DefaultMaxRetries,
			wantTimeout:    DefaultTimeout,
			wantProfile:    DefaultProfile,
			wantSandbox:    true,
		},
		{
			name: "custom retries and timeout",
			config: CadreenConfig{
				APIKey:     "test-key",
				MaxRetries: 5,
				Timeout:    10 * time.Second,
				Profile:    "minimal",
			},
			wantBaseURL:    DefaultBaseURL,
			wantMaxRetries: 5,
			wantTimeout:    10 * time.Second,
			wantProfile:    "minimal",
			wantSandbox:    false,
		},
		{
			name: "trailing slash stripped from base URL",
			config: CadreenConfig{
				APIKey:  "test-key",
				BaseURL: "https://custom.example.com/",
			},
			wantBaseURL:    "https://custom.example.com",
			wantMaxRetries: DefaultMaxRetries,
			wantTimeout:    DefaultTimeout,
			wantProfile:    DefaultProfile,
			wantSandbox:    false,
		},
		{
			name: "zero MaxRetries uses default",
			config: CadreenConfig{
				APIKey:     "test-key",
				MaxRetries: 0,
			},
			wantBaseURL:    DefaultBaseURL,
			wantMaxRetries: DefaultMaxRetries,
			wantTimeout:    DefaultTimeout,
			wantProfile:    DefaultProfile,
			wantSandbox:    false,
		},
		{
			name: "negative MaxRetries uses default",
			config: CadreenConfig{
				APIKey:     "test-key",
				MaxRetries: -1,
			},
			wantBaseURL:    DefaultBaseURL,
			wantMaxRetries: DefaultMaxRetries,
			wantTimeout:    DefaultTimeout,
			wantProfile:    DefaultProfile,
			wantSandbox:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.config)
			if c == nil {
				t.Fatal("NewClient returned nil")
			}
			if c.config.BaseURL != tt.wantBaseURL {
				t.Errorf("BaseURL = %q, want %q", c.config.BaseURL, tt.wantBaseURL)
			}
			if c.config.MaxRetries != tt.wantMaxRetries {
				t.Errorf("MaxRetries = %d, want %d", c.config.MaxRetries, tt.wantMaxRetries)
			}
			if c.config.Timeout != tt.wantTimeout {
				t.Errorf("Timeout = %v, want %v", c.config.Timeout, tt.wantTimeout)
			}
			if c.config.Profile != tt.wantProfile {
				t.Errorf("Profile = %q, want %q", c.config.Profile, tt.wantProfile)
			}
			if c.config.Sandbox != tt.wantSandbox {
				t.Errorf("Sandbox = %v, want %v", c.config.Sandbox, tt.wantSandbox)
			}
			if c.config.APIKey != tt.config.APIKey {
				t.Errorf("APIKey = %q, want %q", c.config.APIKey, tt.config.APIKey)
			}
			if c.httpClient == nil {
				t.Error("httpClient is nil")
			}
		})
	}
}

func TestSandboxRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("fixture matching with METHOD /path", func(t *testing.T) {
		expectValue := "method_path_match"
		client := NewClient(CadreenConfig{
			Sandbox: true,
			Fixtures: map[string]any{
				"POST /api/v1/cadreen/test": map[string]any{
					"value": expectValue,
				},
			},
		})
		var result map[string]any
		err := client.do(ctx, "POST", "/api/v1/cadreen/test", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["value"] != expectValue {
			t.Errorf("value = %q, want %q", result["value"], expectValue)
		}
	})

	t.Run("fixture matching with /path fallback", func(t *testing.T) {
		expectValue := "path_fallback_match"
		client := NewClient(CadreenConfig{
			Sandbox: true,
			Fixtures: map[string]any{
				"/api/v1/cadreen/test": map[string]any{
					"value": expectValue,
				},
			},
		})
		var result map[string]any
		err := client.do(ctx, "DELETE", "/api/v1/cadreen/test", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["value"] != expectValue {
			t.Errorf("value = %q, want %q", result["value"], expectValue)
		}
	})

	t.Run("missing fixture error", func(t *testing.T) {
		client := NewClient(CadreenConfig{
			Sandbox:  true,
			Fixtures: map[string]any{},
		})
		var result map[string]any
		err := client.do(ctx, "POST", "/api/v1/nonexistent", nil, &result)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T: %v", err, err)
		}
		if apiErr.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
		}
		if !strings.Contains(apiErr.Message, "No fixture") {
			t.Errorf("expected 'No fixture' in message, got %q", apiErr.Message)
		}
	})
}

func TestRetryConfiguration(t *testing.T) {
	t.Run("default MaxRetries = 3", func(t *testing.T) {
		c := NewClient(CadreenConfig{APIKey: "key"})
		if c.config.MaxRetries != 3 {
			t.Errorf("Default MaxRetries = %d, want 3", c.config.MaxRetries)
		}
	})

	t.Run("custom MaxRetries", func(t *testing.T) {
		c := NewClient(CadreenConfig{APIKey: "key", MaxRetries: 7})
		if c.config.MaxRetries != 7 {
			t.Errorf("MaxRetries = %d, want 7", c.config.MaxRetries)
		}
	})
}

func TestAPIError(t *testing.T) {
	t.Run("parseError with valid error JSON", func(t *testing.T) {
		body := []byte(`{
			"error": {
				"type": "invalid_request",
				"code": "missing_field",
				"message": "Field 'name' is required",
				"hint": "Include the 'name' field in your request body.",
				"details": [{"field": "name", "message": "required"}],
				"retry_after": 0
			}
		}`)

		err := parseError(400, body)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.StatusCode != 400 {
			t.Errorf("StatusCode = %d, want 400", err.StatusCode)
		}
		if err.Type != "invalid_request" {
			t.Errorf("Type = %q, want 'invalid_request'", err.Type)
		}
		if err.Code != "missing_field" {
			t.Errorf("Code = %q, want 'missing_field'", err.Code)
		}
		if err.Message != "Field 'name' is required" {
			t.Errorf("Message = %q, want 'Field 'name' is required'", err.Message)
		}
		if err.Hint != "Include the 'name' field in your request body." {
			t.Errorf("Hint = %q, want 'Include the 'name' field in your request body.'", err.Hint)
		}
		if len(err.Details) != 1 {
			t.Errorf("Details length = %d, want 1", len(err.Details))
		} else {
			if err.Details[0].Field != "name" {
				t.Errorf("Details[0].Field = %q, want 'name'", err.Details[0].Field)
			}
		}
	})

	t.Run("parseError with invalid JSON falls back", func(t *testing.T) {
		body := []byte("not json")
		err := parseError(502, body)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.StatusCode != 502 {
			t.Errorf("StatusCode = %d, want 502", err.StatusCode)
		}
		if err.Type != "error" {
			t.Errorf("Type = %q, want 'error'", err.Type)
		}
		if err.Code != "unknown" {
			t.Errorf("Code = %q, want 'unknown'", err.Code)
		}
	})

	t.Run("APIError.Error() returns message", func(t *testing.T) {
		err := &APIError{Message: "something went wrong"}
		if err.Error() != "something went wrong" {
			t.Errorf("Error() = %q, want 'something went wrong'", err.Error())
		}
	})

	t.Run("APIError with full fields", func(t *testing.T) {
		err := &APIError{
			StatusCode: 429,
			Type:       "rate_limit",
			Code:       "too_many_requests",
			Message:    "Rate limit exceeded. Try again later.",
			Hint:       "Reduce request frequency.",
			NextAction: &NextAction{Type: "retry", Label: "Wait and retry"},
			Details:    []ErrorDetail{{Field: "rate", Message: "exceeded"}},
			RetryAfter: 30,
		}
		if err.Error() != "Rate limit exceeded. Try again later." {
			t.Errorf("Error() = %q", err.Error())
		}
	})
}

func TestGenerateIdempotencyKey(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 100; i++ {
		key := generateIdempotencyKey()
		if seen[key] {
			t.Errorf("duplicate idempotency key generated: %q", key)
		}
		seen[key] = true

		parts := strings.Split(key, "-")
		if len(parts) != 5 {
			t.Errorf("key %q has %d parts, want 5 (UUID v4 format)", key, len(parts))
		}

		if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
			t.Errorf("key %q has wrong segment lengths", key)
		}

		versionChar := parts[2][0]
		if versionChar != '4' {
			t.Errorf("key %q: version nibble = %c, want '4'", key, versionChar)
		}
	}
}

func TestPopulateResult(t *testing.T) {
	type testDest struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email,omitempty"`
	}

	t.Run("populates struct from map fixture", func(t *testing.T) {
		src := map[string]any{"name": "Alice", "age": 30}
		var dst testDest
		err := populateResult(src, &dst)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "Alice" {
			t.Errorf("Name = %q, want 'Alice'", dst.Name)
		}
		if dst.Age != 30 {
			t.Errorf("Age = %d, want 30", dst.Age)
		}
	})

	t.Run("populates struct from JSON bytes", func(t *testing.T) {
		src := json.RawMessage(`{"name": "Bob", "age": 25}`)
		var dst testDest
		err := populateResult(src, &dst)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "Bob" {
			t.Errorf("Name = %q, want 'Bob'", dst.Name)
		}
		if dst.Age != 25 {
			t.Errorf("Age = %d, want 25", dst.Age)
		}
	})

	t.Run("nil dst returns nil", func(t *testing.T) {
		src := map[string]any{"name": "ignored"}
		err := populateResult(src, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("extra keys are ignored", func(t *testing.T) {
		src := map[string]any{"name": "Charlie", "age": 42, "extra_field": "ignored"}
		var dst testDest
		err := populateResult(src, &dst)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "Charlie" {
			t.Errorf("Name = %q, want 'Charlie'", dst.Name)
		}
		if dst.Age != 42 {
			t.Errorf("Age = %d, want 42", dst.Age)
		}
	})
}

func TestIntentStreamRejectsSandbox(t *testing.T) {
	ctx := context.Background()
	client := NewClient(CadreenConfig{
		Sandbox: true,
	})
	req := IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "test"}},
	}
	_, err := client.IntentStream(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "streaming is not available in sandbox mode" {
		t.Errorf("error = %q, want 'streaming is not available in sandbox mode'", err.Error())
	}
}
