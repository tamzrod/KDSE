package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestEngine_VerifyWorkspace_NoRuntime(t *testing.T) {
	// Create temp directory without .kdse
	dir := t.TempDir()

	engine := NewEngine(dir)
	ctx := context.Background()

	// Verify should fail without .kdse
	result, err := engine.VerifyWorkspace(ctx)
	if err == nil {
		t.Error("expected error when .kdse doesn't exist")
	}
	if result != nil && result.Valid {
		t.Error("expected verification to fail")
	}
}

func TestEngine_InitializeWorkspace(t *testing.T) {
	// Create temp directory
	dir := t.TempDir()

	engine := NewEngine(dir)
	ctx := context.Background()

	opts := InitOptions{
		Path:     dir,
		Type:     RuntimeTypeCLI,
		Version:  "1.0.0",
		Template: "default",
	}

	// Initialize workspace
	ws, err := engine.InitializeWorkspace(ctx, opts)
	if err != nil {
		t.Fatalf("initialization failed: %v", err)
	}

	// Verify workspace was created
	if ws == nil {
		t.Fatal("workspace is nil")
	}

	// Verify .kdse exists
	kdsePath := filepath.Join(dir, ".kdse")
	if _, err := os.Stat(kdsePath); os.IsNotExist(err) {
		t.Error(".kdse directory was not created")
	}

	// Verify runtime files exist
	files := []string{
		"runtime.yaml",
		"workspace.yaml",
		"methodology.yaml",
		"phase.yaml",
		"session.yaml",
		"metadata.yaml",
	}

	for _, file := range files {
		filePath := filepath.Join(kdsePath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("required file %s was not created", file)
		}
	}

	// Verify phase directories exist
	dirs := []string{
		"knowledge",
		"architecture",
		"implementation",
		"verification",
		"reports",
	}

	for _, d := range dirs {
		dirPath := filepath.Join(kdsePath, d)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("required directory %s was not created", d)
		}
	}
}

func TestEngine_VerifyWorkspace_WithRuntime(t *testing.T) {
	// Create temp directory with .kdse
	dir := t.TempDir()
	kdseDir := filepath.Join(dir, ".kdse")

	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse: %v", err)
	}

	// Create minimal runtime files
	runtimeYaml := `runtime:
  type: cli
  version: 1.0.0
`
	phaseYaml := `current: initialization
`

	if err := os.WriteFile(filepath.Join(kdseDir, "runtime.yaml"), []byte(runtimeYaml), 0644); err != nil {
		t.Fatalf("failed to write runtime.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(kdseDir, "phase.yaml"), []byte(phaseYaml), 0644); err != nil {
		t.Fatalf("failed to write phase.yaml: %v", err)
	}

	engine := NewEngine(dir)
	ctx := context.Background()

	result, err := engine.VerifyWorkspace(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("verification result is nil")
	}
	if !result.Valid {
		t.Errorf("verification should pass: %v", result.Errors)
	}
}

func TestEngine_GetPhase(t *testing.T) {
	// Create temp directory with .kdse
	dir := t.TempDir()
	kdseDir := filepath.Join(dir, ".kdse")

	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse: %v", err)
	}

	phaseYaml := `current: knowledge
`
	if err := os.WriteFile(filepath.Join(kdseDir, "phase.yaml"), []byte(phaseYaml), 0644); err != nil {
		t.Fatalf("failed to write phase.yaml: %v", err)
	}

	engine := NewEngine(dir)
	ctx := context.Background()

	phase, err := engine.GetPhase(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if phase != PhaseKnowledge {
		t.Errorf("expected phase knowledge, got %s", phase)
	}
}

func TestEngine_AdvancePhase_InvalidTransition(t *testing.T) {
	// Create temp directory with .kdse
	dir := t.TempDir()
	kdseDir := filepath.Join(dir, ".kdse")

	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse: %v", err)
	}

	// Set current phase to initialization
	phaseYaml := `current: initialization
`
	if err := os.WriteFile(filepath.Join(kdseDir, "phase.yaml"), []byte(phaseYaml), 0644); err != nil {
		t.Fatalf("failed to write phase.yaml: %v", err)
	}

	engine := NewEngine(dir)
	ctx := context.Background()

	// Try to advance to implementation (skipping knowledge)
	_, err := engine.AdvancePhase(ctx, PhaseImplementation)
	if err == nil {
		t.Error("expected error for invalid phase transition")
	}
}

