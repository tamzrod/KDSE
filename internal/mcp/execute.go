// Package mcp provides the KDSE MCP (Model Context Protocol) tools.
// This module implements the execute tool with automatic enforcement.
package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kdse/runtime/internal/enforcer"
	"github.com/kdse/runtime/internal/guard"
	"github.com/kdse/runtime/internal/orchestration"
	"github.com/kdse/runtime/internal/types"
)

// ExecuteRequest represents a request to execute a KDSE operation
type ExecuteRequest struct {
	Objective       string            `json:"objective"`        // The user's request
	Force           bool              `json:"force"`           // Skip enforcement (dangerous)
	AutoFoundation  bool              `json:"auto_foundation"`  // Auto-create foundation
	AutoKnowledge   bool              `json:"auto_knowledge"`  // Auto-collect knowledge
	EnforcementLevel enforcer.EnforcementLevel `json:"enforcement_level"`
}

// ExecuteResponse represents the response from execute
type ExecuteResponse struct {
	SessionID      string                     `json:"session_id"`
	Phase          types.OrchestrationPhase    `json:"phase"`
	Blocked        bool                       `json:"blocked"`
	BlockedReason  string                     `json:"blocked_reason,omitempty"`
	Violations     []enforcer.EnforcementError `json:"violations,omitempty"`
	WorkOrder      *orchestration.WorkOrder   `json:"work_order,omitempty"`
	AuditWarnings  []string                   `json:"audit_warnings,omitempty"`
	AutoActions    []string                   `json:"auto_actions,omitempty"`
	Report         string                     `json:"report"`
	Success        bool                       `json:"success"`
	NextSteps      []string                   `json:"next_steps"`
}

// ExecuteTool is the main KDSE execute tool with automatic enforcement
type ExecuteTool struct {
	repoPath    string
	guard       *guard.SessionGuard
	enforcer    *enforcer.Engine
	orchestrator *orchestration.Manager
}

// NewExecuteTool creates a new execute tool
func NewExecuteTool(repoPath string) *ExecuteTool {
	return &ExecuteTool{
		repoPath:    repoPath,
		guard:       guard.NewSessionGuard(repoPath),
		enforcer:    enforcer.NewEngine(repoPath),
		orchestrator: orchestration.NewManager(repoPath),
	}
}

// Execute processes a user request with automatic enforcement
func (t *ExecuteTool) Execute(req *ExecuteRequest) *ExecuteResponse {
	startTime := time.Now()
	response := &ExecuteResponse{
		Success:       false,
		Blocked:       false,
		Violations:    []enforcer.EnforcementError{},
		AuditWarnings: []string{},
		AutoActions:   []string{},
		NextSteps:     []string{},
	}

	// Get session ID if active
	if sessionID, err := t.guard.GetSessionID(); err == nil {
		response.SessionID = sessionID
	}

	// Step 1: Check guard initialization
	if err := t.guard.EnforceInitialized(); err != nil {
		response.Blocked = true
		response.BlockedReason = fmt.Sprintf("Session not initialized: %v", err)
		response.Report = t.formatBlockedReport(response, startTime)
		return response
	}

	// Load orchestrator state
	state, err := t.orchestrator.Load()
	if err != nil {
		// Initialize if not exists
		state, err = t.orchestrator.Initialize()
		if err != nil {
			response.Blocked = true
			response.BlockedReason = fmt.Sprintf("Failed to initialize session: %v", err)
			response.Report = t.formatBlockedReport(response, startTime)
			return response
		}
	}
	response.Phase = state.CurrentPhase

	// Step 2: Detect if this is a coding/implementation request
	isImplementationRequest := t.detectImplementationRequest(req.Objective)

	// Step 3: If implementation request, enforce KDSE principles
	if isImplementationRequest && !req.Force {
		enforcementResult := t.enforceWithResponse(req)
		
		if enforcementResult.Blocked {
			response.Blocked = true
			response.BlockedReason = enforcementResult.BlockedReason
			response.Violations = enforcementResult.Violations
			response.AuditWarnings = enforcementResult.Warnings
			response.AutoActions = enforcementResult.AutoActions
			response.Report = t.formatBlockedReport(response, startTime)
			return response
		}

		// Record any warnings
		response.Violations = enforcementResult.Violations
		response.AuditWarnings = enforcementResult.Warnings
		response.AutoActions = enforcementResult.AutoActions
	}

	// Step 4: Generate appropriate work order based on phase
	workOrder := t.generateWorkOrderForObjective(req.Objective, state.CurrentPhase)
	response.WorkOrder = workOrder

	// Step 5: Check for phase-appropriate actions
	phaseWarnings := t.checkPhaseAppropriateness(state.CurrentPhase, req.Objective)
	response.AuditWarnings = append(response.AuditWarnings, phaseWarnings...)

	// Step 6: Generate next steps
	response.NextSteps = t.suggestNextSteps(state.CurrentPhase, req.Objective)

	// Step 7: Determine success and format report
	response.Success = !response.Blocked
	response.Report = t.formatSuccessReport(response, startTime)

	return response
}

