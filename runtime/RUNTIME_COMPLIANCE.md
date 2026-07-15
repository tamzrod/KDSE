# KDSE Runtime Compliance Specification

## Overview

KDSE Runtime Compliance ensures that every KDSE project follows the correct engineering contract. This document defines the rules, verification procedures, and acceptance criteria for compliant KDSE implementations.

---

## Engineering Contract

### Pre-Work Requirements

Before any engineering work begins, the following MUST be verified:

1. **Project Root Determination** - Identify the project root directory
2. **Runtime Initialization** - Execute `kdse initialize`
3. **Runtime Verification** - Execute `kdse runtime verify`
4. **Compliance Check** - Run compliance validation

### Expected Structure

After initialization, the project root MUST contain:

```
PROJECT_ROOT/
├── .kdse/
│   ├── foundation/
│   │   ├── PROBLEM.md
│   │   ├── SPEC.md
│   │   ├── ARCHITECTURE.md
│   │   ├── ASSUMPTIONS.md
│   │   ├── REQUIREMENTS.md
│   │   └── README.md
│   ├── knowledge/
│   │   ├── general/
│   │   ├── operational/
│   │   └── developmental/
│   ├── laboratory/
│   │   ├── experiments/
│   │   └── reports/
│   ├── evidence/
│   ├── reports/
│   ├── references/
│   ├── traceability/
│   ├── artifacts/
│   ├── someday/
│   │   ├── ideas/
│   │   ├── archived/
│   │   └── promoted/
│   ├── config/
│   ├── state/
│   ├── runtime.yaml
│   ├── manifest.yaml
│   ├── session-state.yaml
│   ├── knowledge-index.yaml
│   └── artifact-index.yaml
└── [application code outside .kdse]
```

---

## Runtime Verification Protocol

### Step 1: Initialize Runtime

```bash
kdse initialize
```

### Step 2: Verify Runtime

```bash
kdse runtime verify
```

**Expected Output:**
```
Workspace:         PASS
Foundation:        PASS
Knowledge:         PASS
Laboratory:        PASS
Evidence:          PASS
Reports:           PASS
References:        PASS
Traceability:      PASS
Artifacts:         PASS
Configuration:     PASS
State:             PASS
Somday:            PASS

Confidence:       1.00
Status:           OPERATIONAL
```

### Step 3: Check Compliance

Execute compliance validation (see Implementation section).

---

## Ownership Boundaries

### KDSE Runtime Ownership

The following paths belong EXCLUSIVELY to KDSE:

| Path | Owner | AI Action |
|------|-------|-----------|
| `.kdse/foundation/` | Runtime | Modify templates only |
| `.kdse/knowledge/` | Runtime | Add files to subdirs only |
| `.kdse/laboratory/` | Runtime | Add experiments only |
| `.kdse/evidence/` | Runtime | Add evidence only |
| `.kdse/reports/` | Runtime | Add reports only |
| `.kdse/references/` | Runtime | Add references only |
| `.kdse/traceability/` | Runtime | Update matrices only |
| `.kdse/artifacts/` | Runtime | Add artifacts only |
| `.kdse/someday/` | Runtime | Manage ideas only |
| `.kdse/config/` | Runtime | Configuration only |
| `.kdse/state/` | Runtime | State only |

### Application Ownership

Application code MUST exist OUTSIDE `.kdse`:

```
PROJECT_ROOT/
├── .kdse/                    ← Engineering (Runtime)
├── src/                      ← Application (Project)
├── app/                      ← Application (Project)
├── lib/                      ← Application (Project)
├── tests/                    ← Application (Project)
├── docs/                     ← Application (Project)
├── requirements.txt          ← Application (Project)
├── package.json              ← Application (Project)
└── [other application files] ← Application (Project)
```

---

## Phase Compliance Rules

### Problem Phase

**Task:** Define the problem

**Allowed Actions:**
- ✓ Update `.kdse/foundation/PROBLEM.md`
- ✓ Read existing templates
- ✓ Run `kdse status`

**Forbidden Actions:**
- ✗ Create new directories in `.kdse/`
- ✗ Create application code
- ✗ Modify runtime files

### Foundation Phase

**Task:** Define solution scope

**Allowed Actions:**
- ✓ Update `.kdse/foundation/SPEC.md`
- ✓ Update `.kdse/foundation/ASSUMPTIONS.md`
- ✓ Update `.kdse/foundation/REQUIREMENTS.md`
- ✓ Update `.kdse/foundation/PROBLEM.md`

