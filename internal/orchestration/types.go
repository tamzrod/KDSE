package orchestration

import "time"

// OrchestrationState represents the current state of the orchestration engine
type OrchestrationState struct {
	SessionID     string          `json:"session_id"`
	CurrentPhase  OrchestrationPhase `json:"current_phase"`
	PreviousPhase OrchestrationPhase `json:"previous_phase,omitempty"`
	Workspace     WorkspaceInfo   `json:"workspace"`
	Confidence    ConfidenceLevel  `json:"confidence"`
	EvidenceState EvidenceState   `json:"evidence_state"`
	Metrics       OrchestrationMetrics `json:"metrics"`
	History       []PhaseTransition `json:"history,omitempty"`
	Blocked       BlockedState     `json:"blocked,omitempty"`
	UpdatedAt     string           `json:"updated_at"`
}

// WorkspaceInfo captures the current workspace resolution
type WorkspaceInfo struct {
	ResolvedPath    string          `json:"resolved_path"`
	WorkspaceType   WorkspaceType   `json:"workspace_type"`
	RepositoryPath string          `json:"repository_path,omitempty"`
	ProjectPath    string          `json:"project_path,omitempty"`
	TempPath       string          `json:"temp_path,omitempty"`
	KDSEPath       string          `json:"kdse_path"`
}

// WorkspaceType represents the type of workspace resolution
type WorkspaceType string

const (
	WorkspaceTypeRepository WorkspaceType = "repository"
	WorkspaceTypeProject    WorkspaceType = "project"
	WorkspaceTypeTemporary  WorkspaceType = "temporary"
	WorkspaceTypeUnknown    WorkspaceType = "unknown"
)

// OrchestrationPhase represents phases in the state-based orchestration
// Canonical phases defined by MCP orchestration
type OrchestrationPhase string

const (
        PhaseIdle           OrchestrationPhase = "Idle"
        PhaseProblem        OrchestrationPhase = "Problem"
        PhaseKnowledge      OrchestrationPhase = "Knowledge Collection"
        PhaseFoundation     OrchestrationPhase = "Foundation"
        PhaseAudit          OrchestrationPhase = "Audit"
        PhaseAssessment     OrchestrationPhase = "Assessment"
        PhaseArchitecture   OrchestrationPhase = "Architecture"
        PhaseImplementation OrchestrationPhase = "Implementation"
        PhaseComplete       OrchestrationPhase = "Complete"
        PhaseBlocked        OrchestrationPhase = "Blocked"
)

// PhaseTransitions defines valid transitions between phases
var PhaseTransitions = map[OrchestrationPhase][]OrchestrationPhase{
        PhaseIdle:           {PhaseProblem},
        PhaseProblem:        {PhaseKnowledge},
        PhaseKnowledge:      {PhaseFoundation},
        PhaseFoundation:     {PhaseAudit},
        PhaseAudit:          {PhaseAssessment, PhaseArchitecture},
        PhaseAssessment:     {PhaseArchitecture, PhaseFoundation},
        PhaseArchitecture:   {PhaseImplementation},
        PhaseImplementation: {PhaseComplete},
}

// PhaseConfidenceThreshold defines minimum confidence for each phase
var PhaseConfidenceThreshold = map[OrchestrationPhase]float64{
        PhaseIdle:           0.0,
        PhaseProblem:        0.6,
        PhaseKnowledge:      0.7,
        PhaseFoundation:     0.75,
        PhaseAudit:          0.8,
        PhaseAssessment:     0.8,
        PhaseArchitecture:   0.85,
        PhaseImplementation: 0.9,
}

// ConfidenceLevel represents the confidence assessment
type ConfidenceLevel struct {
	Overall      float64          `json:"overall"`
	Foundation   float64          `json:"foundation"`
	Repository   float64          `json:"repository"`
	Evidence     float64          `json:"evidence"`
	Threshold    float64          `json:"threshold"`
	MeetsThreshold bool           `json:"meets_threshold"`
	Dimensions   map[string]float64 `json:"dimensions,omitempty"`
	Assessments  []ConfidenceAssessment `json:"assessments,omitempty"`
}

// ConfidenceAssessment is a single confidence evaluation
type ConfidenceAssessment struct {
	Type      string  `json:"type"`
	Score     float64 `json:"score"`
	Weight    float64 `json:"weight"`
	Details   string  `json:"details"`
	Satisfied bool    `json:"satisfied"`
}

// EvidenceState represents the current evidence evaluation
type EvidenceState struct {
	Required  []EvidenceRequirement `json:"required"`
	Present   []string              `json:"present"`
	Missing   []string              `json:"missing"`
	TotalRequired int               `json:"total_required"`
	TotalPresent int                `json:"total_present"`
	Completeness float64            `json:"completeness"`
}

// EvidenceRequirement defines evidence needed for a phase
type EvidenceRequirement struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Paths       []string `json:"paths"`
	Weight      float64  `json:"weight"`
	Critical    bool     `json:"critical"`
	Satisfied   bool     `json:"satisfied"`
}

// OrchestrationMetrics tracks execution metrics
type OrchestrationMetrics struct {
	CycleCount      int       `json:"cycle_count"`
	PhaseExecutions map[string]int `json:"phase_executions"`
	StartTime       time.Time `json:"start_time"`
	LastCycleTime   time.Time `json:"last_cycle_time"`
	CycleDuration   float64   `json:"cycle_duration_ms"`
}

// PhaseTransition records a phase change
type PhaseTransition struct {
	From        OrchestrationPhase `json:"from"`
	To          OrchestrationPhase `json:"to"`
	Reason      string             `json:"reason"`
	Confidence  float64           `json:"confidence_at_transition"`
	Timestamp   time.Time          `json:"timestamp"`
}

// BlockedState indicates why execution is blocked
type BlockedState struct {
	Blocked    bool     `json:"blocked"`
	Reasons    []string `json:"reasons"`
	Required   []string `json:"required_for_unblock"`
	CanRetry   bool     `json:"can_retry"`
}

// PhaseDecision represents the decision for the next phase
type PhaseDecision struct {
	NextPhase    OrchestrationPhase `json:"next_phase"`
	Reason       string             `json:"reason"`
	Confidence   float64            `json:"confidence"`
	ShouldExecute bool              `json:"should_execute"`
	BlockingReasons []string        `json:"blocking_reasons,omitempty"`
}

// ExecuteCycleResult is the result of a single execute cycle
type ExecuteCycleResult struct {
	CycleNumber    int               `json:"cycle_number"`
	PhaseExecuted  OrchestrationPhase `json:"phase_executed"`
	State          *OrchestrationState `json:"state"`
	Decision       *PhaseDecision     `json:"decision"`
	Success        bool              `json:"success"`
	Error          string            `json:"error,omitempty"`
	Continue       bool              `json:"continue"`
}

// EngineConfig contains configuration for the orchestration engine
type EngineConfig struct {
	FoundationThreshold float64 `json:"foundation_threshold"`
	EvidenceThreshold   float64 `json:"evidence_threshold"`
	MaxCycles          int     `json:"max_cycles"`
	TempWorkspaceBase   string  `json:"temp_workspace_base"`
	EnableMigration     bool    `json:"enable_migration"`
}

// DefaultEngineConfig returns the default engine configuration
func DefaultEngineConfig() *EngineConfig {
	return &EngineConfig{
		FoundationThreshold: 0.7,
		EvidenceThreshold:   0.6,
		MaxCycles:          100,
		TempWorkspaceBase:   "temp",
		EnableMigration:     true,
	}
}
