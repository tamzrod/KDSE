# KDSE SLIM ARCHITECTURE
## The Minimal Implementation That Preserves the Methodology

**Date:** 2026-07-15  
**Method:** Aggressive simplification with Doctor diagnosis as authority  
**Goal:** Simplest architecture that preserves constitutional principles  

---

# PREAMBLE: THE CONSTITUTIONAL BOUNDARY

## What Must Be Preserved (Non-Negotiable)

| Principle | Description |
|-----------|-------------|
| **Knowledge Precedes Architecture** | Every architectural decision traces to knowledge |
| **Authority Flows Downward** | Implementation cannot contradict Architecture; Architecture cannot contradict Knowledge |
| **Evidence Supports Knowledge** | Knowledge is backed by evidence; Evidence ≠ Authority |
| **Verification Is Continuous** | Alignment must be continuously checked |
| **Artifacts Are Living** | Engineering artifacts evolve; not created once and forgotten |
| **Traceability Is Absolute** | Every decision traces to authorized knowledge |

## What Can Be Simplified (Everything Else)

Every concept outside the constitutional boundary must justify its cost.

---

# PART I: CONCEPTS REMOVED

## 1.1 Accidental Complexity From MCP

### REMOVED: STRICT Mode

**What it was:** A mode that blocked LLM from acting outside Work Orders

**Why removed:** 
- Temporary scaffolding for LLM limitations
- Anti-methodology (blocking ≠ enforcement)
- Creates false sense of safety

**What replaces it:**
- Work Orders provide guidance, not barriers
- LLM learns KDSE principles, not rules to bypass
- Runtime enforces through decisions, not blocking

---

### REMOVED: Phase Transitions Map

**What it was:** A static map defining valid phase-to-phase transitions

```
PhaseIdle → PhaseProblem
PhaseProblem → PhaseKnowledge
...
```

**Why removed:**
- Implementation detail exposed as methodology
- Prevents organic workflow
- Creates artificial boundaries

**What replaces it:**
- Decision Engine evaluates current state
- Decisions emerge from evidence, not rules
- Progress is tracked, not mandated

---

### REMOVED: Phase-Based Confidence

**What it was:** Float-based confidence (0.0-1.0) that gated phase transitions

```
PhaseConfidenceThreshold = map[Phase]float64{
    PhaseProblem:        0.6,
    PhaseKnowledge:      0.7,
    ...
}
```

**Why removed:**
- Two different concepts conflated (operational vs. quality)
- Float precision is false accuracy
- Creates gaming opportunity

**What replaces it:**
- Evidence Strength (ternary)
- State Indicator (lifecycle)
- No gating, only tracking

---

### REMOVED: session_status Tool

**What it was:** A separate tool returning detailed session status

**Why removed:**
- Duplicates information available through `status`
- Implementation detail exposed
- Adds complexity without value

**What replaces it:**
- `status` returns complete session information
- One tool, complete information

---

### REMOVED: collect Tool

**What it was:** A standalone tool for artifact collection

**Why removed:**
- Collect is part of derivation, not separate
- Creating a separate phase implies it's optional
- Evidence gathering is continuous, not episodic

**What replaces it:**
- Evidence gathering is embedded in Derivation activity
- Not a standalone tool
- Continuously happens as part of Knowledge work

---

### REMOVED: migrate Tool

**What it was:** A tool for migrating legacy KDSE directories

**Why removed:**
- Legacy compatibility concern
- Not methodology, only migration
- Should be one-time event, not ongoing tool

**What replaces it:**
- Migration documentation (one-time)
- Migration script (run once)
- Remove from ongoing MCP tools

---

## 1.2 Over-Engineered Concepts

### REMOVED: 9 Lifecycle States

**What it was:** Proposed, Experimental, Draft, Reviewed, Approved, Reference, Canonical, Superseded, Deprecated, Archived

**Why removed:**
- 9 states exceed operational need
- Most states never used
- Tracking burden exceeds value

