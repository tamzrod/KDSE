# KDSE Engineering Session Protocol (ESP) Specification

**Document Version:** 1.0  
**Status:** NORMATIVE  
**Type:** Protocol Specification  
**Effective Date:** 2026-07-19  

---

## 1. Purpose

This document defines the **KDSE Engineering Session Protocol (ESP)** — the authoritative specification governing all interactions between AI agents and KDSE Runtime environments.

ESP establishes the complete contractual framework for engineering sessions, encompassing:

- **Identity establishment** — Determination of participant roles and authorities
- **Context establishment** — Creation and verification of engineering context
- **Session lifecycle** — Definition of all valid session states and transitions
- **Failure handling** — Deterministic recovery procedures for all failure modes
- **Completion criteria** — Unambiguous conditions for session activation

This protocol is **normative** — all compliant implementations MUST produce identical engineering contexts regardless of underlying AI model, execution environment, or runtime variant.

---

## 2. Scope

### 2.1 In Scope

This protocol governs:

1. **Engineering session establishment** — The complete handshake from idle state to active engineering session
2. **Bootstrap procedures** — The deterministic process establishing engineering context
3. **Authority determination** — The resolution of conflicting claims through authoritative sources
4. **Session state management** — Legal states, transitions, and terminal conditions
5. **Failure recovery** — Deterministic handling of all defined failure scenarios
6. **Completion verification** — Criteria establishing when engineering reasoning may begin

### 2.2 Out of Scope

This protocol does NOT govern:

1. **Engineering methodology** — Phase definitions, artifact specifications, and engineering rules (see KDSE Methodology)
2. **Runtime implementation** — Internal architecture, algorithms, or data structures
3. **AI reasoning** — The internal processes of the AI agent during active engineering
4. **External integrations** — Interactions with external systems, repositories, or services
5. **Report generation** — The format, content, or production of engineering reports
6. **Knowledge extraction** — Procedures for extracting knowledge from external sources

### 2.3 Compliance Level

| Level | Description | Applicability |
|-------|-------------|---------------|
| **FULL** | Complete protocol implementation | CLI Runtime, MCP Runtime |
| **PARTIAL** | Bootstrap and context only | Assessment Tools, Audit Systems |

---

## 3. Terminology

### 3.1 Core Terms

| Term | Definition |
|------|------------|
| **Engineering Session** | A bounded period of AI-assisted engineering work governed by KDSE methodology |
| **Engineering Context** | The complete set of information defining the current engineering state |
| **KDSE Runtime** | The authoritative runtime environment managing engineering state |
| **Workspace** | The filesystem location containing the project and KDSE runtime |
| **Bootstrap** | The deterministic process establishing engineering context |
| **Handshake** | The structured exchange establishing session identity and authority |
| **Engineering State** | The persisted and verifiable state of an engineering project |

### 3.2 Participant Terms

| Term | Definition |
|------|------------|
| **User** | The human operator initiating and authorizing engineering sessions |
| **AI Agent** | The artificial intelligence system performing engineering reasoning |
| **Runtime Adapter** | The thin adapter layer (CLI, MCP, IDE) mediating AI-agent interaction |
| **Workspace Engine** | The component owning and persisting all engineering state |

### 3.3 State Terms

| Term | Definition |
|------|------------|
| **IDLE** | No active engineering session exists |
| **DISCOVERING** | Runtime is locating and identifying the workspace |
| **AUTHENTICATING** | Runtime is verifying workspace ownership |
| **BOOTSTRAPPING** | Runtime is establishing engineering context |
| **VERIFYING** | Runtime is validating engineering context completeness |
| **ACTIVE** | Engineering session is established and AI reasoning may proceed |
| **SUSPENDED** | Session is paused but may be resumed |
| **TERMINATED** | Session has ended and cannot be resumed |
| **FAILED** | An unrecoverable error has occurred |

### 3.4 Authority Terms

| Term | Definition |
|------|------------|
| **Authoritative Source** | A verified, trusted origin of information taking precedence over inference |
| **Heuristic** | An inference-based determination used only when no authoritative source exists |
| **Fallback** | The predetermined order of sources to consult when higher-priority sources are unavailable |
| **Conflict Resolution** | The deterministic procedure when multiple sources provide contradictory information |

---

## 4. Protocol Goals

### 4.1 Primary Goals

**G-001: Deterministic Context Establishment**  
Every compliant implementation MUST produce the identical engineering context given identical inputs, regardless of AI model or runtime variant.

**G-002: Authority Precedence**  
Information SHALL originate from authoritative sources. The protocol defines an explicit hierarchy ensuring no ambiguity in source priority.

**G-003: Explicit Ownership**  
Every action has exactly one owner. The protocol defines ownership boundaries preventing overlap or gaps.

**G-004: Separation of Concerns**  
Bootstrap establishes engineering context. Engineering performs reasoning. These operations MUST NOT overlap or interfere.

**G-005: Implementation Independence**  
The protocol specifies behavior without dependence on any programming language, framework, or implementation technology.

### 4.2 Secondary Goals

**G-006: Failure Determinism**  
Every failure mode has exactly one valid recovery procedure, eliminating ambiguity in error handling.

**G-007: Session Integrity**  
The protocol ensures engineering sessions maintain consistent state throughout their lifecycle.

**G-008: Auditability**  
All significant state transitions are traceable to authoritative evidence.

---

## 5. Non-Goals

The following are explicitly OUT OF SCOPE for this protocol:

**NG-001: AI Reasoning**  
The internal processes by which an AI agent reaches engineering decisions are not governed by this protocol.

