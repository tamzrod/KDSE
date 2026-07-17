# Architecture Audit: KDSE Ownership Boundary Drift

**Date:** 2026-07-17
**Status:** CRITICAL - Architectural Correction Required
**Severity:** HIGH

## Executive Summary

KDSE has experienced architectural drift where the `.kdse/` directory has been incorrectly treated as the project workspace rather than the engineering runtime. This violates KDSE's core engineering philosophy: **KDSE augments software engineering, it never replaces it.**

## Problem Statement

The current KDSE implementation incorrectly assumes that:
1. `.kdse/` is the authoritative foundation for all project artifacts
2. Engineering artifacts (architecture, implementation, verification) belong in `.kdse/`
3. Standard software engineering practices should be subordinate to KDSE runtime

This creates a dangerous anti-pattern where engineers unfamiliar with KDSE cannot understand, build, or maintain the project.

## Current State Analysis

### Directory Structure (Current)

```
Project/
├── README.md                           ✓ Project layer (correct)
├── LICENSE                             ✓ Project layer (correct)
├── docs/                               ✓ Project layer (correct)
├── src/                                ✓ Project layer (correct)
├── tests/                              ✓ Project layer (correct)
├── laboratory/                         ✗ Should be in .kdse/laboratory/
├── runtime/                            ✗ Should be in .kdse/runtime/
├── deploy/                             ✓ Project layer (correct)
├── cmd/                                ✓ Project layer (correct)
├── internal/                           ✓ Project layer (correct)
├── templates/                          ✓ Project layer (correct)
├── examples/                           ✓ Project layer (correct)
├── go.mod                              ✓ Project layer (correct)
├── Dockerfile                          ✓ Project layer (correct)
└── .kdse/
    ├── architecture/                   ✓ Runtime layer (correct - engineering process)
    ├── assessments/                   ✓ Runtime layer (correct - engineering process)
    ├── bootstrap/                      ✓ Runtime layer (correct - runtime config)
    ├── evidence/                       ✓ Runtime layer (correct - engineering evidence)
    ├── knowledge/                      ✓ Runtime layer (correct - extracted knowledge)
    ├── phase0/                         ✓ Runtime layer (correct - runtime init)
    ├── runtime/                        ✓ Runtime layer (correct - runtime state)
    └── verification/                   ✓ Runtime layer (correct - engineering verification)
```

### Architectural Violations

| # | Violation | Current Location | Correct Location | Severity |
|---|-----------|-----------------|-----------------|----------|
| 1 | Laboratory artifacts | `/laboratory/` | `.kdse/laboratory/` | HIGH |
| 2 | Runtime documentation | `/runtime/` | `docs/runtime/` or `.kdse/runtime/` | MEDIUM |
| 3 | README misleads users | README.md | N/A (update content) | HIGH |
| 4 | Engine creates artifacts in .kdse/ | `internal/workspace/engine.go` | Project root or .kdse/ only for runtime |

## Root Cause Analysis

### Contributing Factors

1. **Misaligned Philosophy**: Early design treated `.kdse/` as the authoritative state, leading to all artifacts being stored there
2. **Missing Ownership Model**: No explicit classification of artifact ownership
3. **Unclear Boundaries**: No enforcement mechanism preventing project artifacts from entering `.kdse/`
4. **Documentation Confusion**: README.md and other docs propagated the wrong mental model

### Evidence of Drift

```go
// internal/workspace/engine.go - Lines 166-193
// Creates knowledge/, architecture/, implementation/, verification/, reports/ in .kdse/

knowledgeDir := filepath.Join(stagingKDSE, "knowledge")
architectureDir := filepath.Join(stagingKDSE, "architecture")
implementationDir := filepath.Join(stagingKDSE, "implementation")
verificationDir := filepath.Join(stagingKDSE, "verification")
reportsDir := filepath.Join(stagingKDSE, "reports")
```

This code creates engineering artifacts inside `.kdse/`, which violates the ownership boundary.

## Target Architecture

### Corrected Directory Structure

