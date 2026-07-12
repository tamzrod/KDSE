package normalize

import (
	"regexp"
	"strings"
)

// Extractor handles knowledge extraction from documentation
type Extractor struct{}

// NewExtractor creates a new knowledge extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractKnowledge extracts structured knowledge from discovered documents
func (e *Extractor) ExtractKnowledge(docs []DiscoveredDocument) KnowledgeExtraction {
	extraction := KnowledgeExtraction{
		Requirements: []ExtractedReq{},
		Decisions:   []ExtractedDecision{},
		Glossary:    []GlossaryTerm{},
	}

	// Extract from all documents
	for _, doc := range docs {
		e.extractFromDocument(&extraction, doc)
	}

	// Deduplicate and set domain
	extraction.Deduplicate()

	return extraction
}

func (e *Extractor) extractFromDocument(extraction *KnowledgeExtraction, doc DiscoveredDocument) {
	// Read the actual content
	content := e.readDocumentContent(doc.Path)

	switch doc.Type {
	case DocsTypeREADME:
		e.extractFromREADME(extraction, content, doc.Path)
	case DocsTypeArchitecture, DocsTypeDesign:
		e.extractFromArchitecture(extraction, content, doc.Path)
	case DocsTypeADR:
		e.extractFromADR(extraction, content, doc.Path)
	case DocsTypeAPI:
		e.extractFromAPI(extraction, content, doc.Path)
	default:
		e.extractGeneralKnowledge(extraction, content, doc.Path)
	}
}

func (e *Extractor) readDocumentContent(path string) string {
	// This would read the file - simplified for now
	return ""
}

func (e *Extractor) extractFromREADME(extraction *KnowledgeExtraction, content, source string) {
	// Extract purpose/overview
	if matches := e.findPattern(content, `(?im)^#?\s*(?:about|overview|introduction|purpose)[:\s]*\n?(.*?)(?:\n\n|\z)`); len(matches) > 0 {
		extraction.Purpose = strings.TrimSpace(matches[0])
	}

	// Extract domain
	if matches := e.findPattern(content, `(?im)(?:language|framework|platform|tool):\s*(\w+)`); len(matches) > 0 {
		extraction.Domain = matches[0]
	}

	// Extract stakeholders
	stakeholders := e.findPattern(content, `(?im)(?:author|maintainer|contributor)s?[:\s]*([^\n]+)`)
	for _, s := range stakeholders {
		parts := strings.Split(s, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				extraction.Stakeholders = append(extraction.Stakeholders, p)
			}
		}
	}

	// Extract requirements
	requirements := e.extractRequirements(content, source)
	extraction.Requirements = append(extraction.Requirements, requirements...)
}

func (e *Extractor) extractFromArchitecture(extraction *KnowledgeExtraction, content, source string) {
	// Extract architectural decisions
	decisions := e.extractDecisions(content, source)
	extraction.Decisions = append(extraction.Decisions, decisions...)

	// Extract constraints
	constraints := e.findPattern(content, `(?im)(?:constraint|limitation)s?[:\s]+([^\n]+)`)
	for _, c := range constraints {
		c = strings.TrimSpace(c)
		if c != "" {
			extraction.Constraints = append(extraction.Constraints, c)
		}
	}

	// Extract glossary terms
	glossary := e.extractGlossary(content, source)
	extraction.Glossary = append(extraction.Glossary, glossary...)
}

func (e *Extractor) extractFromADR(extraction *KnowledgeExtraction, content, source string) {
	decisions := e.extractDecisionsFromADR(content, source)
	extraction.Decisions = append(extraction.Decisions, decisions...)
}

func (e *Extractor) extractFromAPI(extraction *KnowledgeExtraction, content, source string) {
	// Extract endpoint requirements
	endpoints := e.findPattern(content, `(?im)(?:GET|POST|PUT|DELETE|PATCH)\s+([^\s\n]+)`)
	for i, ep := range endpoints {
		extraction.Requirements = append(extraction.Requirements, ExtractedReq{
			ID:          e.generateReqID("API", i),
			Description: "API endpoint: " + ep,
			Source:      source,
			Type:        "API",
		})
	}
}

