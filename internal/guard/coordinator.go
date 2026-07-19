// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kdse/runtime/internal/discover"
	"github.com/kdse/runtime/internal/toolchain"
)

// Coordinator provides a high-level interface for runtime initialization.
// It orchestrates all guards and manages the complete initialization flow.
type Coordinator struct {
	repoPath string
	guard    *RuntimeGuard
}

// NewCoordinator creates a new Coordinator for the given repository path.
// Uses shared discovery to resolve to project root.
func NewCoordinator(projectPath string) *Coordinator {
	// Use shared discovery to resolve to project root
	runtimePaths, err := discover.Resolve(projectPath)
	if err != nil {
		// Fallback to provided path if discovery fails
		return &Coordinator{
			repoPath: projectPath,
			guard:    NewRuntimeGuard(projectPath),
		}
	}
	return &Coordinator{
		repoPath: runtimePaths.RepositoryPath,
		guard:    NewRuntimeGuard(runtimePaths.RepositoryPath),
	}
}

// NewCoordinatorWithRepository creates a coordinator with a pre-resolved repository path.
// This is useful when the repository path has already been resolved.
func NewCoordinatorWithRepository(repoPath string) *Coordinator {
	return &Coordinator{
		repoPath: repoPath,
		guard:    NewRuntimeGuard(repoPath),
	}
}

// Initialize performs a complete runtime initialization.
// This is the primary method for setting up KDSE for the first time.
//
// Before creating the workspace, this method ensures:
// 1. A valid software project exists at the target path (via language-specific files)
// 2. KDSE NEVER creates the project or Git repository
// 3. The .kdse workspace is always created inside the existing project
func (c *Coordinator) Initialize(ctx context.Context) error {
	log.Printf("[COORDINATOR] Starting initialization...")

	// Step 0: Verify a valid project exists
	log.Printf("[COORDINATOR] Step 0: Verifying project exists...")
	projectResult, err := c.EnsureProject(ctx)
	if err != nil {
		return fmt.Errorf("[COORDINATOR] No project detected: %w", err)
	}
	log.Printf("[COORDINATOR] Project verified: %s", projectResult.ProjectName)

	// Step 1: Validate project (should now pass)
	log.Printf("[COORDINATOR] Step 1: Validating project...")
	projectValidation := c.guard.projectGuard.Validate(ctx)
	if !projectValidation.Valid {
		return fmt.Errorf("[COORDINATOR] Project validation failed: %s. %s",
			projectValidation.Error.Message, projectValidation.Error.Hint)
	}
	log.Printf("[COORDINATOR] Project validated: %s", projectValidation.ProjectName)

	// Step 2: Create workspace
	log.Printf("[COORDINATOR] Step 2: Creating workspace...")
	workspacePath := filepath.Join(c.repoPath, ".kdse")
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return fmt.Errorf("[COORDINATOR] Failed to create workspace: %w", err)
	}
	log.Printf("[COORDINATOR] Workspace created: %s", workspacePath)

	// Step 3: Create initial workspace files
	log.Printf("[COORDINATOR] Step 3: Creating workspace files...")
	if err := c.createWorkspaceFiles(workspacePath); err != nil {
		return fmt.Errorf("[COORDINATOR] Failed to create workspace files: %w", err)
	}
	log.Printf("[COORDINATOR] Workspace files created")

	// Step 4: Create session
	log.Printf("[COORDINATOR] Step 4: Creating session...")
	sessionGuard := NewSessionValidationGuard(c.repoPath)
	sessionState, err := sessionGuard.CreateSession()
	if err != nil {
		return fmt.Errorf("[COORDINATOR] Failed to create session: %w", err)
	}
	log.Printf("[COORDINATOR] Session created: %s", sessionState.SessionID)

	// Step 5: Validate complete runtime
	log.Printf("[COORDINATOR] Step 5: Validating runtime...")
	result := c.guard.Validate(ctx)
	if !result.Valid {
		return fmt.Errorf("[COORINATOR] Runtime validation failed: %v", result.Errors)
	}

	log.Printf("[COORDINATOR] Initialization complete!")
	return nil
}

