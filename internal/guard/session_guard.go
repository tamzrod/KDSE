// Package guard provides the KDSE Session Guard.
//
// DEPRECATED: The original SessionGuard has been refactored into separate guards:
//   - RuntimeGuard: Orchestrates all guards (internal/guard/runtime_guard.go)
//   - ProjectGuard: Project discovery and validation (internal/guard/project_guard.go)
//   - WorkspaceGuard: Workspace validation (internal/guard/workspace_guard.go)
//   - SessionValidationGuard: Session validation (internal/guard/session_validation_guard.go)
//   - LifecycleGuard: Lifecycle validation (internal/guard/lifecycle_guard.go)
//
// This file is kept for backward compatibility but delegates to the new architecture.
//
// The SessionGuard is now only responsible for:
//   - Active session detection
//   - Session validity
//   - Session expiration
//   - Session recovery
//
// Session Guard must NEVER:
//   - Discover projects (now ProjectGuard)
//   - Validate projects (now ProjectGuard)
//   - Initialize projects (now RuntimeGuard)
//   - Create workspaces (now RuntimeGuard)
package guard

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// GuardError represents a session guard enforcement error
type GuardError struct {
	Code    string
	Message string
	Hint    string
}

func (e *GuardError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common guard errors (legacy types for backward compatibility)
var (
	ErrWorkspaceNotInitialized = &GuardError{
		Code:    "KDSE_GUARD_001",
		Message: "KDSE workspace not initialized",
		Hint:    "Please run `kdse initialize` first",
	}
	ErrSessionNotActive = &GuardError{
		Code:    "KDSE_GUARD_002",
		Message: "No active KDSE session",
		Hint:    "Please run `kdse initialize` to start a session",
	}
	ErrSessionInvalidLegacy = &GuardError{
		Code:    "KDSE_GUARD_003",
		Message: "Session state is invalid or corrupted",
		Hint:    "Please run `kdse initialize` to reset the session",
	}
	ErrSessionExpiredLegacy = &GuardError{
		Code:    "KDSE_GUARD_004",
		Message: "Session has expired",
		Hint:    "Please run `kdse initialize` to start a new session",
	}
)

// SessionGuard is a thread-safe enforcement layer that validates session state.
// It has been refactored to ONLY manage sessions - all project and workspace
// logic has been moved to dedicated guards.
//
// NOTE: For new code, use RuntimeGuard which orchestrates all guards together.
// This SessionGuard is kept for backward compatibility.
type SessionGuard struct {
	repoPath      string
	workspacePath string
	mu            sync.RWMutex
	autoInit      bool
	logging       bool
}

// SessionState represents the persisted session state used by the guard
type SessionState struct {
	SessionID      string `json:"session_id"`
	SessionToken   string `json:"session_token,omitempty"`
	StartedAt      string `json:"started_at"`
	UpdatedAt      string `json:"updated_at"`
	WorkspaceRoot  string `json:"workspace_root"`
	Version        string `json:"version"`
	Initialized    bool   `json:"initialized"`
}

// GuardResult contains the result of a guard check
type GuardResult struct {
	Valid          bool        `json:"valid"`
	SessionActive  bool        `json:"session_active"`
	WorkspaceReady bool        `json:"workspace_ready"`
	SessionID      string      `json:"session_id,omitempty"`
	SessionAge     string      `json:"session_age,omitempty"`
	Error         *GuardError `json:"error,omitempty"`
}

// NewSessionGuard creates a new SessionGuard for the given repository path.
// NOTE: For new code, prefer RuntimeGuard which orchestrates all guards.
func NewSessionGuard(repoPath string) *SessionGuard {
	return &SessionGuard{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
		autoInit:      false,
		logging:       true,
	}
}

// NewSessionGuardWithAutoInit creates a new SessionGuard with auto-initialization enabled
// NOTE: For new code, prefer RuntimeGuard which orchestrates all guards.
func NewSessionGuardWithAutoInit(repoPath string) *SessionGuard {
	return &SessionGuard{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
		autoInit:      true,
		logging:       true,
	}
}

// SetAutoInit enables or disables auto-initialization
func (g *SessionGuard) SetAutoInit(enabled bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.autoInit = enabled
}

// SetLogging enables or disables guard logging
func (g *SessionGuard) SetLogging(enabled bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.logging = enabled
}

// EnforceInitialized is the primary entry point for session enforcement.
//
// DEPRECATED: This method has been superseded by RuntimeGuard.Validate().
// For new code, use RuntimeGuard which orchestrates all guards together.
//
// This method checks if the workspace and session are properly initialized and active.
// If autoInit is enabled, it will attempt to initialize on first run.
// Returns nil if initialized, otherwise returns a GuardError.
func (g *SessionGuard) EnforceInitialized() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.logging {
		log.Printf("[GUARD] Checking initialization status...")
	}

	// Delegate to RuntimeGuard for full validation
	runtimeGuard := NewRuntimeGuard(g.repoPath)
	result := runtimeGuard.Validate(nil)

	if result.Valid {
		if g.logging {
			log.Printf("[GUARD] Initialization check passed.")
		}
		return nil
	}

	// Map errors from RuntimeGuard to legacy GuardError
	for _, err := range result.Errors {
		switch err.GuardType {
		case GuardTypeProject:
			return &GuardError{
				Code:    "KDSE_GUARD_000",
				Message: err.Message,
				Hint:    err.Hint,
			}
		case GuardTypeWorkspace:
			return ErrWorkspaceNotInitialized
		case GuardTypeSession:
			if err.Code == "SESSION_EXPIRED" {
				return ErrSessionExpiredLegacy
			}
			return ErrSessionNotActive
		}
	}

	return ErrSessionNotActive
}

