# Engineering Independence Test

## Purpose

This document establishes the **Engineering Independence Test** as a validation mechanism in KDSE. Every derived Engineering Knowledge statement shall pass this test. The test ensures that Engineering Knowledge remains valid if the implementation is completely rewritten.

## The Test

### Canonical Formulation

Every Engineering Knowledge statement shall pass the following validation:

> **"If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"**

| Result | Classification | Action |
|--------|---------------|--------|
| YES | Engineering Knowledge | Retain in Knowledge Artifacts |
| NO | Architecture or Implementation | Move to appropriate artifact type |

### The Test Rationale

The test ensures that:

1. **Engineering Knowledge is implementation-independent**: It describes engineering, not technology
2. **Separation of concerns is maintained**: Knowledge, Architecture, and Implementation remain distinct
3. **Knowledge has longevity**: It remains valid across technology changes
4. **Traceability is meaningful**: Links between artifacts represent real relationships

## Applying the Test

### Step 1: State the Knowledge

Clearly articulate the statement to be tested.

### Step 2: Ask the Question

Apply the test question to the statement.

### Step 3: Evaluate Independence

Consider whether the statement would remain true under various changes.

### Step 4: Classify the Result

Classify the statement based on the test outcome.

## Test Examples

### Example 1: Engineering Knowledge Passing

**Statement**: "The inverter shall support grid-forming control mode."

**Test**: "If we rewrote the system in a different language, using different protocols, would this statement remain true?"

**Evaluation**: Yes. Grid-forming is an engineering capability that exists regardless of implementation.

**Classification**: Engineering Knowledge ✓

---

### Example 2: Implementation Failing

**Statement**: "The inverter shall use Modbus TCP for communication."

**Test**: "If we rewrote the system using a different protocol, would this statement remain true?"

**Evaluation**: No. Modbus TCP is a specific technology choice.

**Classification**: Implementation ✗

**Action**: Remove from Engineering Knowledge, document as Implementation decision.

---

### Example 3: Architecture Failing

**Statement**: "The control service shall use message passing for inter-service communication."

**Test**: "If we changed to a different architectural pattern, would this statement remain true?"

**Evaluation**: No. Message passing is an architectural pattern choice.

**Classification**: Architecture ✗

**Action**: Remove from Engineering Knowledge, document as Architecture decision.

---

### Example 4: Gray Area - Requires Judgment

**Statement**: "The historian shall store time-series data with millisecond precision."

**Test**: "If we changed database technology, would this statement remain true?"

**Evaluation**: Partially. The precision requirement is engineering, but the specific precision threshold might be implementation-dependent.

**Analysis**:

- Is millisecond precision an engineering requirement? Yes
- Is millisecond specifically required? Possibly implementation choice

**Resolution**: Refine the statement to separate engineering from implementation:
- Engineering: "The historian shall store time-series data with sufficient precision for analysis."
- Implementation: "Millisecond precision shall be used for time-series storage."

---

### Example 5: System Behavior

**Statement**: "Upon grid fault, the inverter shall enter fault ride-through mode."

**Test**: "If we implemented this in a different language, would this statement remain true?"

**Evaluation**: Yes. The system behavior is an engineering requirement.

**Classification**: Engineering Knowledge ✓

---

### Example 6: Protocol Specification

**Statement**: "Fault detection shall use voltage threshold comparison."

**Test**: "If we used a different detection method, would this statement remain true?"

**Evaluation**: This is ambiguous. "Voltage threshold comparison" describes both engineering intent (detect by voltage) and implementation approach (comparison).

**Resolution**: Refine:
- Engineering: "Fault detection shall identify grid voltage anomalies."
- Architecture/Implementation: "Voltage threshold comparison shall be used for fault detection."

## Independence Dimensions

The test considers independence across multiple dimensions:

### Dimension 1: Programming Language

**Test Change**: Rewrite in Python, Java, Go, Rust, etc.

| Statement | Result |
|-----------|--------|
| "The system shall process measurements in real-time" | PASS |
| "The system shall use Python for processing" | FAIL |

### Dimension 2: Communication Protocol

**Test Change**: Use Modbus, DNP3, IEC 61850, OPC UA, MQTT, REST, etc.

| Statement | Result |
|-----------|--------|
| "The system shall exchange sensor data with controllers" | PASS |
| "The system shall use Modbus TCP for data exchange" | FAIL |

### Dimension 3: Runtime Environment

**Test Change**: Run on bare metal, containers, VMs, cloud, edge, etc.

| Statement | Result |
|-----------|--------|
| "The system shall operate continuously" | PASS |
| "The system shall run in Docker containers" | FAIL |

### Dimension 4: Software Framework

**Test Change**: Use different frameworks, libraries, or platforms.

| Statement | Result |
|-----------|--------|
| "The API shall provide access to historian data" | PASS |
| "The API shall use Node.js Express framework" | FAIL |

### Dimension 5: Vendor/Product

**Test Change**: Use different vendor products or components.

| Statement | Result |
|-----------|--------|
| "Power conversion shall be rated at 500kW" | PASS |
| "ABB inverters shall be used" | FAIL |

### Dimension 6: Platform

**Test Change**: Run on different hardware or infrastructure.

| Statement | Result |
|-----------|--------|
| "The system shall operate in ambient temperatures to 50°C" | PASS |
| "The system shall run on ARM processors" | FAIL |

## Common Test Scenarios

### Scenario 1: Control System

