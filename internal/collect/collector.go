package collect

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Collector discovers and catalogs engineering artifacts
type Collector struct {
	repoPath  string
	sessionID string
}

// NewCollector creates a new artifact collector
func NewCollector(repoPath, sessionID string) *Collector {
	return &Collector{
		repoPath:  repoPath,
		sessionID: sessionID,
	}
}

// Collect discovers artifacts in the artifacts/ directory and generates inventory
func (c *Collector) Collect() (*CollectionResult, error) {
	startTime := time.Now()
	result := NewCollectionResult(c.sessionID, c.repoPath)

	// Discover artifacts in artifacts/ directory
	artifactsDir := filepath.Join(c.repoPath, "artifacts")
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		// No artifacts directory - nothing to collect
		result.SetCompleted()
		result.ProcessingTime = time.Since(startTime).Seconds()
		return result, nil
	}

	// Walk the artifacts directory
	index := 0
	filepath.Walk(artifactsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		relPath, _ := filepath.Rel(c.repoPath, path)
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Calculate hash
		hash, err := CalculateHash(path)
		if err != nil {
			hash = "error"
		}

		// Create artifact record
		artifact := DiscoveredArtifact{
			ID:           GenerateArtifactID(c.repoPath, path, index),
			Path:         path,
			RelativePath: relPath,
			Name:         info.Name(),
			Category:     DetectCategory(path),
			Size:         info.Size(),
			Hash:         hash,
			Modified:     info.ModTime().Format(time.RFC3339),
			Extension:    filepath.Ext(path),
			CollectionID: c.sessionID,
		}

		result.AddArtifact(artifact)
		index++
		return nil
	})

	result.SetCompleted()
	result.ProcessingTime = time.Since(startTime).Seconds()

	// Save inventory
	if err := c.saveInventory(result); err != nil {
		return result, err
	}

	// Save report
	c.saveReport(result)

	return result, nil
}

// saveInventory saves the artifact inventory as JSON
func (c *Collector) saveInventory(result *CollectionResult) error {
	dir := filepath.Join(c.repoPath, ".kdse", "artifacts")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "inventory.json")
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// saveReport saves a human-readable collection report
func (c *Collector) saveReport(result *CollectionResult) {
	dir := filepath.Join(c.repoPath, ".kdse", "reports")
	os.MkdirAll(dir, 0755)

	path := filepath.Join(dir, fmt.Sprintf("artifact-collection-%s.md", c.sessionID))
	content := FormatReport(result)

	os.WriteFile(path, []byte(content), 0644)
}

// FormatReport generates a human-readable collection report
func FormatReport(result *CollectionResult) string {
	var b strings.Builder

	b.WriteString("# KDSE Artifact Collection Report\n\n")
	b.WriteString("| Field | Value |\n")
	b.WriteString("|-------|-------|\n")
	b.WriteString(fmt.Sprintf("| Session ID | %s |\n", result.SessionID))
	b.WriteString(fmt.Sprintf("| Repository | %s |\n", result.Repository))
	b.WriteString(fmt.Sprintf("| Collection Date | %s |\n", result.StartedAt))
	b.WriteString(fmt.Sprintf("| Report Version | 1.0 |\n\n")

	b.WriteString("## Summary\n\n")
	b.WriteString(fmt.Sprintf("**Artifacts Discovered:** %d\n", len(result.ArtifactsFound)))
	b.WriteString(fmt.Sprintf("**Total Size:** %s\n", formatSize(result.TotalSize)))
	b.WriteString(fmt.Sprintf("**Processing Time:** %.2f seconds\n\n", result.ProcessingTime))

	// Group by category
	categories := make(map[ArtifactCategory][]DiscoveredArtifact)
	for _, art := range result.ArtifactsFound {
		categories[art.Category] = append(categories[art.Category], art)
	}

	b.WriteString("## Artifacts by Category\n\n")
	b.WriteString("| Category | Count | Size |\n")
	b.WriteString("|----------|-------|------|\n")

	var categoryOrder = []ArtifactCategory{
		CategoryManual,
		CategoryStandard,
		CategorySpec,
		CategoryDatasheet,
		CategoryDrawing,
		CategoryDocument,
		CategoryImage,
		CategoryVideo,
		CategoryArchive,
		CategoryUnknown,
	}

	for _, cat := range categoryOrder {
		arts, ok := categories[cat]
		if !ok || len(arts) == 0 {
			continue
		}

		var totalSize int64
		for _, a := range arts {
			totalSize += a.Size
		}
		b.WriteString(fmt.Sprintf("| %s | %d | %s |\n", cat, len(arts), formatSize(totalSize)))
	}
	b.WriteString("\n")

	// Artifact list
	b.WriteString("## Artifact Inventory\n\n")
	b.WriteString("| ID | Name | Category | Size | Hash |\n")
	b.WriteString("|----|------|----------|------|------|\n")

	for _, art := range result.ArtifactsFound {
		hash := art.Hash
		if len(hash) > 12 {
			hash = hash[:12] + "..."
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			art.ID, art.Name, art.Category, formatSize(art.Size), hash))
	}
	b.WriteString("\n")

	b.WriteString("---\n\n")
	b.WriteString(fmt.Sprintf("*Report generated by KDSE Collect v1.0*\n"))
	b.WriteString(fmt.Sprintf("*Generated: %s*\n", time.Now().Format(time.RFC3339)))

	return b.String()
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
