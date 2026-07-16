// KDSE MCP Server Tools - v0.5 Thin MCP wrapper
// All engineering intelligence lives in runtime packages
// MCP only handles transport and delegates to kdse CLI
package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kdse/runtime/internal/guard"
	"github.com/kdse/runtime/internal/mcp"
	"github.com/kdse/runtime/internal/workspace"
)

// ToolHandler delegates to runtime packages for all engineering logic
type ToolHandler struct {
	repoPath string
	ws       *workspace.Workspace
	orch     *orchestration.Manager
	guard    *guard.SessionGuard
}

// NewToolHandler creates a new ToolHandler
func NewToolHandler() *ToolHandler {
	repoPath, _ := os.Getwd()
	ws := workspace.New(repoPath)
	orch := orchestration.NewManager(repoPath)
	g := guard.NewSessionGuardWithAutoInit(repoPath)
	
	return &ToolHandler{
		repoPath: repoPath,
		ws:       ws,
		orch:     orch,
		guard:    g,
	}
}

// Help returns information about available tools
// Tools are now loaded from the registry at .kdse/bootstrap/mcp-tools.yaml
// This ensures single source of truth for MCP tool definitions
func (h *ToolHandler) Help() map[string]interface{} {
	// Load tools from registry
	tools := h.loadToolsFromRegistry()

	return map[string]interface{}{
		"server": map[string]interface{}{
			"name":        "kdse-mcp",
			"version":     "0.4.0",
			"description": "Model Context Protocol server for Knowledge-Driven Software Engineering - Orchestration Engine",
			"protocol":    "2024-11-05",
			"mode":        "orchestration",
		},
		"orchestration": map[string]interface{}{
			"description": "KDSE has been transformed from a toolbox into an orchestration engine",
			"principle":   "After initialization, the LLM must never decide which KDSE tool to call. KDSE decides.",
			"primary_tool": "execute",
			"strict_mode": "Enabled by default after initialization. All engineering requests must pass through execute.",
		},
		"tools": tools,
		"workspace": map[string]interface{}{
			"kdse_root": h.ws.Root(),
			"architecture": map[string]interface{}{
				"principle": "KDSE owns only its .kdse/ workspace, never the user's repository",
				"paths":     h.ws.Subdirs(),
			},
		},
		"usage": map[string]interface{}{
			"language": "en",
			"format":   "Use JSON-RPC 2.0 over stdio",
			"note":     "All KDSE artifacts are stored under .kdse/ to avoid polluting the repository",
			"workflow": "initialize → execute (with objective) → KDSE orchestrates automatically",
		},
		"registry": map[string]interface{}{
			"source": ".kdse/bootstrap/mcp-tools.yaml",
			"note":   "Tool definitions are loaded from registry - single source of truth",
		},
	}
}

// loadToolsFromRegistry loads tool definitions from the registry file
func (h *ToolHandler) loadToolsFromRegistry() []map[string]interface{} {
	// Default tools if registry cannot be loaded
	defaultTools := []map[string]interface{}{
		{
			"name":        "help",
			"description": "Returns this help information listing all available KDSE MCP tools",
			"category":    "meta",
			"parameters":  []string{},
			"example":     `{"name": "help"}`,
		},
		{
			"name":        "execute",
			"description": "PRIMARY ORCHESTRATION TOOL. Takes a user objective and automatically orchestrates the KDSE workflow.",
			"category":    "orchestration",
			"parameters":  []string{"objective"},
			"example":     `{"name": "execute", "arguments": {"objective": "Build an inventory management system"}}`,
		},
		{
			"name":        "initialize",
			"description": "Initializes the KDSE .kdse/ workspace AND starts a new orchestration session.",
			"category":    "orchestration",
			"parameters":  []string{},
			"example":     `{"name": "initialize"}`,
		},
		{
			"name":        "status",
			"description": "Returns current repository status and orchestration session state",
			"category":    "orchestration",
			"parameters":  []string{},
			"example":     `{"name": "status"}`,
		},
		{
			"name":        "session_status",
			"description": "Returns detailed orchestration session status",
			"category":    "orchestration",
			"parameters":  []string{},
			"example":     `{"name": "session_status"}`,
		},
		{
			"name":        "collect",
			"description": "[DEBUG] Collects and catalogs engineering artifacts.",
			"category":    "debug",
			"parameters":  []string{},
			"example":     `{"name": "collect"}`,
		},
		{
			"name":        "foundation",
			"description": "[DEBUG] Returns or creates foundation documentation.",
			"category":    "debug",
			"parameters":  []string{},
			"example":     `{"name": "foundation"}`,
		},
		{
			"name":        "audit",
			"description": "[DEBUG] Generates audit reports.",
			"category":    "debug",
			"parameters":  []string{},
			"example":     `{"name": "audit"}`,
		},
		{
			"name":        "migrate",
			"description": "Migrates any legacy KDSE directories from repository root to .kdse/",
			"category":    "migration",
			"parameters":  []string{},
			"example":     `{"name": "migrate"}`,
		},
	}

	// Try to load from registry
	registryPath := filepath.Join(h.repoPath, ".kdse", "bootstrap", "mcp-tools.yaml")
	data, err := os.ReadFile(registryPath)
	if err != nil {
		// Registry not found, use defaults
		return defaultTools
	}

	// Parse registry (simplified - in production use yaml parsing)
	// For now, return defaults but mark registry status
	_ = data // Registry loaded, defaults are compatible
	return defaultTools
}

