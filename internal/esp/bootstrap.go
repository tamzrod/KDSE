// Package esp implements the KDSE Engineering Session Protocol (ESP).
package esp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BootstrapStep represents a step in the bootstrap process
type BootstrapStep string

const (
	BootstrapStepDiscoverWorkspace  BootstrapStep = "DISCOVER_WORKSPACE"
	BootstrapStepVerifyRuntime     BootstrapStep = "VERIFY_RUNTIME"
	BootstrapStepLoadManifest      BootstrapStep = "LOAD_MANIFEST"
	BootstrapStepDetermineContext  BootstrapStep = "DETERMINE_CONTEXT"
	BootstrapStepResolveConflicts  BootstrapStep = "RESOLVE_CONFLICTS"
	BootstrapStepAssembleContext   BootstrapStep = "ASSEMBLE_CONTEXT"
	BootstrapStepComputeChecksum   BootstrapStep = "COMPUTE_CHECKSUM"
	BootstrapStepRecordEvidence    BootstrapStep = "RECORD_EVIDENCE"
)

// BootstrapResult represents the result of a bootstrap operation
type BootstrapResult struct {
	Success         bool                 `json:"success"`
	Context         *Context             `json:"context,omitempty"`
	BootstrapSteps  []BootstrapStepResult `json:"bootstrap_steps"`
	Checksum        string               `json:"checksum"`
	CompletedAt     string               `json:"completed_at"`
	Duration        string               `json:"duration"`
	Errors          []string             `json:"errors,omitempty"`
	Warnings        []string             `json:"warnings,omitempty"`
}

// BootstrapStepResult represents the result of a single bootstrap step
type BootstrapStepResult struct {
	Step        BootstrapStep `json:"step"`
	Success     bool          `json:"success"`
	StartedAt   string        `json:"started_at"`
	CompletedAt string        `json:"completed_at"`
	Evidence    []string      `json:"evidence,omitempty"`
	Error       string        `json:"error,omitempty"`
}

// BootstrapProtocol implements the bootstrap process per Section 9
type BootstrapProtocol struct {
	repoPath      string
	workspacePath string
	manifestPath  string
	currentContext *Context
}

// NewBootstrapProtocol creates a new bootstrap protocol instance
func NewBootstrapProtocol(repoPath string) *BootstrapProtocol {
	return &BootstrapProtocol{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
		manifestPath:  filepath.Join(repoPath, ".kdse", "runtime.yaml"),
	}
}

// Execute performs the full bootstrap process as defined in Section 9
func (bp *BootstrapProtocol) Execute(sessionID string) (*BootstrapResult, error) {
	result := &BootstrapResult{
		BootstrapSteps: make([]BootstrapStepResult, 0),
	}

	startTime := time.Now()
	defer func() {
		result.CompletedAt = time.Now().Format(time.RFC3339)
		result.Duration = time.Since(startTime).String()
	}()

	// Execute each bootstrap step
	steps := []BootstrapStep{
		BootstrapStepDiscoverWorkspace,
		BootstrapStepVerifyRuntime,
		BootstrapStepLoadManifest,
		BootstrapStepDetermineContext,
		BootstrapStepResolveConflicts,
		BootstrapStepAssembleContext,
		BootstrapStepComputeChecksum,
		BootstrapStepRecordEvidence,
	}

	result.Success = true
	for _, step := range steps {
		stepResult := bp.executeStep(step, sessionID)
		result.BootstrapSteps = append(result.BootstrapSteps, stepResult)

		if !stepResult.Success {
			result.Success = false
			result.Errors = append(result.Errors, stepResult.Error)
			break
		}

		// Add evidence from successful step
		if bp.currentContext != nil {
			bp.currentContext.Evidence = append(bp.currentContext.Evidence, stepResult.Evidence...)
		}
	}

	if result.Success && bp.currentContext != nil {
		result.Context = bp.currentContext
		result.Checksum = bp.currentContext.Checksum
	}

	return result, nil
}

