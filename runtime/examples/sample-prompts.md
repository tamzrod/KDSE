# KDSE Runtime Prompt Examples

**Document Version:** 1.0  
**Type:** Illustrative Examples  
**Effective Date:** 2026-07-10

---

## Basic Session

### Start Session

```
Run KDSE

Repository: /workspace/project/myapp
Target Maturity: 6.5
Operator: Jane Developer
```

### View Status

```
KDSE Status
```

### View Report

```
KDSE Report
```

### Close Session

```
Close KDSE

Session ID: KDSE-RT-2026-07-10-001
Reason: Target maturity reached
```

---

## Approval Sequence

### Approve

```
Approve

Recommendation: KDSE-ACT-001
Reason: Aligns with sprint goals
```

### Approve with Modifications

```
Approve with Modifications

Recommendation: KDSE-ACT-001
Modifications:
1. Limit scope to Phase 1 only
2. Use existing template format
Reason: Resource constraints
```

### Reject

```
Reject

Recommendation: KDSE-ACT-001
Reason: Priority shifted to security work
```

### Defer

```
Defer

Recommendation: KDSE-ACT-001
Reason: Waiting for architecture decision
Resume After: Architecture review meeting
```

---

## Query Examples

### View Progress

```
KDSE Progress
```

### View Scores

```
KDSE Scores

Format: comparison
Compare To: KDSE-RT-2026-07-09-001
```

### View Findings

```
KDSE Findings

Filter: critical
Dimension: Knowledge Artifacts
```

---

## Resume Session

```
Continue KDSE

Session ID: KDSE-RT-2026-07-10-001
```

---

*See [PROMPTS.md](../../PROMPTS.md) for the complete command reference.*
