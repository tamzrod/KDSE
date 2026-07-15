# KDSE BERSERK RECONSTRUCTION
## First-Principles Analysis of Knowledge-Driven Software Engineering

**Analysis Date:** 2026-07-15  
**Mission:** Destroy all assumptions, re-derive from zero  
**Verdict:** Complete architectural reconstruction  

---

# PREAMBLE: THE ONE CONSTANT

One principle cannot be questioned. It is the bedrock. Everything else is derived.

> **Principle 1 (ABSOLUTE): Knowledge Precedes Architecture**  
> Architecture derives from knowledge. Architecture decisions must be traceable to specific knowledge artifacts. Architecture that cannot be traced to knowledge is not compliant with KDSE.

Every concept in this document must justify its existence against this single principle.

---

# PART I: WHY KDSE EXISTS

## The Fundamental Problem KDSE Solves

Software engineering fails because:

1. **Decisions are made without traceable basis**
2. **Architecture drifts from original intent**
3. **Implementation contradicts knowledge**
4. **No evidence of why decisions were made**
5. **Systems become incomprehensible over time**

KDSE exists to ensure every engineering decision traces to authorized knowledge.

## The Non-Negotiable Goal

**Every code change, every architectural decision, every implementation choice must trace to authorized knowledge artifacts.**

This is the only goal that matters.

---

# PART II: FIRST-PRINCIPLES CONCEPT ANALYSIS

## Concept: KNOWLEDGE

### Question: Why does Knowledge exist?

**Answer:** Because without authoritative knowledge, architecture has no basis, implementation has no guidance, and decisions cannot be traced.

### Question: What KDSE principle requires it?

- **Principle 1:** Architecture derives from knowledge
- **Principle 4:** Knowledge is the longest-lived artifact
- **Principle 5:** Engineering decisions must be traceable
- **Principle 6:** Code realizes knowledge
- **Principle 12:** Knowledge is implementation-independent

### Question: If removed, what principle breaks?

ALL OF THEM. Knowledge is the foundation of the entire methodology.

### Question: What is Knowledge's correct definition?

**ORIGINAL DEFINITION (CORRECT):**
> Engineering Knowledge is authoritative understanding about the problem domain, requirements, constraints, and context. Knowledge is the highest-authority artifact type.

**MCP DEFINITION (WRONG):**
> Knowledge is a phase in a workflow that produces context documents.

### First-Principles Derivation

Knowledge must:

1. **Exist continuously** - Not created in a phase, then forgotten
2. **Be derived, not collected** - "Collection" implies passive gathering; derivation requires reasoning
3. **Have explicit authority** - Authority comes from derivation process, not collection
4. **Be traceable** - Every knowledge statement traces to evidence
5. **Outlive implementation** - Knowledge persists when code is rewritten
6. **Be implementation-independent** - True if implementation can change without breaking knowledge

### VERDICT: KNOWLEDGE MUST SURVIVE

**Status:** KEEP  
**Changes Required:** 
- Remove "phase" designation
- Restore continuous evolution
- Restore 5-stage derivation lifecycle
- Restore Evidence Strength, not Confidence

---

## Concept: EVIDENCE

### Question: Why does Evidence exist?

**Answer:** Because knowledge must be derived from something. Evidence is the raw material from which knowledge is derived.

### Question: What KDSE principle requires it?

- **Principle 11:** Reference Artifacts Support Engineering Knowledge - evidence supports knowledge
- **Principle 13:** Evidence Strengthens but Does Not Authorize - evidence ≠ authority

### Question: If removed, what principle breaks?

**Principle 11** breaks immediately. Knowledge would have no basis.

### First-Principles Derivation

Evidence is:

1. **Raw facts** - Unprocessed information from any source
2. **Not authoritative** - Evidence supports, never authorizes
3. **Multiple sources** - Single evidence is weak; corroboration strengthens
4. **Classifiable by strength** - Evidence Strength reflects corroboration level

### The Evidence → Knowledge Derivation

```
Reference Artifact 1 ─┐
Reference Artifact 2 ─┼──▶ Evidence ──▶ Knowledge (authoritative)
Reference Artifact 3 ─┘
```

**The Critical Distinction:**
- Evidence is **observed** - can be collected
- Knowledge is **derived** - requires reasoning
- Collection ≠ Derivation

### VERDICT: EVIDENCE MUST SURVIVE

**Status:** KEEP  
**Changes Required:**
- Separate Evidence from Knowledge conceptually
- Restore Evidence Strength (★★★★★) as distinct from Confidence
- Clarify: Evidence is gathered; Knowledge is derived

---

## Concept: CONFIDENCE

### Question: Why does Confidence exist?

**Answer:** To answer: "How much can we trust this knowledge?" But ONLY in the sense of evidence strength, not as a gating mechanism.

### Question: What KDSE principle requires it?

- **Principle 13:** Evidence Strengthens but Does Not Authorize

### Question: If removed, what principle breaks?

Nothing. Confidence is not a core principle; it's a supporting concept.

### The Critical Error in Current Implementation

**MCP treats Confidence as:**
- A float between 0.0 and 1.0
- A gating mechanism for phase transitions
- Something the runtime calculates

**KDSE defines Confidence as:**
- Evidence Strength (★★★★★)
- Derived from corroboration across Reference Artifacts
- Supports trust, never authorizes

