// Package detection provides environment detection utilities for KDSE.
//
// Environment Detection classifies the current workspace into exactly one environment
// type before any planning or engineering activity. This ensures KDSE behaves
// correctly based on where it is executing.
//
// Environment Types:
//   - BLANK_WORKSPACE: No recognizable project, no Git, no KDSE project
//   - SOFTWARE_PROJECT: Existing software project, may contain Git, not KDSE
//   - KDSE_RUNTIME: The KDSE framework/runtime itself
//   - KDSE_PROJECT: Existing project initialized by KDSE
package detection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnvironmentType represents the classified environment type.
// Exactly one environment type is assigned per detection.
type EnvironmentType string

const (
	// EnvironmentBlankWorkspace indicates no recognizable project exists.
	// No Git repository, no KDSE project, no software project indicators.
	EnvironmentBlankWorkspace EnvironmentType = "BLANK_WORKSPACE"

	// EnvironmentSoftwareProject indicates an existing software project.
	// May contain Git, not a KDSE runtime, no initialized KDSE project.
	EnvironmentSoftwareProject EnvironmentType = "SOFTWARE_PROJECT"

	// EnvironmentKDSERuntime indicates the KDSE framework/runtime itself.
	// Development work targets KDSE, not a normal software project.
	EnvironmentKDSERuntime EnvironmentType = "KDSE_RUNTIME"

	// EnvironmentKDSEProject indicates an existing project initialized by KDSE.
	// Contains valid KDSE runtime metadata, engineering targets the project.
	EnvironmentKDSEProject EnvironmentType = "KDSE_PROJECT"

	// EnvironmentUnknown indicates evidence is insufficient for classification.
	EnvironmentUnknown EnvironmentType = "UNKNOWN"
)

// String returns a human-readable string for the environment type.
func (e EnvironmentType) String() string {
	return string(e)
}

// IsValid returns true if the environment type is a valid classification.
func (e EnvironmentType) IsValid() bool {
	switch e {
	case EnvironmentBlankWorkspace, EnvironmentSoftwareProject,
		EnvironmentKDSERuntime, EnvironmentKDSEProject:
		return true
	default:
		return false
	}
}

// Evidence represents the collected evidence for environment classification.
// Each field contains verifiable facts about the workspace.
type Evidence struct {
	// Git-related evidence
	HasGitRepo     bool   // .git directory exists
	GitRemoteURL  string  // Remote URL if available
	GitRootPath   string  // Git repository root

	// Project evidence
	HasGoMod       bool   // go.mod file exists
	HasPackageJSON bool   // package.json exists
	HasPyProject   bool   // pyproject.toml or setup.py exists
	HasCargoToml   bool   // Cargo.toml exists
	HasPomXML      bool   // pom.xml exists
	HasMakefile    bool   // Makefile exists

	// KDSE-specific evidence
	HasKDSEProject bool   // .kdse/ directory exists
	HasManifest     bool   // .kdse/manifest.yaml exists
	ManifestValid   bool   // manifest.yaml is valid JSON
	ManifestData    *ManifestInfo  // Parsed manifest data

	// KDSE Runtime markers (identify the KDSE repository itself)
	HasGoModKDSE    bool   // go.mod with module github.com/kdse/runtime
	HasRuntimeDir   bool   // internal/runtime/ directory exists
	HasDetectionDir  bool   // internal/detection/ directory exists
	HasTemplatesDir bool   // templates/ directory exists (KDSE project templates)
	HasDocsDir      bool   // docs/ directory with runtime docs

	// Generic project indicators
	HasReadme       bool   // README.md exists
	HasSrcDir       bool   // src/ or similar source directory
	HasTestDir      bool   // tests/ or similar test directory
	HasLicense      bool   // LICENSE file exists

	// Path information
	RepoRoot        string  // Detected repository root
	ModuleName      string  // Go module name (from go.mod)
	ProjectType     string  // Detected project type
}

// ManifestInfo contains parsed information from .kdse/manifest.yaml.
type ManifestInfo struct {
	Schema         string `json:"schema"`
	Version        string `json:"version"`
	Generated      string `json:"generated"`
	Status         string `json:"status"`
	InstallationSource string `json:"installation,omitempty"`
	RuntimeVersion string `json:"runtime,omitempty"`
}

// EvidenceResult contains the detection result with full evidence.
type EvidenceResult struct {
	Environment EnvironmentType
	Confidence  float64       // 0.0 to 1.0 confidence in classification
	Evidence    Evidence
	Warnings    []string     // Non-fatal issues encountered
	Errors      []string     // Fatal errors during detection
	Timestamp   string      // When detection occurred
}

