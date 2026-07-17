// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// OwnershipGuard validates ownership boundaries between project and KDSE runtime.
// It ensures:
//   - Project artifacts are in the project layer (not in .kdse/)
//   - Runtime artifacts are in the runtime layer (.kdse/)
//   - No architectural drift occurs
type OwnershipGuard struct {
	repoPath string
}

// NewOwnershipGuard creates a new Ownership Guard for the given repository path
func NewOwnershipGuard(repoPath string) *OwnershipGuard {
	return &OwnershipGuard{
		repoPath: repoPath,
	}
}

// OwnershipViolation represents a violation of ownership boundaries
type OwnershipViolation struct {
	Artifact       string
	ActualPath     string
	ExpectedDomain string
	SuggestedPath  string
	Severity       string // "error" or "warning"
	Message        string
}

// OwnershipGuardResult contains the results of ownership validation
type OwnershipGuardResult struct {
	RuntimeGuardResult
	Violations []OwnershipViolation
}

// Validate checks for ownership boundary violations
func (g *OwnershipGuard) Validate(ctx context.Context) *OwnershipGuardResult {
	result := &OwnershipGuardResult{
		RuntimeGuardResult: RuntimeGuardResult{
			Valid:     true,
			State:     StateProject,
			GuardType: GuardTypeProject,
		},
		Violations: []OwnershipViolation{},
	}

	// Check for project artifacts in .kdse/ (critical violation)
	g.checkProjectArtifactsInRuntime(result)

	// Check for runtime artifacts outside .kdse/ (warning)
	g.checkRuntimeArtifactsOutsideRuntime(result)

	// Check for reference artifacts (should be in .kdse/references/)
	g.checkReferencePlacement(result)

	// Check for knowledge artifacts (should be in .kdse/knowledge/)
	g.checkKnowledgePlacement(result)

	if len(result.Violations) > 0 {
		result.Valid = false
		result.Error = NewRuntimeGuardError(
			GuardTypeProject,
			"OWNERSHIP_VIOLATION",
			"Ownership boundaries violated",
			"Review ownership violations and migrate artifacts to correct locations",
			StateProject,
		)
	}

	return result
}

// checkProjectArtifactsInRuntime checks if project artifacts are incorrectly in .kdse/
func (g *OwnershipGuard) checkProjectArtifactsInRuntime(result *OwnershipGuardResult) {
	kdseDir := filepath.Join(g.repoPath, ".kdse")

	// Project artifacts that should NOT be in .kdse/
	projectArtifacts := []string{
		"README.md",
		"LICENSE",
		"CHANGELOG.md",
		"src/",
		"tests/",
		"cmd/",
		"internal/",
		"docs/",
		"deploy/",
	}

	// Check if any project artifacts exist in .kdse/
	for _, artifact := range projectArtifacts {
		path := filepath.Join(kdseDir, artifact)
		if exists, _ := g.pathExists(path); exists {
			isDir := strings.HasSuffix(artifact, "/")
			artifactName := strings.TrimSuffix(artifact, "/")

			var msg string
			if isDir {
				msg = "Project directory incorrectly placed in .kdse/"
			} else {
				msg = "Project file incorrectly placed in .kdse/"
			}

			violation := OwnershipViolation{
				Artifact:        artifactName,
				ActualPath:      filepath.Join(".kdse", artifact),
				ExpectedDomain:  "project",
				SuggestedPath:   artifact,
				Severity:        "error",
				Message:         msg,
			}

			result.Violations = append(result.Violations, violation)
		}
	}
}

// checkRuntimeArtifactsOutsideRuntime checks if runtime artifacts are outside .kdse/
func (g *OwnershipGuard) checkRuntimeArtifactsOutsideRuntime(result *OwnershipGuardResult) {
	// Runtime artifacts that SHOULD be in .kdse/
	runtimeArtifacts := []string{
		"laboratory",
		"runtime",
	}

	for _, artifact := range runtimeArtifacts {
		path := filepath.Join(g.repoPath, artifact)
		if exists, isDir := g.pathExists(path); exists && isDir {
			violation := OwnershipViolation{
				Artifact:       artifact,
				ActualPath:     artifact,
				ExpectedDomain: "runtime",
				SuggestedPath: filepath.Join(".kdse", artifact),
				Severity:       "warning",
				Message:        "Runtime artifact should be in .kdse/",
			}
			result.Violations = append(result.Violations, violation)
		}
	}
}

