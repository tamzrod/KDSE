# Domain Interfaces

## Purpose

This document establishes **Domain Interface** as a first-class concept in KDSE. Domain Interfaces represent domain responsibilities rather than implementation technologies.

## What Domain Interfaces Are

A Domain Interface is an implementation-independent contract describing the information exchanged between domain concepts.

Domain Interfaces define:

- **What information exists**: The data items exchanged
- **What the information means**: The semantic meaning of each data item
- **Responsibilities**: What each party provides and expects
- **Constraints**: Any limitations or requirements
- **Assumptions**: What is assumed to be true
- **Units**: What units apply (when applicable)
- **Relationships**: How this interface relates to other domain concepts

Domain Interfaces intentionally exclude:

- Programming language
- Software framework
- Communication protocol
- Database technology
- Vendor-specific implementations
- Deployment architecture

## Domain Independence

A Domain Interface shall remain valid regardless of:

- Programming language
- Software framework
- Communication protocol
- Database
- Vendor
- Deployment architecture

Changing implementation technology shall never require changing a Domain Interface.

## Domain Interface vs Implementation Technology

### The Distinction

| Aspect | Domain Interface | Implementation Technology |
|--------|-----------------|-------------------------|
| Focus | What is exchanged | How it is exchanged |
| Technology | Domain-agnostic | Technology-specific |
| Example | Measurement data | Modbus register 40001 |
| Example | Status information | MQTT topic /status |
| Example | Historical records | REST API endpoint |

### Why the Distinction Matters

The distinction between Domain Interfaces and implementation technologies provides:

1. **Technology Independence**: Domain Interfaces remain valid across technology changes
2. **Clear Responsibilities**: Domain responsibilities are clearly defined
3. **Flexibility**: Implementation can change without affecting Domain Interface definitions
4. **Traceability**: Domain decisions trace to Domain Interfaces

## Domain Interface Definition Template

When defining a Domain Interface, document:

### 1. Interface Purpose

Describe what the interface provides or accepts.

### 2. Information Model

Define the data exchanged:

| Data Item | Meaning | Type | Units | Range | Update Rate |

### 3. Responsibilities

Define responsibilities of each party:

| Party | Responsibility |
|-------|----------------|
| Provider | What it must provide |
| Consumer | What it must consume |

### 4. Quality of Service

Define non-functional requirements:

- Latency requirements
- Availability requirements
- Accuracy requirements

### 5. Constraints

Define any constraints:

- Environmental constraints
- Safety constraints
- Operational constraints

## Domain Interface Examples

The following examples illustrate the concept of Domain Interfaces. These are generic illustrations—the implementation or project may define Domain Interfaces appropriate to its own domain.

### Example 1: Observation Interface

> **Observation Interface**
>
> Provides observation data from the system.
>
> | Information | Meaning | Units | Update Rate |
> |-------------|---------|-------|------------|
> | Value | The observed value | (varies) | Per observation |
> | Timestamp | When observed | Time | - |
> | Quality | Data quality indicator | (varies) | - |
>
> **Responsibilities**:
>
> - Provider publishes observations as they occur
> - Consumer receives and processes observations
> - Quality indicator reflects data validity

### Example 2: Command Interface

> **Command Interface**
>
> Accepts commands from authorized sources.
>
> | Command | Meaning | Parameters |
> |---------|---------|------------|
> | Execute | Execute the specified action | Action identifier, parameters |
> | Cancel | Cancel a pending action | Action identifier |
> | Query | Query command status | Action identifier |
>
> **Responsibilities**:
>
> - Interface receives validated commands
> - Commands are executed in priority order
> - Status is reported upon completion

### Example 3: Status Interface

> **Status Interface**
>
> Exposes current system state.
>
> | Status | Meaning | Values |
> |--------|---------|--------|
> | State | Current operational state | (varies by domain) |
> | Mode | Current operating mode | (varies by domain) |
> | Health | System health indicator | (varies by domain) |
>
> **Responsibilities**:
>
> - Status reflects current system state
> - Status updates are event-driven
> - Status is available to all authorized consumers

