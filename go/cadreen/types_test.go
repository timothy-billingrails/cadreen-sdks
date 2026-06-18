package cadreen

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestClarificationQuestions_UnmarshalJSON(t *testing.T) {
	t.Run("valid JSON object array", func(t *testing.T) {
		data := []byte(`[
			{"id": "q1", "question": "What is your budget?", "type": "text", "required": true, "reason": "needed for planning"},
			{"id": "q2", "question": "What is the deadline?", "type": "date", "required": false}
		]`)
		var q ClarificationQuestions
		err := q.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(q) != 2 {
			t.Fatalf("len = %d, want 2", len(q))
		}
		if q[0].ID != "q1" {
			t.Errorf("q[0].ID = %q, want 'q1'", q[0].ID)
		}
		if q[0].Question != "What is your budget?" {
			t.Errorf("q[0].Question = %q", q[0].Question)
		}
		if q[0].Type != "text" {
			t.Errorf("q[0].Type = %q, want 'text'", q[0].Type)
		}
		if !q[0].Required {
			t.Error("q[0].Required should be true")
		}
		if q[0].Reason != "needed for planning" {
			t.Errorf("q[0].Reason = %q", q[0].Reason)
		}
		if q[1].Type != "date" {
			t.Errorf("q[1].Type = %q, want 'date'", q[1].Type)
		}
	})

	t.Run("valid JSON string array (legacy format)", func(t *testing.T) {
		data := []byte(`["First question?", "Second question?"]`)
		var q ClarificationQuestions
		err := q.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(q) != 2 {
			t.Fatalf("len = %d, want 2", len(q))
		}
		if q[0].Question != "First question?" {
			t.Errorf("q[0].Question = %q", q[0].Question)
		}
		if q[0].Type != "open" {
			t.Errorf("q[0].Type = %q, want 'open' (default from legacy)", q[0].Type)
		}
		if q[0].Required {
			t.Error("q[0].Required should be false (default from legacy)")
		}
		if q[1].Question != "Second question?" {
			t.Errorf("q[1].Question = %q", q[1].Question)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid`)
		var q ClarificationQuestions
		err := q.UnmarshalJSON(data)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty array", func(t *testing.T) {
		data := []byte(`[]`)
		var q ClarificationQuestions
		err := q.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(q) != 0 {
			t.Errorf("len = %d, want 0", len(q))
		}
	})

	t.Run("empty array via json.Unmarshal", func(t *testing.T) {
		data := []byte(`[]`)
		var q ClarificationQuestions
		err := json.Unmarshal(data, &q)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(q) != 0 {
			t.Errorf("len = %d, want 0", len(q))
		}
	})
}

func TestAllAtomCategories(t *testing.T) {
	expected := []AtomCategory{
		AtomCategoryReference,
		AtomCategoryPreference,
		AtomCategoryEpisode,
		AtomCategoryPrecedent,
		AtomCategoryNote,
		AtomCategoryProject,
	}
	result := AllAtomCategories()
	if len(result) != len(expected) {
		t.Fatalf("len = %d, want %d", len(result), len(expected))
	}
	for i, cat := range expected {
		if result[i] != cat {
			t.Errorf("result[%d] = %q, want %q", i, result[i], cat)
		}
	}
}

func TestAllMemoryTypes(t *testing.T) {
	expected := AllAtomCategories()
	result := AllMemoryTypes()
	if len(result) != len(expected) {
		t.Fatalf("len = %d, want %d", len(result), len(expected))
	}
	for i, cat := range expected {
		if result[i] != cat {
			t.Errorf("result[%d] = %q, want %q", i, result[i], cat)
		}
	}
}

func TestIntentResultTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value IntentResultType
		want  string
	}{
		{"direct", IntentResultDirect, "direct"},
		{"clarify", IntentResultClarify, "clarify"},
		{"execution", IntentResultExecution, "execution"},
		{"blocked", IntentResultBlocked, "blocked"},
		{"connect_required", IntentResultConnectRequired, "connect_required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.want)
			}
		})
	}
}

func TestGovernanceDecisionTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value GovernanceDecisionType
		want  string
	}{
		{"auto", GovernanceAuto, "auto"},
		{"auto_complete", GovernanceAutoComplete, "auto_complete"},
		{"handoff", GovernanceHandoff, "handoff"},
		{"escalate", GovernanceEscalate, "escalate"},
		{"clarify_requester", GovernanceClarifyRequester, "clarify_requester"},
		{"abstain", GovernanceAbstain, "abstain"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.want)
			}
		})
	}
}

func TestHealthStatusConstants(t *testing.T) {
	tests := []struct {
		name  string
		value HealthStatus
		want  string
	}{
		{"healthy", HealthStatusHealthy, "healthy"},
		{"degraded", HealthStatusDegraded, "degraded"},
		{"unhealthy", HealthStatusUnhealthy, "unhealthy"},
		{"unknown", HealthStatusUnknown, "unknown"},
		{"latent", HealthStatusLatent, "latent"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.want)
			}
		})
	}
}

func TestErrorTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value ErrorType
		want  string
	}{
		{"invalid_request", ErrorTypeInvalidRequest, "invalid_request"},
		{"authentication_error", ErrorTypeAuthenticationError, "authentication_error"},
		{"permission_error", ErrorTypePermissionError, "permission_error"},
		{"not_found", ErrorTypeNotFound, "not_found"},
		{"conflict", ErrorTypeConflict, "conflict"},
		{"validation_error", ErrorTypeValidationError, "validation_error"},
		{"rate_limit", ErrorTypeRateLimit, "rate_limit"},
		{"internal_error", ErrorTypeInternalError, "internal_error"},
		{"service_unavailable", ErrorTypeServiceUnavailable, "service_unavailable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.want)
			}
		})
	}
}

func TestMemoryTypeAliases(t *testing.T) {
	if MemoryTypeReference != AtomCategoryReference {
		t.Error("MemoryTypeReference should equal AtomCategoryReference")
	}
	if MemoryTypePreference != AtomCategoryPreference {
		t.Error("MemoryTypePreference should equal AtomCategoryPreference")
	}
	if MemoryTypeEpisode != AtomCategoryEpisode {
		t.Error("MemoryTypeEpisode should equal AtomCategoryEpisode")
	}
	if MemoryTypePrecedent != AtomCategoryPrecedent {
		t.Error("MemoryTypePrecedent should equal AtomCategoryPrecedent")
	}
	if MemoryTypeNote != AtomCategoryNote {
		t.Error("MemoryTypeNote should equal AtomCategoryNote")
	}
	if MemoryTypeProject != AtomCategoryProject {
		t.Error("MemoryTypeProject should equal AtomCategoryProject")
	}
}

func TestMemoryItemContentAlias(t *testing.T) {
	var a AtomContent
	var b MemoryItemContent
	a = b
	b = a
	_ = a
	_ = b
}

func TestMemTypesEqualAtomCategories(t *testing.T) {
	memTypes := AllMemoryTypes()
	atomCats := AllAtomCategories()
	if !reflect.DeepEqual(memTypes, atomCats) {
		t.Error("AllMemoryTypes() should return the same slice as AllAtomCategories()")
	}
}

func TestRecoveryStatusConstants(t *testing.T) {
	tests := []struct {
		name  string
		value RecoveryStatus
		want  string
	}{
		{"diagnosing", RecoveryDiagnosing, "diagnosing"},
		{"recovering", RecoveryRecovering, "recovering"},
		{"sub_execution", RecoverySubExecution, "sub_execution"},
		{"escalating", RecoveryEscalating, "escalating"},
		{"recovered", RecoveryRecovered, "recovered"},
		{"failed", RecoveryFailed, "failed"},
		{"skipped", RecoverySkipped, "skipped"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.want)
			}
		})
	}
}

func TestIntentModeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value IntentMode
		want  string
	}{
		{"auto", IntentModeAuto, "auto"},
		{"chat", IntentModeChat, "chat"},
		{"execution", IntentModeExecution, "execution"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.want)
			}
		})
	}
}
