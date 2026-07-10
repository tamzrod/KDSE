# KDSE Execution Loop

**Document Version:** 1.1  
**Effective Date:** 2026-07-10  
**Change Note:** Added Phase 3.4 (Repository Phase Detection) and Phase 3.5 (Phase-Appropriate Recommendations) to address KDSE-CASE-001 OBS-001

---

## Purpose

The Execution Loop defines the continuous improvement cycle that drives KDSE sessions. Unlike a linear workflow, the Execution Loop is an iterative cycle that repeats until target maturity is reached or session termination.

---

## Loop Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ┌───────────┐                                                       │
│  │ Run KDSE  │◀──────────────────────────────────────────┐        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │   Load    │                                               │        │
│  │ Standards │                                               │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │ Foundation│                                               │        │
│  │Verification│                                              │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │ Repository│                                               │        │
│  │Assessment │                                               │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │Compliance │                                               │        │
│  │   Audit   │                                               │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │ Generate  │                                               │        │
│  │   KDSE    │                                               │        │
│  │   Report  │                                              │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │ Recommend │                                               │        │
│  │   Highest │                                               │        │
│  │Value Action│                                              │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │   Await   │                                               │        │
│  │   Human   │                                               │        │
│  │ Approval  │                                               │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ├─────────────────────────────────────────────────────┤        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │Implement │                                               │        │
│  │ Approved │                                               │        │
│  │   Work   │                                               │        │
│  └─────┬─────┘                                               │        │
│        │                                                     │        │
│        ▼                                                     │        │
│  ┌───────────┐                                               │        │
│  │  Verify   │─────────────────────────────────────────────┘        │
│  │  Results  │                                                   │
│  └─────┬─────┘                                                   │
│        │                                                          │
│        ▼                                                          │
│  ┌───────────┐                                                     │
│  │ Re-run   │                                                     │
│  │Compliance │                                                     │
│  │   Audit   │                                                     │
│  └─────┬─────┘                                                     │
│        │                                                           │
│        ▼                                                           │
│  ┌───────────────────┐                                             │
│  │Repeat until target│                                             │
│  │    maturity       │                                             │
│  │    reached         │                                             │
│  └───────────────────┘                                              │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Loop Phases

### Phase 1: Initialization

**Duration:** Typically < 1 minute

**Purpose:** Establish session context and load required resources

**Activities:**
1. Receive session parameters (repository, target, owner)
2. Load KDSE standards and templates
3. Retrieve repository metadata
4. Initialize session state

**Standards Loaded:**
- KDSE Foundation documents
- Audit templates and scoring criteria
- Session Protocol
- Agent Specification

**Exit Criteria:**
- All resources loaded successfully
- Session state initialized
- Ready for assessment

---

### Phase 2: Assessment

**Duration:** Typically 5-30 minutes

**Purpose:** Evaluate current state through systematic audits

**Activities:**

#### 2.1 Foundation Verification

Verify that KDSE standards are accessible and consistent:
- Check all foundation documents available
- Verify cross-reference integrity
- Confirm terminology consistency
- Validate audit standards present

Reference: [FOUNDATION_AUDIT.md](../audit/FOUNDATION_AUDIT.md) (dimensions 1-5)

#### 2.2 Repository Assessment

Inventory and map repository contents:
- Catalog existing artifacts by type
- Identify artifact locations and formats
- Map relationships between artifacts
- Verify steward assignments
- Assess artifact lifecycle states

#### 2.3 Compliance Audit

Run full KDSE Compliance Audit:
- Evaluate all 7 dimensions
- Collect evidence per dimension
- Score each dimension per AUDIT_SCORING.md
- Identify findings by severity
- Document all gaps

Reference: [COMPLIANCE_AUDIT.md](../audit/COMPLIANCE_AUDIT.md)

**Exit Criteria:**
- All audits complete
- Evidence documented
- Findings categorized
- Baseline score established

---

### Phase 3: Recommendation

**Duration:** Typically 2-5 minutes

