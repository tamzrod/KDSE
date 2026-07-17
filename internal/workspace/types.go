package workspace

import (
	"errors"
	"time"
)

// RuntimeType represents the type of runtime
type RuntimeType string

const (
	RuntimeTypeCLI RuntimeType = "cli"
	RuntimeTypeMCP RuntimeType = "mcp"
)

// Errors
var (
	ErrRuntimeMissing     = errors.New("runtime not found: .kdse directory missing")
	ErrRuntimeInvalid     = errors.New("runtime configuration invalid")
	ErrVersionUnsupported = errors.New("runtime version not supported")
	ErrPhaseInvalid      = errors.New("invalid phase state")
	ErrInvalidTransition = errors.New("invalid phase transition")
	ErrIncompletePhase   = errors.New("phase has incomplete artifacts")
	ErrVerificationFailed = errors.New("workspace verification failed")
	ErrArtifactMissing   = errors.New("required artifact missing")
	ErrArtifactInvalid   = errors.New("artifact validation failed")
	ErrSessionNotFound   = errors.New("session not found")
	ErrWorkspaceNotFound = errors.New("workspace not found")
)

// InitOptions contains options for workspace initialization
type InitOptions struct {
	Path      string
	Type      RuntimeType
	Version   string
	Template  string
	Metadata  map[string]string
}

// VerificationResult contains the result of workspace verification
type VerificationResult struct {
	Valid       bool
	Phase       Phase
	Errors      []VerificationError
	Warnings    []VerificationWarning
	Timestamp   time.Time
	RuntimeInfo *RuntimeInfo
}

// VerificationError contains error details
type VerificationError struct {
	Code    string
	Message string
	Path    string
}

// VerificationWarning contains warning details
type VerificationWarning struct {
	Code    string
	Message string
	Path    string
}

// RuntimeInfo contains runtime information
type RuntimeInfo struct {
	Type    RuntimeType
	Version string
	Commit  string
}

// Transition represents a phase transition
type Transition struct {
	From      Phase
	To        Phase
	Timestamp time.Time
	Verified  bool
	Evidence  []string
}

// RuntimeContext represents the KDSE workspace runtime state
// This is the authoritative state container owned by the Workspace Engine
type RuntimeContext struct {
	Path    string
	Root    string
	Config  *WorkspaceConfig
	State   *WorkspaceState
	Runtime *RuntimeConfig
	Session *Session
}

// WorkspaceConfig contains workspace configuration
type WorkspaceConfig struct {
	Version string
	Root    string
}

// WorkspaceState contains current workspace state
type WorkspaceState struct {
	CurrentPhase     Phase
	PhaseHistory     []PhaseTransition
	LastVerification time.Time
	ArtifactCount    int
	ReportCount      int
	SessionCount     int
}

// PhaseTransition records a phase change
type PhaseTransition struct {
	From      Phase
	To        Phase
	Timestamp time.Time
}

// RuntimeConfig contains runtime configuration
type RuntimeConfig struct {
	Type            RuntimeType
	Version         string
	Commit          string
	TemplateVersion string
	TemplateCommit  string
}

// Session represents an engineering session
type Session struct {
	ID        string
	Created   time.Time
	Ended     time.Time
	Phase     Phase
	Actions   []SessionAction
	Metadata  map[string]string
}

// SessionAction records an action in a session
type SessionAction struct {
	Timestamp time.Time
	Action    string
	Phase     Phase
	Result    ActionResult
}

// ActionResult represents the result of an action
type ActionResult struct {
	Success bool
	Error  string
}

// ReportOptions contains options for report generation
type ReportOptions struct {
	Type   ReportType
	Output string
}

// ReportType represents the type of report
type ReportType string

const (
	ReportTypePhase         ReportType = "phase"
	ReportTypeVerification  ReportType = "verification"
	ReportTypeSummary       ReportType = "summary"
	ReportTypeProgress      ReportType = "progress"
)

// Report represents a generated report
type Report struct {
	Type      ReportType
	Title     string
	Content   string
	Generated time.Time
}

// ReportSummary contains report metadata
type ReportSummary struct {
	Type      ReportType
	Title     string
	Generated time.Time
}

// Artifact represents an artifact in the workspace
type Artifact struct {
	Path     string
	Type     string
	Size     int64
	Checksum string
}

// KnowledgeCollection contains collected knowledge
type KnowledgeCollection struct {
	Artifacts []KnowledgeArtifact
	Count    int
}

// KnowledgeArtifact represents a knowledge artifact
type KnowledgeArtifact struct {
	Path     string
	Title    string
	Summary  string
	Verified bool
}

// EngineError represents an error from the workspace engine
type EngineError struct {
	Code        string
	Message     string
	Details     map[string]interface{}
	Remediation string
}

func (e *EngineError) Error() string {
	return e.Message
}

// ArtifactSpec defines a required artifact
type ArtifactSpec struct {
	Path     string
	Type     string
	Required bool
}

// ValidationResult contains artifact validation results
type ValidationResult struct {
	Valid    bool
	Verified []Artifact
	Missing  []ArtifactSpec
	Invalid  []Artifact
}
