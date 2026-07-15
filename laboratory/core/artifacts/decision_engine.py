"""
KDSE Slim - Decision Engine
Evaluates current state, identifies gaps, and produces decisions.
Three-step cycle: STATE → GAPS → DECIDE
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional
from .registry import ArtifactRegistry, ArtifactType, LifecycleState


class GapType(Enum):
    """Types of gaps that can be identified."""
    MISSING_KNOWLEDGE = "missing_knowledge"
    MISSING_ARCHITECTURE = "missing_architecture"
    MISSING_IMPLEMENTATION = "missing_implementation"
    MISSING_EVIDENCE = "missing_evidence"
    MISSING_TRACEABILITY = "missing_traceability"
    EVIDENCE_WEAK = "evidence_weak"
    AUTHORITY_VIOLATION = "authority_violation"
    UNAPPROVED_DEPENDENCY = "unapproved_dependency"


class DecisionType(Enum):
    """Types of decisions that can be produced."""
    DERIVE_KNOWLEDGE = "derive_knowledge"
    DERIVE_ARCHITECTURE = "derive_architecture"
    REALIZE_IMPLEMENTATION = "realize_implementation"
    COLLECT_EVIDENCE = "collect_evidence"
    VALIDATE_TRACES = "validate_traces"
    REVIEW_ARTIFACT = "review_artifact"
    APPROVE_ARTIFACT = "approve_artifact"
    CHECK_ALIGNMENT = "check_alignment"
    ENRICH_FOUNDATION = "enrich_foundation"
    IDENTIFY_ASSUMPTIONS = "identify_assumptions"
    IDENTIFY_GAPS = "identify_gaps"
    COMPLETE = "complete"


@dataclass
class Gap:
    """A gap identified in the engineering state."""
    type: GapType
    artifact_id: str = ""
    description: str = ""
    severity: str = "medium"  # low, medium, high, critical
    related_artifacts: List[str] = field(default_factory=list)
    recommendation: str = ""


@dataclass
class WorkOrder:
    """
    A guidance for the next engineering action.
    Work Orders provide guidance, never blocking.
    """
    decision_type: DecisionType
    objective: str
    traces_to: List[str] = field(default_factory=list)
    expected_deliverables: List[str] = field(default_factory=list)
    check: str = ""
    priority: int = 5  # 1-10, higher = more urgent
    rationale: str = ""
    
    def to_dict(self) -> dict:
        return {
            "decision_type": self.decision_type.value,
            "objective": self.objective,
            "traces_to": self.traces_to,
            "expected_deliverables": self.expected_deliverables,
            "check": self.check,
            "priority": self.priority,
            "rationale": self.rationale
        }


@dataclass
class StateAssessment:
    """Assessment of the current engineering state."""
    knowledge_count: int = 0
    architecture_count: int = 0
    implementation_count: int = 0
    verification_count: int = 0
    evidence_count: int = 0
    approved_knowledge: int = 0
    approved_architecture: int = 0
    has_objective: bool = False
    has_problem_statement: bool = False
    has_specification: bool = False
    has_requirements: bool = False
    has_assumptions: bool = False
    has_constraints: bool = False
    has_architecture: bool = False
    traceability_completeness: float = 0.0
    evidence_coverage: float = 0.0


class DecisionEngine:
    """
    The KDSE Slim Decision Engine.
    Implements the three-step cycle: STATE → GAPS → DECIDE
    """
    
    def __init__(self, registry: ArtifactRegistry):
        self.registry = registry
    
    def evaluate_state(self) -> StateAssessment:
        """
        STEP 1: STATE
        Evaluate what exists, its types, states, and traces.
        """
        assessment = StateAssessment()
        
        # Count artifacts by type
        knowledge = self.registry.list_by_type(ArtifactType.KNOWLEDGE)
        architecture = self.registry.list_by_type(ArtifactType.ARCHITECTURE)
        implementation = self.registry.list_by_type(ArtifactType.IMPLEMENTATION)
        verification = self.registry.list_by_type(ArtifactType.VERIFICATION)
        evidence = self.registry.list_by_type(ArtifactType.EVIDENCE)
        
        assessment.knowledge_count = len(knowledge)
        assessment.architecture_count = len(architecture)
        assessment.implementation_count = len(implementation)
        assessment.verification_count = len(verification)
        assessment.evidence_count = len(evidence)
        
        # Count approved artifacts
        assessment.approved_knowledge = len([
            a for a in knowledge if a.state == LifecycleState.APPROVED
        ])
        assessment.approved_architecture = len([
            a for a in architecture if a.state == LifecycleState.APPROVED
        ])
        
        # Check for key artifacts
        for artifact in knowledge:
            if "problem" in artifact.name.lower():
                assessment.has_problem_statement = True
            if "spec" in artifact.name.lower():
                assessment.has_specification = True
            if "requirement" in artifact.name.lower():
                assessment.has_requirements = True
            if "assumption" in artifact.name.lower():
                assessment.has_assumptions = True
            if "constraint" in artifact.name.lower():
                assessment.has_constraints = True
        
        for artifact in architecture:
            if "architecture" in artifact.name.lower():
                assessment.has_architecture = True
        
        # Calculate traceability completeness
        total_traces = sum(len(a.traces_to) for a in self.registry.artifacts.values())
        possible_traces = sum(1 for a in self.registry.artifacts.values() 
                            if a.type != ArtifactType.EVIDENCE)
        if possible_traces > 0:
            assessment.traceability_completeness = min(1.0, total_traces / possible_traces)
        
        # Calculate evidence coverage
        knowledge_with_evidence = len([
            a for a in knowledge if a.evidence_strength
        ])
        if assessment.knowledge_count > 0:
            assessment.evidence_coverage = knowledge_with_evidence / assessment.knowledge_count
        
        assessment.has_objective = assessment.has_problem_statement
        
        return assessment
    
    def identify_gaps(self, assessment: StateAssessment) -> List[Gap]:
        """
        STEP 2: GAPS
        Identify what is missing, weak, or contradictory.
        """
        gaps = []
        
        # Missing Foundation artifacts
        if not assessment.has_problem_statement:
            gaps.append(Gap(
                type=GapType.MISSING_KNOWLEDGE,
                description="Problem statement is missing",
                severity="critical",
                recommendation="Create PROBLEM.md defining the problem scope"
            ))
        
        if not assessment.has_specification:
            gaps.append(Gap(
                type=GapType.MISSING_KNOWLEDGE,
                description="Project specification is missing",
                severity="critical",
                recommendation="Create SPEC.md with project overview"
            ))
        
        if not assessment.has_requirements:
            gaps.append(Gap(
                type=GapType.MISSING_KNOWLEDGE,
                description="Requirements are missing",
                severity="high",
                recommendation="Create REQUIREMENTS.md with functional requirements"
            ))
        
        if not assessment.has_assumptions:
            gaps.append(Gap(
                type=GapType.MISSING_KNOWLEDGE,
                description="Assumptions are not documented",
                severity="medium",
                recommendation="Create ASSUMPTIONS.md with documented assumptions"
            ))
        
        if not assessment.has_constraints:
            gaps.append(Gap(
                type=GapType.MISSING_KNOWLEDGE,
                description="Constraints are not documented",
                severity="medium",
                recommendation="Create CONSTRAINTS.md with project constraints"
            ))
        
        # Missing Architecture
        if not assessment.has_architecture and assessment.has_requirements:
            gaps.append(Gap(
                type=GapType.MISSING_ARCHITECTURE,
                description="Architecture decisions are missing",
                severity="high",
                recommendation="Create ARCHITECTURE.md with structural decisions"
            ))
        
        # Missing evidence
        if assessment.knowledge_count > 0 and assessment.evidence_coverage < 0.5:
            gaps.append(Gap(
                type=GapType.MISSING_EVIDENCE,
                description=f"Only {assessment.evidence_coverage:.0%} of knowledge has evidence",
                severity="high",
                recommendation="Collect evidence sources to support knowledge claims"
            ))
        
        # Weak evidence
        for artifact in self.registry.list_by_type(ArtifactType.KNOWLEDGE):
            if artifact.evidence_strength == "●":
                gaps.append(Gap(
                    type=GapType.EVIDENCE_WEAK,
                    artifact_id=artifact.id,
                    description=f"Weak evidence for knowledge: {artifact.name}",
                    severity="medium",
                    recommendation="Find additional evidence sources to strengthen knowledge"
                ))
        
        # Authority violations
        for artifact in self.registry.artifacts.values():
            violations = self.registry.check_authority_violation(artifact.id)
            for violation in violations:
                gaps.append(Gap(
                    type=GapType.AUTHORITY_VIOLATION,
                    artifact_id=artifact.id,
                    description=violation,
                    severity="critical",
                    recommendation="Address authority violation before proceeding"
                ))
        
        # Unapproved dependencies
        for artifact in self.registry.list_by_type(ArtifactType.ARCHITECTURE):
            if artifact.state != LifecycleState.APPROVED:
                gaps.append(Gap(
                    type=GapType.UNAPPROVED_DEPENDENCY,
                    artifact_id=artifact.id,
                    description=f"Architecture not approved: {artifact.name}",
                    severity="high",
                    recommendation="Approve architecture before implementation"
                ))
        
        # Missing implementation traceability
        for artifact in self.registry.list_by_type(ArtifactType.IMPLEMENTATION):
            if not artifact.traces_to:
                gaps.append(Gap(
                    type=GapType.MISSING_TRACEABILITY,
                    artifact_id=artifact.id,
                    description=f"Implementation has no traceability: {artifact.name}",
                    severity="high",
                    recommendation="Add trace to architecture decision"
                ))
        
        return gaps
    
    def decide(self, assessment: StateAssessment, gaps: List[Gap]) -> List[WorkOrder]:
        """
        STEP 3: DECIDE
        Produce guidance based on state and gaps.
        Decision is guided by principles, not rules.
        """
        work_orders = []
        
        # Critical gaps take priority
        critical_gaps = [g for g in gaps if g.severity == "critical"]
        high_gaps = [g for g in gaps if g.severity == "high"]
        
        # Foundation creation is always first
        if not assessment.has_problem_statement:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.DERIVE_KNOWLEDGE,
                objective="Create Foundation Skeleton",
                expected_deliverables=[
                    ".kdse/knowledge/PROBLEM.md",
                    ".kdse/knowledge/SPEC.md"
                ],
                check="Does PROBLEM.md define clear scope and constraints?",
                priority=10,
                rationale="Foundation must exist before any other work"
            ))
            return work_orders  # Foundation first, then continue
        
        if not assessment.has_specification:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.DERIVE_KNOWLEDGE,
                objective="Create SPEC.md",
                expected_deliverables=[
                    ".kdse/knowledge/SPEC.md"
                ],
                check="Does SPEC.md provide project overview?",
                priority=9,
                rationale="Specification needed to define what to build"
            ))
            return work_orders
        
        if not assessment.has_requirements:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.DERIVE_KNOWLEDGE,
                objective="Document requirements",
                expected_deliverables=[
                    ".kdse/knowledge/REQUIREMENTS.md"
                ],
                check="Are all functional requirements documented?",
                priority=8,
                rationale="Requirements needed before architecture"
            ))
            return work_orders
        
        if not assessment.has_assumptions:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.IDENTIFY_ASSUMPTIONS,
                objective="Document assumptions",
                expected_deliverables=[
                    ".kdse/knowledge/ASSUMPTIONS.md"
                ],
                check="Are all reasonable assumptions documented?",
                priority=7,
                rationale="Assumptions must be explicit, not silent"
            ))
            return work_orders
        
        if not assessment.has_constraints:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.DERIVE_KNOWLEDGE,
                objective="Document constraints",
                expected_deliverables=[
                    ".kdse/knowledge/CONSTRAINTS.md"
                ],
                check="Are all constraints (technical, business, regulatory) documented?",
                priority=7,
                rationale="Constraints limit architecture choices"
            ))
            return work_orders
        
        # Evidence collection if weak coverage
        if assessment.evidence_coverage < 0.8:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.COLLECT_EVIDENCE,
                objective="Enrich knowledge with evidence",
                expected_deliverables=[
                    ".kdse/evidence/sources/ - Evidence source inventory"
                ],
                check="Do all knowledge statements cite evidence?",
                priority=7,
                rationale="Knowledge without evidence is opinion"
            ))
            return work_orders
        
        # Architecture if prerequisites met
        if not assessment.has_architecture and assessment.approved_knowledge >= 2:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.DERIVE_ARCHITECTURE,
                objective="Define system architecture",
                expected_deliverables=[
                    ".kdse/architecture/ARCHITECTURE.md"
                ],
                traces_to=[a.id for a in self.registry.list_approved()[:3]],
                check="Does architecture trace to knowledge?",
                priority=8,
                rationale="Architecture must derive from approved knowledge"
            ))
            return work_orders
        
        # Check alignment if architecture exists
        if assessment.has_architecture and assessment.architecture_count > 0:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.CHECK_ALIGNMENT,
                objective="Verify K→A→I alignment",
                expected_deliverables=[
                    ".kdse/verification/checks/ - Alignment verification"
                ],
                check="Does every architecture decision cite knowledge?",
                priority=6,
                rationale="Continuous verification ensures alignment"
            ))
        
        # Address critical gaps
        for gap in critical_gaps:
            if gap.type == GapType.AUTHORITY_VIOLATION:
                work_orders.append(WorkOrder(
                    decision_type=DecisionType.VALIDATE_TRACES,
                    objective=f"Fix authority violation: {gap.description}",
                    priority=10,
                    rationale="Authority violations must be resolved"
                ))
        
        # Address high gaps
        for gap in high_gaps:
            if gap.type == GapType.MISSING_ARCHITECTURE:
                work_orders.append(WorkOrder(
                    decision_type=DecisionType.DERIVE_ARCHITECTURE,
                    objective="Create architecture decisions",
                    expected_deliverables=[".kdse/architecture/"],
                    priority=8
                ))
        
        # If no gaps, identify enrichment opportunities
        if not work_orders:
            work_orders.append(WorkOrder(
                decision_type=DecisionType.ENRICH_FOUNDATION,
                objective="Enrich existing artifacts",
                expected_deliverables=[
                    "Updated knowledge with stronger evidence",
                    "Refined architecture with more detail"
                ],
                check="Can any knowledge be strengthened?",
                priority=3,
                rationale="Continuous enrichment improves quality"
            ))
        
        return work_orders
    
    def execute_cycle(self) -> dict:
        """
        Execute the full three-step decision cycle.
        """
        # Step 1: Evaluate state
        assessment = self.evaluate_state()
        
        # Step 2: Identify gaps
        gaps = self.identify_gaps(assessment)
        
        # Step 3: Decide
        work_orders = self.decide(assessment, gaps)
        
        return {
            "assessment": {
                "knowledge_count": assessment.knowledge_count,
                "architecture_count": assessment.architecture_count,
                "implementation_count": assessment.implementation_count,
                "approved_knowledge": assessment.approved_knowledge,
                "approved_architecture": assessment.approved_architecture,
                "traceability_completeness": assessment.traceability_completeness,
                "evidence_coverage": assessment.evidence_coverage,
                "foundation_complete": assessment.has_specification and 
                                     assessment.has_requirements and
                                     assessment.has_assumptions and
                                     assessment.has_constraints
            },
            "gaps": [
                {
                    "type": g.type.value,
                    "artifact_id": g.artifact_id,
                    "description": g.description,
                    "severity": g.severity,
                    "recommendation": g.recommendation
                }
                for g in gaps
            ],
            "work_orders": [wo.to_dict() for wo in work_orders],
            "next_action": work_orders[0].to_dict() if work_orders else None
        }
