# Architecture Phase

## Purpose

This document establishes **Architecture** as a separate methodology phase in KDSE. Architecture explains how Engineering Knowledge is organized into software. Architecture shall never be confused with Engineering Knowledge.

## What Architecture Is

Architecture is the organizational structure that explains how Engineering Knowledge is realized in software.

Architecture describes:

- **Modules**: How software is divided into units
- **Services**: How functionality is separated into services
- **Deployment**: How components are deployed
- **Topology**: How components are arranged
- **Redundancy**: How reliability is achieved
- **Communication Architecture**: How components communicate
- **Software Boundaries**: Where boundaries exist between components
- **Runtime Organization**: How the system operates at runtime

## What Architecture Is Not

Architecture is not:

### Engineering Knowledge

Architecture is derived from Engineering Knowledge but is not Engineering Knowledge itself.

| Aspect | Engineering Knowledge | Architecture |
|--------|----------------------|--------------|
| Focus | Engineering intent | Software organization |
| Independence | Implementation-independent | Technology-aware |
| Lifetime | Longest-lived | Shorter than Knowledge |
| Example | "System supports grid-forming mode" | "Grid-forming is a separate service" |

### Implementation

Architecture precedes implementation and guides it.

| Aspect | Architecture | Implementation |
|--------|-------------|---------------|
| Focus | Organization | Realization |
| Content | Structure | Code |
| Technology | Technology-agnostic | Technology-specific |

### Reference Artifacts

Architecture is derived from Reference Artifacts but is not the artifacts themselves.

| Aspect | Reference Artifacts | Architecture |
|--------|--------------------|--------------|
| Role | Evidence | Derived structure |
| Authority | Support Engineering Knowledge | Authorize Implementation |
| Creation | Already existed | Created during Architecture phase |

## The Role of Architecture in KDSE

Architecture occupies the space between Engineering Knowledge and Implementation:

```
Engineering Knowledge
        ↓
Architecture (Organization)
        ↓
Implementation (Realization)
```

### Architecture Derives from Engineering Knowledge

Architecture decisions must trace to Engineering Knowledge:

1. Identify the Engineering Knowledge that needs architectural expression
2. Determine organizational approaches
3. Select the approach that best serves the Engineering Knowledge
4. Document the decision and its rationale

### Architecture Authorizes Implementation

Implementation must conform to Architecture:

1. Implementation traces to Architecture
2. Implementation realizes Architecture
3. Implementation cannot contradict Architecture

## Architecture Categories

### Category 1: Structural Architecture

Describes how software is organized into modules, services, and components.

**Examples**:

- "The system is organized into control, monitoring, and historian services"
- "The control service contains grid-following and grid-forming modules"
- "Each service exposes a well-defined interface"

### Category 2: Communication Architecture

Describes how components communicate.

**Examples**:

- "Services communicate through message passing"
- "The control service subscribes to sensor data"
- "The historian receives data via a publish-subscribe channel"

### Category 3: Deployment Architecture

Describes how components are deployed.

**Examples**:

- "The control service runs on dedicated hardware"
- "Edge devices run the sensor interface module"
- "The historian runs in the cloud data center"

### Category 4: Runtime Architecture

Describes how the system operates at runtime.

**Examples**:

- "The control service operates continuously"
- "Sensor data flows at 100ms intervals"
- "Commands are processed synchronously"

### Category 5: Redundancy Architecture

Describes how reliability is achieved.

**Examples**:

- "The control service runs with 1+1 redundancy"
- "Sensor data is replicated across multiple subscribers"
- "Failover is automatic within 1 second"

## Architecture Decision Records

Architecture decisions are documented in Architecture Decision Records (ADRs).

### ADR Contents

An ADR documents:

1. **Decision**: What was decided
2. **Context**: The situation that prompted the decision
3. **Decision**: The decision that was made
4. **Rationale**: Why this decision was made
5. **Consequences**: What this decision entails
6. **Alternatives**: What alternatives were considered

### ADR Example

```
# ADR-001: Service-Based Architecture for Control System

## Context

The control system must support multiple control modes (grid-following, 
grid-forming) with independent evolution and deployment.

## Decision

We will use a service-based architecture with separate services for:
- Grid-following control
- Grid-forming control
- Monitoring
- Historian

## Rationale

- Independent evolution of control modes
- Independent deployment scaling
- Clear separation of concerns
- Matches Engineering Knowledge structure

## Consequences

- Inter-service communication required
- Latency considerations
- Service discovery needed
- Distributed failure modes

## Alternatives Considered

- Monolithic architecture (rejected - less flexible)
- Microservices per control loop (rejected - too granular)
```

## Architecture and the Engineering Independence Test

Architecture statements shall be assessed against the Engineering Independence Test:

> "If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"

| Result | Classification | Action |
|--------|---------------|--------|
| YES | Engineering Knowledge | Move to Knowledge Artifacts |
| NO | Architecture | Retain in Architecture |

### Example Assessment

**Statement**: "The system shall use a message broker for inter-service communication."