**What replaces it:**
- 4 states: DRAFT, REVIEWED, APPROVED, ARCHIVED
- Capture governance essence
- Maintainable

---

### REMOVED: 5 Derivation Stages

**What it was:** Reference Analysis → Knowledge Derivation → Evidence Correlation → Knowledge Validation → Approved Knowledge

**Why removed:**
- Academic elegance ≠ operational usability
- Creates excessive artifacts and processes
- Engineers will skip or fake

**What replaces it:**
- 3 stages: GATHER → DERIVE → VALIDATE
- Capture derivation principles
- Simpler to follow

---

### REMOVED: 5-Star Evidence Strength

**What it was:** ★★★★★ scale with complex criteria for each level

**Why removed:**
- 5 levels of precision is false accuracy
- Rating burden exceeds value
- Gaming risk (inflation)

**What replaces it:**
- 3 levels: STRONG (●●●), MODERATE (●●), WEAK (●)
- Capture evidence quality
- Simple to assign

---

### REMOVED: Multiple Confidence Systems

**What it was:** Float-based Confidence (workflow), Evidence Strength (knowledge), State Indicator (artifacts)

**Why removed:**
- Three different "confidence" concepts
- Users confused about which to use
- Maintenance overhead

**What replaces it:**
- One concept: Evidence Strength (●●●/●●/●)
- Applied consistently
- Clear interface

---

## 1.3 Duplicate Concepts

### REMOVED: Audit (as standalone concept)

**What it was:** A phase/step called "Audit"

**Why removed:**
- Verification is the concept; "Audit" is a label
- Duplicate of continuous verification
- Creates "audit phase" misconception

**What replaces it:**
- "Check" as the activity
- Verification is continuous
- No standalone "audit phase"

---

### REMOVED: Validation (as standalone concept)

**What it was:** A stage in derivation called "Knowledge Validation"

**Why removed:**
- Overlaps with Verification
- Creates duplicate activity
- Unclear scope

**What replaces it:**
- VALIDATE as third derivation stage
- Validates the derivation quality
- Clear purpose

---

### REMOVED: Proof (as standalone concept)

**What it was:** "Verification Proof" as separate artifact type

**Why removed:**
- All artifacts are proof of something
- Creates unnecessary type
- Verification results are just artifacts

**What replaces it:**
- Verification results are artifacts
- Stored in verification/ directory
- Same lifecycle as other artifacts

---

### REMOVED: Reference Artifacts (as separate concept)

**What it was:** "Reference Artifacts" as distinct from "Evidence"

**Why removed:**
- Reference Artifacts and Evidence are the same thing
- Two names for one concept
- Confusion about distinction

**What replaces it:**
- "Evidence Sources" as single concept
- Evidence that backs knowledge claims
- Clear: sources → evidence → knowledge

---

### REMOVED: ADRs (as separate concept)

**What it was:** Architecture Decision Records as distinct from Architecture

**Why removed:**
- ADRs are how Architecture decisions are documented
- Separate concept creates unnecessary taxonomy
- Architecture IS the collection of decisions

**What replaces it:**
- "Architecture Decisions" as single concept
- ADRs are the document format
- No separate ADR concept

---

### REMOVED: Stewardship (as separate concept)

**What it was:** Formal steward roles with transfer procedures

**Why removed:**
- Over-engineered for most teams
- Creates bureaucracy
- "Owner" is sufficient

**What replaces it:**
- Simple "Owner" attribute on artifacts
- One person responsible
- No formal transfer process

---

## 1.4 Unnecessary Directories

### REMOVED: operational/ Directory

**What it was:** Operational knowledge artifacts

**Why removed:**
- Implementation artifacts, not KDSE concern
- Operational knowledge belongs in implementation/
- No special treatment needed

---

### REMOVED: developmental/ Directory

**What it was:** Development experience artifacts

