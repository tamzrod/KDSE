package detection

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEnvironmentTypeString tests the String method.
func TestEnvironmentTypeString(t *testing.T) {
	tests := []struct {
		env      EnvironmentType
		expected string
	}{
		{EnvironmentBlankWorkspace, "BLANK_WORKSPACE"},
		{EnvironmentSoftwareProject, "SOFTWARE_PROJECT"},
		{EnvironmentKDSERuntime, "KDSE_RUNTIME"},
		{EnvironmentKDSEProject, "KDSE_PROJECT"},
		{EnvironmentUnknown, "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.env.String(); got != tt.expected {
			t.Errorf("EnvironmentType.String() = %v, want %v", got, tt.expected)
		}
	}
}

// TestEnvironmentTypeIsValid tests the IsValid method.
func TestEnvironmentTypeIsValid(t *testing.T) {
	tests := []struct {
		env      EnvironmentType
		expected bool
	}{
		{EnvironmentBlankWorkspace, true},
		{EnvironmentSoftwareProject, true},
		{EnvironmentKDSERuntime, true},
		{EnvironmentKDSEProject, true},
		{EnvironmentUnknown, false},
	}

	for _, tt := range tests {
		if got := tt.env.IsValid(); got != tt.expected {
			t.Errorf("EnvironmentType.IsValid() = %v, want %v", got, tt.expected)
		}
	}
}

// TestDetectionKDSEProject tests detection of a KDSE-initialized project.
func TestDetectionKDSEProject(t *testing.T) {
	// Create a temporary directory simulating a KDSE project
	tmpDir := t.TempDir()

	// Create .kdse/ directory
	kdseDir := filepath.Join(tmpDir, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse dir: %v", err)
	}

	// Create manifest.json
	manifestPath := filepath.Join(kdseDir, "manifest.json")
	manifestContent := `{
		"schema": "https://kdse.dev/schemas/manifest/v1.0",
		"version": "1.0.0",
		"generated": "2026-07-11T00:00:00Z",
		"status": "ACTIVE"
	}`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("failed to create manifest.json: %v", err)
	}

	// Create some project files
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/project\n"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Detect
	detector := NewEnvironmentDetector(tmpDir)
	result := detector.Detect()

	if result.Environment != EnvironmentKDSEProject {
		t.Errorf("Expected EnvironmentKDSEProject, got %v", result.Environment)
	}
}

// TestDetectionSoftwareProject tests detection of a normal software project.
func TestDetectionSoftwareProject(t *testing.T) {
	// Create a temporary directory simulating a software project
	tmpDir := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module example.com/project\n"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Create README.md
	if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# My Project\n"), 0644); err != nil {
		t.Fatalf("failed to create README.md: %v", err)
	}

	// Detect
	detector := NewEnvironmentDetector(tmpDir)
	result := detector.Detect()

	if result.Environment != EnvironmentSoftwareProject {
		t.Errorf("Expected EnvironmentSoftwareProject, got %v", result.Environment)
	}
}

// TestDetectionBlankWorkspace tests detection of a blank workspace.
func TestDetectionBlankWorkspace(t *testing.T) {
	// Create an empty temporary directory
	tmpDir := t.TempDir()

	// Detect
	detector := NewEnvironmentDetector(tmpDir)
	result := detector.Detect()

	if result.Environment != EnvironmentBlankWorkspace {
		t.Errorf("Expected EnvironmentBlankWorkspace, got %v", result.Environment)
	}
}

// TestDetectionKDSERuntime tests detection of the KDSE runtime itself.
func TestDetectionKDSERuntime(t *testing.T) {
	// Create a temporary directory simulating the KDSE runtime
	tmpDir := t.TempDir()

	// Create go.mod with KDSE module
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module github.com/kdse/runtime\n"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Create internal/runtime/ directory
	if err := os.MkdirAll(filepath.Join(tmpDir, "internal", "runtime"), 0755); err != nil {
		t.Fatalf("failed to create internal/runtime: %v", err)
	}

	// Create internal/detection/ directory
	if err := os.MkdirAll(filepath.Join(tmpDir, "internal", "detection"), 0755); err != nil {
		t.Fatalf("failed to create internal/detection: %v", err)
	}

	// Create templates/ directory
	if err := os.MkdirAll(filepath.Join(tmpDir, "templates"), 0755); err != nil {
		t.Fatalf("failed to create templates: %v", err)
	}

	// Create docs/ directory
	if err := os.MkdirAll(filepath.Join(tmpDir, "docs"), 0755); err != nil {
		t.Fatalf("failed to create docs: %v", err)
	}

	// Detect
	detector := NewEnvironmentDetector(tmpDir)
	result := detector.Detect()

	if result.Environment != EnvironmentKDSERuntime {
		t.Errorf("Expected EnvironmentKDSERuntime, got %v", result.Environment)
	}

	// Verify confidence is high
	if result.Confidence < 0.8 {
		t.Errorf("Expected confidence >= 0.8, got %v", result.Confidence)
	}
}

