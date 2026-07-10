# KDSE Report Format

**Document Version:** 1.1  
**Effective Date:** 2026-07-10  
**Change Note:** Added Repository Phase and Assessment vs Compliance Score distinction to address KDSE-CASE-001 OBS-003

---

## Purpose

The KDSE Report extends the existing audit template by combining audit results with actionable engineering guidance. While the Compliance Audit provides diagnostic information, the KDSE Report provides therapeutic direction—the path from current state to target state.

---

## Relationship to Existing Audits

The KDSE Report is built upon existing audit documents:

| Source Document | Usage in KDSE Report |
|-----------------|----------------------|
| [AUDIT_TEMPLATE.md](../audit/AUDIT_TEMPLATE.md) | Base template structure |
| [COMPLIANCE_AUDIT.md](../audit/COMPLIANCE_AUDIT.md) | Audit dimensions, findings, recommendations |
| [AUDIT_SCORING.md](../audit/AUDIT_SCORING.md) | Score presentation, maturity levels |
| [FOUNDATION_AUDIT.md](../audit/FOUNDATION_AUDIT.md) | Foundation verification sections |

The KDSE Report does not duplicate these documents—it references and synthesizes them.

---

## Report Structure Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     KDSE REPORT                             │
├─────────────────────────────────────────────────────────────┤
│ 1. HEADER                                                   │
│    - Session identification                                  │
│    - Timestamps                                             │
│    - Repository information                                 │
├─────────────────────────────────────────────────────────────┤
│ 2. CURRENT STATUS                                            │
│    - Executive summary                                      │
│    - Quick metrics                                          │
│    - Session context                                        │
├─────────────────────────────────────────────────────────────┤
│ 3. AUDIT SUMMARY                                            │
│    - Dimension scores table                                  │
│    - Score visualization                                     │
│    - Comparison to previous (if applicable)                  │
├─────────────────────────────────────────────────────────────┤
│ 4. HIGHEST PRIORITY FINDINGS                                │
│    - Top 3-5 findings by impact                             │
│    - Evidence per finding                                   │
│    - Severity categorization                                │
├─────────────────────────────────────────────────────────────┤
│ 5. RECOMMENDED NEXT ACTION                                  │
│    - Single highest-value action                             │
│    - Detailed action description                            │
│    - Action categorization                                  │
├─────────────────────────────────────────────────────────────┤
│ 6. EXPECTED IMPACT                                          │
│    - Projected score improvement                             │
│    - Timeline estimate                                       │
│    - Risk assessment                                        │
├─────────────────────────────────────────────────────────────┤
│ 7. REQUIRED HUMAN APPROVAL                                  │
│    - Approval request                                       │
│    - Decision options                                       │
│    - Deadline (if applicable)                               │
├─────────────────────────────────────────────────────────────┤
│ 8. SESSION STATE                                            │
│    - Current state                                          │
│    - Progress through loop                                  │
│    - Next steps                                             │
├─────────────────────────────────────────────────────────────┤
│ APPENDIX: Full Audit Results                                │
│    - Complete findings list                                 │
│    - Full dimension breakdown                               │
│    - Evidence index                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## Section Templates

### 1. Header

```markdown
# KDSE Report

## Session Information

| Field | Value |
|-------|-------|
| Report ID | KDSE-REP-{session_id} |
| Generated | {timestamp} |
| Repository | {repository_path} |
| Repository Version | {git_commit_hash} |
| Session Owner | {owner_name} |
| Target Maturity | {target_score}/10 |
| Report Version | 1.0 |
```

**Field Definitions:**

| Field | Description | Source |
|-------|-------------|--------|
| Report ID | Unique identifier following KDSE-REP-{session_id} format | Agent generated |
| Generated | ISO 8601 timestamp of report creation | Agent generated |
| Repository | Path or URL of the audited repository | Session parameter |
| Repository Version | Git commit hash or equivalent | Repository |
| Session Owner | Human responsible for approvals | Session parameter |
| Target Maturity | Desired compliance score | Session parameter |
| Report Version | Format version (always 1.0 for this spec) | Specification |

---

### 2. Current Status

