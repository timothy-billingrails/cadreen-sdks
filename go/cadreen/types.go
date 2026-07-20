package cadreen

import (
	"encoding/json"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
	HealthStatusLatent    HealthStatus = "latent"
)

type GovernanceDecisionType string

const (
	GovernanceAuto             GovernanceDecisionType = "auto"
	GovernanceAutoComplete     GovernanceDecisionType = "auto_complete"
	GovernanceHandoff          GovernanceDecisionType = "handoff"
	GovernanceEscalate         GovernanceDecisionType = "escalate"
	GovernanceClarifyRequester GovernanceDecisionType = "clarify_requester"
	GovernanceAbstain          GovernanceDecisionType = "abstain"
)

type RecoveryStatus string

const (
	RecoveryDiagnosing    RecoveryStatus = "diagnosing"
	RecoveryRecovering    RecoveryStatus = "recovering"
	RecoverySubExecution  RecoveryStatus = "sub_execution"
	RecoveryEscalating    RecoveryStatus = "escalating"
	RecoveryRecovered     RecoveryStatus = "recovered"
	RecoveryFailed        RecoveryStatus = "failed"
	RecoverySkipped       RecoveryStatus = "skipped"
)

type IntentMode string

const (
	IntentModeAuto      IntentMode = "auto"
	IntentModeChat      IntentMode = "chat"
	IntentModeExecution IntentMode = "execution"
)

type IntentResultType string

const (
	IntentResultDirect          IntentResultType = "direct"
	IntentResultClarify         IntentResultType = "clarify"
	IntentResultExecution       IntentResultType = "execution"
	IntentResultBlocked         IntentResultType = "blocked"
	IntentResultConnectRequired IntentResultType = "connect_required"
)

type ErrorType string

const (
	ErrorTypeInvalidRequest      ErrorType = "invalid_request"
	ErrorTypeAuthenticationError ErrorType = "authentication_error"
	ErrorTypePermissionError     ErrorType = "permission_error"
	ErrorTypeNotFound            ErrorType = "not_found"
	ErrorTypeConflict            ErrorType = "conflict"
	ErrorTypeValidationError     ErrorType = "validation_error"
	ErrorTypeRateLimit           ErrorType = "rate_limit"
	ErrorTypeInternalError       ErrorType = "internal_error"
	ErrorTypeServiceUnavailable  ErrorType = "service_unavailable"
)

type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

type ConnectionGroup struct {
	Capability string `json:"capability"`
	Status     string `json:"status"`
}

type ListConnectionsResponse struct {
	Connections       []ConnectionGroup `json:"connections"`
	TotalCapabilities int               `json:"total_capabilities"`
	Pagination        Pagination        `json:"pagination,omitempty"`
}

type AtomContent struct {
	Text        string   `json:"text,omitempty"`
	Source      string   `json:"source,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	Constraint  string   `json:"constraint,omitempty"`
	Query       string   `json:"query,omitempty"`
	ToolsUsed   []string `json:"tools_used,omitempty"`
	Outcome     string   `json:"outcome,omitempty"`
	Situation   string   `json:"situation,omitempty"`
	Action      string   `json:"action,omitempty"`
	Result      string   `json:"result,omitempty"`
	Name        string   `json:"name,omitempty"`
	Constraints []string `json:"constraints,omitempty"`
	Deadline    string   `json:"deadline,omitempty"`
	IsPrivate   bool     `json:"is_private,omitempty"`
}

type Atom struct {
	ID              string            `json:"id"`
	Type            string            `json:"type"`
	Kind            string            `json:"kind,omitempty"`
	Classified      bool              `json:"classified"`
	Domain          string            `json:"domain"`
	Scope           string            `json:"scope,omitempty"`
	Content         AtomContent       `json:"content,omitempty"`
	Authority       int               `json:"authority"`
	Version         int               `json:"version"`
	Tags            []string          `json:"tags,omitempty"`
	Classifications map[string]string `json:"classifications,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
}

type CreateMemoryResponse struct {
	ID              string            `json:"id"`
	Type            string            `json:"type"`
	Kind            string            `json:"kind,omitempty"`
	Classified      bool              `json:"classified"`
	Domain          string            `json:"domain"`
	Scope           string            `json:"scope,omitempty"`
	Content         AtomContent       `json:"content,omitempty"`
	Authority       int               `json:"authority"`
	Version         int               `json:"version"`
	Indexed         bool              `json:"indexed,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Classifications map[string]string `json:"classifications,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
}

type SearchMemoryResponse struct {
	Results []Atom `json:"results"`
	Count   int    `json:"count"`
}

type AtomCategory string

const (
	AtomCategoryReference  AtomCategory = "reference"
	AtomCategoryPreference AtomCategory = "preference"
	AtomCategoryEpisode    AtomCategory = "episode"
	AtomCategoryPrecedent  AtomCategory = "precedent"
	AtomCategoryNote       AtomCategory = "note"
	AtomCategoryProject    AtomCategory = "project"
)

func AllAtomCategories() []AtomCategory {
	return []AtomCategory{
		AtomCategoryReference,
		AtomCategoryPreference,
		AtomCategoryEpisode,
		AtomCategoryPrecedent,
		AtomCategoryNote,
		AtomCategoryProject,
	}
}

type MemoryTypesResponse struct {
	TypeValues  []string `json:"type_values"`
	KindValues  []string `json:"kind_values"`
	Description string   `json:"description"`
}

type MemoryProfileResponse struct {
	UserID     string         `json:"user_id"`
	TotalAtoms int            `json:"total_atoms"`
	Domains    map[string]int `json:"domains"`
	Atoms      []Atom         `json:"atoms"`
}

