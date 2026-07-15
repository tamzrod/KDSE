# KDSE Project Runtime

**Document Version:** 1.0  
**Type:** Normative Runtime Specification  
**Effective Date:** 2026-07-15

---

## Purpose

This document defines the KDSE Project Runtime architecture—redesigning `kdse init` from first principles to create a complete working engineering environment. The .kdse directory is treated as the project's engineering operating system, not a documentation folder.

---

## 1. Global Runtime

### 1.1 Definition

The **Global Runtime** is the installed KDSE system that owns the master template. It is installed on the developer's machine and serves as the source of truth for all project templates.

### 1.2 Location

```
~/.kdse/
    ├── runtime/              # KDSE system files
    ├── template/             # Master template for projects
    └── config/               # Global configuration
```

### 1.3 Global Runtime Responsibilities

| Responsibility | Description |
|---------------|-------------|
| Template Ownership | Maintains the canonical project template |
| Version Management | Tracks KDSE version across projects |
| Global Configuration | Stores user preferences and defaults |
| Reference Knowledge | Houses reusable engineering knowledge |

### 1.4 Installation

The Global Runtime is installed once per machine:

```bash
# Install global runtime
kdse install

# Verify installation
kdse --version
```

### 1.5 Template Synchronization

```
Global Runtime Template          Project .kdse/
        │                               │
        │        kdse init            │
        │  ────────────────────────►   │
        │        (copies template)     │
        │                               │
        │        kdse update          │
        │  ────────────────────────►   │
        │   (syncs knowledge updates) │
        │                               │
```

---

## 2. Project Runtime

### 2.1 Definition

The **Project Runtime** is the `.kdse/` directory within a project repository. It is the project's engineering operating environment—self-describing, immediately usable, and never empty.

### 2.2 Location

```
project/
├── .git/                    # Version control
├── src/                     # Project source
├── .kdse/                   # Project Runtime (THIS IS THE ENGINE)
│   ├── runtime/             # Execution state
│   ├── foundation/          # Project-specific foundations
│   ├── knowledge/           # Project knowledge
│   ├── laboratory/          # Methodology validation
│   ├── evidence/            # Engineering evidence
│   ├── traceability/        # Decision tracking
│   ├── references/           # Domain references
│   ├── reports/             # Generated reports
│   ├── config/              # Project configuration
│   └── state/               # Session state
```

### 2.3 Comparison with Git

| Aspect | Git | KDSE |
|--------|-----|------|
| Initializes | Working repository | Working engineering environment |
| Creates | `.git/` with commits | `.kdse/` with knowledge |
| Empty Start | No—has HEAD, refs | No—has baseline knowledge |
| Immediately Usable | Yes | Yes |

### 2.4 Constitutional Rule

**A KDSE project must be self-describing.**

A newly cloned repository containing `.kdse/` must allow another engineer or AI to understand:

- Why decisions were made
- What knowledge supports them
- What assumptions exist
- What evidence exists
- What standards were adopted
- What experience has been gained

**No engineering knowledge may exist only inside an LLM session.**

---

## 3. Template Architecture

### 3.1 Template Definition

The **Template** is the master copy of the project runtime structure stored in the Global Runtime. It is the source from which all project runtimes are initialized.

### 3.2 Template Structure

