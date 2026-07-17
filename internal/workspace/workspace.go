// Package workspace manages the KDSE .kdse/ workspace directory structure
// and enforces the architectural rule that KDSE owns only its own workspace,
// never the user's repository.
package workspace

import (
	"os"
	"path/filepath"
	"strings"
)

// Workspace manages the .kdse/ directory structure
type Workspace struct {
	repoPath string
	kdsePath string
}

// Subdirectory constants defining the KDSE workspace structure
const (
	SubDirFoundation    = "foundation"
	SubDirKnowledge     = "knowledge"
	SubDirArchitecture  = "architecture"
	SubDirImplementation = "implementation"
	SubDirVerification  = "verification"
	SubDirContext       = "context"
	SubDirArtifacts     = "artifacts"
	SubDirRuntime       = "runtime"
	SubDirSessions      = "sessions"
	SubDirConfidence    = "confidence"
	SubDirOperational   = "operational"
	SubDirDevelopmental = "developmental"
	SubDirReports       = "reports"
	SubDirDocs          = "docs"
	SubDirCache         = "cache"
	SubDirNormalized    = "normalized"
	SubDirEvidence      = "evidence"
)

// Legacy directory names that should be migrated to .kdse/
var LegacyDirs = []string{
	"foundation",
	"knowledge",
	"context",
	"artifacts",
}

// New creates a new Workspace for the given repository path
func New(repoPath string) *Workspace {
	return &Workspace{
		repoPath: repoPath,
		kdsePath: filepath.Join(repoPath, ".kdse"),
	}
}

// Root returns the absolute path to the .kdse/ directory
func (w *Workspace) Root() string {
	return w.kdsePath
}

// RepoPath returns the repository root path
func (w *Workspace) RepoPath() string {
	return w.repoPath
}

// SubPath returns the absolute path to a subdirectory within .kdse/
func (w *Workspace) SubPath(subdir string) string {
	return filepath.Join(w.kdsePath, subdir)
}

// Initialize creates the .kdse/ workspace directory structure
func (w *Workspace) Initialize() error {
	// Create the main .kdse/ directory
	if err := os.MkdirAll(w.kdsePath, 0755); err != nil {
		return err
	}

	// Create required subdirectories only when needed (lazy creation)
	// The workspace is ready to use with just the root directory
	// Subdirectories are created on-demand by individual tools

	return nil
}

// EnsureSubdir creates a subdirectory if it doesn't exist
func (w *Workspace) EnsureSubdir(subdir string) error {
	path := w.SubPath(subdir)
	return os.MkdirAll(path, 0755)
}

// Subdirs returns a map of all standard KDSE subdirectories
func (w *Workspace) Subdirs() map[string]string {
	return map[string]string{
		SubDirFoundation:     w.SubPath(SubDirFoundation),
		SubDirKnowledge:      w.SubPath(SubDirKnowledge),
		SubDirArchitecture:   w.SubPath(SubDirArchitecture),
		SubDirImplementation: w.SubPath(SubDirImplementation),
		SubDirVerification:   w.SubPath(SubDirVerification),
		SubDirContext:        w.SubPath(SubDirContext),
		SubDirArtifacts:      w.SubPath(SubDirArtifacts),
		SubDirRuntime:        w.SubPath(SubDirRuntime),
		SubDirSessions:       w.SubPath(SubDirSessions),
		SubDirConfidence:     w.SubPath(SubDirConfidence),
		SubDirOperational:    w.SubPath(SubDirOperational),
		SubDirDevelopmental:  w.SubPath(SubDirDevelopmental),
		SubDirReports:        w.SubPath(SubDirReports),
		SubDirDocs:           w.SubPath(SubDirDocs),
		SubDirCache:          w.SubPath(SubDirCache),
		SubDirNormalized:     w.SubPath(SubDirNormalized),
		SubDirEvidence:       w.SubPath(SubDirEvidence),
	}
}

