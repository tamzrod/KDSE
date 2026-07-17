package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kdse/runtime/internal/workspace"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.TemplateName != "core" {
		t.Errorf("Expected template name 'core', got %s", cfg.TemplateName)
	}
	if cfg.TemplateOwner != "kdse" {
		t.Errorf("Expected template owner 'kdse', got %s", cfg.TemplateOwner)
	}
	if cfg.TemplateRepo != "workspace-templates" {
		t.Errorf("Expected template repo 'workspace-templates', got %s", cfg.TemplateRepo)
	}
	if cfg.TemplateRef != "main" {
		t.Errorf("Expected template ref 'main', got %s", cfg.TemplateRef)
	}
}

func TestNewBootstrapper(t *testing.T) {
	repoPath := "/test/repo"
	cfg := DefaultConfig()
	
	b := NewBootstrapper(repoPath, cfg)

	if b.repoPath != repoPath {
		t.Errorf("Expected repo path %s, got %s", repoPath, b.repoPath)
	}
	if b.templateName != cfg.TemplateName {
		t.Errorf("Expected template name %s, got %s", cfg.TemplateName, b.templateName)
	}
	if b.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}
	if b.workspace == nil {
		t.Error("Expected workspace to be initialized")
	}
}

func TestNewBootstrapperWithNilConfig(t *testing.T) {
	repoPath := "/test/repo"
	b := NewBootstrapper(repoPath, nil)

	if b.templateName != "core" {
		t.Errorf("Expected default template name 'core', got %s", b.templateName)
	}
}

