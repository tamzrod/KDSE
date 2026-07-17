// Package orchestration provides the KDSE orchestration engine that transforms
// KDSE from a toolbox into an orchestration engine. After initialization, the
// LLM never decides which KDSE tool to call - the execute tool decides.
package orchestration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kdse/runtime/internal/enforcer"
	"github.com/kdse/runtime/internal/types"
)

// Phase is an alias for types.OrchestrationPhase for backward compatibility
// The authoritative definition is in github.com/kdse/runtime/internal/types
type Phase = types.OrchestrationPhase

// Re-export Phase constants from types package for convenience
// These are the authoritative definitions
const (
	PhaseIdle           = types.PhaseIdle
	PhaseProblem        = types.PhaseProblem
	PhaseKnowledge      = types.PhaseKnowledge
	PhaseFoundation     = types.PhaseFoundation
	PhaseAudit          = types.PhaseAudit
	PhaseAssessment     = types.PhaseAssessment
	PhaseArchitecture   = types.PhaseArchitecture
	PhaseImplementation = types.PhaseImplementation
	PhaseComplete       = types.PhaseComplete
	PhaseBlocked        = types.PhaseBlocked
)

// Re-export PhaseTransitions and PhaseConfidenceThreshold from types package
var PhaseTransitions = types.PhaseTransitions
var PhaseConfidenceThreshold = types.PhaseConfidenceThreshold

// ExecutionMode represents the mode of execution
type ExecutionMode string

const (
	ModeToolbox ExecutionMode = "toolbox" // Legacy mode - tools available directly
	ModeStrict  ExecutionMode = "strict"   // Strict mode - all requests through execute
)

// PhaseTransitions and PhaseConfidenceThreshold are imported from types package
// See github.com/kdse/runtime/internal/types for authoritative definitions

// SessionState represents the runtime session state for KDSE orchestration
type SessionState struct {
	SessionID         string            `json:"session_id"`
	StartedAt         string            `json:"started_at"`
	UpdatedAt         string            `json:"updated_at"`
	CurrentPhase      types.OrchestrationPhase `json:"current_phase"`
	Confidence        float64           `json:"confidence"`
	Evidence          []string          `json:"evidence"`
	CompletedPhases   []types.OrchestrationPhase `json:"completed_phases"`
	NextAllowedPhases []types.OrchestrationPhase `json:"next_allowed_phases"`
	ExecutionMode     ExecutionMode     `json:"execution_mode"`
	Workspace         *WorkspaceState   `json:"workspace"`
	Objective         string            `json:"objective,omitempty"`
	BlockedReason     string            `json:"blocked_reason,omitempty"`
	LastAction        string            `json:"last_action,omitempty"`
	PhaseHistory      []PhaseTransition `json:"phase_history"`
}

// WorkspaceState tracks workspace-related state
type WorkspaceState struct {
	Initialized    bool   `json:"initialized"`
	Root           string `json:"root"`
	HasFoundation  bool   `json:"has_foundation"`
	HasArtifacts   bool   `json:"has_artifacts"`
	HasAuditReport bool   `json:"has_audit_report"`
}

// PhaseTransition records a phase change with timestamp
type PhaseTransition struct {
	From      Phase   `json:"from"`
	To        Phase   `json:"to"`
	Timestamp string  `json:"timestamp"`
	Confidence float64 `json:"confidence"`
	Evidence  []string `json:"evidence,omitempty"`
}

// Manager handles session state persistence and transitions
type Manager struct {
	repoPath    string
	guard       *SessionGuard
	enforcer    *enforcer.Engine
	initialized bool
}

// SessionGuard provides workspace and session initialization enforcement
// This is a lightweight version that integrates with the guard package
type SessionGuard struct {
	repoPath string
	wsPath   string
}

// NewSessionGuard creates a new session guard for the given repository path
func NewSessionGuard(repoPath string) *SessionGuard {
	return &SessionGuard{
		repoPath: repoPath,
		wsPath:   filepath.Join(repoPath, ".kdse"),
	}
}

