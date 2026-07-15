// Package runtime implements evidence-driven KDSE runtime management.
package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ComplianceReport represents the result of a compliance check
type ComplianceReport struct {
	Timestamp    string     `json:"timestamp"`
	Compliant    bool       `json:"compliant"`
	Violations   []Violation `json:"violations,omitempty"`
	Warnings     []string   `json:"warnings,omitempty"`
	Evidence     []string   `json:"evidence,omitempty"`
}

// Violation represents a compliance violation
type Violation struct {
	Type        string `json:"type"`
	Path        string `json:"path,omitempty"`
	Message     string `json:"message"`
	Severity    string `json:"severity"`
	CorrectiveAction string `json:"corrective_action,omitempty"`
}

// Severity constants
const (
	SeverityCritical = "CRITICAL"
	SeverityHigh     = "HIGH"
	SeverityMedium   = "MEDIUM"
	SeverityLow      = "LOW"
)

// Violation type constants
const (
	ViolationMissingRuntime    = "MISSING_RUNTIME"
	ViolationRuntimeInvalid    = "RUNTIME_INVALID"
	ViolationAppCodeInRuntime  = "APP_CODE_IN_RUNTIME"
	ViolationDocsOutsideRuntime = "DOCS_OUTSIDE_RUNTIME"
	ViolationManualFolder      = "MANUAL_FOLDER_CREATION"
	ViolationPhaseViolation    = "PHASE_VIOLATION"
	ViolationVerificationMissing = "VERIFICATION_MISSING"
)

// ValidateCompliance checks if the project follows KDSE runtime rules
func ValidateCompliance(projectRoot string) (*ComplianceReport, error) {
	report := &ComplianceReport{
		Timestamp:  time.Now().Format(time.RFC3339),
		Violations: []Violation{},
		Warnings:   []string{},
		Evidence:   []string{},
	}

	// Check 1: .kdse exists
	kdsePath := filepath.Join(projectRoot, ".kdse")
	if _, err := os.Stat(kdsePath); os.IsNotExist(err) {
		report.Violations = append(report.Violations, Violation{
			Type:           ViolationMissingRuntime,
			Message:        ".kdse directory does not exist. Run 'kdse initialize' first.",
			Severity:       SeverityCritical,
			CorrectiveAction: "Execute 'kdse initialize' to create the runtime.",
		})
		report.Compliant = false
		return report, nil
	}
	report.Evidence = append(report.Evidence, fmt.Sprintf(".kdse exists at %s", kdsePath))

	// Check 2: Verify runtime structure
	runtime := New(projectRoot)
	verifyResult := runtime.Verify()
	if !verifyResult.Success {
		report.Violations = append(report.Violations, Violation{
			Type:           ViolationRuntimeInvalid,
			Message:        fmt.Sprintf("Runtime verification failed. Confidence: %.2f", verifyResult.Confidence),
			Severity:       SeverityCritical,
			CorrectiveAction: "Run 'kdse initialize' or 'kdse runtime fix' to restore runtime.",
		})
	} else {
		report.Evidence = append(report.Evidence, fmt.Sprintf("Runtime verification passed (Confidence: %.2f)", verifyResult.Confidence))
	}

	// Check 3: No application code in .kdse
	appCodeViolation := checkNoAppCodeInRuntime(projectRoot, kdsePath)
	if appCodeViolation != nil {
		report.Violations = append(report.Violations, *appCodeViolation)
	} else {
		report.Evidence = append(report.Evidence, "No application code found in .kdse/")
	}

	// Check 4: Foundation templates exist and contain project-specific content
	foundationViolation := checkFoundationCompliance(projectRoot, kdsePath)
	if foundationViolation != nil {
		report.Violations = append(report.Violations, *foundationViolation)
	} else {
		report.Evidence = append(report.Evidence, "Foundation templates exist and are populated")
	}

	// Check 5: Runtime directories not created manually by checking timestamps
	manualFolderViolation := checkNoManualFolderCreation(projectRoot, kdsePath)
	if manualFolderViolation != nil {
		report.Violations = append(report.Violations, *manualFolderViolation)
	} else {
		report.Evidence = append(report.Evidence, "No manual folder creation detected in .kdse/")
	}

	// Check 6: Knowledge phase has content (if applicable)
	if hasContent(filepath.Join(kdsePath, DirKnowledge)) {
		report.Evidence = append(report.Evidence, "Knowledge collection has content")
	} else {
		report.Warnings = append(report.Warnings, "Knowledge collection is empty - consider populating during Knowledge Phase")
	}

	// Check 7: Evidence directory exists and has content (if implementation exists)
	evidencePath := filepath.Join(kdsePath, DirEvidence)
	if hasContent(evidencePath) {
		report.Evidence = append(report.Evidence, "Evidence collection has content")
	}

	// Determine overall compliance
	report.Compliant = len(report.Violations) == 0

	return report, nil
}

