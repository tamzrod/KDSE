// Package guard provides the runtime guard architecture for KDSE.
package guard

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// RuntimeGuard is the entry point for every runtime command.
// It orchestrates guards in deterministic order and stops on first failure.
//
// Guard Execution Order:
//   1. Project Guard - Detect and validate project
//   2. Workspace Guard - Discover and validate workspace
//   3. Session Guard - Detect and validate session
//   4. Lifecycle Guard - Validate lifecycle state
//
// The Runtime Guard itself does not perform validation.
// It delegates all validation to specialized guards.
type RuntimeGuard struct {
	repoPath string

	// Guards - initialized by options
	projectGuard   *ProjectGuard
	workspaceGuard *WorkspaceGuard
	sessionGuard   *SessionValidationGuard
	lifecycleGuard *LifecycleGuard

	mu      sync.RWMutex
	options RuntimeGuardOptions
}

// RuntimeGuardOptions configures the Runtime Guard behavior
type RuntimeGuardOptions struct {
	// AutoInitialize enables automatic initialization on first run
	AutoInitialize bool

	// Logging enables debug logging
	Logging bool

	// SkipSession skips session guard (for status-only operations)
	SkipSession bool

	// SkipLifecycle skips lifecycle guard (for basic operations)
	SkipLifecycle bool
}

// DefaultRuntimeGuardOptions returns the default options
func DefaultRuntimeGuardOptions() RuntimeGuardOptions {
	return RuntimeGuardOptions{
		AutoInitialize: false,
		Logging:        true,
		SkipSession:    false,
		SkipLifecycle: false,
	}
}

// NewRuntimeGuard creates a new Runtime Guard for the given repository path
func NewRuntimeGuard(repoPath string) *RuntimeGuard {
	return &RuntimeGuard{
		repoPath:       repoPath,
		projectGuard:   NewProjectGuard(repoPath),
		workspaceGuard: NewWorkspaceGuard(repoPath),
		sessionGuard:   NewSessionValidationGuard(repoPath),
		lifecycleGuard: NewLifecycleGuard(repoPath),
		options:        DefaultRuntimeGuardOptions(),
	}
}

// NewRuntimeGuardWithOptions creates a new Runtime Guard with custom options
func NewRuntimeGuardWithOptions(repoPath string, opts RuntimeGuardOptions) *RuntimeGuard {
	g := NewRuntimeGuard(repoPath)
	g.options = opts
	return g
}

// SetOptions updates the Runtime Guard options
func (g *RuntimeGuard) SetOptions(opts RuntimeGuardOptions) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.options = opts
}

// Validate executes all guards in deterministic order.
// Returns immediately on first failure.
//
// Execution flow:
//   1. Project Guard → Project detected and valid
//   2. Workspace Guard → Workspace exists and valid
//   3. Session Guard → Session is active and valid
//   4. Lifecycle Guard → Lifecycle is ready for operations
func (g *RuntimeGuard) Validate(ctx context.Context) *RuntimeValidationResult {
	g.mu.Lock()
	opts := g.options
	g.mu.Unlock()

	result := &RuntimeValidationResult{
		Valid:      true,
		FinalState: StateNoProject,
	}

	if opts.Logging {
		log.Printf("[RUNTIME_GUARD] Starting validation sequence...")
	}

	// Step 1: Project Guard
	if opts.Logging {
		log.Printf("[RUNTIME_GUARD] Executing Project Guard...")
	}
	projectResult := g.projectGuard.Validate(ctx)
	result.Project = projectResult

	if !projectResult.Valid {
		if opts.Logging {
			log.Printf("[RUNTIME_GUARD] Project Guard failed: %s", projectResult.Error.Message)
		}
		result.Valid = false
		result.FinalState = StateNoProject
		result.AddError(projectResult.Error)
		return result
	}

	if opts.Logging {
		log.Printf("[RUNTIME_GUARD] Project Guard passed: %s", projectResult.ProjectName)
	}

	// Step 2: Workspace Guard
	if opts.Logging {
		log.Printf("[RUNTIME_GUARD] Executing Workspace Guard...")
	}
	workspaceResult := g.workspaceGuard.Validate(ctx)
	result.Workspace = workspaceResult

	if !workspaceResult.Valid {
		if opts.Logging {
			log.Printf("[RUNTIME_GUARD] Workspace Guard failed: %s", workspaceResult.Error.Message)
		}
		result.Valid = false
		result.FinalState = StateProject
		result.AddError(workspaceResult.Error)
		return result
	}

	if opts.Logging {
		log.Printf("[RUNTIME_GUARD] Workspace Guard passed")
	}

	// Step 3: Session Guard (unless skipped)
	if !opts.SkipSession {
		if opts.Logging {
			log.Printf("[RUNTIME_GUARD] Executing Session Guard...")
		}
		sessionResult := g.sessionGuard.Validate(ctx)
		result.Session = sessionResult

		if !sessionResult.Valid {
			if opts.Logging {
				log.Printf("[RUNTIME_GUARD] Session Guard failed: %s", sessionResult.Error.Message)
			}
			result.Valid = false
			result.FinalState = StateWorkspace
			result.AddError(sessionResult.Error)
			return result
		}

		if opts.Logging {
			log.Printf("[RUNTIME_GUARD] Session Guard passed: %s", sessionResult.SessionID)
		}
	}

	// Step 4: Lifecycle Guard (unless skipped)
	if !opts.SkipLifecycle {
		if opts.Logging {
			log.Printf("[RUNTIME_GUARD] Executing Lifecycle Guard...")
		}
		lifecycleResult := g.lifecycleGuard.Validate(ctx)
		result.Lifecycle = lifecycleResult

		if !lifecycleResult.Valid {
			if opts.Logging {
				log.Printf("[RUNTIME_GUARD] Lifecycle Guard failed: %s", lifecycleResult.Error.Message)
			}
			result.Valid = false
			result.FinalState = StateSession
			result.AddError(lifecycleResult.Error)
			return result
		}

		if opts.Logging {
			log.Printf("[RUNTIME_GUARD] Lifecycle Guard passed")
		}
	}

	// All guards passed
	result.FinalState = StateLifecycleReady

	if opts.Logging {
		log.Printf("[RUNTIME_GUARD] Validation sequence complete: ALL PASSED")
	}

	return result
}

