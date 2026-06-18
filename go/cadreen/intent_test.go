package cadreen

import (
	"context"
	"testing"
)

func newSandboxClient(fixtures map[string]any) *Client {
	return NewClient(CadreenConfig{
		APIKey:   "sandbox_key",
		Sandbox:  true,
		Fixtures: fixtures,
	})
}

func TestIntentInvoke_DirectResult(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "direct",
			"status":   "answered",
			"trace_id": "trace_direct_001",
			"message": map[string]any{
				"role":    "assistant",
				"content": "Your invoice has been refunded.",
			},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "Refund invoice inv_123"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Type != IntentResultDirect {
		t.Errorf("Type = %q, want %q", result.Type, IntentResultDirect)
	}
	if result.Status != "answered" {
		t.Errorf("Status = %q, want 'answered'", result.Status)
	}
	if result.TraceID != "trace_direct_001" {
		t.Errorf("TraceID = %q, want 'trace_direct_001'", result.TraceID)
	}
	if result.Message == nil {
		t.Fatal("Message is nil")
	}
	if result.Message.Role != "assistant" {
		t.Errorf("Message.Role = %q, want 'assistant'", result.Message.Role)
	}
	if result.Message.Content != "Your invoice has been refunded." {
		t.Errorf("Message.Content = %q", result.Message.Content)
	}
	if result.Intelligence == nil {
		t.Error("Intelligence should not be nil (default provided)")
	}
}

func TestIntentInvoke_BlockedResult(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "blocked",
			"trace_id": "trace_blocked_001",
			"meta": map[string]any{
				"governance": map[string]any{
					"reason_code": "exceeds_spend_limit",
					"policy_id":   "pol_spend_001",
				},
			},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "Buy 1000 units"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Type != IntentResultBlocked {
		t.Errorf("Type = %q, want %q", result.Type, IntentResultBlocked)
	}
	if result.Status != "blocked_by_policy" {
		t.Errorf("Status = %q, want 'blocked_by_policy'", result.Status)
	}
	if result.ReasonCode != "exceeds_spend_limit" {
		t.Errorf("ReasonCode = %q, want 'exceeds_spend_limit'", result.ReasonCode)
	}
	if result.PolicyID != "pol_spend_001" {
		t.Errorf("PolicyID = %q, want 'pol_spend_001'", result.PolicyID)
	}
	if result.TraceID != "trace_blocked_001" {
		t.Errorf("TraceID = %q, want 'trace_blocked_001'", result.TraceID)
	}
}

func TestIntentInvoke_ClarifyResult(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "clarify",
			"trace_id": "trace_clarify_001",
			"clarification": map[string]any{
				"conversation_id": "conv_001",
				"questions": []map[string]any{
					{
						"id":       "q1",
						"question": "What is the priority level?",
						"type":     "select",
						"required": true,
					},
				},
			},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "Deploy to production"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Type != IntentResultClarify {
		t.Errorf("Type = %q, want %q", result.Type, IntentResultClarify)
	}
	if result.Status != "needs_input" {
		t.Errorf("Status = %q, want 'needs_input'", result.Status)
	}
	if result.ConversationID != "conv_001" {
		t.Errorf("ConversationID = %q, want 'conv_001'", result.ConversationID)
	}
	if len(result.Questions) != 1 {
		t.Fatalf("len(Questions) = %d, want 1", len(result.Questions))
	}
	if result.Questions[0].ID != "q1" {
		t.Errorf("Questions[0].ID = %q", result.Questions[0].ID)
	}
}

func TestIntentInvoke_ExecutionResult(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "execution",
			"trace_id": "trace_exec_001",
			"execution": map[string]any{
				"id":        "exec_001",
				"status":    "running",
				"stream_url": "https://example.com/stream",
				"poll_url":   "https://example.com/poll",
			},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "Run report"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Type != IntentResultExecution {
		t.Errorf("Type = %q, want %q", result.Type, IntentResultExecution)
	}
	if result.Status != "ready_to_execute" {
		t.Errorf("Status = %q, want 'ready_to_execute'", result.Status)
	}
	if result.Execution == nil {
		t.Fatal("Execution is nil")
	}
	if result.Execution.ID != "exec_001" {
		t.Errorf("Execution.ID = %q, want 'exec_001'", result.Execution.ID)
	}
	if result.Execution.Status != "running" {
		t.Errorf("Execution.Status = %q, want 'running'", result.Execution.Status)
	}
}