**Why removed:**
- Implementation artifacts, not KDSE concern
- Development experience is just knowledge
- No special treatment needed

---

### REMOVED: cache/ Directory

**What it was:** Cached analysis results

**Why removed:**
- Cache is implementation concern
- Not methodology
- Should not be visible in workspace

---

### REMOVED: normalized/ Directory

**What it was:** Normalized documentation

**Why removed:**
- Normalization is runtime concern
- Not methodology
- Should not be visible in workspace

---

### REMOVED: sessions/ Directory

**What it was:** Session management artifacts

**Why removed:**
- Session is runtime state
- Not artifact
- Should not be workspace concern

---

### REMOVED: confidence/ Directory

**What it was:** Confidence metrics

**Why removed:**
- Confidence is property of knowledge
- Not separate artifact type
- No separate directory needed

---

# PART II: CONCEPTS MERGED

## 2.1 Merged: Evidence Sources

**From:**
- Reference Artifacts
- Evidence

**To:**
- **Evidence Sources** (single concept)

**Definition:**
> Evidence Sources are information from which Knowledge is derived. They provide evidence that backs knowledge claims.

**Format:**
```yaml
evidence_sources:
  - type: "code"
    location: "src/auth/*.go"
    finding: "Password hashing uses bcrypt"
  - type: "document"
    location: "docs/requirements.md"
    finding: "Security requirements specify bcrypt"
```

---

## 2.2 Merged: Knowledge Derivation

**From:**
- Reference Analysis
- Knowledge Derivation
- Evidence Correlation
- Knowledge Validation

**To:**
- **Derivation** (single activity with 3 stages)

**Stages:**

```
GATHER → DERIVE → VALIDATE
```

| Stage | Purpose | Output |
|-------|---------|--------|
| GATHER | Find and classify evidence sources | Evidence inventory |
| DERIVE | Transform evidence into knowledge | Knowledge statements |
| VALIDATE | Check derivation quality | Validated knowledge |

---

## 2.3 Merged: Verification

**From:**
- Audit
- Verification
- Proof

**To:**
- **Check** (single activity)

**Definition:**
> Check is the continuous activity of confirming alignment between layers. Lower layers trace to higher layers.

**Types:**

| Check | Confirms |
|-------|----------|
| K→A Check | Architecture traces to Knowledge |
| A→I Check | Implementation traces to Architecture |
| E→K Check | Evidence backs Knowledge claims |

---

## 2.4 Merged: Architecture Decisions

**From:**
- Architecture
- Architecture Decision Records (ADRs)

**To:**
- **Architecture Decisions** (single concept)

**Definition:**
> Architecture Decisions are structural decisions that derive from Knowledge and authorize Implementation.

**Format:**
```yaml
decision:
  id: "ARCH-001"
  title: "Service-Based Architecture"
  traces_to: ["KNOW-001", "KNOW-002"]
  rationale: "Independent evolution of control modes"
```

---

## 2.5 Merged: Artifact Owner

**From:**
- Stewardship
- Owner
- Knowledge Steward, Architecture Steward, etc.

**To:**
- **Owner** (single attribute)

**Definition:**
> Every artifact has one Owner responsible for its quality and evolution.

---

## 2.6 Merged: Progress

**From:**
- Phases
- Stages

**To:**
- **Progress** (single concept)

**Definition:**
> Progress indicates advancement through engineering work. Phases are organizational markers, not authority gates.

**Markers:**
```
PROBLEM → KNOWLEDGE → FOUNDATION → ARCHITECTURE → IMPLEMENTATION
```

**Principle:** Progress markers indicate where we are; they do not gate what we do.

---

# PART III: CONCEPTS SIMPLIFIED

## 3.1 Simplified: Artifact Lifecycle

**From:** 9 states  
**To:** 4 states  

```
DRAFT → REVIEWED → APPROVED → ARCHIVED
```

