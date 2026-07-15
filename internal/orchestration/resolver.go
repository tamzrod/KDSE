package orchestration

import (
	"os"
	"path/filepath"
	"strings"
)

// WorkspaceResolver resolves workspace paths without hardcoding /app or /workspace
type WorkspaceResolver struct {
	config     *EngineConfig
	currentWD  string
}

// NewWorkspaceResolver creates a new workspace resolver
func NewWorkspaceResolver(config *EngineConfig) (*WorkspaceResolver, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &WorkspaceResolver{
		config:    config,
		currentWD: wd,
	}, nil
}

// NewWorkspaceResolverWithPath creates a resolver with a specific working directory
func NewWorkspaceResolverWithPath(workingDir string, config *EngineConfig) *WorkspaceResolver {
	if config == nil {
		config = DefaultEngineConfig()
	}
	return &WorkspaceResolver{
		config:    config,
		currentWD: workingDir,
	}
}

// ResolveWorkspace resolves the workspace starting from the given path
// Hierarchy: Repository → Project Folder → Temporary Workspace
func (r *WorkspaceResolver) ResolveWorkspace(startPath string) (*WorkspaceInfo, error) {
	// Normalize the starting path
	startPath = r.normalizePath(startPath)

	// Determine workspace type and resolve paths
	workspaceType := r.detectWorkspaceType(startPath)

	var repoPath, projectPath, tempPath, kdsePath string

	switch workspaceType {
	case WorkspaceTypeTemporary:
		// In temp workspace, find parent project/repository
		projectPath = r.findProjectPath(startPath)
		repoPath = r.findRepositoryPath(startPath)
		tempPath = startPath
		kdsePath = filepath.Join(tempPath, ".kdse")

	case WorkspaceTypeProject:
		projectPath = startPath
		repoPath = startPath
		tempPath = ""
		kdsePath = filepath.Join(projectPath, ".kdse")

	case WorkspaceTypeRepository:
		repoPath = startPath
		projectPath = ""
		tempPath = ""
		kdsePath = filepath.Join(repoPath, ".kdse")

	default:
		// Unknown type, treat as repository
		repoPath = startPath
		kdsePath = filepath.Join(repoPath, ".kdse")
	}

	return &WorkspaceInfo{
		ResolvedPath:    startPath,
		WorkspaceType:   workspaceType,
		RepositoryPath:  repoPath,
		ProjectPath:     projectPath,
		TempPath:       tempPath,
		KDSEPath:       kdsePath,
	}, nil
}

// ResolveTemporaryWorkspace creates a temporary workspace under ./temp/.kdse
func (r *WorkspaceResolver) ResolveTemporaryWorkspace(projectName string) (*WorkspaceInfo, error) {
	// Build temp path: currentDir/temp/.kdse/<projectName>
	baseTempDir := r.resolveTempBase()
	tempDir := filepath.Join(baseTempDir, ".kdse", projectName)

	// Ensure temp directory exists
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}

	return &WorkspaceInfo{
		ResolvedPath:   tempDir,
		WorkspaceType:  WorkspaceTypeTemporary,
		TempPath:       tempDir,
		KDSEPath:       filepath.Join(tempDir, ".kdse"),
	}, nil
}

// MigrateToProject migrates .kdse from temp workspace to project workspace
func (r *WorkspaceResolver) MigrateToProject(tempKDSEPath, projectPath string) error {
	if !r.config.EnableMigration {
		return nil
	}

	projectKDSEPath := filepath.Join(projectPath, ".kdse")

	// Check if source exists
	if _, err := os.Stat(tempKDSEPath); os.IsNotExist(err) {
		return nil // Nothing to migrate
	}

	// Ensure project directory exists
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	// Check if project .kdse already exists
	if _, err := os.Stat(projectKDSEPath); err == nil {
		// Project already has .kdse, merge or skip based on preference
		return r.mergeKDSE(tempKDSEPath, projectKDSEPath)
	}

	// Move temp .kdse to project
	return os.Rename(tempKDSEPath, projectKDSEPath)
}

// ResolveKDSEPath returns the KDSE path for a given working directory
func (r *WorkspaceResolver) ResolveKDSEPath(workingDir string) string {
	normalized := r.normalizePath(workingDir)
	return filepath.Join(normalized, ".kdse")
}