```
Project/
├── README.md                           ✓ Project layer - Project documentation
├── LICENSE                             ✓ Project layer - Legal
├── CHANGELOG.md                        ✓ Project layer - Version history
├── docs/                               ✓ Project layer - Project documentation
│   ├── architecture/                   ✓ Project layer - Software architecture
│   ├── api/                            ✓ Project layer - API documentation
│   ├── deployment/                     ✓ Project layer - Deployment docs
│   └── runtime/                        ✓ Project layer - Runtime usage docs
├── src/                                ✓ Project layer - Source code
├── tests/                              ✓ Project layer - Test code
├── cmd/                                ✓ Project layer - CLI entrypoints
├── internal/                           ✓ Project layer - Internal packages
├── deploy/                             ✓ Project layer - Deployment configs
├── templates/                         ✓ Project layer - Project templates
├── examples/                           ✓ Project layer - Usage examples
├── .github/                            ✓ Project layer - CI/CD config
└── .kdse/                             ✓ Runtime layer - KDSE runtime
    ├── runtime/                        ✓ Runtime layer - Runtime state
    ├── sessions/                      ✓ Runtime layer - Engineering sessions
    ├── state/                         ✓ Runtime layer - State management
    ├── cache/                         ✓ Runtime layer - Cached data
    ├── reports/                       ✓ Runtime layer - Engineering reports
    ├── evidence/                      ✓ Runtime layer - Engineering evidence
    ├── traceability/                  ✓ Runtime layer - Traceability links
    ├── references/                    ✓ Reference layer - External standards
    ├── knowledge/                     ✓ Knowledge layer - Extracted knowledge
    └── laboratory/                    ✓ Runtime layer - Engineering laboratory
```

### Four Ownership Domains

| Domain | Owner | Contains | Location |
|--------|-------|----------|----------|
| **Project Layer** | Software Project | Deliverables, documentation, code, tests | Project root |
| **Runtime Layer** | KDSE Runtime | State, sessions, evidence, reports | `.kdse/` |
| **Reference Layer** | KDSE Runtime | External standards, vendor docs | `.kdse/references/` |
| **Knowledge Layer** | KDSE Runtime | Extracted knowledge, analysis | `.kdse/knowledge/` |

## Impact Assessment

### Business Impact

| Impact Area | Current State | Target State | Risk |
|-------------|---------------|--------------|------|
| Maintainability | Requires KDSE knowledge | Standard software project | HIGH |
| Onboarding | KDSE training required | Standard onboarding | HIGH |
| Portability | Locked to KDSE | Portable, standard project | HIGH |
| Compliance | Architecture drift | Clear ownership | MEDIUM |

### Technical Impact

| Component | Current Behavior | Required Behavior |
|-----------|-----------------|------------------|
| `internal/workspace/engine.go` | Creates artifacts in `.kdse/` | Creates project layout in root |
| `kdse_shttp_tools.py` | Assumes `.kdse/` is workspace | Enforces ownership boundaries |
| Initialization | Creates `.kdse/` only | Creates project root + `.kdse/` |
| Artifact routing | Undefined | Classified by ownership rules |

## Remediation Plan

### Phase 1: Documentation Update
- [ ] Create this architecture audit
- [ ] Update README.md to reflect correct philosophy
- [ ] Create ownership model documentation
- [ ] Update internal documentation

### Phase 2: Code Updates
- [ ] Update `internal/workspace/engine.go` initialization
- [ ] Add ownership classification to artifact creation
- [ ] Implement enforcement mechanism
- [ ] Update MCP tools for ownership awareness

### Phase 3: Migration
- [ ] Move `/laboratory/` to `.kdse/laboratory/`
- [ ] Move runtime docs to `docs/runtime/`
- [ ] Update all references
- [ ] Verify no project artifacts in `.kdse/`

### Phase 4: Verification
- [ ] Test initialization creates standard project layout
- [ ] Verify artifact routing enforces ownership
- [ ] Confirm deployment works
- [ ] Document verification results

## Recommendations

1. **Immediate**: Update README.md to correct the philosophical direction
2. **Short-term**: Update initialization to create standard project layout
3. **Medium-term**: Implement artifact ownership classification
4. **Long-term**: Build enforcement mechanism into runtime

## Conclusion

KDSE must return to its core mission: **augmenting software engineering**. The current architecture treats KDSE as the project foundation, which is backwards. The project must always remain the foundation, with KDSE providing runtime support and engineering discipline.

---

**Prepared by:** Architecture Audit
**Reviewed by:** [Pending]
**Approved by:** [Pending]
