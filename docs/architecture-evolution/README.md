# KDSE Architecture Evolution

**Document Version:** 1.0  
**Effective Date:** 2026-07-17

---

## Purpose

This directory contains documentation for the KDSE Architecture Evolution initiative. The evolution transforms KDSE from its current state into a Runtime-Centric Engineering Platform.

---

## Evolution Index

| Document | Description | Status |
|----------|-------------|--------|
| [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) | Runtime-Centric Architecture | Current |
| [KAE-002](KAE-002-CLI-MCP-RUNTIME-SEPARATION.md) | CLI/MCP Runtime Separation | Planned |
| [KAE-003](KAE-003-RUNTIME-BOOTSTRAP.md) | Runtime Bootstrap | Planned |
| [KAE-004](KAE-004-EVIDENCE-FIRST-RUNTIME.md) | Evidence-First Runtime | Planned |

---

## Architecture Decision Records

| Document | Decision | Status |
|----------|----------|--------|
| [ADR-001](../decisions/ADR-001-RUNTIME-IS-THE-AUTHORITY.md) | Runtime is the Authority | Accepted |
| [ADR-002](../decisions/ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md) | One Methodology, Multiple Runtimes | Accepted |
| [ADR-003](../decisions/ADR-003-WORKSPACE-ENGINE.md) | Workspace Engine | Accepted |

---

## Evolution Phases

### Phase 1: Architecture Documentation
- [x] Create architecture documents
- [x] Create evolution documents
- [x] Create ADRs

### Phase 2: Repository Refactor
- [ ] Restructure cmd/ directory
- [ ] Restructure internal/ directory
- [ ] Create new package boundaries

### Phase 3: Workspace Engine
- [ ] Create engine package
- [ ] Implement state ownership
- [ ] Implement verification

### Phase 4: Runtime Abstraction
- [ ] Define Runtime interface
- [ ] Implement CLI runtime
- [ ] Implement MCP runtime

### Phase 5: Runtime Bootstrap
- [ ] Remove hardcoded initialization
- [ ] Implement template-based bootstrap
- [ ] Implement verification-first initialization

### Phase 6-15: Additional phases
- [ ] See [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) for details

---

## Quick Start

### Understanding the Architecture

1. Read [PRINCIPLES.md](../architecture/PRINCIPLES.md) to understand core principles
2. Read [CURRENT_ARCHITECTURE.md](../architecture/CURRENT_ARCHITECTURE.md) to understand current state
3. Read [RUNTIME_ARCHITECTURE.md](../architecture/RUNTIME_ARCHITECTURE.md) to understand target state
4. Read [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) for evolution plan

### Key Concepts

| Concept | Description |
|---------|-------------|
| Runtime | The .kdse directory that defines a KDSE project |
| Workspace Engine | The component that owns all project state |
| Methodology | The engineering rules that govern all KDSE projects |
| Thin Adapter | Runtimes (CLI/MCP) that only translate, no business logic |

### Core Principles

```
P-001: Evidence First
P-002: Runtime is the Authority
P-003: One Methodology
P-004: Multiple Runtimes
P-005: Runtime Independence
P-006: Filesystem is Evidence
```

---

## Contributing

### Before Making Changes

1. Review the principles in [PRINCIPLES.md](../architecture/PRINCIPLES.md)
2. Check if an ADR exists for the change
3. If not, create an ADR using [ADR template](../decisions/TEMPLATE.md)

### Architecture Decision Process

1. Identify the problem or opportunity
2. Document the current state
3. Explore possible solutions
4. Evaluate trade-offs
5. Make a decision
6. Document the decision in an ADR
7. Implement the decision
8. Update relevant documentation

---

## Document Relationships

```
docs/
├── architecture-evolution/
│   ├── README.md (this file)
│   ├── KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md
│   ├── KAE-002-CLI-MCP-RUNTIME-SEPARATION.md
│   ├── KAE-003-RUNTIME-BOOTSTRAP.md
│   └── KAE-004-EVIDENCE-FIRST-RUNTIME.md
│
├── architecture/
│   ├── PRINCIPLES.md
│   ├── CURRENT_ARCHITECTURE.md
│   ├── RUNTIME_ARCHITECTURE.md
│   ├── METHODOLOGY.md
│   ├── WORKSPACE_ENGINE.md
│   ├── CLI_RUNTIME.md
│   └── MCP_RUNTIME.md
│
└── decisions/
    ├── ADR-001-RUNTIME-IS-THE-AUTHORITY.md
    ├── ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md
    ├── ADR-003-WORKSPACE-ENGINE.md
    └── TEMPLATE.md
```

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-07-17 | Initial architecture evolution documentation | KDSE Team |

---

*This document is informative. For normative guidance, see [PRINCIPLES.md](../architecture/PRINCIPLES.md).*
