# KDSE Runtime Prompts

**Document Version:** 1.0  
**Type:** Informative Reference Implementation  
**Effective Date:** 2026-07-10

---

## Purpose

This document provides official prompts for operating the KDSE Runtime. These prompts are templates that initiate Runtime actions by referencing the KDSE Standard.

**Note:** These prompts are for human or automated operators. The Runtime interprets these commands and executes the corresponding workflow.

---

## Command Reference

| Command | Purpose | Reference |
|---------|---------|-----------|
| `Run KDSE` | Start a new session | [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) |
| `Continue KDSE` | Resume a session | [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) |
| `KDSE Status` | View session state | [EXECUTION_MODEL.md](EXECUTION_MODEL.md) |
| `KDSE Report` | Generate report | [REPORT_SPEC.md](REPORT_SPEC.md) |
| `Close KDSE` | End a session | [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) |

---

## Session Commands

### Run KDSE

Start a new KDSE development session.

**Prompt:**
```
Run KDSE

Repository: [repository path or URL]
Target Maturity: [target score 0-10]
Operator: [your name]
Scope: [specific dimensions, optional]
Constraints: [resource limits, optional]
```

**Example:**
```
Run KDSE

Repository: /workspace/project/myapp
Target Maturity: 7.0
Operator: Jane Developer
Scope: Knowledge Artifacts, Traceability
```

**What Happens:**
1. Runtime initializes session
2. Loads KDSE Standard
3. Runs Foundation Verification
4. Assesses repository
5. Generates Runtime Report
6. Presents recommendation