// TestDetectionWithGitRepo tests detection with Git repository present.
func TestDetectionWithGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .git directory
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("failed to create .git: %v", err)
	}

	// Create a package.json
	pkgPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(pkgPath, []byte(`{"name": "my-project"}`), 0644); err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	// Detect
	detector := NewEnvironmentDetector(tmpDir)
	result := detector.Detect()

	if result.Environment != EnvironmentSoftwareProject {
		t.Errorf("Expected EnvironmentSoftwareProject, got %v", result.Environment)
	}

	if !result.Evidence.HasGitRepo {
		t.Error("Expected HasGitRepo to be true")
	}
}

// TestDetectionAmbiguousCase tests handling of ambiguous evidence.
func TestDetectionAmbiguousCase(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a KDSE project with KDSE markers - should still be KDSE_PROJECT
	// because KDSE_PROJECT takes precedence over KDSE_RUNTIME when manifest is valid
	kdseDir := filepath.Join(tmpDir, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse dir: %v", err)
	}

	manifestPath := filepath.Join(kdseDir, "manifest.json")
	manifestContent := `{
		"schema": "https://kdse.dev/schemas/manifest/v1.0",
		"version": "1.0.0",
		"status": "ACTIVE"
	}`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("failed to create manifest.json: %v", err)
	}

	// Detect
	detector := NewEnvironmentDetector(tmpDir)
	result := detector.Detect()

	// Should be KDSE_PROJECT because it has a valid manifest
	if result.Environment != EnvironmentKDSEProject {
		t.Errorf("Expected EnvironmentKDSEProject, got %v", result.Environment)
	}
}

// TestDetectEnvironmentFunction tests the convenience function.
func TestDetectEnvironmentFunction(t *testing.T) {
	tmpDir := t.TempDir()

	env := DetectEnvironment(tmpDir)

	if env != EnvironmentBlankWorkspace {
		t.Errorf("Expected EnvironmentBlankWorkspace, got %v", env)
	}
}

// TestEvidenceResultWarnings tests warning handling.
func TestEvidenceResultWarnings(t *testing.T) {
	result := &EvidenceResult{
		Warnings: []string{"warning 1", "warning 2"},
		Errors:   []string{},
	}

	if !result.HasWarnings() {
		t.Error("Expected HasWarnings() to return true")
	}

	if result.HasErrors() {
		t.Error("Expected HasErrors() to return false")
	}
}

// TestEvidenceResultErrors tests error handling.
func TestEvidenceResultErrors(t *testing.T) {
	result := &EvidenceResult{
		Warnings: []string{},
		Errors:   []string{"error 1"},
	}

	if result.HasWarnings() {
		t.Error("Expected HasWarnings() to return false")
	}

	if !result.HasErrors() {
		t.Error("Expected HasErrors() to return true")
	}
}

// TestClassificationRules tests that all classification rules are defined.
func TestClassificationRules(t *testing.T) {
	rules := GetClassificationRules()

	if len(rules) != 4 {
		t.Errorf("Expected 4 classification rules, got %d", len(rules))
	}

	// Verify rule names
	expectedNames := []string{
		"KDSE Runtime Detection",
		"KDSE Project Detection",
		"Software Project Detection",
		"Blank Workspace Detection",
	}

	for i, rule := range rules {
		if rule.Name != expectedNames[i] {
			t.Errorf("Rule %d: expected name %q, got %q", i, expectedNames[i], rule.Name)
		}

		if rule.Description == "" {
			t.Errorf("Rule %d (%s) has empty description", i, rule.Name)
		}

		if len(rule.RequiredEvidence) == 0 {
			t.Errorf("Rule %d (%s) has no required evidence", i, rule.Name)
		}
	}
}

// TestRulePriority tests that rules are in correct priority order.
func TestRulePriority(t *testing.T) {
	rules := GetClassificationRules()

	// Rules should be in priority order (1, 2, 3, 4)
	for i, rule := range rules {
		if rule.Priority != i+1 {
			t.Errorf("Rule %d: expected priority %d, got %d", i, i+1, rule.Priority)
		}
	}
}

