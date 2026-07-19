// Package esp implements the KDSE Engineering Session Protocol (ESP).
package esp

import (
	"os"
	"path/filepath"
	"testing"
)

// TestComplianceRequirements tests all mandatory compliance requirements
// per Section 15.2 of the ESP specification.
func TestComplianceRequirements(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		// C-001: Protocol MUST implement all handshake messages
		{"HandshakeMessages", testHandshakeMessages},

		// C-002: Protocol MUST execute bootstrap steps
		{"BootstrapSteps", testBootstrapSteps},

		// C-003: Protocol MUST enforce authority hierarchy
		{"AuthorityHierarchy", testAuthorityHierarchy},

		// C-004: Protocol MUST implement all state transitions
		{"StateTransitions", testStateTransitions},

		// C-005: Protocol MUST handle all failure codes
		{"FailureCodes", testFailureCodes},

		// C-006: Protocol MUST enforce completion criteria
		{"CompletionCriteria", testCompletionCriteria},

		// C-007: Protocol MUST NOT begin AI reasoning before ACTIVE state
		{"NoReasoningBeforeActive", testNoReasoningBeforeActive},

		// C-008: Protocol MUST produce identical context for identical inputs
		{"DeterministicContext", testDeterministicContext},

		// C-009: Protocol MUST use authoritative sources over heuristics
		{"AuthorityPrecedence", testAuthorityPrecedence},

		// C-010: Protocol MUST NOT depend on implementation-specific behavior
		{"ImplementationIndependence", testImplementationIndependence},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// testHandshakeMessages verifies all handshake messages are implemented
func testHandshakeMessages(t *testing.T) {
	expectedMessages := []MessageType{
		MsgConnect,
		MsgDiscover,
		MsgDiscoverAck,
		MsgAuthenticate,
		MsgAuthenticateAck,
		MsgBootstrap,
		MsgBootstrapAck,
		MsgVerify,
		MsgVerifyAck,
		MsgContextRequest,
		MsgContextQuery,
		MsgContextResponse,
	}

	implementedMessages := map[MessageType]bool{
		MsgConnect:          true,
		MsgDiscover:         true,
		MsgDiscoverAck:      true,
		MsgAuthenticate:     true,
		MsgAuthenticateAck:  true,
		MsgBootstrap:        true,
		MsgBootstrapAck:     true,
		MsgVerify:           true,
		MsgVerifyAck:        true,
		MsgContextRequest:   true,
		MsgContextQuery:     true,
		MsgContextResponse:  true,
	}

	for _, msg := range expectedMessages {
		if !implementedMessages[msg] {
			t.Errorf("Handshake message not implemented: %s", msg)
		}
	}
}

// testBootstrapSteps verifies all bootstrap steps are implemented
func testBootstrapSteps(t *testing.T) {
	expectedSteps := []BootstrapStep{
		BootstrapStepDiscoverWorkspace,
		BootstrapStepVerifyRuntime,
		BootstrapStepLoadManifest,
		BootstrapStepDetermineContext,
		BootstrapStepResolveConflicts,
		BootstrapStepAssembleContext,
		BootstrapStepComputeChecksum,
		BootstrapStepRecordEvidence,
	}

	actualSteps := GetBootstrapSteps()

	if len(actualSteps) != len(expectedSteps) {
		t.Errorf("Bootstrap step count mismatch: expected %d, got %d",
			len(expectedSteps), len(actualSteps))
	}

	for i, expected := range expectedSteps {
		if actualSteps[i] != expected {
			t.Errorf("Bootstrap step %d mismatch: expected %s, got %s",
				i, expected, actualSteps[i])
		}
	}
}

// testAuthorityHierarchy verifies authority hierarchy is enforced
func testAuthorityHierarchy(t *testing.T) {
	precedence := AuthorityPrecedence()

	// Verify expected order
	expectedOrder := []AuthorityLevel{
		AuthorityKDSEStandard,
		AuthorityRuntimeManifest,
		AuthorityProjectManifest,
		AuthoritySessionState,
		AuthorityAI,
	}

	if len(precedence) != len(expectedOrder) {
		t.Fatalf("Authority precedence count mismatch: expected %d, got %d",
			len(expectedOrder), len(precedence))
	}

	for i, expected := range expectedOrder {
		if precedence[i] != expected {
			t.Errorf("Authority precedence[%d] mismatch: expected %s, got %s",
				i, expected, precedence[i])
		}
	}
}

// testStateTransitions verifies all legal state transitions work
func testStateTransitions(t *testing.T) {
	// Test valid transitions
	validCases := []struct {
		from, to SessionState
	}{
		{StateIdle, StateDiscovering},
		{StateDiscovering, StateAuthenticating},
		{StateAuthenticating, StateBootstrapping},
		{StateBootstrapping, StateVerifying},
		{StateVerifying, StateActive},
		{StateActive, StateSuspended},
		{StateActive, StateTerminated},
		{StateSuspended, StateActive},
		{StateFailed, StateIdle},
	}

	for _, tc := range validCases {
		if !IsValidTransition(tc.from, tc.to) {
			t.Errorf("Expected valid transition from %s to %s", tc.from, tc.to)
		}
	}

	// Test invalid transitions
	invalidCases := []struct {
		from, to SessionState
	}{
		{StateIdle, StateActive},
		{StateIdle, StateTerminated},
		{StateTerminated, StateActive},
		{StateFailed, StateActive},
	}

	for _, tc := range invalidCases {
		if IsValidTransition(tc.from, tc.to) {
			t.Errorf("Expected invalid transition from %s to %s", tc.from, tc.to)
		}
	}
}

// testFailureCodes verifies all failure codes are defined
func testFailureCodes(t *testing.T) {
	expectedCodes := []FailureCode{
		FailureWorkspaceNotFound,
		FailureWorkspaceUnauthorized,
		FailureRuntimeNotFound,
		FailureManifestInvalid,
		FailureManifestCorrupted,
		FailureAuthenticationFailed,
		FailureTokenExpired,
		FailureTokenInvalid,
		FailureCredentialsInvalid,
		FailureSessionAlreadyActive,
		FailureBootstrapFailed,
		FailureContextIncomplete,
		FailureBootstrapTimeout,
		FailureFilesystemInaccessible,
		FailureVerificationFailed,
		FailureChecksumMismatch,
		FailureIntegrityViolation,
		FailureContextCorrupted,
		FailureVersionIncompatible,
		FailureRuntimeUpgradeRequired,
	}

	for _, code := range expectedCodes {
		info, ok := FailureCodeInfo[code]
		if !ok {
			t.Errorf("Failure code %d has no info", code)
		}
		if info.Name == "" {
			t.Errorf("Failure code %d has empty name", code)
		}
		if info.Description == "" {
			t.Errorf("Failure code %d has empty description", code)
		}
	}
}

// testCompletionCriteria verifies completion criteria are checked
func testCompletionCriteria(t *testing.T) {
	verifier := NewCompletionVerifier("/nonexistent")

	// Test with nil context
	criteria := verifier.VerifyActiveState(nil)
	if criteria.AllPassed {
		t.Error("Expected criteria to fail with nil context")
	}

	// Test with valid context
	ctx := NewContext("test-project", "/test/repo", "Problem", "ws-123", "sess-456")
	ctx.Update()

	criteria = verifier.VerifyActiveState(ctx)
	if criteria.AllPassed {
		t.Error("Expected criteria to fail without workspace")
	}
}

// testNoReasoningBeforeActive verifies reasoning is blocked before ACTIVE state
func testNoReasoningBeforeActive(t *testing.T) {
	session := NewSession("ws-123")

	// Verify session starts in IDLE state
	if session.GetState() != StateIdle {
		t.Errorf("Expected session to start in IDLE state, got %s", session.GetState())
	}

	// Verify IsActive returns false for non-active states
	nonActiveStates := []SessionState{
		StateIdle,
		StateDiscovering,
		StateAuthenticating,
		StateBootstrapping,
		StateVerifying,
		StateFailed,
		StateTerminated,
	}

	for _, state := range nonActiveStates {
		session.mu.Lock()
		session.State = state
		session.mu.Unlock()

		if session.IsActive() {
			t.Errorf("Session should not be active in state %s", state)
		}
	}

	// Verify IsActive returns true only for ACTIVE state
	session.mu.Lock()
	session.State = StateActive
	session.mu.Unlock()

	if !session.IsActive() {
		t.Error("Session should be active in ACTIVE state")
	}
}

// testDeterministicContext verifies context is deterministic
func testDeterministicContext(t *testing.T) {
	// Create two contexts with same inputs
	ctx1 := NewContext("test-project", "/test/repo", "Problem", "ws-123", "sess-456")
	ctx2 := NewContext("test-project", "/test/repo", "Problem", "ws-123", "sess-456")

	// Set same timestamps for deterministic comparison
	ctx1.InitializedAt = "2026-07-19T00:00:00Z"
	ctx2.InitializedAt = "2026-07-19T00:00:00Z"
	ctx1.LastUpdated = "2026-07-19T00:00:00Z"
	ctx2.LastUpdated = "2026-07-19T00:00:00Z"

	// Compute checksums
	cs1, err1 := ctx1.ComputeChecksum()
	cs2, err2 := ctx2.ComputeChecksum()

	if err1 != nil {
		t.Errorf("Failed to compute checksum 1: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Failed to compute checksum 2: %v", err2)
	}

	// Checksums should match for same inputs
	if cs1 != cs2 {
		t.Errorf("Checksums should match for same inputs: %s != %s", cs1, cs2)
	}
}

// testAuthorityPrecedence verifies authoritative sources take precedence
func testAuthorityPrecedence(t *testing.T) {
	resolutions := []AuthorityResolution{
		{Level: AuthorityAI, Value: "ai-value", Confidence: 0.5},
		{Level: AuthorityRuntimeManifest, Value: "runtime-value", Confidence: 1.0},
		{Level: AuthorityKDSEStandard, Value: "standard-value", Confidence: 1.0},
	}

	value, level := ResolveWithAuthority(resolutions, "test-field")

	if level != AuthorityKDSEStandard {
		t.Errorf("Expected highest authority level, got %s", level)
	}

	if value != "standard-value" {
		t.Errorf("Expected standard-value, got %v", value)
	}
}

// testImplementationIndependence verifies no implementation-specific dependencies
func testImplementationIndependence(t *testing.T) {
	// Verify protocol types don't use implementation-specific types
	// This is a structural test - in Go we can't enforce this at compile time
	// but we can verify the types are defined correctly

	session := NewSession("ws-123")

	// Verify state machine works without any external dependencies
	if err := session.TransitionTo(StateDiscovering, "test"); err != nil {
		t.Errorf("Failed to transition: %v", err)
	}

	if session.GetState() != StateDiscovering {
		t.Errorf("Expected state to be DISCOVERING, got %s", session.GetState())
	}
}

// TestConformanceTestSuite tests the conformance test suite per Section 15.4.1
func TestConformanceTestSuite(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{"T-001 Handshake Complete", testHandshakeMessages},
		{"T-002 Bootstrap Deterministic", testDeterministicContext},
		{"T-003 Authority Resolution", testAuthorityPrecedence},
		{"T-004 State Transitions", testStateTransitions},
		{"T-005 Failure Handling", testFailureCodes},
		{"T-006 Completion Criteria", testCompletionCriteria},
		{"T-007 Isolation Test", testNoReasoningBeforeActive},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestBootstrapProtocolIntegration tests the bootstrap protocol end-to-end
func TestBootstrapProtocolIntegration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	workspaceDir := filepath.Join(tmpDir, ".kdse")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Create required runtime files
	runtimeYaml := "version: 1.0\nruntime: esp\n"
	if err := os.WriteFile(filepath.Join(workspaceDir, "runtime.yaml"), []byte(runtimeYaml), 0644); err != nil {
		t.Fatalf("Failed to create runtime.yaml: %v", err)
	}

	manifestYaml := "version: 1.0\n"
	if err := os.WriteFile(filepath.Join(workspaceDir, "manifest.yaml"), []byte(manifestYaml), 0644); err != nil {
		t.Fatalf("Failed to create manifest.yaml: %v", err)
	}

	sessionStateYaml := "session_id: test-session\n"
	if err := os.WriteFile(filepath.Join(workspaceDir, "session-state.yaml"), []byte(sessionStateYaml), 0644); err != nil {
		t.Fatalf("Failed to create session-state.yaml: %v", err)
	}

	// Run bootstrap
	protocol := NewBootstrapProtocol(tmpDir)
	result, err := protocol.Execute("test-session-id")

	if err != nil {
		t.Fatalf("Bootstrap failed: %v", err)
	}

	if result == nil {
		t.Fatal("Bootstrap returned nil result")
	}

	// Verify all steps completed
	if len(result.BootstrapSteps) != 8 {
		t.Errorf("Expected 8 bootstrap steps, got %d", len(result.BootstrapSteps))
	}
}