// checkNoAppCodeInRuntime verifies no application code exists in .kdse
func checkNoAppCodeInRuntime(projectRoot, kdsePath string) *Violation {
	// Common application file patterns
	appPatterns := []string{
		".go", ".java", ".py", ".js", ".ts", ".tsx", ".jsx",
		".c", ".cpp", ".h", ".rs", ".rb", ".php",
		"main.go", "index.js", "app.js", "server.js",
		"package.json", "requirements.txt", "Gemfile", "Cargo.toml",
	}

	// Common application directory patterns
	appDirs := []string{
		"src/", "app/", "lib/", "cmd/", "internal/",
		"static/", "templates/", "views/", "public/",
		"dist/", "build/", "bin/",
	}

	// Check for application directories in .kdse
	entries, err := os.ReadDir(kdsePath)
	if err != nil {
		return &Violation{
			Type:        ViolationAppCodeInRuntime,
			Path:        kdsePath,
			Message:     fmt.Sprintf("Cannot read .kdse directory: %v", err),
			Severity:    SeverityHigh,
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip known KDSE directories
		if isKDSEPath(entry.Name()) {
			continue
		}

		// Check if it's an application directory
		for _, appDir := range appDirs {
			if strings.TrimSuffix(appDir, "/") == entry.Name() {
				return &Violation{
					Type:           ViolationAppCodeInRuntime,
					Path:           filepath.Join(kdsePath, entry.Name()),
					Message:        fmt.Sprintf("Application directory '%s' found in .kdse. Application code MUST be outside .kdse/", entry.Name()),
					Severity:       SeverityHigh,
					CorrectiveAction: fmt.Sprintf("Move %s outside .kdse/ to %s/", entry.Name(), projectRoot),
				}
			}
		}
	}

	// Check for application files directly in .kdse root
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		for _, pattern := range appPatterns {
			if strings.HasSuffix(name, pattern) || name == pattern {
				return &Violation{
					Type:           ViolationAppCodeInRuntime,
					Path:           filepath.Join(kdsePath, name),
					Message:        fmt.Sprintf("Application file '%s' found in .kdse/. Application code MUST be outside .kdse/", name),
					Severity:       SeverityHigh,
					CorrectiveAction: fmt.Sprintf("Move %s outside .kdse/", name),
				}
			}
		}
	}

	return nil
}

// checkFoundationCompliance verifies foundation documents are populated
func checkFoundationCompliance(projectRoot, kdsePath string) *Violation {
	foundationPath := filepath.Join(kdsePath, DirFoundation)
	
	requiredFiles := []string{"PROBLEM.md", "SPEC.md", "ARCHITECTURE.md", "ASSUMPTIONS.md", "REQUIREMENTS.md"}
	
	for _, file := range requiredFiles {
		filePath := filepath.Join(foundationPath, file)
		
		// Check file exists
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			return &Violation{
				Type:           ViolationRuntimeInvalid,
				Path:           filePath,
				Message:        fmt.Sprintf("Required foundation file '%s' does not exist", file),
				Severity:       SeverityHigh,
				CorrectiveAction: "Run 'kdse initialize' to create foundation templates",
			}
		}

		// Check file is not empty
		if info.Size() < 100 {
			return &Violation{
				Type:           ViolationVerificationMissing,
				Path:           filePath,
				Message:        fmt.Sprintf("Foundation file '%s' appears to be a template (size: %d bytes). Populate it during the appropriate phase.", file, info.Size()),
				Severity:       SeverityMedium,
				CorrectiveAction: fmt.Sprintf("Update %s with project-specific content during the %s Phase", file, getFilePhase(file)),
			}
		}
	}

	return nil
}

// checkNoManualFolderCreation checks if folders were created by AI in .kdse
func checkNoManualFolderCreation(projectRoot, kdsePath string) *Violation {
	// List of directories that should exist after initialization
	expectedDirs := map[string]bool{
		DirRuntime:        true,
		DirFoundation:     true,
		DirKnowledge:      true,
		DirLaboratory:     true,
		DirEvidence:       true,
		DirReferences:     true,
		DirTraceability:   true,
		DirReports:        true,
		DirConfig:         true,
		DirState:          true,
		DirArtifacts:      true,
		DirSessions:       true,
		DirNormalized:     true,
		DirCache:          true,
		DirSomeday:        true,
	}

	entries, err := os.ReadDir(kdsePath)
	if err != nil {
		return &Violation{
			Type:        ViolationManualFolder,
			Message:     fmt.Sprintf("Cannot read .kdse directory: %v", err),
			Severity:    SeverityMedium,
		}
	}

	var unexpectedDirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip if it's a known KDSE directory
		if expectedDirs[entry.Name()] {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		unexpectedDirs = append(unexpectedDirs, entry.Name())
	}

	if len(unexpectedDirs) > 0 {
		return &Violation{
			Type:           ViolationManualFolder,
			Path:           kdsePath,
			Message:        fmt.Sprintf("Unexpected directories found in .kdse/: %v. Only KDSE runtime directories are allowed.", unexpectedDirs),
			Severity:       SeverityMedium,
			CorrectiveAction: "Remove manually created directories from .kdse/ or report if they are intentional KDSE extensions",
		}
	}

	return nil
}