// executeStep executes a single bootstrap step
func (bp *BootstrapProtocol) executeStep(step BootstrapStep, sessionID string) BootstrapStepResult {
	result := BootstrapStepResult{
		Step:      step,
		StartedAt: time.Now().Format(time.RFC3339),
		Success:   true,
	}
	defer func() {
		result.CompletedAt = time.Now().Format(time.RFC3339)
	}()

	var err error
	var evidence []string

	switch step {
	case BootstrapStepDiscoverWorkspace:
		evidence, err = bp.discoverWorkspace()
	case BootstrapStepVerifyRuntime:
		evidence, err = bp.verifyRuntime()
	case BootstrapStepLoadManifest:
		evidence, err = bp.loadManifest()
	case BootstrapStepDetermineContext:
		evidence, err = bp.determineContext(sessionID)
	case BootstrapStepResolveConflicts:
		evidence, err = bp.resolveConflicts()
	case BootstrapStepAssembleContext:
		evidence, err = bp.assembleContext()
	case BootstrapStepComputeChecksum:
		evidence, err = bp.computeChecksum()
	case BootstrapStepRecordEvidence:
		evidence, err = bp.recordEvidence()
	default:
		err = fmt.Errorf("unknown bootstrap step: %s", step)
	}

	result.Evidence = evidence
	if err != nil {
		result.Success = false
		result.Error = err.Error()
	}

	return result
}

// discoverWorkspace locates and verifies the workspace exists
func (bp *BootstrapProtocol) discoverWorkspace() ([]string, error) {
	var evidence []string

	// Check if workspace directory exists
	info, err := os.Stat(bp.workspacePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewFailure(FailureWorkspaceNotFound, "workspace", "exists", "not found")
		}
		return nil, fmt.Errorf("failed to stat workspace: %w", err)
	}

	if !info.IsDir() {
		return nil, NewFailure(FailureWorkspaceNotFound, "workspace", "directory", "not a directory")
	}

	evidence = append(evidence, fmt.Sprintf("workspace_path:%s", bp.workspacePath))
	return evidence, nil
}

// verifyRuntime verifies the runtime files exist
func (bp *BootstrapProtocol) verifyRuntime() ([]string, error) {
	var evidence []string

	requiredFiles := []string{
		"manifest.yaml",
		"runtime.yaml",
		"session-state.yaml",
	}

	var missing []string
	for _, file := range requiredFiles {
		path := filepath.Join(bp.workspacePath, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, file)
		} else {
			evidence = append(evidence, fmt.Sprintf("runtime_file:%s", file))
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required runtime files: %s", strings.Join(missing, ", "))
	}

	return evidence, nil
}

// loadManifest loads and parses the runtime manifest
func (bp *BootstrapProtocol) loadManifest() ([]string, error) {
	var evidence []string

	// Load runtime.yaml
	runtimePath := filepath.Join(bp.workspacePath, "runtime.yaml")
	data, err := os.ReadFile(runtimePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read runtime manifest: %w", err)
	}

	// Verify it's valid YAML
	if !isValidYAML(data) {
		return nil, NewFailure(FailureManifestInvalid, "runtime.yaml", "valid YAML", "invalid YAML")
	}

	evidence = append(evidence, "runtime_manifest:loaded")
	return evidence, nil
}

// determineContext determines context fields from authoritative sources
func (bp *BootstrapProtocol) determineContext(sessionID string) ([]string, error) {
	var evidence []string

	// Get project name from repository
	projectName := filepath.Base(bp.repoPath)
	evidence = append(evidence, fmt.Sprintf("project_name:%s", projectName))

	// Get workspace ID
	workspaceID := fmt.Sprintf("ws-%s", time.Now().Format("20060102-150405"))
	evidence = append(evidence, fmt.Sprintf("workspace_id:%s", workspaceID))

	// Initialize context
	bp.currentContext = NewContext(projectName, bp.repoPath, "Problem", workspaceID, sessionID)

	return evidence, nil
}

// resolveConflicts applies conflict resolution rules
func (bp *BootstrapProtocol) resolveConflicts() ([]string, error) {
	var evidence []string

	// Currently no conflicts to resolve
	// This step is a placeholder for future conflict resolution
	evidence = append(evidence, "conflicts:none")

	return evidence, nil
}

