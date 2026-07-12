package collect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RepositoryProvider collects knowledge from repository documentation
type RepositoryProvider struct {
	repoPath string
}

// NewRepositoryProvider creates a new repository knowledge provider
func NewRepositoryProvider(repoPath string) *RepositoryProvider {
	return &RepositoryProvider{repoPath: repoPath}
}

// Name returns the provider name
func (p *RepositoryProvider) Name() string {
	return "Repository Documentation"
}

// Collect collects knowledge from repository documentation
func (p *RepositoryProvider) Collect(input *CollectionInput) ([]CollectedArtifact, error) {
	var artifacts []CollectedArtifact

	// Scan for documentation files
	docPatterns := []string{".md", ".txt", ".rst", ".adoc"}
	
	err := filepath.Walk(p.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and common non-doc directories
		relPath, _ := filepath.Rel(p.repoPath, path)
		if strings.HasPrefix(relPath, ".") || 
		   strings.Contains(relPath, "node_modules") ||
		   strings.Contains(relPath, "vendor") ||
		   strings.Contains(relPath, "__pycache__") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// Check if it's a documentation file
		ext := strings.ToLower(filepath.Ext(path))
		for _, pattern := range docPatterns {
			if ext == pattern {
				artifact, err := p.extractKnowledge(path, input)
				if err == nil && artifact != nil {
					artifacts = append(artifacts, *artifact)
				}
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("repository scan failed: %w", err)
	}

	return artifacts, nil
}

// CanCollect returns true for all domains
func (p *RepositoryProvider) CanCollect(domain KnowledgeDomain) bool {
	return true
}

// extractKnowledge extracts knowledge from a documentation file
func (p *RepositoryProvider) extractKnowledge(path string, input *CollectionInput) (*CollectedArtifact, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := string(content)
	if len(text) < 50 {
		return nil, nil // Too short to be meaningful
	}

	// Determine domain from file path and content
	domain := p.inferDomain(filepath.Base(path), text)
	
	// Determine authority from file location
	authority := p.determineAuthority(path)

	// Generate title from filename
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Create summary from first paragraph
	summary := p.extractSummary(text)

	return &CollectedArtifact{
		ID:              GenerateArtifactID(domain),
		Domain:          domain,
		Title:           title,
		Content:         text,
		Summary:         summary,
		Source:          SourceRepository,
		SourcePath:      path,
		Authority:       authority,
		CollectionDate:  time.Now().Format(time.RFC3339),
		CollectedBy:     input.OperatorName,
		NormativeRefs:   []string{},
		Dependencies:    []string{},
		Traceability: Traceability{
			ArtifactID:     GenerateArtifactID(domain),
			Source:         "Repository: " + p.repoPath,
			Authority:      string(authority),
			Version:        kdseVersion,
			CollectionDate: time.Now().Format(time.RFC3339),
			CollectedBy:    input.OperatorName,
			NormativeRefs:  []string{},
			Dependencies:    []string{},
			TraceabilityID: fmt.Sprintf("REPO-%s", time.Now().Format("20060102-150405")),
		},
		ConfidenceLevel: 0.7,
		Tags:            p.extractTags(text),
	}, nil
}

// inferDomain determines the knowledge domain from file path and content
func (p *RepositoryProvider) inferDomain(filename, content string) KnowledgeDomain {
	filenameLower := strings.ToLower(filename)
	contentLower := strings.ToLower(content)

	// Check filename first
	if strings.Contains(filenameLower, "transformer") {
		return DomainTransformers
	}
	if strings.Contains(filenameLower, "battery") {
		return DomainBattery
	}
	if strings.Contains(filenameLower, "relay") || strings.Contains(filenameLower, "protection") {
		return DomainRelay
	}
	if strings.Contains(filenameLower, "weather") || strings.Contains(filenameLower, "solar") {
		return DomainWeather
	}
	if strings.Contains(filenameLower, "physics") || strings.Contains(filenameLower, "thermal") {
		return DomainPhysics
	}
	if strings.Contains(filenameLower, "equipment") || strings.Contains(filenameLower, "hardware") {
		return DomainEquipment
	}
	if strings.Contains(filenameLower, "standard") || strings.Contains(filenameLower, "ieee") || 
	   strings.Contains(filenameLower, "iec") || strings.Contains(filenameLower, "iso") {
		return DomainStandards
	}
	if strings.Contains(filenameLower, "protocol") || strings.Contains(filenameLower, "api") {
		return DomainProtocols
	}
	if strings.Contains(filenameLower, "control") || strings.Contains(filenameLower, "algorithm") {
		return DomainControl
	}
	if strings.Contains(filenameLower, "simul") || strings.Contains(filenameLower, "model") {
		return DomainSimulation
	}
	if strings.Contains(filenameLower, "business") || strings.Contains(filenameLower, "requirement") {
		return DomainBusiness
	}
	if strings.Contains(filenameLower, "glossary") || strings.Contains(filenameLower, "vocabulary") {
		return DomainVocabulary
	}
	if strings.Contains(filenameLower, "environment") || strings.Contains(filenameLower, "site") {
		return DomainEnvironment
	}

	// Check content for domain indicators
	if containsKeyword(contentLower, "transformer", "magnetic flux", "inductance") {
		return DomainTransformers
	}
	if containsKeyword(contentLower, "battery", "soc", "soh", "charge state") {
		return DomainBattery
	}
	if containsKeyword(contentLower, "relay", "overcurrent", "protection coordination") {
		return DomainRelay
	}
	if containsKeyword(contentLower, "weather", "irradiance", "solar radiation") {
		return DomainWeather
	}

	return DomainGeneral
}

func containsKeyword(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

// determineAuthority determines the authority level from file location
func (p *RepositoryProvider) determineAuthority(path string) AuthorityLevel {
	relPath, _ := filepath.Rel(p.repoPath, path)
	relPath = strings.ToLower(relPath)

	if strings.Contains(relPath, "vendor") || strings.Contains(relPath, "manufacturer") {
		return AuthorityVendor
	}
	if strings.Contains(relPath, "standard") || strings.Contains(relPath, "spec") {
		return AuthorityNormative
	}
	if strings.Contains(relPath, "src/") || strings.Contains(relPath, "implementation") {
		return AuthorityProject
	}

	return AuthorityProject
}

// extractSummary extracts a summary from content
func (p *RepositoryProvider) extractSummary(content string) string {
	lines := strings.Split(content, "\n")
	var summaryLines []string
	lineCount := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and headers
		if line == "" || strings.HasPrefix(line, "#") || 
		   strings.HasPrefix(line, "===") || strings.HasPrefix(line, "---") {
			continue
		}
		summaryLines = append(summaryLines, line)
		lineCount++
		if lineCount >= 3 {
			break
		}
	}

	summary := strings.Join(summaryLines, " ")
	if len(summary) > 200 {
		summary = summary[:197] + "..."
	}
	return summary
}

// extractTags extracts tags from content
func (p *RepositoryProvider) extractTags(content string) []string {
	var tags []string
	contentLower := strings.ToLower(content)

	domainTags := map[string]KnowledgeDomain{
		"transformer": DomainTransformers,
		"battery":      DomainBattery,
		"relay":        DomainRelay,
		"protection":   DomainRelay,
		"weather":      DomainWeather,
		"solar":        DomainWeather,
		"physics":      DomainPhysics,
		"thermal":      DomainPhysics,
		"standard":     DomainStandards,
		"protocol":     DomainProtocols,
		"control":      DomainControl,
		"simulation":   DomainSimulation,
		"model":        DomainSimulation,
		"business":     DomainBusiness,
		"requirement":  DomainBusiness,
	}

	for keyword, domain := range domainTags {
		if strings.Contains(contentLower, keyword) {
			tags = append(tags, string(domain))
		}
	}

	return tags
}

// OperatorProvider collects knowledge from operator input
type OperatorProvider struct{}

// NewOperatorProvider creates a new operator knowledge provider
func NewOperatorProvider() *OperatorProvider {
	return &OperatorProvider{}
}

// Name returns the provider name
func (p *OperatorProvider) Name() string {
	return "Operator Knowledge"
}

// Collect collects knowledge from operator input
// This is a placeholder - in a real implementation, this would prompt for input
func (p *OperatorProvider) Collect(input *CollectionInput) ([]CollectedArtifact, error) {
	// In a full implementation, this would:
	// 1. Prompt the operator for knowledge
	// 2. Accept uploaded documents
	// 3. Process pasted content
	// For now, return empty - operator can add knowledge via file uploads
	return []CollectedArtifact{}, nil
}

// CanCollect returns true for all domains
func (p *OperatorProvider) CanCollect(domain KnowledgeDomain) bool {
	return true
}

// AIAssistedProvider provides AI-assisted knowledge collection
type AIAssistedProvider struct {
	repoPath string
}

// NewAIAssistedProvider creates a new AI-assisted knowledge provider
func NewAIAssistedProvider(repoPath string) *AIAssistedProvider {
	return &AIAssistedProvider{repoPath: repoPath}
}

// Name returns the provider name
func (p *AIAssistedProvider) Name() string {
	return "AI-Assisted Knowledge"
}

// Collect provides guidance for AI-assisted knowledge collection
func (p *AIAssistedProvider) Collect(input *CollectionInput) ([]CollectedArtifact, error) {
	var artifacts []CollectedArtifact

	// For gaps that need AI assistance, create placeholder artifacts
	// with instructions for AI-assisted collection
	for _, gap := range p.identifyGapsRequiringAI(input) {
		artifact := CollectedArtifact{
			ID:              GenerateArtifactID(gap.Domain),
			Domain:          gap.Domain,
			Title:           fmt.Sprintf("AI-Assisted: %s", gap.Topic),
			Content:         p.generateAICollectionPrompt(gap),
			Summary:         "This artifact requires AI-assisted knowledge collection. " + gap.Description,
			Source:          SourceAI,
			Authority:       AuthorityDerived,
			CollectionDate:  time.Now().Format(time.RFC3339),
			CollectedBy:     "AI-Assisted",
			NormativeRefs:   gap.Recommendations,
			Dependencies:    []string{},
			Traceability: Traceability{
				ArtifactID:     GenerateArtifactID(gap.Domain),
				Source:         "AI-Assisted Collection",
				Authority:      string(AuthorityDerived),
				Version:        kdseVersion,
				CollectionDate: time.Now().Format(time.RFC3339),
				CollectedBy:    "AI-Assisted",
				NormativeRefs:  gap.Recommendations,
				Dependencies:   []string{},
				TraceabilityID: fmt.Sprintf("AI-%s", time.Now().Format("20060102-150405")),
			},
			ConfidenceLevel: 0.5,
			Tags:            []string{string(gap.Domain), "ai-assisted", "pending"},
		}
		artifacts = append(artifacts, artifact)
	}

	return artifacts, nil
}

// CanCollect returns true for all domains
func (p *AIAssistedProvider) CanCollect(domain KnowledgeDomain) bool {
	return true
}

// identifyGapsRequiringAI identifies gaps that need AI-assisted collection
func (p *AIAssistedProvider) identifyGapsRequiringAI(input *CollectionInput) []KnowledgeGap {
	// For now, return gaps that are marked as needing external research
	var gaps []KnowledgeGap
	
	for _, domain := range input.KnowledgeAreas {
		// These domains typically require external research
		if domain == DomainStandards || domain == DomainVendor || 
		   domain == DomainPhysics || domain == DomainSimulation {
			gaps = append(gaps, KnowledgeGap{
				ID:              GenerateGapID(domain, 1),
				Domain:          domain,
				Topic:           string(domain) + " Knowledge",
				Description:     "Requires external research and authoritative sources",
				Severity:        "Medium",
				Impact:          "External sources needed for complete knowledge",
				Recommendations: []string{
					"Research authoritative sources",
					"Consult vendor documentation",
					"Review applicable standards",
				},
			})
		}
	}

	return gaps
}

// generateAICollectionPrompt generates a prompt for AI-assisted collection
func (p *AIAssistedProvider) generateAICollectionPrompt(gap KnowledgeGap) string {
	return fmt.Sprintf(`# AI-Assisted Knowledge Collection

## Topic: %s
## Domain: %s
## Gap ID: %s

## Description
%s

## Impact
%s

## Recommended Actions
%s

## Instructions
1. Research authoritative sources for this knowledge area
2. Document findings in KDSE-standard format
3. Include proper traceability and references
4. Set appropriate confidence level

## Authority Level
Target: %s or higher

---
This artifact was created to guide AI-assisted knowledge collection.
Complete the knowledge collection and update this artifact.
`, gap.Topic, gap.Domain, gap.ID, gap.Description, gap.Impact,
		strings.Join(gap.Recommendations, "\n- "), AuthorityNormative)
}

// StandardsProvider collects knowledge from standards documents
type StandardsProvider struct {
	repoPath string
}

// NewStandardsProvider creates a new standards knowledge provider
func NewStandardsProvider(repoPath string) *StandardsProvider {
	return &StandardsProvider{repoPath: repoPath}
}

// Name returns the provider name
func (p *StandardsProvider) Name() string {
	return "Standards Documentation"
}

// Collect collects knowledge from standards documents
func (p *StandardsProvider) Collect(input *CollectionInput) ([]CollectedArtifact, error) {
	var artifacts []CollectedArtifact

	// Look for standards directories
	standardsDirs := []string{"standards/", "spec/", "reference/", "normative/"}

	for _, dir := range standardsDirs {
		standardsPath := filepath.Join(p.repoPath, dir)
		if info, err := os.Stat(standardsPath); err == nil && info.IsDir() {
			p.scanStandardsDir(standardsPath, input, &artifacts)
		}
	}

	return artifacts, nil
}

// scanStandardsDir scans a standards directory for documents
func (p *StandardsProvider) scanStandardsDir(dirPath string, input *CollectionInput, artifacts *[]CollectedArtifact) {
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" || ext == ".pdf" || ext == ".txt" {
			artifact, err := p.extractStandardsKnowledge(path, input)
			if err == nil && artifact != nil {
				*artifacts = append(*artifacts, *artifact)
			}
		}

		return nil
	})
}

