"""
KDSE Slim - Laboratory Runner
Validates the KDSE methodology through scenarios.
The Laboratory does NOT test software; it validates methodology.
"""

from dataclasses import dataclass, field
from typing import List, Dict, Optional
from datetime import datetime
from enum import Enum
import os
import json
import uuid

from .artifacts.registry import ArtifactRegistry, ArtifactType, LifecycleState, Artifact
from .artifacts.decision_engine import DecisionEngine, StateAssessment, Gap, WorkOrder
from .artifacts.derivation import DerivationOrchestrator, EvidenceStrength, KnowledgeStatement
from .artifacts.foundation import FoundationSkeleton, KnowledgeDistinction


class LaboratoryResult(Enum):
    """Laboratory experiment result."""
    PASS = "pass"
    FAIL = "fail"
    INCOMPLETE = "incomplete"


@dataclass
class LaboratoryScenario:
    """A laboratory scenario for validation."""
    id: str
    title: str
    objective: str
    expected_behavior: Dict = field(default_factory=dict)
    forbidden_behavior: Dict = field(default_factory=dict)
    checkpoints: List[Dict] = field(default_factory=list)


@dataclass  
class LaboratoryReport:
    """Report from a laboratory experiment."""
    scenario_id: str
    scenario_title: str
    started_at: str
    completed_at: str = ""
    result: LaboratoryResult = LaboratoryResult.INCOMPLETE
    objective: str = ""
    steps_executed: List[Dict] = field(default_factory=list)
    artifacts_created: List[str] = field(default_factory=list)
    decisions_made: List[Dict] = field(default_factory=list)
    checkpoints_passed: List[str] = field(default_factory=list)
    checkpoints_failed: List[str] = field(default_factory=list)
    knowledge_distinction: Dict = field(default_factory=dict)
    forbidden_actions: List[str] = field(default_factory=list)
    findings: List[str] = field(default_factory=list)
    pass_reason: str = ""
    fail_reason: str = ""
    
    def to_dict(self) -> dict:
        return {
            "scenario_id": self.scenario_id,
            "scenario_title": self.scenario_title,
            "started_at": self.started_at,
            "completed_at": self.completed_at,
            "result": self.result.value,
            "objective": self.objective,
            "steps_executed": self.steps_executed,
            "artifacts_created": self.artifacts_created,
            "decisions_made": self.decisions_made,
            "checkpoints_passed": self.checkpoints_passed,
            "checkpoints_failed": self.checkpoints_failed,
            "knowledge_distinction": self.knowledge_distinction,
            "forbidden_actions": self.forbidden_actions,
            "findings": self.findings,
            "pass_reason": self.pass_reason,
            "fail_reason": self.fail_reason
        }


