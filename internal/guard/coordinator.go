// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Coordinator provides a high-level interface for runtime initialization.
// It orchestrates all guards and manages the complete initialization flow.
type Coordinator struct {
	repoPath string
	guard    *RuntimeGuard
}

// NewCoordinator creates a new Coordinator for the given repository path
func NewCoordinator(repoPath string) *Coordinator {
	return &Coordinator{
		repoPath: repoPath,
		guard:    NewRuntimeGuard(repoPath),
	}
}

// Initialize performs a complete runtime initialization.
// This is the primary method for setting up KDSE for the first time.
func (c *Coordinator) Initialize(ctx context.Context) error {
	log.Printf("[COORDINATOR] Starting initialization...")

	// Step 1: Validate project first
	log.Printf("[COORDINATOR] Step 1: Validating project...")
	projectResult := c.guard.projectGuard.Validate(ctx)
	if !projectResult.Valid {
		return fmt.Errorf("[COORDINATOR] Project validation failed: %s. %s",
			projectResult.Error.Message, projectResult.Error.Hint)
	}
	log.Printf("[COORDINATOR] Project validated: %s", projectResult.ProjectName)

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
