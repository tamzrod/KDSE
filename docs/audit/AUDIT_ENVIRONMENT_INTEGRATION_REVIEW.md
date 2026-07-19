# Engineering Review: `kdse audit` Environment Detection Integration

**Document Version:** 1.0  
**Type:** Engineering Review  
**Date:** 2026-07-19  
**Reviewer:** OpenHands Agent  
**Status:** Draft  

---

## Executive Summary

This review examines the current implementation of `kdse audit` to determine whether it correctly integrates the new Environment Detection subsystem. The review finds that **Environment Detection is not currently integrated** into any audit functionality, and significant changes are required.

**Key Finding:** There is no `kdse audit` command in the current implementation. The closest equivalent is `kdse compliance` with an `--audit` flag, which does not invoke Environment Detection.

---

## 1. Current Audit Workflow

### 1.1 Command Structure

The current implementation has the following structure:

```
kdse compliance [--json] [--audit]
```

**Location:** `/cmd/kdse/main.go`, lines 1467-1513

### 1.2 Current Flow

```
kdse compliance
    │
    ├─► Parse flags: --json, --audit
    │
    ├─► ValidateCompliance(repoPath)
    │       │
    │       ├─► Check 1: .kdse exists?
    │       ├─► Check 2: Runtime structure valid?
    │       ├─► Check 3: No app code in .kdse?
    │       ├─► Check 4: Foundation templates exist?
    │       ├─► Check 5: No manual folder creation?
    │       ├─► Check 6: Knowledge has content?
    │       └─► Check 7: Evidence has content?
    │
    ├─► Output report (JSON or formatted)
    │
    └─► If --audit flag AND compliant:
            └─► Print "KDSE Pre-Implementation Audit"
                ✓ Project root identified
                ✓ .kdse directory exists
                ✓ Runtime verification passed
                ✓ All templates exist
                ✓ No manual folder creation detected
                ✓ Compliance validated
```

### 1.3 Current Behavior

**What the current implementation does:**

1. **Assumes KDSE_PROJECT**: The entire `ValidateCompliance` function assumes the workspace is a KDSE-initialized project with a `.kdse/` directory.

2. **No Environment Detection**: There is no call to the Environment Detection subsystem (`internal/detection/environment.go`).

3. **No Profile Selection**: There is no logic to select different audit profiles based on workspace type.

4. **No Subject Isolation**: The same compliance checks run regardless of what kind of workspace is being audited.

5. **No Failure Handling for UNKNOWN**: If the workspace cannot be classified, the compliance check will fail with a confusing error.

### 1.4 Code Reference

```go
// cmd/kdse/main.go:1467
func handleCompliance(repoPath string, args []string) {
    jsonOutput := false
    fullAudit := false
    
    // Parse flags...
    
    // Run compliance validation
    report, err := kdseruntime.ValidateCompliance(repoPath)
    
    // Output...
}

// internal/runtime/compliance.go:51
func ValidateCompliance(projectRoot string) (*ComplianceReport, error) {
    // Check 1: .kdse exists
    kdsePath := filepath.Join(projectRoot, ".kdse")
    if _, err := os.Stat(kdsePath); os.IsNotExist(err) {
        report.Violations = append(report.Violations, Violation{
            Type:    ViolationMissingRuntime,
            Message: ".kdse directory does not exist. Run 'kdse initialize' first.",
        })
        report.Compliant = false
        return report, nil
    }
    // ... more checks assuming .kdse exists
}
```

---

## 2. Correct Workflow

### 2.1 Expected Command Structure

```
kdse audit [--json] [--profile <name>]
```

### 2.2 Expected Flow

```
kdse audit
    │
    ├─► Parse flags: --json, --profile
    │
    ├─► Environment Detection (FIRST STEP)
    │       │
    │       └─► DetectEnvironment(repoPath)
    │               │
    │               ├─► KDSE_RUNTIME?
    │               ├─► KDSE_PROJECT?
    │               ├─► SOFTWARE_PROJECT?
    │               ├─► BLANK_WORKSPACE?
    │               └─► UNKNOWN?
    │
    ├─► Display Classification Header
    │       │
    │       ├─► Workspace Classification: <type>
    │       ├─► Detection Confidence: <0.00-1.00>
    │       └─► Evidence Summary: <markers found>
    │
    ├─► Handle UNKNOWN Classification
    │       │
    │       ├─► STOP auditing
    │       ├─► Report missing evidence
    │       └─► Recommend corrective action
    │
    ├─► Select Audit Profile
    │       │
    │       ├─► KDSE_RUNTIME → Runtime Audit Profile
    │       ├─► KDSE_PROJECT → Project Audit Profile
    │       ├─► SOFTWARE_PROJECT → Readiness Audit Profile
    │       └─► BLANK_WORKSPACE → Initialization Readiness Profile
    │
    ├─► Execute Selected Audit
    │       │
    │       ├─► Runtime Audit: Check internal/ package structure
    │       ├─► Project Audit: Check .kdse/ structure and compliance
    │       ├─► Readiness Audit: Check for project readiness
    │       └─► Initialization Audit: Check initialization requirements
    │
    └─► Output Report with Classification Header
            │
            ├─► Workspace Classification
            ├─► Detection Confidence
            ├─► Evidence Used
            ├─► Selected Profile
            ├─► Audit Findings
            └─► Recommendations
```

