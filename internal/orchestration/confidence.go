package orchestration

import (
	"os"
	"path/filepath"
	"strings"
)

// ConfidenceEvaluator evaluates confidence levels for orchestration decisions
type ConfidenceEvaluator struct {
	config *EngineConfig
}

// NewConfidenceEvaluator creates a new confidence evaluator
func NewConfidenceEvaluator(config *EngineConfig) *ConfidenceEvaluator {
	if config == nil {
		config = DefaultEngineConfig()
	}
	return &ConfidenceEvaluator{config: config}
}

// EvaluateConfidence performs a complete confidence evaluation
func (e *ConfidenceEvaluator) EvaluateConfidence(workspace *WorkspaceInfo) (*ConfidenceLevel, error) {
	// Evaluate foundation confidence first (most important)
	foundationConf, err := e.evaluateFoundation(workspace)
	if err != nil {
		foundationConf = 0.0
	}

	// Evaluate repository confidence
	repoConf, err := e.evaluateRepository(workspace)
	if err != nil {
		repoConf = 0.0
	}

	// Evaluate evidence confidence
	evidenceConf, err := e.evaluateEvidenceConfidence(workspace)
	if err != nil {
		evidenceConf = 0.0
	}

	// Calculate overall confidence (weighted average)
	overall := (foundationConf * 0.5) + (repoConf * 0.3) + (evidenceConf * 0.2)

	meetsThreshold := foundationConf >= e.config.FoundationThreshold

	dimensions := map[string]float64{
		"foundation":  foundationConf,
		"repository":  repoConf,
		"evidence":    evidenceConf,
	}

	assessments := e.gatherAssessments(workspace, foundationConf, repoConf, evidenceConf)

	return &ConfidenceLevel{
		Overall:        overall,
		Foundation:     foundationConf,
		Repository:     repoConf,
		Evidence:       evidenceConf,
		Threshold:      e.config.FoundationThreshold,
		MeetsThreshold: meetsThreshold,
		Dimensions:     dimensions,
		Assessments:    assessments,
	}, nil
}

// CanImplement returns true if implementation is allowed (Foundation threshold met)
func (e *ConfidenceEvaluator) CanImplement(confidence *ConfidenceLevel) bool {
	return confidence != nil && confidence.MeetsThreshold
}

// evaluateFoundation evaluates the Foundation document confidence
func (e *ConfidenceEvaluator) evaluateFoundation(workspace *WorkspaceInfo) (float64, error) {
	kdsePath := workspace.KDSEPath

	// Required Foundation documents
	requiredDocs := map[string]float64{
		"foundation/README.md":           1.0,
		"foundation/001-vision.md":       0.8,
		"foundation/002-principles.md":    0.8,
		"foundation/003-standards.md":    0.9,
		"foundation/004-engineering-model.md": 1.0,
		"foundation/005-architecture.md": 0.7,
		"foundation/006-chain-of-authority.md": 0.8,
	}

	// Check for audit directory
	auditDir := filepath.Join(kdsePath, "..", "docs", "audit")
	hasAudit := false
	if info, err := os.Stat(auditDir); err == nil && info.IsDir() {
		hasAudit = true
	}

	// Check for knowledge directory
	knowledgeDir := filepath.Join(kdsePath, "knowledge")
	hasKnowledge := false
	if info, err := os.Stat(knowledgeDir); err == nil && info.IsDir() {
		hasKnowledge = true
	}

	score := 0.0
	maxScore := 0.0

	for doc, weight := range requiredDocs {
		maxScore += weight
		docPath := filepath.Join(kdsePath, doc)
		if _, err := os.Stat(docPath); err == nil {
			score += weight
		}
	}

	// Add bonus for audit directory
	if hasAudit {
		score += 0.5
		maxScore += 0.5
	}

	// Add bonus for knowledge directory
	if hasKnowledge {
		score += 0.3
		maxScore += 0.3
	}

	if maxScore == 0 {
		return 0.0, nil
	}

	return score / maxScore, nil
}

// evaluateRepository evaluates repository readiness confidence
func (e *ConfidenceEvaluator) evaluateRepository(workspace *WorkspaceInfo) (float64, error) {
	repoPath := workspace.RepositoryPath
	if repoPath == "" {
		repoPath = workspace.ResolvedPath
	}

	score := 0.0

	// Check for README
	if _, err := os.Stat(filepath.Join(repoPath, "README.md")); err == nil {
		score += 1.0
	}

	// Check for git repository
	if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
		score += 1.0
	}

	// Check for source directory
	if _, err := os.Stat(filepath.Join(repoPath, "src")); err == nil {
		score += 0.5
	}

	// Check for docs directory
	if _, err := os.Stat(filepath.Join(repoPath, "docs")); err == nil {
		score += 0.5
	}

	return score / 4.0, nil
}

