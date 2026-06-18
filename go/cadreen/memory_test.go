package cadreen

import (
	"context"
	"strings"
	"testing"
)

func TestTeach(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"POST /api/v1/cadreen/memory": map[string]any{
			"id":         "mem_001",
			"type":       "preference",
			"kind":       "user_setting",
			"classified": false,
			"domain":     "settings",
			"scope":      "personal",
			"authority":  5,
			"version":    1,
			"indexed":    true,
			"tags":       []string{"pref", "theme"},
			"created_at": "2025-01-15T10:00:00Z",
			"content": map[string]any{
				"name":    "dark_mode",
				"subject": "ui",
			},
		},
	})

	result, err := client.Teach(ctx, RememberRequest{
		Type:   "preference",
		Domain: "settings",
		Content: map[string]any{
			"name": "dark_mode",
		},
		Tags: []string{"pref", "theme"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "mem_001" {
		t.Errorf("ID = %q, want 'mem_001'", result.ID)
	}
	if result.Type != "preference" {
		t.Errorf("Type = %q, want 'preference'", result.Type)
	}
	if result.Kind != "user_setting" {
		t.Errorf("Kind = %q, want 'user_setting'", result.Kind)
	}
	if result.Domain != "settings" {
		t.Errorf("Domain = %q, want 'settings'", result.Domain)
	}
	if !result.Indexed {
		t.Error("Indexed should be true")
	}
	if len(result.Tags) != 2 {
		t.Errorf("Tags length = %d, want 2", len(result.Tags))
	}
}

func TestSearch(t *testing.T) {
	ctx := context.Background()
	queryKey := "GET /api/v1/cadreen/memory/search?query=dark+mode"
	client := newSandboxClient(map[string]any{
		queryKey: map[string]any{
			"results": []map[string]any{
				{
					"id":         "atom_001",
					"type":       "preference",
					"domain":     "settings",
					"classified": false,
					"authority":  5,
					"version":    1,
					"content": map[string]any{
						"text": "User prefers dark mode",
					},
				},
				{
					"id":         "atom_002",
					"type":       "preference",
					"domain":     "settings",
					"classified": false,
					"authority":  3,
					"version":    1,
				},
			},
			"count": 2,
		},
	})

	result, err := client.Search(ctx, "dark mode", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("Count = %d, want 2", result.Count)
	}
	if len(result.Results) != 2 {
		t.Fatalf("Results length = %d, want 2", len(result.Results))
	}
	if result.Results[0].ID != "atom_001" {
		t.Errorf("Results[0].ID = %q, want 'atom_001'", result.Results[0].ID)
	}
	if result.Results[0].Type != "preference" {
		t.Errorf("Results[0].Type = %q, want 'preference'", result.Results[0].Type)
	}
	if result.Results[0].Content.Text != "User prefers dark mode" {
		t.Errorf("Results[0].Content.Text = %q", result.Results[0].Content.Text)
	}
}

func TestSearchWithOptions(t *testing.T) {
	ctx := context.Background()
	queryKey := "GET /api/v1/cadreen/memory/search?domain=settings&limit=10&query=dark+mode&tag=pref"
	client := newSandboxClient(map[string]any{
		queryKey: map[string]any{
			"results": []map[string]any{
				{
					"id":     "atom_003",
					"type":   "preference",
					"domain": "settings",
				},
			},
			"count": 1,
		},
	})

	result, err := client.Search(ctx, "dark mode", &SearchMemoryOptions{
		Domain: "settings",
		Tag:    "pref",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
}

func TestGetAtom(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/memory/atom_001": map[string]any{
			"id":              "atom_001",
			"type":            "reference",
			"kind":            "doc",
			"classified":      true,
			"domain":          "docs",
			"scope":           "project",
			"authority":       7,
			"version":         3,
			"tags":            []string{"important"},
			"classifications": map[string]string{"level": "internal"},
			"created_at":      "2025-03-01T08:00:00Z",
			"content": map[string]any{
				"text":   "API documentation for v2",
				"source": "internal wiki",
			},
		},
	})

	result, err := client.GetAtom(ctx, "atom_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "atom_001" {
		t.Errorf("ID = %q, want 'atom_001'", result.ID)
	}
	if result.Type != "reference" {
		t.Errorf("Type = %q, want 'reference'", result.Type)
	}
	if result.Kind != "doc" {
		t.Errorf("Kind = %q, want 'doc'", result.Kind)
	}
	if !result.Classified {
		t.Error("Classified should be true")
	}
	if result.Authority != 7 {
		t.Errorf("Authority = %d, want 7", result.Authority)
	}
	if result.Version != 3 {
		t.Errorf("Version = %d, want 3", result.Version)
	}
	if len(result.Classifications) != 1 {
		t.Errorf("Classifications length = %d, want 1", len(result.Classifications))
	}
	if result.Classifications["level"] != "internal" {
		t.Errorf("Classifications[level] = %q, want 'internal'", result.Classifications["level"])
	}
	if result.Content.Text != "API documentation for v2" {
		t.Errorf("Content.Text = %q", result.Content.Text)
	}
}

func TestProfile(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/memory/profile/user_001": map[string]any{
			"user_id":     "user_001",
			"total_atoms": 42,
			"domains": map[string]int{
				"settings": 15,
				"docs":     27,
			},
			"atoms": []map[string]any{
				{
					"id":     "atom_001",
					"type":   "preference",
					"domain": "settings",
				},
			},
		},
	})

	result, err := client.Profile(ctx, "user_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UserID != "user_001" {
		t.Errorf("UserID = %q, want 'user_001'", result.UserID)
	}
	if result.TotalAtoms != 42 {
		t.Errorf("TotalAtoms = %d, want 42", result.TotalAtoms)
	}
	if len(result.Domains) != 2 {
		t.Errorf("Domains length = %d, want 2", len(result.Domains))
	}
	if result.Domains["settings"] != 15 {
		t.Errorf("Domains[settings] = %d, want 15", result.Domains["settings"])
	}
	if result.Domains["docs"] != 27 {
		t.Errorf("Domains[docs] = %d, want 27", result.Domains["docs"])
	}
	if len(result.Atoms) != 1 {
		t.Fatalf("Atoms length = %d, want 1", len(result.Atoms))
	}
	if result.Atoms[0].ID != "atom_001" {
		t.Errorf("Atoms[0].ID = %q, want 'atom_001'", result.Atoms[0].ID)
	}
}

func TestMemoryTypes(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{
		"GET /api/v1/cadreen/memory/types": map[string]any{
			"type_values": []string{"reference", "preference", "episode", "precedent", "note", "project"},
			"kind_values": []string{"user_setting", "doc", "conversation_log"},
			"description": "Available memory types and their sub-kind values.",
		},
	})

	result, err := client.MemoryTypes(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.TypeValues) != 6 {
		t.Errorf("TypeValues length = %d, want 6", len(result.TypeValues))
	}
	if len(result.KindValues) != 3 {
		t.Errorf("KindValues length = %d, want 3", len(result.KindValues))
	}
	if !strings.Contains(result.Description, "memory types") {
		t.Errorf("Description = %q", result.Description)
	}
}

func TestTeach_MissingFixture(t *testing.T) {
	ctx := context.Background()
	client := newSandboxClient(map[string]any{})

	_, err := client.Teach(ctx, RememberRequest{
		Type:   "reference",
		Domain: "test",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "No fixture") {
		t.Errorf("error should contain 'No fixture', got %q", err.Error())
	}
}
