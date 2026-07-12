package normalize

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DiscoveryConfig configures what documentation to discover
type DiscoveryConfig struct {
	IncludeHidden   bool
	MaxFileSize     int64
	IncludePatterns []string
	ExcludePatterns []string
}

// DefaultDiscoveryConfig returns the default discovery configuration
func DefaultDiscoveryConfig() *DiscoveryConfig {
	return &DiscoveryConfig{
		IncludeHidden:   false,
		MaxFileSize:     10 * 1024 * 1024, // 10MB
		IncludePatterns: []string{"*.md", "*.txt", "*.rst", "*.adoc"},
		ExcludePatterns: []string{
			"node_modules",
			".git",
			"vendor",
			"__pycache__",
			".venv",
			"dist",
			"build",
		},
	}
}

// Discovery handles finding documentation in repositories
type Discovery struct {
	repoPath string
	config   *DiscoveryConfig
}

// NewDiscovery creates a new documentation discovery handler
func NewDiscovery(repoPath string) *Discovery {
	return &Discovery{
		repoPath: repoPath,
		config:   DefaultDiscoveryConfig(),
	}
}

// NewDiscoveryWithConfig creates a discovery with custom config
func NewDiscoveryWithConfig(repoPath string, config *DiscoveryConfig) *Discovery {
	return &Discovery{
		repoPath: repoPath,
		config:   config,
	}
}

// Discover finds all documentation in the repository
func (d *Discovery) Discover() ([]DiscoveredDocument, error) {
	var documents []DiscoveredDocument

	err := filepath.Walk(d.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip excluded directories
		if info.IsDir() {
			if d.shouldExcludeDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip non-documentation files
		if !d.isDocumentationFile(path) {
			return nil
		}

		// Skip hidden files if configured
		if !d.config.IncludeHidden && d.isHiddenFile(path) {
			return nil
		}

		// Skip large files
		if info.Size() > d.config.MaxFileSize {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(d.repoPath, path)
		doc := DiscoveredDocument{
			Path:         relPath,
			Type:         d.classifyDocument(relPath, string(content)),
			Title:        d.extractTitle(string(content)),
			Size:         info.Size(),
			LastModified: info.ModTime().Format(time.RFC3339),
			Encoding:     detectEncoding(content),
			Confidence:   d.calculateConfidence(relPath, string(content)),
		}

		documents = append(documents, doc)
		return nil
	})

	return documents, err
}

// DiscoverByType finds documentation of specific types
func (d *Discovery) DiscoverByType(targetTypes ...DocType) ([]DiscoveredDocument, error) {
	all, err := d.Discover()
	if err != nil {
		return nil, err
	}

	var filtered []DiscoveredDocument
	targetMap := make(map[DocType]bool)
	for _, t := range targetTypes {
		targetMap[t] = true
	}

	for _, doc := range all {
		if targetMap[doc.Type] {
			filtered = append(filtered, doc)
		}
	}

	return filtered, nil
}

// DiscoverArchitectureDocs finds architecture-related documentation
func (d *Discovery) DiscoverArchitectureDocs() ([]DiscoveredDocument, error) {
	return d.DiscoverByType(
		DocsTypeArchitecture,
		DocsTypeDesign,
		DocsTypeADR,
		DocsTypeEngineering,
	)
}

// DiscoverKnowledgeDocs finds knowledge/artifacts documentation
func (d *Discovery) DiscoverKnowledgeDocs() ([]DiscoveredDocument, error) {
	return d.DiscoverByType(
		DocsTypeREADME,
		DocsTypeGuide,
		DocsTypeTutorial,
		DocsTypeReference,
	)
}

func (d *Discovery) shouldExcludeDir(name string) bool {
	for _, pattern := range d.config.ExcludePatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}

func (d *Discovery) isDocumentationFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	docExtensions := map[string]bool{
		".md":  true,
		".txt": true,
		".rst": true,
		".adoc": true,
	}

	// Also check for specific filenames
	baseName := filepath.Base(path)
	specialFiles := map[string]bool{
		"README":          true,
		"LICENSE":         true,
		"CHANGELOG":       true,
		"CONTRIBUTING":    true,
		"ARCHITECTURE":    true,
		"DESIGN":          true,
		"TODO":            true,
		"TROUBLESHOOTING": true,
	}

	return docExtensions[ext] || specialFiles[strings.ToUpper(baseName)]
}

func (d *Discovery) isHiddenFile(path string) bool {
	relPath, _ := filepath.Rel(d.repoPath, path)
	parts := strings.Split(relPath, string(filepath.Separator))
	for _, part := range parts {
		if strings.HasPrefix(part, ".") && part != "." {
			return true
		}
	}
	return false
}

