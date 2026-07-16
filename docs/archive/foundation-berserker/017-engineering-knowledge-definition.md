# Engineering Knowledge Definition

## Purpose

This document provides a comprehensive definition of **Engineering Knowledge** in KDSE. Engineering Knowledge is the authoritative, implementation-independent understanding derived from Reference Artifacts through structured analysis.

## Canonical Definition

**Engineering Knowledge** is implementation-independent understanding that describes engineering purpose, behavior, intent, and constraints. Engineering Knowledge remains valid if the implementation is completely rewritten.

## What Engineering Knowledge Describes

Engineering Knowledge shall describe:

### Engineering Purpose

- Why the system or component exists
- What problem it solves
- What value it provides

### Plant Behavior

- How the physical system behaves
- What responses to expect from equipment
- What operational characteristics exist

### Algorithms

- How calculations are performed
- What logic governs behavior
- What decision criteria apply

### Engineering Intent

- What the designer intended
- Why specific approaches were chosen
- What goals the system serves

### Operating Modes

- What modes the system supports
- How modes are entered and exited
- What behavior differs between modes

### Engineering Constraints

- What limits apply to operation
- What boundaries must not be exceeded
- What conditions must be maintained

### Engineering Assumptions

- What is assumed to be true
- What conditions are expected
- What dependencies exist

### Safety Behavior

- What safety functions exist
- How safety is maintained
- What failsafe states apply

### Control Philosophy

- How control is organized
- What control hierarchy exists
- How control loops interact

### Engineering State Machines

- What states are possible
- What transitions are valid
- What triggers state changes

### Domain Interfaces

- What information is exchanged
- What the information means
- What units are used

## What Domain Knowledge Is Not

Domain Knowledge is **not**:

### Implementation Details

- Programming language syntax
- Runtime environment
- Communication protocols
- Framework-specific patterns

### Vendor Implementations

- Specific product features
- Proprietary extensions
- Vendor-specific configurations

### Technology Choices

- Language selection
- Platform decisions
- Tool choices
- Infrastructure decisions

## The Engineering Independence Principle

Engineering Knowledge shall remain valid if the implementation is completely rewritten.

### The Independence Test

Every Engineering Knowledge statement shall pass the following validation:

> "If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"

| Result | Classification | Action |
|--------|---------------|--------|
| YES | Engineering Knowledge | Retain in Knowledge Artifacts |
| NO | Architecture or Implementation | Move to appropriate artifact type |

### Examples

**Engineering Knowledge** (passes the test):

- "The inverter shall support grid-forming control mode"
- "The system shall respond to grid faults within 100ms"
- "Power shall be limited to prevent exceeding equipment ratings"

**Not Engineering Knowledge** (fails the test):

- "The inverter shall use Modbus TCP for communication" → Implementation
- "The control logic shall be implemented in Node-RED" → Implementation
- "The system shall use IEC 61850 for substation communication" → Architecture/Implementation

### Why Implementation Independence Matters

Implementation independence ensures:

1. **Longevity**: Knowledge remains valid across technology changes
2. **Reusability**: Knowledge applies to new implementations
3. **Clarity**: Focus is on engineering, not technology
4. **Authority**: Knowledge is not tied to specific vendors or tools

## Categories of Engineering Knowledge

### Category 1: Functional Knowledge

Describes what the system does.

**Examples**:

- "The system shall measure AC voltage at the point of common coupling"
- "The inverter shall convert DC power to AC power"
- "The controller shall regulate reactive power output"

### Category 2: Behavioral Knowledge

Describes how the system behaves.

**Examples**:

- "When grid voltage drops below 0.9 pu, the inverter shall enter low-voltage ride-through mode"
- "The system shall ramp power output at no more than 10% per minute"
- "Upon loss of grid, the system shall disconnect within 2 seconds"

### Category 3: Constraint Knowledge

Describes limitations and boundaries.

**Examples**:

- "Active power shall not exceed the inverter nameplate rating"
- "The system shall operate within -10°C to +50°C ambient temperature"
- "Reactive power shall be limited to maintain power factor between 0.85 lagging and 0.85 leading"

### Category 4: Intentional Knowledge

Describes why decisions were made.

**Examples**:

- "Grid-forming mode is required to enable black start capability"
- "The 100ms fault response time is specified to meet grid code requirements"
- "Distributed control is chosen to improve system resilience"

### Category 5: Interface Knowledge

Describes what information is exchanged.

**Examples**:

- "The Reference Meter provides normalized electrical measurements"
- "The Plant Status Interface exposes current operating state"
- "The Historian Interface provides historical data access"

## Engineering Knowledge Structure

Engineering Knowledge Artifacts shall contain:

### Statement

The Engineering Knowledge statement itself.

### Category

The category of knowledge (Functional, Behavioral, Constraint, Intentional, Interface).

### Evidence Sources

The Reference Artifacts that support this knowledge.

### Evidence Strength

The strength of supporting evidence (★★★★★ to ★☆☆☆☆).

### Dependencies

Other Engineering Knowledge that this statement depends upon.

### Dependents

Other Engineering Knowledge or Architecture that depends upon this statement.

### Validation Status

Whether the statement has been validated.

### Review Status