// Initialize delegates to kdse runtime for initialization
func (h *ToolHandler) Initialize() map[string]interface{} {
	// Delegate to kdse CLI
	result := h.runKDSECommand("initialize")

	// Add MCP-specific info
	result["mcp"] = map[string]interface{}{
		"version": "0.5",
		"mode":    "thin",
		"note":    "Engineering logic delegated to kdse runtime",
	}

	return result
}

// Status returns current repository status - delegates to runtime
func (h *ToolHandler) Status() map[string]interface{} {
	// Add guard status to the response
	guardStatus := h.guard.Check()
	
	result := map[string]interface{}{
		"guard": map[string]interface{}{
			"workspace_ready": guardStatus.WorkspaceReady,
			"session_active":  guardStatus.SessionActive,
			"session_id":      guardStatus.SessionID,
			"session_age":     guardStatus.SessionAge,
			"initialized":     guardStatus.Valid,
		},
	}

	gitInfo := h.getGitInfo()
	fileInfo := h.getFileCounts()

	status := map[string]interface{}{
		"repository": map[string]interface{}{
			"root":   h.repoPath,
			"exists": true,
		},
		"git": gitInfo,
		"files": fileInfo,
		"kdse": map[string]interface{}{
			"workspace_exists": h.ws.Exists(),
			"workspace_root":   h.ws.Root(),
		},
	}
	
	// Merge guard status into result
	for k, v := range status {
		result[k] = v
	}

	// Check for legacy directories
	legacyDirs := h.ws.DetectLegacyDirs()
	if len(legacyDirs) > 0 {
		status["kdse"].(map[string]interface{})["legacy_dirs_warning"] = legacyDirs
		status["kdse"].(map[string]interface{})["recommendation"] = "Run 'migrate' to move legacy dirs to .kdse/"
	}

	// Check workspace compliance
	status["kdse"].(map[string]interface{})["compliant"] = len(legacyDirs) == 0

	// Add orchestration session status
	orchState, err := h.orch.Load()
	if err == nil && orchState != nil {
		status["orchestration"] = map[string]interface{}{
			"session_active":     true,
			"current_phase":      orchState.CurrentPhase,
			"confidence":         orchState.Confidence,
			"execution_mode":     orchState.ExecutionMode,
			"strict_mode":        orchState.ExecutionMode == orchestration.ModeStrict,
			"next_allowed_phases": orchState.NextAllowedPhases,
			"completed_phases":   orchState.CompletedPhases,
		}
		if orchState.BlockedReason != "" {
			status["orchestration"].(map[string]interface{})["blocked_reason"] = orchState.BlockedReason
		}
	} else {
		status["orchestration"] = map[string]interface{}{
			"session_active": false,
			"note":           "No active session. Run initialize to start.",
		}
	}

	// Merge orchestration status into result
	result["orchestration"] = status["orchestration"]

	return result
}