**Reference:** [SESSION_PROTOCOL.md - Session Start](SESSION_PROTOCOL.md#phase-1-session-start)

---

### Run KDSE (Minimal)

Start a new session with minimal parameters.

**Prompt:**
```
Run KDSE

Repository: [repository]
Target: [score]
```

**Example:**
```
Run KDSE

Repository: /workspace/project/myapp
Target: 7.0
```

**What Happens:**
1. Runtime uses default parameters
2. Operator assumed from context
3. Full scope assessment

---

### Continue KDSE

Resume an existing or paused session.

**Prompt:**
```
Continue KDSE

Session ID: [session ID]
```

**Example:**
```
Continue KDSE

Session ID: KDSE-RT-2026-07-10-001
```

**What Happens:**
1. Runtime restores session state
2. Presents current report
3. Awaits operator decision
4. Continues from paused point

**Reference:** [SESSION_PROTOCOL.md - Continuation](SESSION_PROTOCOL.md#phase-6-session-continuation)

---

### KDSE Status

View the current session state.

**Prompt:**
```
KDSE Status
```

**What Returns:**
```
Session ID: [session ID]
Phase: [current phase]
Iteration: [N]
Duration: [elapsed time]
Last Action: [timestamp]
Next Action: [pending action]
```

**Reference:** [EXECUTION_MODEL.md - States](EXECUTION_MODEL.md#runtime-states)

---

### KDSE Report

Generate or retrieve the current Runtime Report.

**Prompt:**
```
KDSE Report
```

**Prompt with Format:**
```
KDSE Report

Format: [full | summary | delta]
```

**Example:**
```
KDSE Report

Format: full
```

**What Returns:**
Full Runtime Report per [REPORT_SPEC.md](REPORT_SPEC.md)

---

### Close KDSE

End the current session.

**Prompt:**
```
Close KDSE

Session ID: [session ID]
Reason: [optional reason]
```

**Example:**
```
Close KDSE

Session ID: KDSE-RT-2026-07-10-001
Reason: Target maturity reached
```

**What Happens:**
1. Runtime finalizes report
2. Records session metrics
3. Archives session artifacts
4. Presents completion summary

**Reference:** [SESSION_PROTOCOL.md - Completion](SESSION_PROTOCOL.md#phase-7-session-completion)

---

## Approval Commands

### Approve Recommendation

Approve the current recommendation.

**Prompt:**
```
Approve

Recommendation: [action ID]
Reason: [optional rationale]
```

**Example:**
```
Approve

Recommendation: KDSE-ACT-001
Reason: Aligns with Q3 roadmap priorities
```

**What Happens:**
1. Runtime records approval
2. Proceeds to implementation
3. Executes approved action

---

### Approve with Modifications

Approve the recommendation with specified changes.

**Prompt:**
```
Approve with Modifications

Recommendation: [action ID]
Modifications:
1. [change 1]
2. [change 2]
Reason: [rationale]
```

**Example:**
```
Approve with Modifications

Recommendation: KDSE-ACT-001
Modifications:
1. Limit scope to Phase 1 requirements only
2. Defer non-functional requirements to later session
Reason: Resource constraints this quarter
```

**What Happens:**
1. Runtime records approval with changes
2. Proceeds to implementation with modifications
3. Documents changes

---

### Reject Recommendation

Reject the current recommendation.

**Prompt:**
```
Reject

Recommendation: [action ID]
Reason: [why rejected]
```

**Example:**
```
Reject

Recommendation: KDSE-ACT-001
Reason: Priority shifted to security hardening this quarter
```

**What Happens:**
1. Runtime records rejection
2. Identifies next recommendation
3. Presents alternative

---

### Defer Decision

Postpone the decision to a later session.

**Prompt:**
```
Defer

Recommendation: [action ID]
Reason: [why deferring]
Resume After: [optional date or event]
```

**Example:**
```
Defer

Recommendation: KDSE-ACT-001
Reason: Waiting for architecture decision from leadership
Resume After: Architecture review meeting (July 15)
```

**What Happens:**
1. Runtime pauses session
2. Preserves state
3. Awaits resume command

---

## Query Commands

### View Findings

Display current audit findings.

**Prompt:**
```
KDSE Findings

Filter: [all | critical | high | medium | low]
Dimension: [specific dimension, optional]
```

**Example:**
```
KDSE Findings

Filter: high
Dimension: Knowledge Artifacts
```

---

### View Scores

Display current compliance scores.

**Prompt:**
```
KDSE Scores

Format: [table | comparison | delta]
Compare To: [previous session ID, optional]
```

**Example:**
```
KDSE Scores

Format: comparison
Compare To: KDSE-RT-2026-07-09-001
```

---

### View Progress

Display progress toward target maturity.

**Prompt:**
```
KDSE Progress
```

**What Returns:**
```
Target: 7.0/10
Current: 5.2/10
Baseline: 4.8/10
Improvement: +0.4 (since session start)
Remaining: 1.8 points
Progress: 18% of journey complete
```

---

### View Session History

Display session execution history.

**Prompt:**
```
KDSE History

Session ID: [session ID]
```

**Example:**
```
KDSE History

Session ID: KDSE-RT-2026-07-10-001
```

---

## Reference Templates

### Full Session Template

Complete template for starting a KDSE session.

```
Run KDSE

## Session Parameters
Repository: [full path or URL]
Target Maturity: [0-10]
Operator: [name]
Scope: [dimensions, optional]
Constraints: [limits, optional]

## Context
Previous Session: [session ID if continuing]
Last Assessment: [date if recent]
Known Issues: [optional notes]
```

---

### Approval Template

Template for recording approval decisions.

```
## Approval Record

Recommendation ID: [ID]
Recommended Action: [summary]
Expected Impact: [+X points]
Operator: [name]
Decision: [APPROVE | APPROVE WITH MODIFICATIONS | REJECT | DEFER]
Timestamp: [ISO 8601]

[If modified:]
Modifications:
1. [change]
2. [change]

[If rejected:]
Reason: [explanation]

[If deferred:]
Reason: [explanation]
Resume After: [condition]
```

---

## Usage Examples

### Example 1: New Project Assessment

**Step 1: Start Session**
```
Run KDSE

Repository: /workspace/project/new-service
Target Maturity: 6.0
Operator: Jane Developer
```

**Step 2: Review Report**
```
KDSE Report
```

**Step 3: Approve Action**
```
Approve

Recommendation: KDSE-ACT-001
Reason: Foundation work is priority
```

**Step 4: Continue**
```
Continue KDSE

Session ID: [session ID from step 1]
```

**Step 5: Complete**
```
Close KDSE

Session ID: [current session ID]
Reason: Target reached
```

---

### Example 2: Targeted Improvement

**Step 1: Start Focused Session**
```
Run KDSE

Repository: /workspace/project/existing-service
Target Maturity: 7.5
Operator: John Engineer
Scope: Verification Practices
```

**Step 2: Approve**
```
Approve with Modifications

Recommendation: KDSE-ACT-005
Modifications:
1. Focus on unit test coverage only
2. Skip integration tests for now
```

**Step 3: Close**
```
Close KDSE

Session ID: [session ID]
Reason: Focused scope complete
```

---

### Example 3: Pause and Resume

**Step 1: Start Session**
```
Run KDSE

Repository: /workspace/project/large-service
Target Maturity: 8.0
Operator: Jane Developer
```

**Step 2: Defer Decision**
```
Defer

Recommendation: KDSE-ACT-003
Reason: Waiting for architecture review
Resume After: July 20
```

**Step 3: Resume Later**
```
Continue KDSE

Session ID: [session ID]
```

---

## Quick Reference Card

```
╔═══════════════════════════════════════════════════════════════╗
║                    KDSE COMMAND QUICK REFERENCE               ║
╠═══════════════════════════════════════════════════════════════╣
║                                                               ║
║  START SESSION                                               ║
║  Run KDSE                                                    ║
║  Repository: [path] Target: [score] Operator: [name]         ║
║                                                               ║
║  VIEW STATE                                                  ║
║  KDSE Status                                                 ║
║  KDSE Report                                                 ║
║  KDSE Progress                                               ║
║                                                               ║
║  DECISIONS                                                   ║
║  Approve [action ID]                                         ║
║  Approve with Modifications [action ID]                      ║
║  Reject [action ID] Reason: [why]                           ║
║  Defer [action ID] Resume After: [when]                      ║
║                                                               ║
║  END SESSION                                                 ║
║  Close KDSE Session ID: [ID] Reason: [why]                   ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## Notes

1. **Prompts are templates**: Adapt parameters as needed for your context
2. **Reference the Standard**: These prompts invoke workflows defined in KDSE documentation
3. **Human authorization required**: Implementation always requires operator approval
4. **Session persistence**: Use Session ID to resume paused sessions

---

*This document is an informative reference implementation. It provides command templates for operating the KDSE Runtime.*
