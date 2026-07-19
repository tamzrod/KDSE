// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"os"
	"path/filepath"

	"github.com/kdse/runtime/internal/discover"
)

// ProjectGuard owns project discovery and validation.
// It is responsible for:
//   - Detecting engineering projects via language-specific files
//   - Rejecting invalid locations
//   - Rejecting generic directories
//   - Validating project eligibility
//
// Project Guard does NOT require Git - Git is only optional evidence.
// It does NOT initialize KDSE. It does NOT create .kdse.
type ProjectGuard struct {
	repoPath string
}

// NewProjectGuard creates a new Project Guard for the given repository path
func NewProjectGuard(repoPath string) *ProjectGuard {
	return &ProjectGuard{
		repoPath: repoPath,
	}
}

// Validate performs project discovery and validation.
// Returns ProjectGuardResult indicating if a valid project was detected.
//
// This uses the new project-first discovery:
// 1. Detect software project via language-specific files
// 2. Optionally detect Git repository (as evidence only)
func (g *ProjectGuard) Validate(ctx context.Context) *ProjectGuardResult {
	result := &ProjectGuardResult{
		RuntimeGuardResult: RuntimeGuardResult{
			Valid:        false,
			State:        StateNoProject,
			StateBefore:  StateNoProject,
			StateAfter:   StateNoProject,
			GuardType:    GuardTypeProject,
		},
	}

	// Step 1: Check if path is accessible
	if !g.isPathAccessible() {
		result.Error = ErrInvalidProjectLocation
		return result
	}

	// Step 2: Check if it looks like an engineering project using discover package
	isProject, projectType, indicators := g.detectProject()
	if !isProject {
		result.Error = ErrNoProjectDetected
		return result
	}

	// Step 3: Additional validation
	if err := g.validateProjectEligibility(); err != nil {
		result.Error = ErrInvalidProjectLocation
		return result
	}

	// Project is valid
	result.Valid = true
	result.State = StateProject
	result.StateBefore = StateNoProject
	result.StateAfter = StateProject
	result.ProjectPath = g.repoPath
	result.ProjectName = filepath.Base(g.repoPath)
	result.IsGitRepo = g.isGitRepo()
	result.DetectedArtifacts = indicators

	// Store project type in artifacts if available
	if projectType != "" {
		result.DetectedArtifacts = append(result.DetectedArtifacts, "type:"+string(projectType))
	}

	return result
}

// Exists performs a quick check without full validation.
// Returns true if a project appears to exist at the path.
func (g *ProjectGuard) Exists() bool {
	return g.isPathAccessible() && discover.HasProject(g.repoPath)
}

// isPathAccessible checks if the repository path is accessible
func (g *ProjectGuard) isPathAccessible() bool {
	info, err := os.Stat(g.repoPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isGitRepo checks if the directory is a git repository (OPTIONAL evidence)
func (g *ProjectGuard) isGitRepo() bool {
	gitPath := filepath.Join(g.repoPath, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// detectProject uses the discover package to detect a software project
func (g *ProjectGuard) detectProject() (bool, discover.ProjectType, []string) {
	paths, err := discover.Resolve(g.repoPath)
	if err != nil {
		return false, "", nil
	}

	return true, paths.ProjectType, paths.ProjectIndicators
}

// validateProjectEligibility performs additional eligibility checks
func (g *ProjectGuard) validateProjectEligibility() error {
	// Check if directory is readable
	testFile := filepath.Join(g.repoPath, ".kdse.test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return err
	}
	os.Remove(testFile)

	return nil
}
