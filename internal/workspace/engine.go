package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kdse/runtime/internal/methodology/lifecycle"
)

// Use lifecycle.Phase as the authoritative Phase type
type Phase = lifecycle.Phase

// Phase constants from lifecycle package
const (
	PhaseInitialization  = lifecycle.PhaseInitialization
	PhaseKnowledge      = lifecycle.PhaseKnowledge
	PhaseArchitecture   = lifecycle.PhaseArchitecture
	PhaseImplementation = lifecycle.PhaseImplementation
	PhaseVerification   = lifecycle.PhaseVerification
	PhaseReports        = lifecycle.PhaseReports
)

// Engine is the main interface for the Workspace Engine
// It owns all project state and enforces the KDSE methodology
type Engine interface {
	// Workspace lifecycle
	InitializeWorkspace(ctx context.Context, opts InitOptions) (*RuntimeContext, error)
	VerifyWorkspace(ctx context.Context) (*VerificationResult, error)
	LoadWorkspace(ctx context.Context) (*RuntimeContext, error)
	DestroyWorkspace(ctx context.Context) error

	// Phase management
	GetPhase(ctx context.Context) (Phase, error)
	AdvancePhase(ctx context.Context, target Phase) (*Transition, error)
	GetPhaseHistory(ctx context.Context) ([]PhaseTransition, error)

	// Artifact management
	GetArtifacts(ctx context.Context, phase Phase) ([]Artifact, error)
	VerifyArtifacts(ctx context.Context, phase Phase) (*ValidationResult, error)
	ValidateArtifact(ctx context.Context, artifact Artifact) error

	// Reporting
	GenerateReport(ctx context.Context, opts ReportOptions) (*Report, error)
	ListReports(ctx context.Context) ([]ReportSummary, error)

	// Session management
	CreateSession(ctx context.Context) (*Session, error)
	GetSession(ctx context.Context) (*Session, error)
	EndSession(ctx context.Context) error

	// Knowledge management
	CollectKnowledge(ctx context.Context) (*KnowledgeCollection, error)
}

// DefaultEngine implements the Engine interface
type DefaultEngine struct {
	workspacePath string
	lifecycle     lifecycle.Lifecycle
}

// NewEngine creates a new workspace engine
func NewEngine(workspacePath string) *DefaultEngine {
	return &DefaultEngine{
		workspacePath: workspacePath,
		lifecycle:     lifecycle.NewLifecycle(),
	}
}

// VerifyWorkspace verifies the workspace state
// This is the verification gate that MUST be called before any engineering action
func (e *DefaultEngine) VerifyWorkspace(ctx context.Context) (*VerificationResult, error) {
	result := &VerificationResult{
		Valid:     true,
		Timestamp: time.Now(),
	}

	// Step 1: Check .kdse exists (PRIMARY EVIDENCE)
	kdseDir := filepath.Join(e.workspacePath, ".kdse")
	if _, err := os.Stat(kdseDir); os.IsNotExist(err) {
		return &VerificationResult{
			Valid: false,
			Errors: []VerificationError{{
				Code:    "KDSE_MISSING",
				Message: ".kdse directory not found",
			}},
			Timestamp: time.Now(),
		}, ErrRuntimeMissing
	}

	// Step 2: Load and verify runtime configuration
	runtimeConfig, err := e.loadRuntimeConfig()
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, VerificationError{
			Code:    "RUNTIME_INVALID",
			Message: err.Error(),
		})
		result.RuntimeInfo = nil
		return result, nil
	}

	// Step 3: Load current phase
	phase, err := e.loadPhase()
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, VerificationError{
			Code:    "PHASE_INVALID",
			Message: err.Error(),
		})
	} else {
		result.Phase = phase
	}

	// Step 4: Verify required artifacts for current phase
	artifacts, err := e.loadArtifactPaths(phase)
	if err == nil {
		phaseValidation, err := e.lifecycle.ValidatePhase(ctx, phase, artifacts)
		if err == nil && !phaseValidation.Valid {
			result.Valid = false
			for _, vErr := range phaseValidation.Errors {
				result.Errors = append(result.Errors, VerificationError{
					Code:    vErr.Type,
					Message: vErr.Message,
					Path:    vErr.Path,
				})
			}
		}
	}

	result.RuntimeInfo = &RuntimeInfo{
		Type:    runtimeConfig.Type,
		Version: runtimeConfig.Version,
		Commit:  runtimeConfig.Commit,
	}
	return result, nil
}

