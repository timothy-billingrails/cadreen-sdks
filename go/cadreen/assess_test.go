package cadreen

import (
	"context"
	"testing"
)

func TestAssess(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/assess": map[string]any{
			"assessment": map[string]any{
				"task": "Build a customer dashboard",
				"capabilities": []map[string]any{
					{
						"name":        "query_database",
						"human_name":  "Query Database",
						"description": "Run SQL queries against the data warehouse",
						"score":       0.95,
						"health":      "healthy",
						"status":      "connected",
						"category":    "data",
					},
					{
						"name":        "render_chart",
						"human_name":  "Render Chart",
						"description": "Generate visual charts from data",
						"score":       0.85,
						"health":      "healthy",
						"status":      "connected",
						"category":    "visualization",
					},
				},
				"gaps": []map[string]any{},
				"gap_filling_tasks": []any{},
				"blocking_gaps":     0,
				"policies_recommended": []map[string]any{
					{
						"policy":   "Data access review",
						"reason":   "Dashboard queries customer data",
						"action":   "draft",
						"blocking": true,
					},
				},
				"needs_clarification": []string{"What time range should the dashboard cover?"},
				"can_do":              0.85,
				"assessment_quality":  "high",
				"ready_capabilities":  2,
				"total_capabilities":  2,
				"gap_count":           0,
				"ready_for_deployment": false,
				"stack": map[string]any{
					"cadreen": []map[string]any{
						{
							"name":    "governance",
							"type":    "component",
							"status":  "active",
							"functions": []string{"policy_evaluation"},
						},
					},
				},
				"governance_decision": map[string]any{
					"type":       "auto",
					"confidence": 0.9,
					"reason":     "All capabilities available",
				},
			},
		},
	})

	result, err := client.Assess(ctx, "Build a customer dashboard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Task != "Build a customer dashboard" {
		t.Errorf("Task = %q", result.Task)
	}
	if result.CanDo != 0.85 {
		t.Errorf("CanDo = %f, want 0.85", result.CanDo)
	}
	if result.AssessmentQuality != "high" {
		t.Errorf("AssessmentQuality = %q, want 'high'", result.AssessmentQuality)
	}
	if result.ReadyCapabilities != 2 {
		t.Errorf("ReadyCapabilities = %d, want 2", result.ReadyCapabilities)
	}
	if result.TotalCapabilities != 2 {
		t.Errorf("TotalCapabilities = %d, want 2", result.TotalCapabilities)
	}
	if result.GapCount != 0 {
		t.Errorf("GapCount = %d, want 0", result.GapCount)
	}
	if result.ReadyForDeployment {
		t.Error("ReadyForDeployment should be false")
	}
	if result.BlockingGaps != 0 {
		t.Errorf("BlockingGaps = %d, want 0", result.BlockingGaps)
	}

	if len(result.Capabilities) != 2 {
		t.Fatalf("Capabilities length = %d, want 2", len(result.Capabilities))
	}
	if result.Capabilities[0].Name != "query_database" {
		t.Errorf("Capabilities[0].Name = %q, want 'query_database'", result.Capabilities[0].Name)
	}
	if result.Capabilities[0].Score != 0.95 {
		t.Errorf("Capabilities[0].Score = %f, want 0.95", result.Capabilities[0].Score)
	}

	if len(result.PoliciesRecommended) != 1 {
		t.Fatalf("PoliciesRecommended length = %d, want 1", len(result.PoliciesRecommended))
	}
	if result.PoliciesRecommended[0].Policy != "Data access review" {
		t.Errorf("PoliciesRecommended[0].Policy = %q", result.PoliciesRecommended[0].Policy)
	}
	if !result.PoliciesRecommended[0].Blocking {
		t.Error("PoliciesRecommended[0].Blocking should be true")
	}

	if result.GovernanceDecision == nil {
		t.Fatal("GovernanceDecision is nil")
	}
	if result.GovernanceDecision.Type != "auto" {
		t.Errorf("GovernanceDecision.Type = %q, want 'auto'", result.GovernanceDecision.Type)
	}
	if result.GovernanceDecision.Confidence != 0.9 {
		t.Errorf("GovernanceDecision.Confidence = %f, want 0.9", result.GovernanceDecision.Confidence)
	}

	if len(result.NeedsClarification) != 1 {
		t.Errorf("NeedsClarification length = %d, want 1", len(result.NeedsClarification))
	}
	if result.NeedsClarification[0] != "What time range should the dashboard cover?" {
		t.Errorf("NeedsClarification[0] = %q", result.NeedsClarification[0])
	}
}

func TestAssessWithDomain(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/assess": map[string]any{
			"assessment": map[string]any{
				"task":               "Send weekly report",
				"can_do":             0.70,
				"assessment_quality": "medium",
				"ready_capabilities":  1,
				"total_capabilities":  3,
				"gap_count":           2,
				"ready_for_deployment": false,
				"blocking_gaps":        1,
				"gaps": []map[string]any{
					{
						"capability":  "email_service",
						"description": "No email connector configured",
						"blocking":    true,
						"severity":    "high",
					},
				},
				"capabilities": []map[string]any{
					{
						"name":   "generate_report",
						"score":  0.9,
						"health": "healthy",
						"status": "connected",
					},
				},
			},
		},
	})

	result, err := client.AssessWithDomain(ctx, "Send weekly report", "reporting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Task != "Send weekly report" {
		t.Errorf("Task = %q", result.Task)
	}
	if result.CanDo != 0.70 {
		t.Errorf("CanDo = %f, want 0.70", result.CanDo)
	}
	if result.TotalCapabilities != 3 {
		t.Errorf("TotalCapabilities = %d, want 3", result.TotalCapabilities)
	}
	if result.ReadyCapabilities != 1 {
		t.Errorf("ReadyCapabilities = %d, want 1", result.ReadyCapabilities)
	}
	if result.GapCount != 2 {
		t.Errorf("GapCount = %d, want 2", result.GapCount)
	}
	if result.BlockingGaps != 1 {
		t.Errorf("BlockingGaps = %d, want 1", result.BlockingGaps)
	}
	if len(result.Gaps) != 1 {
		t.Fatalf("Gaps length = %d, want 1", len(result.Gaps))
	}
	if result.Gaps[0].Capability != "email_service" {
		t.Errorf("Gaps[0].Capability = %q, want 'email_service'", result.Gaps[0].Capability)
	}
	if !result.Gaps[0].Blocking {
		t.Error("Gaps[0].Blocking should be true")
	}
}