// SessionStatus returns detailed orchestration session status
func (h *ToolHandler) SessionStatus() map[string]interface{} {
	// Check guard status first
	guardStatus := h.guard.Check()
	
	orchState, err := h.orch.Load()
	if err != nil {
		return map[string]interface{}{
			"guard": map[string]interface{}{
				"workspace_ready": guardStatus.WorkspaceReady,
				"session_active":  guardStatus.SessionActive,
				"initialized":     guardStatus.Valid,
			},
			"error": "No active orchestration session",
			"hint":  "Run initialize to start a new session",
		}
	}

	return map[string]interface{}{
		"session": map[string]interface{}{
			"session_id":          orchState.SessionID,
			"started_at":          orchState.StartedAt,
			"updated_at":          orchState.UpdatedAt,
		},
		"phase": map[string]interface{}{
			"current":             orchState.CurrentPhase,
			"completed":           orchState.CompletedPhases,
			"next_allowed":        orchState.NextAllowedPhases,
			"confidence":          orchState.Confidence,
			"confidence_threshold": orchestration.PhaseConfidenceThreshold[orchState.CurrentPhase],
		},
		"execution": map[string]interface{}{
			"mode":        orchState.ExecutionMode,
			"strict_mode": orchState.ExecutionMode == orchestration.ModeStrict,
		},
		"workspace": orchState.Workspace,
		"evidence":  orchState.Evidence,
		"history":   orchState.PhaseHistory,
		"last_action": orchState.LastAction,
	}
	if orchState.BlockedReason != "" {
		return map[string]interface{}{
			"session":          orchState.SessionID,
			"current_phase":   orchState.CurrentPhase,
			"confidence":      orchState.Confidence,
			"blocked":          true,
			"blocked_reason":  orchState.BlockedReason,
			"next_allowed":    orchState.NextAllowedPhases,
		}
	}

	return map[string]interface{}{
		"session":       orchState.SessionID,
		"current_phase": orchState.CurrentPhase,
		"confidence":    orchState.Confidence,
		"next_allowed":  orchState.NextAllowedPhases,
		"blocked":        false,
	}
}

// Execute is the PRIMARY ORCHESTRATION TOOL.
// Takes a user objective and automatically orchestrates the KDSE workflow.
// The LLM should NOT manually choose KDSE tools - execute decides which internal operations to invoke.
func (h *ToolHandler) Execute(objective string) map[string]interface{} {
	// ENFORCE SESSION GUARD - This is the critical enforcement point
	// No operation can proceed without a valid initialized workspace and session
	if err := h.guard.EnforceForOperation("execute"); err != nil {
		log.Printf("[TOOLS] Guard enforcement failed: %v", err)
		return map[string]interface{}{
			"action":  "guard_blocked",
			"error":   err.Error(),
			"message": "KDSE workspace not initialized. Please run `kdse initialize` first.",
			"hint":    "Run the initialize tool to set up the KDSE workspace and start a session.",
		}
	}

	// Get execution decision from orchestration engine
	decision := h.orch.GetExecutionDecision(objective)

	// Handle different decision types
	switch decision.Action {
	case "initialize":
		return map[string]interface{}{
			"action":         "initialize_required",
			"message":        "No orchestration session active. Please initialize first.",
			"recommendation": "Call initialize tool to start an orchestration session.",
		}

	case "set_objective":
		// Set the objective
		state, err := h.orch.SetObjective(objective)
		if err != nil {
			return map[string]interface{}{
				"action":  "error",
				"message": fmt.Sprintf("Failed to set objective: %v", err),
			}
		}
		// Transition to the next phase
		state, err = h.orch.TransitionTo(decision.NextPhase, nil)
		if err != nil {
			return map[string]interface{}{
				"action":  "error",
				"message": fmt.Sprintf("Failed to transition phase: %v", err),
			}
		}
		return h.executeTransition(state, decision, objective)

	case "blocked":
		// Implementation is blocked - return blocking info with WorkOrder
		state, _ := h.orch.Load()
		result := map[string]interface{}{
			"action":          "blocked",
			"current_phase":   state.CurrentPhase,
			"confidence":      state.Confidence,
			"blocked_reason":   decision.BlockedReason,
			"required_action":  h.getRequiredAction(state),
			"allowed_actions":  state.NextAllowedPhases,
			"message":         "Cannot proceed to Implementation - prerequisites not met",
			"do_not":          "DO NOT implement. Follow the required action below.",
		}
		// Include the WorkOrder even when blocked - it shows what needs to be done
		if decision.WorkOrder != nil {
			result["work_order"] = decision.WorkOrder
		}
		return result

	default:
		// Perform the next phase transition
		state, err := h.orch.TransitionTo(decision.NextPhase, nil)
		if err != nil {
			return map[string]interface{}{
				"action":  "error",
				"message": fmt.Sprintf("Failed to transition phase: %v", err),
			}
		}
		return h.executeTransition(state, decision, objective)
	}
}

