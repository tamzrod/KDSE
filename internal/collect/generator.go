package collect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArtifactGenerator generates KDSE-standard knowledge artifacts
type ArtifactGenerator struct {
	repoPath string
	sessionID string
}

// NewArtifactGenerator creates a new artifact generator
func NewArtifactGenerator(repoPath, sessionID string) *ArtifactGenerator {
	return &ArtifactGenerator{
		repoPath:  repoPath,
		sessionID: sessionID,
	}
}

// DomainPaths maps knowledge domains to their directory paths
var DomainPaths = map[KnowledgeDomain]string{
	DomainPhysics:      "knowledge/physics",
	DomainEquipment:    "knowledge/equipment",
	DomainEnvironment:  "knowledge/environment",
	DomainStandards:    "knowledge/standards",
	DomainBusiness:     "knowledge/business",
	DomainSimulation:   "knowledge/simulation",
	DomainControl:      "knowledge/control",
	DomainProtocols:    "knowledge/protocols",
	DomainVocabulary:   "knowledge/vocabulary",
	DomainTransformers: "knowledge/equipment/transformers",
	DomainBattery:      "knowledge/equipment/battery",
	DomainRelay:        "knowledge/equipment/relay",
	DomainWeather:      "knowledge/environment/weather",
	DomainGeneral:      "knowledge/general",
}

// GenerateArtifacts creates KDSE-standard artifacts from collected knowledge
func (g *ArtifactGenerator) GenerateArtifacts(artifacts []CollectedArtifact) ([]CollectedArtifact, error) {
	var generated []CollectedArtifact

	// Ensure knowledge base directory exists
	kbDir := filepath.Join(g.repoPath, ".kdse", "knowledge")
	if err := os.MkdirAll(kbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create knowledge base directory: %w", err)
	}

	// Group artifacts by domain
	domainGroups := g.groupByDomain(artifacts)

	// Generate domain directories and artifacts
	for domain, domainArtifacts := range domainGroups {
		dirPath := g.getDomainPath(domain)
		if err := os.MkdirAll(filepath.Join(g.repoPath, ".kdse", dirPath), 0755); err != nil {
			continue
		}

		// Generate artifact files
		for _, artifact := range domainArtifacts {
			path, err := g.writeArtifact(artifact, dirPath)
			if err == nil {
				artifact.Path = path
				generated = append(generated, artifact)
			}
		}

		// Generate domain index
		g.generateDomainIndex(domain, domainArtifacts, dirPath)
	}

	// Generate master knowledge index
	g.generateMasterIndex(generated)

	return generated, nil
}

// groupByDomain groups artifacts by their knowledge domain
func (g *ArtifactGenerator) groupByDomain(artifacts []CollectedArtifact) map[KnowledgeDomain][]CollectedArtifact {
	groups := make(map[KnowledgeDomain][]CollectedArtifact)
	for _, artifact := range artifacts {
		groups[artifact.Domain] = append(groups[artifact.Domain], artifact)
	}
	return groups
}

// getDomainPath returns the directory path for a domain
func (g *ArtifactGenerator) getDomainPath(domain KnowledgeDomain) string {
	if path, ok := DomainPaths[domain]; ok {
		return path
	}
	return "knowledge/general"
}

