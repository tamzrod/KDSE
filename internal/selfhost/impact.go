// Package selfhost implements KDSE self-hosting capabilities.
package selfhost

import (
	"fmt"
	"strings"
)

// ImpactAnalyzer analyzes the impact of proposed changes
type ImpactAnalyzer struct {
	model *ArchitectureModel
}

// ChangeType represents the type of change being proposed
type ChangeType string

const (
	ChangeTypeAdd       ChangeType = "add"
	ChangeTypeModify    ChangeType = "modify"
	ChangeTypeDelete   ChangeType = "delete"
	ChangeTypeRefactor ChangeType = "refactor"
)

// Change represents a proposed change to the runtime
type Change struct {
	Type        ChangeType `json:"type"`
	Component   string     `json:"component"`
	Path        string     `json:"path,omitempty"`
	Description string     `json:"description"`
	Risk        RiskLevel  `json:"risk"`
}

// RiskLevel represents the risk level of a change
type RiskLevel string

const (
	RiskLow     RiskLevel = "low"
	RiskMedium  RiskLevel = "medium"
	RiskHigh    RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// ImpactReport represents the impact analysis result
type ImpactReport struct {
	Change            *Change                   `json:"change"`
	AffectedComponents []string                 `json:"affected_components"`
	ImpactScore       float64                  `json:"impact_score"` // 0.0 - 1.0
	RiskLevel         RiskLevel                `json:"risk_level"`
	BreakingChanges   []string                 `json:"breaking_changes,omitempty"`
	AffectedDependencies [][]string             `json:"affected_dependencies"`
	Evidence         []string                 `json:"evidence"`
	Recommendations   []string                 `json:"recommendations,omitempty"`
}

// ImpactAnalysisRequest contains the request for impact analysis
type ImpactAnalysisRequest struct {
	Changes       []*Change                  `json:"changes"`
	CurrentModel  *ArchitectureModel         `json:"current_model"`
	ConsiderTests bool                       `json:"consider_tests"`
}

// ImpactAnalysisResult contains the result of analyzing multiple changes
type ImpactAnalysisResult struct {
	Requests   []*ImpactAnalysisRequest      `json:"requests,omitempty"`
	Reports    []*ImpactReport               `json:"reports"`
	TotalRisk  RiskLevel                     `json:"total_risk"`
	CanProceed bool                          `json:"can_proceed"`
	Summary    string                        `json:"summary"`
}

// NewImpactAnalyzer creates a new impact analyzer
func NewImpactAnalyzer(model *ArchitectureModel) *ImpactAnalyzer {
	return &ImpactAnalyzer{
		model: model,
	}
}

// AnalyzeChange analyzes the impact of a single change
func (a *ImpactAnalyzer) AnalyzeChange(change *Change) *ImpactReport {
	report := &ImpactReport{
		Change:    change,
		Evidence:  []string{},
	}

	// Determine affected components based on change type
	switch change.Type {
	case ChangeTypeAdd:
		report.AffectedComponents = a.analyzeAddImpact(change)
		report.ImpactScore = 0.3
		report.RiskLevel = RiskLow

	case ChangeTypeModify:
		report.AffectedComponents = a.analyzeModifyImpact(change)
		report.ImpactScore = 0.5
		report.RiskLevel = a.determineModifyRisk(change)

	case ChangeTypeDelete:
		report.AffectedComponents = a.analyzeDeleteImpact(change)
		report.ImpactScore = 0.8
		report.RiskLevel = RiskHigh

	case ChangeTypeRefactor:
		report.AffectedComponents = a.analyzeRefactorImpact(change)
		report.ImpactScore = 0.6
		report.RiskLevel = RiskMedium
	}

	// Analyze dependencies
	report.AffectedDependencies = a.analyzeDependencyImpact(report.AffectedComponents)

	// Check for breaking changes
	report.BreakingChanges = a.detectBreakingChanges(change, report.AffectedComponents)

	// Generate evidence
	report.Evidence = a.generateEvidence(change, report)

	// Generate recommendations
	report.Recommendations = a.generateRecommendations(report)

	// Update risk level based on breaking changes
	if len(report.BreakingChanges) > 0 {
		report.RiskLevel = RiskCritical
		report.ImpactScore = 1.0
	}

	return report
}

// analyzeAddImpact determines impact of adding a component
func (a *ImpactAnalyzer) analyzeAddImpact(change *Change) []string {
	var affected []string

	// New components affect their dependencies
	for name, comp := range a.model.Components {
		for _, dep := range comp.Dependencies {
			if dep == change.Component {
				affected = append(affected, name)
			}
		}
	}

	return affected
}

// analyzeModifyImpact determines impact of modifying a component
func (a *ImpactAnalyzer) analyzeModifyImpact(change *Change) []string {
	var affected []string

	// Modifying a component affects everything that depends on it
	affected = append(affected, change.Component)

	// Find all components that depend on this one
	for name, comp := range a.model.Components {
		for _, dep := range comp.Dependencies {
			if dep == change.Component {
				affected = append(affected, name)
				// Also include transitive dependents
				affected = append(affected, a.findTransitiveDependents(name)...)
			}
		}
	}

	return uniqueStrings(affected)
}

// analyzeDeleteImpact determines impact of deleting a component
func (a *ImpactAnalyzer) analyzeDeleteImpact(change *Change) []string {
	var affected []string

	// Deleting affects all dependents
	for name, comp := range a.model.Components {
		for _, dep := range comp.Dependencies {
			if dep == change.Component {
				affected = append(affected, name)
			}
		}
	}

	// Also mark the component itself
	affected = append(affected, change.Component)

	return uniqueStrings(affected)
}

// analyzeRefactorImpact determines impact of refactoring
func (a *ImpactAnalyzer) analyzeRefactorImpact(change *Change) []string {
	var affected []string

	// Refactoring affects the component and its direct dependents
	affected = append(affected, change.Component)

	for name, comp := range a.model.Components {
		for _, dep := range comp.Dependencies {
			if dep == change.Component {
				affected = append(affected, name)
			}
		}
	}

	return uniqueStrings(affected)
}

// findTransitiveDependents finds all components that transitively depend on a component
func (a *ImpactAnalyzer) findTransitiveDependents(component string) []string {
	var dependents []string
	visited := make(map[string]bool)

	var findDeps func(name string)
	findDeps = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true

		for _, comp := range a.model.Components {
			for _, dep := range comp.Dependencies {
				if dep == name {
					findDeps(comp.Name)
				}
			}
		}
	}

	findDeps(component)
	delete(visited, component) // Remove the original component

	for name := range visited {
		dependents = append(dependents, name)
	}

	return dependents
}