| State | Meaning | Authority |
|-------|---------|-----------|
| DRAFT | Initial creation | None |
| REVIEWED | Peer reviewed | Partial |
| APPROVED | Authoritative | Full |
| ARCHIVED | Superseded | Historical |

**Simplification Rationale:**
- 4 states capture governance essence
- More states create tracking burden without value
- States can be added if needed, not mandated

---

## 3.2 Simplified: Evidence Strength

**From:** ★★★★★ scale  
**To:** ●●●/●●/● scale  

```
●●● (STRONG)    - Multiple independent sources confirm
●● (MODERATE)  - Some corroboration exists  
● (WEAK)        - Single source or inference only
```

**Simplification Rationale:**
- 3 levels capture evidence quality
- More levels create false precision
- Simple to assign, simple to understand

---

## 3.3 Simplified: Work Orders

**From:** Complex structure with required_work, blocked_actions, completion_criteria  
**To:** Simple structure with intent and boundaries  

```yaml
work_order:
  objective: "Define authentication service architecture"
  traces_to: ["KNOW-001"]
  deliverable: "authentication/ARCHITECTURE.md"
  check: "Service boundaries align with KNOWLEDGE claims"
```

**Simplification Rationale:**
- Intent is clear
- Boundaries are implicit (can't contradict upstream)
- No need to list what's blocked (principle-based)

---

## 3.4 Simplified: Phases

**From:** 8 phases with complex transitions  
**To:** 5 progress markers  

```
PROBLEM → KNOWLEDGE → FOUNDATION → ARCHITECTURE → IMPLEMENTATION
```

**Simplification Rationale:**
- Clear progression through engineering work
- Not authority gates, just markers
- Natural sequence through Knowledge → Architecture → Implementation

---

## 3.5 Simplified: Decision Engine

**From:** 5-step cycle with complex logic  
**To:** 3-step cycle  

```
STATE → GAPS → DECIDE
```

| Step | Question |
|------|----------|
| STATE | What exists? What traces to what? |
| GAPS | What is missing? What contradicts? |
| DECIDE | What should happen next? |

**Simplification Rationale:**
- 3 steps capture decision essence
- Simpler to implement
- Clearer to understand

---

# PART IV: RUNTIME RESPONSIBILITIES

## 4.1 Core Responsibilities (MUST HAVE)

| Responsibility | Description | Why Essential |
|----------------|-------------|--------------|
| **Artifact Registry** | Track all artifacts, their types, states, owners | Without this, nothing is managed |
| **Traceability Graph** | Track what traces to what (K→A→I) | Without this, authority cannot be enforced |
| **Evidence Store** | Store evidence sources and their strength | Without this, knowledge has no basis |
| **Decision Engine** | Evaluate state, identify gaps, decide next action | Without this, no guidance |
| **Derivation Orchestrator** | Guide GATHER → DERIVE → VALIDATE process | Without this, knowledge isn't derived |

## 4.2 Secondary Responsibilities (SHOULD HAVE)

| Responsibility | Description | Why Secondary |
|----------------|-------------|--------------|
| **Check Engine** | Run continuous K→A, A→I, E→K checks | Can be manual initially |
| **Report Generator** | Produce human-readable status | UX enhancement |
| **Workspace Manager** | Manage .kdse/ directory structure | Implementation detail |

## 4.3 Not Runtime Responsibilities (REMOVED)

| Removed | Reason |
|---------|--------|
| Phase enforcement | Not methodology, only workflow |
| Confidence gating | Not evidence-based |
| STRICT mode | Anti-methodology |
| Session tracking | Implementation detail |
| Cache management | Not methodology |

---

# PART V: MCP RESPONSIBILITIES

## 5.1 What MCP SHOULD BE

**Definition:** MCP is a transport protocol for Runtime interaction.

**Principle:** MCP should expose engineering capabilities, not implementation details.

## 5.2 Tools (After Simplification)

### Reduced Tool Set

| Tool | Purpose | Merged From |
|------|---------|-------------|
| **init** | Initialize workspace | initialize |
| **status** | Get current state | status, session_status |
| **decide** | Get next decision | execute |
| **derive** | Run derivation cycle | (was spread across tools) |
| **check** | Run verification checks | audit, collect |
| **help** | Get guidance | help |

**Removed:**
- collect (merged into derive)
- foundation (not a tool, it's an artifact)
- migrate (one-time event, not ongoing)
- session_status (merged into status)

## 5.3 Tool Definitions

### init
```
PURPOSE: Initialize KDSE workspace
INPUT: None (uses current directory)
OUTPUT: Workspace ready at .kdse/
```

### status
```
PURPOSE: Get current engineering state
INPUT: None
OUTPUT: 
  - What artifacts exist
  - Their lifecycle states
  - Traceability graph status
  - Evidence strength
  - Current progress marker
```

### decide
```
PURPOSE: Get next engineering decision
INPUT: Optional context (what you're working on)
OUTPUT:
  - Current objective
  - What should trace to what
  - What gaps exist
  - Next action recommendation
```

### derive
```
PURPOSE: Run knowledge derivation cycle
INPUT: Evidence sources to process
OUTPUT:
  - Knowledge statements
  - Evidence citations
  - Strength assignments
```

### check
```
PURPOSE: Verify alignment
INPUT: Layer to check (K→A, A→I, E→K)
OUTPUT:
  - Traceability status
  - Gaps identified
  - Recommendations
```

### help
```
PURPOSE: Get KDSE guidance
INPUT: Question or topic
OUTPUT: Relevant KDSE principles and guidance
```

## 5.4 What MCP Should NOT Expose

| Not Exposed | Reason |
|-------------|--------|
| Phase transitions | Implementation detail |
| Confidence thresholds | Not methodology |
| Internal state structures | Implementation detail |
| STRICT mode | Anti-methodology |
| Session management | Implementation detail |

---

# PART VI: CLI RESPONSIBILITIES

## 6.1 Commands (After Simplification)

### Reduced Command Set

| Command | Purpose | Merged From |
|---------|---------|-------------|
| **kdse init** | Initialize workspace | initialize |
| **kdse status** | Show current state | status |
| **kdse derive** | Run derivation | collect, derive |
| **kdse check** | Run verification | audit |
| **kdse help** | Get guidance | help |

**Removed:**
- session_status (merged into status)
- migrate (one-time, not command)
- foundation (not a command)
- collect (merged into derive)

## 6.2 Command Definitions

### kdse init
```
kdse init [--project <name>]

Initialize KDSE workspace in current directory.
Creates .kdse/ with minimal structure.
```

### kdse status
```
kdse status [--verbose]

Show current engineering state:
- Artifact count by type
- Lifecycle state summary
- Traceability completeness
- Evidence coverage
- Current progress
```

### kdse derive
```
kdse derive [--evidence <path>]

Run knowledge derivation:
- GATHER: Find evidence sources
- DERIVE: Create knowledge statements  
- VALIDATE: Assign strength

Outputs to .kdse/knowledge/
```

### kdse check
```
kdse check [--layer <K|A|I>]

Verify alignment:
- K→A: Architecture traces to Knowledge
- A→I: Implementation traces to Architecture
- E→K: Evidence backs Knowledge

Outputs to .kdse/verification/
```

### kdse help
```
kdse help [--topic <name>]

Show KDSE principles and guidance.
Topics: principles, artifacts, derivation, verification
```

---

# PART VII: MINIMAL DIRECTORY STRUCTURE

## 7.1 Target Structure

```
.kdse/
├── knowledge/           # Knowledge artifacts (highest authority)
│   ├── *.md            # Knowledge statements
│   └── metadata.yaml    # Evidence citations, strength
│
├── architecture/        # Architecture decisions (traces to knowledge)
│   ├── *.md            # Architecture decisions
│   └── metadata.yaml    # Traceability (traces to KNOWLEDGE)
│
├── implementation/     # Code (traces to architecture)
│   └── (project code)
│
├── verification/       # Verification artifacts
│   ├── checks/         # Check results
│   └── reports/        # Verification reports
│
├── evidence/            # Evidence sources
│   └── sources/        # Evidence source inventory
│
└── .state/
    └── registry.yaml   # Artifact registry (what exists)
```

## 7.2 Directory Principles

| Principle | Application |
|-----------|------------|
| **Organize by artifact type** | Not by phase, not by implementation |
| **Flat is better than nested** | Fewer directories, clear purpose |
| **Empty directories are removed** | Don't create until needed |
| **Authority reflected in order** | Knowledge is first, Implementation is last |

## 7.3 What Was Removed

| Removed | Reason |
|---------|--------|
| foundation/ | Foundation artifacts are in knowledge/ |
| context/ | Context is derived, not stored separately |
| artifacts/ | Evidence is in evidence/ |
| runtime/ | Runtime state is internal |
| sessions/ | Session is internal |
| confidence/ | Confidence is property of knowledge |
| operational/ | Not KDSE concern |
| developmental/ | Not KDSE concern |
| cache/ | Not methodology |
| normalized/ | Not methodology |
| reports/ | Merged into verification/ |

---

# PART VIII: MINIMAL RUNTIME ARCHITECTURE

## 8.1 Core Components

```
┌─────────────────────────────────────────────────────────────────────┐
│                          KDSE RUNTIME                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              ARTIFACT REGISTRY                                │   │
│  │  • What artifacts exist                                      │   │
│  │  • What is each artifact's type?                             │   │
│  │  • What is each artifact's state?                            │   │
│  │  • Who is the owner?                                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              TRACEABILITY GRAPH                              │   │
│  │  • K→A traces                                              │   │
│  │  • A→I traces                                              │   │
│  │  • E→K citations                                           │   │
│  │  • Gap detection                                            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              EVIDENCE STORE                                  │   │
│  │  • Evidence source inventory                                │   │
│  │  • Evidence → Knowledge citations                           │   │
│  │  • Evidence Strength (●●●/●●/●)                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              DECISION ENGINE                                 │   │
│  │  1. STATE: What exists, what traces to what?               │   │
│  │  2. GAPS: What is missing, what contradicts?                │   │
│  │  3. DECIDE: What should happen next?                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              DERIVATION ORCHESTRATOR                         │   │
│  │  GATHER → DERIVE → VALIDATE                                 │   │
│  │  • Guide evidence gathering                                  │   │
│  │  • Support knowledge derivation                              │   │
│  │  • Validate and assign strength                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 8.2 Component Responsibilities

### Artifact Registry

**Responsibilities:**
- Track artifacts by type
- Track lifecycle state
- Track owner
- Provide queries

**Not responsibilities:**
- Enforcement (Decision Engine)
- Derivation (Derivation Orchestrator)

---

### Traceability Graph

**Responsibilities:**
- Store trace relationships
- Detect gaps
- Detect contradictions
- Validate chains

**Not responsibilities:**
- Artifact storage (Artifact Registry)
- Evidence strength (Evidence Store)

---

### Evidence Store

**Responsibilities:**
- Inventory evidence sources
- Link evidence to knowledge
- Track evidence strength

**Not responsibilities:**
- Artifact management (Artifact Registry)
- Traceability (Traceability Graph)

---

### Decision Engine

**Responsibilities:**
- Evaluate current state
- Identify gaps
- Produce decisions
- Generate guidance

**Not responsibilities:**
- Artifact storage (Artifact Registry)
- Derivation (Derivation Orchestrator)

---

### Derivation Orchestrator

**Responsibilities:**
- Guide GATHER stage
- Guide DERIVE stage
- Guide VALIDATE stage
- Assign evidence strength

**Not responsibilities:**
- State evaluation (Decision Engine)
- Traceability (Traceability Graph)

---

## 8.3 What Was Merged

| Merged From | Merged To |
|-------------|-----------|
| State Manager + Workspace Manager | Artifact Registry |
| Decision Engine (complex) | Decision Engine (simplified) |
| Verification Engine + Audit Engine | Check (part of Decision Engine) |

---

# PART IX: FINAL CONCEPT MAP

## 9.1 The Slim KDSE Concept Set

```
CONSTITUTIONAL (Non-Negotiable)
├── Knowledge
├── Architecture
├── Implementation
├── Evidence Sources
├── Traceability
└── Verification

ESSENTIAL (Required for KDSE)
├── Artifact Registry
├── Decision Engine
├── Derivation Orchestrator
└── Evidence Store

ORGANIZATIONAL (For clarity)
├── Progress Markers (PROBLEM → KNOWLEDGE → FOUNDATION → ARCHITECTURE → IMPLEMENTATION)
├── Lifecycle States (DRAFT → REVIEWED → APPROVED → ARCHIVED)
└── Evidence Strength (●●●/●●/●)

TRANSPORT (Implementation)
├── MCP Tools (init, status, decide, derive, check, help)
└── CLI Commands (kdse init, kdse status, kdse derive, kdse check, kdse help)
```

## 9.2 Concept Relationships

```
Evidence Sources
       │
       │ GATHER
       ▼
Knowledge Statements ←── Evidence Strength
       │
       │ TRACE
       ▼
Architecture Decisions ←── Traceability (traces to KNOWLEDGE)
       │
       │ TRACE
       ▼
Implementation ←── Traceability (traces to ARCHITECTURE)
       │
       │ CHECK
       ▼
Verification ←── Continuous alignment confirmation

All artifacts:
├── Have Owner (one person)
├── Have Lifecycle State (DRAFT → REVIEWED → APPROVED → ARCHIVED)
└── Are tracked in Artifact Registry
```

## 9.3 The Complete Slim KDSE

```
┌─────────────────────────────────────────────────────────────────────┐
│                         KDSE SLIM                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  KNOWLEDGE                                                         │
│  └── Evidence Sources → Derive → Knowledge → Strength (●●●/●●/●)    │
│                                                                     │
│  ARCHITECTURE                                                       │
│  └── Knowledge → Derive → Architecture → Traces to KNOWLEDGE       │
│                                                                     │
│  IMPLEMENTATION                                                     │
│  └── Architecture → Realize → Implementation → Traces to ARCH      │
│                                                                     │
│  VERIFICATION                                                       │
│  └── Continuous checks: K→A, A→I, E→K                             │
│                                                                     │
│  ARTIFACT REGISTRY                                                  │
│  └── All artifacts tracked with: type, state, owner, traces        │
│                                                                     │
│  DECISION ENGINE                                                    │
│  └── STATE → GAPS → DECIDE (3-step cycle)                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

# PART X: MIGRATION PLAN

## Phase 1: Document (Week 1)

**Goal:** Document the target state

**Actions:**
1. Write KDSE Slim Architecture document
2. Update KDSE principles (remove over-engineered concepts)
3. Create simplified templates

**Deliverables:**
- KDSE_SLIM_ARCHITECTURE.md
- Updated KDSE principles
- Simplified templates (knowledge.md, architecture.md, verification.md)

---

## Phase 2: Runtime Refactor (Weeks 2-4)

**Goal:** Simplify runtime to essential components

**Actions:**
1. Remove Phase Transitions map
2. Remove STRICT mode
3. Remove Phase Confidence thresholds
4. Implement simplified Decision Engine (STATE → GAPS → DECIDE)
5. Merge State Manager + Workspace Manager → Artifact Registry
6. Simplify Derivation Orchestrator to 3 stages

**Deliverables:**
- Refactored runtime with 5 core components
- No phase enforcement
- Evidence strength (●●●/●●/●) instead of float

---

## Phase 3: MCP Refactor (Weeks 5-6)

**Goal:** Simplify MCP to essential tools

**Actions:**
1. Remove: collect, migrate, session_status, foundation
2. Merge: execute → decide, audit → check, initialize → init
3. Simplify tool responses (no complex state, just essential info)
4. Remove STRICT mode exposure

**Deliverables:**
- 6 tools: init, status, decide, derive, check, help
- No phase enforcement tools
- Clean tool interfaces

---

## Phase 4: CLI Refactor (Weeks 7-8)

**Goal:** Simplify CLI to essential commands

**Actions:**
1. Remove: session_status, migrate, foundation
2. Merge: collect → derive, audit → check
3. Simplify output (progress markers, not phases)

**Deliverables:**
- 5 commands: init, status, derive, check, help
- Clear progress indicators
- Human-readable output

---

## Phase 5: Directory Migration (Weeks 9-10)

**Goal:** Migrate to minimal directory structure

**Actions:**
1. Create new directories (if needed): knowledge/, architecture/, evidence/, verification/
2. Migrate existing artifacts to new structure
3. Remove: foundation/, context/, artifacts/, runtime/, sessions/, confidence/, operational/, developmental/, cache/, normalized/
4. Ensure traceability is maintained

**Deliverables:**
- Minimal directory structure
- All artifacts properly organized
- Traceability preserved

---

## Phase 6: Cleanup (Weeks 11-12)

**Goal:** Remove all accidental complexity

**Actions:**
1. Remove legacy compatibility code
2. Remove duplicate concepts
3. Remove unused directories
4. Update documentation
5. Test simplified workflow

**Deliverables:**
- Clean codebase
- Updated documentation
- Working simplified system

---

# PART XI: SUCCESS CRITERIA

## Definition of Success

| Criterion | Measure |
|-----------|--------|
| **Easier to understand** | New engineer can explain KDSE in 5 minutes |
| **Easier to teach** | KDSE can be taught in 1 day workshop |
| **Easier to implement** | Runtime code reduced by 50%+ |
| **Harder to misuse** | Misuse produces obvious errors, not silent failures |
| **Constitutional principles preserved** | All 6 non-negotiable principles intact |

## Verification

| Check | Method |
|-------|--------|
| Can explain Knowledge → Architecture → Implementation trace? | Teaching exercise |
| Can trace any artifact back to knowledge? | Traceability audit |
| Is evidence strength assigned to all knowledge? | Evidence coverage check |
| Is verification continuous? | Check frequency audit |
| Are artifacts living (updated, not static)? | Artifact history review |

---

# APPENDIX: BEFORE AND AFTER

## Concept Count Reduction

| Category | Before | After | Reduction |
|----------|--------|-------|-----------|
| Directories | 13 | 6 | 54% |
| MCP Tools | 9 | 6 | 33% |
| CLI Commands | ~15 | 5 | 67% |
| Lifecycle States | 9 | 4 | 56% |
| Derivation Stages | 5 | 3 | 40% |
| Evidence Levels | 5 | 3 | 40% |
| Confidence Systems | 3 | 1 | 67% |

## Concept Elimination

| Eliminated | Reason |
|------------|--------|
| STRICT Mode | Anti-methodology |
| Phase Transitions | Implementation detail |
| Phase Confidence | Wrong concept |
| session_status | Duplication |
| collect | Part of derive |
| migrate | One-time event |
| 9 Lifecycle States | Over-engineered |
| 5-Star Evidence | False precision |
| Reference Artifacts | Duplicate of Evidence |
| ADRs | Duplicate of Architecture |
| Stewardship | Over-engineered |
| Multiple Confidence | Confusing |

## Net Result

```
BEFORE: Complex, over-engineered, easy to misuse
AFTER:  Simple, essential, hard to misuse

The same methodology.
Smaller footprint.
Clearer purpose.
```

---

*This document defines the KDSE Slim Architecture. Simplicity serves the methodology.*