// IsInitialized checks if the workspace is properly initialized
func (g *SessionGuard) IsInitialized() bool {
	// Check if .kdse directory exists
	if _, err := os.Stat(g.wsPath); os.IsNotExist(err) {
		return false
	}

	// Check if session state file exists
	sessionPath := filepath.Join(g.wsPath, "session-state.json")
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return false
	}

	// Try to load and validate session state
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return false
	}

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return false
	}

	// Validate required fields
	if state["session_id"] == nil || state["session_id"] == "" {
		return false
	}
	if state["started_at"] == nil || state["started_at"] == "" {
		return false
	}

	return true
}

// NewManager creates a new orchestration manager for the given repository path
func NewManager(repoPath string) *Manager {
	return &Manager{
		repoPath:    repoPath,
		guard:       NewSessionGuard(repoPath),
		enforcer:    enforcer.NewEngine(repoPath),
		initialized: false,
	}
}

// IsInitialized checks if the workspace is properly initialized
func (m *Manager) IsInitialized() bool {
	if !m.guard.IsInitialized() {
		return false
	}
	
	// Also check orchestration state
	_, err := m.Load()
	return err == nil
}

// Initialize creates a new session with default state
func (m *Manager) Initialize() (*SessionState, error) {
	now := time.Now().Format(time.RFC3339)

	state := &SessionState{
		SessionID:         generateSessionID(),
		StartedAt:         now,
		UpdatedAt:         now,
		CurrentPhase:      PhaseIdle,
		Confidence:        0.0,
		Evidence:          []string{},
		CompletedPhases:   []Phase{},
		NextAllowedPhases: []Phase{PhaseProblem},
		ExecutionMode:     ModeStrict, // STRICT mode enabled by default
		Workspace: &WorkspaceState{
			Initialized: false,
		},
	}

	// Ensure .kdse directory exists
	kdseDir := filepath.Join(m.repoPath, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return nil, err
	}

	// Create session-state.json for the guard
	if err := m.createSessionState(state); err != nil {
		log.Printf("[ORCH] Warning: Failed to create session-state.json: %v", err)
	}

	if err := m.Save(state); err != nil {
		return nil, err
	}

	m.initialized = true
	log.Printf("[ORCH] Session initialized: %s", state.SessionID)
	return state, nil
}

// createSessionState creates the session-state.json file for the guard
func (m *Manager) createSessionState(state *SessionState) error {
	sessionState := map[string]interface{}{
		"session_id":      state.SessionID,
		"started_at":       state.StartedAt,
		"updated_at":       state.UpdatedAt,
		"workspace_root":   filepath.Join(m.repoPath, ".kdse"),
		"version":          "1.0.0",
		"initialized":      true,
	}

	data, err := json.MarshalIndent(sessionState, "", "  ")
	if err != nil {
		return err
	}

	sessionPath := filepath.Join(m.repoPath, ".kdse", "session-state.json")
	return os.WriteFile(sessionPath, data, 0644)
}

// Load retrieves the current session state
func (m *Manager) Load() (*SessionState, error) {
	statePath := m.statePath()
	
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no session state found - run initialize first")
		}
		return nil, err
	}

	var state SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// Save persists the session state
func (m *Manager) Save(state *SessionState) error {
	state.UpdatedAt = time.Now().Format(time.RFC3339)
	
	// Ensure .kdse directory exists
	kdseDir := filepath.Join(m.repoPath, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return err
	}

	statePath := m.statePath()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0644)
}

// TransitionTo moves the session to a new phase if valid
func (m *Manager) TransitionTo(newPhase Phase, evidence []string) (*SessionState, error) {
	state, err := m.Load()
	if err != nil {
		return nil, err
	}

	// Check if transition is valid
	if !m.canTransition(state.CurrentPhase, newPhase) {
		return nil, fmt.Errorf("invalid transition from %s to %s", state.CurrentPhase, newPhase)
	}

	// Record transition
	transition := PhaseTransition{
		From:        state.CurrentPhase,
		To:          newPhase,
		Timestamp:   time.Now().Format(time.RFC3339),
		Confidence:  state.Confidence,
		Evidence:    evidence,
	}

	// Update state
	if newPhase != PhaseComplete {
		state.CompletedPhases = append(state.CompletedPhases, state.CurrentPhase)
	}
	state.CurrentPhase = newPhase
	state.PhaseHistory = append(state.PhaseHistory, transition)
	state.NextAllowedPhases = m.getNextAllowedPhases(newPhase)
	state.LastAction = fmt.Sprintf("Transitioned to %s", newPhase)
	
	if evidence != nil {
		state.Evidence = append(state.Evidence, evidence...)
	}

	if err := m.Save(state); err != nil {
		return nil, err
	}

	return state, nil
}