| Statement | Test Result | Classification |
|-----------|-------------|----------------|
| "System supports grid-following and grid-forming modes" | PASS | Engineering Knowledge |
| "Grid-following uses droop control algorithm" | AMBIGUOUS | Refine statement |
| "Control logic implemented in Node-RED" | FAIL | Implementation |
| "Grid-following is a separate service" | FAIL | Architecture |

### Scenario 2: Communication System

| Statement | Test Result | Classification |
|-----------|-------------|----------------|
| "Sensor data is exchanged between components" | PASS | Engineering Knowledge |
| "Components communicate via message passing" | PASS | Architecture |
| "MQTT is used with topic structure /plant/sensors" | FAIL | Implementation |

### Scenario 3: Data Management

| Statement | Test Result | Classification |
|-----------|-------------|----------------|
| "Time-series data is stored for historical analysis" | PASS | Engineering Knowledge |
| "Data is stored with second-level precision" | PASS | Engineering Knowledge (with judgment) |
| "InfluxDB is used for time-series storage" | FAIL | Implementation |

### Scenario 4: User Interface

| Statement | Test Result | Classification |
|-----------|-------------|----------------|
| "Operators can monitor system status" | PASS | Engineering Knowledge |
| "Status is displayed on operator dashboard" | PASS | Architecture |
| "Dashboard is built with React" | FAIL | Implementation |

## The Test and Artifact Hierarchy

The test reinforces the separation between artifact types:

```
┌─────────────────────────────────────────────────────────┐
│                    Engineering Knowledge                  │
│            (Passes Engineering Independence Test)        │
├─────────────────────────────────────────────────────────┤
│                       Architecture                       │
│           (May fail Engineering Independence Test)       │
│          (Describes organization, not technology)        │
├─────────────────────────────────────────────────────────┤
│                      Implementation                      │
│              (Fails Engineering Independence Test)       │
│             (Describes specific technology)             │
└─────────────────────────────────────────────────────────┘
```

## Test Documentation

### Test Log Entry

```
| Statement | Test Question | Evaluation | Result |
|-----------|---------------|------------|--------|
| "Grid-forming supported" | Rewrite in different language? | Yes | PASS |
| "Uses Modbus TCP" | Change protocol? | No | FAIL |
```

### Test Detail Entry

```
## Test: "Grid-forming supported"

**Statement**: The system shall support grid-forming control mode.

**Test Question**: If the implementation were rewritten tomorrow
using a different programming language, communication protocol,
runtime, framework, vendor, or platform, would this statement
still remain true?

**Evaluation**:
- Different language: Yes - grid-forming is not language-specific
- Different protocol: Yes - grid-forming capability exists regardless
- Different vendor: Yes - the capability requirement persists
- Different platform: Yes - the engineering requirement remains

**Result**: PASS - Engineering Knowledge

**Classification**: Engineering Knowledge

**Notes**: None
```

## Test Limitations

### Limitation 1: Gray Areas Exist

Some statements fall in gray areas requiring judgment.

**Example**: "The system shall respond within 100ms."

- Is 100ms an engineering requirement? Yes
- Is 100ms specifically required? Could be implementation constraint

**Resolution**: Apply judgment, refine ambiguous statements.

### Limitation 2: Multiple Interpretations

The same statement can be interpreted differently.

**Example**: "The historian shall use time-series database."

- Interpretation 1: Time-series capability is required → Engineering Knowledge
- Interpretation 2: A specific database type is required → Implementation

**Resolution**: Refine to clarify intent.

### Limitation 3: Context Matters

Context affects test outcomes.

**Example**: "The system shall use IEC 61850."

- In context of engineering requirements: IEC 61850 is a standard, implementation-independent
- In context of specific product: IEC 61850 is a technology choice

**Resolution**: Consider context when applying the test.

## Common Test Errors

### Error 1: Confusing Engineering Intent with Implementation

**Incorrect Test**: "The statement mentions 'Modbus', which is implementation, so FAIL."

**Why Incorrect**: Some mentions of technology represent engineering requirements (e.g., IEC 61850 as a standard).

**Correct Test**: Evaluate whether the statement describes an engineering requirement or a technology choice.

### Error 2: Overly Strict Interpretation

**Incorrect Test**: "Any mention of specific technology FAILs the test."

**Why Incorrect**: Engineering requirements may reference specific standards or capabilities.

**Correct Test**: Evaluate whether the statement would remain true across technology changes.

### Error 3: Ignoring Context

**Incorrect Test**: "The statement mentions 'REST API', which is implementation, so FAIL."

**Why Incorrect**: The need for API access might be engineering; the specific technology is implementation.

**Correct Test**: Separate the engineering requirement from the implementation choice.

## Test Integration

### Integration with Knowledge Derivation

The test is applied during Knowledge Derivation:

```
Evidence Identified
        ↓
Formulate Knowledge Statement
        ↓
Apply Engineering Independence Test
        │
        ├── PASS ──→ Retain as Engineering Knowledge
        │
        └── FAIL ──→ Reclassify (Architecture or Implementation)
```

### Integration with Review

The test is applied during Knowledge Review:

```
Knowledge Artifact Review
        │
        ▼
Verify Independence Test Applied
        │
        ├── Applied and passed ──→ Approve
        │
        └── Not applied or failed ──→ Return for correction
```

## Summary

The Engineering Independence Test ensures that:

- **Engineering Knowledge is implementation-independent**
- **Statements remain valid across technology changes**
- **Separation of concerns is maintained**
- **Knowledge has longevity beyond specific implementations**

The test is a simple question with profound implications for methodology integrity.

---

## Version

- **Document Version**: 1.0
- **Effective Date**: 2026-07-13
- **Change Note**: Initial release establishing Engineering Independence Test
