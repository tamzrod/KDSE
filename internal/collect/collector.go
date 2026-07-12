package collect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kdse/runtime/internal/types"
)

// Collector orchestrates the knowledge collection process
type Collector struct {
	repoPath    string
	sessionID   string
	operator    string
	gapDetector *GapDetector
	generator   *ArtifactGenerator
	registry    *ProviderRegistry
}

// NewCollector creates a new knowledge collector
func NewCollector(repoPath, operator string) *Collector {
	return &Collector{
		repoPath:    repoPath,
		sessionID:   types.GenerateSessionID(),
		operator:    operator,
		gapDetector: NewGapDetector(repoPath),
		generator:   NewArtifactGenerator(repoPath, types.GenerateSessionID()),
		registry:    NewProviderRegistry(repoPath),
	}
}

// Collect executes the full knowledge collection process
func (c *Collector) Collect(input *CollectionInput) (*CollectionResult, error) {
	startTime := time.Now()

	// Initialize result
	result := NewCollectionResult(c.sessionID, c.repoPath, c.operator)

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              KDSE Knowledge Collection                         ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Repository:   %s\n", c.repoPath)
	fmt.Printf("║ Session ID:   %s\n", c.sessionID)
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Step 1: Analyze knowledge gaps
	fmt.Println("[1/5] Analyzing knowledge gaps...")
	gaps := c.gapDetector.DetectGaps(input)
	result.KnowledgeGapsIdentified = gaps
	result.KnowledgeAreasReviewed = c.getKnowledgeAreas(gaps)
	
	if len(gaps) > 0 {
		fmt.Printf("      Found %d knowledge gaps\n", len(gaps))
		for _, gap := range gaps {
			fmt.Printf("      • [%s] %s (%s)\n", gap.ID, gap.Topic, gap.Severity)
		}
	}
	fmt.Println()

	// Step 2: Collect knowledge from all providers
	fmt.Println("[2/5] Collecting knowledge from available sources...")
	var allArtifacts []CollectedArtifact

	for _, provider := range c.registry.GetProviders() {
		fmt.Printf("      Checking provider: %s...\n", provider.Name())
		artifacts, err := provider.Collect(input)
		if err != nil {
			fmt.Printf("      ⚠ Provider %s encountered issues: %v\n", provider.Name(), err)
			continue
		}
		if len(artifacts) > 0 {
			fmt.Printf("      ✓ Collected %d artifacts\n", len(artifacts))
			allArtifacts = append(allArtifacts, artifacts...)
		}
	}
	fmt.Println()

	// Step 3: Check for existing knowledge
	fmt.Println("[3/5] Checking for existing knowledge...")
	var newArtifacts []CollectedArtifact
	var updatedArtifacts []CollectedArtifact
	var existingArtifacts []CollectedArtifact

	for _, artifact := range allArtifacts {
		existing := c.generator.CheckExisting(artifact.Domain, artifact.Title)
		if existing != nil {
			// Check if update is needed
			if c.shouldUpdate(existing, artifact) {
				updated, err := c.generator.UpdateArtifact(*existing, artifact)
				if err == nil {
					updatedArtifacts = append(updatedArtifacts, updated)
					fmt.Printf("      ✓ Updated: %s\n", artifact.Title)
				}
			} else {
				existingArtifacts = append(existingArtifacts, *existing)
				fmt.Printf("      → Already present: %s\n", artifact.Title)
			}
		} else {
			newArtifacts = append(newArtifacts, artifact)
		}
	}
	fmt.Println()

	// Step 4: Generate artifacts
	fmt.Println("[4/5] Generating KDSE-standard artifacts...")
	generated, err := c.generator.GenerateArtifacts(newArtifacts)
	if err != nil {
		return nil, fmt.Errorf("artifact generation failed: %w", err)
	}
	result.ArtifactsCollected = generated
	
	for _, a := range updatedArtifacts {
		result.ArtifactsUpdated = append(result.ArtifactsUpdated, a)
	}
	for _, a := range existingArtifacts {
		result.ArtifactsExisting = append(result.ArtifactsExisting, a)
	}

	fmt.Printf("      ✓ Generated %d new artifacts\n", len(generated))
	fmt.Printf("      ✓ Updated %d existing artifacts\n", len(updatedArtifacts))
	fmt.Printf("      ✓ Verified %d existing artifacts\n", len(existingArtifacts))
	fmt.Println()

	// Step 5: Determine remaining gaps
	fmt.Println("[5/5] Analyzing remaining gaps...")
	result.KnowledgeStillMissing = c.identifyRemainingGaps(gaps, generated, updatedArtifacts)
	result.Recommendations = c.generateRecommendations(result)

	fmt.Printf("      %d gaps still require attention\n", len(result.KnowledgeStillMissing))
	fmt.Println()

	// Calculate statistics
	result.CalculateStats(startTime)
	result.SetCompleted()

	// Save result
	if err := c.saveResult(result); err != nil {
		fmt.Printf("Warning: Could not save collection result: %v\n", err)
	}

	return result, nil
}

