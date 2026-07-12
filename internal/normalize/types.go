package normalize

import (
	"time"

	"github.com/kdse/runtime/internal/types"
)

// DiscoveredDocument represents a document found in the repository
type DiscoveredDocument struct {
	Path           string   `json:"path"`
	Type           DocType  `json:"type"`
	Title          string   `json:"title"`
	Content        string   `json:"content,omitempty"`
	Size           int64    `json:"size"`
	LastModified   string   `json:"last_modified"`
	Encoding       string   `json:"encoding"`
	Confidence     float64  `json:"confidence"`
	KnowledgeTypes []string `json:"knowledge_types"`
}

// DocType represents the type of discovered documentation
type DocType string

const (
	DocsTypeREADME      DocType = "README"
	DocsTypeArchitecture DocType = "Architecture"
	DocsTypeDesign       DocType = "Design"
	DocsTypeAPI          DocType = "API"
	DocsTypeGuide        DocType = "Guide"
	DocsTypeTutorial     DocType = "Tutorial"
	DocsTypeReference    DocType = "Reference"
	DocsTypeChangelog    DocType = "Changelog"
	DocsTypeContributing DocType = "Contributing"
	DocsTypeLicense      DocType = "License"
	DocsTypeADR          DocType = "ADR"
	DocsTypeEngineering DocType = "Engineering"
	DocsTypeWiki         DocType = "Wiki"
	DocsTypeExamples     DocType = "Examples"
	DocsTypeUnknown      DocType = "Unknown"
)

// NormalizedArtifact represents a KDSE-standard document generated from existing docs
type NormalizedArtifact struct {
	ID              string           `json:"id"`
	Path            string           `json:"path"`
	DocType         string           `json:"doc_type"`
	Title           string           `json:"title"`
	Content         string           `json:"content"`
	Traceability    Traceability     `json:"traceability"`
	Metadata        ArtifactMetadata `json:"metadata"`
	CreatedAt       string           `json:"created_at"`
}

// Traceability records the relationship between original and KDSE documents
type Traceability struct {
	DerivedFrom        []DerivedFrom `json:"derived_from"`
	OriginalDocuments  []string      `json:"original_documents"`
	Authority          string        `json:"authority"`
	Version            string        `json:"version"`
	NormativeRefs      []string      `json:"normative_refs"`
	TransformationType string         `json:"transformation_type"`
}

// DerivedFrom records the source document for each piece of knowledge
type DerivedFrom struct {
	Path        string `json:"path"`
	DocType    string `json:"doc_type"`
	Confidence string `json:"confidence"`
	Excerpt    string `json:"excerpt,omitempty"`
}

// ArtifactMetadata contains standard KDSE artifact metadata
type ArtifactMetadata struct {
	Author      string   `json:"author,omitempty"`
	Version     string   `json:"version"`
	KDSEVersion string   `json:"kdse_version"`
	Phase       string   `json:"phase"`
	Tags        []string `json:"tags"`
	Domain      string   `json:"domain,omitempty"`
}

// NormalizationResult contains the results of a normalization operation
type NormalizationResult struct {
	SessionID         string              `json:"session_id"`
	StartedAt         string              `json:"started_at"`
	CompletedAt       string              `json:"completed_at"`
	Repository        string              `json:"repository"`
	DiscoveredDocs    []DiscoveredDocument `json:"discovered_docs"`
	NormalizedArts    []NormalizedArtifact `json:"normalized_artifacts"`
	PreservedDocs     []string            `json:"preserved_docs"`
	SkippedDocs       []string            `json:"skipped_docs"`
	Statistics        NormalizationStats  `json:"statistics"`
	KnowledgeExtract  KnowledgeExtraction  `json:"knowledge_extraction"`
}