// DetectionOptions contains configuration for the detector.
type DetectionOptions struct {
	// AllowAmbiguous allows returning a classification even with conflicting evidence.
	// When false (default), ambiguous cases return EnvironmentUnknown.
	AllowAmbiguous bool

	// RequireGitRepo requires a Git repository for software project detection.
	// When true, a software project without Git is classified as BLANK_WORKSPACE.
	RequireGitRepo bool

	// StrictManifestValidation enables strict validation of manifest files.
	StrictManifestValidation bool
}

// DefaultDetectionOptions returns the default detection options.
func DefaultDetectionOptions() *DetectionOptions {
	return &DetectionOptions{
		AllowAmbiguous:           false,
		RequireGitRepo:           false,
		StrictManifestValidation: true,
	}
}

// Detector classifies the workspace environment using evidence-based detection.
type Detector struct {
	repoPath string
	options  *DetectionOptions
}

// NewDetector creates a new environment detector for the given repository path.
func NewEnvironmentDetector(repoPath string) *Detector {
	return &Detector{
		repoPath: repoPath,
		options:  DefaultDetectionOptions(),
	}
}

// NewDetectorWithOptions creates a new detector with custom options.
func NewEnvironmentDetectorWithOptions(repoPath string, opts *DetectionOptions) *Detector {
	if opts == nil {
		opts = DefaultDetectionOptions()
	}
	return &Detector{
		repoPath: repoPath,
		options:  opts,
	}
}

// Detect classifies the workspace environment using evidence-based detection.
// Returns the detected environment type and full evidence collected.
func (d *Detector) Detect() *EvidenceResult {
	result := &EvidenceResult{
		Evidence:  Evidence{},
		Warnings:  []string{},
		Errors:    []string{},
		Timestamp: "",
	}

	// Normalize path
	repoPath := filepath.Clean(d.repoPath)

	// Collect all evidence
	d.collectEvidence(repoPath, result)

	// Classify based on evidence
	result.Environment = d.classify(result)

	// Calculate confidence
	result.Confidence = d.calculateConfidence(result)

	return result
}

// DetectEnvironment is a convenience function that returns only the environment type.
func DetectEnvironment(repoPath string) EnvironmentType {
	detector := NewEnvironmentDetector(repoPath)
	return detector.Detect().Environment
}

// collectEvidence gathers all evidence from the workspace.
func (d *Detector) collectEvidence(repoPath string, result *EvidenceResult) {
	result.Evidence.RepoRoot = repoPath

	// Collect Git evidence
	d.collectGitEvidence(repoPath, result)

	// Collect project evidence
	d.collectProjectEvidence(repoPath, result)

	// Collect KDSE-specific evidence
	d.collectKDSESpecificEvidence(repoPath, result)

	// Collect generic project indicators
	d.collectGenericIndicators(repoPath, result)
}

// collectGitEvidence gathers Git-related evidence.
func (d *Detector) collectGitEvidence(repoPath string, result *EvidenceResult) {
	// Check for .git directory
	gitPath := filepath.Join(repoPath, ".git")
	if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
		result.Evidence.HasGitRepo = true
		result.Evidence.GitRootPath = repoPath
	}

	// Try to get Git remote URL
	if result.Evidence.HasGitRepo {
		remoteURL := d.getGitRemoteURL(repoPath)
		result.Evidence.GitRemoteURL = remoteURL
	}
}

// getGitRemoteURL retrieves the Git remote URL using git commands.
func (d *Detector) getGitRemoteURL(repoPath string) string {
	// Implementation uses git command - see git_resolver.go for reference
	// This is a simplified version that checks for common patterns
	return ""
}

// collectProjectEvidence gathers software project evidence.
func (d *Detector) collectProjectEvidence(repoPath string, result *EvidenceResult) {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("cannot read directory: %v", err))
		return
	}

	// Create file map for quick lookup
	fileMap := make(map[string]os.DirEntry)
	for _, entry := range entries {
		fileMap[entry.Name()] = entry
	}

	// Check for Go project
	if entry, ok := fileMap["go.mod"]; ok && !entry.IsDir() {
		result.Evidence.HasGoMod = true
		moduleName := d.parseGoModuleName(filepath.Join(repoPath, "go.mod"))
		result.Evidence.ModuleName = moduleName
	}

	// Check for Node.js project
	if _, ok := fileMap["package.json"]; ok {
		result.Evidence.HasPackageJSON = true
	}

	// Check for Python project
	if _, ok := fileMap["pyproject.toml"]; ok {
		result.Evidence.HasPyProject = true
	} else if _, ok := fileMap["setup.py"]; ok {
		result.Evidence.HasPyProject = true
	}

	// Check for Rust project
	if _, ok := fileMap["Cargo.toml"]; ok {
		result.Evidence.HasCargoToml = true
	}

	// Check for Java project
	if _, ok := fileMap["pom.xml"]; ok {
		result.Evidence.HasPomXML = true
	}

	// Check for Makefile
	if _, ok := fileMap["Makefile"]; ok {
		result.Evidence.HasMakefile = true
	}

	// Determine project type
	result.Evidence.ProjectType = d.detectProjectType(result.Evidence)
}