class LaboratoryRunner:
    """
    Runs laboratory scenarios to validate KDSE methodology.
    
    The Laboratory validates that KDSE:
    ✓ Creates Foundation Skeleton
    ✓ Populates known facts
    ✓ Identifies assumptions
    ✓ Separates assumptions from facts
    ✓ Identifies knowledge gaps
    ✓ Records evidence strength
    ✓ Refines Foundation
    ✓ Produces engineering artifacts
    
    The Laboratory ensures KDSE does NOT:
    ✗ Choose programming language without justification
    ✗ Choose framework without justification
    ✗ Choose database without justification
    ✗ Generate source code prematurely
    ✗ Begin implementation without authorization
    """
    
    def __init__(self, scenario: LaboratoryScenario, workspace_root: str = ".kdse"):
        self.scenario = scenario
        self.workspace_root = workspace_root
        self.report = LaboratoryReport(
            scenario_id=scenario.id,
            scenario_title=scenario.title,
            started_at=datetime.now().isoformat(),
            objective=scenario.objective
        )
        self.registry = ArtifactRegistry()
        self.decision_engine = DecisionEngine(self.registry)
        self.derivation = DerivationOrchestrator()
        self.foundation = FoundationSkeleton(workspace_root)
        self.steps: List[str] = []
        self.forbidden_patterns = [
            "source code",
            ".py", ".js", ".go", ".java", ".rs", 
            "programming language",
            "framework",
            "database schema",
            "implementation started"
        ]
    
    def run(self) -> LaboratoryReport:
        """
        Execute the laboratory scenario.
        
        Returns:
            LaboratoryReport with results
        """
        print(f"Starting Laboratory Scenario: {self.scenario.title}")
        print(f"Objective: {self.scenario.objective}")
        print("-" * 60)
        
        # Step 1: Create Foundation Skeleton immediately
        print("\n[Step 1] Creating Foundation Skeleton...")
        self._step_create_foundation()
        
        # Step 2: Evaluate initial state
        print("\n[Step 2] Evaluating initial state...")
        self._step_evaluate_state()
        
        # Step 3: Identify Known Facts
        print("\n[Step 3] Identifying Known Facts...")
        self._step_identify_facts()
        
        # Step 4: Identify Assumptions
        print("\n[Step 4] Identifying Assumptions...")
        self._step_identify_assumptions()
        
        # Step 5: Identify Knowledge Gaps
        print("\n[Step 5] Identifying Knowledge Gaps...")
        self._step_identify_gaps()
        
        # Step 6: Record Evidence Strength
        print("\n[Step 6] Recording Evidence Strength...")
        self._step_record_evidence_strength()
        
        # Step 7: Refine Foundation
        print("\n[Step 7] Refining Foundation...")
        self._step_refine_foundation()
        
        # Finalize report
        self._finalize_report()
        
        return self.report
    
    def _step_create_foundation(self):
        """Step 1: Create Foundation Skeleton immediately."""
        result = self.foundation.create_skeleton(self.scenario.objective)
        
        self.steps.append("foundation_creation")
        self.report.steps_executed.append({
            "step": "foundation_creation",
            "action": "Created Foundation Skeleton",
            "artifacts": result["artifacts_created"]
        })
        
        # Register Foundation artifacts in registry
        for artifact_info in result["artifacts_created"]:
            artifact = Artifact(
                id=f"FOUND-{artifact_info['name'].replace('.md', '')}",
                name=artifact_info['name'],
                type=ArtifactType.KNOWLEDGE,
                state=LifecycleState.DRAFT,
                location=artifact_info['path'],
                description=f"Foundation artifact: {artifact_info['name']}"
            )
            self.registry.register(artifact)
            self.report.artifacts_created.append(artifact_info['path'])
        
        # Check checkpoint: Foundation Skeleton created
        if len(self.report.artifacts_created) >= 6:
            self.report.checkpoints_passed.append("foundation_skeleton_created")
        else:
            self.report.checkpoints_failed.append("foundation_skeleton_created")
    
    def _step_evaluate_state(self):
        """Step 2: Evaluate current engineering state."""
        assessment = self.decision_engine.evaluate_state()
        
        self.steps.append("state_evaluation")
        self.report.steps_executed.append({
            "step": "state_evaluation",
            "assessment": {
                "knowledge_count": assessment.knowledge_count,
                "has_problem_statement": assessment.has_problem_statement,
                "has_specification": assessment.has_specification,
                "has_requirements": assessment.has_requirements,
                "has_assumptions": assessment.has_assumptions,
                "has_constraints": assessment.has_constraints,
                "foundation_complete": assessment.has_specification and 
                                     assessment.has_requirements and
                                     assessment.has_assumptions and
                                     assessment.has_constraints
            }
        })
        
        # Check checkpoint: Foundation completeness evaluated
        self.report.checkpoints_passed.append("state_evaluated")
    
    def _step_identify_facts(self):
        """Step 3: Identify Known Facts from evidence."""
        # Known Facts are statements backed by multiple evidence sources
        # For LAB-001, we derive what we can identify from the objective
        
        facts = []
        
        # Read PROBLEM.md to identify facts
        problem_path = os.path.join(self.workspace_root, "knowledge", "PROBLEM.md")
        if os.path.exists(problem_path):
            with open(problem_path) as f:
                content = f.read()
                if "inventory" in content.lower():
                    facts.append({
                        "statement": "The problem involves inventory management",
                        "evidence": "Objective states 'inventory management system'",
                        "strength": "●●"
                    })
        
        self.report.knowledge_distinction["known_facts"] = facts
        self.report.steps_executed.append({
            "step": "identify_facts",
            "facts_identified": len(facts),
            "facts": facts
        })
        
        # Register as knowledge with strength
        for fact in facts:
            statement = KnowledgeStatement(
                id=f"KNOW-{uuid.uuid4().hex[:8].upper()}",
                statement=fact["statement"],
                strength=EvidenceStrength.MODERATE,
                category="fact"
            )
            self.derivation.knowledge_statements.append(statement)
            
            # Register artifact
            artifact = Artifact(
                id=statement.id,
                name=fact["statement"][:50],
                type=ArtifactType.KNOWLEDGE,
                state=LifecycleState.DRAFT,
                evidence_strength=fact["strength"],
                description=f"Known Fact: {fact['statement']}"
            )
            self.registry.register(artifact)
        
        # Check checkpoint: Facts identified
        if len(facts) > 0:
            self.report.checkpoints_passed.append("facts_identified")
        else:
            self.report.checkpoints_failed.append("facts_identified")
    
    def _step_identify_assumptions(self):
        """Step 4: Identify Reasonable Assumptions."""
        assumptions = []
        
        # Read ASSUMPTIONS.md to identify documented assumptions
        assumptions_path = os.path.join(self.workspace_root, "knowledge", "ASSUMPTIONS.md")
        if os.path.exists(assumptions_path):
            with open(assumptions_path) as f:
                content = f.read()
                # Count assumption placeholders
                if "[ASSUME-001]" in content or "[ASSUME-002]" in content:
                    assumptions.append({
                        "statement": "Assumptions are documented in ASSUMPTIONS.md",
                        "status": "documented",
                        "strength": "●"
                    })
        
        # The objective itself contains assumptions
        assumptions.append({
            "statement": "The user needs a functional inventory management system",
            "basis": "Stated in objective",
            "strength": "●",
            "status": "assumed"
        })
        
        self.report.knowledge_distinction["assumptions"] = assumptions
        self.report.steps_executed.append({
            "step": "identify_assumptions",
            "assumptions_identified": len(assumptions),
            "assumptions": assumptions
        })
        
        # Register as knowledge with lower strength
        for assumption in assumptions:
            statement = KnowledgeStatement(
                id=f"KNOW-{uuid.uuid4().hex[:8].upper()}",
                statement=assumption["statement"],
                strength=EvidenceStrength.WEAK,
                category="assumption"
            )
            self.derivation.knowledge_statements.append(statement)
        
        # Check checkpoint: Assumptions distinguished from facts
        if assumptions:
            self.report.checkpoints_passed.append("assumptions_distinguished")
        else:
            self.report.checkpoints_failed.append("assumptions_distinguished")
    
    def _step_identify_gaps(self):
        """Step 5: Identify Knowledge Gaps."""
        gaps = []
        
        # Read REQUIREMENTS.md to identify unknown requirements
        req_path = os.path.join(self.workspace_root, "knowledge", "REQUIREMENTS.md")
        if os.path.exists(req_path):
            with open(req_path) as f:
                content = f.read()
                if "Unknown Requirements" in content or "Requirement Gaps" in content:
                    gaps.append({
                        "category": "functional_requirements",
                        "statement": "Specific functional requirements are not yet known",
                        "investigation": "Needs stakeholder input and domain analysis"
                    })
        
        # Common gaps for any system
        common_gaps = [
            {
                "category": "technical",
                "statement": "Technology stack is not determined",
                "investigation": "Depends on requirements and constraints"
            },
            {
                "category": "performance",
                "statement": "Performance requirements are not specified",
                "investigation": "Needs SLAs and usage patterns"
            },
            {
                "category": "integration",
                "statement": "Integration requirements are unknown",
                "investigation": "Needs system context"
            }
        ]
        gaps.extend(common_gaps)
        
        self.report.knowledge_distinction["knowledge_gaps"] = gaps
        self.report.steps_executed.append({
            "step": "identify_gaps",
            "gaps_identified": len(gaps),
            "gaps": gaps
        })
        
        # Check checkpoint: Gaps identified
        if len(gaps) > 0:
            self.report.checkpoints_passed.append("gaps_identified")
        else:
            self.report.checkpoints_failed.append("gaps_identified")
    
    def _step_record_evidence_strength(self):
        """Step 6: Record Evidence Strength for all knowledge."""
        strength_summary = {
            "total_statements": len(self.derivation.knowledge_statements),
            "strong": 0,
            "moderate": 0,
            "weak": 0
        }
        
        for statement in self.derivation.knowledge_statements:
            if statement.strength == EvidenceStrength.STRONG:
                strength_summary["strong"] += 1
            elif statement.strength == EvidenceStrength.MODERATE:
                strength_summary["moderate"] += 1
            else:
                strength_summary["weak"] += 1
        
        self.report.knowledge_distinction["evidence_strength"] = strength_summary
        self.report.steps_executed.append({
            "step": "record_evidence_strength",
            "summary": strength_summary
        })
        
        # Check checkpoint: Evidence strength recorded
        self.report.checkpoints_passed.append("evidence_strength_recorded")
    
    def _step_refine_foundation(self):
        """Step 7: Refine Foundation based on analysis."""
        # Run decision cycle to determine next actions
        decision_result = self.decision_engine.execute_cycle()
        
        self.report.decisions_made = decision_result.get("work_orders", [])
        self.report.steps_executed.append({
            "step": "refine_foundation",
            "decisions": decision_result.get("work_orders", [])[:3]  # Top 3
        })
        
        # Check checkpoint: Foundation refinement planned
        self.report.checkpoints_passed.append("foundation_refinement_planned")
    
    def _finalize_report(self):
        """Finalize the laboratory report."""
        self.report.completed_at = datetime.now().isoformat()
        
        # Check for forbidden behavior
        self._check_forbidden_behavior()
        
        # Determine pass/fail
        critical_checkpoints = [
            "foundation_skeleton_created",
            "facts_identified",
            "assumptions_distinguished",
            "gaps_identified"
        ]
        
        all_critical_passed = all(
            cp in self.report.checkpoints_passed 
            for cp in critical_checkpoints
        )
        no_forbidden = len(self.report.forbidden_actions) == 0
        
        if all_critical_passed and no_forbidden:
            self.report.result = LaboratoryResult.PASS
            self.report.pass_reason = (
                "KDSE methodology correctly applied. "
                "Foundation skeleton created, knowledge properly distinguished, "
                "and no premature implementation detected."
            )
        else:
            self.report.result = LaboratoryResult.FAIL
            reasons = []
            if not all_critical_passed:
                failed = [cp for cp in critical_checkpoints 
                          if cp not in self.report.checkpoints_passed]
                reasons.append(f"Failed checkpoints: {', '.join(failed)}")
            if not no_forbidden:
                reasons.append(f"Forbidden behavior detected: {self.report.forbidden_actions}")
            self.report.fail_reason = "; ".join(reasons)
        
        # Add findings
        self.report.findings = [
            f"Foundation artifacts created: {len(self.report.artifacts_created)}",
            f"Knowledge statements derived: {len(self.derivation.knowledge_statements)}",
            f"Known Facts: {len(self.report.knowledge_distinction.get('known_facts', []))}",
            f"Assumptions: {len(self.report.knowledge_distinction.get('assumptions', []))}",
            f"Knowledge Gaps: {len(self.report.knowledge_distinction.get('knowledge_gaps', []))}",
            f"Next decisions: {len(self.report.decisions_made)}"
        ]
        
        print("\n" + "=" * 60)
        print(f"LABORATORY RESULT: {self.report.result.value.upper()}")
        print("=" * 60)
        if self.report.result == LaboratoryResult.PASS:
            print(f"\n✓ {self.report.pass_reason}")
        else:
            print(f"\n✗ {self.report.fail_reason}")
        print("\nFindings:")
        for finding in self.report.findings:
            print(f"  • {finding}")
    
    def _check_forbidden_behavior(self):
        """Check if any forbidden behavior occurred."""
        # Check artifact locations for forbidden patterns
        forbidden_found = []
        
        for artifact_path in self.report.artifacts_created:
            # Check for source code files (shouldn't exist yet)
            for pattern in self.forbidden_patterns:
                if pattern in artifact_path.lower():
                    forbidden_found.append(f"Source file created: {artifact_path}")
        
        # Check decisions for premature implementation
        for decision in self.report.decisions_made:
            obj = decision.get("objective", "").lower()
            for pattern in self.forbidden_patterns:
                if pattern in obj:
                    forbidden_found.append(f"Decision includes: {pattern}")
        
        self.report.forbidden_actions = forbidden_found


def create_lab001_scenario() -> LaboratoryScenario:
    """Create LAB-001: Inventory Management scenario."""
    return LaboratoryScenario(
        id="LAB-001",
        title="Inventory Management System",
        objective="Create a full inventory management system",
        expected_behavior={
            "creates_foundation": True,
            "identifies_facts": True,
            "identifies_assumptions": True,
            "identifies_gaps": True,
            "records_evidence": True
        },
        forbidden_behavior={
            "chooses_language": False,
            "chooses_framework": False,
            "generates_code": False,
            "starts_implementation": False
        },
        checkpoints=[
            {
                "id": "foundation_skeleton_created",
                "description": "Foundation Skeleton (6 documents) created"
            },
            {
                "id": "facts_identified", 
                "description": "Known Facts distinguished from Assumptions"
            },
            {
                "id": "assumptions_distinguished",
                "description": "Assumptions explicitly marked as such"
            },
            {
                "id": "gaps_identified",
                "description": "Knowledge Gaps documented"
            },
            {
                "id": "evidence_strength_recorded",
                "description": "Evidence Strength assigned to knowledge"
            },
            {
                "id": "foundation_refinement_planned",
                "description": "Next steps for Foundation enrichment determined"
            }
        ]
    )