// writeArtifact writes a single artifact to disk
func (g *ArtifactGenerator) writeArtifact(artifact CollectedArtifact, dirPath string) (string, error) {
	content := g.formatArtifact(artifact)
	
	// Create filename from title
	filename := g.sanitizeFilename(artifact.Title) + ".md"
	filePath := filepath.Join(g.repoPath, ".kdse", dirPath, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

// formatArtifact formats an artifact as KDSE-standard markdown
func (g *ArtifactGenerator) formatArtifact(artifact CollectedArtifact) string {
	var b strings.Builder

	b.WriteString("# " + artifact.Title + "\n\n")
	b.WriteString("**Artifact ID:** " + artifact.ID + "\n")
	b.WriteString("**Domain:** " + string(artifact.Domain) + "\n")
	b.WriteString("**Source:** " + string(artifact.Source) + "\n")
	if artifact.SourcePath != "" {
		b.WriteString("**Source Path:** " + artifact.SourcePath + "\n")
	}
	b.WriteString("**Authority Level:** " + string(artifact.Authority) + "\n")
	b.WriteString("**Version:** " + kdseVersion + "\n")
	b.WriteString("**Collection Date:** " + artifact.CollectionDate + "\n")
	b.WriteString("**Collected By:** " + artifact.CollectedBy + "\n")
	b.WriteString("**Confidence Level:** " + fmt.Sprintf("%.0f%%", artifact.ConfidenceLevel*100) + "\n")

	if len(artifact.Tags) > 0 {
		b.WriteString("**Tags:** " + strings.Join(artifact.Tags, ", ") + "\n")
	}

	b.WriteString("\n---\n\n")

	// Summary section
	b.WriteString("## Summary\n\n")
	b.WriteString(artifact.Summary + "\n\n")

	// Content section
	b.WriteString("## Content\n\n")
	b.WriteString(artifact.Content + "\n\n")

	// Traceability section
	b.WriteString("## Traceability\n\n")
	b.WriteString("| Field | Value |\n")
	b.WriteString("|-------|-------|\n")
	b.WriteString(fmt.Sprintf("| Artifact ID | %s |\n", artifact.Traceability.ArtifactID))
	b.WriteString(fmt.Sprintf("| Traceability ID | %s |\n", artifact.Traceability.TraceabilityID))
	b.WriteString(fmt.Sprintf("| Source | %s |\n", artifact.Traceability.Source))
	b.WriteString(fmt.Sprintf("| Authority | %s |\n", artifact.Traceability.Authority))
	b.WriteString(fmt.Sprintf("| Version | %s |\n", artifact.Traceability.Version))
	b.WriteString(fmt.Sprintf("| Collection Date | %s |\n", artifact.Traceability.CollectionDate))
	b.WriteString(fmt.Sprintf("| Collected By | %s |\n", artifact.Traceability.CollectedBy))
	b.WriteString("| KDSE Version | " + kdseVersion + " |\n\n")

	if len(artifact.NormativeRefs) > 0 {
		b.WriteString("### Normative References\n\n")
		for _, ref := range artifact.NormativeRefs {
			b.WriteString("- " + ref + "\n")
		}
		b.WriteString("\n")
	}

	if len(artifact.Dependencies) > 0 {
		b.WriteString("### Dependencies\n\n")
		for _, dep := range artifact.Dependencies {
			b.WriteString("- " + dep + "\n")
		}
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("---\n\n")
	b.WriteString("*Generated by KDSE Collect on " + time.Now().Format(time.RFC3339) + "*\n")

	return b.String()
}

// sanitizeFilename creates a safe filename from a title
func (g *ArtifactGenerator) sanitizeFilename(title string) string {
	// Replace spaces with hyphens
	filename := strings.ReplaceAll(title, " ", "-")
	// Remove special characters
	filename = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, filename)
	// Limit length
	if len(filename) > 50 {
		filename = filename[:50]
	}
	return filename
}

// generateDomainIndex generates an index file for a domain
func (g *ArtifactGenerator) generateDomainIndex(domain KnowledgeDomain, artifacts []CollectedArtifact, dirPath string) error {
	var b strings.Builder

	b.WriteString("# " + strings.Title(string(domain)) + " Knowledge Domain\n\n")
	b.WriteString("**Domain:** " + string(domain) + "\n")
	b.WriteString("**Artifacts:** " + fmt.Sprintf("%d\n", len(artifacts)))
	b.WriteString("**Generated:** " + time.Now().Format(time.RFC3339) + "\n\n")

	b.WriteString("---\n\n")
	b.WriteString("## Overview\n\n")
	b.WriteString(fmt.Sprintf("This domain contains %d knowledge artifacts related to %s.\n\n", len(artifacts), domain))

	b.WriteString("## Artifacts\n\n")
	b.WriteString("| Artifact | Authority | Confidence | Source |\n")
	b.WriteString("|----------|----------|------------|--------|\n")

	for _, artifact := range artifacts {
		title := artifact.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		b.WriteString(fmt.Sprintf("| [%s](%s) | %s | %.0f%% | %s |\n",
			title,
			g.sanitizeFilename(artifact.Title)+".md",
			artifact.Authority,
			artifact.ConfidenceLevel*100,
			artifact.Source))
	}

	b.WriteString("\n---\n\n")
	b.WriteString("*Generated by KDSE Collect*\n")

	content := b.String()
	path := filepath.Join(g.repoPath, ".kdse", dirPath, "DOMAIN_INDEX.md")
	return os.WriteFile(path, []byte(content), 0644)
}

// generateMasterIndex generates the master knowledge base index
func (g *ArtifactGenerator) generateMasterIndex(artifacts []CollectedArtifact) error {
	var b strings.Builder

	b.WriteString("# KDSE Knowledge Base\n\n")
	b.WriteString("**Session ID:** " + g.sessionID + "\n")
	b.WriteString("**Generated:** " + time.Now().Format(time.RFC3339) + "\n")
	b.WriteString("**Total Artifacts:** " + fmt.Sprintf("%d\n", len(artifacts)) + "\n\n")

	b.WriteString("---\n\n")

	// Group by domain for the index
	domainGroups := make(map[KnowledgeDomain][]CollectedArtifact)
	for _, artifact := range artifacts {
		domainGroups[artifact.Domain] = append(domainGroups[artifact.Domain], artifact)
	}

	b.WriteString("## Knowledge Domains\n\n")
	b.WriteString("| Domain | Artifacts | Authority Distribution |\n")
	b.WriteString("|--------|-----------|----------------------|\n")

	for domain, domainArtifacts := range domainGroups {
		authorityCounts := g.countAuthorities(domainArtifacts)
		authStr := formatAuthorityCounts(authorityCounts)
		dirPath := g.getDomainPath(domain)
		b.WriteString(fmt.Sprintf("| [%s](%s/DOMAIN_INDEX.md) | %d | %s |\n",
			strings.Title(string(domain)), dirPath, len(domainArtifacts), authStr))
	}

	b.WriteString("\n---\n\n")

	// Summary statistics
	b.WriteString("## Collection Summary\n\n")
	b.WriteString("### Authority Levels\n\n")
	allAuthorityCounts := g.countAllAuthorities(artifacts)
	for auth, count := range allAuthorityCounts {
		b.WriteString(fmt.Sprintf("- %s: %d artifacts\n", auth, count))
	}

	b.WriteString("\n### Confidence Distribution\n\n")
	confDist := g.confidenceDistribution(artifacts)
	for level, count := range confDist {
		b.WriteString(fmt.Sprintf("- %s: %d artifacts\n", level, count))
	}

	b.WriteString("\n### Source Distribution\n\n")
	sourceCounts := g.countSources(artifacts)
	for source, count := range sourceCounts {
		b.WriteString(fmt.Sprintf("- %s: %d artifacts\n", source, count))
	}

	b.WriteString("\n---\n\n")
	b.WriteString("*Knowledge base generated by KDSE Collect v" + kdseVersion + "*\n")

	content := b.String()
	path := filepath.Join(g.repoPath, ".kdse", "knowledge", "INDEX.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func (g *ArtifactGenerator) countAuthorities(artifacts []CollectedArtifact) map[AuthorityLevel]int {
	counts := make(map[AuthorityLevel]int)
	for _, a := range artifacts {
		counts[a.Authority]++
	}
	return counts
}

func formatAuthorityCounts(counts map[AuthorityLevel]int) string {
	var parts []string
	for auth, count := range counts {
		parts = append(parts, fmt.Sprintf("%s:%d", auth, count))
	}
	return strings.Join(parts, ", ")
}

func (g *ArtifactGenerator) countAllAuthorities(artifacts []CollectedArtifact) map[AuthorityLevel]int {
	return g.countAuthorities(artifacts)
}

func (g *ArtifactGenerator) confidenceDistribution(artifacts []CollectedArtifact) map[string]int {
	dist := map[string]int{
		"High (>=90%)":   0,
		"Medium (70-89%)": 0,
		"Low (<70%)":     0,
	}

	for _, a := range artifacts {
		conf := a.ConfidenceLevel * 100
		if conf >= 90 {
			dist["High (>=90%)"]++
		} else if conf >= 70 {
			dist["Medium (70-89%)"]++
		} else {
			dist["Low (<70%)"]++
		}
	}

	return dist
}

func (g *ArtifactGenerator) countSources(artifacts []CollectedArtifact) map[KnowledgeSource]int {
	counts := make(map[KnowledgeSource]int)
	for _, a := range artifacts {
		counts[a.Source]++
	}
	return counts
}

// UpdateArtifact updates an existing artifact with new content
func (g *ArtifactGenerator) UpdateArtifact(existing, new CollectedArtifact) (CollectedArtifact, error) {
	updated := new
	updated.ID = existing.ID
	updated.CollectionDate = time.Now().Format(time.RFC3339)
	updated.Traceability.TraceabilityID = fmt.Sprintf("UPD-%s", time.Now().Format("20060102-150405"))
	updated.Traceability.CollectionDate = updated.CollectionDate

	if existing.Path != "" {
		if err := g.writeArtifact(updated, g.getDomainPath(updated.Domain)); err != nil {
			return updated, err
		}
		updated.Path = existing.Path
	}

	return updated, nil
}

// CheckExisting checks if an artifact already exists for the given domain and title
func (g *ArtifactGenerator) CheckExisting(domain KnowledgeDomain, title string) *CollectedArtifact {
	dirPath := g.getDomainPath(domain)
	checkPath := filepath.Join(g.repoPath, ".kdse", dirPath, g.sanitizeFilename(title)+".md")

	if _, err := os.Stat(checkPath); err == nil {
		// Found existing artifact
		content, err := os.ReadFile(checkPath)
		if err != nil {
			return nil
		}

		return &CollectedArtifact{
			ID:       GenerateArtifactID(domain),
			Domain:   domain,
			Title:    title,
			Path:     checkPath,
			Content:  string(content),
			Source:   SourceRepository,
			Authority: AuthorityProject,
		}
	}

	return nil
}
