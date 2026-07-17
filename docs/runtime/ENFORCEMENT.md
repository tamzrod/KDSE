# KDSE Enforcement - Runtime Guard for KDSE Principles

## Overview

KDSE Enforcement is the **runtime guard** that makes KDSE principles automatic, not just suggestions. It blocks premature coding, enforces foundation-first methodology, and provides automatic audit warnings.

## Core Principle

> **Knowledge Base + Foundation MUST be built first**
> **Architecture before Implementation**
> **Repository-first Analysis**

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         KDSE MCP Server                                 │
├─────────────────────────────────────────────────────────────────────────┤
│  execute() ──► ┌────────────────┐ ──► ┌────────────────────────────┐  │
│                │  Session Guard │     │     Enforcement Engine     │  │
│                │                │     │                            │  │
│                │ • IsInit?      │     │ • Foundation Check         │  │
│                │ • Session OK?  │     │ • Knowledge Check          │  │
│                │ • Not Expired? │     │ • Architecture Check       │  │
│                └───────┬────────┘     │ • Repository Analysis      │  │
│                        │              └─────────────┬──────────────┘  │
│                        ▼                            ▼                  │
│               ┌────────────────┐          ┌───────────────────────┐    │
│               │    BLOCKED     │          │   Allowed / Warning   │    │
│               │                │          │                       │    │
│               │ Return error   │          │ Generate WorkOrder    │    │
│               │ Show violations│          │ Auto-create if needed│    │
│               │ Suggest fix    │          │ Show audit warnings  │    │
│               └────────────────┘          └───────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Enforcement Levels

| Level | Behavior |
|-------|----------|
| `off` | No enforcement - all operations allowed |
| `warning` | Warn but allow - violations shown but not blocking |
| `strict` | Block violations - prevents implementation without foundation |
| `hard` | Block everything - even --force flag is ignored |

## Error Codes

| Code | Message | Severity |
|------|---------|----------|
| `KDSE_ENF_001` | Foundation directory does not exist | BLOCKER |
| `KDSE_ENF_002` | Knowledge directory does not exist | BLOCKER |
| `KDSE_ENF_003` | Architecture document does not exist | BLOCKER |
| `KDSE_ENF_004` | Premature implementation detected | BLOCKER |
| `KDSE_ENF_005` | Repository has not been analyzed | MEDIUM |
| `KDSE_ENF_006` | Foundation incomplete (missing documents) | BLOCKER |
| `KDSE_ENF_007` | Knowledge incomplete (empty categories) | HIGH |
| `KDSE_ENF_008` | Phase violation | Depends |
| `KDSE_ENF_009` | Compliance violation | Depends |

## Required Artifacts

### Foundation Documents (Must Exist)

```
.kdse/foundation/
├── PROBLEM.md         # Problem statement and impact
├── SPEC.md            # Project specification
├── REQUIREMENTS.md    # Functional requirements
├── ASSUMPTIONS.md     # Key assumptions
├── CONSTRAINTS.md     # Project constraints
└── ARCHITECTURE.md    # System architecture (required before implementation)
```

### Knowledge Categories (Must Have Content)

```
.kdse/knowledge/
├── general/           # General domain knowledge
├── operational/       # Operational procedures and patterns
├── developmental/     # Development-related knowledge
└── repository-analysis.md  # Repository structure analysis
```

## Phase Enforcement

### Phase Order

```
Problem → Knowledge → Foundation → Audit → Assessment → Architecture → Implementation → Complete
```

### Implementation Gate

Implementation is **BLOCKED** unless ALL of the following are true:

1. ✅ `.kdse/` workspace exists
2. ✅ All foundation documents exist and are populated (>100 bytes)
3. ✅ All knowledge categories exist and have content
4. ✅ `ARCHITECTURE.md` exists and is populated (>200 bytes)
5. ✅ `repository-analysis.md` exists

### Phase-Specific Checks

| Phase | Required | Checked |
|-------|----------|---------|
| Problem | Workspace initialized | Yes |
| Knowledge | Problem defined | Yes |
| Foundation | Knowledge collected | Yes |
| Audit | Foundation complete | Yes |
| Assessment | Audit complete | Yes |
| Architecture | Assessment complete | Yes |
| Implementation | All above + Architecture | Yes (STRICT) |

## Automatic Actions

When `autoCreate: true` (default):