// CanCollect returns true for standards domain
func (p *StandardsProvider) CanCollect(domain KnowledgeDomain) bool {
	return domain == DomainStandards
}

// extractStandardsKnowledge extracts knowledge from a standards document
func (p *StandardsProvider) extractStandardsKnowledge(path string, input *CollectionInput) (*CollectedArtifact, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := string(content)
	if len(text) < 100 {
		return nil, nil
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	domain := DomainStandards

	return &CollectedArtifact{
		ID:              GenerateArtifactID(domain),
		Domain:          domain,
		Title:           title,
		Content:         text,
		Summary:         p.extractSummary(text),
		Source:          SourceStandards,
		SourcePath:      path,
		Authority:       AuthorityNormative,
		CollectionDate:  time.Now().Format(time.RFC3339),
		CollectedBy:     input.OperatorName,
		NormativeRefs:   p.extractStandardRefs(text),
		Dependencies:    []string{},
		Traceability: Traceability{
			ArtifactID:     GenerateArtifactID(domain),
			Source:         "Standards: " + path,
			Authority:      string(AuthorityNormative),
			Version:        kdseVersion,
			CollectionDate: time.Now().Format(time.RFC3339),
			CollectedBy:    input.OperatorName,
			NormativeRefs:  p.extractStandardRefs(text),
			Dependencies:   []string{},
			TraceabilityID: fmt.Sprintf("STD-%s", time.Now().Format("20060102-150405")),
		},
		ConfidenceLevel: 0.95,
		Tags:            []string{"standards", "normative"},
	}, nil
}

// extractSummary extracts a summary from standards content
func (p *StandardsProvider) extractSummary(content string) string {
	lines := strings.Split(content, "\n")
	var summaryLines []string
	lineCount := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		summaryLines = append(summaryLines, line)
		lineCount++
		if lineCount >= 2 {
			break
		}
	}

	summary := strings.Join(summaryLines, " ")
	if len(summary) > 200 {
		summary = summary[:197] + "..."
	}
	return summary
}