**NG-002: Methodology Specification**  
Engineering phases, artifact requirements, and engineering rules are defined in the KDSE Methodology, not this protocol.

**NG-003: Report Format**  
The content and structure of engineering reports are not defined by this protocol.

**NG-004: Knowledge Extraction**  
Procedures for extracting knowledge from external references are outside protocol scope.

**NG-005: Runtime Internals**  
Implementation details, internal data structures, and algorithm choices are not constrained.

**NG-006: Performance Requirements**  
Timing, latency, and throughput specifications are not part of this protocol.

---

## 6. Participants (Actors)

### 6.1 User

| Attribute | Specification |
|-----------|---------------|
| **Role** | Human operator initiating and authorizing engineering sessions |
| **Authority** | Ultimate authority for engineering decisions; can override protocol constraints |
| **Responsibility** | Define session scope, authorize phase transitions, approve implementations |
| **Ownership** | Owns the engineering intent and business requirements |
| **Boundary** | Cannot directly modify runtime state; operates through protocol-defined interfaces |
| **Input Forms** | Direct command, configuration file, session parameters |

### 6.2 AI Agent

| Attribute | Specification |
|-----------|---------------|
| **Role** | Autonomous reasoning system performing engineering work |
| **Authority** | Limited to reasoning and recommendation; requires User authorization for implementation |
| **Responsibility** | Follow KDSE methodology, produce evidence-backed recommendations, maintain traceability |
| **Ownership** | Owns the reasoning process and intermediate artifacts |
| **Boundary** | Cannot bypass authority hierarchy; cannot modify authoritative state without protocol |
| **Constraint** | MUST NOT begin engineering reasoning until session reaches ACTIVE state |

### 6.3 KDSE Runtime

| Attribute | Specification |
|-----------|---------------|
| **Role** | Authoritative keeper of engineering state |
| **Authority** | Highest authority for engineering context; supersedes all claims by AI agents |
| **Responsibility** | Establish and maintain engineering context, enforce methodology, validate evidence |
| **Ownership** | Owns all persisted engineering state (`.kdse/`) |
| **Boundary** | Does not perform reasoning; does not interpret business requirements |
| **Variants** | CLI Runtime, MCP Runtime, IDE Runtime — all must produce identical context |

### 6.4 Repository

| Attribute | Specification |
|-----------|---------------|
| **Role** | Container for project artifacts and KDSE runtime state |
| **Authority** | Provides authoritative filesystem-level evidence of project state |
| **Responsibility** | Persist engineering artifacts, provide verifiable evidence of state |
| **Ownership** | Jointly owned: Project Layer by project, Runtime Layer by KDSE |
| **Boundary** | Does not interpret or validate content; only provides storage |
| **Constraint** | Must support atomic operations for state consistency |

### 6.5 Operating System

| Attribute | Specification |
|-----------|---------------|
| **Role** | Provider of filesystem access and process execution environment |
| **Authority** | Controls resource access, file permissions, and process isolation |
| **Responsibility** | Execute runtime processes, maintain filesystem consistency, enforce permissions |
| **Ownership** | Owns system resources and kernel-level state |
| **Boundary** | Does not understand or interpret KDSE concepts |

### 6.6 Git

| Attribute | Specification |
|-----------|---------------|
| **Role** | Version control system providing historical evidence and branch state |
| **Authority** | Authoritative for version history, commit state, and branch relationships |
| **Responsibility** | Track artifact changes, provide historical evidence, validate consistency |
| **Ownership** | Owns version control metadata and history |
| **Boundary** | Not required for protocol compliance; available as authoritative fallback |
| **Constraint** | Git metadata is authoritative ONLY when filesystem evidence is insufficient |

### 6.7 External Knowledge Sources

| Attribute | Specification |
|-----------|---------------|
| **Role** | Provider of domain-specific knowledge (standards, specifications, references) |
| **Authority** | Authoritative for domain knowledge within their scope |
| **Responsibility** | Provide verifiable knowledge artifacts, maintain reference integrity |
| **Ownership** | Owned by external bodies (standards organizations, vendors, etc.) |
| **Boundary** | Cannot directly modify engineering context; accessed through explicit protocols |
| **Constraint** | Not required for protocol compliance; provides authoritative evidence when referenced |

---

## 7. Authority Hierarchy

### 7.1 Hierarchy Definition

Authority descends through the following hierarchy. Higher authority ALWAYS supersedes lower authority:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        KDSE STANDARD (Highest)                          │
│                                                                         │
│  The authoritative specification of KDSE methodology and requirements  │
│  Defined in: KDSE Standard Documents                                  │
│  Authority: Establishes all lower-level authority                     │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           RUNTIME MANIFEST                             │
│                                                                         │
│  The authoritative declaration of runtime configuration and state      │
│  Stored in: .kdse/runtime.yaml                                         │
│  Authority: Defines what is and is not a KDSE project                │
│  Supersedes: All claims by AI agents or users                         │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         PROJECT MANIFEST                               │
│                                                                         │
│  The authoritative declaration of project configuration and state      │
│  Stored in: .kdse/project.yaml                                        │
│  Authority: Defines project-specific parameters                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          SESSION STATE                                 │
│                                                                         │
│  The authoritative record of current session parameters and state      │
│  Stored in: .kdse/session.yaml                                        │
│  Authority: Defines current session context                           │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          AI REASONING (Lowest)                          │
│                                                                         │
│  The AI agent's interpretations, recommendations, and decisions       │
│  Authority: NONE until authorized by higher authorities               │
│  Constraint: MUST NOT contradict higher authorities                   │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Conflict Resolution Rules