```markdown
## Current Status

### Executive Summary

{3-5 sentence summary of current state, key findings, and recommended direction.}

### Quick Metrics

| Metric | Value | Change |
|--------|-------|--------|
| Overall Score | {score}/10 ({maturity_level}) | {+/-delta from baseline} |
| Highest Dimension | {dimension_name} ({score}/10) | - |
| Lowest Dimension | {dimension_name} ({score}/10) | - |
| Critical Findings | {count} | {change} |
| High Findings | {count} | {change} |

### Session Context

- **Session Type**: {Initial | Continuation | Targeted}
- **Session Count**: {sequential number for this repository}
- **Session Duration**: {estimated time to target}
- **Constraints**: {any session constraints}
```

**Maturity Level Reference:**

| Score Range | Level |
|-------------|-------|
| 0-2 | Concept |
| 2-4 | Defined |
| 4-6 | Structured |
| 6-8 | Usable |
| 8-9 | Validated |
| 9-10 | Proven |

---

### 3. Audit Summary

```markdown
## Audit Summary

### Repository Phase

| Field | Value |
|-------|-------|
| Detected Phase | {phase} |
| Primary Metric | {Assessment Score or Compliance Score} |

**Phase Context:** Per [COMPLIANCE_AUDIT.md](../docs/audit/COMPLIANCE_AUDIT.md), the primary metric is:
- **Assessment Score** for repositories in Research, Knowledge Development, or Architecture phases
- **Compliance Score** for repositories in Implementation, Verification, or Evolution phases

### Dimension Scores

| Dimension | Score | Change | Status |
|-----------|-------|--------|--------|
| 1. Knowledge Artifacts | {X}/10 | {+/-Y} | {Status} |
| 2. Architecture Artifacts | {X}/10 | {+/-Y} | {Status} |
| 3. Implementation Artifacts | {X}/10 | {+/-Y} | {Status} |
| 4. Verification Practices | {X}/10 | {+/-Y} | {Status} |
| 5. Traceability | {X}/10 | {+/-Y} | {Status} |
| 6. Authority Hierarchy | {X}/10 | {+/-Y} | {Status} |
| 7. Governance | {X}/10 | {+/-Y} | {Status} |
| **Overall Score** | **{X}/10** | **{+/-Y}** | **{Status}** |

### Score Key

- **Status Indicators:**
  - ✅ Meeting target (within 1 point of target)
  - ⚠️ Below target (1-2 points below)
  - ❌ Significantly below target (>2 points below)
  - 🎯 Target reached

### Score Comparison (If Previous Audit Available)

{Visual representation of score changes, e.g.,}

```
Dimension          Previous   Current   Change
───────────────────────────────────────────────
Knowledge              5.0       5.5     +0.5
Architecture           4.0       4.5     +0.5
Implementation         5.5       5.5      —
Verification           3.0       4.0     +1.0
Traceability           4.5       5.0     +0.5
Authority              6.0       6.0      —
Governance             5.0       5.5     +0.5
───────────────────────────────────────────────
Overall                4.7       5.1     +0.4
```

### Gap Analysis Summary

| Category | Count | Impact |
|----------|-------|--------|
| Critical Gaps | {n} | Blocks progress |
| Major Gaps | {n} | Significant impact |
| Minor Gaps | {n} | Improvement opportunity |

### Compliance Level

**Current Level:** {Level} (Score: {X}/10)

{Level definition from COMPLIANCE_AUDIT.md}

| Level | Score Range | Description |
|-------|-------------|-------------|
| Compliant | 8-10 | Fully compliant with KDSE |
| Substantially Compliant | 6-8 | Compliant with minor gaps |
| Partially Compliant | 4-6 | Significant gaps exist |
| Minimally Compliant | 2-4 | Basic structure, major gaps |
| Non-Compliant | 0-2 | Does not meet KDSE standards |
```

---

### 4. Highest Priority Findings

