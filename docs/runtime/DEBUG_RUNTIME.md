# KDSE Debug Runtime Architecture

**Document Version:** 1.0  
**Type:** Informative Implementation  
**Effective Date:** 2026-07-11

---

## Purpose

This document defines the architecture for the KDSE Debug Runtime—a deterministic engineering debugging system that transforms debugging from endless hypothesis generation into evidence-driven root cause analysis.

---

## Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         KDSE Runtime                                    │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │                    Debug Engine (Core)                            │  │
│  │                                                                  │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │  │
│  │  │  Evidence    │  │  Hypothesis  │  │  Confidence  │           │  │
│  │  │  Collector   │  │  Manager     │  │  Tracker     │           │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘           │  │
│  │                                                                  │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │  │
│  │  │  Loop       │  │  Root Cause  │  │  Report      │           │  │
│  │  │  Detector   │  │  Selector    │  │  Generator   │           │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘           │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  Commands: kdse debug | kdse verify | kdse repair | kdse audit        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Design Principles

### 1. Evidence Over Speculation

The Debug Runtime prioritizes evidence collection over hypothesis generation. No hypothesis shall be generated without supporting evidence.

### 2. Confidence-Driven Decision Making

Root cause selection occurs when confidence exceeds a configurable threshold (default: 90%). The runtime SHALL NOT proceed to implementation until confidence threshold is met.

### 3. Loop Prevention

Repeated investigation of the same component is detected and prevented. Previous evidence is referenced, not re-collected.

### 4. Structured Artifacts

All debugging artifacts are stored as structured data for traceability and reproducibility.

---

## Debug Workflow

```
┌─────────────────┐
│  Failure        │
│  Detected       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Evidence       │
│  Collection     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Hypothesis     │
│  Generation     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Evidence       │
│  Evaluation     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Confidence     │◄────────┐
│  Assessment     │         │
└────────┬────────┘         │
         │                  │
    ┌────┴────┐           │
    │ >= 90%? │           │
    └────┬────┘           │
    YES  │  NO            │
    │    │                │
    ▼    └────────────────┘
┌─────────────────┐
│  Root Cause     │
│  Selection      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Implementation │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Verification   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Regression     │
│  Tests          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Runtime        │
│  Report         │
└─────────────────┘
```

---

## Components

### 1. Evidence Collector

**Purpose:** Collect and store evidence systematically.

**Evidence Types:**

| Type | Description | Example |
|------|-------------|---------|
| `exception` | Exception messages and traces | `NullPointerException at line 42` |
| `test_failure` | Test failure output | `TestCase.assertEquals failed: expected X, got Y` |
| `log` | Runtime logs | `ERROR: Connection timeout after 30s` |
| `source` | Source code inspection | `File X imports Y which creates Z` |
| `state` | Repository state | `git diff shows changes in file A` |
| `config` | Configuration values | `DATABASE_URL=sqlite:///:memory:` |
| `dependency` | Dependency versions | `package.json: "lodash": "^4.17.21"` |

**Evidence Schema:**

```yaml
evidence:
  id: "E-001"
  type: "exception"
  timestamp: "2026-07-11T12:00:00Z"
  source:
    file: "src/repository.py"
    line: 42
    function: "get_user"
  content: |
    SQLite BusyError: database is locked
  collected_by: "evidence_collector"
  tags: ["database", "concurrency", "sqlite"]
```

### 2. Hypothesis Manager

**Purpose:** Manage hypotheses with evidence tracking.

**Hypothesis Schema:**

```yaml
hypothesis:
  id: "H-001"
  description: "Nested SQLite connections create database lock"
  status: "active"  # active | evaluated | rejected | selected
  confidence:
    initial: 40
    current: 92
    threshold: 90
  supporting_evidence:
    - "E-001"
    - "E-003"
  contradicting_evidence:
    - "E-005"
  affected_components:
    - "BookRepository"
    - "AuthorRepository"
  experiments:
    - id: "X-001"
      action: "Added connection timeout"
      result: "pass"
      confidence_delta: +15
  created_at: "2026-07-11T12:05:00Z"
  updated_at: "2026-07-11T12:30:00Z"
```

**Hypothesis Lifecycle:**

```
[NEW] ──Evaluate──▶ [ACTIVE] ──Confidence >= 90%──▶ [SELECTED]
                     │                                  │
                     │                        ┌─────────┘
                     │                        │
                     └──Conflicting Evidence──▶ [REJECTED]
                     │                        │
                     └──Max Hyps Exceeded──────┘
```

### 3. Confidence Tracker

**Purpose:** Calculate and track confidence scores.

**Confidence Rules:**

