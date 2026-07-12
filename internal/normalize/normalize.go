package normalize

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/kdse/runtime/internal/types"
)

// Normalizer orchestrates the documentation normalization process
type Normalizer struct {
	repoPath  string
	sessionID string
}

// NewNormalizer creates a new normalizer
func NewNormalizer(repoPath string) *Normalizer {
	return &Normalizer{
		repoPath:  repoPath,
		sessionID: types.GenerateSessionID(),
	}
}

// Normalize executes the full normalization process
func (n *Normalizer) Normalize() (*NormalizationResult, error) {
	startTime := time.Now()

	// Initialize result
	result := NewNormalizationResult(n.sessionID, n.repoPath)

	// Step 1: Discover documentation
	fmt.Println("Discovering documentation...")
	docs, err := n.discoverDocumentation()
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}
	result.DiscoveredDocs = docs

	// Step 2: Extract knowledge
	fmt.Printf("Analyzing %d documentation files...\n", len(docs))
	extractor := NewExtractor()
	knowledge := extractor.ExtractKnowledge(docs)
	result.KnowledgeExtract = knowledge

	// Step 3: Mark preserved documents
	for _, doc := range docs {
		result.AddPreserved(doc.Path)
	}

	// Step 4: Generate KDSE artifacts
	fmt.Println("Generating KDSE-standard documentation...")
	generator := NewGenerator(n.repoPath, n.sessionID)
	artifacts, err := generator.Generate(result)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}
	result.NormalizedArts = artifacts

	// Step 5: Calculate statistics
	result.CalculateStats(startTime)
	result.SetCompleted()

	return result, nil
}

// NormalizeWithOptions executes normalization with specific options
func (n *Normalizer) NormalizeWithOptions(options *NormalizeOptions) (*NormalizationResult, error) {
	startTime := time.Now()

	// Initialize result
	result := NewNormalizationResult(n.sessionID, n.repoPath)

	// Step 1: Discover documentation with options
	fmt.Println("Discovering documentation...")
	var docs []DiscoveredDocument
	var err error

	if options != nil && len(options.IncludeTypes) > 0 {
		docs, err = n.discoverByTypes(options.IncludeTypes)
	} else {
		docs, err = n.discoverDocumentation()
	}
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}
	result.DiscoveredDocs = docs

	// Step 2: Extract knowledge
	fmt.Printf("Analyzing %d documentation files...\n", len(docs))
	extractor := NewExtractor()
	knowledge := extractor.ExtractKnowledge(docs)
	result.KnowledgeExtract = knowledge

	// Step 3: Mark preserved documents
	for _, doc := range docs {
		result.AddPreserved(doc.Path)
	}

	// Step 4: Generate KDSE artifacts
	fmt.Println("Generating KDSE-standard documentation...")
	generator := NewGenerator(n.repoPath, n.sessionID)
	artifacts, err := generator.Generate(result)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}
	result.NormalizedArts = artifacts

	// Step 5: Calculate statistics
	result.CalculateStats(startTime)
	result.SetCompleted()

	return result, nil
}

func (n *Normalizer) discoverDocumentation() ([]DiscoveredDocument, error) {
	discovery := NewDiscovery(n.repoPath)
	return discovery.Discover()
}

func (n *Normalizer) discoverByTypes(types []DocType) ([]DiscoveredDocument, error) {
	discovery := NewDiscovery(n.repoPath)
	return discovery.DiscoverByType(types...)
}

// NormalizeOptions contains configuration for normalization
type NormalizeOptions struct {
	IncludeTypes []DocType
	ExcludeTypes []DocType
	MaxFiles     int
}

// ReportFormatter formats the normalization report
type ReportFormatter struct{}