```markdown
## Highest Priority Findings

{For each of top 3-5 findings:}

### Finding {N}: {Finding Title}

**Dimension:** {Affected dimension from COMPLIANCE_AUDIT.md}
**Severity:** {Critical | High | Medium | Low}
**Score Impact:** {Estimated score reduction if unaddressed}

**Description:**
{Brief description of the finding}

**Evidence:**
```
{Reference to specific evidence observed during audit}
```

**Impact:**
{What this finding prevents or causes}

**Recommendation:**
{How to address this finding}

**Actionable:** {Yes | No}
**Action Priority:** {1 | 2 | 3 | 4} (if actionable)

---

### Findings Requiring Immediate Attention

{Critical findings that block progress, if any}

### Findings Ready for Action

{Actionable findings prioritized by value, if any}

### Findings Requiring Further Analysis

{Non-actionable findings that need more study, if any}
```

---

### 5. Recommended Next Action

```markdown
## Recommended Next Action

### Action Summary

| Field | Value |
|-------|-------|
| Action ID | KDSE-ACT-{sequential_id} |
| Priority | {1 (Critical) | 2 (High) | 3 (Medium)} |
| Target Dimension | {Dimension this action addresses} |
| Expected Score Impact | +{X.Y} points |

### Action Description

**What:**
{Explicit description of the action to take}

**Why:**
{Rationale connecting action to findings and impact}

**How:**
{Implementation approach in 2-5 steps}

### Derived From

This action is derived from:
- Finding(s): {Related finding IDs}
- Dimension: {Target dimension}
- Evidence: {Reference to audit evidence}

### Constraints

- **Authority Level:** {Knowledge | Architecture | Implementation}
- **Traceability Required:** {Yes | No}
- **Steward Approval Needed:** {Yes | No, and from whom}

### Preconditions

Before implementing this action:
1. {Prerequisite 1}
2. {Prerequisite 2}

### Postconditions

After implementing this action:
1. {Expected outcome 1}
2. {Expected outcome 2}

### Alternative Actions Considered

| Alternative | Rejected Because |
|-------------|------------------|
| {Alternative 1} | {Reason for rejection} |
| {Alternative 2} | {Reason for rejection} |

### Backlog Items

Actions identified but not recommended now:

| Action | Reason for Deferral |
|--------|---------------------|
| {Deferred action 1} | {Reason} |
| {Deferred action 2} | {Reason} |
```

---

### 6. Expected Impact

```markdown
## Expected Impact

### Score Projection

| Dimension | Current | Projected | Change |
|-----------|---------|-----------|--------|
| {Dimension 1} | {X}/10 | {Y}/10 | +{Z} |
| {Dimension 2} | {X}/10 | {Y}/10 | +{Z} |
| {Target Dimension} | {X}/10 | {Y}/10 | +{Z} |
| **Overall** | **{X}/10** | **{Y}/10** | **+{Z}** |

### Progress to Target

```
Current: {X}/10 ({maturity_level})
Target:   {Y}/10
Progress: {percentage}% complete
Remaining: {Y-X} points
```

### Timeline Estimate

| Phase | Estimated Duration |
|-------|-------------------|
| Implementation | {X} hours/days |
| Verification | {X} hours/days |
| **Total** | **{X} hours/days** |

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-----------|
| {Risk 1} | {Low/Medium/High} | {Low/Medium/High} | {Mitigation} |
| {Risk 2} | {Low/Medium/High} | {Low/Medium/High} | {Mitigation} |

### Dependencies

This action depends on:
- {Dependency 1}
- {Dependency 2}

This action enables:
- {Enabling action 1}
- {Enabling action 2}

### Value Assessment

**Value Score:** {Calculated value/impact ratio}

{Explanation of value calculation}

---

### 7. Required Human Approval

```markdown
## Required Human Approval

### Approval Request

**To:** {Session Owner}
**From:** KDSE Agent
**Request:** Implement recommended action

**Recommended Action:**
{Repeat action summary from Section 5}

**Expected Impact:**
{Repeat impact summary from Section 6}

### Decision Options

| Option | Action |
|--------|--------|
| **APPROVE** | Proceed with recommended action as described |
| **APPROVE WITH MODIFICATIONS** | Proceed with specified changes |
| **REJECT** | Decline this action; next best alternative will be recommended |
| **DEFER** | Postpone decision to later session |
| **TERMINATE** | End session without implementing |

### Approval Deadline

{If deadline applies: "This decision is requested by {deadline} to maintain session timeline."}
{If no deadline: "No deadline; decision at your discretion."}

