# KDSE Session Protocol

**Document Version:** 1.0  
**Effective Date:** 2026-07-10

---

## Purpose

The Session Protocol defines the lifecycle of a KDSE development session. A session is a bounded period of engineering work guided by KDSE principles, anchored in audit evidence, and proceeding through defined states until completion or termination.

---

## Session Definition

A **KDSE Session** is a discrete unit of engineering work that:

1. Begins with established context (repository, target maturity, constraints)
2. Proceeds through defined states following the Execution Loop
3. Produces artifacts (KDSE Report, implemented work, verification results)
4. Ends with either completion or termination
5. Contributes to overall maturity progress

---

## Session Lifecycle States

```
┌─────────────┐
│ INITIALIZING│
└──────┬──────┘
       │
       ▼
┌─────────────┐     ┌───────────┐
│ ASSESSING   │────▶│RECOMMENDING│
└──────┬──────┘     └─────┬─────┘
       │                  │
       │                  ▼
       │          ┌──────────────┐
       │          │AWAITING_     │
       │          │APPROVAL       │
       │          └──────┬───────┘
       │                 │
       ▼                 ▼
┌─────────────┐   ┌─────────────┐
│ VERIFYING   │◀──│IMPLEMENTING │
└──────┬──────┘   └─────────────┘
       │
       ▼
┌─────────────┐     ┌─────────────┐
│ COMPLETED   │     │ TERMINATED  │
└─────────────┘     └─────────────┘
```

---

## State Definitions

### INITIALIZING

**Purpose:** Establish session context and load required resources.

**Entry Conditions:**
- User command to start session received
- Repository location identified
- Session parameters defined (target maturity, scope, constraints)

**Activities:**
1. Load KDSE standards and templates
2. Retrieve repository metadata
3. Identify relevant knowledge artifacts
4. Establish baseline context
5. Record session parameters

**Exit Conditions:**
- All resources loaded successfully
- Session parameters validated
- Context established for assessment

**Outputs:**
- Loaded standards cache
- Repository metadata summary
- Session parameter record

---

### ASSESSING

**Purpose:** Evaluate current state through audits and gather evidence.

**Entry Conditions:**
- Session initialized
- Resources loaded
- Clear assessment scope defined

**Activities:**
1. Run Foundation Verification (standards check)
2. Run Repository Assessment (structure and content)
3. Run Compliance Audit (full KDSE evaluation)
4. Collect audit evidence
5. Identify findings and gaps

**Audit Reference:**
- Foundation Verification: [FOUNDATION_AUDIT.md](../audit/FOUNDATION_AUDIT.md) dimensions 1-5
- Repository Assessment: Structure and artifact inventory
- Compliance Audit: [COMPLIANCE_AUDIT.md](../audit/COMPLIANCE_AUDIT.md) dimensions 1-7

**Exit Conditions:**
- Audits completed successfully
- Evidence collected and documented
- Findings categorized by severity

**Outputs:**
- Foundation Verification results
- Repository Assessment report
- Compliance Audit report
- Findings inventory

---

### RECOMMENDING

**Purpose:** Generate KDSE Report and identify highest-value next actions.

**Entry Conditions:**
- Assessment complete
- Findings documented
- Evidence available

**Activities:**
1. Generate KDSE Report (see [REPORT_FORMAT.md](REPORT_FORMAT.md))
2. Analyze findings for actionability
3. Prioritize actions by value and impact
4. Identify single highest-value next action
5. Document recommendation rationale

**Exit Conditions:**
- KDSE Report generated
- Next action recommended with clear rationale
- Expected impact documented
- Human approval requested

**Outputs:**
- KDSE Report
- Recommended next action
- Expected impact assessment
- Approval request

---

### AWAITING_APPROVAL

**Purpose:** Pause execution pending human authorization.

**Entry Conditions:**
- Recommendation generated
- Human notification sent
- Decision pending

**Activities:**
1. Present KDSE Report to human
2. Highlight recommended action
3. Display expected impact
4. Await explicit approval or rejection
5. Record decision and any modifications