// extractStandardRefs extracts standard references from content
func (p *StandardsProvider) extractStandardRefs(content string) []string {
	var refs []string
	patterns := []string{"IEC ", "IEEE ", "ISO ", "NIST ", "ANSI ", "UL ", "CSA ", "EN "}

	contentLower := strings.ToLower(content)
	for _, pattern := range patterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			// Extract context around the standard reference
			idx := strings.Index(contentLower, strings.ToLower(pattern))
			start := idx - 10
			if start < 0 {
				start = 0
			}
			end := idx + 50
			if end > len(content) {
				end = len(content)
			}
			refs = append(refs, strings.TrimSpace(content[start:end]))
		}
	}

	return refs
}

// ProviderRegistry manages available knowledge providers
type ProviderRegistry struct {
	providers []KnowledgeProvider
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry(repoPath string) *ProviderRegistry {
	registry := &ProviderRegistry{
		providers: []KnowledgeProvider{
			NewRepositoryProvider(repoPath),
			NewOperatorProvider(),
			NewAIAssistedProvider(repoPath),
			NewStandardsProvider(repoPath),
		},
	}
	return registry
}

// Register adds a provider to the registry
func (r *ProviderRegistry) Register(provider KnowledgeProvider) {
	r.providers = append(r.providers, provider)
}

// GetProviders returns all registered providers
func (r *ProviderRegistry) GetProviders() []KnowledgeProvider {
	return r.providers
}

// GetProvidersForDomain returns providers that can collect for a domain
func (r *ProviderRegistry) GetProvidersForDomain(domain KnowledgeDomain) []KnowledgeProvider {
	var result []KnowledgeProvider
	for _, p := range r.providers {
		if p.CanCollect(domain) {
			result = append(result, p)
		}
	}
	return result
}