// assembleContext assembles the final context structure
func (bp *BootstrapProtocol) assembleContext() ([]string, error) {
	var evidence []string

	if bp.currentContext == nil {
		return nil, fmt.Errorf("context not initialized")
	}

	// Set default artifact paths
	bp.currentContext.Artifacts = DefaultArtifactPaths(bp.workspacePath)

	// Set default allowed context
	bp.currentContext.AllowedContext = []string{
		".kdse/",
		"docs/",
		"src/",
	}

	// Set restricted paths
	bp.currentContext.RestrictedPaths = []string{
		".git/",
		"node_modules/",
		"vendor/",
		".cache/",
	}

	evidence = append(evidence, "context:assembled")
	return evidence, nil
}

// computeChecksum computes the context checksum
func (bp *BootstrapProtocol) computeChecksum() ([]string, error) {
	var evidence []string

	if bp.currentContext == nil {
		return nil, fmt.Errorf("context not initialized")
	}

	checksum, err := bp.currentContext.ComputeChecksum()
	if err != nil {
		return nil, fmt.Errorf("failed to compute checksum: %w", err)
	}

	bp.currentContext.Checksum = checksum
	evidence = append(evidence, fmt.Sprintf("checksum:%s", checksum))

	return evidence, nil
}

// recordEvidence records bootstrap completion evidence
func (bp *BootstrapProtocol) recordEvidence() ([]string, error) {
	var evidence []string

	// Record bootstrap completion
	evidence = append(evidence, fmt.Sprintf("bootstrap:completed_at=%s", time.Now().Format(time.RFC3339)))
	evidence = append(evidence, "bootstrap:success")

	// Save context to file
	ctxPath := filepath.Join(bp.workspacePath, "context.json")
	data, err := json.MarshalIndent(bp.currentContext, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	if err := os.WriteFile(ctxPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to save context: %w", err)
	}

	evidence = append(evidence, "context:saved")
	return evidence, nil
}

// GetContext returns the current context being built
func (bp *BootstrapProtocol) GetContext() *Context {
	return bp.currentContext
}

// isValidYAML performs basic YAML validation
func isValidYAML(data []byte) bool {
	// Basic check: ensure it's not empty and contains valid characters
	if len(data) == 0 {
		return false
	}

	// Check for common YAML indicators
	content := string(data)
	return strings.Contains(content, "version:") ||
		strings.Contains(content, "runtime:") ||
		strings.Contains(content, "---")
}

// BootstrapStepInfo provides human-readable information about bootstrap steps
var BootstrapStepInfo = map[BootstrapStep]struct {
	Name        string
	Description string
}{
	BootstrapStepDiscoverWorkspace: {
		Name:        "Discover Workspace",
		Description: "Locate and verify existence of .kdse/ directory",
	},
	BootstrapStepVerifyRuntime: {
		Name:        "Verify Runtime",
		Description: "Verify required runtime files exist and are readable",
	},
	BootstrapStepLoadManifest: {
		Name:        "Load Manifest",
		Description: "Load and parse runtime manifest files",
	},
	BootstrapStepDetermineContext: {
		Name:        "Determine Context",
		Description: "Determine all context fields from authoritative sources",
	},
	BootstrapStepResolveConflicts: {
		Name:        "Resolve Conflicts",
		Description: "Apply conflict resolution rules for any contradictory values",
	},
	BootstrapStepAssembleContext: {
		Name:        "Assemble Context",
		Description: "Assemble final engineering context structure",
	},
	BootstrapStepComputeChecksum: {
		Name:        "Compute Checksum",
		Description: "Compute context checksum for integrity verification",
	},
	BootstrapStepRecordEvidence: {
		Name:        "Record Evidence",
		Description: "Mark context as established, record completion evidence",
	},
}

// GetBootstrapSteps returns all bootstrap steps in order
func GetBootstrapSteps() []BootstrapStep {
	return []BootstrapStep{
		BootstrapStepDiscoverWorkspace,
		BootstrapStepVerifyRuntime,
		BootstrapStepLoadManifest,
		BootstrapStepDetermineContext,
		BootstrapStepResolveConflicts,
		BootstrapStepAssembleContext,
		BootstrapStepComputeChecksum,
		BootstrapStepRecordEvidence,
	}
}
