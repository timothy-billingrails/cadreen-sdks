package cadreen

import (
	"context"
	"testing"
)

func TestCatalog(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/connections/catalog": map[string]any{
			"categories": []map[string]any{
				{
					"name":        "CRM",
					"description": "Customer relationship management tools",
					"integrations": []map[string]any{
						{
							"id":           "salesforce",
							"name":         "Salesforce",
							"description":  "Cloud CRM platform",
							"category":     "CRM",
							"capabilities": []string{"read_contacts", "create_opportunity"},
							"tags":         []string{"crm", "enterprise"},
							"provider":     "salesforce",
							"status":       "available",
							"auth_type":    "oauth2",
							"install_time": "30s",
							"popularity":   95,
							"featured":     true,
						},
					},
				},
			},
			"installed":       []string{"slack"},
			"total_available": 1,
		},
	})

	result, err := client.Catalog(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Categories) != 1 {
		t.Fatalf("Categories length = %d, want 1", len(result.Categories))
	}
	if result.Categories[0].Name != "CRM" {
		t.Errorf("Categories[0].Name = %q, want 'CRM'", result.Categories[0].Name)
	}
	if len(result.Categories[0].Integrations) != 1 {
		t.Fatalf("Integrations length = %d, want 1", len(result.Categories[0].Integrations))
	}
	integration := result.Categories[0].Integrations[0]
	if integration.ID != "salesforce" {
		t.Errorf("Integration.ID = %q, want 'salesforce'", integration.ID)
	}
	if integration.Popularity != 95 {
		t.Errorf("Integration.Popularity = %d, want 95", integration.Popularity)
	}
	if !integration.Featured {
		t.Error("Featured should be true")
	}
	if len(result.Installed) != 1 {
		t.Errorf("Installed length = %d, want 1", len(result.Installed))
	}
	if result.Installed[0] != "slack" {
		t.Errorf("Installed[0] = %q, want 'slack'", result.Installed[0])
	}
	if result.TotalAvailable != 1 {
		t.Errorf("TotalAvailable = %d, want 1", result.TotalAvailable)
	}
}

func TestInstall(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/connections/install": map[string]any{
			"status":         "pending_auth",
			"auth_url":       "https://accounts.google.com/o/oauth2/auth",
			"provider":       "google",
			"estimated_time": "45s",
		},
	})

	result, err := client.Install(ctx, "google_calendar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "pending_auth" {
		t.Errorf("Status = %q, want 'pending_auth'", result.Status)
	}
	if result.AuthURL != "https://accounts.google.com/o/oauth2/auth" {
		t.Errorf("AuthURL = %q", result.AuthURL)
	}
	if result.Provider != "google" {
		t.Errorf("Provider = %q, want 'google'", result.Provider)
	}
	if result.EstimatedTime != "45s" {
		t.Errorf("EstimatedTime = %q, want '45s'", result.EstimatedTime)
	}
}

func TestRegisterOpenAPI(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/connections/openapi": map[string]any{
			"id":                "conn_openapi_001",
			"name":              "Stripe API",
			"type":              "openapi",
			"tools_generated":   5,
			"tools_registered":  5,
			"functions":         []string{"create_charge", "refund_charge", "list_customers", "create_customer", "get_invoice"},
			"spec_url":          "https://api.stripe.com/openapi.json",
			"status":            "connected",
		},
	})

	result, err := client.RegisterOpenAPI(ctx, RegisterOpenAPIRequest{
		Name:    "Stripe API",
		SpecURL: "https://api.stripe.com/openapi.json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "conn_openapi_001" {
		t.Errorf("ID = %q, want 'conn_openapi_001'", result.ID)
	}
	if result.Name != "Stripe API" {
		t.Errorf("Name = %q", result.Name)
	}
	if result.Type != "openapi" {
		t.Errorf("Type = %q, want 'openapi'", result.Type)
	}
	if result.ToolsGenerated != 5 {
		t.Errorf("ToolsGenerated = %d, want 5", result.ToolsGenerated)
	}
	if result.ToolsRegistered != 5 {
		t.Errorf("ToolsRegistered = %d, want 5", result.ToolsRegistered)
	}
	if result.Status != "connected" {
		t.Errorf("Status = %q, want 'connected'", result.Status)
	}
	if len(result.Functions) != 5 {
		t.Errorf("Functions length = %d, want 5", len(result.Functions))
	}
}