**Purpose:** Transform audit results into actionable guidance

**Activities:**

#### 3.1 Generate KDSE Report

Produce standard report per [REPORT_FORMAT.md](REPORT_FORMAT.md):
- Current Status summary
- Audit Summary with dimension scores
- Highest Priority Findings
- Recommended Next Action
- Expected Impact
- Required Human Approval
- Session State

#### 3.2 Identify Highest-Value Action

Analyze findings to determine highest-value next action:

**Value Calculation:**
```
Action Value = Impact × Feasibility × Priority

Where:
- Impact = Expected score improvement
- Feasibility = Probability of successful implementation
- Priority = Severity of addressed finding(s)
```

**Selection Criteria:**
1. Highest calculated value
2. Directly addresses Critical or High findings
3. Enables downstream improvements
4. Within session constraints
5. Human-approvable

#### 3.3 Document Rationale

For each recommendation, document:
- Which finding(s) it addresses
- Why it has highest value
- What evidence supports it
- What alternatives were considered

**Exit Criteria:**
- KDSE Report generated
- Recommendation documented
- Human approval requested

#### 3.4 Repository Phase Detection

Before generating recommendations, determine the current repository phase to ensure recommendations respect the KDSE Chain of Authority.

**Phase Detection Method:**

| Phase Indicators | Repository Phase |
|-------------------|------------------|
| No knowledge artifacts identified | Research |
| Knowledge artifacts exist, no architecture | Knowledge Development |
| Knowledge and architecture artifacts present, no implementation | Architecture |
| Implementation artifacts present, limited verification | Implementation |
| Verification artifacts with test execution evidence present | Verification |
| Ongoing maintenance and evolution activities | Evolution |

**Phase Detection Process:**
1. Inventory artifact types present in repository
2. Assess relative maturity of each artifact type
3. Identify highest-maturity artifact type
4. Map to corresponding repository phase

**Phase Context Recording:**
Record detected phase in session state for:
- Report generation (Section 4)
- Recommendation filtering (Section 3.5)
- Score presentation with phase context

#### 3.5 Phase-Appropriate Recommendations

Filter recommendations to include only phase-appropriate actions that respect the Chain of Authority.

**Phase-Appropriate Action Matrix:**

| Repository Phase | Appropriate Actions | Excluded/Deprioritized Actions |
|------------------|--------------------|------------------------------|
| Research | Discover, analyze, map problem domain | Architecture-specific work |
| Knowledge Development | Create, validate, structure knowledge artifacts | Implementation-specific work |
| Architecture | Create architecture, derive from knowledge | Implementation recommendations |
| Implementation | Create implementation, maintain traceability | Verification-only recommendations |
| Verification | Verify alignment, execute tests, document results | New implementation work |
| Evolution | Evolve artifacts, maintain relevance, retire obsolete | Fundamental restructuring |

**Recommendation Filtering Rules:**

1. **Primary Filter:** Recommendations must target dimensions applicable to current phase
2. **Secondary Filter:** Recommendations must not violate phase prerequisites
3. **Tertiary Filter:** Recommendations should enable next phase progression

**Example Filtering:**

For a repository in Architecture phase with:
- Knowledge artifacts present (mature)
- Architecture artifacts present (developing)
- No implementation artifacts

**Appropriate recommendations:**
- Create architecture artifacts deriving from knowledge
- Review architecture against knowledge artifacts
- Document architectural decisions

**Excluded/deprioritized recommendations:**
- Create implementation (Architecture phase prerequisite not met)
- Execute verification tests (No implementation to verify)

**Chain of Authority Compliance:**
Recommendations that would violate the Chain of Authority must not be presented as primary recommendations. Implementation can only be recommended when:
- Architecture artifacts exist and are mature
- Implementation derives from and traces to Architecture

**Exit Criteria:**
- Repository phase detected and recorded
- Recommendations filtered to phase-appropriate actions
- Chain of Authority compliance verified
- Phase context included in report

---

### Phase 4: Approval

**Duration:** Variable (human-dependent)

