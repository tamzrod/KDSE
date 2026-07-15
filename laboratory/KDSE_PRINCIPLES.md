# KDSE Slim Principles
## Constitutional Principles (Non-Negotiable)

These principles cannot be modified, overridden, or bypassed by any implementation.

---

## 1. Knowledge Precedes Architecture

**Principle:** Architecture derives from knowledge. Every architectural decision must trace to specific knowledge artifacts.

**Implication:** 
- No architecture without derived knowledge
- Every decision cites its knowledge basis
- Architecture without knowledge citation is not KDSE-compliant

---

## 2. Authority Flows Downward

**Principle:** Lower layers cannot contradict higher layers. Implementation cannot contradict Architecture. Architecture cannot contradict Knowledge.

**Implication:**
- Knowledge is highest authority
- Architecture derives authority from Knowledge
- Implementation derives authority from Architecture
- Violations must be flagged

---

## 3. Evidence Supports, Never Authorizes

**Principle:** Evidence strengthens knowledge but does not create authority. Authority derives from derivation process, not evidence quantity.

**Implication:**
- Evidence Strength (●●●/●●/●) indicates trust level
- Strong evidence ≠ authoritative
- Derivation process creates authority
- Evidence enables judgment, not determination

---

## 4. Traceability Is Absolute

**Principle:** Every engineering decision must trace to authorized knowledge artifacts. Decisions without traceability are not KDSE-compliant.

**Implication:**
- Architecture traces to Knowledge
- Implementation traces to Architecture
- All artifacts must cite their authority basis
- Gaps in traceability must be identified

---

## 5. Knowledge Is Implementation-Independent

**Principle:** Knowledge remains valid if implementation is completely rewritten. Knowledge describes purpose, behavior, and constraints—not technology.

**Implication:**
- Knowledge statements survive technology changes
- Technology choices are Architecture, not Knowledge
- Language-agnostic statements are Knowledge
- Platform-specific statements are Architecture

---

## 6. Verification Is Continuous

**Principle:** Alignment between Knowledge → Architecture → Implementation must be continuously verified.

**Implication:**
- Verification is not a phase
- Checks run continuously
- Alignment must be confirmed at all times
- Violations must be reported immediately

---

## 7. Derivation Requires Reasoning

**Principle:** Knowledge cannot be extracted; it must be derived through analysis and judgment.

**Implication:**
- Collection is not derivation
- GATHER → DERIVE → VALIDATE process is mandatory
- Evidence must be interpreted, not just listed
- Derivation must be documented

---

## 8. Artifacts Are Living

**Principle:** Engineering artifacts evolve continuously. They are not created once and forgotten.

**Implication:**
- Artifacts have lifecycle states
- Artifacts are updated as understanding grows
- Archived ≠ deleted
- History is preserved

---

## Slim Architecture Rules

### Artifact Organization

Artifacts are organized by type:
- **knowledge/** - Highest authority
- **architecture/** - Traces to knowledge
- **implementation/** - Traces to architecture
- **verification/** - Confirms alignment

### Lifecycle States

All artifacts follow lifecycle:
- **DRAFT** - Initial creation (no authority)
- **REVIEWED** - Peer reviewed (partial authority)
- **APPROVED** - Authoritative (full authority)
- **ARCHIVED** - Superseded (historical)

### Evidence Strength

All knowledge has evidence strength:
- **●●● (STRONG)** - Multiple independent sources confirm
- **●● (MODERATE)** - Some corroboration exists
- **● (WEAK)** - Single source or inference only

### Progress Markers

Progress indicates advancement:
- **PROBLEM** - Problem defined
- **KNOWLEDGE** - Knowledge derived
- **FOUNDATION** - Foundation established
- **ARCHITECTURE** - Architecture decisions made
- **IMPLEMENTATION** - Implementation realized

### Decision Engine

Runtime decisions follow:
1. **STATE** - What exists? What traces to what?
2. **GAPS** - What is missing? What contradicts?
3. **DECIDE** - What should happen next?

### Work Orders

Work Orders provide guidance, never blocking:
- Current objective
- What should trace to what
- Expected deliverables
- What to check

---

## What Is NOT KDSE

The following are explicitly excluded:

- **Phases as authority gates** - Phases are progress markers
- **Confidence as gating** - Evidence Strength is separate
- **STRICT mode** - Blocking ≠ enforcement
- **Phase transitions** - Evidence drives decisions
- **Audit as phase** - Verification is continuous
- **Collect as tool** - Evidence gathering is part of derivation

---

*These principles define KDSE Slim. Every implementation must comply.*
