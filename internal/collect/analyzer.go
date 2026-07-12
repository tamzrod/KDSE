package collect

import (
	"fmt"
	"strings"

	"github.com/kdse/runtime/internal/types"
)

// GapDetector analyzes repository state and identifies knowledge gaps
type GapDetector struct {
	repoPath string
}

// NewGapDetector creates a new gap detector
func NewGapDetector(repoPath string) *GapDetector {
	return &GapDetector{repoPath: repoPath}
}

// DetectGaps analyzes the input and returns identified knowledge gaps
func (d *GapDetector) DetectGaps(input *CollectionInput) []KnowledgeGap {
	var gaps []KnowledgeGap
	seq := 1

	// Analyze from audit findings
	if len(input.AuditFindings) > 0 {
		for _, finding := range input.AuditFindings {
			domain := d.inferDomainFromFinding(finding)
			gap := KnowledgeGap{
				ID:           GenerateGapID(domain, seq),
				Domain:       domain,
				Topic:        d.extractTopic(finding.Title),
				Description:  finding.Details,
				Severity:     d.mapSeverity(finding.Severity),
				Impact:       d.assessImpact(finding),
				Recommendations: d.suggestRecommendations(domain, finding),
			}
			gaps = append(gaps, gap)
			seq++
		}
	}

	// Analyze from session state dimensions
	if input.SessionState != nil && len(input.SessionState.Dimensions) > 0 {
		for dim, score := range input.SessionState.Dimensions {
			if score < 5.0 {
				domain := d.domainFromDimension(dim)
				gap := KnowledgeGap{
					ID:              GenerateGapID(domain, seq),
					Domain:          domain,
					Topic:           fmt.Sprintf("%s Knowledge", dim),
					Description:     fmt.Sprintf("The %s dimension has a low score (%.1f/10), indicating missing or incomplete knowledge artifacts.", dim, score),
					Severity:        d.severityFromScore(score),
					Impact:          fmt.Sprintf("Low %s score affects overall project maturity assessment", dim),
					Recommendations: d.recommendForDimension(domain, dim),
				}
				gaps = append(gaps, gap)
				seq++
			}
		}
	}

	// Add domain-specific gaps based on knowledge areas requested
	if len(input.KnowledgeAreas) > 0 {
		for _, area := range input.KnowledgeAreas {
			if !d.hasGapForDomain(gaps, area) {
				gap := KnowledgeGap{
					ID:              GenerateGapID(area, seq),
					Domain:          area,
					Topic:           string(area) + " Knowledge",
					Description:     fmt.Sprintf("Requested knowledge collection for %s domain", area),
					Severity:        "Medium",
					Impact:          fmt.Sprintf("Missing %s knowledge may affect implementation quality", area),
					Recommendations: d.recommendForDomain(area),
				}
				gaps = append(gaps, gap)
				seq++
			}
		}
	}

	// Add common engineering knowledge gaps if none identified
	if len(gaps) == 0 {
		gaps = append(gaps, KnowledgeGap{
			ID:              GenerateGapID(DomainGeneral, seq),
			Domain:          DomainGeneral,
			Topic:           "General Engineering Knowledge",
			Description:     "No specific knowledge gaps identified. Collecting foundational engineering knowledge.",
			Severity:        "Low",
			Impact:          "General knowledge collection to establish baseline",
			Recommendations: []string{"Review existing documentation", "Identify domain-specific requirements"},
		})
	}

	return gaps
}

// inferDomainFromFinding infers the knowledge domain from a finding
func (d *GapDetector) inferDomainFromFinding(finding types.Finding) KnowledgeDomain {
	titleLower := strings.ToLower(finding.Title)
	detailsLower := strings.ToLower(finding.Details)

	// Check for specific domain indicators
	if containsAny(titleLower, "transformer", "magnetic", "inductance", "flux") {
		return DomainTransformers
	}
	if containsAny(titleLower, "battery", "soc", "soh", "charge", "discharge") {
		return DomainBattery
	}
	if containsAny(titleLower, "relay", "protection", "overcurrent", "circuit breaker") {
		return DomainRelay
	}
	if containsAny(titleLower, "weather", "solar", "irradiance", "temperature", "wind") {
		return DomainWeather
	}
	if containsAny(titleLower, "physics", "thermal", "mechanical", "electrical") {
		return DomainPhysics
	}
	if containsAny(titleLower, "equipment", "device", "hardware", "component") {
		return DomainEquipment
	}
	if containsAny(titleLower, "standard", "iec", "ieee", "iso", "regulation") {
		return DomainStandards
	}
	if containsAny(titleLower, "protocol", "communication", "interface", "api") {
		return DomainProtocols
	}
	if containsAny(titleLower, "control", "algorithm", "logic", "pid") {
		return DomainControl
	}
	if containsAny(titleLower, "simulation", "model", "digital twin") {
		return DomainSimulation
	}
	if containsAny(titleLower, "business", "requirement", "process") {
		return DomainBusiness
	}
	if containsAny(titleLower, "glossary", "vocabulary", "terminology") {
		return DomainVocabulary
	}

	// Check dimension
	dimLower := strings.ToLower(finding.Dimension)
	if containsAny(dimLower, "knowledge", "documentation", "artifacts") {
		return DomainGeneral
	}
	if containsAny(dimLower, "architecture", "design") {
		return DomainEquipment
	}
	if containsAny(dimLower, "verification", "testing", "validation") {
		return DomainStandards
	}

	return DomainGeneral
}

func containsAny(s string, parts ...string) bool {
	for _, part := range parts {
		if strings.Contains(s, part) {
			return true
		}
	}
	return false
}