// determineModifyRisk determines the risk level for a modification
func (a *ImpactAnalyzer) determineModifyRisk(change *Change) RiskLevel {
	// High-risk components: runtime, orchestration, guard
	highRiskComponents := map[string]bool{
		"runtime":       true,
		"orchestration": true,
		"guard":         true,
		"discover":      true,
	}

	// Medium-risk components: knowledge, toolchain
	mediumRiskComponents := map[string]bool{
		"knowledge": true,
		"toolchain": true,
		"normalize": true,
	}

	if highRiskComponents[change.Component] {
		return RiskHigh
	}
	if mediumRiskComponents[change.Component] {
		return RiskMedium
	}

	return RiskLow
}

// analyzeDependencyImpact analyzes how dependencies are affected
func (a *ImpactAnalyzer) analyzeDependencyImpact(affected []string) [][]string {
	var impactedDeps [][]string

	for _, affectedComp := range affected {
		for _, dep := range a.model.Dependencies {
			if dep.Source == affectedComp || dep.Target == affectedComp {
				impactedDeps = append(impactedDeps, []string{dep.Source, dep.Target})
			}
		}
	}

	return impactedDeps
}

// detectBreakingChanges identifies breaking changes
func (a *ImpactAnalyzer) detectBreakingChanges(change *Change, affected []string) []string {
	var breaking []string

	switch change.Type {
	case ChangeTypeDelete:
		if change.Component == "runtime" {
			breaking = append(breaking, "Deleting runtime component breaks all other components")
		}
		if change.Component == "discover" {
			breaking = append(breaking, "Deleting discover breaks project detection")
		}

	case ChangeTypeModify:
		// Check if modification removes provided interfaces
		if comp, ok := a.model.Components[change.Component]; ok {
			if len(comp.Provides) == 0 && len(affected) > 3 {
				breaking = append(breaking, "Modification may break multiple dependent components")
			}
		}
	}

	return breaking
}

