// Package esp implements the KDSE Engineering Session Protocol (ESP).
package esp

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CompletionCriteria represents the conditions required to reach ACTIVE state
// as defined in Section 14 of the ESP specification.
type CompletionCriteria struct {
	// Verification results
	Results []CriterionResult `json:"results"`

	// Summary
	AllPassed      bool    `json:"all_passed"`
	PassedCount    int     `json:"passed_count"`
	FailedCount    int     `json:"failed_count"`
	Confidence     float64 `json:"confidence"`
	CompletionTime string  `json:"completion_time"`
}

// CriterionResult represents the result of a single completion criterion
type CriterionResult struct {
	Condition    string `json:"condition"`
	Method       string `json:"verification_method"`
	Evidence     string `json:"evidence"`
	Status       string `json:"status"` // "PASS" or "FAIL"
	ErrorMessage string `json:"error,omitempty"`
}

// CompletionVerifier verifies completion criteria for ACTIVE state
type CompletionVerifier struct {
	repoPath      string
	workspacePath string
}

// NewCompletionVerifier creates a new completion verifier
func NewCompletionVerifier(repoPath string) *CompletionVerifier {
	return &CompletionVerifier{
		repoPath:      repoPath,
		workspacePath: filepath.Join(repoPath, ".kdse"),
	}
}

// VerifyActiveState verifies all conditions required for ACTIVE state
// per Section 14.1 of the ESP specification.
func (cv *CompletionVerifier) VerifyActiveState(ctx *Context) *CompletionCriteria {
	criteria := &CompletionCriteria{
		Results:       make([]CriterionResult, 0),
		CompletionTime: time.Now().Format(time.RFC3339),
	}

	// Verify each condition
	criteria.Results = append(criteria.Results, cv.verifyRuntimeManifestExists())
	criteria.Results = append(criteria.Results, cv.verifyRuntimeManifestValid())
	criteria.Results = append(criteria.Results, cv.verifyWorkspacePathValid())
	criteria.Results = append(criteria.Results, cv.verifySessionStateCreated(ctx))
	criteria.Results = append(criteria.Results, cv.verifyPhaseStateValid(ctx))
	criteria.Results = append(criteria.Results, cv.verifyContextChecksumComputed(ctx))
	criteria.Results = append(criteria.Results, cv.verifyContextCompleteness(ctx))
	criteria.Results = append(criteria.Results, cv.verifyBootstrapComplete(ctx))

	// Count results
	for _, result := range criteria.Results {
		if result.Status == "PASS" {
			criteria.PassedCount++
		} else {
			criteria.FailedCount++
		}
	}

	// Determine overall pass/fail
	criteria.AllPassed = criteria.FailedCount == 0
	criteria.Confidence = float64(criteria.PassedCount) / float64(len(criteria.Results))

	return criteria
}

// verifyRuntimeManifestExists checks .kdse/runtime.yaml is present
func (cv *CompletionVerifier) verifyRuntimeManifestExists() CriterionResult {
	result := CriterionResult{
		Condition: "Runtime manifest exists",
		Method:    "Filesystem check",
	}

	manifestPath := filepath.Join(cv.workspacePath, "runtime.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		result.Status = "FAIL"
		result.Evidence = ".kdse/runtime.yaml not found"
		result.ErrorMessage = err.Error()
	} else {
		result.Status = "PASS"
		result.Evidence = ".kdse/runtime.yaml present"
	}

	return result
}

// verifyRuntimeManifestValid checks the manifest can be parsed
func (cv *CompletionVerifier) verifyRuntimeManifestValid() CriterionResult {
	result := CriterionResult{
		Condition: "Runtime manifest valid",
		Method:    "YAML parse success",
	}

	manifestPath := filepath.Join(cv.workspacePath, "runtime.yaml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		result.Status = "FAIL"
		result.Evidence = "Cannot read runtime.yaml"
		result.ErrorMessage = err.Error()
		return result
	}

	if !isValidYAML(data) {
		result.Status = "FAIL"
		result.Evidence = "Invalid YAML syntax"
		result.ErrorMessage = "YAML parse errors detected"
	} else {
		result.Status = "PASS"
		result.Evidence = "YAML parsed successfully"
	}

	return result
}