// InitializeWorkspace initializes a new KDSE workspace
// This creates both the standard project layout AND the KDSE runtime
func (e *DefaultEngine) InitializeWorkspace(ctx context.Context, opts InitOptions) (*RuntimeContext, error) {
	// Create staging directory for atomic initialization
	stagingDir, err := os.MkdirTemp("", "kdse-init-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(stagingDir)

	// Step 1: Create standard project layout in staging
	// This ensures KDSE augments the project, not replaces it
	if err := e.createProjectLayout(stagingDir); err != nil {
		return nil, err
	}

	// Step 2: Create .kdse structure in staging
	stagingKDSE := filepath.Join(stagingDir, ".kdse")
	if err := os.MkdirAll(stagingKDSE, 0755); err != nil {
		return nil, err
	}

	// Create required runtime files
	if err := e.createRuntimeFiles(stagingKDSE, opts); err != nil {
		return nil, err
	}

	// Initialize phase to initialization
	if err := e.initializePhase(stagingKDSE); err != nil {
		return nil, err
	}

	// Create KDSE runtime directories (runtime layer)
	if err := e.createRuntimeDirectories(stagingKDSE); err != nil {
		return nil, err
	}

	// Atomic move to final location
	if err := atomicMove(stagingDir, e.workspacePath); err != nil {
		return nil, err
	}

	// CRITICAL: Verify runtime before returning success
	// This implements the verification-first principle
	verifyResult, err := e.VerifyWorkspace(ctx)
	if err != nil || !verifyResult.Valid {
		// Rollback on verification failure
		os.RemoveAll(filepath.Join(e.workspacePath, ".kdse"))
		return nil, ErrVerificationFailed
	}

	return e.LoadWorkspace(ctx)
}

// createProjectLayout creates the standard project directory structure
// This is the project layer - owned by the software project, not KDSE
func (e *DefaultEngine) createProjectLayout(stagingDir string) error {
	// Project directories - these belong to the project, not KDSE
	projectDirs := []string{
		"docs",
		"docs/architecture",
		"docs/api",
		"docs/deployment",
		"docs/design",
		"src",
		"tests",
	}

	for _, dir := range projectDirs {
		fullPath := filepath.Join(stagingDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create README.md if it doesn't exist in project root
	readmePath := filepath.Join(stagingDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		readme := "# Project\n\nThis is a KDSE-enabled software project.\n\n## Project Structure\n\n- docs/ - Project documentation\n- src/ - Source code\n- tests/ - Test code\n\n## KDSE Runtime\n\nThe .kdse/ directory contains the KDSE engineering runtime.\nSee docs/architecture/OWNERSHIP_MODEL.md for details.\n\n## Quick Start\n\n1. Install dependencies\n2. Build the project\n3. Run tests\n"
		if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
			return err
		}
	}

	// Create README.md for project directories (idempotent)
	for _, dir := range projectDirs {
		readmePath := filepath.Join(stagingDir, dir, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			dirName := filepath.Base(dir)
			readme := fmt.Sprintf(`# %s

This directory is part of the project layer (owned by the software project).

`, dirName)
			if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

// createRuntimeDirectories creates KDSE runtime layer directories
// These belong to KDSE, not the project
func (e *DefaultEngine) createRuntimeDirectories(stagingKDSE string) error {
	// Runtime directories - these belong to KDSE, not the project
	runtimeDirs := []string{
		"runtime",
		"sessions",
		"state",
		"cache",
		"reports",
		"evidence",
		"traceability",
		"references",
		"references/modbus",
		"references/iec61850",
		"references/vendor",
		"knowledge",
		"laboratory",
	}

	for _, dir := range runtimeDirs {
		fullPath := filepath.Join(stagingKDSE, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create README files in runtime directories (idempotent)
	for _, dir := range runtimeDirs {
		readmePath := filepath.Join(stagingKDSE, dir, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			dirName := filepath.Base(dir)
			var readme string
			switch dirName {
			case "references":
				readme = `# References

This directory contains external authoritative references.
These are external sources, NOT project documentation.

`
			case "knowledge":
				readme = `# Knowledge

This directory contains knowledge extracted from references.
Knowledge must maintain traceability back to its references.

`
			case "laboratory":
				readme = `# Laboratory

This directory contains the engineering laboratory for experiments.
This is a runtime artifact, not project documentation.

`
			default:
				readme = fmt.Sprintf(`# %s

This directory is part of the KDSE runtime layer.

`, dirName)
			}
			if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadWorkspace loads the current workspace state
func (e *DefaultEngine) LoadWorkspace(ctx context.Context) (*RuntimeContext, error) {
	kdseDir := filepath.Join(e.workspacePath, ".kdse")

	// Verify .kdse exists
	if _, err := os.Stat(kdseDir); os.IsNotExist(err) {
		return nil, ErrRuntimeMissing
	}

	// Load runtime config
	runtimeConfig, err := e.loadRuntimeConfig()
	if err != nil {
		return nil, err
	}

	// Load phase
	phase, err := e.loadPhase()
	if err != nil {
		return nil, err
	}

	// Load workspace state
	state := &WorkspaceState{
		CurrentPhase:     phase,
		LastVerification: time.Now(),
	}

	// Load phase history
	history, _ := e.loadPhaseHistory()

	// Count artifacts
	artifactCount := e.countArtifacts()

	// Count reports
	reportCount := e.countReports()

	state.PhaseHistory = history
	state.ArtifactCount = artifactCount
	state.ReportCount = reportCount

	return &RuntimeContext{
		Path:    e.workspacePath,
		Root:    e.workspacePath,
		Config:  &WorkspaceConfig{Version: "1.0.0"},
		State:   state,
		Runtime: runtimeConfig,
	}, nil
}

// DestroyWorkspace removes the workspace
func (e *DefaultEngine) DestroyWorkspace(ctx context.Context) error {
	kdseDir := filepath.Join(e.workspacePath, ".kdse")
	return os.RemoveAll(kdseDir)
}

// GetPhase returns the current phase
func (e *DefaultEngine) GetPhase(ctx context.Context) (Phase, error) {
	return e.loadPhase()
}

// AdvancePhase advances to the next phase
func (e *DefaultEngine) AdvancePhase(ctx context.Context, target Phase) (*Transition, error) {
	// Get current phase
	current, err := e.loadPhase()
	if err != nil {
		return nil, err
	}

	// Validate transition
	if !e.lifecycle.IsValidTransition(current, target) {
		return nil, &EngineError{
			Code:        "INVALID_TRANSITION",
			Message:     "Invalid phase transition",
			Details:     map[string]interface{}{"from": current, "to": target},
			Remediation: "Use 'kdse phase show' to see valid transitions",
		}
	}

	// Verify current phase completion
	result, err := e.VerifyArtifacts(ctx, current)
	if err != nil {
		return nil, err
	}
	if !result.Valid {
		return nil, &EngineError{
			Code:        "INCOMPLETE_PHASE",
			Message:     "Current phase has incomplete artifacts",
			Details:     map[string]interface{}{"phase": current, "missing": result.Missing},
			Remediation: "Complete required artifacts before advancing",
		}
	}

	// Create transition record
	transition := &Transition{
		From:      current,
		To:        target,
		Timestamp: time.Now(),
		Verified:  true,
	}

	// Persist new phase
	if err := e.persistPhase(target); err != nil {
		return nil, err
	}

	// Record transition in history
	if err := e.recordTransition(transition); err != nil {
		return nil, err
	}

	return transition, nil
}

// GetPhaseHistory returns the phase transition history
func (e *DefaultEngine) GetPhaseHistory(ctx context.Context) ([]PhaseTransition, error) {
	return e.loadPhaseHistory()
}

// GetArtifacts returns artifacts for a phase
func (e *DefaultEngine) GetArtifacts(ctx context.Context, phase Phase) ([]Artifact, error) {
	kdseDir := filepath.Join(e.workspacePath, ".kdse")
	dirMap := map[Phase]string{
		PhaseInitialization: kdseDir,
		PhaseKnowledge:      filepath.Join(kdseDir, "knowledge"),
		PhaseArchitecture:   filepath.Join(kdseDir, "architecture"),
		PhaseImplementation: filepath.Join(kdseDir, "implementation"),
		PhaseVerification:   filepath.Join(kdseDir, "verification"),
		PhaseReports:        filepath.Join(kdseDir, "reports"),
	}

	dir, ok := dirMap[phase]
	if !ok {
		return nil, ErrPhaseInvalid
	}

	return e.listArtifacts(dir)
}

// VerifyArtifacts verifies artifacts for a phase
func (e *DefaultEngine) VerifyArtifacts(ctx context.Context, phase Phase) (*ValidationResult, error) {
	// Get required artifacts from lifecycle
	required := e.lifecycle.GetRequiredArtifacts(phase)

	// Get existing artifacts
	existing, err := e.GetArtifacts(ctx, phase)
	if err != nil {
		return nil, err
	}

	result := &ValidationResult{
		Valid:    true,
		Verified: []Artifact{},
		Missing:  []ArtifactSpec{},
		Invalid:  []Artifact{},
	}

	// Build existing path map
	existingPaths := make(map[string]Artifact)
	for _, art := range existing {
		existingPaths[art.Path] = art
	}

	// Check required artifacts
	for _, spec := range required {
		if art, exists := existingPaths[spec.Path]; exists {
			result.Verified = append(result.Verified, art)
		} else if spec.Required {
			result.Missing = append(result.Missing, ArtifactSpec{
				Path:     spec.Path,
				Type:     spec.Type,
				Required: spec.Required,
			})
			result.Valid = false
		}
	}

	return result, nil
}

// ValidateArtifact validates a single artifact
func (e *DefaultEngine) ValidateArtifact(ctx context.Context, artifact Artifact) error {
	path := filepath.Join(e.workspacePath, artifact.Path)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ErrArtifactMissing
	}
	if err != nil {
		return err
	}

	if info.IsDir() {
		return ErrArtifactInvalid
	}

	// Check size
	if info.Size() == 0 {
		return ErrArtifactInvalid
	}

	return nil
}

// GenerateReport generates a report
func (e *DefaultEngine) GenerateReport(ctx context.Context, opts ReportOptions) (*Report, error) {
	ws, err := e.LoadWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	var content string
	switch opts.Type {
	case ReportTypePhase:
		content = e.formatPhaseReport(ws)
	case ReportTypeVerification:
		content = e.formatVerificationReport(ws)
	case ReportTypeSummary:
		content = e.formatSummaryReport(ws)
	case ReportTypeProgress:
		content = e.formatProgressReport(ws)
	default:
		return nil, ErrVerificationFailed
	}

	// Save report if output path specified
	if opts.Output != "" {
		if err := os.WriteFile(opts.Output, []byte(content), 0644); err != nil {
			return nil, err
		}
	}

	return &Report{
		Type:      opts.Type,
		Title:     string(opts.Type) + " Report",
		Content:   content,
		Generated: time.Now(),
	}, nil
}

// ListReports returns all generated reports
func (e *DefaultEngine) ListReports(ctx context.Context) ([]ReportSummary, error) {
	reportsDir := filepath.Join(e.workspacePath, ".kdse", "reports")

	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		return nil, err
	}

	var reports []ReportSummary
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			info, _ := entry.Info()
			reports = append(reports, ReportSummary{
				Type:      ReportTypeSummary,
				Title:     entry.Name(),
				Generated: info.ModTime(),
			})
		}
	}

	return reports, nil
}

// CreateSession creates a new session
func (e *DefaultEngine) CreateSession(ctx context.Context) (*Session, error) {
	session := &Session{
		ID:       generateSessionID(),
		Created:  time.Now(),
		Phase:    PhaseInitialization,
		Metadata: make(map[string]string),
	}

	if err := e.persistSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession returns the current session
func (e *DefaultEngine) GetSession(ctx context.Context) (*Session, error) {
	return e.loadSession()
}

// EndSession ends the current session
func (e *DefaultEngine) EndSession(ctx context.Context) error {
	session, err := e.loadSession()
	if err != nil {
		return err
	}

	session.Ended = time.Now()
	return e.persistSession(session)
}

// CollectKnowledge collects knowledge artifacts
func (e *DefaultEngine) CollectKnowledge(ctx context.Context) (*KnowledgeCollection, error) {
	artifacts, err := e.GetArtifacts(ctx, PhaseKnowledge)
	if err != nil {
		return nil, err
	}

	var knowledgeArtifacts []KnowledgeArtifact
	for _, art := range artifacts {
		knowledgeArtifacts = append(knowledgeArtifacts, KnowledgeArtifact{
			Path:     art.Path,
			Verified: true,
		})
	}

	return &KnowledgeCollection{
		Artifacts: knowledgeArtifacts,
		Count:    len(knowledgeArtifacts),
	}, nil
}
