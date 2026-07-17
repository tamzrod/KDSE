// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionValidationGuard is responsible ONLY for session management.
// It validates:
//   - Active session detection
//   - Session validity
//   - Session expiration
//   - Session recovery
//
// Session Guard must NEVER:
//   - Discover projects
//   - Validate projects
//   - Initialize projects
//   - Create workspaces
//
// Session Guard assumes Project Guard and Workspace Guard already succeeded.
type SessionValidationGuard struct {
	repoPath      string
	workspacePath string
	mu            sync.RWMutex
}

// SessionState represents the persisted session state
type SessionState struct {
	SessionID     string `json:"session_id"`
	SessionToken  string `json:"session_token,omitempty"`
	StartedAt     string `json:"started_at"`
	UpdatedAt     string `json:"updated_at"`
	WorkspaceRoot string `json:"workspace_root"`
	Version       string `json:"version"`
	Initialized   bool   `json:"initialized"`
}

// NewSessionValidationGuard creates a new Session Validation Guard
func NewSessionValidationGuard(repoPath string) *SessionValidationGuard {
	return &SessionValidationGuard{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
	}
}

// Validate validates session state.
// Returns SessionGuardResult indicating if session is valid.
func (g *SessionValidationGuard) Validate(ctx context.Context) *SessionGuardResult {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := &SessionGuardResult{
		GuardResult: GuardResult{
			Valid:        false,
			State:        StateWorkspace,
			StateBefore:  StateWorkspace,
			StateAfter:   StateWorkspace,
			GuardType:    GuardTypeSession,
		},
	}

	// Step 1: Load session state
	state, err := g.loadSessionState()
	if err != nil {
		if os.IsNotExist(err) {
			result.Error = ErrNoActiveSession
		} else {
			result.Error = ErrSessionInvalid
		}
		return result
	}

	// Step 2: Validate session state structure
	if err := g.validateSessionState(state); err != nil {
		result.Error = ErrSessionInvalid
		return result
	}

	// Step 3: Check session expiration (24-hour default)
	if g.isSessionExpired(state) {
		result.Error = ErrSessionExpired
		return result
	}

	// Session is valid
	result.Valid = true
	result.State = StateSession
	result.StateBefore = StateWorkspace
	result.StateAfter = StateSession
	result.SessionID = state.SessionID
	result.StartedAt = state.StartedAt
	result.WorkspaceRoot = state.WorkspaceRoot

	return result
}

// HasActiveSession checks if there is an active session without full validation
func (g *SessionValidationGuard) HasActiveSession() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	state, err := g.loadSessionState()
	if err != nil {
		return false
	}

	if err := g.validateSessionState(state); err != nil {
		return false
	}

	return !g.isSessionExpired(state)
}

// GetSessionID returns the current session ID if active
func (g *SessionValidationGuard) GetSessionID() (string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	state, err := g.loadSessionState()
	if err != nil {
		return "", err
	}

	return state.SessionID, nil
}

// GetSessionState returns the current session state
func (g *SessionValidationGuard) GetSessionState() (*SessionState, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.loadSessionState()
}

// UpdateSessionTimestamp updates the session's last activity timestamp
func (g *SessionValidationGuard) UpdateSessionTimestamp() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	state, err := g.loadSessionState()
	if err != nil {
		return err
	}

	state.UpdatedAt = time.Now().Format(time.RFC3339)
	return g.saveSessionState(state)
}

// sessionStatePath returns the path to the session state file
func (g *SessionValidationGuard) sessionStatePath() string {
	return filepath.Join(g.workspacePath, "session.yaml")
}

// loadSessionState loads the session state from disk
func (g *SessionValidationGuard) loadSessionState() (*SessionState, error) {
	path := g.sessionStatePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// saveSessionState persists the session state to disk
func (g *SessionValidationGuard) saveSessionState(state *SessionState) error {
	path := g.sessionStatePath()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// validateSessionState checks if the session state is valid
func (g *SessionValidationGuard) validateSessionState(state *SessionState) error {
	if state == nil {
		return fmt.Errorf("session state is nil")
	}

	if state.SessionID == "" {
		return fmt.Errorf("session ID is empty")
	}

	if state.StartedAt == "" {
		return fmt.Errorf("session start time is empty")
	}

	if _, err := time.Parse(time.RFC3339, state.StartedAt); err != nil {
		return fmt.Errorf("invalid session start time format: %w", err)
	}

	if state.WorkspaceRoot == "" {
		return fmt.Errorf("workspace root is empty")
	}

	return nil
}

// isSessionExpired checks if the session has expired (24-hour default expiry)
func (g *SessionValidationGuard) isSessionExpired(state *SessionState) bool {
	const sessionExpiryDuration = 24 * time.Hour

	startedAt, err := time.Parse(time.RFC3339, state.StartedAt)
	if err != nil {
		return true // Invalid start time means expired
	}

	expiryTime := startedAt.Add(sessionExpiryDuration)
	return time.Now().After(expiryTime)
}

// FormatSessionAge returns a human-readable session age string
func (g *SessionValidationGuard) FormatSessionAge(startedAt string) string {
	started, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		return "unknown"
	}

	age := time.Since(started)
	if age < time.Minute {
		return "just now"
	}
	if age < time.Hour {
		minutes := int(age.Minutes())
		return fmt.Sprintf("%d minute(s) ago", minutes)
	}
	hours := int(age.Hours())
	return fmt.Sprintf("%d hour(s) ago", hours)
}

// CreateSession creates a new session
func (g *SessionValidationGuard) CreateSession() (*SessionState, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	sessionState := &SessionState{
		SessionID:     g.generateSessionID(),
		StartedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
		WorkspaceRoot: g.workspacePath,
		Version:       "1.0.0",
		Initialized:   true,
	}

	if err := g.saveSessionState(sessionState); err != nil {
		return nil, err
	}

	return sessionState, nil
}

// EndSession clears the current session
func (g *SessionValidationGuard) EndSession() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	sessionPath := g.sessionStatePath()
	if _, err := os.Stat(sessionPath); err == nil {
		return os.Remove(sessionPath)
	}

	return nil
}

// generateSessionID creates a unique session identifier
func (g *SessionValidationGuard) generateSessionID() string {
	return fmt.Sprintf("KDSE-GUARD-%s", time.Now().Format("20060102-150405"))
}
