// KDSE MCP Server Tools - Read-only access to KDSE repository knowledge
package tools

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// DefaultRepositoryRoot is the default path to the KDSE repository
	DefaultRepositoryRoot = "/workspace/project/KDSE"
)

type ToolHandler struct {
	repoRoot string
}

func NewToolHandler() *ToolHandler {
	repoRoot := os.Getenv("KDSE_REPO_ROOT")
	if repoRoot == "" {
		repoRoot = DefaultRepositoryRoot
	}
	return &ToolHandler{repoRoot: repoRoot}
}

// Help returns information about available tools
func (h *ToolHandler) Help() map[string]interface{} {
	tools := []map[string]interface{}{
		{
			"name":        "help",
			"description": "Returns this help information listing all available KDSE MCP tools",
			"category":    "meta",
			"parameters":  []string{},
			"example":     `{"name": "help"}`,
		},
		{
			"name":        "initialize",
			"description": "Returns KDSE repository initialization information including module, version, and features",
			"category":    "repository",
			"parameters":  []string{},
			"example":     `{"name": "initialize"}`,
		},
		{
			"name":        "status",
			"description": "Returns current KDSE repository status including git state, file counts, and compliance indicators",
			"category":    "repository",
			"parameters":  []string{},
			"example":     `{"name": "status"}`,
		},
	}

	return map[string]interface{}{
		"server": map[string]interface{}{
			"name":        "kdse-mcp-server",
			"version":     "0.1.0",
			"description": "Model Context Protocol server for Knowledge-Driven Software Engineering",
			"protocol":    "2024-11-05",
		},
		"tools": tools,
		"usage": map[string]interface{}{
			"language":    "en",
			"format":      "Use JSON-RPC 2.0 over stdio",
			"note":        "This server reads from the KDSE repository. It does not own knowledge.",
			"repository":  h.repoRoot,
		},
	}
}

// Initialize returns repository initialization information
func (h *ToolHandler) Initialize() map[string]interface{} {
	return h.collectInitializeInfo()
}

// Status returns repository status information
func (h *ToolHandler) Status() map[string]interface{} {
	return h.collectStatusInfo()
}

func (h *ToolHandler) collectInitializeInfo() map[string]interface{} {
	info := map[string]interface{}{
		"repository": map[string]interface{}{
			"root":      h.repoRoot,
			"exists":    false,
		},
		"module":      nil,
		"version":     nil,
		"goVersion":   nil,
		"features":    []string{},
		"components":  map[string]interface{}{},
	}

	if _, err := os.Stat(h.repoRoot); os.IsNotExist(err) {
		return info
	}
	info["repository"].(map[string]interface{})["exists"] = true

	// Read go.mod
	goModPath := filepath.Join(h.repoRoot, "go.mod")
	if data, err := os.ReadFile(goModPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				info["module"] = strings.TrimPrefix(line, "module ")
			} else if strings.HasPrefix(line, "go ") {
				info["goVersion"] = strings.TrimPrefix(line, "go ")
			}
		}
	}

	// Discover components
	components := map[string]interface{}{}

	// Check cmd directory
	cmdPath := filepath.Join(h.repoRoot, "cmd")
	if entries, err := os.ReadDir(cmdPath); err == nil {
		cmds := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				cmds = append(cmds, entry.Name())
			}
		}
		if len(cmds) > 0 {
			components["commands"] = cmds
		}
	}

	// Check internal directory
	internalPath := filepath.Join(h.repoRoot, "internal")
	if entries, err := os.ReadDir(internalPath); err == nil {
		packages := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				packages = append(packages, entry.Name())
			}
		}
		if len(packages) > 0 {
			components["packages"] = packages
		}
	}

	// Check runtime directory
	runtimePath := filepath.Join(h.repoRoot, "runtime")
	if _, err := os.Stat(runtimePath); err == nil {
		if entries, err := os.ReadDir(runtimePath); err == nil {
			docs := []string{}
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					docs = append(docs, strings.TrimSuffix(entry.Name(), ".md"))
				}
			}
			if len(docs) > 0 {
				components["runtime_docs"] = docs
			}
		}
	}

	// Check docs directory
	docsPath := filepath.Join(h.repoRoot, "docs")
	if _, err := os.Stat(docsPath); err == nil {
		if entries, err := os.ReadDir(docsPath); err == nil {
			sections := []string{}
			for _, entry := range entries {
				if entry.IsDir() {
					sections = append(sections, entry.Name())
				}
			}
			if len(sections) > 0 {
				components["docs_sections"] = sections
			}
		}
	}

	info["components"] = components

	// Define v0.1 features
	info["features"] = []string{
		"help",
		"initialize",
		"status",
	}

	info["version"] = "0.1.0"
	info["supported_protocols"] = []string{"2024-11-05"}

	return info
}

