// Package esp implements the KDSE Engineering Session Protocol (ESP).
//
// This package provides a complete implementation of the ESP specification
// as defined in the KDSE Engineering Session Protocol document.
//
// Key features:
//   - Session state machine with all defined states and transitions
//   - Handshake protocol with all message types
//   - Bootstrap protocol with all required steps
//   - Authority hierarchy and conflict resolution
//   - Failure codes and recovery procedures
//   - Completion criteria verification
//
// Example usage:
//
//	// Create a new session
//	session := esp.NewSession("workspace-123")
//
//	// Transition through states
//	if err := session.TransitionTo(esp.StateDiscovering, "Starting session"); err != nil {
//	    log.Fatalf("Invalid transition: %v", err)
//	}
//
//	// Bootstrap the context
//	protocol := esp.NewBootstrapProtocol("/path/to/repo")
//	result, err := protocol.Execute(session.ID)
//	if err != nil {
//	    log.Fatalf("Bootstrap failed: %v", err)
//	}
//
//	// Verify completion criteria
//	verifier := esp.NewCompletionVerifier("/path/to/repo")
//	criteria := verifier.VerifyActiveState(result.Context)
//
//	fmt.Printf("Ready for AI reasoning: %v\n", criteria.AllPassed)
package esp

// Version information
var (
	// Version is the protocol implementation version
	Version = "1.0.0"

	// APIVersion is the supported API version
	APIVersion = "v1"
)