### Missing Foundation Documents
Automatically creates template documents:
- `PROBLEM.md` - Template with sections for description, impact, criteria
- `SPEC.md` - Template with overview, scope, deliverables
- `REQUIREMENTS.md` - Template with FR-XXX format
- `ASSUMPTIONS.md` - Template with assumption categories
- `CONSTRAINTS.md` - Template with constraint categories

### Missing Knowledge Categories
Automatically creates:
- `general/` directory with README
- `operational/` directory with README
- `developmental/` directory with README

## Audit Warnings

Non-blocking warnings issued for:

1. **Premature Phase Request**
   - Implementation requested during Problem/Knowledge/Foundation phase
   - Example: "WARNING: Implementation requested during Foundation phase"

2. **Missing Repository Analysis**
   - Repository structure not documented
   - Severity: MEDIUM (blocks in hard mode)

3. **Empty Documents**
   - Foundation documents exist but are templates
   - Severity: HIGH

## Usage Examples

### Python

```python
from kdse_shttp_tools import shttp_initialize, shttp_execute, shttp_status

# Initialize workspace
shttp_initialize()

# Execute with enforcement
result = shttp_execute("Build a payroll system with biometrics")

if result['blocked']:
    print("❌ BLOCKED")
    print(f"Reason: {result['blocked_reason']}")
    print("Violations:")
    for v in result['violations']:
        print(f"  - {v['code']}: {v['message']}")
    print("\nNext steps:")
    for step in result['next_steps']:
        print(f"  → {step}")
else:
    print("✓ Proceed with implementation")
```

### HTTP

```bash
# Initialize
curl -X POST http://localhost:8080/initialize

# Execute (blocked without foundation)
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"objective": "Build payroll system"}'

# Check foundation status
curl http://localhost:8080/foundation

# Create foundation
curl -X POST http://localhost:8080/foundation \
  -d '{"documents": ["PROBLEM.md", "SPEC.md"]}'

# Run audit
curl http://localhost:8080/audit
```

## Test Bench: Payroll System

To test enforcement with the payroll system example:

```bash
# 1. Start fresh
kdse initialize

# 2. Try to implement (should block)
kdse execute "Build a payroll system with biometrics"
# Expected: BLOCKED - Foundation incomplete

# 3. Complete phases properly
kdse foundation create
kdse knowledge create
kdse analyze --repository
kdse phase transition foundation

# 4. Try again
kdse execute "Build a payroll system with biometrics"
# Expected: Allowed (or warnings if documents not populated)
```

## Compliance Report Format

```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "enforcement_level": "strict",
  "blocked": true,
  "can_proceed": false,
  "violations": [
    {
      "code": "KDSE_ENF_006",
      "message": "Foundation incomplete: missing documents: [PROBLEM.md, SPEC.md]",
      "severity": "BLOCKER",
      "required_phase": "Foundation",
      "corrective_action": "Create: [PROBLEM.md, SPEC.md]",
      "blocked": true
    }
  ],
  "warnings": [],
  "auto_actions": ["Created: PROBLEM.md", "Created: SPEC.md"]
}
```

## Integration Points

### With Session Guard
- Enforcer uses guard to check initialization
- Blocked if workspace not initialized

### With Orchestration Engine
- Enforcer reports to orchestrator
- Phase transitions trigger re-enforcement

### With Compliance
- Enforcer violations feed into compliance report
- Compliance check includes enforcement status

## Best Practices

1. **Always Initialize First**
   ```bash
   kdse initialize
   ```

2. **Check Status Before Implementation**
   ```bash
   kdse status
   ```

3. **Use Audit to Find Gaps**
   ```bash
   kdse audit
   ```

4. **Let Auto-Create Handle Templates**
   - Don't manually create empty documents
   - Let the system create templates with proper structure

5. **Populate Templates Before Implementation**
   - Auto-created documents are templates
   - Fill them with real content before proceeding

## Troubleshooting

### "Foundation directory does not exist"
```bash
kdse initialize
kdse foundation create
```

### "Knowledge base is empty"
```bash
kdse knowledge create
kdse analyze --repository
```

### "Architecture document does not exist"
```bash
kdse phase transition architecture
# Then create ARCHITECTURE.md
```

### "Still blocked after creating files"
- Check file sizes - files must be >100 bytes
- Run `kdse compliance` for detailed report

## Future Enhancements

- [ ] LLM-based document quality checking
- [ ] Automatic evidence collection
- [ ] Phase transition approval workflow
- [ ] Custom enforcement rules per project
- [ ] Integration with CI/CD pipelines

---

*This document is part of the KDSE Runtime specification.*