**Purpose:** Pause for human authorization before implementation

**Activities:**

#### 4.1 Present Recommendation

Deliver KDSE Report to human with:
- Clear action description
- Expected impact
- Risk assessment
- Decision options

#### 4.2 Await Decision

Wait for human response with options:
- **APPROVE**: Proceed with action
- **APPROVE WITH MODIFICATIONS**: Proceed with changes
- **REJECT**: Decline, next action will be recommended
- **DEFER**: Postpone to later session
- **TERMINATE**: End session

#### 4.3 Record Decision

Document human decision:
- Decision type
- Rationale (if provided)
- Modifications (if any)
- Timestamp

**Exit Criteria:**
- Human provides explicit decision
- Decision recorded
- Appropriate transition determined

---

### Phase 5: Implementation

**Duration:** Variable (work-dependent)

**Purpose:** Execute approved work according to KDSE principles

**Activities:**

#### 5.1 Prepare Implementation

- Review implementation constraints
- Confirm preconditions met
- Prepare implementation environment
- Assign responsibilities

#### 5.2 Execute Work

Follow KDSE Engineering Model:
- Derive from knowledge artifacts
- Maintain authority hierarchy
- Document all decisions
- Create/update traceability links
- Follow artifact lifecycle

#### 5.3 Document Changes

- Record all modifications
- Document deviations (if any)
- Update artifact references
- Maintain evidence trail

**Constraints:**
- Must respect authority hierarchy
- Must maintain traceability
- Must document all decisions
- Cannot proceed without approval

**Exit Criteria:**
- Implementation complete
- All artifacts updated
- Traceability verified
- Ready for verification

---

### Phase 6: Verification

**Duration:** Typically 5-30 minutes

**Purpose:** Confirm results through re-audit

**Activities:**

#### 6.1 Re-run Compliance Audit

Execute full or targeted Compliance Audit:
- Compare current scores to baseline
- Document improvements
- Identify new gaps (if any)
- Verify implementation correctness

#### 6.2 Measure Progress

Calculate:
- Score improvement (delta)
- Findings addressed
- Progress to target
- Diminishing returns assessment

#### 6.3 Generate Updated Report

Produce new KDSE Report with:
- Updated scores
- Progress summary
- Next recommendations (if applicable)

**Exit Criteria:**
- Verification complete
- Progress measured
- Decision point reached

---

## Loop Decision Points

After each iteration, the loop reaches a decision point:

```
┌─────────────────────────────────────────────────────────────────┐
│                     DECISION POINT                               │
│                                                                  │
│  Target maturity reached? ──────YES─────▶ COMPLETE SESSION       │
│         │                                                         │
│         NO                                                        │
│         │                                                         │
│         ▼                                                         │
│  More high-value actions available?                              │
│         │                                                         │
│         ├───────────YES──────────▶ RETURN TO PHASE 2            │
│         │                      (Continue Loop)                   │
│         │                                                         │
│         NO                                                        │
│         │                                                         │
│         ▼                                                         │
│  Diminishing returns threshold reached?                          │
│         │                                                         │
│         ├───────────YES──────────▶ COMPLETE SESSION              │
│         │                      (No more value)                    │
│         │                                                         │
│         NO                                                        │
│         │                                                         │
│         ▼                                                         │
│  ────────────▶ RETURN TO PHASE 2                                │
│  (Unlikely case)                                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Decision Criteria

| Condition | Action |
|-----------|--------|
| Score ≥ Target | Complete session |
| Value gained < Threshold | Complete session (diminishing returns) |
| No actionable findings | Complete session |
| Human terminates | Terminate session |
| Otherwise | Continue loop |

---

## Loop Metrics

Track key metrics across iterations:

### Per-Iteration Metrics

| Metric | Description |
|--------|-------------|
| Iteration Duration | Time for this iteration |
| Score Improvement | Delta from previous score |
| Actions Completed | Number of actions approved |
| Findings Addressed | Number of findings resolved |
| Human Decisions | Approval/rejection count |

### Cumulative Metrics

| Metric | Description |
|--------|-------------|
| Total Duration | Session total time |
| Total Improvement | Score delta from baseline |
| Total Actions | Actions completed this session |
| Progress Rate | Improvement per hour |
| ROI Assessment | Value gained vs. effort invested |

---

## Loop Variants

### Standard Loop

Default execution for general use.

```
INITIALIZING → ASSESSING → RECOMMENDING → APPROVAL → IMPLEMENTING → VERIFYING → (repeat or complete)
```

### Targeted Loop

Focused on specific dimension(s).

```
INITIALIZING → TARGETED_ASSESSING → RECOMMENDING → APPROVAL → IMPLEMENTING → TARGETED_VERIFYING → (repeat or complete)
```

### Quick Loop

Minimal assessment for known state.

```
INITIALIZING → QUICK_ASSESSING → RECOMMENDING → APPROVAL → IMPLEMENTING → VERIFYING → (repeat or complete)
```

---

## Loop Integration with KDSE Engineering Model

The Execution Loop maps to KDSE Engineering Model stages:

```
┌──────────────────────────────────────────────────────────────┐
│              KDSE Engineering Model                          │
│                                                              │
│  Knowledge                                                   │
│      │                                                       │
│      ▼                                                       │
│  Architecture                                                │
│      │                                                       │
│      ▼                                                       │
│  Implementation                                              │
│      │                                                       │
│      ▼                                                       │
│  Verification                                                │
│      │                                                       │
│      ▼                                                       │
│  Evolution                                                   │
└──────────────────────────────────────────────────────────────┘
                          │
                          │ Maps to
                          ▼