type Policy struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	Priority      int    `json:"priority"`
	RequiresHuman bool   `json:"requires_human"`
	ApproverRole  string `json:"approver_role,omitempty"`
	SLAHours      int    `json:"sla_hours,omitempty"`
	Rationale     string `json:"rationale,omitempty"`
}

type PolicyBundle struct {
	ID        string   `json:"id"`
	Version   int      `json:"version"`
	Name      string   `json:"name"`
	Policies  []Policy `json:"policies"`
	CreatedAt string   `json:"created_at,omitempty"`
}

type GovernanceDecision struct {
	Type       string  `json:"type"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

type EvaluatePolicyResponse struct {
	Action           string             `json:"action"`
	Domain           string             `json:"domain"`
	GovernanceResult GovernanceDecision `json:"result"`
}

type Escalation struct {
	ID           string   `json:"id"`
	Intent       string   `json:"intent,omitempty"`
	Status       string   `json:"status"`
	Category     string   `json:"category,omitempty"`
	ExecutionID  string   `json:"execution_id,omitempty"`
	ToolName     string   `json:"tool_name,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
	Severity     string   `json:"severity,omitempty"`
	HumanPrompt  string   `json:"human_prompt,omitempty"`
	Suggestions  []string `json:"suggestions,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	ResolvedAt   string   `json:"resolved_at,omitempty"`
	ResolvedBy   string   `json:"resolved_by,omitempty"`
	Resolution   string   `json:"resolution,omitempty"`
}

type ListEscalationsResponse struct {
	Escalations []Escalation `json:"escalations"`
	Count       int          `json:"count"`
	Pagination  Pagination   `json:"pagination,omitempty"`
}

type CredentialMetadata struct {
	ID                string `json:"id"`
	Provider          string `json:"provider"`
	CredentialName    string `json:"credential_name"`
	Type              string `json:"type,omitempty"`
	IsActive          bool   `json:"is_active"`
	HasCredentialData bool   `json:"has_credential_data"`
}

type ListCredentialsResponse struct {
	Credentials []CredentialMetadata `json:"credentials"`
	Count       int                  `json:"count"`
	Pagination  Pagination           `json:"pagination,omitempty"`
}

type CapabilityMatch struct {
	Name        string   `json:"name"`
	HumanName   string   `json:"human_name,omitempty"`
	Description string   `json:"description,omitempty"`
	Score       float64  `json:"score,omitempty"`
	MatchedOn   []string `json:"matched_on,omitempty"`
	Health      string   `json:"health,omitempty"`
	Source      string   `json:"source,omitempty"`
	Status      string   `json:"status,omitempty"`
	Functions   []string `json:"functions,omitempty"`
	Category    string   `json:"category,omitempty"`
}

type ListCapabilitiesResponse struct {
	Available  []CapabilityMatch `json:"available"`
	Gaps       []Gap             `json:"gaps,omitempty"`
	Count      int               `json:"count"`
	Pagination Pagination        `json:"pagination,omitempty"`
}

type ListPoliciesResponse struct {
	Policies   []Policy   `json:"policies"`
	Version    int        `json:"version,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type Outcome struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	Ready       bool     `json:"ready"`
	BlockedBy   []string `json:"blocked_by,omitempty"`
}

