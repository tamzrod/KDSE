// Package detection provides project and repository detection utilities.
package detection

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GitResolver resolves the Git repository root using git commands.
// This ensures KDSE always works relative to the actual repository, not container paths.
type GitResolver struct {
	workingDir string
}

// NewGitResolver creates a new Git resolver starting from the given directory.
func NewGitResolver(workingDir string) *GitResolver {
	return &GitResolver{workingDir: workingDir}
}

// ResolveRoot finds the root of the Git repository.
// Returns the absolute path to the repository root if found.
// Returns an error if no Git repository is found.
func (r *GitResolver) ResolveRoot() (string, error) {
	// First try: git rev-parse --show-toplevel from the working directory
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = r.workingDir

	output, err := cmd.Output()
	if err == nil {
		// Success - trim whitespace and return
		root := filepath.Clean(string(output))
		root = filepath.ToSlash(root)
		// Remove trailing slashes
		for len(root) > 0 && root[len(root)-1] == '/' {
			root = root[:len(root)-1]
		}
		return root, nil
	}

	// Fallback: walk up the directory tree looking for .git
	root, err := r.findGitRoot(r.workingDir)
	if err != nil {
		return "", err
	}

	return root, nil
}

// findGitRoot walks up the directory tree looking for a .git directory or file.
func (r *GitResolver) findGitRoot(dir string) (string, error) {
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
		return "", &GitNotFoundError{Path: r.workingDir}
	}

	return r.findGitRoot(parent)
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
	_, err := NewGitResolver(dir).ResolveRoot()
	return err == nil
}

// GitNotFoundError indicates no Git repository was found.
type GitNotFoundError struct {
	Path string
}

func (e *GitNotFoundError) Error() string {
	return "no Git repository found starting from: " + e.Path
}