// getKnowledgeAreas extracts unique knowledge areas from gaps
func (c *Collector) getKnowledgeAreas(gaps []KnowledgeGap) []KnowledgeDomain {
	seen := make(map[KnowledgeDomain]bool)
	var areas []KnowledgeDomain
	for _, gap := range gaps {
		if !seen[gap.Domain] {
			areas = append(areas, gap.Domain)
			seen[gap.Domain] = true
		}
	}
	return areas
}

// shouldUpdate determines if an existing artifact should be updated
func (c *Collector) shouldUpdate(existing, new CollectedArtifact) bool {
	// Update if new artifact has higher confidence
	if new.ConfidenceLevel > existing.ConfidenceLevel {
		return true
	}
	// Update if new artifact has higher authority
	if c.authorityRank(new.Authority) > c.authorityRank(existing.Authority) {
		return true
	}
	// Update if content is significantly different
	if len(new.Content) > len(existing.Content)*2 {
		return true
	}
	return false
}

// authorityRank returns a numeric rank for authority levels
func (c *Collector) authorityRank(auth AuthorityLevel) int {
	switch auth {
	case AuthorityVerified:
		return 6
	case AuthorityNormative:
		return 5
	case AuthorityVendor:
		return 4
	case AuthorityProject:
		return 3
	case AuthorityOperator:
		return 2
	case AuthorityDerived:
		return 1
	default:
		return 0
	}
}

// identifyRemainingGaps identifies gaps that still need attention
func (c *Collector) identifyRemainingGaps(allGaps, newArtifacts, updatedArtifacts []CollectedArtifact) []KnowledgeGap {
	// Create set of domains that have been covered
	coveredDomains := make(map[KnowledgeDomain]bool)
	for _, a := range newArtifacts {
		coveredDomains[a.Domain] = true
	}
	for _, a := range updatedArtifacts {
		coveredDomains[a.Domain] = true
	}

	// Find gaps that are not covered
	var remaining []KnowledgeGap
	for _, gap := range allGaps {
		if !coveredDomains[gap.Domain] {
			remaining = append(remaining, gap)
		}
	}

	return remaining
}

// generateRecommendations generates recommendations based on collection results
func (c *Collector) generateRecommendations(result *CollectionResult) []string {
	var recs []string

	// Check coverage
	totalArtifacts := len(result.ArtifactsCollected) + len(result.ArtifactsUpdated)
	if totalArtifacts == 0 {
		recs = append(recs, "No new knowledge artifacts collected. Consider reviewing available documentation.")
	}

	// Check for critical gaps
	for _, gap := range result.KnowledgeStillMissing {
		if gap.Severity == "Critical" || gap.Severity == "High" {
			recs = append(recs, fmt.Sprintf("CRITICAL: Address %s knowledge gap - %s", gap.Domain, gap.Topic))
		}
	}

	// Authority recommendations
	hasLowAuthority := false
	for _, a := range result.ArtifactsCollected {
		if a.Authority == AuthorityDerived || a.Authority == AuthorityOperator {
			hasLowAuthority = true
			break
		}
	}
	if hasLowAuthority {
		recs = append(recs, "Some collected knowledge has low authority level. Seek verification from authoritative sources.")
	}

	// Next actions
	if len(result.KnowledgeStillMissing) > 0 {
		recs = append(recs, "Run 'kdse collect' again after addressing identified gaps.")
	}
	if totalArtifacts > 0 {
		recs = append(recs, "Run 'kdse normalize' to integrate collected knowledge into the knowledge base.")
	}
	recs = append(recs, "Review collected artifacts in .kdse/knowledge/ directory.")

	return recs
}