// enforceWithResponse runs enforcement and collects all feedback
func (t *ExecuteTool) enforceWithResponse(req *ExecuteRequest) *enforcementResult {
	result := &enforcementResult{
		Blocked:      false,
		BlockedReason: "",
		Violations:   []enforcer.EnforcementError{},
		Warnings:     []string{},
		AutoActions:  []string{},
	}

	// Configure enforcer
	if req.EnforcementLevel != "" {
		t.enforcer.SetEnforcementLevel(req.EnforcementLevel)
	}
	if req.AutoFoundation {
		t.enforcer.SetAutoCreate(true)
	}
	if req.AutoKnowledge {
		t.enforcer.SetAutoCreate(true)
	}

	// Run enforcement
	err := t.enforcer.EnforceImplementation()
	if err != nil {
		// Check if it's an enforcement error
		if enfErr, ok := err.(*enforcer.EnforcementError); ok {
			result.Blocked = enfErr.Blocked
			result.BlockedReason = enfErr.Message
			if enfErr.Blocked {
				result.Violations = t.enforcer.GetViolations()
				result.AutoActions = t.collectAutoActions()
				return result
			}
		}
		
		// Non-blocking error - add as warning
		result.Warnings = append(result.Warnings, err.Error())
	}

	// Check for any remaining violations
	if t.enforcer.HasViolations() {
		result.Violations = t.enforcer.GetViolations()
		for _, v := range result.Violations {
			if !v.Blocked {
				result.Warnings = append(result.Warnings, fmt.Sprintf("[%s] %s", v.Code, v.Message))
			}
		}
	}

	// Record auto-actions taken
	result.AutoActions = t.collectAutoActions()

	return result
}

type enforcementResult struct {
	Blocked      bool
	BlockedReason string
	Violations   []enforcer.EnforcementError
	Warnings     []string
	AutoActions  []string
}

// collectAutoActions collects actions that were auto-performed
func (t *ExecuteTool) collectAutoActions() []string {
	var actions []string
	
	kdsePath := filepath.Join(t.repoPath, ".kdse")
	foundationPath := filepath.Join(kdsePath, "foundation")
	
	// Check if foundation docs were created
	for _, doc := range enforcer.RequiredFoundationDocs {
		docPath := filepath.Join(foundationPath, doc)
		if _, err := os.Stat(docPath); err == nil {
			actions = append(actions, fmt.Sprintf("Created: %s", doc))
		}
	}
	
	return actions
}

// detectImplementationRequest determines if the objective is implementation-focused
func (t *ExecuteTool) detectImplementationRequest(objective string) bool {
	objective = strings.ToLower(objective)
	
	implementationKeywords := []string{
		"build", "create", "implement", "develop", "write", "code",
		"make", "add", "create function", "create class", "create file",
		"write test", "add feature", "build system", "develop api",
	}
	
	for _, keyword := range implementationKeywords {
		if strings.Contains(objective, keyword) {
			return true
		}
	}
	
	return false
}

// generateWorkOrderForObjective creates appropriate work order based on objective
func (t *ExecuteTool) generateWorkOrderForObjective(objective string, phase types.OrchestrationPhase) *orchestration.WorkOrder {
	// For now, delegate to orchestrator's GenerateWorkOrder
	// In a full implementation, this would be more sophisticated
	return t.orchestrator.GenerateWorkOrder(phase)
}