// evaluateEvidenceConfidence evaluates evidence completeness
func (e *ConfidenceEvaluator) evaluateEvidenceConfidence(workspace *WorkspaceInfo) (float64, error) {
	evidencePath := filepath.Join(workspace.KDSEPath, "evidence")
	
	if _, err := os.Stat(evidencePath); os.IsNotExist(err) {
		return 0.0, nil
	}

	// Count evidence files
	var count int
	filepath.Walk(evidencePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && !strings.HasPrefix(filepath.Base(path), ".") {
			count++
		}
		return nil
	})

	// Normalize to 0-1 scale (5+ files = full score)
	if count >= 5 {
		return 1.0, nil
	}
	return float64(count) / 5.0, nil
}

// gatherAssessments collects detailed confidence assessments
func (e *ConfidenceEvaluator) gatherAssessments(workspace *WorkspaceInfo, foundation, repository, evidence float64) []ConfidenceAssessment {
	return []ConfidenceAssessment{
		{
			Type:      "foundation_documents",
			Score:     foundation,
			Weight:    0.5,
			Details:   "Foundation document completeness and accessibility",
			Satisfied: foundation >= e.config.FoundationThreshold,
		},
		{
			Type:      "repository_structure",
			Score:     repository,
			Weight:    0.3,
			Details:   "Repository has basic structure and documentation",
			Satisfied: repository >= 0.5,
		},
		{
			Type:      "evidence_gathered",
			Score:     evidence,
			Weight:    0.2,
			Details:   "Evidence collected for current phase",
			Satisfied: evidence >= e.config.EvidenceThreshold,
		},
	}
}

// IsFoundationReady checks if Foundation is ready for implementation
func (e *ConfidenceEvaluator) IsFoundationReady(workspace *WorkspaceInfo) (bool, float64, error) {
	foundationConf, err := e.evaluateFoundation(workspace)
	if err != nil {
		return false, 0.0, err
	}

	return foundationConf >= e.config.FoundationThreshold, foundationConf, nil
}

// GetRequiredFoundationDocs returns the list of required Foundation documents
func (e *ConfidenceEvaluator) GetRequiredFoundationDocs() []string {
	return []string{
		"foundation/README.md",
		"foundation/001-vision.md",
		"foundation/002-principles.md",
		"foundation/003-standards.md",
		"foundation/004-engineering-model.md",
		"foundation/005-architecture.md",
		"foundation/006-chain-of-authority.md",
	}
}

// AssessFoundationCompleteness returns detailed Foundation assessment
func (e *ConfidenceEvaluator) AssessFoundationCompleteness(workspace *WorkspaceInfo) (*FoundationAssessment, error) {
	kdsePath := workspace.KDSEPath
	requiredDocs := e.GetRequiredFoundationDocs()

	assessment := &FoundationAssessment{
		Required:  len(requiredDocs),
		Present:   0,
		Missing:   []string{},
		Documents: []DocumentStatus{},
	}

	for _, doc := range requiredDocs {
		docPath := filepath.Join(kdsePath, doc)
		exists, _ := exists(docPath)

		status := DocumentStatus{
			Path:     doc,
			Required: true,
			Present:  exists,
		}

		if exists {
			assessment.Present++
			status.Content, _ = e.assessDocumentContent(docPath)
		} else {
			assessment.Missing = append(assessment.Missing, doc)
		}

		assessment.Documents = append(assessment.Documents, status)
	}

	assessment.Completeness = float64(assessment.Present) / float64(assessment.Required)
	assessment.Ready = assessment.Completeness >= e.config.FoundationThreshold

	return assessment, nil
}

// FoundationAssessment contains detailed Foundation evaluation
type FoundationAssessment struct {
	Required     int               `json:"required"`
	Present      int               `json:"present"`
	Missing      []string          `json:"missing"`
	Completeness float64           `json:"completeness"`
	Ready        bool              `json:"ready"`
	Documents    []DocumentStatus  `json:"documents"`
}

// DocumentStatus represents the status of a single document
type DocumentStatus struct {
	Path          string `json:"path"`
	Required      bool   `json:"required"`
	Present       bool   `json:"present"`
	ContentScore  float64 `json:"content_score,omitempty"`
	Content       string `json:"content,omitempty"`
}

// assessDocumentContent does basic content quality assessment
func (e *ConfidenceEvaluator) assessDocumentContent(path string) (string, float64) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0.0
	}

	content := string(data)
	
	// Basic quality checks
	score := 0.0
	if len(content) > 100 {
		score += 0.3
	}
	if len(content) > 500 {
		score += 0.2
	}

	return content, score
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