### First-Principles Derivation

Confidence (Evidence Strength) is:

1. **Qualitative** - ★★★★★ scale, not float
2. **Evidence-based** - Derived from corroboration count
3. **Authority-independent** - Strength ≠ Authority
4. **Knowledge-specific** - Applied to knowledge statements, not phases

### The Dangerous Confusion

```
MCP: "Confidence = 0.85 → proceed to Implementation"
KDSE: "Evidence Strength = ★★★☆☆ → knowledge is moderately supported"
```

**MCP has conflated two completely different concepts:**
- Evidence-based trust (★★★★★)
- Phase gating mechanism (float)

### VERDICT: CONFIDENCE (as Evidence Strength) MUST SURVIVE, MCP CONFIDENCE MUST DIE

**Status:** SPLIT  
**Evidence Strength:** KEEP (as part of Knowledge derivation)  
**Phase Gating:** REMOVE (anti-KDSE, procedural convenience)

---

## Concept: ARTIFACT

### Question: Why does Artifact exist?

**Answer:** Because engineering produces tangible outputs that must be tracked, traced, and governed.

### Question: What KDSE principle requires it?

- **Principle 5:** Engineering decisions must be traceable
- **Principle 6:** Code realizes knowledge
- **Principle 8:** Authority flows downward

### Question: If removed, what principle breaks?

**Principle 5** breaks immediately. Without artifacts, decisions cannot be traced.

### First-Principles Derivation

Artifacts are:

1. **The outputs of engineering work** - Everything produced during engineering
2. **Classified by authority level** - Knowledge > Architecture > Implementation
3. **Tracked through lifecycle states** - Proposed → Approved → Superseded
4. **Governed by stewards** - Someone responsible for quality
5. **Traceable to higher authorities** - Implementation traces to Architecture, which traces to Knowledge

### The MCP Error

**MCP treats artifacts as:**
- Phase outputs
- Documents that exist when created
- Not continuously evolving
- Without lifecycle states

**KDSE defines artifacts as:**
- Living documents that continuously evolve
- Tracked through explicit lifecycle states
- With authority levels that gate their use
- Under governance stewardship

### Artifact Types in KDSE (First-Principles)

1. **Knowledge Artifacts** - Authoritative understanding (highest authority)
2. **Architecture Artifacts** - Structural decisions (derives from Knowledge)
3. **Implementation Artifacts** - Code realizing architecture (lowest authority)
4. **Verification Artifacts** - Proof of alignment (applies authority)

### VERDICT: ARTIFACT MUST SURVIVE, BUT RESTORED PROPERLY

**Status:** KEEP (with corrections)  
**Changes Required:**
- Restore lifecycle states
- Restore authority hierarchy
- Restore stewardship
- Remove "phase output" treatment

---

## Concept: ARCHITECTURE

### Question: Why does Architecture exist?

**Answer:** Because Knowledge (intent) must be translated into structure (architecture) before it can be realized (implementation).

### Question: What KDSE principle requires it?

- **Principle 1:** Knowledge Precedes Architecture
- **Principle 2:** Architecture Precedes Implementation
- **Principle 8:** Authority flows downward

### Question: If removed, what principle breaks?

**Principle 2** breaks immediately.

### First-Principles Derivation

Architecture is:

1. **The translation layer** - Converts knowledge intent into structural decisions
2. **Authored by derivation** - Every architectural decision traces to knowledge
3. **Authorizes implementation** - Implementation cannot proceed without architectural authorization
4. **Technology-aware** - Unlike knowledge, architecture makes technology choices

### The Critical Question: Is Architecture a Phase?

**NO. Architecture is not a phase.**

Architecture is:
- An artifact type (like Knowledge, Implementation)
- Continuous (not bounded by time)
- Derivation-based (not sequential)

**Phases are an implementation convenience, not a KDSE concept.**

### Architecture Activities (Not Phases)

1. **Identify architectural drivers** - Which knowledge requires architectural expression?
2. **Explore approaches** - What structural options exist?
3. **Make decisions** - Which approach best serves the knowledge?
4. **Document decisions** - ADRs, diagrams, specifications
5. **Validate against knowledge** - Does architecture satisfy derived knowledge?

### VERDICT: ARCHITECTURE MUST SURVIVE, BUT NOT AS A PHASE

**Status:** KEEP (as artifact type)  
**Changes Required:**
- Remove "phase" designation
- Restore as continuous derivation activity
- Remove temporal boundaries

---

## Concept: IMPLEMENTATION

### Question: Why does Implementation exist?

**Answer:** Because knowledge must be realized in code.

### Question: What KDSE principle requires it?

- **Principle 2:** Architecture Precedes Implementation
- **Principle 6:** Code realizes knowledge
- **Principle 8:** Authority flows downward

### Question: If removed, what principle breaks?

**Principle 6** breaks immediately.

### First-Principles Derivation

Implementation is:

1. **The realization layer** - Code manifests decisions
2. **Authored by architecture** - Implementation traces to architecture
3. **Verified against architecture** - Implementation must conform
4. **Subject to verification** - Alignment must be confirmed

### The Critical Question: Is Implementation a Phase?

**NO. Implementation is not a phase.**

Implementation is:
- An artifact type (like Knowledge, Architecture)
- Continuous (not bounded by time)
- Realization-based (not sequential)