// isKDSEPath checks if a path name is a known KDSE directory
func isKDSEPath(name string) bool {
	kdsePaths := map[string]bool{
		"runtime": true, "foundation": true, "knowledge": true,
		"laboratory": true, "evidence": true, "reports": true,
		"references": true, "traceability": true, "artifacts": true,
		"config": true, "state": true, "sessions": true,
		"normalized": true, "cache": true, "someday": true,
		"general": true, "operational": true, "developmental": true,
		"experiments": true, "ideas": true, "archived": true, "promoted": true,
	}
	return kdsePaths[name]
}

// getFilePhase returns the appropriate phase for a foundation file
func getFilePhase(filename string) string {
	switch filename {
	case "PROBLEM.md":
		return "Problem"
	case "SPEC.md", "ASSUMPTIONS.md", "REQUIREMENTS.md":
		return "Foundation"
	case "ARCHITECTURE.md":
		return "Architecture"
	default:
		return "Unknown"
	}
}

// hasContent checks if a directory has any files
func hasContent(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return true
		}
		// Check subdirectory
		if hasContent(filepath.Join(dirPath, entry.Name())) {
			return true
		}
	}

	return false
}

// FormatComplianceReport formats a compliance report for display
func FormatComplianceReport(report *ComplianceReport) string {
	var output string

	output += "╔═══════════════════════════════════════════════════════════════════════╗\n"
	output += "║              KDSE RUNTIME COMPLIANCE REPORT                        ║\n"
	output += "╠═══════════════════════════════════════════════════════════════════════╣\n"
	output += fmt.Sprintf("║ Timestamp:   %s\n", report.Timestamp)
	output += fmt.Sprintf("║ Compliant:  %s\n", boolToYesNo(report.Compliant))
	output += "╠═══════════════════════════════════════════════════════════════════════╣\n"

	if len(report.Violations) > 0 {
		output += "║ VIOLATIONS                                                     ║\n"
		for _, v := range report.Violations {
			severityIcon := getSeverityIcon(v.Severity)
			output += fmt.Sprintf("║ %s %s [%s]\n", severityIcon, v.Type, v.Severity)
			output += fmt.Sprintf("║   %s\n", truncate(v.Message, 60))
			if v.Path != "" {
				output += fmt.Sprintf("║   Path: %s\n", truncate(v.Path, 55))
			}
			if v.CorrectiveAction != "" {
				output += fmt.Sprintf("║   Fix: %s\n", truncate(v.CorrectiveAction, 55))
			}
		}
		output += "╠═══════════════════════════════════════════════════════════════════════╣\n"
	}

	if len(report.Warnings) > 0 {
		output += "║ WARNINGS                                                       ║\n"
		for _, w := range report.Warnings {
			output += fmt.Sprintf("║ ⚠ %s\n", truncate(w, 62))
		}
		output += "╠═══════════════════════════════════════════════════════════════════════╣\n"
	}

	if len(report.Evidence) > 0 {
		output += "║ EVIDENCE                                                       ║\n"
		for _, e := range report.Evidence {
			output += fmt.Sprintf("║ ✓ %s\n", truncate(e, 62))
		}
		output += "╠═══════════════════════════════════════════════════════════════════════╣\n"
	}

	if report.Compliant {
		output += "║ Status:    COMPLIANT                                           ║\n"
		output += "║ KDSE Runtime Contract: SATISFIED                              ║\n"
	} else {
		output += "║ Status:    NON-COMPLIANT                                       ║\n"
		output += "║ KDSE Runtime Contract: VIOLATED                               ║\n"
		output += "║                                                               ║\n"
		output += "║ STOP. Do not continue until compliance is restored.           ║\n"
	}

	output += "╚═══════════════════════════════════════════════════════════════════════╝\n"

	return output
}

func boolToYesNo(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

func getSeverityIcon(severity string) string {
	switch severity {
	case SeverityCritical:
		return "🔴"
	case SeverityHigh:
		return "🟠"
	case SeverityMedium:
		return "🟡"
	case SeverityLow:
		return "⚪"
	default:
		return "❓"
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ToJSON returns the report as JSON
func (r *ComplianceReport) ToJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}
