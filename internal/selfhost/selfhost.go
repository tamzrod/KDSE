// Package selfhost implements KDSE self-hosting capabilities.
// This module enables KDSE to analyze and evolve its own architecture using
// the same evidence-driven methodology it uses for software projects.
package selfhost

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager is the main entry point for self-hosting operations
type Manager struct {
	repoPath string
	kdsePath string
}

// NewManager creates a new self-hosting manager
func NewManager(repoPath string) *Manager {
	return &Manager{
		repoPath: repoPath,
		kdsePath: filepath.Join(repoPath, ".kdse"),
	}
}

// Initialize prepares the self-hosting runtime
func (m *Manager) Initialize() error {
	// Ensure self-hosting directories exist
	dirs := []string{
		filepath.Join(m.kdsePath, "runtime"),
		filepath.Join(m.kdsePath, "runtime", "promotion"),
		filepath.Join(m.kdsePath, "runtime", "promotion", "snapshots"),
		filepath.Join(m.kdsePath, "reports"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// RunSelfAssessment performs a complete self-assessment
func (m *Manager) RunSelfAssessment() (*SelfAssessmentReport, error) {
	analyzer := NewAnalyzer(m.repoPath)
	return analyzer.Analyze()
}

// AnalyzeImpact analyzes the impact of proposed changes
func (m *Manager) AnalyzeImpact(changes []*Change) (*ImpactAnalysisResult, error) {
	report, err := m.RunSelfAssessment()
	if err != nil {
		return nil, err
	}

	analyzer := NewImpactAnalyzer(report.Architecture)
	return analyzer.AnalyzeMultiple(changes), nil
}

// StartEvolution starts a new evolution workflow
func (m *Manager) StartEvolution() (*EvolutionWorkflow, error) {
	if err := m.Initialize(); err != nil {
		return nil, err
	}

	workflow := NewEvolutionWorkflow(m.repoPath)
	if err := workflow.Start(); err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetPromotionManager returns the promotion manager
func (m *Manager) GetPromotionManager() *PromotionManager {
	manager := NewPromotionManager(m.repoPath)
	manager.Initialize()
	return manager
}

// LoadArchitectureModel loads the saved architecture model
func (m *Manager) LoadArchitectureModel() (*ArchitectureModel, error) {
	path := filepath.Join(m.kdsePath, "runtime", "architecture-model.json")
	return LoadArchitectureModel(path)
}

// SaveArchitectureModel saves the architecture model
func (m *Manager) SaveArchitectureModel(model *ArchitectureModel) error {
	path := filepath.Join(m.kdsePath, "runtime", "architecture-model.json")
	return model.Save(path)
}

// Version returns the self-hosting module version
const Version = "1.0"