type Gap struct {
	Capability  string `json:"capability"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
	Blocking    bool   `json:"blocking"`
	Severity    string `json:"severity"`
	Source      string `json:"source,omitempty"`
}

type StackItem struct {
	Name        string   `json:"name"`
	Type        string   `json:"type,omitempty"`
	Source      string   `json:"source,omitempty"`
	Status      string   `json:"status,omitempty"`
	Description string   `json:"description,omitempty"`
	Contains    []string `json:"contains,omitempty"`
	Functions   []string `json:"functions,omitempty"`
}

type StackBreakdown struct {
	UserData   []StackItem `json:"user_data,omitempty"`
	Cadreen    []StackItem `json:"cadreen,omitempty"`
	Connectors []StackItem `json:"connectors,omitempty"`
	Gaps       []StackItem `json:"gaps,omitempty"`
}

type PolicyRecommendation struct {
	Policy   string `json:"policy"`
	Reason   string `json:"reason"`
	Action   string `json:"action"`
	Blocking bool   `json:"blocking"`
}

type Assessment struct {
	Task                string                 `json:"task"`
	Capabilities        []CapabilityMatch      `json:"capabilities,omitempty"`
	Gaps                []Gap                  `json:"gaps,omitempty"`
	GapFillingTasks     []any          `json:"gap_filling_tasks,omitempty"`
	BlockingGaps        int                    `json:"blocking_gaps,omitempty"`
	PoliciesRecommended []PolicyRecommendation `json:"policies_recommended,omitempty"`
	NeedsClarification  []string               `json:"needs_clarification,omitempty"`
	CanDo               float64                `json:"can_do"`
	AssessmentQuality   string                 `json:"assessment_quality"`
	ReadyCapabilities   int                    `json:"ready_capabilities"`
	TotalCapabilities   int                    `json:"total_capabilities"`
	GapCount            int                    `json:"gap_count"`
	ReadyForDeployment  bool                   `json:"ready_for_deployment"`
	Stack               StackBreakdown         `json:"stack,omitempty"`
	GovernanceDecision  *GovernanceDecision    `json:"governance_decision,omitempty"`
	Outcomes            []Outcome              `json:"outcomes,omitempty"`
	Intelligence        *IntelligenceMeta      `json:"intelligence,omitempty"`
}

type CapabilityTrace struct {
	TotalAvailable     int      `json:"total_available"`
	HealthyCount       int      `json:"healthy_count"`
	ActiveIntegrations []string `json:"active_integrations,omitempty"`
}

type ReasoningTrace struct {
	CapabilityMatches int `json:"capability_matches,omitempty"`
}

type MemoryTrace struct {
	Healthy          bool `json:"healthy"`
	KnowledgeQueried int  `json:"knowledge_queried,omitempty"`
}

type GovernanceTrace struct {
	Active      bool         `json:"active"`
	Decision    string       `json:"decision,omitempty"`
	Confidence  float64      `json:"confidence,omitempty"`
	ReasonCode  string       `json:"reason_code,omitempty"`
	PolicyID    string       `json:"policy_id,omitempty"`
	NextActions []NextAction `json:"next_actions,omitempty"`
}

type HumilityTrace struct {
	GapsDetected int `json:"gaps_detected,omitempty"`
	Blocking     int `json:"blocking,omitempty"`
}

type ProcessTrace struct {
	StartedAt  string          `json:"started_at"`
	DurationMs int64           `json:"duration_ms"`
	Components map[string]bool `json:"components,omitempty"`
}

type NextAction struct {
	Type     string `json:"type"`
	Label    string `json:"label"`
	Endpoint string `json:"endpoint,omitempty"`
	Reason   string `json:"reason"`
}

type IntelligenceStage struct {
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	DurationMs int64                  `json:"duration_ms,omitempty"`
	Detail     string                 `json:"detail,omitempty"`
	Inputs     map[string]any `json:"inputs,omitempty"`
	Outputs    map[string]any `json:"outputs,omitempty"`
}

type FieldStability struct {
	Stable   []string `json:"stable"`
	Evolving []string `json:"evolving"`
	Internal []string `json:"internal"`
}

type IntelligenceMeta struct {
	Version        string              `json:"version"`
	Summary        string              `json:"summary,omitempty"`
	Capability     CapabilityTrace     `json:"capability"`
	Reasoning      ReasoningTrace      `json:"reasoning"`
	Memory         MemoryTrace         `json:"memory"`
	Governance     GovernanceTrace     `json:"governance"`
	Humility       HumilityTrace       `json:"humility"`
	Process        ProcessTrace        `json:"process"`
	NextAction     *NextAction         `json:"next_action,omitempty"`
	Stages         []IntelligenceStage `json:"stages,omitempty"`
	FieldStability FieldStability      `json:"field_stability"`
}

type IntentMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type IntentContext struct {
	ExistingConnectors []string               `json:"existing_connectors,omitempty"`
	Constraints        map[string]any `json:"constraints,omitempty"`
	Domain             string                 `json:"domain,omitempty"`
}

type IntentRequest struct {
	Messages       []IntentMessage `json:"messages"`
	ConversationID string          `json:"conversation_id,omitempty"`
	Context        *IntentContext  `json:"context,omitempty"`
	Mode           string          `json:"mode,omitempty"`
	Stream         bool            `json:"stream,omitempty"`
	DryRun         bool            `json:"dry_run,omitempty"`
	UserID         string          `json:"user_id,omitempty"`
}

type ResponseMessage struct {
	Role     string                 `json:"role"`
	Content  string                 `json:"content"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type ResponseExecution struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	StreamURL string `json:"stream_url,omitempty"`
	PollURL   string `json:"poll_url,omitempty"`
}

type ResponseClarificationQuestion struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Reason   string `json:"reason,omitempty"`
}

type ClarificationQuestions []ResponseClarificationQuestion

func (q *ClarificationQuestions) UnmarshalJSON(data []byte) error {
	var objects []ResponseClarificationQuestion
	if err := json.Unmarshal(data, &objects); err == nil {
		*q = objects
		return nil
	}
	var strings []string
	if err := json.Unmarshal(data, &strings); err != nil {
		return err
	}
	result := make([]ResponseClarificationQuestion, len(strings))
	for i, s := range strings {
		result[i] = ResponseClarificationQuestion{Question: s, Type: "open", Required: false}
	}
	*q = result
	return nil
}

type IntentResult struct {
	Type           IntentResultType              `json:"type"`
	Status         string                        `json:"status"`
	NextAction     *NextAction                   `json:"next_action,omitempty"`
	TraceID        string                        `json:"trace_id"`
	Message        *ResponseMessage              `json:"message,omitempty"`
	Questions      []ResponseClarificationQuestion `json:"questions,omitempty"`
	ConversationID string                        `json:"conversation_id,omitempty"`
	Execution      *ResponseExecution            `json:"execution,omitempty"`
	ReasonCode     string                        `json:"reason_code,omitempty"`
	PolicyID       string                        `json:"policy_id,omitempty"`
	Intelligence   *IntelligenceMeta             `json:"intelligence,omitempty"`
}

type RegisterOpenAPIRequest struct {
	Name         string `json:"name"`
	SpecURL      string `json:"spec_url,omitempty"`
	SpecContent  string `json:"spec_content,omitempty"`
	CredentialID string `json:"credential_id,omitempty"`
}

type RegisterOpenAPIResponse struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	ToolsGenerated  int      `json:"tools_generated,omitempty"`
	ToolsRegistered int      `json:"tools_registered,omitempty"`
	Functions       []string `json:"functions,omitempty"`
	SpecURL         string   `json:"spec_url,omitempty"`
	Status          string   `json:"status"`
}

type CreatePolicyResponse struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Version              int    `json:"version"`
	Status               string `json:"status"`
	ConfirmationRequired bool   `json:"confirmation_required,omitempty"`
	ApproveURL           string `json:"approve_url,omitempty"`
}

