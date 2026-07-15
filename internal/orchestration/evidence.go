package orchestration

import (
	"os"
	"path/filepath"
)

// EvidenceEvaluator evaluates evidence requirements for each phase
type EvidenceEvaluator struct {
	config *EngineConfig
}

// NewEvidenceEvaluator creates a new evidence evaluator
func NewEvidenceEvaluator(config *EngineConfig) *EvidenceEvaluator {
	if config == nil {
		config = DefaultEngineConfig()
	}
	return &EvidenceEvaluator{config: config}
}

// GetRequiredEvidenceForPhase returns evidence requirements for a specific phase
func (e *EvidenceEvaluator) GetRequiredEvidenceForPhase(phase OrchestrationPhase) []EvidenceRequirement {
	switch phase {
	case PhaseResolve:
		return []EvidenceRequirement{
			{ID: "repo-root", Type: "repository", Description: "Repository root accessible", Paths: []string{"README.md", "go.mod", "package.json"}, Weight: 1.0, Critical: true},
		}
	case PhaseAssess:
		return []EvidenceRequirement{
			{ID: "knowledge-artifacts", Type: "documentation", Description: "Knowledge artifacts (README, SPEC)", Paths: []string{"README.md", "SPEC.md"}, Weight: 1.0, Critical: true},
			{ID: "architecture", Type: "documentation", Description: "Architecture documentation", Paths: []string{"ARCHITECTURE.md", "architecture/"}, Weight: 0.7, Critical: false},
		}
	case PhaseFoundation:
		return []EvidenceRequirement{
			{ID: "foundation-readme", Type: "foundation", Description: "Foundation README", Paths: []string{".kdse/foundation/README.md"}, Weight: 1.0, Critical: true},
			{ID: "vision", Type: "foundation", Description: "Vision document", Paths: []string{".kdse/foundation/001-vision.md"}, Weight: 0.8, Critical: true},
			{ID: "principles", Type: "foundation", Description: "Principles document", Paths: []string{".kdse/foundation/002-principles.md"}, Weight: 0.8, Critical: true},
			{ID: "standards", Type: "foundation", Description: "Standards document", Paths: []string{".kdse/foundation/003-standards.md"}, Weight: 0.9, Critical: true},
			{ID: "engineering-model", Type: "foundation", Description: "Engineering model", Paths: []string{".kdse/foundation/004-engineering-model.md"}, Weight: 1.0, Critical: true},
			{ID: "chain-of-authority", Type: "foundation", Description: "Chain of authority", Paths: []string{".kdse/foundation/006-chain-of-authority.md"}, Weight: 0.8, Critical: false},
		}
	case PhaseCollect:
		return []EvidenceRequirement{
			{ID: "evidence-directory", Type: "evidence", Description: "Evidence collection directory", Paths: []string{".kdse/evidence/"}, Weight: 1.0, Critical: true},
			{ID: "screenshots", Type: "evidence", Description: "Screenshot evidence", Paths: []string{".kdse/evidence/screenshots/"}, Weight: 0.5, Critical: false},
			{ID: "tests", Type: "evidence", Description: "Test evidence", Paths: []string{".kdse/evidence/tests/"}, Weight: 0.5, Critical: false},
		}
	case PhaseAnalyze:
		return []EvidenceRequirement{
			{ID: "collected-evidence", Type: "evidence", Description: "Previously collected evidence", Paths: []string{".kdse/evidence/"}, Weight: 1.0, Critical: true},
			{ID: "audit-results", Type: "audit", Description: "Audit results", Paths: []string{".kdse/reports/audit-*.md"}, Weight: 0.7, Critical: false},
		}
	case PhaseDesign:
		return []EvidenceRequirement{
			{ID: "analysis-results", Type: "analysis", Description: "Analysis results", Paths: []string{".kdse/knowledge/"}, Weight: 1.0, Critical: true},
			{ID: "architecture", Type: "architecture", Description: "Architecture design", Paths: []string{"ARCHITECTURE.md", ".kdse/foundation/005-architecture.md"}, Weight: 0.8, Critical: true},
		}
	case PhaseImplement:
		return []EvidenceRequirement{
			{ID: "design-spec", Type: "specification", Description: "Design specification", Paths: []string{".kdse/artifacts/design.md"}, Weight: 1.0, Critical: true},
			{ID: "context-handoff", Type: "context", Description: "Context handoff document", Paths: []string{".kdse/context.json"}, Weight: 0.8, Critical: true},
		}
	case PhaseVerify:
		return []EvidenceRequirement{
			{ID: "implementation", Type: "implementation", Description: "Implementation artifacts", Paths: []string{"src/", "lib/", "main.go"}, Weight: 1.0, Critical: true},
			{ID: "tests", Type: "tests", Description: "Test files", Paths: []string{"tests/", "*_test.go"}, Weight: 0.8, Critical: true},
			{ID: "documentation", Type: "documentation", Description: "Updated documentation", Paths: []string{"docs/", "README.md"}, Weight: 0.6, Critical: false},
		}
	default:
		return []EvidenceRequirement{}
	}
}

