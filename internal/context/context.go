package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kdse/runtime/internal/types"
)

// HandoffContext represents the primary context handoff artifact
// per CONTEXT_HANDOFF.md specification
type HandoffContext struct {
	Project         string          `json:"project"`
	ProjectVersion  string          `json:"project_version,omitempty"`
	Schema          string          `json:"schema"`
	CurrentStage    string          `json:"current_stage"`
	PreviousStage   *string         `json:"previous_stage,omitempty"`
	StageHistory    []StageEntry    `json:"stage_history,omitempty"`
	Evidence        []string        `json:"evidence,omitempty"`
	NextAction      string          `json:"next_action"`
	AllowedContext  []string        `json:"allowed_context,omitempty"`
	RestrictedPaths []string        `json:"restricted_paths,omitempty"`
	Session         *SessionInfo    `json:"session,omitempty"`
	Artifacts       ArtifactPaths   `json:"artifacts"`
	Metadata        ContextMetadata `json:"metadata"`
}

// StageEntry records a completed stage transition
type StageEntry struct {
	Stage       string `json:"stage"`
	CompletedAt string `json:"completed_at"`
}

// SessionInfo captures current session metadata
type SessionInfo struct {
	SessionID    string  `json:"session_id"`
	StartedAt    string  `json:"started_at"`
	LastUpdated  string  `json:"last_updated"`
	SandboxID    string  `json:"sandbox_id,omitempty"`
}

// ArtifactPaths defines artifact directory locations
type ArtifactPaths struct {
	Reports     string `json:"reports"`
	Screenshots string `json:"screenshots"`
	Tests       string `json:"tests"`
	Benchmarks  string `json:"benchmarks"`
}

// ContextMetadata tracks lifecycle information
type ContextMetadata struct {
	InitializedAt    string `json:"initialized_at"`
	LastTransition   string `json:"last_transition,omitempty"`
	TransitionsCount int    `json:"transitions_count"`
}

// RuntimeContext maintains backward compatibility with existing code
type RuntimeContext struct {
	SessionID   string                  `json:"session_id"`
	Timestamp  string                  `json:"timestamp"`
	Repository *types.Repository       `json:"repository"`
	Phase      string                  `json:"phase"`
	State      string                  `json:"state"`
	Dimensions map[string]float64      `json:"dimensions"`
	Findings   []types.Finding         `json:"findings"`
	NextAction string                  `json:"next_action"`
	Guidance   string                  `json:"guidance"`
}

// Builder for RuntimeContext (backward compatible)
type Builder struct {
	repoPath   string
	repository *types.Repository
}

func NewBuilder(repoPath string) *Builder {
	return &Builder{repoPath: repoPath}
}

func (b *Builder) WithRepository(repo *types.Repository) *Builder {
	b.repository = repo
	return b
}

func (b *Builder) Build() *RuntimeContext {
	ctx := &RuntimeContext{
		SessionID:  types.GenerateSessionID(),
		Timestamp:  time.Now().Format(time.RFC3339),
		State:      string(types.StateIdle),
		Dimensions: make(map[string]float64),
		Findings:   []types.Finding{},
	}

	if b.repository != nil {
		ctx.Repository = b.repository
		ctx.Phase = b.repository.Phase
		ctx.Dimensions = b.assessDimensions()
		ctx.Findings = b.generateFindings()
		ctx.NextAction = b.determineNextAction()
		ctx.Guidance = b.generateGuidance()
	}

	return ctx
}

// NewHandoffContext creates a new HandoffContext with defaults
func NewHandoffContext(project string, currentStage string, nextAction string) *HandoffContext {
	now := time.Now().Format(time.RFC3339)
	return &HandoffContext{
		Project:        project,
		Schema:         "https://kdse.dev/schemas/context-handoff/v1",
		CurrentStage:   currentStage,
		Evidence:       []string{},
		NextAction:     nextAction,
		AllowedContext: []string{},
		RestrictedPaths: []string{},
		Artifacts: ArtifactPaths{
			Reports:     ".kdse/reports/",
			Screenshots: ".kdse/evidence/screenshots/",
			Tests:       ".kdse/evidence/tests/",
			Benchmarks:  ".kdse/evidence/benchmarks/",
		},
		Metadata: ContextMetadata{
			InitializedAt:    now,
			TransitionsCount: 0,
		},
	}
}

// TransitionStage moves to a new stage, recording history
func (h *HandoffContext) TransitionStage(newStage string, evidence ...string) {
	now := time.Now().Format(time.RFC3339)

	// Record current stage in history
	h.StageHistory = append(h.StageHistory, StageEntry{
		Stage:       h.CurrentStage,
		CompletedAt: now,
	})

	// Update stage
	h.PreviousStage = &h.CurrentStage
	h.CurrentStage = newStage

	// Add evidence
	h.Evidence = append(h.Evidence, evidence...)

	// Update metadata
	h.Metadata.LastTransition = now
	h.Metadata.TransitionsCount++
}

// SetNextAction updates the next action directive
func (h *HandoffContext) SetNextAction(action string) {
	h.NextAction = action
}

// AddEvidence appends evidence file paths
func (h *HandoffContext) AddEvidence(evidence ...string) {
	h.Evidence = append(h.Evidence, evidence...)
}

// AddAllowedContext appends paths the AI may read
func (h *HandoffContext) AddAllowedContext(paths ...string) {
	h.AllowedContext = append(h.AllowedContext, paths...)
}

// StartSession initializes session metadata
func (h *HandoffContext) StartSession(sessionID, sandboxID string) {
	h.Session = &SessionInfo{
		SessionID:   sessionID,
		StartedAt:   time.Now().Format(time.RFC3339),
		LastUpdated: time.Now().Format(time.RFC3339),
		SandboxID:   sandboxID,
	}
}