// createWorkspaceFiles creates the initial workspace configuration files
func (c *Coordinator) createWorkspaceFiles(workspacePath string) error {
	timestamp := time.Now().Format(time.RFC3339)

	// Create manifest.yaml
	manifest := map[string]interface{}{
		"version":     "1.0.0",
		"created_at":  timestamp,
		"last_updated": timestamp,
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "manifest.yaml"), manifest); err != nil {
		return err
	}

	// Create session.yaml
	session := map[string]interface{}{
		"session_id":     "initial",
		"started_at":     timestamp,
		"updated_at":     timestamp,
		"workspace_root": workspacePath,
		"version":        "1.0.0",
		"initialized":    false,
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "session.yaml"), session); err != nil {
		return err
	}

	// Create phase.yaml
	phase := map[string]interface{}{
		"current":  "initialization",
		"history":  []string{},
		"version":  "1.0.0",
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "phase.yaml"), phase); err != nil {
		return err
	}

	// Create runtime.yaml
	runtime := map[string]interface{}{
		"version":     "1.0.0",
		"type":        "evidence-driven",
		"created_at":  timestamp,
		"strict_mode": true,
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "runtime.yaml"), runtime); err != nil {
		return err
	}

	// Create metadata.yaml
	metadata := map[string]interface{}{
		"version":    "1.0.0",
		"created_at": timestamp,
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "metadata.yaml"), metadata); err != nil {
		return err
	}

	// Create workspace.yaml
	wsConfig := map[string]interface{}{
		"version":    "1.0.0",
		"root":       workspacePath,
		"created_at": timestamp,
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "workspace.yaml"), wsConfig); err != nil {
		return err
	}

	// Create methodology.yaml
	methodology := map[string]interface{}{
		"version":    "1.0.0",
		"name":       "KDSE",
		"created_at": timestamp,
	}
	if err := c.writeYAML(filepath.Join(workspacePath, "methodology.yaml"), methodology); err != nil {
		return err
	}

	// Create required directories
	requiredDirs := []string{
		"runtime",
		"foundation",
		"knowledge",
		"evidence",
		"reports",
		"references",
		"traceability",
		"artifacts",
		"config",
		"state",
		"sessions",
		"normalized",
		"cache",
		"someday",
		"knowledge/general",
		"knowledge/operational",
		"knowledge/developmental",
	}

	for _, dir := range requiredDirs {
		path := filepath.Join(workspacePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// writeYAML writes a map as YAML to a file
func (c *Coordinator) writeYAML(path string, data map[string]interface{}) error {
	// Simple YAML format writer
	content := "version: " + getString(data, "version", "1.0.0") + "\n"

	for key, value := range data {
		if key == "version" {
			continue
		}

		switch v := value.(type) {
		case string:
			content += fmt.Sprintf("%s: %s\n", key, v)
		case []string:
			if len(v) == 0 {
				content += fmt.Sprintf("%s: []\n", key)
			} else {
				content += fmt.Sprintf("%s:\n", key)
				for _, item := range v {
					content += fmt.Sprintf("  - %s\n", item)
				}
			}
		default:
			content += fmt.Sprintf("%s: %v\n", key, v)
		}
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// getString safely gets a string from a map
func getString(data map[string]interface{}, key, defaultVal string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return defaultVal
}

// Validate runs a full runtime validation
func (c *Coordinator) Validate(ctx context.Context) *RuntimeValidationResult {
	return c.guard.Validate(ctx)
}

// Status returns a human-readable status of the runtime
func (c *Coordinator) Status() string {
	state := c.guard.GetCurrentState()

	status := map[RuntimeState]string{
		StateNoProject:       "No project detected",
		StateProject:         "Project detected, workspace not initialized",
		StateWorkspace:       "Workspace exists, no active session",
		StateSession:         "Session active, lifecycle not ready",
		StateLifecycleReady:  "Runtime ready for operations",
	}

	msg, ok := status[state]
	if !ok {
		msg = "Unknown state"
	}

	return fmt.Sprintf("Runtime State: %s\nStatus: %s", state, msg)
}

// ProjectInfo contains information about a detected project
type ProjectInfo struct {
	ProjectPath  string
	ProjectName  string
	ProjectType  discover.ProjectType
	IsGitRepo    bool
	Indicators   []string
}

// EnsureProject verifies a valid software project exists at the path.
// KDSE does NOT create projects - it only attaches to existing ones.
//
// This method:
// 1. Uses project-first discovery (detects via language-specific files)
// 2. Optionally checks for Git repository (as evidence only)
// 3. Returns project info if found
// 4. Returns error if no project exists
//
// Error message clearly instructs user to initialize their project first.
func (c *Coordinator) EnsureProject(ctx context.Context) (*ProjectInfo, error) {
	// Use shared discovery to check for software project
	runtimePaths, err := discover.Resolve(c.repoPath)
	if err != nil {
		// No software project found
		log.Printf("[COORDINATOR] No software project detected at: %s", c.repoPath)
		return nil, &NoProjectError{
			Path: c.repoPath,
			Message: "No software project detected.\n\n" +
				"KDSE requires a software project before initialization.\n\n" +
				"Please initialize your project first:\n" +
				"  - Go:       go mod init\n" +
				"  - Node.js:  npm init\n" +
				"  - Python:   python -m pip freeze > requirements.txt (or create pyproject.toml)\n" +
				"  - Rust:     cargo init\n" +
				"  - Java:     mvn archetype:generate or create pom.xml\n" +
				"  - .NET:     dotnet new\n" +
				"  - Or create any project structure with source files\n\n" +
				"Then run kdse initialize.",
		}
	}

	log.Printf("[COORDINATOR] Project detected at: %s (type: %s, git: %v)",
		runtimePaths.RepositoryPath, runtimePaths.ProjectType, runtimePaths.IsGitRepo)

	return &ProjectInfo{
		ProjectPath: runtimePaths.RepositoryPath,
		ProjectName: filepath.Base(runtimePaths.RepositoryPath),
		ProjectType: runtimePaths.ProjectType,
		IsGitRepo:   runtimePaths.IsGitRepo,
		Indicators:  runtimePaths.ProjectIndicators,
	}, nil
}

// NoProjectError indicates no software project was found.
type NoProjectError struct {
	Path    string
	Message string
}

func (e *NoProjectError) Error() string {
	return e.Message
}

// VerifyToolchains verifies that all required development toolchains are available.
// This is a critical verification step that must pass before implementation can begin.
//
// KDSE never silently skips verification due to missing tooling.
// If toolchains are missing, this method returns an error with actionable instructions.
func (c *Coordinator) VerifyToolchains() (*ToolchainVerificationResult, error) {
	log.Printf("[COORDINATOR] Verifying development toolchains...")

	result := &ToolchainVerificationResult{
		ProjectPath: c.repoPath,
		Success:     true,
	}

	// Perform runtime verification
	verification := toolchain.VerifyRuntime(c.repoPath, toolchain.VerificationLevelRequired)

	result.Verification = verification
	result.RequiredTools = verification.RequiredTools
	result.MissingTools = verification.MissingTools
	result.InstallBase = verification.InstallBase

	if !verification.Success {
		result.Success = false
		log.Printf("[COORDINATOR] Toolchain verification FAILED: %d missing toolchains", len(verification.MissingTools))

		// Build error message with actionable instructions
		var errMsg strings.Builder
		errMsg.WriteString("Required toolchains are missing:\n\n")
		for _, tcType := range verification.MissingTools {
			errMsg.WriteString(fmt.Sprintf("  %s:\n", tcType))
			errMsg.WriteString(getToolchainInstructions(tcType))
			errMsg.WriteString("\n")
		}
		errMsg.WriteString("\nPlease install the missing toolchains and re-run kdse initialize.\n")
		errMsg.WriteString("Toolchains will be verified automatically.\n")

		result.ErrorMessage = errMsg.String()
		return result, fmt.Errorf(result.ErrorMessage)
	}

	log.Printf("[COORDINATOR] Toolchain verification PASSED: all required toolchains available")
	return result, nil
}

// ToolchainVerificationResult contains the result of toolchain verification
type ToolchainVerificationResult struct {
	ProjectPath   string
	Success       bool
	RequiredTools  []toolchain.ToolchainType
	MissingTools  []toolchain.ToolchainType
	InstallBase   string
	Verification  *toolchain.RuntimeVerification
	ErrorMessage  string
}

// getToolchainInstructions returns installation instructions for a toolchain
func getToolchainInstructions(tcType toolchain.ToolchainType) string {
	switch tcType {
	case toolchain.ToolchainGo:
		return `    Download from: https://go.dev/dl/
    Extract to: /workspace/.tools/go
    Add to PATH: export PATH=/workspace/.tools/go/bin:$PATH`
	case toolchain.ToolchainNode:
		return `    Download from: https://nodejs.org/
    Extract to: /workspace/.tools/node
    Add to PATH: export PATH=/workspace/.tools/node/bin:$PATH`
	case toolchain.ToolchainPython:
		return `    Linux:   apt install python3 python3-pip
    macOS:   brew install python3
    Windows: Download from python.org`
	case toolchain.ToolchainJava:
		return `    Download from: https://adoptium.net/ or https://aws.amazon.com/corretto/
    Set JAVA_HOME: export JAVA_HOME=/path/to/java`
	case toolchain.ToolchainDotNet:
		return `    Download from: https://dotnet.microsoft.com/download
    Add to PATH: export PATH=/path/to/dotnet:$PATH`
	case toolchain.ToolchainRust:
		return `    Install: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
    Toolchain will be installed to ~/.cargo/bin`
	default:
		return `    Installation instructions not available`
	}
}
