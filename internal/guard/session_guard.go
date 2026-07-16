// Package guard provides the KDSE Session Guard - a critical enforcement layer
// that ensures all KDSE operations require a valid initialized workspace and session state.
// No operation should bypass this guard.
package guard

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kdse/runtime/internal/workspace"
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

// Common guard errors
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
	ErrSessionInvalid = &GuardError{
		Code:    "KDSE_GUARD_003",
		Message: "Session state is invalid or corrupted",
		Hint:    "Please run `kdse initialize` to reset the session",
	}
	ErrSessionExpired = &GuardError{
		Code:    "KDSE_GUARD_004",
		Message: "Session has expired",
		Hint:    "Please run `kdse initialize` to start a new session",
	}
)

// SessionGuard is a thread-safe enforcement layer that validates workspace and session state
// before allowing any KDSE operations to proceed.
type SessionGuard struct {
	repoPath  string
	ws        *workspace.Workspace
	mu        sync.RWMutex
	autoInit  bool
	logging   bool
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

// NewSessionGuard creates a new SessionGuard for the given repository path
func NewSessionGuard(repoPath string) *SessionGuard {
	return &SessionGuard{
		repoPath: repoPath,
		ws:       workspace.New(repoPath),
		autoInit: false,
		logging:  true,
	}
}

// NewSessionGuardWithAutoInit creates a new SessionGuard with auto-initialization enabled
func NewSessionGuardWithAutoInit(repoPath string) *SessionGuard {
	return &SessionGuard{
		repoPath: repoPath,
		ws:       workspace.New(repoPath),
		autoInit: true,
		logging:  true,
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
// It checks if the workspace and session are properly initialized and active.
// If autoInit is enabled, it will attempt to initialize on first run.
// Returns nil if initialized, otherwise returns a GuardError.
func (g *SessionGuard) EnforceInitialized() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.logging {
		log.Printf("[GUARD] Checking initialization status...")
	}

	// Step 1: Check if .kdse/ workspace directory exists
	if !g.ws.Exists() {
		if g.logging {
			log.Printf("[GUARD] Workspace not found: %s", g.ws.Root())
		}
		
		if g.autoInit {
			return g.autoInitialize()
		}
		
		return ErrWorkspaceNotInitialized
	}

	// Step 2: Check if session state file exists and is valid
	sessionState, err := g.loadSessionState()
	if err != nil {
		if g.logging {
			log.Printf("[GUARD] Session state error: %v", err)
		}
		
		if g.autoInit && os.IsNotExist(err) {
			return g.autoInitialize()
		}
		
		return ErrSessionNotActive
	}

	// Step 3: Verify session state is valid
	if err := g.validateSessionState(sessionState); err != nil {
		if g.logging {
			log.Printf("[GUARD] Session validation failed: %v", err)
		}
		
		if g.autoInit {
			return g.autoInitialize()
		}
		
		return ErrSessionInvalid
	}

	// Step 4: Verify session hasn't expired (optional 24-hour expiry)
	if g.isSessionExpired(sessionState) {
		if g.logging {
			log.Printf("[GUARD] Session expired: %s", sessionState.SessionID)
		}
		
		if g.autoInit {
			return g.autoInitialize()
		}
		
		return ErrSessionExpired
	}

	if g.logging {
		log.Printf("[GUARD] Initialization check passed. Session: %s", sessionState.SessionID)
	}

	return nil
}

// Check performs a non-enforcing check of initialization status
// Returns GuardResult without returning an error
func (g *SessionGuard) Check() *GuardResult {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := &GuardResult{
		Valid:         false,
		SessionActive: false,
		WorkspaceReady: g.ws.Exists(),
	}

	if !g.ws.Exists() {
		result.Error = ErrWorkspaceNotInitialized
		return result
	}

	sessionState, err := g.loadSessionState()
	if err != nil {
		result.Error = ErrSessionNotActive
		return result
	}

	if err := g.validateSessionState(sessionState); err != nil {
		result.Error = ErrSessionInvalid
		return result
	}

	if g.isSessionExpired(sessionState) {
		result.Error = ErrSessionExpired
		return result
	}

	result.Valid = true
	result.SessionActive = true
	result.SessionID = sessionState.SessionID
	result.SessionAge = g.formatSessionAge(sessionState.StartedAt)

	return result
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
func (g *SessionGuard) Initialize() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.logging {
		log.Printf("[GUARD] Initializing KDSE workspace...")
	}

	// Step 1: Create workspace directory
	if err := g.ws.Initialize(); err != nil {
		if g.logging {
			log.Printf("[GUARD] Workspace initialization failed: %v", err)
		}
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Step 2: Create session state
	sessionState := &SessionState{
		SessionID:     g.generateSessionID(),
		StartedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
		WorkspaceRoot: g.ws.Root(),
		Version:       "1.0.0",
		Initialized:   true,
	}

	if err := g.saveSessionState(sessionState); err != nil {
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
	return filepath.Join(g.repoPath, ".kdse", "session-state.json")
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
func (g *SessionGuard) autoInitialize() error {
	if g.logging {
		log.Printf("[GUARD] Auto-initializing workspace...")
	}

	// Create workspace directory
	if err := g.ws.Initialize(); err != nil {
		return fmt.Errorf("auto-initialization failed: %w", err)
	}

	// Create session state
	sessionState := &SessionState{
		SessionID:     g.generateSessionID(),
		StartedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
		WorkspaceRoot: g.ws.Root(),
		Version:       "1.0.0",
		Initialized:   true,
	}

	if err := g.saveSessionState(sessionState); err != nil {
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
func (g *SessionGuard) IsInitialized() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.ws.Exists() {
		return false
	}

	_, err := g.loadSessionState()
	return err == nil
}

// Reset clears all session state and workspace
// This should only be used for testing or recovery scenarios
func (g *SessionGuard) Reset() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.logging {
		log.Printf("[GUARD] Resetting session state...")
	}

	// Remove session state file
	sessionPath := g.sessionStatePath()
	if _, err := os.Stat(sessionPath); err == nil {
		if err := os.Remove(sessionPath); err != nil {
			return fmt.Errorf("failed to remove session state: %w", err)
		}
	}

	// Note: We don't remove the .kdse/ directory as it may contain other data
	// The next Initialize() will overwrite the session state

	if g.logging {
		log.Printf("[GUARD] Session state reset complete")
	}

	return nil
}

// GetSessionID returns the current session ID if a session is active
func (g *SessionGuard) GetSessionID() (string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	state, err := g.loadSessionState()
	if err != nil {
		return "", err
	}

	return state.SessionID, nil
}

// UpdateSessionTimestamp updates the session's last activity timestamp
func (g *SessionGuard) UpdateSessionTimestamp() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	state, err := g.loadSessionState()
	if err != nil {
		return err
	}

	state.UpdatedAt = time.Now().Format(time.RFC3339)
	return g.saveSessionState(state)
}