// ResolveSubPath returns a subdirectory path within the KDSE workspace
func (r *WorkspaceResolver) ResolveSubPath(workingDir, subdir string) string {
	kdsePath := r.ResolveKDSEPath(workingDir)
	return filepath.Join(kdsePath, subdir)
}

// detectWorkspaceType determines what type of workspace the path represents
func (r *WorkspaceResolver) detectWorkspaceType(path string) WorkspaceType {
	normalized := r.normalizePath(path)

	// Check if path contains temp directory marker
	if strings.Contains(normalized, filepath.Join(r.config.TempWorkspaceBase, ".kdse")) {
		return WorkspaceTypeTemporary
	}

	// Check if path contains project marker or is a known project
	if r.isProjectDirectory(normalized) {
		return WorkspaceTypeProject
	}

	// Check if it's a git repository
	if r.isGitRepository(normalized) {
		return WorkspaceTypeRepository
	}

	return WorkspaceTypeUnknown
}

// isProjectDirectory checks if the directory is a project root
func (r *WorkspaceResolver) isProjectDirectory(path string) bool {
	// Check for common project indicators
	projectMarkers := []string{
		"package.json",
		"go.mod",
		"Cargo.toml",
		"pom.xml",
		"build.gradle",
		"Makefile",
		"CMakeLists.txt",
		".project",
	}

	for _, marker := range projectMarkers {
		if _, err := os.Stat(filepath.Join(path, marker)); err == nil {
			return true
		}
	}

	return false
}

// isGitRepository checks if the directory is a git repository
func (r *WorkspaceResolver) isGitRepository(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// findProjectPath finds the parent project directory from a temp workspace
func (r *WorkspaceResolver) findProjectPath(tempPath string) string {
	// Navigate up from temp/.kdse/<project> to find project root
	parts := strings.Split(tempPath, string(filepath.Separator))
	
	// Find index of temp directory
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == r.config.TempWorkspaceBase {
			// Return everything up to and including temp
			if i+1 < len(parts) {
				return filepath.Join(parts[:i+1]...)
			}
			break
		}
	}

	return filepath.Dir(tempPath)
}

// findRepositoryPath finds the parent repository from a temp workspace
func (r *WorkspaceResolver) findRepositoryPath(tempPath string) string {
	// For now, same as project path
	// In future, could check for .git in parent directories
	return r.findProjectPath(tempPath)
}

// resolveTempBase returns the base temp directory path
func (r *WorkspaceResolver) resolveTempBase() string {
	return filepath.Join(r.currentWD, r.config.TempWorkspaceBase)
}

// normalizePath normalizes a file path
func (r *WorkspaceResolver) normalizePath(path string) string {
	if path == "" {
		return r.currentWD
	}

	// Handle relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(r.currentWD, path)
	}

	// Clean the path
	return filepath.Clean(path)
}

// GetCurrentWorkingDirectory returns the current working directory
func (r *WorkspaceResolver) GetCurrentWorkingDirectory() string {
	return r.currentWD
}

// SetCurrentWorkingDirectory sets the current working directory
func (r *WorkspaceResolver) SetCurrentWorkingDirectory(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Verify directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return os.ErrNotExist
	}

	r.currentWD = absPath
	return nil
}

// mergeKDSE merges source KDSE into destination (basic implementation)
func (r *WorkspaceResolver) mergeKDSE(src, dst string) error {
	// For now, just copy files that don't exist in destination
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		srcFile := filepath.Join(src, entry.Name())
		dstFile := filepath.Join(dst, entry.Name())

		if _, err := os.Stat(dstFile); os.IsNotExist(err) {
			data, err := os.ReadFile(srcFile)
			if err != nil {
				continue
			}
			os.WriteFile(dstFile, data, 0644)
		}
	}

	return nil
}

// ResolvePath resolves any path relative to the workspace
func (r *WorkspaceResolver) ResolvePath(relativePath string) string {
	return filepath.Join(r.currentWD, relativePath)
}

// ResolveFromWorkspace resolves a path relative to the KDSE workspace
func (r *WorkspaceResolver) ResolveFromWorkspace(workspaceInfo *WorkspaceInfo, relativePath string) string {
	return filepath.Join(workspaceInfo.KDSEPath, relativePath)
}
