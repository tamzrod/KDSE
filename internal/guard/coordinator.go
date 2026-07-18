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
//
// Before creating the workspace, this method ensures:
// 1. A valid engineering project exists at the target path
// 2. If no project exists, a minimal project is created
// 3. The .kdse workspace is always created inside the project
func (c *Coordinator) Initialize(ctx context.Context) error {
	log.Printf("[COORDINATOR] Starting initialization...")

	// Step 0: Ensure a valid project exists before creating workspace
	log.Printf("[COORDINATOR] Step 0: Ensuring valid project exists...")
	projectPath, err := c.EnsureProject(ctx)
	if err != nil {
		return fmt.Errorf("[COORDINATOR] Failed to ensure project: %w", err)
	}
	// Update repoPath if project was created in a subdirectory
	if projectPath != c.repoPath {
		log.Printf("[COORDINATOR] Project created in subdirectory: %s", projectPath)
		log.Printf("[COORDINATOR] Updating working directory to project root...")
		if err := os.Chdir(projectPath); err != nil {
			return fmt.Errorf("[COORDINATOR] Failed to change to project directory: %w", err)
		}
		c.repoPath = projectPath
		// Re-create guard with new path
		c.guard = NewRuntimeGuard(c.repoPath)
	}

	// Step 1: Validate project (should now pass)
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

// EnsureProject ensures a valid engineering project exists at the target path.
// If the current directory is already a valid project, returns the current path.
// If not, creates a minimal project in a subdirectory and returns that path.
//
// Minimal project requirements (per ProjectGuard):
// - At least 2 distinct artifact categories
// This implementation creates a README.md (documentation) and a src/ directory (source_code)
func (c *Coordinator) EnsureProject(ctx context.Context) (string, error) {
	// Check if current directory is already a valid project
	if c.guard.projectGuard.Exists() {
		log.Printf("[COORDINATOR] Valid project already exists at: %s", c.repoPath)
		return c.repoPath, nil
	}

	// Check if .kdse already exists (workspace was initialized without project)
	kdsePath := filepath.Join(c.repoPath, ".kdse")
	if _, err := os.Stat(kdsePath); err == nil {
		// .kdse exists but no project - this is the problematic case
		// We need to create a minimal project
		log.Printf("[COORDINATOR] .kdse exists but no valid project found")
	}

	// Check if directory has any content at all
	entries, err := os.ReadDir(c.repoPath)
	if err != nil {
		return "", fmt.Errorf("cannot read directory: %w", err)
	}

	// Count non-hidden entries
	hasContent := false
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			hasContent = true
			break
		}
	}

	if !hasContent {
		// Empty directory - create minimal project directly in current directory
		log.Printf("[COORDINATOR] Creating minimal project in current directory: %s", c.repoPath)
		if err := c.createMinimalProject(c.repoPath); err != nil {
			return "", fmt.Errorf("failed to create minimal project: %w", err)
		}
		return c.repoPath, nil
	}

	// Directory has content but no valid project
	// Create project in subdirectory
	log.Printf("[COORDINATOR] No valid project found, creating one...")

	// Generate project name from current directory or use default
	projectName := c.generateProjectName(c.repoPath)
	projectPath := filepath.Join(c.repoPath, projectName)

	// Ensure unique path
	counter := 1
	for {
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			break
		}
		projectPath = filepath.Join(c.repoPath, fmt.Sprintf("%s-%d", projectName, counter))
		counter++
	}

	log.Printf("[COORDINATOR] Creating project at: %s", projectPath)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create project directory: %w", err)
	}

	if err := c.createMinimalProject(projectPath); err != nil {
		return "", fmt.Errorf("failed to create minimal project: %w", err)
	}

	return projectPath, nil
}

// createMinimalProject creates the minimum artifacts required for ProjectGuard validation.
// ProjectGuard requires at least 2 distinct artifact categories.
// This creates: README.md (documentation) + src/ directory (source_code)
func (c *Coordinator) createMinimalProject(projectPath string) error {
	// Create README.md
	readmeContent := "# Engineering Project\n\n" +
		"This is an auto-generated engineering project initialized with KDSE.\n\n" +
		"## Project Structure\n\n" +
		"- `src/` - Source code directory\n" +
		"- `.kdse/` - KDSE workspace (created by KDSE runtime)\n\n" +
		"## Getting Started\n\n" +
		"Initialize the KDSE workspace:\n" +
		"```bash\nkdse initialize\n```\n"

	readmePath := filepath.Join(projectPath, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Create src/ directory
	srcPath := filepath.Join(projectPath, "src")
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		return fmt.Errorf("failed to create src directory: %w", err)
	}

	// Create a .gitkeep in src to ensure directory is tracked
	gitkeepPath := filepath.Join(srcPath, ".gitkeep")
	if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create .gitkeep: %w", err)
	}

	// Create .gitignore
	gitignoreContent := `# KDSE
.kdse/

# Dependencies
node_modules/
vendor/
__pycache__/

# Build outputs
dist/
build/
*.o
*.so

# IDE
.vscode/
.idea/
`
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	log.Printf("[COORDINATOR] Minimal project created with README.md and src/")
	return nil
}

// generateProjectName creates a valid directory name from the current path
func (c *Coordinator) generateProjectName(dirPath string) string {
	baseName := filepath.Base(dirPath)

	// If it's a home directory or root, use a default name
	if baseName == "~" || baseName == "" {
		return "engineering-project"
	}

	// Clean the name - remove special characters, lowercase
	name := strings.ToLower(baseName)
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r + 32 // lowercase
		}
		if r >= '0' && r <= '9' {
			return r
		}
		return '-' // replace special chars with hyphen
	}, name)

	// Remove consecutive hyphens
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	name = strings.Trim(name, "-")

	// Ensure non-empty
	if name == "" {
		return "engineering-project"
	}

	// Limit length
	if len(name) > 50 {
		name = name[:50]
	}

	return name
}