### Implementation Activities (Not Phases)

1. **Realize architecture** - Write code that conforms to architecture
2. **Maintain traceability** - Every code change traces to architecture
3. **Document deviations** - When implementation differs, document and justify
4. **Verify alignment** - Confirm implementation matches architecture

### VERDICT: IMPLEMENTATION MUST SURVIVE, BUT NOT AS A PHASE

**Status:** KEEP (as artifact type)  
**Changes Required:**
- Remove "phase" designation
- Restore as continuous realization activity
- Remove temporal boundaries

---

## Concept: FOUNDATION

### Question: Why does Foundation exist?

**Answer:** To capture authoritative knowledge that governs engineering decisions.

### Question: What KDSE principle requires it?

- **Principle 1:** Knowledge Precedes Architecture
- **Principle 5:** Engineering decisions must be traceable

### Question: If removed, what principle breaks?

Nothing directly, but without Foundation, there would be no place to capture authoritative knowledge.

### The Critical Question: Is Foundation a Phase?

**NO. Foundation is not a phase.**

**Foundation is a COLLECTION OF ARTIFACTS.**

### MCP Error

MCP treats Foundation as:
- A phase in linear workflow
- 5 documents created once
- Then forgotten

**This is fundamentally wrong.**

### First-Principles Derivation

Foundation is:

1. **A set of living artifacts** - SPEC.md, REQUIREMENTS.md, ASSUMPTIONS.md, CONSTRAINTS.md, GLOSSARY.md
2. **Continuously evolving** - Not created once
3. **Under governance** - Lifecycle states, steward assignments
4. **Authoritative** - These documents authorize downstream decisions

### Foundation Artifacts

1. **SPEC.md** - Project specification (authoritative)
2. **REQUIREMENTS.md** - Functional requirements (authoritative)
3. **ASSUMPTIONS.md** - Documented assumptions (authoritative)
4. **CONSTRAINTS.md** - Documented constraints (authoritative)
5. **GLOSSARY.md** - Domain terminology (authoritative)

### VERDICT: FOUNDATION MUST SURVIVE, BUT AS ARTIFACTS, NOT PHASE

**Status:** KEEP (as artifact collection)  
**Changes Required:**
- Remove "phase" designation
- Restore as living artifacts
- Add lifecycle state tracking

---

## Concept: AUDIT

### Question: Why does Audit exist?

**Answer:** To confirm that implementation aligns with architecture, and architecture aligns with knowledge.

### Question: What KDSE principle requires it?

- **Principle 9:** Verification confirms alignment
- **Principle 5:** Engineering decisions must be traceable

### Question: If removed, what principle breaks?

**Principle 9** breaks.

### The Critical Question: Is Audit a Phase?

**NO. Audit is not a phase.**

**Audit is a CONTINUOUS ACTIVITY.**

### MCP Error

MCP treats Audit as:
- A phase between Foundation and Assessment
- One-time evaluation
- Check if documents exist

**This is fundamentally wrong.**

### First-Principles Derivation

Audit/Verification is:

1. **Continuous** - Not bounded by time
2. **Applies to all artifact types** - Knowledge, Architecture, Implementation
3. **Confirms alignment** - Lower traces to higher
4. **Produces evidence** - Verification artifacts prove alignment

### What Audit Evaluates

1. **Knowledge → Architecture alignment** - Does architecture trace to knowledge?
2. **Architecture → Implementation alignment** - Does implementation trace to architecture?
3. **Knowledge completeness** - Is knowledge sufficient to authorize decisions?
4. **Traceability integrity** - Are all decisions traceable?

### VERDICT: AUDIT MUST SURVIVE, BUT AS CONTINUOUS ACTIVITY

**Status:** KEEP (as continuous evaluation)  
**Changes Required:**
- Remove "phase" designation
- Restore as continuous activity
- Remove temporal boundaries

---

## Concept: SESSION

### Question: Why does Session exist?

**Answer:** To bound engineering work temporally and provide context.

### Question: What KDSE principle requires it?

None directly. Session is an operational concept, not a KDSE principle.

### First-Principles Derivation

Session is:

1. **A bounded context** - Engineering work within a time period
2. **Tracks state** - What artifacts exist, what decisions made
3. **Produces reports** - Documents the work done
4. **Enables continuity** - Sessions can continue previous work

### Is Session a Phase?

**NO. Session is not a phase.**

Session is a container for activities.

### VERDICT: SESSION MUST SURVIVE, BUT NOT AS PHASE

**Status:** KEEP (as container)  
**Changes Required:**
- Remove phase connotations
- Restore as temporal context

---

## Concept: WORKSPACE

### Question: Why does Workspace exist?

**Answer:** To provide a physical location for artifacts.

### Question: What KDSE principle requires it?

**Principle 5** indirectly - decisions must be traceable, which requires artifacts to exist somewhere.

### First-Principles Derivation

Workspace is:

1. **A directory structure** - Where artifacts live
2. **Organized by artifact type** - Not by phase
3. **Versioned** - Tracked like code
4. **Transportable** - Can be moved between systems

### Is Workspace the Repository?

**NO.**

- **Repository**: Where source code lives
- **Workspace**: Where KDSE artifacts live

They are separate.

### The MCP Workspace Error

