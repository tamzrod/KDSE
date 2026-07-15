"""
KDSE Slim - Foundation Skeleton Generator
Immediately creates Foundation artifacts after receiving an objective.
Foundation is the living artifact collection that guides all engineering.
"""

from dataclasses import dataclass, field
from typing import List, Optional, Dict
from datetime import datetime
import os


@dataclass
class FoundationArtifact:
    """A Foundation artifact."""
    name: str
    path: str
    content: str = ""
    created: bool = False
    category: str = ""  # problem, spec, requirement, assumption, constraint


class FoundationSkeleton:
    """
    Creates the Foundation Skeleton immediately after receiving an objective.
    
    Foundation includes:
    - PROBLEM.md: What problem are we solving?
    - SPEC.md: What is the project specification?
    - REQUIREMENTS.md: What are the functional requirements?
    - ASSUMPTIONS.md: What are the reasonable assumptions?
    - CONSTRAINTS.md: What are the constraints?
    - GLOSSARY.md: What domain terminology is used?
    """
    
    def __init__(self, workspace_root: str = ".kdse"):
        self.workspace_root = workspace_root
        self.knowledge_dir = os.path.join(workspace_root, "knowledge")
        self.artifacts: List[FoundationArtifact] = []
    
    def create_skeleton(self, objective: str) -> Dict:
        """
        Create the Foundation Skeleton immediately after receiving objective.
        
        Args:
            objective: The user's objective
            
        Returns:
            Report of what was created
        """
        created = []
        
        # Create directories
        os.makedirs(self.knowledge_dir, exist_ok=True)
        
        # 1. PROBLEM.md - What problem are we solving?
        problem_path = os.path.join(self.knowledge_dir, "PROBLEM.md")
        problem_content = self._generate_problem(objective)
        self._write_artifact(problem_path, problem_content)
        created.append({"name": "PROBLEM.md", "path": problem_path})
        
        # 2. SPEC.md - What is the project specification?
        spec_path = os.path.join(self.knowledge_dir, "SPEC.md")
        spec_content = self._generate_spec(objective)
        self._write_artifact(spec_path, spec_content)
        created.append({"name": "SPEC.md", "path": spec_path})
        
        # 3. REQUIREMENTS.md - What are the functional requirements?
        req_path = os.path.join(self.knowledge_dir, "REQUIREMENTS.md")
        req_content = self._generate_requirements(objective)
        self._write_artifact(req_path, req_content)
        created.append({"name": "REQUIREMENTS.md", "path": req_path})
        
        # 4. ASSUMPTIONS.md - What are the reasonable assumptions?
        assumptions_path = os.path.join(self.knowledge_dir, "ASSUMPTIONS.md")
        assumptions_content = self._generate_assumptions(objective)
        self._write_artifact(assumptions_path, assumptions_content)
        created.append({"name": "ASSUMPTIONS.md", "path": assumptions_path})
        
        # 5. CONSTRAINTS.md - What are the constraints?
        constraints_path = os.path.join(self.knowledge_dir, "CONSTRAINTS.md")
        constraints_content = self._generate_constraints(objective)
        self._write_artifact(constraints_path, constraints_content)
        created.append({"name": "CONSTRAINTS.md", "path": constraints_path})
        
        # 6. GLOSSARY.md - What domain terminology is used?
        glossary_path = os.path.join(self.knowledge_dir, "GLOSSARY.md")
        glossary_content = self._generate_glossary(objective)
        self._write_artifact(glossary_path, glossary_content)
        created.append({"name": "GLOSSARY.md", "path": glossary_path})
        
        return {
            "objective": objective,
            "created_at": datetime.now().isoformat(),
            "artifacts_created": created,
            "knowledge_dir": self.knowledge_dir,
            "status": "foundation_skeleton_created"
        }
    
    def _write_artifact(self, path: str, content: str):
        """Write artifact content to file."""
        with open(path, 'w') as f:
            f.write(content)
    
    def _generate_problem(self, objective: str) -> str:
        """Generate PROBLEM.md content."""
        return f"""# PROBLEM.md
## Problem Statement

**Objective:** {objective}

### Problem Definition

[Describe the problem being solved. What user need exists?]

### Scope

**In Scope:**
- [What will be addressed]

**Out of Scope:**
- [What will NOT be addressed]

### Stakeholders

- [Who has interest in this problem being solved?]

### Success Criteria

- [What constitutes successful problem resolution?]

---

*This document defines the problem scope. It is the foundation of all engineering work.*
*Created: {datetime.now().isoformat()}*
*State: DRAFT*
*Evidence Strength: ● (Initial assessment, may change)*
"""
    
    def _generate_spec(self, objective: str) -> str:
        """Generate SPEC.md content."""
        return f"""# SPEC.md
## Project Specification

**Objective:** {objective}

### Overview

[High-level description of what will be built.]

### Goals

1. [Primary goal]
2. [Secondary goals]

### Non-Goals

- [What this project will NOT do]

### Background

[Why does this problem exist? What context is relevant?]

### Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk 1] | [Impact] | [Mitigation] |

---

*This document provides project overview. It guides but does not prescribe implementation.*
*Created: {datetime.now().isoformat()}*
*State: DRAFT*
*Evidence Strength: ● (Initial specification, may change)*
"""
    
    def _generate_requirements(self, objective: str) -> str:
        """Generate REQUIREMENTS.md content."""
        return f"""# REQUIREMENTS.md
## Functional Requirements

**Objective:** {objective}

### Known Requirements

#### Must Have

1. **[REQ-001]** [Requirement statement]
   - *Evidence:* [What evidence supports this?]
   - *Status:* [Known/Assumed/Gap]

2. **[REQ-002]** [Requirement statement]
   - *Evidence:* [What evidence supports this?]
   - *Status:* [Known/Assumed/Gap]

#### Should Have

1. **[REQ-101]** [Requirement statement]
   - *Evidence:* [What evidence supports this?]
   - *Status:* [Known/Assumed/Gap]

### Unknown Requirements

[List requirements that are not yet known and need investigation.]

### Requirement Gaps

| Gap | Priority | Investigation Needed |
|-----|----------|---------------------|
| [Gap 1] | [High/Med/Low] | [What information is needed?] |

---

*This document lists functional requirements. Requirements with "Assumed" status must be validated.*
*Created: {datetime.now().isoformat()}*
*State: DRAFT*
"""
    
    def _generate_assumptions(self, objective: str) -> str:
        """Generate ASSUMPTIONS.md content."""
        return f"""# ASSUMPTIONS.md
## Documented Assumptions

**Objective:** {objective}

### Critical Assumptions

> ⚠️ **IMPORTANT:** Assumptions are NOT facts. They must be validated before becoming requirements.

1. **[ASSUME-001]** [Assumption statement]
   - *Rationale:* [Why is this assumed?]
   - *Validation:* [How will this be validated?]
   - *Status:* [Assumed/Validated/Invalidated]
   - *Evidence:* [What supports this assumption?]

2. **[ASSUME-002]** [Assumption statement]
   - *Rationale:* [Why is this assumed?]
   - *Validation:* [How will this be validated?]
   - *Status:* [Assumed/Validated/Invalidated]
   - *Evidence:* [What supports this assumption?]

### Assumptions Awaiting Validation

| Assumption | Priority | Validation Status |
|------------|----------|-------------------|
| [ASSUME-001] | High | Not validated |

### Invalidated Assumptions

[None yet documented]

---

*This document tracks assumptions. Nothing silently becomes a fact.*
*Created: {datetime.now().isoformat()}*
*State: DRAFT*
"""
    
    def _generate_constraints(self, objective: str) -> str:
        """Generate CONSTRAINTS.md content."""
        return f"""# CONSTRAINTS.md
## Project Constraints

**Objective:** {objective}

### Technical Constraints

| Constraint | Source | Justification |
|------------|--------|---------------|
| [Tech Constraint 1] | [Source] | [Why is this a constraint?] |

### Business Constraints

| Constraint | Source | Justification |
|------------|--------|---------------|
| [Business Constraint 1] | [Source] | [Why is this a constraint?] |

### Regulatory Constraints

| Constraint | Source | Justification |
|------------|--------|---------------|
| [Regulatory Constraint 1] | [Source] | [Why is this a constraint?] |

### Known Flexibility

- [Where can constraints be relaxed?]

---

*This document defines what limits the solution space. Constraints must be respected.*
*Created: {datetime.now().isoformat()}*
*State: DRAFT*
"""
    
    def _generate_glossary(self, objective: str) -> str:
        """Generate GLOSSARY.md content."""
        return f"""# GLOSSARY.md
## Domain Terminology

**Objective:** {objective}

### Terms

| Term | Definition | Source |
|------|------------|--------|
| [Term 1] | [Definition] | [Source] |
| [Term 2] | [Definition] | [Source] |

### Acronyms

| Acronym | Expansion |
|---------|-----------|
| [ACR-1] | [Full form] |

### Ambiguous Terms

[List terms that have multiple meanings and how they are used in this context.]

---

*This document standardizes domain terminology. Ambiguous terms are clarified.*
*Created: {datetime.now().isoformat()}*
*State: DRAFT*
"""


class KnowledgeDistinction:
    """
    Ensures proper distinction between:
    - Known Facts (backed by evidence)
    - Reasonable Assumptions (documented, to be validated)
    - Knowledge Gaps (unknown, need investigation)
    - Evidence (sources that support knowledge)
    """
    
    @staticmethod
    def is_fact(evidence_count: int, validated: bool) -> bool:
        """Determine if something should be treated as a fact."""
        return evidence_count >= 2 and validated
    
    @staticmethod
    def is_assumption(evidence_count: int, validated: bool) -> bool:
        """Determine if something should be treated as an assumption."""
        return evidence_count > 0 and not validated
    
    @staticmethod
    def is_gap(evidence_count: int) -> bool:
        """Determine if something is a knowledge gap."""
        return evidence_count == 0
    
    @staticmethod
    def classify_statement(evidence_count: int, validated: bool) -> str:
        """Classify a statement as fact, assumption, or gap."""
        if evidence_count == 0:
            return "GAP"
        elif evidence_count >= 2 and validated:
            return "FACT"
        elif evidence_count >= 1 and validated:
            return "STRONG_ASSUMPTION"
        else:
            return "ASSUMPTION"
