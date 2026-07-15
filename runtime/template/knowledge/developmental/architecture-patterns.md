# Architecture Patterns

**Type:** Developmental Knowledge  
**Template:** Baseline  
**Version:** 1.0

---

## Purpose

This document catalogs common architecture patterns applicable to software engineering. These patterns provide proven solutions to recurring architectural challenges.

---

## Patterns Catalog

### Layered Architecture

**Problem:** How to organize code into logical layers?

**Solution:** Separate concerns into horizontal layers.

```
┌─────────────────────────────────────────────────────────────────────┐
│                      LAYERED ARCHITECTURE                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Presentation Layer                                          │   │
│  │  - User interface                                           │   │
│  │  - API controllers                                         │   │
│  │  - View models                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Business Logic Layer                                         │   │
│  │  - Domain services                                          │   │
│  │  - Business rules                                          │   │
│  │  - Validation                                               │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Data Access Layer                                            │   │
│  │  - Repositories                                             │   │
│  │  - Data mappers                                             │   │
│  │  - Query builders                                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Infrastructure Layer                                         │   │
│  │  - Database                                                 │   │
│  │  - External services                                       │   │
│  │  - File system                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**When to use:**
- Standard business applications
- CRUD-heavy applications
- Teams new to architecture patterns

---

### Microservices Architecture

**Problem:** How to scale independently and deploy components separately?

**Solution:** Decompose into small, autonomous services.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    MICROSERVICES ARCHITECTURE                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Service  │  │ Service  │  │ Service  │  │ Service  │          │
│  │    A     │  │    B     │  │    C     │  │    D     │          │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘          │
│       │             │             │             │                 │
│       └─────────────┴──────┬──────┴─────────────┘                 │
│                            │                                        │
│                     ┌──────┴──────┐                                │
│                     │   Service   │                                │
│                     │   Mesh /    │                                │
│                     │   Gateway   │                                │
│                     └─────────────┘                                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**When to use:**
- High scalability requirements
- Independent deployment needs
- Multiple team ownership
- Technology diversity

---

### Event-Driven Architecture

**Problem:** How to handle asynchronous operations and loose coupling?

**Solution:** Use events to communicate between components.

```
┌─────────────────────────────────────────────────────────────────────┐
│                   EVENT-DRIVEN ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────┐         ┌──────────┐         ┌──────────┐          │
│  │ Producer │────────▶│  Event   │◀────────│ Consumer │          │
│  │          │  Event  │  Broker  │  Event  │          │          │
│  └──────────┘         └──────────┘         └──────────┘          │
│                                          │                          │
│                                          ▼                          │
│                                    ┌──────────┐                     │
│                                    │ Consumer │                     │
│                                    └──────────┘                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**When to use:**
- Real-time processing
- Loose coupling requirements
- High throughput needs
- Complex workflows

---

### Hexagonal Architecture

**Problem:** How to keep business logic isolated from external concerns?

**Solution:** Ports and adapters architecture.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    HEXAGONAL ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│                           ┌─────────┐                               │
│                           │ Domain  │                               │
│                           │  Core   │                               │
│                           │  Logic  │                               │
│                           └────┬────┘                               │
│                                │                                     │
│                    ┌───────────┴───────────┐                       │
│                    │                       │                        │
│              ┌─────┴─────┐           ┌─────┴─────┐                  │
│              │   Input   │           │  Output   │                  │
│              │   Ports   │           │   Ports   │                  │
│              └─────┬─────┘           └─────┬─────┘                  │
│                    │                       │                        │
│              ┌─────┴─────┐           ┌─────┴─────┐                  │
│              │  Input    │           │  Output   │                  │
│              │ Adapters  │           │ Adapters  │                  │
│              │ (REST,CLI)│           │ (DB,Queue)│                  │
│              └───────────┘           └───────────┘                  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**When to use:**
- Testability priority
- Technology flexibility
- Long-term maintainability
- Domain-driven design

---

## Pattern Selection Guide

| Scenario | Recommended Pattern |
|----------|-------------------|
| Small team, simple app | Layered |
| High scale, multiple teams | Microservices |
| Real-time processing | Event-Driven |
| Testability critical | Hexagonal |
| Domain complexity | Domain-Driven Design + Hexagonal |

---

## Pattern Combinations

Patterns can be combined:

| Primary | Secondary | Use Case |
|---------|-----------|----------|
| Microservices | Event-Driven | Distributed real-time systems |
| Layered | Hexagonal | Enterprise applications |
| Hexagonal | Event-Driven | Domain-driven event sourcing |

---

## Anti-Patterns to Avoid

| Anti-Pattern | Description | Solution |
|-------------|-------------|----------|
| Big Ball of Mud | No clear boundaries | Apply layering or microservices |
| Service-Oriented Angst | Over-engineered services | Right-size service boundaries |
| Gold Plating | Unnecessary complexity | Start simple, evolve as needed |
| Copy-Paste Reuse | Duplicated code across services | Extract shared libraries or services |

---

*This document is developmental knowledge. Project-specific patterns may be added in .kdse/knowledge/developmental/patterns/.*