**Human Decision Options:**
- **APPROVE**: Proceed with recommended action
- **APPROVE_WITH_MODIFICATIONS**: Proceed with changes specified
- **REJECT**: Decline recommended action
- **DEFER**: Postpone decision to later session
- **TERMINATE**: End session entirely

**Exit Conditions:**
- Human provides decision
- Decision recorded
- Appropriate next state determined

**Outputs:**
- Decision record
- Modifications (if any)
- Next state determination

---

### IMPLEMENTING

**Purpose:** Execute approved work according to KDSE principles.

**Entry Conditions:**
- Human approval received
- Work scope clearly defined
- Implementation plan available

**Activities:**
1. Derive implementation from knowledge artifacts
2. Maintain traceability throughout implementation
3. Document all decisions
4. Create or update artifacts as needed
5. Maintain authority hierarchy compliance

**Constraints:**
- Must not contradict higher-authority artifacts
- Must maintain traceability links
- Must document all deviations
- Must preserve knowledge precedence

**Exit Conditions:**
- Implementation complete
- Artifacts created or updated
- Traceability links established
- Ready for verification

**Outputs:**
- Implemented artifacts
- Updated traceability records
- Decision log
- Deviation documentation (if any)

---

### VERIFYING

**Purpose:** Confirm implementation results and measure progress.

**Entry Conditions:**
- Implementation complete
- Verification criteria defined
- Verification scope identified

**Activities:**
1. Re-run Compliance Audit (or relevant dimensions)
2. Compare current scores to baseline
3. Document improvements achieved
4. Identify any new gaps introduced
5. Assess progress toward target maturity

**Verification Standards:**
- [COMPLIANCE_AUDIT.md](../audit/COMPLIANCE_AUDIT.md)
- [AUDIT_SCORING.md](../audit/AUDIT_SCORING.md)

**Exit Conditions:**
- Verification complete
- Progress measured
- Decision made on continuation

**Outputs:**
- Verification report
- Score comparison (baseline vs. current)
- Progress assessment
- Continuation recommendation

---

### COMPLETED

**Purpose:** End session successfully with full metrics and records.

**Entry Conditions:**
- Target achieved OR no further high-value actions identified
- All artifacts in final state
- Metrics recorded

**Activities:**
1. Finalize KDSE Report
2. Record session metrics
3. Archive session artifacts
4. Update maturity records
5. Generate completion summary

**Exit Conditions:**
- Session formally closed
- All records complete
- Metrics available for analysis

**Outputs:**
- Final KDSE Report
- Session metrics
- Maturity progression record
- Completion summary

---

### TERMINATED

**Purpose:** End session without completion.

**Entry Conditions:**
- User termination request
- Irrecoverable error
- Session timeout
- Decision to end without completion

**Activities:**
1. Record termination reason
2. Preserve partial work
3. Document incomplete items
4. Generate termination summary
5. Update session status

**Exit Conditions:**
- Session formally closed
- Termination reason recorded
- Partial artifacts preserved

**Outputs:**
- Termination report
- Partial KDSE Report (if applicable)
- Termination reason
- Resume guidance (if applicable)

---

## Session Parameters

Each session must define:

| Parameter | Description | Required |
|-----------|-------------|----------|
| repository | Target repository location | Yes |
| target_maturity | Desired compliance level (0-10) | Yes |
| scope | Specific audit dimensions to focus on | No |
| constraints | Resource limits, deadlines, boundaries | No |
| owner | Human responsible for approvals | Yes |
| max_iterations | Session iteration limit | No |

---

## Session Transitions

### Normal Flow

```
INITIALIZING → ASSESSING → RECOMMENDING → AWAITING_APPROVAL → IMPLEMENTING → VERIFYING → COMPLETED
```

### Early Termination Flows

```
INITIALIZING → TERMINATED (error during initialization)
ASSESSING → TERMINATED (assessment failure)
RECOMMENDING → TERMINATED (user request)
AWAITING_APPROVAL → TERMINATED (user decision)
IMPLEMENTING → TERMINATED (implementation failure)
VERIFYING → TERMINATED (verification failure)
```

### Continuation Loop

```
VERIFYING → ASSESSING (if progress insufficient, continue)
VERIFYING → RECOMMENDING (if more actions available)
VERIFYING → COMPLETED (if target reached or no more actions)
```