### Accountability

By approving this action, you acknowledge:
1. The action has been derived from audit evidence
2. The expected impact has been assessed
3. You accept responsibility for this engineering decision

### Approval Response Format

```
APPROVAL: [APPROVE | APPROVE_WITH_MODIFICATIONS | REJECT | DEFER | TERMINATE]
REASON: [Your rationale, required for non-approve decisions]
MODIFICATIONS: [If APPROVE_WITH_MODIFICATIONS, describe changes]
SIGNATURE: [Your name and timestamp]
```

---

### 8. Session State

```markdown
## Session State

### Current State

| Field | Value |
|-------|-------|
| State | {Current session state} |
| Iteration | {N} of {max} |
| Duration | {Elapsed time} |
| Next Deadline | {If applicable} |

### Progress Through Execution Loop

```
[✓] Run KDSE
    │
[✓] Load Standards
    │
[✓] Foundation Verification
    │
[✓] Repository Assessment
    │
[✓] Compliance Audit
    │
[✓] Generate KDSE Report
    │
[✓] Recommend Action
    │
[→] Await Human Approval
    │
[ ] Implement Approved Work
    │
[ ] Verify Results
    │
[ ] {Continue or Complete}
```

### Session History

| Iteration | Action | Approved | Score Change |
|-----------|--------|----------|--------------|
| 1 | {Action 1} | Yes | +0.5 |
| 2 | {Action 2} | Yes | +0.8 |
| ... | ... | ... | ... |

### Next Steps

**If Approved:**
1. Implement recommended action
2. Re-run Compliance Audit
3. Generate updated KDSE Report
4. Continue or complete based on results

**If Rejected:**
1. Identify next highest-value action
2. Generate updated recommendation
3. Request new approval

**If Terminated:**
1. Preserve session state
2. Generate termination report
3. Provide resume guidance

### Human Decision Required

**Action Needed:** Approval decision on recommended action

**Deadline:** {If applicable}

**Consequence of No Response:** {If deadline applies: "Session will auto-continue with next best action" or "Session will terminate"}
```

---

### Appendix: Full Audit Results

```markdown
## Appendix: Full Audit Results

### Complete Findings List

| ID | Dimension | Severity | Title | Actionable |
|----|-----------|----------|-------|------------|
| F001 | Dimension | Critical | {Title} | Yes |
| F002 | Dimension | High | {Title} | Yes |
| ... | ... | ... | ... | ... |

### Full Dimension Breakdown

#### 1. Knowledge Artifacts

**Score:** {X}/10

**Strengths:**
- {What is working well}

**Gaps:**
- {What needs improvement}

**Evidence:**
```
{Detailed evidence}
```

**Recommendation:**
{If score below target}

---

[Repeat for each dimension]

### Evidence Index

| Evidence ID | Type | Source | Finding(s) |
|-------------|------|--------|------------|
| E001 | Document | {path} | F001, F002 |
| E002 | Code | {path} | F003 |
| ... | ... | ... | ... |

### Audit Methodology

- **Audit Standard:** KDSE Compliance Audit 1.0
- **Audit Date:** {timestamp}
- **Auditor:** KDSE Agent
- **Repository Version:** {hash}

---

## Document Control

| Field | Value |
|-------|-------|
| Report Version | 1.0 |
| Generated By | KDSE Agent |
| Generation Date | {ISO 8601 timestamp} |
| Template Version | 1.0 |

### Change History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | {date} | Initial report |

---

*This KDSE Report was generated following the KDSE Execution Model. Audit results reference COMPLIANCE_AUDIT.md. Scoring follows AUDIT_SCORING.md.*
```

---

## Report Completion Checklist

Before finalizing a KDSE Report:

- [ ] All required sections present
- [ ] Header complete with all fields
- [ ] Current status reflects audit results accurately
- [ ] Audit summary includes all dimension scores
- [ ] Findings prioritized by impact
- [ ] Recommendation clearly justified
- [ ] Expected impact calculated
- [ ] Approval request explicit
- [ ] Session state accurate
- [ ] Appendix complete with full audit results
- [ ] Evidence references verified
- [ ] Terminology consistent with glossary

---

## Example KDSE Report (Abbreviated)