**Forbidden Actions:**
- ✗ Create new directories in `.kdse/`
- ✗ Create application code
- ✗ Modify runtime files

### Architecture Phase

**Task:** Design system architecture

**Allowed Actions:**
- ✓ Update `.kdse/foundation/ARCHITECTURE.md`
- ✓ Create diagrams (in project, not .kdse)
- ✓ Reference architecture decisions

**Forbidden Actions:**
- ✗ Create new directories in `.kdse/`
- ✗ Create application code
- ✗ Modify runtime files

### Knowledge Phase

**Task:** Collect engineering knowledge

**Allowed Actions:**
- ✓ Add files to `.kdse/knowledge/general/`
- ✓ Add files to `.kdse/knowledge/operational/`
- ✓ Add files to `.kdse/knowledge/developmental/`
- ✓ Update knowledge indices

**Forbidden Actions:**
- ✗ Create new directories in `.kdse/`
- ✗ Create application code
- ✗ Modify runtime files

### Implementation Phase

**Task:** Implement application

**Allowed Actions:**
- ✓ Create application code OUTSIDE `.kdse/`
- ✓ Create project structure OUTSIDE `.kdse/`
- ✓ Update `.kdse/artifacts/` with build artifacts
- ✓ Update traceability

**Forbidden Actions:**
- ✗ Create directories inside `.kdse/` (except artifacts)
- ✗ Put application code in `.kdse/`

### Verification Phase

**Task:** Verify implementation

**Allowed Actions:**
- ✓ Add test results to `.kdse/evidence/`
- ✓ Add performance metrics to `.kdse/evidence/`
- ✓ Add validation reports to `.kdse/evidence/`

**Forbidden Actions:**
- ✗ Create new directories in `.kdse/`
- ✗ Modify application code without justification

### Documentation Phase

**Task:** Document engineering work

**Allowed Actions:**
- ✓ Add reports to `.kdse/reports/`
- ✓ Update documentation OUTSIDE `.kdse/`
- ✓ Generate normalized docs

**Forbidden Actions:**
- ✗ Create new directories in `.kdse/`
- ✗ Move application code to `.kdse/`

---

## Violation Detection

### Automatic Violations

The following are automatically detected as violations:

| Violation | Detection Method |
|-----------|------------------|
| Directories created in .kdse manually | Compare creation timestamps |
| Application code in .kdse | Path analysis |
| Runtime structure missing | `kdse runtime verify` |
| Engineering docs outside .kdse | Path analysis |
| Manual folder creation in .kdse | Git history analysis |

### Manual Violations

The following require manual review:

| Violation | Detection Method |
|-----------|------------------|
| Runtime folders created by AI | Code review |
| Premature implementation | Phase analysis |
| Skipped phases | Artifact analysis |

---

## Acceptance Criteria

### Pre-Implementation

| Criterion | Verification | Status |
|-----------|--------------|--------|
| `.kdse/` exists in project root | `ls -la .kdse/` | REQUIRED |
| Runtime verification succeeds | `kdse runtime verify` | REQUIRED |
| Confidence equals 1.00 | Verification output | REQUIRED |
| All templates exist | File existence check | REQUIRED |
| No manual runtime folders | Git history check | REQUIRED |

### During Implementation

| Criterion | Verification | Status |
|-----------|--------------|--------|
| Application code outside .kdse | Path analysis | REQUIRED |
| Engineering docs inside .kdse | Path analysis | REQUIRED |
| No runtime folder creation | Code review | REQUIRED |
| Phase artifacts in correct location | Path analysis | REQUIRED |

### Post-Implementation

| Criterion | Verification | Status |
|-----------|--------------|--------|
| Evidence in `.kdse/evidence/` | Path analysis | REQUIRED |
| Reports in `.kdse/reports/` | Path analysis | REQUIRED |
| Application code verified | Functional test | REQUIRED |
| KDSE compliance maintained | Final audit | REQUIRED |

---

## Compliance Validation Implementation

### Function: ValidateRuntimeCompliance