**Rule CR-001: Runtime Manifest Priority**  
When the Runtime Manifest (.kdse/) contradicts any other source, the Runtime Manifest is authoritative.

**Rule CR-002: Explicit User Input Priority**  
When User provides explicit input through protocol-defined interfaces, it supersedes AI inference.

**Rule CR-003: Filesystem Evidence Priority**  
When filesystem artifacts contradict claims (verbal or documented), filesystem evidence prevails.

**Rule CR-004: Git Metadata Fallback**  
Git metadata (commit, branch, tag) is authoritative ONLY when filesystem evidence is unavailable or insufficient.

**Rule CR-005: Timestamp Ordering**  
When multiple authoritative sources conflict and no other resolution rule applies, the most recent timestamp wins.

**Rule CR-006: No Resolution by Authority Level**  
When two sources at the same authority level conflict, the source with stronger evidence (more specific, more verifiable) wins.

### 7.3 Authority Boundaries

| Authority Level | Can Override | Cannot Override |
|-----------------|--------------|-----------------|
| KDSE Standard | None | All lower authorities |
| Runtime Manifest | AI claims, User claims | KDSE Standard |
| Project Manifest | AI claims | Runtime Manifest, KDSE Standard |
| Session State | None | All higher authorities |
| AI Reasoning | None | All higher authorities |

---

## 8. Engineering Session Handshake

### 8.1 Handshake Overview

The Engineering Session Handshake is a multi-phase exchange inspired by established industrial protocols (TCP, Modbus, OPC UA). It establishes identity, authority, ownership, and engineering context before permitting engineering work.

### 8.2 Handshake Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      ENGINEERING SESSION HANDSHAKE                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   AI Agent              Runtime Adapter              KDSE Runtime         │
│      │                        │                          │              │
│      │──── CONNECT ─────────▶│                          │              │
│      │                        │──── DISCOVER ──────────▶│              │
│      │                        │◀─── DISCOVER_ACK ───────│              │
│      │                        │                          │              │
│      │                        │──── AUTHENTICATE ───────▶│              │
│      │                        │◀─── AUTHENTICATE_ACK ────│              │
│      │                        │                          │              │
│      │◀─── CONTEXT_REQUEST ──│                          │              │
│      │──── CONTEXT_QUERY ───▶│                          │              │
│      │                        │──── BOOTSTRAP ─────────▶│              │
│      │                        │◀─── BOOTSTRAP_ACK ───────│              │
│      │                        │                          │              │
│      │                        │──── VERIFY ────────────▶│              │
│      │                        │◀─── VERIFY_ACK ──────────│              │
│      │                        │                          │              │
│      │◀─── CONTEXT_RESPONSE ──│                          │              │
│      │                        │                          │              │
│      │◀══════════ SESSION ESTABLISHED ═══════════════════│              │
│      │                                                         
│      │                        │                          │              │
│      │ (Engineering reasoning may proceed)                │              │
│      │                                                         
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Handshake Messages

#### 8.3.1 CONNECT

| Attribute | Value |
|-----------|-------|
| **Direction** | AI Agent → Runtime Adapter |
| **Purpose** | Initiate session establishment |
| **Precondition** | None |
| **Payload** | `{ "client_id": "<string>", "capabilities": ["<string>"] }` |
| **Response** | DISCOVER_ACK |

#### 8.3.2 DISCOVER

| Attribute | Value |
|-----------|-------|
| **Direction** | Runtime Adapter → KDSE Runtime |
| **Purpose** | Locate and identify workspace |
| **Precondition** | CONNECT received |
| **Payload** | `{ "workspace_path": "<path>", "workspace_type": "local|remote" }` |
| **Response** | DISCOVER_ACK |
| **Verification** | Confirm .kdse/ exists and is readable |

#### 8.3.3 DISCOVER_ACK

| Attribute | Value |
|-----------|-------|
| **Direction** | KDSE Runtime → Runtime Adapter |
| **Purpose** | Acknowledge workspace discovery |
| **Payload** | `{ "status": "found|not_found|unauthorized", "workspace_id": "<uuid>", "runtime_version": "<version>" }` |
| **On Failure** | Handshake terminates; return error code |

#### 8.3.4 AUTHENTICATE

| Attribute | Value |
|-----------|-------|
| **Direction** | Runtime Adapter → KDSE Runtime |
| **Purpose** | Verify workspace ownership and permissions |
| **Precondition** | DISCOVER_ACK received with status "found" |
| **Payload** | `{ "workspace_id": "<uuid>", "authentication_token": "<token>" }` |
| **Response** | AUTHENTICATE_ACK |
| **Verification** | Confirm read/write permissions on .kdse/ |

#### 8.3.5 AUTHENTICATE_ACK

| Attribute | Value |
|-----------|-------|
| **Direction** | KDSE Runtime → Runtime Adapter |
| **Purpose** | Acknowledge authentication |
| **Payload** | `{ "status": "authorized|unauthorized|readonly", "session_id": "<uuid>", "permissions": ["<permission>"] }` |
| **On Unauthorized** | Handshake terminates; return error code |

#### 8.3.6 CONTEXT_REQUEST

| Attribute | Value |
|-----------|-------|
| **Direction** | KDSE Runtime → AI Agent |
| **Purpose** | Request engineering context specification |
| **Precondition** | AUTHENTICATE_ACK received with status "authorized" or "readonly" |
| **Payload** | `{ "required_context": ["<field>"], "optional_context": ["<field>"] }` |
| **Response** | CONTEXT_QUERY |