type ConfirmPolicyResponse struct {
	ID              string `json:"id"`
	Version         int    `json:"version"`
	PreviousVersion int    `json:"previous_version,omitempty"`
	Status          string `json:"status"`
	AlreadyActive   bool   `json:"already_active,omitempty"`
	ConfirmedAt     string `json:"confirmed_at,omitempty"`
}

type SetupConnection struct {
	Capability string `json:"capability"`
}

type SetupCredential struct {
	Provider string                 `json:"provider"`
	Name     string                 `json:"name"`
	KeyData  map[string]any `json:"key_data"`
}

type SetupMemory struct {
	Type      string                 `json:"type,omitempty"`
	Content   map[string]any `json:"content"`
	Domain    string                 `json:"domain,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Authority int                    `json:"authority,omitempty"`
}

type SetupPolicy struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Rule        string `json:"rule"`
	Severity    string `json:"severity,omitempty"`
}

type SetupRequest struct {
	WorkspaceID string            `json:"workspace_id,omitempty"`
	Purpose     string            `json:"purpose,omitempty"`
	Examples    []string          `json:"examples,omitempty"`
	Constraints []string          `json:"constraints,omitempty"`
	Connections []SetupConnection `json:"connections,omitempty"`
	Credentials []SetupCredential `json:"credentials,omitempty"`
	Memory      []SetupMemory     `json:"memory,omitempty"`
	Policies    []SetupPolicy     `json:"policies,omitempty"`
	Confirm     bool              `json:"confirm,omitempty"`
	DryRun      bool              `json:"dry_run,omitempty"`
}

type SetupProposal struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
}

type SetupConnectionResult struct {
	Capability string          `json:"capability"`
	Status     string          `json:"status"`
	Detail     json.RawMessage `json:"detail,omitempty"`
	Error      string          `json:"error,omitempty"`
}

type SetupCredentialResult struct {
	Provider string `json:"provider"`
	Name     string `json:"name"`
	ID       string `json:"id,omitempty"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

type SetupMemoryResult struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Kind       string `json:"kind,omitempty"`
	Classified bool   `json:"classified"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
}

type SetupPolicyResult struct {
	Name   string `json:"name"`
	ID     string `json:"id,omitempty"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type SetupResult struct {
	Connections []SetupConnectionResult `json:"connections"`
	Credentials []SetupCredentialResult `json:"credentials"`
	Memory      []SetupMemoryResult     `json:"memory"`
	Policies    []SetupPolicyResult     `json:"policies"`
	Proposals   []SetupProposal         `json:"proposals,omitempty"`
	Notice      string                  `json:"notice,omitempty"`
	Applied     int                     `json:"applied"`
	Failed      int                     `json:"failed"`
	DryRun      bool                    `json:"dry_run,omitempty"`
}

type HealingDiagnosis struct {
	ErrorCategory       string  `json:"error_category,omitempty"`
	SemanticReason      string  `json:"semantic_reason,omitempty"`
	RootCause           string  `json:"root_cause,omitempty"`
	CanRetry            bool    `json:"can_retry,omitempty"`
	NeedsSubExecution   bool    `json:"needs_sub_execution,omitempty"`
	NeedsHuman          bool    `json:"needs_human,omitempty"`
	ShouldSkip          bool    `json:"should_skip,omitempty"`
	NeedsReDecide       bool    `json:"needs_re_decide,omitempty"`
	NeedsTryAlternative bool    `json:"needs_try_alternative,omitempty"`
	RetryDelayMs        int64   `json:"retry_delay_ms,omitempty"`
	Confidence          float64 `json:"confidence,omitempty"`
}

type HealingStatsResponse struct {
	TotalPrecedents      int                    `json:"total_precedents,omitempty"`
	SuccessfulRecoveries int                    `json:"successful_recoveries,omitempty"`
	FailedRecoveries     int                    `json:"failed_recoveries,omitempty"`
	SuccessRate          float64                `json:"success_rate,omitempty"`
	AvgDurationMs        int                    `json:"avg_duration_ms,omitempty"`
	CommonStrategies     []StrategyCount        `json:"common_strategies,omitempty"`
	TopTools             []ToolHealingStats     `json:"top_tools,omitempty"`
	ByCategory           map[string]any `json:"by_category,omitempty"`
	TimeRange            TimeRange              `json:"time_range,omitempty"`
}

type StrategyCount struct {
	Strategy string `json:"strategy"`
	Count    int    `json:"count"`
}

type ToolHealingStats struct {
	ToolName    string  `json:"tool_name"`
	Total       int     `json:"total"`
	Successful  int     `json:"successful"`
	Failed      int     `json:"failed"`
	SuccessRate float64 `json:"success_rate"`
	TopStrategy string  `json:"top_strategy,omitempty"`
}

type TimeRange struct {
	FirstPrecedent string `json:"first_precedent,omitempty"`
	LastPrecedent  string `json:"last_precedent,omitempty"`
}

type HealingPrecedent struct {
	ID               string   `json:"id"`
	ToolName         string   `json:"tool_name,omitempty"`
	ErrorType        string   `json:"error_type"`
	ErrorCategory    string   `json:"error_category,omitempty"`
	SemanticReason   string   `json:"semantic_reason,omitempty"`
	RootCause        string   `json:"root_cause,omitempty"`
	RecoveryStrategy string   `json:"recovery_strategy,omitempty"`
	Success          bool     `json:"success"`
	WhatWorked       string   `json:"what_worked,omitempty"`
	WhatFailed       string   `json:"what_failed,omitempty"`
	Attempts         int      `json:"attempts"`
	DurationMs       int64    `json:"duration_ms,omitempty"`
	Confidence       float64  `json:"confidence"`
	CreatedAt        string   `json:"created_at,omitempty"`
	Domain           string   `json:"domain,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

type IntelligenceTraceEntry struct {
	ID          string            `json:"id"`
	Domain      string            `json:"domain"`
	RequestPath string            `json:"request_path"`
	Method      string            `json:"request_method"`
	Meta        IntelligenceMeta  `json:"meta"`
	CreatedAt   string            `json:"created_at,omitempty"`
}

type ListIntelligenceResponse struct {
	Traces     []IntelligenceTraceEntry `json:"traces"`
	Count      int                      `json:"count"`
	Pagination Pagination               `json:"pagination,omitempty"`
}

type IntelligenceStats struct {
	Traces24h             int                `json:"traces_24h"`
	Traces7d              int                `json:"traces_7d"`
	Traces30d             int                `json:"traces_30d"`
	AvgConfidenceByDomain map[string]float64 `json:"avg_confidence_by_domain,omitempty"`
	GapDetectionRate      float64            `json:"gap_detection_rate"`
	GovernanceDecisions   map[string]int     `json:"governance_decisions,omitempty"`
}

type AuditAction struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	Name       string  `json:"name,omitempty"`
	Timestamp  string  `json:"timestamp,omitempty"`
	Status     string  `json:"status,omitempty"`
	Details    string  `json:"details,omitempty"`
	CostUSD    float64 `json:"cost_usd,omitempty"`
	DurationMs int     `json:"duration_ms,omitempty"`
}

type ListAuditResponse struct {
	Actions    []AuditAction `json:"actions"`
	Count      int           `json:"count"`
	Pagination Pagination    `json:"pagination,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorInner `json:"error"`
}

type ErrorInner struct {
	Type       string        `json:"type"`
	Code       string        `json:"code"`
	Message    string        `json:"message"`
	Hint       string        `json:"hint,omitempty"`
	NextAction *NextAction   `json:"next_action,omitempty"`
	Details    []ErrorDetail `json:"details,omitempty"`
	DocURL     string        `json:"doc_url,omitempty"`
	Param      string        `json:"param,omitempty"`
	RetryAfter int           `json:"retry_after,omitempty"`
}

type APIError struct {
	StatusCode int
	Type       string
	Code       string
	Message    string
	Hint       string
	NextAction *NextAction
	Details    []ErrorDetail
	RetryAfter int
}

func (e *APIError) Error() string {
	return e.Message
}

// ---------------------------------------------------------------------------
// Abstraction rename: Atom → MemoryItem
// ---------------------------------------------------------------------------

// MemoryItemContent is the content of a memory item.
type MemoryItemContent = AtomContent

// MemoryType is a user-writable memory type category.
type MemoryType = AtomCategory

const (
	MemoryTypeReference  = AtomCategoryReference
	MemoryTypePreference = AtomCategoryPreference
	MemoryTypeEpisode    = AtomCategoryEpisode
	MemoryTypePrecedent  = AtomCategoryPrecedent
	MemoryTypeNote       = AtomCategoryNote
	MemoryTypeProject    = AtomCategoryProject
)

func AllMemoryTypes() []AtomCategory { return AllAtomCategories() }

// ---------------------------------------------------------------------------
// Marketplace catalog types
// ---------------------------------------------------------------------------

type CatalogIntegration struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Capabilities []string `json:"capabilities,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Provider     string   `json:"provider"`
	Status       string   `json:"status"`
	AuthType     string   `json:"auth_type"`
	InstallTime  string   `json:"install_time"`
	Popularity   int      `json:"popularity"`
	Featured     bool     `json:"featured"`
}

type CatalogCategory struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Integrations []CatalogIntegration `json:"integrations"`
}

type CatalogResponse struct {
	Categories     []CatalogCategory `json:"categories"`
	Installed      []string          `json:"installed"`
	TotalAvailable int               `json:"total_available"`
}

type InstallResponse struct {
	Status        string `json:"status"`
	AuthURL       string `json:"auth_url,omitempty"`
	Provider      string `json:"provider"`
	EstimatedTime string `json:"estimated_time,omitempty"`
}

// ---------------------------------------------------------------------------
// MCP registration types
// ---------------------------------------------------------------------------

type RegisterMCPRequest struct {
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Transport string            `json:"transport,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type RegisterMCPResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Transport string `json:"transport,omitempty"`
	URL       string `json:"url,omitempty"`
}