func TestEngine_AdvancePhase_ValidTransition(t *testing.T) {
	// Create temp directory with .kdse
	dir := t.TempDir()
	kdseDir := filepath.Join(dir, ".kdse")

	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse: %v", err)
	}

	// Create required knowledge artifacts
	knowledgeDir := filepath.Join(kdseDir, "knowledge")
	if err := os.MkdirAll(knowledgeDir, 0755); err != nil {
		t.Fatalf("failed to create knowledge dir: %v", err)
	}

	// Create required files
	requiredFiles := []string{
		"requirements.md",
		"stakeholders.md",
		"constraints.md",
	}
	for _, file := range requiredFiles {
		if err := os.WriteFile(filepath.Join(knowledgeDir, file), []byte("# Test"), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", file, err)
		}
	}

	// Set current phase to initialization
	phaseYaml := `current: initialization
`
	if err := os.WriteFile(filepath.Join(kdseDir, "phase.yaml"), []byte(phaseYaml), 0644); err != nil {
		t.Fatalf("failed to write phase.yaml: %v", err)
	}

	engine := NewEngine(dir)
	ctx := context.Background()

	// Advance to knowledge
	transition, err := engine.AdvancePhase(ctx, PhaseKnowledge)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if transition == nil {
		t.Fatal("transition is nil")
	}
	if transition.From != PhaseInitialization {
		t.Errorf("expected from initialization, got %s", transition.From)
	}
	if transition.To != PhaseKnowledge {
		t.Errorf("expected to knowledge, got %s", transition.To)
	}
}

func TestEngine_VerifyArtifacts_Incomplete(t *testing.T) {
	// Create temp directory with .kdse
	dir := t.TempDir()
	kdseDir := filepath.Join(dir, ".kdse")

	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("failed to create .kdse: %v", err)
	}

	// Create knowledge directory but no artifacts
	knowledgeDir := filepath.Join(kdseDir, "knowledge")
	if err := os.MkdirAll(knowledgeDir, 0755); err != nil {
		t.Fatalf("failed to create knowledge dir: %v", err)
	}

	// Set current phase to knowledge
	phaseYaml := `current: knowledge
`
	if err := os.WriteFile(filepath.Join(kdseDir, "phase.yaml"), []byte(phaseYaml), 0644); err != nil {
		t.Fatalf("failed to write phase.yaml: %v", err)
	}

	engine := NewEngine(dir)
	ctx := context.Background()

	result, err := engine.VerifyArtifacts(ctx, PhaseKnowledge)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("verification result is nil")
	}
	if result.Valid {
		t.Error("verification should fail for incomplete phase")
	}
	if len(result.Missing) == 0 {
		t.Error("expected missing artifacts")
	}
}

func TestLifecycle_IsValidTransition(t *testing.T) {
	lc := lifecycle.NewLifecycle()

	tests := []struct {
		from Phase
		to   Phase
		want bool
	}{
		{PhaseInitialization, PhaseKnowledge, true},
		{PhaseKnowledge, PhaseArchitecture, true},
		{PhaseArchitecture, PhaseImplementation, true},
		{PhaseImplementation, PhaseVerification, true},
		{PhaseVerification, PhaseReports, true},
		{PhaseInitialization, PhaseArchitecture, false}, // Skip
		{PhaseKnowledge, PhaseImplementation, false},     // Skip
		{PhaseReports, PhaseInitialization, false},       // Regression
		{PhaseInitialization, PhaseImplementation, false}, // Skip multiple
	}

	for _, tt := range tests {
		got := lc.IsValidTransition(tt.from, tt.to)
		if got != tt.want {
			t.Errorf("IsValidTransition(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestLifecycle_GetNextPhase(t *testing.T) {
	lc := lifecycle.NewLifecycle()

	tests := []struct {
		current Phase
		want    Phase
		wantErr bool
	}{
		{PhaseInitialization, PhaseKnowledge, false},
		{PhaseKnowledge, PhaseArchitecture, false},
		{PhaseArchitecture, PhaseImplementation, false},
		{PhaseImplementation, PhaseVerification, false},
		{PhaseVerification, PhaseReports, false},
		{PhaseReports, "", true}, // Terminal phase
	}

	for _, tt := range tests {
		got, err := lc.GetNextPhase(tt.current)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetNextPhase(%s) error = %v, wantErr %v", tt.current, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("GetNextPhase(%s) = %v, want %v", tt.current, got, tt.want)
		}
	}
}

func TestLifecycle_IsTerminalPhase(t *testing.T) {
	lc := lifecycle.NewLifecycle()

	if !lc.IsTerminalPhase(PhaseReports) {
		t.Error("expected reports to be terminal phase")
	}

	if lc.IsTerminalPhase(PhaseInitialization) {
		t.Error("expected initialization to not be terminal phase")
	}

	if lc.IsTerminalPhase(PhaseImplementation) {
		t.Error("expected implementation to not be terminal phase")
	}
}

func TestRuntimeType_Constants(t *testing.T) {
	if RuntimeTypeCLI != "cli" {
		t.Errorf("RuntimeTypeCLI = %s, want cli", RuntimeTypeCLI)
	}
	if RuntimeTypeMCP != "mcp" {
		t.Errorf("RuntimeTypeMCP = %s, want mcp", RuntimeTypeMCP)
	}
}

func TestPhase_Constants(t *testing.T) {
	phases := []Phase{
		PhaseInitialization,
		PhaseKnowledge,
		PhaseArchitecture,
		PhaseImplementation,
		PhaseVerification,
		PhaseReports,
	}

	expected := []string{
		"initialization",
		"knowledge",
		"architecture",
		"implementation",
		"verification",
		"reports",
	}

	for i, phase := range phases {
		if string(phase) != expected[i] {
			t.Errorf("Phase[%d] = %s, want %s", i, phase, expected[i])
		}
	}
}
