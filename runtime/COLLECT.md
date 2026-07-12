# KDSE Knowledge Collection Runtime

**Document Version:** 1.0  
**Type:** Informative Reference Implementation  
**Effective Date:** 2026-07-12

---

## Purpose

The KDSE Knowledge Collection Runtime provides a structured mechanism for collecting missing engineering knowledge before implementation. Engineering projects are knowledge-driven, and implementation quality depends on the completeness of the engineering knowledge available to the project.

---

## Design Philosophy

### What kdse collect IS

`kdse collect` is a **Knowledge Acquisition workflow**. Its purpose is to help the operator:

1. **Identify** missing engineering knowledge
2. **Collect** knowledge from available sources
3. **Normalize** knowledge into KDSE Documentation Standard artifacts
4. **Integrate** knowledge into the KDSE knowledge base

### What kdse collect IS NOT

- `kdse collect` is NOT a code generator
- `kdse collect` is NOT an internet scraper
- `kdse collect` does NOT replace domain expertise

---

## Workflow

### When to Run

The operator may execute `kdse collect` at any time. It is completely independent.

After `kdse audit` or `kdse normalize`, the runtime may recommend running `kdse collect` when knowledge gaps are detected. The recommendation is informational only—the operator always decides.

### Collection Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  1. ANALYZE KNOWLEDGE GAPS                                  │   │
│  │     • Audit findings                                        │   │
│  │     • Normalization results                                 │   │
│  │     • Existing knowledge artifacts                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  2. COLLECT FROM SOURCES                                    │   │
│  │     • Repository documentation                              │   │
│  │     • Standards documents                                   │   │
│  │     • Vendor manuals                                        │   │
│  │     • Operator knowledge                                    │   │
│  │     • AI-assisted acquisition                               │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  3. NORMALIZE TO KDSE STANDARD                              │   │
│  │     • Generate KDSE-standard artifacts                      │   │
│  │     • Include full traceability                             │   │
│  │     • Assign authority levels                               │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  4. GENERATE REPORT                                         │   │
│  │     • Collection summary                                    │   │
│  │     • Recommendations                                      │   │
│  │     • Remaining gaps                                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Knowledge Domains

KDSE supports the following knowledge domains:

| Domain | Description | Examples |
|--------|-------------|----------|
| physics | Physics principles and calculations | Thermal models, electrical equations |
| equipment | Equipment specifications and behavior | Transformer specs, battery characteristics |
| environment | Environmental conditions | Weather models, site conditions |
| standards | Standards and regulations | IEC, IEEE, OSHA requirements |
| business | Business rules and processes | Operating procedures, control logic |
| simulation | Simulation models | Digital twins, test scenarios |
| control | Control algorithms | PID tuning, protection logic |
| protocols | Communication protocols | Modbus, DNP3, IEC 61850 |
| vocabulary | Domain terminology | Glossary, acronyms |
| transformers | Transformer behavior | Magnetic flux, thermal limits |
| battery | Battery behavior | SOC/SOH, charge curves |
| relay | Relay protection | Coordination, settings |
| weather | Weather models | Solar irradiance, wind patterns |
| general | General engineering knowledge | Project overview, stakeholders |

---

## Knowledge Sources

### Supported Sources

| Source | Description | Authority Level |
|--------|-------------|----------------|
| Repository | Existing repository documentation | Project |
| Upload | Uploaded documents | Varies |
| Standards | Standards documents (IEC, IEEE, etc.) | Normative |
| Vendor | Vendor documentation | Vendor |
| Operator | Operator-supplied knowledge | Operator |
| AI-Assisted | AI-assisted knowledge acquisition | Derived |

### Extensibility

The implementation is extensible for future knowledge providers:

```go
type KnowledgeProvider interface {
    Name() string
    Collect(input *CollectionInput) ([]CollectedArtifact, error)
    CanCollect(domain KnowledgeDomain) bool
}
```

---

## Authority Levels

KDSE distinguishes knowledge by authority level:

| Level | Symbol | Description | Example |
|-------|--------|-------------|---------|
| Verified | ✓ | Tested and validated | Validated test results |
| Normative | ★ | KDSE standard or specification | KDSE Standard documents |
| Vendor | ◆ | Vendor documentation | Manufacturer specs |
| Project | ● | Project-specific knowledge | Design decisions |
| Operator | ○ | Operator-provided knowledge | Operating experience |
| Derived | ◐ | Derived from other sources | Calculated values |

The authority level is preserved throughout the knowledge base.

---

## Traceability

Every collected artifact includes full traceability:

| Field | Description |
|-------|-------------|
| Artifact ID | Unique identifier |
| Source | Original source document |
| Authority | Authority level |
| Version | Version information |
| Collection Date | When knowledge was collected |
| Collected By | Who collected the knowledge |
| Normative References | Related standards |
| Dependencies | Related artifacts |
| Traceability ID | Unique trace identifier |

---

## Normalization

Collected knowledge is NOT stored as raw notes. Every collected item is normalized into KDSE Documentation Standard artifacts.

### Output Structure

```
knowledge/
├── physics/
│   ├── DOMAIN_INDEX.md
│   └── thermal-model.md
├── equipment/
│   ├── DOMAIN_INDEX.md
│   ├── transformer-specs.md
│   └── battery-characteristics.md
├── standards/
│   ├── DOMAIN_INDEX.md
│   └── iec-61850-reference.md
└── INDEX.md  (master index)
```

### Artifact Format

Each artifact follows KDSE Documentation Standard:

```markdown
# Artifact Title

**Artifact ID:** equipment-ART-20260712
**Domain:** equipment
**Source:** vendor
**Authority Level:** Vendor
**Version:** 1.0
**Collection Date:** 2026-07-12T10:00:00Z
**Collected By:** Jane Developer
**Confidence Level:** 95%

---

## Summary

Brief description of the collected knowledge.

## Content

Detailed knowledge content...

## Traceability

| Field | Value |
|-------|-------|
| Traceability ID | TRC-20260712-001 |
| Source | Vendor Manual v2.1 |
| ... | ... |
```

---

## Reporting

### Collection Report

Generated at `.kdse/reports/collection-report-<session>.md`:

```markdown
# KDSE Knowledge Collection Report

**Report ID:** KDSE-RT-2026-07-12-001
**Repository:** /workspace/project/myapp
**Collection Date:** 2026-07-12

---

## Executive Summary

This report documents the knowledge collection performed...

## Collection Statistics

| Metric | Value |
|--------|-------|
| Knowledge Areas Reviewed | 5 |
| Knowledge Gaps Identified | 3 |
| Artifacts Collected | 12 |
| Artifacts Updated | 2 |
| Success Rate | 85% |

## Knowledge Still Missing

Detailed list of remaining gaps with recommendations...
```

### Terminal Output

Running `kdse collect` provides immediate summary:

```
📊 Knowledge Collection Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Knowledge Areas Reviewed:    5
  Knowledge Gaps Found:       3
  Artifacts Collected:        12
  Artifacts Updated:          2
  Artifacts Already Present:   8
  Knowledge Still Missing:     1
  Processing Time:           2.34s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 Recommendations:
  → Run 'kdse collect' again after addressing identified gaps
  → Run 'kdse normalize' to integrate collected knowledge

📁 Collected artifacts are available in: .kdse/knowledge/
📄 Full report: .kdse/collection-result.json
```

---

## Command Interface

### Basic Usage

```bash
kdse collect
```

### With Options

```bash
# Specify operator
kdse collect --operator "Jane Developer"

# Focus on specific domains
kdse collect --domain equipment --domain physics --domain standards

# Set priority
kdse collect --priority high --domain transformers

# Combine options
kdse collect --operator "Jane" --domain battery --domain relay --priority critical
```

### Integration Commands

```bash
# After audit - collect based on findings
kdse audit && kdse collect

# After normalization - integrate with normalized docs
kdse normalize && kdse collect

# Full workflow
kdse audit && kdse normalize && kdse collect && kdse run
```

---

## Configuration

### Provider Configuration

Knowledge providers can be registered in the collector:

```go
registry := NewProviderRegistry(repoPath)
registry.Register(NewCustomProvider())
```

### Domain Configuration

Default domain paths can be customized:

```go
DomainPaths = map[KnowledgeDomain]string{
    DomainPhysics:    "knowledge/physics",
    DomainEquipment:  "knowledge/equipment",
    // ...
}
```

---

## Success Criteria

The KDSE Runtime successfully implements `kdse collect` when:

1. ✅ Knowledge gaps are identified from audit findings and normalization results
2. ✅ Knowledge is collected from multiple extensible sources
3. ✅ Collected knowledge is normalized into KDSE-standard artifacts
4. ✅ Full traceability is maintained for every artifact
5. ✅ Authority levels are preserved and documented
6. ✅ Collection reports are generated with statistics and recommendations
7. ✅ The command integrates naturally with `kdse audit`, `kdse normalize`, and `kdse run`
8. ✅ The command remains completely optional and independent

---

## Examples

### Example 1: Initial Knowledge Collection

```bash
$ kdse collect
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Session ID:   KDSE-RT-2026-07-12-001
╚═══════════════════════════════════════════════════════════════╝

[1/5] Analyzing knowledge gaps...
      Found 5 knowledge gaps

[2/5] Collecting knowledge from available sources...
      ✓ Collected 15 artifacts

[3/5] Checking for existing knowledge...
      ✓ Verified 3 existing artifacts

[4/5] Generating KDSE-standard artifacts...
      ✓ Generated 12 new artifacts

[5/5] Analyzing remaining gaps...
      1 gap still requires attention

📊 Knowledge Collection Summary
  Knowledge Areas Reviewed:    5
  Artifacts Collected:        12
  Knowledge Still Missing:     1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📁 Collected artifacts are available in: .kdse/knowledge/
```

### Example 2: Domain-Focused Collection

```bash
$ kdse collect --domain battery --domain transformer --operator "Battery Team"
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Operator:     Battery Team
║ Areas:        battery, transformer
╚═══════════════════════════════════════════════════════════════╝

[1/5] Analyzing knowledge gaps...
      Found 2 knowledge gaps

[2/5] Collecting knowledge from available sources...
      ✓ Collected 8 artifacts

...
```

---

## Related Documents

| Document | Relationship |
|----------|-------------|
| [COMMANDS.md](COMMANDS.md) | Command interface definition |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Runtime architecture context |
| [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) | Session lifecycle |
| [NORMALIZE.md](NORMALIZE.md) | Documentation normalization |

---

## Document Relationships

```
COLLECT.md (this document)
    │
    ├── Defines: Knowledge collection workflow
    │
    ├── Referenced by:
    │   ├── COMMANDS.md (command reference)
    │   └── ARCHITECTURE.md (architecture context)
    │
    └── Related to:
        ├── NORMALIZE.md (normalization workflow)
        ├── SESSION_PROTOCOL.md (session context)
        └── WORKFLOW.md (overall workflow)
```

---

*This document is an informative reference implementation. It describes the Knowledge Collection Runtime, not KDSE requirements.*
