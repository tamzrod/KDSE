# KDSE Knowledge Collection Examples

**Document Version:** 1.0  
**Type:** Example  
**Effective Date:** 2026-07-12

---

## Basic Usage

### Minimal Collection

```bash
$ kdse collect
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Session ID:   KDSE-RT-2026-07-12-001
╚═══════════════════════════════════════════════════════════════╝

[1/5] Analyzing knowledge gaps...
      Found 3 knowledge gaps

[2/5] Collecting knowledge from available sources...
      Checking provider: Repository Documentation...
      ✓ Collected 5 artifacts
      Checking provider: AI-Assisted Knowledge...
      ✓ Collected 2 artifacts

[3/5] Checking for existing knowledge...
      ✓ Verified 3 existing artifacts

[4/5] Generating KDSE-standard artifacts...
      ✓ Generated 7 new artifacts

[5/5] Analyzing remaining gaps...
      0 gaps still require attention

📊 Knowledge Collection Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Knowledge Areas Reviewed:    3
  Knowledge Gaps Found:       3
  Artifacts Collected:        7
  Artifacts Updated:          0
  Artifacts Already Present:   3
  Knowledge Still Missing:     0
  Processing Time:           1.45s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📁 Collected artifacts are available in: .kdse/knowledge/
```

---

## With Operator Name

```bash
$ kdse collect --operator "Jane Developer"
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Operator:     Jane Developer
╚═══════════════════════════════════════════════════════════════╝
...
```

---

## Domain-Focused Collection

### Single Domain

```bash
$ kdse collect --domain battery
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Areas:        battery
╚═══════════════════════════════════════════════════════════════╝

[1/5] Analyzing knowledge gaps...
      Found 1 knowledge gaps
      • [battery-GAP-20260712-001] Battery Behavior (High)

[2/5] Collecting knowledge from available sources...
      Checking provider: Repository Documentation...
      ✓ Collected 3 battery-related artifacts
...
```

### Multiple Domains

```bash
$ kdse collect --domain equipment --domain physics --domain standards
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Areas:        equipment, physics, standards
╚═══════════════════════════════════════════════════════════════╝
...
```

---

## Priority-Based Collection

### High Priority

```bash
$ kdse collect --priority high
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Priority:     high
╚═══════════════════════════════════════════════════════════════╝
...
```

### Critical Priority with Domain

```bash
$ kdse collect --priority critical --domain transformers
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Areas:        transformers
║ Priority:     critical
╚═══════════════════════════════════════════════════════════════╝
...
```

---

## Combined Options

```bash
$ kdse collect --operator "Power Systems Team" --domain relay --domain transformer --domain battery --priority high
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/solar-inverter
║ Operator:     Power Systems Team
║ Areas:        relay, transformer, battery
║ Priority:     high
╚═══════════════════════════════════════════════════════════════╝
...
```

---

## Integration with Other Commands

### Full Workflow

```bash
# Run audit to identify gaps
kdse audit

# Run normalization to standardize existing docs
kdse normalize

# Collect missing knowledge based on findings
kdse collect

# Start session with comprehensive knowledge base
kdse run
```

### Collect After Audit

```bash
$ kdse audit && kdse collect
[Audit output...]
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Knowledge Collection                       ║
...
```

---

## Knowledge Domains Reference

| Domain | Flag | Example |
|--------|------|---------|
| physics | `--domain physics` | Thermal models, electrical equations |
| equipment | `--domain equipment` | Equipment specs, behavior |
| environment | `--domain environment` | Site conditions, weather |
| standards | `--domain standards` | IEC, IEEE requirements |
| business | `--domain business` | Operating procedures |
| simulation | `--domain simulation` | Digital twins, models |
| control | `--domain control` | PID, protection logic |
| protocols | `--domain protocols` | Modbus, IEC 61850 |
| vocabulary | `--domain vocabulary` | Glossary, acronyms |
| transformers | `--domain transformers` | Magnetic flux, thermal limits |
| battery | `--domain battery` | SOC/SOH, charge curves |
| relay | `--domain relay` | Protection coordination |
| weather | `--domain weather` | Solar irradiance |
| general | `--domain general` | Project overview |

---

## Output Files

After running `kdse collect`, the following files are generated:

### Directory Structure

```
.kdse/
├── knowledge/
│   ├── INDEX.md                    # Master knowledge index
│   ├── physics/
│   │   ├── DOMAIN_INDEX.md
│   │   └── thermal-model.md
│   ├── equipment/
│   │   ├── DOMAIN_INDEX.md
│   │   ├── transformer-specs.md
│   │   └── battery-characteristics.md
│   ├── standards/
│   │   ├── DOMAIN_INDEX.md
│   │   └── iec-61850-reference.md
│   └── ...
├── reports/
│   └── collection-report-<session>.md
└── collection-result.json
```

### Artifact Format

```markdown
# Transformer Specifications

**Artifact ID:** transformers-ART-20260712
**Domain:** transformers
**Source:** repository
**Authority Level:** Project
**Version:** 1.0
**Collection Date:** 2026-07-12T10:00:00Z
**Collected By:** Jane Developer
**Confidence Level:** 85%

---

## Summary

This document contains transformer specifications...

## Content

Detailed specifications...

## Traceability

| Field | Value |
|-------|-------|
| Traceability ID | TRC-20260712-001 |
| Source | docs/transformer-specs.md |
| ... | ... |
```

---

## Common Use Cases

### Initial Project Setup

```bash
# Fresh repository - establish baseline knowledge
kdse collect --domain general --operator "Project Lead"
```

### Power Systems Project

```bash
# Collect power engineering knowledge
kdse collect \
  --operator "Power Engineer" \
  --domain transformers \
  --domain battery \
  --domain relay \
  --domain standards \
  --priority high
```

### Control Systems Project

```bash
# Collect control systems knowledge
kdse collect \
  --operator "Controls Team" \
  --domain control \
  --domain protocols \
  --domain simulation
```

---

*This document provides examples of using the kdse collect command.*
