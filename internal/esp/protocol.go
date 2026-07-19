// Package esp implements the KDSE Engineering Session Protocol (ESP).
// ESP establishes the complete contractual framework for engineering sessions,
// encompassing identity establishment, context establishment, session lifecycle,
// failure handling, and completion criteria.
//
// This implementation is compliant with KDSE Engineering Session Protocol
// Specification version 1.0.
package esp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Protocol version constants
const (
	ProtocolVersion   = "1.0"
	ProtocolName      = "KDSE Engineering Session Protocol"
	ProtocolShortName = "ESP"
)

// =============================================================================
// Session States (Section 11)
// =============================================================================

// SessionState represents the complete set of valid session states
// as defined in Section 11 of the ESP specification.
type SessionState string

const (
	// StateIdle - No active engineering session exists
	StateIdle SessionState = "IDLE"

	// StateDiscovering - Runtime is locating and identifying the workspace
	StateDiscovering SessionState = "DISCOVERING"

	// StateAuthenticating - Runtime is verifying workspace ownership
	StateAuthenticating SessionState = "AUTHENTICATING"

	// StateBootstrapping - Runtime is establishing engineering context
	StateBootstrapping SessionState = "BOOTSTRAPPING"

	// StateVerifying - Runtime is validating engineering context completeness
	StateVerifying SessionState = "VERIFYING"

	// StateActive - Engineering session is established and AI reasoning may proceed
	StateActive SessionState = "ACTIVE"

	// StateSuspended - Session is paused but may be resumed
	StateSuspended SessionState = "SUSPENDED"

	// StateTerminated - Session has ended and cannot be resumed
	StateTerminated SessionState = "TERMINATED"

	// StateFailed - An unrecoverable error has occurred
	StateFailed SessionState = "FAILED"
)

// IsTerminal returns true if the state is a terminal state
func (s SessionState) IsTerminal() bool {
	return s == StateTerminated || s == StateFailed
}

// IsActive returns true if the session is in an active or resumable state
func (s SessionState) IsActive() bool {
	return s == StateActive || s == StateSuspended
}

// String returns the string representation of the state
func (s SessionState) String() string {
	return string(s)
}

// ValidStates returns all valid session states
func ValidStates() []SessionState {
	return []SessionState{
		StateIdle,
		StateDiscovering,
		StateAuthenticating,
		StateBootstrapping,
		StateVerifying,
		StateActive,
		StateSuspended,
		StateTerminated,
		StateFailed,
	}
}

// =============================================================================
// State Transitions (Section 11.1)
// =============================================================================

// StateTransition represents a valid state transition
type StateTransition struct {
	From   SessionState
	To     SessionState
	Reason string
}

// ValidTransitions defines all legal state transitions per Section 11.1
var ValidTransitions = []StateTransition{
	// From IDLE
	{From: StateIdle, To: StateDiscovering, Reason: "Session initiation"},
	{From: StateIdle, To: StateFailed, Reason: "Fatal error during initialization"},

	// From DISCOVERING
	{From: StateDiscovering, To: StateAuthenticating, Reason: "Workspace located"},
	{From: StateDiscovering, To: StateFailed, Reason: "Workspace not found or inaccessible"},

	// From AUTHENTICATING
	{From: StateAuthenticating, To: StateBootstrapping, Reason: "Ownership verified"},
	{From: StateAuthenticating, To: StateFailed, Reason: "Authentication failed"},

	// From BOOTSTRAPPING
	{From: StateBootstrapping, To: StateVerifying, Reason: "Context established"},
	{From: StateBootstrapping, To: StateFailed, Reason: "Bootstrap failed"},

	// From VERIFYING
	{From: StateVerifying, To: StateActive, Reason: "Verification successful"},
	{From: StateVerifying, To: StateFailed, Reason: "Verification failed"},

	// From ACTIVE
	{From: StateActive, To: StateSuspended, Reason: "Session suspended by user"},
	{From: StateActive, To: StateTerminated, Reason: "Session completed normally"},
	{From: StateActive, To: StateFailed, Reason: "Fatal error during session"},

	// From SUSPENDED
	{From: StateSuspended, To: StateActive, Reason: "Session resumed"},
	{From: StateSuspended, To: StateTerminated, Reason: "Session abandoned"},
	{From: StateSuspended, To: StateFailed, Reason: "Fatal error during resumption"},

	// From FAILED
	{From: StateFailed, To: StateIdle, Reason: "Retry session"},

	// From TERMINATED
	// No transitions - terminal state
}