// canTransition checks if a phase transition is valid
func (m *Manager) canTransition(from, to Phase) bool {
	allowed, exists := PhaseTransitions[from]
	if !exists {
		return false
	}
	
	for _, p := range allowed {
		if p == to {
			return true
		}
	}
	return false
}

// getNextAllowedPhases returns the phases that can follow the current phase
func (m *Manager) getNextAllowedPhases(current Phase) []Phase {
	if phases, exists := PhaseTransitions[current]; exists {
		return phases
	}
	return []Phase{}
}

// UpdateConfidence updates the confidence level
func (m *Manager) UpdateConfidence(confidence float64) (*SessionState, error) {
	state, err := m.Load()
	if err != nil {
		return nil, err
	}

	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}

	state.Confidence = confidence
	state.LastAction = fmt.Sprintf("Updated confidence to %.2f", confidence)

	if err := m.Save(state); err != nil {
		return nil, err
	}

	return state, nil
}

// UpdateWorkspace updates workspace state
func (m *Manager) UpdateWorkspace(ws *WorkspaceState) (*SessionState, error) {
	state, err := m.Load()
	if err != nil {
		return nil, err
	}

	state.Workspace = ws
	state.LastAction = "Updated workspace state"

	if err := m.Save(state); err != nil {
		return nil, err
	}

	return state, nil
}

// SetObjective sets the user objective for this session
func (m *Manager) SetObjective(objective string) (*SessionState, error) {
	state, err := m.Load()
	if err != nil {
		return nil, err
	}

	state.Objective = objective
	state.LastAction = fmt.Sprintf("Set objective: %s", objective)

	if err := m.Save(state); err != nil {
		return nil, err
	}

	return state, nil
}

// CanProceedTo checks if we can proceed to a specific phase
func (m *Manager) CanProceedTo(targetPhase Phase) (bool, string) {
	state, err := m.Load()
	if err != nil {
		return false, fmt.Sprintf("failed to load session state: %v", err)
	}

	// Check if in strict mode
	if state.ExecutionMode != ModeStrict {
		return true, ""
	}

	// Check if phase is in allowed next phases
	for _, allowed := range state.NextAllowedPhases {
		if allowed == targetPhase {
			return true, ""
		}
	}

	// Check if transitioning to a later phase without completing prerequisites
	if targetPhase == PhaseImplementation {
		if state.CurrentPhase != PhaseArchitecture {
			return false, fmt.Sprintf("Cannot implement: current phase is %s, must complete Architecture first", state.CurrentPhase)
		}
		if state.Confidence < PhaseConfidenceThreshold[PhaseImplementation] {
			return false, fmt.Sprintf("Cannot implement: confidence %.2f below threshold %.2f", 
				state.Confidence, PhaseConfidenceThreshold[PhaseImplementation])
		}
	}

	return false, fmt.Sprintf("Phase %s not allowed from current state. Allowed: %v", targetPhase, state.NextAllowedPhases)
}

// CheckImplementationReadiness verifies if the project is ready for implementation
// This method now integrates with the enforcer engine for strict KDSE principle enforcement
func (m *Manager) CheckImplementationReadiness() (bool, *SessionState) {
	state, err := m.Load()
	if err != nil {
		return false, nil
	}

	// FIRST: Run enforcer checks for strict KDSE principle enforcement
	// This blocks premature implementation if foundation/knowledge/architecture is missing
	enforcementErr := m.enforcer.EnforceImplementation()
	
	if enforcementErr != nil {
		// Get violations from enforcer
		violations := m.enforcer.GetViolations()
		var reasons []string
		for _, v := range violations {
			if v.Blocked {
				reasons = append(reasons, fmt.Sprintf("[%s] %s", v.Code, v.Message))
			}
		}
		state.BlockedReason = fmt.Sprintf("KDSE Enforcement blocked: %v", reasons)
		
		// Also include in session state for the LLM to see
		if state.Workspace == nil {
			state.Workspace = &WorkspaceState{}
		}
		state.Workspace.HasFoundation = !m.hasEnforcementViolation(violations, enforcer.CodeNoFoundation, enforcer.CodeFoundationIncomplete)
		state.Workspace.HasArtifacts = !m.hasEnforcementViolation(violations, enforcer.CodeNoKnowledge, enforcer.CodeKnowledgeIncomplete)
		
		return false, state
	}

	// SECOND: Check orchestration-level prerequisites
	prerequisites := []struct {
		check  bool
		reason string
	}{
		{state.CurrentPhase == PhaseArchitecture, fmt.Sprintf("Current phase: %s (need Architecture)", state.CurrentPhase)},
		{state.Confidence >= PhaseConfidenceThreshold[PhaseImplementation], fmt.Sprintf("Confidence: %.2f (need %.2f)", state.Confidence, PhaseConfidenceThreshold[PhaseImplementation])},
		{state.Workspace != nil && state.Workspace.HasFoundation, "Foundation exists"},
		{state.Workspace != nil && state.Workspace.HasAuditReport, "Audit report exists"},
	}

	blocked := false
	var reasons []string
	for _, prereq := range prerequisites {
		if !prereq.check {
			blocked = true
			reasons = append(reasons, prereq.reason)
		}
	}

	if blocked {
		state.BlockedReason = fmt.Sprintf("Implementation blocked: %v", reasons)
		return false, state
	}

	state.BlockedReason = ""
	return true, state
}