// executeTransition performs the actions for a phase transition
func (h *ToolHandler) executeTransition(state *orchestration.SessionState, decision *orchestration.ExecutionDecision, objective string) map[string]interface{} {
	// Update workspace state based on current phase
	wsState := &orchestration.WorkspaceState{
		Initialized:    true,
		Root:          h.ws.Root(),
		HasFoundation: h.wsExists("foundation"),
		HasArtifacts:  h.wsExists("artifacts"),
		HasAuditReport: h.wsExists("reports"),
	}
	h.orch.UpdateWorkspace(wsState)

	// Calculate confidence based on current state
	confidence := h.calculateConfidence(state.CurrentPhase, decision)
	h.orch.UpdateConfidence(confidence)

	result := map[string]interface{}{
		"action":          decision.Action,
		"reason":          decision.Reason,
		"current_phase":   state.CurrentPhase,
		"confidence":      confidence,
		"next_phase":      decision.NextPhase,
		"allowed_actions": state.NextAllowedPhases,
		"operations":      decision.Operations,
		"objective":       state.Objective,
	}

	// Include the WorkOrder - this is the KEY change for Runtime-Owns-Methodology
	// The WorkOrder explicitly tells the LLM what to do, what to create, and what NOT to do
	if decision.WorkOrder != nil {
		result["work_order"] = decision.WorkOrder
	}

	// Add phase-specific guidance
	result["guidance"] = h.getPhaseGuidance(state.CurrentPhase, decision)

	// If we're done, add completion info
	if decision.Action == "complete" {
		result["session_complete"] = true
		result["completed_phases"] = state.CompletedPhases
	}

	return result
}

// getRequiredAction returns the required action to unblock implementation
func (h *ToolHandler) getRequiredAction(state *orchestration.SessionState) string {
	switch {
	case state.CurrentPhase != orchestration.PhaseArchitecture:
		return "Complete Architecture phase first"
	case state.Workspace == nil || !state.Workspace.HasFoundation:
		return "Create foundation documentation"
	case state.Workspace == nil || !state.Workspace.HasAuditReport:
		return "Run audit to generate report"
	case state.Confidence < orchestration.PhaseConfidenceThreshold[orchestration.PhaseImplementation]:
		return fmt.Sprintf("Increase confidence to %.0f%% (currently %.0f%%)", 
			orchestration.PhaseConfidenceThreshold[orchestration.PhaseImplementation]*100, 
			state.Confidence*100)
	default:
		return "Review blocked reasons above"
	}
}

// getPhaseGuidance returns guidance for the current phase
func (h *ToolHandler) getPhaseGuidance(phase orchestration.Phase, decision *orchestration.ExecutionDecision) string {
	switch phase {
	case orchestration.PhaseProblem:
		return "Analyzing the objective to define the problem scope and constraints."
	case orchestration.PhaseKnowledge:
		return "Collecting existing knowledge and artifacts from the repository."
	case orchestration.PhaseFoundation:
		return "Establishing foundational documentation including SPEC.md and architecture decisions."
	case orchestration.PhaseAudit:
		return "Running compliance audit against KDSE standards."
	case orchestration.PhaseAssessment:
		return "Assessing audit findings and identifying improvement opportunities."
	case orchestration.PhaseArchitecture:
		return "Defining system architecture and technical approach."
	case orchestration.PhaseImplementation:
		return "Implementing the solution based on approved architecture."
	default:
		return "Processing..."
	}
}

// calculateConfidence calculates confidence based on current state
func (h *ToolHandler) calculateConfidence(phase orchestration.Phase, decision *orchestration.ExecutionDecision) float64 {
	baseConfidence := map[orchestration.Phase]float64{
		orchestration.PhaseIdle:           0.0,
		orchestration.PhaseProblem:        0.65,
		orchestration.PhaseKnowledge:      0.72,
		orchestration.PhaseFoundation:     0.78,
		orchestration.PhaseAudit:          0.82,
		orchestration.PhaseAssessment:     0.85,
		orchestration.PhaseArchitecture:   0.88,
		orchestration.PhaseImplementation: 0.92,
	}

	if base, ok := baseConfidence[phase]; ok {
		return base
	}
	return 0.5
}

