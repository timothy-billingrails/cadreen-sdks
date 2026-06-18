package cadreen

import (
	"context"
	"testing"
)

func TestDraft(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/policies": map[string]any{
			"id":                   "pol_001",
			"name":                 "Spend Limit Policy",
			"version":              1,
			"status":               "draft",
			"confirmation_required": true,
			"approve_url":          "https://example.com/approve/pol_001",
		},
	})

	result, err := client.Draft(ctx, CreatePolicyRequest{
		Name:   "Spend Limit Policy",
		Domain: "finance",
		Rules: []map[string]any{
			{"field": "amount", "operator": "lte", "value": 1000},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "pol_001" {
		t.Errorf("ID = %q, want 'pol_001'", result.ID)
	}
	if result.Name != "Spend Limit Policy" {
		t.Errorf("Name = %q", result.Name)
	}
	if result.Version != 1 {
		t.Errorf("Version = %d, want 1", result.Version)
	}
	if result.Status != "draft" {
		t.Errorf("Status = %q, want 'draft'", result.Status)
	}
	if !result.ConfirmationRequired {
		t.Error("ConfirmationRequired should be true")
	}
	if result.ApproveURL != "https://example.com/approve/pol_001" {
		t.Errorf("ApproveURL = %q", result.ApproveURL)
	}
}

func TestDraft_AutoDraft(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/policies": map[string]any{
			"id":      "pol_002",
			"name":    "Auto-generated Policy",
			"version": 1,
			"status":  "draft",
		},
	})

	result, err := client.Draft(ctx, CreatePolicyRequest{
		Domain:    "operations",
		AutoDraft: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "pol_002" {
		t.Errorf("ID = %q, want 'pol_002'", result.ID)
	}
	if result.Name != "Auto-generated Policy" {
		t.Errorf("Name = %q", result.Name)
	}
}

func TestEvaluate(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/policies/evaluate": map[string]any{
			"action": "allow",
			"domain": "finance",
			"result": map[string]any{
				"type":       "auto",
				"confidence": 0.95,
				"reason":     "Within spend limit policy",
			},
		},
	})

	result, err := client.Evaluate(ctx, "transfer_funds", "finance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "allow" {
		t.Errorf("Action = %q, want 'allow'", result.Action)
	}
	if result.Domain != "finance" {
		t.Errorf("Domain = %q, want 'finance'", result.Domain)
	}
	if result.GovernanceResult.Type != "auto" {
		t.Errorf("GovernanceResult.Type = %q, want 'auto'", result.GovernanceResult.Type)
	}
	if result.GovernanceResult.Confidence != 0.95 {
		t.Errorf("GovernanceResult.Confidence = %f, want 0.95", result.GovernanceResult.Confidence)
	}
	if result.GovernanceResult.Reason != "Within spend limit policy" {
		t.Errorf("GovernanceResult.Reason = %q", result.GovernanceResult.Reason)
	}
}

func TestEvaluatePolicy(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/policies/evaluate": map[string]any{
			"action": "block",
			"domain": "security",
			"result": map[string]any{
				"type":       "escalate",
				"confidence": 0.8,
				"reason":     "Requires human approval",
			},
		},
	})

	result, err := client.EvaluatePolicy(ctx, EvaluatePolicyRequest{
		Action: "delete_data",
		Domain: "security",
		Context: map[string]any{
			"data_type": "PII",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "block" {
		t.Errorf("Action = %q, want 'block'", result.Action)
	}
	if result.Domain != "security" {
		t.Errorf("Domain = %q, want 'security'", result.Domain)
	}
	if result.GovernanceResult.Type != "escalate" {
		t.Errorf("GovernanceResult.Type = %q, want 'escalate'", result.GovernanceResult.Type)
	}
	if result.GovernanceResult.Confidence != 0.8 {
		t.Errorf("GovernanceResult.Confidence = %f, want 0.8", result.GovernanceResult.Confidence)
	}
}

func TestConfirm(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/policies/pol_001/confirm": map[string]any{
			"id":               "pol_001",
			"version":          2,
			"previous_version": 1,
			"status":           "active",
			"confirmed_at":     "2025-06-15T14:30:00Z",
		},
	})

	result, err := client.Confirm(ctx, "pol_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "pol_001" {
		t.Errorf("ID = %q, want 'pol_001'", result.ID)
	}
	if result.Version != 2 {
		t.Errorf("Version = %d, want 2", result.Version)
	}
	if result.PreviousVersion != 1 {
		t.Errorf("PreviousVersion = %d, want 1", result.PreviousVersion)
	}
	if result.Status != "active" {
		t.Errorf("Status = %q, want 'active'", result.Status)
	}
	if result.ConfirmedAt != "2025-06-15T14:30:00Z" {
		t.Errorf("ConfirmedAt = %q", result.ConfirmedAt)
	}
}

func TestConfirm_AlreadyActive(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/policies/pol_002/confirm": map[string]any{
			"id":             "pol_002",
			"version":        1,
			"status":         "active",
			"already_active": true,
		},
	})

	result, err := client.Confirm(ctx, "pol_002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.AlreadyActive {
		t.Error("AlreadyActive should be true")
	}
}