#### 8.3.7 CONTEXT_QUERY

| Attribute | Value |
|-----------|-------|
| **Direction** | AI Agent → Runtime Adapter |
| **Purpose** | Provide context parameters |
| **Precondition** | CONTEXT_REQUEST received |
| **Payload** | `{ "session_type": "<type>", "scope": "<scope>", "constraints": {}, "metadata": {} }` |
| **Response** | BOOTSTRAP |

#### 8.3.8 BOOTSTRAP

| Attribute | Value |
|-----------|-------|
| **Direction** | Runtime Adapter → KDSE Runtime |
| **Purpose** | Initiate engineering context establishment |
| **Precondition** | CONTEXT_QUERY received |
| **Payload** | `{ "session_id": "<uuid>", "context_parameters": {...}, "template_version": "<version>" }` |
| **Response** | BOOTSTRAP_ACK |

#### 8.3.9 BOOTSTRAP_ACK

| Attribute | Value |
|-----------|-------|
| **Direction** | KDSE Runtime → Runtime Adapter |
| **Purpose** | Acknowledge bootstrap completion |
| **Payload** | `{ "status": "success|incomplete|failed", "engineering_context": {...}, "warnings": [] }` |
| **On Incomplete** | Proceed to VERIFY phase |
| **On Failed** | Handshake terminates; return error code |

#### 8.3.10 VERIFY

| Attribute | Value |
|-----------|-------|
| **Direction** | Runtime Adapter → KDSE Runtime |
| **Purpose** | Verify engineering context completeness |
| **Precondition** | BOOTSTRAP_ACK received with status "incomplete" |
| **Payload** | `{ "verification_scope": "full|partial", "verification_items": ["<item>"] }` |
| **Response** | VERIFY_ACK |

#### 8.3.11 VERIFY_ACK

| Attribute | Value |
|-----------|-------|
| **Direction** | KDSE Runtime → Runtime Adapter |
| **Purpose** | Acknowledge verification |
| **Payload** | `{ "status": "verified|failed", "verification_report": {...}, "issues": [] }` |
| **On Failed** | Handshake terminates; return error code |

#### 8.3.12 CONTEXT_RESPONSE

| Attribute | Value |
|-----------|-------|
| **Direction** | KDSE Runtime → AI Agent |
| **Purpose** | Deliver established engineering context |
| **Precondition** | VERIFY_ACK received with status "verified" |
| **Payload** | `{ "engineering_context": {...}, "session_token": "<token>", "expires_at": "<timestamp>" }` |

---

## 9. Bootstrap Protocol

### 9.1 Bootstrap Overview

The Bootstrap Protocol is the deterministic process by which the KDSE Runtime establishes the Engineering Context. This protocol is executed during the BOOTSTRAP phase of the Engineering Session Handshake.

### 9.2 Bootstrap State Machine

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           BOOTSTRAP PROTOCOL                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│    ┌──────────────┐                                                     │
│    │   START     │                                                     │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │  LOCATE     │◀── Locate .kdse/ directory                          │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │  VALIDATE   │◀── Verify .kdse/ structure                          │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │   LOAD      │◀── Load runtime manifest                            │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │  DETERMINE  │◀── Determine engineering context                    │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │   RESOLVE   │◀── Resolve authority conflicts                      │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │  ASSEMBLE   │◀── Assemble final context                           │
│    └──────┬───────┘                                                     │
│           │                                                             │
│           ▼                                                             │
│    ┌──────────────┐                                                     │
│    │   COMPLETE  │                                                     │
│    └──────────────┘                                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 9.3 Bootstrap Steps

#### Step B-001: LOCATE

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-001 |
| **Name** | LOCATE |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime owns filesystem discovery |
| **Preconditions** | Authentication completed |
| **Action** | Locate and verify existence of .kdse/ directory |
| **Output** | Confirmed .kdse/ path |
| **Next State** | VALIDATE |
| **Failure** | Return error: RUNTIME_NOT_FOUND |

#### Step B-002: VALIDATE

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-002 |
| **Name** | VALIDATE |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime owns structure validation |
| **Preconditions** | .kdse/ path confirmed |
| **Action** | Verify required runtime files exist and are readable |
| **Output** | Validation result with missing/invalid files |
| **Next State** | LOAD (if valid) or FAIL |
| **Required Files** | runtime.yaml, workspace.yaml, phase.yaml, session.yaml |
| **Failure** | Return error: RUNTIME_INVALID |

#### Step B-003: LOAD

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-003 |
| **Name** | LOAD |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime owns manifest loading |
| **Preconditions** | Validation successful |
| **Action** | Load and parse runtime manifest files |
| **Output** | Parsed runtime configuration |
| **Next State** | DETERMINE |
| **Failure** | Return error: MANIFEST_PARSE_ERROR |

#### Step B-004: DETERMINE

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-004 |
| **Name** | DETERMINE |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime determines context |
| **Preconditions** | Manifest loaded |
| **Action** | Determine all context fields from authoritative sources |
| **Output** | Raw context data from all sources |
| **Next State** | RESOLVE |
| **Authority Order** | 1. Runtime Manifest 2. Project Manifest 3. Session State 4. Heuristics |
| **Failure** | Return error: CONTEXT_INCOMPLETE |

#### Step B-005: RESOLVE

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-005 |
| **Name** | RESOLVE |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime resolves conflicts |
| **Preconditions** | Context data collected |
| **Action** | Apply conflict resolution rules for any contradictory values |
| **Output** | Resolved context with conflict log |
| **Next State** | ASSEMBLE |
| **Conflict Log** | Record all resolved conflicts for audit |