**Assessment**:
- Would this be true if we changed message broker technology? NO
- Would this be true if we changed programming language? YES
- Conclusion: This is Architecture (communication organization), not Implementation (specific broker)

**Statement**: "The system shall use Apache Kafka for message brokering."

**Assessment**:
- Would this be true if we changed to RabbitMQ? NO
- Conclusion: This is Implementation, not Architecture

## Architecture Activities

### Activity 1: Identify Architectural Drivers

Identify the Domain Knowledge that drives architectural decisions:

1. Review approved Domain Knowledge
2. Identify knowledge that requires architectural expression
3. Prioritize by significance and complexity

### Activity 2: Explore Architectural Approaches

Explore possible organizational approaches:

1. Generate alternative architectures
2. Evaluate against Domain Knowledge
3. Consider non-functional requirements
4. Assess trade-offs

### Activity 3: Make Architecture Decisions

Make decisions that establish structure:

1. Select architectural approach
2. Document rationale
3. Create ADRs for significant decisions
4. Ensure traceability to Domain Knowledge

### Activity 4: Document Architecture

Document the resulting architecture:

1. Create architecture diagrams
2. Define interfaces (Domain Interfaces, not implementation)
3. Document component responsibilities
4. Establish communication patterns

### Activity 5: Validate Against Domain Knowledge

Verify that architecture satisfies Domain Knowledge:

1. Trace each architectural element to Domain Knowledge
2. Confirm all Domain Knowledge is addressed
3. Identify any gaps or conflicts

## Architecture and Domain Interfaces

Domain Interfaces represent domain responsibilities at the architecture level.

### What Domain Interfaces Are

Domain Interfaces define:

- What information is exchanged
- What the information means
- Responsibilities of each party
- Units (when applicable)

### What Domain Interfaces Are Not

Domain Interfaces intentionally exclude:

- Communication technologies (REST, gRPC, MQTT)
- Protocol details
- Technology-specific implementations

### Example

**Domain Interface**: Observation Interface

> "The sensor provides observation data with timestamps and quality indicators."

**Architecture Expression**: Communication Architecture

> "The sensor component exposes a Domain Interface for observation data."

**Implementation**: Specific technology choices

> "Observation data arrives via MQTT topic /sensors/observations."

## Architecture Boundaries

Architecture establishes boundaries between components.

### Boundary Types

1. **Service Boundaries**: Where services are separated
2. **Module Boundaries**: Where modules are separated
3. **Domain Boundaries**: Where domains are separated
4. **Deployment Boundaries**: Where deployments are separated

### Boundary Principles

- Boundaries align with Engineering Knowledge boundaries
- Boundaries minimize cross-cutting concerns
- Boundaries enable independent evolution
- Boundaries respect Engineering Knowledge dependencies

## Common Errors

### Error 1: Including Implementation Details

**Incorrect Architecture**:

> "The historian shall use InfluxDB with 1-second retention."

**Correct Architecture**:

> "The historian shall store time-series data with configurable retention."

### Error 2: Confusing Architecture with Engineering Knowledge

**Incorrect**:

> "The inverter shall support grid-forming mode" ← This is Engineering Knowledge

**Correct**:

> "Grid-forming control shall be implemented as a separate service" ← This is Architecture

### Error 3: Over-Constraining Implementation

**Incorrect Architecture**:

> "The control service shall use a specific PID tuning algorithm."

**Correct Architecture**:

> "The control service shall implement closed-loop control with configurable parameters."

## Architecture Validation Checklist

When creating or reviewing Architecture, verify:

- [ ] Does each architectural decision trace to Domain Knowledge?
- [ ] Are there any gaps where Domain Knowledge is not addressed?
- [ ] Are there any conflicts between architectural decisions?
- [ ] Are ADRs created for significant decisions?
- [ ] Is the architecture technology-agnostic where appropriate?
- [ ] Are Domain Interfaces defined without implementation specifics?
- [ ] Does implementation can conform to this architecture?

## Architecture and Authority

Architecture derives authority from Domain Knowledge:

```
Domain Knowledge (Authority Source)
        │
        │ Derivation
        ↓
Architecture (Authority Bearer)
        │
        │ Authorization
        ↓
Implementation (Authority Consumer)
```

Architecture may not contradict Domain Knowledge. If architecture cannot satisfy Domain Knowledge, the situation must be resolved through clarification or redesign.

## Summary

Architecture is a separate methodology phase that:

- **Explains** how Domain Knowledge is organized into software
- **Follows** the Knowledge Derivation Lifecycle
- **Produces** architectural decisions and documentation
- **Authorizes** implementation through traceable decisions
- **Is distinct from** Domain Knowledge (intent) and Implementation (realization)
- **Uses** Domain Interfaces to define boundaries

Understanding Architecture as a distinct phase is essential for maintaining separation of concerns.

---

## Version

- **Document Version**: 2.0
- **Effective Date**: 2026-07-13
- **Change Note**: Updated to Domain Interface terminology and generalized examples
