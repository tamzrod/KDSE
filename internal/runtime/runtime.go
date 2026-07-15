// Package runtime implements evidence-driven KDSE runtime management.
// Every operation follows: Execute → Verify → Report
package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VerificationResult represents the verification status of an artifact
type VerificationResult struct {
	Artifact   string `json:"artifact"`
	Path       string `json:"path"`
	Status     string `json:"status"` // "PASS" or "FAIL"
	Evidence   string `json:"evidence,omitempty"`
	Error      string `json:"error,omitempty"`
	Timestamp  string `json:"timestamp"`
}

// RuntimeManifest defines the complete runtime structure
type RuntimeManifest struct {
	Version     string                 `json:"version"`
	CreatedAt   string                 `json:"created_at"`
	Directories []ManifestDirectory    `json:"directories"`
	Files       []ManifestFile         `json:"files"`
	Invariants  []RuntimeInvariant     `json:"invariants"`
}

// ManifestDirectory describes a required directory
type ManifestDirectory struct {
	Path     string `json:"path"`
	Required bool   `json:"required"`
	Purpose  string `json:"purpose"`
}

// ManifestFile describes a required file
type ManifestFile struct {
	Path     string `json:"path"`
	Required bool   `json:"required"`
	Purpose  string `json:"purpose"`
	Template string `json:"template,omitempty"`
}

// RuntimeInvariant defines a phase transition requirement
type RuntimeInvariant struct {
	Phase       string   `json:"phase"`
	Requires    []string `json:"requires"`
	Description string   `json:"description"`
}

// InitializeResult is the result of runtime initialization
type InitializeResult struct {
	Success       bool                 `json:"success"`
	WorkspacePath string               `json:"workspace_path"`
	Confidence    float64              `json:"confidence"`
	Verification  []VerificationResult `json:"verification"`
	Evidence      []string             `json:"evidence"`
	Errors        []string             `json:"errors,omitempty"`
	Timestamp     string               `json:"timestamp"`
}

// VerificationReport is the result of runtime verification
type VerificationReport struct {
	Success      bool                 `json:"success"`
	Confidence   float64              `json:"confidence"`
	Components   []VerificationResult `json:"components"`
	Missing      []string             `json:"missing,omitempty"`
	Failed       []string             `json:"failed,omitempty"`
	Timestamp    string               `json:"timestamp"`
}

// Standard runtime directories
const (
	DirRuntime        = "runtime"
	DirFoundation     = "foundation"
	DirKnowledge      = "knowledge"
	DirLaboratory     = "laboratory"
	DirEvidence       = "evidence"
	DirReferences     = "references"
	DirTraceability   = "traceability"
	DirReports        = "reports"
	DirConfig         = "config"
	DirState          = "state"
	DirArtifacts      = "artifacts"
	DirSessions       = "sessions"
	DirNormalized     = "normalized"
	DirCache          = "cache"
	DirSomeday       = "someday"
)

// Standard runtime files
const (
	FileManifest        = "manifest.yaml"
	FileSessionState    = "session-state.yaml"
	FileRuntimeConfig   = "runtime.yaml"
	FileKnowledgeIndex  = "knowledge-index.yaml"
	FileArtifactIndex   = "artifact-index.yaml"
)

// Runtime defines the evidence-driven runtime manager
type Runtime struct {
	repoPath  string
	kdsePath  string
	manifest  *RuntimeManifest
	verified  bool
	invariant *InvariantEngine
}

// New creates a new Runtime for the given repository path
func New(repoPath string) *Runtime {
	return &Runtime{
		repoPath:  repoPath,
		kdsePath:  filepath.Join(repoPath, ".kdse"),
		manifest:  DefaultManifest(),
		verified:  false,
		invariant: NewInvariantEngine(),
	}
}

// Initialize creates a full operational KDSE runtime
// Returns evidence of what was created and verified
func (r *Runtime) Initialize() *InitializeResult {
	result := &InitializeResult{
		WorkspacePath: r.kdsePath,
		Timestamp:     time.Now().Format(time.RFC3339),
		Verification:  []VerificationResult{},
		Evidence:      []string{},
	}

	// Execute: Create all directories
	r.executeDirectoryCreation(result)

	// Execute: Create all required files
	r.executeFileCreation(result)

	// Verify: Check every artifact
	r.verifyAllArtifacts(result)

	// Calculate confidence based on verification
	result.Confidence = r.calculateConfidence(result.Verification)

	// Determine success
	result.Success = r.determineSuccess(result)

	return result
}