```
~/.kdse/template/
    ├── runtime/                     # Execution infrastructure
    │   ├── state/                   # State management
    │   ├── logs/                    # Execution logs
    │   └── temp/                    # Temporary workspace
    │
    ├── foundation/                  # Baseline foundations
    │   ├── README.md               # Foundation overview
    │   ├── principles/             # Engineering principles
    │   ├── requirements/           # Requirements engineering
    │   ├── traceability/           # Traceability standards
    │   └── documentation/          # Doc standards
    │
    ├── knowledge/                   # Core knowledge base
    │   ├── general/                # General engineering
    │   │   ├── engineering-principles.md
    │   │   ├── requirements-engineering.md
    │   │   ├── traceability.md
    │   │   └── documentation-standards.md
    │   ├── operational/            # Operational knowledge
    │   │   ├── knowledge-derivation.md
    │   │   ├── foundation-management.md
    │   │   └── verification-workflow.md
    │   └── developmental/          # Developmental knowledge
    │       ├── architecture-patterns.md
    │       ├── documentation-patterns.md
    │       └── common-practices.md
    │
    ├── laboratory/                  # Methodology validation
    │   ├── README.md               # Laboratory overview
    │   ├── protocols/              # Validation protocols
    │   └── templates/              # Lab templates
    │
    ├── evidence/                    # Evidence collection
    │   ├── README.md               # Evidence guide
    │   ├── collection/             # Collection templates
    │   └── classification/         # Evidence taxonomy
    │
    ├── traceability/                 # Decision tracking
    │   ├── README.md               # Traceability guide
    │   ├── decision-log/           # Decision records
    │   └── requirement-map/       # Requirement mapping
    │
    ├── references/                   # Domain references
    │   ├── README.md               # Reference guide
    │   └── domains/               # Domain knowledge
    │
    ├── reports/                     # Report templates
    │   ├── session-template.md
    │   ├── audit-template.md
    │   └── finding-template.md
    │
    ├── config/                      # Configuration templates
    │   ├── manifest.yaml           # Environment manifest
    │   └── preferences.yaml        # User preferences
    │
    ├── state/                       # State templates
    │   └── session-state.yaml
    │
    └── artifacts/                   # Artifact inventory
        └── base-inventory.json
```

### 3.3 Template Contents

#### 3.3.1 General Knowledge

Baseline engineering knowledge applicable to all projects:

| Document | Purpose |
|----------|---------|
| engineering-principles.md | Core engineering principles |
| requirements-engineering.md | Requirements gathering standards |
| traceability.md | Traceability requirements |
| documentation-standards.md | Documentation guidelines |

#### 3.3.2 Operational Knowledge

Knowledge for running KDSE sessions:

| Document | Purpose |
|----------|---------|
| knowledge-derivation.md | How to derive new knowledge |
| foundation-management.md | Managing project foundations |
| verification-workflow.md | Verification procedures |

#### 3.3.3 Developmental Knowledge

Pattern-based knowledge for engineering work:

| Document | Purpose |
|----------|---------|
| architecture-patterns.md | Common architectural patterns |
| documentation-patterns.md | Documentation templates |
| common-practices.md | Engineering best practices |

### 3.4 Template Characteristics

| Characteristic | Value |
|---------------|-------|
| Content | Baseline engineering knowledge |
| State | Ready to use, not empty |
| Version | Synced with KDSE version |
| Customizable | Yes—project adds its own knowledge |

---

## 4. Runtime Initialization Sequence

### 4.1 First Principles

`kdse init` must behave like `git init`:

| Aspect | Git | KDSE |
|--------|-----|------|
| Creates | Complete working repository | Complete working environment |
| Copies | Default files, HEAD | Template contents |
| Immediately Usable | Yes | Yes |
| No Empty State | Has commits | Has knowledge |

### 4.2 Initialization Command

```bash
kdse init [OPTIONS]

OPTIONS:
  --template NAME     Template name (default: default)
  --no-gitignore      Skip .gitignore update
  --minimal           Initialize with minimal template
  --force            Overwrite existing .kdse/
```

### 4.3 Initialization Sequence

