package orchestration

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWorkspaceResolver tests the workspace resolver
func TestWorkspaceResolver(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "kdse-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := DefaultEngineConfig()
	resolver := NewWorkspaceResolverWithPath(tempDir, config)

	// Test workspace type detection for empty directory
	workspace, err := resolver.ResolveWorkspace(tempDir)
	if err != nil {
		t.Fatalf("Failed to resolve workspace: %v", err)
	}

	if workspace.WorkspaceType != WorkspaceTypeUnknown {
		t.Errorf("Expected WorkspaceTypeUnknown for empty dir, got %s", workspace.WorkspaceType)
	}

	// Test creating KDSE directory
	kdsePath := workspace.KDSEPath
	if kdsePath != filepath.Join(tempDir, ".kdse") {
		t.Errorf("Expected KDSE path %s, got %s", filepath.Join(tempDir, ".kdse"), kdsePath)
	}

	// Test temp workspace creation
	tempWorkspace, err := resolver.ResolveTemporaryWorkspace("test-project")
	if err != nil {
		t.Fatalf("Failed to create temp workspace: %v", err)
	}

	if tempWorkspace.WorkspaceType != WorkspaceTypeTemporary {
		t.Errorf("Expected WorkspaceTypeTemporary, got %s", tempWorkspace.WorkspaceType)
	}

	expectedPath := filepath.Join(tempDir, config.TempWorkspaceBase, ".kdse", "test-project")
	if tempWorkspace.ResolvedPath != expectedPath {
		t.Errorf("Expected temp path %s, got %s", expectedPath, tempWorkspace.ResolvedPath)
	}

	// Verify temp directory was created
	if _, err := os.Stat(tempWorkspace.ResolvedPath); os.IsNotExist(err) {
		t.Errorf("Temp directory was not created at %s", tempWorkspace.ResolvedPath)
	}
}

