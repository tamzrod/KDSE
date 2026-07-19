// Package selfhost implements KDSE self-hosting capabilities.
// Enables KDSE to analyze and evolve its own architecture using evidence-driven methodology.
package selfhost

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// ComponentType represents the type of a runtime component
type ComponentType string

const (
	ComponentTypeModule    ComponentType = "module"
	ComponentTypeCommand   ComponentType = "command"
	ComponentTypeRuntime   ComponentType = "runtime"
	ComponentTypeToolchain ComponentType = "toolchain"
	ComponentTypeKnowledge ComponentType = "knowledge"
	ComponentTypeOrchestration ComponentType = "orchestration"
)

// DependencyType represents the type of dependency relationship
type DependencyType string

const (
	DependencyTypeImport   DependencyType = "import"
	DependencyTypeInvoke  DependencyType = "invoke"
	DependencyTypeConfig  DependencyType = "config"
	DependencyTypeData    DependencyType = "data"
)

// Component represents a single component in the runtime architecture
type Component struct {
	Name         string           `json:"name"`
	Type         ComponentType    `json:"type"`
	Purpose      string           `json:"purpose"`
	Path         string           `json:"path"`
	Dependencies []string         `json:"dependencies,omitempty"`
	Provides     []string         `json:"provides,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// Dependency represents a dependency relationship between components
type Dependency struct {
	Source      string        `json:"source"`
	Target      string        `json:"target"`
	Type        DependencyType `json:"type"`
	Description string        `json:"description,omitempty"`
}

// DataFlow represents a data flow between components
type DataFlow struct {
	From        string `json:"from"`
	To          string `json:"to"`
	DataType    string `json:"data_type"`
	Description string `json:"description,omitempty"`
}

// ArchitectureModel represents the complete architecture of KDSE
type ArchitectureModel struct {
	Version       string                 `json:"version"`
	GeneratedAt   string                 `json:"generated_at"`
	Components    map[string]*Component  `json:"components"`
	Dependencies  []*Dependency          `json:"dependencies"`
	DataFlows     []*DataFlow            `json:"data_flows,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Summary       *ArchitectureSummary  `json:"summary,omitempty"`
}

// ArchitectureSummary provides statistics about the architecture
type ArchitectureSummary struct {
	TotalComponents   int            `json:"total_components"`
	TotalDependencies int            `json:"total_dependencies"`
	ByType            map[string]int `json:"by_type"`
	Cycles            [][]string     `json:"cycles,omitempty"`
	Depth             int            `json:"max_depth"`
}

// SelfAssessmentReport represents the result of self-assessment
type SelfAssessmentReport struct {
	Timestamp       string              `json:"timestamp"`
	Version         string              `json:"version"`
	Architecture    *ArchitectureModel  `json:"architecture"`
	HealthStatus    *HealthStatus       `json:"health_status"`
	Dependencies    *DependencyAnalysis `json:"dependencies"`
	Recommendations []string            `json:"recommendations,omitempty"`
}

// HealthStatus represents the health assessment of the runtime
type HealthStatus struct {
	Overall        string             `json:"overall"`
	Score          float64            `json:"score"` // 0.0 - 1.0
	Checks         []*HealthCheck     `json:"checks"`
	CriticalIssues []string           `json:"critical_issues,omitempty"`
	Warnings       []string           `json:"warnings,omitempty"`
}

// HealthCheck represents a single health check
type HealthCheck struct {
	Name        string `json:"name"`
	Status      string `json:"status"` // "PASS", "FAIL", "WARN"
	Description string `json:"description"`
	Evidence    string `json:"evidence,omitempty"`
}

// DependencyAnalysis represents the dependency analysis results
type DependencyAnalysis struct {
	Total         int                  `json:"total"`
	Direct        int                  `json:"direct"`
	Transitive    int                  `json:"transitive"`
	Violations    []*DependencyViolation `json:"violations,omitempty"`
	Duplications  []string             `json:"duplications,omitempty"`
}

// DependencyViolation represents a dependency rule violation
type DependencyViolation struct {
	Rule        string `json:"rule"`
	Source      string `json:"source"`
	Target      string `json:"target"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "error", "warning"
}

// NewArchitectureModel creates a new empty architecture model
func NewArchitectureModel() *ArchitectureModel {
	return &ArchitectureModel{
		Version:      "1.0",
		GeneratedAt: time.Now().Format(time.RFC3339),
		Components:  make(map[string]*Component),
		Dependencies: []*Dependency{},
		DataFlows:   []*DataFlow{},
		Metadata:    make(map[string]interface{}),
	}
}

// AddComponent adds a component to the model
func (m *ArchitectureModel) AddComponent(c *Component) {
	m.Components[c.Name] = c
}

// AddDependency adds a dependency to the model
func (m *ArchitectureModel) AddDependency(d *Dependency) {
	m.Dependencies = append(m.Dependencies, d)
}

// AddDataFlow adds a data flow to the model
func (m *ArchitectureModel) AddDataFlow(f *DataFlow) {
	m.DataFlows = append(m.DataFlows, f)
}

// CalculateSummary computes statistics for the model
func (m *ArchitectureModel) CalculateSummary() {
	summary := &ArchitectureSummary{
		ByType: make(map[string]int),
	}

	for _, c := range m.Components {
		summary.TotalComponents++
		summary.ByType[string(c.Type)]++
	}

	summary.TotalDependencies = len(m.Dependencies)
	summary.Depth = m.calculateMaxDepth()

	m.Summary = summary
}

// calculateMaxDepth calculates the maximum dependency depth
func (m *ArchitectureModel) calculateMaxDepth() int {
	maxDepth := 0
	visited := make(map[string]bool)

	var dfs func(name string, depth int)
	dfs = func(name string, depth int) {
		if visited[name] {
			return
		}
		visited[name] = true
		if depth > maxDepth {
			maxDepth = depth
		}

		for _, d := range m.Dependencies {
			if d.Source == name {
				dfs(d.Target, depth+1)
			}
		}
	}

	for name := range m.Components {
		visited = make(map[string]bool)
		dfs(name, 0)
	}

	return maxDepth
}

// Save saves the architecture model to a file
func (m *ArchitectureModel) Save(path string) error {
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

// Load loads an architecture model from a file
func LoadArchitectureModel(path string) (*ArchitectureModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var model ArchitectureModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// FormatModelSummary returns a human-readable summary of the model
func FormatModelSummary(m *ArchitectureModel) string {
	summary := ""
	summary += "═══════════════════════════════════════════════════════════════\n"
	summary += "                    KDSE Architecture Model                     \n"
	summary += "═══════════════════════════════════════════════════════════════\n"
	summary += "  Version: " + m.Version + "\n"
	summary += "  Generated: " + m.GeneratedAt + "\n"
	summary += "───────────────────────────────────────────────────────────────\n"
	summary += "  Components: " + string(rune('0'+len(m.Components)%10)) + "\n"
	for t, count := range m.Summary.ByType {
		summary += "    • " + t + ": " + string(rune('0'+count%10)) + "\n"
	}
	summary += "  Dependencies: " + string(rune('0'+m.Summary.TotalDependencies%10)) + "\n"
	summary += "  Max Depth: " + string(rune('0'+m.Summary.Depth%10)) + "\n"
	summary += "═══════════════════════════════════════════════════════════════\n"
	return summary
}