// checkPhaseAppropriateness checks if the requested action matches current phase
func (t *ExecuteTool) checkPhaseAppropriateness(phase types.OrchestrationPhase, objective string) []string {
	var warnings []string
	
	isImpl := t.detectImplementationRequest(objective)
	
	switch phase {
	case types.PhaseProblem:
		if isImpl {
			warnings = append(warnings, "WARNING: Implementation requested during Problem phase. Consider completing problem definition first.")
		}
	case types.PhaseKnowledge:
		if isImpl {
			warnings = append(warnings, "WARNING: Implementation requested during Knowledge phase. Knowledge base should be built first.")
		}
	case types.PhaseFoundation:
		if isImpl {
			warnings = append(warnings, "WARNING: Implementation requested during Foundation phase. Foundation documents should be complete.")
		}
	case types.PhaseArchitecture:
		if isImpl {
			warnings = append(warnings, "WARNING: Implementation requested during Architecture phase. Architecture should be finalized first.")
		}
	}
	
	return warnings
}

// suggestNextSteps suggests what to do next based on current state
func (t *ExecuteTool) suggestNextSteps(phase types.OrchestrationPhase, objective string) []string {
	var steps []string
	
	// Load actual workspace state to provide accurate suggestions
	kdsePath := filepath.Join(t.repoPath, ".kdse")
	
	switch phase {
	case types.PhaseIdle, types.PhaseProblem:
		steps = append(steps, "1. Complete problem definition in .kdse/foundation/PROBLEM.md")
		steps = append(steps, "2. Analyze repository structure with: kdse analyze --repository")
		steps = append(steps, "3. Collect domain knowledge with: kdse knowledge collect")
	case types.PhaseKnowledge:
		steps = append(steps, "1. Complete knowledge collection in .kdse/knowledge/")
		steps = append(steps, "2. Create/verify SPEC.md with: kdse spec create")
		steps = append(steps, "3. Define requirements with: kdse requirements create")
	case types.PhaseFoundation:
		// Check which docs are missing
		foundationPath := filepath.Join(kdsePath, "foundation")
		for _, doc := range enforcer.RequiredFoundationDocs {
			docPath := filepath.Join(foundationPath, doc)
			if _, err := os.Stat(docPath); os.IsNotExist(err) {
				steps = append(steps, fmt.Sprintf("4. Create %s", doc))
			}
		}
		steps = append(steps, "5. Run audit with: kdse audit")
	case types.PhaseAudit:
		steps = append(steps, "1. Review audit findings")
		steps = append(steps, "2. Address any critical/high issues")
		steps = append(steps, "3. Proceed to architecture with: kdse phase architecture")
	case types.PhaseArchitecture:
		steps = append(steps, "1. Design system architecture in .kdse/foundation/ARCHITECTURE.md")
		steps = append(steps, "2. Review architecture decisions")
		steps = append(steps, "3. Proceed to implementation with: kdse phase implementation")
	case types.PhaseImplementation:
		steps = append(steps, "1. Implement according to ARCHITECTURE.md")
		steps = append(steps, "2. Run tests with: kdse test")
		steps = append(steps, "3. Verify with: kdse verify")
	}
	
	// Always add KDSE principles reminder
	steps = append(steps, "")
	steps = append(steps, "Remember: KDSE requires Knowledge → Foundation → Architecture → Implementation")
	
	return steps
}

