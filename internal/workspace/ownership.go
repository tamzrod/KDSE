package workspace

import (
	"path/filepath"
	"strings"
)

// OwnershipDomain represents the domain that owns an artifact
type OwnershipDomain string

const (
	// DomainProject represents the project layer (owned by the software project)
	DomainProject OwnershipDomain = "project"

	// DomainRuntime represents the KDSE runtime layer
	DomainRuntime OwnershipDomain = "runtime"

	// DomainReference represents the reference layer (external authorities)
	DomainReference OwnershipDomain = "reference"

	// DomainKnowledge represents the knowledge layer (extracted knowledge)
	DomainKnowledge OwnershipDomain = "knowledge"

	// DomainUnknown represents an unknown ownership domain
	DomainUnknown OwnershipDomain = "unknown"
)

// ArtifactClassifier classifies artifacts based on ownership rules
type ArtifactClassifier struct{}

// NewArtifactClassifier creates a new artifact classifier
func NewArtifactClassifier() *ArtifactClassifier {
	return &ArtifactClassifier{}
}

// Classify determines the ownership domain for an artifact based on its path
// and content classification hints.
func (c *ArtifactClassifier) Classify(path string) OwnershipDomain {
	path = filepath.Clean(path)

	// Rule 1: If artifact is in .kdse/ directory, determine subdomain
	if strings.HasPrefix(path, ".kdse/") || strings.HasPrefix(path, ".kdse\\") {
		return c.classifyRuntimeArtifact(path)
	}

	// Rule 2: If artifact is in project directories, it's project-owned
	if c.isProjectDirectory(path) {
		return DomainProject
	}

	// Rule 3: Check known project files at root
	if c.isProjectRootFile(path) {
		return DomainProject
	}

	// Default to project layer (conservative assumption)
	return DomainProject
}

// classifyRuntimeArtifact determines the subdomain for .kdse/ artifacts
func (c *ArtifactClassifier) classifyRuntimeArtifact(path string) OwnershipDomain {
	// Remove .kdse/ prefix for analysis
	relPath := strings.TrimPrefix(path, ".kdse/")
	relPath = strings.TrimPrefix(relPath, ".kdse\\")

	// Split path into components
	parts := strings.Split(filepath.ToSlash(relPath), "/")

	if len(parts) == 0 {
		return DomainRuntime
	}

	firstDir := parts[0]

	switch firstDir {
	case "references":
		return DomainReference
	case "knowledge":
		return DomainKnowledge
	default:
		// Everything else in .kdse/ is runtime layer
		return DomainRuntime
	}
}

// isProjectDirectory checks if the path is in a project-owned directory
func (c *ArtifactClassifier) isProjectDirectory(path string) bool {
	projectDirs := []string{
		"docs/",
		"docs\\",
		"src/",
		"src\\",
		"tests/",
		"tests\\",
		"cmd/",
		"cmd\\",
		"internal/",
		"internal\\",
		"deploy/",
		"deploy\\",
		"templates/",
		"templates\\",
		"examples/",
		"examples\\",
		".github/",
		".github\\",
	}

	for _, dir := range projectDirs {
		if strings.HasPrefix(path, dir) {
			return true
		}
	}

	return false
}

// isProjectRootFile checks if the file is a project-owned root file
func (c *ArtifactClassifier) isProjectRootFile(path string) bool {
	projectRootFiles := []string{
		"README.md",
		"LICENSE",
		"CHANGELOG.md",
		"go.mod",
		"go.sum",
		"package.json",
		"package-lock.json",
		"Cargo.toml",
		"Cargo.lock",
		"pyproject.toml",
		"requirements.txt",
		"Dockerfile",
		".gitignore",
		".gitattributes",
		"Makefile",
		"CMakeLists.txt",
		"setup.py",
		"setup.cfg",
		"pyproject.toml",
	}

	fileName := filepath.Base(path)
	for _, f := range projectRootFiles {
		if fileName == f {
			return true
		}
	}

	return false
}

// GetSuggestedPath returns the suggested path for an artifact based on its classification
// This helps enforce ownership boundaries by suggesting correct locations
func (c *ArtifactClassifier) GetSuggestedPath(artifactName string, domain OwnershipDomain) string {
	switch domain {
	case DomainProject:
		return filepath.Join("docs", artifactName)
	case DomainRuntime:
		return filepath.Join(".kdse", "runtime", artifactName)
	case DomainReference:
		return filepath.Join(".kdse", "references", artifactName)
	case DomainKnowledge:
		return filepath.Join(".kdse", "knowledge", artifactName)
	default:
		return artifactName
	}
}

