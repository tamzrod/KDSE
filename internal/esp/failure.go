// Package esp implements the KDSE Engineering Session Protocol (ESP).
package esp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Failure codes as defined in Section 13 of the ESP specification
type FailureCode int

const (
	// Discovery failures (1001-1099)
	FailureWorkspaceNotFound    FailureCode = 1001
	FailureWorkspaceUnauthorized FailureCode = 1002
	FailureRuntimeNotFound      FailureCode = 1003
	FailureManifestInvalid      FailureCode = 1004
	FailureManifestCorrupted    FailureCode = 1005

	// Authentication failures (1101-1199)
	FailureAuthenticationFailed   FailureCode = 1101
	FailureTokenExpired           FailureCode = 1102
	FailureTokenInvalid           FailureCode = 1103
	FailureCredentialsInvalid     FailureCode = 1104
	FailureSessionAlreadyActive   FailureCode = 1105

	// Bootstrap failures (1201-1299)
	FailureBootstrapFailed        FailureCode = 1201
	FailureContextIncomplete      FailureCode = 1202
	FailureBootstrapTimeout       FailureCode = 1203
	FailureFilesystemInaccessible FailureCode = 1204

	// Verification failures (1301-1399)
	FailureVerificationFailed    FailureCode = 1301
	FailureChecksumMismatch       FailureCode = 1302
	FailureIntegrityViolation     FailureCode = 1303
	FailureContextCorrupted       FailureCode = 1304

	// Version/Compatibility failures (1501-1599)
	FailureVersionIncompatible    FailureCode = 1501
	FailureRuntimeUpgradeRequired FailureCode = 1502
)

