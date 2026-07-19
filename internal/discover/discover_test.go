package discover

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Helper to create a temp git repo with a Go project for testing
// Project files are now required for discovery, not just Git
func setupTestRepo(tb testing.TB) (string, func()) {
	tb.Helper()

	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		tb.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo using real git commands
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		tb.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git for test
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	_ = cmd.Run() // Ignore errors for config

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	_ = cmd.Run() // Ignore errors for config

	// Create a minimal Go project structure (project files required for discovery)
	projectFiles := []string{
		"go.mod",
		"main.go",
		"README.md",
	}
	for _, file := range projectFiles {
		path := filepath.Join(tmpDir, file)
		content := "package main\n"
		if file == "go.mod" {
			content = "module github.com/test/project\n\ngo 1.21\n"
		} else if file == "README.md" {
			content = "# Test Project\n"
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			os.RemoveAll(tmpDir)
			tb.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// TestResolveFromRepoRoot tests resolution from repository root
func TestResolveFromRepoRoot(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if !paths.IsGitRepo {
		t.Error("Expected IsGitRepo to be true")
	}

	if paths.RepositoryPath != tmpDir {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir, paths.RepositoryPath)
	}

	expectedRuntime := filepath.Join(tmpDir, ".kdse")
	if paths.RuntimePath != expectedRuntime {
		t.Errorf("Expected RuntimePath to be %s, got %s", expectedRuntime, paths.RuntimePath)
	}
}

// TestResolveFromNestedDirectory tests resolution from a nested directory
func TestResolveFromNestedDirectory(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create nested directories
	nestedDir := filepath.Join(tmpDir, "cmd", "server", "internal")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	paths, err := Resolve(nestedDir)
	if err != nil {
		t.Fatalf("Resolve failed from nested dir: %v", err)
	}

	if paths.RepositoryPath != tmpDir {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir, paths.RepositoryPath)
	}

	expectedRuntime := filepath.Join(tmpDir, ".kdse")
	if paths.RuntimePath != expectedRuntime {
		t.Errorf("Expected RuntimePath to be %s, got %s", expectedRuntime, paths.RuntimePath)
	}
}

// TestResolveFromEmptyPath tests resolution with empty path (uses cwd)
func TestResolveFromEmptyPath(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Save original cwd
	origCwd, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get current working directory")
	}
	defer os.Chdir(origCwd)

	// Create temp git repo and change to it
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	if err := os.Chdir(tmpDir); err != nil {
		t.Skip("Cannot change directory")
	}

	paths, err := Resolve("")
	if err != nil {
		t.Fatalf("Resolve with empty path failed: %v", err)
	}

	if paths.RepositoryPath != tmpDir {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir, paths.RepositoryPath)
	}
}

