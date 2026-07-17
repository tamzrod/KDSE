// Package bootstrap provides GitHub-based template bootstrapping for KDSE.
// It downloads official workspace templates from GitHub and initializes
// the .kdse/ workspace directory structure.
package bootstrap

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kdse/runtime/internal/workspace"
)

// Bootstrapper handles the initialization of KDSE workspaces from GitHub templates.
type Bootstrapper struct {
	repoPath      string
	templateName  string
	templateOwner string
	templateRepo  string
	templateRef   string
	httpClient    *http.Client
	workspace     *workspace.Workspace
}

// Config holds bootstrap configuration options.
type Config struct {
	// TemplateName is the name of the template to use (e.g., "core", "web").
	TemplateName string
	// TemplateOwner is the GitHub owner/organization (default: "kdse").
	TemplateOwner string
	// TemplateRepo is the GitHub repository name (default: "workspace-templates").
	TemplateRepo string
	// TemplateRef is the git ref to download (branch, tag, or commit).
	TemplateRef string
}

// DefaultConfig returns the default bootstrap configuration.
func DefaultConfig() *Config {
	return &Config{
		TemplateName: "core",
		TemplateOwner: "kdse",
		TemplateRepo:  "workspace-templates",
		TemplateRef:   "main",
	}
}

// Result contains the outcome of a bootstrap operation.
type Result struct {
	Success        bool
	KDSEVersion    string
	TemplateVersion string
	TemplateCommit  string
	WorkspacePath   string
	Errors         []error
	CreatedPaths   []string
}

// NewBootstrapper creates a new bootstrapper for the given repository path.
func NewBootstrapper(repoPath string, cfg *Config) *Bootstrapper {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Bootstrapper{
		repoPath:      repoPath,
		templateName:  cfg.TemplateName,
		templateOwner: cfg.TemplateOwner,
		templateRepo:  cfg.TemplateRepo,
		templateRef:   cfg.TemplateRef,
		httpClient:    &http.Client{},
		workspace:     workspace.New(repoPath),
	}
}

// Initialize performs the full bootstrap process:
// 1. Create temporary directory
// 2. Download template from GitHub
// 3. Extract template
// 4. Copy to workspace
// 5. Write session metadata
// 6. Verify installation
// 7. Clean up temporary files
func (b *Bootstrapper) Initialize() *Result {
	result := &Result{
		Success:      true,
		WorkspacePath: filepath.Join(b.repoPath, ".kdse"),
	}

	// Create temporary directory for download
	tmpDir, err := os.MkdirTemp("", "kdse-bootstrap-*")
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Errorf("failed to create temp directory: %w", err))
		return result
	}
	defer os.RemoveAll(tmpDir)

	// Track created paths for rollback
	var createdPaths []string
	createdPaths = append(createdPaths, b.workspace.Root())

	// Step 1: Download template
	archivePath := filepath.Join(tmpDir, "template.tar.gz")
	commit, err := b.downloadTemplate(archivePath)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Errorf("download failed: %w", err))
		b.rollback(createdPaths)
		return result
	}
	result.TemplateCommit = commit

	// Step 2: Extract template
	extractDir := filepath.Join(tmpDir, "extracted")
	if err := b.extractTemplate(archivePath, extractDir); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Errorf("extraction failed: %w", err))
		b.rollback(createdPaths)
		return result
	}

	// Step 3: Copy template to workspace
	if err := b.copyTemplate(extractDir, &createdPaths); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Errorf("copy failed: %w", err))
		b.rollback(createdPaths)
		return result
	}

	// Step 4: Write session metadata
	sessionPath := filepath.Join(b.workspace.Root(), "session.yaml")
	kdseVersion := b.getKDSEVersion()
	if err := b.writeSessionMetadata(sessionPath, kdseVersion, result.TemplateCommit); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Errorf("session metadata write failed: %w", err))
		b.rollback(createdPaths)
		return result
	}
	createdPaths = append(createdPaths, sessionPath)

	result.KDSEVersion = kdseVersion
	result.TemplateVersion = b.templateName
	result.CreatedPaths = createdPaths

	// Step 5: Verify installation
	if err := b.verifyInstallation(); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Errorf("verification failed: %w", err))
		b.rollback(createdPaths)
		return result
	}

	return result
}