func TestRegisterMCP(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/connections/mcp": map[string]any{
			"id":        "mcp_001",
			"name":      "Filesystem Server",
			"type":      "mcp",
			"status":    "connected",
			"transport": "stdio",
			"url":       "http://localhost:8080/mcp",
		},
	})

	result, err := client.RegisterMCP(ctx, RegisterMCPRequest{
		Name:      "Filesystem Server",
		URL:       "http://localhost:8080/mcp",
		Transport: "stdio",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "mcp_001" {
		t.Errorf("ID = %q, want 'mcp_001'", result.ID)
	}
	if result.Name != "Filesystem Server" {
		t.Errorf("Name = %q", result.Name)
	}
	if result.Type != "mcp" {
		t.Errorf("Type = %q, want 'mcp'", result.Type)
	}
	if result.Status != "connected" {
		t.Errorf("Status = %q, want 'connected'", result.Status)
	}
	if result.Transport != "stdio" {
		t.Errorf("Transport = %q, want 'stdio'", result.Transport)
	}
}

func TestListConnections(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/connections": map[string]any{
			"connections": []map[string]any{
				{
					"capability": "email",
					"pathways": []map[string]any{
						{
							"id":         "pw_001",
							"capability": "email",
							"connector":  "gmail",
							"transport":  "oauth2",
							"health":     "healthy",
							"tool_id":    "tool_send_email",
						},
					},
					"status": "active",
				},
			},
			"total_capabilities": 1,
			"total_pathways":     1,
			"pagination": map[string]any{
				"limit":    20,
				"offset":   0,
				"has_more": false,
			},
		},
	})

	result, err := client.ListConnections(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalCapabilities != 1 {
		t.Errorf("TotalCapabilities = %d, want 1", result.TotalCapabilities)
	}
	if result.TotalPathways != 1 {
		t.Errorf("TotalPathways = %d, want 1", result.TotalPathways)
	}
	if len(result.Connections) != 1 {
		t.Fatalf("Connections length = %d, want 1", len(result.Connections))
	}
	conn := result.Connections[0]
	if conn.Capability != "email" {
		t.Errorf("Capability = %q, want 'email'", conn.Capability)
	}
	if conn.Status != "active" {
		t.Errorf("Status = %q, want 'active'", conn.Status)
	}
	if len(conn.Pathways) != 1 {
		t.Fatalf("Pathways length = %d, want 1", len(conn.Pathways))
	}
	if conn.Pathways[0].ID != "pw_001" {
		t.Errorf("Pathways[0].ID = %q, want 'pw_001'", conn.Pathways[0].ID)
	}
	if conn.Pathways[0].Health != "healthy" {
		t.Errorf("Pathways[0].Health = %q, want 'healthy'", conn.Pathways[0].Health)
	}
	if result.Pagination.Limit != 20 {
		t.Errorf("Pagination.Limit = %d, want 20", result.Pagination.Limit)
	}
}

func TestListCapabilities(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/capabilities": map[string]any{
			"available": []map[string]any{
				{
					"name":        "send_email",
					"human_name":  "Send Email",
					"description": "Send an email via connected provider",
					"score":       0.95,
					"matched_on":  []string{"name", "description"},
					"health":      "healthy",
					"source":      "gmail",
					"status":      "connected",
					"functions":   []string{"send", "schedule"},
					"category":    "communication",
				},
			},
			"gaps": []map[string]any{
				{
					"capability":  "create_invoice",
					"reason":      "No connector installed",
					"description": "Invoice creation requires a billing connector",
					"blocking":    true,
					"severity":    "high",
					"source":      "assessment",
				},
			},
			"count": 1,
			"pagination": map[string]any{
				"limit":    20,
				"offset":   0,
				"has_more": false,
			},
		},
	})

	result, err := client.ListCapabilities(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
	if len(result.Available) != 1 {
		t.Fatalf("Available length = %d, want 1", len(result.Available))
	}
	cap := result.Available[0]
	if cap.Name != "send_email" {
		t.Errorf("Name = %q, want 'send_email'", cap.Name)
	}
	if cap.Score != 0.95 {
		t.Errorf("Score = %f, want 0.95", cap.Score)
	}
	if len(cap.Functions) != 2 {
		t.Errorf("Functions length = %d, want 2", len(cap.Functions))
	}

	if len(result.Gaps) != 1 {
		t.Fatalf("Gaps length = %d, want 1", len(result.Gaps))
	}
	if result.Gaps[0].Capability != "create_invoice" {
		t.Errorf("Gaps[0].Capability = %q, want 'create_invoice'", result.Gaps[0].Capability)
	}
	if !result.Gaps[0].Blocking {
		t.Error("Gaps[0].Blocking should be true")
	}
	if result.Gaps[0].Severity != "high" {
		t.Errorf("Gaps[0].Severity = %q, want 'high'", result.Gaps[0].Severity)
	}
}
