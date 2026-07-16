package orchestration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Engine is the state-based orchestration engine
type Engine struct {
	config        *EngineConfig
	resolver      *WorkspaceResolver
	confidence    *ConfidenceEvaluator
	evidence      *EvidenceEvaluator
	state         *OrchestrationState
	sessionActive bool
}

// NewEngine creates a new orchestration engine
func NewEngine(config *EngineConfig) (*Engine, error) {
	if config == nil {
		config = DefaultEngineConfig()
	}

	resolver, err := NewWorkspaceResolver(config)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		config:     config,
		resolver:    resolver,
		confidence:  NewConfidenceEvaluator(config),
		evidence:    NewEvidenceEvaluator(config),
		state:       nil,
	}

	return engine, nil
}

// NewEngineWithResolver creates an engine with a specific resolver
func NewEngineWithResolver(config *EngineConfig, resolver *WorkspaceResolver) *Engine {
	if config == nil {
		config = DefaultEngineConfig()
	}

	return &Engine{
		config:     config,
		resolver:    resolver,
		confidence:  NewConfidenceEvaluator(config),
		evidence:    NewEvidenceEvaluator(config),
		state:       nil,
	}
}

// Initialize initializes the orchestration session
func (e *Engine) Initialize(workingPath string) error {
	// Step 1: Resolve workspace
	workspace, err := e.resolver.ResolveWorkspace(workingPath)
	if err != nil {
		return fmt.Errorf("failed to resolve workspace: %w", err)
	}

	// Ensure KDSE directory exists
	if err := os.MkdirAll(workspace.KDSEPath, 0755); err != nil {
		return fmt.Errorf("failed to create KDSE directory: %w", err)
	}

	// Initialize state
	e.state = &OrchestrationState{
		SessionID:    generateSessionID(),
		CurrentPhase: PhaseProblem,
		Workspace:    *workspace,
		Metrics: OrchestrationMetrics{
			CycleCount:      0,
			PhaseExecutions: make(map[string]int),
			StartTime:       time.Now(),
			LastCycleTime:   time.Now(),
		},
		History: []PhaseTransition{},
		Blocked: BlockedState{
			Blocked:    false,
			Reasons:    []string{},
			Required:   []string{},
			CanRetry:   true,
		},
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Evaluate initial state
	if err := e.evaluateCurrentState(); err != nil {
		return fmt.Errorf("failed to evaluate initial state: %w", err)
	}

	e.sessionActive = true
	return nil
}

// ExecuteCycle runs a single orchestration cycle
// Each cycle: Problem → Knowledge → Foundation → Audit → Assessment → Architecture → Implementation → Complete
func (e *Engine) ExecuteCycle() (*ExecuteCycleResult, error) {
	if !e.sessionActive || e.state == nil {
		return nil, fmt.Errorf("session not initialized")
	}

	e.state.Metrics.CycleCount++
	cycleStart := time.Now()

	result := &ExecuteCycleResult{
		CycleNumber: e.state.Metrics.CycleCount,
	}

	// Step 1: Re-resolve workspace (in case it changed)
	if err := e.resolveWorkspace(); err != nil {
		result.Error = fmt.Sprintf("workspace resolution failed: %v", err)
		result.Success = false
		return result, err
	}

	// Step 2: Evaluate current state
	if err := e.evaluateCurrentState(); err != nil {
		result.Error = fmt.Sprintf("state evaluation failed: %v", err)
		result.Success = false
		return result, err
	}

	// Step 3: Evaluate confidence
	if err := e.evaluateConfidence(); err != nil {
		result.Error = fmt.Sprintf("confidence evaluation failed: %v", err)
		result.Success = false
		return result, err
	}

	// Step 4: Evaluate missing evidence
	if err := e.evaluateMissingEvidence(); err != nil {
		result.Error = fmt.Sprintf("evidence evaluation failed: %v", err)
		result.Success = false
		return result, err
	}

	// Step 5: Decide next phase
	decision := e.decideNextPhase()
	result.Decision = decision

	// Step 6: Execute only the decided phase
	if decision.ShouldExecute {
		if err := e.executePhase(decision.NextPhase); err != nil {
			result.Error = fmt.Sprintf("phase execution failed: %v", err)
			result.Success = false
			return result, err
		}
	}

	// Step 7: Re-evaluate after execution
	if err := e.reevaluate(); err != nil {
		result.Error = fmt.Sprintf("re-evaluation failed: %v", err)
		// Non-fatal, continue
	}

	// Update metrics
	e.state.Metrics.LastCycleTime = time.Now()
	e.state.Metrics.CycleDuration = float64(time.Since(cycleStart).Milliseconds())
	e.state.UpdatedAt = time.Now().Format(time.RFC3339)

	result.PhaseExecuted = e.state.CurrentPhase
	result.State = e.state
	result.Success = true
	result.Continue = e.shouldContinue()

	return result, nil
}

// resolveWorkspace resolves/refreshes the workspace
func (e *Engine) resolveWorkspace() error {
	workspace, err := e.resolver.ResolveWorkspace(e.state.Workspace.ResolvedPath)
	if err != nil {
		return err
	}
	e.state.Workspace = *workspace
	return nil
}

// evaluateCurrentState evaluates the current orchestration state
func (e *Engine) evaluateCurrentState() error {
	// Check if we need to block based on evidence
	canProceed, blockedReasons := e.evidence.CanProceedToPhase(&e.state.Workspace, e.state.CurrentPhase)

	if !canProceed {
		e.state.Blocked = BlockedState{
			Blocked:    true,
			Reasons:    blockedReasons,
			Required:   blockedReasons,
			CanRetry:   true,
		}
	} else {
		e.state.Blocked = BlockedState{
			Blocked:    false,
			Reasons:    []string{},
			Required:   []string{},
			CanRetry:   true,
		}
	}

	return nil
}

// evaluateConfidence evaluates confidence for the current state
func (e *Engine) evaluateConfidence() error {
	conf, err := e.confidence.EvaluateConfidence(&e.state.Workspace)
	if err != nil {
		return err
	}
	e.state.Confidence = *conf
	return nil
}

// evaluateMissingEvidence evaluates missing evidence for current phase
func (e *Engine) evaluateMissingEvidence() error {
	state, err := e.evidence.EvaluateEvidence(&e.state.Workspace, e.state.CurrentPhase)
	if err != nil {
		return err
	}
	e.state.EvidenceState = *state
	return nil
}

// decideNextPhase decides which phase to execute next
func (e *Engine) decideNextPhase() *PhaseDecision {
	decision := &PhaseDecision{
		NextPhase:    e.state.CurrentPhase,
		Reason:       "Current phase not complete",
		Confidence:   e.state.Confidence.Overall,
		ShouldExecute: false,
	}

	// Check if blocked
	if e.state.Blocked.Blocked {
		decision.ShouldExecute = false
		decision.Reason = "Blocked by missing evidence"
		decision.BlockingReasons = e.state.Blocked.Reasons
		return decision
	}

	// Check if Foundation threshold is met (blocks Implementation)
	if e.state.CurrentPhase == PhaseImplement && !e.state.Confidence.MeetsThreshold {
		decision.ShouldExecute = false
		decision.NextPhase = PhaseFoundation
		decision.Reason = fmt.Sprintf(
			"Foundation threshold not met (%.2f < %.2f), must improve Foundation before Implementation",
			e.state.Confidence.Foundation,
			e.state.Confidence.Threshold,
		)
		decision.BlockingReasons = []string{"Foundation confidence below threshold"}
		return decision
	}

	// Phase completion logic - MCP canonical phases
	switch e.state.CurrentPhase {
	case PhaseProblem:
		// Problem definition complete
		if e.state.Workspace.KDSEPath != "" {
			decision.NextPhase = PhaseKnowledge
			decision.Reason = "Problem definition complete"
			decision.ShouldExecute = true
		}

	case PhaseKnowledge:
		// Knowledge collection complete
		if e.state.EvidenceState.Completeness >= e.config.EvidenceThreshold {
			decision.NextPhase = PhaseFoundation
			decision.Reason = "Knowledge collection complete"
			decision.ShouldExecute = true
		}

	case PhaseFoundation:
		// Foundation complete
		if e.state.Confidence.MeetsThreshold {
			decision.NextPhase = PhaseAudit
			decision.Reason = "Foundation meets threshold"
			decision.ShouldExecute = true
		} else {
			decision.ShouldExecute = false
			decision.Reason = fmt.Sprintf("Foundation at %.2f, need %.2f",
				e.state.Confidence.Foundation, e.state.Confidence.Threshold)
		}

	case PhaseAudit:
		// Audit complete
		if e.state.EvidenceState.Completeness >= 0.7 {
			decision.NextPhase = PhaseAssessment
			decision.Reason = "Audit complete"
			decision.ShouldExecute = true
		}

	case PhaseAssessment:
		// Assessment complete
		decision.NextPhase = PhaseArchitecture
		decision.Reason = "Assessment complete"
		decision.ShouldExecute = true

	case PhaseArchitecture:
		// Architecture complete
		decision.NextPhase = PhaseImplementation
		decision.Reason = "Architecture complete"
		decision.ShouldExecute = true

	case PhaseImplementation:
		// Implementation complete
		if e.state.Confidence.MeetsThreshold {
			decision.NextPhase = PhaseComplete
			decision.Reason = "Implementation verified"
			decision.ShouldExecute = true
		} else {
			decision.NextPhase = PhaseFoundation
			decision.Reason = "Need more work"
			decision.ShouldExecute = true
		}

	case PhaseComplete, PhaseBlocked:
		decision.ShouldExecute = false
		decision.Reason = "Session ending"
	}
	return decision
}

// executePhase executes a single phase
func (e *Engine) executePhase(phase OrchestrationPhase) error {
	if phase == e.state.CurrentPhase {
		// Already in this phase
		return nil
	}

	// Record transition
	transition := PhaseTransition{
		From:        e.state.CurrentPhase,
		To:          phase,
		Reason:      "Phase decision",
		Confidence:  e.state.Confidence.Overall,
		Timestamp:   time.Now(),
	}
	e.state.History = append(e.state.History, transition)

	// Update phase
	e.state.PreviousPhase = e.state.CurrentPhase
	e.state.CurrentPhase = phase

	// Update metrics
	phaseName := string(phase)
	e.state.Metrics.PhaseExecutions[phaseName]++

	return nil
}

// reevaluate performs a final re-evaluation after phase execution
func (e *Engine) reevaluate() error {
	// Re-evaluate everything after phase change
	if err := e.evaluateCurrentState(); err != nil {
		return err
	}

	if err := e.evaluateConfidence(); err != nil {
		return err
	}

	return e.evaluateMissingEvidence()
}

// shouldContinue determines if the session should continue
func (e *Engine) shouldContinue() bool {
	if e.state == nil {
		return false
	}

	// Check max cycles
	if e.state.Metrics.CycleCount >= e.config.MaxCycles {
		return false
	}

	// Check if complete
	if e.state.CurrentPhase == PhaseComplete {
		return false
	}

	return true
}

// GetState returns the current orchestration state
func (e *Engine) GetState() *OrchestrationState {
	return e.state
}

// IsSessionActive returns whether a session is active
func (e *Engine) IsSessionActive() bool {
	return e.sessionActive
}

// CanImplement returns whether implementation is allowed
func (e *Engine) CanImplement() bool {
	return e.confidence.CanImplement(&e.state.Confidence)
}

// GetFoundationStatus returns Foundation readiness status
func (e *Engine) GetFoundationStatus() (ready bool, confidence float64, missing []string, err error) {
	ready, confidence, err = e.confidence.IsFoundationReady(&e.state.Workspace)
	if err != nil {
		return
	}

	assessment, err := e.confidence.AssessFoundationCompleteness(&e.state.Workspace)
	if err != nil {
		return
	}

	missing = assessment.Missing
	return
}

// MigrateToProject migrates temporary workspace to project workspace
func (e *Engine) MigrateToProject(projectPath string) error {
	if e.state == nil || e.state.Workspace.WorkspaceType != WorkspaceTypeTemporary {
		return nil // Nothing to migrate
	}

	return e.resolver.MigrateToProject(e.state.Workspace.KDSEPath, projectPath)
}

// CreateTemporaryWorkspace creates a temporary workspace for the session
func (e *Engine) CreateTemporaryWorkspace(projectName string) error {
	workspace, err := e.resolver.ResolveTemporaryWorkspace(projectName)
	if err != nil {
		return err
	}

	e.state.Workspace = *workspace
	return nil
}

// SaveState persists the orchestration state
func (e *Engine) SaveState() error {
	if e.state == nil {
		return nil
	}

	statePath := filepath.Join(e.state.Workspace.KDSEPath, "orchestration-state.json")
	return writeJSON(statePath, e.state)
}

// LoadState loads a persisted orchestration state
func (e *Engine) LoadState(workspacePath string) error {
	resolved, err := e.resolver.ResolveWorkspace(workspacePath)
	if err != nil {
		return err
	}

	statePath := filepath.Join(resolved.KDSEPath, "orchestration-state.json")
	
	state := &OrchestrationState{}
	if err := readJSON(statePath, state); err != nil {
		return err
	}

	e.state = state
	e.sessionActive = true
	return nil
}

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	return fmt.Sprintf("KDSE-ORCH-%s", time.Now().Format("20060102-150405"))
}

// Helper functions for JSON serialization
func writeJSON(path string, v interface{}) error {
	data, err := jsonMarshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return jsonUnmarshal(data, v)
}

// jsonMarshal wraps encoding/json.Marshal with error handling
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