```
┌─────────────────────────────────────────────────────────────────────┐
│                    KDSE INIT SEQUENCE                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  PHASE 1: DISCOVER                                                   │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 1.1 Check current directory                                  │   │
│  │ 1.2 Detect existing .kdse/                                  │   │
│  │ 1.3 Detect legacy directories                               │   │
│  │ 1.4 Check git repository                                    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 2: LOCATE TEMPLATE                                            │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 2.1 Find global runtime (~/.kdse/)                          │   │
│  │ 2.2 Locate template directory                               │   │
│  │ 2.3 Validate template integrity                            │   │
│  │ 2.4 Load template manifest                                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 3: CREATE RUNTIME                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 3.1 Create .kdse/ root directory                           │   │
│  │ 3.2 Create all subdirectories                              │   │
│  │ 3.3 Set appropriate permissions                            │   │
│  │ 3.4 Initialize state files                                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 4: COPY TEMPLATE                                              │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 4.1 Copy foundation/ contents                               │   │
│  │ 4.2 Copy knowledge/ contents                               │   │
│  │ 4.3 Copy laboratory/ contents                              │   │
│  │ 4.4 Copy evidence/ contents                                │   │
│  │ 4.5 Copy traceability/ contents                            │   │
│  │ 4.6 Copy references/ contents                              │   │
│  │ 4.7 Copy reports/ templates                                │   │
│  │ 4.8 Copy config/ templates                                  │   │
│  │ 4.9 Copy runtime/ structure                                │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 5: INITIALIZE STATE                                          │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 5.1 Create project manifest                                 │   │
│  │ 5.2 Initialize session state                               │   │
│  │ 5.3 Create project digest                                   │   │
│  │ 5.4 Initialize confidence scores                           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 6: INITIALIZE PROJECT KNOWLEDGE                              │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 6.1 Create project-specific directories                    │   │
│  │ 6.2 Initialize knowledge inventory                         │   │
│  │ 6.3 Create domain knowledge stubs                         │   │
│  │ 6.4 Initialize reference index                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 7: INITIALIZE TRACEABILITY                                   │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 7.1 Create decision log                                     │   │
│  │ 7.2 Initialize requirement map                             │   │
│  │ 7.3 Create assumption registry                             │   │
│  │ 7.4 Initialize evidence index                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 8: INITIALIZE LABORATORY                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 8.1 Copy laboratory protocols                              │   │
│  │ 8.2 Create validation workspace                             │   │
│  │ 8.3 Initialize lab results directory                       │   │
│  │ 8.4 Create lab report template                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 9: INITIALIZE EVIDENCE                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 9.1 Create evidence collection directory                    │   │
│  │ 9.2 Initialize evidence taxonomy                            │   │
│  │ 9.3 Create evidence templates                               │   │
│  │ 9.4 Initialize evidence index                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  PHASE 10: VERIFY INSTALLATION                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 10.1 Verify all directories created                         │   │
│  │ 10.2 Verify template contents copied                        │   │
│  │ 10.3 Verify state files initialized                        │   │
│  │ 10.4 Verify configuration valid                             │   │
│  │ 10.5 Generate verification report                           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 4.4 Verification Output

```
KDSE Initialization Complete
============================

Version:    1.0
Template:   default
Directory:  /project/.kdse/
Created:    2026-07-15T14:08:53Z

Runtime Structure:
  ✓ foundation/       12 files
  ✓ knowledge/       15 files
  ✓ laboratory/      5 files
  ✓ evidence/        4 files
  ✓ traceability/    3 files
  ✓ references/     2 files
  ✓ reports/         3 templates
  ✓ config/          2 files
  ✓ runtime/         3 directories

State:
  ✓ Project manifest created
  ✓ Session state initialized
  ✓ Project digest created
  ✓ Confidence scores initialized

Status: ✓ READY

The project is now a KDSE-enabled engineering environment.
Run 'kdse status' to view project state.
Run 'kdse execute' to start engineering.
```

---

## 5. Project Digest Lifecycle

### 5.1 Definition

The **Project Digest** is the living record of all engineering knowledge within a project. It captures every piece of knowledge consumed, created, refined, validated, or lost during engineering work.

### 5.2 Constitutional Requirement

Every engineering activity either:
- **Consumes knowledge** — Uses existing knowledge
- **Creates knowledge** — Produces new insights
- **Refines knowledge** — Improves existing knowledge
- **Validates knowledge** — Verifies knowledge correctness

**Nothing disappears.**

### 5.3 Digest Structure

```yaml
# .kdse/state/project-digest.yaml
project-digest:
  version: "1.0"
  created: "2026-07-15T14:08:53Z"
  updated: "2026-07-15T14:08:53Z"
  
  knowledge:
    baseline:
      count: 15
      last-updated: "2026-07-15"
    project:
      count: 0
      last-updated: null
      
  consumption:
    sessions: 0
    entries: []
    
  creation:
    sessions: 0
    entries: []
    
  refinement:
    sessions: 0
    entries: []
    
  validation:
    sessions: 0
    entries: []