// QuickCheck performs a non-blocking validation without enforcement.
// Returns results without stopping on first failure.
func (g *RuntimeGuard) QuickCheck(ctx context.Context) *RuntimeValidationResult {
	result := &RuntimeValidationResult{
		Valid:      true,
		FinalState: StateNoProject,
	}

	// Project Guard - always runs
	projectResult := g.projectGuard.Validate(ctx)
	result.Project = projectResult
	if !projectResult.Valid {
		result.Valid = false
		result.FinalState = StateNoProject
		result.AddError(projectResult.Error)
		return result
	}

	// Workspace Guard - always runs
	workspaceResult := g.workspaceGuard.Validate(ctx)
	result.Workspace = workspaceResult
	if !workspaceResult.Valid {
		result.Valid = false
		result.FinalState = StateProject
		result.AddError(workspaceResult.Error)
		return result
	}

	// Session Guard
	sessionResult := g.sessionGuard.Validate(ctx)
	result.Session = sessionResult
	if !sessionResult.Valid {
		result.Valid = false
		result.FinalState = StateWorkspace
		result.AddError(sessionResult.Error)
		return result
	}

	// Lifecycle Guard
	lifecycleResult := g.lifecycleGuard.Validate(ctx)
	result.Lifecycle = lifecycleResult
	if !lifecycleResult.Valid {
		result.Valid = false
		result.FinalState = StateSession
		result.AddError(lifecycleResult.Error)
		return result
	}

	result.FinalState = StateLifecycleReady
	return result
}

// EnforceForCommand validates the runtime and returns an error with command context.
// This is the primary entry point for all runtime commands.
func (g *RuntimeGuard) EnforceForCommand(ctx context.Context, commandName string) error {
	result := g.Validate(ctx)

	if result.Valid {
		return nil
	}

	// Get the first error for context
	var firstErr *RuntimeGuardError
	for _, err := range result.Errors {
		firstErr = err
		break
	}

	if firstErr == nil {
		return fmt.Errorf("[RUNTIME_GUARD] Cannot execute '%s': validation failed with no specific error", commandName)
	}

	return fmt.Errorf("[%s:%s] Cannot %s: %s. %s",
		firstErr.GuardType,
		firstErr.Code,
		commandName,
		firstErr.Message,
		firstErr.Hint,
	)
}

// GetCurrentState returns the current runtime state without validation
func (g *RuntimeGuard) GetCurrentState() RuntimeState {
	// Quick check without full validation
	projectValid := g.projectGuard.Exists()

	if !projectValid {
		return StateNoProject
	}

	workspaceValid := g.workspaceGuard.Exists()

	if !workspaceValid {
		return StateProject
	}

	sessionValid := g.sessionGuard.HasActiveSession()

	if !sessionValid {
		return StateWorkspace
	}

	lifecycleValid := g.lifecycleGuard.IsValid()

	if !lifecycleValid {
		return StateSession
	}

	return StateLifecycleReady
}

// FormatResult formats the guard result for display
func FormatGuardResult(result *RuntimeValidationResult) string {
	if result.Valid {
		return fmt.Sprintf("✓ Runtime Guard: ALL PASSED (State: %s)", result.FinalState)
	}

	var output string
	output += fmt.Sprintf("✗ Runtime Guard: FAILED (Final State: %s)\n", result.FinalState)
	output += "\nGuard Errors:\n"

	for _, err := range result.Errors {
		output += fmt.Sprintf("  [%s:%s] %s\n", err.GuardType, err.Code, err.Message)
		output += fmt.Sprintf("    Hint: %s\n", err.Hint)
	}

	return output
}
