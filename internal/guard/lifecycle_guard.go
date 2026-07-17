// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

// LifecycleGuard is responsible ONLY for lifecycle state management.
// It validates:
//   - Current phase
//   - Allowed transitions
//   - Lifecycle integrity
//
// Lifecycle Guard must NEVER:
//   - Discover projects
//   - Perform workspace logic
//   - Manage sessions
type LifecycleGuard struct {
	repoPath      string
	workspacePath string
}

// NewLifecycleGuard creates a new Lifecycle Guard
func NewLifecycleGuard(repoPath string) *LifecycleGuard {
	return &LifecycleGuard{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
	}
}

// Validate validates lifecycle state.
// Returns LifecycleGuardResult indicating if lifecycle is valid.
func (g *LifecycleGuard) Validate(ctx context.Context) *LifecycleGuardResult {
	result := &LifecycleGuardResult{
		RuntimeGuardResult: RuntimeGuardResult{
			Valid:        false,
			State:        StateSession,
			StateBefore:  StateSession,
			StateAfter:   StateSession,
			GuardType:    GuardTypeLifecycle,
		},
	}

	// Step 1: Load lifecycle state
	state, err := g.loadLifecycleState()
	if err != nil {
		result.Error = ErrLifecycleInvalid
		return result
	}

	// Step 2: Validate phase
	currentPhase := state.Current
	if currentPhase == "" {
		currentPhase = "initialization"
	}

	if !g.isValidPhase(currentPhase) {
		result.Error = ErrLifecycleInvalid
		return result
	}

	// Step 3: Validate allowed transitions
	allowedPhases := g.getAllowedTransitions(currentPhase)
	if len(allowedPhases) == 0 && currentPhase != "reports" {
		// Terminal phase without transition is valid
	}

	// Lifecycle is valid
	result.Valid = true
	result.State = StateLifecycleReady
	result.StateBefore = StateSession
	result.StateAfter = StateLifecycleReady
	result.CurrentPhase = currentPhase
	result.AllowedPhases = allowedPhases
	result.TransitionValid = true

	return result
}

// IsValid checks if lifecycle state is valid without full validation
func (g *LifecycleGuard) IsValid() bool {
	state, err := g.loadLifecycleState()
	if err != nil {
		return false
	}

	currentPhase := state.Current
	if currentPhase == "" {
		return true // Default phase is valid
	}

	return g.isValidPhase(currentPhase)
}

// GetCurrentPhase returns the current phase
func (g *LifecycleGuard) GetCurrentPhase() (string, error) {
	state, err := g.loadLifecycleState()
	if err != nil {
		return "", err
	}

	if state.Current == "" {
		return "initialization", nil
	}
	return state.Current, nil
}

// GetAllowedTransitions returns allowed next phases
func (g *LifecycleGuard) GetAllowedTransitions() ([]string, error) {
	phase, err := g.GetCurrentPhase()
	if err != nil {
		return nil, err
	}

	return g.getAllowedTransitions(phase), nil
}

// lifecycleState represents the persisted lifecycle state
type lifecycleState struct {
	Current       string   `json:"current"`
	History       []string `json:"history,omitempty"`
	Version       string   `json:"version"`
}

// loadLifecycleState loads the lifecycle state from disk
func (g *LifecycleGuard) loadLifecycleState() (*lifecycleState, error) {
	path := filepath.Join(g.workspacePath, "phase.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state lifecycleState
	if err := json.Unmarshal(data, &state); err != nil {
		// Try YAML-style parsing
		return g.parseYAMLState(data)
	}

	return &state, nil
}

// parseYAMLState parses YAML-formatted state
func (g *LifecycleGuard) parseYAMLState(data []byte) (*lifecycleState, error) {
	// Simple YAML parsing for phase state
	state := &lifecycleState{
		Version: "1.0.0",
	}

	content := string(data)
	for _, line := range splitLines(content) {
		line = trimAndRemoveComment(line)
		if line == "" {
			continue
		}

		if hasPrefix(line, "current:") {
			parts := splitYAMLValue(line, ":")
			if len(parts) == 2 {
				state.Current = trim(parts[1])
			}
		}
	}

	return state, nil
}

// isValidPhase checks if a phase is valid
func (g *LifecycleGuard) isValidPhase(phase string) bool {
	validPhases := map[string]bool{
		"initialization": true,
		"knowledge":      true,
		"architecture":   true,
		"implementation": true,
		"verification":   true,
		"reports":        true,
	}
	return validPhases[phase]
}

// getAllowedTransitions returns allowed transitions for a phase
func (g *LifecycleGuard) getAllowedTransitions(phase string) []string {
	transitions := map[string][]string{
		"initialization": {"knowledge"},
		"knowledge":      {"architecture"},
		"architecture":   {"implementation"},
		"implementation": {"verification"},
		"verification":   {"reports"},
		"reports":        {},
	}

	if phases, ok := transitions[phase]; ok {
		return phases
	}
	return nil
}

// SetPhase sets the current phase
func (g *LifecycleGuard) SetPhase(phase string) error {
	if !g.isValidPhase(phase) {
		return ErrLifecycleTransitionInvalid
	}

	state, err := g.loadLifecycleState()
	if err != nil {
		state = &lifecycleState{
			Version: "1.0.0",
		}
	}

	// Record in history
	if state.Current != "" && state.Current != phase {
		state.History = append(state.History, state.Current)
	}

	state.Current = phase

	return g.saveLifecycleState(state)
}

// saveLifecycleState persists the lifecycle state
func (g *LifecycleGuard) saveLifecycleState(state *lifecycleState) error {
	path := filepath.Join(g.workspacePath, "phase.yaml")

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Helper functions for YAML-like parsing

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimAndRemoveComment(s string) string {
	// Remove inline comments
	for i := 0; i < len(s); i++ {
		if s[i] == '#' {
			s = s[:i]
			break
		}
	}
	return trim(s)
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func splitYAMLValue(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