// generateEvidence generates evidence for the impact analysis
func (a *ImpactAnalyzer) generateEvidence(change *Change, report *ImpactReport) []string {
	var evidence []string

	evidence = append(evidence, fmt.Sprintf("Change type: %s", change.Type))
	evidence = append(evidence, fmt.Sprintf("Component: %s", change.Component))

	if len(report.AffectedComponents) > 0 {
		evidence = append(evidence, fmt.Sprintf("Directly affected components: %d", len(report.AffectedComponents)))
		evidence = append(evidence, fmt.Sprintf("Components: %s", strings.Join(report.AffectedComponents, ", ")))
	}

	if len(report.AffectedDependencies) > 0 {
		evidence = append(evidence, fmt.Sprintf("Affected dependencies: %d", len(report.AffectedDependencies)))
	}

	return evidence
}

// generateRecommendations generates recommendations based on impact
func (a *ImpactAnalyzer) generateRecommendations(report *ImpactReport) []string {
	var recommendations []string

	if report.RiskLevel == RiskCritical {
		recommendations = append(recommendations, "This change is critical and requires extensive review")
		recommendations = append(recommendations, "Consider a staged rollout with immediate rollback capability")
		recommendations = append(recommendations, "Ensure all tests pass before proceeding")
	}

	if report.RiskLevel == RiskHigh {
		recommendations = append(recommendations, "High-risk change - ensure backward compatibility")
		recommendations = append(recommendations, "Add deprecation warnings if applicable")
		recommendations = append(recommendations, "Document the change thoroughly")
	}

	if len(report.BreakingChanges) > 0 {
		recommendations = append(recommendations, "Breaking changes detected - coordinate with all affected teams")
	}

	if len(report.AffectedComponents) > 5 {
		recommendations = append(recommendations, "Many components affected - consider incremental changes")
	}

	return recommendations
}

// AnalyzeMultiple analyzes the impact of multiple changes
func (a *ImpactAnalyzer) AnalyzeMultiple(changes []*Change) *ImpactAnalysisResult {
	result := &ImpactAnalysisResult{
		Reports: []*ImpactReport{},
	}

	var totalImpact float64
	var maxRisk RiskLevel

	for _, change := range changes {
		report := a.AnalyzeChange(change)
		result.Reports = append(result.Reports, report)

		totalImpact += report.ImpactScore
		if report.RiskLevel == RiskCritical {
			maxRisk = RiskCritical
		} else if report.RiskLevel == RiskHigh && maxRisk != RiskCritical {
			maxRisk = RiskHigh
		} else if report.RiskLevel == RiskMedium && maxRisk == RiskLow {
			maxRisk = RiskMedium
		}
	}

	result.TotalRisk = maxRisk
	result.CanProceed = maxRisk == RiskLow || maxRisk == RiskMedium
	result.Summary = fmt.Sprintf(
		"Analyzed %d changes with %.2f total impact score and %s risk level",
		len(changes), totalImpact, maxRisk,
	)

	return result
}

// FormatImpactReport formats an impact report for display
func FormatImpactReport(report *ImpactReport) string {
	var output string

	output += "═══════════════════════════════════════════════════════════════\n"
	output += "                    Impact Analysis Report                     \n"
	output += "═══════════════════════════════════════════════════════════════\n"
	output += fmt.Sprintf("  Change: %s %s\n", report.Change.Type, report.Change.Component)
	output += fmt.Sprintf("  Risk Level: %s\n", report.Change.Risk)
	output += fmt.Sprintf("  Impact Score: %.2f\n", report.ImpactScore)
	output += "───────────────────────────────────────────────────────────────\n"

	if len(report.AffectedComponents) > 0 {
		output += "  Affected Components:\n"
		for _, comp := range report.AffectedComponents {
			output += fmt.Sprintf("    • %s\n", comp)
		}
	}

	if len(report.BreakingChanges) > 0 {
		output += "  Breaking Changes:\n"
		for _, bc := range report.BreakingChanges {
			output += fmt.Sprintf("    ⚠ %s\n", bc)
		}
	}

	if len(report.Evidence) > 0 {
		output += "  Evidence:\n"
		for _, ev := range report.Evidence {
			output += fmt.Sprintf("    %s\n", ev)
		}
	}

	if len(report.Recommendations) > 0 {
		output += "  Recommendations:\n"
		for _, rec := range report.Recommendations {
			output += fmt.Sprintf("    → %s\n", rec)
		}
	}

	output += "═══════════════════════════════════════════════════════════════\n"

	return output
}

// uniqueStrings removes duplicates from a string slice
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