// UpdateSession updates session timestamp
func (h *HandoffContext) UpdateSession() {
	if h.Session != nil {
		h.Session.LastUpdated = time.Now().Format(time.RFC3339)
	}
}

func (b *Builder) assessDimensions() map[string]float64 {
	dimensions := map[string]float64{
		"Knowledge Artifacts":     5.0,
		"Architecture Artifacts":   5.0,
		"Implementation Artifacts": 5.0,
		"Verification Practices":   5.0,
		"Traceability":             5.0,
		"Authority Hierarchy":      5.0,
		"Governance":               5.0,
	}

	if b.repository == nil {
		return dimensions
	}

	artifacts := b.repository.Artifacts

	dimensions["Knowledge Artifacts"] = b.scoreDimension(artifacts, []string{"docs/", "README.md", "SPEC.md"}, 2.0)
	dimensions["Architecture Artifacts"] = b.scoreDimension(artifacts, []string{"architecture/", "ARCHITECTURE.md"}, 2.0)
	dimensions["Implementation Artifacts"] = b.scoreDimension(artifacts, []string{"src/"}, 1.5)
	dimensions["Verification Practices"] = b.scoreDimension(artifacts, []string{"tests/"}, 1.5)
	dimensions["Traceability"] = b.scoreDimension(artifacts, []string{"requirements/"}, 1.0)
	dimensions["Authority Hierarchy"] = 7.0
	dimensions["Governance"] = b.scoreDimension(artifacts, []string{"CONTRIBUTING.md", "LICENSE", "commands/"}, 1.0)

	return dimensions
}

func (b *Builder) scoreDimension(artifacts []string, indicators []string, weight float64) float64 {
	score := 2.0
	for _, artifact := range artifacts {
		for _, indicator := range indicators {
			if artifact == indicator {
				score += weight
			}
		}
	}
	if score > 10 {
		score = 10
	}
	return score
}

func (b *Builder) generateFindings() []types.Finding {
	var findings []types.Finding

	if b.repository == nil {
		return findings
	}

	if b.repository.Phase == string(types.PhaseConcept) {
		findings = append(findings, types.Finding{
			ID:        "F001",
			Severity:  "High",
			Dimension: "Knowledge Artifacts",
			Title:     "Repository lacks foundational documentation",
			Details:   "Create README.md and SPEC.md to establish knowledge baseline",
		})
	}

	if !b.repository.IsGitRepo {
		findings = append(findings, types.Finding{
			ID:        "F002",
			Severity:  "Medium",
			Dimension: "Governance",
			Title:     "Repository is not version controlled",
			Details:   "Initialize git repository for version control",
		})
	}

	dimensions := b.assessDimensions()
	for dim, score := range dimensions {
		if score < 4 {
			findings = append(findings, types.Finding{
				ID:        "F003",
				Severity:  "Medium",
				Dimension: dim,
				Title:     "Dimension below threshold",
				Details:   "This dimension requires attention to reach target maturity",
			})
		}
	}

	return findings
}

func (b *Builder) determineNextAction() string {
	if b.repository == nil {
		return "Initialize repository with README.md and basic structure"
	}

	phase := b.repository.Phase

	switch phase {
	case string(types.PhaseConcept):
		return "Create foundational documentation: README.md, SPEC.md"
	case string(types.PhaseDefined):
		return "Establish architecture approach and document design decisions"
	case string(types.PhaseStructured):
		return "Implement core functionality and establish testing practices"
	case string(types.PhaseUsable):
		return "Enhance verification practices and establish traceability"
	case string(types.PhaseValidated):
		return "Document governance policies and establish operational procedures"
	case string(types.PhaseProven):
		return "Continue monitoring and incremental improvements"
	default:
		return "Review and enhance engineering artifacts"
	}
}

func (b *Builder) generateGuidance() string {
	if b.repository == nil {
		return "Initialize your repository with KDSE-compliant artifacts"
	}

	return "Continue building engineering context based on detected phase: " + b.repository.Phase
}

func (c *RuntimeContext) Save(repoPath string) error {
	kdseDir := filepath.Join(repoPath, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return err
	}

	ctxPath := filepath.Join(kdseDir, "context.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ctxPath, data, 0644)
}

func Load(repoPath string) (*RuntimeContext, error) {
	ctxPath := filepath.Join(repoPath, ".kdse", "context.json")
	data, err := os.ReadFile(ctxPath)
	if err != nil {
		return nil, err
	}

	var ctx RuntimeContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

// SaveHandoff saves the HandoffContext to .kdse/context.json
func (h *HandoffContext) SaveHandoff(repoPath string) error {
	// Ensure .kdse directory exists
	kdseDir := filepath.Join(repoPath, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return err
	}

	// Ensure evidence subdirectories exist
	for _, subdir := range []string{"reports", "evidence/screenshots", "evidence/tests", "evidence/benchmarks"} {
		if err := os.MkdirAll(filepath.Join(kdseDir, subdir), 0755); err != nil {
			return err
		}
	}

	ctxPath := filepath.Join(kdseDir, "context.json")
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ctxPath, data, 0644)
}

// LoadHandoff loads the HandoffContext from .kdse/context.json
func LoadHandoff(repoPath string) (*HandoffContext, error) {
	ctxPath := filepath.Join(repoPath, ".kdse", "context.json")
	data, err := os.ReadFile(ctxPath)
	if err != nil {
		return nil, err
	}

	var ctx HandoffContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

// MustLoadHandoff loads HandoffContext or panics
func MustLoadHandoff(repoPath string) *HandoffContext {
	ctx, err := LoadHandoff(repoPath)
	if err != nil {
		panic("failed to load handoff context: " + err.Error())
	}
	return ctx
}