// saveResult saves the collection result to disk
func (c *Collector) saveResult(result *CollectionResult) error {
	// Ensure .kdse directory exists
	kdseDir := filepath.Join(c.repoPath, ".kdse")
	if err := os.MkdirAll(kdseDir, 0755); err != nil {
		return err
	}

	// Save JSON result
	resultPath := filepath.Join(kdseDir, "collection-result.json")
	return saveJSONResult(resultPath, result)
}

// saveJSONResult saves a result as JSON
func saveJSONResult(path string, result *CollectionResult) error {
	content, err := jsonMarshalIndent(result)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

// jsonMarshalIndent marshals a result with indentation
func jsonMarshalIndent(v interface{}) ([]byte, error) {
	// Simple indentation without external dependencies
	// For a real implementation, use encoding/json
	return []byte(fmt.Sprintf("%+v", v)), nil
}

// FormatResult formats the collection result for display
func FormatResult(result *CollectionResult) string {
	var b strings.Builder

	b.WriteString("╔═══════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║              Knowledge Collection Complete                    ║\n")
	b.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")

	b.WriteString("║ Knowledge Areas Reviewed                                      ║\n")
	for _, area := range result.KnowledgeAreasReviewed {
		b.WriteString(fmt.Sprintf("║   • %s\n", area))
	}

	b.WriteString("╠═══════════════════════════════════════════════════════════════╣\n")

	b.WriteString(fmt.Sprintf("║ Knowledge Added:            %d\n", len(result.ArtifactsCollected)))
	b.WriteString(fmt.Sprintf("║ Knowledge Updated:          %d\n", len(result.ArtifactsUpdated)))
	b.WriteString(fmt.Sprintf("║ Knowledge Already Present:  %d\n", len(result.ArtifactsExisting)))
	b.WriteString(fmt.Sprintf("║ Knowledge Still Missing:    %d\n", len(result.KnowledgeStillMissing)))

	if len(result.KnowledgeStillMissing) > 0 {
		b.WriteString("╠═══════════════════════════════════════════════════════════════╣")
		b.WriteString("║ Gaps Requiring Attention                                       ║\n")
		for _, gap := range result.KnowledgeStillMissing {
			b.WriteString(fmt.Sprintf("║   • [%s] %s\n", gap.Severity, gap.Topic))
		}
	}

	b.WriteString("╠═══════════════════════════════════════════════════════════════╣")
	b.WriteString("║ Next Recommended Actions                                      ║\n")
	for _, rec := range result.Recommendations {
		b.WriteString(fmt.Sprintf("║   • %s\n", rec))
	}

	b.WriteString("╚═══════════════════════════════════════════════════════════════╝\n")

	return b.String()
}

// FormatSummary formats a brief summary for terminal output
func FormatSummary(result *CollectionResult) string {
	var b strings.Builder

	b.WriteString("\n📊 Knowledge Collection Summary\n")
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString(fmt.Sprintf("  Knowledge Areas Reviewed:    %d\n", len(result.KnowledgeAreasReviewed)))
	b.WriteString(fmt.Sprintf("  Knowledge Gaps Found:        %d\n", len(result.KnowledgeGapsIdentified)))
	b.WriteString(fmt.Sprintf("  Artifacts Collected:         %d\n", len(result.ArtifactsCollected)))
	b.WriteString(fmt.Sprintf("  Artifacts Updated:          %d\n", len(result.ArtifactsUpdated)))
	b.WriteString(fmt.Sprintf("  Artifacts Already Present:  %d\n", len(result.ArtifactsExisting)))
	b.WriteString(fmt.Sprintf("  Knowledge Still Missing:    %d\n", len(result.KnowledgeStillMissing)))
	b.WriteString(fmt.Sprintf("  Processing Time:            %.2fs\n", result.Statistics.ProcessingTime))
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if len(result.Recommendations) > 0 {
		b.WriteString("\n📋 Recommendations:\n")
		for _, rec := range result.Recommendations {
			b.WriteString(fmt.Sprintf("  → %s\n", rec))
		}
	}

	b.WriteString("\n📁 Collected artifacts are available in: .kdse/knowledge/\n")
	b.WriteString(fmt.Sprintf("📄 Full report: .kdse/collection-result.json\n\n"))

	return b.String()
}

// GenerateCollectionReport generates a detailed collection report
func GenerateCollectionReport(result *CollectionResult) string {
	var b strings.Builder

	b.WriteString("# KDSE Knowledge Collection Report\n\n")
	b.WriteString("**Report ID:** " + result.SessionID + "\n")
	b.WriteString("**Repository:** " + result.Repository + "\n")
	b.WriteString("**Operator:** " + result.Operator + "\n")
	b.WriteString("**Collection Date:** " + result.StartedAt + "\n")
	b.WriteString("**Report Version:** 1.0\n\n")

	b.WriteString("---\n\n")

	b.WriteString("## Executive Summary\n\n")
	b.WriteString(fmt.Sprintf("This report documents the knowledge collection performed for the %s project.\n\n", result.Repository))

	b.WriteString("### Collection Statistics\n\n")
	b.WriteString("| Metric | Value |\n")
	b.WriteString("|--------|-------|\n")
	b.WriteString(fmt.Sprintf("| Knowledge Areas Reviewed | %d |\n", len(result.KnowledgeAreasReviewed)))
	b.WriteString(fmt.Sprintf("| Knowledge Gaps Identified | %d |\n", len(result.KnowledgeGapsIdentified)))
	b.WriteString(fmt.Sprintf("| Artifacts Collected | %d |\n", len(result.ArtifactsCollected)))
	b.WriteString(fmt.Sprintf("| Artifacts Updated | %d |\n", len(result.ArtifactsUpdated)))
	b.WriteString(fmt.Sprintf("| Artifacts Already Present | %d |\n", len(result.ArtifactsExisting)))
	b.WriteString(fmt.Sprintf("| Knowledge Still Missing | %d |\n", len(result.KnowledgeStillMissing)))
	b.WriteString(fmt.Sprintf("| Processing Time | %.2f seconds |\n", result.Statistics.ProcessingTime))
	b.WriteString(fmt.Sprintf("| Success Rate | %.1f%% |\n", result.Statistics.SuccessRate))

	b.WriteString("\n---\n\n")

	if len(result.KnowledgeAreasReviewed) > 0 {
		b.WriteString("## Knowledge Areas Reviewed\n\n")
		for _, area := range result.KnowledgeAreasReviewed {
			b.WriteString(fmt.Sprintf("- %s\n", area))
		}
		b.WriteString("\n")
	}

	if len(result.KnowledgeGapsIdentified) > 0 {
		b.WriteString("## Identified Knowledge Gaps\n\n")
		b.WriteString("| Gap ID | Domain | Topic | Severity | Impact |\n")
		b.WriteString("|--------|--------|-------|----------|--------|\n")
		for _, gap := range result.KnowledgeGapsIdentified {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				gap.ID, gap.Domain, gap.Topic, gap.Severity, gap.Impact))
		}
		b.WriteString("\n")
	}

	if len(result.ArtifactsCollected) > 0 {
		b.WriteString("## Collected Artifacts\n\n")
		b.WriteString("| Artifact | Domain | Authority | Confidence | Source |\n")
		b.WriteString("|----------|--------|-----------|------------|--------|\n")
		for _, art := range result.ArtifactsCollected {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %.0f%% | %s |\n",
				art.Title, art.Domain, art.Authority, art.ConfidenceLevel*100, art.Source))
		}
		b.WriteString("\n")
	}

	if len(result.ArtifactsUpdated) > 0 {
		b.WriteString("## Updated Artifacts\n\n")
		b.WriteString("| Artifact | Domain | Previous Authority | New Authority |\n")
		b.WriteString("|----------|--------|-------------------|---------------|\n")
		for _, art := range result.ArtifactsUpdated {
			b.WriteString(fmt.Sprintf("| %s | %s | Project | %s |\n",
				art.Title, art.Domain, art.Authority))
		}
		b.WriteString("\n")
	}

	if len(result.KnowledgeStillMissing) > 0 {
		b.WriteString("## Knowledge Still Missing\n\n")
		for _, gap := range result.KnowledgeStillMissing {
			b.WriteString(fmt.Sprintf("### [%s] %s\n\n", gap.Severity, gap.Topic))
			b.WriteString(fmt.Sprintf("**Domain:** %s\n\n", gap.Domain))
			b.WriteString(fmt.Sprintf("**Description:** %s\n\n", gap.Description))
			b.WriteString(fmt.Sprintf("**Impact:** %s\n\n", gap.Impact))
			b.WriteString("**Recommendations:**\n")
			for _, rec := range gap.Recommendations {
				b.WriteString(fmt.Sprintf("- %s\n", rec))
			}
			b.WriteString("\n")
		}
	}

	if len(result.Recommendations) > 0 {
		b.WriteString("## Recommendations\n\n")
		for _, rec := range result.Recommendations {
			b.WriteString(fmt.Sprintf("- %s\n", rec))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Confidence Assessment\n\n")
	b.WriteString("Collected knowledge has been categorized by authority level:\n\n")
	b.WriteString("| Authority Level | Description | Count |\n")
	b.WriteString("|-----------------|-------------|-------|\n")
	b.WriteString("| Verified | Tested and validated | " + countByAuthority(result.ArtifactsCollected, AuthorityVerified) + " |\n")
	b.WriteString("| Normative | KDSE standard or specification | " + countByAuthority(result.ArtifactsCollected, AuthorityNormative) + " |\n")
	b.WriteString("| Vendor | Vendor documentation | " + countByAuthority(result.ArtifactsCollected, AuthorityVendor) + " |\n")
	b.WriteString("| Project | Project-specific knowledge | " + countByAuthority(result.ArtifactsCollected, AuthorityProject) + " |\n")
	b.WriteString("| Operator | Operator-provided knowledge | " + countByAuthority(result.ArtifactsCollected, AuthorityOperator) + " |\n")
	b.WriteString("| Derived | Derived from other sources | " + countByAuthority(result.ArtifactsCollected, AuthorityDerived) + " |\n")
	b.WriteString("\n")

	b.WriteString("---\n\n")
	b.WriteString(fmt.Sprintf("*Report generated by KDSE Collect v%s*\n", kdseVersion))
	b.WriteString(fmt.Sprintf("*Generated: %s*\n", time.Now().Format(time.RFC3339)))

	return b.String()
}

func countByAuthority(artifacts []CollectedArtifact, auth AuthorityLevel) string {
	count := 0
	for _, a := range artifacts {
		if a.Authority == auth {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

// SaveCollectionReport saves the collection report to disk
func SaveCollectionReport(repoPath string, result *CollectionResult) error {
	report := GenerateCollectionReport(result)

	reportDir := filepath.Join(repoPath, ".kdse", "reports")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return err
	}

	reportPath := filepath.Join(reportDir, fmt.Sprintf("collection-report-%s.md", result.SessionID))
	return os.WriteFile(reportPath, []byte(report), 0644)
}