MCP organizes workspace by:
- Phase output (foundation/, knowledge/, context/)
- Session state (sessions/)
- Runtime artifacts (runtime/)

**Should be organized by:**
- Artifact type (knowledge/, architecture/, implementation/, verification/)
- Authority level

### VERDICT: WORKSPACE MUST SURVIVE, BUT RESTRUCTURE

**Status:** KEEP (with restructuring)  
**Changes Required:**
- Organize by artifact type, not phase
- Reflect authority hierarchy

---

## Concept: PROJECT

### Question: Why does Project exist?

**Answer:** To define the scope of engineering work.

### Question: What KDSE principle requires it?

None directly. Project is an organizational concept.

### First-Principles Derivation

Project is:

1. **A scope definition** - What problem are we solving?
2. **Contains workspace** - KDSE artifacts live in project
3. **May contain repository** - Source code may live here
4. **Has objective** - What outcome is desired?

### Is Project a Phase?

**NO. Project is not a phase.**

Project is a container for artifacts and activities.

### VERDICT: PROJECT MUST SURVIVE, BUT NOT AS PHASE

**Status:** KEEP (as container)  
**Changes Required:**
- Remove phase connotations

---

## Concept: REPOSITORY

### Question: Why does Repository exist?

**Answer:** To store source code.

### Question: What KDSE principle requires it?

**None.** Repository is not a KDSE concept. It's a source control concept.

### First-Principles Derivation

Repository is:

1. **Source code storage** - Where code lives
2. **Version controlled** - Git, etc.
3. **Separate from workspace** - KDSE artifacts ≠ source code
4. **Referenced by workspace** - Artifacts may analyze repository

### The Critical Distinction

```
Repository: Source code (git-managed)
Workspace: KDSE artifacts (.kdse-managed)
```

### VERDICT: REPOSITORY IS NOT A KDSE CONCEPT

**Status:** EXTERNAL CONCEPT (not KDSE scope)

---

## Concept: RUNTIME

### Question: Why does Runtime exist?

**Answer:** To orchestrate engineering activities and enforce methodology.

### Question: What KDSE principle requires it?

**Indirectly all of them.** Without something to enforce methodology, KDSE is just documentation.

### First-Principles Derivation

Runtime is:

1. **The methodology engine** - Enforces KDSE principles
2. **Owns decisions** - Not LLM, not human - Runtime
3. **Tracks state** - What exists, what's authorized, what's missing
4. **Produces outputs** - Reports, recommendations, decisions

### Runtime Responsibilities

1. **Enforce authority hierarchy** - Lower cannot contradict higher
2. **Maintain traceability** - Every decision traces to knowledge
3. **Evaluate evidence** - Assess Evidence Strength
4. **Produce decisions** - What to do next based on state

### Runtime Non-Responsibilities

1. **Writing code** - LLM does this
2. **Making architectural decisions** - Human does this with Runtime guidance
3. **Defining principles** - KDSE Standard does this

### Is Runtime a Phase?

**NO. Runtime is not a phase.**

Runtime is an execution engine.

### VERDICT: RUNTIME MUST SURVIVE AS METHODOLOGY ENGINE

**Status:** KEEP  
**Changes Required:**
- Own decisions, not LLM
- Not phase-based

---

## Concept: EXECUTE

### Question: Why does Execute exist?

**Answer:** To run the Runtime and produce engineering decisions.

### Question: What KDSE principle requires it?

**Indirectly.** Without execution, methodology is dormant.

### The MCP Error

MCP treats Execute as:
- A tool that generates Work Orders
- A state machine that advances phases
- A way to tell LLM what to do

**This is fundamentally wrong.**

### First-Principles Derivation

Execute is:

1. **NOT a phase transition** - Phases don't exist
2. **NOT Work Order generation** - That's procedural convenience
3. **NOT LLM instruction** - That's anti-methodology

Execute is:

1. **A decision cycle** - Evaluate state → produce decision
2. **Driven by evidence** - Not by phase rules
3. **Owned by Runtime** - Not LLM

### The Correct Execute Cycle

```
1. LOAD current state
   - What artifacts exist?
   - What is their authority level?
   - What traces exist?

2. EVALUATE evidence
   - What new evidence exists?
   - How does evidence strengthen knowledge?
   - Are there contradictions?

3. CHECK alignment
   - Does Architecture trace to Knowledge?
   - Does Implementation trace to Architecture?
   - Are there gaps?

4. PRODUCE decision
   - What should happen next?
   - Based on evidence, not phase rules
   - Runtime owns this decision
```

### VERDICT: EXECUTE MUST SURVIVE, BUT COMPLETELY REDESIGNED

**Status:** KEEP (redesign)  
**Changes Required:**
- Not phase-based
- Evidence-driven, not rule-based
- Runtime owns decisions

---

## Concept: COLLECT

### Question: Why does Collect exist?

**Answer:** To gather evidence from Reference Artifacts.

### Question: What KDSE principle requires it?

**Principle 11:** Reference Artifacts Support Engineering Knowledge
**Principle 14:** Repository First

### The MCP Error

MCP treats Collect as:
- A tool that catalogs files
- An optional phase
- A debugging tool

**This is fundamentally wrong.**

### First-Principles Derivation

Collect is:

1. **Evidence gathering** - Finding Reference Artifacts
2. **Part of derivation** - Not separate, not optional
3. **Required for knowledge** - Without evidence, no knowledge

Collect is part of the Knowledge Derivation Lifecycle, not a separate phase.

### Collect Activities

1. **Discover Reference Artifacts** - Find all evidence sources
2. **Catalog artifacts** - What exists, where, provenance
3. **Extract evidence** - Factual statements, assertions
4. **Classify evidence** - By type, strength, relevance

### VERDICT: COLLECT MUST SURVIVE, BUT AS PART OF DERIVATION

**Status:** KEEP (as derivation sub-activity)  
**Changes Required:**
- Not a phase
- Part of Knowledge Derivation
- Not optional

---

## Concept: MCP

### Question: Why does MCP exist?

**Answer:** To provide a transport mechanism for Runtime communication.

### Question: What KDSE principle requires it?

**None.** MCP is not a KDSE concept. It's a transport protocol.

### First-Principles Derivation

MCP is:

1. **A transport layer** - How Runtime communicates
2. **An implementation detail** - Not architectural
3. **One option among many** - Could be CLI, API, GUI, etc.
4. **Transparent** - User should not know MCP exists

### The MCP Error

MCP became:
- The architecture
- The methodology definition
- The authority source

**This is backwards.**

### Correct Relationship

```
KDSE Runtime ──uses──▶ Transport (MCP, CLI, etc.)
                            │
                            ├── MCP Server
                            ├── CLI Interface
                            └── API Interface
```

### VERDICT: MCP IS NOT A KDSE CONCEPT

**Status:** EXTERNAL (implementation detail)  
**Changes Required:**
- Not the architecture
- Just one transport option

---

## Concept: CLI

### Question: Why does CLI exist?

**Answer:** To provide a command-line interface to the Runtime.

### Question: What KDSE principle requires it?

**None.** CLI is not a KDSE concept. It's an interface option.

### First-Principles Derivation

CLI is:

1. **An interface layer** - How user interacts with Runtime
2. **An implementation detail** - Not architectural
3. **One option among many** - Could be GUI, API, MCP, etc.

### VERDICT: CLI IS NOT A KDSE CONCEPT

**Status:** EXTERNAL (interface option)  

---

# PART III: CONCEPTS THAT MUST DIE

## 1. PHASES (as currently defined)

**What Dies:** Problem, Knowledge, Foundation, Audit, Assessment, Architecture, Implementation as sequential phases.

**Why:** 
- No KDSE principle requires phases
- Phases are a procedural convenience
- They create artificial boundaries
- They prevent continuous evolution

**What Replaces:** Continuous activities driven by evidence

---

## 2. WORK ORDERS

**What Dies:** Explicit instructions telling LLM what to do.

**Why:**
- Anti-methodology
- LLM should execute KDSE principles, not instructions
- Work Orders are procedural crutches
- They bypass Runtime ownership

**What Replaces:** Evidence-driven decision engine

---

## 3. STRICT MODE

**What Dies:** Mode that blocks LLM from acting outside Work Orders.

**Why:**
- Anti-methodology
- LLM should understand KDSE, not be blocked
- Blocking is not enforcement
- It's a workaround for bad architecture

**What Replaces:** LLM understanding of KDSE principles

---

## 4. PHASE-BASED CONFIDENCE

**What Dies:** Float-based confidence (0.0-1.0) gating phase transitions.

**Why:**
- No KDSE principle supports this
- Confidence is Evidence Strength (★★★★★)
- Phase gating is procedural convenience
- It prevents evidence-driven decisions

**What Replaces:** Evidence Strength evaluation

---

## 5. FOUNDATION AS PHASE

**What Dies:** Foundation as a workflow step that creates documents.

**Why:**
- Foundation is a collection of artifacts
- Not a phase
- Creates documents once, then forgets
- No continuous evolution

**What Replaces:** Living Foundation artifacts under governance

---

## 6. AUDIT AS PHASE

**What Dies:** Audit as a workflow step between Foundation and Assessment.

**Why:**
- Audit is continuous evaluation
- Not a phase
- Evaluates all artifact types
- Not one-time check

**What Replaces:** Continuous verification activity

---

## 7. COLLECT AS PHASE/TOOL

**What Dies:** Collect as an optional tool.

**Why:**
- Collect is part of derivation
- Not a phase
- Not optional
- Required for Knowledge

**What Replaces:** Collect as derivation sub-activity

---

## 8. SESSION-BASED WORKSPACE

**What Dies:** Workspace organized by session.

**Why:**
- Workspace should reflect artifact hierarchy
- Not session state
- Artifacts exist across sessions
- Session is container, not organizer

**What Replaces:** Artifact-type organization

---

# PART IV: CONCEPTS THAT MUST SPLIT

## 1. CONFIDENCE → Evidence Strength + State Indicator

**Split:**
- **Evidence Strength (★★★★★):** Belongs to Knowledge derivation
- **State Indicator:** Belongs to Artifact lifecycle

**Reason:** MCP conflated two different concepts.

---

## 2. FOUNDATION → Artifacts + Governance

**Split:**
- **Foundation Artifacts:** SPEC.md, REQUIREMENTS.md, etc. (belong in workspace)
- **Foundation Governance:** Lifecycle states, stewards (belong in governance)

**Reason:** Foundation is both an artifact collection AND a governance concept.

