package collect

import (
	"time"

	"github.com/kdse/runtime/internal/types"
)

const kdseVersion = "1.0.0"

// AuthorityLevel represents the level of authority for collected knowledge
type AuthorityLevel string

const (
	AuthorityVerified   AuthorityLevel = "Verified"    // Tested and validated
	AuthorityNormative  AuthorityLevel = "Normative"   // KDSE standard or specification
	AuthorityVendor     AuthorityLevel = "Vendor"      // From vendor documentation
	AuthorityProject    AuthorityLevel = "Project"     // Project-specific knowledge
	AuthorityOperator   AuthorityLevel = "Operator"    // Operator-provided knowledge
	AuthorityDerived    AuthorityLevel = "Derived"     // Derived from other sources
)

// KnowledgeDomain represents the domain category for collected knowledge
type KnowledgeDomain string

const (
	DomainPhysics      KnowledgeDomain = "physics"
	DomainEquipment    KnowledgeDomain = "equipment"
	DomainEnvironment  KnowledgeDomain = "environment"
	DomainStandards    KnowledgeDomain = "standards"
	DomainBusiness     KnowledgeDomain = "business"
	DomainSimulation   KnowledgeDomain = "simulation"
	DomainControl      KnowledgeDomain = "control"
	DomainProtocols    KnowledgeDomain = "protocols"
	DomainVocabulary   KnowledgeDomain = "vocabulary"
	DomainTransformers KnowledgeDomain = "transformers"
	DomainBattery      KnowledgeDomain = "battery"
	DomainRelay       KnowledgeDomain = "relay"
	DomainWeather     KnowledgeDomain = "weather"
	DomainGeneral     KnowledgeDomain = "general"
)

// KnowledgeSource represents the source of collected knowledge
type KnowledgeSource string

const (
	SourceRepository   KnowledgeSource = "repository"
	SourceUpload      KnowledgeSource = "upload"
	SourceStandards    KnowledgeSource = "standards"
	SourceVendor      KnowledgeSource = "vendor"
	SourceOperator    KnowledgeSource = "operator"
	SourceAI         KnowledgeSource = "ai-assisted"
	SourceExternal    KnowledgeSource = "external"
)

// KnowledgeGap represents an identified gap in engineering knowledge
type KnowledgeGap struct {
	ID           string          `json:"id"`
	Domain       KnowledgeDomain `json:"domain"`
	Topic        string          `json:"topic"`
	Description  string          `json:"description"`
	Severity     string          `json:"severity"` // Critical, High, Medium, Low
	Impact       string          `json:"impact"`    // How this gap affects implementation
	Recommendations []string      `json:"recommendations"`
}

// CollectedArtifact represents a piece of collected and normalized knowledge
type CollectedArtifact struct {
	ID              string            `json:"id"`
	Domain          KnowledgeDomain   `json:"domain"`
	Title           string            `json:"title"`
	Content         string            `json:"content"`
	Summary         string            `json:"summary"`
	Source          KnowledgeSource   `json:"source"`
	SourcePath      string            `json:"source_path,omitempty"`
	Authority       AuthorityLevel    `json:"authority"`
	Version         string            `json:"version,omitempty"`
	CollectionDate  string            `json:"collection_date"`
	CollectedBy     string            `json:"collected_by,omitempty"`
	NormativeRefs   []string          `json:"normative_refs,omitempty"`
	Dependencies    []string          `json:"dependencies,omitempty"`
	Traceability    Traceability      `json:"traceability"`
	ConfidenceLevel float64           `json:"confidence_level"` // 0.0-1.0
	Tags            []string          `json:"tags,omitempty"`
	Path            string            `json:"path,omitempty"`   // File path where stored
}

// Traceability contains traceability information for collected artifacts
type Traceability struct {
	ArtifactID      string   `json:"artifact_id"`
	Source          string   `json:"source"`
	Authority       string   `json:"authority"`
	Version         string   `json:"version"`
	CollectionDate  string   `json:"collection_date"`
	CollectedBy     string   `json:"collected_by"`
	NormativeRefs   []string `json:"normative_refs"`
	Dependencies    []string `json:"dependencies"`
	TraceabilityID  string   `json:"traceability_id"`
}

// CollectionInput contains input data for knowledge collection
type CollectionInput struct {
	RepositoryPath  string
	AuditFindings   []types.Finding
	NormResult      interface{} // Can be *normalize.NormalizationResult
	SessionState    *types.SessionState
	OperatorName    string
	KnowledgeAreas  []KnowledgeDomain
	PriorityLevel   string // Critical, High, Medium, Low
}