func TestIntentInvoke_ConnectRequired(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "connect_required",
			"trace_id": "trace_connect_001",
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "Check Salesforce"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Type != IntentResultConnectRequired {
		t.Errorf("Type = %q, want %q", result.Type, IntentResultConnectRequired)
	}
	if result.Status != "blocked_by_missing_connection" {
		t.Errorf("Status = %q, want 'blocked_by_missing_connection'", result.Status)
	}
	if result.NextAction == nil {
		t.Fatal("NextAction is nil for connect_required")
	}
	if result.NextAction.Endpoint != "/api/v1/cadreen/connections" {
		t.Errorf("NextAction.Endpoint = %q", result.NextAction.Endpoint)
	}
}

func TestIntentInvoke_MissionMappedToExecution(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "mission",
			"trace_id": "trace_mission_001",
			"mission": map[string]any{
				"id":        "mission_001",
				"status":    "pending",
				"stream_url": "https://example.com/mission-stream",
			},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "Launch campaign"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Type != IntentResultExecution {
		t.Errorf("Type = %q, want %q (mission maps to execution)", result.Type, IntentResultExecution)
	}
	if result.Execution == nil {
		t.Fatal("Execution is nil")
	}
	if result.Execution.ID != "mission_001" {
		t.Errorf("Execution.ID = %q, want 'mission_001'", result.Execution.ID)
	}
}

func TestMapIntentResponse(t *testing.T) {
	tests := []struct {
		name      string
		raw       unifiedIntentResponse
		traceID   string
		wantType  IntentResultType
		wantStatus string
	}{
		{
			name: "direct type",
			raw: unifiedIntentResponse{
				Type:    "direct",
				Status:  "completed",
				TraceID: "tr1",
				Message: &ResponseMessage{Role: "assistant", Content: "Done"},
			},
			traceID:   "",
			wantType:  IntentResultDirect,
			wantStatus: "completed",
		},
		{
			name: "clarify type",
			raw: unifiedIntentResponse{
				Type:    "clarify",
				TraceID: "tr2",
				Clarification: &unifiedClarification{
					Questions: []ResponseClarificationQuestion{
						{ID: "q1", Question: "What?", Type: "text", Required: true},
					},
				},
			},
			traceID:   "",
			wantType:  IntentResultClarify,
			wantStatus: "needs_input",
		},
		{
			name: "execution type",
			raw: unifiedIntentResponse{
				Type:    "execution",
				TraceID: "tr3",
				Execution: &ResponseExecution{ID: "exec1", Status: "running"},
			},
			traceID:   "",
			wantType:  IntentResultExecution,
			wantStatus: "ready_to_execute",
		},
		{
			name: "blocked type",
			raw: unifiedIntentResponse{
				Type:    "blocked",
				TraceID: "tr4",
				Meta: &unifiedMeta{
					Governance: &unifiedGovernance{
						ReasonCode: "policy_violation",
						PolicyID:   "pol_001",
					},
				},
			},
			traceID:   "",
			wantType:  IntentResultBlocked,
			wantStatus: "blocked_by_policy",
		},
		{
			name: "connect_required type",
			raw: unifiedIntentResponse{
				Type:    "connect_required",
				TraceID: "tr5",
			},
			traceID:   "",
			wantType:  IntentResultConnectRequired,
			wantStatus: "blocked_by_missing_connection",
		},
		{
			name: "explicit traceID overrides raw",
			raw: unifiedIntentResponse{
				Type:    "direct",
				TraceID: "raw_trace_id",
				Message: &ResponseMessage{Role: "assistant", Content: "OK"},
			},
			traceID:   "explicit_trace_id",
			wantType:  IntentResultDirect,
			wantStatus: "answered",
		},
		{
			name: "unknown type defaults to direct",
			raw: unifiedIntentResponse{
				Type:    "weird_type",
				ID:      "some_id",
				Message: &ResponseMessage{Role: "assistant", Content: "Fallback"},
			},
			traceID:   "",
			wantType:  IntentResultDirect,
			wantStatus: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapIntentResponse(tt.raw, tt.traceID)

			if result.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", result.Type, tt.wantType)
			}
			if result.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", result.Status, tt.wantStatus)
			}

			if tt.traceID != "" {
				if result.TraceID != tt.traceID {
					t.Errorf("TraceID = %q, want %q", result.TraceID, tt.traceID)
				}
			}

			if result.Intelligence == nil {
				t.Error("Intelligence should not be nil (default provided)")
			}
		})
	}
}