// executeDirectoryCreation creates all required directories
func (r *Runtime) executeDirectoryCreation(result *InitializeResult) {
	directories := []string{
		DirRuntime,
		DirFoundation,
		DirKnowledge,
		DirLaboratory,
		DirEvidence,
		DirReferences,
		DirTraceability,
		DirReports,
		DirConfig,
		DirState,
		DirArtifacts,
		DirSessions,
		DirNormalized,
		DirCache,
		DirSomeday,
	}

	for _, dir := range directories {
		path := filepath.Join(r.kdsePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  dir,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create directory %s: %v", dir, err))
		}
	}

	// Create someday subdirectories
	somedaySubdirs := []string{"ideas", "archived", "promoted"}
	for _, subdir := range somedaySubdirs {
		path := filepath.Join(r.kdsePath, DirSomeday, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  DirSomeday + "/" + subdir,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create directory %s/%s: %v", DirSomeday, subdir, err))
		}
	}
}

// executeFileCreation creates all required runtime files
func (r *Runtime) executeFileCreation(result *InitializeResult) {
	files := map[string]string{
		FileManifest:       r.generateManifestContent(),
		FileSessionState:   r.generateSessionStateContent(),
		FileRuntimeConfig:  r.generateRuntimeConfigContent(),
		FileKnowledgeIndex: r.generateKnowledgeIndexContent(),
		FileArtifactIndex:  r.generateArtifactIndexContent(),
	}

	for filename, content := range files {
		path := filepath.Join(r.kdsePath, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			result.Verification = append(result.Verification, VerificationResult{
				Artifact:  filename,
				Path:      path,
				Status:    "FAIL",
				Error:     err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create file %s: %v", filename, err))
		}
	}

	// Create someday manifest
	somedayManifestPath := filepath.Join(r.kdsePath, DirSomeday, "someday.yaml")
	somedayContent := r.generateSomedayManifestContent()
	if err := os.WriteFile(somedayManifestPath, []byte(somedayContent), 0644); err != nil {
		result.Verification = append(result.Verification, VerificationResult{
			Artifact:  DirSomeday + "/someday.yaml",
			Path:      somedayManifestPath,
			Status:    "FAIL",
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create file %s: %v", DirSomeday+"/someday.yaml", err))
	}
}

// verifyAllArtifacts verifies every created artifact
func (r *Runtime) verifyAllArtifacts(result *InitializeResult) {
	// Verify directories
	directories := []string{
		DirRuntime, DirFoundation, DirKnowledge, DirLaboratory,
		DirEvidence, DirReferences, DirTraceability, DirReports,
		DirConfig, DirState, DirArtifacts, DirSessions,
		DirNormalized, DirCache, DirSomeday,
	}

	for _, dir := range directories {
		path := filepath.Join(r.kdsePath, dir)
		result.Verification = append(result.Verification, r.verifyDirectory(path, dir))
	}

	// Verify someday subdirectories
	somedaySubdirs := []string{"ideas", "archived", "promoted"}
	for _, subdir := range somedaySubdirs {
		path := filepath.Join(r.kdsePath, DirSomeday, subdir)
		result.Verification = append(result.Verification, r.verifyDirectory(path, DirSomeday+"/"+subdir))
	}

	// Verify files
	files := []string{
		FileManifest, FileSessionState, FileRuntimeConfig,
		FileKnowledgeIndex, FileArtifactIndex,
	}

	for _, file := range files {
		path := filepath.Join(r.kdsePath, file)
		result.Verification = append(result.Verification, r.verifyFile(path, file))
	}

	// Verify someday manifest
	somedayManifestPath := filepath.Join(r.kdsePath, DirSomeday, "someday.yaml")
	result.Verification = append(result.Verification, r.verifyFile(somedayManifestPath, DirSomeday+"/someday.yaml"))

	// Add evidence for successful verifications
	for _, v := range result.Verification {
		if v.Status == "PASS" {
			result.Evidence = append(result.Evidence, fmt.Sprintf("%s: %s", v.Artifact, v.Path))
		}
	}
}

// verifyDirectory checks if a directory exists and is accessible
func (r *Runtime) verifyDirectory(path, name string) VerificationResult {
	info, err := os.Stat(path)
	if err != nil {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     fmt.Sprintf("Directory does not exist: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	if !info.IsDir() {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     "Path exists but is not a directory",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return VerificationResult{
		Artifact:  name,
		Path:      path,
		Status:    "PASS",
		Evidence:  fmt.Sprintf("Directory exists, mode: %o", info.Mode()),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// verifyFile checks if a file exists and is readable
func (r *Runtime) verifyFile(path, name string) VerificationResult {
	info, err := os.Stat(path)
	if err != nil {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     fmt.Sprintf("File does not exist: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	if info.IsDir() {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     "Path exists but is a directory, not a file",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	// Verify file is readable
	content, err := os.ReadFile(path)
	if err != nil {
		return VerificationResult{
			Artifact:  name,
			Path:      path,
			Status:    "FAIL",
			Error:     fmt.Sprintf("File exists but is not readable: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return VerificationResult{
		Artifact:  name,
		Path:      path,
		Status:    "PASS",
		Evidence:  fmt.Sprintf("File readable, size: %d bytes", len(content)),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// calculateConfidence calculates verification confidence
func (r *Runtime) calculateConfidence(verifications []VerificationResult) float64 {
	if len(verifications) == 0 {
		return 0.0
	}

	passed := 0
	for _, v := range verifications {
		if v.Status == "PASS" {
			passed++
		}
	}

	return float64(passed) / float64(len(verifications))
}

// determineSuccess checks if all required artifacts passed
func (r *Runtime) determineSuccess(result *InitializeResult) bool {
	for _, v := range result.Verification {
		if v.Status == "FAIL" {
			return false
		}
	}
	return len(result.Errors) == 0
}

// Verify performs a complete runtime verification
func (r *Runtime) Verify() *VerificationReport {
	report := &VerificationReport{
		Components: []VerificationResult{},
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Check workspace
	report.Components = append(report.Components, r.verifyWorkspace())

	// Check runtime directory
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirRuntime), DirRuntime))

	// Check foundation
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirFoundation), DirFoundation))

	// Check knowledge
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirKnowledge), DirKnowledge))

	// Check laboratory
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirLaboratory), DirLaboratory))

	// Check configuration
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirConfig), DirConfig))

	// Check state
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirState), DirState))

	// Check someday
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirSomeday), DirSomeday))
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirSomeday, "ideas"), DirSomeday+"/ideas"))
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirSomeday, "archived"), DirSomeday+"/archived"))
	report.Components = append(report.Components, r.verifyDirectory(filepath.Join(r.kdsePath, DirSomeday, "promoted"), DirSomeday+"/promoted"))
	report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, DirSomeday, "someday.yaml"), DirSomeday+"/someday.yaml"))

	// Check manifest
	report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, FileManifest), FileManifest))

	// Check session state
	report.Components = append(report.Components, r.verifyFile(filepath.Join(r.kdsePath, FileSessionState), FileSessionState))

	// Calculate confidence
	report.Confidence = r.calculateConfidence(report.Components)

	// Determine overall status
	report.Success = r.determineVerificationSuccess(report)

	return report
}