func (d *Discovery) classifyDocument(path, content string) DocType {
	baseName := strings.ToUpper(filepath.Base(path))

	// Check specific filenames first
	switch {
	case strings.HasPrefix(baseName, "README"):
		return DocsTypeREADME
	case strings.HasPrefix(baseName, "ARCHITECTURE"):
		return DocsTypeArchitecture
	case strings.HasPrefix(baseName, "DESIGN"):
		return DocsTypeDesign
	case strings.HasPrefix(baseName, "CONTRIBUTING"):
		return DocsTypeContributing
	case strings.HasPrefix(baseName, "LICENSE"):
		return DocsTypeLicense
	case strings.HasPrefix(baseName, "CHANGELOG"):
		return DocsTypeChangelog
	case strings.HasPrefix(baseName, "ADR-") || strings.HasPrefix(baseName, "ADR_"):
		return DocsTypeADR
	case strings.HasPrefix(baseName, "API"):
		return DocsTypeAPI
	case strings.HasPrefix(baseName, "GUIDE"):
		return DocsTypeGuide
	case strings.HasPrefix(baseName, "TUTORIAL"):
		return DocsTypeTutorial
	case strings.HasPrefix(baseName, "EXAMPLE"):
		return DocsTypeExamples
	}

	// Check path-based classification
	lowerPath := strings.ToLower(path)
	switch {
	case strings.Contains(lowerPath, "/architecture/"):
		return DocsTypeArchitecture
	case strings.Contains(lowerPath, "/design/"):
		return DocsTypeDesign
	case strings.Contains(lowerPath, "/docs/"):
		if strings.Contains(lowerPath, "api") {
			return DocsTypeAPI
		}
		return DocsTypeReference
	case strings.Contains(lowerPath, "/wiki/"):
		return DocsTypeWiki
	case strings.Contains(lowerPath, "/examples/"):
		return DocsTypeExamples
	case strings.Contains(lowerPath, "/adr/") || strings.Contains(lowerPath, "/decisions/"):
		return DocsTypeADR
	case strings.Contains(lowerPath, "/engineering/"):
		return DocsTypeEngineering
	}

	// Content-based classification
	return d.classifyByContent(content)
}

func (d *Discovery) classifyByContent(content string) DocType {
	upper := strings.ToUpper(content)

	// Look for content patterns
	if strings.Contains(upper, "API") && (strings.Contains(upper, "ENDPOINT") || strings.Contains(upper, "REST") || strings.Contains(upper, "GRPC")) {
		return DocsTypeAPI
	}
	if strings.Contains(upper, "ARCHITECTURE") || strings.Contains(upper, "COMPONENT") || strings.Contains(upper, "SYSTEM DESIGN") {
		return DocsTypeArchitecture
	}
	if strings.Contains(upper, "TUTORIAL") || strings.Contains(upper, "GETTING STARTED") || strings.Contains(upper, "QUICK START") {
		return DocsTypeTutorial
	}
	if strings.Contains(upper, "HOW TO") || strings.Contains(upper, "GUIDE") {
		return DocsTypeGuide
	}
	if strings.Contains(upper, "CHANGE LOG") || strings.Contains(upper, "RELEASE NOTE") {
		return DocsTypeChangelog
	}

	return DocsTypeReference
}

func (d *Discovery) extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Check for markdown heading
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
		// Check for asciidoc heading
		if strings.HasPrefix(line, "= ") {
			return strings.TrimPrefix(line, "= ")
		}
		// Check for rst heading (underlined)
		// This would require more context, skip for now
	}
	return ""
}

func (d *Discovery) calculateConfidence(path, content string) float64 {
	baseName := strings.ToUpper(filepath.Base(path))
	lowerPath := strings.ToLower(path)

	confidence := 0.5

	// High confidence for known patterns
	switch {
	case strings.HasPrefix(baseName, "README"):
		confidence = 0.95
	case strings.HasPrefix(baseName, "ARCHITECTURE"):
		confidence = 0.95
	case strings.HasPrefix(baseName, "DESIGN"):
		confidence = 0.95
	case strings.HasPrefix(baseName, "ADR-") || strings.HasPrefix(baseName, "ADR_"):
		confidence = 0.95
	case strings.Contains(lowerPath, "/architecture/"):
		confidence = 0.90
	case strings.Contains(lowerPath, "/design/"):
		confidence = 0.90
	}

	// Boost for substantial content
	if len(content) > 1000 {
		confidence += 0.05
	}
	if len(content) > 5000 {
		confidence += 0.05
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

func detectEncoding(content []byte) string {
	// Simple detection - in production would use proper encoding detection
	if isUTF8(content) {
		return "UTF-8"
	}
	return "Unknown"
}

func isUTF8(content []byte) bool {
	// Basic UTF-8 validation
	for i := 0; i < len(content); {
		r, size := decodeRune(content[i:])
		if r == 0xFFFD && size == 1 {
			return false
		}
		i += size
	}
	return true
}

func decodeRune(b []byte) (rune, int) {
	if len(b) == 0 {
		return 0, 0
	}
	switch {
	case b[0]&0x80 == 0:
		return rune(b[0]), 1
	case b[0]&0xE0 == 0xC0 && len(b) >= 2:
		return rune(b[0]&0x1F)<<6 | rune(b[1]&0x3F), 2
	case b[0]&0xF0 == 0xE0 && len(b) >= 3:
		return rune(b[0]&0x0F)<<12 | rune(b[1]&0x3F)<<6 | rune(b[2]&0x3F), 3
	case b[0]&0xF8 == 0xF0 && len(b) >= 4:
		return rune(b[0]&0x07)<<18 | rune(b[1]&0x3F)<<12 | rune(b[2]&0x3F)<<6 | rune(b[3]&0x3F), 4
	default:
		return 0xFFFD, 1
	}
}