// verifyWorkspacePathValid checks the workspace directory is accessible
func (cv *CompletionVerifier) verifyWorkspacePathValid() CriterionResult {
	result := CriterionResult{
		Condition: "Workspace path valid",
		Method:    "Directory accessible",
	}

	info, err := os.Stat(cv.workspacePath)
	if err != nil {
		result.Status = "FAIL"
		result.Evidence = "Directory inaccessible"
		result.ErrorMessage = err.Error()
		return result
	}

	if !info.IsDir() {
		result.Status = "FAIL"
		result.Evidence = "Path is not a directory"
		result.ErrorMessage = "Workspace path exists but is not a directory"
	} else {
		result.Status = "PASS"
		result.Evidence = "Directory exists and is accessible"
	}

	return result
}

// verifySessionStateCreated checks session ID has been generated
func (cv *CompletionVerifier) verifySessionStateCreated(ctx *Context) CriterionResult {
	result := CriterionResult{
		Condition: "Session state created",
		Method:    "Session ID generated",
	}

	if ctx == nil || ctx.SessionID == "" {
		result.Status = "FAIL"
		result.Evidence = "Session ID not generated"
		result.ErrorMessage = "Session ID is empty"
	} else {
		result.Status = "PASS"
		result.Evidence = fmt.Sprintf("Session ID: %s", ctx.SessionID)
	}

	return result
}

// verifyPhaseStateValid checks phase is from approved list
func (cv *CompletionVerifier) verifyPhaseStateValid(ctx *Context) CriterionResult {
	result := CriterionResult{
		Condition: "Phase state valid",
		Method:    "Phase from approved list",
	}

	validPhases := map[string]bool{
		"Problem":         true,
		"Knowledge":        true,
		"Foundation":       true,
		"Architecture":     true,
		"Implementation":   true,
		"Verification":     true,
		"Documentation":    true,
		"Audit":           true,
	}

	if ctx == nil || ctx.Phase == "" {
		result.Status = "FAIL"
		result.Evidence = "Phase not set"
		result.ErrorMessage = "Phase is empty"
	} else if !validPhases[ctx.Phase] {
		result.Status = "FAIL"
		result.Evidence = fmt.Sprintf("Invalid phase: %s", ctx.Phase)
		result.ErrorMessage = "Phase not in approved list"
	} else {
		result.Status = "PASS"
		result.Evidence = fmt.Sprintf("Phase: %s", ctx.Phase)
	}

	return result
}

// verifyContextChecksumComputed checks checksum has been computed
func (cv *CompletionVerifier) verifyContextChecksumComputed(ctx *Context) CriterionResult {
	result := CriterionResult{
		Condition: "Context checksum computed",
		Method:    "Hash computation",
	}

	if ctx == nil || ctx.Checksum == "" {
		result.Status = "FAIL"
		result.Evidence = "Checksum not computed"
		result.ErrorMessage = "Context checksum is empty"
	} else {
		result.Status = "PASS"
		result.Evidence = fmt.Sprintf("Checksum: %s", ctx.Checksum)
	}

	return result
}

// verifyContextCompleteness checks all required fields are present
func (cv *CompletionVerifier) verifyContextCompleteness(ctx *Context) CriterionResult {
	result := CriterionResult{
		Condition: "Context completeness verified",
		Method:    "All required fields present",
	}

	if ctx == nil {
		result.Status = "FAIL"
		result.Evidence = "Context is nil"
		result.ErrorMessage = "Context object is nil"
		return result
	}

	missing := ctx.Validate()
	if len(missing) > 0 {
		result.Status = "FAIL"
		result.Evidence = fmt.Sprintf("Missing fields: %v", missing)
		result.ErrorMessage = "Completion report indicates missing fields"
	} else {
		result.Status = "PASS"
		result.Evidence = "All required fields present"
	}

	return result
}

// verifyBootstrapComplete checks all bootstrap steps have finished
func (cv *CompletionVerifier) verifyBootstrapComplete(ctx *Context) CriterionResult {
	result := CriterionResult{
		Condition: "Bootstrap complete",
		Method:    "All bootstrap steps finished",
	}

	if ctx == nil || ctx.InitializedAt == "" {
		result.Status = "FAIL"
		result.Evidence = "Bootstrap not started"
		result.ErrorMessage = "Initialization timestamp not set"
	} else {
		result.Status = "PASS"
		result.Evidence = fmt.Sprintf("Initialized at: %s", ctx.InitializedAt)
	}

	return result
}