---

## 3. VERIFICATION → Audit + Proof

**Split:**
- **Audit:** Evaluation of alignment (continuous activity)
- **Proof:** Evidence that evaluation occurred (artifact)

**Reason:** These are distinct but related.

---

## 4. KNOWLEDGE → Derivation Process + Authoritative Output

**Split:**
- **Knowledge Derivation:** The 5-stage process
- **Knowledge Artifacts:** The authoritative outputs

**Reason:** Process and output are distinct.

---

# PART V: CONCEPTS THAT MUST MERGE

## 1. Reference Artifacts + Evidence → Evidence Sources

**Merge:**
- Reference Artifacts and Evidence are the same concept
- Evidence Sources = raw information with provenance

---

## 2. Architecture Decision Record + Architecture → Architectural Decisions

**Merge:**
- ADRs are how architectural decisions are documented
- Architecture is the collection of decisions
- Merge into single concept: Architectural Decision

---

## 3. Implementation Artifact + Code + Configuration → Implementation

**Merge:**
- Code, configuration, deployment artifacts = Implementation
- Implementation is the artifact type, not individual files

---

# PART VI: NEW ARCHITECTURE

## The First-Principles Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│                     KDSE STANDARD (Normative)                       │
│                                                                     │
│  • Core Principles                                                 │
│  • Chain of Authority                                              │
│  • Artifact Definitions                                            │
│  • Derivation Lifecycle                                            │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Governs
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│                      KDSE RUNTIME (Engine)                          │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              STATE MANAGER                                    │   │
│  │  • Artifact Registry (what exists)                          │   │
│  │  • Authority Tracker (who authorized)                       │   │
│  │  • Traceability Graph (what traces to what)                 │   │
│  │  • Evidence Store (what supports what)                      │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              DECISION ENGINE                                  │   │
│  │  • Evaluate current state                                    │   │
│  │  • Assess evidence sufficiency                               │   │
│  │  • Check alignment (K→A, A→I)                               │   │
│  │  • Produce decisions                                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              DERIVATION ORCHESTRATOR                          │   │
│  │  • Reference Analysis                                        │   │
│  │  • Knowledge Derivation                                      │   │
│  │  • Evidence Correlation                                      │   │
│  │  • Knowledge Validation                                      │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              VERIFICATION ENGINE                              │   │
│  │  • Continuous alignment checking                             │   │
│  │  • Traceability verification                                 │   │
│  │  • Gap identification                                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Orchestrates
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│                        ARTIFACT HIERARCHY                            │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  KNOWLEDGE ARTIFACTS (Highest Authority)                    │   │
│  │  • SPEC.md, REQUIREMENTS.md, ASSUMPTIONS.md, etc.           │   │
│  │  • Evidence Strength: ★★★★★                                  │   │
│  │  • Lifecycle: Proposed → Approved → Reference                │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              │ Authorizes                           │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  ARCHITECTURE ARTIFACTS (Authorizes Implementation)        │   │
│  │  • ADRs, Component Specs, Interface Definitions             │   │
│  │  • Traces to Knowledge                                      │   │
│  │  • Lifecycle: Proposed → Approved                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              │ Realizes                             │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  IMPLEMENTATION ARTIFACTS (Realizes Architecture)          │   │
│  │  • Source code, Configuration, Scripts                       │   │
│  │  • Traces to Architecture                                  │   │
│  │  • Lifecycle: Draft → Implemented → Verified                │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              │ Confirms                             │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  VERIFICATION ARTIFACTS (Confirms Alignment)                │   │
│  │  • Test results, Audit reports, Review findings             │   │
│  │  • Evidence of alignment                                   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## The Continuous Derivation Model

```
┌─────────────────────────────────────────────────────────────────────┐
│                     REFERENCE ARTIFACTS (Evidence)                   │
│              Project Docs, Code, Vendor Docs, Standards              │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Discover & Extract
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        REFERENCE ANALYSIS                            │
│                    Evidence identification                           │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Derive
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     KNOWLEDGE DERIVATION                             │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  1. Concept Analysis                                        │   │
│  │     Identify concepts with engineering implications          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  2. Responsibility Identification                            │   │
│  │     Define what system must do                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  3. Architectural Decision Making                            │   │
│  │     Make structural decisions                                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  4. Architecture Documentation                                │   │
│  │     Document decisions                                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  5. Verification Against Knowledge                           │   │
│  │     Confirm architecture satisfies knowledge                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Authorize
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     KNOWLEDGE ARTIFACTS                             │
│              Authoritative understanding (★★★ to ★★★★★)             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Derive
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     ARCHITECTURE ARTIFACTS                           │
│              Structural decisions that trace to Knowledge             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Realize
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   IMPLEMENTATION ARTIFACTS                          │
│           Code that traces to Architecture (which traces to Knowledge) │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Verify
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   VERIFICATION ARTIFACTS                             │
│              Proof of alignment (continuous activity)               │
└─────────────────────────────────────────────────────────────────────┘
```

---

## The Evidence-Driven Decision Model