func TestMapIntentResponse_IntelligenceHandling(t *testing.T) {
	customIntel := &IntelligenceMeta{
		Summary: "custom summary",
		Governance: GovernanceTrace{
			Active:     true,
			Decision:   "auto",
			Confidence: 0.95,
		},
	}

	t.Run("uses provided intelligence", func(t *testing.T) {
		raw := unifiedIntentResponse{
			Type:         "direct",
			Intelligence: customIntel,
			Message:      &ResponseMessage{Role: "assistant", Content: "OK"},
		}
		result := mapIntentResponse(raw, "")
		if result.Intelligence != customIntel {
			t.Error("Intelligence should be the provided value")
		}
		if result.Intelligence.Summary != "custom summary" {
			t.Errorf("Summary = %q", result.Intelligence.Summary)
		}
	})

	t.Run("falls back to default intelligence when nil", func(t *testing.T) {
		raw := unifiedIntentResponse{
			Type:    "direct",
			Message: &ResponseMessage{Role: "assistant", Content: "OK"},
		}
		result := mapIntentResponse(raw, "")
		if result.Intelligence == nil {
			t.Fatal("Intelligence should use default, not nil")
		}
		if result.Intelligence != defaultIntelligence {
			t.Error("Intelligence should be the default instance")
		}
	})
}

func TestIntentStream_RejectsSandbox(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(nil)
	req := IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "test"}},
	}
	_, err := client.IntentStream(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "streaming is not available in sandbox mode" {
		t.Errorf("error = %q", err.Error())
	}
}

func TestDeriveStatus(t *testing.T) {
	tests := []struct {
		rawType string
		want    string
	}{
		{"direct", "answered"},
		{"clarify", "needs_input"},
		{"mission", "ready_to_execute"},
		{"execution", "ready_to_execute"},
		{"blocked", "blocked_by_policy"},
		{"connect_required", "blocked_by_missing_connection"},
		{"", "failed"},
		{"unknown_type", "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.rawType+" → "+tt.want, func(t *testing.T) {
			got := deriveStatus(tt.rawType)
			if got != tt.want {
				t.Errorf("deriveStatus(%q) = %q, want %q", tt.rawType, got, tt.want)
			}
		})
	}
}

func TestIntentInvoke_TraceIDFromFixture(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":     "direct",
			"status":   "answered",
			"trace_id": "trace_from_fixture",
			"message":  map[string]any{"role": "assistant", "content": "OK"},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "test"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TraceID != "trace_from_fixture" {
		t.Errorf("TraceID = %q, want 'trace_from_fixture'", result.TraceID)
	}
}

func TestIntentInvoke_TraceIDFallbackToID(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intent": map[string]any{
			"type":    "direct",
			"status":  "answered",
			"id":      "fallback_id_123",
			"message": map[string]any{"role": "assistant", "content": "OK"},
		},
	})

	result, err := client.IntentInvoke(ctx, IntentRequest{
		Messages: []IntentMessage{{Role: "user", Content: "test"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TraceID != "fallback_id_123" {
		t.Errorf("TraceID = %q, want 'fallback_id_123' (fallback to ID)", result.TraceID)
	}
}