// hasEnforcementViolation checks if violations contain a specific code
func (m *Manager) hasEnforcementViolation(violations []enforcer.EnforcementError, codes ...string) bool {
	violationMap := make(map[string]bool)
	for _, v := range violations {
		violationMap[v.Code] = true
	}
	for _, code := range codes {
		if violationMap[code] {
			return true
		}
	}
	return false
}

// GetEnforcementReport returns the current enforcement status
func (m *Manager) GetEnforcementReport() *enforcer.EnforcementReport {
	return m.enforcer.GenerateEnforcementReport()
}

// GetExecutionDecision determines what KDSE operations to invoke based on session state.
// Returns an ExecutionDecision with an attached WorkOrder that explicitly defines what the LLM must do.
func (m *Manager) GetExecutionDecision(objective string) *ExecutionDecision {
	state, err := m.Load()
	if err != nil {
		return &ExecutionDecision{
			Action:   "initialize",
			Reason:   "No session state found",
			NextPhase: PhaseProblem,
		}
	}

	decision := &ExecutionDecision{
		SessionState: state,
	}

	// If no objective set, set it first
	if state.Objective == "" && objective != "" {
		decision.Action = "set_objective"
		decision.NextPhase = PhaseProblem // Always transition to Problem first
		decision.Reason = "Setting objective before workflow"
		// Generate WorkOrder for the next phase (Problem)
		decision.WorkOrder = m.GenerateWorkOrder(PhaseProblem, objective)
		return decision
	}

	// Determine next action based on current phase
	switch state.CurrentPhase {
	case PhaseIdle:
		decision.Action = "problem"
		decision.NextPhase = PhaseProblem
		decision.Reason = "Starting new KDSE workflow"
		decision.Operations = []string{"analyze_objective"}
		decision.WorkOrder = m.GenerateWorkOrder(PhaseProblem, state.Objective)

	case PhaseProblem:
		decision.Action = "knowledge"
		decision.NextPhase = PhaseKnowledge
		decision.Reason = "Problem defined, collecting knowledge"
		decision.Operations = []string{"collect"}
		decision.WorkOrder = m.GenerateWorkOrder(PhaseKnowledge, state.Objective)

	case PhaseKnowledge:
		decision.Action = "foundation"
		decision.NextPhase = PhaseFoundation
		decision.Reason = "Knowledge collected, establishing foundation"
		decision.Operations = []string{"foundation"}
		decision.WorkOrder = m.GenerateWorkOrder(PhaseFoundation, state.Objective)

	case PhaseFoundation:
		decision.Action = "audit"
		decision.NextPhase = PhaseAudit
		decision.Reason = "Foundation established, running audit"
		decision.Operations = []string{"audit"}
		decision.WorkOrder = m.GenerateWorkOrder(PhaseAudit, state.Objective)

	case PhaseAudit:
		decision.Action = "assessment"
		decision.NextPhase = PhaseAssessment
		decision.Reason = "Audit complete, assessing findings"
		decision.Operations = []string{"assess"}
		decision.WorkOrder = m.GenerateWorkOrder(PhaseAssessment, state.Objective)

	case PhaseAssessment:
		decision.Action = "architecture"
		decision.NextPhase = PhaseArchitecture
		decision.Reason = "Assessment complete, defining architecture"
		decision.Operations = []string{"architecture"}
		decision.WorkOrder = m.GenerateWorkOrder(PhaseArchitecture, state.Objective)

	case PhaseArchitecture:
		// Check if ready for implementation
		if ready, _ := m.CheckImplementationReadiness(); ready {
			decision.Action = "implement"
			decision.NextPhase = PhaseImplementation
			decision.Reason = "All prerequisites met, ready for implementation"
			decision.Operations = []string{"implement"}
			decision.WorkOrder = m.GenerateWorkOrder(PhaseImplementation, state.Objective)
		} else {
			decision.Action = "blocked"
			decision.NextPhase = PhaseImplementation
			decision.Blocked = true
			decision.BlockedReason = state.BlockedReason
			decision.Reason = "Implementation blocked - prerequisites not met"
			// Include partial work order showing what blocks implementation
			decision.WorkOrder = m.GenerateWorkOrder(PhaseArchitecture, state.Objective)
			decision.WorkOrder.BlockedActions = append(decision.WorkOrder.BlockedActions, 
				"IMPLEMENTATION BLOCKED: "+state.BlockedReason)
		}

	case PhaseImplementation:
		decision.Action = "complete"
		decision.NextPhase = PhaseComplete
		decision.Reason = "Implementation complete"
		decision.Operations = []string{"verify", "report"}

	case PhaseComplete:
		decision.Action = "idle"
		decision.NextPhase = PhaseIdle
		decision.Reason = "Session complete"
	}

	return decision
}