// Check performs a non-enforcing check of initialization status
// Returns GuardResult without returning an error
//
// DEPRECATED: For new code, use RuntimeGuard.QuickCheck().
func (g *SessionGuard) Check() *GuardResult {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Delegate to RuntimeGuard for full quick check
	runtimeGuard := NewRuntimeGuard(g.repoPath)
	result := runtimeGuard.QuickCheck(nil)

	return &GuardResult{
		Valid:         result.Valid,
		SessionActive: result.Session != nil && result.Session.Valid,
		WorkspaceReady: result.Workspace != nil && result.Workspace.Valid,
		SessionID:     "",
		SessionAge:    "",
		Error:        nil,
	}
}

// EnforceForOperation is a convenience method that enforces initialization
// and includes the operation name in error messages
func (g *SessionGuard) EnforceForOperation(operationName string) error {
	err := g.EnforceInitialized()
	if err != nil {
		if ge, ok := err.(*GuardError); ok {
			return fmt.Errorf("[%s] Cannot %s: %s. %s", ge.Code, operationName, ge.Message, ge.Hint)
		}
		return fmt.Errorf("[GUARD] Cannot %s: %v", operationName, err)
	}
	return nil
}

// Initialize performs a full initialization of the KDSE workspace and session
//
// DEPRECATED: For new code, use RuntimeGuard with initialization.
func (g *SessionGuard) Initialize() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.logging {
		log.Printf("[GUARD] Initializing KDSE workspace...")
	}

	// Create workspace directory
	if err := os.MkdirAll(g.workspacePath, 0755); err != nil {
		if g.logging {
			log.Printf("[GUARD] Workspace initialization failed: %v", err)
		}
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Create session state using the new SessionValidationGuard
	sessionGuard := NewSessionValidationGuard(g.repoPath)
	sessionState, err := sessionGuard.CreateSession()
	if err != nil {
		if g.logging {
			log.Printf("[GUARD] Session state save failed: %v", err)
		}
		return fmt.Errorf("failed to save session state: %w", err)
	}

	if g.logging {
		log.Printf("[GUARD] Initialization complete. Session ID: %s", sessionState.SessionID)
	}

	return nil
}