// FormatReport generates a human-readable normalization report
func FormatReport(result *NormalizationResult) string {
	var b strings.Builder

	b.WriteString("# KDSE Normalization Report\n\n")
	b.WriteString("| Field | Value |\n")
	b.WriteString("|-------|-------|\n")
	b.WriteString(fmt.Sprintf("| Session ID | %s |\n", result.SessionID))
	b.WriteString(fmt.Sprintf("| Repository | %s |\n", result.Repository))
	b.WriteString(fmt.Sprintf("| Started | %s |\n", result.StartedAt))
	b.WriteString(fmt.Sprintf("| Completed | %s |\n", result.CompletedAt))
	b.WriteString(fmt.Sprintf("| Report Version | 1.0 |\n\n"))

	b.WriteString("## Summary\n\n")
	b.WriteString(fmt.Sprintf("**Documents Found:** %d\n", result.Statistics.TotalDocsFound))
	b.WriteString(fmt.Sprintf("**Artifacts Generated:** %d\n", result.Statistics.TotalArtifactsGen))
	b.WriteString(fmt.Sprintf("**Processing Time:** %.2f seconds\n", result.Statistics.ProcessingTime))
	b.WriteString(fmt.Sprintf("**Success Rate:** %.1f%%\n\n", result.Statistics.SuccessRate))

	if len(result.DiscoveredDocs) > 0 {
		b.WriteString("## Discovered Documentation\n\n")
		b.WriteString("| Path | Type | Confidence | Size |\n")
		b.WriteString("|------|------|------------|------|\n")
		for _, doc := range result.DiscoveredDocs {
			b.WriteString(fmt.Sprintf("| %s | %s | %.0f%% | %d bytes |\n",
				doc.Path, doc.Type, doc.Confidence*100, doc.Size))
		}
		b.WriteString("\n")
	}

	if len(result.NormalizedArts) > 0 {
		b.WriteString("## Generated KDSE Artifacts\n\n")
		b.WriteString("| Artifact | Path | Traceability |\n")
		b.WriteString("|----------|------|-------------|\n")
		for _, art := range result.NormalizedArts {
			b.WriteString(fmt.Sprintf("| %s | %s | %d sources |\n",
				art.Title, filepath.Base(art.Path), len(art.Traceability.DerivedFrom)))
		}
		b.WriteString("\n")
	}

	if len(result.KnowledgeExtract.Requirements) > 0 {
		b.WriteString("## Extracted Requirements\n\n")
		for i, req := range result.KnowledgeExtract.Requirements[:min(10, len(result.KnowledgeExtract.Requirements))] {
			b.WriteString(fmt.Sprintf("%d. %s\n", i+1, req.Description))
			b.WriteString(fmt.Sprintf("   - Source: %s\n\n", req.Source))
		}
		if len(result.KnowledgeExtract.Requirements) > 10 {
			b.WriteString(fmt.Sprintf("*... and %d more requirements\n\n", len(result.KnowledgeExtract.Requirements)-10))
		}
	}

	if len(result.KnowledgeExtract.Decisions) > 0 {
		b.WriteString("## Extracted Decisions\n\n")
		for _, dec := range result.KnowledgeExtract.Decisions {
			b.WriteString(fmt.Sprintf("### %s: %s\n\n", dec.ID, dec.Title))
			b.WriteString(fmt.Sprintf("**Status:** %s\n\n", dec.Status))
			b.WriteString(fmt.Sprintf("%s\n\n", dec.Decision))
		}
	}

	if len(result.KnowledgeExtract.Glossary) > 0 {
		b.WriteString("## Extracted Glossary Terms\n\n")
		for _, term := range result.KnowledgeExtract.Glossary[:min(20, len(result.KnowledgeExtract.Glossary))] {
			b.WriteString(fmt.Sprintf("**%s:** %s\n\n", term.Term, term.Definition))
		}
		if len(result.KnowledgeExtract.Glossary) > 20 {
			b.WriteString(fmt.Sprintf("*... and %d more terms\n", len(result.KnowledgeExtract.Glossary)-20))
		}
	}

	b.WriteString("## Traceability\n\n")
	b.WriteString("Every normalized artifact includes:\n\n")
	b.WriteString("- **Derived From:** Original documents used\n")
	b.WriteString("- **Authority:** Source attribution\n")
	b.WriteString("- **Version:** KDSE version information\n")
	b.WriteString("- **Normative References:** KDSE standard references\n\n")

	b.WriteString("## Next Steps\n\n")
	b.WriteString("1. Review generated KDSE artifacts in `.kdse/normalized/`\n")
	b.WriteString("2. Verify traceability information is correct\n")
	b.WriteString("3. Run `kdse run` to start a KDSE session with normalized context\n")
	b.WriteString("4. Original documentation remains unchanged\n\n")

	b.WriteString("---\n\n")
	b.WriteString("*Generated by KDSE Normalization v1.0.0*\n")

	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
