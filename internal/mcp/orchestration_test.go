package orchestration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_Initialize(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	m := NewManager(tmpDir)

	// Test initialization
	state, err := m.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Verify initial state
	if state.CurrentPhase != PhaseIdle {
		t.Errorf("Expected initial phase to be PhaseIdle, got %s", state.CurrentPhase)
	}
	if state.Confidence != 0.0 {
		t.Errorf("Expected initial confidence to be 0.0, got %f", state.Confidence)
	}
	if state.ExecutionMode != ModeStrict {
		t.Errorf("Expected execution mode to be ModeStrict, got %s", state.ExecutionMode)
	}
	if len(state.NextAllowedPhases) != 1 || state.NextAllowedPhases[0] != PhaseProblem {
		t.Errorf("Expected NextAllowedPhases to contain only PhaseProblem, got %v", state.NextAllowedPhases)
	}

	// Verify state was persisted
	statePath := filepath.Join(tmpDir, ".kdse", "orchestration_state.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("State file was not created")
	}
}

func TestManager_TransitionTo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	m := NewManager(tmpDir)
	m.Initialize()

	// Test valid transition: Idle -> Problem
	state, err := m.TransitionTo(PhaseProblem, nil)
	if err != nil {
		t.Fatalf("Transition to PhaseProblem failed: %v", err)
	}
	if state.CurrentPhase != PhaseProblem {
		t.Errorf("Expected phase to be PhaseProblem, got %s", state.CurrentPhase)
	}

	// Test invalid transition: Problem -> Complete (should fail)
	_, err = m.TransitionTo(PhaseComplete, nil)
	if err == nil {
		t.Error("Expected error for invalid transition, got nil")
	}

	// Test valid transition chain
	m.TransitionTo(PhaseKnowledge, nil)
	m.TransitionTo(PhaseFoundation, nil)
	m.TransitionTo(PhaseAudit, nil)

	// Verify completed phases
	state, _ = m.Load()
	if len(state.CompletedPhases) != 4 {
		t.Errorf("Expected 4 completed phases, got %d", len(state.CompletedPhases))
	}
}

func TestManager_UpdateConfidence(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "kdse-test-*")
	defer os.RemoveAll(tmpDir)

	m := NewManager(tmpDir)
	m.Initialize()

	// Test confidence update
	state, err := m.UpdateConfidence(0.75)
	if err != nil {
		t.Fatalf("UpdateConfidence failed: %v", err)
	}
	if state.Confidence != 0.75 {
		t.Errorf("Expected confidence to be 0.75, got %f", state.Confidence)
	}

	// Test confidence clamping
	state, _ = m.UpdateConfidence(1.5)
	if state.Confidence != 1.0 {
		t.Errorf("Expected confidence to be clamped to 1.0, got %f", state.Confidence)
	}

	state, _ = m.UpdateConfidence(-0.5)
	if state.Confidence != 0.0 {
		t.Errorf("Expected confidence to be clamped to 0.0, got %f", state.Confidence)
	}
}

func TestManager_GetExecutionDecision(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "kdse-test-*")
	defer os.RemoveAll(tmpDir)

	m := NewManager(tmpDir)
	m.Initialize()

	// Test with no objective set
	decision := m.GetExecutionDecision("Build inventory system")
	if decision.Action != "set_objective" {
		t.Errorf("Expected action 'set_objective', got '%s'", decision.Action)
	}

	// Set objective
	m.SetObjective("Build inventory system")

	// Test phase progression
	decision = m.GetExecutionDecision("Build inventory system")
	if decision.Action != "problem" {
		t.Errorf("Expected action 'problem', got '%s'", decision.Action)
	}

	// Transition to Problem
	m.TransitionTo(PhaseProblem, nil)
	decision = m.GetExecutionDecision("Build inventory system")
	if decision.Action != "knowledge" {
		t.Errorf("Expected action 'knowledge', got '%s'", decision.Action)
	}
}

func TestManager_CheckImplementationReadiness(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "kdse-test-*")
	defer os.RemoveAll(tmpDir)

	m := NewManager(tmpDir)
	m.Initialize()
	m.SetObjective("Build inventory system")

	// Before Architecture phase
	ready, _ := m.CheckImplementationReadiness()
	if ready {
		t.Error("Should not be ready for implementation before Architecture phase")
	}

	// Progress to Architecture
	m.TransitionTo(PhaseProblem, nil)
	m.TransitionTo(PhaseKnowledge, nil)
	m.TransitionTo(PhaseFoundation, nil)
	m.TransitionTo(PhaseAudit, nil)
	m.TransitionTo(PhaseAssessment, nil)
	m.TransitionTo(PhaseArchitecture, nil)

	// Check readiness with workspace state
	ws := &WorkspaceState{
		Initialized:    true,
		HasFoundation:  true,
		HasAuditReport: true,
	}
	m.UpdateWorkspace(ws)
	m.UpdateConfidence(0.92)

	ready, _ = m.CheckImplementationReadiness()
	if !ready {
		t.Error("Should be ready for implementation after completing all prerequisites")
	}
}

func TestManager_CanProceedTo(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "kdse-test-*")
	defer os.RemoveAll(tmpDir)

	m := NewManager(tmpDir)
	m.Initialize()
	m.SetObjective("Build system")

	// Should be able to proceed to Problem
	canProceed, _ := m.CanProceedTo(PhaseProblem)
	if !canProceed {
		t.Error("Should be able to proceed to PhaseProblem")
	}

	// Should NOT be able to jump to Implementation
	canProceed, reason := m.CanProceedTo(PhaseImplementation)
	if canProceed {
		t.Error("Should NOT be able to proceed directly to PhaseImplementation")
	}

	// Verify blocking reason is provided
	if reason == "" {
		t.Error("Expected blocking reason to be provided")
	}
}

func TestPhaseTransitions(t *testing.T) {
	// Test that all phases have valid transitions defined
	for phase := range PhaseTransitions {
		if len(PhaseTransitions[phase]) == 0 {
			t.Errorf("Phase %s has no transitions defined", phase)
		}
	}
}

func TestPhaseConfidenceThreshold(t *testing.T) {
	// Test that confidence thresholds are reasonable
	for phase, threshold := range PhaseConfidenceThreshold {
		if threshold < 0 || threshold > 1 {
			t.Errorf("Invalid confidence threshold for phase %s: %f", phase, threshold)
		}
		// Implementation should have highest threshold
		if phase == PhaseImplementation && threshold < 0.9 {
			t.Errorf("Implementation phase should have high confidence threshold, got %f", threshold)
		}
	}
}