// verifyWorkspace checks if the .kdse workspace exists
func (r *Runtime) verifyWorkspace() VerificationResult {
	info, err := os.Stat(r.kdsePath)
	if err != nil {
		return VerificationResult{
			Artifact:  "Workspace",
			Path:      r.kdsePath,
			Status:    "FAIL",
			Error:     "Workspace does not exist",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	if !info.IsDir() {
		return VerificationResult{
			Artifact:  "Workspace",
			Path:      r.kdsePath,
			Status:    "FAIL",
			Error:     "Workspace path exists but is not a directory",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return VerificationResult{
		Artifact:  "Workspace",
		Path:      r.kdsePath,
		Status:    "PASS",
		Evidence:  "Workspace directory exists",
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// determineVerificationSuccess checks if all components passed
func (r *Runtime) determineVerificationSuccess(report *VerificationReport) bool {
	for _, c := range report.Components {
		if c.Status == "FAIL" {
			report.Failed = append(report.Failed, c.Artifact)
		}
	}
	return len(report.Failed) == 0
}

// CheckInvariant verifies if a phase transition is allowed
func (r *Runtime) CheckInvariant(phase, required string) (bool, string) {
	return r.invariant.Check(phase, required)
}

// DefaultManifest returns the default runtime manifest
func DefaultManifest() *RuntimeManifest {
	return &RuntimeManifest{
		Version:   "1.0.0",
		CreatedAt: time.Now().Format(time.RFC3339),
		Directories: []ManifestDirectory{
			{Path: DirRuntime, Required: true, Purpose: "Runtime execution state and logs"},
			{Path: DirFoundation, Required: true, Purpose: "Project foundation documents"},
			{Path: DirKnowledge, Required: true, Purpose: "Collected engineering knowledge"},
			{Path: DirLaboratory, Required: true, Purpose: "Testing and experimentation"},
			{Path: DirEvidence, Required: true, Purpose: "Evidence artifacts"},
			{Path: DirReferences, Required: true, Purpose: "Reference materials"},
			{Path: DirTraceability, Required: true, Purpose: "Requirement traceability"},
			{Path: DirReports, Required: true, Purpose: "Generated reports"},
			{Path: DirConfig, Required: true, Purpose: "Runtime configuration"},
			{Path: DirState, Required: true, Purpose: "Session and runtime state"},
			{Path: DirArtifacts, Required: true, Purpose: "Artifact inventory"},
			{Path: DirSessions, Required: true, Purpose: "Session history"},
			{Path: DirNormalized, Required: true, Purpose: "Normalized documentation"},
			{Path: DirCache, Required: true, Purpose: "Cached computations"},
			{Path: DirSomeday, Required: true, Purpose: "Someday/Maybe knowledge repository"},
		},
		Files: []ManifestFile{
			{Path: FileManifest, Required: true, Purpose: "Runtime manifest definition"},
			{Path: FileSessionState, Required: true, Purpose: "Current session state"},
			{Path: FileRuntimeConfig, Required: true, Purpose: "Runtime configuration"},
			{Path: FileKnowledgeIndex, Required: true, Purpose: "Knowledge artifact index"},
			{Path: FileArtifactIndex, Required: true, Purpose: "Artifact inventory"},
		},
		Invariants: DefaultInvariants(),
	}
}

// DefaultInvariants returns the default runtime invariants
func DefaultInvariants() []RuntimeInvariant {
	return []RuntimeInvariant{
		{Phase: "Problem", Requires: []string{"Runtime initialized"}, Description: "Runtime must be initialized before problem phase"},
		{Phase: "Foundation", Requires: []string{"Runtime initialized"}, Description: "Foundation requires initialized runtime"},
		{Phase: "Knowledge", Requires: []string{"Foundation exists"}, Description: "Knowledge collection requires foundation"},
		{Phase: "Architecture", Requires: []string{"Knowledge collected"}, Description: "Architecture requires knowledge"},
		{Phase: "Implementation", Requires: []string{"Architecture approved"}, Description: "Implementation requires approved architecture"},
		{Phase: "Verification", Requires: []string{"Implementation complete"}, Description: "Verification requires implementation"},
		{Phase: "Documentation", Requires: []string{"Verification complete"}, Description: "Documentation requires verification"},
		{Phase: "Audit", Requires: []string{"All phases complete"}, Description: "Audit requires all phases complete"},
	}
}

// File generation templates

func (r *Runtime) generateManifestContent() string {
	m := r.manifest
	m.CreatedAt = time.Now().Format(time.RFC3339)
	data, _ := json.MarshalIndent(m, "", "  ")
	return string(data)
}

func (r *Runtime) generateSessionStateContent() string {
	state := map[string]interface{}{
		"version":       "1.0.0",
		"session_id":    fmt.Sprintf("KDSE-SESSION-%s", time.Now().Format("20060102-150405")),
		"status":        "Initialized",
		"phase":         "Problem",
		"confidence":    0.0,
		"evidence":      []string{},
		"created_at":    time.Now().Format(time.RFC3339),
		"last_verified": nil,
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	return string(data)
}

func (r *Runtime) generateRuntimeConfigContent() string {
	config := map[string]interface{}{
		"version":               "1.0.0",
		"runtime":               "evidence-driven",
		"strict_mode":           true,
		"confidence_threshold":  0.7,
		"evidence_threshold":    0.6,
		"max_cycles":            100,
		"auto_verify":           true,
		"enforce_invariants":    true,
		"created_at":            time.Now().Format(time.RFC3339),
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	return string(data)
}

func (r *Runtime) generateKnowledgeIndexContent() string {
	index := map[string]interface{}{
		"version":       "1.0.0",
		"last_updated":  time.Now().Format(time.RFC3339),
		"artifacts":     []map[string]interface{}{},
		"categories": map[string]int{
			"architecture":   0,
			"design":         0,
			"implementation": 0,
			"testing":        0,
			"documentation":  0,
		},
		"total_count": 0,
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	return string(data)
}

func (r *Runtime) generateArtifactIndexContent() string {
	index := map[string]interface{}{
		"version":      "1.0.0",
		"last_updated": time.Now().Format(time.RFC3339),
		"artifacts":    []map[string]interface{}{},
		"categories": map[string]int{
			"foundation":    0,
			"evidence":      0,
			"reference":     0,
			"traceability":  0,
			"report":        0,
		},
		"total_count": 0,
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	return string(data)
}

// generateSomedayManifestContent generates the initial someday manifest
func (r *Runtime) generateSomedayManifestContent() string {
	manifest := map[string]interface{}{
		"version":      "1.0.0",
		"created_at":   time.Now().Format(time.RFC3339),
		"last_updated": time.Now().Format(time.RFC3339),
		"ideas":        []string{},
		"total_count":  0,
		"by_status": map[string]int{
			"SOMEDAY":     0,
			"LABORATORY":  0,
			"PROMOTED":    0,
			"IMPLEMENTED": 0,
			"ARCHIVED":    0,
			"REJECTED":    0,
		},
		"by_priority": map[string]int{
			"priority-1": 0,
			"priority-2": 0,
			"priority-3": 0,
			"priority-4": 0,
			"priority-5": 0,
		},
		"next_idea_id": 1,
	}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	return string(data)
}

// FormatInitializeResult formats initialization result for display
func FormatInitializeResult(result *InitializeResult) string {
	var output string

	output += "╔═══════════════════════════════════════════════════════════════╗\n"
	output += "║           KDSE Evidence-Driven Initialization                 ║\n"
	output += "╠═══════════════════════════════════════════════════════════════╣\n"
	output += fmt.Sprintf("║ Workspace: %s\n", result.WorkspacePath)
	output += fmt.Sprintf("║ Confidence: %.2f\n", result.Confidence)
	output += fmt.Sprintf("║ Status: %s\n", boolToStatus(result.Success))
	output += "╠═══════════════════════════════════════════════════════════════╣\n"
	output += "║ Verification Results                                         ║\n"

	for _, v := range result.Verification {
		statusIcon := "✓"
		if v.Status == "FAIL" {
			statusIcon = "✗"
		}
		output += fmt.Sprintf("║ %s %-12s %s\n", statusIcon, v.Artifact, v.Status)
	}

	output += "╠═══════════════════════════════════════════════════════════════╣"
	if result.Success {
		output += "\n║ Status: OPERATIONAL                                          ║"
	} else {
		output += "\n║ Status: FAILED                                               ║"
		if len(result.Errors) > 0 {
			output += "\n╠═══════════════════════════════════════════════════════════════╣"
			output += "\n║ Errors                                                       ║"
			for _, err := range result.Errors {
				output += fmt.Sprintf("\n║   • %s\n", err)
			}
		}
	}

	output += "\n╚═══════════════════════════════════════════════════════════════╝\n"
	return output
}

// FormatVerificationReport formats verification report for display
func FormatVerificationReport(report *VerificationReport) string {
	var output string

	output += "╔═══════════════════════════════════════════════════════════════╗\n"
	output += "║              KDSE Runtime Self-Audit                         ║\n"
	output += "╠═══════════════════════════════════════════════════════════════╣"

	for _, c := range report.Components {
		statusIcon := "PASS"
		if c.Status == "FAIL" {
			statusIcon = "FAIL"
		}
		output += fmt.Sprintf("║ %-12s %-8s %s\n", c.Artifact, statusIcon, c.Path)
	}

	output += "╠═══════════════════════════════════════════════════════════════╣"
	output += fmt.Sprintf("║ Confidence: %.2f\n", report.Confidence)

	if report.Success {
		output += "║ Status: OPERATIONAL                                          ║\n"
	} else {
		output += "║ Status: FAILED                                               ║\n"
		if len(report.Failed) > 0 {
			output += "║ Failed Components:                                           ║\n"
			for _, f := range report.Failed {
				output += fmt.Sprintf("║   • %s\n", f)
			}
		}
	}

	output += "╚═══════════════════════════════════════════════════════════════╝\n"
	return output
}

func boolToStatus(b bool) string {
	if b {
		return "SUCCESS"
	}
	return "FAILED"
}