// ExecutionDecision contains the orchestration decision
type ExecutionDecision struct {
	Action         string       `json:"action"`
	Reason         string      `json:"reason"`
	NextPhase      Phase        `json:"next_phase"`
	Operations     []string    `json:"operations,omitempty"`
	Blocked        bool        `json:"blocked,omitempty"`
	BlockedReason  string      `json:"blocked_reason,omitempty"`
	SessionState   *SessionState `json:"session_state,omitempty"`
	WorkOrder      *WorkOrder   `json:"work_order,omitempty"`
}

// WorkOrder is an explicit engineering directive from KDSE runtime to the LLM.
// The runtime owns the methodology; the LLM executes only the Work Order.
type WorkOrder struct {
	Phase               Phase              `json:"phase"`
	PhaseDescription    string             `json:"phase_description"`
	RequiredWork        []string           `json:"required_work"`
	ExpectedDeliverables []string          `json:"expected_deliverables"`
	CompletionCriteria  []string           `json:"completion_criteria"`
	BlockedActions      []string           `json:"blocked_actions"`
	NextPhase           Phase              `json:"next_phase"`
	StrictModeEnforced  bool               `json:"strict_mode_enforced"`
	ConfidenceGate      float64            `json:"confidence_gate,omitempty"`
}

// GenerateWorkOrder creates an explicit work order for the given phase
func (m *Manager) GenerateWorkOrder(phase Phase, objective string) *WorkOrder {
	switch phase {
	case PhaseProblem:
		return &WorkOrder{
			Phase:               PhaseProblem,
			PhaseDescription:    "Define and understand the problem scope and constraints",
			RequiredWork: []string{
				"Analyze the user objective: " + objective,
				"Identify the core problem being solved",
				"Define explicit scope boundaries",
				"Identify stakeholders and users",
				"Document known constraints and requirements",
			},
			ExpectedDeliverables: []string{
				".kdse/context/PROBLEM.md - Problem statement and scope",
			},
			CompletionCriteria: []string{
				"PROBLEM.md exists in .kdse/context/",
				"Problem statement clearly articulates the user's need",
				"Scope boundaries are explicitly defined",
				"At least one stakeholder/user is identified",
			},
			BlockedActions: []string{
				"DO NOT generate any code or implementation",
				"DO NOT create project structure or folders outside .kdse/",
				"DO NOT design architecture or technical solutions",
				"DO NOT write tests or configuration files",
			},
			NextPhase:           PhaseKnowledge,
			StrictModeEnforced:  true,
		}

	case PhaseKnowledge:
		return &WorkOrder{
			Phase:               PhaseKnowledge,
			PhaseDescription:    "Collect existing knowledge and artifacts from the repository",
			RequiredWork: []string{
				"Scan the repository for existing documentation",
				"Identify README files and project descriptions",
				"Collect any existing specifications or requirements",
				"Document the current state of the codebase",
				"Note any existing patterns or conventions",
			},
			ExpectedDeliverables: []string{
				".kdse/context/KNOWLEDGE.md - Repository knowledge summary",
			},
			CompletionCriteria: []string{
				"KNOWLEDGE.md exists in .kdse/context/",
				"Repository structure is documented",
				"Existing documentation is cataloged",
				"Relevant patterns/conventions are noted",
			},
			BlockedActions: []string{
				"DO NOT modify any existing code",
				"DO NOT create new source files",
				"DO NOT design solutions or architecture",
				"DO NOT write implementation code",
			},
			NextPhase:           PhaseFoundation,
			StrictModeEnforced:  true,
		}

	case PhaseFoundation:
		return &WorkOrder{
			Phase:               PhaseFoundation,
			PhaseDescription:    "Establish foundational documentation for the project",
			RequiredWork: []string{
				"Create SPEC.md with detailed specifications",
				"Create REQUIREMENTS.md with functional requirements",
				"Create ASSUMPTIONS.md documenting key assumptions",
				"Create CONSTRAINTS.md listing project constraints",
				"Create GLOSSARY.md defining domain terminology",
			},
			ExpectedDeliverables: []string{
				".kdse/foundation/SPEC.md - Project specification",
				".kdse/foundation/REQUIREMENTS.md - Functional requirements",
				".kdse/foundation/ASSUMPTIONS.md - Key assumptions",
				".kdse/foundation/CONSTRAINTS.md - Project constraints",
				".kdse/foundation/GLOSSARY.md - Domain terminology",
			},
			CompletionCriteria: []string{
				"All 5 foundation documents exist in .kdse/foundation/",
				"SPEC.md contains detailed project specification",
				"REQUIREMENTS.md lists functional requirements",
				"ASSUMPTIONS.md documents at least 3 key assumptions",
				"CONSTRAINTS.md lists technical and business constraints",
				"GLOSSARY.md defines domain terminology",
			},
			BlockedActions: []string{
				"DO NOT generate any source code",
				"DO NOT create project directory structure",
				"DO NOT write implementation files",
				"DO NOT create configuration files",
				"DO NOT write tests or build scripts",
			},
			NextPhase:           PhaseAudit,
			StrictModeEnforced:  true,
		}

	case PhaseAudit:
		return &WorkOrder{
			Phase:               PhaseAudit,
			PhaseDescription:    "Run compliance audit against KDSE standards",
			RequiredWork: []string{
				"Verify all foundation documents are complete",
				"Check that SPEC.md meets quality standards",
				"Validate REQUIREMENTS.md completeness",
				"Identify any gaps or missing information",
				"Generate audit report with findings",
			},
			ExpectedDeliverables: []string{
				".kdse/reports/audit-[timestamp].md - Audit report",
			},
			CompletionCriteria: []string{
				"Audit report exists in .kdse/reports/",
				"All foundation documents are verified",
				"Any gaps are identified with remediation steps",
				"Audit score/rating is provided",
			},
			BlockedActions: []string{
				"DO NOT modify any foundation documents during audit",
				"DO NOT generate code or implementation",
				"DO NOT create project structure",
			},
			NextPhase:           PhaseAssessment,
			StrictModeEnforced:  true,
		}

	case PhaseAssessment:
		return &WorkOrder{
			Phase:               PhaseAssessment,
			PhaseDescription:    "Assess audit findings and identify improvement opportunities",
			RequiredWork: []string{
				"Review all audit findings",
				"Prioritize identified issues by severity",
				"Document assessment of foundation quality",
				"Identify specific improvements needed",
				"Create remediation plan if issues found",
			},
			ExpectedDeliverables: []string{
				".kdse/reports/assessment-[timestamp].md - Assessment report",
			},
			CompletionCriteria: []string{
				"Assessment report exists in .kdse/reports/",
				"All critical issues are documented",
				"Improvement recommendations are provided",
				"Remediation plan is documented if needed",
			},
			BlockedActions: []string{
				"DO NOT implement any fixes during assessment",
				"DO NOT generate code or implementation",
				"DO NOT modify foundation documents",
			},
			NextPhase:           PhaseArchitecture,
			StrictModeEnforced:  true,
		}

	case PhaseArchitecture:
		return &WorkOrder{
			Phase:               PhaseArchitecture,
			PhaseDescription:    "Define system architecture and technical approach",
			RequiredWork: []string{
				"Design system architecture based on SPEC.md",
				"Define component structure and responsibilities",
				"Document technology choices and rationale",
				"Create architectural diagrams (if applicable)",
				"Define interfaces and data flows",
			},
			ExpectedDeliverables: []string{
				".kdse/foundation/ARCHITECTURE.md - System architecture document",
			},
			CompletionCriteria: []string{
				"ARCHITECTURE.md exists in .kdse/foundation/",
				"System components are defined",
				"Technology stack is documented",
				"Data flows are described",
				"Architectural decisions are justified",
			},
			BlockedActions: []string{
				"DO NOT generate implementation code",
				"DO NOT create source files",
				"DO NOT set up project build structure",
				"DO NOT write configuration files",
			},
			NextPhase:           PhaseImplementation,
			StrictModeEnforced:  true,
			ConfidenceGate:      PhaseConfidenceThreshold[PhaseImplementation],
		}

	case PhaseImplementation:
		return &WorkOrder{
			Phase:               PhaseImplementation,
			PhaseDescription:    "Implement the solution based on approved architecture",
			RequiredWork: []string{
				"Follow ARCHITECTURE.md for implementation",
				"Implement components as specified",
				"Write unit tests for core functionality",
				"Ensure code follows defined patterns",
				"Document implementation decisions",
			},
			ExpectedDeliverables: []string{
				"Source code files in project directory",
				"Unit tests for core functionality",
				"Configuration files as needed",
			},
			CompletionCriteria: []string{
				"Code compiles without errors",
				"Core functionality is implemented per SPEC.md",
				"Unit tests cover critical paths",
				"Code follows project conventions",
			},
			BlockedActions: []string{
				"DO NOT deviate from ARCHITECTURE.md",
				"DO NOT implement features not in REQUIREMENTS.md",
				"DO NOT skip testing",
			},
			NextPhase:           PhaseComplete,
			StrictModeEnforced:  true,
		}

	default:
		return &WorkOrder{
			Phase:               phase,
			PhaseDescription:    "Unknown phase",
			RequiredWork:        []string{},
			ExpectedDeliverables: []string{},
			CompletionCriteria:  []string{},
			BlockedActions:      []string{"DO NOT perform any actions"},
			NextPhase:           PhaseComplete,
			StrictModeEnforced:  true,
		}
	}
}

