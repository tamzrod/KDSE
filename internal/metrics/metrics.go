// Package metrics provides basic metrics collection for KDSE workspace.
package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kdse/runtime/internal/knowledge"
	"github.com/kdse/runtime/internal/workspace"
)

// Metrics captures workspace metrics for reporting
type Metrics struct {
	GeneratedAt  string          `json:"generated_at"`
	Repository   string          `json:"repository"`
	WorkspaceOK  bool            `json:"workspace_ok"`
	Knowledge    KnowledgeMetrics `json:"knowledge"`
	Traceability TraceMetrics    `json:"traceability"`
}

// KnowledgeMetrics captures knowledge-related metrics
type KnowledgeMetrics struct {
	Total       int `json:"total"`
	Notebook    int `json:"notebook"`
	Candidate   int `json:"candidate"`
	Promoted    int `json:"promoted"`
	Rejected    int `json:"rejected"`
	AvgStrength int `json:"avg_evidence_strength"`
}

// TraceMetrics captures traceability metrics
type TraceMetrics struct {
	KnowledgeWithEvidence int `json:"knowledge_with_evidence"`
	CoveragePercent      int `json:"coverage_percent"`
}

// Collector collects metrics from the workspace
type Collector struct {
	repoPath string
}

// NewCollector creates a new metrics collector
func NewCollector(repoPath string) *Collector {
	return &Collector{repoPath: repoPath}
}

// Collect gathers all metrics from the workspace
func (c *Collector) Collect() (*Metrics, error) {
	m := &Metrics{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Repository:  c.repoPath,
	}

	ws := workspace.New(c.repoPath)
	m.WorkspaceOK = ws.Exists()

	// Collect knowledge metrics
	km := knowledge.NewManager(c.repoPath)
	if err := km.Load(); err == nil {
		stats := km.Stats()
		m.Knowledge.Total = stats["total"]
		m.Knowledge.Notebook = stats["notebook"]
		m.Knowledge.Candidate = stats["candidate"]
		m.Knowledge.Promoted = stats["promoted"]
		m.Knowledge.Rejected = stats["rejected"]

		// Calculate average evidence strength
		var totalStrength, count int
		for _, e := range km.List("") {
			if e.EvidenceStrength > 0 {
				totalStrength += int(e.EvidenceStrength)
				count++
			}
		}
		if count > 0 {
			m.Knowledge.AvgStrength = totalStrength / count
		}

		// Traceability coverage
		for _, e := range km.List("") {
			if e.Status == knowledge.StatusPromoted && len(e.EvidenceRefs) > 0 {
				m.Traceability.KnowledgeWithEvidence++
			}
		}
		if m.Knowledge.Promoted > 0 {
			m.Traceability.CoveragePercent = (m.Traceability.KnowledgeWithEvidence * 100) / m.Knowledge.Promoted
		}
	}

	return m, nil
}

// Save saves metrics to the reports directory
func (m *Metrics) Save(repoPath string) error {
	path := filepath.Join(repoPath, ".kdse", "reports", "metrics.json")

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Format formats metrics for display
func (m *Metrics) Format() string {
	return fmt.Sprintf(`KDSE Metrics Report
==================

Workspace Status: %s

Knowledge Pipeline:
  Total Entries:   %d
  Notebook:        %d
  Candidate:       %d
  Promoted:        %d
  Rejected:        %d
  Avg Strength:    %d/5

Traceability:
  With Evidence:  %d
  Coverage:        %d%%

Generated: %s`, boolToStatus(m.WorkspaceOK),
		m.Knowledge.Total,
		m.Knowledge.Notebook,
		m.Knowledge.Candidate,
		m.Knowledge.Promoted,
		m.Knowledge.Rejected,
		m.Knowledge.AvgStrength,
		m.Traceability.KnowledgeWithEvidence,
		m.Traceability.CoveragePercent,
		m.GeneratedAt)
}

func boolToStatus(b bool) string {
	if b {
		return "Operational"
	}
	return "Not Initialized"
}