### 2.3 Expected Output Format

```
╔═══════════════════════════════════════════════════════════════════════╗
║                    KDSE AUDIT REPORT                                  ║
╠═══════════════════════════════════════════════════════════════════════╣
║ WORKSPACE CLASSIFICATION                                              ║
║ ────────────────────────                                              ║
║ Environment  : KDSE_RUNTIME                                           ║
║ Confidence   : 0.98                                                   ║
║ Evidence     : go.mod(github.com/kdse/runtime), internal/runtime/,    ║
║               internal/detection/, templates/, docs/                  ║
╠═══════════════════════════════════════════════════════════════════════╣
║ SELECTED PROFILE : Runtime Audit                                      ║
╠═══════════════════════════════════════════════════════════════════════╣
║ AUDIT FINDINGS                                                       ║
...
```

---

## 3. Missing Integration Points

### 3.1 Missing: `kdse audit` Command Entry Point

**Current State:** No `audit` command exists in the switch statement.

**Location:** `/cmd/kdse/main.go`, lines 132-179

**Missing Code:**
```go
case "audit":
    handleAudit(repoPath, args)
```

### 3.2 Missing: Environment Detection Import

**Current State:** The detection package is not imported in main.go.

**Current imports:**
```go
import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    
    "github.com/kdse/runtime/internal/agreement"
    "github.com/kdse/runtime/internal/bootstrap"
    // ... other imports
)
```

**Missing import:**
```go
"github.com/kdse/runtime/internal/detection"
```

### 3.3 Missing: Classification Header Display

**Current State:** No code displays the workspace classification before audit findings.

**Expected location:** After Environment Detection, before audit execution.

### 3.4 Missing: UNKNOWN Handling

**Current State:** No handling for `UNKNOWN` environment type.

**Expected behavior:** Stop and report missing evidence.

### 3.5 Missing: Audit Profile Selection

**Current State:** No logic to select different audit profiles.

**Expected logic:**
```go
func selectAuditProfile(env detection.EnvironmentType) string {
    switch env {
    case detection.EnvironmentKDSERuntime:
        return "runtime"
    case detection.EnvironmentKDSEProject:
        return "project"
    case detection.EnvironmentSoftwareProject:
        return "readiness"
    case detection.EnvironmentBlankWorkspace:
        return "initialization"
    default:
        return "unknown"
    }
}
```

### 3.6 Missing: Subject Isolation

**Current State:** `ValidateCompliance` always assumes KDSE_PROJECT.

**Required:** Different compliance checks for different workspace types.

---

## 4. Required Code Changes

### 4.1 New Command Handler

**File:** `cmd/kdse/main.go`

Add new function:
```go
func handleAudit(repoPath string, args []string) {
    jsonOutput := false
    profile := ""
    
    // Parse flags
    for _, arg := range args {
        switch arg {
        case "--json":
            jsonOutput = true
        case "--profile":
            // Handle profile override
        }
    }
    
    // STEP 1: Environment Detection
    detector := detection.NewEnvironmentDetector(repoPath)
    result := detector.Detect()
    
    // Display Classification Header
    displayClassificationHeader(result)
    
    // STEP 2: Handle UNKNOWN
    if result.Environment == detection.EnvironmentUnknown {
        handleUnknownEnvironment(result)
        os.Exit(1)
    }
    
    // STEP 3: Select and Execute Audit Profile
    profile := selectAuditProfile(result.Environment)
    auditReport := executeAudit(repoPath, profile, result)
    
    // Output
    if jsonOutput {
        fmt.Println(auditReport.ToJSON())
    } else {
        fmt.Println(formatAuditReport(auditReport))
    }
}
```