func (e *Extractor) extractGeneralKnowledge(extraction *KnowledgeExtraction, content, source string) {
	// Extract requirements from general documentation
	requirements := e.extractRequirements(content, source)
	extraction.Requirements = append(extraction.Requirements, requirements...)

	// Extract glossary from any document
	glossary := e.extractGlossary(content, source)
	extraction.Glossary = append(extraction.Glossary, glossary...)
}

func (e *Extractor) extractRequirements(content, source string) []ExtractedReq {
	var requirements []ExtractedReq

	// Look for requirement patterns
	patterns := []struct {
		pattern string
		prefix  string
	}{
		{`(?im)^\s*[-*]\s*\[?\s*(?:shall|must|will|should|can)\s+(.+)`, "IMPERATIVE"},
		{`(?im)^\s*\d+[.)]\s*(?:shall|must|will|should|can)\s+(.+)`, "NUMBERED"},
		{`(?im)(?:requirement)s?[:\s]+([^\n]+)`, "LABELED"},
	}

	for i, p := range patterns {
		matches := e.findPattern(content, p.pattern)
		for _, m := range matches {
			requirements = append(requirements, ExtractedReq{
				ID:          e.generateReqID("REQ", i),
				Description: strings.TrimSpace(m),
				Source:      source,
				Type:        p.prefix,
			})
		}
	}

	return requirements
}

func (e *Extractor) extractDecisions(content, source string) []ExtractedDecision {
	var decisions []ExtractedDecision

	// Look for decision patterns
	contextMatches := e.findPattern(content, `(?im)(?:context|background)[:\s]+([^\n]+)`)
	decisionMatches := e.findPattern(content, `(?im)(?:decision|solution)[:\s]+([^\n]+)`)
	rationaleMatches := e.findPattern(content, `(?im)(?:rationale|reason)[:\s]+([^\n]+)`)

	for i := range contextMatches {
		decision := ExtractedDecision{
			ID:        e.generateDecisionID(i),
			Context:   strings.TrimSpace(contextMatches[i]),
			Decision:  getOrDefault(decisionMatches, i),
			Rationale: getOrDefault(rationaleMatches, i),
			Source:    source,
			Status:    "Proposed",
		}
		decisions = append(decisions, decision)
	}

	return decisions
}

func (e *Extractor) extractDecisionsFromADR(content, source string) []ExtractedDecision {
	var decisions []ExtractedDecision

	// Extract ADR components
	title := e.findFirstPattern(content, `(?im)^#\s+(.+)`)
	context := e.findFirstPattern(content, `(?im)(?:##?\s*)?context[:\s]+([^\n]+)`)
	decision := e.findFirstPattern(content, `(?im)(?:##?\s*)?decision[:\s]+([^\n]+)`)
	rationale := e.findFirstPattern(content, `(?im)(?:##?\s*)?rationale[:\s]+([^\n]+)`)
	status := e.findFirstPattern(content, `(?im)(?:##?\s*)?status[:\s]*\n?(\w+)`)

	if context != "" || decision != "" {
		decision := ExtractedDecision{
			ID:        e.extractADRID(source),
			Title:     title,
			Context:   strings.TrimSpace(context),
			Decision:  strings.TrimSpace(decision),
			Rationale: strings.TrimSpace(rationale),
			Source:    source,
			Status:    e.normalizeStatus(status),
		}
		decisions = append(decisions, decision)
	}

	return decisions
}

