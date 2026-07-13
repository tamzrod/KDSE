# Knowledge-Driven Software Engineering (KDSE)

**KDSE** is an engineering methodology in which structured knowledge serves as the authoritative source from which all other software artifacts are derived, maintained, and verified throughout the software lifecycle.

## Quick Reference

| Document | Purpose |
|----------|---------|
| [Foundation](docs/foundation/) | The authoritative definition of KDSE |
| [Evolution](docs/evolution/) | Evidence-driven methodology improvements |
| [Audit](docs/audit/) | KDSE Audit System standards |
| [Runtime](runtime/) | Official reference implementation |
| [000-what-is-kdse.md](docs/foundation/000-what-is-kdse.md) | What is KDSE? |
| [001-why-kdse-exists.md](docs/foundation/001-why-kdse-exists.md) | Why KDSE exists |
| [002-scope.md](docs/foundation/002-scope.md) | What KDSE is and is not |
| [003-core-principles.md](docs/foundation/003-core-principles.md) | Core principles |
| [004-engineering-model.md](docs/foundation/004-engineering-model.md) | Engineering lifecycle |
| [005-engineering-artifacts.md](docs/foundation/005-engineering-artifacts.md) | Artifact definitions and stewardship |
| [006-chain-of-authority.md](docs/foundation/006-chain-of-authority.md) | Authority hierarchy |
| [007-glossary.md](docs/foundation/007-glossary.md) | Terminology |
| [008-future-vision.md](docs/foundation/008-future-vision.md) | Long-term vision and maturity model |
| [009-engineering-knowledge.md](docs/foundation/009-engineering-knowledge.md) | Engineering knowledge definition |
| [010-knowledge-derivation.md](docs/foundation/010-knowledge-derivation.md) | Derivation mechanics |
| [011-adoption-model.md](docs/foundation/011-adoption-model.md) | Adoption path |
| [012-traceability.md](docs/foundation/012-traceability.md) | Traceability framework |
| [013-authority-resolution.md](docs/foundation/013-authority-resolution.md) | Authority and conflict resolution |
| [014-engineering-review-process.md](docs/foundation/014-engineering-review-process.md) | Methodology review process |
| [015-reference-artifacts.md](docs/foundation/015-reference-artifacts.md) | Reference Artifacts as engineering evidence |
| [016-reference-analysis-knowledge-derivation.md](docs/foundation/016-reference-analysis-knowledge-derivation.md) | Knowledge Derivation Lifecycle |
| [017-engineering-knowledge-definition.md](docs/foundation/017-engineering-knowledge-definition.md) | Domain Knowledge definition |
| [018-architecture-phase.md](docs/foundation/018-architecture-phase.md) | Architecture as distinct phase |
| [019-implementation-phase.md](docs/foundation/019-implementation-phase.md) | Implementation as distinct phase |
| [020-domain-interfaces.md](docs/foundation/020-domain-interfaces.md) | Domain Interfaces |
| [021-evidence-and-strength.md](docs/foundation/021-evidence-and-strength.md) | Evidence and Evidence Strength |
| [022-collector-philosophy.md](docs/foundation/022-collector-philosophy.md) | Collector philosophy |
| [023-question-classification.md](docs/foundation/023-question-classification.md) | Question classification |
| [024-engineering-independence-test.md](docs/foundation/024-engineering-independence-test.md) | Engineering Independence Test |

## Core Principles

1. Knowledge precedes architecture
2. Architecture precedes implementation
3. Implementation precedes verification
4. Knowledge is the longest-lived artifact
5. Engineering decisions must be traceable
6. Code realizes knowledge
7. Knowledge is language-independent
8. Authority flows downward
9. Verification confirms alignment
10. Evolution maintains authority
11. Reference Artifacts support Domain Knowledge (not replace)
12. Domain Knowledge is implementation-independent
13. Evidence Strength strengthens but does not authorize
14. Repository First: analyze artifacts before asking operator
15. Contradictions are preserved, never silently resolved

## Engineering Lifecycle

```
Reference Artifacts → Reference Analysis → Domain Knowledge Derivation → Architecture → Implementation → Verification → Evolution
```

## Key Concepts

- **Reference Artifacts**: Existing sources of domain information (evidence, not authority)
- **Domain Knowledge**: Implementation-independent understanding that remains valid if implementation is rewritten
- **Domain Interfaces**: Domain responsibilities that exclude implementation technologies
- **Evidence Strength**: Domain support measure (★★★★★ to ★☆☆☆☆) replacing AI confidence
- **Engineering Independence Test**: Validation ensuring knowledge remains valid across technology changes
- **Collector**: Methodology component defined by responsibility, not artifact type
- **Question Classification**: Routing unresolved items to correct phase (Knowledge/Architecture/Implementation)
- **Repository First Principle**: Analyze all artifacts before asking operator
- **Artifact Lifecycle**: Progression through defined states (Draft, Reviewed, Approved, etc.)
- **Stewardship**: Responsibility without dominion; knowledge is stewarded, not owned
- **Traceability**: The ability to follow relationships between artifacts
- **Separation of Concerns**: Reference Artifacts, Domain Knowledge, Architecture, and Implementation remain independent

## Separation of Concerns

KDSE maintains strict separation between:

| Concern | Description |
|---------|-------------|
| Reference Artifact | Domain evidence (project docs, implementation, vendor docs) |
| Domain Knowledge | Implementation-independent understanding |
| Architecture | Organization of Domain Knowledge into software |
| Implementation | Realization of Architecture using specific technologies |

## License

Apache 2.0