### 4.2 Classification Header Function

**File:** `cmd/kdse/main.go`

```go
func displayClassificationHeader(result *detection.EvidenceResult) {
    fmt.Println()
    fmt.Println("╔═══════════════════════════════════════════════════════════════════════╗")
    fmt.Println("║ WORKSPACE CLASSIFICATION                                              ║")
    fmt.Println("╠═══════════════════════════════════════════════════════════════════════╣")
    fmt.Printf("║ Environment  : %s\n", result.Environment)
    fmt.Printf("║ Confidence   : %.2f\n", result.Confidence)
    fmt.Printf("║ Evidence     : %s\n", formatEvidence(result.Evidence))
    fmt.Println("╚═══════════════════════════════════════════════════════════════════════╝")
    fmt.Println()
}
```

### 4.3 UNKNOWN Handler Function

**File:** `cmd/kdse/main.go`

```go
func handleUnknownEnvironment(result *detection.EvidenceResult) {
    fmt.Println()
    fmt.Println("╔═══════════════════════════════════════════════════════════════════════╗")
    fmt.Println("║ AUDIT CANNOT PROCEED                                                 ║")
    fmt.Println("╠═══════════════════════════════════════════════════════════════════════╣")
    fmt.Println("║ Environment could not be determined. Evidence is insufficient.      ║")
    fmt.Println("║                                                                       ║")
    fmt.Println("║ Missing Evidence:                                                     ║")
    for _, warning := range result.Warnings {
        fmt.Printf("║   • %s\n", warning)
    }
    fmt.Println("║                                                                       ║")
    fmt.Println("║ Corrective Actions:                                                   ║")
    fmt.Println("║   • Ensure you are in a valid project directory                       ║")
    fmt.Println("║   • For KDSE projects: Ensure .kdse/manifest.json exists             ║")
    fmt.Println("║   • For software projects: Ensure project files exist (go.mod, etc.) ║")
    fmt.Println("║   • For KDSE runtime: Ensure you are in the KDSE repository         ║")
    fmt.Println("╚═══════════════════════════════════════════════════════════════════════╝")
}
```

### 4.4 Audit Profile Selection

**File:** `cmd/kdse/main.go`

```go
func selectAuditProfile(env detection.EnvironmentType) string {
    switch env {
    case detection.EnvironmentKDSERuntime:
        return "runtime"
    case detection.EnvironmentKDSEProject:
        return "project"
    case detection.EnvironmentSoftwareProject:
        return "readiness"
    case detection.EnvironmentBlankWorkspace:
        return "initialization"
    default:
        return "unknown"
    }
}
```

### 4.5 Audit Report Struct with Classification

**File:** `internal/runtime/compliance.go`

Extend `ComplianceReport`:
```go
type ComplianceReport struct {
    Timestamp    string     `json:"timestamp"`
    Compliant    bool       `json:"compliant"`
    Violations   []Violation `json:"violations,omitempty"`
    Warnings     []string   `json:"warnings,omitempty"`
    Evidence     []string   `json:"evidence,omitempty"`
    
    // NEW: Workspace classification
    WorkspaceClassification *detection.EnvironmentType `json:"workspace_classification,omitempty"`
    ClassificationConfidence float64                   `json:"classification_confidence,omitempty"`
    ClassificationEvidence  *detection.Evidence        `json:"classification_evidence,omitempty"`
    SelectedProfile         string                    `json:"selected_profile,omitempty"`
}
```

### 4.6 New Audit Profiles

**File:** `internal/runtime/audit_profiles.go` (new file)

```go
package runtime

import "github.com/kdse/runtime/internal/detection"

// AuditProfile defines an audit profile for a workspace type
type AuditProfile struct {
    Name        string
    Description string
    Checks      []AuditCheck
}

// AuditCheck defines a single audit check
type AuditCheck struct {
    Name        string
    Execute     func(repoPath string) *AuditResult
    Severity    string
    Description string
}

// GetAuditProfile returns the audit profile for a workspace type
func GetAuditProfile(env detection.EnvironmentType) *AuditProfile {
    switch env {
    case detection.EnvironmentKDSERuntime:
        return GetRuntimeAuditProfile()
    case detection.EnvironmentKDSEProject:
        return GetProjectAuditProfile()
    case detection.EnvironmentSoftwareProject:
        return GetReadinessAuditProfile()
    case detection.EnvironmentBlankWorkspace:
        return GetInitializationAuditProfile()
    default:
        return nil
    }
}
```

---

## 5. Architectural Issues Discovered