```

### 5.4 Digest Operations

| Operation | Trigger | Effect |
|-----------|---------|--------|
| Consume | Use knowledge in session | Log consumption entry |
| Create | Derive new knowledge | Add to project knowledge |
| Refine | Improve existing knowledge | Update and log change |
| Validate | Verify knowledge | Record validation result |

### 5.5 Digest Lifecycle Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                      PROJECT DIGEST LIFECYCLE                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│    Engineering Session                                               │
│         │                                                           │
│         ├──► CONSUME knowledge                                      │
│         │      │                                                    │
│         │      └──► Log entry in digest                            │
│         │           (what was used, why, outcome)                   │
│         │                                                           │
│         ├──► CREATE knowledge                                       │
│         │      │                                                    │
│         │      └──► Add to project knowledge                        │
│         │           └──► Update digest (count++, entries++)        │
│         │                                                           │
│         ├──► REFIN E knowledge                                      │
│         │      │                                                    │
│         │      └──► Update knowledge file                          │
│         │           └──► Log refinement in digest                  │
│         │                (before, after, reason)                    │
│         │                                                           │
│         └──► VALIDATE knowledge                                     │
│                │                                                    │
│                └──► Record validation                              │
│                     └──► Update digest                              │
│                          (validated, invalidated, needs-work)       │
│                                                                     │
│    Session Complete                                                  │
│         │                                                           │
│         └──► Digest persisted to .kdse/state/project-digest.yaml   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.6 Digest Enforcement

The KDSE runtime enforces the digest lifecycle:

| Enforcement | Mechanism |
|------------|----------|
| Nothing disappears | Digest captures all activities |
| Successful work enriches | Digest count increases |
| Failed validation | Marked in digest |
| Knowledge provenance | Every entry has session reference |

---

## 6. Knowledge Discovery

### 6.1 Definition

**Knowledge Discovery** is the process of finding relevant reference knowledge when an objective becomes available. It searches installed knowledge bases to identify candidate references.

### 6.2 Discovery Trigger

Knowledge Discovery is triggered when:

| Trigger | Description |
|---------|-------------|
| New objective | A new engineering objective is introduced |
| Objective evolution | An existing objective changes scope |
| Knowledge gap | Session identifies missing knowledge |
| Reference request | Engineer explicitly requests references |

### 6.3 Discovery Process

```
┌─────────────────────────────────────────────────────────────────────┐
│                      KNOWLEDGE DISCOVERY                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  OBJECTIVE: "Build an Inventory Management System"                   │
│                                                                     │
│         │                                                           │
│         ▼                                                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 1. PARSE OBJECTIVE                                          │   │
│  │    Extract key concepts:                                    │   │
│  │    - Inventory                                              │   │
│  │    - Management                                             │   │
│  │    - System                                                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 2. SEARCH GLOBAL KNOWLEDGE                                  │   │
│  │    Query: inventory                                         │   │
│  │    Results:                                                 │   │
│  │    - Inventory Domain (high relevance)                      │   │
│  │    - Stock Movement (medium relevance)                      │   │
│  │    - Barcode Concepts (low relevance)                       │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 3. SEARCH PROJECT KNOWLEDGE                                 │   │
│  │    Query: inventory                                         │   │
│  │    Results: [] (no project knowledge yet)                   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 4. BUILD CANDIDATE LIST                                     │   │
│  │    Candidates:                                               │   │
│  │    - Inventory Domain (global)                              │   │
│  │    - Stock Movement (global)                                │   │
│  │    - Barcode Concepts (global)                              │   │
│  │    - Retail Workflow (contextual)                            │   │
│  │    - GS1 Standards (contextual)                            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 5. RANK CANDIDATES                                          │   │
│  │    1. Inventory Domain (relevance: 0.95)                     │   │
│  │    2. Stock Movement (relevance: 0.85)                      │   │
│  │    3. Retail Workflow (relevance: 0.70)                      │   │
│  │    4. Barcode Concepts (relevance: 0.65)                     │   │
│  │    5. GS1 Standards (relevance: 0.50)                       │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 6. PRESENT ENGINEERING PROPOSALS                            │   │
│  │    Proposed References:                                     │   │
│  │    1. Inventory Domain (HIGH - include)                      │   │
│  │    2. Stock Movement (HIGH - include)                       │   │
│  │    3. Retail Workflow (MEDIUM - consider)                    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 6.4 Engineering Proposals

