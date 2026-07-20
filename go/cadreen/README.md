# cadreen

Go SDK for [Cadreen](https://accomplishanything.today/infra/docs) — Intelligence as a Service.

> **⚠️ Compatibility Advisory:** Version 0.7.0 has known contract mismatches with the server. Version 0.7.0 is unsupported. **Upgrade to 0.7.2.** See [CHANGELOG](#v072) for details.

Cadreen is a cognitive operating system. Send messages describing what you want done, and Cadreen reasons, connects tools, recalls knowledge, governs actions, and escalates to humans when needed. The SDK handles authentication, retries, idempotency, streaming, sandbox mode, and error classification.

## Install

```bash
go get github.com/timothy-billingrails/cadreen-sdks/go/cadreen@latest
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"
)

func main() {
	c := cadreen.NewClient(cadreen.CadreenConfig{
		APIKey: os.Getenv("CADREEN_API_KEY"),
	})

	result, err := c.IntentInvoke(context.Background(), cadreen.IntentRequest{
		Messages: []cadreen.IntentMessage{
			{Role: "user", Content: "Handle refund for invoice inv_123"},
		},
	})
	if err != nil {
		panic(err)
	}

	switch result.Type {
	case cadreen.IntentResultDirect:
		fmt.Println(result.Message.Content)
	case cadreen.IntentResultExecution:
		stream, _ := c.ExecutionStream(context.Background(), result.Execution.ID)
		for event := range stream {
			fmt.Println(event.Type, event.Data)
		}
	case cadreen.IntentResultBlocked:
		fmt.Printf("Blocked by %s: %s\n", result.PolicyID, result.ReasonCode)
	case cadreen.IntentResultClarify:
		for _, q := range result.Questions {
			fmt.Println(q.Question)
		}
	case cadreen.IntentResultConnectRequired:
		fmt.Println("Connect:", result.Connection.Endpoint)
	}
}
```

## Chat Completions

OpenAI-compatible chat completions with built-in governance. Every tool call goes through governance before execution — auto-approved calls execute silently, blocked calls become conversation.

```go
// Basic completion
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "Hello!"},
	},
})
fmt.Println(resp.Choices[0].Message.Content)

// With tool calling
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "Refund order 456"},
	},
	Tools: []cadreen.ChatToolDefinition{{
		Type: "function",
		Function: cadreen.ChatFunctionDefinition{
			Name:        "process_refund",
			Description: "Process a refund for an order",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]any{"order_id": map[string]any{"type": "string"}},
				"required":   []string{"order_id"},
			},
		},
	}},
})

// Check for tool calls
if len(resp.Choices[0].Message.ToolCalls) > 0 {
	for _, tc := range resp.Choices[0].Message.ToolCalls {
		fmt.Printf("%s(%s)\n", tc.Function.Name, tc.Function.Arguments)
	}
}

// Resume a conversation
followUp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "What about order 789?"},
	},
	ConversationID: resp.ConversationID, // use conversation_id from prior response
})
```

### Streaming

```go
ch, err := c.ChatCompletionsStream(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "Hello!"},
	},
})
for event := range ch {
	if event.Error != nil {
		log.Fatal(event.Error)
	}
	if event.Type == "reasoning_delta" {
		// Model is thinking — render in a collapsible accordion
		fmt.Printf("[thinking] %s", event.Reasoning)
	}
	if len(event.Chunk.Choices) > 0 {
		fmt.Print(event.Chunk.Choices[0].Delta.Content)
	}
}
```

### Tool Discovery

```go
tools, err := c.ListTools(ctx)
for _, tool := range tools.Data {
	fmt.Printf("%s: %s\n", tool.Function.Name, tool.Function.Description)
}
```

### Tool Chaining

When the model proposes tool calls, send results back for follow-up:

```go
resp, err := c.ChatCompletions(ctx, cadreen.ChatCompletionRequest{
	Messages: []cadreen.ChatMessage{
		{Role: "user", Content: "What's the weather in NYC?"},
		{Role: "assistant", ToolCalls: []cadreen.ChatToolCall{{
			ID:       "tc_1",
			Type:     "function",
			Function: cadreen.ChatFunctionCall{Name: "get_weather", Arguments: `{"city":"NYC"}`},
		}}},
		{Role: "tool", ToolCallID: "tc_1", Content: `{"temp": 72, "condition": "sunny"}`},
	},
})
// Model may propose more tools or return a final text response
```

## Configuration

```go
c := cadreen.NewClient(cadreen.CadreenConfig{
	APIKey:     "sk_cadreen_...",                    // required
	BaseURL:    "https://accomplishanything.today",  // optional, default shown
	MaxRetries: 3,                                   // optional, default 3
	Timeout:    30 * time.Second,                    // optional, default 30s
	Profile:    "lean",                              // optional: "lean" | "audit" | "full" (default "full")
})
```

### Response Profiles

| Profile | What you get | Use when |
|---------|-------------|----------|
| `"full"` (default) | Full response metadata | You want full transparency |
| `"audit"` | Only governance decision + confidence + blocking gaps | You need to react to gates |
| `"lean"` | No envelope. Just `trace_id` | Hot-looping, minimal payload |

### Sandbox Mode

Test without hitting the API:

```go
c := cadreen.NewClient(cadreen.CadreenConfig{
	Sandbox: true,
	Fixtures: map[string]any{
		"POST /api/v1/cadreen/intent": cadreen.IntentResult{
			Type:    cadreen.IntentResultDirect,
			Status:  "answered",
			TraceID: "sandbox-trace",
			Message: &cadreen.ResponseMessage{Role: "assistant", Content: "It's done."},
		},
	},
})
```

## Marketplace

```go
catalog, _ := c.Catalog(ctx)
for _, category := range catalog.Categories {
	fmt.Printf("%s: %d integrations\n", category.Name, len(category.Integrations))
}

install, _ := c.Install(ctx, "slack")
if install.Status == "pending_auth" {
	fmt.Printf("Authenticate at: %s\n", install.AuthURL)
}
```

## Proposals

Cadreen watches your usage and suggests improvements — actions to automate, schedules to set, rules to relax. You decide what runs.

```go
// List proposals waiting for your decision
proposals, err := c.ListProposals(ctx, cadreen.ListProposalsOptions{
	Status: "proposed",
	Limit:  50,
})
for _, p := range proposals.Proposals {
	fmt.Printf("[%s] %s (%.0f%% confidence)\n", p.Status, p.Title, p.Confidence*100)
}

// Accept — executes via the intent engine
result, err := c.AcceptProposal(ctx, "550e8400-...")
fmt.Printf("Execution: %s, Action: %s\n", result.ExecutionID, result.Action)

// Dismiss — teaches Cadreen what you don't want
err = c.DismissProposal(ctx, "550e8400-...", cadreen.DismissProposalRequest{
	Reason: "We handle this manually",
})

// See counts by status
stats, err := c.ProposalStats(ctx)
fmt.Printf("Waiting: %d, Accepted: %d\n", stats.Proposed, stats.Accepted)
```

## Documents

```go
// List documents
docs, _ := c.ListDocuments(ctx)
for _, d := range docs.Documents {
	fmt.Printf("%s (%s)\n", d.Name, d.ContentType)
}

// Upload a document
result, _ := c.UploadDocument(ctx, "/path/to/report.pdf")
fmt.Printf("Uploaded: %s (ID: %s)\n", result.Name, result.ID)

// Get document details
doc, _ := c.GetDocument(ctx, result.ID)
fmt.Printf("Status: %s, Size: %d bytes\n", doc.Status, doc.Size)
```

## Resources

| Resource | Methods |
|----------|---------|
| `Client.IntentInvoke()` | Intent — primary entry point |
| `Client.ChatCompletions()` | Non-streaming chat completions |
| `Client.ChatCompletionsStream()` | Streaming chat completions (SSE) |
| `Client.ListTools()` | Tool discovery |
| `Client.ExecutionStream()` | Stream execution events |
| `Client.GetExecutionStatus()` | Get execution status |
| `Client.Setup()` | One-time setup |
| `Client.Catalog()` | Browse marketplace |
| `Client.Install()` | Install integration |
| `Client.ListDocuments()` | List uploaded documents |
| `Client.GetDocument()` | Get document by ID |
| `Client.UploadDocument()` | Upload a document from file path |
| `Client.ListEscalations()` | List escalations |
| `Client.GetEscalation()` | Get escalation by ID |
| `Client.ResolveEscalation()` | Resolve an escalation |
| `Client.HealingStats()` | Healing statistics |
| `Client.ListHealingPrecedents()` | List healing precedents |
| `Client.Diagnose()` | Diagnose a failure |
| `Client.CreateWebhook()` | Create webhook subscription |
| `Client.ListWebhooks()` | List webhook subscriptions |
| `Client.DeleteWebhook()` | Delete webhook |
| `Client.ListLearningPatterns()` | Detected patterns |
| `Client.ListLearningEpisodes()` | Learning episodes |
| `Client.ListLearningSuggestions()` | Improvement suggestions |
| `Client.ListCredentials()` | List credentials |
| `Client.CreateCredential()` | Create credential |
| `Client.DeleteCredential()` | Delete credential |
| `Client.ListProposals()` | List task proposals |
| `Client.GetProposal()` | Get proposal by ID |
| `Client.AcceptProposal()` | Accept a proposal |
| `Client.DismissProposal()` | Dismiss a proposal |
| `Client.ProposalStats()` | Proposal statistics |
| `Client.CreateSetupSession()` | Create a setup session |
| `Client.ListSetupSessions()` | List setup sessions |
| `Client.GetSetupSession()` | Get setup session by ID |
| `Client.AddToSetupSession()` | Add resources to a session |
| `Client.ApplySetupSession()` | Apply a session atomically |
| `Client.ListPolicies()` | List policies |
| `Client.GetPolicy()` | Get policy bundle |
| `Client.ListCapabilities()` | List capabilities |
| `Client.Assess()` | Assess task readiness |
| `Client.CreateAgent()` | Create an agent |
| `Client.ListAgents()` | List agents |
| `Client.GetAgent()` | Get agent details |
| `Client.UpdateAgent()` | Update agent |
| `Client.DeleteAgent()` | Delete agent |
| `Client.DeployAgent()` | Deploy agent |
| `Client.GetAgentConfig()` | Get agent config |
| `Client.GetAgentCapabilities()` | Get agent capabilities |
| `Client.SendAgentMessage()` | Send message to agent |
| `Client.ListAgentMessages()` | List agent messages |
| `Client.CreateAgentExecution()` | Create agent execution |
| `Client.ListAgentExecutions()` | List agent executions |
| `Client.CreateAgentKnowledge()` | Add knowledge |
| `Client.SearchAgentKnowledge()` | Search knowledge |
| `Client.CreateAgentGovernance()` | Create governance policy |
| `Client.ListAgentGovernance()` | List governance policies |
| `Client.StartNegotiation()` | Start negotiation |
| `Client.ListNegotiations()` | List negotiations |
| `Client.CreateFederation()` | Create federation link |
| `Client.ListFederations()` | List federation links |
| `Client.ApproveFederation()` | Approve federation |
| `Client.ConnectExternalAgent()` | Connect to external agent |
| `Client.ListExternalConnections()` | List external connections |
| `Client.ApproveExternalConnection()` | Approve external connection |
| `Client.CreateResponse()` | Create response (OpenAI-compatible) |
| `Client.CreateResponseStreaming()` | Create streaming response |
| `Client.ListDevices()` | List devices |
| `Client.GetDevice()` | Get device by ID |
| `Client.CreateDevice()` | Create device |
| `Client.DeleteDevice()` | Delete device |
| `Client.GetDeviceStatus()` | Get device status |
| `Client.DiagnoseDevice()` | Diagnose device |
| `Client.AskDevice()` | Ask device a question |
| `Client.GetDeviceMapStats()` | Get device map stats |
| `Client.ListDeviceTasks()` | List device tasks |
| `Client.CreateDeviceTask()` | Create device task |

## Agents

Agents are autonomous workers that handle tasks, follow rules, and learn from outcomes.

```go
// Create an agent
agent, _ := c.CreateAgent(ctx, cadreen.CreateAgentRequest{
    Name:        "Support Agent",
    Description: "Handles customer support requests",
})
fmt.Printf("Agent created: %s\n", agent.ID)

// List agents
result, _ := c.ListAgents(ctx, cadreen.ListAgentsParams{})
for _, a := range result.Agents {
    fmt.Printf("%s (%s)\n", a.Name, a.Status)
}

// Deploy an agent
_, _ = c.DeployAgent(ctx, "agent_123")

// Send a message
msg, _ := c.SendAgentMessage(ctx, "agent_123", cadreen.SendAgentMessageRequest{
    From:    "agent_456",
    Content: "What's the refund policy?",
})

// Create an execution
exec, _ := c.CreateAgentExecution(ctx, "agent_123", cadreen.CreateAgentExecutionRequest{
    Task: "Process refund for order #1234",
})

// Add knowledge
_, _ = c.CreateAgentKnowledge(ctx, "agent_123", cadreen.CreateAgentKnowledgeRequest{
    Type:    "reference",
    Content: "Refunds require manager approval for amounts over $100",
})

// Search knowledge
results, _ := c.SearchAgentKnowledge(ctx, "agent_123", cadreen.SearchAgentKnowledgeRequest{
    Query: "refund policy",
})

// Create governance policy
_, _ = c.CreateAgentGovernance(ctx, "agent_123", cadreen.CreateAgentGovernanceRequest{
    Name:  "Refund Approval",
    Rules: []map[string]any{{"action": "refund", "condition": "amount > 100", "decision": "handoff"}},
})
```

## Federation

Federation lets workspaces share agents and knowledge across boundaries.

```go
// Create a federation link
link, _ := c.CreateFederation(ctx, cadreen.CreateFederationRequest{
    TargetWorkspaceID: "ws_456",
})
fmt.Printf("Link created: %s (status: %s)\n", link.ID, link.Status)

// List federation links
result, _ := c.ListFederations(ctx)
for _, link := range result.Federations {
    fmt.Printf("%s → %s\n", link.Name, link.Status)
}

// Approve a pending link
_, _ = c.ApproveFederation(ctx, "link_123")
```

## External Agents (A2A)

Connect to agents from other systems (LangChain, CrewAI, etc.) using the A2A protocol.

```go
// Enable external agents
_, _ = c.UpdateExternalAgentSettings(ctx, true)

// Connect to an external agent
conn, _ := c.ConnectExternalAgent(ctx, "agent_123", cadreen.ConnectExternalAgentRequest{
    AgentCardURL: "https://example.com/.well-known/agent.json",
})
fmt.Printf("Connection: %s (status: %s)\n", conn.ID, conn.Status)

// List connections
result, _ := c.ListExternalConnections(ctx, "agent_123", cadreen.ListExternalConnectionsParams{})
for _, c := range result.Connections {
    fmt.Printf("%s → %s (%s)\n", c.AgentName, c.Status, c.Health)
}

// Approve a pending connection
_, _ = c.ApproveExternalConnection(ctx, "agent_123", "conn_456")

// List interactions
interactions, _ := c.ListExternalInteractions(ctx, "agent_123", "conn_456", cadreen.ListExternalInteractionsParams{})
for _, i := range interactions.Interactions {
    fmt.Printf("%s %s: %s\n", i.Direction, i.Operation, i.Status)
}
```

## Responses API

OpenAI-compatible responses API with built-in governance and memory.

```go
// Create a response
response, _ := c.CreateResponse(ctx, cadreen.ResponseRequest{
    Model: "cadreen",
    Input: "What tools do I have?",
})
fmt.Println(response.OutputText)

// Conversation state (server-managed)
response2, _ := c.CreateResponse(ctx, cadreen.ResponseRequest{
    Model:              "cadreen",
    Input:              "What about refund tools?",
    PreviousResponseID: &response.ID,
})

// Streaming
iter, _ := c.CreateResponseStreaming(ctx, cadreen.ResponseRequest{
    Model: "cadreen",
    Input: "Explain quantum computing",
})
for iter.Next() {
    event := iter.Current()
    if event.Type == "response.output_text.delta" {
        fmt.Print(event.Delta)
    }
}
```

## Changelog

### v0.7.2
- Added `devices.go` — full device lifecycle (CreateDevice, ListDevices, GetDevice, DeleteDevice, GetDeviceStatus, UpdateDeviceState, GetDeviceMap, etc.)
- `GetAgentConfig` now returns `*AgentConfig` instead of `map[string]any`
- `CreateAgentExecution` now returns `*AgentExecution` instead of `map[string]any`
- `GetFederationPermissions` and `UpdateFederationPermissions` now return `*FederationPermissions` instead of `map[string]any`
- Added typed structs: `AgentConfig`, `AgentExecution`, `FederationPermissions`

### v0.7.1
- Fix: sandbox mode now returns `ErrSandboxStreaming` for `PostStream` and `Stream` instead of making network requests
- Fix: release workflow Homebrew formula update ordering
- Documented new server-side execution events: `execution_steps_complete`, `mission_completed_with_gaps`

### v0.7.0
- **BREAKING:** Removed `Pathways` and `TotalPathways` from connection responses. `ConnectionGroup` now returns only `Capability` and `Status`.
- **BREAKING:** Removed `Pathway` type. Internal routing details (connector, transport, tool_id) are no longer exposed.
- **BREAKING:** Changed `ConnectManualDetail` from `{Pathways: [...]}` to `{Capability, Available, Health}`.
- **BREAKING:** Removed `WorkspaceID` from response types: `SetupResult`, `SetupSession`, `WebhookSubscription`, `WebhookPayload`. (Still accepted on request types.)
- **BREAKING:** Removed `AuthScheme` from `ExternalAgentConnection` responses.
- **BREAKING:** Removed `AtomsConsulted`, `EpisodesMatched`, `PrecedentsApplied` from memory trace in intelligence metadata.
- **BREAKING:** `SourcesConsulted` renamed to `KnowledgeQueried` in MemoryTrace.
- **BREAKING:** All entity responses (Agent, Knowledge, Governance, Federation, Negotiation, ExternalAgentConnection) no longer include `WorkspaceID`.
- New MCP SSE endpoint: `GET /api/v1/cadreen/mcp/sse` + `POST /api/v1/cadreen/mcp/message`. Connect Cadreen as an MCP server without installing the npm package.
- Added `agents.go` — full agent lifecycle (CreateAgent, ListAgents, GetAgent, UpdateAgent, DeleteAgent, DeployAgent, GetAgentConfig, GetAgentCapabilities)
- Added agent messaging (SendAgentMessage, ListAgentMessages)
- Added agent executions (CreateAgentExecution, ListAgentExecutions)
- Added agent knowledge (CreateAgentKnowledge, ListAgentKnowledge, SearchAgentKnowledge, DeleteAgentKnowledge)
- Added agent governance (CreateAgentGovernance, ListAgentGovernance, UpdateAgentGovernance, DeleteAgentGovernance)
- Added agent audit (ListAgentAudit)
- Added negotiations (StartNegotiation, ListNegotiations, GetNegotiation, RespondToNegotiation)
- Added `federation.go` — cross-workspace federation (CreateFederation, ListFederations, GetFederation, ApproveFederation, SuspendFederation, RevokeFederation)
- Added federation permissions (GetFederationPermissions, UpdateFederationPermissions)
- Added federation agent linking (LinkFederationAgent, ListFederationAgents, UnlinkFederationAgent)
- Added `external_agents.go` — A2A external agent connections (ConnectExternalAgent, ListExternalConnections, GetExternalConnection, ApproveExternalConnection, SuspendExternalConnection, RevokeExternalConnection, DeleteExternalConnection)
- Added external agent interactions (ListExternalInteractions)
- Added external agent settings (GetExternalAgentSettings, UpdateExternalAgentSettings, ListAllExternalConnections)
- Added `responses.go` — OpenAI-compatible responses API (CreateResponse, GetResponse, CreateResponseStreaming)
- Added 31 new types: Agent, AgentKnowledge, AgentGovernancePolicy, AgentAuditEntry, AgentNegotiation, FederationLink, FederationAgent, ExternalAgentConnection, ExternalAgentInteraction, ExternalAgentSettings, ExternalAgentSkill, ExternalAgentCapabilities, etc.

### v0.6.3
- Fix: retry body reuse bug — retries were sending empty body because `bytes.Reader` was consumed on first attempt
- Add: `VerifyWebhookSignature()` function for HMAC-SHA256 webhook signature verification

### v0.6.2
- Added `Intelligence` field to `ChatCompletionResponse` — full intelligence metadata
- Added `ConversationID` field to `ChatCompletionResponse` — conversation continuity
- Added `Reasoning` field to `ChatDelta` — streaming reasoning from thinking models
- Added `ReasoningTokens`, `CacheWriteTokens`, `PromptTokensDetails` to `ChatUsage`
- Added `Reasoning` field to `ChatStreamEvent` — non-empty on reasoning_delta events
- Document `reasoning_delta` streaming event in README
- IntelligenceMeta shape aligned with server

### v0.6.1
- Added `UserID` field to `IntentRequest` and `ChatCompletionRequest` — pass end-user identity for per-user context and memory filtering
- Added `ListWorkspaceUsers()`, `InviteUser()`, `UpdateUserRole()`, `RemoveUser()` — manage workspace team members
- Added `WorkspaceUser`, `WorkspaceRole`, `InviteUserRequest`, `UpdateRoleRequest` types

### v0.6.0 (BREAKING)
- **BREAKING:** `ResolveEscalationRequest.Resolution` renamed to `Decision` — aligns SDK with API contract
- **BREAKING:** `DiagnoseRequest.Error` renamed to `ErrorMessage` — aligns SDK with API contract
- Retry default was already at 3 — no change needed

### v0.5.5
- Added `CreateSetupSession()`, `ListSetupSessions()`, `GetSetupSession()`, `AddToSetupSession()`, `ApplySetupSession()` — stateful setup sessions
- Added `SetupSession`, `SetupSessionCreateRequest`, `SetupSessionAddRequest`, `SetupSessionApplyRequest`, `SetupSessionApplyResult` types

### v0.5.4
- Added `ListProposals()`, `GetProposal()`, `AcceptProposal()`, `DismissProposal()`, `ProposalStats()` — task proposals
- Added `TaskProposal`, `ProposalEvidence`, `ListProposalsResponse`, `AcceptProposalResponse`, `DismissProposalResponse`, `ProposalStatsResponse` types

### v0.5.3
- Added `DryRun` field to `IntentRequest` — preview intent classification, governance, and capability assessment without creating a mission or persisting conversation
- Fixed `Version` constant (was stuck at 0.5.0)

### v0.5.2
- Added `DryRun` field to `SetupRequest` — preview what would be created without persisting
- Added `DryRun` and `Notice` fields to `SetupResult`
- Added `ListBlueprints()`, `GetBlueprint()`, `CreateBlueprint()`, `DeleteBlueprint()`, `RunBlueprint()`, `ListBlueprintRuns()`
- Added `ListSchedules()`, `GetSchedule()`, `CreateSchedule()`, `PauseSchedule()`, `ResumeSchedule()`

### v0.5.1
- Added `UploadDocument()` — upload a document from file path via multipart POST
- Added `ListDocuments()`, `GetDocument()` — document management
- Added `ListEscalations()`, `GetEscalation()`, `ResolveEscalation()` — escalation management
- Added `HealingStats()`, `ListHealingPrecedents()`, `Diagnose()` — self-healing
- Added `CreateWebhook()`, `ListWebhooks()`, `DeleteWebhook()` — webhook subscriptions
- Added `ListLearningPatterns()`, `ListLearningEpisodes()`, `ListLearningSuggestions()` — learning insights
- Added `ListCredentials()`, `CreateCredential()`, `DeleteCredential()` — credential management
- Added `ListPolicies()`, `GetPolicy()` — policy listing (was missing, TS/Python had it)

### v0.5.0 (BREAKING)
- **BREAKING:** API endpoints moved to Cadreen surface (`/api/v1/cadreen/`):
  - `/api/v1/chat/completions` → `/api/v1/cadreen/chat/completions`
  - `/api/v1/tools` → `/api/v1/cadreen/tools`
- All external API calls now route through the Cadreen surface
- Renamed response profile levels for clarity

### v0.4.1
- Added `ChatStreamEvent.RawJSON` — raw JSON bytes before typed parsing; enables extracting `pending_actions`, `conversation_id`, and `intelligence` fields without SDK type changes
- Fixed unused `bytes` import in `executions.go`

### v0.4.0
- Added `ChatCompletions()` — OpenAI-compatible chat completions with governance
- Added `ChatCompletionsStream()` — streaming chat completions via SSE
- Added `ListTools()` — discover available tools as OpenAI function definitions
- Added tool calling support: `Tools` param, `ToolCalls` in responses, tool chaining
- Added `ConversationID` for persistent conversations across requests
- Added `ChatMessage`, `ChatToolDefinition`, `ChatToolCall` and related types
- Added `ExecutionStream()` and `GetExecutionStatus()` for execution monitoring

### v0.3.0
- Added `Catalog()` — browse the unified integration marketplace
- Added `Install()` — one-click install with OAuth flow

### v0.2.0
- Initial public release

## License

MIT