### 5.1 Issue: Compliance Assumes KDSE_PROJECT

**Severity:** High

**Description:** `ValidateCompliance()` in `internal/runtime/compliance.go` assumes the workspace is a KDSE_PROJECT with an existing `.kdse/` directory. It will fail or produce misleading results for other workspace types.

**Current behavior:**
```go
// compliance.go:59-70
if _, err := os.Stat(kdsePath); os.IsNotExist(err) {
    report.Violations = append(report.Violations, Violation{
        Type:        ViolationMissingRuntime,
        Message:     ".kdse directory does not exist. Run 'kdse initialize' first.",
        Severity:    SeverityCritical,
    })
    // This is correct for KDSE_PROJECT, but wrong for:
    // - KDSE_RUNTIME: Should not require .kdse/
    // - SOFTWARE_PROJECT: Should suggest initialization
    // - BLANK_WORKSPACE: Should explain no project exists
}
```

**Recommendation:** The compliance check should be workspace-aware, producing different messages based on the detected environment type.

### 5.2 Issue: No Separate `kdse audit` Command

**Severity:** Medium

**Description:** The "audit" functionality is buried as a flag (`--audit`) to the `compliance` command. This makes it harder to discover and use.

**Current:**
```
kdse compliance --audit
```

**Expected:**
```
kdse audit
```

**Recommendation:** Create a dedicated `kdse audit` command that wraps the appropriate functionality.

### 5.3 Issue: No Environment Detection in Main Entry Point

**Severity:** High

**Description:** The main.go file doesn't call Environment Detection before routing commands. This means commands like `audit`, `report`, `status`, etc. cannot be workspace-aware.

**Current:**
```go
switch cmd {
case "audit":
    handleAudit(repoPath, args)  // No environment context passed
// ...
}
```

**Recommendation:** Either:
1. Call Environment Detection at the start of each command handler that needs it, OR
2. Call Environment Detection once in main() and pass the result to handlers

### 5.4 Issue: Audit Reports Don't Include Classification

**Severity:** Medium

**Description:** Current audit/compliance reports don't display the workspace classification, making it unclear what kind of workspace was audited.

**Current report header:**
```
╔═══════════════════════════════════════════════════════════════════════╗
║              KDSE RUNTIME COMPLIANCE REPORT                        ║
```

**Expected header:**
```
╔═══════════════════════════════════════════════════════════════════════╗
║ WORKSPACE CLASSIFICATION                                              ║
║ Environment  : KDSE_PROJECT                                          ║
║ Confidence   : 0.95                                                  ║
╠═══════════════════════════════════════════════════════════════════════╣
║              KDSE PROJECT COMPLIANCE REPORT                          ║
```

---

## 6. Implementation Plan

### Phase 1: Add `kdse audit` Command (Minimal)

1. Add `"audit"` case to switch statement
2. Create `handleAudit()` function
3. Call `detection.DetectEnvironment()` 
4. Display classification header
5. Handle UNKNOWN case
6. Select and call existing `ValidateCompliance()`

### Phase 2: Extend Compliance Report

1. Add classification fields to `ComplianceReport`
2. Populate classification in `ValidateCompliance()`
3. Display classification in report output

### Phase 3: Add Audit Profiles

1. Create `audit_profiles.go` with profile definitions
2. Implement runtime audit checks
3. Implement readiness audit checks
4. Implement initialization audit checks

### Phase 4: Subject Isolation

1. Modify `ValidateCompliance()` to be workspace-aware
2. Implement separate checks for each workspace type
3. Add appropriate messages for each workspace type

---

## 7. Summary

| Finding | Severity | Status |
|---------|----------|--------|
| No `kdse audit` command | Medium | Needs implementation |
| No Environment Detection integration | High | Needs implementation |
| No classification header in reports | Medium | Needs implementation |
| No UNKNOWN handling | High | Needs implementation |
| No audit profile selection | High | Needs implementation |
| Compliance assumes KDSE_PROJECT | High | Needs architectural fix |
| No subject isolation | High | Needs implementation |

**Total Issues:** 7  
**Critical (blocks audit):** 4  
**High (affects correctness):** 3  
**Medium (usability):** 1  

---

## 8. Recommendations

1. **Immediate:** Add `kdse audit` command that invokes Environment Detection first
2. **Short-term:** Extend compliance reports to include workspace classification
3. **Medium-term:** Implement separate audit profiles for each workspace type
4. **Long-term:** Make all KDSE commands workspace-aware

---

*End of Engineering Review*
