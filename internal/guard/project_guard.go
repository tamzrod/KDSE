// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// ProjectGuard owns project discovery and validation.
// It is responsible for:
//   - Detecting engineering projects
//   - Rejecting invalid locations
//   - Rejecting generic directories
//   - Validating project eligibility
//
// Project Guard returns only structured results.
// It does NOT initialize KDSE.
// It does NOT create .kdse.
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
func (g *ProjectGuard) Validate(ctx context.Context) *ProjectGuardResult {
	result := &ProjectGuardResult{
		GuardResult: RuntimeGuardResult{
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

	// Step 2: Check if it looks like an engineering project
	isProject, artifacts := g.detectProjectArtifacts()
	if !isProject {
		result.Error = ErrGenericDirectory
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
	result.DetectedArtifacts = artifacts

	return result
}

// Exists performs a quick check without full validation.
// Returns true if a project appears to exist at the path.
func (g *ProjectGuard) Exists() bool {
	return g.isPathAccessible() && g.hasProjectIndicators()
}

// isPathAccessible checks if the repository path is accessible
func (g *ProjectGuard) isPathAccessible() bool {
	info, err := os.Stat(g.repoPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isGitRepo checks if the directory is a git repository
func (g *ProjectGuard) isGitRepo() bool {
	gitPath := filepath.Join(g.repoPath, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// detectProjectArtifacts detects engineering project artifacts
func (g *ProjectGuard) detectProjectArtifacts() (bool, []string) {
	var artifacts []string
	seen := make(map[string]bool)

	// Project indicator patterns
	projectIndicators := []struct {
		name      string
		patterns  []string
	}{
		{
			name:     "documentation",
			patterns: []string{"README.md", "docs/", "CONTRIBUTING.md", "LICENSE", "CHANGELOG.md"},
		},
		{
			name:     "source_code",
			patterns: []string{"src/", "lib/", "cmd/", "internal/", "main.go", "index.js", "app/"},
		},
		{
			name:     "testing",
			patterns: []string{"tests/", "test/", "__tests__/", ".test.", ".spec."},
		},
		{
			name:     "dependencies",
			patterns: []string{"package.json", "requirements.txt", "go.mod", "Cargo.toml", "pom.xml", "build.gradle", "Gemfile"},
		},
		{
			name:     "architecture",
			patterns: []string{"ARCHITECTURE.md", "architecture/", "SPEC.md"},
		},
		{
			name:     "configuration",
			patterns: []string{".gitignore", ".editorconfig", "Makefile", "Dockerfile"},
		},
	}

	// Scan directory tree
	err := filepath.Walk(g.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			// Skip hidden directories and common non-project directories
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" || name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(g.repoPath, path)
		if err != nil {
			return nil
		}

		if seen[relPath] {
			return nil
		}

		// Check against patterns
		for _, indicator := range projectIndicators {
			for _, pattern := range indicator.patterns {
				if g.matchesPattern(relPath, info, pattern) {
					dir := filepath.Dir(relPath)
					if !seen[dir] {
						artifacts = append(artifacts, dir)
						seen[dir] = true
					}
					seen[relPath] = true
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return false, nil
	}

	// Require at least 2 distinct artifact categories for a valid project
	categoryCount := 0
	categories := make(map[string]bool)
	for _, indicator := range projectIndicators {
		for _, pattern := range indicator.patterns {
			for _, artifact := range artifacts {
				if g.patternMatches(pattern, artifact) {
					categories[indicator.name] = true
				}
			}
		}
	}
	for range categories {
		categoryCount++
	}

	return categoryCount >= 2, artifacts
}

// matchesPattern checks if a path matches a pattern
func (g *ProjectGuard) matchesPattern(relPath string, info os.FileInfo, pattern string) bool {
	// Directory pattern (ends with /)
	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(relPath, pattern)
	}

	// Exact match
	if relPath == pattern {
		return true
	}

	// File extension pattern
	if strings.HasPrefix(pattern, ".") {
		return strings.HasSuffix(info.Name(), pattern)
	}

	// Prefix match
	return strings.HasPrefix(relPath, pattern)
}

// patternMatches checks if an artifact path matches a pattern
func (g *ProjectGuard) patternMatches(pattern, artifact string) bool {
	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(artifact, pattern) || artifact == strings.TrimSuffix(pattern, "/")
	}
	return artifact == pattern
}

// hasProjectIndicators performs a quick check for project indicators
func (g *ProjectGuard) hasProjectIndicators() bool {
	entries, err := os.ReadDir(g.repoPath)
	if err != nil {
		return false
	}

	// Quick heuristics for common project files
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files/dirs and common non-project items
		if strings.HasPrefix(name, ".") {
			continue
		}

		// Common project indicators
		switch name {
		case "README.md", "package.json", "go.mod", "Cargo.toml", "requirements.txt",
			"Gemfile", "pom.xml", "build.gradle", "Makefile", "Dockerfile":
			return true
		case "src", "lib", "cmd", "internal", "app", "docs":
			if !entry.IsDir() {
				continue
			}
			return true
		}
	}

	return false
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