// EvaluateEvidence evaluates evidence state for the current phase
func (e *EvidenceEvaluator) EvaluateEvidence(workspace *WorkspaceInfo, phase OrchestrationPhase) (*EvidenceState, error) {
	requirements := e.GetRequiredEvidenceForPhase(phase)
	
	state := &EvidenceState{
		Required:  requirements,
		Present:   []string{},
		Missing:   []string{},
		TotalRequired: len(requirements),
		TotalPresent:  0,
		Completeness: 0.0,
	}

	basePath := workspace.ResolvedPath
	if workspace.ProjectPath != "" {
		basePath = workspace.ProjectPath
	}

	for i := range requirements {
		req := &requirements[i]
		if satisfied := e.checkRequirement(basePath, *req); satisfied {
			state.Present = append(state.Present, req.ID)
			req.Satisfied = true
			state.TotalPresent++
		} else {
			state.Missing = append(state.Missing, req.ID)
			req.Satisfied = false
		}
	}

	if state.TotalRequired > 0 {
		state.Completeness = float64(state.TotalPresent) / float64(state.TotalRequired)
	}

	return state, nil
}

// checkRequirement verifies if a single evidence requirement is satisfied
func (e *EvidenceEvaluator) checkRequirement(basePath string, req EvidenceRequirement) bool {
	for _, pathPattern := range req.Paths {
		fullPath := filepath.Join(basePath, pathPattern)
		if info, err := os.Stat(fullPath); err == nil {
			// If it's a directory, check if it has contents
			if info.IsDir() {
				entries, err := os.ReadDir(fullPath)
				if err == nil && len(entries) > 0 {
					return true
				}
			} else {
				return true
			}
		}
	}
	return false
}

// GetMissingEvidence returns list of missing evidence for a phase
func (e *EvidenceEvaluator) GetMissingEvidence(workspace *WorkspaceInfo, phase OrchestrationPhase) []EvidenceRequirement {
	state, err := e.EvaluateEvidence(workspace, phase)
	if err != nil {
		return e.GetRequiredEvidenceForPhase(phase)
	}

	var missing []EvidenceRequirement
	for _, req := range state.Required {
		if !req.Satisfied {
			missing = append(missing, req)
		}
	}
	return missing
}

// GetCriticalMissingEvidence returns only critical missing evidence
func (e *EvidenceEvaluator) GetCriticalMissingEvidence(workspace *WorkspaceInfo, phase OrchestrationPhase) []EvidenceRequirement {
	allMissing := e.GetMissingEvidence(workspace, phase)
	
	var critical []EvidenceRequirement
	for _, req := range allMissing {
		if req.Critical {
			critical = append(critical, req)
		}
	}
	return critical
}

// CanProceedToPhase checks if evidence requirements are met for the target phase
func (e *EvidenceEvaluator) CanProceedToPhase(workspace *WorkspaceInfo, targetPhase OrchestrationPhase) (bool, []string) {
	state, err := e.EvaluateEvidence(workspace, targetPhase)
	if err != nil {
		return false, []string{"Failed to evaluate evidence"}
	}

	// Check critical requirements
	var blockedReasons []string
	for _, req := range state.Required {
		if req.Critical && !req.Satisfied {
			blockedReasons = append(blockedReasons, req.Description)
		}
	}

	return len(blockedReasons) == 0, blockedReasons
}

