// Package selfhost implements KDSE self-hosting capabilities.
package selfhost

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PromotionStage represents the promotion stage
type PromotionStage string

const (
	StageStaging PromotionStage = "staging"
	StageCanary  PromotionStage = "canary"
	StageStable  PromotionStage = "stable"
	StageRolledBack PromotionStage = "rolled_back"
)

// Promotion represents a promotion event
type Promotion struct {
	ID            string          `json:"id"`
	FromStage     PromotionStage  `json:"from_stage"`
	ToStage       PromotionStage  `json:"to_stage"`
	Timestamp     string          `json:"timestamp"`
	Reason        string          `json:"reason"`
	Status        string          `json:"status"` // "pending", "success", "failed", "rolled_back"
	RollbackOf    string          `json:"rollback_of,omitempty"`
	Evidence      []string        `json:"evidence,omitempty"`
	ApprovedBy    string          `json:"approved_by,omitempty"`
}

// PromotionManager manages staged promotions
type PromotionManager struct {
	repoPath    string
	kdsePath     string
	history      []*Promotion
	currentStage PromotionStage
	snapshots    map[string]*Snapshot
}

// Snapshot represents a point-in-time state
type Snapshot struct {
	ID          string                 `json:"id"`
	Stage       PromotionStage         `json:"stage"`
	Timestamp   string                 `json:"timestamp"`
	Model       *ArchitectureModel     `json:"model,omitempty"`
	State       *WorkflowState         `json:"state,omitempty"`
	Changes     []*Change              `json:"changes,omitempty"`
	Description string                 `json:"description"`
}

// PromotionResult represents the result of a promotion attempt
type PromotionResult struct {
	Success      bool           `json:"success"`
	Promotion    *Promotion     `json:"promotion,omitempty"`
	FromStage    PromotionStage `json:"from_stage"`
	ToStage      PromotionStage `json:"to_stage"`
	Snapshot     *Snapshot      `json:"snapshot,omitempty"`
	Errors       []string       `json:"errors,omitempty"`
	Warnings     []string       `json:"warnings,omitempty"`
}

// NewPromotionManager creates a new promotion manager
func NewPromotionManager(repoPath string) *PromotionManager {
	return &PromotionManager{
		repoPath:     repoPath,
		kdsePath:     filepath.Join(repoPath, ".kdse"),
		history:      []*Promotion{},
		currentStage: StageStable,
		snapshots:    make(map[string]*Snapshot),
	}
}

// Initialize initializes the promotion system
func (m *PromotionManager) Initialize() error {
	dir := filepath.Join(m.kdsePath, "runtime", "promotion")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Load existing history
	return m.loadHistory()
}

// Promote attempts to promote to the next stage
func (m *PromotionManager) Promote(toStage PromotionStage, reason string, approvedBy string) *PromotionResult {
	result := &PromotionResult{
		FromStage: m.currentStage,
		ToStage:   toStage,
	}

	// Validate promotion path
	if !m.isValidPromotionPath(m.currentStage, toStage) {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf(
			"Invalid promotion path from %s to %s", m.currentStage, toStage,
		))
		return result
	}

	// Create snapshot before promotion
	snapshot := m.createSnapshot(toStage)
	result.Snapshot = snapshot

	// Create promotion record
	promotion := &Promotion{
		ID:         generatePromotionID(),
		FromStage:  m.currentStage,
		ToStage:    toStage,
		Timestamp:  time.Now().Format(time.RFC3339),
		Reason:     reason,
		Status:     "pending",
		ApprovedBy: approvedBy,
	}

	// Validate promotion
	validationErrors := m.validatePromotion(toStage)
	if len(validationErrors) > 0 {
		result.Success = false
		result.Errors = validationErrors
		promotion.Status = "failed"
		m.history = append(m.history, promotion)
		return result
	}

	// Execute promotion
	if err := m.executePromotion(toStage); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, err.Error())
		promotion.Status = "failed"
		m.history = append(m.history, promotion)
		return result
	}

	// Mark success
	promotion.Status = "success"
	promotion.Evidence = []string{
		fmt.Sprintf("Promoted from %s to %s", m.currentStage, toStage),
		fmt.Sprintf("Snapshot: %s", snapshot.ID),
	}

	m.history = append(m.history, promotion)
	m.currentStage = toStage
	result.Success = true
	result.Promotion = promotion

	// Save state
	m.saveHistory()
	m.saveSnapshot(snapshot)

	return result
}

