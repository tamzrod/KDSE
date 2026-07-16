# Implementation Phase

## Purpose

This document establishes **Implementation** as a separate methodology phase in KDSE. Implementation explains how Architecture is realized using specific technologies. Implementation details shall never appear inside Engineering Knowledge Artifacts.

## What Implementation Is

Implementation is the realization of Architecture using specific technologies.

Implementation includes:

- **Source Code**: Program code in specific languages
- **Node-RED Flows**: Visual programming flows
- **Go Programs**: Go language implementations
- **PLC Logic**: Programmable logic controller code
- **JavaScript**: JavaScript/TypeScript implementations
- **Modbus Registers**: Specific register mappings
- **MQTT Topics**: Specific topic hierarchies
- **REST Endpoints**: Specific API definitions
- **OPC UA Nodes**: Specific node configurations
- **IEC 61850 Models**: Specific data models
- **DNP3 Objects**: Specific object definitions
- **Configuration Files**: Specific configuration formats

## What Implementation Is Not

Implementation is not:

### Engineering Knowledge

Implementation is derived from Architecture, which derives from Engineering Knowledge.

| Aspect | Engineering Knowledge | Implementation |
|--------|----------------------|----------------|
| Focus | Engineering intent | Technology realization |
| Independence | Implementation-independent | Technology-specific |
| Lifetime | Longest-lived | Shorter than Architecture |

### Architecture

Implementation realizes Architecture.

| Aspect | Architecture | Implementation |
|--------|--------------|----------------|
| Focus | Organization | Realization |
| Content | Structure | Code/Configuration |
| Technology | Technology-agnostic | Technology-specific |

### Reference Artifacts

Implementation may become Reference Artifacts for future knowledge derivation.

| Aspect | Reference Artifacts | Implementation |
|--------|--------------------|----------------|
| Role | Evidence for derivation | Result of derivation |
| Creation | Pre-exist derivation | Created during implementation |

## The Role of Implementation in KDSE

Implementation occupies the final step in the artifact hierarchy:

```
Engineering Knowledge
        ↓
Architecture (Organization)
        ↓
Implementation (Realization)
```

### Architecture Guides Implementation

Implementation must conform to Architecture:

1. Identify the Architecture that authorizes this implementation
2. Realize the architectural decisions in code
3. Ensure implementation does not contradict Architecture
4. Maintain traceability from implementation to Architecture

### Implementation Traces to Architecture

Every implementation artifact must trace to Architecture:

```
Implementation Artifact
        │
        │ Realizes
        ↓
Architecture Decision
        │
        │ Derives from
        ↓
Engineering Knowledge
```

## Implementation Technology Examples

### Programming Languages

| Technology | Description | Is Implementation |
|------------|-------------|-------------------|
| Go | Go language programs | Yes |
| JavaScript | JavaScript/TypeScript code | Yes |
| Python | Python programs | Yes |
| Java | Java programs | Yes |

### Visual Programming

| Technology | Description | Is Implementation |
|------------|-------------|-------------------|
| Node-RED | Visual flow programming | Yes |
| LabVIEW | Graphical programming | Yes |
| FBD | Function block diagrams | Yes |

### Industrial Control

| Technology | Description | Is Implementation |
|------------|-------------|-------------------|
| IEC 61131-3 | PLC programming standards | Yes |
| Structured Text | PLC text language | Yes |
| Ladder Logic | PLC graphical language | Yes |

### Communication Protocols

| Technology | Description | Is Implementation |
|------------|-------------|-------------------|
| Modbus RTU | Serial communication | Yes |
| Modbus TCP | Ethernet communication | Yes |
| DNP3 | SCADA protocol | Yes |
| IEC 61850 | Substation protocol | Yes |
| OPC UA | Industrial interoperability | Yes |

### Data Registrations

| Technology | Description | Is Implementation |
|------------|-------------|-------------------|
| Modbus Register 40001 | Specific register | Yes |
| MQTT Topic /power/voltage | Specific topic | Yes |
| OPC UA Node ID | Specific node | Yes |

### Runtime Environments

| Technology | Description | Is Implementation |
|------------|-------------|-------------------|
| Docker | Container runtime | Yes |
| Kubernetes | Orchestration | Yes |
| Node.js | JavaScript runtime | Yes |

## Implementation and the Engineering Independence Test

Implementation statements fail the Engineering Independence Test:

> "If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"

If NO: The statement is Implementation.

### Examples

**Implementation** (fails the test):

- "The historian uses InfluxDB for storage" (InfluxDB is specific)
- "Sensor data arrives via Modbus TCP register 40001" (Modbus is specific)
- "The API is implemented as REST endpoints" (REST is specific)