After Knowledge Discovery, candidate references become **Engineering Proposals**:

```yaml
proposals:
  - reference: "Inventory Domain"
    relevance: 0.95
    status: "proposed"
    action: "include"
  - reference: "Stock Movement"
    relevance: 0.85
    status: "proposed"
    action: "include"
  - reference: "Retail Workflow"
    relevance: 0.70
    status: "proposed"
    action: "consider"
```

### 6.5 Proposal Approval

After approval, proposals become **project knowledge**:

```
Proposal          Approval          Project Knowledge
────────          ────────          ────────────────
Inventory Domain  ──APPROVE──►     Added to knowledge/
Stock Movement    ──APPROVE──►     Added to knowledge/
Retail Workflow   ──DEFER────────►  Stored for later
```

---

## 7. Reverse Pull

### 7.1 Definition

**Reverse Pull** is the process of pulling knowledge from the Global Runtime into a project. Unlike forward installation, Reverse Pull is demand-driven by project needs.

### 7.2 Reverse Pull Triggers

| Trigger | Description |
|---------|-------------|
| Knowledge Discovery | Found candidate in global knowledge |
| Session Request | Session requests specific knowledge |
| Gap Analysis | Identified knowledge gap |
| Reference Import | Explicit import request |

### 7.3 Reverse Pull Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                        REVERSE PULL                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  PROJECT NEEDS                                                      │
│       │                                                             │
│       │  "We need Inventory Domain knowledge"                        │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 1. LOCATE IN GLOBAL RUNTIME                                  │   │
│  │    Path: ~/.kdse/knowledge/domains/inventory-domain/        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 2. FETCH KNOWLEDGE                                           │   │
│  │    Files:                                                    │   │
│  │    - inventory-domain.md                                    │   │
│  │    - stock-levels.md                                        │   │
│  │    - reorder-points.md                                       │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 3. COPY TO PROJECT                                          │   │
│  │    Destination: .kdse/knowledge/domains/inventory-domain/  │   │
│  │    Action: Copy (not move)                                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 4. UPDATE PROJECT DIGEST                                    │   │
│  │    Entry: "Pulled Inventory Domain from global runtime"      │   │
│  │    Count: project.knowledge.count++                          │   │
│  │    Source: global::domains/inventory-domain                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 5. UPDATE TRACEABILITY                                      │   │
│  │    Decision: "Include Inventory Domain reference"            │   │
│  │    Reason: "Required for inventory system"                   │   │
│  │    Evidence: "KDSE Knowledge Discovery"                     │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 7.4 Reverse Pull Commands

```bash
# Pull specific knowledge
kdse pull inventory-domain

# Pull with preview
kdse pull --preview inventory-domain

# Pull all recommended references
kdse pull --recommended

# Pull and approve
kdse pull inventory-domain --approve
```

### 7.5 Pull vs. Push

| Direction | Mechanism | Trigger |
|-----------|-----------|---------|
| Push | Global → Projects | KDSE update |
| Pull | Project → Global | Project need |

---

## 8. Experience Capture

### 8.1 Definition

**Experience Capture** is the process of extracting reusable engineering knowledge from completed work and storing it for future use.

### 8.2 Constitutional Rule

**Engineering experience must accumulate. Never disappear.**

### 8.3 Capture Trigger

After every completed engineering task:

```
┌─────────────────────────────────────────────────────────────────────┐
│                   EXPERIENCE CAPTURE TRIGGER                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Engineering Task Complete                                          │
│         │                                                           │
│         ▼                                                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ CAPTURE QUESTION:                                           │   │
│  │ "Did we discover reusable engineering knowledge?"          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│               ┌──────────────┴──────────────┐                      │
│               │                             │                      │
│              YES                            NO                       │
│               │                             │                      │
│               ▼                             ▼                       │
│  ┌────────────────────────┐    ┌────────────────────────┐          │
│  │ CAPTURE EXPERIENCE     │    │ CLOSE TASK             │          │
│  │ 1. Normalize            │    │ Digest: completed      │          │
│  │ 2. Classify             │    │ No knowledge captured  │          │
│  │ 3. Store                │    └────────────────────────┘          │
│  │ 4. Index                │                                          │
│  │ 5. Update references    │                                          │
│  └────────────────────────┘                                         │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 8.4 Capture Process

```yaml
# Experience Capture Sequence

1. IDENTIFY
   - What worked well?
   - What patterns emerged?
   - What was learned?

2. NORMALIZE
   - Format according to knowledge template
   - Ensure consistent structure
   - Remove project-specific details

3. CLASSIFY
   - General Knowledge?
   - Operational Knowledge?
   - Developmental Knowledge?

4. STORE
   - Project knowledge (specific)
   - Global knowledge (reusable)

5. INDEX
   - Keywords
   - Domain tags
   - Relationship links

6. UPDATE REFERENCES
   - Add to knowledge index
   - Update global runtime (if reusable)
   - Update project digest
```

### 8.5 Knowledge Types for Capture

| Type | Scope | Example |
|------|-------|---------|
| General | Reusable globally | "Always validate inputs" |
| Operational | KDSE workflow | "Discovery improves with context" |
| Developmental | Project-specific | "Inventory uses FIFO" |

### 8.6 Capture Integration

Experience Capture integrates with:

| System | Integration |
|--------|-------------|
| Project Digest | Records capture event |
| Traceability | Links to originating session |
| Global Runtime | Pushes reusable knowledge |
| Laboratory | Validates captured knowledge |

---

## 9. Laboratory Integration

### 9.1 Definition

The **Laboratory** is the KDSE methodology validation environment. Every initialized project contains `.kdse/laboratory/` to validate methodology before production use.

### 9.2 Laboratory Mandate

**Nothing enters the Runtime unless it first survives the Laboratory.**

### 9.3 Laboratory Structure

```
.kdse/laboratory/
├── README.md                    # Laboratory overview
├── protocols/                   # Validation protocols
│   ├── knowledge-validation.md
│   ├── pattern-validation.md
│   └── process-validation.md
├── templates/                   # Lab templates
│   ├── lab-report.md
│   └── experiment-design.md
├── results/                     # Lab results
│   └── {experiment-id}.json
├── scenarios/                   # Test scenarios
│   └── {scenario-id}.md
└── core/                        # Lab engine
    └── laboratory.py
```

### 9.4 Laboratory in Initialization

During `kdse init`, the Laboratory is initialized as follows:

```
┌─────────────────────────────────────────────────────────────────────┐
│               LABORATORY INITIALIZATION                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  kdse init                                                          │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ PHASE 8: INITIALIZE LABORATORY                               │   │
│  │                                                             │   │
│  │ 8.1 Copy laboratory protocols                              │   │
│  │      Source: ~/.kdse/template/laboratory/protocols/        │   │
│  │      Dest:   .kdse/laboratory/protocols/                   │   │
│  │                                                             │   │
│  │ 8.2 Create validation workspace                             │   │
│  │      Directory: .kdse/laboratory/validation/               │   │
│  │      Purpose: Run experiments                               │   │
│  │                                                             │   │
│  │ 8.3 Initialize lab results directory                         │   │
│  │      Directory: .kdse/laboratory/results/                  │   │
│  │      Purpose: Store experiment results                       │   │
│  │                                                             │   │
│  │ 8.4 Create lab report template                              │   │
│  │      File: .kdse/laboratory/templates/lab-report.md         │   │
│  │      Purpose: Standardize lab reports                       │   │
│  │                                                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ LABORATORY READY                                            │   │
│  │                                                             │   │
│  │ Protocols:     3 files                                     │   │
│  │ Templates:      2 files                                    │   │
│  │ Workspace:     validation/ (ready)                          │   │
│  │ Results:       results/ (empty, ready)                       │   │
│  │                                                             │   │
│  │ Run 'kdse lab validate --all' to test protocols            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 9.5 Laboratory Validation Scope

