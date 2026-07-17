# KDSE Implementation Plan

**Document Version:** 1.0  
**Type:** Implementation Guide  
**Status:** Draft  
**Author:** KDSE Engineering Team  
**Date:** 2026-07-17

---

## 1. Project Overview

### 1.1 What is the Project?

KDSE (Knowledge-Driven Software Engineering) is a Runtime-Centric Engineering Platform where the KDSE Runtime (`.kdse/`) is the authoritative foundation of every engineering project.

**Key Characteristics:**
- Evidence-based engineering methodology
- Single source of truth: the `.kdse/` directory
- Multiple runtimes (CLI, MCP) as thin adapters
- Sequential phase enforcement
- Workspace Engine as state owner

### 1.2 Architectural Principles

| Principle | Code | Description |
|-----------|------|-------------|
| Evidence First | P-001 | No engineering work without evidence |
| Runtime is Authority | P-002 | `.kdse/` is the authoritative state |
| One Methodology | P-003 | Only one engineering methodology |
| Multiple Runtimes | P-004 | CLI, MCP are execution adapters |
| Runtime Independence | P-005 | Methodology never depends on runtime |
| Filesystem is Evidence | P-006 | Without `.kdse/`, not a KDSE project |

### 1.3 Target Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                       KDSE Methodology                                │
│                      (Engineering Rules)                              │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Workspace Engine                                 │
│                       (State Owner)                                  │
└─────────────────────────────────────────────────────────────────────┘
                    ┌───────────────────┬───────────────────┐
                    ▼                   ▼                   ▼
            ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
            │ CLI Runtime │     │ MCP Runtime │     │ IDE Runtime │
            │(Thin Adapter│     │(Thin Adapter│     │  (Future)  │
            └─────────────┘     └─────────────┘     └─────────────┘
```

---

## 2. What Already Exists

### 2.1 Architecture Documents (Complete ✓)

| Document | Location | Status |
|---------|----------|--------|
| PRINCIPLES.md | docs/architecture/ | Complete |
| RUNTIME_ARCHITECTURE.md | docs/architecture/ | Complete |
| METHODOLOGY.md | docs/architecture/ | Complete |
| WORKSPACE_ENGINE.md | docs/architecture/ | Complete |
| CLI_RUNTIME.md | docs/architecture/ | Complete |
| MCP_RUNTIME.md | docs/architecture/ | Complete |
| CURRENT_ARCHITECTURE.md | docs/architecture/ | Complete |

### 2.2 Evolution Documents (Complete ✓)

| Document | Location | Status |
|---------|----------|--------|
| README.md | docs/architecture-evolution/ | Complete |
| KAE-001 | docs/architecture-evolution/ | Complete |
| KAE-002 | docs/architecture-evolution/ | Complete |
| KAE-003 | docs/architecture-evolution/ | Complete |
| KAE-004 | docs/architecture-evolution/ | Complete |
| MIGRATION.md | docs/architecture-evolution/ | Complete |

### 2.3 Architecture Decisions (Complete ✓)

| ADR | Decision | Status |
|-----|----------|--------|
| ADR-001 | Runtime is the Authority | Accepted |
| ADR-002 | One Methodology, Multiple Runtimes | Accepted |
| ADR-003 | Workspace Engine | Accepted |

### 2.4 Core Implementation (Partial ✓)

| Component | Location | Status |
|-----------|----------|--------|
| Phase definitions | internal/methodology/lifecycle/phases.go | Complete |
| Lifecycle interface | internal/methodology/lifecycle/lifecycle.go | Complete |
| Workspace types | internal/workspace/types.go | Complete |
| Engine interface | internal/workspace/engine.go | Partial |
| Engine helpers | internal/workspace/engine_helpers.go | Partial |
| Basic tests | internal/workspace/engine_test.go | Partial |

---

## 3. What is Missing

### 3.1 Workspace Engine Sub-packages (CRITICAL)

According to WORKSPACE_ENGINE.md section "Implementation Structure":

```
internal/workspace/
├── engine/
│   ├── engine.go           ← EXISTS (partial)
│   ├── verify.go           ← MISSING
│   ├── init.go             ← MISSING
│   ├── phase.go            ← MISSING
│   ├── artifact.go         ← MISSING
│   ├── session.go          ← MISSING
│   └── report.go           ← MISSING
├── loader/                  ← MISSING
│   ├── loader.go
│   └── config.go
├── validator/               ← MISSING
│   ├── validator.go
│   └── phase.go
└── state/                   ← MISSING
    ├── state.go
    └── persistence.go
```

### 3.2 Methodology Packages (CRITICAL)

According to METHODOLOGY.md section "Package Structure":

```
internal/methodology/
├── lifecycle/               ← EXISTS
│   ├── lifecycle.go        ← EXISTS
│   └── phases.go           ← EXISTS
├── phases/                  ← MISSING
│   ├── knowledge/
│   ├── architecture/
│   ├── implementation/
│   ├── verification/
│   └── reports/
├── authority/               ← MISSING
│   ├── authority.go
│   └── verification.go
├── verification/            ← MISSING
│   ├── verification.go
│   └── validators/
└── knowledge/               ← MISSING
    ├── knowledge.go
    └── promotion.go
```

### 3.3 Thin Adapter Runtimes (HIGH)

According to CLI_RUNTIME.md and MCP_RUNTIME.md:

**CLI Runtime (cmd/kdse/):**
- Current state: Contains mixed business logic
- Required: Thin adapter only
- Missing: commands/, formatter/, parser/ packages
- Refactor: Remove all business logic

**MCP Runtime (cmd/kdse-mcp/):**
- Current state: Needs verification
- Required: Thin adapter only
- Missing: server/, tools/, formatter/ structure

### 3.4 Runtime Interface Abstraction (HIGH)

Per RUNTIME_ARCHITECTURE.md:
- Runtime interface not fully abstracted
- CLI and MCP should implement same interface
- Missing: Runtime interface definition with implementations

### 3.5 Template System (MEDIUM)

Per I-003 (Template-Based Bootstrap):
- templates/ directory exists
- Missing: Template-based initialization
- Missing: GitHub bootstrap integration
- Missing: Template versioning

### 3.6 Comprehensive Testing (MEDIUM)

Per KAE-001 Phase 13:
- Basic tests exist
- Missing: CLI/MCP equivalence tests
- Missing: Phase transition tests
- Missing: Artifact validation tests

---

## 4. Implementation Priorities

### Priority 1: Complete Workspace Engine (Critical)

**Rationale:** The Workspace Engine is the single source of truth. All other work depends on it.

| Task | Estimated | Dependencies |
|------|-----------|--------------|
| Extract verify logic to verify.go | 2h | None |
| Extract init logic to init.go | 2h | None |
| Extract phase logic to phase.go | 2h | lifecycle package |
| Extract artifact logic to artifact.go | 2h | lifecycle package |
| Extract session logic to session.go | 2h | None |
| Extract report logic to report.go | 2h | None |
| Create loader package | 4h | types package |
| Create validator package | 4h | lifecycle package |
| Create state package | 4h | None |
| Integration tests | 8h | All above |

**Total: ~32 hours**

### Priority 2: Methodology Packages (Critical)

**Rationale:** The methodology must be standalone and never depend on runtime.

| Task | Estimated | Dependencies |
|------|-----------|--------------|
| Create authority package | 4h | None |
| Create verification package | 8h | authority package |
| Create phases/* packages | 16h | lifecycle package |
| Create knowledge package | 8h | verification package |
| Tests | 8h | All above |

**Total: ~44 hours**

### Priority 3: Thin Adapter Runtimes (High)

**Rationale:** CLI and MCP must be thin adapters to ensure consistency.

| Task | Estimated | Dependencies |
|------|-----------|--------------|
| Refactor cmd/kdse/ to thin adapter | 16h | Workspace Engine |
| Create cmd/kdse/commands/ package | 8h | Workspace Engine |
| Create cmd/kdse/formatter/ package | 4h | None |
| Create cmd/kdse/parser/ package | 4h | None |
| Refactor cmd/kdse-mcp/ to thin adapter | 16h | Workspace Engine |
| Create cmd/kdse-mcp/server/ package | 8h | Workspace Engine |
| Create cmd/kdse-mcp/tools/ package | 8h | Workspace Engine |
| Integration tests | 8h | All above |

**Total: ~72 hours**

### Priority 4: Runtime Interface Abstraction (High)

**Rationale:** Both runtimes must implement the same interface for consistency.

| Task | Estimated | Dependencies |
|------|-----------|--------------|
| Define Runtime interface | 4h | Workspace Engine |
| CLI implements Runtime | 8h | CLI refactoring |
| MCP implements Runtime | 8h | MCP refactoring |
| Equivalence tests | 16h | Both runtimes |

**Total: ~36 hours**

### Priority 5: Template System (Medium)

**Rationale:** Initialization must be template-based, not hardcoded.

| Task | Estimated | Dependencies |
|------|-----------|--------------|
| Design template structure | 8h | None |
| Implement template loader | 8h | loader package |
| Implement GitHub bootstrap | 16h | Template loader |
| Template versioning | 8h | Bootstrap |
| Tests | 8h | All above |

**Total: ~48 hours**

### Priority 6: Comprehensive Testing (Medium)

**Rationale:** Prevent regressions and ensure equivalence.

| Task | Estimated | Dependencies |
|------|-----------|--------------|
| Phase transition tests | 8h | Workspace Engine |
| Artifact validation tests | 8h | Workspace Engine |
| CLI/MCP equivalence tests | 16h | Both runtimes |
| Integration tests | 16h | All packages |
| Performance benchmarks | 8h | Integration tests |

**Total: ~56 hours**

---

## 5. Implementation Sequence

### Phase 1: Workspace Engine Completion (Weeks 1-2)

```
Week 1:
├── Day 1-2: Extract verify.go, init.go
├── Day 3-4: Extract phase.go, artifact.go
└── Day 5: Review and test

Week 2:
├── Day 1-2: Create loader package
├── Day 3-4: Create validator package
└── Day 5: Create state package + integration
```

### Phase 2: Methodology Packages (Weeks 3-4)

```
Week 3:
├── Day 1-2: Create authority package
├── Day 3-4: Create verification package
└── Day 5: Create phases/knowledge

Week 4:
├── Day 1-2: Create remaining phases/*
├── Day 3-4: Create knowledge package
└── Day 5: Testing
```

### Phase 3: Thin Adapter Runtimes (Weeks 5-6)

```
Week 5:
├── Day 1-3: Refactor cmd/kdse/
├── Day 4-5: Create commands/, formatter/, parser/
└── Ongoing: Test thin adapter pattern

Week 6:
├── Day 1-3: Refactor cmd/kdse-mcp/
├── Day 4-5: Create server/, tools/
└── Ongoing: Integration testing
```

### Phase 4: Runtime Abstraction (Week 7)

```
Day 1-2: Define Runtime interface
Day 3-4: CLI implements Runtime
Day 5: MCP implements Runtime
```

### Phase 5: Template System (Weeks 8-9)

```
Week 8:
├── Day 1-2: Design template structure
├── Day 3-4: Implement template loader
└── Day 5: Basic tests

Week 9:
├── Day 1-3: GitHub bootstrap
├── Day 4: Template versioning
└── Day 5: Full integration
```

### Phase 6: Testing and Polish (Week 10)

```
Day 1-2: Comprehensive tests
Day 3-4: Equivalence tests
Day 5: Performance benchmarks
```

---

## 6. Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing functionality | High | High | Phased rollout, migration guide |
| CLI/MCP behavior divergence | Medium | High | Equivalence tests, shared interface |
| Performance regression | Low | Medium | Benchmarks, optimization pass |
| Template versioning conflicts | Medium | Medium | Careful versioning strategy |
| Insufficient test coverage | Medium | High | Mandatory test coverage gates |
| Circular dependencies | Medium | High | Dependency graph enforcement |
| Migration complexity | High | Medium | Comprehensive migration guide |

---

## 7. Success Criteria

The implementation is complete when:

| # | Criterion | Verification Method |
|---|-----------|---------------------|
| 1 | Workspace Engine owns all state | Package inspection |
| 2 | Methodology packages are independent | Dependency graph analysis |
| 3 | CLI is thin adapter | No business logic imports |
| 4 | MCP is thin adapter | No business logic imports |
| 5 | Runtime interface is consistent | Equivalence tests pass |
| 6 | Template bootstrap works | Integration tests |
| 7 | All phases are enforced | Phase transition tests |
| 8 | Verification gate works | Pre-action verification tests |
| 9 | CLI/MCP produce identical workspaces | Equivalence tests |
| 10 | No regressions | Full test suite passes |

---

## 8. Immediate Next Steps

### Week 1: Start with Workspace Engine

1. Review existing engine.go
2. Extract verify.go from engine.go
3. Extract init.go from engine.go
4. Write tests for extracted code
5. Create loader package

### Pre-Work Checklist

- [ ] Review WORKSPACE_ENGINE.md thoroughly
- [ ] Review METHODOLOGY.md for package structure
- [ ] Review existing internal/workspace/ code
- [ ] Set up CI/CD for new packages
- [ ] Create dependency graph for validation

---

## 9. Document Relationships

```
IMPLEMENTATION_PLAN.md
    │
    ├── References:
    │   ├── PRINCIPLES.md
    │   ├── RUNTIME_ARCHITECTURE.md
    │   ├── WORKSPACE_ENGINE.md
    │   ├── METHODOLOGY.md
    │   ├── CLI_RUNTIME.md
    │   └── MCP_RUNTIME.md
    │
    ├── Referenced By:
    │   └── KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md
    │
    └── Related:
        ├── CURRENT_ARCHITECTURE.md
        └── MIGRATION.md
```

---

## 10. Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-07-17 | Initial implementation plan | KDSE Engineering |

---

*This document is actionable. Each task should have associated GitHub issues for tracking.*