// ---------------------------------------------------------------------------
// Connect result (discriminated union)
// ---------------------------------------------------------------------------

type ConnectResult struct {
	Type       string         `json:"type"`
	Capability string         `json:"capability,omitempty"`
	Detail     map[string]any `json:"detail,omitempty"`
}

// ---------------------------------------------------------------------------
// Document types
// ---------------------------------------------------------------------------

type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type,omitempty"`
	Size        int64  `json:"size,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type ListDocumentsResponse struct {
	Documents  []Document `json:"documents"`
	Count      int        `json:"count"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type UploadDocumentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type,omitempty"`
	Size        int64  `json:"size,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ---------------------------------------------------------------------------
// Webhook types
// ---------------------------------------------------------------------------

type Webhook struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Events    []string `json:"events,omitempty"`
	Secret    string   `json:"secret,omitempty"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at,omitempty"`
}

type ListWebhooksResponse struct {
	Webhooks   []Webhook  `json:"webhooks"`
	Count      int        `json:"count"`
	Pagination Pagination `json:"pagination,omitempty"`
}

// ---------------------------------------------------------------------------
// Proposal types
// ---------------------------------------------------------------------------

type ProposalEvidence struct {
	Description string  `json:"description"`
	Source      string  `json:"source,omitempty"`
	Count       int     `json:"count,omitempty"`
	Confidence  float64 `json:"confidence,omitempty"`
}