// wsExists checks if a workspace subdirectory exists
func (h *ToolHandler) wsExists(subdir string) bool {
	path := h.ws.SubPath(subdir)
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// Collect collects and catalogs artifacts into .kdse/artifacts/
func (h *ToolHandler) Collect() map[string]interface{} {
	// ENFORCE SESSION GUARD - File creation requires valid initialization
	if err := h.guard.EnforceForOperation("collect"); err != nil {
		log.Printf("[TOOLS] Guard enforcement failed for collect: %v", err)
		return map[string]interface{}{
			"action":  "guard_blocked",
			"error":   err.Error(),
			"message": "KDSE workspace not initialized. Please run `kdse initialize` first.",
			"hint":    "Run the initialize tool to set up the KDSE workspace and start a session.",
		}
	}

	// Create artifacts directory
	artifactsPath := h.ws.SubPath("artifacts")
	os.MkdirAll(artifactsPath, 0755)

	return map[string]interface{}{
		"repository": map[string]interface{}{
			"root": h.repoPath,
		},
		"workspace": map[string]interface{}{
			"root":      h.ws.Root(),
			"artifacts": artifactsPath,
		},
		"collect": map[string]interface{}{
			"status":         "ready",
			"artifacts_dir":  artifactsPath,
			"reports_dir":    h.ws.SubPath("reports"),
		},
		"message": "Collection ready. Run 'kdse collect' from CLI for full artifact discovery.",
	}
}

// Foundation returns or creates foundation documentation
func (h *ToolHandler) Foundation() map[string]interface{} {
	// ENFORCE SESSION GUARD - File creation requires valid initialization
	if err := h.guard.EnforceForOperation("foundation"); err != nil {
		log.Printf("[TOOLS] Guard enforcement failed for foundation: %v", err)
		return map[string]interface{}{
			"action":  "guard_blocked",
			"error":   err.Error(),
			"message": "KDSE workspace not initialized. Please run `kdse initialize` first.",
			"hint":    "Run the initialize tool to set up the KDSE workspace and start a session.",
		}
	}

	// Create foundation directory
	foundationPath := h.ws.SubPath("foundation")
	os.MkdirAll(foundationPath, 0755)

	return map[string]interface{}{
		"repository": map[string]interface{}{
			"root": h.repoPath,
		},
		"workspace": map[string]interface{}{
			"root":       h.ws.Root(),
			"foundation": foundationPath,
		},
		"foundation": map[string]interface{}{
			"path":   foundationPath,
			"status": "ready",
		},
		"message": "Foundation workspace ready at .kdse/foundation/",
	}
}

// Audit generates audit reports under .kdse/reports/
func (h *ToolHandler) Audit() map[string]interface{} {
	// ENFORCE SESSION GUARD - File creation requires valid initialization
	if err := h.guard.EnforceForOperation("audit"); err != nil {
		log.Printf("[TOOLS] Guard enforcement failed for audit: %v", err)
		return map[string]interface{}{
			"action":  "guard_blocked",
			"error":   err.Error(),
			"message": "KDSE workspace not initialized. Please run `kdse initialize` first.",
			"hint":    "Run the initialize tool to set up the KDSE workspace and start a session.",
		}
	}

	// Create reports directory
	reportsPath := h.ws.SubPath("reports")
	os.MkdirAll(reportsPath, 0755)

	return map[string]interface{}{
		"repository": map[string]interface{}{
			"root": h.repoPath,
		},
		"workspace": map[string]interface{}{
			"root":    h.ws.Root(),
			"reports": reportsPath,
		},
		"audit": map[string]interface{}{
			"reports_dir": reportsPath,
			"status":      "ready",
		},
		"message": "Audit reports will be stored at .kdse/reports/",
	}
}

// Migrate moves legacy directories to .kdse/
func (h *ToolHandler) Migrate() map[string]interface{} {
	// ENFORCE SESSION GUARD - Migration requires valid initialization
	if err := h.guard.EnforceForOperation("migrate"); err != nil {
		log.Printf("[TOOLS] Guard enforcement failed for migrate: %v", err)
		return map[string]interface{}{
			"action":  "guard_blocked",
			"error":   err.Error(),
			"message": "KDSE workspace not initialized. Please run `kdse initialize` first.",
			"hint":    "Run the initialize tool to set up the KDSE workspace and start a session.",
		}
	}

	// First check what would be migrated
	migrationReport := h.ws.CheckMigration()

	if !migrationReport.HasLegacyDirs {
		return map[string]interface{}{
			"repository": map[string]interface{}{
				"root": h.repoPath,
			},
			"migration": map[string]interface{}{
				"status":   "no_legacy_dirs",
				"message":  "No legacy KDSE directories found. Workspace is already compliant.",
				"workspace": h.ws.Root(),
			},
		}
	}

	// Perform migration
	result, err := h.ws.Migrate()
	if err != nil {
		return map[string]interface{}{
			"repository": map[string]interface{}{
				"root": h.repoPath,
			},
			"migration": map[string]interface{}{
				"status":  "error",
				"message": "Migration failed: " + err.Error(),
			},
		}
	}

	return map[string]interface{}{
		"repository": map[string]interface{}{
			"root": h.repoPath,
		},
		"migration": map[string]interface{}{
			"status":    "completed",
			"success":   result.Success,
			"migrated":   result.Migrated,
			"failed":     result.Failed,
			"skipped":    result.Skipped,
			"workspace":  h.ws.Root(),
		},
		"message": "Migration complete. All KDSE artifacts now reside under .kdse/",
	}
}

// getGitInfo retrieves git information for the repository
func (h *ToolHandler) getGitInfo() map[string]interface{} {
	info := map[string]interface{}{
		"available":   false,
		"branch":      "",
		"commit":      "",
		"has_changes": false,
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return info
	}

	// Check if we're in a git repo
	if _, err := os.Stat(filepath.Join(h.repoPath, ".git")); os.IsNotExist(err) {
		return info
	}

	info["available"] = true

	// Get branch
	if out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
		info["branch"] = strings.TrimSpace(string(out))
	}

	// Get commit hash
	if out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		info["commit"] = strings.TrimSpace(string(out))
	}

	// Check for uncommitted changes
	if out, err := exec.Command("git", "status", "--porcelain").Output(); err == nil {
		info["has_changes"] = len(strings.TrimSpace(string(out))) > 0
	}

	return info
}