// formatBlockedReport formats the response when blocked
func (t *ExecuteTool) formatBlockedReport(resp *ExecuteResponse, startTime time.Time) string {
	var sb strings.Builder
	
	sb.WriteString("\n╔════════════════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║              KDSE EXECUTE - BLOCKED                                 ║\n")
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ Objective: %s\n", truncate(resp.BlockedReason, 50)))
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	sb.WriteString("║ ██ KDSE PRINCIPLES VIOLATED                                       ║\n")
	sb.WriteString("║                                                               ║\n")
	sb.WriteString("║ Implementation is BLOCKED because:                              ║\n")
	sb.WriteString(fmt.Sprintf("║   %s\n", truncate(resp.BlockedReason, 62)))
	sb.WriteString("║                                                               ║\n")
	
	if len(resp.Violations) > 0 {
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
		sb.WriteString("║ VIOLATIONS                                                     ║\n")
		for _, v := range resp.Violations {
			icon := "🔴"
			if !v.Blocked {
				icon = "🟡"
			}
			sb.WriteString(fmt.Sprintf("║ %s [%s] %s\n", icon, v.Code, v.Severity))
			sb.WriteString(fmt.Sprintf("║   %s\n", truncate(v.Message, 62)))
			if v.RequiredPhase != "" {
				sb.WriteString(fmt.Sprintf("║   → Complete %s phase first\n", v.RequiredPhase))
			}
		}
	}
	
	if len(resp.AutoActions) > 0 {
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
		sb.WriteString("║ AUTO-ACTIONS (auto-created):                                     ║\n")
		for _, action := range resp.AutoActions {
			sb.WriteString(fmt.Sprintf("║   ✓ %s\n", action))
		}
	}
	
	if len(resp.AuditWarnings) > 0 {
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
		sb.WriteString("║ AUDIT WARNINGS                                                  ║\n")
		for _, w := range resp.AuditWarnings {
			sb.WriteString(fmt.Sprintf("║   ⚠ %s\n", truncate(w, 62)))
		}
	}
	
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	sb.WriteString("║ REQUIRED PHASES                                                 ║\n")
	sb.WriteString("║ 1. Problem → 2. Knowledge → 3. Foundation → 4. Architecture → 5. Implement ║\n")
	
	if len(resp.NextSteps) > 0 {
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
		sb.WriteString("║ NEXT STEPS                                                      ║\n")
		for _, step := range resp.NextSteps {
			sb.WriteString(fmt.Sprintf("║   %s\n", truncate(step, 62)))
		}
	}
	
	sb.WriteString("╚════════════════════════════════════════════════════════════════════════╝\n")
	
	return sb.String()
}

// formatSuccessReport formats the response when successful
func (t *ExecuteTool) formatSuccessReport(resp *ExecuteResponse, startTime time.Time) string {
	var sb strings.Builder
	duration := time.Since(startTime)
	
	sb.WriteString("\n╔════════════════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║              KDSE EXECUTE - READY                                  ║\n")
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ Session: %s\n", truncate(resp.SessionID, 45)))
	sb.WriteString(fmt.Sprintf("║ Phase: %s\n", resp.Phase))
	sb.WriteString(fmt.Sprintf("║ Duration: %v\n", duration))
	sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	
	if len(resp.Violations) > 0 {
		sb.WriteString("║ WARNINGS (non-blocking):                                         ║\n")
		for _, v := range resp.Violations {
			if !v.Blocked {
				sb.WriteString(fmt.Sprintf("║   ⚠ [%s] %s\n", v.Code, truncate(v.Message, 50)))
			}
		}
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	}
	
	if len(resp.AutoActions) > 0 {
		sb.WriteString("║ AUTO-ACTIONS COMPLETED                                          ║\n")
		for _, action := range resp.AutoActions {
			sb.WriteString(fmt.Sprintf("║   ✓ %s\n", action))
		}
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	}
	
	if resp.WorkOrder != nil {
		sb.WriteString("║ WORK ORDER                                                      ║\n")
		sb.WriteString(fmt.Sprintf("║   Phase: %s\n", resp.WorkOrder.Phase))
		sb.WriteString(fmt.Sprintf("║   %s\n", truncate(resp.WorkOrder.PhaseDescription, 62)))
		
		if len(resp.WorkOrder.BlockedActions) > 0 {
			sb.WriteString("║   BLOCKED ACTIONS:                                               ║\n")
			for _, blocked := range resp.WorkOrder.BlockedActions {
				sb.WriteString(fmt.Sprintf("║     ✗ %s\n", truncate(blocked, 55)))
			}
		}
	}
	
	if len(resp.NextSteps) > 0 {
		sb.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
		sb.WriteString("║ NEXT STEPS                                                      ║\n")
		for _, step := range resp.NextSteps {
			sb.WriteString(fmt.Sprintf("║   %s\n", truncate(step, 62)))
		}
	}
	
	sb.WriteString("╚════════════════════════════════════════════════════════════════════════╝\n")
	
	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ToJSON returns the response as JSON
func (r *ExecuteResponse) ToJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}
