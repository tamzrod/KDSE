# KDSE Runtime

**Document Version:** 1.0  
**Type:** Informative Reference Implementation  
**Effective Date:** 2026-07-10

---

## Purpose

The KDSE Runtime is the official reference implementation of Knowledge-Driven Software Engineering. It provides a standardized operational workflow for applying KDSE principles during day-to-day engineering work.

The Runtime answers the question: **"How do I operate a software engineering session using KDSE?"**

---

## Relationship to the KDSE Standard

```
┌─────────────────────────────────────────────────────────────┐
│                     KDSE Standard                           │
│                                                             │
│  • Principles (003-core-principles.md)                     │
│  • Authority Hierarchy (006-chain-of-authority.md)          │
│  • Engineering Model (004-engineering-model.md)            │
│  • Audit Standards (docs/audit/)                           │
│  • Glossary (007-glossary.md)                              │
│                                                             │
│  ⚠️ NORMATIVE: Defines what must be true                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Consumes
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      KDSE Runtime                           │
│                                                             │
│  • Execution Model                                         │
│  • Session Protocol                                        │
│  • Runtime Reports                                         │
│  • Workflows                                               │
│  • Prompts                                                 │
│                                                             │
│  ℹ️ INFORMATIVE: Reference implementation, not requirement   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Operates
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Human Engineer / Operator                  │
│                  / AI Assistant / CI Pipeline               │
└─────────────────────────────────────────────────────────────┘
```

---

## Normative vs Informative

### Normative (The KDSE Standard)

The following are **normative**—they define requirements that KDSE-compliant systems must meet:

| Document | Purpose |
|----------|---------|
| [docs/foundation/003-core-principles.md](../docs/foundation/003-core-principles.md) | Core principles |
| [docs/foundation/006-chain-of-authority.md](../docs/foundation/006-chain-of-authority.md) | Authority hierarchy |
| [docs/audit/COMPLIANCE_AUDIT.md](../docs/audit/COMPLIANCE_AUDIT.md) | Compliance evaluation |
| [docs/audit/FOUNDATION_AUDIT.md](../docs/audit/FOUNDATION_AUDIT.md) | Foundation evaluation |
| [docs/audit/AUDIT_SCORING.md](../docs/audit/AUDIT_SCORING.md) | Scoring methodology |

### Informative (The KDSE Runtime)

The following are **informative**—they are reference implementations, not requirements:

| Document | Purpose |
|----------|---------|
| [EXECUTION_MODEL.md](EXECUTION_MODEL.md) | Operational workflow reference |
| [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) | Session lifecycle reference |
| [REPORT_SPEC.md](REPORT_SPEC.md) | Report format reference |
| [PROMPTS.md](PROMPTS.md) | Command prompts reference |
| [WORKFLOW.md](WORKFLOW.md) | Workflow diagrams reference |
| [examples/](examples/) | Example usage |

---

## What the Runtime Is NOT

The Runtime does NOT:

- ❌ Redefine KDSE principles
- ❌ Redefine audits (references them)
- ❌ Redefine scoring (references it)
- ❌ Redefine governance
- ❌ Redefine chain of authority
- ❌ Introduce new normative requirements
- ❌ Modify the KDSE Standard

The Runtime does:

- ✅ Orchestrate existing KDSE artifacts
- ✅ Provide an operational workflow
- ✅ Generate actionable recommendations
- ✅ Standardize engineering sessions
- ✅ Reference the Standard

---

## Quick Start

### Option 1: Normalize Existing Documentation

If your repository already has documentation:

1. **Open your repository** in your preferred environment

2. **Enter the command:**
   ```
   kdse normalize
   ```

3. **The Runtime will:**
   - Discover existing documentation
   - Analyze and extract knowledge
   - Generate KDSE-standard artifacts
   - Build full traceability
   - Produce a normalization report

4. **Review the generated artifacts** in `.kdse/normalized/`

5. **Continue with** `kdse run` or other commands

### Option 2: Start Immediately

If you want to start working right away:

1. **Open your repository** in your preferred environment

2. **Enter the command:**
   ```
   kdse run
   ```

3. **The Runtime will:**
   - Load the KDSE Standard
   - Run Foundation Verification
   - Run Repository Assessment
   - Execute Compliance Audit
   - Generate Runtime Report
   - Recommend next action

4. **Review the Runtime Report** and approve or modify the recommended action

5. **Implement** the approved action

6. **Verify** the results

7. **Repeat** or **close** the session

### Session Commands

| Command | Purpose |
|---------|---------|
| `kdse normalize` | Normalize existing documentation to KDSE standard |
| `kdse run` | Start a new KDSE session |
| `kdse continue` | Resume an existing session |
| `kdse status` | View current session state |
| `kdse report` | Generate runtime report |

---

## Key Concepts

### Runtime

The operational component that orchestrates KDSE sessions. The Runtime consumes the Standard and produces actionable guidance.

### Session

A bounded period of KDSE activity. Sessions have defined start, execution, and completion phases.

### Operator

The human or system executing KDSE sessions. Operators work with the Runtime to achieve engineering goals.

### Workflow

The sequence of operations the Runtime performs. Workflows reference, not replace, the Standard.

### Runtime Report

A summary document produced by the Runtime. Reports summarize audit findings, not replace audits.

---

## Document Overview

| Document | Description |
|----------|-------------|
| [README.md](README.md) | This overview |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Runtime architecture and structure |
| [EXECUTION_MODEL.md](EXECUTION_MODEL.md) | Runtime lifecycle and states |
| [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) | Session lifecycle details |
| [REPORT_SPEC.md](REPORT_SPEC.md) | Runtime Report specification |
| [COMMANDS.md](COMMANDS.md) | Command interface definition |
| [CONFORMANCE.md](CONFORMANCE.md) | Conformance criteria |
| [VERSIONING.md](VERSIONING.md) | Versioning strategy |
| [PROMPTS.md](PROMPTS.md) | Copy-paste command prompts |
| [WORKFLOW.md](WORKFLOW.md) | Visual workflow diagrams |
| [examples/](examples/) | Example sessions and reports |

---

## Design Principles

1. **Standard-First**: The Runtime always references the Standard, never replaces it
2. **Evidence-Based**: All recommendations trace to audit evidence
3. **Human-Authorized**: Operators approve all implementations
4. **Technology-Neutral**: Works with any repository, tooling, or environment
5. **AI-Independent**: Operable by humans or AI assistants without change
6. **Repository-Independent**: Applies to any KDSE-compliant repository

---

## Terminology Note

The Runtime uses neutral terminology to remain AI-independent:

| Avoid | Use Instead |
|-------|------------|
| Agent | Runtime |
| Autonomous | Operated |
| LLM / Claude / GPT | AI Assistant |
| Chat | Conversation |

---

## Standards Reference

The Runtime references these authoritative KDSE documents:

- [KDSE Foundation](../docs/foundation/)
- [KDSE Audit System](../docs/audit/)
- [KDSE Engineering Model](../docs/foundation/004-engineering-model.md)
- [KDSE Chain of Authority](../docs/foundation/006-chain-of-authority.md)

---

*This Runtime is an informative reference implementation. It demonstrates how KDSE can be applied in practice while preserving the normative authority of the KDSE Standard.*