The Laboratory validates:

| Validation | Target | Purpose |
|------------|--------|---------|
| Knowledge | New knowledge | Correctness, consistency |
| Patterns | Engineering patterns | Applicability, completeness |
| Processes | KDSE processes | Effectiveness, efficiency |
| Integration | Cross-component | Coherence, dependencies |

### 9.6 Laboratory Workflow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    LABORATORY WORKFLOW                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  New Knowledge Created                                              │
│         │                                                           │
│         ▼                                                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 1. SUBMIT TO LABORATORY                                     │   │
│  │    Command: kdse lab submit --knowledge <file>             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 2. RUN VALIDATION PROTOCOLS                                 │   │
│  │    - Syntax validation                                      │   │
│  │    - Semantic validation                                    │   │
│  │    - Pattern matching                                       │   │
│  │    - Cross-reference validation                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 3. GENERATE LAB REPORT                                      │   │
│  │    - Validation results                                     │   │
│  │    - Findings                                               │   │
│  │    - Recommendations                                        │   │
│  │    - Decision: APPROVE / REJECT / REVISE                   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 4. PROCESS DECISION                                         │   │
│  │    APPROVE  → Add to knowledge base                         │   │
│  │    REJECT   → Return to author                              │   │
│  │    REVISE   → Return with feedback                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 5. UPDATE TRACEABILITY                                      │   │
│  │    - Lab results archived                                    │   │
│  │    - Decision logged                                         │   │
│  │    - Evidence linked                                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 10. Migration Plan

### 10.1 Purpose

The Migration Plan defines how to transition existing KDSE projects from the old implementation to the redesigned Project Runtime.

### 10.2 Migration Types

| Type | Scope | When |
|------|-------|------|
| Directory Migration | Legacy → .kdse/ | Always |
| Template Migration | Old template → New template | Major version |
| Knowledge Migration | Ad-hoc → Structured | On upgrade |
| State Migration | Session → Project Digest | On upgrade |

### 10.3 Directory Migration

```
┌─────────────────────────────────────────────────────────────────────┐
│                    DIRECTORY MIGRATION                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  BEFORE (Legacy Structure)                                          │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ project/                                                    │   │
│  │ ├── foundation/           ← Legacy                          │   │
│  │ ├── knowledge/            ← Legacy                          │   │
│  │ ├── context/              ← Legacy                          │   │
│  │ └── artifacts/            ← Legacy                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│                          ↓ kdse migrate                             │
│                                                                     │
│  AFTER (KDSE Runtime)                                                │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ project/                                                    │   │
│  │ ├── .kdse/               ← Unified Runtime                  │   │
│  │ │   ├── foundation/       ← Migrated                        │   │
│  │ │   ├── knowledge/        ← Migrated                        │   │
│  │ │   ├── context/          ← Migrated                        │   │
│  │ │   ├── artifacts/        ← Migrated                        │   │
│  │ │   ├── laboratory/       ← New                             │   │
│  │ │   ├── evidence/         ← New                             │   │
│  │ │   ├── traceability/     ← New                             │   │
│  │ │   ├── references/       ← New                             │   │
│  │ │   └── ...                                                 │   │
│  │ └── foundation/           ← Removed (now in .kdse/)         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 10.4 Migration Commands

```bash
# Check migration status
kdse migrate --check

# Perform directory migration
kdse migrate --directories

# Perform full migration
kdse migrate --all

# Perform template update
kdse migrate --template

