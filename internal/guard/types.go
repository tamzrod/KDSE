// Package guard provides the runtime guard architecture for KDSE.
// Guards enforce preconditions in deterministic order before any runtime operation.
//
// Architecture Overview:
//
//	Runtime Guard
//	    ├── Project Guard
//	    │      ├── Project Discovery
//	    │      └── Project Validation
//	    ├── Workspace Guard
//	    ├── Session Guard
//	    └── Lifecycle Guard
//
// Each guard validates one state transition in the following state model:
//
//	No Project → Project → Workspace → Session → Lifecycle Ready
package guard

import (
	"fmt"
)

// RuntimeState represents the deterministic state of the runtime
type RuntimeState string

const (
	// StateNoProject indicates no valid project has been detected
	StateNoProject RuntimeState = "NO_PROJECT"
	// StateProject indicates a valid project has been detected
	StateProject RuntimeState = "PROJECT"
	// StateWorkspace indicates KDSE workspace exists
	StateWorkspace RuntimeState = "WORKSPACE"
	// StateSession indicates an active session is running
	StateSession RuntimeState = "SESSION"
	// StateLifecycleReady indicates lifecycle is ready for operations
	StateLifecycleReady RuntimeState = "LIFECYCLE_READY"
)

// GuardType identifies which guard produced a result
type GuardType string

const (
	GuardTypeRuntime   GuardType = "RUNTIME"
	GuardTypeProject   GuardType = "PROJECT"
	GuardTypeWorkspace GuardType = "WORKSPACE"
	GuardTypeSession   GuardType = "SESSION"
	GuardTypeLifecycle GuardType = "LIFECYCLE"
)

// RuntimeGuardError represents an error from any guard
type RuntimeGuardError struct {
	GuardType GuardType
	Code      string
	Message   string
	Hint      string
	State     RuntimeState
}

func (e *RuntimeGuardError) Error() string {
	return fmt.Sprintf("[%s:%s] %s", e.GuardType, e.Code, e.Message)
}

// NewRuntimeGuardError creates a new RuntimeGuardError
func NewRuntimeGuardError(guardType GuardType, code, message, hint string, state RuntimeState) *RuntimeGuardError {
	return &RuntimeGuardError{
		GuardType: guardType,
		Code:      code,
		Message:   message,
		Hint:      hint,
		State:     state,
	}
}

// Common runtime guard errors
var (
	// Project-level errors
	ErrNoProjectDetected = NewRuntimeGuardError(
		GuardTypeProject,
		"NO_PROJECT",
		"No engineering project detected",
		"Ensure you are in a directory containing source code or project files",
		StateNoProject,
	)

	ErrInvalidProjectLocation = NewRuntimeGuardError(
		GuardTypeProject,
		"INVALID_LOCATION",
		"Project location is invalid or inaccessible",
		"Ensure the directory exists and is accessible",
		StateNoProject,
	)

	ErrGenericDirectory = NewRuntimeGuardError(
		GuardTypeProject,
		"GENERIC_DIRECTORY",
		"Directory does not appear to be an engineering project",
		"Ensure you are in a valid project directory with source code or documentation",
		StateNoProject,
	)

	// Workspace-level errors
	ErrWorkspaceMissing = NewRuntimeGuardError(
		GuardTypeWorkspace,
		"WORKSPACE_MISSING",
		"KDSE workspace not found",
		"Run 'kdse initialize' to create the workspace",
		StateProject,
	)

	ErrWorkspaceCorrupted = NewRuntimeGuardError(
		GuardTypeWorkspace,
		"WORKSPACE_CORRUPTED",
		"KDSE workspace is corrupted or incomplete",
		"Run 'kdse initialize' to recreate the workspace",
		StateProject,
	)

	ErrWorkspaceVersionMismatch = NewRuntimeGuardError(
		GuardTypeWorkspace,
		"VERSION_MISMATCH",
		"Workspace version is incompatible with runtime",
		"Consider running 'kdse migrate' or 'kdse initialize'",
		StateProject,
	)

	// Session-level errors
	ErrNoActiveSession = NewRuntimeGuardError(
		GuardTypeSession,
		"NO_SESSION",
		"No active KDSE session",
		"Run 'kdse initialize' to start a session",
		StateWorkspace,
	)

	ErrSessionExpired = NewRuntimeGuardError(
		GuardTypeSession,
		"SESSION_EXPIRED",
		"Session has expired",
		"Run 'kdse initialize' to start a new session",
		StateWorkspace,
	)

	ErrSessionInvalid = NewRuntimeGuardError(
		GuardTypeSession,
		"SESSION_INVALID",
		"Session state is invalid or corrupted",
		"Run 'kdse initialize' to reset the session",
		StateWorkspace,
	)

	// Lifecycle-level errors
	ErrLifecycleInvalid = NewRuntimeGuardError(
		GuardTypeLifecycle,
		"LIFECYCLE_INVALID",
		"Lifecycle state is invalid",
		"Check lifecycle configuration",
		StateSession,
	)

	ErrLifecycleTransitionInvalid = NewRuntimeGuardError(
		GuardTypeLifecycle,
		"TRANSITION_INVALID",
		"Invalid lifecycle phase transition",
		"Use 'kdse phase show' to see valid transitions",
		StateSession,
	)
)

// RuntimeGuardResult is the common result type returned by all guards in the new architecture.
// NOTE: The legacy SessionGuard has its own GuardResult type for backward compatibility.
// This type is used by RuntimeGuard, ProjectGuard, WorkspaceGuard, etc.
type RuntimeGuardResult struct {
	Valid       bool
	State       RuntimeState
	GuardType   GuardType
	StateBefore RuntimeState
	StateAfter  RuntimeState
	Error       *RuntimeGuardError
}

// ProjectGuardResult contains project-specific guard results
type ProjectGuardResult struct {
	RuntimeGuardResult
	ProjectPath  string
	ProjectName  string
	IsGitRepo    bool
	DetectedArtifacts []string
}

// WorkspaceGuardResult contains workspace-specific guard results
type WorkspaceGuardResult struct {
	RuntimeGuardResult
	WorkspaceRoot string
	Version       string
	Integrity     bool
}

// SessionGuardResult contains session-specific guard results
type SessionGuardResult struct {
	RuntimeGuardResult
	SessionID   string
	StartedAt   string
	WorkspaceRoot string
}

// LifecycleGuardResult contains lifecycle-specific guard results
type LifecycleGuardResult struct {
	RuntimeGuardResult
	CurrentPhase    string
	AllowedPhases   []string
	TransitionValid bool
}

// RuntimeValidationResult contains the combined results from all guards during a validation run.
type RuntimeValidationResult struct {
	Valid        bool
	FinalState   RuntimeState
	Project      *ProjectGuardResult
	Workspace    *WorkspaceGuardResult
	Session      *SessionGuardResult
	Lifecycle    *LifecycleGuardResult
	Errors       []*RuntimeGuardError
}

// AddError adds an error to the result
func (r *RuntimeValidationResult) AddError(err *RuntimeGuardError) {
	r.Errors = append(r.Errors, err)
	r.Valid = false
}
