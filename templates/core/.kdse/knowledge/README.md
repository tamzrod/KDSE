# Knowledge Directory

This directory contains captured knowledge artifacts derived from evidence.

## Purpose

The knowledge directory stores validated knowledge entries that have been:
- Discovered from evidence
- Submitted as candidates
- Reviewed and accepted

## Structure

- `001-domain-overview.md` - High-level understanding of the problem domain
- `002-terminology.md` - Glossary of domain-specific terms
- `003-requirements.md` - Captured requirements and constraints

## Usage

Knowledge entries are managed through the KDSE notebook:

```bash
# Add new knowledge
kdse notebook add "New insight from analysis"

# Submit for review
kdse promote submit <entry-id>

# Review and accept
kdse promote review <entry-id> --accept
```

## Guidelines

- All knowledge must be derived from evidence
- Claims should cite source artifacts
- Include confidence level when known
- Update when new evidence contradicts existing knowledge
