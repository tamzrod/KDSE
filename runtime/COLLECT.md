# KDSE Artifact Collection Runtime

**Document Version:** 2.0  
**Type:** Informative Reference Implementation  
**Effective Date:** 2026-07-12

---

## Purpose

The KDSE Artifact Collection Runtime discovers and catalogs engineering evidence in the repository. It answers the question:

> **"What engineering evidence exists in this repository?"**

---

## Design Philosophy

### KDSE Runtime Role

The KDSE Runtime is an **Engineering Runtime**. It orchestrates engineering. It does not perform engineering.

| Responsibility | Belongs To |
|--------------|-----------|
| Repository awareness | Runtime |
| Session management | Runtime |
| Engineering context | Runtime |
| Artifact inventory | Runtime |
| Traceability | Runtime |
| Reporting | Runtime |

| Non-Responsibility | Belongs To |
|-------------------|-----------|
| AI frameworks | Executor |
| Knowledge engines | Executor |
| PDF parsing | Executor |
| RAG systems | Executor |
| Inference engines | Executor |
| Knowledge interpretation | Executor |

### What kdse collect IS

`kdse collect` discovers and catalogs engineering evidence. It:
- Scans the `artifacts/` directory
- Records file metadata
- Calculates integrity hashes
- Generates artifact inventory
- Produces collection reports

### What kdse collect IS NOT

`kdse collect` never:
- Interprets evidence content
- Extracts knowledge
- Assigns authority levels
- Identifies knowledge gaps
- Generates normalized knowledge

These responsibilities belong to executors.

---

## Evidence Pipeline

```
Engineering Evidence
        ↓
   kdse collect
        ↓
  Artifact Inventory
        ↓
     Executor
        ↓
Knowledge Extraction
        ↓
 Normalized Markdown
        ↓
 Operator Review
        ↓
Approved Knowledge
        ↓
  KDSE Runtime
```

The runtime consumes normalized knowledge only. It does not produce it.

---

## Directory Structure

### Raw Evidence

Raw evidence remains immutable in the `artifacts/` directory:

```
artifacts/
├── manuals/
│   ├── user-manual.pdf
│   └── maintenance-guide.pdf
├── standards/
│   ├── iec-61850.pdf
│   └── ieee-1549.pdf
├── specifications/
│   ├── transformer-spec.md
│   └── battery-requirements.md
├── datasheets/
│   ├── converter-datasheet.pdf
│   └── pcb-specs.md
├── drawings/
│   ├── system-architecture.png
│   └── circuit-schematic.svg
├── images/
├── videos/
└── archives/
```

### Runtime Output

```
.kdse/
├── artifacts/
│   └── inventory.json      # Artifact inventory
└── reports/
    └── artifact-collection-<session>.md  # Collection report
```

Raw artifacts remain unchanged. The runtime never modifies them.

---

## Artifact Categories

Artifacts are categorized by filename patterns and extensions:

| Category | Patterns | Extensions |
|----------|----------|------------|
| manual | manual, guide, handbook | - |
| standard | standard, iec, ieee, iso, nist | - |
| specification | spec, requirement | - |
| datasheet | datasheet, data-sheet | - |
| drawing | drawing, diagram, schematic | - |
| image | - | .jpg, .png, .gif, .svg, .bmp |
| video | - | .mp4, .avi, .mov, .mkv |
| archive | - | .zip, .tar, .gz, .rar, .7z |
| document | - | .pdf, .md, .txt, .doc, .rst, .adoc |
| unknown | - | Other extensions |

---

## Traceability

Each discovered artifact includes:

| Field | Description |
|-------|-------------|
| ID | Unique artifact identifier |
| Path | Absolute file path |
| RelativePath | Path relative to repository |
| Name | Filename |
| Category | Artifact category |
| Size | File size in bytes |
| Hash | SHA-256 hash for integrity |
| Modified | Last modification timestamp |
| Extension | File extension |
| CollectionID | Collection session identifier |