#### Step B-006: ASSEMBLE

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-006 |
| **Name** | ASSEMBLE |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime assembles final context |
| **Preconditions** | All conflicts resolved |
| **Action** | Assemble final engineering context structure |
| **Output** | Complete engineering context |
| **Next State** | COMPLETE |
| **Verification** | Verify assembled context meets completion criteria |

#### Step B-007: COMPLETE

| Attribute | Value |
|-----------|-------|
| **Step ID** | B-007 |
| **Name** | COMPLETE |
| **Actor** | KDSE Runtime |
| **Authority** | Runtime finalizes context |
| **Preconditions** | Context assembled and verified |
| **Action** | Mark context as established, record completion evidence |
| **Output** | Engineering context with establishment proof |
| **Next State** | END (bootstrap complete) |
| **Evidence** | Timestamp, checksum, establishment proof |

---

## 10. Engineering Context

### 10.1 Context Definition

The **Engineering Context** is the complete, verifiable, authoritative record of the current engineering state. It is the foundation upon which all engineering reasoning is performed.

### 10.2 Required Context Fields

| Field | Type | Description | Authority |
|-------|------|-------------|-----------|
| **workspace_id** | UUID | Unique identifier for the workspace | Runtime Manifest |
| **workspace_path** | Path | Absolute path to workspace root | Runtime Manifest |
| **runtime_version** | Version | KDSE Runtime version | Runtime Manifest |
| **runtime_type** | Enum | CLI, MCP, or IDE | Runtime Manifest |
| **session_id** | UUID | Current session identifier | Session State |
| **session_type** | Enum | New, Resumed, Assessment | Session State |
| **current_phase** | Enum | Current engineering phase | Phase State |
| **phase_history** | Array | Chronological phase transitions | Phase State |
| **project_id** | UUID | Project identifier | Project Manifest |
| **project_name** | String | Human-readable project name | Project Manifest |
| **template_version** | Version | KDSE template version | Runtime Manifest |
| **established_at** | Timestamp | Context establishment time | Runtime |
| **established_by** | String | Runtime instance identifier | Runtime |
| **context_checksum** | Hash | Integrity verification hash | Runtime |

### 10.3 Optional Context Fields

| Field | Type | Description | Authority |
|-------|------|-------------|-----------|
| **user_id** | UUID | Authenticated user identifier | User Input |
| **user_name** | String | Human-readable user name | User Input |
| **session_scope** | Object | Session scope parameters | User Input |
| **constraints** | Object | Engineering constraints | User Input |
| **external_refs** | Array | External knowledge references | External Sources |
| **previous_session** | UUID | Prior session for resume | Session State |
| **branch** | String | Git branch name | Git (Fallback) |
| **commit** | Hash | Git commit hash | Git (Fallback) |

### 10.4 Forbidden Assumptions

The following MUST NOT be assumed by any AI agent:

| Forbidden Assumption | Correct Behavior |
|--------------------|------------------|
| Project language is X | Determine from authoritative project files |
| Framework is Y | Determine from authoritative configuration |
| Build system is Z | Determine from authoritative build files |
| Current phase is P | Read from Runtime Manifest |
| User wants feature F | Receive explicit User authorization |
| Environment is production | Determine from explicit configuration |
| Dependencies are X, Y, Z | Determine from authoritative dependency files |

### 10.5 Context Completion Criteria

The Engineering Context is COMPLETE when ALL of the following are true:

| Criterion | Verification |
|-----------|--------------|
| **Required fields present** | All fields in Section 10.2 have non-null values |
| **Runtime manifest valid** | runtime.yaml parses successfully |
| **Phase valid** | Phase value is from approved phase list |
| **Session ID unique** | No duplicate session ID exists |
| **Workspace path exists** | Directory at specified path exists and is accessible |
| **Context checksum valid** | Computed checksum matches recorded checksum |

---

## 11. Decision Rules

### 11.1 Decision Authority

Every protocol decision MUST reference an authoritative source. The following rules define the decision hierarchy.

### 11.2 Project Identification

| Scenario | Decision Rule | Authority Source |
|----------|---------------|------------------|
| Determine workspace root | Always use .kdse/ parent directory | Filesystem |
| Verify KDSE project | Check .kdse/ exists | Filesystem |
| Project name | Read from project manifest | Project Manifest |
| Project type | Read from project manifest | Project Manifest |

### 11.3 Phase Determination

| Scenario | Decision Rule | Authority Source |
|----------|---------------|------------------|
| Current phase | Always read from phase.yaml | Runtime Manifest |
| Phase validity | Compare against approved phases | KDSE Standard |
| Phase history | Read from phase.yaml | Runtime Manifest |
| Phase transition | Validate against transition rules | KDSE Standard |

### 11.4 Session Determination

| Scenario | Decision Rule | Authority Source |
|----------|---------------|------------------|
| Session existence | Check for session.yaml | Filesystem |
| Session type | Read from session.yaml | Session State |
| Session parameters | Read from session.yaml | Session State |
| Session continuity | Verify session ID match | Session State |

### 11.5 Authority Resolution

| Scenario | Decision Rule | Priority Order |
|----------|---------------|----------------|
| Conflicting manifest values | Use highest-priority manifest | 1. Runtime, 2. Project, 3. Session |
| Manifest vs. inference | Always use manifest | Manifest |
| Manifest vs. user claim | Use manifest unless user provides explicit override | Manifest |
| User vs. AI claim | User claim takes precedence | User Input |
| Multiple user inputs | Use most recent timestamp | Timestamp |

