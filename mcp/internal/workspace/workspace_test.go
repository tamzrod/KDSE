package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	ws := New("/repo")
	if ws.RepoPath() != "/repo" {
		t.Errorf("Expected repo path /repo, got %s", ws.RepoPath())
	}
	if ws.Root() != "/repo/.kdse" {
		t.Errorf("Expected .kdse path /repo/.kdse, got %s", ws.Root())
	}
}

func TestSubPath(t *testing.T) {
	ws := New("/repo")
	expected := "/repo/.kdse/foundation"
	if got := ws.SubPath("foundation"); got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestInitialize(t *testing.T) {
	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ws := New(tmpDir)

	// Initialize should create .kdse/
	if err := ws.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Check .kdse/ exists
	if !ws.Exists() {
		t.Error(".kdse/ directory was not created")
	}
}

func TestEnsureSubdir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ws := New(tmpDir)
	ws.Initialize()

	// Ensure a subdirectory
	if err := ws.EnsureSubdir("knowledge"); err != nil {
		t.Fatalf("EnsureSubdir failed: %v", err)
	}

	// Check it exists
	expectedPath := filepath.Join(tmpDir, ".kdse", "knowledge")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Subdirectory was not created at %s", expectedPath)
	}
}

func TestDetectLegacyDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create legacy directories
	os.MkdirAll(filepath.Join(tmpDir, "foundation"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "knowledge"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "not-legacy"), 0755) // Should not be detected

	ws := New(tmpDir)
	legacy := ws.DetectLegacyDirs()

	if len(legacy) != 2 {
		t.Errorf("Expected 2 legacy dirs, got %d: %v", len(legacy), legacy)
	}

	found := map[string]bool{}
	for _, d := range legacy {
		found[d] = true
	}

	if !found["foundation"] {
		t.Error("foundation should be detected")
	}
	if !found["knowledge"] {
		t.Error("knowledge should be detected")
	}
	if found["not-legacy"] {
		t.Error("not-legacy should not be detected")
	}
}

func TestCheckMigration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create legacy directory
	os.MkdirAll(filepath.Join(tmpDir, "foundation"), 0755)

	ws := New(tmpDir)
	report := ws.CheckMigration()

	if !report.HasLegacyDirs {
		t.Error("Report should indicate legacy dirs exist")
	}
	if len(report.LegacyDirs) != 1 {
		t.Errorf("Expected 1 legacy dir, got %d", len(report.LegacyDirs))
	}
}

func TestMigrate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create legacy directory with content
	legacyPath := filepath.Join(tmpDir, "foundation")
	os.MkdirAll(legacyPath, 0755)
	os.WriteFile(filepath.Join(legacyPath, "test.md"), []byte("# Test"), 0644)

	ws := New(tmpDir)
	result, err := ws.Migrate()

	if err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	if !result.Success {
		t.Error("Migration should succeed")
	}

	if len(result.Migrated) != 1 {
		t.Errorf("Expected 1 migrated dir, got %d", len(result.Migrated))
	}

	// Check content moved
	newPath := filepath.Join(tmpDir, ".kdse", "foundation", "test.md")
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("Content not migrated to %s", newPath)
	}

	// Check old location gone
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Error("Old directory should be removed after migration")
	}
}

func TestGetPaths(t *testing.T) {
	ws := New("/repo")
	paths := ws.GetPaths()

	if paths.Root != "/repo/.kdse" {
		t.Errorf("Expected root /repo/.kdse, got %s", paths.Root)
	}
	if paths.Foundation != "/repo/.kdse/foundation" {
		t.Errorf("Expected foundation /repo/.kdse/foundation, got %s", paths.Foundation)
	}
	if paths.Knowledge != "/repo/.kdse/knowledge" {
		t.Errorf("Expected knowledge /repo/.kdse/knowledge, got %s", paths.Knowledge)
	}
}

func TestResolvePath(t *testing.T) {
	ws := New("/repo")

	tests := []struct {
		input    string
		expected string
	}{
		{".kdse/foundation", "/repo/.kdse/foundation"},
		{"foundation", "/repo/.kdse/foundation"},
		{"knowledge", "/repo/.kdse/knowledge"},
		{".kdse/reports/audit.md", "/repo/.kdse/reports/audit.md"},
	}

	for _, tt := range tests {
		got := ws.ResolvePath(tt.input)
		if got != tt.expected {
			t.Errorf("ResolvePath(%s): expected %s, got %s", tt.input, tt.expected, got)
		}
	}
}