// CollectionResult contains the results of a knowledge collection operation
type CollectionResult struct {
	SessionID          string             `json:"session_id"`
	StartedAt          string             `json:"started_at"`
	CompletedAt        string              `json:"completed_at"`
	Repository         string             `json:"repository"`
	Operator           string             `json:"operator"`
	KnowledgeAreasReviewed []KnowledgeDomain `json:"knowledge_areas_reviewed"`
	KnowledgeGapsIdentified []KnowledgeGap  `json:"knowledge_gaps_identified"`
	ArtifactsCollected []CollectedArtifact `json:"artifacts_collected"`
	ArtifactsUpdated   []CollectedArtifact `json:"artifacts_updated"`
	ArtifactsExisting  []CollectedArtifact `json:"artifacts_existing"`
	KnowledgeStillMissing []KnowledgeGap   `json:"knowledge_still_missing"`
	Statistics         CollectionStats    `json:"statistics"`
	Recommendations    []string           `json:"recommendations"`
}

// CollectionStats contains statistics about the collection process
type CollectionStats struct {
	TotalGapsIdentified   int     `json:"total_gaps_identified"`
	TotalArtifactsAdded    int     `json:"total_artifacts_added"`
	TotalArtifactsUpdated  int     `json:"total_artifacts_updated"`
	TotalArtifactsExisting int     `json:"total_artifacts_existing"`
	TotalKnowledgeCollected int    `json:"total_knowledge_collected"`
	ProcessingTime        float64 `json:"processing_time_seconds"`
	SuccessRate           float64 `json:"success_rate"`
}

// CollectionReport is a human-readable collection report
type CollectionReport struct {
	ReportID         string             `json:"report_id"`
	CollectionResult *CollectionResult   `json:"collection_result"`
	GeneratedAt      string             `json:"generated_at"`
}

// KnowledgeProvider is an interface for extensible knowledge sources
type KnowledgeProvider interface {
	Name() string
	Collect(input *CollectionInput) ([]CollectedArtifact, error)
	CanCollect(domain KnowledgeDomain) bool
}

// GapAnalyzer analyzes knowledge gaps from audit findings and normalization results
type GapAnalyzer interface {
	Analyze(input *CollectionInput) []KnowledgeGap
}

// NewCollectionResult creates a new collection result
func NewCollectionResult(sessionID, repoPath, operator string) *CollectionResult {
	return &CollectionResult{
		SessionID:         sessionID,
		StartedAt:         time.Now().Format(time.RFC3339),
		Repository:        repoPath,
		Operator:          operator,
		KnowledgeAreasReviewed: []KnowledgeDomain{},
		KnowledgeGapsIdentified: []KnowledgeGap{},
		ArtifactsCollected: []CollectedArtifact{},
		ArtifactsUpdated:   []CollectedArtifact{},
		ArtifactsExisting: []CollectedArtifact{},
		KnowledgeStillMissing: []KnowledgeGap{},
		Recommendations:   []string{},
	}
}

// SetCompleted marks the collection as complete
func (r *CollectionResult) SetCompleted() {
	r.CompletedAt = time.Now().Format(time.RFC3339)
}

// CalculateStats computes statistics from the collection
func (r *CollectionResult) CalculateStats(startTime time.Time) {
	r.Statistics = CollectionStats{
		TotalGapsIdentified:    len(r.KnowledgeGapsIdentified),
		TotalArtifactsAdded:     len(r.ArtifactsCollected),
		TotalArtifactsUpdated:   len(r.ArtifactsUpdated),
		TotalArtifactsExisting: len(r.ArtifactsExisting),
		TotalKnowledgeCollected: len(r.ArtifactsCollected) + len(r.ArtifactsUpdated),
		ProcessingTime:          time.Since(startTime).Seconds(),
		SuccessRate:            calculateSuccessRate(len(r.ArtifactsCollected), len(r.KnowledgeGapsIdentified)),
	}
}

func calculateSuccessRate(collected, gaps int) float64 {
	if gaps == 0 {
		return 100.0
	}
	return float64(collected) / float64(gaps) * 100
}

// GenerateArtifactID generates a unique artifact ID
func GenerateArtifactID(domain KnowledgeDomain) string {
	timestamp := time.Now().Format("20060102150405")
	return string(domain) + "-ART-" + timestamp
}

// GenerateGapID generates a unique gap ID
func GenerateGapID(domain KnowledgeDomain, sequence int) string {
	return string(domain) + "-GAP-" + time.Now().Format("20060102") + "-" + formatSeq(sequence)
}

func formatSeq(n int) string {
	if n < 10 {
		return "00" + string(rune('0'+n))
	}
	if n < 100 {
		return "0" + string(rune('0'+n/10)) + string(rune('0'+n%10))
	}
	return string(rune('0'+n/100)) + string(rune('0'+(n%100)/10)) + string(rune('0'+n%10))
}