// IsValidTransition checks if a transition is legally valid
func IsValidTransition(from, to SessionState) bool {
	for _, t := range ValidTransitions {
		if t.From == from && t.To == to {
			return true
		}
	}
	return false
}

// GetValidNextStates returns all valid states that can be transitioned to from the given state
func GetValidNextStates(from SessionState) []SessionState {
	var states []SessionState
	seen := make(map[SessionState]bool)
	for _, t := range ValidTransitions {
		if t.From == from && !seen[t.To] {
			states = append(states, t.To)
			seen[t.To] = true
		}
	}
	return states
}

// =============================================================================
// Handshake Messages (Section 8)
// =============================================================================

// MessageType represents the type of handshake message
type MessageType string

const (
	// Runtime Adapter → KDSE Runtime
	MsgDiscover        MessageType = "DISCOVER"
	MsgAuthenticate    MessageType = "AUTHENTICATE"
	MsgBootstrap       MessageType = "BOOTSTRAP"
	MsgVerify          MessageType = "VERIFY"

	// KDSE Runtime → Runtime Adapter
	MsgDiscoverAck        MessageType = "DISCOVER_ACK"
	MsgAuthenticateAck    MessageType = "AUTHENTICATE_ACK"
	MsgBootstrapAck       MessageType = "BOOTSTRAP_ACK"
	MsgVerifyAck         MessageType = "VERIFY_ACK"

	// AI Agent ↔ Runtime Adapter
	MsgConnect          MessageType = "CONNECT"
	MsgContextRequest   MessageType = "CONTEXT_REQUEST"
	MsgContextQuery     MessageType = "CONTEXT_QUERY"
	MsgContextResponse  MessageType = "CONTEXT_RESPONSE"
)

// HandshakeMessage represents a protocol message
type HandshakeMessage struct {
	Type      MessageType `json:"type"`
	Timestamp string      `json:"timestamp"`
	ID        string      `json:"id"`
	Payload   interface{} `json:"payload,omitempty"`
}

// NewMessage creates a new handshake message with generated ID and timestamp
func NewMessage(msgType MessageType, payload interface{}) *HandshakeMessage {
	return &HandshakeMessage{
		Type:      msgType,
		Timestamp: time.Now().Format(time.RFC3339),
		ID:        uuid.New().String(),
		Payload:   payload,
	}
}

// DISCOVER message (Section 8.3.2)
type DiscoverPayload struct {
	WorkspacePath string `json:"workspace_path"`
	WorkspaceType string `json:"workspace_type"` // "local" or "remote"
}

// DISCOVER_ACK message (Section 8.3.3)
type DiscoverAckPayload struct {
	Status         string `json:"status"` // "found", "not_found", "unauthorized"
	WorkspaceID    string `json:"workspace_id,omitempty"`
	RuntimeVersion string `json:"runtime_version"`
}

// AUTHENTICATE message (Section 8.3.4)
type AuthenticatePayload struct {
	WorkspaceID  string `json:"workspace_id"`
	AuthMethod   string `json:"auth_method"`
	Credentials  string `json:"credentials,omitempty"`
}

// AUTHENTICATE_ACK message (Section 8.3.5)
type AuthenticateAckPayload struct {
	Status   string `json:"status"` // "authenticated", "failed", "expired"
	AuthToken string `json:"auth_token,omitempty"`
}