// sessionStatePath returns the path to the session state file
func (g *SessionGuard) sessionStatePath() string {
	return filepath.Join(g.workspacePath, "session.yaml")
}

// loadSessionState loads the session state from disk
func (g *SessionGuard) loadSessionState() (*SessionState, error) {
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
func (g *SessionGuard) saveSessionState(state *SessionState) error {
	path := g.sessionStatePath()
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// validateSessionState checks if the session state is valid
func (g *SessionGuard) validateSessionState(state *SessionState) error {
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
func (g *SessionGuard) isSessionExpired(state *SessionState) bool {
	const sessionExpiryDuration = 24 * time.Hour

	startedAt, err := time.Parse(time.RFC3339, state.StartedAt)
	if err != nil {
		return true // Invalid start time means expired
	}

	expiryTime := startedAt.Add(sessionExpiryDuration)
	return time.Now().After(expiryTime)
}

// formatSessionAge returns a human-readable session age string
func (g *SessionGuard) formatSessionAge(startedAt string) string {
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

// autoInitialize performs automatic initialization when autoInit is enabled
//
// DEPRECATED: This method has been refactored.
func (g *SessionGuard) autoInitialize() error {
	if g.logging {
		log.Printf("[GUARD] Auto-initializing workspace...")
	}

	// Create workspace directory
	if err := os.MkdirAll(g.workspacePath, 0755); err != nil {
		return fmt.Errorf("auto-initialization failed: %w", err)
	}

	// Create session state using SessionValidationGuard
	sessionGuard := NewSessionValidationGuard(g.repoPath)
	sessionState, err := sessionGuard.CreateSession()
	if err != nil {
		return fmt.Errorf("auto-initialization failed: %w", err)
	}

	if g.logging {
		log.Printf("[GUARD] Auto-initialization complete. Session ID: %s", sessionState.SessionID)
	}

	return nil
}

// generateSessionID creates a unique session identifier
func (g *SessionGuard) generateSessionID() string {
	return fmt.Sprintf("KDSE-GUARD-%s", time.Now().Format("20060102-150405"))
}

// IsInitialized is a quick check if the workspace is initialized
// without enforcing initialization
//
// DEPRECATED: For new code, use RuntimeGuard.GetCurrentState().
func (g *SessionGuard) IsInitialized() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Quick check - delegate to RuntimeGuard
	runtimeGuard := NewRuntimeGuard(g.repoPath)
	state := runtimeGuard.GetCurrentState()
	return state == StateLifecycleReady
}

// Reset clears all session state and workspace
// This should only be used for testing or recovery scenarios
//
// DEPRECATED: For new code, use SessionValidationGuard.EndSession().
func (g *SessionGuard) Reset() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.logging {
		log.Printf("[GUARD] Resetting session state...")
	}

	// Remove session state file using SessionValidationGuard
	sessionGuard := NewSessionValidationGuard(g.repoPath)
	if err := sessionGuard.EndSession(); err != nil {
		return fmt.Errorf("failed to remove session state: %w", err)
	}

	if g.logging {
		log.Printf("[GUARD] Session state reset complete")
	}

	return nil
}

// GetSessionID returns the current session ID if a session is active
//
// DEPRECATED: For new code, use SessionValidationGuard.GetSessionID().
func (g *SessionGuard) GetSessionID() (string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	sessionGuard := NewSessionValidationGuard(g.repoPath)
	return sessionGuard.GetSessionID()
}

// UpdateSessionTimestamp updates the session's last activity timestamp
//
// DEPRECATED: For new code, use SessionValidationGuard.UpdateSessionTimestamp().
func (g *SessionGuard) UpdateSessionTimestamp() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	sessionGuard := NewSessionValidationGuard(g.repoPath)
	return sessionGuard.UpdateSessionTimestamp()
}
