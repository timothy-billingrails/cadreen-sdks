package cadreen

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// recordedRequest captures what the SDK sends over the wire.
type recordedRequest struct {
	Method string
	Path   string
	Body   map[string]json.RawMessage
}

// contractServer creates an httptest server that records requests and returns
// minimal valid responses. Used to verify SDK HTTP shapes match server expectations.
func contractServer(t *testing.T, responseMap map[string]any) (*httptest.Server, *[]recordedRequest) {
	t.Helper()
	var recorded []recordedRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var bodyMap map[string]json.RawMessage
		if len(body) > 0 {
			json.Unmarshal(body, &bodyMap)
		}
		recorded = append(recorded, recordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Body:   bodyMap,
		})

		key := r.Method + " " + r.URL.Path
		resp, ok := responseMap[key]
		if !ok {
			// Try pattern matching for parameterized routes
			for pattern, val := range responseMap {
				if matchRoute(pattern, key) {
					resp = val
					ok = true
					break
				}
			}
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"route not registered","method":"` + r.Method + `","path":"` + r.URL.Path + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	return srv, &recorded
}

// matchRoute does simple pattern matching for parameterized routes like
// "GET /api/v1/cadreen/agents/{id}" matching "GET /api/v1/cadreen/agents/abc123"
func matchRoute(pattern, actual string) bool {
	pParts := splitPath(pattern)
	aParts := splitPath(actual)
	if len(pParts) != len(aParts) {
		return false
	}
	for i, p := range pParts {
		if len(p) > 2 && p[0] == '{' && p[len(p)-1] == '}' {
			continue // wildcard segment
		}
		if p != aParts[i] {
			return false
		}
	}
	return true
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range splitBySlash(path) {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitBySlash(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

// minimalAgent is the minimal valid agent response.
var minimalAgent = map[string]any{
	"id":          "550e8400-e29b-41d4-a716-446655440000",
	"workspaceId": "550e8400-e29b-41d4-a716-446655440001",
	"name":        "test-agent",
	"status":      "draft",
	"health":      "unknown",
}

var minimalAgentList = map[string]any{
	"agents": []map[string]any{minimalAgent},
	"total":  1,
}

var minimalFederationLink = map[string]any{
	"id":                "550e8400-e29b-41d4-a716-446655440002",
	"workspaceId":       "550e8400-e29b-41d4-a716-446655440001",
	"targetWorkspaceId": "550e8400-e29b-41d4-a716-446655440003",
	"status":            "pending_approval",
}

var minimalKnowledge = map[string]any{
	"id":        "550e8400-e29b-41d4-a716-446655440004",
	"factType":  "reference",
	"subject":   "test subject",
	"confidence": 1.0,
}

var minimalPolicy = map[string]any{
	"id":      "550e8400-e29b-41d4-a716-446655440005",
	"name":    "test-policy",
	"enabled": true,
}

var minimalConnection = map[string]any{
	"id":           "550e8400-e29b-41d4-a716-446655440006",
	"agentId":      "550e8400-e29b-41d4-a716-446655440000",
	"agentCardUrl": "https://example.com/.well-known/agent.json",
	"status":       "pending",
}

var minimalMessage = map[string]any{
	"id":      "550e8400-e29b-41d4-a716-446655440007",
	"content": "test response",
}

var minimalExecution = map[string]any{
	"id":     "550e8400-e29b-41d4-a716-446655440008",
	"status": "planning",
}

var minimalResponse = map[string]any{
	"id":     "resp_abc123",
	"object": "response",
	"status": "completed",
	"model":  "gpt-4.1",
	"output": []map[string]any{},
}

func TestContract_Agents(t *testing.T) {
	srv, recorded := contractServer(t, map[string]any{
		"POST /api/v1/cadreen/agents":                minimalAgent,
		"GET /api/v1/cadreen/agents":                 minimalAgentList,
		"GET /api/v1/cadreen/agents/{id}":            minimalAgent,
		"POST /api/v1/cadreen/agents/{id}/deploy":    minimalAgent,
	})
	defer srv.Close()

	client := NewClient(CadreenConfig{APIKey: "test-key", BaseURL: srv.URL})
	ctx := context.Background()

	// Create agent
	t.Run("CreateAgent", func(t *testing.T) {
		agent, err := client.CreateAgent(ctx, CreateAgentRequest{
			Name:        "test-agent",
			Description: "A test agent",
		})
		if err != nil {
			t.Fatalf("CreateAgent: %v", err)
		}
		if agent.Name != "test-agent" {
			t.Errorf("name = %q, want %q", agent.Name, "test-agent")
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Path != "/api/v1/cadreen/agents" {
			t.Errorf("path = %q, want /api/v1/cadreen/agents", r.Path)
		}
	})

	// List agents
	t.Run("ListAgents", func(t *testing.T) {
		_, err := client.ListAgents(ctx, ListAgentsParams{})
		if err != nil {
			t.Fatalf("ListAgents: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "GET" {
			t.Errorf("method = %q, want GET", r.Method)
		}
	})

	// Get agent
	t.Run("GetAgent", func(t *testing.T) {
		_, err := client.GetAgent(ctx, "550e8400-e29b-41d4-a716-446655440000")
		if err != nil {
			t.Fatalf("GetAgent: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "GET" {
			t.Errorf("method = %q, want GET", r.Method)
		}
	})

	// Deploy agent
	t.Run("DeployAgent", func(t *testing.T) {
		_, err := client.DeployAgent(ctx, "550e8400-e29b-41d4-a716-446655440000", DeployAgentRequest{
			ConfigSnapshot: json.RawMessage(`{"model":"gpt-4.1"}`),
			ChangeSummary:  "initial deploy",
		})
		if err != nil {
			t.Fatalf("DeployAgent: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Body == nil {
			t.Fatal("expected body for deploy")
		}
		if _, ok := r.Body["configSnapshot"]; !ok {
			t.Error("expected configSnapshot in body")
		}
	})
}

func TestContract_Knowledge(t *testing.T) {
	srv, recorded := contractServer(t, map[string]any{
		"POST /api/v1/cadreen/agents/{id}/knowledge": minimalKnowledge,
	})
	defer srv.Close()

	client := NewClient(CadreenConfig{APIKey: "test-key", BaseURL: srv.URL})
	ctx := context.Background()

	t.Run("CreateKnowledge", func(t *testing.T) {
		_, err := client.CreateAgentKnowledge(ctx, "550e8400-e29b-41d4-a716-446655440000", CreateAgentKnowledgeRequest{
			FactType:  "reference",
			Subject:   "GDPR Article 17",
			Predicate: "states",
			Object:    "Right to erasure",
		})
		if err != nil {
			t.Fatalf("CreateAgentKnowledge: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Body == nil {
			t.Fatal("expected body")
		}
		// Verify camelCase field names
		if _, ok := r.Body["factType"]; !ok {
			t.Error("expected factType in body (camelCase)")
		}
		if _, ok := r.Body["subject"]; !ok {
			t.Error("expected subject in body")
		}
		if _, ok := r.Body["predicate"]; !ok {
			t.Error("expected predicate in body")
		}
		if _, ok := r.Body["object"]; !ok {
			t.Error("expected object in body")
		}
	})
}

func TestContract_Messaging(t *testing.T) {
	srv, recorded := contractServer(t, map[string]any{
		"POST /api/v1/cadreen/agents/{id}/send": minimalMessage,
	})
	defer srv.Close()

	client := NewClient(CadreenConfig{APIKey: "test-key", BaseURL: srv.URL})
	ctx := context.Background()

	t.Run("SendMessage", func(t *testing.T) {
		_, err := client.SendAgentMessage(ctx, "550e8400-e29b-41d4-a716-446655440000", SendAgentMessageRequest{
			FromAgentID: "550e8400-e29b-41d4-a716-446655440001",
			Content:     "Hello from the test",
		})
		if err != nil {
			t.Fatalf("SendAgentMessage: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Body == nil {
			t.Fatal("expected body")
		}
		if _, ok := r.Body["fromAgentId"]; !ok {
			t.Error("expected fromAgentId in body (camelCase)")
		}
		if _, ok := r.Body["content"]; !ok {
			t.Error("expected content in body")
		}
	})
}

func TestContract_Execution(t *testing.T) {
	srv, recorded := contractServer(t, map[string]any{
		"POST /api/v1/cadreen/agents/{id}/executions": minimalExecution,
	})
	defer srv.Close()

	client := NewClient(CadreenConfig{APIKey: "test-key", BaseURL: srv.URL})
	ctx := context.Background()

	t.Run("CreateExecution", func(t *testing.T) {
		_, err := client.CreateAgentExecution(ctx, "550e8400-e29b-41d4-a716-446655440000", CreateAgentExecutionRequest{
			Intent: "Analyze the quarterly report",
		})
		if err != nil {
			t.Fatalf("CreateAgentExecution: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Body == nil {
			t.Fatal("expected body")
		}
		if _, ok := r.Body["intent"]; !ok {
			t.Error("expected intent in body")
		}
		// Verify old field names are NOT sent
		if _, ok := r.Body["task"]; ok {
			t.Error("should NOT send 'task' (deprecated)")
		}
		if _, ok := r.Body["input"]; ok {
			t.Error("should NOT send 'input' (deprecated)")
		}
	})
}

func TestContract_Federation(t *testing.T) {
	srv, recorded := contractServer(t, map[string]any{
		"POST /api/v1/cadreen/federation": minimalFederationLink,
	})
	defer srv.Close()

	client := NewClient(CadreenConfig{APIKey: "test-key", BaseURL: srv.URL})
	ctx := context.Background()

	t.Run("CreateFederation", func(t *testing.T) {
		targetID := "550e8400-e29b-41d4-a716-446655440003"
		_, err := client.CreateFederation(ctx, CreateFederationRequest{
			TargetWorkspaceID: targetID,
		})
		if err != nil {
			t.Fatalf("CreateFederationLink: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Body == nil {
			t.Fatal("expected body")
		}
		// Verify camelCase field name
		if _, ok := r.Body["targetWorkspaceId"]; !ok {
			t.Error("expected targetWorkspaceId in body (camelCase)")
		}
		// Verify old field name is NOT sent
		if _, ok := r.Body["target_workspace_id"]; ok {
			t.Error("should NOT send target_workspace_id (deprecated)")
		}
	})
}

func TestContract_Responses(t *testing.T) {
	srv, recorded := contractServer(t, map[string]any{
		"POST /api/v1/cadreen/responses": minimalResponse,
	})
	defer srv.Close()

	client := NewClient(CadreenConfig{APIKey: "test-key", BaseURL: srv.URL})
	ctx := context.Background()

	t.Run("CreateResponse", func(t *testing.T) {
		_, err := client.CreateResponse(ctx, ResponseRequest{
			Model: "gpt-4.1",
			Input: "What is the capital of France?",
		})
		if err != nil {
			t.Fatalf("CreateResponse: %v", err)
		}
		r := (*recorded)[len(*recorded)-1]
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Body == nil {
			t.Fatal("expected body")
		}
		if _, ok := r.Body["model"]; !ok {
			t.Error("expected model in body")
		}
		if _, ok := r.Body["input"]; !ok {
			t.Error("expected input in body")
		}
	})
}