// statePath returns the path to the session state file
func (m *Manager) statePath() string {
	return filepath.Join(m.repoPath, ".kdse", "orchestration_state.json")
}

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	return fmt.Sprintf("KDSE-ORCH-%s", time.Now().Format("20060102-150405"))
}

// GetSessionStatus returns a human-readable status summary
func (m *Manager) GetSessionStatus() (string, error) {
	state, err := m.Load()
	if err != nil {
		return "", err
	}

	status := fmt.Sprintf(`KDSE Orchestration Session
========================
Session ID: %s
Current Phase: %s
Confidence: %.2f
Execution Mode: %s
Completed: %v
Next Allowed: %v
`,
		state.SessionID,
		state.CurrentPhase,
		state.Confidence,
		state.ExecutionMode,
		state.CompletedPhases,
		state.NextAllowedPhases,
	)

	if state.BlockedReason != "" {
		status += fmt.Sprintf("\nBlocked: %s\n", state.BlockedReason)
	}

	return status, nil
}

// SetExecutionMode changes the execution mode
func (m *Manager) SetExecutionMode(mode ExecutionMode) (*SessionState, error) {
	state, err := m.Load()
	if err != nil {
		return nil, err
	}

	state.ExecutionMode = mode
	state.LastAction = fmt.Sprintf("Changed execution mode to %s", mode)

	if err := m.Save(state); err != nil {
		return nil, err
	}

	return state, nil
}

// IsStrictMode returns true if strict mode is enabled
func (m *Manager) IsStrictMode() bool {
	state, err := m.Load()
	if err != nil {
		return true // Default to strict mode
	}
	return state.ExecutionMode == ModeStrict
}
