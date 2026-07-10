# Glossary

## Purpose

This glossary establishes canonical terminology for KDSE. Each term has exactly one definition within the KDSE context. When the same term has different meanings in different contexts, this glossary establishes the authoritative meaning.

## A

### Artifact

A tangible or intangible object created, maintained, or used during the engineering process.

In KDSE, artifacts are typed (Knowledge, Architecture, Implementation, Verification, Governance) with defined purposes, owners, lifetimes, and authority levels.

### Authority

The legitimate power to make decisions and the legitimacy to require compliance.

In KDSE, authority flows downward through the artifact hierarchy. Higher-authority artifacts constrain and authorize lower-authority artifacts.

## B

### Body of Knowledge (BoK)

The complete set of concepts, terms, and activities that constitute a discipline.

The KDSE Body of Knowledge encompasses all artifacts, principles, processes, and practices that define the discipline.

## C

### Change Management

The process by which changes to artifacts are proposed, evaluated, approved, implemented, and documented.

In KDSE, changes must respect the authority hierarchy. Changes to lower-layer artifacts do not automatically propagate, but changes to higher-layer artifacts may require corresponding lower-layer changes.

### Context

The circumstances, environment, and conditions within which a system operates.

In KDSE, context is captured as knowledge and informs architectural decisions.

## D

### Derivation

The process by which lower-layer artifacts are created based on higher-layer artifacts.

In KDSE, derivation is mandatory: Architecture derives from Knowledge; Implementation derives from Architecture.

### Discipline

A field of study or professional practice governed by established principles and practices.

KDSE is a discipline within software engineering.

## E

### Evolution

The process by which artifacts change over time while maintaining traceability and authority alignment.

Evolution is the fifth stage of the KDSE lifecycle and is continuous throughout the system lifecycle.

## G

### Governance

The system of rules, practices, and processes by which authority is exercised and decisions are made.

In KDSE, Governance artifacts establish ownership, authority delegation, and process compliance requirements.

## I

### Implementation

The physical realization of a system through code, configuration, and related artifacts.

Implementation derives authority from Architecture and must not contradict Architecture or Knowledge.

### Intent

The purpose or goal behind a decision, requirement, or artifact.

In KDSE, understanding intent is essential for maintaining alignment when deriving artifacts.

## K

### Knowledge

The authoritative understanding about a problem domain, including requirements, constraints, context, and assumptions.

Knowledge is the highest-authority artifact type in KDSE. All other artifacts derive authority from Knowledge.

### Knowledge Artifact

An artifact of type Knowledge, capturing authoritative understanding.

Knowledge artifacts are created during the Knowledge stage and persist throughout the system lifecycle.

## M

### Methodology

A systematic approach to a discipline, defining principles, processes, and practices.

KDSE is a methodology for software engineering.

## N

### Non-Functional Requirement

A requirement that specifies criteria for system operation rather than specific behavior.

Examples include performance, security, reliability, and scalability requirements. Non-functional requirements are captured as Knowledge artifacts.

## P

### Pattern

A reusable solution to a recurring problem within a specific context.

Patterns may inform architecture and implementation but are not authoritative in KDSE. Authority derives from Knowledge artifacts.

### Practice

An activity or technique performed as part of engineering work.

Practices in KDSE are derived from principles and may vary by context. Practices are not authoritative; principles are.

### Principle

A fundamental truth or proposition that guides decision-making and action.

KDSE defines core principles that are timeless and context-independent.

### Problem Domain

The area of expertise or application to which a system addresses.

Understanding the problem domain is captured as Knowledge.

## R

### Requirement

A statement of need or constraint that a system must satisfy.

Requirements are captured as Knowledge artifacts in KDSE.

### Repository

A storage location for artifacts, typically version-controlled.

In KDSE, a repository may contain artifacts of any type. The repository maintains traceability relationships between artifacts.

### Resolution

The process of addressing a contradiction or gap between artifacts.

Resolution may involve clarification at a higher layer or correction at a lower layer.

## T

### Traceability

The ability to follow the relationship between artifacts across the hierarchy.

In KDSE, every artifact below Knowledge must be traceable to authorized Knowledge artifacts. Traceability enables impact analysis and compliance verification.

### Traceability Path

The sequence of artifacts connecting two related artifacts.

For example, an Implementation artifact traces to an Architecture artifact through an explicit relationship. The Architecture artifact traces to a Knowledge artifact.

## V

### Verification

The process of confirming that implementation aligns with architecture and that architecture aligns with knowledge.

Verification artifacts document the results of verification activities. Verification derives authority from Knowledge artifacts.

### Verification Artifact

An artifact of type Verification, documenting verification activities and results.

Verification artifacts include verification plans, test cases, test results, and reports.

## W

### Working Memory

The temporary state used during reasoning and derivation.

Working memory in KDSE refers to intermediate artifacts and reasoning processes used during derivation. Working memory is distinct from persistent artifacts.