// NormalizationStats contains statistics about the normalization process
type NormalizationStats struct {
	TotalDocsFound    int     `json:"total_docs_found"`
	TotalArtifactsGen int     `json:"total_artifacts_generated"`
	TotalBytesWritten int64   `json:"total_bytes_written"`
	ProcessingTime    float64 `json:"processing_time_seconds"`
	SuccessRate       float64 `json:"success_rate"`
}

// KnowledgeExtraction contains extracted engineering knowledge
type KnowledgeExtraction struct {
	Domain       string             `json:"domain"`
	Purpose      string             `json:"purpose"`
	Stakeholders []string           `json:"stakeholders"`
	Requirements []ExtractedReq     `json:"requirements"`
	Decisions    []ExtractedDecision `json:"decisions"`
	Constraints  []string           `json:"constraints"`
	Glossary     []GlossaryTerm     `json:"glossary"`
}

// ExtractedReq represents an extracted requirement
type ExtractedReq struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Type        string `json:"type"`
}

// ExtractedDecision represents an extracted architectural decision
type ExtractedDecision struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Context     string   `json:"context"`
	Decision    string   `json:"decision"`
	Rationale   string   `json:"rationale"`
	Source      string   `json:"source"`
	Status      string   `json:"status"`
	Alternatives []string `json:"alternatives,omitempty"`
}

// GlossaryTerm represents an extracted glossary term
type GlossaryTerm struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
	Source     string `json:"source"`
}

// Analyzer handles documentation analysis
type Analyzer struct {
	repoPath string
}

// NewAnalyzer creates a new documentation analyzer
func NewAnalyzer(repoPath string) *Analyzer {
	return &Analyzer{repoPath: repoPath}
}

// NewNormalizationResult creates a new normalization result
func NewNormalizationResult(sessionID, repoPath string) *NormalizationResult {
	return &NormalizationResult{
		SessionID:         sessionID,
		StartedAt:         time.Now().Format(time.RFC3339),
		Repository:        repoPath,
		DiscoveredDocs:    []DiscoveredDocument{},
		NormalizedArts:    []NormalizedArtifact{},
		PreservedDocs:     []string{},
		SkippedDocs:       []string{},
		KnowledgeExtract:  KnowledgeExtraction{},
	}
}

// SetCompleted marks the normalization as complete
func (r *NormalizationResult) SetCompleted() {
	r.CompletedAt = time.Now().Format(time.RFC3339)
}

// CalculateStats computes statistics from the normalization
func (r *NormalizationResult) CalculateStats(startTime time.Time) {
	r.Statistics = NormalizationStats{
		TotalDocsFound:    len(r.DiscoveredDocs),
		TotalArtifactsGen: len(r.NormalizedArts),
		ProcessingTime:    time.Since(startTime).Seconds(),
		SuccessRate:       calculateSuccessRate(len(r.NormalizedArts), len(r.DiscoveredDocs)),
	}
}

func calculateSuccessRate(generated, discovered int) float64 {
	if discovered == 0 {
		return 0
	}
	return float64(generated) / float64(discovered) * 100
}

// FormatReport generates a human-readable normalization report
func (r *NormalizationResult) FormatReport() string {
	return FormatReport(r)
}

// GetPhase returns the engineering phase based on normalized artifacts
func (r *NormalizationResult) GetPhase() types.EngineeringPhase {
	count := len(r.NormalizedArts)
	switch {
	case count >= 8:
		return types.PhaseValidated
	case count >= 5:
		return types.PhaseUsable
	case count >= 3:
		return types.PhaseStructured
	case count >= 1:
		return types.PhaseDefined
	default:
		return types.PhaseConcept
	}
}

// AddPreserved adds a document to the preserved list
func (r *NormalizationResult) AddPreserved(path string) {
	r.PreservedDocs = append(r.PreservedDocs, path)
}

// AddSkipped adds a document to the skipped list
func (r *NormalizationResult) AddSkipped(path string) {
	r.SkippedDocs = append(r.SkippedDocs, path)
}
