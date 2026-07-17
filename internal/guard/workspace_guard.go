// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WorkspaceGuard validates KDSE workspace existence and integrity.
// It is responsible for:
//   - Discovering existing KDSE workspace
//   - Validating workspace integrity
//
// Workspace Guard must NEVER discover projects.
// Workspace Guard assumes Project Guard already succeeded.
type WorkspaceGuard struct {
	repoPath    string
	workspacePath string
}

// NewWorkspaceGuard creates a new Workspace Guard for the given repository path
func NewWorkspaceGuard(repoPath string) *WorkspaceGuard {
	return &WorkspaceGuard{
		repoPath:     repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
	}
}

// Validate validates workspace existence and integrity.
// Returns WorkspaceGuardResult indicating if workspace is valid.
func (g *WorkspaceGuard) Validate(ctx context.Context) *WorkspaceGuardResult {
	result := &WorkspaceGuardResult{
		RuntimeGuardResult: RuntimeGuardResult{
			Valid:        false,
			State:        StateProject,
			StateBefore:  StateProject,
			StateAfter:   StateProject,
			GuardType:    GuardTypeWorkspace,
		},
		WorkspaceRoot: g.workspacePath,
	}

	// Step 1: Check if .kdse directory exists
	if !g.Exists() {
		result.Error = ErrWorkspaceMissing
		return result
	}

	// Step 2: Validate workspace structure
	if !g.hasValidStructure() {
		result.Error = ErrWorkspaceCorrupted
		return result
	}

	// Step 3: Validate workspace version
	version, err := g.getWorkspaceVersion()
	if err != nil || !g.isVersionCompatible(version) {
		result.Error = ErrWorkspaceVersionMismatch
		return result
	}

	// Step 4: Validate workspace integrity
	if !g.validateIntegrity() {
		result.Error = ErrWorkspaceCorrupted
		return result
	}

	// Workspace is valid
	result.Valid = true
	result.State = StateWorkspace
	result.StateBefore = StateProject
	result.StateAfter = StateWorkspace
	result.Version = version
	result.Integrity = true

	return result
}

// Exists checks if the KDSE workspace directory exists
func (g *WorkspaceGuard) Exists() bool {
	info, err := os.Stat(g.workspacePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Root returns the workspace root path
func (g *WorkspaceGuard) Root() string {
	return g.workspacePath
}

// hasValidStructure checks if workspace has required directories
func (g *WorkspaceGuard) hasValidStructure() bool {
	requiredDirs := []string{
		"runtime",
		"foundation",
		"knowledge",
		"evidence",
		"reports",
	}

	for _, dir := range requiredDirs {
		path := filepath.Join(g.workspacePath, dir)
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			return false
		}
	}

	return true
}

// getWorkspaceVersion reads the workspace version
func (g *WorkspaceGuard) getWorkspaceVersion() (string, error) {
	configPath := filepath.Join(g.workspacePath, "runtime.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Try manifest.yaml as fallback
		configPath = filepath.Join(g.workspacePath, "manifest.yaml")
		data, err = os.ReadFile(configPath)
		if err != nil {
			return "", err
		}
	}

	var config struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		// Try YAML parsing for legacy formats
		return g.parseYAMLVersion(data)
	}

	return config.Version, nil
}

// parseYAMLVersion parses version from YAML content
func (g *WorkspaceGuard) parseYAMLVersion(data []byte) (string, error) {
	content := string(data)

	// Simple version extraction (handles common YAML formats)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "version:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	// Default version for legacy workspaces
	return "1.0.0", nil
}

// isVersionCompatible checks if workspace version is compatible
func (g *WorkspaceGuard) isVersionCompatible(version string) bool {
	// For now, accept any version 1.x
	if strings.HasPrefix(version, "1.") {
		return true
	}

	// Accept explicit "1.0.0" or "1.0"
	if version == "1.0.0" || version == "1.0" {
		return true
	}

	return false
}

// validateIntegrity checks workspace internal consistency
func (g *WorkspaceGuard) validateIntegrity() bool {
	// Check that required files exist and are readable
	requiredFiles := []string{
		"manifest.yaml",
		"session.yaml",
		"phase.yaml",
	}

	for _, file := range requiredFiles {
		path := filepath.Join(g.workspacePath, file)
		data, err := os.ReadFile(path)
		if err != nil {
			return false
		}

		// Basic validation: file should not be empty
		if len(data) == 0 {
			return false
		}
	}

	return true
}

// WorkspaceInfo returns detailed workspace information
type WorkspaceInfo struct {
	Root       string
	Version    string
	CreatedAt  time.Time
	Phase      string
	Integrity  bool
}

// GetInfo returns detailed workspace information
func (g *WorkspaceGuard) GetInfo() (*WorkspaceInfo, error) {
	info := &WorkspaceInfo{
		Root: g.workspacePath,
	}

	// Get version
	version, err := g.getWorkspaceVersion()
	if err != nil {
		return nil, err
	}
	info.Version = version

	// Get phase
	phasePath := filepath.Join(g.workspacePath, "phase.yaml")
	data, err := os.ReadFile(phasePath)
	if err == nil {
		var phaseData struct {
			Current string `json:"current"`
		}
		if json.Unmarshal(data, &phaseData) == nil {
			info.Phase = phaseData.Current
		}
	}

	// Check integrity
	info.Integrity = g.validateIntegrity()

	// Get creation time from manifest
	manifestPath := filepath.Join(g.workspacePath, "manifest.yaml")
	stat, err := os.Stat(manifestPath)
	if err == nil {
		info.CreatedAt = stat.ModTime()
	}

	return info, nil
}
