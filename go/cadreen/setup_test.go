package cadreen

import (
	"context"
	"testing"
)

func TestSetup(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/setup": map[string]any{
			"workspace_id": "ws_001",
			"connections": []map[string]any{
				{
					"capability": "email",
					"status":     "configured",
					"detail":     map[string]any{"provider": "gmail"},
				},
				{
					"capability": "database",
					"status":     "failed",
					"error":      "No database connectors available",
				},
			},
			"credentials": []map[string]any{
				{
					"provider": "google",
					"name":     "Gmail API Key",
					"id":       "cred_001",
					"status":   "created",
				},
			},
			"memory": []map[string]any{
				{
					"id":         "mem_setup_001",
					"type":       "preference",
					"kind":       "workspace_config",
					"classified": false,
					"status":     "created",
				},
			},
			"policies": []map[string]any{
				{
					"name":   "Default Spend Limit",
					"id":     "pol_setup_001",
					"status": "draft",
				},
			},
			"applied": 3,
			"failed":  1,
			"proposals": []map[string]any{
				{
					"type":        "connection",
					"description": "Connect to Gmail for email capabilities",
					"detail":      "OAuth2 authentication required",
				},
				{
					"type":        "policy",
					"description": "Create spend limit policy",
					"detail":      "Default limit: $1000",
				},
			},
		},
	})

	result, err := client.Setup(ctx, SetupRequest{
		Purpose: "Set up workspace for customer support team",
		Connections: []SetupConnection{
			{Capability: "email"},
			{Capability: "database"},
		},
		Credentials: []SetupCredential{
			{
				Provider: "google",
				Name:     "Gmail API Key",
				KeyData:  map[string]any{"api_key": "xxx"},
			},
		},
		Memory: []SetupMemory{
			{
				Type:   "preference",
				Domain: "workspace",
				Content: map[string]any{
					"theme": "dark",
				},
			},
		},
		Policies: []SetupPolicy{
			{
				Name: "Default Spend Limit",
				Rule: "spend < 1000",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.WorkspaceID != "ws_001" {
		t.Errorf("WorkspaceID = %q, want 'ws_001'", result.WorkspaceID)
	}
	if result.Applied != 3 {
		t.Errorf("Applied = %d, want 3", result.Applied)
	}
	if result.Failed != 1 {
		t.Errorf("Failed = %d, want 1", result.Failed)
	}

	if len(result.Connections) != 2 {
		t.Fatalf("Connections length = %d, want 2", len(result.Connections))
	}
	if result.Connections[0].Capability != "email" {
		t.Errorf("Connections[0].Capability = %q, want 'email'", result.Connections[0].Capability)
	}
	if result.Connections[0].Status != "configured" {
		t.Errorf("Connections[0].Status = %q, want 'configured'", result.Connections[0].Status)
	}
	if result.Connections[1].Capability != "database" {
		t.Errorf("Connections[1].Capability = %q, want 'database'", result.Connections[1].Capability)
	}
	if result.Connections[1].Status != "failed" {
		t.Errorf("Connections[1].Status = %q, want 'failed'", result.Connections[1].Status)
	}
	if result.Connections[1].Error != "No database connectors available" {
		t.Errorf("Connections[1].Error = %q", result.Connections[1].Error)
	}

	if len(result.Credentials) != 1 {
		t.Fatalf("Credentials length = %d, want 1", len(result.Credentials))
	}
	if result.Credentials[0].Provider != "google" {
		t.Errorf("Credentials[0].Provider = %q, want 'google'", result.Credentials[0].Provider)
	}
	if result.Credentials[0].ID != "cred_001" {
		t.Errorf("Credentials[0].ID = %q, want 'cred_001'", result.Credentials[0].ID)
	}
	if result.Credentials[0].Status != "created" {
		t.Errorf("Credentials[0].Status = %q, want 'created'", result.Credentials[0].Status)
	}

	if len(result.Memory) != 1 {
		t.Fatalf("Memory length = %d, want 1", len(result.Memory))
	}
	if result.Memory[0].ID != "mem_setup_001" {
		t.Errorf("Memory[0].ID = %q, want 'mem_setup_001'", result.Memory[0].ID)
	}
	if result.Memory[0].Type != "preference" {
		t.Errorf("Memory[0].Type = %q, want 'preference'", result.Memory[0].Type)
	}

	if len(result.Policies) != 1 {
		t.Fatalf("Policies length = %d, want 1", len(result.Policies))
	}
	if result.Policies[0].Name != "Default Spend Limit" {
		t.Errorf("Policies[0].Name = %q, want 'Default Spend Limit'", result.Policies[0].Name)
	}
	if result.Policies[0].ID != "pol_setup_001" {
		t.Errorf("Policies[0].ID = %q, want 'pol_setup_001'", result.Policies[0].ID)
	}

	if len(result.Proposals) != 2 {
		t.Fatalf("Proposals length = %d, want 2", len(result.Proposals))
	}
	if result.Proposals[0].Type != "connection" {
		t.Errorf("Proposals[0].Type = %q, want 'connection'", result.Proposals[0].Type)
	}
	if result.Proposals[1].Type != "policy" {
		t.Errorf("Proposals[1].Type = %q, want 'policy'", result.Proposals[1].Type)
	}
}

func TestSetup_Minimal(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/setup": map[string]any{
			"workspace_id": "ws_minimal",
			"connections":  []map[string]any{},
			"credentials":  []map[string]any{},
			"memory":       []map[string]any{},
			"policies":     []map[string]any{},
			"applied":      0,
			"failed":       0,
		},
	})

	result, err := client.Setup(ctx, SetupRequest{
		Purpose: "minimal workspace",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.WorkspaceID != "ws_minimal" {
		t.Errorf("WorkspaceID = %q, want 'ws_minimal'", result.WorkspaceID)
	}
	if result.Applied != 0 {
		t.Errorf("Applied = %d, want 0", result.Applied)
	}
	if result.Failed != 0 {
		t.Errorf("Failed = %d, want 0", result.Failed)
	}
	if len(result.Connections) != 0 {
		t.Errorf("Connections length = %d, want 0", len(result.Connections))
	}
}