// Rollback rolls back to a previous stage
func (m *PromotionManager) Rollback(toStage PromotionStage, reason string) *PromotionResult {
	result := &PromotionResult{
		FromStage: m.currentStage,
		ToStage:   toStage,
	}

	// Find snapshot for the target stage
	var snapshot *Snapshot
	for _, s := range m.snapshots {
		if s.Stage == toStage {
			snapshot = s
			break
		}
	}

	if snapshot == nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf(
			"No snapshot found for stage %s", toStage,
		))
		return result
	}

	// Create rollback promotion
	promotion := &Promotion{
		ID:         generatePromotionID(),
		FromStage:  m.currentStage,
		ToStage:    StageRolledBack,
		Timestamp:  time.Now().Format(time.RFC3339),
		Reason:    reason,
		Status:    "pending",
		RollbackOf: snapshot.ID,
	}

	// Execute rollback
	if err := m.executeRollback(snapshot); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, err.Error())
		promotion.Status = "failed"
		m.history = append(m.history, promotion)
		return result
	}

	promotion.Status = "success"
	promotion.Evidence = []string{
		fmt.Sprintf("Rolled back from %s to %s", m.currentStage, toStage),
		fmt.Sprintf("Restored snapshot: %s", snapshot.ID),
	}

	m.history = append(m.history, promotion)
	m.currentStage = toStage
	result.Success = true
	result.Promotion = promotion
	result.Snapshot = snapshot

	m.saveHistory()

	return result
}

// isValidPromotionPath validates the promotion path
func (m *PromotionManager) isValidPromotionPath(from, to PromotionStage) bool {
	// Define valid promotion paths
	validPaths := map[PromotionStage][]PromotionStage{
		StageStable:     {StageStaging, StageCanary},
		StageStaging:    {StageCanary, StageStable},
		StageCanary:     {StageStable, StageStaging},
	}

	allowed, ok := validPaths[from]
	if !ok {
		return false
	}

	for _, stage := range allowed {
		if stage == to {
			return true
		}
	}

	return false
}

// validatePromotion validates that a promotion can proceed
func (m *PromotionManager) validatePromotion(toStage PromotionStage) []string {
	var errors []string

	// Check health status if we have an assessment
	// In practice, this would run a quick health check

	// Validate based on target stage
	switch toStage {
	case StageStable:
		// Stable requires canary or staging validation
		if m.currentStage != StageCanary && m.currentStage != StageStaging {
			errors = append(errors, "Must pass through staging or canary before stable")
		}

	case StageCanary:
		// Canary requires staging
		if m.currentStage == StageStable {
			errors = append(errors, "Cannot go directly to canary from stable")
		}
	}

	return errors
}

// executePromotion executes the actual promotion
func (m *PromotionManager) executePromotion(toStage PromotionStage) error {
	// In practice, this would:
	// 1. Update configuration files
	// 2. Update the runtime state
	// 3. Notify any observers

	return nil
}

// executeRollback restores a previous snapshot
func (m *PromotionManager) executeRollback(snapshot *Snapshot) error {
	// Restore state from snapshot
	if snapshot.State != nil {
		m.currentStage = snapshot.Stage
	}

	return nil
}

// createSnapshot creates a snapshot of the current state
func (m *PromotionManager) createSnapshot(stage PromotionStage) *Snapshot {
	snapshot := &Snapshot{
		ID:          generateSnapshotID(),
		Stage:       stage,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: fmt.Sprintf("Snapshot for promotion to %s", stage),
	}

	m.snapshots[snapshot.ID] = snapshot
	return snapshot
}

// GetHistory returns the promotion history
func (m *PromotionManager) GetHistory() []*Promotion {
	return m.history
}