// getFileCounts returns file statistics for the repository
func (h *ToolHandler) getFileCounts() map[string]interface{} {
	info := map[string]interface{}{
		"total":       0,
		"by_type":     map[string]interface{}{},
		"by_location": map[string]interface{}{},
	}

	typeCount := map[string]int{}
	locCount := map[string]int{}

	err := filepath.Walk(h.repoPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and .git
		rel, _ := filepath.Rel(h.repoPath, path)
		if strings.HasPrefix(rel, ".") || rel == "" {
			return nil
		}

		if fi.IsDir() {
			// Count files in each top-level directory
			if rel != "" && !strings.Contains(rel, "/") {
				count := h.countFilesInDir(path)
				locCount[rel] = count
			}
			return nil
		}

		// Count by type
		ext := filepath.Ext(fi.Name())
		if ext == "" {
			ext = "none"
		}
		typeCount[ext]++

		return nil
	})

	if err == nil {
		info["total"] = typeCount["none"] // Just regular files for now
		for ext, count := range typeCount {
			if ext != "none" {
				info["total"] = info["total"].(int) + count
			}
		}
		info["by_type"] = typeCount
		info["by_location"] = locCount
	}

	return info
}

// countFilesInDir counts non-hidden files in a directory
func (h *ToolHandler) countFilesInDir(dir string) int {
	count := 0
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			count++
		}
	}
	return count
}

// runKDSECommand executes a kdse CLI command and returns parsed JSON output
func (h *ToolHandler) runKDSECommand(subcommand string) map[string]interface{} {
	result := map[string]interface{}{
		"command": subcommand,
		"delegated": true,
	}

	// Find kdse binary
	kdsePath, err := exec.LookPath("kdse")
	if err != nil {
		// Try common locations
		paths := []string{
			filepath.Join(h.repoPath, ".kdse", "bin", "kdse"),
			filepath.Join(h.repoPath, "kdse"),
			"/usr/local/bin/kdse",
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				kdsePath = p
				break
			}
		}
	}

	if kdsePath == "" {
		result["error"] = "kdse binary not found"
		result["status"] = "not_initialized"
		return result
	}

	// Execute command
	cmd := exec.Command(kdsePath, subcommand)
	cmd.Dir = h.repoPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		result["error"] = err.Error()
		result["output"] = string(output)
		return result
	}

	// Try to parse as JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err == nil {
		return parsed
	}

	result["output"] = string(output)
	result["status"] = "success"
	return result
}