The review state (Draft, Reviewed, Approved).

## Characteristics of Engineering Knowledge

### Characteristic 1: Implementation Independence

Engineering Knowledge does not depend on:

- Programming language
- Runtime environment
- Communication protocol
- Vendor implementation
- Software framework
- Platform technology

### Characteristic 2: Traceability

Engineering Knowledge traces to:

- Reference Artifacts (source evidence)
- Derived Architecture (downstream usage)
- Verification Criteria (validation basis)

### Characteristic 3: Authority

Engineering Knowledge carries authority because:

- It is validated against evidence
- It follows the Knowledge Derivation Lifecycle
- It is reviewed and approved
- It is independent of implementation

### Characteristic 4: Completeness

Engineering Knowledge is complete when:

- All necessary evidence has been analyzed
- All identified gaps have been addressed
- Dependencies are documented
- Contradictions are preserved

## Engineering Knowledge and Other Artifacts

### Relationship to Reference Artifacts

Reference Artifacts are the evidence. Engineering Knowledge is the derived understanding.

```
Reference Artifact → Evidence → Engineering Knowledge
```

Reference Artifacts support Engineering Knowledge. They do not replace it.

### Relationship to Architecture

Architecture explains how Engineering Knowledge is organized into software.

```
Engineering Knowledge → Derivation → Architecture
```

Architecture decisions must trace to Engineering Knowledge.

### Relationship to Implementation

Implementation realizes Architecture using specific technologies.

```
Architecture → Realization → Implementation
```

Implementation must not contradict Engineering Knowledge.

## Examples of Engineering Knowledge

### Example 1: Power Rating

**Engineering Knowledge**:

> "The power conversion system shall be rated for 500kW continuous operation."

**Analysis**:

- Implementation-independent: ✓ (Rating is an engineering specification)
- Traces to evidence: ✓ (RA-001 P&ID, RA-003 Vendor Manual)
- Evidence Strength: ★★★★☆ (Project doc + Vendor manual)

**Not this**:

> "The inverter shall be model XYZ-500 from vendor ABC." → Vendor Implementation

### Example 2: Control Mode

**Engineering Knowledge**:

> "The system shall support grid-following and grid-forming control modes."

**Analysis**:

- Implementation-independent: ✓ (Modes are engineering capabilities)
- Traces to evidence: ✓ (RA-001 Specification, RA-002 Node-RED, RA-003 Manual)
- Evidence Strength: ★★★★★ (All sources agree)

**Not this**:

> "Grid-forming shall be implemented using Droop Control algorithm in the DSP." → Implementation

### Example 3: Fault Response

**Engineering Knowledge**:

> "Upon sensing a grid fault, the inverter shall enter the fault ride-through mode specified by the applicable grid code."

**Analysis**:

- Implementation-independent: ✓ (Fault response is an engineering requirement)
- Traces to evidence: ✓ (RA-001 Grid code reference, RA-002 Implementation)
- Evidence Strength: ★★★☆☆ (Project doc only)

**Not this**:

> "The fault detection shall use voltage threshold comparison with 5-sample debounce." → Implementation

### Example 4: Data Interface

**Engineering Knowledge**:

> "The Reference Meter provides normalized electrical measurements including voltage, current, frequency, and power."

**Analysis**:

- Implementation-independent: ✓ (What information is provided, not how)
- Traces to evidence: ✓ (RA-001 Documentation, RA-002 Implementation)
- Evidence Strength: ★★★★☆ (Multiple sources)

**Not this**:

> "Measurements arrive through Modbus register 40001." → Implementation

## Common Errors

### Error 1: Embedding Implementation Details

**Incorrect**:

> "The controller shall use MQTT to publish status updates to the SCADA topic."

**Correct**:

> "The controller shall provide real-time status updates to the operator interface."

### Error 2: Including Vendor Specifics

**Incorrect**:

> "The inverter shall use ABB's ProQ flow control algorithm."

**Correct**:

> "The inverter shall regulate reactive power output to meet grid code requirements."

### Error 3: Specifying Technology

**Incorrect**:

> "The historian shall use InfluxDB time-series database."

**Correct**:

> "The historian shall store time-series data with millisecond precision."

## Validation Checklist

When creating or reviewing Engineering Knowledge, verify:

- [ ] Does this statement remain true if the implementation is rewritten?
- [ ] Does this statement describe engineering, not technology?
- [ ] Is this statement traceable to Reference Artifacts?
- [ ] Is the Evidence Strength assessed?
- [ ] Are dependencies documented?
- [ ] Are contradictions preserved if any exist?

## Summary

Engineering Knowledge is implementation-independent understanding that:

- **Describes** engineering purpose, behavior, intent, and constraints
- **Remains valid** even if implementation is completely rewritten
- **Does not depend on** programming language, runtime, protocol, or vendor
- **Traces to** Reference Artifacts as evidence
- **Carries authority** through validation and review
- **Is distinct from** Architecture and Implementation

Understanding the definition of Engineering Knowledge is essential for maintaining the separation of concerns that makes KDSE effective.

---

## Version

- **Document Version**: 1.0
- **Effective Date**: 2026-07-13
- **Change Note**: Initial release establishing formal Engineering Knowledge definition
