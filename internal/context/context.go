package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kdse/runtime/internal/types"
)

type RuntimeContext struct {
	SessionID   string                 `json:"session_id"`
	Timestamp   string                 `json:"timestamp"`
	Repository  *types.Repository       `json:"repository"`
	Phase       string                 `json:"phase"`
	State       string                 `json:"state"`
	Dimensions  map[string]float64     `json:"dimensions"`
	Findings    []types.Finding        `json:"findings"`
	NextAction  string                 `json:"next_action"`
	Guidance    string                 `json:"guidance"`
}

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
			if artifact == indicator || artifact == indicator {
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
