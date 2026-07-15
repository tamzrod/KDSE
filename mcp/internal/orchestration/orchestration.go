// Package orchestration provides the KDSE orchestration engine that transforms
// KDSE from a toolbox into an orchestration engine. After initialization, the
// LLM never decides which KDSE tool to call - the execute tool decides.
package orchestration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Phase represents the current phase in the KDSE workflow
type Phase string

const (
	PhaseIdle             Phase = "Idle"
	PhaseProblem          Phase = "Problem"
	PhaseKnowledge        Phase = "Knowledge Collection"
	PhaseFoundation       Phase = "Foundation"
	PhaseAudit            Phase = "Audit"
	PhaseAssessment       Phase = "Assessment"
	PhaseArchitecture     Phase = "Architecture"
	PhaseImplementation   Phase = "Implementation"
	PhaseComplete         Phase = "Complete"
	PhaseBlocked          Phase = "Blocked"
)

// ExecutionMode represents the mode of execution
type ExecutionMode string

const (
	ModeToolbox ExecutionMode = "toolbox" // Legacy mode - tools available directly
	ModeStrict  ExecutionMode = "strict"   // Strict mode - all requests through execute
)

// PhaseTransition defines valid transitions between phases
var PhaseTransitions = map[Phase][]Phase{
	PhaseIdle:           {PhaseProblem},
	PhaseProblem:        {PhaseKnowledge},
	PhaseKnowledge:      {PhaseFoundation},
	PhaseFoundation:     {PhaseAudit},
	PhaseAudit:          {PhaseAssessment, PhaseArchitecture},
	PhaseAssessment:     {PhaseArchitecture, PhaseFoundation},
	PhaseArchitecture:   {PhaseImplementation},
	PhaseImplementation: {PhaseComplete},
}

// Minimum confidence thresholds for each phase
var PhaseConfidenceThreshold = map[Phase]float64{
	PhaseIdle:           0.0,
	PhaseProblem:        0.6,
	PhaseKnowledge:      0.7,
	PhaseFoundation:     0.75,
	PhaseAudit:          0.8,
	PhaseAssessment:     0.8,
	PhaseArchitecture:   0.85,
	PhaseImplementation: 0.9,
}

// SessionState represents the runtime session state for KDSE orchestration
type SessionState struct {
	SessionID         string            `json:"session_id"`
	StartedAt         string            `json:"started_at"`
	UpdatedAt         string            `json:"updated_at"`
	CurrentPhase      Phase             `json:"current_phase"`
	Confidence        float64           `json:"confidence"`
	Evidence          []string          `json:"evidence"`
	CompletedPhases   []Phase           `json:"completed_phases"`
	NextAllowedPhases []Phase           `json:"next_allowed_phases"`
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
	repoPath string
}

// NewManager creates a new orchestration manager for the given repository path
func NewManager(repoPath string) *Manager {
	return &Manager{repoPath: repoPath}
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

	if err := m.Save(state); err != nil {
		return nil, err
	}

	return state, nil
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
func (m *Manager) CheckImplementationReadiness() (bool, *SessionState) {
	state, err := m.Load()
	if err != nil {
		return false, nil
	}

	// Check all prerequisites for implementation
	prerequisites := []struct {
		check    bool
		reason   string
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

// GetExecutionDecision determines what KDSE operations to invoke based on session state
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
		decision.NextPhase = state.CurrentPhase
		decision.Reason = "Setting objective before workflow"
		return decision
	}

	// Determine next action based on current phase
	switch state.CurrentPhase {
	case PhaseIdle:
		decision.Action = "problem"
		decision.NextPhase = PhaseProblem
		decision.Reason = "Starting new KDSE workflow"
		decision.Operations = []string{"analyze_objective"}

	case PhaseProblem:
		decision.Action = "knowledge"
		decision.NextPhase = PhaseKnowledge
		decision.Reason = "Problem defined, collecting knowledge"
		decision.Operations = []string{"collect"}

	case PhaseKnowledge:
		decision.Action = "foundation"
		decision.NextPhase = PhaseFoundation
		decision.Reason = "Knowledge collected, establishing foundation"
		decision.Operations = []string{"foundation"}

	case PhaseFoundation:
		decision.Action = "audit"
		decision.NextPhase = PhaseAudit
		decision.Reason = "Foundation established, running audit"
		decision.Operations = []string{"audit"}

	case PhaseAudit:
		decision.Action = "assessment"
		decision.NextPhase = PhaseAssessment
		decision.Reason = "Audit complete, assessing findings"
		decision.Operations = []string{"assess"}

	case PhaseAssessment:
		decision.Action = "architecture"
		decision.NextPhase = PhaseArchitecture
		decision.Reason = "Assessment complete, defining architecture"
		decision.Operations = []string{"architecture"}

	case PhaseArchitecture:
		// Check if ready for implementation
		if ready, _ := m.CheckImplementationReadiness(); ready {
			decision.Action = "implement"
			decision.NextPhase = PhaseImplementation
			decision.Reason = "All prerequisites met, ready for implementation"
			decision.Operations = []string{"implement"}
		} else {
			decision.Action = "blocked"
			decision.NextPhase = PhaseImplementation
			decision.Blocked = true
			decision.BlockedReason = state.BlockedReason
			decision.Reason = "Implementation blocked - prerequisites not met"
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