---

## Command Interface

### Basic Usage

```bash
kdse collect
```

### Output

```
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Artifact Collection                        ║
╠═══════════════════════════════════════════════════════════════╣
║ Repository:   /workspace/project/myapp
╚═══════════════════════════════════════════════════════════════╝

Discovering engineering artifacts in artifacts/ directory...

╔═══════════════════════════════════════════════════════════════╗
║              Collection Complete                             ║
╠═══════════════════════════════════════════════════════════════╣
║ Artifacts Discovered: 12
║ Total Size:          4.2 MB
║ Processing Time:     0.34s
╚═══════════════════════════════════════════════════════════════╝

Artifact inventory: .kdse/artifacts/inventory.json
Collection report:  .kdse/reports/

Artifacts by category:
  document: 5
  datasheet: 3
  drawing: 2
  manual: 1
  standard: 1

The runtime discovers and catalogs evidence.
Interpretation belongs to executors.
```

---

## Inventory Format

The artifact inventory is stored as JSON:

```json
{
  "session_id": "kdse-collect-20260712",
  "started_at": "2026-07-12T10:00:00Z",
  "completed_at": "2026-07-12T10:00:00Z",
  "repository": "/workspace/project/myapp",
  "artifacts_found": [
    {
      "id": "ART-a1b2-AA",
      "path": "/workspace/project/myapp/artifacts/standards/iec-61850.pdf",
      "relative_path": "artifacts/standards/iec-61850.pdf",
      "name": "iec-61850.pdf",
      "category": "standard",
      "size": 1048576,
      "hash": "abc123...",
      "modified": "2026-07-10T15:30:00Z",
      "extension": ".pdf",
      "collection_id": "kdse-collect-20260712"
    }
  ],
  "total_size": 4404019,
  "processing_time_seconds": 0.34
}
```

---

## Collection Report

A human-readable report is also generated:

```markdown
# KDSE Artifact Collection Report

| Field | Value |
|-------|-------|
| Session ID | kdse-collect-20260712 |
| Repository | /workspace/project/myapp |
| Collection Date | 2026-07-12T10:00:00Z |
| Report Version | 1.0 |

## Summary

**Artifacts Discovered:** 12
**Total Size:** 4.2 MB
**Processing Time:** 0.34 seconds

## Artifacts by Category

| Category | Count | Size |
|----------|-------|------|
| document | 5 | 1.2 MB |
| datasheet | 3 | 890 KB |
| drawing | 2 | 2.1 MB |
| manual | 1 | 120 KB |
| standard | 1 | 45 KB |

## Artifact Inventory

| ID | Name | Category | Size | Hash |
|----|------|----------|------|------|
| ART-a1b2-AA | iec-61850.pdf | standard | 1.0 MB | abc123... |
| ART-c3d4-BB | transformer-spec.md | specification | 12 KB | def456... |

---

*Report generated by KDSE Collect v1.0*
```

---

## Success Criteria

After refactoring, `kdse collect` successfully answers:

> ✅ **"What engineering evidence exists in this repository?"**

It never answers:

> ❌ **"What does this evidence mean?"**

---

## Simplification Summary

The refactored implementation:

| Before | After |
|--------|-------|
| ~2086 lines | ~250 lines |
| Knowledge providers | Removed |
| AI-assisted collection | Removed |
| Domain inference | Removed |
| Authority assignment | Removed |
| Gap analysis | Removed |
| Content interpretation | Removed |

The runtime remains lightweight, deterministic, repository-aware, and consistent with the KDSE philosophy.

---

## Related Documents

| Document | Relationship |
|----------|-------------|
| [COMMANDS.md](COMMANDS.md) | Command interface definition |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Runtime architecture context |
| [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) | Session lifecycle |

---

*This document is an informative reference implementation. It describes the Artifact Collection Runtime, not KDSE requirements.*
