package workspace

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"kdse/internal/methodology/lifecycle"
)

// Engine is the main interface for the Workspace Engine
// It owns all project state and enforces the KDSE methodology
type Engine interface {
	// Workspace lifecycle
	InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error)
	VerifyWorkspace(ctx context.Context) (*VerificationResult, error)
	LoadWorkspace(ctx context.Context) (*Workspace, error)
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

	result.RuntimeInfo = runtimeConfig
	return result, nil
}

// InitializeWorkspace initializes a new KDSE workspace
func (e *DefaultEngine) InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error) {
	kdseDir := filepath.Join(e.workspacePath, ".kdse")

	// Create staging directory for atomic initialization
	stagingDir, err := os.MkdirTemp("", "kdse-init-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(stagingDir)

	// Create .kdse structure in staging
	stagingKDSE := filepath.Join(stagingDir, ".kdse")
	if err := os.MkdirAll(stagingKDSE, 0755); err != nil {
		return nil, err
	}

	// Create required files
	if err := e.createRuntimeFiles(stagingKDSE, opts); err != nil {
		return nil, err
	}

	// Initialize phase to initialization
	if err := e.initializePhase(stagingKDSE); err != nil {
		return nil, err
	}

	// Create knowledge directory
	knowledgeDir := filepath.Join(stagingKDSE, "knowledge")
	if err := os.MkdirAll(knowledgeDir, 0755); err != nil {
		return nil, err
	}

	// Create architecture directory
	architectureDir := filepath.Join(stagingKDSE, "architecture")
	if err := os.MkdirAll(architectureDir, 0755); err != nil {
		return nil, err
	}

	// Create implementation directory
	implementationDir := filepath.Join(stagingKDSE, "implementation")
	if err := os.MkdirAll(implementationDir, 0755); err != nil {
		return nil, err
	}

	// Create verification directory
	verificationDir := filepath.Join(stagingKDSE, "verification")
	if err := os.MkdirAll(verificationDir, 0755); err != nil {
		return nil, err
	}

	// Create reports directory
	reportsDir := filepath.Join(stagingKDSE, "reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return nil, err
	}

	// Create README files in directories
	e.createDirectoryReadme(knowledgeDir)
	e.createDirectoryReadme(architectureDir)
	e.createDirectoryReadme(implementationDir)
	e.createDirectoryReadme(verificationDir)
	e.createDirectoryReadme(reportsDir)

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

// LoadWorkspace loads the current workspace state
func (e *DefaultEngine) LoadWorkspace(ctx context.Context) (*Workspace, error) {
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

	return &Workspace{
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
		validNext, _ := e.lifecycle.GetNextPhase(current)
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
			result.Missing = append(result.Missing, spec)
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