// TestResolveNoProject tests that resolution fails without a software project
// Note: Git is now OPTIONAL - the test verifies project detection, not Git detection
func TestResolveNoProject(t *testing.T) {
	// Create a temp directory that is NOT a software project (no language files)
	tmpDir, err := os.MkdirTemp("", "kdse-no-project-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = Resolve(tmpDir)
	if err == nil {
		t.Error("Expected error for non-project directory")
	}
	if err != ErrNoProject {
		t.Errorf("Expected ErrNoProject, got %v", err)
	}
}

// TestResolveWithProjectButNoGit tests that resolution succeeds with a project
// even if there's no Git repository (Git is now optional)
func TestResolveWithProjectButNoGit(t *testing.T) {
	// Create a temp directory with project files but NO git repo
	tmpDir, err := os.MkdirTemp("", "kdse-project-no-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create project indicators (Go project)
	projectFiles := []string{"go.mod", "main.go"}
	for _, file := range projectFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("package main\n"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Expected success for project without git, got: %v", err)
	}

	// Should succeed with project detection
	if paths.RepositoryPath != tmpDir {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir, paths.RepositoryPath)
	}

	// Git should be false (we didn't create a git repo)
	if paths.IsGitRepo {
		t.Error("Expected IsGitRepo to be false")
	}

	// Project type should be detected
	if paths.ProjectType != ProjectTypeGo {
		t.Errorf("Expected ProjectType to be %s, got %s", ProjectTypeGo, paths.ProjectType)
	}
}

// TestResolveFromAnotherGitRepo tests that we find the correct git repo
// when starting from a directory that is itself a git repo
func TestResolveFromAnotherGitRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create a temp git repo
	tmpDir1, cleanup1 := setupTestRepo(t)
	defer cleanup1()

	// Create a nested directory in tmpDir1
	nestedDir := filepath.Join(tmpDir1, "submodule", "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Resolve from nested directory in tmpDir1
	paths, err := Resolve(nestedDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Should resolve to tmpDir1's git root
	if paths.RepositoryPath != tmpDir1 {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir1, paths.RepositoryPath)
	}
}

// TestResolveRuntime tests the convenience function
func TestResolveRuntime(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	runtimePath, err := ResolveRuntime(tmpDir)
	if err != nil {
		t.Fatalf("ResolveRuntime failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".kdse")
	if runtimePath != expected {
		t.Errorf("Expected %s, got %s", expected, runtimePath)
	}
}

// TestResolveRepository tests the convenience function
func TestResolveRepository(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	repoPath, err := ResolveRepository(tmpDir)
	if err != nil {
		t.Fatalf("ResolveRepository failed: %v", err)
	}

	if repoPath != tmpDir {
		t.Errorf("Expected %s, got %s", tmpDir, repoPath)
	}
}

// TestMustResolvePanics tests that MustResolve panics on error
func TestMustResolvePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustResolve to panic")
		}
	}()

	// Create a temp directory that is NOT a git repo
	tmpDir, err := os.MkdirTemp("", "kdse-no-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	MustResolve(tmpDir)
}

// TestEnsureRuntime tests runtime directory creation
func TestEnsureRuntime(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	runtimePath, err := EnsureRuntime(tmpDir)
	if err != nil {
		t.Fatalf("EnsureRuntime failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".kdse")
	if runtimePath != expected {
		t.Errorf("Expected %s, got %s", expected, runtimePath)
	}

	// Check that directory was created
	info, err := os.Stat(runtimePath)
	if err != nil {
		t.Errorf("Runtime directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Runtime path is not a directory")
	}
}

// TestRuntimePathsString tests the String method
func TestRuntimePathsString(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if paths.String() != paths.RuntimePath {
		t.Errorf("String() should return RuntimePath")
	}
}

// TestInvalidProjectPath tests handling of invalid paths
func TestInvalidProjectPath(t *testing.T) {
	// Non-existent directory
	_, err := Resolve("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

// TestHasGitRepository tests the helper function
func TestHasGitRepository(t *testing.T) {
	// Create temp git repo
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	if !HasGitRepository(tmpDir) {
		t.Error("Expected HasGitRepository to return true for git repo")
	}

	// Create temp non-git directory
	tmpDir2, err := os.MkdirTemp("", "kdse-no-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir2)

	if HasGitRepository(tmpDir2) {
		t.Error("Expected HasGitRepository to return false for non-git directory")
	}
}

// TestSubmoduleGitRepo tests handling of git submodules
func TestSubmoduleGitRepo(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Remove the .git directory and create a .git file (submodule reference)
	gitDir := filepath.Join(tmpDir, ".git")
	os.RemoveAll(gitDir)

	gitFile := filepath.Join(tmpDir, ".git")
	gitContent := "gitdir: /some/path/to/.git/modules/submodule"
	if err := os.WriteFile(gitFile, []byte(gitContent), 0644); err != nil {
		t.Fatalf("Failed to create .git file: %v", err)
	}

	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed for submodule: %v", err)
	}

	if !paths.IsGitRepo {
		t.Error("Expected IsGitRepo to be true for submodule")
	}
}

// TestDifferentCwdThanProject tests that cwd doesn't affect resolution
// when project path is explicitly provided
func TestDifferentCwdThanProject(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Save original cwd
	origCwd, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get current working directory")
	}
	defer os.Chdir(origCwd)

	// Create temp git repo
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Change to a different directory (simulating container cwd)
	otherDir, err := os.MkdirTemp("", "kdse-other-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(otherDir)

	if err := os.Chdir(otherDir); err != nil {
		t.Skip("Cannot change directory")
	}

	// Resolve with explicit project path - should ignore cwd
	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Should resolve to tmpDir, not otherDir
	if paths.RepositoryPath != tmpDir {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir, paths.RepositoryPath)
	}
}

// TestRepositoryWithGitHub tests repository with .github directory
func TestRepositoryWithGitHub(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create .github directory
	githubDir := filepath.Join(tmpDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatalf("Failed to create .github dir: %v", err)
	}

	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Should resolve correctly regardless of .github presence
	if paths.RepositoryPath != tmpDir {
		t.Errorf("Expected RepositoryPath to be %s, got %s", tmpDir, paths.RepositoryPath)
	}

	expectedRuntime := filepath.Join(tmpDir, ".kdse")
	if paths.RuntimePath != expectedRuntime {
		t.Errorf("Expected RuntimePath to be %s, got %s", expectedRuntime, paths.RuntimePath)
	}
}

// TestRepositoryAlreadyHasKDSE tests repository that already has .kdse
func TestRepositoryAlreadyHasKDSE(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Pre-create .kdse directory
	kdseDir := filepath.Join(tmpDir, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("Failed to create .kdse dir: %v", err)
	}

	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Should still resolve correctly
	if paths.RuntimePath != kdseDir {
		t.Errorf("Expected RuntimePath to be %s, got %s", kdseDir, paths.RuntimePath)
	}

	// Verify .kdse exists
	info, err := os.Stat(kdseDir)
	if err != nil || !info.IsDir() {
		t.Error("Expected .kdse to exist")
	}
}

// TestRuntimeIndependentOfServerCwd tests that runtime resolution is independent
// of the server's working directory (the core bug fix)
func TestRuntimeIndependentOfServerCwd(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Save original cwd
	origCwd, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get current working directory")
	}
	defer os.Chdir(origCwd)

	// Create a git repo at a specific location
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Change to a completely different directory (simulating server at /app)
	otherDir := "/tmp" // Use /tmp as different location
	if err := os.Chdir(otherDir); err != nil {
		t.Skip("Cannot change to /tmp")
	}

	// Verify cwd is now different
	newCwd, _ := os.Getwd()
	if newCwd == tmpDir {
		t.Skip("Could not change directory away from tmpDir")
	}

	// Resolve with explicit project path - should still resolve to tmpDir
	paths, err := Resolve(tmpDir)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// THE CRITICAL TEST: RepositoryPath should be tmpDir, NOT newCwd (/tmp)
	// This proves the bug fix works
	if paths.RepositoryPath != tmpDir {
		t.Errorf("BUG: RepositoryPath resolved to %s instead of %s (server cwd: %s)",
			paths.RepositoryPath, tmpDir, newCwd)
	}

	// Runtime should be at tmpDir/.kdse, NOT /tmp/.kdse
	expectedRuntime := filepath.Join(tmpDir, ".kdse")
	if paths.RuntimePath != expectedRuntime {
		t.Errorf("BUG: RuntimePath resolved to %s instead of %s (server cwd: %s)",
			paths.RuntimePath, expectedRuntime, newCwd)
	}
}

// BenchmarkResolve benchmarks the Resolve function
func BenchmarkResolve(b *testing.B) {
	tmpDir, cleanup := setupTestRepo(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Resolve(tmpDir)
	}
}
