// Package enforcer provides strict KDSE principle enforcement.
// It blocks premature coding and ensures knowledge base + foundation are built first.
// This is the core runtime guard that makes KDSE principles automatic, not just suggestions.
package enforcer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EnforcementLevel represents how strictly to enforce KDSE principles
type EnforcementLevel string

const (
	EnforcementOff     EnforcementLevel = "off"      // No enforcement
	EnforcementWarning EnforcementLevel = "warning"  // Warn but allow
	EnforcementStrict  EnforcementLevel = "strict"   // Block violations
	EnforcementHard    EnforcementLevel = "hard"    // Block even with --force flag
)

// EnforcementError represents a KDSE principle violation
type EnforcementError struct {
	Code             string   `json:"code"`
	Message          string   `json:"message"`
	Hint             string   `json:"hint"`
	RequiredPhase    string   `json:"required_phase,omitempty"`
	CurrentPhase     string   `json:"current_phase,omitempty"`
	MissingArtifacts []string `json:"missing_artifacts,omitempty"`
	Violations       []string `json:"violations,omitempty"`
	Severity         string   `json:"severity"`
	CorrectiveAction string   `json:"corrective_action,omitempty"`
	Blocked          bool     `json:"blocked"`
}

func (e *EnforcementError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Error codes for KDSE enforcement
const (
	CodeNoFoundation            = "KDSE_ENF_001"
	CodeNoKnowledge             = "KDSE_ENF_002"
	CodeNoArchitecture          = "KDSE_ENF_003"
	CodePrematureImplementation = "KDSE_ENF_004"
	CodeRepositoryNotAnalyzed   = "KDSE_ENF_005"
	CodeFoundationIncomplete    = "KDSE_ENF_006"
	CodeKnowledgeIncomplete     = "KDSE_ENF_007"
	CodePhaseViolation          = "KDSE_ENF_008"
	CodeComplianceViolation     = "KDSE_ENF_009"
)

// Severity constants
const (
	SeverityBlocker  = "BLOCKER"
	SeverityCritical = "CRITICAL"
	SeverityHigh     = "HIGH"
	SeverityMedium   = "MEDIUM"
	SeverityLow      = "LOW"
)

// Required foundation documents
var RequiredFoundationDocs = []string{
	"PROBLEM.md",
	"SPEC.md",
	"REQUIREMENTS.md",
	"ASSUMPTIONS.md",
	"CONSTRAINTS.md",
}

// Required knowledge categories
var RequiredKnowledgeCategories = []string{
	"general",
	"operational",
	"developmental",
}

// Engine is the strict KDSE principle enforcement engine
type Engine struct {
	repoPath         string
	kdsePath         string
	enforcementLevel EnforcementLevel
	autoCreate       bool  // Auto-create missing foundation/knowledge
	autoTransition   bool  // Auto-transition to required phase
	logging          bool
	violations       []EnforcementError
}

// NewEngine creates a new enforcement engine
func NewEngine(repoPath string) *Engine {
	return &Engine{
		repoPath:         repoPath,
		kdsePath:         filepath.Join(repoPath, ".kdse"),
		enforcementLevel:  EnforcementStrict,
		autoCreate:       true,
		autoTransition:   true,
		logging:          true,
		violations:       []EnforcementError{},
	}
}

// NewEngineWithLevel creates an engine with specific enforcement level
func NewEngineWithLevel(repoPath string, level EnforcementLevel) *Engine {
	engine := NewEngine(repoPath)
	engine.enforcementLevel = level
	return engine
}

// SetEnforcementLevel changes the enforcement level
func (e *Engine) SetEnforcementLevel(level EnforcementLevel) {
	e.enforcementLevel = level
}

// SetAutoCreate enables/disables auto-creation of missing artifacts
func (e *Engine) SetAutoCreate(enabled bool) {
	e.autoCreate = enabled
}

// SetAutoTransition enables/disables auto-phase transitions
func (e *Engine) SetAutoTransition(enabled bool) {
	e.autoTransition = enabled
}

// EnforceImplementation checks if implementation is allowed
// Returns error if foundation/knowledge/architecture is missing
func (e *Engine) EnforceImplementation() error {
	e.violations = []EnforcementError{}

	// Step 1: Check workspace exists
	if err := e.checkWorkspaceExists(); err != nil {
		return err
	}

	// Step 2: Check foundation exists and is complete
	if err := e.checkFoundation(); err != nil {
		return err
	}

	// Step 3: Check knowledge base exists
	if err := e.checkKnowledge(); err != nil {
		return err
	}

	// Step 4: Check architecture exists
	if err := e.checkArchitecture(); err != nil {
		return err
	}

	// Step 5: Check repository analysis was done
	if err := e.checkRepositoryAnalysis(); err != nil {
		return err
	}

	// If we have violations and are in warning mode, return warning
	if len(e.violations) > 0 && e.enforcementLevel == EnforcementWarning {
		return e.violations[0]
	}

	return nil
}

// EnforceFoundation checks if foundation phase is complete
func (e *Engine) EnforceFoundation() error {
	e.violations = []EnforcementError{}

	if err := e.checkWorkspaceExists(); err != nil {
		return err
	}

	return e.checkFoundation()
}

// EnforceKnowledge checks if knowledge phase has content
func (e *Engine) EnforceKnowledge() error {
	e.violations = []EnforcementError{}

	if err := e.checkWorkspaceExists(); err != nil {
		return err
	}

	return e.checkKnowledge()
}

// EnforceArchitecture checks if architecture document exists
func (e *Engine) EnforceArchitecture() error {
	e.violations = []EnforcementError{}

	if err := e.checkWorkspaceExists(); err != nil {
		return err
	}

	return e.checkArchitecture()
}

// EnforceForOperation checks if a specific operation is allowed
func (e *Engine) EnforceForOperation(operation string) error {
	switch operation {
	case "implement", "build", "code", "create-files", "write-code":
		return e.EnforceImplementation()
	case "foundation", "spec", "requirements":
		return e.EnforceFoundation()
	case "knowledge", "research", "learn":
		return e.EnforceKnowledge()
	case "architecture", "design", "architecture-doc":
		return e.EnforceArchitecture()
	default:
		return nil
	}
}

// checkWorkspaceExists verifies .kdse/ directory exists
func (e *Engine) checkWorkspaceExists() error {
	if _, err := os.Stat(e.kdsePath); os.IsNotExist(err) {
		err := &EnforcementError{
			Code:             CodeNoFoundation,
			Message:          "KDSE workspace (.kdse/) does not exist",
			Hint:             "Run 'kdse initialize' to create the workspace",
			Severity:         SeverityBlocker,
			Blocked:          true,
			CorrectiveAction: "Execute: kdse initialize",
		}
		e.violations = append(e.violations, *err)
		return err
	}
	return nil
}

// checkFoundation verifies all foundation documents exist and are populated
func (e *Engine) checkFoundation() error {
	foundationPath := filepath.Join(e.kdsePath, "foundation")

	// Check directory exists
	if _, err := os.Stat(foundationPath); os.IsNotExist(err) {
		err := &EnforcementError{
			Code:             CodeNoFoundation,
			Message:          "Foundation directory does not exist",
			Hint:             "Foundation must be created before implementation",
			Severity:         SeverityBlocker,
			Blocked:          true,
			RequiredPhase:    "Foundation",
			CorrectiveAction: "Execute: kdse foundation create",
		}
		e.violations = append(e.violations, *err)
		return err
	}

	// Check required documents
	var missing []string
	var empty []string

	for _, doc := range RequiredFoundationDocs {
		docPath := filepath.Join(foundationPath, doc)
		info, err := os.Stat(docPath)
		if os.IsNotExist(err) {
			missing = append(missing, doc)
			continue
		}
		// Check if file is empty or just template
		if info.Size() < 100 {
			empty = append(empty, doc)
		}
	}

	if len(missing) > 0 {
		err := &EnforcementError{
			Code:             CodeFoundationIncomplete,
			Message:          fmt.Sprintf("Foundation incomplete: missing documents: %v", missing),
			Hint:             "All foundation documents must exist before implementation",
			Severity:         SeverityBlocker,
			Blocked:          true,
			RequiredPhase:    "Foundation",
			MissingArtifacts: missing,
			CorrectiveAction: fmt.Sprintf("Create: %v", missing),
		}
		e.violations = append(e.violations, *err)

		// Auto-create if enabled
		if e.autoCreate {
			e.createFoundationDocuments(missing)
		}

		return err
	}

	if len(empty) > 0 {
		err := &EnforcementError{
			Code:             CodeFoundationIncomplete,
			Message:          fmt.Sprintf("Foundation incomplete: empty documents: %v", empty),
			Hint:             "Foundation documents must be populated with project-specific content",
			Severity:         SeverityHigh,
			Blocked:          true,
			RequiredPhase:    "Foundation",
			MissingArtifacts: empty,
			CorrectiveAction: "Populate: " + strings.Join(empty, ", "),
		}
		e.violations = append(e.violations, *err)
		return err
	}

	return nil
}

// checkKnowledge verifies knowledge base has content
func (e *Engine) checkKnowledge() error {
	knowledgePath := filepath.Join(e.kdsePath, "knowledge")

	// Check directory exists
	if _, err := os.Stat(knowledgePath); os.IsNotExist(err) {
		err := &EnforcementError{
			Code:             CodeNoKnowledge,
			Message:          "Knowledge directory does not exist",
			Hint:             "Knowledge base must be built before implementation",
			Severity:         SeverityBlocker,
			Blocked:          true,
			RequiredPhase:    "Knowledge Collection",
			CorrectiveAction: "Execute: kdse knowledge collect",
		}
		e.violations = append(e.violations, *err)
		return err
	}

	// Check required categories exist
	var missing []string
	for _, cat := range RequiredKnowledgeCategories {
		catPath := filepath.Join(knowledgePath, cat)
		if _, err := os.Stat(catPath); os.IsNotExist(err) {
			missing = append(missing, cat)
		}
	}

	if len(missing) > 0 {
		err := &EnforcementError{
			Code:             CodeKnowledgeIncomplete,
			Message:          fmt.Sprintf("Knowledge incomplete: missing categories: %v", missing),
			Hint:             "All knowledge categories must have content",
			Severity:         SeverityHigh,
			Blocked:          true,
			RequiredPhase:    "Knowledge Collection",
			MissingArtifacts: missing,
			CorrectiveAction: "Execute: kdse knowledge collect --all",
		}
		e.violations = append(e.violations, *err)

		// Auto-create if enabled
		if e.autoCreate {
			e.createKnowledgeCategories(missing)
		}

		return err
	}

	// Check if any content exists
	hasContent := false
	for _, cat := range RequiredKnowledgeCategories {
		catPath := filepath.Join(knowledgePath, cat)
		if e.dirHasContent(catPath) {
			hasContent = true
			break
		}
	}

	if !hasContent {
		err := &EnforcementError{
			Code:             CodeNoKnowledge,
			Message:          "Knowledge base is empty",
			Hint:             "Collect domain knowledge before implementation",
			Severity:         SeverityBlocker,
			Blocked:          true,
			RequiredPhase:    "Knowledge Collection",
			CorrectiveAction: "Execute: kdse knowledge collect",
		}
		e.violations = append(e.violations, *err)
		return err
	}

	return nil
}

// checkArchitecture verifies architecture document exists
func (e *Engine) checkArchitecture() error {
	archPath := filepath.Join(e.kdsePath, "foundation", "ARCHITECTURE.md")

	if _, err := os.Stat(archPath); os.IsNotExist(err) {
		err := &EnforcementError{
			Code:             CodeNoArchitecture,
			Message:          "Architecture document does not exist",
			Hint:             "Architecture must be defined before implementation",
			Severity:         SeverityBlocker,
			Blocked:          true,
			RequiredPhase:    "Architecture",
			CorrectiveAction: "Execute: kdse architecture design",
		}
		e.violations = append(e.violations, *err)
		return err
	}

	// Check if populated
	info, err := os.Stat(archPath)
	if err != nil || info.Size() < 200 {
		err := &EnforcementError{
			Code:             CodeNoArchitecture,
			Message:          "Architecture document is empty or incomplete",
			Hint:             "Define system architecture before implementation",
			Severity:         SeverityHigh,
			Blocked:          true,
			RequiredPhase:    "Architecture",
			CorrectiveAction: "Populate: .kdse/foundation/ARCHITECTURE.md",
		}
		e.violations = append(e.violations, *err)
		return err
	}

	return nil
}

// checkRepositoryAnalysis verifies repository has been analyzed
func (e *Engine) checkRepositoryAnalysis() error {
	analysisPath := filepath.Join(e.kdsePath, "knowledge", "repository-analysis.md")

	if _, err := os.Stat(analysisPath); os.IsNotExist(err) {
		err := &EnforcementError{
			Code:             CodeRepositoryNotAnalyzed,
			Message:          "Repository has not been analyzed",
			Hint:             "Analyze repository structure before implementation",
			Severity:         SeverityMedium,
			Blocked:          e.enforcementLevel == EnforcementHard,
			RequiredPhase:    "Problem",
			CorrectiveAction: "Execute: kdse analyze --repository",
		}
		e.violations = append(e.violations, *err)

		if e.enforcementLevel == EnforcementHard {
			return err
		}
	}

	return nil
}

// createFoundationDocuments creates missing foundation documents
func (e *Engine) createFoundationDocuments(missing []string) {
	foundationPath := filepath.Join(e.kdsePath, "foundation")
	os.MkdirAll(foundationPath, 0755)

	templates := map[string]string{
		"PROBLEM.md": "# Problem Statement\n\n## Problem Description\n\n[Describe the problem to be solved]\n\n## Impact\n\n[Describe the impact of not solving this problem]\n\n## Success Criteria\n\n- [ ] Criterion 1\n- [ ] Criterion 2\n",
		"SPEC.md": "# Project Specification\n\n## Overview\n\n[High-level description of the project]\n\n## Scope\n\n### In Scope\n\n- Item 1\n- Item 2\n\n### Out of Scope\n\n- Item 1\n\n## Deliverables\n\n- Deliverable 1\n",
		"REQUIREMENTS.md": "# Functional Requirements\n\n## FR-001: [Title]\n\n**Description:** [What the requirement is]\n\n**Acceptance Criteria:**\n- [ ] Criterion 1\n- [ ] Criterion 2\n\n## FR-002: [Title]\n\n**Description:** [What the requirement is]\n\n**Acceptance Criteria:**\n- [ ] Criterion 1\n",
		"ASSUMPTIONS.md": "# Key Assumptions\n\n## Technical Assumptions\n\n1. [Assumption 1]\n2. [Assumption 2]\n\n## Business Assumptions\n\n1. [Assumption 1]\n2. [Assumption 2]\n\n## Constraints\n\n1. [Constraint 1]\n2. [Constraint 2]\n",
		"CONSTRAINTS.md": "# Project Constraints\n\n## Technical Constraints\n\n- Constraint 1\n- Constraint 2\n\n## Schedule Constraints\n\n- Constraint 1\n\n## Resource Constraints\n\n- Constraint 1\n",
	}

	for _, doc := range missing {
		content, ok := templates[doc]
		if !ok {
			content = fmt.Sprintf("# %s\n\n[Content to be filled]\n", strings.TrimSuffix(doc, ".md"))
		}

		docPath := filepath.Join(foundationPath, doc)
		if _, err := os.Stat(docPath); os.IsNotExist(err) {
			os.WriteFile(docPath, []byte(content), 0644)
		}
	}
}

// createKnowledgeCategories creates missing knowledge categories
func (e *Engine) createKnowledgeCategories(missing []string) {
	knowledgePath := filepath.Join(e.kdsePath, "knowledge")
	os.MkdirAll(knowledgePath, 0755)

	for _, cat := range missing {
		catPath := filepath.Join(knowledgePath, cat)
		os.MkdirAll(catPath, 0755)

		// Create README in each category
		readme := fmt.Sprintf("# %s Knowledge\n\n[Collect knowledge artifacts for %s]\n",
			strings.Title(cat), cat)
		os.WriteFile(filepath.Join(catPath, "README.md"), []byte(readme), 0644)
	}
}

// dirHasContent checks if a directory has any files
func (e *Engine) dirHasContent(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if e.dirHasContent(filepath.Join(dirPath, entry.Name())) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

// GetViolations returns all recorded violations
func (e *Engine) GetViolations() []EnforcementError {
	return e.violations
}

// HasViolations returns true if any violations were recorded
func (e *Engine) HasViolations() bool {
	return len(e.violations) > 0
}

// GenerateEnforcementReport creates a detailed enforcement report
func (e *Engine) GenerateEnforcementReport() *EnforcementReport {
	report := &EnforcementReport{
		Timestamp:       time.Now().Format(time.RFC3339),
		EnforcementLevel: e.enforcementLevel,
		Violations:      []EnforcementError{},
		Warnings:        []string{},
		Blocked:         false,
		CanProceed:      true,
	}

	for _, v := range e.violations {
		if v.Blocked && (e.enforcementLevel == EnforcementStrict || e.enforcementLevel == EnforcementHard) {
			report.Blocked = true
			report.CanProceed = false
		}
		report.Violations = append(report.Violations, v)
	}

	return report
}

// FormatEnforcementReport formats the report as a string
func (e *Engine) FormatEnforcementReport() string {
	report := e.GenerateEnforcementReport()

	var sb strings.Builder
	sb.WriteString("\n╔════════════════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║              KDSE ENFORCEMENT REPORT                                ║\n")
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ Enforcement Level: %s\n", e.enforcementLevel))
	sb.WriteString(fmt.Sprintf("║ Timestamp: %s\n", report.Timestamp))
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")

	if len(report.Violations) == 0 {
		sb.WriteString("║ ✓ All KDSE principles satisfied\n")
	} else {
		for _, v := range report.Violations {
			icon := "✗"
			if !v.Blocked {
				icon = "⚠"
			}
			sb.WriteString(fmt.Sprintf("║ %s [%s] %s\n", icon, v.Code, v.Severity))
			sb.WriteString(fmt.Sprintf("║   %s\n", v.Message))
			if v.RequiredPhase != "" {
				sb.WriteString(fmt.Sprintf("║   Required Phase: %s\n", v.RequiredPhase))
			}
			if v.CorrectiveAction != "" {
				sb.WriteString(fmt.Sprintf("║   Fix: %s\n", v.CorrectiveAction))
			}
		}
	}

	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	if report.Blocked {
		sb.WriteString("║ ██ BLOCKED: Premature implementation detected                      ║\n")
		sb.WriteString("║ ██ Foundation and Knowledge phase MUST be completed first         ║\n")
	} else if report.CanProceed {
		sb.WriteString("║ ✓ Can proceed with implementation\n")
	} else {
		sb.WriteString("║ ⚠ Warnings issued - proceed with caution\n")
	}
	sb.WriteString("╚════════════════════════════════════════════════════════════════════════╝\n")

	return sb.String()
}

// ToJSON returns the enforcement report as JSON
func (e *Engine) ToJSON() string {
	report := e.GenerateEnforcementReport()
	data, _ := json.MarshalIndent(report, "", "  ")
	return string(data)
}

// EnforcementReport is the structured report output
type EnforcementReport struct {
	Timestamp        string             `json:"timestamp"`
	EnforcementLevel EnforcementLevel   `json:"enforcement_level"`
	Violations       []EnforcementError `json:"violations"`
	Warnings         []string           `json:"warnings"`
	Blocked          bool               `json:"blocked"`
	CanProceed       bool               `json:"can_proceed"`
}

// ValidateOperation checks if an operation is allowed in the current state
func (e *Engine) ValidateOperation(operation string) *ValidationResult {
	result := &ValidationResult{
		Operation:   operation,
		Allowed:    true,
		Violations: []EnforcementError{},
		AutoActions: []string{},
	}

	// Run appropriate check
	var err error
	switch operation {
	case "implement", "build", "code", "create-files", "write-code":
		err = e.EnforceImplementation()
	case "foundation", "spec", "requirements":
		err = e.EnforceFoundation()
	case "knowledge", "research", "learn":
		err = e.EnforceKnowledge()
	case "architecture", "design", "architecture-doc":
		err = e.EnforceArchitecture()
	}

	if err != nil {
		result.Allowed = false
		result.Violations = e.violations

		// Suggest auto-actions
		if e.autoCreate {
			for _, v := range e.violations {
				if v.CorrectiveAction != "" {
					result.AutoActions = append(result.AutoActions, v.CorrectiveAction)
				}
			}
		}
	}

	return result
}

// ValidationResult contains the result of an operation validation
type ValidationResult struct {
	Operation   string             `json:"operation"`
	Allowed     bool               `json:"allowed"`
	Violations  []EnforcementError `json:"violations"`
	AutoActions []string           `json:"auto_actions,omitempty"`
}