// downloadTemplate downloads the template archive from GitHub.
func (b *Bootstrapper) downloadTemplate(destPath string) (string, error) {
	// Construct the download URL using GitHub archive
	url := fmt.Sprintf(
		"https://github.com/%s/%s/archive/%s.tar.gz",
		b.templateOwner,
		b.templateRepo,
		b.templateRef,
	)

	resp, err := b.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub returned status %d: %s", resp.StatusCode, url)
	}

	// Check if we're redirected to get the actual commit
	commit := b.templateRef
	if resp.Request != nil && resp.Request.URL != nil {
		// Extract commit from final URL after redirects
		urlPath := resp.Request.URL.Path
		parts := strings.Split(urlPath, "/")
		if len(parts) >= 2 {
			commit = parts[len(parts)-1]
		}
	}

	// Save the archive
	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write archive: %w", err)
	}

	return commit, nil
}

// extractTemplate extracts the downloaded archive.
func (b *Bootstrapper) extractTemplate(archivePath, destDir string) error {
	// Create extraction directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}

	// Use tar command to extract
	return extractTarGz(archivePath, destDir)
}

// copyTemplate copies the extracted template to the workspace.
func (b *Bootstrapper) copyTemplate(extractDir string, createdPaths *[]string) error {
	// Find the template directory inside the extracted archive
	// The archive structure is typically: repo-ref/x/y/z files
	templateBase := filepath.Join(extractDir, b.templateRepo+"-"+b.templateRef)
	
	// If template is in a subdirectory, look for it
	templateDir := filepath.Join(templateBase, "templates", b.templateName)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		// Try without template subdirectory
		templateDir = filepath.Join(templateBase, b.templateName)
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			// Use the extracted base directly (archive contains template files at root)
			templateDir = templateBase
		}
	}

	// Copy everything from template directory to workspace root
	return copyDirContents(templateDir, b.workspace.Root(), createdPaths)
}

// verifyInstallation checks that all required workspace artifacts exist.
func (b *Bootstrapper) verifyInstallation() error {
	requiredDirs := []string{
		"knowledge",
		"architecture",
		"implementation",
		"verification",
		"reports",
		"docs",
	}

	var missing []string
	for _, dir := range requiredDirs {
		path := filepath.Join(b.workspace.Root(), dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, dir)
		}
	}

	// Also verify .kdse directory itself
	if _, err := os.Stat(b.workspace.Root()); os.IsNotExist(err) {
		missing = append(missing, ".kdse")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required artifacts: %s", strings.Join(missing, ", "))
	}

	return nil
}

// rollback removes all created paths on error.
func (b *Bootstrapper) rollback(createdPaths []string) {
	for _, path := range createdPaths {
		os.RemoveAll(path)
	}
}

// writeSessionMetadata creates the session.yaml file with initialization info.
func (b *Bootstrapper) writeSessionMetadata(sessionPath, kdseVersion, templateCommit string) error {
	content := fmt.Sprintf(`# KDSE Session Metadata
# This file is automatically generated - do not edit manually

kdse_version: %s
template_version: %s
template_commit: %s
initialized: %s
`,
		kdseVersion,
		b.templateName,
		templateCommit,
		getCurrentTimestamp(),
	)

	return os.WriteFile(sessionPath, []byte(content), 0644)
}

// getKDSEVersion returns the current KDSE version.
func (b *Bootstrapper) getKDSEVersion() string {
	return "2.0"
}

// copyDirContents recursively copies all contents from src to dest.
func copyDirContents(src, dest string, createdPaths *[]string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			*createdPaths = append(*createdPaths, destPath)
			if err := copyDirContents(srcPath, destPath, createdPaths); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", srcPath, err)
			}
			*createdPaths = append(*createdPaths, destPath)
		}
	}

	return nil
}

// copyFile copies a single file from src to dest.
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