// Failure represents a protocol failure
type Failure struct {
	Code        FailureCode `json:"code"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Details     FailureDetails `json:"details"`
	Remediation *Remediation `json:"remediation,omitempty"`
	Timestamp   string      `json:"timestamp"`
	RequestID   string      `json:"request_id"`
}

// FailureDetails contains diagnostic information about the failure
type FailureDetails struct {
	Diagnostic string      `json:"diagnostic"`
	Source     string      `json:"source"`
	Expected   string      `json:"expected"`
	Actual     string      `json:"actual"`
}

// Remediation describes how to recover from a failure
type Remediation struct {
	Action string   `json:"action"`
	Owner  string   `json:"owner"`
	Steps  []string `json:"steps"`
}

// FailureCodeInfo maps failure codes to their descriptions
var FailureCodeInfo = map[FailureCode]struct {
	Name        string
	Description string
}{
	FailureWorkspaceNotFound:    {"WORKSPACE_NOT_FOUND", "The specified workspace does not exist"},
	FailureWorkspaceUnauthorized: {"WORKSPACE_UNAUTHORIZED", "The user does not have access to the workspace"},
	FailureRuntimeNotFound:      {"RUNTIME_NOT_FOUND", "The KDSE runtime is not installed or not accessible"},
	FailureManifestInvalid:       {"MANIFEST_INVALID", "The runtime manifest is invalid or malformed"},
	FailureManifestCorrupted:    {"MANIFEST_CORRUPTED", "The runtime manifest has been corrupted"},

	FailureAuthenticationFailed:  {"AUTHENTICATION_FAILED", "Workspace authentication failed"},
	FailureTokenExpired:          {"TOKEN_EXPIRED", "The session token has expired"},
	FailureTokenInvalid:          {"TOKEN_INVALID", "The session token is invalid"},
	FailureCredentialsInvalid:    {"CREDENTIALS_INVALID", "The provided credentials are invalid"},
	FailureSessionAlreadyActive:   {"SESSION_ALREADY_ACTIVE", "A session is already active for this workspace"},

	FailureBootstrapFailed:        {"BOOTSTRAP_FAILED", "The bootstrap process failed"},
	FailureContextIncomplete:      {"CONTEXT_INCOMPLETE", "The engineering context is incomplete"},
	FailureBootstrapTimeout:       {"BOOTSTRAP_TIMEOUT", "The bootstrap process timed out"},
	FailureFilesystemInaccessible: {"FILESYSTEM_INACCESSIBLE", "The workspace filesystem is not accessible"},

	FailureVerificationFailed:    {"VERIFICATION_FAILED", "Context verification failed"},
	FailureChecksumMismatch:      {"CHECKSUM_MISMATCH", "The context checksum does not match"},
	FailureIntegrityViolation:    {"INTEGRITY_VIOLATION", "An integrity violation was detected"},
	FailureContextCorrupted:      {"CONTEXT_CORRUPTED", "The engineering context is corrupted"},

	FailureVersionIncompatible:    {"VERSION_INCOMPATIBLE", "The runtime version is incompatible"},
	FailureRuntimeUpgradeRequired: {"UPGRADE_REQUIRED", "A runtime upgrade is required"},
}

// Recovery procedures for each failure code
var RecoveryProcedures = map[FailureCode]Remediation{
	FailureWorkspaceNotFound: {
		Action: "Create or specify a valid workspace",
		Owner:  "User",
		Steps: []string{
			"Verify the workspace path is correct",
			"Initialize a new workspace with 'kdse init'",
			"Specify the correct project path",
		},
	},
	FailureWorkspaceUnauthorized: {
		Action: "Obtain proper workspace access",
		Owner:  "User",
		Steps: []string{
			"Verify you have read permissions for the workspace",
			"Request access from the workspace owner",
			"Use correct credentials for authentication",
		},
	},
	FailureRuntimeNotFound: {
		Action: "Install or locate KDSE runtime",
		Owner:  "Runtime",
		Steps: []string{
			"Verify KDSE is installed: 'kdse version'",
			"Install KDSE if not present",
			"Ensure KDSE is in your PATH",
		},
	},
	FailureManifestInvalid: {
		Action: "Fix or regenerate runtime manifest",
		Owner:  "Runtime",
		Steps: []string{
			"Check .kdse/manifest.yaml for syntax errors",
			"Regenerate manifest with 'kdse init'",
			"Validate manifest against schema",
		},
	},
	FailureManifestCorrupted: {
		Action: "Restore manifest from backup or regenerate",
		Owner:  "Runtime",
		Steps: []string{
			"Check for backup manifest in .kdse/backups/",
			"Restore from backup if available",
			"Regenerate manifest with 'kdse init'",
		},
	},

	FailureAuthenticationFailed: {
		Action: "Re-authenticate with valid credentials",
		Owner:  "User",
		Steps: []string{
			"Provide valid workspace credentials",
			"Verify authentication token is current",
			"Retry authentication",
		},
	},
	FailureTokenExpired: {
		Action: "Refresh session token",
		Owner:  "Runtime",
		Steps: []string{
			"Obtain a new session token",
			"Re-authenticate if necessary",
			"Resume session with new token",
		},
	},
	FailureTokenInvalid: {
		Action: "Obtain valid session token",
		Owner:  "Runtime",
		Steps: []string{
			"Verify token format and encoding",
			"Obtain a new valid token",
			"Retry with correct token",
		},
	},
	FailureCredentialsInvalid: {
		Action: "Provide valid credentials",
		Owner:  "User",
		Steps: []string{
			"Verify credentials are correct",
			"Update credentials if necessary",
			"Retry authentication",
		},
	},
	FailureSessionAlreadyActive: {
		Action: "Use existing session or terminate it",
		Owner:  "User",
		Steps: []string{
			"Check for existing session: 'kdse status'",
			"Terminate existing session: 'kdse session end'",
			"Resume existing session if appropriate",
		},
	},

	FailureBootstrapFailed: {
		Action: "Retry bootstrap with corrected inputs",
		Owner:  "Runtime",
		Steps: []string{
			"Identify specific bootstrap error",
			"Fix or provide missing information",
			"Retry bootstrap process",
		},
	},
	FailureContextIncomplete: {
		Action: "Provide missing context fields",
		Owner:  "User",
		Steps: []string{
			"Identify missing required fields",
			"Provide missing information",
			"Retry context assembly",
		},
	},
	FailureBootstrapTimeout: {
		Action: "Retry bootstrap with extended timeout",
		Owner:  "Runtime",
		Steps: []string{
			"Retry bootstrap process",
			"Check system resources",
			"Reduce workspace complexity if needed",
		},
	},
	FailureFilesystemInaccessible: {
		Action: "Restore filesystem access",
		Owner:  "System",
		Steps: []string{
			"Verify filesystem is mounted",
			"Check file permissions",
			"Ensure sufficient disk space",
		},
	},

	FailureVerificationFailed: {
		Action: "Fix context issues and re-verify",
		Owner:  "Runtime",
		Steps: []string{
			"Identify verification errors",
			"Fix context issues",
			"Retry verification",
		},
	},
	FailureChecksumMismatch: {
		Action: "Recompute context checksum",
		Owner:  "Runtime",
		Steps: []string{
			"Verify context integrity",
			"Recompute checksum",
			"Update stored checksum if needed",
		},
	},
	FailureIntegrityViolation: {
		Action: "Investigate and resolve integrity issue",
		Owner:  "Runtime",
		Steps: []string{
			"Identify the integrity violation",
			"Restore corrupted artifacts from backup",
			"Re-establish integrity",
		},
	},
	FailureContextCorrupted: {
		Action: "Restore context from backup or regenerate",
		Owner:  "Runtime",
		Steps: []string{
			"Check for context backup",
			"Restore from backup if available",
			"Regenerate context if no backup",
		},
	},

	FailureVersionIncompatible: {
		Action: "Update or configure compatibility",
		Owner:  "User",
		Steps: []string{
			"Check required runtime version",
			"Upgrade runtime if possible",
			"Configure compatibility mode if available",
		},
	},
	FailureRuntimeUpgradeRequired: {
		Action: "Upgrade KDSE runtime",
		Owner:  "User",
		Steps: []string{
			"Check current runtime version",
			"Download and install latest version",
			"Verify installation",
		},
	},
}

// NewFailure creates a new failure with the given code and details
func NewFailure(code FailureCode, source, expected, actual string) *Failure {
	info, ok := FailureCodeInfo[code]
	if !ok {
		info = FailureCodeInfo[FailureBootstrapFailed]
	}

	details := FailureDetails{
		Source:   source,
		Expected: expected,
		Actual:   actual,
	}

	// Generate diagnostic message
	details.Diagnostic = fmt.Sprintf("%s: expected %s, got %s", info.Description, expected, actual)

	recovery, hasRecovery := RecoveryProcedures[code]

	return &Failure{
		Code:        code,
		Name:        info.Name,
		Description: info.Description,
		Details:     details,
		Timestamp:   time.Now().Format(time.RFC3339),
		RequestID:   uuid.New().String(),
		Remediation: func() *Remediation {
			if hasRecovery {
				return &recovery
			}
			return nil
		}(),
	}
}

// Error implements the error interface
func (f *Failure) Error() string {
	return fmt.Sprintf("[%d] %s: %s", f.Code, f.Name, f.Details.Diagnostic)
}

// ToResponse converts the failure to a JSON response format
func (f *Failure) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"failure": map[string]interface{}{
			"code":        f.Code,
			"name":        f.Name,
			"description": f.Description,
			"details": map[string]string{
				"diagnostic": f.Details.Diagnostic,
				"source":     f.Details.Source,
				"expected":   f.Details.Expected,
				"actual":     f.Details.Actual,
			},
			"timestamp":  f.Timestamp,
			"request_id": f.RequestID,
		},
	}

	if f.Remediation != nil {
		response["failure"].(map[string]interface{})["remediation"] = map[string]interface{}{
			"action": f.Remediation.Action,
			"owner":  f.Remediation.Owner,
			"steps":  f.Remediation.Steps,
		}
	}

	return response
}

// FailureRecoveryHandler defines the interface for handling failures
type FailureRecoveryHandler interface {
	HandleFailure(failure *Failure) error
	CanRecover(code FailureCode) bool
	RecoverySteps(code FailureCode) []string
}

// DefaultRecoveryHandler provides default failure recovery logic
type DefaultRecoveryHandler struct{}

// HandleFailure handles a failure based on its code
func (h *DefaultRecoveryHandler) HandleFailure(failure *Failure) error {
	// Log the failure
	fmt.Printf("Failure detected: %s\n", failure.Error())

	// Check if recovery is possible
	if !h.CanRecover(failure.Code) {
		return fmt.Errorf("unrecoverable failure: %s", failure.Name)
	}

	// Return recovery steps
	fmt.Printf("Recovery steps: %v\n", h.RecoverySteps(failure.Code))
	return nil
}

// CanRecover returns true if the failure code is recoverable
func (h *DefaultRecoveryHandler) CanRecover(code FailureCode) bool {
	_, hasRecovery := RecoveryProcedures[code]
	return hasRecovery && code != FailureRuntimeNotFound
}

// RecoverySteps returns the recovery steps for a failure code
func (h *DefaultRecoveryHandler) RecoverySteps(code FailureCode) []string {
	if recovery, ok := RecoveryProcedures[code]; ok {
		return recovery.Steps
	}
	return []string{"Manual intervention required"}
}
