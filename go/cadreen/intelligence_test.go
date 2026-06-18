package cadreen

import (
	"context"
	"testing"
)

func TestGetTrace(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/intelligence/trace_001": map[string]any{
			"id":            "trace_001",
			"domain":        "finance",
			"request_path":  "/api/v1/cadreen/intent",
			"request_method": "POST",
			"created_at":    "2025-06-01T12:00:00Z",
			"meta": map[string]any{
				"version": "1.0",
				"summary": "Refund processed successfully",
				"capability": map[string]any{
					"total_available":     5,
					"healthy_count":       5,
					"active_integrations": []string{"stripe"},
				},
				"reasoning": map[string]any{
					"capability_matches": 3,
				},
				"memory": map[string]any{
					"healthy":          true,
					"knowledge_queried": 2,
				},
				"governance": map[string]any{
					"active":     true,
					"decision":   "auto",
					"confidence": 0.92,
				},
				"humility": map[string]any{
					"gaps_detected": 0,
					"blocking":      0,
				},
				"process": map[string]any{
					"started_at":  "2025-06-01T12:00:00Z",
					"duration_ms": 350,
				},
				"field_stability": map[string]any{
					"stable":   []string{"version", "summary"},
					"evolving": []string{"governance.confidence"},
					"internal": []string{},
				},
			},
		},
	})

	result, err := client.GetTrace(ctx, "trace_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "trace_001" {
		t.Errorf("ID = %q, want 'trace_001'", result.ID)
	}
	if result.Domain != "finance" {
		t.Errorf("Domain = %q, want 'finance'", result.Domain)
	}
	if result.RequestPath != "/api/v1/cadreen/intent" {
		t.Errorf("RequestPath = %q", result.RequestPath)
	}
	if result.Method != "POST" {
		t.Errorf("Method = %q, want 'POST'", result.Method)
	}
	if result.Meta.Summary != "Refund processed successfully" {
		t.Errorf("Meta.Summary = %q, want 'Refund processed successfully'", result.Meta.Summary)
	}
	if result.Meta.Capability.TotalAvailable != 5 {
		t.Errorf("Meta.Capability.TotalAvailable = %d, want 5", result.Meta.Capability.TotalAvailable)
	}
	if result.Meta.Governance.Decision != "auto" {
		t.Errorf("Meta.Governance.Decision = %q, want 'auto'", result.Meta.Governance.Decision)
	}
	if result.Meta.Process.DurationMs != 350 {
		t.Errorf("Meta.Process.DurationMs = %d, want 350", result.Meta.Process.DurationMs)
	}
}

func TestListTraces(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/intelligence?domain=finance&limit=5": map[string]any{
			"traces": []map[string]any{
				{
					"id":            "trace_001",
					"domain":        "finance",
					"request_path":  "/api/v1/cadreen/intent",
					"request_method": "POST",
					"meta": map[string]any{
						"version": "1.0",
						"capability": map[string]any{
							"total_available": 3,
							"healthy_count":   3,
						},
						"reasoning": map[string]any{},
						"memory": map[string]any{
							"healthy": true,
						},
						"governance":   map[string]any{"active": false},
						"humility":     map[string]any{},
						"process":      map[string]any{"started_at": "", "duration_ms": 0},
						"field_stability": map[string]any{
							"stable":   []string{},
							"evolving": []string{},
							"internal": []string{},
						},
					},
				},
			},
			"count": 1,
			"pagination": map[string]any{
				"limit":    5,
				"offset":   0,
				"has_more": false,
			},
		},
	})

	result, err := client.ListTraces(ctx, &ListTracesOptions{
		Domain: "finance",
		Limit:  5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
	if len(result.Traces) != 1 {
		t.Fatalf("Traces length = %d, want 1", len(result.Traces))
	}
	if result.Traces[0].ID != "trace_001" {
		t.Errorf("Traces[0].ID = %q, want 'trace_001'", result.Traces[0].ID)
	}
	if result.Traces[0].Domain != "finance" {
		t.Errorf("Traces[0].Domain = %q", result.Traces[0].Domain)
	}
}

func TestListTraces_NoOptions(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/intelligence": map[string]any{
			"traces": []map[string]any{},
			"count":  0,
		},
	})

	result, err := client.ListTraces(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 0 {
		t.Errorf("Count = %d, want 0", result.Count)
	}
	if len(result.Traces) != 0 {
		t.Errorf("Traces length = %d, want 0", len(result.Traces))
	}
}

func TestIntelligenceStats(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/intelligence/stats": map[string]any{
			"traces_24h":         150,
			"traces_7d":          980,
			"traces_30d":         4200,
			"avg_confidence_by_domain": map[string]float64{
				"finance":  0.91,
				"support":  0.87,
				"ops":      0.93,
			},
			"gap_detection_rate": 0.12,
			"governance_decisions": map[string]int{
				"auto":             3500,
				"handoff":          400,
				"escalate":         200,
				"clarify_requester": 100,
			},
		},
	})

	result, err := client.IntelligenceStats(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Traces24h != 150 {
		t.Errorf("Traces24h = %d, want 150", result.Traces24h)
	}
	if result.Traces7d != 980 {
		t.Errorf("Traces7d = %d, want 980", result.Traces7d)
	}
	if result.Traces30d != 4200 {
		t.Errorf("Traces30d = %d, want 4200", result.Traces30d)
	}
	if result.GapDetectionRate != 0.12 {
		t.Errorf("GapDetectionRate = %f, want 0.12", result.GapDetectionRate)
	}
	if result.AvgConfidenceByDomain["finance"] != 0.91 {
		t.Errorf("AvgConfidenceByDomain[finance] = %f, want 0.91", result.AvgConfidenceByDomain["finance"])
	}
	if result.GovernanceDecisions["auto"] != 3500 {
		t.Errorf("GovernanceDecisions[auto] = %d, want 3500", result.GovernanceDecisions["auto"])
	}
}