// BOOTSTRAP message (Section 8.3.6)
type BootstrapPayload struct {
	WorkspaceID string   `json:"workspace_id"`
	AuthToken   string   `json:"auth_token"`
	Context     *Context `json:"context,omitempty"`
}

// BOOTSTRAP_ACK message (Section 8.3.7)
type BootstrapAckPayload struct {
	Status         string   `json:"status"` // "success", "partial", "failed"
	Context        *Context `json:"context,omitempty"`
	MissingFields  []string `json:"missing_fields,omitempty"`
	Checksum       string   `json:"checksum"`
}

// VERIFY message (Section 8.3.8)
type VerifyPayload struct {
	ContextChecksum string `json:"context_checksum"`
	RequiredFields []string `json:"required_fields"`
}

// VERIFY_ACK message (Section 8.3.9)
type VerifyAckPayload struct {
	Status       string   `json:"status"` // "valid", "invalid", "incomplete"
	Errors       []string `json:"errors,omitempty"`
	MissingFields []string `json:"missing_fields,omitempty"`
	SessionToken string   `json:"session_token,omitempty"`
}

// CONNECT message (Section 8.3.1)
type ConnectPayload struct {
	ClientID     string   `json:"client_id"`
	Capabilities []string `json:"capabilities"`
}

// CONTEXT_REQUEST message (Section 8.3.10)
type ContextRequestPayload struct {
	SessionToken string `json:"session_token,omitempty"`
	Fields       []string `json:"fields,omitempty"`
}

// CONTEXT_QUERY message (Section 8.3.11)
type ContextQueryPayload struct {
	SessionToken string   `json:"session_token"`
	Query        []string `json:"query"`
}

// CONTEXT_RESPONSE message (Section 8.3.12)
type ContextResponsePayload struct {
	Status  string   `json:"status"`
	Context *Context `json:"context,omitempty"`
	Query   []string `json:"query,omitempty"`
}

// =============================================================================
// Engineering Context (Section 9.2)
// =============================================================================

// Context represents the complete engineering context
type Context struct {
	// Required fields
	Project      string `json:"project"`
	Repository   string `json:"repository"`
	Phase        string `json:"phase"`
	WorkspaceID  string `json:"workspace_id"`
	SessionID    string `json:"session_id"`

	// Optional fields
	ProjectVersion  string   `json:"project_version,omitempty"`
	SessionToken    string   `json:"session_token,omitempty"`
	SandboxID       string   `json:"sandbox_id,omitempty"`
	Evidence        []string `json:"evidence,omitempty"`
	AllowedContext  []string `json:"allowed_context,omitempty"`
	RestrictedPaths []string `json:"restricted_paths,omitempty"`
	Artifacts       *ArtifactPaths `json:"artifacts,omitempty"`

	// Lifecycle metadata
	InitializedAt  string `json:"initialized_at"`
	StartedAt     string `json:"started_at"`
	LastUpdated   string `json:"last_updated"`
	Checksum      string `json:"checksum,omitempty"`
}

// ArtifactPaths defines artifact directory locations
type ArtifactPaths struct {
	Reports     string `json:"reports"`
	Screenshots string `json:"screenshots"`
	Tests       string `json:"tests"`
	Benchmarks  string `json:"benchmarks"`
}

// DefaultArtifactPaths returns the standard artifact paths
func DefaultArtifactPaths(workspacePath string) *ArtifactPaths {
	return &ArtifactPaths{
		Reports:     workspacePath + "/reports/",
		Screenshots: workspacePath + "/evidence/screenshots/",
		Tests:       workspacePath + "/evidence/tests/",
		Benchmarks:  workspacePath + "/evidence/benchmarks/",
	}
}

