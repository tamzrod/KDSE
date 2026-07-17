# System Context

## Purpose

This document defines the system boundaries and external interactions.

## Context Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         SYSTEM                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                                                         │    │
│  │                    [System Name]                        │    │
│  │                                                         │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
          │                        │                    │
          │                        │                    │
    ┌─────┴─────┐            ┌─────┴─────┐        ┌─────┴─────┐
    │  User A   │            │  User B   │        │ External  │
    │           │            │           │        │ Service   │
    └───────────┘            └───────────┘        └───────────┘
```

## External Entities

| Entity | Description | Interface |
|--------|-------------|-----------|
| [Entity] | [Description] | [API/File/UI] |

## Data Flows

### Input Flows

| Flow | Source | Format | Frequency |
|------|--------|--------|-----------|
| [Flow] | [Source] | [Format] | [Frequency] |

### Output Flows

| Flow | Destination | Format | Frequency |
|------|-------------|--------|-----------|
| [Flow] | [Destination] | [Format] | [Frequency] |

## External Dependencies

| Dependency | Version | Critical | Fallback |
|------------|---------|----------|----------|
| [Dep] | [Ver] | [Y/N] | [Fallback] |

## Trust Boundaries

[Describe security boundaries and trust relationships]

---

**Related Documents**:
- [Knowledge: Domain Overview](../knowledge/001-domain-overview.md)
- [Knowledge: Requirements](../knowledge/003-requirements.md)

**Evidence Sources**: [List source artifacts]