| Evidence Type | Confidence Impact |
|---------------|------------------|
| Direct error trace | +20% |
| Test failure confirming hypothesis | +15% |
| Log message indicating failure | +10% |
| Source code confirms mechanism | +10% |
| Configuration enables failure | +5% |
| Contradicting evidence | -25% |
| Partial evidence | +5% |

**Confidence Calculation:**

```
confidence = min(100, initial_confidence + sum(supporting) - sum(contradicting))
```

### 4. Loop Detector

**Purpose:** Detect and prevent repeated investigation.

**Detectable Patterns:**

| Pattern | Detection Method | Action |
|---------|-----------------|--------|
| Repeated file inspection | Hash of file + content | Reference cached evidence |
| Repeated module reload | Module import timestamp | Skip reload |
| Repeated schema check | Schema version + hash | Use cached result |
| Repeated database inspection | Query + result hash | Return cached data |
| Repeated cache clear | Action timestamp | Warn and skip |

**Loop Detection Schema:**

```yaml
loop_record:
  pattern_id: "L-001"
  pattern_type: "file_inspection"
  key: "src/repository.py:abc123"
  first_occurrence: "2026-07-11T12:00:00Z"
  last_occurrence: "2026-07-11T12:30:00Z"
  count: 5
  cached_evidence: "E-001"
```

### 5. Root Cause Selector

**Purpose:** Select root cause when confidence threshold is met.

**Selection Criteria:**

1. Confidence >= threshold (default 90%)
2. No active hypotheses with higher confidence
3. Evidence has been evaluated
4. Operator approval received

**Root Cause Report Schema:**

```yaml
root_cause_report:
  id: "RC-001"
  session_id: "DEBUG-20260711-120000"
  failure:
    summary: "Database locked during concurrent access"
    evidence_count: 5
  selected_hypothesis: "H-001"
  confidence: 92
  evidence_summary:
    supporting: 3
    contradicting: 1
  recommended_fix:
    description: "Implement connection pooling with timeout"
    files: ["src/repository.py", "src/config.py"]
    changes: [...]
  alternative_explanations:
    - id: "H-002"
      confidence: 45
      reason_not_selected: "Lower confidence than H-001"
  created_at: "2026-07-11T12:35:00Z"
  operator_approved: true
  approved_at: "2026-07-11T12:36:00Z"
```

### 6. Report Generator

**Purpose:** Produce structured debugging reports.

**Debug Session Report Schema:**

```yaml
debug_session_report:
  session_id: "DEBUG-20260711-120000"
  started_at: "2026-07-11T12:00:00Z"
  completed_at: "2026-07-11T12:45:00Z"
  duration_seconds: 2700
  
  failure:
    summary: "Database locked during concurrent access"
    severity: "high"
    reproducibility: "consistent"
  
  root_cause:
    id: "RC-001"
    description: "Nested SQLite connections create database lock"
    confidence: 92
  
  evidence:
    collected: 7
    evaluated: 7
    supporting: 4
    contradicting: 2
  
  hypotheses:
    generated: 3
    active: 0
    rejected: 2
    selected: 1
  
  experiments:
    performed: 4
    passed: 3
    failed: 1
  
  implementation:
    files_modified: 2
    lines_added: 45
    lines_removed: 12
  
  verification:
    unit_tests: "passed"
    integration_tests: "passed"
    regression_tests: "passed"
    original_failure: "resolved"
  
  engineering_time:
    evidence_collection: 600
    hypothesis_evaluation: 900
    implementation: 600
    verification: 600
    total: 2700
  
  status: "completed"
```

---

## Directory Structure

```
.dkdse/
├── debug/
│   ├── engine.sh              # Core debug engine
│   ├── evidence/
│   │   ├── collector.sh       # Evidence collection
│   │   └── store.json         # Evidence store
│   ├── hypotheses/
│   │   ├── manager.sh         # Hypothesis management
│   │   └── registry.json      # Hypothesis registry
│   ├── confidence/
│   │   └── tracker.sh         # Confidence calculation
│   ├── loops/
│   │   ├── detector.sh        # Loop detection
│   │   └── history.json       # Loop history
│   ├── reports/
│   │   ├── generator.sh       # Report generation
│   │   └── root-cause/        # Root cause reports
│   └── sessions/
│       └── DEBUG-*.json       # Debug session files
└── bootstrap/
    └── debug-config.yaml      # Debug configuration
```

---

## Debug Configuration

**Location:** `.kdse/bootstrap/debug-config.yaml`