type TaskProposal struct {
	ID              string             `json:"id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Intent          string             `json:"intent"`
	Domain          string             `json:"domain,omitempty"`
	ProposalType    string             `json:"proposal_type"`
	MissionIntent   string             `json:"mission_intent,omitempty"`
	TriggerType     string             `json:"trigger_type"`
	TriggerSource   string             `json:"trigger_source"`
	TriggerDetails  string             `json:"trigger_details,omitempty"`
	Evidence        []ProposalEvidence `json:"evidence,omitempty"`
	Confidence      float64            `json:"confidence"`
	Priority        int                `json:"priority"`
	Status          string             `json:"status"`
	CreatedAt       string             `json:"created_at"`
	ExpiresAt       string             `json:"expires_at,omitempty"`
	AcceptedAt      string             `json:"accepted_at,omitempty"`
	DismissedAt     string             `json:"dismissed_at,omitempty"`
	DismissalReason string             `json:"dismissal_reason,omitempty"`
	ExecutionID     string             `json:"execution_id,omitempty"`
	DedupKey        string             `json:"dedup_key,omitempty"`
	RequiresReview  bool               `json:"requires_review,omitempty"`
}

// ---------------------------------------------------------------------------
// Learning types
// ---------------------------------------------------------------------------

type LearningPattern struct {
	ID          string   `json:"id"`
	Pattern     string   `json:"pattern"`
	Confidence  float64  `json:"confidence"`
	Occurrences int      `json:"occurrences,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
}

type LearningEpisode struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Outcome     string `json:"outcome,omitempty"`
	TraceID     string `json:"trace_id,omitempty"`
	Domain      string `json:"domain,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type LearningSuggestion struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Impact      string `json:"impact,omitempty"`
	Domain      string `json:"domain,omitempty"`
}

type ListLearningPatternsResponse struct {
	Patterns   []LearningPattern `json:"patterns"`
	Count      int               `json:"count"`
	Pagination Pagination        `json:"pagination,omitempty"`
}

type ListLearningEpisodesResponse struct {
	Episodes   []LearningEpisode `json:"episodes"`
	Count      int               `json:"count"`
	Pagination Pagination        `json:"pagination,omitempty"`
}

type ListLearningSuggestionsResponse struct {
	Suggestions []LearningSuggestion `json:"suggestions"`
	Count       int                  `json:"count"`
	Pagination  Pagination           `json:"pagination,omitempty"`
}

// ── Workspace Users ──

type WorkspaceRole string

const (
	WorkspaceRoleAdmin    WorkspaceRole = "admin"
	WorkspaceRoleOperator WorkspaceRole = "operator"
	WorkspaceRoleMember   WorkspaceRole = "member"
	WorkspaceRoleViewer   WorkspaceRole = "viewer"
)

type WorkspaceUser struct {
	ID          string        `json:"id"`
	WorkspaceID string        `json:"workspace_id"`
	UserID      string        `json:"user_id"`
	Role        WorkspaceRole `json:"role"`
	InvitedBy   string        `json:"invited_by,omitempty"`
	InvitedAt   string        `json:"invited_at"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

type InviteUserRequest struct {
	Email string        `json:"email"`
	Role  WorkspaceRole `json:"role,omitempty"`
}

type UpdateRoleRequest struct {
	Role WorkspaceRole `json:"role"`
}

type ListWorkspaceUsersResponse struct {
	Users []WorkspaceUser `json:"users"`
	Count int             `json:"count"`
}

// ---------------------------------------------------------------------------
// Agent types
// ---------------------------------------------------------------------------

type Agent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Status      string         `json:"status,omitempty"`
	Type        string         `json:"type,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	CreatedAt   string         `json:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
}