# Perform knowledge restructuring
kdse migrate --knowledge
```

### 10.5 Migration Sequence

```
┌─────────────────────────────────────────────────────────────────────┐
│                      MIGRATION SEQUENCE                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  kdse migrate                                                       │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 1. DETECT LEGACY STRUCTURE                                   │   │
│  │    Check for: foundation/, knowledge/, context/, artifacts/ │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 2. BACKUP CURRENT STATE                                     │   │
│  │    Create: .kdse/.backup/pre-migration/                    │   │
│  │    Copy all existing content                                │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 3. CREATE .KDSE/ STRUCTURE                                  │   │
│  │    Initialize full runtime structure                        │   │
│  │    Copy template contents                                   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 4. MIGRATE LEGACY CONTENT                                   │   │
│  │    foundation/ → .kdse/foundation/                          │   │
│  │    knowledge/  → .kdse/knowledge/                            │   │
│  │    context/    → .kdse/context/                              │   │
│  │    artifacts/  → .kdse/artifacts/                            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 5. INITIALIZE NEW COMPONENTS                                │   │
│  │    - Create laboratory/ structure                           │   │
│  │    - Create evidence/ structure                             │   │
│  │    - Create traceability/ structure                         │   │
│  │    - Create references/ structure                            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 6. CREATE PROJECT DIGEST                                    │   │
│  │    - Record migration event                                 │   │
│  │    - Import existing knowledge count                        │   │
│  │    - Initialize digest metadata                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 7. UPDATE TRACEABILITY                                     │   │
│  │    - Log migration decision                                  │   │
│  │    - Record source locations                                 │   │
│  │    - Link to backup                                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 8. REMOVE LEGACY DIRECTORIES                                │   │
│  │    After confirmation:                                       │   │
│  │    - Remove foundation/                                     │   │
│  │    - Remove knowledge/                                      │   │
│  │    - Remove context/                                        │   │
│  │    - Remove artifacts/                                      │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 9. VERIFY MIGRATION                                         │   │
│  │    - All content accessible                                  │   │
│  │    - No content lost                                         │   │
│  │    - New structure functional                                │   │
│  │    - Digest updated                                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 10. GENERATE MIGRATION REPORT                                │   │
│  │     - Files migrated                                         │   │
│  │     - Files created                                          │   │
│  │     - Warnings                                                │   │
│  │     - Recommendations                                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 10.6 Migration Report Format

```yaml
# Migration Report: 2026-07-15
migration:
  id: "MIG-2026-07-15-001"
  timestamp: "2026-07-15T14:08:53Z"
  from-version: "0.x"
  to-version: "1.0"
  
  backup:
    location: ".kdse/.backup/pre-migration/"
    size: "2.3 MB"
    
  directories:
    migrated:
      - foundation/ (15 files)
      - knowledge/ (23 files)
      - context/ (4 files)
      - artifacts/ (12 files)
    created:
      - laboratory/
      - evidence/
      - traceability/
      - references/
      
  knowledge:
    baseline-count: 15
    project-count: 12
    total-count: 27
    
  warnings:
    - "Some context references may need updating"
    
  status: "COMPLETE"
  
  next-steps:
    - "Run 'kdse status' to verify"
    - "Run 'kdse lab validate' to test laboratory"
    - "Commit migration to version control"
```

---

## Summary

The KDSE Project Runtime redesign transforms `kdse init` from a directory creator into a complete engineering environment initializer:

| Component | Before | After |
|-----------|--------|-------|
| `kdse init` | Creates empty `.kdse/` | Copies full template |
| Project state | Empty, needs setup | Ready to engineer |
| Knowledge | Must add manually | Baseline included |
| Laboratory | Not present | Initialized |
| Evidence | Not present | Initialized |
| Traceability | Not present | Initialized |
| References | Not present | Initialized |

### Key Principles

1. **Git-like behavior**: `kdse init` creates a complete working environment, not an empty folder
2. **Template-driven**: The Global Runtime owns the template; projects copy from it
3. **Nothing disappears**: The Project Digest captures all engineering activities
4. **Demand-driven**: Reverse Pull pulls knowledge when projects need it
5. **Experience accumulates**: Every session captures reusable knowledge
6. **Laboratory validates**: Nothing enters the Runtime without lab validation
7. **Self-describing**: Every project explains itself to new engineers or AI

### Migration Path

Existing projects migrate using:
```bash
kdse migrate --all
```

This ensures continuity while adopting the new architecture.

---

*This document defines the KDSE Project Runtime architecture. It is normative for KDSE-enabled repositories.*