// TestProjectIndicators tests various project type detection.
func TestProjectIndicators(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]bool
		expected string
	}{
		{
			name:     "Go project",
			files:    map[string]bool{"go.mod": true, "main.go": true},
			expected: "go",
		},
		{
			name:     "Node project",
			files:    map[string]bool{"package.json": true, "index.js": true},
			expected: "node",
		},
		{
			name:     "Python project",
			files:    map[string]bool{"pyproject.toml": true},
			expected: "python",
		},
		{
			name:     "Rust project",
			files:    map[string]bool{"Cargo.toml": true},
			expected: "rust",
		},
		{
			name:     "Java project",
			files:    map[string]bool{"pom.xml": true},
			expected: "java",
		},
		{
			name:     "Makefile project",
			files:    map[string]bool{"Makefile": true},
			expected: "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evidence := Evidence{}

			switch tt.expected {
			case "go":
				evidence.HasGoMod = true
			case "node":
				evidence.HasPackageJSON = true
			case "python":
				evidence.HasPyProject = true
			case "rust":
				evidence.HasCargoToml = true
			case "java":
				evidence.HasPomXML = true
			case "c":
				evidence.HasMakefile = true
			}

			detector := &EnvDetector{}
			projectType := detector.detectProjectType(evidence)

			if projectType != tt.expected {
				t.Errorf("Expected project type %s, got %s", tt.expected, projectType)
			}
		})
	}
}

// TestParseGoModuleName tests parsing of Go module names.
func TestParseGoModuleName(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"module github.com/kdse/runtime\n", "github.com/kdse/runtime"},
		{"module example.com/my-project\n", "example.com/my-project"},
		{"module local-module\n", "local-module"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			goModPath := filepath.Join(tmpDir, "go.mod")
			if err := os.WriteFile(goModPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create go.mod: %v", err)
			}

			detector := &EnvDetector{}
			moduleName := detector.parseGoModuleName(goModPath)

			if moduleName != tt.expected {
				t.Errorf("Expected module %q, got %q", tt.expected, moduleName)
			}
		})
	}
}

// TestParseManifestJSON tests JSON manifest parsing.
func TestParseManifestJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid manifest
	manifestContent := `{
		"schema": "https://kdse.dev/schemas/manifest/v1.0",
		"version": "1.0.0",
		"generated": "2026-07-11T00:00:00Z",
		"status": "ACTIVE",
		"runtime": {
			"version": "2.0"
		},
		"installation": {
			"source_repository": "https://github.com/kdse/project"
		}
	}`

	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("failed to create manifest.json: %v", err)
	}

	detector := &EnvDetector{}
	info, valid := detector.parseManifestJSON(manifestPath)

	if !valid {
		t.Error("Expected valid manifest")
	}

	if info.Schema != "https://kdse.dev/schemas/manifest/v1.0" {
		t.Errorf("Expected schema %q, got %q", "https://kdse.dev/schemas/manifest/v1.0", info.Schema)
	}

	if info.RuntimeVersion != "2.0" {
		t.Errorf("Expected runtime version %q, got %q", "2.0", info.RuntimeVersion)
	}
}

// TestInvalidManifestJSON tests handling of invalid JSON.
func TestInvalidManifestJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an invalid manifest
	manifestContent := `{invalid json}`
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("failed to create manifest.json: %v", err)
	}

	detector := &EnvDetector{}
	_, valid := detector.parseManifestJSON(manifestPath)

	if valid {
		t.Error("Expected invalid manifest to return false")
	}
}

// TestDefaultDetectionOptions tests default options.
func TestDefaultDetectionOptions(t *testing.T) {
	opts := DefaultDetectionOptions()

	if opts.AllowAmbiguous {
		t.Error("Expected AllowAmbiguous to be false")
	}

	if opts.RequireGitRepo {
		t.Error("Expected RequireGitRepo to be false")
	}

	if !opts.StrictManifestValidation {
		t.Error("Expected StrictManifestValidation to be true")
	}
}

// TestConfidenceCalculation tests confidence calculation for each environment.
func TestConfidenceCalculation(t *testing.T) {
	tests := []struct {
		name     string
		evidence Evidence
		env      EnvironmentType
		minConf  float64
	}{
		{
			name: "KDSE Runtime with all markers",
			evidence: Evidence{
				HasGoModKDSE:    true,
				HasRuntimeDir:   true,
				HasDetectionDir: true,
				HasTemplatesDir: true,
				HasDocsDir:      true,
			},
			env:     EnvironmentKDSERuntime,
			minConf: 0.8,
		},
		{
			name: "KDSE Project with valid manifest",
			evidence: Evidence{
				HasKDSEProject: true,
				HasManifest:    true,
				ManifestValid:  true,
			},
			env:     EnvironmentKDSEProject,
			minConf: 0.9,
		},
		{
			name: "Software Project with multiple indicators",
			evidence: Evidence{
				HasGoMod:       true,
				HasPackageJSON: true,
				HasPyProject:   true,
			},
			env:     EnvironmentSoftwareProject,
			minConf: 0.7,
		},
		{
			name:     "Blank Workspace",
			evidence: Evidence{},
			env:      EnvironmentBlankWorkspace,
			minConf:  0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &EnvDetector{}
			confidence := detector.calculateConfidence(&EvidenceResult{
				Environment: tt.env,
				Evidence:    tt.evidence,
			})

			if confidence < tt.minConf {
				t.Errorf("Expected confidence >= %v, got %v", tt.minConf, confidence)
			}
		})
	}
}
