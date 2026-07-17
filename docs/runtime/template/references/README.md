# References

**Type:** Template Structure  
**Version:** 1.0

---

## Purpose

The References directory contains domain-specific reference knowledge pulled from the Global Runtime or added during project work.

---

## Structure

```
references/
├── README.md                    # This file
├── domains/                     # Domain knowledge
│   └── README.md
├── external/                    # External references
│   └── README.md
└── index.yaml                   # Reference index
```

---

## Domain References

Domain-specific knowledge is stored here:

```
references/domains/
├── {domain-name}/
│   ├── README.md
│   ├── concepts/
│   ├── patterns/
│   └── glossary/
```

---

## Reference Sources

| Source | Description |
|--------|-------------|
| Global Runtime | Pulled via `kdse pull` |
| External | Standards, documentation |
| Project | Created during engineering |

---

## Discovery

References are discovered during Knowledge Discovery:

```bash
# Pull specific domain
kdse pull {domain-name}

# Pull recommended references
kdse pull --recommended
```

---

*This directory is copied from the KDSE template during initialization.*