### 11.6 Fallback Hierarchy

When an authoritative source is unavailable, use the following fallback order:

| Missing Source | Fallback Source | Condition |
|----------------|-----------------|-----------|
| Runtime Manifest | Error — runtime is mandatory | Non-negotiable |
| Project Manifest | Runtime Manifest defaults | Optional project config |
| Session State | Create new session | Session required |
| Git metadata | Skip Git-specific context | Git not required |
| User input | Use runtime defaults | Input optional |

### 11.7 Heuristic Boundaries

Heuristics (inference-based determinations) MAY be used ONLY when:

1. **No authoritative source exists** for the required information
2. **All authoritative sources have been exhausted**
3. **The heuristic is explicitly documented** and logged
4. **The heuristic is flagged as inferred** in the context

---

## 12. State Machine

### 12.1 State Definitions

| State | Code | Description |
|-------|------|-------------|
| **IDLE** | 0 | No engineering session exists; no active context |
| **CONNECTING** | 10 | Session establishment initiated |
| **DISCOVERING** | 20 | Runtime locating workspace |
| **AUTHENTICATING** | 30 | Verifying workspace access |
| **BOOTSTRAPPING** | 40 | Establishing engineering context |
| **VERIFYING** | 50 | Verifying context completeness |
| **ACTIVE** | 100 | Engineering session established; AI reasoning permitted |
| **SUSPENDED** | 110 | Session paused; context preserved |
| **RESUMING** | 120 | Resuming from suspended state |
| **TERMINATING** | 130 | Session ending gracefully |
| **TERMINATED** | 200 | Session ended; context archived |
| **FAILED** | 900 | Unrecoverable error; session cannot continue |

### 12.2 Legal Transitions

| From State | To State | Trigger | Precondition |
|------------|----------|---------|--------------|
| IDLE | CONNECTING | User initiates session | None |
| CONNECTING | DISCOVERING | CONNECT received | Valid client ID |
| CONNECTING | FAILED | Connection error | Invalid client |
| DISCOVERING | AUTHENTICATING | Workspace found | .kdse/ exists |
| DISCOVERING | FAILED | Workspace not found | .kdse/ missing |
| AUTHENTICATING | BOOTSTRAPPING | Access verified | Read/write granted |
| AUTHENTICATING | FAILED | Access denied | Insufficient permissions |
| BOOTSTRAPPING | VERIFYING | Bootstrap incomplete | Context partial |
| BOOTSTRAPPING | ACTIVE | Bootstrap complete | Context complete |
| BOOTSTRAPPING | FAILED | Bootstrap error | Unrecoverable error |
| VERIFYING | ACTIVE | Verification passed | All checks pass |
| VERIFYING | FAILED | Verification failed | Required checks fail |
| ACTIVE | SUSPENDED | User suspends | Active session |
| ACTIVE | TERMINATING | User terminates | Active session |
| SUSPENDED | RESUMING | User resumes | Suspended session |
| RESUMING | ACTIVE | Resume successful | Context restored |
| RESUMING | FAILED | Resume failed | Context corrupted |
| TERMINATING | TERMINATED | Termination complete | Graceful shutdown |
| FAILED | IDLE | Reset requested | Error cleared |

### 12.3 Terminal States

| State | Terminal? | Exit Condition |
|-------|-----------|---------------|
| TERMINATED | Yes | Session archived; context preserved |
| FAILED | Yes | Error cleared; new session required |

### 12.4 State Transition Diagram

```
                         ┌─────────────────────────────────────────┐
                         │                                         │
                         │              ┌─────────┐                │
    ┌────────┐           │              │  IDLE   │◀───┐          │
    │ Reset  │───────────▶│              └────┬────┘    │          │
    └────────┘           │                   │         │          │
                         │                   │ User    │          │
                         │                   │ Init    │          │
                         │                   ▼         │          │
                         │              ┌─────────┐    │          │
                         │              │CONNECTING│───┘          │
                         │              └────┬────┘              │
                         │                   │                  │
                         │                   ▼                  │
                         │              ┌─────────┐              │
                         │              │DISCOVERING│             │
                         │              └────┬────┘              │
                         │                   │                  │
                         │                   ▼                  │
                         │              ┌─────────┐              │
                         │              │AUTHENTICATING│          │
                         │              └────┬────┘              │
                         │                   │                  │
                         │                   ▼                  │
                         │              ┌─────────┐              │
                         │              │BOOTSTRAPPING│          │
                         │              └────┬────┘              │
                         │                   │                  │
                         │         ┌─────────┴─────────┐       │
                         │         │                   │       │
                         │         ▼                   ▼       │
                         │  ┌─────────────┐     ┌─────────┐   │
                         │  │  VERIFYING  │     │ ACTIVE  │   │
                         │  └──────┬──────┘     └────┬────┘   │
                         │         │               │         │
                         │         │               ├────┐    │
                         │         │               │    │    │
                         │         ▼               ▼    │    │
                         │  ┌─────────────┐  ┌──────────┐ │    │
                         │  │   FAILED    │  │SUSPENDED │ │    │
                         │  └─────────────┘  └────┬─────┘ │    │
                         │                       │       │    │
                         │                       ▼       │    │
                         │                  ┌─────────┐   │    │
                         │                  │RESUMING │───┘    │
                         │                  └─────────┘        │
                         │                                         │
                         │                       ▲               │
                         │                       │               │
                         │                  ┌─────────┐         │
                         │                  │TERMINATING│─────────┘
                         │                  └────┬────┘
                         │                       │
                         │                       ▼
                         │                  ┌─────────┐
                         │                  │TERMINATED│
                         │                  └─────────┘
                         │                                         │
                         └─────────────────────────────────────────┘
```

