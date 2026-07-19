// Package discover provides centralized runtime discovery for KDSE.
//
// KDSE is project-centric: the runtime ALWAYS lives at <project root>/.kdse
//
// Project Discovery (NEW ARCHITECTURE):
// KDSE is an engineering runtime that ATTACHES to an existing software project.
// It does NOT create projects or Git repositories.
//
// Discovery Order:
// 1. Detect if a software project exists (language-specific files)
// 2. If project found, optionally detect Git repository (as evidence only)
// 3. Resolve to project root (not Git root - they may differ)
// 4. Runtime path = <project root>/.kdse
//
// Project Evidence (language-specific files):
// - go.mod (Go)
// - package.json (Node.js)
// - pyproject.toml / setup.py / requirements.txt (Python)
// - Cargo.toml (Rust)
// - pom.xml / build.gradle (Java)
// - *.sln / *.csproj (.NET)
// - composer.json (PHP)
// - Makefile / CMakeLists.txt (C/C++)
//
// Git is ONLY optional evidence of project validity.
// KDSE never requires Git to function.
package discover

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Common errors
var (
	// ErrNoProject indicates no software project was found.
	// KDSE requires a software project to exist before initialization.
	ErrNoProject = errors.New("no software project detected: KDSE requires a software project")

	// ErrNoGitRepository is kept for backward compatibility but Git is now optional.
	// Deprecated: Use ErrNoProject instead.
	ErrNoGitRepository = errors.New("no git repository found (Git is now optional evidence)")

	// ErrInvalidProjectPath indicates the provided project path is invalid.
	ErrInvalidProjectPath = errors.New("invalid project path")
)

// ProjectType identifies the detected project language/framework
type ProjectType string

const (
	ProjectTypeGo        ProjectType = "go"
	ProjectTypeNode      ProjectType = "node"
	ProjectTypePython    ProjectType = "python"
	ProjectTypeRust      ProjectType = "rust"
	ProjectTypeJava      ProjectType = "java"
	ProjectTypeDotNet    ProjectType = "dotnet"
	ProjectTypePHP       ProjectType = "php"
	ProjectTypeC         ProjectType = "c"
	ProjectTypeUnknown   ProjectType = "unknown"
)

// ProjectIndicators defines the files/directories that indicate a project type
var ProjectIndicators = map[ProjectType][]string{
	ProjectTypeGo:     {"go.mod", "go.sum", "main.go", "cmd/", "internal/"},
	ProjectTypeNode:    {"package.json", "package-lock.json", "yarn.lock", "node_modules/"},
	ProjectTypePython:  {"pyproject.toml", "setup.py", "requirements.txt", "Pipfile", "venv/", ".venv/"},
	ProjectTypeRust:    {"Cargo.toml", "Cargo.lock", "src/"},
	ProjectTypeJava:    {"pom.xml", "build.gradle", "build.gradle.kts", "src/main/java/"},
	ProjectTypeDotNet:  {".sln", ".csproj", ".fsproj", "Program.cs"},
	ProjectTypePHP:     {"composer.json", "artisan", "public/index.php"},
	ProjectTypeC:      {"Makefile", "CMakeLists.txt", "configure", "*.c", "*.h"},
}

// RuntimePaths contains all resolved paths for KDSE runtime
type RuntimePaths struct {
	// ProjectPath is the original project path provided or resolved from cwd
	ProjectPath string
	// RepositoryPath is the project root (may or may not be Git root)
	RepositoryPath string
	// RuntimePath is the KDSE runtime path (<project>/.kdse)
	RuntimePath string
	// ProjectRoot is the resolved project root (same as RepositoryPath)
	ProjectRoot string
	// IsGitRepo indicates if a Git repository was found (optional evidence)
	IsGitRepo bool
	// ProjectType indicates the detected project language/framework
	ProjectType ProjectType
	// ProjectIndicators lists the detected project evidence files
	ProjectIndicators []string
}

// String returns a string representation of the runtime paths
func (r *RuntimePaths) String() string {
	return r.RuntimePath
}

// Resolve discovers the KDSE runtime paths from a project path.
//
// This function implements the NEW KDSE project-first discovery rules:
// 1. If projectPath is provided, start from there
// 2. If projectPath is empty, use current working directory
// 3. Detect if a software project exists (language-specific files)
// 4. Optionally detect Git repository (as evidence only)
// 5. Runtime path = <project root>/.kdse
//
// Parameters:
//   - projectPath: The project path provided by the client (can be empty)
//
// Returns:
//   - RuntimePaths containing all resolved paths
//   - Error if no software project is found
func Resolve(projectPath string) (*RuntimePaths, error) {
	// Step 1: Determine the starting path
	startPath := projectPath
	if startPath == "" {
		// Use current working directory as fallback
		cwd, err := os.Getwd()
		if err != nil {
			return nil, ErrInvalidProjectPath
		}
		startPath = cwd
	}

	// Step 2: Normalize the path
	startPath = filepath.Clean(startPath)

	// Step 3: Check if path is accessible
	if err := validatePathAccessible(startPath); err != nil {
		return nil, ErrInvalidProjectPath
	}

	// Step 4: Detect software project
	projectRoot, projectType, indicators, err := detectProject(startPath)
	if err != nil {
		return nil, ErrNoProject
	}

	// Step 5: Optionally detect Git repository (as evidence only)
	isGitRepo := hasGitRepository(projectRoot)

	// Step 6: Compute runtime path
	runtimePath := filepath.Join(projectRoot, ".kdse")

	return &RuntimePaths{
		ProjectPath:        projectPath,
		RepositoryPath:     projectRoot,
		RuntimePath:        runtimePath,
		ProjectRoot:        projectRoot,
		IsGitRepo:          isGitRepo,
		ProjectType:        projectType,
		ProjectIndicators:  indicators,
	}, nil
}