// GetCurrentStage returns the current promotion stage
func (m *PromotionManager) GetCurrentStage() PromotionStage {
	return m.currentStage
}

// GetSnapshots returns all snapshots
func (m *PromotionManager) GetSnapshots() []*Snapshot {
	snapshots := make([]*Snapshot, 0, len(m.snapshots))
	for _, s := range m.snapshots {
		snapshots = append(snapshots, s)
	}
	return snapshots
}

// loadHistory loads promotion history from disk
func (m *PromotionManager) loadHistory() error {
	path := filepath.Join(m.kdsePath, "runtime", "promotion", "history.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var history []*Promotion
	if err := json.Unmarshal(data, &history); err != nil {
		return err
	}

	m.history = history

	// Set current stage based on last promotion
	if len(history) > 0 {
		last := history[len(history)-1]
		if last.Status == "success" {
			m.currentStage = last.ToStage
		}
	}

	return nil
}

// saveHistory saves promotion history to disk
func (m *PromotionManager) saveHistory() error {
	dir := filepath.Join(m.kdsePath, "runtime", "promotion")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "history.json")
	data, err := json.MarshalIndent(m.history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// saveSnapshot saves a snapshot to disk
func (m *PromotionManager) saveSnapshot(snapshot *Snapshot) error {
	dir := filepath.Join(m.kdsePath, "runtime", "promotion", "snapshots")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, snapshot.ID+".json")
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// generatePromotionID generates a unique promotion ID
func generatePromotionID() string {
	return fmt.Sprintf("KDSE-PROMO-%s", time.Now().Format("20060102-150405"))
}

// generateSnapshotID generates a unique snapshot ID
func generateSnapshotID() string {
	return fmt.Sprintf("KDSE-SNAP-%s", time.Now().Format("20060102-150405"))
}

// FormatPromotionHistory formats the promotion history for display
func FormatPromotionHistory(history []*Promotion) string {
	var output string

	output += "═══════════════════════════════════════════════════════════════\n"
	output += "                    Promotion History                          \n"
	output += "═══════════════════════════════════════════════════════════════\n"

	if len(history) == 0 {
		output += "  No promotions recorded.\n"
	} else {
		for _, p := range history {
			statusIcon := map[string]string{
				"pending":    "○",
				"success":    "●",
				"failed":     "✗",
				"rolled_back": "↩",
			}[p.Status]

			output += fmt.Sprintf("  %s %s → %s\n", statusIcon, p.FromStage, p.ToStage)
			output += fmt.Sprintf("     %s\n", p.Timestamp)
			if p.Reason != "" {
				output += fmt.Sprintf("     Reason: %s\n", p.Reason)
			}
		}
	}

	output += "═══════════════════════════════════════════════════════════════\n"

	return output
}

// boolToStatus converts a boolean to status string
func boolToStatus(b bool) string {
	if b {
		return "SUCCESS"
	}
	return "FAILED"
}

// FormatPromotionResult formats a promotion result for display
func FormatPromotionResult(result *PromotionResult) string {
	var output string

	output += "═══════════════════════════════════════════════════════════════\n"
	output += "                    Promotion Result                          \n"
	output += "═══════════════════════════════════════════════════════════════\n"
	output += fmt.Sprintf("  Status: %s\n", boolToStatus(result.Success))
	output += fmt.Sprintf("  From: %s → To: %s\n", result.FromStage, result.ToStage)

	if result.Success && result.Promotion != nil {
		output += fmt.Sprintf("  Promotion ID: %s\n", result.Promotion.ID)
		output += fmt.Sprintf("  Timestamp: %s\n", result.Promotion.Timestamp)
	}

	if len(result.Errors) > 0 {
		output += "  Errors:\n"
		for _, err := range result.Errors {
			output += fmt.Sprintf("    ✗ %s\n", err)
		}
	}

	if len(result.Warnings) > 0 {
		output += "  Warnings:\n"
		for _, warn := range result.Warnings {
			output += fmt.Sprintf("    ⚠ %s\n", warn)
		}
	}

	output += "═══════════════════════════════════════════════════════════════\n"

	return output
}