// NewContext creates a new engineering context
func NewContext(project, repository, phase, workspaceID, sessionID string) *Context {
	now := time.Now().Format(time.RFC3339)
	return &Context{
		Project:        project,
		Repository:     repository,
		Phase:          phase,
		WorkspaceID:    workspaceID,
		SessionID:      sessionID,
		InitializedAt:  now,
		StartedAt:      now,
		LastUpdated:    now,
		Evidence:       []string{},
		AllowedContext: []string{},
		RestrictedPaths: []string{},
	}
}

// Validate validates the context and returns any missing required fields
func (c *Context) Validate() []string {
	var missing []string

	if c.Project == "" {
		missing = append(missing, "project")
	}
	if c.Repository == "" {
		missing = append(missing, "repository")
	}
	if c.Phase == "" {
		missing = append(missing, "phase")
	}
	if c.WorkspaceID == "" {
		missing = append(missing, "workspace_id")
	}
	if c.SessionID == "" {
		missing = append(missing, "session_id")
	}

	return missing
}

// IsValid returns true if all required fields are present
func (c *Context) IsValid() bool {
	return len(c.Validate()) == 0
}

// Update updates the context timestamp and returns an error if the checksum fails to compute
func (c *Context) Update() error {
	c.LastUpdated = time.Now().Format(time.RFC3339)
	checksum, err := c.ComputeChecksum()
	if err != nil {
		return err
	}
	c.Checksum = checksum
	return nil
}

// ComputeChecksum computes a checksum of the context for integrity verification
func (c *Context) ComputeChecksum() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context: %w", err)
	}
	return fmt.Sprintf("%x", uuid.NewMD5(uuid.NameSpaceOID, data)), nil
}

// =============================================================================
// Authority Hierarchy (Section 7)
// =============================================================================

// AuthorityLevel defines the authority hierarchy
type AuthorityLevel string

const (
	// AuthorityKDSEStandard is the highest authority
	AuthorityKDSEStandard AuthorityLevel = "KDSE_STANDARD"

	// AuthorityRuntimeManifest supersedes AI claims
	AuthorityRuntimeManifest AuthorityLevel = "RUNTIME_MANIFEST"

	// AuthorityProjectManifest defines project-specific parameters
	AuthorityProjectManifest AuthorityLevel = "PROJECT_MANIFEST"

	// AuthoritySessionState defines current session context
	AuthoritySessionState AuthorityLevel = "SESSION_STATE"

	// AuthorityAI is the lowest authority
	AuthorityAI AuthorityLevel = "AI"
)

// AuthorityPrecedence returns the precedence order of authority levels
func AuthorityPrecedence() []AuthorityLevel {
	return []AuthorityLevel{
		AuthorityKDSEStandard,
		AuthorityRuntimeManifest,
		AuthorityProjectManifest,
		AuthoritySessionState,
		AuthorityAI,
	}
}

// CompareAuthority compares two authority levels
// Returns negative if a < b, zero if a == b, positive if a > b
func CompareAuthority(a, b AuthorityLevel) int {
	precedence := AuthorityPrecedence()
	aIdx, bIdx := -1, -1

	for i, level := range precedence {
		if level == a {
			aIdx = i
		}
		if level == b {
			bIdx = i
		}
	}

	if aIdx < 0 || bIdx < 0 {
		return 0
	}

	return bIdx - aIdx // Lower index = higher authority
}

// AuthorityResolution represents an authoritative source
type AuthorityResolution struct {
	Level      AuthorityLevel `json:"level"`
	Source     string         `json:"source"`
	Value      interface{}    `json:"value"`
	Confidence float64        `json:"confidence"`
	Timestamp  string         `json:"timestamp"`
}

// ResolveWithAuthority resolves a value using authority precedence
func ResolveWithAuthority(resolutions []AuthorityResolution, field string) (interface{}, AuthorityLevel) {
	if len(resolutions) == 0 {
		return nil, ""
	}

	// Find the resolution with highest authority
	var best AuthorityResolution
	for _, r := range resolutions {
		if best.Level == "" || CompareAuthority(r.Level, best.Level) > 0 {
			best = r
		}
	}

	return best.Value, best.Level
}