// IsValidOwnership determines if placing an artifact at the given path is valid
// based on the artifact's classification
func (c *ArtifactClassifier) IsValidOwnership(artifactName string, targetPath string) (bool, OwnershipDomain, string) {
	domain := c.Classify(artifactName)

	// Determine expected base directory for the domain
	var expectedBase string
	switch domain {
	case DomainProject:
		expectedBase = ""
	case DomainRuntime:
		expectedBase = ".kdse"
	case DomainReference:
		expectedBase = ".kdse/references"
	case DomainKnowledge:
		expectedBase = ".kdse/knowledge"
	}

	// Check if the target path is valid for the domain
	targetDir := filepath.Dir(targetPath)
	if strings.HasPrefix(targetDir, expectedBase) || targetDir == expectedBase || expectedBase == "" {
		return true, domain, ""
	}

	// Return suggestion for correct path
	suggested := c.GetSuggestedPath(filepath.Base(targetPath), domain)
	return false, domain, suggested
}

// OwnershipViolation represents a violation of ownership boundaries
type OwnershipViolation struct {
	Artifact    string
	ActualPath  string
	ExpectedDomain OwnershipDomain
	SuggestedPath string
	Message      string
}

// ValidateArtifactPlacement checks if an artifact is in the correct ownership domain
func (c *ArtifactClassifier) ValidateArtifactPlacement(artifactPath string, artifactName string) *OwnershipViolation {
	domain := c.Classify(artifactName)

	// Determine expected base path for the domain
	var expectedBase string
	switch domain {
	case DomainProject:
		expectedBase = ""
	case DomainRuntime:
		expectedBase = ".kdse"
	case DomainReference:
		expectedBase = ".kdse/references"
	case DomainKnowledge:
		expectedBase = ".kdse/knowledge"
	}

	// Check if artifact is in the correct location
	pathDir := filepath.Dir(artifactPath)
	isValid := strings.HasPrefix(pathDir, expectedBase) || pathDir == expectedBase || expectedBase == ""

	if !isValid {
		suggested := c.GetSuggestedPath(artifactName, domain)
		return &OwnershipViolation{
			Artifact:       artifactName,
			ActualPath:     artifactPath,
			ExpectedDomain: domain,
			SuggestedPath: suggested,
			Message:       "Artifact placed in wrong ownership domain",
		}
	}

	return nil
}

// ProjectOwnedPaths returns all paths that are owned by the project layer
func ProjectOwnedPaths() []string {
	return []string{
		"README.md",
		"LICENSE",
		"CHANGELOG.md",
		"go.mod",
		"go.sum",
		"Dockerfile",
		".gitignore",
		".github/",
		"docs/",
		"src/",
		"tests/",
		"cmd/",
		"internal/",
		"deploy/",
		"templates/",
		"examples/",
	}
}

// RuntimeOwnedPaths returns all paths that are owned by the runtime layer
func RuntimeOwnedPaths() []string {
	return []string{
		".kdse/runtime/",
		".kdse/sessions/",
		".kdse/state/",
		".kdse/cache/",
		".kdse/reports/",
		".kdse/evidence/",
		".kdse/traceability/",
		".kdse/laboratory/",
	}
}

// ReferencePaths returns all paths that are owned by the reference layer
func ReferencePaths() []string {
	return []string{
		".kdse/references/",
	}
}

// KnowledgePaths returns all paths that are owned by the knowledge layer
func KnowledgePaths() []string {
	return []string{
		".kdse/knowledge/",
	}
}

// IsProjectOwned checks if a path is owned by the project layer
func IsProjectOwned(path string) bool {
	for _, p := range ProjectOwnedPaths() {
		if strings.HasPrefix(path, p) || path == p {
			return true
		}
	}
	return false
}

// IsRuntimeOwned checks if a path is owned by the runtime layer
func IsRuntimeOwned(path string) bool {
	for _, p := range RuntimeOwnedPaths() {
		if strings.HasPrefix(path, p) || path == p {
			return true
		}
	}
	return false
}

// IsReferenceOwned checks if a path is owned by the reference layer
func IsReferenceOwned(path string) bool {
	for _, p := range ReferencePaths() {
		if strings.HasPrefix(path, p) || path == p {
			return true
		}
	}
	return false
}

// IsKnowledgeOwned checks if a path is owned by the knowledge layer
func IsKnowledgeOwned(path string) bool {
	for _, p := range KnowledgePaths() {
		if strings.HasPrefix(path, p) || path == p {
			return true
		}
	}
	return false
}