```go
// ValidateRuntimeCompliance checks if the project follows KDSE runtime rules
func ValidateRuntimeCompliance(projectRoot string) (*ComplianceReport, error) {
    report := &ComplianceReport{
        Timestamp: time.Now().Format(time.RFC3339),
        Violations: []Violation{},
    }

    // Check 1: .kdse exists
    kdsePath := filepath.Join(projectRoot, ".kdse")
    if _, err := os.Stat(kdsePath); os.IsNotExist(err) {
        report.Violations = append(report.Violations, Violation{
            Type: "MISSING_RUNTIME",
            Message: ".kdse directory does not exist. Run 'kdse initialize' first.",
        })
        return report, nil
    }

    // Check 2: Verify runtime
    runtime := New(kdsePath)
    verifyResult := runtime.Verify()
    if !verifyResult.Success {
        report.Violations = append(report.Violations, Violation{
            Type: "RUNTIME_INVALID",
            Message: fmt.Sprintf("Runtime verification failed. Confidence: %.2f", verifyResult.Confidence),
        })
    }

    // Check 3: No application code in .kdse
    if err := checkNoAppCodeInRuntime(projectRoot); err != nil {
        report.Violations = append(report.Violations, Violation{
            Type: "APP_CODE_IN_RUNTIME",
            Message: err.Error(),
        })
    }

    // Check 4: Engineering docs in .kdse
    if err := checkEngineeringDocsInRuntime(projectRoot); err != nil {
        report.Violations = append(report.Violations, Violation{
            Type: "DOCS_OUTSIDE_RUNTIME",
            Message: err.Error(),
        })
    }

    // Check 5: No manual folder creation
    if err := checkNoManualFolderCreation(projectRoot); err != nil {
        report.Violations = append(report.Violations, Violation{
            Type: "MANUAL_FOLDER_CREATION",
            Message: err.Error(),
        })
    }

    report.Compliant = len(report.Violations) == 0
    return report, nil
}
```

### Violation Types

| Type | Severity | Description |
|------|----------|-------------|
| `MISSING_RUNTIME` | CRITICAL | .kdse does not exist |
| `RUNTIME_INVALID` | CRITICAL | Runtime verification failed |
| `APP_CODE_IN_RUNTIME` | HIGH | Application code found in .kdse |
| `DOCS_OUTSIDE_RUNTIME` | HIGH | Engineering docs outside .kdse |
| `MANUAL_FOLDER_CREATION` | MEDIUM | Folders created by AI in .kdse |
| `PHASE_VIOLATION` | MEDIUM | Wrong phase artifacts |
| `VERIFICATION_MISSING` | LOW | Evidence/reports missing |

---

## Enforcement Actions

### On Violation Detection

1. **STOP** - Do not continue engineering work
2. **REPORT** - Display violation details
3. **FIX** - Correct the violation
4. **VERIFY** - Re-run compliance check
5. **CONTINUE** - Resume only when compliant

### Violation Response Matrix

| Violation | Response |
|-----------|----------|
| MISSING_RUNTIME | Execute `kdse initialize` |
| RUNTIME_INVALID | Fix runtime or reinitialize |
| APP_CODE_IN_RUNTIME | Move code outside .kdse |
| DOCS_OUTSIDE_RUNTIME | Move docs to .kdse |
| MANUAL_FOLDER_CREATION | Remove manually created folders |
| PHASE_VIOLATION | Review phase requirements |
| VERIFICATION_MISSING | Add missing evidence/reports |

---

## Audit Checklist

### Pre-Implementation Audit

```
□ Project root identified
□ .kdse directory exists
□ Runtime initialized
□ Runtime verification passed (Confidence = 1.00)
□ All templates exist
□ No manual folder creation detected
□ Compliance validated
```

### During Engineering Audit

```
□ Working in correct phase
□ Phase artifacts in correct location
□ No runtime folder creation
□ Application code outside .kdse
□ Engineering docs inside .kdse
□ Phase progress tracked
```

### Post-Implementation Audit

```
□ Evidence collected in .kdse/evidence/
□ Reports generated in .kdse/reports/
□ Application code functional
□ Traceability matrix complete
□ Final compliance check passed
□ KDSE contract satisfied
```

---

## Implementation

The compliance validation is integrated into the KDSE runtime:

```bash
# Check runtime compliance
kdse runtime verify

# Check project compliance
kdse compliance check

# Full compliance audit
kdse audit full
```

---

## References

- [KDSE Runtime Architecture](RUNTIME_ARCHITECTURE.md)
- [Someday/Maybe Management](SOMEDAY_ARCHITECTURE.md)
- [Evidence-Driven Verification](EVIDENCE_MODEL.md)
- [Initialization Specification](INITIALIZATION_SPEC.md)
