// KDSE MCP Server Tools - v0.3 MCP communication with .kdse/ workspace support
package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kdse/runtime/internal/workspace"
)

// ToolHandler provides MCP tools with workspace-aware responses
type ToolHandler struct {
	repoPath string
	ws       *workspace.Workspace
}

// NewToolHandler creates a new ToolHandler
func NewToolHandler() *ToolHandler {
	repoPath, _ := os.Getwd()
	ws := workspace.New(repoPath)
	return &ToolHandler{
		repoPath: repoPath,
		ws:       ws,
	}
}

// Help returns information about available tools
func (h *ToolHandler) Help() map[string]interface{} {
	return map[string]interface{}{
		"server": map[string]interface{}{
			"name":        "kdse-mcp",
			"version":     "0.3.0",
			"description": "Model Context Protocol server for Knowledge-Driven Software Engineering",
			"protocol":    "2024-11-05",
		},
		"tools": []map[string]interface{}{
			{
				"name":        "help",
				"description": "Returns this help information listing all available KDSE MCP tools",
				"category":    "meta",
				"parameters":  []string{},
				"example":     `{"name": "help"}`,
			},
			{
				"name":        "initialize",
				"description": "Initializes the KDSE .kdse/ workspace. Creates the workspace directory structure if it doesn't exist. Returns workspace information.",
				"category":    "workspace",
				"parameters":  []string{},
				"example":     `{"name": "initialize"}`,
			},
			{
				"name":        "status",
				"description": "Returns current repository status including git state, file counts, KDSE workspace state, and compliance indicators",
				"category":    "workspace",
				"parameters":  []string{},
				"example":     `{"name": "status"}`,
			},
			{
				"name":        "collect",
				"description": "Collects and catalogs engineering artifacts into .kdse/artifacts/",
				"category":    "workspace",
				"parameters":  []string{},
				"example":     `{"name": "collect"}`,
			},
			{
				"name":        "foundation",
				"description": "Returns or creates foundation documentation under .kdse/foundation/",
				"category":    "workspace",
				"parameters":  []string{},
				"example":     `{"name": "foundation"}`,
			},
			{
				"name":        "audit",
				"description": "Generates audit reports under .kdse/reports/",
				"category":    "workspace",
				"parameters":  []string{},
				"example":     `{"name": "audit"}`,
			},
			{
				"name":        "migrate",
				"description": "Migrates any legacy KDSE directories (foundation/, knowledge/, context/, artifacts/) from repository root to .kdse/",
				"category":    "migration",
				"parameters":  []string{},
				"example":     `{"name": "migrate"}`,
			},
		},
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
		},
	}
}

// Initialize initializes the KDSE .kdse/ workspace
func (h *ToolHandler) Initialize() map[string]interface{} {
	// Initialize workspace
	h.ws.Initialize()

	// Check for legacy directories
	migrationReport := h.ws.CheckMigration()

	return map[string]interface{}{
		"repository": map[string]interface{}{
			"root":   h.repoPath,
			"exists": true,
		},
		"workspace": map[string]interface{}{
			"initialized":    h.ws.Exists(),
			"root":           h.ws.Root(),
			"paths":          h.ws.GetPaths(),
		},
		"migration": map[string]interface{}{
			"has_legacy_dirs": migrationReport.HasLegacyDirs,
			"legacy_dirs":     migrationReport.LegacyDirs,
			"can_migrate":      migrationReport.CanMigrate,
			"recommendations": migrationReport.Recommendations,
		},
		"module":              "github.com/kdse/runtime",
		"version":             "0.3.0",
		"goVersion":           runtime.Version(),
		"features":            []string{"help", "initialize", "status", "collect", "foundation", "audit", "migrate"},
		"supported_protocols": []string{"2024-11-05"},
		"components": map[string]interface{}{
			"commands":    []string{"kdse"},
			"packages":    []string{"collect", "config", "context", "detection", "normalize", "report", "state", "types", "workspace"},
			"docs_sections": []string{"audit", "evolution", "execution", "foundation", "runtime"},
			"runtime_docs":  []string{"ARCHITECTURE", "BUILD", "COMMANDS", "EXECUTION_MODEL", "SESSION_PROTOCOL", "WORKFLOW"},
		},
	}
}

// Status returns current repository status information
func (h *ToolHandler) Status() map[string]interface{} {
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

	// Check for legacy directories
	legacyDirs := h.ws.DetectLegacyDirs()
	if len(legacyDirs) > 0 {
		status["kdse"].(map[string]interface{})["legacy_dirs_warning"] = legacyDirs
		status["kdse"].(map[string]interface{})["recommendation"] = "Run 'migrate' to move legacy dirs to .kdse/"
	}

	// Check workspace compliance
	status["kdse"].(map[string]interface{})["compliant"] = len(legacyDirs) == 0

	return status
}

// Collect collects and catalogs artifacts into .kdse/artifacts/
func (h *ToolHandler) Collect() map[string]interface{} {
	// Ensure workspace exists
	h.ws.Initialize()

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
	// Ensure workspace exists
	h.ws.Initialize()

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
	// Ensure workspace exists
	h.ws.Initialize()

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