```yaml
debug_runtime:
  version: "1.0"
  
confidence:
  threshold: 90
  minimum_initial: 20
  maximum_initial: 60
  
hypothesis:
  max_active: 5
  auto_reject_below: 20
  
evidence:
  collection_timeout_seconds: 300
  auto_categorize: true
  
loops:
  detection_enabled: true
  max_repetitions: 3
  cooldown_seconds: 60
  
reports:
  output_format: "json"
  include_evidence: true
  include_hypotheses: true
  
implementation:
  require_operator_approval: true
  auto_backup: true
```

---

## Command Integration

### `kdse debug`

Primary debugging command using the full Debug Engine workflow.

```
kdse debug --failure "test failure" --evidence-type exception,log
```

### `kdse verify`

Reuses Debug Engine for verification failures.

```
kdse verify --test "test_name" --expected "value"
```

### `kdse repair`

Reuses Debug Engine for automated repairs.

```
kdse repair --issue "configuration" --auto-fix
```

### `kdse audit`

Reuses Debug Engine for audit finding investigations.

```
kdse audit --finding "F-001" --investigate
```

---

## State Machine

```
┌──────────────────────────────────────────────────────────────────────┐
│                        DEBUG STATES                                   │
├──────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  [INITIAL] ──Start Debug──▶ [EVIDENCE_COLLECTION]                     │
│                                       │                               │
│                                       ▼                               │
│                              [HYPOTHESIS_GENERATION]                   │
│                                       │                               │
│                                       ▼                               │
│                              [EVIDENCE_EVALUATION]                    │
│                                       │                               │
│                                       ▼                               │
│                              [CONFIDENCE_ASSESSMENT]                  │
│                                       │                               │
│                         ┌─────────────┴─────────────┐                  │
│                         │                           │                  │
│                    Confidence >= 90%          Confidence < 90%        │
│                         │                           │                  │
│                         ▼                           │                  │
│                    [ROOT_CAUSE_SELECTED]            │                  │
│                         │                           │                  │
│                    Operator Approved?                 │                  │
│                    ┌───┴───┐                        │                  │
│                    │       │                        │                  │
│                   YES     NO                        │                  │
│                    │       │                        │                  │
│                    ▼       ▼                        │                  │
│            [IMPLEMENTING]  [WAITING_APPROVAL]──────┘                  │
│                    │                           │                        │
│                    ▼                           │                        │
│             [VERIFICATION]                      │                        │
│                    │                           │                        │
│              Tests Pass?                         │                        │
│              ┌───┴───┐                         │                        │
│              │       │                         │                        │
│             YES     NO                         │                        │
│              │       │                         │                        │
│              ▼       └─────────────────────────┘                        │
│       [REGRESSION_TESTS]      [EVIDENCE_COLLECTION]                   │
│              │                                                     │
│              ▼                                                     │
│       [COMPLETED]                                                   │
│                                                                       │
└──────────────────────────────────────────────────────────────────────┘
```

---

## API Reference

### Core Functions

```bash
# Initialize debug session
debug_init --failure "description" [--severity high|medium|low]

# Collect evidence
debug_collect_evidence --type exception --content "..." --source "file:line"

# Generate hypothesis
debug_new_hypothesis --description "..." --confidence 40 --components "A,B"

# Evaluate evidence against hypothesis
debug_evaluate --hypothesis-id H-001 --evidence-id E-001 --impact positive|negative

# Check confidence threshold
debug_check_confidence --hypothesis-id H-001

# Select root cause
debug_select_root_cause --hypothesis-id H-001

# Detect loops
debug_check_loop --pattern file_inspection --key "file.py"

# Generate report
debug_generate_report [--format json|markdown]

# Complete session
debug_complete [--status success|failed]
```

---

## Implementation Notes

### Phase 1: Core Engine
- Evidence collection with categorization
- Hypothesis creation and management
- Basic confidence tracking

### Phase 2: Intelligence
- Loop detection
- Evidence correlation
- Confidence optimization

### Phase 3: Integration
- Command integration (debug, verify, repair, audit)
- Report generation
- Operator approval workflow

---

## Success Criteria

| Criterion | Measurement |
|-----------|-------------|
| Evidence drives decisions | 0 hypotheses without evidence |
| Confidence-based selection | Root cause selected at >= 90% |
| Loop prevention | 0 repeated investigations |
| Operator approval | 100% of implementations approved |
| Verification required | 0 fixes without verification |

---

## References

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Runtime architecture
- [SESSION_PROTOCOL.md](./SESSION_PROTOCOL.md) - Session lifecycle
- [REPORT_SPEC.md](./REPORT_SPEC.md) - Report format

---

*This document describes the Debug Runtime implementation architecture.*