func TestReplay(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intelligence/trace_001/replay": map[string]any{
			"trace_id":            "trace_001",
			"mode":               "current",
			"domain":             "finance",
			"original_gate":      "auto",
			"original_confidence": 0.92,
			"current_gate":       "auto",
			"current_confidence":  0.94,
			"gate_changed":       false,
			"change_summary":     "Confidence increased slightly due to new memory",
			"current_capability": map[string]interface{}{"total_available": 6},
			"current_memory":     map[string]interface{}{"healthy": true},
			"current_gaps":       map[string]interface{}{"blocking": 0},
			"replay_note":        "Replayed with current state. No gate change.",
		},
	})

	result, err := client.Replay(ctx, "trace_001", "current")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TraceID != "trace_001" {
		t.Errorf("TraceID = %q, want 'trace_001'", result.TraceID)
	}
	if result.Mode != "current" {
		t.Errorf("Mode = %q, want 'current'", result.Mode)
	}
	if result.Domain != "finance" {
		t.Errorf("Domain = %q, want 'finance'", result.Domain)
	}
	if result.OriginalGate != "auto" {
		t.Errorf("OriginalGate = %q, want 'auto'", result.OriginalGate)
	}
	if result.OriginalConfidence != 0.92 {
		t.Errorf("OriginalConfidence = %f, want 0.92", result.OriginalConfidence)
	}
	if result.CurrentGate != "auto" {
		t.Errorf("CurrentGate = %q, want 'auto'", result.CurrentGate)
	}
	if result.CurrentConfidence != 0.94 {
		t.Errorf("CurrentConfidence = %f, want 0.94", result.CurrentConfidence)
	}
	if result.GateChanged {
		t.Error("GateChanged should be false")
	}
}

func TestReplay_EmptyModeDefaultsToCurrent(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intelligence/trace_002/replay": map[string]any{
			"trace_id": "trace_002",
			"mode":     "current",
			"domain":   "ops",
		},
	})

	result, err := client.Replay(ctx, "trace_002", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != "current" {
		t.Errorf("Mode = %q, want 'current' (default)", result.Mode)
	}
}

func TestHandoff(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/intelligence/trace_001/handoff": map[string]any{
			"trace_id":   "trace_001",
			"domain":     "finance",
			"created_at": "2025-06-01T12:05:00Z",
			"governance": map[string]interface{}{
				"decision":   "handoff",
				"confidence": 0.65,
			},
			"what_the_system_knew": map[string]interface{}{
				"user_intent": "Refund invoice inv_123",
			},
			"what_the_system_didnt_know": map[string]interface{}{
				"refund_policy": "unknown",
			},
			"what_happened": map[string]interface{}{
				"attempted": "refund",
				"result":    "blocked by policy",
			},
			"suggested_actions": []map[string]interface{}{
				{"action": "review refund policy", "priority": "high"},
			},
			"next_action": map[string]interface{}{
				"type":  "human_review",
				"label": "Review refund decision",
			},
			"trace_url": "https://accomplishanything.today/traces/trace_001",
		},
	})

	result, err := client.Handoff(ctx, "trace_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TraceID != "trace_001" {
		t.Errorf("TraceID = %q, want 'trace_001'", result.TraceID)
	}
	if result.Domain != "finance" {
		t.Errorf("Domain = %q, want 'finance'", result.Domain)
	}
	if len(result.SuggestedActions) != 1 {
		t.Errorf("SuggestedActions length = %d, want 1", len(result.SuggestedActions))
	}
	if result.TraceURL != "https://accomplishanything.today/traces/trace_001" {
		t.Errorf("TraceURL = %q", result.TraceURL)
	}
}

func TestPromote(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/intelligence/trace_001/promote": map[string]any{
			"id":              "promo_001",
			"kind":            "tool",
			"status":          "active",
			"tool_name":       "refund_invoice",
			"tool_sequence":   []string{"verify_invoice", "process_refund", "notify_customer"},
			"source_trace_id": "trace_001",
		},
	})

	result, err := client.Promote(ctx, "trace_001", "tool", "refund_invoice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "promo_001" {
		t.Errorf("ID = %q, want 'promo_001'", result.ID)
	}
	if result.Kind != "tool" {
		t.Errorf("Kind = %q, want 'tool'", result.Kind)
	}
	if result.Status != "active" {
		t.Errorf("Status = %q, want 'active'", result.Status)
	}
	if result.ToolName != "refund_invoice" {
		t.Errorf("ToolName = %q, want 'refund_invoice'", result.ToolName)
	}
	if result.SourceTraceID != "trace_001" {
		t.Errorf("SourceTraceID = %q, want 'trace_001'", result.SourceTraceID)
	}
	if len(result.ToolSequence) != 3 {
		t.Errorf("ToolSequence length = %d, want 3", len(result.ToolSequence))
	}
	if result.ToolSequence[0] != "verify_invoice" {
		t.Errorf("ToolSequence[0] = %q, want 'verify_invoice'", result.ToolSequence[0])
	}
}