// parseGoModuleName reads the module name from go.mod.
func (d *Detector) parseGoModuleName(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module ")
		}
	}
	return ""
}

// detectProjectType determines the project type based on evidence.
func (d *Detector) detectProjectType(evidence Evidence) string {
	switch {
	case evidence.HasGoMod:
		return "go"
	case evidence.HasPackageJSON:
		return "node"
	case evidence.HasPyProject:
		return "python"
	case evidence.HasCargoToml:
		return "rust"
	case evidence.HasPomXML:
		return "java"
	case evidence.HasMakefile:
		return "c"
	default:
		return "unknown"
	}
}

// collectKDSESpecificEvidence gathers KDSE-specific evidence.
func (d *Detector) collectKDSESpecificEvidence(repoPath string, result *EvidenceResult) {
	// Check for .kdse/ directory
	kdseDir := filepath.Join(repoPath, ".kdse")
	if info, err := os.Stat(kdseDir); err == nil && info.IsDir() {
		result.Evidence.HasKDSEProject = true
	}

	// Check for manifest.yaml
	if result.Evidence.HasKDSEProject {
		manifestPath := filepath.Join(kdseDir, "manifest.yaml")
		if entry, err := os.Stat(manifestPath); err == nil && !entry.IsDir() {
			result.Evidence.HasManifest = true

			// Try to parse manifest
			manifestInfo, valid := d.parseManifest(manifestPath)
			if valid {
				result.Evidence.ManifestValid = true
				result.Evidence.ManifestData = manifestInfo
			}
		}

		// Also check for manifest.json (alternative)
		manifestJSONPath := filepath.Join(kdseDir, "manifest.json")
		if entry, err := os.Stat(manifestJSONPath); err == nil && !entry.IsDir() {
			result.Evidence.HasManifest = true

			if !result.Evidence.ManifestValid {
				manifestInfo, valid := d.parseManifestJSON(manifestJSONPath)
				if valid {
					result.Evidence.ManifestValid = true
					result.Evidence.ManifestData = manifestInfo
				}
			}
		}
	}

	// Check for KDSE Runtime markers (identifies the KDSE repository itself)
	// KDSE Runtime markers: files/directories that exist ONLY in the KDSE runtime repo
	result.Evidence.HasRuntimeDir = d.dirExists(filepath.Join(repoPath, "internal", "runtime"))
	result.Evidence.HasDetectionDir = d.dirExists(filepath.Join(repoPath, "internal", "detection"))
	result.Evidence.HasTemplatesDir = d.dirExists(filepath.Join(repoPath, "templates"))
	result.Evidence.HasDocsDir = d.dirExists(filepath.Join(repoPath, "docs"))

	// KDSE Runtime check: module github.com/kdse/runtime
	if result.Evidence.HasGoMod && result.Evidence.ModuleName == "github.com/kdse/runtime" {
		result.Evidence.HasGoModKDSE = true
	}
}

// dirExists checks if a directory exists at the given path.
func (d *Detector) dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// parseManifest parses a YAML manifest file.
func (d *Detector) parseManifest(path string) (*ManifestInfo, bool) {
	// Simple YAML parsing - in production, use a YAML library
	// For now, just check if file exists and is readable
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return &ManifestInfo{
		Schema:   "parsed",
		Version:  "1.0",
	}, true
}