```markdown
# KDSE Report

## Session Information

| Field | Value |
|-------|-------|
| Report ID | KDSE-REP-2026-07-10-001 |
| Generated | 2026-07-10T14:30:00Z |
| Repository | /workspace/project/example |
| Repository Version | a1b2c3d |
| Session Owner | Jane Developer |
| Target Maturity | 7.0/10 |
| Report Version | 1.0 |

## Current Status

### Executive Summary

The example repository currently achieves a compliance score of 5.2/10 (Structured 
maturity), with significant gaps in Knowledge Artifacts and Verification Practices. 
The highest-impact finding is the absence of structured knowledge for core 
requirements, which blocks progress in three downstream dimensions. Implementing 
knowledge artifacts is recommended as the highest-value next action.

### Quick Metrics

| Metric | Value | Change |
|--------|-------|--------|
| Overall Score | 5.2/10 (Structured) | +0.0 (baseline) |
| Highest Dimension | Authority Hierarchy (6.5/10) | — |
| Lowest Dimension | Knowledge Artifacts (3.5/10) | — |
| Critical Findings | 1 | — |
| High Findings | 3 | — |

## Audit Summary

| Dimension | Score | Status |
|-----------|-------|--------|
| 1. Knowledge Artifacts | 3.5/10 | ❌ |
| 2. Architecture Artifacts | 5.0/10 | ⚠️ |
| 3. Implementation Artifacts | 5.5/10 | ⚠️ |
| 4. Verification Practices | 4.0/10 | ❌ |
| 5. Traceability | 5.5/10 | ⚠️ |
| 6. Authority Hierarchy | 6.5/10 | ✅ |
| 7. Governance | 5.5/10 | ⚠️ |
| **Overall Score** | **5.2/10** | **⚠️** |

## Highest Priority Findings

### Finding 1: Missing Core Knowledge Artifacts

**Dimension:** Knowledge Artifacts  
**Severity:** Critical  
**Score Impact:** Estimated -1.5 points

**Description:** No structured knowledge artifacts exist for core requirements.

**Evidence:** Directory /docs/knowledge is empty. No requirements documentation.

**Recommendation:** Create knowledge artifacts for all identified requirements.

---

## Recommended Next Action

| Field | Value |
|-------|-------|
| Action ID | KDSE-ACT-001 |
| Priority | 1 (Critical) |
| Target Dimension | Knowledge Artifacts |
| Expected Score Impact | +1.2 points |

**What:** Create structured knowledge artifacts documenting core requirements, 
constraints, and context.

**Why:** Missing knowledge artifacts block architecture derivation, prevent 
effective verification, and break the authority chain.

**How:**
1. Identify all core requirements from existing sources
2. Structure requirements as Knowledge artifacts per KDSE standards
3. Document provenance and assumptions
4. Assign stewardship
5. Establish knowledge-architecture traceability

---

## Expected Impact

| Dimension | Current | Projected | Change |
|-----------|---------|-----------|--------|
| Knowledge Artifacts | 3.5/10 | 5.0/10 | +1.5 |
| Architecture Artifacts | 5.0/10 | 5.5/10 | +0.5 |
| Traceability | 5.5/10 | 6.0/10 | +0.5 |
| **Overall** | **5.2/10** | **6.0/10** | **+0.8** |

---

## Required Human Approval

**To:** Jane Developer  
**From:** KDSE Agent  
**Request:** Implement recommended action

**Decision Options:**
- **APPROVE**: Proceed with knowledge artifact creation
- **REJECT**: Decline; alternative will be recommended
- **DEFER**: Postpone to later session

**Approval Response:**

```
APPROVAL: [APPROVE | REJECT | DEFER]
REASON: [Your rationale]
SIGNATURE: [Your name]
```

---

## Session State

**State:** AWAITING_APPROVAL  
**Iteration:** 1 of unlimited  
**Progress:** 7/10 steps complete

**Next Steps:**
- If approved: Implement knowledge artifacts → Verify → Continue
- If rejected: Recommend alternative → Request new approval

**Human Decision Required:** Approval on recommended action
```

---

*This format extends the audit template to provide actionable guidance while maintaining evidence-based rigor.*