---

## Session Timing

| Phase | Typical Duration | Notes |
|-------|------------------|-------|
| INITIALIZING | < 1 minute | Depends on repository size |
| ASSESSING | 5-30 minutes | Includes full Compliance Audit |
| RECOMMENDING | 2-5 minutes | Report generation |
| AWAITING_APPROVAL | Variable | Human-dependent |
| IMPLEMENTING | Variable | Work-dependent |
| VERIFYING | 5-30 minutes | Re-run relevant audits |
| COMPLETED | < 1 minute | Record finalization |

---

## Session Metrics

Each session records:

- **Duration**: Total session time
- **State Transitions**: Count and sequence
- **Audit Scores**: Before and after comparison
- **Actions Completed**: Count and list
- **Findings Addressed**: Count and categorization
- **Maturity Progress**: Score improvement
- **Human Decisions**: Approval/rejection count

---

## Session Boundaries

### Session Start

A session starts when:

1. User issues "Run KDSE" command
2. Scheduled execution triggers
3. Repository event initiates session
4. External API receives start request

### Session Isolation

Sessions are isolated:
- Each session has unique ID
- Session state is preserved
- Sessions can run concurrently (different repositories)
- Sessions can run sequentially (same repository)

### Session Resumption

A terminated session may be resumed:
1. New session initiated with same parameters
2. Previous session state retrieved
3. Resume from appropriate state
4. Continuation guidance followed

---

## Relationship to KDSE Engineering Model

The Session Protocol maps to the KDSE Engineering Model stages:

| Session State | Engineering Stage | Focus |
|---------------|-------------------|-------|
| INITIALIZING | Entry | Context establishment |
| ASSESSING | Verification | Current state evaluation |
| RECOMMENDING | Knowledge | Gap analysis and prioritization |
| AWAITING_APPROVAL | Governance | Human authorization |
| IMPLEMENTING | Implementation | Approved work execution |
| VERIFYING | Verification | Result confirmation |
| COMPLETED/TERMINATED | Evolution | State recording |

---

## Examples

### Example 1: New Repository Session

**Context:** Repository at maturity level 3.5, target 6.0.

```
INITIALIZING
  → Load KDSE standards
  → Identify repository artifacts
  → Establish baseline

ASSESSING
  → Foundation Verification (quick check)
  → Repository Assessment (artifact inventory)
  → Compliance Audit (full evaluation: 3.5/10)

RECOMMENDING
  → Generate KDSE Report
  → Identify: "Implement Knowledge Artifacts for Core Requirements"
  → Expected impact: +1.0 maturity points

AWAITING_APPROVAL
  → Human reviews report
  → APPROVE

IMPLEMENTING
  → Create knowledge artifacts
  → Document traceability
  → Update governance

VERIFYING
  → Re-run Compliance Audit
  → New score: 4.5/10
  → Progress confirmed

VERIFYING → RECOMMENDING (more actions available)
  → Identify: "Establish Traceability Links"
  → Expected impact: +0.8 maturity points

AWAITING_APPROVAL
  → APPROVE

IMPLEMENTING
  → Create traceability links
  → Verify completeness

VERIFYING
  → Re-run Compliance Audit
  → New score: 5.3/10

VERIFYING → COMPLETED (target not reached but diminishing returns)
```

### Example 2: Targeted Fix Session

**Context:** Repository at 6.8/10. Known gap in Verification Practices.

```
INITIALIZING
  → Load standards
  → Set scope to dimension 4 (Verification Practices)

ASSESSING
  → Targeted assessment only
  → Identify specific gap: Missing verification criteria derivation

RECOMMENDING
  → Generate focused KDSE Report
  → Recommend: "Derive verification criteria from knowledge artifacts"

AWAITING_APPROVAL
  → APPROVE

IMPLEMENTING
  → Create verification criteria documents
  → Link to knowledge artifacts

VERIFYING
  → Check dimension 4 only
  → Score improved: 6.8 → 7.4

COMPLETED
```

---

*This protocol defines the canonical lifecycle for KDSE development sessions. All implementations must follow this state machine.*