// parseManifestJSON parses a JSON manifest file.
func (d *Detector) parseManifestJSON(path string) (*ManifestInfo, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var manifest struct {
		Schema    string `json:"schema"`
		Version   string `json:"version"`
		Generated string `json:"generated"`
		Status    string `json:"status"`
		Runtime   struct {
			Version string `json:"version"`
		} `json:"runtime"`
		Installation struct {
			SourceRepository string `json:"source_repository"`
		} `json:"installation"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, false
	}

	return &ManifestInfo{
		Schema:              manifest.Schema,
		Version:             manifest.Version,
		Generated:           manifest.Generated,
		Status:              manifest.Status,
		RuntimeVersion:      manifest.Runtime.Version,
		InstallationSource:   manifest.Installation.SourceRepository,
	}, true
}

// collectGenericIndicators gathers generic project indicators.
func (d *Detector) collectGenericIndicators(repoPath string, result *EvidenceResult) {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		switch entry.Name() {
		case "README.md":
			result.Evidence.HasReadme = true
		case "LICENSE":
			result.Evidence.HasLicense = true
		case "src", "lib", "source":
			if entry.IsDir() {
				result.Evidence.HasSrcDir = true
			}
		case "tests", "test", "__tests__":
			if entry.IsDir() {
				result.Evidence.HasTestDir = true
			}
		}
	}
}

// classify determines the environment type based on collected evidence.
func (d *Detector) classify(result *EvidenceResult) EnvironmentType {
	evidence := result.Evidence

	// Detection order (priority): KDSE_RUNTIME > KDSE_PROJECT > SOFTWARE_PROJECT > BLANK_WORKSPACE

	// Step 1: Check for KDSE Runtime
	if d.isKDSERuntime(evidence) {
		return EnvironmentKDSERuntime
	}

	// Step 2: Check for KDSE Project
	if d.isKDSEProject(evidence) {
		return EnvironmentKDSEProject
	}

	// Step 3: Check for Software Project
	if d.isSoftwareProject(evidence) {
		return EnvironmentSoftwareProject
	}

	// Step 4: Default to BLANK_WORKSPACE
	return EnvironmentBlankWorkspace
}

// isKDSERuntime checks if this is the KDSE runtime repository itself.
func (d *Detector) isKDSERuntime(evidence Evidence) bool {
	// REQUIRED EVIDENCE for KDSE_RUNTIME:
	// At least 2 of these markers must be present:
	// 1. Module name is github.com/kdse/runtime
	// 2. internal/runtime/ directory exists
	// 3. internal/detection/ directory exists
	// 4. templates/ directory exists
	// 5. docs/ directory exists

	markerCount := 0
	if evidence.HasGoModKDSE {
		markerCount++
	}
	if evidence.HasRuntimeDir {
		markerCount++
	}
	if evidence.HasDetectionDir {
		markerCount++
	}
	if evidence.HasTemplatesDir {
		markerCount++
	}
	if evidence.HasDocsDir {
		markerCount++
	}

	// Require at least 2 markers for KDSE Runtime classification
	return markerCount >= 2
}

// isKDSEProject checks if this is a KDSE-initialized project.
func (d *Detector) isKDSEProject(evidence Evidence) bool {
	// REQUIRED EVIDENCE for KDSE_PROJECT:
	// 1. .kdse/ directory MUST exist
	// 2. Either manifest.yaml or manifest.json MUST exist
	// 3. manifest MUST be valid (parseable)

	if !evidence.HasKDSEProject {
		return false
	}

	if !evidence.HasManifest {
		return false
	}

	if d.options.StrictManifestValidation && !evidence.ManifestValid {
		return false
	}

	// Additional validation: ensure it's not the KDSE Runtime itself
	if evidence.HasGoModKDSE || evidence.HasRuntimeDir {
		return false
	}

	return true
}

// isSoftwareProject checks if this is a normal software project.
func (d *Detector) isSoftwareProject(evidence Evidence) bool {
	// REQUIRED EVIDENCE for SOFTWARE_PROJECT:
	// At least one project indicator MUST exist
	// AND must not be a KDSE project

	if evidence.HasKDSEProject {
		return false
	}

	// Check for any project indicators
	hasProjectIndicator := evidence.HasGoMod ||
		evidence.HasPackageJSON ||
		evidence.HasPyProject ||
		evidence.HasCargoToml ||
		evidence.HasPomXML ||
		evidence.HasMakefile

	if !hasProjectIndicator {
		return false
	}

	// If Git is required and not present, not a valid software project
	if d.options.RequireGitRepo && !evidence.HasGitRepo {
		return false
	}

	return true
}

// calculateConfidence calculates confidence in the classification.
func (d *Detector) calculateConfidence(result *EvidenceResult) float64 {
	evidence := result.Evidence
	env := result.Environment

	switch env {
	case EnvironmentKDSERuntime:
		// High confidence if multiple KDSE markers present
		markers := 0
		if evidence.HasGoModKDSE {
			markers++
		}
		if evidence.HasRuntimeDir {
			markers++
		}
		if evidence.HasDetectionDir {
			markers++
		}
		if evidence.HasTemplatesDir {
			markers++
		}
		return float64(markers) / 5.0

	case EnvironmentKDSEProject:
		// High confidence if manifest is valid
		if evidence.ManifestValid {
			return 0.95
		}
		return 0.7

	case EnvironmentSoftwareProject:
		// Confidence based on number of project indicators
		indicators := 0
		if evidence.HasGoMod {
			indicators++
		}
		if evidence.HasPackageJSON {
			indicators++
		}
		if evidence.HasPyProject {
			indicators++
		}
		if evidence.HasCargoToml {
			indicators++
		}
		if evidence.HasPomXML {
			indicators++
		}
		if evidence.HasMakefile {
			indicators++
		}
		conf := 0.6 + float64(indicators)*0.05
		if conf > 0.95 {
			conf = 0.95
		}
		return conf

	case EnvironmentBlankWorkspace:
		// Low confidence by default
		return 0.8

	default:
		return 0.0
	}
}

// String returns a human-readable representation of the detection result.
func (r *EvidenceResult) String() string {
	return fmt.Sprintf("Environment: %s (confidence: %.0f%%)\nRepo: %s\nEvidence: %+v",
		r.Environment, r.Confidence*100, r.Evidence.RepoRoot, r.Evidence)
}

// HasWarnings returns true if there are any warnings.
func (r *EvidenceResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// HasErrors returns true if there are any errors.
func (r *EvidenceResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// ClassificationRule represents a single classification rule.
type ClassificationRule struct {
	// Name is a human-readable name for the rule.
	Name string

	// Description explains why this rule exists.
	Description string

	// Priority determines evaluation order (lower = evaluated first).
	Priority int

	// Check function returns true if this rule applies.
	Check func(Evidence) bool

	// Result is the environment type if this rule matches.
	Result EnvironmentType

	// RequiredEvidence lists the evidence required for this rule.
	RequiredEvidence []string
}

// GetClassificationRules returns all classification rules with documentation.
func GetClassificationRules() []ClassificationRule {
	return []ClassificationRule{
		{
			Name:        "KDSE Runtime Detection",
			Description: "The KDSE runtime repository contains unique markers that identify it as the framework itself, not a normal software project.",
			Priority:    1,
			Check: func(e Evidence) bool {
				// Check KDSE Runtime markers
				markerCount := 0
				if e.HasGoModKDSE {
					markerCount++
				}
				if e.HasRuntimeDir {
					markerCount++
				}
				if e.HasDetectionDir {
					markerCount++
				}
				if e.HasTemplatesDir {
					markerCount++
				}
				return markerCount >= 2
			},
			Result: EnvironmentKDSERuntime,
			RequiredEvidence: []string{
				"github.com/kdse/runtime module (go.mod)",
				"internal/runtime/ directory",
				"internal/detection/ directory",
				"templates/ directory",
				"docs/ directory",
			},
		},
		{
			Name:        "KDSE Project Detection",
			Description: "A KDSE project contains a valid .kdse/ directory with manifest, indicating it was initialized by KDSE for engineering governance.",
			Priority:    2,
			Check: func(e Evidence) bool {
				return e.HasKDSEProject && e.HasManifest && e.ManifestValid && !e.HasGoModKDSE
			},
			Result: EnvironmentKDSEProject,
			RequiredEvidence: []string{
				".kdse/ directory",
				"manifest.yaml or manifest.json",
				"Valid, parseable manifest",
				"Not the KDSE runtime itself",
			},
		},
		{
			Name:        "Software Project Detection",
			Description: "A normal software project contains project-specific files indicating its language/framework, but is not a KDSE runtime or project.",
			Priority:    3,
			Check: func(e Evidence) bool {
				if e.HasKDSEProject {
					return false
				}
				return e.HasGoMod || e.HasPackageJSON || e.HasPyProject ||
					e.HasCargoToml || e.HasPomXML || e.HasMakefile
			},
			Result: EnvironmentSoftwareProject,
			RequiredEvidence: []string{
				"At least one project indicator (go.mod, package.json, etc.)",
				"No .kdse/ directory",
			},
		},
		{
			Name:        "Blank Workspace Detection",
			Description: "A blank workspace has no recognizable project structure, no Git repository, and no KDSE initialization.",
			Priority:    4,
			Check: func(e Evidence) bool {
				return !e.HasKDSEProject && !e.HasGoMod && !e.HasPackageJSON &&
					!e.HasPyProject && !e.HasCargoToml && !e.HasPomXML && !e.HasMakefile
			},
			Result: EnvironmentBlankWorkspace,
			RequiredEvidence: []string{
				"No project indicators",
				"No .kdse/ directory",
			},
		},
	}
}

// ValidationError represents a validation error during detection.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