// checkReferencePlacement validates reference artifacts are in .kdse/references/
func (g *OwnershipGuard) checkReferencePlacement(result *OwnershipGuardResult) {
	entries, err := os.ReadDir(g.repoPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			// Check if this looks like a reference directory outside .kdse/references/
			if g.looksLikeReference(entry.Name()) {
				violation := OwnershipViolation{
					Artifact:       entry.Name(),
					ActualPath:     entry.Name(),
					ExpectedDomain: "reference",
					SuggestedPath: filepath.Join(".kdse", "references", entry.Name()),
					Severity:       "warning",
					Message:        "Reference directory should be in .kdse/references/",
				}
				result.Violations = append(result.Violations, violation)
			}
		}
	}
}

// checkKnowledgePlacement validates knowledge artifacts are in .kdse/knowledge/
func (g *OwnershipGuard) checkKnowledgePlacement(result *OwnershipGuardResult) {
	// Knowledge files should be in .kdse/knowledge/
	// Check for knowledge-related files outside the knowledge directory
	entries, err := os.ReadDir(g.repoPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "knowledge-") {
			violation := OwnershipViolation{
				Artifact:       entry.Name(),
				ActualPath:     entry.Name(),
				ExpectedDomain: "knowledge",
				SuggestedPath: filepath.Join(".kdse", "knowledge", entry.Name()),
				Severity:       "warning",
				Message:        "Knowledge artifact should be in .kdse/knowledge/",
			}
			result.Violations = append(result.Violations, violation)
		}
	}
}

// looksLikeReference checks if a directory name looks like a reference directory
func (g *OwnershipGuard) looksLikeReference(name string) bool {
	referenceIndicators := []string{
		"iec61850",
		"modbus",
		"ieee",
		"rfc",
		"vendor",
		"standard",
		"spec",
		"reference",
	}

	name = strings.ToLower(name)
	for _, indicator := range referenceIndicators {
		if strings.Contains(name, indicator) {
			return true
		}
	}
	return false
}

// pathExists checks if a path exists and returns (exists, isDir)
func (g *OwnershipGuard) pathExists(path string) (bool, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return false, false
	}
	return true, info.IsDir()
}

// ValidateArtifactPlacement validates a single artifact placement
func (g *OwnershipGuard) ValidateArtifactPlacement(artifactPath string) *OwnershipViolation {
	artifactName := filepath.Base(artifactPath)

	// Check if it's a project artifact in .kdse/
	if strings.HasPrefix(artifactPath, ".kdse/") {
		projectArtifacts := []string{
			"README.md",
			"LICENSE",
			"CHANGELOG.md",
			"src",
			"tests",
			"cmd",
			"internal",
			"docs",
			"deploy",
		}

		for _, pa := range projectArtifacts {
			if artifactName == pa || strings.HasPrefix(artifactPath, filepath.Join(".kdse", pa)) {
				return &OwnershipViolation{
					Artifact:       artifactName,
					ActualPath:     artifactPath,
					ExpectedDomain: "project",
					SuggestedPath:  filepath.Join(pa, filepath.Base(artifactPath)),
					Severity:       "error",
					Message:        "Project artifact should not be in .kdse/",
				}
			}
		}
	}

	// Check if it's a runtime artifact outside .kdse/
	runtimeArtifacts := []string{"laboratory", "runtime"}
	for _, ra := range runtimeArtifacts {
		if artifactName == ra || strings.HasPrefix(artifactPath, ra) {
			if !strings.HasPrefix(artifactPath, ".kdse/") {
				return &OwnershipViolation{
					Artifact:       artifactName,
					ActualPath:     artifactPath,
					ExpectedDomain: "runtime",
					SuggestedPath:  filepath.Join(".kdse", artifactPath),
					Severity:       "warning",
					Message:        "Runtime artifact should be in .kdse/",
				}
			}
		}
	}

	return nil
}