```
┌─────────────────────────────────────────────────────────────────────┐
│                         DECISION CYCLE                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. LOAD STATE                                                      │
│     ┌─────────────────────────────────────────────────────────┐     │
│     │  • What artifacts exist?                                │     │
│     │  • What is each artifact's authority level?            │     │
│     │  • What traces exist between artifacts?                 │     │
│     │  • What evidence supports current knowledge?            │     │
│     └─────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  2. EVALUATE EVIDENCE                                                │
│     ┌─────────────────────────────────────────────────────────┐     │
│     │  • Is new evidence available?                            │     │
│     │  • Does evidence strengthen or weaken knowledge?          │     │
│     │  • Are there contradictions to resolve?                   │     │
│     │  • Does Evidence Strength justify decisions?             │     │
│     └─────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  3. CHECK ALIGNMENT                                                  │
│     ┌─────────────────────────────────────────────────────────┐     │
│     │  • Does Architecture trace to Knowledge?                │     │
│     │  • Does Implementation trace to Architecture?            │     │
│     │  • Are there gaps in the traceability chain?            │     │
│     │  • Does implementation contradict architecture?         │     │
│     └─────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  4. IDENTIFY GAPS                                                    │
│     ┌─────────────────────────────────────────────────────────┐     │
│     │  • What knowledge is missing or weak?                   │     │
│     │  • What architecture is missing derivation?              │     │
│     │  • What implementation lacks authorization?             │     │
│     │  • What verification is incomplete?                       │     │
│     └─────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  5. PRODUCE DECISION                                                 │
│     ┌─────────────────────────────────────────────────────────┐     │
│     │  Runtime decides:                                        │     │
│     │  • What activity to perform next                         │     │
│     │  • Based on evidence, not rules                          │     │
│     │  • What authorization is needed                          │     │
│     │  • What risks exist                                      │     │
│     └─────────────────────────────────────────────────────────┘     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

# PART VII: RUNTIME RESPONSIBILITIES

## What Runtime MUST Do

| Responsibility | Description | Required By |
|---------------|-------------|-------------|
| **Maintain Artifact Registry** | Track all artifacts, their types, states, authorities | Traceability |
| **Enforce Authority Hierarchy** | Lower cannot contradict higher | Principle 8 |
| **Orchestrate Derivation** | Run the 5-stage Knowledge Derivation | Principle 1 |
| **Evaluate Evidence** | Assess Evidence Strength (not Confidence) | Principle 13 |
| **Verify Alignment** | Continuous K→A→I alignment checking | Principle 9 |
| **Produce Decisions** | Decide next action based on state | Runtime purpose |
| **Track Traceability** | Maintain the trace graph | Principle 5 |

## What Runtime MUST NOT Do

| Prohibited | Reason |
|-----------|--------|
| **Write code** | LLM does this |
| **Make architectural decisions** | Human does this |
| **Generate Work Orders** | Anti-methodology |
| **Block LLM actions** | Blocks KDSE understanding |
| **Gate on phase transitions** | No phases exist |
| **Calculate float Confidence** | Confidence = Evidence Strength |

---

# PART VIII: MCP RESPONSIBILITIES

## What MCP IS

**MCP is a transport protocol for Runtime communication.**

## What MCP SHOULD Do

| Responsibility | Description |
|---------------|-------------|
| **Receive requests** | Pass user input to Runtime |
| **Transmit decisions** | Pass Runtime decisions to user/LLM |
| **Handle transport** | STDIO, HTTP, etc. |
| **Protocol translation** | Convert between formats |

## What MCP SHOULD NOT Do

| Prohibited | Reason |
|-----------|--------|
| **Define methodology** | That's KDSE Standard |
| **Make decisions** | That's Runtime |
| **Own phases** | Phases don't exist |
| **Generate Work Orders** | Anti-methodology |
| **Enforce blocking** | That's not transport |

---

# PART IX: CLI RESPONSIBILITIES

## What CLI IS

**CLI is a user interface for Runtime interaction.**

## What CLI SHOULD Do

| Responsibility | Description |
|---------------|-------------|
| **Accept commands** | Pass user commands to Runtime |
| **Display output** | Show Runtime decisions |
| **Format reports** | Present information readably |
| **Manage workspace** | Create/manage .kdse/ directory |

## What CLI SHOULD NOT Do

| Prohibited | Reason |
|-----------|--------|
| **Implement methodology** | That's Runtime |
| **Make decisions** | That's Runtime |
| **Own phases** | Phases don't exist |
| **Calculate Confidence** | That's Runtime |

---

# PART X: THE KDSE CONSTITUTION

## Non-Negotiable Principles

These principles cannot be modified, overridden, or bypassed by any implementation:

### 1. KNOWLEDGE IS AUTHORITY
> Architecture derives from knowledge. Every architectural decision must trace to specific knowledge artifacts.

### 2. AUTHORITY FLOWS DOWNWARD
> Implementation cannot contradict Architecture. Architecture cannot contradict Knowledge. Lower layers cannot authorize upper layers.

### 3. EVIDENCE SUPPORTS, NEVER AUTHORIZES
> Evidence strengthens knowledge but does not create authority. Authority derives from derivation process, not evidence quantity.

### 4. TRACEABILITY IS ABSOLUTE
> Every engineering decision must trace to authorized knowledge. Decisions without traceability are not KDSE-compliant.

### 5. KNOWLEDGE IS IMPLEMENTATION-INDEPENDENT
> Knowledge remains valid if implementation is completely rewritten. Knowledge describes purpose, not technology.

### 6. VERIFICATION IS CONTINUOUS
> Alignment between Knowledge → Architecture → Implementation must be continuously verified. Not checked once.

### 7. DERIVATION REQUIRES REASONING
> Knowledge cannot be extracted; it must be derived through analysis and judgment. Collection is not derivation.

### 8. ARTIFACTS ARE LIVING
> Engineering artifacts evolve continuously. They are not created once and forgotten.

### 9. GOVERNANCE IS INTEGRATED
> Artifact lifecycle states, stewardship, and review requirements are part of KDSE, not optional.

### 10. METHODOLOGY IS OWNED BY RUNTIME
> The Runtime owns methodology decisions. Not LLM, not human operator, not transport protocol.

---

# PART XI: MIGRATION STRATEGY

## Phase 0: Acknowledge the Poison

**Current state:** MCP implementation has silently changed KDSE meaning.

**Required actions:**
1. Acknowledge that "Knowledge Collection" ≠ Knowledge Derivation
2. Acknowledge that "Foundation phase" ≠ Foundation artifacts
3. Acknowledge that "Audit phase" ≠ Continuous verification
4. Acknowledge that "Confidence (float)" ≠ Evidence Strength (★★★★★)
5. Acknowledge that "Work Orders" = Anti-methodology

## Phase 1: Separate Concepts

**Goal:** Separate MCP transport from KDSE methodology.

**Required actions:**
1. Remove MCP from KDSE architecture definition
2. Define KDSE as transport-agnostic
3. MCP becomes one implementation, not the architecture
4. CLI becomes one interface, not the methodology

## Phase 2: Restore Derivation

**Goal:** Re-implement the 5-stage Knowledge Derivation Lifecycle.

**Required actions:**
1. Remove "Knowledge Collection" phase
2. Implement Reference Analysis as required activity
3. Implement Evidence Correlation as required activity
4. Implement Knowledge Validation as required activity
5. Restore Evidence Strength (★★★★★) scale

## Phase 3: Restore Artifacts

**Goal:** Make artifacts living entities under governance.

**Required actions:**
1. Add lifecycle states to all artifacts
2. Implement steward assignments
3. Add review and approval workflows
4. Remove "phase output" treatment
5. Make artifacts continuously evolving

## Phase 4: Restore Authority

**Goal:** Enforce the authority hierarchy.

**Required actions:**
1. Implement K→A traceability requirements
2. Implement A→I traceability requirements
3. Block implementation that doesn't trace to architecture
4. Block architecture that doesn't trace to knowledge
5. Remove phase-based gating

## Phase 5: Restore Decision Engine

**Goal:** Replace state machine with evidence-driven decisions.

**Required actions:**
1. Remove PhaseTransitions map
2. Implement evidence evaluation in decision cycle
3. Runtime owns decisions, not Work Orders
4. Remove STRICT mode
5. Remove blocked actions

## Phase 6: Clean Workspace

**Goal:** Organize workspace by artifact type, not phase.

**Required actions:**
1. Replace phase directories with artifact directories
2. Organize by authority level
3. Add lifecycle state tracking
4. Remove session-based organization

---

# PART XII: SUMMARY

## Principles That Survived

1. **Knowledge Precedes Architecture** - ABSOLUTE
2. **Architecture Precedes Implementation** - ABSOLUTE
3. **Authority Flows Downward** - ABSOLUTE
4. **Traceability Is Absolute** - ABSOLUTE
5. **Evidence Supports, Never Authorizes** - ABSOLUTE
6. **Knowledge Is Implementation-Independent** - ABSOLUTE
7. **Derivation Requires Reasoning** - ABSOLUTE
8. **Verification Is Continuous** - ABSOLUTE

## Concepts That Must Die

1. Phases (as currently defined)
2. Work Orders
3. STRICT mode
4. Phase-based Confidence
5. Foundation as phase
6. Audit as phase
7. Collect as phase/tool
8. Session-based workspace

## Concepts That Must Split

1. Confidence → Evidence Strength + State Indicator
2. Foundation → Artifacts + Governance
3. Verification → Audit + Proof
4. Knowledge → Process + Output

## Concepts That Must Merge

1. Reference Artifacts + Evidence → Evidence Sources
2. Architecture + ADRs → Architectural Decisions
3. Implementation artifacts → Implementation

## The New KDSE

```
KDSE = Knowledge-Driven Software Engineering

KDSE = Evidence Sources → Knowledge Derivation → Knowledge Artifacts
                                          ↓
                                   Architecture Artifacts
                                          ↓
                                   Implementation Artifacts
                                          ↓
                                   Verification Artifacts

RUNTIME = Methodology Engine (owns decisions)
MCP = Transport (one option)
CLI = Interface (one option)
```

---

## Final Verdict

**The current MCP implementation is NOT KDSE.**

It is a procedural workflow orchestrator that:
- Uses KDSE terminology without KDSE meaning
- Implements phases where KDSE has continuous activities
- Uses Work Orders where KDSE has evidence-driven decisions
- Blocks LLM where KDSE should teach LLM KDSE principles

**KDSE must be rebuilt from first principles.**

The methodology survives intact. The implementation must be destroyed and rebuilt.

---

*This document represents a first-principles re-derivation of Knowledge-Driven Software Engineering. It destroys the current implementation to preserve the methodology.*