---

## 13. Failure Handling

### 13.1 Failure Categories

| Category | Code Range | Recoverable? |
|----------|------------|--------------|
| **Runtime Missing** | 1000-1099 | Yes (may initialize) |
| **Runtime Invalid** | 1100-1199 | Conditional |
| **Authentication Failed** | 1200-1299 | Yes (may retry) |
| **Context Incomplete** | 1300-1399 | Yes (may retry) |
| **Permission Denied** | 1400-1499 | Conditional |
| **Version Mismatch** | 1500-1599 | Conditional |
| **Corruption Detected** | 1600-1699 | Conditional |
| **Session Error** | 1700-1799 | Conditional |

### 13.2 Failure Codes

| Code | Name | Description | Owner | Recovery Action |
|------|------|-------------|-------|----------------|
| 1001 | RUNTIME_NOT_FOUND | .kdse/ directory missing | User | Initialize runtime |
| 1002 | RUNTIME_UNREADABLE | .kdse/ not accessible | OS | Fix permissions |
| 1101 | MANIFEST_INVALID | runtime.yaml malformed | Runtime | Restore from backup |
| 1102 | MANIFEST_MISSING | Required file absent | Runtime | Restore from backup |
| 1103 | PHASE_INVALID | Phase value invalid | Runtime | Set to valid phase |
| 1201 | AUTH_FAILED | Workspace access denied | User | Grant permissions |
| 1202 | AUTH_EXPIRED | Session token expired | Runtime | Re-authenticate |
| 1301 | CONTEXT_INCOMPLETE | Required fields missing | Runtime | Retry bootstrap |
| 1302 | CONTEXT_CORRUPT | Checksum mismatch | Runtime | Restore context |
| 1401 | PERMISSION_READ | Cannot read .kdse/ | OS | Fix permissions |
| 1402 | PERMISSION_WRITE | Cannot write .kdse/ | OS | Fix permissions |
| 1501 | VERSION_INCOMPATIBLE | Runtime version mismatch | Runtime | Upgrade/downgrade |
| 1502 | TEMPLATE_INCOMPATIBLE | Template version mismatch | Runtime | Update templates |
| 1601 | CHECKSUM_MISMATCH | Integrity check failed | Runtime | Verify/restore |
| 1602 | STATE_CORRUPT | State files damaged | Runtime | Restore from backup |
| 1701 | SESSION_EXPIRED | Session too old | Runtime | Start new session |
| 1702 | SESSION_CONFLICT | Concurrent session | Runtime | Resolve conflict |

### 13.3 Deterministic Recovery Procedures

#### 13.3.1 RUNTIME_NOT_FOUND (1001)

| Step | Action | Authority |
|------|--------|-----------|
| 1 | Verify workspace path is correct | Runtime |
| 2 | Check if .kdse/ exists at alternate locations | Runtime |
| 3 | If user requests initialization, proceed to Initialize | User |
| 4 | Initialize new runtime with template | Runtime |
| 5 | Verify initialized runtime | Runtime |
| 6 | Return to BOOTSTRAP | Runtime |

#### 13.3.2 MANIFEST_INVALID (1101)

| Step | Action | Authority |
|------|--------|-----------|
| 1 | Identify specific validation errors | Runtime |
| 2 | Check for backup manifest | Runtime |
| 3 | If backup exists and valid, restore | Runtime |
| 4 | If no backup, report corruption | Runtime |
| 5 | User may reinitialize | User |
| 6 | Return to BOOTSTRAP | Runtime |

#### 13.3.3 CONTEXT_INCOMPLETE (1301)

| Step | Action | Authority |
|------|--------|-----------|
| 1 | Identify missing required fields | Runtime |
| 2 | For each missing field, apply fallback rules | Runtime |
| 3 | If still incomplete, report missing sources | Runtime |
| 4 | User provides required information | User |
| 5 | Reassemble context | Runtime |
| 6 | Return to BOOTSTRAP | Runtime |

#### 13.3.4 VERSION_INCOMPATIBLE (1501)

| Step | Action | Authority |
|------|--------|-----------|
| 1 | Determine required runtime version | Runtime Manifest |
| 2 | Determine current runtime version | Runtime |
| 3 | If downgrade possible, offer compatibility mode | Runtime |
| 4 | If upgrade required, report to user | Runtime |
| 5 | User decides: upgrade, continue, or abort | User |
| 6 | If continuing, proceed with compatibility mode | Runtime |
| 7 | Return to BOOTSTRAP | Runtime |

### 13.4 Failure Response Format

All failures MUST return the following structured response:

```json
{
  "failure": {
    "code": "<integer>",
    "name": "<string>",
    "description": "<string>",
    "details": {
      "diagnostic": "<string>",
      "source": "<string>",
      "expected": "<string>",
      "actual": "<string>"
    },
    "remediation": {
      "action": "<string>",
      "owner": "<string>",
      "steps": ["<string>"]
    },
    "timestamp": "<iso8601>",
    "request_id": "<uuid>"
  }
}
```

---

## 14. Completion Criteria

### 14.1 Session ACTIVE Conditions

The Engineering Session reaches **ACTIVE** state ONLY when ALL of the following conditions are met:

| Condition | Verification Method | Evidence Required |
|-----------|--------------------|--------------------|
| Runtime manifest exists | Filesystem check | .kdse/runtime.yaml present |
| Runtime manifest valid | YAML parse success | No parse errors |
| Workspace path valid | Directory accessible | Directory exists |
| Session state created | Session ID generated | .kdse/session.yaml created |
| Phase state valid | Phase from approved list | Phase value verified |
| Context checksum computed | Hash computation | Checksum recorded |
| Context completeness verified | All required fields present | Completion report |
| Bootstrap complete | All bootstrap steps finished | Bootstrap log |
| Verification passed (if required) | All verification checks pass | Verification report |

### 14.2 Engineering Authorization

Upon reaching ACTIVE state, the following authorizations are granted:

| Authorization | Scope | Constraint |
|---------------|-------|------------|
| Read engineering context | Full context access | Read-only |
| Recommend engineering actions | Within session scope | User approval required |
| Generate reports | Session reports only | Per report specification |
| Read project artifacts | Per permissions | No modification |
| Access external references | Referenced sources only | Per reference permissions |

### 14.3 AI Reasoning Authorization

AI reasoning MAY proceed ONLY after:

1. Session has reached **ACTIVE** state
2. Engineering context has been delivered to AI agent
3. AI agent has confirmed context receipt
4. Session token is valid and non-expired

### 14.4 Active Session Requirements

An ACTIVE session MUST maintain:

| Requirement | Verification |
|-------------|--------------|
| Session token valid | Token not expired |
| Context integrity | Checksum matches |
| Phase consistency | Phase unchanged |
| Workspace accessible | Directory accessible |
| Runtime operational | Runtime responsive |

---

## 15. Compliance Requirements

### 15.1 Conformance Declaration

A KDSE Runtime implementation MAY declare compliance with this protocol by satisfying ALL mandatory requirements defined in this section.

### 15.2 Mandatory Requirements

| ID | Requirement | Verification |
|----|-------------|--------------|
| C-001 | Protocol MUST implement all handshake messages defined in Section 8 | Message conformance test |
| C-002 | Protocol MUST execute bootstrap steps defined in Section 9 | Bootstrap compliance test |
| C-003 | Protocol MUST enforce authority hierarchy defined in Section 7 | Authority resolution test |
| C-004 | Protocol MUST implement all state transitions defined in Section 12 | State machine test |
| C-005 | Protocol MUST handle all failure codes defined in Section 13 | Failure handling test |
| C-006 | Protocol MUST enforce completion criteria defined in Section 14 | Completion criteria test |
| C-007 | Protocol MUST NOT begin AI reasoning before ACTIVE state | State enforcement test |
| C-008 | Protocol MUST produce identical context for identical inputs | Determinism test |
| C-009 | Protocol MUST use authoritative sources over heuristics | Authority precedence test |
| C-010 | Protocol MUST NOT depend on implementation-specific behavior | Independence test |

### 15.3 Optional Features

The following features enhance but are not required for compliance:

| Feature | Description | Compliance Impact |
|---------|-------------|------------------|
| Session resumption | Ability to resume suspended sessions | Optional |
| Parallel sessions | Multiple simultaneous sessions | Optional |
| Remote runtime | Runtime on remote host | Optional |
| Offline mode | Operation without network | Optional |

### 15.4 Compliance Testing

#### 15.4.1 Conformance Test Suite

| Test ID | Test Name | Purpose |
|---------|-----------|---------|
| T-001 | Handshake Complete | Verify all handshake messages implemented |
| T-002 | Bootstrap Deterministic | Verify identical context for identical inputs |
| T-003 | Authority Resolution | Verify correct authority precedence |
| T-004 | State Transitions | Verify all legal transitions work |
| T-005 | Failure Handling | Verify all failure codes handled |
| T-006 | Completion Criteria | Verify ACTIVE state requirements |
| T-007 | Isolation Test | Verify bootstrap and reasoning separation |

#### 15.4.2 Compliance Declaration Format

```markdown
# KDSE Engineering Session Protocol Compliance Declaration

**Implementation Name:** <name>
**Implementation Version:** <version>
**Protocol Version:** 1.0
**Declaration Date:** <date>

## Conformance Level

[ ] FULL CONFORMANCE
[ ] PARTIAL CONFORMANCE

## Mandatory Requirements

| Requirement | Status | Evidence |
|-------------|--------|----------|
| C-001 | [PASS/FAIL] | <test results> |
| C-002 | [PASS/FAIL] | <test results> |
| ... | ... | ... |

## Known Deviations

<description of any deviations from mandatory requirements>

## Testing Evidence

<test results, logs, or other evidence>

## Attestation

We attest that this implementation conforms to the KDSE Engineering Session Protocol 
specification version 1.0 as declared above.

Signed: <signature>
Date: <date>
```

---

## Appendix A: Normative References

| Reference | Title | Authority |
|-----------|-------|-----------|
| KDSE Standard | KDSE Engineering Standard | Highest |
| ADR-001 | Runtime is the Authority | Normative |
| ADR-002 | One Methodology, Multiple Runtimes | Normative |
| OWNERSHIP_MODEL | KDSE Ownership Model | Normative |
| PRINCIPLES | KDSE Architectural Principles | Normative |

## Appendix B: Informative References

| Reference | Title | Purpose |
|-----------|-------|---------|
| TCP Protocol | RFC 793 - TCP | Handshake inspiration |
| Modbus | Modbus Application Protocol | Protocol structure |
| OPC UA | OPC Unified Architecture | State machine inspiration |

---

*This document is NORMATIVE. All compliant KDSE Runtime implementations MUST conform to this specification.*