### Example 4: Record Interface

> **Record Interface**
>
> Provides access to historical records.
>
> | Information | Meaning | Resolution | Retention |
> |-------------|---------|------------|-----------|
> | Records | Historical records | Per record | Per policy |
> | Events | Event history | Per event | Per policy |
> | Audit | Audit trail | Per action | Per policy |
>
> **Responsibilities**:
>
> - Records are accurately timestamped
> - Data integrity is maintained
> - Queries return data within specified latency

## Domain Interfaces in the Artifact Hierarchy

### Relationship to Domain Knowledge

Domain Interfaces are derived from Domain Knowledge:

```
Domain Knowledge
        │
        │ "What domain behavior requires"
        ↓
Domain Interface Definition
        │
        │ "What information must be exchanged"
        ↓
Architecture (Service boundaries)
        │
        │ "How is information exchanged"
        ↓
Implementation (Specific technology)
```

### Domain Interfaces and Architecture

Architecture uses Domain Interfaces to define boundaries:

```
Component A
        │
        │ Domain Interface
        ↓
Component B
```

The Domain Interface defines what crosses the boundary. Architecture determines which components expose which interfaces.

### Domain Interfaces and Implementation

Implementation realizes Domain Interfaces using specific technologies:

| Domain Interface | Implementation Technology |
|-----------------|-------------------------|
| Observation data | REST endpoint, MQTT topic, database query |
| Command data | gRPC, message queue, event bus |
| Status data | WebSocket, polling endpoint, pub/sub |
| Historical data | SQL database, time-series DB, file system |

## Common Errors

### Error 1: Embedding Implementation Details

**Incorrect**:

> "The sensor provides data via HTTP at 9600 baud."

**Why Incorrect**: Communication details are implementation.

**Correct**:

> "The sensor provides measurement values with timestamps."

### Error 2: Specifying Protocol Details

**Incorrect**:

> "Commands are accepted as JSON payloads over HTTP POST."

**Why Incorrect**: Protocol details are implementation.

**Correct**:

> "Commands include action identifiers and parameters."

### Error 3: Including Technology Choices

**Incorrect**:

> "The record interface uses MongoDB with specific collections."

**Why Incorrect**: Technology choice is implementation.

**Correct**:

> "The record interface provides time-series record access."

## Domain Interface Validation Checklist

When creating or reviewing Domain Interfaces, verify:

- [ ] Is the interface technology-agnostic?
- [ ] Is the information meaning clearly defined?
- [ ] Are responsibilities clearly assigned?
- [ ] Is the interface traceable to Domain Knowledge?
- [ ] Does implementation respect the interface without dictating it?

## Domain Generalization

KDSE is a methodology for deriving Domain Knowledge. Domain Interfaces are part of this methodology.

The concept of a Domain Interface applies to any domain:

- **Engineering domains**: Sensor interfaces, control interfaces, status interfaces
- **Business domains**: Customer interfaces, order interfaces, notification interfaces
- **Healthcare domains**: Patient interfaces, clinical interfaces, billing interfaces
- **Financial domains**: Transaction interfaces, account interfaces, reporting interfaces

The methodology remains the same; only the domain specifics change.

## Summary

Domain Interfaces represent domain responsibilities rather than implementation technologies:

- **Define** what information is exchanged, not how
- **Specify** meaning, units, and responsibilities, not protocol details
- **Remain valid** across technology changes
- **Support** separation between Domain Knowledge and Implementation

Understanding Domain Interfaces is essential for maintaining the separation of concerns that makes KDSE effective.

---

## Version

- **Document Version**: 2.0
- **Effective Date**: 2026-07-13
- **Change Note**: Generalized from Engineering Interface to Domain Interface to make KDSE domain-agnostic