// EvaluateAllPhases returns evidence state for all phases
func (e *EvidenceEvaluator) EvaluateAllPhases(workspace *WorkspaceInfo) map[OrchestrationPhase]*EvidenceState {
	phases := []OrchestrationPhase{
		PhaseResolve,
		PhaseAssess,
		PhaseFoundation,
		PhaseCollect,
		PhaseAnalyze,
		PhaseDesign,
		PhaseImplement,
		PhaseVerify,
	}

	results := make(map[OrchestrationPhase]*EvidenceState)
	for _, phase := range phases {
		state, err := e.EvaluateEvidence(workspace, phase)
		if err != nil {
			results[phase] = &EvidenceState{TotalRequired: 0, Completeness: 0.0}
		} else {
			results[phase] = state
		}
	}

	return results
}

// GetEvidenceSummary returns a summary of evidence across all phases
func (e *EvidenceEvaluator) GetEvidenceSummary(workspace *WorkspaceInfo) *EvidenceSummary {
	allStates := e.EvaluateAllPhases(workspace)
	
	summary := &EvidenceSummary{
		PhaseStates:      allStates,
		TotalRequired:    0,
		TotalPresent:     0,
		CriticalMissing:   []string{},
		CompletenessByPhase: map[OrchestrationPhase]float64{},
	}

	for phase, state := range allStates {
		summary.TotalRequired += state.TotalRequired
		summary.TotalPresent += state.TotalPresent
		summary.CompletenessByPhase[phase] = state.Completeness

		for _, req := range state.Required {
			if req.Critical && !req.Satisfied {
				found := false
				for _, existing := range summary.CriticalMissing {
					if existing == req.Description {
						found = true
						break
					}
				}
				if !found {
					summary.CriticalMissing = append(summary.CriticalMissing, req.Description)
				}
			}
		}
	}

	if summary.TotalRequired > 0 {
		summary.OverallCompleteness = float64(summary.TotalPresent) / float64(summary.TotalRequired)
	}

	return summary
}

// EvidenceSummary contains aggregated evidence information
type EvidenceSummary struct {
	PhaseStates          map[OrchestrationPhase]*EvidenceState `json:"phase_states"`
	TotalRequired        int                                   `json:"total_required"`
	TotalPresent         int                                   `json:"total_present"`
	OverallCompleteness  float64                               `json:"overall_completeness"`
	CriticalMissing      []string                              `json:"critical_missing"`
	CompletenessByPhase  map[OrchestrationPhase]float64        `json:"completeness_by_phase"`
}

// GenerateEvidenceCollectionPlan creates a plan for collecting missing evidence
func (e *EvidenceEvaluator) GenerateEvidenceCollectionPlan(workspace *WorkspaceInfo, targetPhase OrchestrationPhase) *EvidenceCollectionPlan {
	missing := e.GetMissingEvidence(workspace, targetPhase)
	
	plan := &EvidenceCollectionPlan{
		TargetPhase: targetPhase,
		Steps:       []EvidenceCollectionStep{},
	}

	for _, req := range missing {
		step := EvidenceCollectionStep{
			RequirementID: req.ID,
			Description:   req.Description,
			Type:          req.Type,
			Paths:         req.Paths,
			Priority:      "high",
			Critical:      req.Critical,
		}

		if req.Critical {
			step.Priority = "critical"
		} else {
			step.Priority = "medium"
		}

		plan.Steps = append(plan.Steps, step)
	}

	return plan
}

// EvidenceCollectionPlan defines how to collect missing evidence
type EvidenceCollectionPlan struct {
	TargetPhase OrchestrationPhase        `json:"target_phase"`
	Steps       []EvidenceCollectionStep  `json:"steps"`
}

// EvidenceCollectionStep is a single evidence collection action
type EvidenceCollectionStep struct {
	RequirementID string   `json:"requirement_id"`
	Description   string   `json:"description"`
	Type          string   `json:"type"`
	Paths         []string `json:"paths"`
	Priority      string   `json:"priority"`
	Critical      bool     `json:"critical"`
}
