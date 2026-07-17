# Components

## Purpose

This document defines the system component structure and responsibilities.

## Component Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         [System]                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │ Component A │  │ Component B │  │ Component C │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

## Component Details

### Component A

**Responsibility**: [What this component does]

**Public API**:
```yaml
endpoints:
  - GET /resource
  - POST /resource
```

**Dependencies**:
- [Component B]
- [External Service X]

**Data**: [What data it manages]

---

### Component B

**Responsibility**: [What this component does]

**Public API**:
```yaml
endpoints:
  - GET /items
```

**Dependencies**:
- [Database]

**Data**: [What data it manages]

---

## Component Interactions

| From | To | Protocol | Description |
|------|----|----------|-------------|
| [Comp] | [Comp] | [Protocol] | [Description] |

## Module Structure

```
src/
├── components/
│   ├── component-a/
│   │   ├── index.ts
│   │   ├── service.ts
│   │   └── types.ts
│   └── component-b/
│       └── ...
```

---

**Related Documents**:
- [System Context](001-system-context.md)
- [Data Model](003-data-model.md)

**Evidence Sources**: [List source artifacts]