// Exists checks if the .kdse/ workspace exists
func (w *Workspace) Exists() bool {
	info, err := os.Stat(w.kdsePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// DetectLegacyDirs returns any legacy KDSE directories found in the repository root
func (w *Workspace) DetectLegacyDirs() []string {
	var found []string
	for _, dir := range LegacyDirs {
		path := filepath.Join(w.repoPath, dir)
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			found = append(found, dir)
		}
	}
	return found
}

// MigrationReport contains information about legacy directories that need migration
type MigrationReport struct {
	HasLegacyDirs   bool     `json:"has_legacy_dirs"`
	LegacyDirs      []string `json:"legacy_dirs"`
	CanMigrate      bool     `json:"can_migrate"`
	Recommendations []string `json:"recommendations"`
}

// CheckMigration returns a report on any legacy directories that need migration
func (w *Workspace) CheckMigration() *MigrationReport {
	legacy := w.DetectLegacyDirs()
	report := &MigrationReport{
		HasLegacyDirs:   len(legacy) > 0,
		LegacyDirs:      legacy,
		CanMigrate:      len(legacy) > 0,
		Recommendations: []string{},
	}

	if report.HasLegacyDirs {
		report.Recommendations = append(report.Recommendations,
			"Legacy KDSE directories detected in repository root",
			"Run 'kdse migrate' to move them under .kdse/",
		)
	}

	return report
}

// Migrate moves legacy directories from repository root to .kdse/
func (w *Workspace) Migrate() (*MigrationResult, error) {
	legacy := w.DetectLegacyDirs()
	if len(legacy) == 0 {
		return &MigrationResult{
			Success:    true,
			Migrated:   []string{},
			Failed:     []string{},
			Skipped:   []string{},
		}, nil
	}

	// Ensure .kdse/ exists
	if err := w.Initialize(); err != nil {
		return nil, err
	}

	result := &MigrationResult{
		Success:    true,
		Migrated:   []string{},
		Failed:     []string{},
		Skipped:   []string{},
	}

	// Create a mapping from legacy names to new names
	// Most legacy dirs map directly, but we need to be careful about conflicts
	dirMapping := map[string]string{
		"foundation": SubDirFoundation,
		"knowledge":  SubDirKnowledge,
		"context":    SubDirContext,
		"artifacts":  SubDirArtifacts,
	}

	for _, legacyDir := range legacy {
		srcPath := filepath.Join(w.repoPath, legacyDir)
		destDir := legacyDir
		if newName, ok := dirMapping[legacyDir]; ok {
			destDir = newName
		}
		destPath := w.SubPath(destDir)

		// Check if destination already exists
		if _, err := os.Stat(destPath); err == nil {
			result.Skipped = append(result.Skipped, legacyDir+" -> .kdse/"+destDir+" (destination exists)")
			continue
		}

		// Move the directory
		if err := os.Rename(srcPath, destPath); err != nil {
			result.Failed = append(result.Failed, legacyDir+" -> .kdse/"+destDir+": "+err.Error())
			result.Success = false
		} else {
			result.Migrated = append(result.Migrated, legacyDir+" -> .kdse/"+destDir)
		}
	}

	return result, nil
}

// MigrationResult contains the outcome of a migration operation
type MigrationResult struct {
	Success      bool     `json:"success"`
	Migrated     []string `json:"migrated"`
	Failed       []string `json:"failed"`
	Skipped      []string `json:"skipped"`
}

// ResolvePath resolves a KDSE workspace path from a short name
func (w *Workspace) ResolvePath(shortPath string) string {
	// Handle both .kdse/path and direct subdirectory names
	if strings.HasPrefix(shortPath, ".kdse/") {
		return filepath.Join(w.repoPath, shortPath)
	}
	// Direct subdirectory name
	return w.SubPath(shortPath)
}

// Paths returns all important workspace paths as a structure
type Paths struct {
	Root          string `json:"root"`
	Foundation    string `json:"foundation"`
	Knowledge     string `json:"knowledge"`
	Architecture  string `json:"architecture"`
	Implementation string `json:"implementation"`
	Verification  string `json:"verification"`
	Context       string `json:"context"`
	Artifacts     string `json:"artifacts"`
	Runtime       string `json:"runtime"`
	Sessions      string `json:"sessions"`
	Confidence    string `json:"confidence"`
	Operational   string `json:"operational"`
	Developmental string `json:"developmental"`
	Reports       string `json:"reports"`
	Docs          string `json:"docs"`
	Cache         string `json:"cache"`
	Normalized    string `json:"normalized"`
	Evidence      string `json:"evidence"`
}

// GetPaths returns all workspace paths
func (w *Workspace) GetPaths() *Paths {
	return &Paths{
		Root:           w.kdsePath,
		Foundation:     w.SubPath(SubDirFoundation),
		Knowledge:      w.SubPath(SubDirKnowledge),
		Architecture:   w.SubPath(SubDirArchitecture),
		Implementation: w.SubPath(SubDirImplementation),
		Verification:   w.SubPath(SubDirVerification),
		Context:        w.SubPath(SubDirContext),
		Artifacts:      w.SubPath(SubDirArtifacts),
		Runtime:        w.SubPath(SubDirRuntime),
		Sessions:       w.SubPath(SubDirSessions),
		Confidence:     w.SubPath(SubDirConfidence),
		Operational:    w.SubPath(SubDirOperational),
		Developmental:  w.SubPath(SubDirDevelopmental),
		Reports:        w.SubPath(SubDirReports),
		Docs:           w.SubPath(SubDirDocs),
		Cache:          w.SubPath(SubDirCache),
		Normalized:     w.SubPath(SubDirNormalized),
		Evidence:       w.SubPath(SubDirEvidence),
	}
}