// Authorization represents authorization granted upon reaching ACTIVE state
type Authorization struct {
	Scope      string   `json:"scope"`
	Constraint string   `json:"constraint"`
	Granted    bool     `json:"granted"`
}

// GetActiveSessionAuthorizations returns authorizations granted at ACTIVE state
// per Section 14.2 of the ESP specification.
func GetActiveSessionAuthorizations() []Authorization {
	return []Authorization{
		{
			Scope:      "Read engineering context",
			Constraint: "Full context access",
			Granted:    true,
		},
		{
			Scope:      "Recommend engineering actions",
			Constraint: "Within session scope, User approval required",
			Granted:    true,
		},
		{
			Scope:      "Generate reports",
			Constraint: "Session reports only, Per report specification",
			Granted:    true,
		},
		{
			Scope:      "Read project artifacts",
			Constraint: "Per permissions, No modification",
			Granted:    true,
		},
		{
			Scope:      "Access external references",
			Constraint: "Referenced sources only, Per reference permissions",
			Granted:    true,
		},
	}
}

// ActiveSessionRequirements defines requirements for maintaining ACTIVE state
// per Section 14.4 of the ESP specification.
type ActiveSessionRequirements struct {
	SessionTokenValid   bool `json:"session_token_valid"`
	ContextIntegrity    bool `json:"context_integrity"`
	PhaseConsistency    bool `json:"phase_consistency"`
	WorkspaceAccessible bool `json:"workspace_accessible"`
	RuntimeOperational  bool `json:"runtime_operational"`
}

// VerifyActiveSessionRequirements verifies requirements for maintaining ACTIVE state
func (cv *CompletionVerifier) VerifyActiveSessionRequirements(ctx *Context) *ActiveSessionRequirements {
	reqs := &ActiveSessionRequirements{
		SessionTokenValid:   ctx != nil && ctx.SessionToken != "",
		ContextIntegrity:    ctx != nil && ctx.Checksum != "",
		PhaseConsistency:    ctx != nil && ctx.Phase != "",
		WorkspaceAccessible: cv.isWorkspaceAccessible(),
		RuntimeOperational:  cv.isRuntimeOperational(),
	}

	return reqs
}

// isWorkspaceAccessible checks if the workspace directory is accessible
func (cv *CompletionVerifier) isWorkspaceAccessible() bool {
	info, err := os.Stat(cv.workspacePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isRuntimeOperational checks if the runtime is operational
func (cv *CompletionVerifier) isRuntimeOperational() bool {
	// Check if required runtime files exist
	requiredFiles := []string{
		"runtime.yaml",
		"session-state.yaml",
	}

	for _, file := range requiredFiles {
		path := filepath.Join(cv.workspacePath, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// FormatCompletionCriteria formats the completion criteria for display
func (c *CompletionCriteria) Format() string {
	result := "╔═══════════════════════════════════════════════════════════════╗\n"
	result += "║           KDSE Engineering Session Completion              ║\n"
	result += "╠═══════════════════════════════════════════════════════════════╣\n"
	result += fmt.Sprintf("║ Passed: %d  |  Failed: %d  |  Confidence: %.0f%%\n", c.PassedCount, c.FailedCount, c.Confidence*100)
	result += "╠═══════════════════════════════════════════════════════════════╣\n"
	result += "║ Verification Results                                      ║\n"

	for _, r := range c.Results {
		statusIcon := "✓"
		if r.Status == "FAIL" {
			statusIcon = "✗"
		}
		result += fmt.Sprintf("║ %s %-20s %s\n", statusIcon, r.Condition, r.Status)
	}

	result += "╠═══════════════════════════════════════════════════════════════╣"
	if c.AllPassed {
		result += "\n║ Status: READY FOR ACTIVE STATE                            ║"
	} else {
		result += "\n║ Status: INCOMPLETE                                        ║"
	}

	result += "\n╚═══════════════════════════════════════════════════════════════╝\n"
	return result
}