// validatePathAccessible checks if a path is accessible
func validatePathAccessible(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return ErrInvalidProjectPath
	}
	return nil
}

// detectProject detects if a directory contains a software project.
// Returns the project root, detected project type, and indicators found.
func detectProject(startPath string) (string, ProjectType, []string, error) {
	// Scan upward for project indicators
	currentDir := startPath

	for {
		// Check for project indicators at this level
		projectType, indicators := scanForProjectIndicators(currentDir)
		if len(indicators) > 0 {
			return currentDir, projectType, indicators, nil
		}

		// Move to parent directory
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached filesystem root
			break
		}
		currentDir = parent
	}

	// No project found
	return "", ProjectTypeUnknown, nil, ErrNoProject
}

// scanForProjectIndicators scans a directory for project indicators
func scanForProjectIndicators(dir string) (ProjectType, []string) {
	var foundIndicators []string
	detectedType := ProjectTypeUnknown

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ProjectTypeUnknown, nil
	}

	// Create a map of files/dirs for quick lookup
	fileMap := make(map[string]bool)
	for _, entry := range entries {
		fileMap[entry.Name()] = true
	}

	// Check each project type
	for projectType, indicators := range ProjectIndicators {
		matched := 0
		var matchedIndicators []string

		for _, indicator := range indicators {
			if fileMap[indicator] {
				matched++
				matchedIndicators = append(matchedIndicators, indicator)
			}
		}

		// Require at least 2 indicators for a valid project
		if matched >= 2 && matched > len(foundIndicators) {
			detectedType = projectType
			foundIndicators = matchedIndicators
		}
	}

	// Also check for common documentation/source files as additional evidence
	commonFiles := []string{"README.md", "LICENSE", "CONTRIBUTING.md", "CHANGELOG.md", "docs/", "src/", "lib/"}
	for _, file := range commonFiles {
		if fileMap[file] {
			foundIndicators = append(foundIndicators, file)
		}
	}

	return detectedType, foundIndicators
}

// hasGitRepository checks if the directory is inside a Git repository.
// This is OPTIONAL evidence - KDSE does not require Git.
func hasGitRepository(dir string) bool {
	// Method 1: Use git rev-parse --show-toplevel
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "true"
	}

	// Method 2: Walk up looking for .git
	gitPath := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return true
	}

	// Check parent directories
	parent := filepath.Dir(dir)
	if parent != dir {
		return hasGitRepository(parent)
	}

	return false
}

// ResolveRuntime is a convenience function that returns only the runtime path.
// Returns an error if no software project is found.
func ResolveRuntime(projectPath string) (string, error) {
	paths, err := Resolve(projectPath)
	if err != nil {
		return "", err
	}
	return paths.RuntimePath, nil
}

// ResolveRepository is a convenience function that returns only the repository path.
// Returns an error if no software project is found.
func ResolveRepository(projectPath string) (string, error) {
	paths, err := Resolve(projectPath)
	if err != nil {
		return "", err
	}
	return paths.RepositoryPath, nil
}

// MustResolve is like Resolve but panics on error.
// Use only when you're certain the path is valid.
func MustResolve(projectPath string) *RuntimePaths {
	paths, err := Resolve(projectPath)
	if err != nil {
		panic("discover.MustResolve: " + err.Error())
	}
	return paths
}

// HasGitRepository checks if the given directory is inside a Git repository.
// DEPRECATED: Git is now optional evidence. Use Resolve() and check RuntimePaths.IsGitRepo instead.
func HasGitRepository(dir string) bool {
	return hasGitRepository(dir)
}

// EnsureRuntime ensures the .kdse runtime directory exists.
// Returns the runtime path and an error if no project exists or runtime cannot be created.
func EnsureRuntime(projectPath string) (string, error) {
	paths, err := Resolve(projectPath)
	if err != nil {
		return "", err
	}

	// Create the runtime directory if it doesn't exist
	if err := os.MkdirAll(paths.RuntimePath, 0755); err != nil {
		return "", err
	}

	return paths.RuntimePath, nil
}

// DetectProjectType returns the detected project type for a path.
// This is a convenience function for quick project type detection.
func DetectProjectType(projectPath string) (ProjectType, error) {
	if projectPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return ProjectTypeUnknown, err
		}
		projectPath = cwd
	}

	paths, err := Resolve(projectPath)
	if err != nil {
		return ProjectTypeUnknown, err
	}

	return paths.ProjectType, nil
}

// HasProject checks if a valid software project exists at the given path.
// This is a quick check without full resolution.
func HasProject(projectPath string) bool {
	if projectPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return false
		}
		projectPath = cwd
	}

	_, _, indicators, err := detectProject(projectPath)
	return err == nil && len(indicators) > 0
}