┌──────────────────────────────────────────────────────────────┐
│              KDSE Execution Loop                             │
│                                                              │
│  Load Standards ────▶ Knowledge (establishing context)       │
│                                                              │
│  Foundation Verification ────▶ Knowledge (verification)      │
│                                                              │
│  Compliance Audit ────▶ Verification (assessing state)        │
│                                                              │
│  Recommend Action ────▶ Architecture (prioritization)        │
│                                                              │
│  Await Approval ────▶ Governance (authorization)           │
│                                                              │
│  Implement Work ────▶ Implementation (doing work)            │
│                                                              │
│  Verify Results ────▶ Verification (confirming alignment)    │
│                                                              │
│  Repeat ────▶ Evolution (continuous improvement)            │
└──────────────────────────────────────────────────────────────┘
```

---

## Loop Governance

### Session Governance

Each session maintains:

| Element | Description |
|---------|-------------|
| Owner | Human responsible for approvals |
| Target | Desired maturity level |
| Constraints | Session limitations |
| Timeout | Maximum session duration |
| Iteration Limit | Maximum iterations (optional) |

### Loop Governance

The loop itself is governed by:

| Principle | Application |
|-----------|-------------|
| Audit-First | Always assess before acting |
| Evidence-Based | All recommendations backed by evidence |
| Human Approval | No implementation without authorization |
| Traceability | Maintain links throughout |
| Progress Measurement | Verify results after each action |

---

## Loop Examples

### Example 1: New Repository Maturation

**Context:** Empty repository, target 6.0/10

**Iteration 1:**
- ASSESSING: Baseline score 2.0/10 (Concept)
- RECOMMENDING: "Create initial knowledge artifacts"
- APPROVAL: APPROVED
- IMPLEMENTING: Create requirements knowledge
- VERIFYING: Score 3.2/10
- DECISION: Continue

**Iteration 2:**
- ASSESSING: New score 3.2/10
- RECOMMENDING: "Define architecture approach"
- APPROVAL: APPROVED
- IMPLEMENTING: Create architecture documents
- VERIFYING: Score 4.5/10
- DECISION: Continue

**Iteration 3-6:**
- Continue loop
- Each iteration adds artifacts, improves scores

**Iteration 7:**
- ASSESSING: Score 5.9/10
- RECOMMENDING: "Establish verification practices"
- APPROVAL: APPROVED
- IMPLEMENTING: Create verification artifacts
- VERIFYING: Score 6.1/10
- DECISION: Target reached → COMPLETE

### Example 2: Gap Remediation

**Context:** Repository at 6.8/10, single gap in Verification

**Iteration 1:**
- ASSESSING: Targeted assessment confirms gap
- RECOMMENDING: "Implement verification criteria derivation"
- APPROVAL: APPROVED
- IMPLEMENTING: Add verification criteria
- VERIFYING: Gap addressed
- DECISION: Continue

**Iteration 2:**
- ASSESSING: Full audit confirms improvement
- RECOMMENDING: "No more high-value actions"
- DECISION: Diminishing returns → COMPLETE

---

## Loop Anti-Patterns

Avoid these patterns:

| Anti-Pattern | Problem | Correct Approach |
|---------------|---------|------------------|
| Skip Assessment | Acting without evidence | Always assess first |
| Skip Approval | Implementing without authorization | Always await approval |
| Skip Verification | Assuming improvement | Always verify results |
| Over-Engineer | More work than value | Stop at diminishing returns |
| Ignore Diminishing Returns | Continued effort without benefit | Complete when value < effort |

---

## Loop Quality Gates

Each phase has quality gates:

### Phase 2 Quality Gate: Assessment Complete

- [ ] All required audits executed
- [ ] Evidence documented for each finding
- [ ] Findings categorized by severity
- [ ] Baseline score established

### Phase 3 Quality Gate: Recommendation Ready

- [ ] KDSE Report follows standard format
- [ ] Recommendation backed by evidence
- [ ] Value calculation documented
- [ ] Alternatives considered

### Phase 4 Quality Gate: Approval Obtained

- [ ] Human provided explicit decision
- [ ] Decision recorded with timestamp
- [ ] Modifications documented (if any)

### Phase 5 Quality Gate: Implementation Verified

- [ ] All artifacts created/updated
- [ ] Traceability links established
- [ ] Decisions documented
- [ ] Deviations noted (if any)

### Phase 6 Quality Gate: Results Confirmed

- [ ] Audit re-run and complete
- [ ] Score improvement documented
- [ ] Progress to target calculated

---

## Loop Termination Conditions

A session terminates when any condition is met:

| Condition | Action | Next Session |
|-----------|--------|--------------|
| Target reached | Complete | New session for new target |
| No actionable findings | Complete | New session if new gaps emerge |
| Diminishing returns | Complete | Resume later if value returns |
| Human terminates | Terminate | Resume with same parameters |
| Timeout | Terminate | Resume with timeout adjustment |
| Error | Terminate | Fix error, resume session |

---

## Continuous Improvement Principles

The Execution Loop embodies KDSE's continuous improvement philosophy:

### 1. Audit-First Principle

Every action begins with assessment. The Compliance Audit is the source of truth, not assumptions.

### 2. Evidence-Based Action

Recommendations flow from audit evidence. Findings → Analysis → Action.

### 3. Human Authorization

Humans remain in control. The Agent recommends; humans decide.

### 4. Measurement-Based Progress

Progress is measured, not assumed. Each iteration ends with verification.

### 5. Diminishing Returns Recognition

The loop knows when to stop. Value must exceed effort to continue.

### 6. Repository Independence

The loop applies to any repository. Content varies; process remains constant.

---

## Loop Summary

The KDSE Execution Loop is the operational heart of the KDSE Execution Model:

| Element | Description |
|---------|-------------|
| Initiation | Session starts with parameters |
| Assessment | Audits establish current state |
| Recommendation | Evidence generates guidance |
| Approval | Humans authorize action |
| Implementation | Work follows KDSE principles |
| Verification | Audits confirm results |
| Decision | Loop continues or completes |
| Evolution | Maturity increases through iterations |

**The loop continues until target maturity is reached or no further high-value actions exist.**

---

*This execution loop is the continuous improvement engine of KDSE. Each iteration increases maturity until targets are achieved or value is exhausted.*