func TestExtractTarGz(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "kdse-extract-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source directory with test files
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	testContent := []byte("test content")
	if err := os.WriteFile(filepath.Join(srcDir, "test.txt"), testContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a subdirectory with a file
	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.txt"), testContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create tar.gz archive using system tar command
	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	archive, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	archive.Close()
	
	// Use os/exec to create tar.gz - skip test if tar not available
	// This is a basic structure test; actual integration tests would use pre-created archives

	// Test with invalid archive
	invalidPath := filepath.Join(tmpDir, "invalid.tar.gz")
	if err := os.WriteFile(invalidPath, []byte("not a valid tar.gz"), 0644); err != nil {
		t.Fatal(err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := extractTarGz(invalidPath, extractDir); err == nil {
		t.Error("Expected error for invalid archive, got nil")
	}
}

func TestExtractTarGzInvalidArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-extract-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an invalid archive
	invalidPath := filepath.Join(tmpDir, "invalid.tar.gz")
	if err := os.WriteFile(invalidPath, []byte("not a valid tar.gz"), 0644); err != nil {
		t.Fatal(err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := extractTarGz(invalidPath, extractDir); err == nil {
		t.Error("Expected error for invalid archive, got nil")
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-copy-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content for copy")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Copy to destination
	destPath := filepath.Join(tmpDir, "dest.txt")
	if err := copyFile(srcPath, destPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify copied content
	copied, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copied) != string(content) {
		t.Errorf("Copied content mismatch: expected %s, got %s", content, copied)
	}

	// Verify permissions
	srcInfo, _ := os.Stat(srcPath)
	destInfo, _ := os.Stat(destPath)
	if srcInfo.Mode() != destInfo.Mode() {
		t.Errorf("Permissions mismatch: expected %v, got %v", srcInfo.Mode(), destInfo.Mode())
	}
}

func TestCopyDirContents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-copydir-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source directory structure
	srcDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create files and directories
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Copy to destination
	destDir := filepath.Join(tmpDir, "dest")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatal(err)
	}

	var createdPaths []string
	if err := copyDirContents(srcDir, destDir, &createdPaths); err != nil {
		t.Fatalf("copyDirContents failed: %v", err)
	}

	// Verify all files were copied
	if _, err := os.Stat(filepath.Join(destDir, "file1.txt")); os.IsNotExist(err) {
		t.Error("file1.txt was not copied")
	}

	if _, err := os.Stat(filepath.Join(destDir, "subdir", "file2.txt")); os.IsNotExist(err) {
		t.Error("subdir/file2.txt was not copied")
	}

	// Verify createdPaths includes all copied paths
	if len(createdPaths) == 0 {
		t.Error("createdPaths should not be empty")
	}
}

func TestVerifyInstallation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-verify-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create workspace with required directories using workspace.New
	ws := workspace.New(tmpDir)
	wsRoot := ws.Root()
	
	if err := os.MkdirAll(wsRoot, 0755); err != nil {
		t.Fatal(err)
	}

	requiredDirs := []string{
		"knowledge",
		"architecture",
		"implementation",
		"verification",
		"reports",
		"docs",
	}

	for _, dir := range requiredDirs {
		if err := os.MkdirAll(filepath.Join(wsRoot, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create bootstrapper and verify
	b := &Bootstrapper{
		repoPath:  tmpDir,
		workspace: ws,
	}

	if err := b.verifyInstallation(); err != nil {
		t.Errorf("verifyInstallation failed: %v", err)
	}
}

func TestVerifyInstallationMissingDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-verify-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create workspace with only some directories
	ws := workspace.New(tmpDir)
	wsRoot := ws.Root()
	
	if err := os.MkdirAll(wsRoot, 0755); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Join(wsRoot, "knowledge"), 0755)

	b := &Bootstrapper{
		repoPath:  tmpDir,
		workspace: ws,
	}

	if err := b.verifyInstallation(); err == nil {
		t.Error("Expected error for missing directories, got nil")
	}
}

func TestRollback(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-rollback-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some test paths
	path1 := filepath.Join(tmpDir, "created1")
	path2 := filepath.Join(tmpDir, "created2")

	if err := os.MkdirAll(path1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(path2, 0755); err != nil {
		t.Fatal(err)
	}

	// Create bootstrapper and rollback
	b := &Bootstrapper{repoPath: tmpDir}
	b.rollback([]string{path1, path2})

	// Verify paths were removed
	if _, err := os.Stat(path1); !os.IsNotExist(err) {
		t.Error("path1 should have been removed")
	}
	if _, err := os.Stat(path2); !os.IsNotExist(err) {
		t.Error("path2 should have been removed")
	}
}

func TestWriteSessionMetadata(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-session-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ws := workspace.New(tmpDir)
	wsRoot := ws.Root()
	
	if err := os.MkdirAll(wsRoot, 0755); err != nil {
		t.Fatal(err)
	}

	b := &Bootstrapper{
		repoPath:     tmpDir,
		templateName: "core",
		workspace:    ws,
	}

	sessionPath := filepath.Join(wsRoot, "session.yaml")
	if err := b.writeSessionMetadata(sessionPath, "2.0", "abc1234"); err != nil {
		t.Fatalf("writeSessionMetadata failed: %v", err)
	}

	// Verify session file was created
	content, err := os.ReadFile(sessionPath)
	if err != nil {
		t.Fatalf("Failed to read session file: %v", err)
	}

	// Check for expected content
	contentStr := string(content)
	if !contains(contentStr, "kdse_version: 2.0") {
		t.Error("session.yaml should contain kdse_version: 2.0")
	}
	if !contains(contentStr, "template_version: core") {
		t.Error("session.yaml should contain template_version: core")
	}
	if !contains(contentStr, "template_commit: abc1234") {
		t.Error("session.yaml should contain template_commit: abc1234")
	}
	if !contains(contentStr, "initialized:") {
		t.Error("session.yaml should contain initialized timestamp")
	}
}

func TestGetKDSEVersion(t *testing.T) {
	b := &Bootstrapper{}
	version := b.getKDSEVersion()

	if version != "2.0" {
		t.Errorf("Expected KDSE version 2.0, got %s", version)
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	ts := getCurrentTimestamp()

	// Verify it's not empty
	if ts == "" {
		t.Error("getCurrentTimestamp returned empty string")
	}

	// Verify it contains expected format elements
	if len(ts) < 10 {
		t.Errorf("Timestamp seems too short: %s", ts)
	}
}

func TestGetCurrentTimestampShort(t *testing.T) {
	ts := getCurrentTimestampShort()

	// Verify format is YYYY-MM-DD
	if len(ts) != 10 {
		t.Errorf("Expected format YYYY-MM-DD, got %s", ts)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