func (h *ToolHandler) collectStatusInfo() map[string]interface{} {
	status := map[string]interface{}{
		"repository": map[string]interface{}{
			"root":   h.repoRoot,
			"exists": false,
		},
		"git":      nil,
		"files":    nil,
		"kdse":     nil,
	}

	if _, err := os.Stat(h.repoRoot); os.IsNotExist(err) {
		return status
	}
	status["repository"].(map[string]interface{})["exists"] = true

	// Git status
	gitStatus := h.getGitStatus()
	status["git"] = gitStatus

	// File counts
	fileCounts := h.countFiles()
	status["files"] = fileCounts

	// KDSE compliance indicators
	kdseStatus := h.checkKDSEStatus()
	status["kdse"] = kdseStatus

	return status
}

func (h *ToolHandler) getGitStatus() map[string]interface{} {
	gitInfo := map[string]interface{}{
		"available":    false,
		"branch":      nil,
		"commit":      nil,
		"has_changes": false,
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return gitInfo
	}

	gitInfo["available"] = true

	// Get branch
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = h.repoRoot
	if out, err := cmd.Output(); err == nil {
		branch := strings.TrimSpace(string(out))
		if branch != "" {
			gitInfo["branch"] = branch
		}
	}

	// Get latest commit
	cmd = exec.Command("git", "rev-parse", "--short=8", "HEAD")
	cmd.Dir = h.repoRoot
	if out, err := cmd.Output(); err == nil {
		gitInfo["commit"] = strings.TrimSpace(string(out))
	}

	// Check for uncommitted changes
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = h.repoRoot
	if out, err := cmd.Output(); err == nil {
		gitInfo["has_changes"] = len(strings.TrimSpace(string(out))) > 0
	}

	return gitInfo
}

func (h *ToolHandler) countFiles() map[string]interface{} {
	counts := map[string]interface{}{
		"total":       0,
		"by_type":     map[string]interface{}{},
		"by_location": map[string]interface{}{},
	}

	typeCounts := map[string]int{}
	locationCounts := map[string]int{}
	total := 0

	filepath.Walk(h.repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and common non-source paths
		relPath, _ := filepath.Rel(h.repoRoot, path)
		if strings.HasPrefix(relPath, ".") ||
			strings.Contains(relPath, "node_modules") ||
			strings.Contains(relPath, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		total++
		ext := filepath.Ext(info.Name())
		if ext == "" {
			ext = "no_extension"
		}
		typeCounts[ext]++

		// Count by top-level location
		parts := strings.SplitN(relPath, string(filepath.Separator), 2)
		if len(parts) > 0 {
			locationCounts[parts[0]]++
		}

		return nil
	})

	counts["total"] = total
	counts["by_type"] = typeCounts
	counts["by_location"] = locationCounts

	return counts
}

func (h *ToolHandler) checkKDSEStatus() map[string]interface{} {
	kdseStatus := map[string]interface{}{
		"compliant":         false,
		"has_readme":        false,
		"has_foundation":    false,
		"has_runtime":       false,
		"has_go_mod":        false,
		"has_docker":        false,
		"has_mcp_server":    false,
		"missing_artifacts": []string{},
	}

	requiredArtifacts := []string{}

	// Check README
	if _, err := os.Stat(filepath.Join(h.repoRoot, "README.md")); err == nil {
		kdseStatus["has_readme"] = true
	} else {
		requiredArtifacts = append(requiredArtifacts, "README.md")
	}

	// Check docs/foundation
	if _, err := os.Stat(filepath.Join(h.repoRoot, "docs", "foundation")); err == nil {
		kdseStatus["has_foundation"] = true
	} else {
		requiredArtifacts = append(requiredArtifacts, "docs/foundation/")
	}

	// Check runtime
	if _, err := os.Stat(filepath.Join(h.repoRoot, "runtime")); err == nil {
		kdseStatus["has_runtime"] = true
	} else {
		requiredArtifacts = append(requiredArtifacts, "runtime/")
	}

	// Check go.mod
	if _, err := os.Stat(filepath.Join(h.repoRoot, "go.mod")); err == nil {
		kdseStatus["has_go_mod"] = true
	} else {
		requiredArtifacts = append(requiredArtifacts, "go.mod")
	}

	// Check Dockerfile
	if _, err := os.Stat(filepath.Join(h.repoRoot, "Dockerfile")); err == nil {
		kdseStatus["has_docker"] = true
	}

	// Check mcp-server
	if _, err := os.Stat(filepath.Join(h.repoRoot, "mcp-server")); err == nil {
		kdseStatus["has_mcp_server"] = true
	}

	kdseStatus["missing_artifacts"] = requiredArtifacts
	kdseStatus["compliant"] = kdseStatus["has_readme"].(bool) &&
		kdseStatus["has_foundation"].(bool) &&
		kdseStatus["has_runtime"].(bool) &&
		kdseStatus["has_go_mod"].(bool)

	return kdseStatus
}

// MarshalJSON helper for pretty printing
func (h *ToolHandler) HelpJSON() string {
	data, _ := json.MarshalIndent(h.Help(), "", "  ")
	return string(data)
}

func (h *ToolHandler) InitializeJSON() string {
	data, _ := json.MarshalIndent(h.Initialize(), "", "  ")
	return string(data)
}

func (h *ToolHandler) StatusJSON() string {
	data, _ := json.MarshalIndent(h.Status(), "", "  ")
	return string(data)
}