// extractTopic extracts a topic from a finding title
func (d *GapDetector) extractTopic(title string) string {
	// Clean up the title to form a topic
	topic := strings.TrimSpace(title)
	topic = strings.TrimPrefix(topic, "[")
	topic = strings.TrimSuffix(topic, "]")
	
	// Capitalize first letter
	if len(topic) > 0 {
		topic = strings.ToUpper(string(topic[0])) + topic[1:]
	}
	
	return topic
}

// mapSeverity maps a finding severity to a gap severity
func (d *GapDetector) mapSeverity(severity string) string {
	switch strings.ToLower(severity) {
	case "critical":
		return "Critical"
	case "high":
		return "High"
	case "medium":
		return "Medium"
	default:
		return "Low"
	}
}

// assessImpact assesses the impact of a finding
func (d *GapDetector) assessImpact(finding types.Finding) string {
	switch strings.ToLower(finding.Severity) {
	case "critical":
		return "Critical impact on implementation quality and safety"
	case "high":
		return "High impact on project delivery and reliability"
	case "medium":
		return "Moderate impact on development efficiency"
	default:
		return "Low impact on overall project"
	}
}

// suggestRecommendations suggests recommendations for a finding
func (d *GapDetector) suggestRecommendations(domain KnowledgeDomain, finding types.Finding) []string {
	recs := []string{
		"Collect relevant documentation from authoritative sources",
		"Document the knowledge in KDSE-standard format",
	}
	
	switch domain {
	case DomainTransformers:
		recs = append(recs, "Research transformer specifications and behavior models")
	case DomainBattery:
		recs = append(recs, "Obtain battery manufacturer documentation and safety guidelines")
	case DomainRelay:
		recs = append(recs, "Document relay protection coordination and settings")
	case DomainWeather:
		recs = append(recs, "Collect historical weather data and forecasting models")
	case DomainPhysics:
		recs = append(recs, "Document applicable physics principles and calculations")
	case DomainStandards:
		recs = append(recs, "Obtain relevant standards documentation (IEC, IEEE, etc.)")
	case DomainProtocols:
		recs = append(recs, "Document communication protocols and interface specifications")
	case DomainControl:
		recs = append(recs, "Collect control algorithm documentation and tuning parameters")
	case DomainSimulation:
		recs = append(recs, "Obtain simulation models and validation data")
	}
	
	return recs
}

// domainFromDimension maps a dimension name to a knowledge domain
func (d *GapDetector) domainFromDimension(dimension string) KnowledgeDomain {
	dimLower := strings.ToLower(dimension)
	
	if containsAny(dimLower, "knowledge", "documentation") {
		return DomainGeneral
	}
	if containsAny(dimLower, "architecture", "design") {
		return DomainEquipment
	}
	if containsAny(dimLower, "verification", "testing") {
		return DomainStandards
	}
	if containsAny(dimLower, "implementation", "code") {
		return DomainControl
	}
	if containsAny(dimLower, "traceability") {
		return DomainBusiness
	}
	if containsAny(dimLower, "governance") {
		return DomainStandards
	}
	
	return DomainGeneral
}

// severityFromScore determines severity based on score
func (d *GapDetector) severityFromScore(score float64) string {
	switch {
	case score < 2:
		return "Critical"
	case score < 4:
		return "High"
	case score < 6:
		return "Medium"
	default:
		return "Low"
	}
}

// recommendForDimension gives recommendations for a dimension gap
func (d *GapDetector) recommendForDimension(domain KnowledgeDomain, dimension string) []string {
	return []string{
		fmt.Sprintf("Collect or create knowledge artifacts for %s", dimension),
		"Review existing documentation for relevant content",
		"Document implicit knowledge in KDSE-standard format",
	}
}

// recommendForDomain gives recommendations for a domain
func (d *GapDetector) recommendForDomain(domain KnowledgeDomain) []string {
	switch domain {
	case DomainPhysics:
		return []string{
			"Document fundamental physics principles",
			"Collect applicable calculation methods",
			"Gather safety factors and design margins",
		}
	case DomainEquipment:
		return []string{
			"Collect equipment specifications",
			"Document operating procedures",
			"Gather maintenance requirements",
		}
	case DomainEnvironment:
		return []string{
			"Document environmental conditions",
			"Collect site-specific data",
			"Gather regulatory requirements",
		}
	case DomainStandards:
		return []string{
			"Obtain applicable standards documents",
			"Document compliance requirements",
			"Collect certification evidence",
		}
	case DomainBusiness:
		return []string{
			"Document business rules and logic",
			"Collect process definitions",
			"Gather requirement specifications",
		}
	case DomainSimulation:
		return []string{
			"Obtain simulation models",
			"Document validation procedures",
			"Collect test scenarios",
		}
	case DomainControl:
		return []string{
			"Document control algorithms",
			"Collect tuning parameters",
			"Gather operating modes",
		}
	case DomainProtocols:
		return []string{
			"Document communication protocols",
			"Collect interface specifications",
			"Gather error handling procedures",
		}
	case DomainVocabulary:
		return []string{
			"Document domain terminology",
			"Collect glossary entries",
			"Gather acronym definitions",
		}
	default:
		return []string{
			"Review existing documentation",
			"Identify knowledge sources",
			"Document findings in KDSE format",
		}
	}
}

// hasGapForDomain checks if a gap already exists for the given domain
func (d *GapDetector) hasGapForDomain(gaps []KnowledgeGap, domain KnowledgeDomain) bool {
	for _, gap := range gaps {
		if gap.Domain == domain {
			return true
		}
	}
	return false
}
