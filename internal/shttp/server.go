// Package shttp provides the KDSE MCP (Model Context Protocol) HTTP server.
// This server exposes KDSE tools as HTTP endpoints with automatic enforcement.
package shttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kdse/runtime/internal/enforcer"
	"github.com/kdse/runtime/internal/guard"
	"github.com/kdse/runtime/internal/mcp"
	"github.com/kdse/runtime/internal/orchestration"
	"github.com/kdse/runtime/internal/runtime"
)

// Server is the KDSE MCP HTTP server
type Server struct {
	repoPath    string
	port        int
	mux         *http.ServeMux
	guard       *guard.SessionGuard
	enforcer    *enforcer.Engine
	orchestrator *orchestration.Manager
	executeTool *mcp.ExecuteTool
	kdseRuntime *runtime.Runtime
	mu          sync.RWMutex
	server      *http.Server
}

// NewServer creates a new KDSE MCP server
func NewServer(repoPath string, port int) *Server {
	s := &Server{
		repoPath:    repoPath,
		port:        port,
		mux:         http.NewServeMux(),
		guard:       guard.NewSessionGuard(repoPath),
		enforcer:    enforcer.NewEngine(repoPath),
		orchestrator: orchestration.NewManager(repoPath),
		executeTool: mcp.NewExecuteTool(repoPath),
		kdseRuntime: runtime.New(repoPath),
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Status and health
	s.mux.HandleFunc("/status", s.handleStatus)
	s.mux.HandleFunc("/health", s.handleHealth)

	// Core KDSE tools (MCP protocol)
	s.mux.HandleFunc("/initialize", s.handleInitialize)
	s.mux.HandleFunc("/execute", s.handleExecute)
	s.mux.HandleFunc("/foundation", s.handleFoundation)
	s.mux.HandleFunc("/knowledge", s.handleKnowledge)
	s.mux.HandleFunc("/audit", s.handleAudit)
	s.mux.HandleFunc("/phase", s.handlePhase)

	// Enforcement tools
	s.mux.HandleFunc("/enforce", s.handleEnforce)
	s.mux.HandleFunc("/compliance", s.handleCompliance)

	// Repository analysis
	s.mux.HandleFunc("/analyze", s.handleAnalyze)
	s.mux.HandleFunc("/collect", s.handleCollect)

	// Report generation
	s.mux.HandleFunc("/report", s.handleReport)
}

// Start begins the HTTP server
func (s *Server) Start() error {
	s.mu.Lock()
	addr := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	s.mu.Unlock()

	log.Printf("[KDSE MCP] Starting server on port %d...", s.port)
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// Response is a standard API response wrapper
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// handleStatus returns server status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	guardResult := s.guard.Check()
	
	resp := Response{
		Success: true,
		Data: map[string]interface{}{
			"kdse": map[string]interface{}{
				"workspace_exists": guardResult.WorkspaceReady,
				"compliant":        guardResult.Valid,
			},
			"guard": map[string]interface{}{
				"initialized":    guardResult.Valid,
				"session_active": guardResult.SessionActive,
				"session_id":     guardResult.SessionID,
				"session_age":    guardResult.SessionAge,
			},
			"orchestration": map[string]interface{}{
				"session_active": s.orchestrator.IsInitialized(),
			},
		},
	}

	s.writeJSON(w, resp)
}

// handleHealth returns health check
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := Response{
		Success: true,
		Data: map[string]interface{}{
			"healthy":   true,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	s.writeJSON(w, resp)
}

// handleInitialize initializes the KDSE workspace and session
func (s *Server) handleInitialize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize workspace
	if err := s.guard.Initialize(); err != nil {
		s.writeError(w, fmt.Sprintf("Initialization failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Initialize orchestration
	state, err := s.orchestrator.Initialize()
	if err != nil {
		s.writeError(w, fmt.Sprintf("Orchestration init failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Initialize runtime
	initResult := s.kdseRuntime.Initialize()

	resp := Response{
		Success: true,
		Data: map[string]interface{}{
			"session_id":     state.SessionID,
			"phase":          state.CurrentPhase,
			"runtime_result": initResult,
			"workspace_path": filepath.Join(s.repoPath, ".kdse"),
		},
	}

	s.writeJSON(w, resp)
}

// handleExecute is the main KDSE execute tool with automatic enforcement
func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req mcp.ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Execute with enforcement
	execResp := s.executeTool.Execute(&req)

	resp := Response{
		Success: !execResp.Blocked,
		Data:    execResp,
	}

	if execResp.Blocked {
		resp.Error = execResp.BlockedReason
		w.WriteHeader(http.StatusForbidden)
	}

	s.writeJSON(w, resp)
}

// handleFoundation handles foundation document operations
func (s *Server) handleFoundation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	foundationPath := filepath.Join(s.repoPath, ".kdse", "foundation")

	switch r.Method {
	case http.MethodGet:
		// Return foundation status
		status := s.getFoundationStatus(foundationPath)
		s.writeJSON(w, Response{Success: true, Data: status})

	case http.MethodPost:
		// Create foundation documents
		var req struct {
			Documents []string `json:"documents"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			created := s.createFoundationDocuments(foundationPath, req.Documents)
			s.writeJSON(w, Response{
				Success: true,
				Data:    map[string]interface{}{"created": created},
			})
		} else {
			// Create all missing documents
			created := s.createAllFoundationDocuments(foundationPath)
			s.writeJSON(w, Response{
				Success: true,
				Data:    map[string]interface{}{"created": created},
			})
		}
	}
}

// getFoundationStatus returns the status of foundation documents
func (s *Server) getFoundationStatus(foundationPath string) map[string]interface{} {
	requiredDocs := enforcer.RequiredFoundationDocs
	status := make(map[string]interface{})
	var missing []string
	var empty []string
	var present []string

	for _, doc := range requiredDocs {
		docPath := filepath.Join(foundationPath, doc)
		info, err := os.Stat(docPath)
		if os.IsNotExist(err) {
			missing = append(missing, doc)
			status[doc] = "missing"
		} else if info.Size() < 100 {
			empty = append(empty, doc)
			status[doc] = "empty"
		} else {
			present = append(present, doc)
			status[doc] = "present"
		}
	}

	return map[string]interface{}{
		"path":     foundationPath,
		"complete": len(missing) == 0 && len(empty) == 0,
		"present":  present,
		"missing":  missing,
		"empty":    empty,
		"status":   status,
	}
}

// createFoundationDocuments creates specific foundation documents
func (s *Server) createFoundationDocuments(foundationPath string, docs []string) []string {
	os.MkdirAll(foundationPath, 0755)
	var created []string

	templates := map[string]string{
		"PROBLEM.md":     getProblemTemplate(),
		"SPEC.md":        getSpecTemplate(),
		"REQUIREMENTS.md": getRequirementsTemplate(),
		"ASSUMPTIONS.md": getAssumptionsTemplate(),
		"CONSTRAINTS.md": getConstraintsTemplate(),
	}

	for _, doc := range docs {
		content, ok := templates[doc]
		if !ok {
			content = fmt.Sprintf("# %s\n\n[Content to be filled]\n", strings.TrimSuffix(doc, ".md"))
		}
		docPath := filepath.Join(foundationPath, doc)
		if _, err := os.Stat(docPath); os.IsNotExist(err) {
			os.WriteFile(docPath, []byte(content), 0644)
			created = append(created, doc)
		}
	}

	return created
}

// createAllFoundationDocuments creates all missing foundation documents
func (s *Server) createAllFoundationDocuments(foundationPath string) []string {
	return s.createFoundationDocuments(foundationPath, enforcer.RequiredFoundationDocs)
}

// handleKnowledge handles knowledge base operations
func (s *Server) handleKnowledge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	knowledgePath := filepath.Join(s.repoPath, ".kdse", "knowledge")

	switch r.Method {
	case http.MethodGet:
		// Return knowledge status
		status := s.getKnowledgeStatus(knowledgePath)
		s.writeJSON(w, Response{Success: true, Data: status})

	case http.MethodPost:
		// Create knowledge categories
		created := s.createKnowledgeCategories(knowledgePath)
		s.writeJSON(w, Response{
			Success: true,
			Data:    map[string]interface{}{"created": created},
		})
	}
}

// getKnowledgeStatus returns the status of knowledge base
func (s *Server) getKnowledgeStatus(knowledgePath string) map[string]interface{} {
	categories := enforcer.RequiredKnowledgeCategories
	status := make(map[string]interface{})
	var missing []string
	var present []string

	for _, cat := range categories {
		catPath := filepath.Join(knowledgePath, cat)
		if _, err := os.Stat(catPath); os.IsNotExist(err) {
			missing = append(missing, cat)
			status[cat] = "missing"
		} else {
			present = append(present, cat)
			// Check if category has content
			hasContent := s.dirHasContent(catPath)
			if hasContent {
				status[cat] = "has_content"
			} else {
				status[cat] = "empty"
			}
		}
	}

	return map[string]interface{}{
		"path":     knowledgePath,
		"complete": len(missing) == 0,
		"present":  present,
		"missing":  missing,
		"status":   status,
	}
}

// createKnowledgeCategories creates knowledge categories
func (s *Server) createKnowledgeCategories(knowledgePath string) []string {
	os.MkdirAll(knowledgePath, 0755)
	var created []string

	for _, cat := range enforcer.RequiredKnowledgeCategories {
		catPath := filepath.Join(knowledgePath, cat)
		if _, err := os.Stat(catPath); os.IsNotExist(err) {
			os.MkdirAll(catPath, 0755)
			readme := fmt.Sprintf("# %s Knowledge\n\n[Collect knowledge artifacts for %s]\n",
				strings.Title(cat), cat)
			os.WriteFile(filepath.Join(catPath, "README.md"), []byte(readme), 0644)
			created = append(created, cat)
		}
	}

	return created
}

// dirHasContent checks if a directory has any files
func (s *Server) dirHasContent(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if s.dirHasContent(filepath.Join(dirPath, entry.Name())) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

// handleAudit runs the compliance audit
func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Run compliance check
	report, err := runtime.ValidateCompliance(s.repoPath)
	if err != nil {
		s.writeError(w, fmt.Sprintf("Audit failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := Response{
		Success: true,
		Data:    report,
	}

	s.writeJSON(w, resp)
}

// handlePhase handles phase transitions
func (s *Server) handlePhase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Return current phase
		state, err := s.orchestrator.Load()
		if err != nil {
			s.writeError(w, "No active session", http.StatusNotFound)
			return
		}

		s.writeJSON(w, Response{
			Success: true,
			Data: map[string]interface{}{
				"current_phase":      state.CurrentPhase,
				"completed_phases":    state.CompletedPhases,
				"next_allowed_phases": state.NextAllowedPhases,
				"confidence":          state.Confidence,
			},
		})

	case http.MethodPost:
		// Transition to new phase
		var req struct {
			Phase    string   `json:"phase"`
			Evidence []string `json:"evidence"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Parse phase
		var phase orchestration.Phase
		switch strings.ToLower(req.Phase) {
		case "problem":
			phase = orchestration.PhaseProblem
		case "knowledge":
			phase = orchestration.PhaseKnowledge
		case "foundation":
			phase = orchestration.PhaseFoundation
		case "audit":
			phase = orchestration.PhaseAudit
		case "assessment":
			phase = orchestration.PhaseAssessment
		case "architecture":
			phase = orchestration.PhaseArchitecture
		case "implementation":
			phase = orchestration.PhaseImplementation
		case "complete":
			phase = orchestration.PhaseComplete
		default:
			s.writeError(w, "Unknown phase", http.StatusBadRequest)
			return
		}

		// Check if transition is allowed
		state, err := s.orchestrator.Load()
		if err != nil {
			s.writeError(w, "No active session", http.StatusNotFound)
			return
		}

		// Check transition validity
		allowed, _ := state.NextAllowedPhases, true // Simplified check
		validTransition := false
		for _, p := range allowed {
			if p == phase {
				validTransition = true
				break
			}
		}

		if !validTransition {
			s.writeJSON(w, Response{
				Success: false,
				Error:   fmt.Sprintf("Invalid transition to %s from current phase", phase),
				Data: map[string]interface{}{
					"current_phase":      state.CurrentPhase,
					"next_allowed_phases": state.NextAllowedPhases,
				},
			})
			return
		}

		// Perform transition
		newState, err := s.orchestrator.TransitionTo(phase, req.Evidence)
		if err != nil {
			s.writeError(w, fmt.Sprintf("Transition failed: %v", err), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, Response{
			Success: true,
			Data:    newState,
		})
	}
}

// handleEnforce runs enforcement check
func (s *Server) handleEnforce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Operation      string `json:"operation"`
		EnforcementLevel string `json:"enforcement_level"`
		AutoCreate     bool   `json:"auto_create"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Default to implementation check
		req.Operation = "implement"
	}

	if req.EnforcementLevel != "" {
		s.enforcer.SetEnforcementLevel(enforcer.EnforcementLevel(req.EnforcementLevel))
	}
	s.enforcer.SetAutoCreate(req.AutoCreate)

	// Run enforcement
	validation := s.enforcer.ValidateOperation(req.Operation)
	report := s.enforcer.GenerateEnforcementReport()

	resp := Response{
		Success: !validation.Allowed,
		Data: map[string]interface{}{
			"validation": validation,
			"report":      report,
			"formatted":   s.enforcer.FormatEnforcementReport(),
		},
	}

	if !validation.Allowed {
		resp.Error = "Enforcement blocked"
		w.WriteHeader(http.StatusForbidden)
	}

	s.writeJSON(w, resp)
}

// handleCompliance runs compliance check
func (s *Server) handleCompliance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := runtime.ValidateCompliance(s.repoPath)
	if err != nil {
		s.writeError(w, fmt.Sprintf("Compliance check failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := Response{
		Success: report.Compliant,
		Data: map[string]interface{}{
			"report":    report,
			"formatted": runtime.FormatComplianceReport(report),
		},
	}

	s.writeJSON(w, resp)
}

// handleAnalyze analyzes the repository
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create repository analysis document
	analysis := s.analyzeRepository()

	resp := Response{
		Success: true,
		Data:    analysis,
	}

	s.writeJSON(w, resp)
}

// analyzeRepository creates a repository analysis document
func (s *Server) analyzeRepository() map[string]interface{} {
	analysisPath := filepath.Join(s.repoPath, ".kdse", "knowledge", "repository-analysis.md")
	os.MkdirAll(filepath.Dir(analysisPath), 0755)

	// Count files by type
	fileTypes := make(map[string]int)
	var totalFiles int
	var totalDirs int

	filepath.Walk(s.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip .kdse and hidden directories
		relPath, _ := filepath.Rel(s.repoPath, path)
		if strings.HasPrefix(relPath, ".kdse") || strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			totalDirs++
		} else {
			totalFiles++
			ext := filepath.Ext(info.Name())
			if ext != "" {
				fileTypes[ext]++
			} else {
				fileTypes["[no ext]"]++
			}
		}
		return nil
	})

	analysis := map[string]interface{}{
		"total_files":   totalFiles,
		"total_dirs":    totalDirs,
		"file_types":    fileTypes,
		"analyzed_at":   time.Now().Format(time.RFC3339),
		"analysis_path": analysisPath,
	}

	// Write analysis document
	content := fmt.Sprintf("# Repository Analysis\n\nGenerated: %s\n\n## Summary\n\n- Total Files: %d\n- Total Directories: %d\n\n## File Types\n\n%s\n", 
		analysis["analyzed_at"],
		totalFiles,
		totalDirs,
		s.formatFileTypes(fileTypes),
	)
	os.WriteFile(analysisPath, []byte(content), 0644)

	return analysis
}

func (s *Server) formatFileTypes(fileTypes map[string]int) string {
	var lines []string
	for ext, count := range fileTypes {
		lines = append(lines, fmt.Sprintf("- %s: %d", ext, count))
	}
	return strings.Join(lines, "\n")
}

// handleCollect collects evidence
func (s *Server) handleCollect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Placeholder for evidence collection
	resp := Response{
		Success: true,
		Data: map[string]interface{}{
			"collected": []string{},
			"message":   "Evidence collection not yet implemented",
		},
	}

	s.writeJSON(w, resp)
}

// handleReport generates a status report
func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guardResult := s.guard.Check()
	
	// Get foundation status
	foundationPath := filepath.Join(s.repoPath, ".kdse", "foundation")
	foundationStatus := s.getFoundationStatus(foundationPath)

	// Get knowledge status
	knowledgePath := filepath.Join(s.repoPath, ".kdse", "knowledge")
	knowledgeStatus := s.getKnowledgeStatus(knowledgePath)

	// Get compliance
	compliance, _ := runtime.ValidateCompliance(s.repoPath)

	resp := Response{
		Success: true,
		Data: map[string]interface{}{
			"session": guardResult,
			"foundation": foundationStatus,
			"knowledge": knowledgeStatus,
			"compliance": compliance,
		},
	}

	s.writeJSON(w, resp)
}

// writeJSON writes a JSON response
func (s *Server) writeJSON(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// writeError writes an error response
func (s *Server) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}

// Template functions
func getProblemTemplate() string {
	return `# Problem Statement

**Created:** ` + time.Now().Format("2006-01-02") + `

## Problem Description

[Describe the problem to be solved in detail]

## Impact

[Describe the impact of not solving this problem]

## Root Causes

1. [Root cause 1]
2. [Root cause 2]

## Success Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Stakeholders

- [Stakeholder 1]
- [Stakeholder 2]
`
}

func getSpecTemplate() string {
	return `# Project Specification

**Created:** ` + time.Now().Format("2006-01-02") + `

## Overview

[High-level description of the project]

## Scope

### In Scope

- Item 1
- Item 2

### Out of Scope

- Item 1
- Item 2

## Deliverables

1. Deliverable 1
2. Deliverable 2

## Quality Standards

[Define quality standards for this project]

## Timeline

[Tentative timeline if known]
`
}

func getRequirementsTemplate() string {
	return `# Functional Requirements

**Created:** ` + time.Now().Format("2006-01-02") + `

## FR-001: [Title]

**Description:** [What the requirement is]

**Priority:** [High/Medium/Low]

**Acceptance Criteria:**
- [ ] Criterion 1
- [ ] Criterion 2

**Dependencies:** [Any dependencies]

---

## FR-002: [Title]

**Description:** [What the requirement is]

**Priority:** [High/Medium/Low]

**Acceptance Criteria:**
- [ ] Criterion 1

**Dependencies:** [Any dependencies]
`
}

func getAssumptionsTemplate() string {
	return `# Key Assumptions

**Created:** ` + time.Now().Format("2006-01-02") + `

## Technical Assumptions

1. [Assumption 1]
2. [Assumption 2]

## Business Assumptions

1. [Assumption 1]
2. [Assumption 2]

## Environment Assumptions

1. [Assumption 1]
2. [Assumption 2]

## Risks from Assumptions

1. [Risk if assumption is wrong]
`
}

func getConstraintsTemplate() string {
	return `# Project Constraints

**Created:** ` + time.Now().Format("2006-01-02") + `

## Technical Constraints

- Constraint 1
- Constraint 2

## Schedule Constraints

- Constraint 1
- Constraint 2

## Resource Constraints

- Constraint 1
- Constraint 2

## Compliance Constraints

- Constraint 1
`
}