// TestConfidenceEvaluator tests confidence evaluation
func TestConfidenceEvaluator(t *testing.T) {
	config := DefaultEngineConfig()
	evaluator := NewConfidenceEvaluator(config)

	// Create a test workspace
	tempDir, err := os.MkdirTemp("", "kdse-confidence-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create KDSE structure
	kdseDir := filepath.Join(tempDir, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("Failed to create KDSE dir: %v", err)
	}

	workspace := &WorkspaceInfo{
		ResolvedPath: tempDir,
		WorkspaceType: WorkspaceTypeRepository,
		KDSEPath: kdseDir,
	}

	// Test confidence evaluation with empty workspace
	conf, err := evaluator.EvaluateConfidence(workspace)
	if err != nil {
		t.Fatalf("Failed to evaluate confidence: %v", err)
	}

	if conf.Foundation >= config.FoundationThreshold {
		t.Errorf("Expected Foundation confidence < threshold for empty workspace, got %.2f", conf.Foundation)
	}

	if conf.MeetsThreshold {
		t.Errorf("Expected MeetsThreshold=false for empty workspace")
	}

	// Create some foundation documents
	foundationDir := filepath.Join(kdseDir, "foundation")
	if err := os.MkdirAll(foundationDir, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}

	// Create README
	os.WriteFile(filepath.Join(foundationDir, "README.md"), []byte("# Foundation\n"), 0644)
	os.WriteFile(filepath.Join(foundationDir, "004-engineering-model.md"), []byte("# Engineering Model\n"), 0644)

	// Re-evaluate
	conf, err = evaluator.EvaluateConfidence(workspace)
	if err != nil {
		t.Fatalf("Failed to re-evaluate confidence: %v", err)
	}

	// Should have improved but still not meet threshold
	if conf.Foundation <= 0 {
		t.Errorf("Expected Foundation confidence > 0 after adding docs, got %.2f", conf.Foundation)
	}
}

// TestEvidenceEvaluator tests evidence evaluation
func TestEvidenceEvaluator(t *testing.T) {
	config := DefaultEngineConfig()
	evaluator := NewEvidenceEvaluator(config)

	// Create a test workspace
	tempDir, err := os.MkdirTemp("", "kdse-evidence-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	workspace := &WorkspaceInfo{
		ResolvedPath: tempDir,
		WorkspaceType: WorkspaceTypeRepository,
		KDSEPath: filepath.Join(tempDir, ".kdse"),
	}

	// Test evidence evaluation for Resolve phase
	state, err := evaluator.EvaluateEvidence(workspace, PhaseResolve)
	if err != nil {
		t.Fatalf("Failed to evaluate evidence: %v", err)
	}

	if state.TotalRequired != 1 {
		t.Errorf("Expected 1 required evidence for Resolve phase, got %d", state.TotalRequired)
	}

	// Test evidence evaluation for Foundation phase
	state, err = evaluator.EvaluateEvidence(workspace, PhaseFoundation)
	if err != nil {
		t.Fatalf("Failed to evaluate evidence for Foundation: %v", err)
	}

	if state.TotalRequired != 6 {
		t.Errorf("Expected 6 required evidence for Foundation phase, got %d", state.TotalRequired)
	}

	// Create some evidence
	evidenceDir := filepath.Join(workspace.KDSEPath, "evidence", "screenshots")
	if err := os.MkdirAll(evidenceDir, 0755); err != nil {
		t.Fatalf("Failed to create evidence dir: %v", err)
	}

	// Test that evidence was found
	state, err = evaluator.EvaluateEvidence(workspace, PhaseCollect)
	if err != nil {
		t.Fatalf("Failed to evaluate evidence for Collect: %v", err)
	}

	// Evidence directory should satisfy some requirements
	if state.TotalPresent < 1 {
		t.Errorf("Expected at least 1 present evidence after creating directory, got %d", state.TotalPresent)
	}
}

// TestEngineInitialization tests engine initialization
func TestEngineInitialization(t *testing.T) {
	config := DefaultEngineConfig()
	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create a test workspace
	tempDir, err := os.MkdirTemp("", "kdse-engine-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize engine
	if err := engine.Initialize(tempDir); err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	// Verify initial state
	state := engine.GetState()
	if state == nil {
		t.Fatal("Expected non-nil state after initialization")
	}

	if state.CurrentPhase != PhaseResolve {
		t.Errorf("Expected initial phase %s, got %s", PhaseResolve, state.CurrentPhase)
	}

	if !engine.IsSessionActive() {
		t.Error("Expected session to be active after initialization")
	}
}

// TestEngineCycle tests a single orchestration cycle
func TestEngineCycle(t *testing.T) {
	config := &EngineConfig{
		FoundationThreshold: 0.3, // Lower threshold for testing
		EvidenceThreshold:   0.2,
		MaxCycles:          10,
		TempWorkspaceBase:  "temp",
		EnableMigration:    false,
	}

	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create a test workspace
	tempDir, err := os.MkdirTemp("", "kdse-cycle-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create minimal KDSE structure
	kdseDir := filepath.Join(tempDir, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("Failed to create KDSE dir: %v", err)
	}

	// Initialize engine
	if err := engine.Initialize(tempDir); err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	// Execute one cycle
	result, err := engine.ExecuteCycle()
	if err != nil {
		t.Fatalf("Failed to execute cycle: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected cycle to succeed, got error: %s", result.Error)
	}

	if result.CycleNumber != 1 {
		t.Errorf("Expected cycle number 1, got %d", result.CycleNumber)
	}

	// Verify state was updated
	state := engine.GetState()
	if state.Metrics.CycleCount != 1 {
		t.Errorf("Expected Metrics.CycleCount=1, got %d", state.Metrics.CycleCount)
	}
}

// TestFoundationThresholdBlocking tests that implementation is blocked until Foundation threshold met
func TestFoundationThresholdBlocking(t *testing.T) {
	config := &EngineConfig{
		FoundationThreshold: 0.9, // High threshold
		EvidenceThreshold:   0.2,
		MaxCycles:          10,
		TempWorkspaceBase:  "temp",
		EnableMigration:    false,
	}

	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create a test workspace without Foundation docs
	tempDir, err := os.MkdirTemp("", "kdse-blocking-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create empty KDSE structure
	kdseDir := filepath.Join(tempDir, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		t.Fatalf("Failed to create KDSE dir: %v", err)
	}

	// Initialize engine
	if err := engine.Initialize(tempDir); err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	// Verify implementation is blocked
	if engine.CanImplement() {
		t.Error("Expected CanImplement=false with empty workspace")
	}

	// Verify state shows blocking
	state := engine.GetState()
	if state.Blocked.Blocked {
		t.Log("State is blocked as expected (Foundation threshold not met)")
	}

	// Even if we manually set phase to Implement, it should still block
	state.CurrentPhase = PhaseImplement
	state.Confidence.MeetsThreshold = false

	decision := engine.decideNextPhase()
	if decision.ShouldExecute {
		t.Error("Should not allow execution when Foundation threshold not met")
	}
}

// TestWorkspaceMigration tests temporary workspace migration
func TestWorkspaceMigration(t *testing.T) {
	config := &EngineConfig{
		FoundationThreshold: 0.5,
		EvidenceThreshold:  0.2,
		MaxCycles:          10,
		TempWorkspaceBase:  "temp",
		EnableMigration:    true,
	}

	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create temp directory structure
	tempDir, err := os.MkdirTemp("", "kdse-migrate-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create project directory
	projectDir := filepath.Join(tempDir, "my-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	// Initialize engine in temp workspace
	tempKDSE := filepath.Join(tempDir, "temp", ".kdse", "my-project")
	if err := os.MkdirAll(tempKDSE, 0755); err != nil {
		t.Fatalf("Failed to create temp KDSE: %v", err)
	}

	// Write something to temp workspace
	os.WriteFile(filepath.Join(tempKDSE, "context.json"), []byte(`{"project":"test"}`), 0644)

	// Initialize engine state
	engine.state = &OrchestrationState{
		Workspace: WorkspaceInfo{
			ResolvedPath: tempKDSE,
			WorkspaceType: WorkspaceTypeTemporary,
			TempPath:      tempKDSE,
			KDSEPath:      tempKDSE,
		},
	}

	// Migrate to project
	if err := engine.MigrateToProject(projectDir); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Verify migration
	projectKDSE := filepath.Join(projectDir, ".kdse")
	if _, err := os.Stat(projectKDSE); os.IsNotExist(err) {
		t.Error("Expected .kdse to be migrated to project directory")
	}

	// Verify context file exists in project
	contextPath := filepath.Join(projectKDSE, "context.json")
	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		t.Error("Expected context.json to be migrated")
	}
}

// TestPhaseTransitions tests phase transition logic
func TestPhaseTransitions(t *testing.T) {
	config := DefaultEngineConfig()
	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create test workspace
	tempDir, err := os.MkdirTemp("", "kdse-phases-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := engine.Initialize(tempDir); err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	// Test Resolve phase
	state := engine.GetState()
	if state.CurrentPhase != PhaseResolve {
		t.Errorf("Expected initial phase %s, got %s", PhaseResolve, state.CurrentPhase)
	}

	// Execute cycle - should move to Assess
	result, err := engine.ExecuteCycle()
	if err != nil {
		t.Fatalf("Failed to execute cycle: %v", err)
	}

	// Decision should allow progression
	if result.Decision == nil {
		t.Fatal("Expected non-nil decision")
	}
}

// TestNoHardcodedPaths tests that paths don't contain hardcoded values
func TestNoHardcodedPaths(t *testing.T) {
	config := DefaultEngineConfig()
	resolver := NewWorkspaceResolverWithPath("/tmp", config)

	// Test with various paths
	testPaths := []string{
		"/home/user/projects/myapp",
		"/var/www/html",
		"/tmp/test",
		"relative/path",
	}

	for _, path := range testPaths {
		workspace, err := resolver.ResolveWorkspace(path)
		if err != nil {
			t.Errorf("Failed to resolve %s: %v", path, err)
			continue
		}

		// Verify no hardcoded /app or /workspace in resolved paths
		resolvedPaths := []string{
			workspace.ResolvedPath,
			workspace.KDSEPath,
			workspace.RepositoryPath,
		}

		for _, rp := range resolvedPaths {
			if rp == "" {
				continue
			}
			// These should not appear in resolved paths
			if containsHardcodedPath(rp) {
				t.Errorf("Found hardcoded path in %s: %s", path, rp)
			}
		}
	}
}

// containsHardcodedPath checks for hardcoded path patterns
func containsHardcodedPath(path string) bool {
	// Check for common hardcoded patterns
	hardcoded := []string{"/app", "/workspace"}
	for _, h := range hardcoded {
		if len(path) >= len(h) && path[:len(h)] == h && len(path) != len(h) {
			// Only flag if it's a prefix match followed by more path
			if path[len(h)] == '/' || path[len(h)] == '\\' {
				return true
			}
		}
	}
	return false
}