type AgentKnowledge struct {
	ID        string   `json:"id"`
	AgentID   string   `json:"agent_id,omitempty"`
	FactType  string   `json:"factType"`
	Subject   string   `json:"subject"`
	Predicate string   `json:"predicate,omitempty"`
	Object    string   `json:"object,omitempty"`
	Source    string   `json:"source,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Visibility string  `json:"visibility,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

type AgentGovernancePolicy struct {
	ID        string           `json:"id"`
	AgentID   string           `json:"agent_id"`
	Name      string           `json:"name"`
	Rules     []map[string]any `json:"rules,omitempty"`
	Priority  int              `json:"priority,omitempty"`
	Active    bool             `json:"active"`
	CreatedAt string           `json:"created_at,omitempty"`
}

type AgentAuditEntry struct {
	ID         string  `json:"id"`
	AgentID    string  `json:"agent_id"`
	Action     string  `json:"action"`
	Details    string  `json:"details,omitempty"`
	Status     string  `json:"status,omitempty"`
	CostUSD    float64 `json:"cost_usd,omitempty"`
	DurationMs int     `json:"duration_ms,omitempty"`
	Timestamp  string  `json:"timestamp,omitempty"`
}

type AgentNegotiation struct {
	ID         string         `json:"id"`
	AgentID    string         `json:"agent_id"`
	Status     string         `json:"status"`
	Initiator  string         `json:"initiator,omitempty"`
	Responder  string         `json:"responder,omitempty"`
	Topic      string         `json:"topic,omitempty"`
	Proposal   map[string]any `json:"proposal,omitempty"`
	Response   map[string]any `json:"response,omitempty"`
	CreatedAt  string         `json:"created_at,omitempty"`
	ResolvedAt string         `json:"resolved_at,omitempty"`
}

type AgentMessage struct {
	ID        string `json:"id"`
	AgentID   string `json:"agent_id"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Content   string `json:"content"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type AgentCapabilitiesResponse struct {
	AgentID      string   `json:"agent_id"`
	Capabilities []string `json:"capabilities"`
}

type AgentConfig struct {
	AgentID            string         `json:"agent_id"`
	Model              string         `json:"model,omitempty"`
	SystemPrompt       string         `json:"system_prompt,omitempty"`
	Temperature        *float64       `json:"temperature,omitempty"`
	MaxTokens          *int           `json:"max_tokens,omitempty"`
	Tools              []string       `json:"tools,omitempty"`
	Connections        []string       `json:"connections,omitempty"`
	MemoryScope        string         `json:"memory_scope,omitempty"`
	GovernancePolicyID string         `json:"governance_policy_id,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

type AgentExecution struct {
	ID          string         `json:"id"`
	AgentID     string         `json:"agent_id"`
	Status      string         `json:"status"`
	Input       map[string]any `json:"input,omitempty"`
	Output      map[string]any `json:"output,omitempty"`
	Error       string         `json:"error,omitempty"`
	StartedAt   string         `json:"started_at"`
	CompletedAt string         `json:"completed_at,omitempty"`
}

type ListAgentsResponse struct {
	Agents     []Agent    `json:"agents"`
	Count      int        `json:"count"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type ListAgentMessagesResponse struct {
	Messages   []AgentMessage `json:"messages"`
	Count      int            `json:"count"`
	Pagination Pagination     `json:"pagination,omitempty"`
}

type ListAgentExecutionsResponse struct {
	Executions []AgentExecution `json:"executions"`
	Count      int              `json:"count"`
	Pagination Pagination       `json:"pagination,omitempty"`
}

type ListAgentKnowledgeResponse struct {
	Knowledge  []AgentKnowledge `json:"knowledge"`
	Count      int              `json:"count"`
	Pagination Pagination       `json:"pagination,omitempty"`
}

type SearchAgentKnowledgeResponse struct {
	Results []AgentKnowledge `json:"results"`
	Count   int              `json:"count"`
}

type ListAgentGovernanceResponse struct {
	Policies []AgentGovernancePolicy `json:"policies"`
	Count    int                     `json:"count"`
}

type ListAgentAuditResponse struct {
	Entries    []AgentAuditEntry `json:"entries"`
	Count      int               `json:"count"`
	Pagination Pagination        `json:"pagination,omitempty"`
}

type ListNegotiationsResponse struct {
	Negotiations []AgentNegotiation `json:"negotiations"`
	Count        int                `json:"count"`
	Pagination   Pagination         `json:"pagination,omitempty"`
}

// ---------------------------------------------------------------------------
// Responses API types
// ---------------------------------------------------------------------------

type ResponseRequest struct {
	Model              string                 `json:"model"`
	Input              interface{}            `json:"input"` // string or []ResponseInputItem
	Instructions       string                 `json:"instructions,omitempty"`
	Tools              []ResponseTool         `json:"tools,omitempty"`
	PreviousResponseID string                 `json:"previous_response_id,omitempty"` // Deprecated: server accepts but ignores; no continuation logic yet
	Store              *bool                  `json:"store,omitempty"`
	Stream             bool                   `json:"stream,omitempty"`
	MaxOutputTokens    int                    `json:"max_output_tokens,omitempty"`
	Temperature        *float64               `json:"temperature,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

type ResponseInputItem struct {
	Type    string      `json:"type"`              // "message", "function_call_output"
	Role    string      `json:"role,omitempty"`    // "user", "assistant", "system"
	Content interface{} `json:"content,omitempty"` // string or []ResponseContentPart
	CallID  string      `json:"call_id,omitempty"` // for function_call_output
	Output  string      `json:"output,omitempty"`  // for function_call_output
}

type ResponseContentPart struct {
	Type     string `json:"type"` // "input_text", "output_text", "image_url"
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type ResponseTool struct {
	Type        string      `json:"type"` // "function"
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
	Strict      *bool       `json:"strict,omitempty"`
}

type Response struct {
	ID                 string                 `json:"id"`
	Object             string                 `json:"object"` // "response"
	CreatedAt          int64                  `json:"created_at"`
	Model              string                 `json:"model"`
	Output             []ResponseOutputItem   `json:"output"`
	OutputText         string                 `json:"output_text,omitempty"`
	Usage              *ResponseUsage         `json:"usage,omitempty"`
	Status             string                 `json:"status"` // "completed", "failed", "in_progress", "cancelled"
	PreviousResponseID string                 `json:"previous_response_id,omitempty"` // Deprecated: server does not persist responses; continuation not implemented
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

type ResponseOutputItem struct {
	ID        string                  `json:"id"`
	Type      string                  `json:"type"`               // "message", "reasoning", "function_call"
	Status    string                  `json:"status,omitempty"`   // "completed", "in_progress", "incomplete"
	Role      string                  `json:"role,omitempty"`     // "assistant"
	Content   []ResponseOutputContent `json:"content,omitempty"`
	Name      string                  `json:"name,omitempty"`      // for function_call
	CallID    string                  `json:"call_id,omitempty"`   // for function_call
	Arguments string                  `json:"arguments,omitempty"` // for function_call
}

type ResponseOutputContent struct {
	Type        string               `json:"type"` // "output_text", "refusal"
	Text        string               `json:"text,omitempty"`
	Annotations []ResponseAnnotation `json:"annotations,omitempty"`
}

type ResponseAnnotation struct {
	Type string `json:"type"` // "url_citation", "file_citation"
	URL  string `json:"url,omitempty"`
}

type ResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type ResponseStreamEvent struct {
	Type         string              `json:"type"`
	Sequence     int                 `json:"sequence,omitempty"`
	Response     *Response           `json:"response,omitempty"`
	Item         *ResponseOutputItem `json:"item,omitempty"`
	OutputIndex  int                 `json:"output_index,omitempty"`
	ContentIndex int                 `json:"content_index,omitempty"`
	Delta        string              `json:"delta,omitempty"`
}

// ---------------------------------------------------------------------------
// Federation types
// ---------------------------------------------------------------------------

type FederationLink struct {
	ID                         string          `json:"id"`
	TargetWorkspaceID          string          `json:"targetWorkspaceId,omitempty"`
	Status                     string          `json:"status,omitempty"`
	Permissions                json.RawMessage `json:"permissions,omitempty"`
	CreatedByUserID            string          `json:"createdByUserId,omitempty"`
	ApprovedByUserID           string          `json:"approvedByUserId,omitempty"`
	SuspendedByUserID          string          `json:"suspendedByUserId,omitempty"`
	SuspensionReason           string          `json:"suspensionReason,omitempty"`
	LastActivityAt             string          `json:"lastActivityAt,omitempty"`
	RevokedAt                  string          `json:"revokedAt,omitempty"`
	RevokeReason               string          `json:"revokeReason,omitempty"`
	SourceWorkspaceName        string          `json:"sourceWorkspaceName,omitempty"`
	SourceWorkspaceSlug        string          `json:"sourceWorkspaceSlug,omitempty"`
	TargetWorkspaceName        string          `json:"targetWorkspaceName,omitempty"`
	CreatedAt                  string          `json:"createdAt,omitempty"`
	UpdatedAt                  string          `json:"updatedAt,omitempty"`
}

type FederationAgent struct {
	ID           string `json:"id"`
	FederationLinkID string `json:"federationLinkId"`
	LocalAgentID     string `json:"localAgentId"`
	RemoteAgentID    string `json:"remoteAgentId"`
	LocalAgentName   string `json:"localAgentName,omitempty"`
	RemoteAgentName  string `json:"remoteAgentName,omitempty"`
	Status           string `json:"status,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
}

type FederationPermissions struct {
	FederationLinkID  string   `json:"federation_link_id"`
	Permissions       []string `json:"permissions"`
	AllowedAgents     []string `json:"allowed_agents,omitempty"`
	AllowedTools      []string `json:"allowed_tools,omitempty"`
	AllowedKnowledge  []string `json:"allowed_knowledge,omitempty"`
	MaxRequestsPerDay *int     `json:"max_requests_per_day,omitempty"`
	UpdatedAt         string   `json:"updated_at"`
}

type ListFederationsResponse struct {
	Federations []FederationLink `json:"federations"`
	Count       int              `json:"count"`
}

type ListFederationAgentsResponse struct {
	Agents []FederationAgent `json:"agents"`
	Count  int               `json:"count"`
}

// ---------------------------------------------------------------------------
// External Agent (A2A) types
// ---------------------------------------------------------------------------

type ExternalAgentSkill struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ExternalAgentCapabilities struct {
	Streaming              bool `json:"streaming"`
	PushNotifications      bool `json:"pushNotifications"`
	StateTransitionHistory bool `json:"stateTransitionHistory"`
}

type ExternalAgentConnection struct {
	ID               string                  `json:"id"`
	AgentID          string                  `json:"agentId"`
	AgentCardURL     string                  `json:"agentCardUrl"`
	AgentName        string                  `json:"agentName"`
	AgentDescription string                  `json:"agentDescription"`
	AgentSystem      string                  `json:"agentSystem"`
	AgentVersion     string                  `json:"agentVersion"`
	AgentCardJSON    map[string]any          `json:"agentCardJson"`
	Skills           []ExternalAgentSkill    `json:"skills"`
	Capabilities     ExternalAgentCapabilities `json:"capabilities"`
	Status           string                    `json:"status"`
	Health           string                  `json:"health"`
	LastUsedAt       string                  `json:"lastUsedAt,omitempty"`
	LastHealthCheckAt string                 `json:"lastHealthCheckAt,omitempty"`
	ErrorMessage     string                  `json:"errorMessage,omitempty"`
	ApprovedBy       string                  `json:"approvedBy,omitempty"`
	ApprovedAt       string                  `json:"approvedAt,omitempty"`
	CreatedAt        string                  `json:"createdAt"`
	UpdatedAt        string                  `json:"updatedAt"`
}

type ExternalAgentInteraction struct {
	ID               string `json:"id"`
	ConnectionID     string `json:"connectionId"`
	AgentID          string `json:"agentId"`
	Direction        string `json:"direction"`
	Operation        string `json:"operation"`
	TaskID           string `json:"taskId,omitempty"`
	Message          string `json:"message,omitempty"`
	Status           string `json:"status"`
	DurationMs       int    `json:"durationMs,omitempty"`
	ErrorMessage     string `json:"errorMessage,omitempty"`
	GovernanceResult string `json:"governanceResult,omitempty"`
	CreatedAt        string `json:"createdAt"`
}

type ExternalAgentSettings struct {
	Enabled bool `json:"enabled"`
}

type ListExternalConnectionsResponse struct {
	Connections []ExternalAgentConnection `json:"connections"`
	Total       int                       `json:"total"`
	Limit       int                       `json:"limit"`
	Offset      int                       `json:"offset"`
}

type ListExternalInteractionsResponse struct {
	Interactions []ExternalAgentInteraction `json:"interactions"`
	Total        int                        `json:"total"`
	Limit        int                        `json:"limit"`
	Offset       int                        `json:"offset"`
}