**Engineering Knowledge** (passes the test):

- "Time-series data is stored with millisecond precision" (precision is engineering)
- "Sensor data provides voltage measurements" (measurement is engineering)
- "Status is accessible via an API" (API concept is architecture, not specific technology)

## Implementation Activities

### Activity 1: Identify Implementation Scope

Identify what needs to be implemented:

1. Review Architecture
2. Identify components requiring implementation
3. Prioritize implementation sequence
4. Identify dependencies

### Activity 2: Select Technologies

Select specific technologies:

1. Review Architecture for technology guidance
2. Evaluate technology options
3. Document technology selections
4. Create Implementation Decisions (if significant)

### Activity 3: Implement Components

Implement each component:

1. Write code according to Architecture
2. Configure systems according to Architecture
3. Integrate components according to Architecture
4. Maintain traceability to Architecture

### Activity 4: Verify Conformance

Verify implementation conforms to Architecture:

1. Review implementation against Architecture
2. Ensure no contradictions
3. Document any deviations requiring Architecture change
4. Update traceability records

## Implementation Documentation

Implementation documentation includes:

### Source Code

- Code files in specific languages
- Comments explaining implementation
- Inline documentation

### Configuration

- Configuration files
- Environment variables
- Deployment specifications

### Build Artifacts

- Compiled binaries
- Container images
- Deployment packages

### Implementation-Specific Documentation

- API documentation
- Protocol documentation
- Integration guides

## Implementation and Domain Interfaces

Domain Interfaces define what information is exchanged at the domain level. Implementation defines how that exchange is realized technologically.

### Example: Observation Interface

**Domain Interface** (Domain Knowledge/Architecture):

> "The sensor provides observation data with timestamps and quality indicators."

**Implementation** (Implementation):

> "Observation data is received via MQTT topic /sensors/observations with JSON payload."

### Implementation Examples by Technology

#### REST Implementation

```
Domain Interface: Observation data
Implementation: GET /api/v1/observations, Response JSON
```

#### Message Queue Implementation

```
Domain Interface: Status data
Implementation: Topic /domain/status, Payload JSON
```

#### Database Implementation

```
Domain Interface: Historical records
Implementation: SQL table with timestamp and value columns
```

#### Event Bus Implementation

```
Domain Interface: Event notifications
Implementation: Event bus with event type and payload
```

## Implementation Validation Checklist

When creating or reviewing Implementation, verify:

- [ ] Does this implementation trace to Architecture?
- [ ] Does this implementation conform to Architecture?
- [ ] Are technology choices documented?
- [ ] Is implementation-specific documentation complete?
- [ ] Are Domain Interfaces respected (not violated)?
- [ ] Is traceability maintained from implementation to Architecture?

## Common Errors

### Error 1: Including Domain Knowledge in Implementation

**Incorrect Implementation**:

> "The system shall support mode X for capability Y."

**Why Incorrect**: This is Domain Knowledge, not Implementation.

**Correct**: Implementation documents how mode X is realized, not that it exists.

### Error 2: Including Architecture in Implementation

**Incorrect Implementation**:

> "The component uses message passing for inter-component communication."

**Why Incorrect**: This is Architecture (organization), not Implementation (specific technology).

**Correct**: Implementation specifies which message broker and topic structure.

### Error 3: Confusing Technology Choice with Domain Requirement

**Incorrect**:

> "The system shall use MongoDB for data storage because MongoDB is the best choice."

**Why Incorrect**: Technology preference is not a domain requirement.

**Correct**:

> "The system shall store time-series data with configurable retention." (Domain Knowledge)
> "MongoDB shall be used for data storage." (Implementation - but document the rationale)

## Implementation and Authority

Implementation derives authority from Architecture:

```
Domain Knowledge (Authorizes)
        │
        │ Derivation
        ↓
Architecture (Guides)
        │
        │ Authorization
        ↓
Implementation (Authorized)
        │
        │ Subject to
        ↓
Verification (Confirms alignment)
```

Implementation may not contradict Architecture. If implementation cannot conform to Architecture, the situation must be resolved through clarification or Architecture change.

## Summary

Implementation is a separate methodology phase that:

- **Explains** how Architecture is realized using specific technologies
- **Follows** Architecture and precedes Verification
- **Includes** source code, configurations, and technology-specific details
- **Traces** to Architecture for authorization
- **Is distinct from** Domain Knowledge (intent) and Architecture (organization)
- **Realizes** Domain Interfaces using specific technologies

Understanding Implementation as a distinct phase is essential for maintaining separation of concerns.

---

## Version

- **Document Version**: 2.0
- **Effective Date**: 2026-07-13
- **Change Note**: Updated to Domain Interface terminology and generalized examples
