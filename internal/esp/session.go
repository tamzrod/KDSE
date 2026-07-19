// Package esp implements the KDSE Engineering Session Protocol (ESP).
package esp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Session represents an active engineering session
type Session struct {
	// Identity
	ID          string       `json:"id"`
	Token       string       `json:"token"`
	State       SessionState `json:"state"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	WorkspaceID string       `json:"workspace_id"`

	// Context
	Context *Context `json:"context,omitempty"`

	// Configuration
	SessionTimeout  time.Duration `json:"session_timeout"`
	MaxIdleDuration time.Duration `json:"max_idle_duration"`

	// State machine lock
	mu sync.RWMutex
}

// NewSession creates a new session with the given workspace ID
func NewSession(workspaceID string) *Session {
	now := time.Now().Format(time.RFC3339)
	return &Session{
		ID:             uuid.New().String(),
		Token:          uuid.New().String(),
		State:          StateIdle,
		CreatedAt:      now,
		UpdatedAt:      now,
		WorkspaceID:    workspaceID,
		SessionTimeout: 24 * time.Hour,
		MaxIdleDuration: 1 * time.Hour,
	}
}

// GetState returns the current session state
func (s *Session) GetState() SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.State
}

// SetContext sets the engineering context
func (s *Session) SetContext(ctx *Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Context = ctx
}

// GetContext returns the current engineering context
func (s *Session) GetContext() *Context {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Context
}

// TransitionTo transitions the session to a new state
// Returns an error if the transition is invalid
func (s *Session) TransitionTo(newState SessionState, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldState := s.State

	// Check if transition is valid
	if !IsValidTransition(oldState, newState) {
		return fmt.Errorf("invalid state transition from %s to %s: %s", oldState, newState, reason)
	}

	s.State = newState
	s.UpdatedAt = time.Now().Format(time.RFC3339)

	return nil
}

// IsActive returns true if the session is in an active state
func (s *Session) IsActive() bool {
	return s.GetState() == StateActive
}

// IsExpired returns true if the session has expired
func (s *Session) IsExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	createdAt, err := time.Parse(time.RFC3339, s.CreatedAt)
	if err != nil {
		return true
	}

	return time.Since(createdAt) > s.SessionTimeout
}

// CanResume returns true if the session can be resumed
func (s *Session) CanResume() bool {
	state := s.GetState()
	return state == StateSuspended
}

// =============================================================================
// Session Manager
// =============================================================================

// SessionManager manages engineering sessions
type SessionManager struct {
	repoPath      string
	workspacePath string
	activeSession *Session
	mu            sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager(repoPath string) *SessionManager {
	return &SessionManager{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
		activeSession: nil,
	}
}

// CreateSession creates a new engineering session
func (sm *SessionManager) CreateSession(workspaceID string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.activeSession != nil && sm.activeSession.GetState() != StateTerminated && sm.activeSession.GetState() != StateFailed {
		return nil, fmt.Errorf("session already exists: %s", sm.activeSession.ID)
	}

	session := NewSession(workspaceID)
	session.State = StateDiscovering
	sm.activeSession = session

	// Save session state
	if err := sm.saveSession(session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

// GetSession returns the current active session
func (sm *SessionManager) GetSession() *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.activeSession
}

// TransitionSession transitions the active session to a new state
func (sm *SessionManager) TransitionSession(newState SessionState, reason string) error {
	sm.mu.RLock()
	session := sm.activeSession
	sm.mu.RUnlock()

	if session == nil {
		return fmt.Errorf("no active session")
	}

	if err := session.TransitionTo(newState, reason); err != nil {
		return err
	}

	return sm.saveSession(session)
}

// SetContext sets the engineering context for the active session
func (sm *SessionManager) SetContext(ctx *Context) error {
	sm.mu.RLock()
	session := sm.activeSession
	sm.mu.RUnlock()

	if session == nil {
		return fmt.Errorf("no active session")
	}

	session.SetContext(ctx)

	// Update session with checksum
	if err := ctx.Update(); err != nil {
		return fmt.Errorf("failed to update context checksum: %w", err)
	}

	return sm.saveSession(session)
}

// SuspendSession suspends the active session
func (sm *SessionManager) SuspendSession() error {
	return sm.TransitionSession(StateSuspended, "Session suspended by user")
}

// ResumeSession resumes a suspended session
func (sm *SessionManager) ResumeSession() error {
	sm.mu.RLock()
	session := sm.activeSession
	sm.mu.RUnlock()

	if session == nil {
		return fmt.Errorf("no session to resume")
	}

	if !session.CanResume() {
		return fmt.Errorf("session cannot be resumed from state: %s", session.GetState())
	}

	return sm.TransitionSession(StateActive, "Session resumed")
}

// TerminateSession terminates the active session
func (sm *SessionManager) TerminateSession() error {
	return sm.TransitionSession(StateTerminated, "Session terminated")
}

// FailSession marks the session as failed
func (sm *SessionManager) FailSession(code FailureCode, message string) error {
	sm.mu.RLock()
	session := sm.activeSession
	sm.mu.RUnlock()

	if session == nil {
		return fmt.Errorf("no session to fail")
	}

	session.mu.Lock()
	session.State = StateFailed
	session.UpdatedAt = time.Now().Format(time.RFC3339)
	session.mu.Unlock()

	return sm.saveSession(session)
}

// LoadSession loads the session state from disk
func (sm *SessionManager) LoadSession() (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionPath := sm.sessionPath()
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read session: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	sm.activeSession = &session
	return &session, nil
}

// sessionPath returns the path to the session state file
func (sm *SessionManager) sessionPath() string {
	return filepath.Join(sm.workspacePath, "esp-session.json")
}

// saveSession persists the session state to disk
func (sm *SessionManager) saveSession(session *Session) error {
	// Ensure workspace directory exists
	if err := os.MkdirAll(sm.workspacePath, 0755); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	sessionPath := sm.sessionPath()
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(sessionPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session: %w", err)
	}

	return nil
}

// EndSession clears the session state
func (sm *SessionManager) EndSession() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.activeSession = nil

	sessionPath := sm.sessionPath()
	if err := os.Remove(sessionPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove session: %w", err)
	}

	return nil
}

// GetSessionState returns the current session state as a string
func (sm *SessionManager) GetSessionState() SessionState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.activeSession == nil {
		return StateIdle
	}

	return sm.activeSession.State
}