func (e *Extractor) extractGlossary(content, source string) []GlossaryTerm {
	var terms []GlossaryTerm

	// Look for glossary section
	glossaryContent := e.findFirstPattern(content, `(?im)(?:glossary|terminology|definitions)[:\s]*\n((?:.|\n)*?)(?:\n##|\n#|$)`)
	if glossaryContent == "" {
		// Look for term definitions inline
		matches := e.findPattern(content, `(?im)^\s*[*_-]\s*\*\*([^*]+)\*\*[:\s]+([^\n]+)`)
		for _, m := range matches {
			parts := strings.SplitN(m, "**", 2)
			if len(parts) == 2 {
				terms = append(terms, GlossaryTerm{
					Term:       strings.TrimSpace(parts[0]),
					Definition: strings.TrimSpace(parts[1]),
					Source:     source,
				})
			}
		}
	} else {
		// Parse glossary entries
		entries := e.findPattern(glossaryContent, `(?im)^\s*[*_-]?\s*([^\s:]+)[:\s]+([^\n]+)`)
		for _, entry := range entries {
			parts := strings.SplitN(entry, ":", 2)
			if len(parts) == 2 {
				terms = append(terms, GlossaryTerm{
					Term:       strings.TrimSpace(parts[0]),
					Definition: strings.TrimSpace(parts[1]),
					Source:     source,
				})
			}
		}
	}

	return terms
}

func (e *Extractor) findPattern(content, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}
	return results
}

func (e *Extractor) findFirstPattern(content, pattern string) string {
	matches := e.findPattern(content, pattern)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func (e *Extractor) generateReqID(prefix string, index int) string {
	return "NORM-REQ-" + prefix + "-" + formatIndex(index)
}

func (e *Extractor) generateDecisionID(index int) string {
	return "NORM-DEC-" + formatIndex(index)
}

func (e *Extractor) extractADRID(source string) string {
	base := strings.TrimSuffix(strings.TrimPrefix(source, "adr/"), ".md")
	base = strings.TrimPrefix(base, "adr_")
	base = strings.TrimPrefix(base, "ADR-")
	base = strings.TrimPrefix(base, "ADR_")
	base = strings.ToUpper(base)
	base = regexp.MustCompile(`[^A-Z0-9]+`).ReplaceAllString(base, "-")
	return "ADR-" + base
}

func (e *Extractor) normalizeStatus(status string) string {
	status = strings.ToLower(strings.TrimSpace(status))
	switch {
	case strings.Contains(status, "accepted"):
		return "Accepted"
	case strings.Contains(status, "deprecated"):
		return "Deprecated"
	case strings.Contains(status, "superseded"):
		return "Superseded"
	case strings.Contains(status, "proposed"):
		return "Proposed"
	default:
		return "Proposed"
	}
}

func formatIndex(index int) string {
	return strings.ToUpper(string(rune('A'+index%26))) + string(rune('0'+index/26))
}

func getOrDefault(slice []string, index int) string {
	if index < len(slice) {
		return slice[index]
	}
	return ""
}

// Deduplicate removes duplicate entries from the extraction
func (k *KnowledgeExtraction) Deduplicate() {
	// Deduplicate stakeholders
	stakeholderMap := make(map[string]bool)
	var uniqueStakeholders []string
	for _, s := range k.Stakeholders {
		lower := strings.ToLower(s)
		if !stakeholderMap[lower] {
			stakeholderMap[lower] = true
			uniqueStakeholders = append(uniqueStakeholders, s)
		}
	}
	k.Stakeholders = uniqueStakeholders

	// Deduplicate requirements by description
	reqMap := make(map[string]bool)
	var uniqueReqs []ExtractedReq
	for _, r := range k.Requirements {
		lower := strings.ToLower(r.Description)
		if !reqMap[lower] {
			reqMap[lower] = true
			uniqueReqs = append(uniqueReqs, r)
		}
	}
	k.Requirements = uniqueReqs

	// Deduplicate decisions by context
	decisionMap := make(map[string]bool)
	var uniqueDecisions []ExtractedDecision
	for _, d := range k.Decisions {
		lower := strings.ToLower(d.Context + d.Decision)
		if !decisionMap[lower] {
			decisionMap[lower] = true
			uniqueDecisions = append(uniqueDecisions, d)
		}
	}
	k.Decisions = uniqueDecisions

	// Deduplicate glossary by term
	glossaryMap := make(map[string]bool)
	var uniqueGlossary []GlossaryTerm
	for _, g := range k.Glossary {
		lower := strings.ToLower(g.Term)
		if !glossaryMap[lower] {
			glossaryMap[lower] = true
			uniqueGlossary = append(uniqueGlossary, g)
		}
	}
	k.Glossary = uniqueGlossary
}
