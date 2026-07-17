# KDSE Ownership Model

**Version:** 1.0.0
**Status:** AUTHORITATIVE
**Effective Date:** 2026-07-17

## Overview

KDSE operates on a strict ownership model that establishes clear boundaries between the software project and the KDSE runtime. This model ensures that:

1. **KDSE augments software engineering** - it never replaces standard practices
2. **The project remains the foundation** - KDSE provides runtime support
3. **Ownership is explicit** - every artifact has a clear owner
4. **Boundaries are enforced** - artifacts cannot cross ownership lines

## Core Principle

> **"A KDSE-enabled repository must remain a standard software project. Any engineer unfamiliar with KDSE should still be able to understand, build, deploy, and maintain the project without understanding KDSE internals."**

## Four Ownership Domains

### 1. Project Layer (Owned by Software Project)

**Purpose:** Contains all software deliverables and project documentation.

**Location:** Project root (/)

**Contains:**
- `README.md` - Project overview
- `LICENSE` - Legal documentation
- `CHANGELOG.md` - Version history
- `docs/` - Project documentation
  - `architecture/` - Software architecture
  - `api/` - API documentation
  - `deployment/` - Deployment documentation
  - `design/` - Design documentation
- `src/` - Source code
- `tests/` - Test code
- `cmd/` - CLI entrypoints
- `internal/` - Internal packages
- `deploy/` - Deployment configurations
- `templates/` - Project templates
- `examples/` - Usage examples
- `.github/` - CI/CD configuration

**Characteristics:**
- Portable - works without KDSE
- Standard - follows industry conventions
- Maintainable - familiar to any engineer

### 2. Runtime Layer (Owned by KDSE Runtime)

**Purpose:** Contains engineering runtime artifacts and state management.

**Location:** `.kdse/`

**Contains:**
- `runtime/` - Runtime state and configuration
- `sessions/` - Engineering session data
- `state/` - State management
- `cache/` - Cached data
- `reports/` - Engineering reports
- `evidence/` - Engineering evidence
- `traceability/` - Traceability links
- `laboratory/` - Engineering laboratory
- `bootstrap/` - Runtime bootstrap files

**Characteristics:**
- Transient - can be regenerated
- Runtime-specific - requires KDSE runtime
- Non-portable - not part of project deliverable

### 3. Reference Layer (Owned by KDSE Runtime)

**Purpose:** Contains external authoritative references.

**Location:** `.kdse/references/`

**Contains:**
- `modbus/` - Modbus TCP specifications
- `iec61850/` - IEC 61850 standards
- `ieee/` - IEEE standards
- `vendor/` - Vendor manuals
- `rfc/` - RFC documents
- `regulatory/` - Regulatory documents

**Characteristics:**
- External - sourced from authoritative bodies
- Immutable - references don't change
- External knowledge - not project knowledge

### 4. Knowledge Layer (Owned by KDSE Runtime)

**Purpose:** Contains knowledge extracted from references.

**Location:** `.kdse/knowledge/`

**Contains:**
- `modbus/` - Extracted Modbus knowledge
- `iec61850/` - Extracted IEC 61850 knowledge
- Domain-specific subdirectories

**Characteristics:**
- Derived - extracted from references
- Traceable - must maintain reference links
- Internal knowledge - understanding of external references

## Artifact Classification Rules

### Rule 1: Project Artifacts (Project Layer)

**If the artifact describes the software:**
→ Store it inside the project

Examples:
- `README.md` - describes the project
- `docs/architecture/` - describes software architecture
- `docs/api/` - describes the API
- `docs/deployment/` - describes deployment
- `src/` - is the software
- `tests/` - tests the software

### Rule 2: Engineering Process Artifacts (Runtime Layer)

**If the artifact describes the engineering process:**
→ Store it inside `.kdse/`

Examples:
- `sessions/` - engineering sessions
- `reports/` - engineering reports
- `evidence/` - engineering evidence
- `traceability/` - traceability data
- `laboratory/` - engineering laboratory
- `cache/` - runtime cache

### Rule 3: External References (Reference Layer)

**If the artifact is an external authority:**
→ Store it in `.kdse/references/`

Examples:
- IEC 61850 standard documents
- Modbus TCP specifications
- IEEE standards
- Vendor manuals
- RFC documents

### Rule 4: Extracted Knowledge (Knowledge Layer)

**If the artifact is extracted knowledge:**
→ Store it in `.kdse/knowledge/`

Examples:
- Modbus protocol knowledge
- IEC 61850 implementation knowledge
- Domain-specific understanding

## Engineering Lifecycle

```
Reference (External Authority)
        │
        ▼
Knowledge (Extracted Understanding)
        │
        ▼
Architecture (Software Design)
        │
        ▼
Implementation (Software Build)
        │
        ▼
Verification (Testing & Validation)
        │
        ▼
Project Documentation (Produced during engineering)
```

**Note:** Project documentation is produced during engineering but belongs to the project repository. References and extracted knowledge belong to the KDSE runtime.

## Boundary Enforcement

### What KDSE Runtime SHALL NOT Own

- `README.md`
- `docs/`
- `src/`
- `tests/`
- `cmd/`
- `internal/`
- Deployment documentation
- API documentation
- Software architecture documentation
- `LICENSE`
- `CHANGELOG.md`

### What KDSE Runtime SHALL Own

- `runtime/` - Runtime state
- `sessions/` - Engineering sessions
- `state/` - State management
- `cache/` - Runtime cache
- `reports/` - Engineering reports
- `evidence/` - Engineering evidence
- `traceability/` - Traceability
- `references/` - External references
- `knowledge/` - Extracted knowledge
- `laboratory/` - Engineering laboratory
- `bootstrap/` - Bootstrap configuration

## Initialization Requirements

When initializing a KDSE workspace, the following standard project layout MUST be created:

```
Project/
├── README.md        (created if missing)
├── docs/            (created if missing)
│   └── README.md
├── src/             (created if missing)
│   └── README.md
├── tests/           (created if missing)
│   └── README.md
└── .kdse/           (always created)
    ├── runtime/
    ├── sessions/
    ├── state/
    ├── cache/
    ├── reports/
    ├── evidence/
    ├── traceability/
    ├── references/
    ├── knowledge/
    └── laboratory/
```

**Idempotency:** Initialization MUST be idempotent. Existing files MUST NOT be overwritten.

## Implementation Checklist

- [ ] Project layer directories created during initialization
- [ ] Runtime layer directories created during initialization
- [ ] Artifact classification implemented in artifact creation
- [ ] Ownership enforcement in runtime guard
- [ ] Verification report confirms boundary compliance

## Anti-Patterns (Forbidden)

1. ❌ Storing project documentation in `.kdse/`
2. ❌ Storing source code in `.kdse/`
3. ❌ Storing tests in `.kdse/`
4. ❌ Using `.kdse/` as the project workspace
5. ❌ Creating project artifacts in runtime layer
6. ❌ Mixing ownership domains

## Correct Patterns (Required)

1. ✅ Project artifacts in project root
2. ✅ Runtime artifacts in `.kdse/`
3. ✅ References in `.kdse/references/`
4. ✅ Knowledge in `.kdse/knowledge/`
5. ✅ Standard project layout preserved
6. ✅ Ownership boundaries respected

---

**Document Version:** 1.0.0
**Last Updated:** 2026-07-17
**Owner:** KDSE Architecture Team
