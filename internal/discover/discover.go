// Package discover provides centralized runtime discovery for KDSE.
// 
// KDSE is repository-centric: the runtime ALWAYS lives at <git root>/.kdse
// 
// This package ensures consistent runtime discovery across all KDSE components:
// - CLI commands
// - MCP server tools
// - HTTP server endpoints
// - Internal packages
//
// Discovery Rules:
// 1. Runtime must NEVER be resolved using os.Getwd() or filepath.Abs(".")
// 2. Runtime MUST be resolved from the project path provided by the client
// 3. If no project path is provided, resolve from current working directory
// 4. Always resolve to Git repository root (git rev-parse --show-toplevel equivalent)
// 5. Runtime path = <git root>/.kdse
// 6. If no Git repository exists, return an error (NEVER create runtime outside repo)
package discover

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// Common errors
var (
	// ErrNoGitRepository indicates no Git repository was found.
	// KDSE requires a Git repository to function.
	ErrNoGitRepository = errors.New("no git repository found: KDSE requires a Git repository")

	// ErrInvalidProjectPath indicates the provided project path is invalid.
	ErrInvalidProjectPath = errors.New("invalid project path")
)

// RuntimePaths contains all resolved paths for KDSE runtime
type RuntimePaths struct {
	// ProjectPath is the original project path provided or resolved from cwd
	ProjectPath string
	// RepositoryPath is the Git repository root (equivalent to git rev-parse --show-toplevel)
	RepositoryPath string
	// RuntimePath is the KDSE runtime path (<repo>/.kdse)
	RuntimePath string
	// GitRoot is the resolved Git root (same as RepositoryPath)
	GitRoot string
	// IsGitRepo indicates if a Git repository was found
	IsGitRepo bool
}

// String returns a string representation of the runtime paths
func (r *RuntimePaths) String() string {
	return r.RuntimePath
}

// Resolve discovers the KDSE runtime paths from a project path.
//
// This function implements the core KDSE runtime discovery rules:
// 1. If projectPath is provided, start from there
// 2. If projectPath is empty, use current working directory
// 3. Resolve to Git repository root
// 4. Runtime path = <git root>/.kdse
//
// Parameters:
//   - projectPath: The project path provided by the client (can be empty)
//
// Returns:
//   - RuntimePaths containing all resolved paths
//   - Error if no Git repository is found
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

	// Step 3: Resolve to Git repository root
	gitRoot, err := resolveGitRoot(startPath)
	if err != nil {
		return nil, err
	}

	// Step 4: Compute runtime path
	runtimePath := filepath.Join(gitRoot, ".kdse")

	return &RuntimePaths{
		ProjectPath:    projectPath,
		RepositoryPath: gitRoot,
		RuntimePath:    runtimePath,
		GitRoot:        gitRoot,
		IsGitRepo:      true,
	}, nil
}

// ResolveRuntime is a convenience function that returns only the runtime path.
// Returns an error if no Git repository is found.
func ResolveRuntime(projectPath string) (string, error) {
	paths, err := Resolve(projectPath)
	if err != nil {
		return "", err
	}
	return paths.RuntimePath, nil
}

// ResolveRepository is a convenience function that returns only the repository path.
// Returns an error if no Git repository is found.
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

// resolveGitRoot finds the Git repository root starting from the given path.
// Uses git rev-parse --show-toplevel as the primary method.
func resolveGitRoot(startPath string) (string, error) {
	// Method 1: Use git rev-parse --show-toplevel
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = startPath

	output, err := cmd.Output()
	if err == nil {
		// Success - parse the output
		root := filepath.Clean(string(output))
		root = filepath.ToSlash(root)
		// Remove trailing slashes
		for len(root) > 0 && root[len(root)-1] == '/' {
			root = root[:len(root)-1]
		}
		return root, nil
	}

	// Method 2: Walk up the directory tree looking for .git
	root, err := findGitRoot(startPath)
	if err != nil {
		return "", ErrNoGitRepository
	}

	return root, nil
}

// findGitRoot walks up the directory tree looking for a .git directory or file.
func findGitRoot(dir string) (string, error) {
	// Normalize the directory
	dir = filepath.Clean(dir)

	// Check if this directory has a .git
	gitPath := filepath.Join(dir, ".git")
	if info, err := os.Stat(gitPath); err == nil {
		// .git exists - verify it's a directory or file (submodule)
		if info.IsDir() || isGitFile(gitPath) {
			return dir, nil
		}
	}

	// Move to parent directory
	parent := filepath.Dir(dir)
	if parent == dir {
		// Reached the filesystem root
		return "", ErrNoGitRepository
	}

	return findGitRoot(parent)
}

// isGitFile checks if path is a gitlink file (submodule reference).
func isGitFile(path string) bool {
	// Git submodules have .git as a file containing "gitdir: /path/to/actual/.git"
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	content := string(data)
	return len(content) > 7 && content[:7] == "gitdir:"
}

// HasGitRepository checks if the given directory is inside a Git repository.
func HasGitRepository(dir string) bool {
	_, err := resolveGitRoot(dir)
	return err == nil
}

// EnsureRuntime ensures the .kdse runtime directory exists.
// Returns the runtime path and an error if the runtime cannot be created.
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
